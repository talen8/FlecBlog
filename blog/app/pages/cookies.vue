<script lang="ts" setup>
definePageMeta({
  showSidebar: false,
});

useSeoMeta({
  title: 'Cookie',
  description: '了解本站在浏览器中存储的数据类型及如何管理这些数据',
});
</script>

<template>
  <div id="page">
    <div class="page-title">Cookie 与本地存储</div>

    <div class="content">
      <p>本页面说明本站在你的浏览器中存储了哪些数据，以及如何管理这些数据。</p>

      <p>
        本站不使用
        Cookie，而是使用浏览器本地存储（LocalStorage）。这些数据仅保存在你的设备上，不会自动发送到服务器。
      </p>

      <p>
        关于数据收集和使用的完整说明，请查看 <router-link to="/privacy">隐私政策</router-link>。
      </p>

      <h1>存储的数据</h1>

      <ul>
        <li><strong>主题偏好</strong>（isDark）：记住你选择的深色/浅色模式</li>
        <li><strong>游客信息</strong>（guest_info）：未登录时评论所填写的昵称、邮箱、网址</li>
        <li><strong>登录令牌</strong>（token）：登录后的身份凭证，用于保持登录状态</li>
      </ul>

      <h1>为什么使用本地存储</h1>

      <p><strong>更好的隐私保护：</strong></p>
      <p>
        与传统 Cookie 不同，LocalStorage
        中的数据不会在每次请求时自动发送到服务器。这意味着你的主题偏好等个人设置完全保存在本地，服务器无权访问。
      </p>

      <p><strong>更快的加载速度：</strong></p>
      <p>主题设置等数据直接从本地读取，无需等待服务器响应，页面加载更快。</p>

      <p><strong>离线也能工作：</strong></p>
      <p>即使网络断开，你的主题设置仍然有效，不会因为无法连接服务器而丢失。</p>

      <h1>你的数据控制权</h1>

      <ul>
        <li><strong>完全控制</strong>：这些数据存储在你的设备上，你可以随时查看、修改或删除</li>
        <li><strong>不会跟踪</strong>：我们不使用这些数据进行跨站跟踪或用户画像</li>
        <li><strong>不会出售</strong>：你的本地数据归你所有，我们无权访问也不会出售</li>
        <li><strong>随时清除</strong>：可以通过浏览器设置随时清除，无需联系我们</li>
      </ul>

      <h1>如何管理这些数据</h1>

      <p><strong>查看存储的数据：</strong></p>
      <p>
        按 <strong>F12</strong> 打开开发者工具 → Application/Storage → Local Storage →
        选择本站域名。
      </p>

      <p><strong>清除存储的数据：</strong></p>
      <ul>
        <li>在浏览器地址栏点击锁图标 → 网站设置 → 清除数据</li>
        <li>
          或按 <strong>Ctrl+Shift+Delete</strong>（Mac:
          <strong>⌘+Shift+Delete</strong>）清除浏览数据
        </li>
      </ul>

      <p><strong>清除后的影响：</strong></p>
      <ul>
        <li>主题设置恢复为默认</li>
        <li>游客评论需重新填写信息</li>
        <li>登录状态丢失，需重新登录</li>
      </ul>

      <blockquote>清除这些数据不会影响你的账号信息（如已注册），只是需要重新登录。</blockquote>

      <h1>常见场景说明</h1>

      <p><strong>更换设备：</strong></p>
      <p>
        本地存储只保存在当前设备的当前浏览器中。如果你更换了设备或浏览器，需要重新设置主题偏好和登录账号。但你的账号数据（如发表的评论）保存在服务器上，不会丢失。
      </p>

      <p><strong>隐私模式/无痕浏览：</strong></p>
      <p>
        在隐私模式下，本地存储功能正常工作。但关闭隐私窗口后，所有数据会被自动清除。如果你经常使用隐私模式，建议注册账号以保存你的部分偏好设置。
      </p>

      <p><strong>清除浏览数据：</strong></p>
      <p>
        如果你清除了浏览器数据，主题会恢复默认，登录状态会丢失。但你的账号本身不受影响，重新登录即可。
      </p>

      <p><strong>同步到其他设备：</strong></p>
      <p>
        浏览器本地存储的数据不会通过浏览器同步功能（如 Chrome
        同步）传输到其他设备。如果你在多个设备上使用本站，建议在每个设备上分别登录账号。
      </p>

      <h1>免责声明</h1>

      <p>本站使用的浏览器存储技术遵循 Web 标准。我们不对以下情况承担责任：</p>

      <ul>
        <li>因浏览器设置、扩展程序或其他软件导致的存储失效</li>
        <li>因用户主动清除数据或浏览器隐私模式导致的数据丢失</li>
        <li>第三方服务（如 Cravatar）的可用性和隐私政策变更</li>
      </ul>

      <p>继续使用本站即表示你理解并接受上述说明。</p>

      <p class="update-time">最后更新时间：2025年10月</p>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#page {
  @extend .cardHover;
  align-self: flex-start;
  padding: 40px;

  .page-title {
    margin: 0 0 10px;
    font-weight: bold;
    font-size: 2rem;
  }

  .content {
    color: var(--font-color);
    line-height: 1.8;
    font-size: 1rem;

    h1 {
      font-size: 1.5rem;
      font-weight: 600;
      margin: 2.5rem 0 1rem;
      color: var(--font-color);
      line-height: 1.3;

      &:first-child {
        margin-top: 1.5rem;
      }
    }

    h2 {
      font-size: 1.25rem;
      font-weight: 600;
      margin: 1.75rem 0 0.875rem;
      color: var(--font-color);
      line-height: 1.3;
    }

    p {
      margin: 1rem 0;
      text-align: justify;
    }

    a {
      color: var(--theme-color);
      text-decoration: none;
      font-weight: 500;
      border-bottom: 1px solid transparent;
      transition: all 0.3s ease;

      &:hover {
        border-bottom-color: var(--theme-color);
      }
    }

    strong {
      font-weight: 600;
      color: var(--font-color);
    }

    ul,
    ol {
      padding-left: 2rem;
      margin: 1rem 0;

      li {
        margin: 0.5rem 0;
      }
    }

    ul {
      list-style: disc;
    }

    ol {
      list-style: decimal;
    }

    blockquote {
      margin: 1.5rem 0;
      padding: 0.75em 1.25em;
      background: var(--flec-heavy-bg);
      border-left: 4px solid var(--theme-color);
      border-radius: 0 6px 6px 0;
      font-size: 0.95em;
      color: var(--theme-meta-color);
      line-height: 1.7;
    }

    .update-time {
      margin-top: 3rem;
      color: var(--theme-meta-color);
      font-size: 0.875rem;
      text-align: right;
      font-weight: 500;
    }
  }
}

// 响应式设计
@media screen and (max-width: 1024px) {
  #page {
    padding: 30px;

    .page-title {
      font-size: 1.75rem;
    }

    .content {
      h1 {
        font-size: 1.35rem;
        margin: 2rem 0 0.85rem;
      }

      h2 {
        font-size: 1.15rem;
        margin: 1.5rem 0 0.75rem;
      }

      p {
        font-size: 0.95rem;
      }

      ul,
      ol {
        li {
          font-size: 0.95rem;
        }
      }

      blockquote {
        font-size: 0.9em;
        padding: 0.65em 1.1em;
      }
    }
  }
}

@media screen and (max-width: 768px) {
  #page {
    padding: 18px;

    .page-title {
      font-size: 1.4rem;
    }

    .content {
      font-size: 0.9rem;

      h1 {
        font-size: 1.25rem;
        margin: 1.75rem 0 0.75rem;
      }

      h2 {
        font-size: 1.1rem;
        margin: 1.35rem 0 0.65rem;
      }

      p {
        font-size: 0.9rem;
      }

      ul,
      ol {
        padding-left: 1.5rem;

        li {
          margin: 0.4rem 0;
          font-size: 0.9rem;
        }
      }

      blockquote {
        font-size: 0.88em;
        padding: 0.65em 1em;
        margin: 1.25rem 0;
      }

      .update-time {
        font-size: 0.8rem;
      }
    }
  }
}
</style>
