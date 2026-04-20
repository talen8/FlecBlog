/**
 * 用户角色枚举
 */
export type UserRole = 'super_admin' | 'admin' | 'user';

/**
 * 用户基本信息
 */
export interface UserInfo {
  id: number;
  email: string;
  email_hash: string;
  is_virtual_email: boolean; // 是否为虚拟邮箱（需绑定真实邮箱）
  avatar?: string;
  badge?: string;
  nickname: string;
  website?: string;
  last_login?: string;
  created_at: string;
  role: UserRole;
  has_password: boolean;
  linked_oauths: string[];
}

/**
 * 用户资料更新参数（所有字段均为可选）
 */
export interface UpdateProfileParams {
  nickname?: string;
  email?: string;
  avatar?: string;
  badge?: string;
  website?: string;
}

/**
 * 登录请求参数
 */
export interface LoginParams {
  email: string;
  password: string;
}

/**
 * 登录响应数据
 */
export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: UserInfo;
}

/**
 * 注册请求参数
 */
export interface RegisterParams {
  email: string;
  nickname: string;
  password: string;
  website?: string;
}

/**
 * 注册响应数据
 */
export interface RegisterResponse {
  access_token: string;
  refresh_token: string;
  user: UserInfo;
}

/**
 * 忘记密码请求参数
 */
export interface ForgotPasswordParams {
  email: string;
}

/**
 * 重置密码请求参数
 */
export interface ResetPasswordParams {
  email: string;
  code: string;
  password: string;
}

/**
 * 修改密码请求参数
 */
export interface ChangePasswordParams {
  old_password: string;
  new_password: string;
}

/**
 * 注销账户请求参数
 */
export interface DeactivateAccountParams {
  password: string;
}

/**
 * 刷新Token请求参数
 */
export interface RefreshTokenParams {
  refresh_token: string;
}

/**
 * 刷新Token响应数据
 */
export interface RefreshTokenResponse {
  access_token: string;
  refresh_token: string;
}
