<script setup lang="ts">
import mediumZoom from 'medium-zoom';
import mermaid from 'mermaid';
import { loadEmojiMap } from '@/composables/useEmojis';
import { useSysConfig } from '@/composables/useStores';

const { blogConfig } = useSysConfig();

interface Props {
  content: string;
}

const props = defineProps<Props>();

const initMermaid = () => {
  mermaid.initialize({
    startOnLoad: false,
    theme: 'default',
    securityLevel: 'loose',
  });
};

const renderMermaidDiagrams = async () => {
  const elements = document.querySelectorAll('.mermaid:not(:has(svg))');

  for (const element of elements) {
    try {
      const { svg } = await mermaid.render(`mermaid-${Date.now()}`, element.textContent || '');
      element.innerHTML = svg;
    } catch (error) {
      console.error('Mermaid 渲染失败:', error);
    }
  }
};

// 表情数据 ref
const emojiMap = ref<Map<string, string> | null>(null);

// 渲染内容
const renderedContent = computed(() => {
  if (!props.content) return '';
  void emojiMap.value;
  return renderMarkdown(props.content);
});

let zoom: ReturnType<typeof mediumZoom> | null = null;

const initZoom = () => {
  const contentEl = document.querySelector('.markdown-content');
  if (!contentEl) return;

  const images = contentEl.querySelectorAll('img');
  if (images.length === 0) return;
  if (zoom) zoom.detach();
  zoom = mediumZoom(images, { margin: 24, background: 'rgba(0, 0, 0, 0.9)', scrollOffset: 48 });
};

watch(
  () => renderedContent.value,
  async () => {
    await nextTick();
    initZoom();
    await renderMermaidDiagrams();
  }
);

onMounted(() => {
  initMermaid();

  // 加载表情数据
  const emojisUrl = blogConfig.value.emojis;
  if (emojisUrl) loadEmojiMap(emojisUrl).then(map => (emojiMap.value = map));
  nextTick(async () => {
    initZoom();
    await renderMermaidDiagrams();
  });
});

onUnmounted(() => {
  if (zoom) {
    zoom.detach();
    zoom = null;
  }
});
</script>

<template>
  <article class="post-content">
    <!-- eslint-disable-next-line vue/no-v-html -- 内容经过 DOMPurify 净化处理 -->
    <div class="markdown-content" v-html="renderedContent" />
  </article>
</template>

<style lang="scss">
@use '@/assets/css/prose' as *;
</style>

<style>
@import 'highlight.js/styles/atom-one-dark.css';
@import 'katex/dist/katex.min.css';

.medium-zoom-overlay {
  z-index: 9999 !important;
}

.medium-zoom-image {
  z-index: 10000 !important;
}
</style>

<style lang="scss" scoped>
.post-content {
  line-height: 1.8;
  font-size: 1rem;
  color: var(--theme-text-color);
  word-wrap: break-word;

  :deep(.markdown-content) {
    img {
      cursor: zoom-in;
      transition: transform 0.2s ease;
      max-width: 100%;
      height: auto;

      &:hover {
        transform: scale(1.02);
      }
    }

    // 表情图片不可点击
    .emoji-image {
      cursor: default !important;
      pointer-events: none;

      &:hover {
        transform: none;
      }
    }

    .mermaid {
      display: flex;
      justify-content: center;
      align-items: center;
      margin: 1.5rem 0;
      padding: 1rem;
      background: var(--theme-bg-color-secondary, #f5f5f5);
      border-radius: 8px;
      overflow-x: auto;

      svg {
        max-width: 100%;
        height: auto;
      }
    }

    .mermaid-error {
      color: #f56c6c;
      padding: 1rem;
      background: #fef0f0;
      border-radius: 4px;
      border-left: 4px solid #f56c6c;
    }
  }
}
</style>
