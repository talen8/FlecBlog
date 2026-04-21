// ========== 友链类型相关 ==========

// 友链类型实体
export interface FriendType {
  id: number;
  name: string;
  sort: number;
  is_visible: boolean;
  count: number; // 该类型下的友链数量
}

// 友链类型列表响应数据
export interface FriendTypeListData {
  list: FriendType[];
  total: number;
  page: number;
  page_size: number;
}

// 创建友链类型请求
export interface CreateFriendTypeRequest {
  name: string;
  sort?: number;
  is_visible?: boolean;
}

// 更新友链类型请求
export interface UpdateFriendTypeRequest {
  name?: string;
  sort?: number;
  is_visible?: boolean;
}

// ========== 友链相关 ==========

// 友链实体
export interface Friend {
  id: number;
  name: string;
  url: string;
  description: string;
  avatar: string;
  screenshot: string; // 网站截图
  sort: number; // 排序值，范围1-10，默认5
  type_id: number | null; // 友链类型ID
  type_name?: string; // 友链类型名称
  is_invalid: boolean; // 是否失效
  is_pending: boolean; // 是否为待审核申请
  rss_url: string; // RSS订阅地址
  rss_latime?: string; // RSS订阅最后更新时间
  accessible: number; // 可访问性状态: 0=正常, -1=忽略检查, >0=连续失败次数
}

// 友链列表查询参数
export interface FriendQuery {
  page?: number;
  page_size?: number;
  keyword?: string; // 关键词搜索（名称、链接、描述）
  type_id?: number; // 友链类型ID
  is_invalid?: boolean; // 是否失效
  is_pending?: boolean; // 是否待审核
  accessible_status?: 'normal' | 'abnormal' | 'ignored'; // 可访问性状态
  rss_status?: 'no_rss' | 'normal' | 'warning' | 'danger'; // RSS状态
  has_screenshot?: boolean; // 是否包含截图
}

// 友链列表响应数据
export interface FriendListData {
  list: Friend[];
  total: number;
  page: number;
  page_size: number;
}

// 创建友链请求
export interface CreateFriendRequest {
  name: string;
  url: string;
  description?: string;
  avatar?: string;
  screenshot?: string;
  sort?: number; // 排序值，范围1-10，默认5
  type_id: number; // 友链类型ID（必选）
  rss_url?: string; // RSS订阅地址
}

// 更新友链请求
export interface UpdateFriendRequest {
  name?: string;
  url?: string;
  description?: string;
  avatar?: string;
  screenshot?: string;
  sort?: number; // 排序值，范围1-10
  type_id?: number | null;
  is_invalid?: boolean; // 是否失效
  is_pending?: boolean; // 是否为待审核申请
  rss_url?: string; // RSS订阅地址
  accessible?: number; // 可访问性状态
}
