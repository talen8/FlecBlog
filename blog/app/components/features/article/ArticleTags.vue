<script lang="ts" setup>
defineProps<{ article: { tags: Array<{ id: number; name: string; url: string }> } }>();
</script>

<template>
  <div v-if="article.tags?.length" class="article-tags">
    <a v-for="tag in article.tags" :key="tag.id" :href="tag.url">{{ tag.name }}</a>
  </div>
</template>

<style lang="scss" scoped>
.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;

  a {
    padding: 1px 8px;
    border: 1px solid var(--theme-color);
    border-radius: 8px;
    color: var(--theme-color);
    text-decoration: none;
    font-size: 0.85rem;

    &:hover {
      background: var(--theme-color);
      color: #fff;
    }
  }
}
</style>
