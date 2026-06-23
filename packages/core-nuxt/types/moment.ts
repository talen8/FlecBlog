/**
 * 视频内容
 */
export interface MomentVideo {
  url: string;
  platform?: 'youtube' | 'bilibili' | 'local';
  video_id?: string;
}

/**
 * 音频内容
 */
export interface MomentAudio {
  url: string;
}

/**
 * 音乐内容
 */
export interface MomentMusic {
  server: 'netease' | 'tencent';
  type: 'song' | 'playlist' | 'album' | 'artist';
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

/**
 * 链接内容
 */
export interface MomentLink {
  url: string;
  title: string;
  favicon?: string;
}

/**
 * 动态内容
 */
export interface MomentContent {
  text?: string;
  images?: string[];
  location?: string;
  tags?: string;
  video?: MomentVideo;
  audio?: MomentAudio;
  music?: MomentMusic;
  link?: MomentLink;
}

/**
 * 动态
 */
export interface Moment {
  id: number;
  content: MomentContent;
  publish_time: string;
}

/**
 * 动态列表响应
 */
export interface MomentListResponse {
  list: Moment[];
  total: number;
  page: number;
  page_size: number;
}
