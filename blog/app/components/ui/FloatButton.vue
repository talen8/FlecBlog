<template>
  <!-- 悬浮按钮组 -->
  <Transition name="fade">
    <div v-show="visible" class="float-button-group">
      <!-- 主题切换按钮 -->
      <div class="float-button theme-toggle" title="切换主题" @click="toggleTheme">
        <i class="ri-moon-line theme-icon-moon" />
        <i class="ri-sun-line theme-icon-sun" />
      </div>

      <!-- 阅读模式按钮（仅文章页显示） -->
      <div v-if="isArticlePage" class="float-button" title="阅读模式" @click="toggleReadingMode">
        <i class="ri-book-open-line" />
      </div>

      <!-- 目录按钮（仅移动端文章页显示） -->
      <div
        v-if="isArticlePage && isMobile"
        class="float-button"
        title="目录"
        @click="showMobileToc = true"
      >
        <i class="ri-list-unordered" />
      </div>

      <!-- 跳转到评论区按钮（仅文章页显示） -->
      <div
        v-if="isArticlePage"
        class="float-button"
        title="跳转评论区"
        @click="scrollToElement('.comment-input')"
      >
        <i class="ri-message-3-line" />
      </div>

      <!-- 回到顶部按钮 -->
      <div
        class="float-button scroll-to-top"
        title="回到顶部"
        @click="scrollToTop"
        @mouseenter="isHovering = true"
        @mouseleave="isHovering = false"
      >
        <Transition name="content-fade" mode="out-in">
          <i v-if="isHovering" key="arrow" class="ri-arrow-up-line" />
          <span v-else key="progress" class="progress-text">{{ readingProgress }}</span>
        </Transition>
      </div>
    </div>
  </Transition>

  <!-- 移动端目录抽屉 -->
  <UiMobileToc :visible="showMobileToc" @close="showMobileToc = false" />

  <!-- 阅读模式退出按钮 -->
  <Transition name="fade">
    <div v-if="isReadingMode" class="reading-exit" title="退出阅读模式" @click="toggleReadingMode">
      <i class="ri-logout-box-r-line" />
    </div>
  </Transition>
</template>

<script setup lang="ts">
const route = useRoute();
const visible = ref(false);
const isHovering = ref(false);
const readingProgress = ref(0);
const isReadingMode = ref(false);
const showMobileToc = ref(false);

const isArticlePage = computed(() => route.name === 'posts-slug');

const isMobile = ref(false);

const checkMobile = () => {
  isMobile.value = window.innerWidth <= 900;
};

const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth',
  });
};

const toggleReadingMode = () => {
  isReadingMode.value = !isReadingMode.value;
  document.documentElement.setAttribute(
    'data-reading-mode',
    isReadingMode.value ? 'true' : 'false'
  );
};

const handleScroll = () => {
  const currentScrollY = window.scrollY;

  // 显示/隐藏按钮
  visible.value = currentScrollY > 100;

  // 计算阅读进度
  const windowHeight = window.innerHeight;
  const documentHeight = document.documentElement.scrollHeight;
  const scrollableHeight = documentHeight - windowHeight;

  if (scrollableHeight > 0) {
    const progress = Math.round((currentScrollY / scrollableHeight) * 100);
    readingProgress.value = Math.min(100, Math.max(0, progress));
  } else {
    readingProgress.value = 0;
  }
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);
  handleScroll();
  checkMobile();
  window.addEventListener('resize', checkMobile);
});

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll);
  window.removeEventListener('resize', checkMobile);
});
</script>

<style scoped lang="scss">
.float-button-group {
  position: fixed;
  right: 10px;
  bottom: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  z-index: 1000;
}

.float-button {
  width: 35px;
  height: 35px;
  background: var(--flec-btn);
  border-radius: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  color: white;
  font-size: 16px;

  &:hover {
    background: var(--flec-btn-hover);
  }

  &.scroll-to-top {
    position: relative;
    overflow: hidden;
  }

  .progress-text {
    line-height: 1;
    user-select: none;
  }

  // 主题切换按钮图标
  &.theme-toggle {
    position: relative;

    .theme-icon-moon,
    .theme-icon-sun {
      position: absolute;
      transition: opacity 0.2s ease;
    }
  }
}

.fade-enter-active,
.fade-leave-active {
  transition:
    opacity 0.3s ease,
    transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateX(60px);
}

.content-fade-enter-active,
.content-fade-leave-active {
  transition: opacity 0.15s ease;
}

.content-fade-enter-from,
.content-fade-leave-to {
  opacity: 0;
}

// 阅读模式退出按钮
.reading-exit {
  position: fixed;
  top: 20px;
  right: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 50px;
  background: var(--flec-btn);
  color: #ffffff;
  border-radius: 5px;
  cursor: pointer;
  z-index: 1000;
  transition: all 0.3s ease;

  &:hover {
    background: var(--flec-btn-hover);
  }

  i {
    font-size: 25px;
  }
}
</style>

<style lang="scss">
// 阅读模式
html[data-reading-mode='true'] {
  // 隐藏非必要元素
  nav,
  footer,
  aside,
  .ai-summary,
  .post-copyright,
  .comments-section,
  .float-button-group {
    display: none !important;
  }

  // 阅读模式配色变量（浅色主题）
  &[data-theme='light'] {
    --reading-bg: #faf9f6;
    --reading-text: #2c2c2c;
    --reading-title: #2c2c2c;
    --reading-meta: #666;
  }

  // 阅读模式配色变量（深色主题）
  &[data-theme='dark'] {
    --reading-bg: #1a1a1a;
    --reading-text: #e8e8e8;
    --reading-title: #f0f0f0;
    --reading-meta: #999;
  }

  body,
  .layout-wrapper,
  .page-main {
    background: var(--reading-bg) !important;
    color: var(--reading-text) !important;
  }

  // 简化头部
  .post-header {
    background: none !important;
    margin-top: 0 !important;
    padding-top: 0 !important;
    pointer-events: none !important;

    &::before {
      display: none !important;
    }

    .post-title {
      color: var(--reading-title) !important;
    }

    .post-meta {
      color: var(--reading-meta) !important;
    }
  }

  // 简化内容区
  .page-main {
    padding: 0 !important;
  }

  .main-layout {
    max-width: 800px !important;
    padding: 0 40px !important;
  }

  .main-content {
    width: 100% !important;
  }

  #post {
    background: none !important;
    box-shadow: none !important;
    border: none !important;
    padding: 0 0 60px !important;
  }

  // 响应式
  @media (max-width: 900px) {
    .main-layout {
      padding: 0 20px !important;
    }
  }

  @media (max-width: 600px) {
    .main-layout {
      padding: 0 15px !important;
    }
  }
}
</style>
