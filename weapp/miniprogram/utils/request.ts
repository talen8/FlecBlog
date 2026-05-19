/**
 * API 请求封装
 * 基于 wx.request 封装为 Promise API，支持请求/响应拦截器
 */

import type { ApiResponse, ApiError, RequestConfig, RequestInterceptor, ErrorInterceptor } from '../types';
import { API_CONFIG, ERROR_MESSAGES } from '../config';

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
/** 等待刷新的请求队列 */
let refreshSubscribers: Array<(token: string) => void> = [];

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
 * 添加到 Token 刷新等待队列
 */
function subscribeTokenRefresh(callback: (token: string) => void): void {
  refreshSubscribers.push(callback);
}

/**
 * 通知所有等待的请求 Token 已刷新
 */
function notifyTokenRefresh(token: string): void {
  refreshSubscribers.forEach((callback) => callback(token));
  refreshSubscribers = [];
}

// ==================== URL 构建 ====================

/**
 * 构建完整请求 URL
 * @param url - 请求路径
 * @param params - 查询参数
 * @returns 完整 URL
 */
function buildUrl(url: string, params?: Record<string, unknown>): string {
  let fullUrl = url.startsWith('http') ? url : `${API_CONFIG.BASE_URL}${url}`;

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
        timeout: API_CONFIG.TIMEOUT,
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
                message: apiResponse.message || ERROR_MESSAGES.REQUEST_FAILED,
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
              message: ERROR_MESSAGES.SERVER_ERROR,
              url: processedConfig.url,
              statusCode,
            };
            runErrorInterceptors(apiError, processedConfig).catch(() => {
              reject(new Error(ERROR_MESSAGES.SERVER_ERROR));
            });
          } else {
            const apiError: ApiError = {
              code: statusCode,
              message: ERROR_MESSAGES.REQUEST_FAILED,
              url: processedConfig.url,
              statusCode,
            };
            runErrorInterceptors(apiError, processedConfig).catch(() => {
              reject(new Error(ERROR_MESSAGES.REQUEST_FAILED));
            });
          }
        },
        fail: (err) => {
          if (loading) wx.hideLoading();

          let errorMsg: string = ERROR_MESSAGES.NETWORK_ERROR;
          if (err.errMsg && err.errMsg.includes('timeout')) {
            errorMsg = ERROR_MESSAGES.TIMEOUT_ERROR;
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
 */
function handleUnauthorized<T>(
  config: RequestConfig,
  resolve: (value: T) => void,
  reject: (reason?: unknown) => void
): void {
  // 如果已经在刷新 Token，将请求加入队列等待
  if (isRefreshing) {
    subscribeTokenRefresh(() => {
      // Token 刷新后，重新执行请求
      request<T>(config).then(resolve).catch(reject);
    });
    return;
  }

  isRefreshing = true;

  // TODO: 调用刷新 Token 接口
  // 这里需要根据实际后端接口实现
  // refreshToken()
  //   .then((newToken) => {
  //     setToken(newToken);
  //     isRefreshing = false;
  //     notifyTokenRefresh(newToken);
  //     // 重试当前请求
  //     request<T>(config).then(resolve).catch(reject);
  //   })
  //   .catch(() => {
  //     isRefreshing = false;
  //     notifyTokenRefresh('');
  //     clearToken();
  //     wx.showToast({ title: ERROR_MESSAGES.UNAUTHORIZED, icon: 'none' });
  //     reject(new Error(ERROR_MESSAGES.UNAUTHORIZED));
  //   });

  // 临时方案：直接提示登录
  wx.showToast({ title: ERROR_MESSAGES.UNAUTHORIZED, icon: 'none' });
  reject(new Error(ERROR_MESSAGES.UNAUTHORIZED));
}

/**
 * 设置基础认证 Token
 */
function setAuthHeader(token: string | null): void {
  if (token) {
    setToken(token);
  } else {
    clearToken();
  }
}

// ==================== 便捷请求方法 ====================

/** GET 请求 */
export function get<T>(
  url: string,
  params?: Record<string, unknown>,
  config?: Omit<RequestConfig, 'url' | 'method' | 'params'>
): Promise<T> {
  return request<T>({ url, method: 'GET', params, ...config });
}

/** POST 请求 */
export function post<T>(
  url: string,
  data?: Record<string, unknown>,
  config?: Omit<RequestConfig, 'url' | 'method' | 'data'>
): Promise<T> {
  return request<T>({ url, method: 'POST', data, ...config });
}

/** PUT 请求 */
export function put<T>(
  url: string,
  data?: Record<string, unknown>,
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
  data?: Record<string, unknown>,
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
