import type {
  LoginParams,
  LoginResponse,
  RegisterParams,
  RegisterResponse,
  ForgotPasswordParams,
  ResetPasswordParams,
  RefreshTokenResponse,
} from '../../../types/auth';
import { createApi } from './createApi';

const authApi = createApi<LoginResponse>('/auth');

/** 用户登录 */
export const login = async (data: LoginParams) => {
  return authApi.post<LoginResponse>('/login', data);
};

/** 用户注册 */
export const register = async (data: RegisterParams) => {
  return authApi.post<RegisterResponse>('/register', data);
};

/** 刷新Token */
export const refreshToken = async () => {
  return authApi.post<RefreshTokenResponse>('/refresh');
};

/** 忘记密码 */
export const forgotPassword = async (data: ForgotPasswordParams) => {
  await authApi.post('/forgot-password', data);
};

/** 重置密码 */
export const resetPassword = async (data: ResetPasswordParams) => {
  await authApi.post('/reset-password', data);
};

/** 获取微信登录二维码 */
export const getWechatQrcode = async (): Promise<{ blob: Blob; scene: string }> => {
  const apiUrl = useRuntimeConfig().public.apiUrl as string;
  const resp = await $fetch.raw(`${apiUrl}/auth/wechat/qrcode`, { responseType: 'blob' });
  return { blob: resp._data as Blob, scene: resp.headers.get('X-Scene') || '' };
};

/** 查询微信扫码状态 */
export const getWechatScene = async (scene: string) => {
  return authApi.get<{ status: string; access_token?: string }>(`/wechat/scene/${scene}`);
};

/** 登出（通知后端清除 session） */
export const logout = async () => {
  await authApi.post('/logout');
};

/** 获取 OAuth 登录跳转地址 */
export const getOAuthUrl = (provider: string, redirect: string) => {
  const config = useRuntimeConfig();
  return `${config.public.apiUrl}/auth/${provider}?redirect=${encodeURIComponent(redirect)}`;
};

/** 获取 OAuth 绑定跳转地址 */
export const getOAuthBindUrl = (provider: string, redirect: string, token: string) => {
  const config = useRuntimeConfig();
  return `${config.public.apiUrl}/auth/${provider}?action=bind&token=${token}&redirect=${encodeURIComponent(redirect)}`;
};
