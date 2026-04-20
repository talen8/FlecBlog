import axios from 'axios';
import type { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import { getAccessToken, getRefreshToken, setTokens, redirectToLogin } from '@/utils/auth';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
interface ApiResponse<T = any> {
  code: number;
  data: T;
  message: string;
}

// 获取 API URL（优先使用运行时配置）
const getApiUrl = () => {
  return (
    window.__APP_CONFIG__?.apiUrl || import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
  );
};

// 创建 axios 实例
const request = axios.create({
  baseURL: getApiUrl(),
  timeout: 15000,
  headers: { 'Content-Type': 'application/json' },
});

// 是否正在刷新token的标志
let isRefreshing = false;
// 存储待重试的请求
let failedQueue: Array<{
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  resolve: (value?: any) => void;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  reject: (reason?: any) => void;
}> = [];

// 处理队列中的请求
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const processQueue = (error: any = null) => {
  failedQueue.forEach(promise => {
    if (error) {
      promise.reject(error);
    } else {
      promise.resolve();
    }
  });
  failedQueue = [];
};

// 请求拦截器
request.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  // refresh接口不需要带Authorization header（在body中发送refresh_token）
  if (config.url === '/auth/refresh') {
    return config;
  }

  // 其他接口带上access token
  const token = getAccessToken();
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    // blob 类型直接返回
    if (response.config.responseType === 'blob') {
      return response.data;
    }
    const { code, message, data } = response.data;
    return code === 0 ? data : Promise.reject(new Error(message || '请求失败'));
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & {
      _retry?: boolean;
    };

    // 处理 blob 请求的错误响应（后端返回 JSON 错误）
    if (originalRequest.responseType === 'blob' && error.response?.data instanceof Blob) {
      const text = await (error.response.data as Blob).text();
      try {
        const json = JSON.parse(text);
        return Promise.reject(new Error(json.message || '请求失败'));
      } catch {
        return Promise.reject(error);
      }
    }

    // 处理401未授权 - 尝试刷新token
    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        // 如果正在刷新，将请求加入队列
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then(() => {
            return request(originalRequest);
          })
          .catch(err => {
            return Promise.reject(err);
          });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const refresh = getRefreshToken();
      if (!refresh) {
        // 没有refresh token，直接跳转登录页
        redirectToLogin();
        return Promise.reject(error);
      }

      try {
        // 调用refresh接口（返回的已经是data，不是整个response）
        const data: { access_token: string; refresh_token: string } = await request.post(
          '/auth/refresh',
          {
            refresh_token: refresh,
          }
        );

        // 更新token
        setTokens(data.access_token, data.refresh_token);

        // 处理队列中的请求
        processQueue();

        // 重试原请求
        return request(originalRequest);
      } catch (refreshError) {
        // 刷新失败，清空队列并跳转登录页
        processQueue(refreshError);
        redirectToLogin();
        return Promise.reject(refreshError);
      } finally {
        isRefreshing = false;
      }
    }

    // 其他错误直接返回
    return Promise.reject(error);
  }
);

export default request;
