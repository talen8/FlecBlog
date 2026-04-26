import request from '@/utils/request';
import type {
  LoginParams,
  LoginResponse,
  User,
  UserListData,
  UserListQuery,
  ResetPasswordRequest,
  CreateUserRequest,
  UpdateUserRequest,
  RefreshTokenResponse,
} from '@/types/user';

/**
 * 用户登录
 * @param data 登录参数
 * @returns Promise<LoginResponse>
 */
export function login(data: LoginParams): Promise<LoginResponse> {
  return request.post('/auth/login', data);
}

/**
 * 获取当前登录用户信息
 * @returns Promise<User>
 */
export function getProfile(): Promise<User> {
  return request.get('/user/profile');
}

/**
 * 刷新Token
 * @returns Promise<RefreshTokenResponse>
 */
export function refreshToken(): Promise<RefreshTokenResponse> {
  return request.post('/auth/refresh');
}

/**
 * 用户登出
 * @returns Promise<void>
 */
export function logout(): Promise<void> {
  return request.post('/auth/logout');
}

/**
 * 获取用户列表
 * @param params 查询参数
 * @returns Promise<UserListData>
 */
export function getUsers(params: UserListQuery): Promise<UserListData> {
  return request.get('/admin/users', { params });
}

/**
 * 获取用户详情
 * @param id 用户ID
 * @returns Promise<User>
 */
export function getUserById(id: number): Promise<User> {
  return request.get(`/admin/users/${id}`);
}

/**
 * 删除用户
 * @param id 用户ID
 * @returns Promise<void>
 */
export function deleteUser(id: number): Promise<void> {
  return request.delete(`/admin/users/${id}`);
}

/**
 * 管理员重置用户密码
 * @param id 用户ID
 * @param data 重置密码数据
 * @returns Promise<void>
 */
export function resetUserPassword(id: number, data: ResetPasswordRequest): Promise<void> {
  return request.put(`/admin/users/${id}/password`, data);
}

/**
 * 创建用户
 * @param data 用户数据
 * @returns Promise<User>
 */
export function createUser(data: CreateUserRequest): Promise<User> {
  return request.post('/admin/users', data);
}

/**
 * 更新用户
 * @param id 用户ID
 * @param data 更新数据
 * @returns Promise<User>
 */
export function updateUser(id: number, data: UpdateUserRequest): Promise<User> {
  return request.put(`/admin/users/${id}`, data);
}
