import { ref } from 'vue';
import type { User } from '@/types/user';

const currentUser = ref<User | null>(null);
let userInfoPromise: Promise<User | null> | null = null;
let userInfoRequestId = 0;
let redirectingToLogin = false;

const ACCESS_TOKEN_KEY = 'access_token';

/**
 * 获取内存中的 access token
 * @returns {string | null} access token字符串或null
 */
export const getAccessToken = (): string | null => {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
};

/**
 * 设置 access token
 * @param {string} token access token字符串
 */
export const setAccessToken = (token: string): void => {
  localStorage.setItem(ACCESS_TOKEN_KEY, token);
};

/**
 * 清除 access token
 */
export const clearAccessToken = (): void => {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
};

/**
 * 检查用户是否已经登录
 * @returns {boolean} true表示已登录，false表示未登录
 */
export const checkAuth = (): boolean => {
  const token = getAccessToken();
  return token !== null && token !== '';
};

/**
 * 清除本地认证状态
 */
export const clearAuthState = (): void => {
  clearAccessToken();
  clearUserInfo();
};

/**
 * 设置当前登录用户信息
 */
export const setUserInfo = (user: User | null): void => {
  currentUser.value = user;
};

/**
 * 清空当前登录用户信息
 */
export const clearUserInfo = (): void => {
  currentUser.value = null;
  userInfoRequestId += 1;
  userInfoPromise = null;
};

/**
 * 获取当前登录用户信息
 */
export const getUserInfo = (): User | null => {
  return currentUser.value;
};

/**
 * 从服务端拉取当前登录用户信息
 */
export const fetchUserInfo = async (): Promise<User | null> => {
  if (!checkAuth()) {
    clearUserInfo();
    return null;
  }

  if (!userInfoPromise) {
    const requestId = ++userInfoRequestId;

    userInfoPromise = import('@/api/user')
      .then(({ getProfile }) => getProfile())
      .then(user => {
        if (requestId !== userInfoRequestId || !checkAuth()) {
          return currentUser.value;
        }

        setUserInfo(user);
        return user;
      })
      .catch(error => {
        if (requestId === userInfoRequestId) {
          clearUserInfo();
        }
        throw error;
      })
      .finally(() => {
        if (requestId === userInfoRequestId) {
          userInfoPromise = null;
        }
      });
  }

  return userInfoPromise;
};

/**
 * 确保当前登录用户信息已初始化
 */
export const ensureUserInfo = async (): Promise<User | null> => {
  if (currentUser.value) return currentUser.value;
  return fetchUserInfo();
};

/**
 * 获取当前登录用户角色
 */
export const getCurrentUserRole = (): string => {
  return currentUser.value?.role || '';
};

/**
 * 当前用户是否为超级管理员
 */
export const isSuperAdmin = (): boolean => {
  return currentUser.value?.role === 'super_admin';
};

/**
 * 跳转到登录页并清理认证状态
 */
export const redirectToLogin = (): void => {
  clearAuthState();

  if (window.location.pathname === '/login') {
    return;
  }

  if (redirectingToLogin) return;

  redirectingToLogin = true;
  window.location.replace('/login');
};

/**
 * 注销用户，清除token和用户信息
 */
export const logout = (): void => {
  clearAuthState();
};
