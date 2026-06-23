<script lang="ts" setup>
definePageMeta({
  showSidebar: false,
});

const { categories } = useCategories();

useSeoMeta({
  title: '分类',
  description: '浏览所有文章分类，探索不同主题的技术与生活内容',
});
</script>

<template>
  <!-- 内容区域 -->
  <div id="page">
    <h1 class="page-title">分类</h1>
    <div class="category-lists">
      <ul class="category-list">
        <li v-for="category in categories" :key="category.id" class="category-list-item">
          <router-link class="category-list-link" :to="category.url">{{
            category.name
          }}</router-link>
          <span class="category-list-count">{{ category.count }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<style lang="scss">
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

  .category-lists {
    .category-list {
      margin-bottom: 0;
      padding-left: 20px;

      .category-list-item {
        position: relative;
        margin: 6px 0;
        padding: 0.12em 0.4em 0.12em 1.4em;

        &::before {
          position: absolute;
          top: 0.8em;
          left: 0;
          width: 0.8em;
          height: 0.8em;
          border: 0.215em solid #49b1f5;
          border-radius: 50%;
          background: transparent;
          content: '';
          cursor: pointer;
          transition: all 0.3s ease-out;
        }

        .category-list-count {
          margin-left: 8px;
          color: var(--theme-meta-color);

          &::after {
            content: ')';
          }

          &::before {
            content: '(';
          }
        }
      }
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
  }
}

@media screen and (max-width: 768px) {
  #page {
    padding: 18px;

    .page-title {
      font-size: 1.4rem;
    }

    .category-lists {
      .category-list {
        padding-left: 12px;

        .category-list-item {
          padding: 0.1em 0.3em 0.1em 1.2em;
          font-size: 0.92rem;

          &::before {
            top: 0.7em;
            width: 0.7em;
            height: 0.7em;
          }
        }
      }
    }
  }
}
</style>
