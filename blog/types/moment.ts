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
