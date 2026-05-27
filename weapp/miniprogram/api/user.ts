/**
 * 用户相关 API 接口
 * 对接后端实际存在的用户接口
 */

import type { UserInfo } from '../types';
import { get, post } from '../utils/request';

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
  wechat_code?: string;
}

/**
 * 用户登录
 */
export function login(params: LoginParams): Promise<LoginResponse> {
  return post<LoginResponse>('/auth/login', params);
}

/**
 * 微信小程序登录
 */
export function wechatLogin(code: string): Promise<LoginResponse> {
  return post<LoginResponse>('/auth/wechat', { code });
}

/**
 * 确认微信扫码登录授权
 */
export function confirmWechatLogin(scene: string, code: string): Promise<void> {
  return post<void>('/auth/wechat/confirm', { scene, code });
}

/**
 * 退出登录
 */
export function logout(): Promise<void> {
  return post<void>('/auth/logout', {}, { needAuth: true });
}

/**
 * 获取用户信息
 */
export function getUserProfile(): Promise<UserInfo> {
  return get<UserInfo>('/user/profile', {}, { needAuth: true });
}
