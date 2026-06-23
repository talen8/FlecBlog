import { ref, nextTick, onMounted, onUnmounted } from 'vue';
import type { Ref } from 'vue';
import { useEventListener, useResizeObserver } from '@vueuse/core';

export interface WaterfallOptions {
  containerSelector: string;
  columns?: number;
  gap?: number;
  debounceDelay?: number;
  waitForImages?: boolean;
  breakpoints?: {
    mobile: number;
    tablet: number;
  };
}

export function useWaterfall(options: WaterfallOptions) {
  const {
    containerSelector,
    columns = 3,
    gap = 15,
    debounceDelay = 150,
    waitForImages = true,
    breakpoints = { mobile: 768, tablet: 1200 },
  } = options;

  const isLayoutReady: Ref<boolean> = ref(false);
  let debounceTimer: ReturnType<typeof setTimeout> | null = null;
  let imageLoadListeners: Array<{ img: HTMLImageElement; listener: () => void }> = [];

  // 根据窗口宽度获取列数
  const getColumns = (): number => {
    const width = window.innerWidth;
    if (width <= breakpoints.mobile) return 1;
    if (width <= breakpoints.tablet) return 2;
    return columns;
  };

  // 等待所有图片加载完成
  const waitForImagesLoad = async (): Promise<void> => {
    await nextTick();
    const images = document.querySelectorAll(`${containerSelector} img`);
    if (images.length === 0) {
      return Promise.resolve();
    }
    const imagePromises = Array.from(images).map(img => {
      return new Promise<void>(resolve => {
        if ((img as HTMLImageElement).complete) {
          resolve();
        } else {
          img.addEventListener('load', () => resolve());
          img.addEventListener('error', () => resolve());
        }
      });
    });
    return Promise.all(imagePromises).then(() => {});
  };

  // 执行瀑布流布局计算
  const waterfallLayout = () => {
    const container = document.querySelector(containerSelector);
    if (!container) return;

    const items = Array.from(container.children) as HTMLElement[];
    if (items.length === 0) return;

    const containerWidth = container.clientWidth;
    const cols = getColumns();
    const columnWidth = (containerWidth - gap * (cols - 1)) / cols;

    const columnsHeight = new Array(cols).fill(0);
    const itemHeights = items.map(item => item.offsetHeight);

    items.forEach((item, index) => {
      const minHeight = Math.min(...columnsHeight);
      const columnIndex = columnsHeight.indexOf(minHeight);

      item.style.width = `${columnWidth}px`;
      item.style.position = 'absolute';
      item.style.left = `${columnIndex * (columnWidth + gap)}px`;
      item.style.top = `${minHeight}px`;

      columnsHeight[columnIndex] = minHeight + (itemHeights[index] || 0) + gap;
    });

    const containerHeight = Math.max(...columnsHeight);
    (container as HTMLElement).style.height = `${containerHeight}px`;

    if (!isLayoutReady.value) {
      isLayoutReady.value = true;
    }
  };

  // 防抖布局函数
  const debouncedLayout = () => {
    if (debounceTimer) clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
      waterfallLayout();
    }, debounceDelay);
  };

  // 清理图片加载监听器
  const cleanupImageListeners = () => {
    imageLoadListeners.forEach(({ img, listener }) => {
      img.removeEventListener('load', listener);
    });
    imageLoadListeners = [];
  };

  // 设置图片加载监听，图片加载完成后自动重新布局
  const setupImageLoadListeners = () => {
    cleanupImageListeners();

    const container = document.querySelector(containerSelector);
    if (!container) return;

    const images = container.querySelectorAll('img');
    images.forEach(img => {
      const htmlImg = img as HTMLImageElement;
      if (!htmlImg.complete) {
        const listener = () => {
          debouncedLayout();
        };
        htmlImg.addEventListener('load', listener);
        imageLoadListeners.push({ img: htmlImg, listener });
      }
    });
  };

  // 主瀑布流函数
  const waterfall = async () => {
    isLayoutReady.value = false;
    cleanupImageListeners();

    if (waitForImages) {
      await waitForImagesLoad();
    } else {
      await nextTick();
    }

    waterfallLayout();
    setupImageLoadListeners();
  };

  onMounted(() => {
    const container = document.querySelector(containerSelector) as HTMLElement;
    if (container) {
      useResizeObserver(container, debouncedLayout);
    }
    useEventListener(window, 'load', waterfall);
    useEventListener(window, 'resize', debouncedLayout);
  });

  onUnmounted(() => {
    cleanupImageListeners();
    if (debounceTimer) {
      clearTimeout(debounceTimer);
    }
  });

  return {
    waterfall,
    waterfallLayout,
    debouncedLayout,
    isLayoutReady,
    waitForImagesLoad,
  };
}
