/**
 * 友链类型
 */
export interface FriendType {
  id: number;
  name: string;
  is_visible: boolean;
  sort: number;
}

/**
 * 友链数据结构
 */
export interface Friend {
  id: number;
  name: string;
  url: string;
  description?: string;
  avatar?: string;
  screenshot?: string;
  sort: number;
  is_invalid: boolean;
  type_id?: number;
  type?: FriendType;
}

/**
 * 友链分组（用于展示）
 */
export interface FriendGroup {
  type_id: number | null;
  type_name: string;
  type_sort: number;
  friends: Friend[];
}

/**
 * 友链分组响应
 */
export interface FriendGroupedResponse {
  groups: FriendGroup[];
  total_groups: number;
  total_friends: number;
}

/**
 * 友链查询参数
 */
export interface FriendQueryParams {
  page?: number;
  page_size?: number;
}

/**
 * 友圈文章作者
 */
export interface FriendCircleAuthor {
  id: number;
  name: string;
  avatar?: string;
  url: string;
}

/**
 * 友圈文章条目
 */
export interface FriendCircleItem {
  id: number;
  title: string;
  link: string;
  published_at?: string;
  author?: FriendCircleAuthor;
}

/**
 * 友圈统计
 */
export interface FriendCircleStats {
  total: number;
  site_count: number;
  today_count: number;
}

/**
 * 友圈响应
 */
export interface FriendCircleResponse {
  list: FriendCircleItem[];
  total: number;
  page: number;
  page_size: number;
}

/**
 * 友圈查询参数
 */
export interface FriendCircleQueryParams {
  page?: number;
  page_size?: number;
  keyword?: string;
}

/**
 * 友链申请请求
 */
export interface FriendApplyRequest {
  name: string; // 网站名称
  url: string; // 网站链接
  description: string; // 网站描述
  avatar: string; // 网站头像/logo
  screenshot?: string; // 网站截图（可选）
}
