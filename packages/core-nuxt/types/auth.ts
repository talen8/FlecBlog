import type { UserInfo } from './user';

/** 登录请求参数 */
export interface LoginParams {
  email: string;
  password: string;
}

/** 登录响应数据 */
export interface LoginResponse {
  access_token: string;
  user: UserInfo;
}

/** 注册请求参数 */
export interface RegisterParams {
  email: string;
  nickname: string;
  password: string;
  website?: string;
}

/** 注册响应数据 */
export interface RegisterResponse {
  access_token: string;
  user: UserInfo;
}

/** 忘记密码请求参数 */
export interface ForgotPasswordParams {
  email: string;
}

/** 重置密码请求参数 */
export interface ResetPasswordParams {
  email: string;
  code: string;
  password: string;
}

/** 刷新Token响应数据 */
export interface RefreshTokenResponse {
  access_token: string;
}
