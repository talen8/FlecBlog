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
