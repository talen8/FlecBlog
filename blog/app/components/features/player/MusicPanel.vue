<script setup lang="ts">
/**
 * 全局音乐播放器面板
 * 全部状态来自全局播放引擎，不持有私有播放状态
 */
const {
  state,
  currentTrack,
  progress,
  playAt,
  toggle,
  next,
  prev,
  seekRatio,
  setVolume,
  cycleMode,
} = useAudioEngine();

const { panel, setHovered, expand } = useMusicPanel();

const isMinimized = computed(() => state.value.playing && !panel.value.hovered);

/** 最小化态歌词文本元素 */
const miniLyricText = ref<HTMLElement | null>(null);
/** 最小化态歌词的溢出宽度（px），大于 0 时才需要滚动 */
const marqueeShift = ref(0);
/** 歌单加载中或音轨缓冲中 */
const isLoading = computed(() => state.value.loading || state.value.buffering);

/** 播放顺序对应的图标与文案 */
const modeIcons: Record<string, { icon: string; label: string }> = {
  list: { icon: 'ri-list-ordered', label: '顺序播放' },
  random: { icon: 'ri-shuffle-line', label: '随机播放' },
  single: { icon: 'ri-repeat-one-line', label: '单曲循环' },
};

const currentMode = computed(() => modeIcons[state.value.mode] ?? modeIcons.list);

const currentLyric = computed(() => state.value.lyrics[state.value.lyricIndex]?.text ?? '');
const nextLyric = computed(() => state.value.lyrics[state.value.lyricIndex + 1]?.text ?? '');

/** 最小化态文本：优先当前歌词，无歌词时回退为曲名与歌手 */
const miniText = computed(() => {
  const track = currentTrack.value;
  if (!track) return currentLyric.value;
  return currentLyric.value || `${track.name} - ${track.artist}`;
});

/** 末句或缺失时长信息时的歌词显示窗口兜底（秒） */
const LYRIC_WINDOW_FALLBACK = 6;
/** 滚动开始前的停留时长（秒），便于读出句首 */
const START_HOLD = 0.8;
/** 滚动结束后的停留时长（秒），便于读出句尾 */
const END_HOLD = 0.6;
/** 最慢滚动速度（px/s），歌词停留时间充裕时不超过该速度，避免拖沓 */
const MIN_SCROLL_SPEED = 16;
/** 滚动时长下限（秒），快歌窗口不足时仍给一个可感知的滚动过程 */
const MIN_SCROLL_DURATION = 0.5;

/**
 * 当前歌词行的可用显示时长（秒）
 * 取下一行与本行的时间差；末句退化为歌曲剩余时长
 */
const lyricWindow = computed(() => {
  const current = state.value.lyrics[state.value.lyricIndex];
  if (!current) return LYRIC_WINDOW_FALLBACK;
  const nextLine = state.value.lyrics[state.value.lyricIndex + 1];
  if (nextLine) return Math.max(1, nextLine.time - current.time);
  return state.value.duration
    ? Math.max(1, state.value.duration - current.time)
    : LYRIC_WINDOW_FALLBACK;
});

/**
 * 歌词超长时的横向滚动参数
 * 位移量由实测溢出宽度决定；本句可用时长扣除首尾停留后即为滚动时长，
 * 句子越短滚动越快，但始终保证句首、句尾各有一段可读的停留时间
 */
const marqueeStyle = computed(() => {
  const available = lyricWindow.value - START_HOLD - END_HOLD;
  const duration = Math.max(
    MIN_SCROLL_DURATION,
    Math.min(marqueeShift.value / MIN_SCROLL_SPEED, available)
  );
  return {
    '--marquee-shift': `-${marqueeShift.value}px`,
    'animation-duration': `${duration}s`,
    'animation-delay': `${START_HOLD}s`,
  };
});

/**
 * 测量最小化态歌词的溢出宽度
 * 面板不可见时容器宽度为 0，测得结果会在下次展开 / 歌词变化时重测覆盖
 */
const measureMiniLyric = () => {
  const el = miniLyricText.value;
  const box = el?.parentElement;
  if (!el || !box) return;
  marqueeShift.value = Math.max(0, el.offsetWidth - box.clientWidth);
};

// 歌词变化与形态切换后重新测量溢出宽度，挂载后先测一次决定是否需要滚动
watch(
  [miniText, isMinimized],
  async () => {
    await nextTick();
    measureMiniLyric();
  },
  { immediate: true }
);

/**
 * 点击进度条跳转
 * @param event - 鼠标事件
 */
const onSeek = (event: MouseEvent) => {
  const rect = (event.currentTarget as HTMLElement).getBoundingClientRect();
  if (!rect.width) return;
  seekRatio((event.clientX - rect.left) / rect.width);
};

/**
 * 调整音量
 * @param event - 输入事件
 */
const onVolume = (event: Event) => {
  setVolume(Number((event.target as HTMLInputElement).value));
};

/** 点击最小化态面板任意处展开 */
const onPanelClick = () => {
  if (isMinimized.value) expand();
};

/**
 * 展开动画起始：高度归零，强制回流后再过渡到实测内容高度
 * @param el - 过渡元素
 */
const onBodyEnter = (el: Element) => {
  const body = el as HTMLElement;
  body.style.height = '0px';
  void body.offsetHeight;
  body.style.height = `${body.scrollHeight}px`;
};

/**
 * 展开动画结束：交还高度控制权，内容变化时可自由伸缩
 * @param el - 过渡元素
 */
const onBodyAfterEnter = (el: Element) => {
  (el as HTMLElement).style.height = '';
};

/**
 * 收起动画：先固定当前高度，再过渡到 0
 * @param el - 过渡元素
 */
const onBodyLeave = (el: Element) => {
  const body = el as HTMLElement;
  body.style.height = `${body.scrollHeight}px`;
  void body.offsetHeight;
  body.style.height = '0px';
};
</script>

<template>
  <div
    class="music-panel"
    :class="{ 'is-minimized': isMinimized }"
    @mouseenter="expand"
    @mouseleave="setHovered(false)"
    @click="onPanelClick"
  >
    <template v-if="currentTrack">
      <!-- 封面、曲目信息与窗口控制 -->
      <div class="panel-head">
        <div class="panel-cover" :class="{ spinning: state.playing }">
          <img v-if="currentTrack.cover" :src="currentTrack.cover" alt="cover" />
          <i v-else class="ri-music-2-fill" />
        </div>
        <div class="panel-meta">
          <!-- 最小化：有歌词时滚动显示歌词，否则显示曲名 -->
          <div v-if="isMinimized" class="mini-lyric" :class="{ marquee: marqueeShift > 0 }">
            <span
              ref="miniLyricText"
              :key="miniText"
              class="mini-lyric-text"
              :style="marqueeStyle"
              >{{ miniText }}</span
            >
          </div>
          <template v-else>
            <div class="panel-title" :title="currentTrack.name">{{ currentTrack.name }}</div>
            <div class="panel-artist" :title="currentTrack.artist">{{ currentTrack.artist }}</div>
          </template>
        </div>
      </div>

      <!-- 展开部分：形态切换时高度与透明度过渡，高度由 JS 钩子按内容实测 -->
      <Transition
        name="panel-body"
        @enter="onBodyEnter"
        @after-enter="onBodyAfterEnter"
        @leave="onBodyLeave"
      >
        <div v-if="!isMinimized" class="panel-body">
          <!-- 歌词 -->
          <div v-if="currentLyric" class="panel-lyric">
            <div class="lyric-current">{{ currentLyric }}</div>
            <div v-if="nextLyric" class="lyric-next">{{ nextLyric }}</div>
          </div>

          <!-- 进度条 -->
          <div class="panel-progress">
            <span class="time">{{ formatDuration(state.currentTime) }}</span>
            <div class="progress-track" @click="onSeek">
              <div class="progress-played" :style="{ width: `${progress}%` }" />
            </div>
            <span class="time">{{ formatDuration(state.duration) }}</span>
          </div>

          <!-- 控制栏 -->
          <div class="panel-controls">
            <button class="ctrl" title="上一首" @click="prev">
              <i class="ri-skip-back-fill" />
            </button>
            <button
              class="ctrl ctrl-play"
              :title="isLoading ? '加载中' : state.playing ? '暂停' : '播放'"
              @click="toggle"
            >
              <i v-if="isLoading" class="ri-loader-4-line spinning" />
              <i v-else :class="state.playing ? 'ri-pause-fill' : 'ri-play-fill'" />
            </button>
            <button class="ctrl" title="下一首" @click="next()">
              <i class="ri-skip-forward-fill" />
            </button>
            <button class="ctrl" :title="currentMode?.label" @click="cycleMode">
              <i :class="currentMode?.icon" />
            </button>
            <div class="volume">
              <i class="ri-volume-up-line" />
              <input
                type="range"
                min="0"
                max="1"
                step="0.01"
                :value="state.volume"
                title="音量"
                @input="onVolume"
              />
            </div>
          </div>

          <!-- 播放列表 -->
          <div v-if="state.tracks.length > 1" class="panel-list">
            <div
              v-for="(track, index) in state.tracks"
              :key="`${track.url}-${index}`"
              class="list-item"
              :class="{ active: index === state.index }"
              @click="playAt(index, true)"
            >
              <span class="item-index">
                <i v-if="index === state.index && state.playing" class="ri-play-fill" />
                <template v-else>{{ index + 1 }}</template>
              </span>
              <span class="item-name">{{ track.name }}</span>
              <span class="item-artist">{{ track.artist }}</span>
            </div>
          </div>
        </div>
      </Transition>
    </template>
  </div>
</template>

<style lang="scss" scoped>
.music-panel {
  position: fixed;
  right: 57px;
  bottom: 20px;
  z-index: 1001;
  width: min(280px, calc(100vw - 80px));
  padding: 12px;
  background: var(--flec-card-bg);
  border: 1px solid var(--flec-border);
  border-radius: 8px;
  box-shadow: 0 4px 18px rgba(0, 0, 0, 0.12);
  backdrop-filter: blur(8px);
  color: var(--flec-moment-font);
  font-size: 0.8rem;
  transition:
    width 0.25s ease,
    padding 0.25s ease;

  &.is-minimized {
    width: min(200px, calc(100vw - 80px));
    padding: 8px 10px;
    cursor: pointer;

    .panel-cover {
      width: 32px;
      height: 32px;
      font-size: 1rem;
    }

    .mini-lyric {
      position: relative;
      overflow: hidden;
      font-size: 0.78rem;
      color: var(--flec-moment-font);
      white-space: nowrap;

      .mini-lyric-text {
        display: inline-block;
        width: max-content;
        white-space: nowrap;
      }

      &.marquee .mini-lyric-text {
        animation-name: mini-lyric-marquee;
        animation-timing-function: linear;
        animation-fill-mode: both;
      }
    }
  }
}

.panel-body {
  overflow: hidden;
}

.panel-body-enter-active,
.panel-body-leave-active {
  transition:
    height 0.25s ease,
    opacity 0.25s ease,
    transform 0.25s ease;
}

.panel-body-enter-from,
.panel-body-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

.panel-head {
  display: flex;
  align-items: center;
  gap: 10px;

  .panel-cover {
    flex-shrink: 0;
    width: 48px;
    height: 48px;
    transition:
      width 0.25s ease,
      height 0.25s ease;
    border-radius: 50%;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: #fff;
    font-size: 1.3rem;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    &.spinning {
      animation: rotate 14s linear infinite;
    }
  }

  .panel-meta {
    min-width: 0;
    flex: 1;

    .panel-title {
      font-size: 0.85rem;
      font-weight: 500;
      color: var(--flec-moment-title);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .panel-artist {
      margin-top: 2px;
      font-size: 0.75rem;
      color: var(--flec-moment-date);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
}

@keyframes rotate {
  to {
    transform: rotate(360deg);
  }
}

@keyframes mini-lyric-marquee {
  from {
    transform: translateX(0);
  }

  to {
    transform: translateX(var(--marquee-shift, 0px));
  }
}

.panel-lyric {
  margin-top: 10px;
  text-align: center;
  overflow: hidden;

  .lyric-current {
    font-size: 0.75rem;
    color: var(--theme-color);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .lyric-next {
    margin-top: 2px;
    font-size: 0.68rem;
    color: var(--flec-moment-date);
    opacity: 0.6;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

.panel-progress {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 10px;

  .time {
    font-size: 0.65rem;
    color: var(--flec-moment-date);
    flex-shrink: 0;
  }

  .progress-track {
    flex: 1;
    height: 3px;
    border-radius: 2px;
    background: rgba(128, 128, 128, 0.25);
    cursor: pointer;

    .progress-played {
      height: 100%;
      border-radius: 2px;
      background: var(--theme-color);
      transition: width 0.1s linear;
    }
  }
}

.panel-controls {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 8px;

  .ctrl {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    padding: 0;
    border: none;
    border-radius: 4px;
    background: transparent;
    color: var(--flec-moment-font);
    cursor: pointer;
    transition: all 0.2s;

    i {
      font-size: 0.9rem;
    }

    .spinning {
      animation: rotate 1s linear infinite;
    }

    &:hover {
      color: var(--theme-color);
      background: rgba(128, 128, 128, 0.12);
    }

    &.active {
      color: var(--theme-color);
    }

    &.ctrl-play {
      width: 28px;
      height: 28px;
      background: var(--flec-btn);
      color: #fff;

      &:hover {
        background: var(--flec-btn-hover);
        color: #fff;
      }
    }
  }

  .volume {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-left: auto;
    color: var(--flec-moment-date);

    i {
      font-size: 0.85rem;
    }

    input[type='range'] {
      width: 56px;
      height: 3px;
      accent-color: var(--theme-color);
      cursor: pointer;
    }
  }
}

.panel-list {
  margin-top: 8px;
  padding-top: 6px;
  max-height: 130px;
  overflow-y: auto;
  border-top: 1px solid var(--flec-border);

  &::-webkit-scrollbar {
    width: 3px;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(128, 128, 128, 0.35);
    border-radius: 2px;
  }

  .list-item {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 4px 6px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.75rem;
    transition: background 0.2s;

    &:hover {
      background: rgba(128, 128, 128, 0.12);
    }

    &.active {
      color: var(--theme-color);
      background: rgba(128, 128, 128, 0.1);
    }

    .item-index {
      width: 16px;
      flex-shrink: 0;
      text-align: center;
      font-size: 0.7rem;
      color: var(--flec-moment-date);
    }

    .item-name {
      flex: 1;
      min-width: 0;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .item-artist {
      max-width: 76px;
      color: var(--flec-moment-date);
      font-size: 0.7rem;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
}

// 移动端
@media (max-width: 768px) {
  .music-panel {
    right: 52px;
    left: 16px;
    bottom: 20px;
    width: auto;
    max-width: none;
    max-height: calc(100vh - 40px);
    overflow-y: auto;
    padding: 14px;

    &.is-minimized {
      width: auto;
      max-width: none;
      overflow: visible;
      padding: 10px 12px;
    }

    .panel-lyric {
      .lyric-current {
        font-size: 0.8rem;
      }

      .lyric-next {
        font-size: 0.72rem;
      }
    }

    .panel-progress {
      .progress-track {
        height: 5px;
      }
    }

    .panel-controls {
      gap: 8px;
      margin-top: 12px;

      .ctrl {
        width: 34px;
        height: 34px;

        i {
          font-size: 1.05rem;
        }

        &.ctrl-play {
          width: 38px;
          height: 38px;
        }
      }

      .volume input[type='range'] {
        width: 72px;
      }
    }

    .panel-list {
      max-height: 180px;

      .list-item {
        padding: 8px 6px;
        font-size: 0.8rem;
      }
    }
  }
}
</style>
