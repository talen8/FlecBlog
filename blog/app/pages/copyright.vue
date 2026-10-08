<script lang="ts" setup>
const { basicConfig } = useSysConfig();
const { themeConfig } = useTheme();

const license = computed(() => getLicense(themeConfig.value.copyright_license));

const contactEmail = computed(() => basicConfig.value.author_email || '');

definePageMeta({
  showSidebar: false,
});

const { title } = usePageSeo({
  path: '/copyright',
  title: '版权协议',
  description: () => `了解本站内容的版权归属和使用许可，包括 ${license.value.short} 协议`,
});
</script>

<template>
  <div id="page">
    <div class="page-title">{{ title }}</div>

    <div class="content">
      <p>本页面说明本站内容的版权归属和使用许可。</p>

      <p>
        本站所有原创内容采用
        <a :href="license.url" target="_blank" rel="noopener noreferrer">{{ license.name }}</a
        >（{{ license.short }}）进行许可。
      </p>

      <h1>版权所有者</h1>

      <p v-if="license.clauses.includes('cc0')">
        本站原创内容已通过
        {{ license.short }}
        贡献至公有领域，作者在法律允许的最大范围内放弃全部著作权，任何人可自由使用。用户发表的评论版权归用户本人所有，但授予本站展示和使用的权利。
      </p>
      <p v-else>
        本站所有原创文章、图片、设计的版权归本站所有。用户发表的评论版权归用户本人所有，但授予本站展示和使用的权利。
      </p>

      <h1>许可协议说明</h1>

      <p>{{ license.short }} 协议包含以下要点：</p>

      <template v-if="license.clauses.includes('by')">
        <h2>署名（BY）</h2>
        <p>转载、引用或改编时必须注明作者与原文链接。</p>
        <ul>
          <li>可在文章开头或明显位置添加超链接，格式如「本文转载自 [原文标题](原文链接)」</li>
          <li>不得删除或隐匿作者署名与原文出处</li>
          <li>不得让人误以为内容由转载者原创</li>
        </ul>
      </template>

      <template v-if="license.clauses.includes('nc')">
        <h2>非商业性使用（NC）</h2>
        <p>禁止将本站内容用于任何商业目的。</p>
        <ul>
          <li>个人学习、研究与非盈利性质的分享不受限制</li>
          <li>不得在转载页面插入广告（如 Google AdSense、百度联盟等）获取收益</li>
          <li>不得要求付费、关注公众号或下载 App 后才能查看</li>
          <li>不得用于商业培训、出版、营销等盈利行为</li>
        </ul>
      </template>

      <template v-if="license.clauses.includes('nd')">
        <h2>禁止演绎（ND）</h2>
        <p>只能原样传播，不得修改、改编或基于原文创作衍生作品。</p>
        <ul>
          <li>可完整转载，或作为参考资料引用片段并注明出处</li>
          <li>不得修改、删减原文后发布</li>
          <li>不得翻译成其他语言后发布</li>
          <li>不得基于原文创作衍生作品</li>
        </ul>
      </template>

      <template v-if="license.clauses.includes('sa')">
        <h2>相同方式共享（SA）</h2>
        <p>改编后的作品必须以相同协议开放共享。</p>
        <ul>
          <li>可修改、混编或基于原文创作衍生作品，但需注明原始出处</li>
          <li>衍生作品必须沿用相同协议开放</li>
          <li>不得为衍生作品附加更严格的授权限制</li>
        </ul>
      </template>

      <template v-if="license.clauses.includes('cc0')">
        <h2>公有领域（CC0）</h2>
        <p>作者在法律允许的最大范围内放弃全部著作权，内容进入公有领域。</p>
        <ul>
          <li>可任意复制、修改、发行与演绎，包括商业用途</li>
          <li>无需署名，但仍建议注明来源</li>
          <li>不得主张对内容本身的著作权</li>
          <li>商标权、专利权不在放弃范围内</li>
        </ul>
      </template>

      <h1>受保护的内容</h1>

      <p v-if="license.clauses.includes('cc0')">
        以下内容已随 {{ license.short }} 一并贡献至公有领域，可自由使用：
      </p>
      <p v-else>以下内容受版权保护：</p>
      <ul>
        <li>本站所有原创文章（标题、正文、代码示例）</li>
        <li>文章配图和封面图片</li>
        <li>网站设计和页面布局</li>
        <li>原创图标和素材</li>
      </ul>

      <p>
        <strong>不受保护的内容：</strong
        >用户发表的评论归用户本人所有，但用户授予本站展示和使用的权利。
      </p>

      <h1>特殊许可</h1>

      <p>在以下情况下，可能获得额外授权：</p>

      <ul>
        <li><strong>学术研究</strong>：用于非营利性学术研究和教育目的，不做任何限制，可自由使用</li>
        <li><strong>友链博主</strong>：被本站友链收录的博客可享有更宽松的引用权限</li>
        <li><strong>个别授权</strong>：如有其他特殊需求，可通过评论或邮件联系协商</li>
      </ul>

      <p>
        联系方式：<strong>{{ contactEmail }}</strong>
      </p>

      <h1>侵权处理</h1>

      <p>
        如发现未经授权的使用，我们保留追究法律责任的权利。同时，我们也尊重他人的版权，如认为本站内容侵犯了你的权利，请联系我们删除。
      </p>

      <h1>免责声明</h1>

      <p>
        本站内容仅供学习和参考，不对内容的准确性、完整性和时效性做任何保证。因使用本站内容造成的任何损失，本站不承担责任。
      </p>

      <p class="update-time">最后更新时间：2026年10月</p>
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
      margin: 2rem 0 1rem;
      color: var(--font-color);

      &:first-child {
        margin-top: 1.5rem;
      }
    }

    h2 {
      font-size: 1.25rem;
      font-weight: 600;
      margin: 1.75rem 0 0.875rem;
      color: var(--font-color);
    }

    p {
      margin: 1rem 0;
      text-align: justify;
    }

    a {
      color: var(--theme-color);
      text-decoration: none;
      font-weight: 500;

      &:hover {
        text-decoration: underline;
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
        margin: 1.75rem 0 0.85rem;
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
        margin: 1.5rem 0 0.75rem;
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
