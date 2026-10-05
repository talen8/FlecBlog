import { effectScope } from 'vue';
import type { Ref } from 'vue';
import type {
  AudioTrack,
  LyricLine,
  MusicServer,
  MusicSource,
  MusicType,
  PlayMode,
} from '../../types/music';
import { fetchLrc } from './useMeting';

/** 播放状态持久化键 */
const STORAGE_KEY = 'flec:global-music';
/** 默认音量 */
const DEFAULT_VOLUME = 0.7;
/** 持久化写入的最小间隔（毫秒） */
const PERSIST_INTERVAL = 2000;
/** 判定不可用后切换到下一首前的等待时长（毫秒），留出可被察觉的停顿 */
const SWITCH_DELAY = 1500;
/** 全局播放器在发声权仲裁中的标识 */
export const GLOBAL_PLAYBACK_SLOT = 'global-music';
/** 音乐平台缺省值 */
const DEFAULT_SERVER: MusicServer = 'netease';
/** 歌单类型固定为歌单 */
const MUSIC_TYPE: MusicType = 'playlist';
/** 默认记忆进度并续播 */
const DEFAULT_AUTOPLAY = true;
/** 默认顺序播放 */
const DEFAULT_ORDER: PlayMode = 'list';
/** 已加载过的歌单标识，避免重复请求 */
let loadedSourceId = '';

/** 全局播放器状态 */
interface EngineState {
  /** 音轨列表 */
  tracks: AudioTrack[];
  /** 当前音轨下标 */
  index: number;
  /** 是否正在播放 */
  playing: boolean;
  /** 当前播放秒数 */
  currentTime: number;
  /** 当前音轨总秒数 */
  duration: number;
  /** 音量 0~1 */
  volume: number;
  /** 播放顺序 */
  mode: PlayMode;
  /** 当前音轨歌词 */
  lyrics: LyricLine[];
  /** 当前歌词行下标 */
  lyricIndex: number;
  /** 歌单加载中 */
  loading: boolean;
  /** 当前音轨缓冲中 */
  buffering: boolean;
  /** 加载失败（歌单为空或 Meting 不可用） */
  failed: boolean;
}

/** 加载歌单的可选项 */
interface LoadOptions {
  /** 是否按上次状态续播 */
  autoplay?: boolean;
  /** 默认播放顺序（无持久化记录时生效） */
  order?: PlayMode;
}

/** 播放器操作集合 */
interface EngineApi {
  load: (source: MusicSource, options?: LoadOptions) => Promise<void>;
  playAt: (index: number, autoplay?: boolean) => Promise<void>;
  play: () => Promise<void>;
  pause: () => void;
  toggle: () => void;
  next: (auto?: boolean) => void;
  prev: () => void;
  seekRatio: (ratio: number) => void;
  setVolume: (volume: number) => void;
  cycleMode: () => void;
}

/** 已持久化的播放状态 */
interface PersistedState {
  index: number;
  time: number;
  volume: number;
  mode: PlayMode;
  playing: boolean;
}

/** 当前活跃的引擎运行时（audio 事件回调通过它访问全局状态） */
let runtime: {
  state: Ref<EngineState>;
  owner: Ref<string | null>;
  api: EngineApi;
} | null = null;

/** 全局唯一的音频元素 */
let audioEl: HTMLAudioElement | null = null;
/** 等待写入的续播进度（元数据就绪后生效） */
let pendingSeek = 0;
/** 失败后延迟切换音轨的计时器 */
let failureTimer: ReturnType<typeof setTimeout> | null = null;
/** 本轮已判定为不可播放的音轨下标 */
const failedTracks = new Set<number>();
/** 是否已请求过播放，决定自动切歌后是否续播 */
let playbackRequested = false;

/** 创建初始状态 */
const createState = (): EngineState => ({
  tracks: [],
  index: 0,
  playing: false,
  currentTime: 0,
  duration: 0,
  volume: DEFAULT_VOLUME,
  mode: 'list',
  lyrics: [],
  lyricIndex: -1,
  loading: false,
  buffering: false,
  failed: false,
});

/** 持久化监听是否已绑定 */
let persistBound = false;

/**
 * 绑定播放状态到本地存储的节流写入
 * 引擎会被多个组件调用，此处保证监听全站只注册一次
 * 使用独立作用域绑定，避免首个调用方组件卸载后监听失效
 * @param state - 播放器状态
 */
function bindPersist(state: Ref<EngineState>) {
  if (persistBound || !import.meta.client) return;
  persistBound = true;

  let timer: ReturnType<typeof setTimeout> | null = null;
  let lastWrite = 0;

  /** 写入当前播放状态到本地存储 */
  const persist = () => {
    const s = state.value;
    localStorage.setItem(
      STORAGE_KEY,
      JSON.stringify({
        index: s.index,
        time: s.currentTime,
        volume: s.volume,
        mode: s.mode,
        playing: s.playing,
      } satisfies PersistedState)
    );
    lastWrite = Date.now();
  };

  effectScope(true).run(() => {
    watch(
      () => [
        state.value.index,
        state.value.currentTime,
        state.value.volume,
        state.value.mode,
        state.value.playing,
      ],
      () => {
        if (timer) return;
        const wait = Math.max(0, PERSIST_INTERVAL - (Date.now() - lastWrite));
        timer = setTimeout(() => {
          timer = null;
          persist();
        }, wait);
      }
    );
  });
}

/**
 * 根据当前播放时间更新高亮歌词行
 * @param state - 播放器状态
 */
const syncLyric = (state: EngineState) => {
  if (!state.lyrics.length) return;
  let index = -1;
  for (let i = 0; i < state.lyrics.length; i++) {
    const lyric = state.lyrics[i];
    if (lyric && state.currentTime >= lyric.time) index = i;
    else break;
  }
  if (index !== state.lyricIndex) state.lyricIndex = index;
};

/** 取消音轨可用性监控 */
function clearFailureTimer() {
  if (failureTimer) clearTimeout(failureTimer);
  failureTimer = null;
}

/**
 * 判断播放失败是否源于浏览器拦截自动播放
 * 只有等待用户手势的情况需要保留当前音轨，其余失败按不可用处理
 * @param error - el.play() 抛出的异常
 * @returns 属于自动播放拦截时返回 true
 */
function isAutoplayBlocked(error: unknown): boolean {
  return error instanceof DOMException && error.name === 'NotAllowedError';
}

/**
 * 查找下一个未被判定为不可播放的音轨
 * @param total - 音轨总数
 * @param from - 起始下标
 * @returns 下一个可播放下标，全部不可播时返回 -1
 */
function nextPlayableIndex(total: number, from: number): number {
  for (let step = 1; step <= total; step++) {
    const index = (from + step) % total;
    if (!failedTracks.has(index)) return index;
  }
  return -1;
}

/**
 * 延迟切换到下一个可播音轨
 * 失败后留出一段停顿再切，避免瞬间跳走让用户以为歌被跳过
 * @param index - 目标音轨下标
 */
function scheduleFailureSwitch(index: number) {
  clearFailureTimer();
  failureTimer = setTimeout(() => {
    failureTimer = null;
    runtime?.api.playAt(index, true);
  }, SWITCH_DELAY);
}

/**
 * 当前音轨不可播放：标记为失败并切到下一个可播音轨
 * 全部音轨均不可播时停止播放并置失败态，避免无限重试
 */
function handleTrackFailure() {
  const rt = runtime;
  if (!rt) return;
  clearFailureTimer();

  const s = rt.state.value;
  if (!s.tracks.length) return;
  // 用户尚未请求播放（如仅恢复续播位置）时不切歌，等下次播放请求时再判定
  if (!playbackRequested) return;

  failedTracks.add(s.index);
  const nextIndex = nextPlayableIndex(s.tracks.length, s.index);
  if (nextIndex < 0) {
    s.playing = false;
    s.failed = true;
    return;
  }
  scheduleFailureSwitch(nextIndex);
}

/**
 * 获取（惰性创建）全局唯一音频元素
 * @returns 音频元素
 */
const getAudioElement = (): HTMLAudioElement => {
  if (audioEl) return audioEl;

  const el = new Audio();
  el.preload = 'metadata';

  el.addEventListener('timeupdate', () => {
    if (!runtime) return;
    runtime.state.value.currentTime = el.currentTime;
    runtime.state.value.buffering = false;
    syncLyric(runtime.state.value);
  });

  el.addEventListener('loadedmetadata', () => {
    if (!runtime) return;
    runtime.state.value.duration = el.duration;
    if (pendingSeek > 0 && Number.isFinite(el.duration)) {
      el.currentTime = Math.min(pendingSeek, el.duration);
      pendingSeek = 0;
    }
  });

  el.addEventListener('play', () => {
    if (!runtime) return;
    runtime.state.value.playing = true;
    runtime.owner.value = GLOBAL_PLAYBACK_SLOT;
  });

  el.addEventListener('pause', () => {
    if (!runtime) return;
    runtime.state.value.playing = false;
    runtime.state.value.buffering = false;
    if (runtime.owner.value === GLOBAL_PLAYBACK_SLOT) runtime.owner.value = null;
  });

  el.addEventListener('waiting', () => {
    if (runtime) runtime.state.value.buffering = true;
  });

  el.addEventListener('playing', () => {
    if (runtime) runtime.state.value.buffering = false;
  });

  el.addEventListener('ended', () => {
    clearFailureTimer();
    runtime?.api.next(true);
  });

  el.addEventListener('error', () => {
    if (!runtime) return;
    runtime.state.value.playing = false;
  });

  audioEl = el;
  return el;
};

/**
 * 读取持久化的播放状态
 * @returns 持久化状态，不存在或已损坏时返回 null
 */
const readPersisted = (): PersistedState | null => {
  if (!import.meta.client) return null;
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw) as Partial<PersistedState>;
    if (typeof parsed.index !== 'number' || typeof parsed.time !== 'number') return null;
    return {
      index: Math.max(0, parsed.index),
      time: Math.max(0, parsed.time),
      volume:
        typeof parsed.volume === 'number' && parsed.volume >= 0 && parsed.volume <= 1
          ? parsed.volume
          : DEFAULT_VOLUME,
      mode: parsed.mode === 'random' || parsed.mode === 'single' ? parsed.mode : 'list',
      playing: parsed.playing === true,
    };
  } catch {
    return null;
  }
};

/**
 * 全局音乐播放引擎（单例）
 * 持有全站唯一的音频元素，负责歌单、进度、音量、播放顺序与歌词
 * @returns state - 播放器状态；currentTrack / progress / hasPlaylist - 派生状态
 * @returns load / playAt / play / pause / toggle / next / prev / seekRatio / setVolume / cycleMode - 操作
 */
export function useAudioEngine() {
  const state = useState<EngineState>('flec-audio-engine', createState);
  const owner = usePlaybackOwner();
  const { fetchTracks } = useMeting();

  const currentTrack = computed<AudioTrack | null>(
    () => state.value.tracks[state.value.index] ?? null
  );
  const hasPlaylist = computed(() => state.value.tracks.length > 1);
  const progress = computed(() =>
    state.value.duration ? (state.value.currentTime / state.value.duration) * 100 : 0
  );

  /**
   * 播放并接管发声权
   * 自动播放被拦截时静默转为暂停等待用户手势；播放失败则按音轨不可用处理
   */
  const play = async () => {
    const el = getAudioElement();
    if (!el.src) return;
    playbackRequested = true;
    state.value.buffering = true;
    try {
      await el.play();
    } catch (error) {
      state.value.playing = false;
      state.value.buffering = false;
      if (isAutoplayBlocked(error)) clearFailureTimer();
      else handleTrackFailure();
    }
  };

  /** 暂停播放并释放发声权 */
  const pause = () => {
    clearFailureTimer();
    playbackRequested = false;
    getAudioElement().pause();
  };

  /** 切换播放 / 暂停 */
  const toggle = () => {
    if (state.value.playing) pause();
    else void play();
  };

  /**
   * 切换到指定音轨
   * @param index - 音轨下标
   * @param autoplay - 是否立即播放
   */
  const playAt = async (index: number, autoplay = true) => {
    const s = state.value;
    if (index < 0 || index >= s.tracks.length) return;

    clearFailureTimer();
    failedTracks.delete(index);

    s.index = index;
    s.currentTime = 0;
    s.duration = 0;
    s.lyrics = [];
    s.lyricIndex = -1;

    const track = s.tracks[index];
    if (!track) return;

    const el = getAudioElement();
    el.src = track.url;
    el.volume = s.volume;
    el.load();

    if (track.lrc) {
      const lyrics = await fetchLrc(track.lrc);
      // 切歌过快时丢弃过期结果
      if (state.value.index === index) state.value.lyrics = lyrics;
    }

    if (autoplay) await play();
  };

  /**
   * 计算下一首下标
   * @param auto - 是否由播放结束触发（单曲循环模式下重播当前曲目）
   * @returns 下一首下标，无音轨时返回 -1
   */
  const resolveNextIndex = (auto: boolean): number => {
    const s = state.value;
    const total = s.tracks.length;
    if (!total) return -1;
    if (s.mode === 'single' && auto) return s.index;
    if (s.mode === 'random') {
      if (total === 1) return 0;
      let next = s.index;
      while (next === s.index) next = Math.floor(Math.random() * total);
      return next;
    }
    return (s.index + 1) % total;
  };

  /**
   * 下一首
   * @param auto - 是否由播放结束触发（单曲循环模式下重播当前曲目）
   */
  const next = (auto = false) => {
    const s = state.value;
    const index = resolveNextIndex(auto);
    if (index >= 0) void playAt(index, auto || s.playing);
  };

  /** 上一首 */
  const prev = () => {
    const s = state.value;
    if (!s.tracks.length) return;
    void playAt((s.index - 1 + s.tracks.length) % s.tracks.length, s.playing);
  };

  /**
   * 按进度比例跳转
   * @param ratio - 0~1 的进度比例
   */
  const seekRatio = (ratio: number) => {
    const el = getAudioElement();
    if (!Number.isFinite(el.duration)) return;
    const clamped = Math.min(1, Math.max(0, ratio));
    el.currentTime = clamped * el.duration;
    state.value.currentTime = el.currentTime;
    syncLyric(state.value);
  };

  /**
   * 设置音量
   * @param volume - 0~1 的音量
   */
  const setVolume = (volume: number) => {
    const clamped = Math.min(1, Math.max(0, volume));
    state.value.volume = clamped;
    getAudioElement().volume = clamped;
  };

  /** 循环切换播放顺序：顺序 → 随机 → 单曲 */
  const cycleMode = () => {
    const order: PlayMode[] = ['list', 'random', 'single'];
    const current = order.indexOf(state.value.mode);
    state.value.mode = order[(current + 1) % order.length] ?? 'list';
  };

  /**
   * 加载歌单并恢复上次播放状态
   * @param source - Meting 音乐来源
   * @param options - 续播与默认播放顺序
   */
  const load = async (source: MusicSource, options: LoadOptions = {}) => {
    const s = state.value;
    if (!source?.id) {
      s.failed = true;
      return;
    }

    s.loading = true;
    s.failed = false;
    try {
      const tracks = await fetchTracks(source);
      if (!tracks.length) {
        s.failed = true;
        return;
      }

      const saved = readPersisted();
      failedTracks.clear();
      s.tracks = tracks;
      s.mode = saved?.mode ?? options.order ?? 'list';
      s.volume = saved?.volume ?? DEFAULT_VOLUME;
      pendingSeek = saved?.time ?? 0;

      const index = saved && saved.index < tracks.length ? saved.index : 0;
      await playAt(index, false);
      if (options.autoplay && saved?.playing) await play();
    } catch {
      s.failed = true;
    } finally {
      s.loading = false;
    }
  };

  const api: EngineApi = {
    load,
    playAt,
    play,
    pause,
    toggle,
    next,
    prev,
    seekRatio,
    setVolume,
    cycleMode,
  };

  // 注册为当前活跃运行时，供 audio 事件回调使用
  // 仅保存全局单例状态，无需随组件销毁清理
  if (import.meta.client) {
    runtime = { state, owner, api };
  }

  // 状态变化时节流写入本地存储（全站仅注册一次）
  bindPersist(state);

  return {
    state,
    currentTrack,
    hasPlaylist,
    progress,
    ...api,
  };
}

/**
 * 全站发声权状态（全局唯一）
 * @returns 当前持有发声权的播放器标识
 */
function usePlaybackOwner() {
  return useState<string | null>('flec-playback-owner', () => null);
}

/**
 * 申请一个播放位，用于保证同一时刻只有一个播放器发声
 * 当其他播放器取得发声权时，会通过 onEvict 通知当前播放器暂停
 * @param id - 当前播放器标识
 * @param onEvict - 失去发声权时的回调（应暂停播放）
 * @returns claim - 取得发声权；release - 释放发声权
 */
export function usePlaybackSlot(id: string, onEvict: () => void) {
  const owner = usePlaybackOwner();

  watch(owner, value => {
    if (value && value !== id) onEvict();
  });

  onScopeDispose(() => {
    if (owner.value === id) owner.value = null;
  });

  return {
    /** 取得发声权 */
    claim: () => {
      owner.value = id;
    },
    /** 释放发声权（仅当自己仍持有） */
    release: () => {
      if (owner.value === id) owner.value = null;
    },
  };
}

/**
 * 读取枚举型主题配置，非法值回落为缺省值
 * @param config - 主题配置
 * @param key - 配置键
 * @param candidates - 合法取值
 * @param fallback - 缺省值
 */
const readEnum = <T extends string>(
  config: Record<string, unknown>,
  key: string,
  candidates: readonly T[],
  fallback: T
): T => {
  const value = config[key];
  return typeof value === 'string' && candidates.includes(value as T) ? (value as T) : fallback;
};

/**
 * 读取歌单 ID（后台可能存为字符串或数字）
 * @param config - 主题配置
 * @returns 歌单 ID，缺失或非法时返回空字符串
 */
const readId = (config: Record<string, unknown>): string => {
  const value = config.music_id;
  if (typeof value === 'number' && Number.isFinite(value)) return String(value).trim();
  if (typeof value === 'string') return value.trim();
  return '';
};

/**
 * 全局音乐播放器：解析主题配置并驱动播放引擎
 * 每个页面只应在悬浮按钮处调用一次，配置变化时重新拉取歌单
 * @returns source - 歌单来源；enabled - 是否已配置歌单
 * @returns 其余为播放引擎的状态与操作
 */
export function useGlobalMusic() {
  const { themeConfig } = useTheme();
  const engine = useAudioEngine();

  const source = computed<MusicSource>(() => {
    const raw = (themeConfig.value ?? {}) as Record<string, unknown>;
    return {
      server: readEnum(raw, 'music_server', ['netease', 'tencent'] as const, DEFAULT_SERVER),
      type: MUSIC_TYPE,
      id: readId(raw),
    };
  });

  const enabled = computed(() => !!source.value.id);

  /** 加载歌单，同一歌单只请求一次 */
  const bootstrap = async () => {
    if (!enabled.value) return;
    const target = source.value;
    if (loadedSourceId === target.id) return;
    loadedSourceId = target.id;
    await engine.load(target, { autoplay: DEFAULT_AUTOPLAY, order: DEFAULT_ORDER });
  };

  onMounted(bootstrap);
  watch(
    () => source.value.id,
    () => void bootstrap()
  );

  return { source, enabled, ...engine };
}
