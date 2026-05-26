<script setup lang="ts">
const route = useRoute();

// 判断当前页面类型
const isArticlePage = computed(() => route.meta.typeHeader === 'post');
</script>

<template>
  <aside id="sidebar">
    <LayoutsSidebarAuthorCard />
    <LayoutsSidebarAnnouncementCard />
    <div class="sticky-sidebar">
      <template v-if="isArticlePage">
        <ClientOnly>
          <LayoutsSidebarTocCard />
        </ClientOnly>
      </template>
      <template v-else>
        <LayoutsSidebarCategoriesCard />
        <LayoutsSidebarTagsCard />
        <LayoutsSidebarArchivesCard />
        <LayoutsSidebarWebInfoCard />
      </template>
    </div>
  </aside>
</template>

<style lang="scss">
@use '@/assets/css/mixins' as *;

#sidebar {
  width: 26%;
  padding-left: 15px;

  .card-widget {
    position: relative;
    overflow: hidden;
    margin-bottom: 20px;
    padding: 20px 24px;
    @extend .cardHover;
  }

  .sticky-sidebar {
    position: sticky;
    top: 70px;
    transition: top 0.3s;
    display: flex;
    flex-direction: column;
  }

  .item-headline {
    padding-bottom: 6px;
    font-size: 1.2rem;
    display: flex;
    align-items: center;

    i {
      margin-right: 10px;
    }

    .more-link {
      margin-left: auto;
      color: var(--font-color);
      opacity: 0.6;
      transition: all 0.3s;
      font-size: 1.3em;

      &:hover {
        opacity: 1;
        color: var(--flec-btn-hover);
      }
    }
  }

  .card-list {
    margin: 0;
    padding: 0;

    .card-list-link {
      display: flex;
      flex-direction: row;
      margin: 2px 0;
      padding: 2px 8px;
      transition: all 0.3s;
      border-radius: 6px;

      &:hover {
        padding: 2px 12px;
        background-color: var(--flec-btn-hover);
        color: #fff;
      }

      span {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .card-list-name {
        flex: 1;
      }
    }
  }

  // 共享的折叠图标样式
  .collapse-icon {
    display: none;
    margin-left: 8px;
    cursor: pointer;
    font-size: 1.2em;
    transition: transform 0.3s ease;

    &.is-expanded {
      transform: rotate(-90deg);
    }
  }

  // 响应式设置
  @media screen and (max-width: 900px) {
    .collapse-icon {
      display: inline-block;
    }

    .item-headline {
      padding-bottom: 0 !important;
      transition: padding-bottom 0.3s ease;

      &.is-expanded {
        padding-bottom: 8px !important;
      }
    }

    .card-list,
    .card-tag-cloud {
      max-height: 0;
      overflow: hidden;
      transition: max-height 0.3s ease;

      &.is-expanded {
        max-height: 1000px;
      }
    }
  }
}

// 响应式设置
@media screen and (max-width: 900px) {
  #sidebar {
    width: 100%;
    margin-top: 20px;
    padding-left: 0;

    .card-widget {
      margin-bottom: 12px;
      padding: 16px 18px;
    }

    .item-headline {
      padding-bottom: 8px;
      font-size: 1.1rem;

      i {
        margin-right: 8px;
      }
    }

    .sticky-sidebar {
      position: static;
    }
  }
}
</style>
