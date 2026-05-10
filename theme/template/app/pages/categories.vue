<script lang="ts" setup>
const { categories, fetchCategories } = useCategories();

await useAsyncData('categories-page', () => fetchCategories(true));

useSeoMeta({
  title: '分类',
  description: '浏览所有文章分类',
});
</script>

<template>
  <div id="page">
    <h1 class="page-title">分类</h1>
    <div v-if="categories.length > 0" class="category-list">
      <NuxtLink
        v-for="category in categories"
        :key="category.id"
        :to="category.url"
        class="category-item"
      >
        {{ category.name }}
        <span class="count">{{ category.count }}</span>
      </NuxtLink>
    </div>
    <div v-else class="empty-state">暂无分类</div>
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

.category-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: var(--flec-heavy-bg);
  border-radius: 8px;
  color: var(--font-color);
  text-decoration: none;
  transition: all 0.2s;

  &:hover {
    background: var(--theme-color);
    color: #fff;

    .count {
      color: #fff;
    }
  }

  .count {
    color: var(--theme-meta-color);
    font-size: 0.9rem;
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
