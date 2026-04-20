<script lang="ts" setup>
import type { Moment } from '@@/types/moment';
import { getMoments } from '@/composables/api/moment';

// 获取动态数据（支持 SSR）
const { data: moments } = await useAsyncData('moments-widget', async () => {
  try {
    const { list } = await getMoments({ page: 1, page_size: 10 });
    return list || [];
  } catch (error) {
    console.error('获取动态列表失败:', error);
    return [];
  }
});

// 当前显示的动态索引
const currentIndex = ref(0);

// 当前显示的动态
const currentMoment = computed(() => moments.value?.[currentIndex.value]);

// 获取动态包含的内容类型
const getContentTypes = (moment: Moment) => {
  const types: string[] = [];
  if (moment.content.images && moment.content.images.length > 0) {
    types.push('image');
  }
  if (moment.content.video) {
    types.push('video');
  }
  if (moment.content.link) {
    types.push('link');
  }
  if (moment.content.music) {
    types.push('music');
  }
  return types;
};

// 下一条
const nextMoment = () => {
  currentIndex.value = (currentIndex.value + 1) % moments.value!.length;
};

// 自动轮播
let timer: number | null = null;

const startAutoPlay = () => {
  stopAutoPlay();
  timer = window.setInterval(nextMoment, 3000);
};

const stopAutoPlay = () => {
  if (timer) {
    clearInterval(timer);
    timer = null;
  }
};

// 初始化
onMounted(() => {
  if (moments.value && moments.value.length > 1) {
    startAutoPlay();
  }
});

onUnmounted(() => {
  stopAutoPlay();
});

// 鼠标悬停停止轮播
const onMouseEnter = stopAutoPlay;
const onMouseLeave = () => {
  if (moments.value && moments.value.length > 1) {
    startAutoPlay();
  }
};
</script>

<template>
  <div
    v-if="moments?.length"
    class="moment-widget"
    @mouseenter="onMouseEnter"
    @mouseleave="onMouseLeave"
  >
    <NuxtLink to="/moment" class="moment-container">
      <!-- 左侧图标 -->
      <div class="widget-icon">
        <i class="ri-send-ins-line" />
      </div>

      <!-- 中间滚动内容 -->
      <div class="widget-center">
        <Transition name="slide" mode="out-in">
          <div v-if="currentMoment" :key="currentIndex" class="moment-content-wrapper">
            <span class="moment-content">
              {{ currentMoment.content.text }}
            </span>
            <span class="content-icons">
              <template v-for="type in getContentTypes(currentMoment)" :key="type">
                <i v-if="type === 'image'" class="ri-image-fill" />
                <i v-if="type === 'video'" class="ri-video-fill" />
                <i v-if="type === 'link'" class="ri-link" />
                <i v-if="type === 'music'" class="ri-music-2-fill" />
              </template>
            </span>
          </div>
        </Transition>
      </div>

      <!-- 右侧箭头 -->
      <div class="widget-icon">
        <i class="ri-arrow-right-s-line" />
      </div>
    </NuxtLink>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

.moment-widget {
  margin: 0 auto 20px;
  padding: 0 15px;
  max-width: 1200px;

  .moment-container {
    @extend .cardHover;
    width: 100%;
    height: 45px;
    display: flex;
    align-items: center;
  }

  .widget-center {
    flex: 1;
    padding: 0 100px;
    display: flex;
    align-items: center;
    justify-content: center;
    min-width: 0;

    .moment-content-wrapper {
      display: inline-flex;
      align-items: center;
      gap: 8px;
      overflow: hidden;
      max-width: 100%;
    }

    .moment-content {
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      max-width: 100%;
    }

    .content-icons {
      display: flex;
      align-items: center;
      gap: 4px;
      flex-shrink: 0;

      i {
        display: flex;
        align-items: center;
      }
    }
  }

  .widget-icon {
    flex-shrink: 0;
    width: 40px;
    display: flex;
    align-items: center;
    justify-content: center;

    i {
      font-size: 1.3rem;
    }
  }
}

.slide-enter-active,
.slide-leave-active {
  transition: all 0.3s ease;
}

.slide-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.slide-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

@media screen and (max-width: 900px) {
  .moment-widget {
    padding: 0 12px;

    .moment-container {
      height: 40px;
    }

    .widget-icon {
      width: 36px;

      i {
        font-size: 1.2rem;
      }
    }

    .widget-center {
      padding: 0 8px;

      .moment-content {
        font-size: 0.85rem;
      }

      .content-icons {
        gap: 3px;

        i {
          font-size: 0.8rem;
        }
      }
    }
  }
}
</style>
