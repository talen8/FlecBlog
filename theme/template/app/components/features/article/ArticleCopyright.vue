<script lang="ts" setup>
import type { Article } from '@@/types';

const props = defineProps<{ article: Article }>();
const { basicConfig } = useSysConfig();

const articleUrl = computed(() => {
  const blogUrl = basicConfig.value?.['basic.blog_url'];
  return blogUrl ? `${blogUrl}${props.article.url}` : props.article.url;
});

const author = computed(() => basicConfig.value?.['basic.author'] || '作者');
</script>

<template>
  <div class="post-copyright">
    <div class="copyright-title">{{ article.title }}</div>
    <div class="copyright-link">{{ articleUrl }}</div>
    <div class="copyright-info">
      <div class="info-item">
        <span class="label">作者</span>
        <span class="value">{{ author }}</span>
      </div>
      <div class="info-item">
        <span class="label">发布时间</span>
        <span class="value">{{ formatDate(article.publish_time) }}</span>
      </div>
      <div class="info-item">
        <span class="label">许可协议</span>
        <span class="value">
          <a
            href="https://creativecommons.org/licenses/by-nc-sa/4.0/"
            target="_blank"
            rel="noopener noreferrer"
          >
            CC BY-NC-SA 4.0
          </a>
        </span>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.post-copyright {
  margin: 30px 0;
  padding: 16px 20px;
  background: var(--flec-heavy-bg);
  border-left: 4px solid var(--theme-color);
}

.copyright-title {
  font-size: 1.1rem;
  font-weight: 500;
  margin-bottom: 6px;
}

.copyright-link {
  color: var(--theme-meta-color);
  font-size: 0.85rem;
  margin-bottom: 12px;
  word-break: break-all;
}

.copyright-info {
  display: flex;
  gap: 24px;
  flex-wrap: wrap;
}

.info-item {
  font-size: 0.9rem;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.label {
  color: var(--theme-meta-color);
  font-size: 0.8rem;
}

.value a {
  color: var(--theme-color);
  text-decoration: none;
}
</style>
