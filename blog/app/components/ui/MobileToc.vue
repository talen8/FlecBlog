<script setup lang="ts">
import type { TocItem } from '@/utils/markdown';

interface Props {
  visible: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits<{
  close: [];
}>();

const { currentArticle } = useCurrentArticle();
const activeId = ref<string>('');
const tocListRef = ref<HTMLElement | null>(null);
const tocPopoverRef = ref<HTMLElement | null>(null);

const toc = computed<TocItem[]>(() => {
  if (!currentArticle.value?.content) return [];
  return extractToc(currentArticle.value.content);
});

const hasToc = computed(() => toc.value.length > 0);

onClickOutside(tocPopoverRef, () => {
  if (props.visible) {
    emit('close');
  }
});

const scrollTocToActive = (id: string) => {
  if (!tocListRef.value) return;

  nextTick(() => {
    const activeButton = tocListRef.value?.querySelector(`[data-toc-id="${id}"]`) as HTMLElement;
    if (!activeButton) return;

    const container = tocListRef.value!;
    const containerHeight = container.clientHeight;
    const buttonTop = activeButton.offsetTop;
    const buttonHeight = activeButton.clientHeight;

    const targetScroll = buttonTop - containerHeight / 2 + buttonHeight / 2;

    container.scrollTo({
      top: targetScroll,
      behavior: 'smooth',
    });
  });
};

const scrollToHeading = (id: string) => {
  scrollToElement(`#${id}`, { block: 'start' });
  emit('close');
};

const handleScroll = () => {
  const referencePoint = 64;
  const headings = toc.value;

  if (headings.length === 0) return;

  let closestHeading: TocItem | undefined = undefined;
  let closestDistance = Infinity;

  for (const heading of headings) {
    const element = document.getElementById(heading.id);
    if (!element) continue;

    const rect = element.getBoundingClientRect();
    const distanceToReference = Math.abs(rect.top - referencePoint);

    if (rect.top <= referencePoint + 50 && distanceToReference < closestDistance) {
      closestDistance = distanceToReference;
      closestHeading = heading;
    }
  }

  const targetHeading = closestHeading || headings[0];
  if (targetHeading && targetHeading.id !== activeId.value) {
    activeId.value = targetHeading.id;
    scrollTocToActive(targetHeading.id);
  }
};

watch(
  () => props.visible,
  newVal => {
    if (newVal) {
      handleScroll();
    }
  }
);

onMounted(() => {
  useEventListener(window, 'scroll', handleScroll, { passive: true });
  handleScroll();
});
</script>

<template>
  <Teleport to="body">
    <Transition name="toc-popover">
      <div v-if="visible && hasToc" ref="tocPopoverRef" class="mobile-toc-popover">
        <div class="toc-header">
          <div class="header-left">
            <i class="ri-menu-line" />
            <span>目录</span>
          </div>
          <span class="toc-count">{{ toc.length }}</span>
        </div>

        <nav ref="tocListRef" class="toc-list" aria-label="文章目录">
          <button
            v-for="item in toc"
            :key="item.id"
            :data-toc-id="item.id"
            :class="['toc-item', `toc-level-${item.level}`, { active: activeId === item.id }]"
            :aria-label="`跳转到 ${item.text}`"
            :aria-current="activeId === item.id ? 'location' : undefined"
            @click="scrollToHeading(item.id)"
          >
            <span class="toc-text">{{ item.text }}</span>
          </button>
        </nav>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
.mobile-toc-popover {
  position: fixed;
  right: 12px;
  bottom: 90px;
  width: 280px;
  max-width: calc(100vw - 24px);
  max-height: 60vh;
  background: var(--theme-bg-color, #fff);
  border-radius: 12px;
  box-shadow:
    0 6px 30px rgba(0, 0, 0, 0.12),
    0 2px 8px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  z-index: 2000;
}

.toc-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 18px;
  border-bottom: 1px solid var(--theme-border-color, rgba(0, 0, 0, 0.06));
  flex-shrink: 0;

  .header-left {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 15px;
    font-weight: 600;
    color: var(--theme-text-color, #333);

    i {
      font-size: 16px;
      color: var(--flec-btn-hover, #49b1f5);
    }
  }

  .toc-count {
    font-size: 13px;
    font-weight: 500;
    color: #999;
    min-width: 20px;
    text-align: right;
  }
}

.toc-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px 0;
  scroll-behavior: smooth;

  &::-webkit-scrollbar {
    width: 3px;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.15);
    border-radius: 3px;
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
  padding: 9px 18px;
  margin: 1px 0;
  cursor: pointer;
  transition: all 0.25s ease;
  line-height: 1.5;
  color: var(--theme-text-color, #555);
  font-family: inherit;
  font-size: inherit;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 3px;
    height: 0;
    background: var(--flec-btn-hover, #49b1f5);
    border-radius: 0 2px 2px 0;
    transition: height 0.25s ease;
  }

  &:hover {
    background-color: rgba(73, 177, 245, 0.06);
    color: var(--theme-text-color, #333);

    &::before {
      height: 16px;
    }
  }

  &.active {
    background-color: var(--flec-btn-hover, #49b1f5);
    color: #fff;
    font-weight: 500;

    &::before {
      display: none;
    }
  }

  .toc-text {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    display: block;
    font-size: 0.88rem;
  }
}

.toc-level-1 {
  padding-left: 18px;
  font-weight: 500;
}

.toc-level-2 {
  padding-left: 28px;
  font-size: 0.95em;
}

.toc-level-3 {
  padding-left: 38px;
  font-size: 0.9em;
  opacity: 0.92;
}

.toc-level-4 {
  padding-left: 48px;
  font-size: 0.85em;
  opacity: 0.85;
}

.toc-level-5 {
  padding-left: 58px;
  font-size: 0.8em;
  opacity: 0.8;
}

.toc-level-6 {
  padding-left: 68px;
  font-size: 0.75em;
  opacity: 0.75;
}

.toc-popover-enter-active,
.toc-popover-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: bottom right;
}

.toc-popover-enter-from,
.toc-popover-leave-to {
  opacity: 0;
  transform: scale(0.92) translate(10px, 10px);
}
</style>
