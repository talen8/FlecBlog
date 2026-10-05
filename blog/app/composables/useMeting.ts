import type { AudioTrack, LyricLine, MusicApiResponse, MusicSource } from '../../types/music';

/** Meting API 兜底地址 */
const DEFAULT_METING_API = 'https://meting.flec.top/api';

/**
 * 解析 LRC 歌词文本为按时间升序的歌词行
 * @param lrcText - LRC 文本内容
 * @returns 歌词行数组，按时间升序排列
 */
export function parseLyrics(lrcText: string): LyricLine[] {
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
}

/**
 * 拉取并解析歌词，支持 LRC 文本或歌词文件 URL
 * @param lrc - LRC 文本内容或歌词文件 URL
 * @returns 歌词行数组，失败时返回空数组
 */
export async function fetchLrc(lrc: string): Promise<LyricLine[]> {
  if (!lrc) return [];
  try {
    const text = lrc.startsWith('http') ? await (await fetch(lrc)).text() : lrc;
    return parseLyrics(text);
  } catch {
    return [];
  }
}

/**
 * Meting 音乐源访问
 * @returns fetchTracks(source) - 拉取音轨列表
 * @returns fetchLyrics(lrc) - 拉取并解析歌词，支持 LRC 文本或歌词 URL
 */
export function useMeting() {
  const { basicConfig } = useSysConfig();

  const apiBase = computed(() => basicConfig.value.meting_api?.trim() || DEFAULT_METING_API);

  /**
   * 拉取音乐源对应的音轨列表
   * @param source - Meting 三元组 { server, type, id }
   * @throws 请求或解析失败时抛出
   */
  const fetchTracks = async (source: MusicSource): Promise<AudioTrack[]> => {
    const { server, type, id } = source;
    const res = await fetch(`${apiBase.value}?server=${server}&type=${type}&id=${id}`);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const data = await res.json();
    const list = (Array.isArray(data) ? data : [data]) as MusicApiResponse[];
    return list
      .filter(item => !!item?.url)
      .map(item => ({
        name: item.name || item.title || '未知歌曲',
        artist: item.artist || item.author || '未知艺术家',
        url: item.url,
        cover: item.pic || item.cover || '',
        lrc: item.lrc || '',
      }));
  };

  return { fetchTracks, fetchLyrics: fetchLrc };
}
