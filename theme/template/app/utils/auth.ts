/**
 * 认证工具函数
 *
 * 管理用户认证状态，包括 Access Token 的存储、获取、清除等操作。
 * Token 存储在 localStorage 中，并在 Vue 响应式状态中同步。
 *
 * @module utils/auth
 */

import { ref, computed } from 'vue';
import type { ApiResponse } from '@@/types';

const ACCESS_TOKEN_KEY = 'access_token';

/**
 * 从 localStorage 获取存储的 Token
 *
 * @returns 存储的 Token 字符串，服务端返回 null
 */
const getStoredToken = (): string | null => {
  if (import.meta.server) return null;
  return localStorage.getItem(ACCESS_TOKEN_KEY);
};

/**
 * 当前 Access Token（响应式）
 */
export const accessToken = ref<string | null>(getStoredToken());

/**
 * 登录状态（计算属性）
 *
 * 根据是否存在有效的 Access Token 判断用户是否已登录。
 */
export const isLoggedIn = computed(() => !!accessToken.value && accessToken.value !== '');

/**
 * 设置 Access Token
 *
 * 将 Token 存储到 localStorage 并更新响应式状态。
 *
 * @param access - Access Token 字符串
 */
export const setAccessToken = (access: string): void => {
  if (import.meta.client) {
    localStorage.setItem(ACCESS_TOKEN_KEY, access);
  }
  accessToken.value = access;
};

/**
 * 获取 Access Token
 *
 * 优先从 localStorage 获取，服务端从响应式状态获取。
 *
 * @returns Access Token 字符串或 null
 */
export const getAccessToken = (): string | null => {
  if (import.meta.client) {
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  }
  return accessToken.value;
};

/**
 * 清除 Access Token
 *
 * 从 localStorage 和响应式状态中清除 Token。
 */
export const clearAccessToken = (): void => {
  if (import.meta.client) {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
  }
  accessToken.value = null;
};

/**
 * 用户登出
 *
 * 清除本地 Token 并通知后端登出。
 * 即使后端请求失败也会清除本地状态。
 */
export const logout = (): void => {
  const token = accessToken.value;
  clearAccessToken();

  if (token) {
    const config = useRuntimeConfig();
    $fetch<ApiResponse<void>>('/auth/logout', {
      method: 'POST',
      baseURL: config.public.apiUrl as string,
      headers: {
        Authorization: `Bearer ${token}`,
      },
      credentials: 'include',
    }).catch(() => {});
  }
};

/**
 * 认证状态 Composable
 *
 * 返回登录状态的响应式引用。
 *
 * @returns 登录状态的 ComputedRef
 *
 * @example
 * ```ts
 * const isLoggedIn = useAuth();
 * if (isLoggedIn.value) {
 *   // 用户已登录
 * }
 * ```
 */
export const useAuth = () => isLoggedIn;
