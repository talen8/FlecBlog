import type { MusicSource, AudioTrack } from '../../types/music';
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
export function useMusic(music: MaybeRefOrGetter<MusicSource>) {
  const { fetchTracks, fetchLyrics } = useMeting();

  const tracks = ref<AudioTrack[]>([]);
  const loading = ref(true);
  const error = ref(false);

  const load = async () => {
    loading.value = true;
    try {
      tracks.value = await fetchTracks(toValue(music));
      error.value = false;
    } catch {
      error.value = true;
    } finally {
      loading.value = false;
    }
  };

  return { tracks, loading, error, load, fetchLyrics };
}
