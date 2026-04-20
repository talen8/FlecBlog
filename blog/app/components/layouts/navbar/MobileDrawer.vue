<script setup lang="ts">
const { basicConfig } = useSysConfig();
const avatarUrl = computed(() => basicConfig.value.author_avatar || '/avatar.webp');

interface Props {
  modelValue: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>();

const { navigationMenus, aggregateMenus } = useMenus();
const { total: articleCount } = useArticles();
const { total: categoryCount } = useCategories();
const { total: tagCount } = useTags();
const expandedMenus = ref<Set<number>>(new Set());
const currentSlide = ref(0);
const slideWrapper = ref<HTMLElement>();

const topAggregateMenus = computed(() =>
  aggregateMenus.value.filter((menu: { parent_id: number | null }) => !menu.parent_id)
);

const isIconUrl = (icon: string) => {
  return (
    icon && (icon.startsWith('http://') || icon.startsWith('https://') || icon.startsWith('/'))
  );
};

const close = () => emit('update:modelValue', false);

const toggleSubmenu = (menuId: number) => {
  if (expandedMenus.value.has(menuId)) {
    expandedMenus.value.delete(menuId);
  } else {
    expandedMenus.value.add(menuId);
  }
};

const switchSlide = (index: number) => {
  currentSlide.value = index;
};

const { direction } = useSwipe(slideWrapper, {
  threshold: 50,
  onSwipeEnd() {
    if (direction.value === 'left' && currentSlide.value < 1) {
      currentSlide.value++;
    } else if (direction.value === 'right' && currentSlide.value > 0) {
      currentSlide.value--;
    }
  },
});

watch(
  () => props.modelValue,
  isOpen => {
    document.body.style.overflow = isOpen ? 'hidden' : '';
  }
);
</script>

<template>
  <Teleport to="body">
    <Transition name="drawer">
      <div v-if="modelValue" class="drawer-overlay" @click="close">
        <div class="drawer-container" @click.stop>
          <!-- 头像 -->
          <div class="avatar-img">
            <NuxtImg :src="avatarUrl" alt="avatar" loading="lazy" />
          </div>

          <!-- 站点数据 -->
          <div class="site-data">
            <a href="/archive" @click="close">
              <div class="headline">文章</div>
              <div class="length-num">{{ articleCount }}</div>
            </a>
            <a href="/tags" @click="close">
              <div class="headline">标签</div>
              <div class="length-num">{{ tagCount }}</div>
            </a>
            <a href="/categories" @click="close">
              <div class="headline">分类</div>
              <div class="length-num">{{ categoryCount }}</div>
            </a>
          </div>

          <!-- 滑动菜单容器 -->
          <div class="sidebar-menu">
            <div class="slide-wrapper">
              <div
                ref="slideWrapper"
                class="slide-box"
                :style="{ transform: `translateX(-${currentSlide * 50}%)` }"
              >
                <!-- 第一页：导航菜单 -->
                <div class="menus-wrapper">
                  <template v-for="menu in navigationMenus" :key="menu.id">
                    <template v-if="menu.children?.length">
                      <div class="nav-item parent-item" @click="toggleSubmenu(menu.id)">
                        <i v-if="menu.icon && !isIconUrl(menu.icon)" :class="menu.icon" />
                        <span>{{ menu.title }}</span>
                        <i
                          class="ri-arrow-right-s-line"
                          :class="{ rotate: expandedMenus.has(menu.id) }"
                        />
                      </div>
                      <Transition name="submenu">
                        <div v-show="expandedMenus.has(menu.id)" class="submenu">
                          <a
                            v-for="child in menu.children"
                            :key="child.id"
                            :href="child.url"
                            class="nav-item"
                            @click="close"
                          >
                            <span>{{ child.title }}</span>
                          </a>
                        </div>
                      </Transition>
                    </template>
                    <a v-else :href="menu.url" class="nav-item" @click="close">
                      <i v-if="menu.icon && !isIconUrl(menu.icon)" :class="menu.icon" />
                      <span>{{ menu.title }}</span>
                    </a>
                  </template>
                </div>

                <!-- 第二页：聚合菜单 -->
                <div class="aggregate-wrapper">
                  <div
                    v-for="menu in topAggregateMenus"
                    v-show="menu.children?.length"
                    :key="menu.id"
                  >
                    <div class="section-title">{{ menu.title }}</div>
                    <div class="aggregate-grid">
                      <a
                        v-for="child in menu.children"
                        :key="child.id"
                        :href="child.url"
                        class="aggregate-item"
                        @click="close"
                      >
                        <NuxtImg
                          v-if="child.icon && isIconUrl(child.icon)"
                          :src="child.icon"
                          :alt="child.title"
                          loading="lazy"
                        />
                        <i v-else-if="child.icon" :class="child.icon" />
                        <span>{{ child.title }}</span>
                      </a>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 滑动指示器 -->
              <div class="slide-indicator">
                <div
                  class="indicator-dot"
                  :class="{ active: currentSlide === 0 }"
                  @click="switchSlide(0)"
                />
                <div
                  class="indicator-dot"
                  :class="{ active: currentSlide === 1 }"
                  @click="switchSlide(1)"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
.drawer-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 9999;
  backdrop-filter: blur(4px);
}

.drawer-container {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  width: 350px;
  max-width: 80vw;
  background-color: var(--flec-card-bg);
  box-shadow: -2px 0 8px rgba(0, 0, 0, 0.15);
  overflow-y: auto;
  transition: transform 0.3s ease-in;

  .avatar-img {
    width: 110px;
    margin: 20px auto;
    text-align: center;

    img {
      border-radius: 50%;
      object-fit: cover;
    }
  }

  .site-data {
    display: flex;
    justify-content: space-around;
    padding: 0 10px;

    a {
      text-align: center;
      text-decoration: none;
      transition: all 0.3s;

      &:hover {
        color: var(--theme-color);
      }

      .headline {
        font-size: 0.85rem;
      }

      .length-num {
        font-size: 1.2rem;
      }
    }
  }

  .sidebar-menu {
    margin: 20px;
    padding: 15px;
    background: var(--flec-card-bg);
    box-shadow: 0 0 1px 1px rgba(7, 17, 27, 0.05);
    border-radius: 10px;

    .slide-wrapper {
      position: relative;
      overflow: hidden;

      .slide-box {
        display: flex;
        transition: transform 0.3s ease;
        width: 200%;

        .menus-wrapper,
        .aggregate-wrapper {
          width: 50%;
          flex-shrink: 0;
        }
      }

      .slide-indicator {
        display: flex;
        justify-content: center;
        margin-top: 15px;
        gap: 8px;

        .indicator-dot {
          width: 8px;
          height: 8px;
          border-radius: 50%;
          background: var(--font-color);
          opacity: 0.3;
          transition: opacity 0.3s;
          cursor: pointer;

          &.active {
            opacity: 1;
          }
        }
      }
    }

    .section-title {
      padding: 6px 0;
      color: var(--font-color);
      font-weight: bold;
      font-size: 0.9rem;
    }

    .nav-item {
      display: flex;
      align-items: center;
      padding: 0 12px;
      margin: 2px 0;
      color: var(--font-color);
      text-decoration: none;
      border-radius: 6px;
      transition: all 0.2s;
      cursor: pointer;

      i {
        margin-right: 10px;
        font-size: 1.1rem;
        width: 20px;
      }

      span {
        flex: 1;
      }

      &:hover {
        background: var(--flec-nav-menu-bg-hover);
        color: #fff;
      }

      &.parent-item {
        i:last-child {
          margin-left: auto;
          margin-right: 0;
          transition: transform 0.3s;

          &.rotate {
            transform: rotate(90deg);
          }
        }
      }
    }

    .submenu {
      margin: 0;
      padding-left: 25px;
      list-style: none;

      .nav-item {
        padding: 0 23px 0 15px;
      }
    }

    .aggregate-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 4px;

      .aggregate-item {
        display: flex;
        align-items: center;
        padding: 2px 8px;
        color: var(--font-color);
        text-decoration: none;
        border-radius: 8px;
        transition: all 0.3s;

        &:hover {
          background: var(--flec-nav-menu-bg-hover);

          span {
            color: #fff;
          }
        }

        img,
        i {
          margin-right: 8px;
          flex-shrink: 0;
        }

        img {
          border-radius: 4px;
          object-fit: cover;
        }

        i {
          font-size: 1.2rem;
        }

        span {
          font-size: 0.9rem;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }
  }
}

// 子菜单动画
.submenu-enter-active,
.submenu-leave-active {
  transition: all 0.3s;
  overflow: hidden;
}

.submenu-enter-from,
.submenu-leave-to {
  opacity: 0;
  max-height: 0;
}

.submenu-enter-to,
.submenu-leave-from {
  opacity: 1;
  max-height: 500px;
}

// 抽屉动画
.drawer-enter-active,
.drawer-leave-active {
  transition: opacity 0.3s ease;
}

.drawer-enter-from,
.drawer-leave-to {
  opacity: 0;
}

.drawer-enter-from .drawer-container,
.drawer-leave-to .drawer-container {
  transform: translateX(100%);
}

.drawer-enter-to .drawer-container,
.drawer-leave-from .drawer-container {
  transform: translateX(0);
}
</style>
