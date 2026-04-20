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
 * 友链申请请求
 */
export interface FriendApplyRequest {
  name: string; // 网站名称
  url: string; // 网站链接
  description: string; // 网站描述
  avatar: string; // 网站头像/logo
  screenshot?: string; // 网站截图（可选）
}
