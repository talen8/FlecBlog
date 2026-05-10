<script lang="ts" setup>
defineProps<{ article: { tags: Array<{ id: number; name: string; url: string }> } }>();
</script>

<template>
  <div v-if="article.tags?.length" class="article-tags">
    <NuxtLink v-for="tag in article.tags" :key="tag.id" :to="tag.url">
      {{ tag.name }}
    </NuxtLink>
  </div>
</template>

<style lang="scss" scoped>
.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid var(--flec-border);

  a {
    padding: 4px 12px;
    border: 1px solid var(--theme-color);
    color: var(--theme-color);
    text-decoration: none;
    font-size: 0.85rem;
  }
}
</style>
