// 用户实体
export interface User {
  id: number;
  email: string;
  nickname: string;
  avatar: string;
  badge?: string;
  website?: string;
  role: string; // super_admin | admin | user | guest
  is_enabled: boolean; // 是否启用
  last_login: string;
  deleted_at?: string;
  has_password: boolean; // 是否设置了密码
  github_id: string; // GitHub ID
  google_id: string; // Google ID
  qq_id: string; // QQ ID
  microsoft_id: string; // Microsoft ID
  feishu_open_id: string; // 飞书 OpenID
}

// 登录请求
export interface LoginParams {
  email: string;
  password: string;
}

// 登录响应
export interface LoginResponse {
  access_token: string;
  user?: User;
}

// 重置密码请求
export interface ResetPasswordRequest {
  new_password: string;
}

// 创建用户请求
export interface CreateUserRequest {
  password: string;
  email: string;
  nickname: string;
  avatar?: string;
  badge?: string;
  website?: string;
  role: 'super_admin' | 'admin' | 'user' | 'guest';
}

// 更新用户请求
export interface UpdateUserRequest {
  password?: string;
  email?: string;
  nickname?: string;
  avatar?: string;
  badge?: string;
  website?: string;
  role?: 'super_admin' | 'admin' | 'user' | 'guest';
  is_enabled?: boolean;
}

// 用户列表查询参数
export interface UserListQuery {
  page: number;
  page_size: number;
  keyword?: string; // 搜索关键词（邮箱、昵称）
  role?: string; // 角色筛选
  is_enabled?: boolean; // 状态筛选
  is_deleted?: boolean; // 是否已删除
  login_method?: string; // 登录方式筛选
  last_login_start?: string; // 最后登录开始时间
  last_login_end?: string; // 最后登录结束时间
  start_time?: string; // 注册开始时间
  end_time?: string; // 注册结束时间
}

// 分页数据
export interface UserListData {
  list: User[];
  total: number;
  page: number;
  page_size: number;
}

// 刷新Token响应
export interface RefreshTokenResponse {
  access_token: string;
}
