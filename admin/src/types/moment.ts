// 动态内容类型
export interface MomentContent {
  text?: string;
  tags?: string;
  images?: string[];
  video?: {
    url: string;
    platform?: string;
    video_id?: string;
  };
  music?: {
    server: string;
    type: string;
    id: string;
  };
  link?: {
    url: string;
    title?: string;
    favicon?: string;
  };
  location?: string;
  book?: Record<string, unknown>;
  movie?: Record<string, unknown>;
}

// 动态实体
export interface Moment {
  id: number;
  content: MomentContent;
  is_publish: boolean; // 是否发布
  publish_time: string;
  created_at: string;
  updated_at: string;
}

// 创建动态请求
export interface CreateMomentRequest {
  content: MomentContent;
  publish_time?: string;
  is_publish?: boolean;
}

// 更新动态请求
export interface UpdateMomentRequest {
  content?: MomentContent;
  publish_time?: string;
  is_publish?: boolean;
}

// 获取链接信息请求
export interface FetchLinkRequest {
  url: string;
}

// 链接信息响应
export interface LinkInfo {
  title: string;
  favicon: string;
}

// 解析视频请求
export interface ParseVideoRequest {
  url: string;
}

// 视频信息响应
export interface VideoInfo {
  platform: string;
  video_id: string;
}

// 分页数据
export interface MomentListData {
  list: Moment[];
  total: number;
  page: number;
  page_size: number;
}
