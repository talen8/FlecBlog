import { ref, computed } from 'vue';
import type { ApiResponse } from '@@/types/request';

const ACCESS_TOKEN_KEY = 'access_token';

/**
 * 获取存储的 access token
 * @returns {string | null} access token字符串或null
 */
const getStoredToken = (): string | null => {
  if (import.meta.server) return null;
  return localStorage.getItem(ACCESS_TOKEN_KEY);
};

export const accessToken = ref<string | null>(getStoredToken());

// 响应式登录状态
export const isLoggedIn = computed(() => !!accessToken.value && accessToken.value !== '');

/**
 * 设置 access token
 */
export const setAccessToken = (access: string): void => {
  if (import.meta.client) {
    localStorage.setItem(ACCESS_TOKEN_KEY, access);
  }
  accessToken.value = access;
};

/**
 * 获取 access token
 */
export const getAccessToken = (): string | null => {
  if (import.meta.client) {
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  }
  return accessToken.value;
};

/**
 * 清除 access token
 */
export const clearAccessToken = (): void => {
  if (import.meta.client) {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
  }
  accessToken.value = null;
};

/**
 * 登出操作
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
 * 获取响应式的登录状态
 */
export const useAuth = () => isLoggedIn;
