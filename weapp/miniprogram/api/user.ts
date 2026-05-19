/**
 * 用户相关 API 接口
 * 对接后端实际存在的用户接口
 */

import type { UserInfo } from '../types';
import { get, post, patch, put, del } from '../utils/request';

/**
 * 用户登录响应
 */
export interface LoginResponse {
  access_token: string;
  user: UserInfo;
}

/**
 * 登录请求参数
 */
export interface LoginParams {
  email: string;
  password: string;
}

/**
 * 更新用户信息参数
 */
export interface UpdateUserParams {
  nickname?: string;
  email?: string;
  avatar?: string;
  badge?: string;
  website?: string;
}

/**
 * 修改密码参数
 */
export interface ChangePasswordParams {
  old_password: string;
  new_password: string;
}

/**
 * 用户登录
 * @param params - 登录参数
 * @returns 登录响应
 */
export function login(params: LoginParams): Promise<LoginResponse> {
  return post<LoginResponse>('/auth/login', params);
}

/**
 * 退出登录
 * @returns 操作结果
 */
export function logout(): Promise<void> {
  return post<void>('/auth/logout', {}, { needAuth: true });
}

/**
 * 刷新 Token
 * @returns 新的 access_token
 */
export function refreshToken(): Promise<{ access_token: string }> {
  return post<{ access_token: string }>('/auth/refresh');
}

/**
 * 获取用户信息
 * @returns 用户信息
 */
export function getUserProfile(): Promise<UserInfo> {
  return get<UserInfo>('/user/profile', {}, { needAuth: true });
}

/**
 * 更新用户信息
 * @param params - 用户信息
 * @returns 更新后的用户信息
 */
export function updateUserProfile(params: UpdateUserParams): Promise<UserInfo> {
  return patch<UserInfo>('/user/profile', params, { needAuth: true });
}

/**
 * 修改密码
 * @param params - 密码参数
 * @returns 操作结果
 */
export function changePassword(params: ChangePasswordParams): Promise<void> {
  return put<void>('/user/password', params, { needAuth: true });
}

/**
 * 设置密码（OAuth 用户首次设置）
 * @param params - 密码参数
 * @returns 操作结果
 */
export function setPassword(params: { password: string; confirm_password: string }): Promise<void> {
  return post<void>('/user/password', params, { needAuth: true });
}

/**
 * 注销账户
 * @param params - 密码参数
 * @returns 操作结果
 */
export function deactivateAccount(params: { password: string }): Promise<void> {
  return del<void>('/user/deactivate', { needAuth: true, data: params });
}

/**
 * 解绑第三方账号
 * @param provider - 提供商名称
 * @returns 操作结果
 */
export function unbindOAuth(provider: string): Promise<void> {
  return del<void>(`/user/oauth/${provider}`, { needAuth: true });
}
