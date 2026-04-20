<script setup lang="ts">
interface Props {
  isScrollingDown: boolean;
  isFixed: boolean;
}

defineProps<Props>();

const { flatNavigationMenus } = useMenus();
const { blogConfig } = useSysConfig();
const { currentArticle } = useCurrentArticle();

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' });
};

const displayTitle = computed(() => {
  return currentArticle.value?.title || blogConfig.value.title;
});
</script>

<template>
  <div class="nav-menu">
    <div class="menu-items" :class="{ hide: isScrollingDown && isFixed }">
      <template v-for="menu in flatNavigationMenus" :key="menu.id">
        <!-- 有子菜单的菜单项 -->
        <div v-if="menu.children && menu.children.length > 0" class="menu-item dropdown">
          <a v-if="menu.url" :href="menu.url" class="brighten" :aria-label="menu.title">
            <i v-if="menu.icon" :class="menu.icon" />
            <span>{{ menu.title }}</span>
            <i class="ri-arrow-down-s-line arrow-icon" />
          </a>

          <span v-else class="brighten menu-label">
            <i v-if="menu.icon" :class="menu.icon" />
            <span>{{ menu.title }}</span>
            <i class="ri-arrow-down-s-line arrow-icon" />
          </span>

          <!-- 下拉菜单 -->
          <ul class="dropdown-menu">
            <li v-for="child in menu.children" :key="child.id">
              <a :href="child.url" :aria-label="child.title">
                <i v-if="child.icon" :class="child.icon" />
                <span>{{ child.title }}</span>
              </a>
            </li>
          </ul>
        </div>

        <!-- 无子菜单的菜单项 -->
        <a v-else :href="menu.url" class="brighten" :aria-label="menu.title">
          <i v-if="menu.icon" :class="menu.icon" />
          <span>{{ menu.title }}</span>
        </a>
      </template>
    </div>
    <div class="scroll-title" :class="{ show: isScrollingDown && isFixed }">
      <a
        href="#"
        class="scroll-to-top brighten no-after"
        aria-label="回到顶部"
        @click.prevent="scrollToTop"
      >
        <span class="title" aria-hidden="true">{{ displayTitle }}</span>
      </a>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

.nav-menu {
  flex: 3;
  display: flex;
  justify-content: center;
  gap: 2rem;
  position: relative;

  .menu-items {
    display: flex;
    align-items: center;
    gap: 1rem;
    opacity: 1;
    transform: translateY(0);
    transition: all 0.3s ease;

    &.hide {
      opacity: 0;
      transform: translateY(-20px);
      pointer-events: none;
    }

    a,
    .menu-label {
      margin: 0 0.5rem;
      display: flex;
      align-items: center;
      gap: 0.3rem;
      white-space: nowrap;
      cursor: pointer;

      i {
        font-size: 1rem;
      }

      .arrow-icon {
        font-size: 1.1rem;
        transition: transform 0.3s ease;
      }
    }

    // 下拉菜单容器
    .menu-item.dropdown {
      position: relative;

      &:hover {
        .dropdown-menu {
          visibility: visible;
          opacity: 1;
          transform: translateX(-50%) translateY(0);
          pointer-events: auto;

          li {
            opacity: 1;
            transform: translateY(0);
          }
        }

        .arrow-icon {
          transform: rotate(180deg);
        }
      }

      .dropdown-menu {
        @extend .cardHover;
        visibility: hidden;
        backdrop-filter: blur(30px);
        position: absolute;
        left: 50%;
        margin-top: 15px;
        padding: 6px;
        min-width: max-content;
        opacity: 0;
        transform: translateX(-50%) translateY(-10px);
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
        pointer-events: none;

        &::before {
          position: absolute;
          top: -20px;
          left: 50%;
          width: 80%;
          height: 30px;
          content: '';
          transform: translateX(-50%);
        }

        li {
          float: left;
          list-style: none;
          white-space: nowrap;
          opacity: 0;
          transform: translateY(-5px);
          transition: all 0.2s ease;

          @for $i from 1 through 10 {
            &:nth-child(#{$i}) {
              transition-delay: #{$i * 0.03}s;
            }
          }

          a {
            display: inline-block;
            padding: 4px 14px;
            margin: 0;
            width: 100%;
            color: var(--flec-nav-fixed-font);
            text-shadow: none !important;
            transition: all 0.2s ease;

            &:hover {
              color: var(--flec-nav-fixed-font-hover);
              background: var(--flec-nav-menu-bg-hover);
              border-radius: 12px;
            }

            i {
              margin-right: 6px;
            }
          }
        }
      }
    }
  }

  .scroll-title {
    width: 100%;
    position: absolute;
    opacity: 0;
    pointer-events: none;
    transform: translateY(20px);
    transition: all 0.3s ease;

    &.show {
      opacity: 1;
      pointer-events: auto;
      transform: translateY(0);
    }

    .scroll-to-top {
      display: flex;
      justify-content: center;
      width: 100%;
      text-align: center;

      .title {
        display: inline;
      }

      .top {
        display: none;
      }
    }
  }
}

@media screen and (max-width: 900px) {
  .nav-menu {
    display: none;
  }
}
</style>
