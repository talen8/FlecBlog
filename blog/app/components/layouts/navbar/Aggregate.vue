<script setup lang="ts">
const { aggregateMenus } = useMenus();

// 获取顶级聚合菜单（已包含子菜单）
const topAggregateMenus = computed(() => {
  return aggregateMenus.value.filter(menu => !menu.parent_id);
});

// 判断 icon 是否为图片 URL
const isImageIcon = (icon: string) => {
  return (
    icon && (icon.startsWith('http://') || icon.startsWith('https://') || icon.startsWith('/'))
  );
};
</script>

<template>
  <div v-if="topAggregateMenus.length > 0" class="nav-aggregate">
    <div class="aggregate-trigger brighten">
      <i class="ri-apps-line ri-lg" />
    </div>

    <!-- 聚合下拉菜单 -->
    <div class="aggregate-dropdown">
      <div
        v-for="menu in topAggregateMenus"
        v-show="menu.children && menu.children.length > 0"
        :key="menu.id"
        class="aggregate-group"
      >
        <!-- 主菜单标题 -->
        <div class="group-title">
          <span>{{ menu.title }}</span>
        </div>

        <!-- 子菜单列 -->
        <div class="group-children">
          <a
            v-for="child in menu.children"
            :key="child.id"
            :href="child.url"
            :aria-label="child.title"
          >
            <NuxtImg
              v-if="child.icon && isImageIcon(child.icon)"
              :src="child.icon"
              :alt="child.title"
              class="icon-img"
              loading="lazy"
            />
            <i v-else-if="child.icon" :class="child.icon + ' ri-lg'" />
            <span>{{ child.title }}</span>
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

.nav-aggregate {
  position: relative;
  margin-right: 30px;

  .aggregate-trigger {
    cursor: pointer;
  }

  .aggregate-dropdown {
    @extend .cardHover;
    backdrop-filter: blur(30px);
    position: absolute;
    left: 0;
    margin-top: 15px;
    padding: 8px;
    width: 300px;
    opacity: 0;
    visibility: hidden;
    transform: translateY(-10px);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    pointer-events: none;

    &::before {
      position: absolute;
      top: -20px;
      left: 0;
      width: 100%;
      height: 30px;
      content: '';
    }

    .aggregate-group {
      margin-bottom: 12px;

      &:last-child {
        margin-bottom: 0;
      }

      .group-title {
        padding: 6px 10px;
        color: var(--flec-nav-fixed-font);
        font-weight: bold;
        font-size: 0.8rem;
        transition: color 0.2s ease;

        &:hover {
          color: var(--flec-nav-fixed-font-hover);
        }
      }

      .group-children {
        display: grid;
        grid-template-columns: repeat(2, 1fr);
        gap: 4px;

        a {
          height: 40px;
          display: flex;
          align-items: center;
          padding: 8px 10px;
          color: var(--flec-nav-fixed-font);
          font-size: 0.9rem;
          opacity: 0;
          transform: translateY(-5px);
          transition: all 0.2s ease;

          &:hover {
            background: var(--flec-nav-menu-bg-hover);
            border-radius: 8px;

            span {
              color: var(--flec-nav-fixed-font-hover);
            }
          }

          img {
            width: 20px;
            height: 20px;
            margin-right: 8px;
            border-radius: 4px;
            object-fit: cover;
          }

          i {
            margin-right: 8px;
          }

          span {
            color: var(--flec-nav-fixed-font);
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
          }
        }
      }
    }
  }

  &:hover {
    .aggregate-dropdown {
      opacity: 1;
      visibility: visible;
      transform: translateY(0);
      pointer-events: auto;

      .group-children a {
        opacity: 1;
        transform: translateY(0);
      }
    }
  }
}

@media screen and (max-width: 768px) {
  .nav-aggregate {
    display: none;
  }
}
</style>
