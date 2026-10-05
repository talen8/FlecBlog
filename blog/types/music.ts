/** 音乐平台 */
export type MusicServer = 'netease' | 'tencent';

/** 音乐源类型 */
export type MusicType = 'song' | 'playlist' | 'album' | 'artist';

/** 播放顺序 */
export type PlayMode = 'list' | 'random' | 'single';

/**
 * 音乐来源（Meting 三元组）
 */
export interface MusicSource {
  server: MusicServer;
  type: MusicType;
  id: string;
}

/** meting api 原始响应项 */
export interface MusicApiResponse {
  name?: string;
  title?: string;
  artist?: string;
  author?: string;
  url: string;
  pic?: string;
  cover?: string;
  lrc?: string;
}

/** 播放器解析后的音轨 */
export interface AudioTrack {
  name: string;
  artist: string;
  url: string;
  cover: string;
  lrc?: string;
}

/** 歌词行 */
export interface LyricLine {
  time: number;
  text: string;
}
