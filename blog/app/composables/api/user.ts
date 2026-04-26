import type {
  LoginParams,
  LoginResponse,
  RegisterParams,
  RegisterResponse,
  UserInfo,
  UpdateProfileParams,
  ForgotPasswordParams,
  ResetPasswordParams,
  ChangePasswordParams,
  DeactivateAccountParams,
  RefreshTokenResponse,
} from '@@/types/user';
import { createApi } from './createApi';

const authApi = createApi<LoginResponse>('/auth');
const userApi = createApi<UserInfo>('');

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

/** 获取当前用户信息 */
export const getUserProfile = async () => {
  return userApi.get<UserInfo>('/user/profile');
};

/** 更新用户资料 */
export const updateUserProfile = async (data: UpdateProfileParams) => {
  return userApi.patchRequest<UserInfo>('/user/profile', data);
};

/** 忘记密码 */
export const forgotPassword = async (data: ForgotPasswordParams) => {
  await authApi.post('/forgot-password', data);
};

/** 重置密码 */
export const resetPassword = async (data: ResetPasswordParams) => {
  await authApi.post('/reset-password', data);
};

/** 修改密码 */
export const changePassword = async (data: ChangePasswordParams) => {
  await userApi.put('/user/password', data);
};

/** 设置密码（OAuth 用户首次设置密码） */
export const setPassword = async (data: { password: string; confirm_password: string }) => {
  await userApi.post('/user/password', data);
};

/** 注销账户 */
export const deactivateAccount = async (data: DeactivateAccountParams) => {
  await userApi.deleteRequest('/user/deactivate', data);
};

/** 解绑第三方账号 */
export const unbindOAuth = async (provider: string) => {
  await userApi.deleteRequest(`/user/oauth/${provider}`);
};
