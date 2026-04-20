import { ref, computed } from 'vue';

// Token 状态
export const accessToken = ref<string | null>(null);
export const refreshToken = ref<string | null>(null);

// 响应式登录状态
export const isLoggedIn = computed(() => !!accessToken.value && accessToken.value !== '');

// 从 localStorage 初始化 token（仅客户端）
if (import.meta.client) {
  accessToken.value = localStorage.getItem('access_token');
  refreshToken.value = localStorage.getItem('refresh_token');
}

/**
 * 设置双token
 */
export const setTokens = (access: string, refresh: string): void => {
  accessToken.value = access;
  refreshToken.value = refresh;

  // 同步到 localStorage（仅客户端）
  if (import.meta.client) {
    localStorage.setItem('access_token', access);
    localStorage.setItem('refresh_token', refresh);
  }
};

/**
 * 设置access token（用于token刷新）
 */
export const setAccessToken = (access: string): void => {
  accessToken.value = access;

  // 同步到 localStorage（仅客户端）
  if (import.meta.client) {
    localStorage.setItem('access_token', access);
  }
};

/**
 * 登出操作（清除双token）
 */
export const logout = (): void => {
  accessToken.value = null;
  refreshToken.value = null;
  if (import.meta.client) {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  }
};

/**
 * 获取响应式的登录状态
 */
export const useAuth = () => isLoggedIn;
