<script lang="ts" setup>
import type { MomentMusic } from '@@/types/moment';

interface AudioTrack {
  name: string;
  artist: string;
  url: string;
  cover: string;
  lrc?: string;
}

interface MusicApiResponse {
  name?: string;
  title?: string;
  artist?: string;
  author?: string;
  url: string;
  pic?: string;
  cover?: string;
  lrc?: string;
}

interface LyricLine {
  time: number;
  text: string;
}

const props = defineProps<{
  music: MomentMusic;
}>();

// 播放器状态
const audioRef = ref<HTMLAudioElement | null>(null);
const isPlaying = ref(false);
const currentTime = ref(0);
const duration = ref(0);
const currentIndex = ref(0);
const isLoading = ref(true);
const showPlaylist = ref(false);
const audioList = ref<AudioTrack[]>([]);
const loadError = ref(false);

// 歌词相关
const lyrics = ref<LyricLine[]>([]);
const currentLyricIndex = ref(-1);

// 计算属性
const currentTrack = computed(() => audioList.value[currentIndex.value] || null);
const hasPlaylist = computed(() => audioList.value.length > 1);
const progress = computed(() => (duration.value ? (currentTime.value / duration.value) * 100 : 0));

// 格式化时间为 mm:ss
const formatTime = (seconds: number) => {
  if (!isFinite(seconds) || isNaN(seconds)) return '00:00';
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
};

// 解析 LRC 格式歌词
const parseLyrics = (lrcText: string): LyricLine[] => {
  if (!lrcText) return [];
  const lines = lrcText.split('\n');
  const lyricLines: LyricLine[] = [];

  for (const line of lines) {
    const match = line.match(/\[(\d{2}):(\d{2})(?:\.(\d{2,3}))?\](.*)/);
    if (match && match[1] && match[2] && match[4]) {
      const minutes = parseInt(match[1]);
      const seconds = parseInt(match[2]);
      const milliseconds = match[3] ? parseInt(match[3].padEnd(3, '0')) : 0;
      const text = match[4].trim();
      if (text) {
        lyricLines.push({ time: minutes * 60 + seconds + milliseconds / 1000, text });
      }
    }
  }
  return lyricLines.sort((a, b) => a.time - b.time);
};

// 获取歌词文本（从 URL）
const fetchLyrics = async (lrcUrl: string): Promise<string> => {
  try {
    const response = await fetch(lrcUrl);
    return await response.text();
  } catch {
    return '';
  }
};

// 获取音乐数据
const fetchMusicData = async () => {
  try {
    const { server, type, id } = props.music;
    const response = await fetch(
      `https://api.injahow.cn/meting/?server=${server}&type=${type}&id=${id}`
    );
    const data = await response.json();
    const list = Array.isArray(data) ? data : [data];

    audioList.value = list.map((item: MusicApiResponse) => ({
      name: item.name || item.title || '未知歌曲',
      artist: item.artist || item.author || '未知艺术家',
      url: item.url,
      cover: item.pic || item.cover || '',
      lrc: item.lrc || '',
    }));

    // 获取并解析第一首歌的歌词
    if (audioList.value[0]?.lrc) {
      const lrcUrl = audioList.value[0].lrc;
      const lrcText = lrcUrl.startsWith('http') ? await fetchLyrics(lrcUrl) : lrcUrl;
      lyrics.value = parseLyrics(lrcText);
    }
    loadError.value = false;
  } catch {
    loadError.value = true;
  } finally {
    isLoading.value = false;
  }
};

// 播放/暂停切换
const togglePlay = () => {
  if (!audioRef.value || !currentTrack.value?.url) return;
  if (isPlaying.value) {
    audioRef.value.pause();
  } else {
    audioRef.value.play().catch(() => (loadError.value = true));
  }
};

// 播放指定歌曲
const playTrack = async (index: number) => {
  if (index < 0 || index >= audioList.value.length) return;
  currentIndex.value = index;
  showPlaylist.value = false;

  // 更新歌词
  const track = audioList.value[index];
  if (track?.lrc) {
    const lrcText = track.lrc.startsWith('http') ? await fetchLyrics(track.lrc) : track.lrc;
    lyrics.value = parseLyrics(lrcText);
  } else {
    lyrics.value = [];
  }
  currentLyricIndex.value = -1;
  nextTick(() => audioRef.value?.play().catch(() => {}));
};

// 进度条点击跳转
const seekTo = (e: MouseEvent) => {
  if (!audioRef.value || !duration.value) return;
  const target = e.currentTarget as HTMLElement;
  const rect = target.getBoundingClientRect();
  audioRef.value.currentTime =
    Math.max(0, Math.min(1, (e.clientX - rect.left) / rect.width)) * duration.value;
};

// 更新当前歌词行
const updateCurrentLyric = (time: number) => {
  if (lyrics.value.length === 0) return;
  let index = -1;
  for (let i = 0; i < lyrics.value.length; i++) {
    const lyric = lyrics.value[i];
    if (lyric && time >= lyric.time) index = i;
    else break;
  }
  if (index !== currentLyricIndex.value) currentLyricIndex.value = index;
};

// 音频事件处理
const onTimeUpdate = () => {
  if (audioRef.value) {
    currentTime.value = audioRef.value.currentTime;
    updateCurrentLyric(currentTime.value);
  }
};

const onLoadedMetadata = () => {
  if (audioRef.value) duration.value = audioRef.value.duration;
};

// 播放结束处理
const onEnded = () => {
  if (hasPlaylist.value) {
    playTrack((currentIndex.value + 1) % audioList.value.length);
  } else {
    isPlaying.value = false;
    currentTime.value = 0;
  }
};

/**
 * 处理音频加载错误
 */
const handleAudioError = () => {
  loadError.value = true;
  isPlaying.value = false;
};

// 初始化
onMounted(() => fetchMusicData());
</script>

<template>
  <div class="flec-music-player">
    <div v-if="isLoading" class="player-status">
      <i class="ri-loader-4-line spin" />
      <span>加载中...</span>
    </div>

    <div v-else-if="loadError || !currentTrack" class="player-status">
      <i class="ri-error-warning-line" />
      <span>音乐加载失败</span>
    </div>

    <template v-else>
      <audio
        ref="audioRef"
        :src="currentTrack.url"
        preload="metadata"
        @timeupdate="onTimeUpdate"
        @loadedmetadata="onLoadedMetadata"
        @play="isPlaying = true"
        @pause="isPlaying = false"
        @ended="onEnded"
        @error="handleAudioError"
      />

      <div class="player-main">
        <div class="player-left">
          <div class="player-cover">
            <img v-if="currentTrack.cover" :src="currentTrack.cover" alt="cover" />
            <div v-else class="cover-placeholder">
              <i class="ri-music-2-fill" />
            </div>
          </div>
          <button class="play-btn" @click="togglePlay">
            <i :class="isPlaying ? 'ri-pause-fill' : 'ri-play-fill'" />
          </button>
        </div>

        <div class="player-right">
          <div class="player-info">
            <span class="track-name">{{ currentTrack.name }}</span>
            <span class="separator"> - </span>
            <span class="track-artist">{{ currentTrack.artist }}</span>
          </div>

          <div v-if="lyrics.length > 0" class="player-lyrics">
            <div
              v-for="(line, index) in lyrics"
              :key="index"
              class="lyric-line"
              :class="{
                active: index === currentLyricIndex,
                next: index === currentLyricIndex + 1,
              }"
            >
              {{ line.text }}
            </div>
          </div>

          <div class="player-progress">
            <div class="progress-bar" @click="seekTo">
              <div class="progress-played" :style="{ width: `${progress}%` }" />
            </div>
            <div class="progress-time">
              {{ formatTime(currentTime) }} / {{ formatTime(duration) }}
            </div>
            <button v-if="hasPlaylist" class="ctrl-btn" @click="showPlaylist = !showPlaylist">
              <i class="ri-menu-fill" />
            </button>
          </div>
        </div>
      </div>

      <div v-if="showPlaylist && hasPlaylist" class="player-playlist">
        <div
          v-for="(track, index) in audioList"
          :key="index"
          class="playlist-item"
          :class="{ active: index === currentIndex }"
          @click="playTrack(index)"
        >
          <span class="item-index">
            <i v-if="index === currentIndex && isPlaying" class="ri-equalizer-fill" />
            <template v-else>{{ index + 1 }}</template>
          </span>
          <span class="item-name">{{ track.name }}</span>
          <span class="item-artist">{{ track.artist }}</span>
        </div>
      </div>
    </template>
  </div>
</template>

<style lang="scss" scoped>
.flec-music-player {
  background: var(--flec-moment-card-bg);
  border-radius: 6px;
  overflow: hidden;
}

.player-status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 80px;
  color: var(--flec-moment-date);
  font-size: 0.85rem;

  i {
    font-size: 1.2rem;
  }
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.player-main {
  display: flex;
  height: 80px;
}

.player-left {
  position: relative;
  flex-shrink: 0;
  height: 100%;

  .player-cover {
    height: 100%;
    aspect-ratio: 1;
    background: #eee;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .cover-placeholder {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      color: #fff;
      font-size: 1.6rem;
    }
  }

  .play-btn {
    position: absolute;
    left: 50%;
    top: 50%;
    transform: translate(-50%, -50%);
    width: 26px;
    height: 26px;
    border: none;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.9);
    color: #49b1f5;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
    transition: all 0.2s;

    i {
      font-size: 14px;
      margin-left: 1px;
    }

    &:hover {
      transform: translate(-50%, -50%) scale(1.1);
      background: #fff;
    }
  }
}

.player-right {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0 12px;
}

.player-info {
  display: flex;
  align-items: baseline;
  margin-bottom: 6px;
  overflow: hidden;

  .track-name {
    font-size: 0.9rem;
    font-weight: 500;
    color: var(--flec-moment-title);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .separator {
    color: var(--flec-moment-date);
    margin: 0 4px;
    flex-shrink: 0;
  }

  .track-artist {
    font-size: 0.85rem;
    color: var(--flec-moment-date);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

.player-lyrics {
  height: 32px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  .lyric-line {
    text-align: center;
    font-size: 0.7rem;
    color: var(--flec-moment-date);
    transition: all 0.3s ease;
    line-height: 1.3;
    opacity: 0;
    display: none;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
    padding: 0 4px;

    &.active {
      color: #49b1f5;
      opacity: 1;
      display: block;
      font-weight: 500;
    }

    &.next {
      color: var(--flec-moment-date);
      opacity: 0.5;
      display: block;
      font-size: 0.65rem;
    }
  }
}

.ctrl-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border: none;
  background: transparent;
  color: var(--flec-moment-date);
  cursor: pointer;
  border-radius: 3px;
  transition: all 0.2s;
  padding: 0;

  i {
    font-size: 0.85rem;
  }

  &:hover {
    color: #49b1f5;
  }
}

.player-progress {
  display: flex;
  align-items: center;
  gap: 6px;

  .progress-bar {
    flex: 1;
    height: 3px;
    background: rgba(0, 0, 0, 0.1);
    border-radius: 2px;
    cursor: pointer;

    .progress-played {
      height: 100%;
      background: #49b1f5;
      border-radius: 2px;
      transition: width 0.1s linear;
    }
  }

  .progress-time {
    font-size: 0.7rem;
    color: var(--flec-moment-date);
    white-space: nowrap;
  }
}

.player-playlist {
  max-height: 120px;
  overflow-y: auto;
  border-top: 1px solid var(--flec-moment-divider, rgba(0, 0, 0, 0.06));
  background: rgba(0, 0, 0, 0.02);

  &::-webkit-scrollbar {
    width: 3px;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(0, 0, 0, 0.15);
    border-radius: 2px;

    &:hover {
      background: rgba(0, 0, 0, 0.25);
    }
  }

  .playlist-item {
    display: flex;
    align-items: center;
    padding: 6px 10px;
    gap: 8px;
    cursor: pointer;
    font-size: 0.8rem;
    transition: background 0.2s;

    &:hover {
      background: rgba(73, 177, 245, 0.08);
    }

    &.active {
      background: rgba(73, 177, 245, 0.12);

      .item-name,
      .item-index {
        color: #49b1f5;
      }
    }

    .item-index {
      width: 18px;
      text-align: center;
      color: var(--flec-moment-date);
      flex-shrink: 0;
      font-size: 0.75rem;

      i {
        font-size: 0.85rem;
      }
    }

    .item-name {
      flex: 1;
      min-width: 0;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      color: var(--flec-moment-font);
    }

    .item-artist {
      color: var(--flec-moment-date);
      font-size: 0.75rem;
      white-space: nowrap;
      max-width: 100px;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }
}

@media screen and (max-width: 768px) {
  .player-left .player-cover {
    width: 56px;
    height: 56px;
  }

  .player-main {
    min-height: 56px;
  }

  .player-status {
    height: 56px;
    font-size: 0.8rem;

    i {
      font-size: 1.1rem;
    }
  }

  .player-right {
    padding: 0 10px;
  }

  .player-info {
    .track-name {
      font-size: 0.85rem;
    }

    .track-artist {
      font-size: 0.8rem;
    }
  }

  .player-lyrics {
    height: 28px;

    .lyric-line {
      font-size: 0.65rem;
      line-height: 1.2;

      &.active {
        font-size: 0.65rem;
      }

      &.next {
        font-size: 0.6rem;
      }
    }
  }

  .player-playlist {
    max-height: 100px;

    .playlist-item {
      padding: 5px 8px;
      gap: 6px;
      font-size: 0.75rem;

      .item-index {
        width: 16px;
        font-size: 0.7rem;

        i {
          font-size: 0.8rem;
        }
      }

      .item-artist {
        font-size: 0.7rem;
        max-width: 80px;
      }
    }
  }
}
</style>
