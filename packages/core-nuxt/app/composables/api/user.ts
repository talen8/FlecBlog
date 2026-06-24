import type {
  UserInfo,
  UpdateProfileParams,
  ChangePasswordParams,
  DeactivateAccountParams,
} from '../../../types/user';
import { createApi } from './createApi';

const userApi = createApi<UserInfo>('/user');

/** 获取当前用户信息 */
export const getUserProfile = async () => {
  return userApi.get<UserInfo>('/profile');
};

/** 更新用户资料 */
export const updateUserProfile = async (data: UpdateProfileParams) => {
  return userApi.patchRequest<UserInfo>('/profile', data);
};

/** 修改密码 */
export const changePassword = async (data: ChangePasswordParams) => {
  await userApi.put('/password', data);
};

/** 设置密码（OAuth 用户首次设置密码） */
export const setPassword = async (data: { password: string; confirm_password: string }) => {
  await userApi.post('/password', data);
};

/** 注销账户 */
export const deactivateAccount = async (data: DeactivateAccountParams) => {
  await userApi.deleteRequest('/deactivate', data);
};

/** 解绑第三方账号 */
export const unbindOAuth = async (provider: string) => {
  await userApi.deleteRequest(`/oauth/${provider}`);
};
