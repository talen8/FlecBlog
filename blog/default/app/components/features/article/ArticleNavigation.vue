<script lang="ts" setup>
import type { ArticleNav } from '@@/types/article';

defineProps<{
  prev?: ArticleNav;
  next?: ArticleNav;
}>();
</script>

<template>
  <div v-if="prev || next" class="article-nav">
    <a v-if="prev" :href="prev.url" class="nav-item prev">
      <div class="label">← 上一篇</div>
      <div class="title">{{ prev.title }}</div>
    </a>
    <a v-if="next" :href="next.url" class="nav-item next">
      <div class="label">下一篇 →</div>
      <div class="title">{{ next.title }}</div>
    </a>
  </div>
</template>

<style lang="scss" scoped>
.article-nav {
  display: flex;
  gap: 12px;
  margin: 20px 0;
}

.nav-item {
  flex: 1;
  padding: 12px 16px;
  border: 1px solid var(--flec-border);
  border-radius: 6px;
  text-decoration: none;
  color: var(--font-color);
  transition: border-color 0.2s;

  &:hover {
    border-color: var(--theme-color);
  }

  &.prev {
    text-align: left;
  }

  &.next {
    text-align: right;
  }

  .label {
    font-size: 0.85rem;
    opacity: 0.7;
    margin-bottom: 6px;
  }

  .title {
    font-size: 0.95rem;
    font-weight: 500;
  }
}

@media screen and (max-width: 600px) {
  .article-nav {
    flex-direction: column;
  }
}
</style>
