/**
 * HTTP 请求封装
 *
 * 基于 ofetch 提供统一的 HTTP 请求方法，支持：
 * - 自动添加 Authorization 头
 * - Token 过期自动刷新
 * - 统一错误处理
 *
 * @module utils/request
 */

import { $fetch } from 'ofetch';
import type { FetchOptions } from 'ofetch';
import { accessToken, setAccessToken, logout } from './auth';
import type { ApiResponse } from '@@/types';

/** HTTP 方法类型 */
type HttpMethod =
  | 'GET'
  | 'HEAD'
  | 'PATCH'
  | 'POST'
  | 'PUT'
  | 'DELETE'
  | 'CONNECT'
  | 'OPTIONS'
  | 'TRACE';

/** 获取 API 基础 URL */
const getBaseURL = () => useRuntimeConfig().public.apiUrl as string;

/** Token 刷新状态 */
let isRefreshing = false;
/** Token 刷新 Promise */
let refreshPromise: Promise<boolean> | null = null;

/**
 * 执行 Token 刷新
 *
 * @returns 刷新是否成功
 */
const doRefreshToken = async (): Promise<boolean> => {
  try {
    const res = await $fetch<ApiResponse<{ access_token: string }>>('/auth/refresh', {
      method: 'POST',
      baseURL: getBaseURL(),
      credentials: 'include',
    });
    if (res.code === 0 && res.data) {
      setAccessToken(res.data.access_token);
      return true;
    }
    return false;
  } catch {
    return false;
  }
};

/**
 * 通用 API 请求方法
 *
 * 所有 API 请求的基础方法，处理认证、Token 刷新和错误。
 *
 * @template T - 响应数据类型
 * @param url - 请求 URL
 * @param options - 请求选项
 * @returns 响应数据
 * @throws 请求错误
 *
 * @example
 * ```ts
 * const data = await apiRequest<UserInfo>('/user/profile');
 * ```
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，提供灵活性
export async function apiRequest<T = any>(
  url: string,
  options: Omit<FetchOptions, 'method'> & { method?: HttpMethod; _retry?: boolean } = {}
): Promise<T> {
  const config = useRuntimeConfig();
  const headers: Record<string, string> = {
    ...((options.headers as Record<string, string>) || {}),
  };

  if (accessToken.value && url !== '/auth/refresh') {
    headers['Authorization'] = `Bearer ${accessToken.value}`;
  }

  try {
    return await $fetch<T>(url, {
      ...options,
      baseURL: config.public.apiUrl,
      headers,
      credentials: 'include',
      // eslint-disable-next-line @typescript-eslint/no-explicit-any -- ofetch 类型兼容
    } as any);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 错误处理，需要访问 response.status
  } catch (error: any) {
    if (error?.response?.status === 401 && !options._retry) {
      if (!isRefreshing) {
        isRefreshing = true;
        refreshPromise = doRefreshToken().finally(() => {
          isRefreshing = false;
        });
      }

      const success = await refreshPromise;
      if (success) {
        return apiRequest<T>(url, { ...options, _retry: true });
      }
      logout();
    }
    throw error;
  }
}

/**
 * GET 请求
 *
 * @template T - 响应数据类型
 * @param url - 请求 URL
 * @param options - 请求选项
 * @returns 响应数据
 *
 * @example
 * ```ts
 * const users = await get<User[]>('/users');
 * const user = await get<User>('/users/1');
 * ```
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，提供灵活性
export async function get<T = any>(
  url: string,
  options: Omit<FetchOptions, 'method'> = {}
): Promise<T> {
  return await apiRequest<T>(url, { ...options, method: 'GET' });
}

/**
 * POST 请求
 *
 * @template T - 响应数据类型
 * @param url - 请求 URL
 * @param body - 请求体
 * @param options - 请求选项
 * @returns 响应数据
 *
 * @example
 * ```ts
 * const result = await post<ApiResponse<User>>('/users', { name: 'John' });
 * ```
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，提供灵活性
export async function post<T = any>(
  url: string,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 请求体类型不确定
  body?: any,
  options: Omit<FetchOptions, 'method'> = {}
): Promise<T> {
  return await apiRequest<T>(url, { ...options, method: 'POST', body });
}

/**
 * PUT 请求（全量更新）
 *
 * @template T - 响应数据类型
 * @param url - 请求 URL
 * @param body - 请求体
 * @param options - 请求选项
 * @returns 响应数据
 *
 * @example
 * ```ts
 * const result = await put<User>('/users/1', { name: 'John', email: 'john@example.com' });
 * ```
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，提供灵活性
export async function put<T = any>(
  url: string,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 请求体类型不确定
  body?: any,
  options: Omit<FetchOptions, 'method'> = {}
): Promise<T> {
  return await apiRequest<T>(url, { ...options, method: 'PUT', body });
}

/**
 * PATCH 请求（部分更新）
 *
 * @template T - 响应数据类型
 * @param url - 请求 URL
 * @param body - 请求体
 * @param options - 请求选项
 * @returns 响应数据
 *
 * @example
 * ```ts
 * const result = await patch<User>('/users/1', { name: 'John' });
 * ```
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，提供灵活性
export async function patch<T = any>(
  url: string,
  // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 请求体类型不确定
  body?: any,
  options: Omit<FetchOptions, 'method'> = {}
): Promise<T> {
  return await apiRequest<T>(url, { ...options, method: 'PATCH', body });
}

/**
 * DELETE 请求
 *
 * @template T - 响应数据类型
 * @param url - 请求 URL
 * @param options - 请求选项
 * @returns 响应数据
 *
 * @example
 * ```ts
 * await del('/users/1');
 * ```
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，提供灵活性
export async function del<T = any>(
  url: string,
  options: Omit<FetchOptions, 'method'> = {}
): Promise<T> {
  return await apiRequest<T>(url, { ...options, method: 'DELETE' });
}
