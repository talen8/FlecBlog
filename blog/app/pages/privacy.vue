<script lang="ts" setup>
const { basicConfig } = useSysConfig();

definePageMeta({
  showSidebar: false,
});

useSeoMeta({
  title: '隐私政策',
  description: '了解我们如何收集、使用和保护您的个人信息',
});

const contactEmail = computed(() => basicConfig.value.author_email || '');
</script>

<template>
  <div id="page">
    <div class="page-title">隐私政策</div>

    <div class="content">
      <p>本隐私政策说明我们如何收集、使用和保护你的个人信息。</p>

      <p>
        继续使用本站即表示你同意本政策。如有疑问，请通过
        <strong>{{ contactEmail }}</strong> 联系我们。
      </p>

      <h1>一、我们收集的信息</h1>

      <h2>1.1 自动收集的信息</h2>

      <p>当你访问本站时，我们会自动收集以下技术信息用于统计分析：</p>

      <ul>
        <li>访问的页面、停留时长、来源网站</li>
        <li>IP 地址、地理位置（国家/城市级别）</li>
        <li>设备类型、操作系统、浏览器类型</li>
      </ul>

      <h2>1.2 用户主动提供的信息</h2>

      <p><strong>注册账号时：</strong></p>
      <ul>
        <li>用户名、邮箱地址、密码（加密存储）</li>
      </ul>

      <p><strong>发表评论时：</strong></p>
      <ul>
        <li>游客评论：昵称、邮箱、个人网站（可选）</li>
        <li>用户评论：使用账号关联的信息</li>
        <li>评论内容、IP 地址</li>
      </ul>

      <h1>二、信息的使用目的</h1>

      <p>我们收集的信息仅用于以下目的：</p>

      <ul>
        <li>分析访问统计，优化网站内容和结构</li>
        <li>监控网站性能，提升用户体验</li>
        <li>识别异常访问，防范安全威胁</li>
        <li>提供评论功能，发送回复通知</li>
        <li>管理用户账号，提供个性化服务</li>
      </ul>

      <h1>三、信息的存储</h1>

      <p><strong>服务器存储：</strong>用户账号信息、评论内容、访问日志等存储在服务器数据库中。</p>

      <p>
        <strong>本地存储：</strong>主题偏好、登录令牌等保存在你的浏览器中，详见
        <router-link to="/cookies">Cookie 与本地存储</router-link>。
      </p>

      <h1>四、第三方服务</h1>

      <p>
        本站使用 Cravatar 头像服务，仅发送邮箱的 MD5 哈希值（不发送原始邮箱），用于获取用户头像。
      </p>

      <h1>五、信息的共享</h1>

      <p>我们不会出售或向第三方披露你的个人信息，除非：</p>

      <ul>
        <li>获得你的明确同意</li>
        <li>法律法规要求</li>
        <li>保护网站和用户安全所必需</li>
      </ul>

      <h1>六、数据安全</h1>

      <ul>
        <li>使用 HTTPS 加密传输</li>
        <li>密码经过加密存储</li>
        <li>严格的服务器访问控制</li>
        <li>定期安全审查和备份</li>
      </ul>

      <h1>七、你的权利</h1>

      <p>你有权：</p>

      <ul>
        <li>访问、更正、删除你的个人信息</li>
        <li>删除账号及所有相关数据</li>
        <li>导出你发表的评论</li>
        <li>拒绝特定数据处理</li>
      </ul>

      <p>
        如需行使上述权利，请联系 <strong>{{ contactEmail }}</strong
        >。
      </p>

      <h1>八、免责声明</h1>

      <p>
        <strong>数据安全：</strong
        >我们采取合理措施保护你的信息，但无法保证绝对安全。互联网传输和电子存储方式均存在固有风险。
      </p>

      <p>
        <strong>用户责任：</strong
        >你应对提供给本站的信息的准确性和合法性负责。因提供虚假、违法信息造成的后果由你自行承担。
      </p>

      <p>
        <strong>服务中断：</strong
        >本站服务可能因维护、故障或不可抗力等原因中断或变更，我们不承担由此产生的损失。
      </p>

      <p>
        <strong>第三方服务：</strong
        >本站可能包含第三方链接或服务，我们不对第三方的隐私政策或内容负责。
      </p>

      <h1>九、内容许可</h1>

      <p>你发表的评论和内容，即表示你授予本站以下权利：</p>

      <ul>
        <li>在本站展示、存储和使用你的内容</li>
        <li>对违反法律或本站规则的内容进行删除或修改</li>
        <li>在必要时向相关部门提供你的信息</li>
      </ul>

      <p>你保留对自己内容的所有权，但应确保内容不侵犯他人权利。</p>

      <h1>十、争议解决</h1>

      <p>
        本隐私政策的解释、效力及争议解决均适用中华人民共和国法律。如发生争议，双方应友好协商解决；协商不成的，任何一方可向本站所在地人民法院提起诉讼。
      </p>

      <h1>十一、政策更新</h1>

      <p>
        本隐私政策可能不定期更新，更新后的政策将在本页面公布。继续使用本站即表示你接受更新后的政策。
      </p>

      <p>
        如有任何疑问，请联系 <strong>{{ contactEmail }}</strong
        >。
      </p>

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

    ul {
      list-style: disc;
      padding-left: 2rem;
      margin: 1rem 0;

      li {
        margin: 0.5rem 0;
      }
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

      ul {
        li {
          font-size: 0.95rem;
        }
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

      ul {
        padding-left: 1.5rem;

        li {
          margin: 0.4rem 0;
          font-size: 0.9rem;
        }
      }

      .update-time {
        font-size: 0.8rem;
      }
    }
  }
}
</style>
