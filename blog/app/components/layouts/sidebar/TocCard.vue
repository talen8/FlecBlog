<script setup lang="ts">
const { currentArticle } = useCurrentArticle();
const activeId = ref<string>('');
const tocListRef = ref<HTMLElement | null>(null);

// 从当前文章中提取目录
const toc = computed<TocItem[]>(() => {
  if (!currentArticle.value?.content) return [];
  return extractToc(currentArticle.value.content);
});

// 判断是否有目录项
const hasToc = computed(() => toc.value.length > 0);

// 滚动目录列表，使激活项居中
const scrollTocToActive = (id: string) => {
  if (!tocListRef.value) return;

  nextTick(() => {
    const activeButton = tocListRef.value?.querySelector(`[data-toc-id="${id}"]`) as HTMLElement;
    if (!activeButton) return;

    const container = tocListRef.value!;
    const containerHeight = container.clientHeight;
    const buttonTop = activeButton.offsetTop;
    const buttonHeight = activeButton.clientHeight;

    // 计算让按钮居中的滚动位置
    const targetScroll = buttonTop - containerHeight / 2 + buttonHeight / 2;

    // 平滑滚动到目标位置
    container.scrollTo({
      top: targetScroll,
      behavior: 'smooth',
    });
  });
};

// 滚动到指定标题
const scrollToHeading = (id: string) => {
  scrollToElement(`#${id}`, { block: 'start' });
};

// 监听滚动，高亮当前阅读项
const handleScroll = () => {
  const referencePoint = 64; // 参考线位置（距视口顶部64px）
  const headings = toc.value;

  if (headings.length === 0) return;

  let closestHeading: TocItem | undefined = undefined;
  let closestDistance = Infinity;

  // 找到距离参考线最近的标题
  for (const heading of headings) {
    const element = document.getElementById(heading.id);
    if (!element) continue;

    const rect = element.getBoundingClientRect();
    const distanceToReference = Math.abs(rect.top - referencePoint);

    // 只考虑在参考线上方或稍微下方的标题（不超过50px�?    // 且距离参考线最近的
    if (rect.top <= referencePoint + 50 && distanceToReference < closestDistance) {
      closestDistance = distanceToReference;
      closestHeading = heading;
    }
  }

  // 如果找到了，激活它；否则激活第一个标题
  const targetHeading = closestHeading || headings[0];
  if (targetHeading && targetHeading.id !== activeId.value) {
    activeId.value = targetHeading.id;
    scrollTocToActive(targetHeading.id);
  }
};

onMounted(() => {
  // 使用 VueUse 自动管理事件监听（自动清理）
  useEventListener(window, 'scroll', handleScroll, { passive: true });
  handleScroll(); // 初始化当前阅读项
});
</script>

<template>
  <div class="card-widget" v-if="hasToc">
    <div class="item-headline">
      <i class="ri-list-unordered"></i>
      <span>目录</span>
    </div>

    <nav ref="tocListRef" class="toc-list" aria-label="文章目录" data-lenis-prevent>
      <button
        v-for="item in toc"
        :key="item.id"
        :data-toc-id="item.id"
        :class="['toc-item', `toc-level-${item.level}`, { active: activeId === item.id }]"
        @click="scrollToHeading(item.id)"
        :aria-label="`跳转到 ${item.text}`"
        :aria-current="activeId === item.id ? 'location' : undefined"
      >
        <span class="toc-text">{{ item.text }}</span>
      </button>
    </nav>
  </div>
</template>

<style lang="scss" scoped>
.toc-list {
  margin: 10px 0 0;
  padding: 0;
  max-height: calc(100vh - 176px);
  overflow-y: auto;
  scroll-behavior: smooth;

  // 自定义滚动条样式
  &::-webkit-scrollbar {
    width: 3px;
  }

  &::-webkit-scrollbar-thumb {
    background: color-mix(in srgb, var(--flec-btn-hover) 50%, transparent);
    border-radius: 3px;

    &:hover {
      background: color-mix(in srgb, var(--flec-btn-hover) 70%, transparent);
    }
  }

  &::-webkit-scrollbar-track {
    background: transparent;
  }
}

.toc-item {
  width: 100%;
  text-align: left;
  background: transparent;
  border: none;
  padding: 6px 8px;
  margin: 2px 0;
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 6px;
  border-left: 2px solid transparent;
  line-height: 1.5;
  color: inherit;
  font-family: inherit;
  font-size: inherit;

  &:hover {
    background-color: rgba(73, 177, 245, 0.1);
    border-left-color: var(--flec-btn-hover);
  }

  &.active {
    background-color: var(--flec-btn-hover);
    color: #fff;
    border-left-color: var(--flec-btn-hover);

    .toc-text {
      font-weight: 500;
    }
  }

  .toc-text {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    display: block;
    font-size: 0.9rem;
  }
}

// 不同级别的标题缩进
.toc-level-1 {
  padding-left: 8px;
  font-weight: 500;

  &:hover {
    padding-left: 4px; // 向左偏移4px
  }
}

.toc-level-2 {
  padding-left: 16px;
  font-size: 0.95em;

  &:hover {
    padding-left: 12px; // 向左偏移4px
  }
}

.toc-level-3 {
  padding-left: 24px;
  font-size: 0.9em;
  opacity: 0.9;

  &:hover {
    padding-left: 20px; // 向左偏移4px
  }
}

.toc-level-4 {
  padding-left: 32px;
  font-size: 0.85em;
  opacity: 0.85;

  &:hover {
    padding-left: 28px; // 向左偏移4px
  }
}

.toc-level-5 {
  padding-left: 40px;
  font-size: 0.8em;
  opacity: 0.8;

  &:hover {
    padding-left: 36px; // 向左偏移4px
  }
}

.toc-level-6 {
  padding-left: 48px;
  font-size: 0.75em;
  opacity: 0.75;

  &:hover {
    padding-left: 44px; // 向左偏移4px
  }
}

@media screen and (max-width: 900px) {
  .card-widget {
    display: none;
  }
}
</style>
