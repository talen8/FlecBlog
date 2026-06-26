const INSTALL_SCRIPT =
  '#!/bin/bash\n' +
  '# FlecBlog 一键部署与管理脚本\n' +
  '# 部署: curl -fsSL https://get.flec.top | bash\n' +
  '# 主题: curl -fsSL https://get.flec.top | bash -s theme <镜像名称>\n' +
  '# 卸载: curl -fsSL https://get.flec.top | bash -s uninstall\n' +
  '# 升级: curl -fsSL https://get.flec.top | bash -s upgrade\n' +
  '# 状态: curl -fsSL https://get.flec.top | bash -s status\n' +
  '# 日志: curl -fsSL https://get.flec.top | bash -s logs\n' +
  '# 帮助: curl -fsSL https://get.flec.top | bash -s help\n' +
  '\n' +
  'set -e\n' +
  '\n' +
  "RED='\\033[0;31m'\n" +
  "GREEN='\\033[0;32m'\n" +
  "YELLOW='\\033[1;33m'\n" +
  "CYAN='\\033[0;36m'\n" +
  "NC='\\033[0m'\n" +
  '\n' +
  'info() { echo -e "${CYAN}[INFO]${NC} $1"; }\n' +
  'success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }\n' +
  'warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }\n' +
  'error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }\n' +
  '\n' +
  'ACTION="deploy"\n' +
  'THEME=""\n' +
  '\n' +
  'while [[ $# -gt 0 ]]; do\n' +
  '  case $1 in\n' +
  '    theme)\n' +
  '      THEME="$2"\n' +
  '      ACTION="theme"\n' +
  '      shift 2\n' +
  '      ;;\n' +
  '    uninstall)\n' +
  '      ACTION="uninstall"\n' +
  '      shift\n' +
  '      ;;\n' +
  '    upgrade)\n' +
  '      ACTION="upgrade"\n' +
  '      shift\n' +
  '      ;;\n' +
  '    status)\n' +
  '      ACTION="status"\n' +
  '      shift\n' +
  '      ;;\n' +
  '    logs)\n' +
  '      ACTION="logs"\n' +
  '      shift\n' +
  '      ;;\n' +
  '    help)\n' +
  '      ACTION="help"\n' +
  '      shift\n' +
  '      ;;\n' +
  '    *)\n' +
  '      shift\n' +
  '      ;;\n' +
  '  esac\n' +
  'done\n' +
  '\n' +
  'deploy() {\n' +
  '  if [ -f docker-compose.yml ] || [ -f .env ]; then\n' +
  '    error "当前目录已存在 docker-compose.yml 或 .env，请确认安装位置"\n' +
  '  fi\n' +
  '\n' +
  '  info "正在创建配置文件..."\n' +
  '\n' +
  "  cat > docker-compose.yml << 'COMPOSE_EOF'\n" +
  'services:\n' +
  '  postgres:\n' +
  '    image: postgres:15-alpine\n' +
  '    container_name: flec_postgres\n' +
  '    restart: unless-stopped\n' +
  '    environment:\n' +
  '      POSTGRES_DB: postgres\n' +
  '      POSTGRES_USER: postgres\n' +
  '      POSTGRES_PASSWORD: ${DB_PASSWORD}\n' +
  '    volumes:\n' +
  '      - postgres_data:/var/lib/postgresql/data\n' +
  '    networks:\n' +
  '      - flec-network\n' +
  '    healthcheck:\n' +
  '      test: ["CMD-SHELL", "pg_isready -U postgres"]\n' +
  '      interval: 10s\n' +
  '      timeout: 5s\n' +
  '      retries: 5\n' +
  '\n' +
  '  server:\n' +
  '    image: talen8/flec-server:latest\n' +
  '    container_name: flec_server\n' +
  '    restart: unless-stopped\n' +
  '    environment:\n' +
  '      DB_HOST: postgres\n' +
  '      DB_PORT: 5432\n' +
  '      DB_NAME: postgres\n' +
  '      DB_USER: postgres\n' +
  '      DB_PASSWORD: ${DB_PASSWORD}\n' +
  '      JWT_SECRET: ${JWT_SECRET}\n' +
  '    ports:\n' +
  '      - "${SERVER_PORT}:8080"\n' +
  '    volumes:\n' +
  '      - ./data:/app/data\n' +
  '    networks:\n' +
  '      - flec-network\n' +
  '    depends_on:\n' +
  '      postgres:\n' +
  '        condition: service_healthy\n' +
  '\n' +
  '  blog:\n' +
  '    image: talen8/flec-blog:latest\n' +
  '    container_name: flec_blog\n' +
  '    restart: unless-stopped\n' +
  '    environment:\n' +
  '      NUXT_PUBLIC_API_URL: ${API_URL}\n' +
  '    ports:\n' +
  '      - "${BLOG_PORT}:3000"\n' +
  '    networks:\n' +
  '      - flec-network\n' +
  '    depends_on:\n' +
  '      - server\n' +
  '\n' +
  '  admin:\n' +
  '    image: talen8/flec-admin:latest\n' +
  '    container_name: flec_admin\n' +
  '    restart: unless-stopped\n' +
  '    environment:\n' +
  '      API_URL: ${API_URL}\n' +
  '    ports:\n' +
  '      - "${ADMIN_PORT}:4000"\n' +
  '    networks:\n' +
  '      - flec-network\n' +
  '    depends_on:\n' +
  '      - server\n' +
  '\n' +
  'networks:\n' +
  '  flec-network:\n' +
  '    driver: bridge\n' +
  '\n' +
  'volumes:\n' +
  '  postgres_data:\n' +
  'COMPOSE_EOF\n' +
  '\n' +
  "  cat > .env << 'ENV_EOF'\n" +
  '# 数据库密码\nDB_PASSWORD=your_database_password\n' +
  '# JWT 密钥\nJWT_SECRET=your_jwt_secret_key\n' +
  '# API 对外访问地址\nAPI_URL=http://localhost:8080/api/v1\n' +
  '# 后端服务端口\nSERVER_PORT=8080\n' +
  '# 博客端端口\nBLOG_PORT=3000\n' +
  '# 管理端端口\nADMIN_PORT=4000\n' +
  'ENV_EOF\n' +
  '\n' +
  '  success "创建完成！配置文件已生成到当前目录"\n' +
  '  echo ""\n' +
  '  echo "  接下来："\n' +
  '  echo -e "  ${YELLOW}1.${NC} 修改 .env 中的配置"\n' +
  '  echo -e "  ${YELLOW}2.${NC} 执行 docker compose up -d"\n' +
  '  echo -e "  ${YELLOW}3.${NC} 配置外部访问（可参考文档）"\n' +
  '  echo ""\n' +
  '}\n' +
  '\n' +
  'uninstall() {\n' +
  '  if [ ! -f docker-compose.yml ]; then\n' +
  '    error "当前目录未找到 docker-compose.yml，请确认安装位置"\n' +
  '  fi\n' +
  '\n' +
  '  info "正在停止并移除容器..."\n' +
  '  docker compose down\n' +
  '\n' +
  '  success "卸载完成！容器已移除，数据卷已保留"\n' +
  '  echo ""\n' +
  '  echo "  如需删除数据卷，执行: docker compose down -v"\n' +
  '  echo "  如需删除配置文件，执行: rm docker-compose.yml .env"\n' +
  '  echo "  如需删除服务数据，执行: rm -rf ./data"\n' +
  '  echo ""\n' +
  '}\n' +
  '\n' +
  'upgrade() {\n' +
  '  if [ ! -f docker-compose.yml ]; then\n' +
  '    error "当前目录未找到 docker-compose.yml，请确认安装位置"\n' +
  '  fi\n' +
  '\n' +
  '  info "正在拉取最新镜像并重启服务..."\n' +
  '  docker compose pull\n' +
  '  docker compose up -d\n' +
  '\n' +
  '  success "升级完成！"\n' +
  '  echo ""\n' +
  '}\n' +
  '\n' +
  'status() {\n' +
  '  if [ ! -f docker-compose.yml ]; then\n' +
  '    error "当前目录未找到 docker-compose.yml，请确认安装位置"\n' +
  '  fi\n' +
  '\n' +
  '  docker compose ps\n' +
  '}\n' +
  '\n' +
  'logs() {\n' +
  '  if [ ! -f docker-compose.yml ]; then\n' +
  '    error "当前目录未找到 docker-compose.yml，请确认安装位置"\n' +
  '  fi\n' +
  '\n' +
  '  docker compose logs -f\n' +
  '}\n' +
  '\n' +
  'help() {\n' +
  '  echo ""\n' +
  '  echo -e "${CYAN}FlecBlog 管理脚本${NC}"\n' +
  '  echo ""\n' +
  '  echo "用法: curl -fsSL https://get.flec.top | bash -s <命令>"\n' +
  '  echo ""\n' +
  '  echo "命令:"\n' +
  '  echo "  (无参数)    创建配置文件"\n' +
  '  echo "  theme       更换博客主题"\n' +
  '  echo "  upgrade     拉取最新镜像并重启"\n' +
  '  echo "  status      查看服务状态"\n' +
  '  echo "  logs        查看服务日志"\n' +
  '  echo "  uninstall   停止并移除容器"\n' +
  '  echo "  help        显示此帮助信息"\n' +
  '  echo ""\n' +
  '}\n' +
  '\n' +
  'change_theme() {\n' +
  '  if [ -z "$THEME" ]; then\n' +
  '    error "请指定主题镜像名称，例如: bash -s theme talen8/flec-blog"\n' +
  '  fi\n' +
  '\n' +
  '  if [ ! -f docker-compose.yml ]; then\n' +
  '    error "当前目录未找到 docker-compose.yml，请确认安装位置"\n' +
  '  fi\n' +
  '\n' +
  '  [[ "$THEME" == *:* ]] || THEME="$THEME:latest"\n' +
  "  sed -i '/blog:/,/^[^ ]/{s|image: .*|image: '\"$THEME\"'|;}' docker-compose.yml\n" +
  '\n' +
  '  info "正在重启博客服务..."\n' +
  '  docker compose up -d blog\n' +
  '\n' +
  '  success "主题已更换为 $THEME"\n' +
  '  echo ""\n' +
  '}\n' +
  '\n' +
  'case $ACTION in\n' +
  '  deploy)\n' +
  '    deploy\n' +
  '    ;;\n' +
  '  theme)\n' +
  '    change_theme\n' +
  '    ;;\n' +
  '  uninstall)\n' +
  '    uninstall\n' +
  '    ;;\n' +
  '  upgrade)\n' +
  '    upgrade\n' +
  '    ;;\n' +
  '  status)\n' +
  '    status\n' +
  '    ;;\n' +
  '  logs)\n' +
  '    logs\n' +
  '    ;;\n' +
  '  help)\n' +
  '    help\n' +
  '    ;;\n' +
  'esac\n';

export default {
  fetch(request) {
    return new Response(INSTALL_SCRIPT, {
      headers: {
        'Content-Type': 'text/plain; charset=utf-8',
      },
    });
  },
};
