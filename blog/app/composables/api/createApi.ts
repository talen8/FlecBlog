import type { ApiResponse, PaginationData, PaginationQuery } from '@@/types/request';
import { get, post, put, patch, del } from '@/utils/request';

interface ApiFactoryOptions {
  stringifyTargetKey?: boolean;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 转换函数，输入输出类型由调用方决定
  transformParams?: (params: any) => any;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 转换函数，输入输出类型由调用方决定
  transformBody?: (body: any) => any;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 处理任意数据结构
function processData(data: any, options: ApiFactoryOptions, isBody: boolean) {
  let processed = { ...data };

  if (options.stringifyTargetKey && processed.target_key !== undefined) {
    processed.target_key = String(processed.target_key);
  }

  const transformFn = isBody ? options.transformBody : options.transformParams;
  if (transformFn) {
    processed = transformFn(processed);
  }

  return processed;
}

export function createApi<T>(endpoint: string, options: ApiFactoryOptions = {}) {
  return {
    getList: async (params?: Partial<PaginationQuery>): Promise<PaginationData<T>> => {
      const response = await get<ApiResponse<PaginationData<T>>>(endpoint, {
        params: processData(params, options, false),
      });
      return response.data;
    },

    getOne: async (id: number | string): Promise<T> => {
      const response = await get<ApiResponse<T>>(`${endpoint}/${id}`);
      return response.data;
    },

    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 数据类型由调用方决定
    create: async (data: any): Promise<T> => {
      const response = await post<ApiResponse<T>>(endpoint, processData(data, options, true));
      return response.data;
    },

    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 数据类型由调用方决定
    update: async (id: number | string, data: any): Promise<T> => {
      const response = await put<ApiResponse<T>>(
        `${endpoint}/${id}`,
        processData(data, options, true)
      );
      return response.data;
    },

    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 数据类型由调用方决定
    patch: async (id: number | string, data: any): Promise<T> => {
      const response = await patch<ApiResponse<T>>(
        `${endpoint}/${id}`,
        processData(data, options, true)
      );
      return response.data;
    },

    delete: async (id: number | string): Promise<void> => {
      await del(`${endpoint}/${id}`);
    },

    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 参数类型由调用方决定
    get: async <R = T>(url: string, params?: any): Promise<R> => {
      const response = await get<ApiResponse<R>>(`${endpoint}${url}`, {
        params: processData(params, options, false),
      });
      return response.data;
    },

    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 数据类型由调用方决定
    post: async <R = T>(url: string, data?: any): Promise<R> => {
      const response = await post<ApiResponse<R>>(
        `${endpoint}${url}`,
        processData(data, options, true)
      );
      return response.data;
    },

    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 数据类型由调用方决定
    put: async <R = T>(url: string, data?: any): Promise<R> => {
      const response = await put<ApiResponse<R>>(
        `${endpoint}${url}`,
        processData(data, options, true)
      );
      return response.data;
    },

    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 数据类型由调用方决定
    patchRequest: async <R = T>(url: string, data?: any): Promise<R> => {
      const response = await patch<ApiResponse<R>>(
        `${endpoint}${url}`,
        processData(data, options, true)
      );
      return response.data;
    },

    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 数据类型由调用方决定
    deleteRequest: async (url: string, data?: any): Promise<void> => {
      await del(`${endpoint}${url}`, {
        body: processData(data, options, true),
      });
    },
  };
}
