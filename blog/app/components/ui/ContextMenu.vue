<script setup lang="ts">
type TargetType = 'link' | 'media' | 'text' | 'input' | 'none';

interface MenuItem {
  id: string;
  label: string;
  icon: string;
  action?: () => void;
  route?: string;
}

const router = useRouter();
const { info, success, error } = useToast();

// ========== 状态管理 ==========
const isVisible = ref(false);
const menuRef = ref<HTMLElement | null>(null);
const pos = ref({ x: 0, y: 0 });
const target = ref<{
  type: TargetType;
  url?: string;
  text?: string;
  element?: HTMLElement;
  mediaType?: 'image' | 'video';
  hasComment?: boolean;
}>({ type: 'none' });

// ========== 工具函数 ==========
const copy = async (text: string, msg: string) => {
  try {
    await navigator.clipboard.writeText(text);
    success(msg);
  } catch {
    error('复制失败');
  }
};

// 向上查找指定标签的元素
const getElement = <T extends HTMLElement>(el: HTMLElement, tag: string) =>
  (el.tagName === tag ? el : el.closest?.(tag.toLowerCase())) as T | null;

// 检测目标元素
const detect = (el: EventTarget | null) => {
  if (!el || !(el instanceof HTMLElement)) return { type: 'none' as TargetType };

  // 输入框
  if (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA' || el.isContentEditable)
    return { type: 'input' as TargetType, element: el };

  // 选中文本
  const text = window.getSelection()?.toString().trim();
  if (text)
    return {
      type: 'text' as TargetType,
      text,
      hasComment: !!document.querySelector('.comment-input'),
    };

  // 视频
  const video = getElement<HTMLVideoElement>(el, 'VIDEO');
  if (video)
    return {
      type: 'media' as TargetType,
      mediaType: 'video' as const,
      url: video.src || video.currentSrc,
      element: video,
    };

  // 图片
  const img = getElement<HTMLImageElement>(el, 'IMG');
  if (img)
    return { type: 'media' as TargetType, mediaType: 'image' as const, url: img.src, element: img };

  // 链接
  const link = getElement<HTMLAnchorElement>(el, 'A');
  if (link?.href) return { type: 'link' as TargetType, url: link.href };

  return { type: 'none' as TargetType };
};

// ========== 菜单控制 ==========
// 首次提示的本地存储键
const TOOLTIP_KEY = 'ctx-tip';
const TOOLTIP_DURATION = 3 * 60 * 1000; // 3分钟

const shouldUseNativeContextMenu = () => {
  if (typeof window === 'undefined') return true;

  const prefersTouchInteraction = window.matchMedia('(hover: none) and (pointer: coarse)').matches;
  const isMobileViewport = window.innerWidth <= 768;

  return prefersTouchInteraction || isMobileViewport;
};

// 显示首次提示
const showFirstTimeTooltip = () => {
  const lastShown = localStorage.getItem(TOOLTIP_KEY);
  const now = Date.now();

  // 如果从未显示过，或者距离上次显示已超过3分钟
  if (!lastShown || now - parseInt(lastShown) > TOOLTIP_DURATION) {
    info('按住 Ctrl 再点击右键，即可恢复原界面哦');
    localStorage.setItem(TOOLTIP_KEY, now.toString());
  }
};

const showMenu = (e: MouseEvent) => {
  // 如果按住Ctrl键，显示原生右键菜单
  if (e.ctrlKey) {
    return; // 不阻止默认行为，显示原生菜单
  }

  e.preventDefault();

  // 首次显示时给出提示
  showFirstTimeTooltip();

  pos.value = { x: e.clientX, y: e.clientY };
  target.value = detect(e.target);
  isVisible.value = true;
  setTimeout(adjust, 0);
};

const close = () => {
  isVisible.value = false;
  target.value = { type: 'none' };
};

// ========== 静态菜单配置 ==========
const navBtns = [
  {
    id: 'back',
    icon: 'ri-arrow-left-line',
    tooltip: '后退',
    action: () => {
      router.back();
      close();
    },
  },
  {
    id: 'forward',
    icon: 'ri-arrow-right-line',
    tooltip: '前进',
    action: () => {
      router.forward();
      close();
    },
  },
  { id: 'refresh', icon: 'ri-refresh-line', tooltip: '刷新', action: () => location.reload() },
  {
    id: 'top',
    icon: 'ri-arrow-up-line',
    tooltip: '回到顶部',
    action: () => {
      scrollTo({ top: 0, behavior: 'smooth' });
      close();
    },
  },
];

const jumpItems = [
  { id: 'archive', label: '归档', icon: 'ri-archive-line', route: '/archive' },
  { id: 'tag', label: '标签', icon: 'ri-price-tag-line', route: '/tags' },
  { id: 'category', label: '分类', icon: 'ri-book-shelf-line', route: '/categories' },
];

const toolItems = [
  { id: 'theme', label: '昼夜切换', icon: 'ri-contrast-2-line', action: toggleTheme },
  {
    id: 'share',
    label: '分享页面',
    icon: 'ri-share-line',
    action: () => copy(location.href, '页面地址已复制到剪贴板'),
  },
  { id: 'copyright', label: '版权声明', icon: 'ri-copyright-line', route: '/copyright' },
];

// ========== 动态菜单 ==========
// 链接菜单
const linkItems = computed(() => {
  if (target.value.type !== 'link' || !target.value.url) return [];
  const url = target.value.url;
  return [
    {
      id: 'open',
      label: '新标签页',
      icon: 'ri-external-link-line',
      action: () => open(url, '_blank'),
    },
    {
      id: 'copy',
      label: '复制链接',
      icon: 'ri-file-copy-line',
      action: () => copy(url, '链接已复制到剪贴板'),
    },
  ];
});

// 媒体菜单（图片/视频）
const mediaItems = computed(() => {
  if (target.value.type !== 'media' || !target.value.url) return [];

  const { url, mediaType, element } = target.value;
  const isVideo = mediaType === 'video';
  const name = isVideo ? '视频' : '图片';
  const items = [];

  // 视频播放/暂停
  if (isVideo && element instanceof HTMLVideoElement) {
    const paused = element.paused;
    items.push({
      id: 'play',
      label: paused ? '播放' : '暂停',
      icon: paused ? 'ri-play-line' : 'ri-pause-line',
      action: () => {
        if (paused) {
          element.play();
        } else {
          element.pause();
        }
      },
    });
  }

  // 保存媒体文件
  const download = async () => {
    try {
      const blob = await fetch(url!).then(r => r.blob());
      const a = Object.assign(document.createElement('a'), {
        href: URL.createObjectURL(blob),
        download: url!.split('/').pop() || (isVideo ? 'video.mp4' : 'image.jpg'),
      });
      a.click();
      URL.revokeObjectURL(a.href);
    } catch {
      error('保存失败');
    }
  };

  items.push(
    {
      id: 'open',
      label: '新标签页',
      icon: isVideo ? 'ri-video-line' : 'ri-image-line',
      action: () => open(url!, '_blank'),
    },
    {
      id: 'copy',
      label: '复制地址',
      icon: 'ri-file-copy-line',
      action: () => copy(url!, `${name}地址已复制到剪贴板`),
    },
    { id: 'download', label: `保存${name}`, icon: 'ri-download-line', action: download }
  );

  return items;
});

// 文本菜单
const textItems = computed(() => {
  if (target.value.type !== 'text' || !target.value.text) return [];

  const { text, hasComment } = target.value;
  const items = [
    {
      id: 'copy',
      label: '复制文本',
      icon: 'ri-file-copy-line',
      action: () => copy(text!, '文本已复制到剪贴板'),
    },
  ];

  // 如果页面有评论区，添加引用功能
  if (hasComment) {
    items.push({
      id: 'quote',
      label: '引用评论',
      icon: 'ri-chat-quote-line',
      action: () => {
        fillComment(`> ${text!.length > 100 ? text!.substring(0, 100) + '...' : text}\n\n`);
        return Promise.resolve();
      },
    });
  }

  items.push({
    id: 'search',
    label: 'Bing 搜索',
    icon: 'ri-search-line',
    action: () => {
      window.open(`https://www.bing.com/search?q=${encodeURIComponent(text!)}`, '_blank');
      return Promise.resolve();
    },
  });

  return items;
});

// 输入框菜单
const inputItems = computed(() => {
  if (target.value.type !== 'input' || !target.value.element) return [];

  const el = target.value.element as HTMLInputElement | HTMLTextAreaElement;
  const hasSel = () => (el.selectionStart || 0) !== (el.selectionEnd || 0);
  const emit = () => el.dispatchEvent(new Event('input', { bubbles: true }));

  // 剪切选中文本
  const cut = async () => {
    if (!hasSel()) return error('请先选中文本');
    try {
      const s = el.selectionStart || 0,
        e = el.selectionEnd || 0,
        v = el.value || '';
      await navigator.clipboard.writeText(v.substring(s, e));
      el.value = v.substring(0, s) + v.substring(e);
      el.setSelectionRange(s, s);
      emit();
    } catch {
      error('剪切失败');
    }
  };

  const copyText = async () => {
    if (!hasSel()) return error('请先选中文本');
    try {
      await navigator.clipboard.writeText(
        (el.value || '').substring(el.selectionStart || 0, el.selectionEnd || 0)
      );
      success('已复制文本');
    } catch {
      error('复制失败');
    }
  };

  // 粘贴剪贴板内容
  const paste = async () => {
    try {
      const text = await navigator.clipboard.readText();
      el.focus();
      const s = el.selectionStart || 0,
        e = el.selectionEnd || 0,
        v = el.value || '';
      el.value = v.substring(0, s) + text + v.substring(e);
      el.setSelectionRange(s + text.length, s + text.length);
      emit();
    } catch {
      error('粘贴失败，请检查剪贴板权限');
    }
  };

  return [
    { id: 'cut', label: '剪切', icon: 'ri-scissors-line', action: cut },
    { id: 'copy', label: '复制', icon: 'ri-file-copy-line', action: copyText },
    { id: 'paste', label: '粘贴', icon: 'ri-clipboard-line', action: paste },
    {
      id: 'all',
      label: '全选',
      icon: 'ri-text',
      action: () => {
        el.focus();
        el.select();
      },
    },
  ];
});

// ========== 位置和交互 ==========
const style = computed(() => ({ left: `${pos.value.x}px`, top: `${pos.value.y}px` }));

// 调整菜单位置，确保不超出视口边界
const adjust = () => {
  if (!menuRef.value) return;
  const { width: w, height: h } = menuRef.value.getBoundingClientRect();
  const { innerWidth: vw, innerHeight: vh } = window;
  let { x, y } = pos.value;
  if (x + w > vw) x = vw - w - 10;
  if (y + h > vh) y = vh - h - 10;
  if (x < 0) x = 10;
  if (y < 0) y = 10;
  pos.value = { x, y };
};

const click = (item: MenuItem) => {
  if (item.action) {
    item.action();
  } else if (item.route) {
    router.push(item.route);
  }
  close();
};

// 事件监听器
const onClickOut = (e: MouseEvent) =>
  menuRef.value && !menuRef.value.contains(e.target as Node) && e.button === 0 && close();
const onEsc = (e: KeyboardEvent) => e.key === 'Escape' && close();
const addListeners = () => {
  addEventListener('click', onClickOut);
  addEventListener('keydown', onEsc);
};
const rmListeners = () => {
  removeEventListener('click', onClickOut);
  removeEventListener('keydown', onEsc);
};

onMounted(() => {
  if (shouldUseNativeContextMenu()) return;
  addEventListener('contextmenu', showMenu);
});
onUnmounted(() => {
  removeEventListener('contextmenu', showMenu);
  rmListeners();
});
</script>

<template>
  <Teleport to="body">
    <Transition name="ctx" @after-enter="addListeners" @after-leave="rmListeners">
      <div v-if="isVisible" ref="menuRef" class="ctx-menu" :style="style" @click.stop>
        <!-- 导航按钮组 -->
        <div class="group small">
          <button
            v-for="b in navBtns"
            :key="b.id"
            class="item icon"
            :title="b.tooltip"
            @click="b.action"
          >
            <i :class="b.icon" />
          </button>
        </div>
        <!-- 动态菜单组（链接、媒体、文本、输入框） -->
        <div
          v-for="items in [linkItems, mediaItems, textItems, inputItems]"
          v-show="items.length"
          :key="items[0]?.id"
          class="group line"
        >
          <div v-for="i in items" :key="i.id" class="item" @click="click(i)">
            <i :class="i.icon" /><span>{{ i.label }}</span>
          </div>
        </div>
        <!-- 快速跳转 -->
        <div class="group line">
          <div v-for="i in jumpItems" :key="i.id" class="item" @click="click(i)">
            <i :class="i.icon" /><span>{{ i.label }}</span>
          </div>
        </div>
        <!-- 工具菜单 -->
        <div class="group line">
          <div v-for="i in toolItems" :key="i.id" class="item" @click="click(i)">
            <i :class="i.icon" /><span>{{ i.label }}</span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

/* 菜单容器 */
.ctx-menu {
  @extend .cardHover;
  position: fixed;
  z-index: 10000;
  width: 160px;
  user-select: none;
}

/* 菜单分组 */
.group {
  display: flex;

  /* 垂直列表样式（带分隔线） */
  &.line {
    padding: 7px 6px;
    flex-direction: column;
    border-top: 1px solid var(--flec-border);
  }

  /* 横向按钮组样式 */
  &.small {
    padding: 7px 6px;
    justify-content: space-between;
  }
}

/* 菜单项 */
.item {
  display: flex;
  align-items: center;
  padding: 4px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--font-color);

  i {
    margin: 0 5px;
  }

  span {
    flex: 1;
    white-space: nowrap;
  }

  &:hover {
    background: var(--flec-btn);
    color: #fff;
  }

  &:active {
    transform: scale(0.98);
  }

  /* 图标按钮样式（导航按钮） */
  &.icon {
    width: 36px;
    height: 36px;
    padding: 0;
    justify-content: center;

    i {
      width: auto;
      font-weight: 600;
    }
  }

  &:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }
}

/* 菜单出现/消失动画 */
.ctx-enter-active,
.ctx-leave-active {
  transition:
    opacity 0.15s,
    transform 0.15s;
}

.ctx-enter-from,
.ctx-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

/* 移动端隐藏右键菜单 */
@media (max-width: 768px) {
  .ctx-menu {
    display: none;
  }
}
</style>
