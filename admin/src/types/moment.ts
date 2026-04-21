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

// 动态列表查询参数
export interface MomentListQuery {
  page: number;
  page_size: number;
  keyword?: string; // 搜索关键词（文本内容）
  tags?: string; // 标签
  location?: string; // 发布地点
  is_publish?: boolean; // 是否发布
  has_images?: boolean; // 是否有图片
  has_video?: boolean; // 是否有视频
  has_music?: boolean; // 是否有音乐
  has_link?: boolean; // 是否有链接
  start_time?: string; // 发布开始时间
  end_time?: string; // 发布结束时间
}

// 分页数据
export interface MomentListData {
  list: Moment[];
  total: number;
  page: number;
  page_size: number;
}
