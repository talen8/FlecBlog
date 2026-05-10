<script lang="ts" setup>
const { tags, fetchTags } = useTags();

await useAsyncData('tags-page', () => fetchTags(true));

useSeoMeta({
  title: '标签',
  description: '浏览所有文章标签',
});
</script>

<template>
  <div id="page">
    <h1 class="page-title">标签</h1>
    <div v-if="tags.length > 0" class="tag-cloud">
      <NuxtLink v-for="tag in tags" :key="tag.id" :to="tag.url" class="tag-item">
        {{ tag.name }}
      </NuxtLink>
    </div>
    <div v-else class="empty-state">暂无标签</div>
  </div>
</template>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border-radius: 12px;
  border: 1px solid var(--flec-border);
  padding: 40px;
}

.page-title {
  margin: 0 0 30px;
  font-size: 2rem;
  font-weight: 700;
}

.tag-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.tag-item {
  padding: 8px 16px;
  background: var(--flec-heavy-bg);
  border-radius: 20px;
  color: var(--font-color);
  text-decoration: none;
  font-size: 0.95rem;
  transition: all 0.2s;

  &:hover {
    background: var(--theme-color);
    color: #fff;
  }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--theme-meta-color);
}

@media screen and (max-width: 768px) {
  #page {
    padding: 20px;
  }

  .page-title {
    font-size: 1.5rem;
  }
}
</style>
