<script lang="ts" setup>
import { renderMarkdown } from '@/utils/markdown';

defineProps<{ content: string }>();
</script>

<template>
  <article class="post-content">
    <!-- eslint-disable-next-line vue/no-v-html -- Markdown 渲染内容，已通过 DOMPurify 净化 -->
    <div class="markdown-content" v-html="renderMarkdown(content || '')" />
  </article>
</template>

<style lang="scss" scoped>
.post-content {
  margin: 20px 0;
}
</style>
