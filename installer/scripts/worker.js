/**
 * Cloudflare Worker for FlecBlog Installer
 * 提供安装脚本的动态分发服务
 *
 * 访问地址:
 * - https://install.flec.top/         -> 返回 Shell 脚本 (Linux)
 */

const REPO = "talen8/FlecBlog";
const BINARY = "flecb";

/**
 * 主入口函数 - 处理所有 HTTP 请求
 */
export default {
  async fetch(request, env) {
    const url = new URL(request.url);
    const pathname = url.pathname;

    const headers = {
      "Content-Type": "text/plain; charset=utf-8",
      "Cache-Control": "no-cache, no-store, must-revalidate",
    };

    return new Response(getShellScript(), { headers });
  },
};

/**
 * 生成 Shell 脚本 (Linux)
 * 功能：下载 flecb 二进制文件并启动安装向导
 */
function getShellScript() {
  return `#!/bin/sh
set -eu

# 配置变量
REPO="${REPO}"
BINARY="${BINARY}"
INSTALL_DIR="/usr/local/bin"

# 检测 root 权限
if [ "$(id -u)" -ne 0 ]; then
  echo "错误: 此脚本需要 root 权限，请使用 sudo 运行"
  echo "  sudo bash $(basename "$0")"
  exit 1
fi

# 检测操作系统和架构
# 支持的系统: linux
# 支持的架构: amd64, arm64
detect_platform() {
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"
  case "$os" in
    linux) os="linux" ;;
    *) echo "不支持的操作系统: $os，仅支持 Linux"; exit 1 ;;
  esac
  case "$arch" in
    x86_64|amd64) arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *) echo "不支持的架构: $arch"; exit 1 ;;
  esac
  echo "\${os}_\${arch}"
}

# 从 GitHub API 获取最新安装器版本号
# 返回格式: x.x.x (不包含 v 前缀)
# 注意: /releases/latest 返回的是所有 release 中最新的，可能是 FlecBlog 版本
# 所以需要从所有 releases 中筛选 installer/ 开头的 tag
get_version() {
  curl -fsSL "https://api.github.com/repos/\${REPO}/releases" | \\
    grep '"tag_name":' | \\
    grep '"installer/v' | \\
    head -1 | \\
    grep -oP '(?<=installer/v)[^"]+'
}

# 主函数
main() {
  echo "FlecBlog 安装器"
  echo "==============="
  echo ""

  # 获取平台和版本信息
  platform=$(detect_platform)
  version=$(get_version)

  echo "平台: \${platform}"
  echo "版本: v\${version}"
  echo ""

  # 创建临时目录
  tmp_dir=$(mktemp -d)
  trap 'rm -rf "$tmp_dir"' EXIT

  # GitHub 镜像列表
  mirrors="
https://ghfast.top/https://github.com
https://ghproxy.net/https://github.com
https://github.com
"

  # 尝试多个镜像下载
  echo "下载中..."
  downloaded=false
  for mirror in $mirrors; do
    [ -z "$mirror" ] && continue
    url="$mirror/\${REPO}/releases/download/installer/v\${version}/\${BINARY}_\${platform}"
    echo "尝试: $mirror"
    if curl -fSL --connect-timeout 10 --max-time 60 --progress-bar -o "\${tmp_dir}/\${BINARY}" "$url"; then
      downloaded=true
      break
    fi
    echo "失败，尝试下一个..."
  done

  if [ "$downloaded" != "true" ]; then
    echo "错误: 所有镜像下载失败"
    exit 1
  fi

  chmod +x "\${tmp_dir}/\${BINARY}"

  # 安装到系统目录
  mv "\${tmp_dir}/\${BINARY}" "\${INSTALL_DIR}/\${BINARY}"

  echo "flecb 已安装到: \${INSTALL_DIR}/\${BINARY}"
  echo ""

  # 启动 FlecBlog 安装向导（重定向 stdin 到终端以支持交互）
  "\${INSTALL_DIR}/\${BINARY}" install < /dev/tty
}

main
`;
}
