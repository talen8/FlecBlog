<script lang="ts" setup>
import type { ArticleNav } from '@@/types';

defineProps<{
  prev?: ArticleNav;
  next?: ArticleNav;
}>();
</script>

<template>
  <div v-if="prev || next" class="article-nav">
    <NuxtLink v-if="prev" :to="prev.url" class="nav-item prev">
      <div class="label">← 上一篇</div>
      <div class="title">{{ prev.title }}</div>
    </NuxtLink>
    <NuxtLink v-if="next" :to="next.url" class="nav-item next">
      <div class="label">下一篇 →</div>
      <div class="title">{{ next.title }}</div>
    </NuxtLink>
  </div>
</template>

<style lang="scss" scoped>
.article-nav {
  display: flex;
  gap: 12px;
  margin: 30px 0;
}

.nav-item {
  flex: 1;
  padding: 16px;
  border: 1px solid var(--flec-border);
  text-decoration: none;
  color: var(--font-color);

  &.prev {
    text-align: left;
  }
  &.next {
    text-align: right;
  }

  .label {
    font-size: 0.85rem;
    color: var(--theme-meta-color);
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
