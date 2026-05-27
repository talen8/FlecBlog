/**
 * API 请求封装
 * 基于 wx.request 封装为 Promise API，支持请求/响应拦截器
 */

import type { ApiResponse, ApiError, RequestConfig, RequestInterceptor, ErrorInterceptor } from '../types';
import { APP_CONFIG } from '../config';

// ==================== 拦截器管理器 ====================

/** 请求拦截器队列 */
const requestInterceptors: RequestInterceptor[] = [];

/** 错误拦截器队列 */
const errorInterceptors: ErrorInterceptor[] = [];

/**
 * 添加请求拦截器
 * @param interceptor - 请求拦截器函数
 * @returns 移除该拦截器的函数
 */
export function addRequestInterceptor(interceptor: RequestInterceptor): () => void {
  requestInterceptors.push(interceptor);
  return () => {
    const index = requestInterceptors.indexOf(interceptor);
    if (index > -1) requestInterceptors.splice(index, 1);
  };
}

/**
 * 添加错误拦截器
 * @param interceptor - 错误拦截器函数
 * @returns 移除该拦截器的函数
 */
export function addErrorInterceptor(interceptor: ErrorInterceptor): () => void {
  errorInterceptors.push(interceptor);
  return () => {
    const index = errorInterceptors.indexOf(interceptor);
    if (index > -1) errorInterceptors.splice(index, 1);
  };
}

/**
 * 执行错误拦截器链
 */
async function runErrorInterceptors(error: ApiError, config: RequestConfig): Promise<never> {
  for (const interceptor of errorInterceptors) {
    await interceptor(error, config);
  }
  throw new Error(error.message);
}

// ==================== Token 管理 ====================

/** Token 刷新状态锁 */
let isRefreshing = false;
/** 刷新 Token 的 Promise，用于并发请求共享 */
let refreshPromise: Promise<boolean> | null = null;

/**
 * 获取当前 Token
 */
export function getToken(): string | null {
  return wx.getStorageSync('access_token');
}

/**
 * 设置 Token
 */
export function setToken(token: string): void {
  wx.setStorageSync('access_token', token);
}

/**
 * 清除 Token
 */
export function clearToken(): void {
  wx.removeStorageSync('access_token');
}

/**
 * 清除登录状态
 */
export function clearAuth(): void {
  clearToken();
  wx.removeStorageSync('user_info');
}

/**
 * 执行 Token 刷新
 * @returns 是否刷新成功
 */
async function doRefreshToken(): Promise<boolean> {
  try {
    const token = getToken();
    if (!token) return false;

    const res = await new Promise<{ code: number; data?: { access_token: string } }>((resolve, reject) => {
      wx.request({
        url: `${APP_CONFIG.API_URL}/auth/refresh`,
        method: 'POST',
        header: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        timeout: 30000,
        success: (res) => resolve(res.data as { code: number; data?: { access_token: string } }),
        fail: (err) => reject(err),
      });
    });

    if (res.code === 0 && res.data?.access_token) {
      setToken(res.data.access_token);
      return true;
    }
    return false;
  } catch {
    return false;
  }
}

/**
 * 登出操作
 */
export function logout(): void {
  const token = getToken();
  clearAuth();

  // 尝试调用后端登出接口（静默失败）
  if (token) {
    wx.request({
      url: `${APP_CONFIG.API_URL}/auth/logout`,
      method: 'POST',
      header: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
    });
  }
}

// ==================== URL 构建 ====================

/**
 * 构建完整请求 URL
 * @param url - 请求路径
 * @param params - 查询参数
 * @returns 完整 URL
 */
function buildUrl(url: string, params?: Record<string, any>): string {
  let fullUrl = url.startsWith('http') ? url : `${APP_CONFIG.API_URL}${url}`;

  if (params && Object.keys(params).length > 0) {
    const queryString = Object.entries(params)
      .filter(([, value]) => value !== undefined && value !== null)
      .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(String(value))}`)
      .join('&');

    if (queryString) {
      fullUrl += `${fullUrl.includes('?') ? '&' : '?'}${queryString}`;
    }
  }

  return fullUrl;
}

// ==================== 请求发送 ====================

/**
 * 核心请求方法
 */
export function request<T>(config: RequestConfig): Promise<T> {
  const {
    url,
    method = 'GET',
    data,
    params,
    header,
    loading = false,
    loadingText = '加载中...',
    needAuth = false,
  } = config;

  // 1. 显示加载提示
  if (loading) {
    wx.showLoading({ title: loadingText, mask: true });
  }

  // 2. 构建最终配置
  const finalConfig: RequestConfig = {
    url,
    method,
    data,
    params,
    header: {
      'Content-Type': 'application/json',
      ...header,
    },
    loading,
    loadingText,
    needAuth,
  };

  // 3. 执行请求拦截器链
  return (async () => {
    let processedConfig = finalConfig;
    for (const interceptor of requestInterceptors) {
      processedConfig = await interceptor(processedConfig);
    }

    // 4. 添加认证 Token
    if (processedConfig.needAuth) {
      const token = getToken();
      if (token) {
        processedConfig.header = {
          ...processedConfig.header,
          Authorization: `Bearer ${token}`,
        };
      }
    }

    // 5. 发送请求
    return new Promise<T>((resolve, reject) => {
      wx.request({
        url: buildUrl(processedConfig.url, processedConfig.params),
        method: processedConfig.method as WechatMiniprogram.RequestOption['method'],
        data: processedConfig.data,
        header: processedConfig.header,
        timeout: 30000,
        success: (res) => {
          if (loading) wx.hideLoading();

          const { statusCode, data: responseData } = res;

          // 6. 处理 HTTP 状态码
          if (statusCode >= 200 && statusCode < 300) {
            const apiResponse = responseData as ApiResponse<T>;

            if (apiResponse.code === 0) {
              resolve(apiResponse.data as T);
            } else {
              const apiError: ApiError = {
                code: apiResponse.code,
                message: apiResponse.message || '请求失败，请稍后重试',
                url: processedConfig.url,
                statusCode,
              };

              // 处理 401 未授权 - 尝试刷新 Token
              if (apiResponse.code === 401 && needAuth) {
                handleUnauthorized(processedConfig, resolve, reject);
                return;
              }

              runErrorInterceptors(apiError, processedConfig).catch(() => {
                reject(new Error(apiError.message));
              });
            }
          } else if (statusCode === 401) {
            handleUnauthorized(processedConfig, resolve, reject);
          } else if (statusCode >= 500) {
            const apiError: ApiError = {
              code: statusCode,
              message: '服务器繁忙，请稍后重试',
              url: processedConfig.url,
              statusCode,
            };
            runErrorInterceptors(apiError, processedConfig).catch(() => {
              reject(new Error('服务器繁忙，请稍后重试'));
            });
          } else {
            const apiError: ApiError = {
              code: statusCode,
              message: '请求失败，请稍后重试',
              url: processedConfig.url,
              statusCode,
            };
            runErrorInterceptors(apiError, processedConfig).catch(() => {
              reject(new Error('请求失败，请稍后重试'));
            });
          }
        },
        fail: (err) => {
          if (loading) wx.hideLoading();

          let errorMsg: string = '网络连接失败，请检查网络设置';
          if (err.errMsg && err.errMsg.includes('timeout')) {
            errorMsg = '请求超时，请稍后重试';
          }

          const apiError: ApiError = {
            code: -1,
            message: errorMsg,
            url: processedConfig.url,
          };

          runErrorInterceptors(apiError, processedConfig).catch(() => {
            reject(new Error(errorMsg));
          });
        },
      });
    });
  })();
}

/**
 * 处理 401 未授权
 * 参考 blog 端实现，支持 Token 自动刷新
 */
async function handleUnauthorized<T>(
  config: RequestConfig,
  resolve: (value: T) => void,
  reject: (reason?: unknown) => void
): Promise<void> {
  // 如果是刷新 Token 的请求本身失败，直接登出
  if (config.url === '/auth/refresh') {
    handleLogout();
    reject(new Error('登录已过期，请重新登录'));
    return;
  }

  // 避免并发刷新：如果已经在刷新，等待刷新结果
  if (isRefreshing && refreshPromise) {
    const success = await refreshPromise;
    if (success) {
      // 刷新成功，重试当前请求
      request<T>(config).then(resolve).catch(reject);
    } else {
      reject(new Error('登录已过期，请重新登录'));
    }
    return;
  }

  // 开始刷新 Token
  isRefreshing = true;
  refreshPromise = doRefreshToken().finally(() => {
    isRefreshing = false;
    refreshPromise = null;
  });

  const success = await refreshPromise;
  if (success) {
    // 刷新成功，重试当前请求
    request<T>(config).then(resolve).catch(reject);
  } else {
    // 刷新失败，登出
    handleLogout();
    reject(new Error('登录已过期，请重新登录'));
  }
}

/**
 * 处理登出逻辑
 */
function handleLogout(): void {
  clearAuth();
  wx.showToast({ title: '登录已过期，请重新登录', icon: 'none' });
  wx.navigateTo({ url: '/pages/login/login' });
}

// ==================== 便捷请求方法 ====================

/** GET 请求 */
export function get<T>(
  url: string,
  params?: Record<string, any>,
  config?: Omit<RequestConfig, 'url' | 'method' | 'params'>
): Promise<T> {
  return request<T>({ url, method: 'GET', params, ...config });
}

/** POST 请求 */
export function post<T>(
  url: string,
  data?: Record<string, any>,
  config?: Omit<RequestConfig, 'url' | 'method' | 'data'>
): Promise<T> {
  return request<T>({ url, method: 'POST', data, ...config });
}

/** PUT 请求 */
export function put<T>(
  url: string,
  data?: Record<string, any>,
  config?: Omit<RequestConfig, 'url' | 'method' | 'data'>
): Promise<T> {
  return request<T>({ url, method: 'PUT', data, ...config });
}

/** DELETE 请求 */
export function del<T>(
  url: string,
  config?: Omit<RequestConfig, 'url' | 'method'>
): Promise<T> {
  return request<T>({ url, method: 'DELETE', ...config });
}

/** PATCH 请求 */
export function patch<T>(
  url: string,
  data?: Record<string, any>,
  config?: Omit<RequestConfig, 'url' | 'method' | 'data'>
): Promise<T> {
  return request<T>({ url, method: 'PATCH', data, ...config });
}

// ==================== 初始化默认拦截器 ====================

/**
 * 初始化请求拦截器
 */
export function initInterceptors(): void {
  // 请求拦截器：添加时间戳防止缓存
  addRequestInterceptor((config) => {
    // 非 GET 请求添加随机参数
    if (config.method !== 'GET') {
      config.params = {
        ...config.params,
        _t: Date.now(),
      };
    }
    return config;
  });

  // 错误拦截器：统一错误日志
  addErrorInterceptor((error) => {
    console.error('[API Error]', {
      url: error.url,
      code: error.code,
      message: error.message,
      time: new Date().toISOString(),
    });
  });
}
