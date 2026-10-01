import type { MomentMusic, AudioTrack, MusicApiResponse, LyricLine } from '../../types/moment';
import { getMoments } from './api/moment';

/**
 * 获取动态列表（支持 SSR）
 * @param pageSize - 每页数量，支持 ref / computed / getter，默认 30
 * @returns data - 动态数组，refresh - 刷新方法
 */
export function useMomentList(pageSize: MaybeRefOrGetter<number> = () => 30) {
  return useAsyncData(
    `moments:${toValue(pageSize)}`,
    async () => {
      const { list } = await getMoments({ page: 1, page_size: toValue(pageSize) });
      return list ?? [];
    },
    { watch: [() => toValue(pageSize)] }
  );
}

/**
 * 音乐播放器状态管理（通过 Meting API 加载音频信息）
 * @param music - 音乐数据 { server, type, id }，支持 ref / computed / getter
 * @returns tracks - 音轨列表
 * @returns loading / error - 加载状态
 * @returns load - 加载音乐
 * @returns fetchLyrics(lrc) - 解析歌词
 */
export function useMusic(music: MaybeRefOrGetter<MomentMusic>) {
  const { basicConfig } = useSysConfig();
  const metingApi = computed(() => basicConfig.value.meting_api || 'https://meting.flec.top/api');

  const tracks = ref<AudioTrack[]>([]);
  const loading = ref(true);
  const error = ref(false);

  const parseLyrics = (lrcText: string): LyricLine[] => {
    if (!lrcText) return [];
    const result: LyricLine[] = [];
    for (const line of lrcText.split('\n')) {
      const match = line.match(/\[(\d{2}):(\d{2})(?:\.(\d{2,3}))?\](.*)/);
      if (match && match[1] && match[2] && match[4]) {
        const text = match[4].trim();
        if (text) {
          const ms = match[3] ? parseInt(match[3].padEnd(3, '0')) : 0;
          result.push({ time: parseInt(match[1]) * 60 + parseInt(match[2]) + ms / 1000, text });
        }
      }
    }
    return result.sort((a, b) => a.time - b.time);
  };

  const fetchLyrics = async (lrc: string): Promise<LyricLine[]> => {
    if (!lrc) return [];
    try {
      const text = lrc.startsWith('http') ? await (await fetch(lrc)).text() : lrc;
      return parseLyrics(text);
    } catch {
      return [];
    }
  };

  const load = async () => {
    loading.value = true;
    try {
      const { server, type, id } = toValue(music);
      const res = await fetch(`${metingApi.value}?server=${server}&type=${type}&id=${id}`);
      const data = await res.json();
      const list = (Array.isArray(data) ? data : [data]) as MusicApiResponse[];
      tracks.value = list.map(item => ({
        name: item.name || item.title || '未知歌曲',
        artist: item.artist || item.author || '未知艺术家',
        url: item.url,
        cover: item.pic || item.cover || '',
        lrc: item.lrc || '',
      }));
      error.value = false;
    } catch {
      error.value = true;
    } finally {
      loading.value = false;
    }
  };

  return { tracks, loading, error, load, fetchLyrics };
}
