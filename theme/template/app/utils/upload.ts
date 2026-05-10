/**
 * 文件上传工具函数
 *
 * 提供文件上传功能，包括文件验证、上传请求等。
 *
 * @module utils/upload
 */

import type { ApiResponse } from '@@/types';

/**
 * 上传类型
 */
export type UploadType = '用户头像' | '评论贴图' | '反馈投诉';

/**
 * 上传响应数据
 */
export interface UploadResponse {
  original_name: string;
  file_url: string;
}

/**
 * 获取最大文件大小限制（MB）
 *
 * 从系统配置中读取上传文件大小限制。
 *
 * @returns 最大文件大小（MB），默认 5MB
 */
export function getMaxFileSizeMB(): number {
  try {
    const { uploadConfig } = useSysConfig();
    const configValue = uploadConfig.value['upload.max_file_size'] || '5';
    const parsed = parseInt(configValue, 10);
    return isNaN(parsed) || parsed <= 0 ? 5 : parsed;
  } catch {
    return 5;
  }
}

/**
 * 获取允许的文件类型
 *
 * 根据上传类型返回允许的 MIME 类型列表和描述。
 *
 * @param type - 上传类型
 * @returns 允许的类型列表和描述
 */
export function getAllowedFileTypes(type: UploadType): {
  allowedTypes: string[];
  typeDescription: string;
} {
  if (type === '反馈投诉') {
    return {
      allowedTypes: [
        'image/jpeg',
        'image/jpg',
        'image/png',
        'image/gif',
        'image/webp',
        'application/pdf',
        'application/msword',
        'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
      ],
      typeDescription: 'JPG、PNG、GIF、WebP 格式的图片或 PDF、DOC、DOCX 格式的文档',
    };
  }
  return {
    allowedTypes: ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'],
    typeDescription: 'JPG、PNG、GIF、WebP 格式的图片',
  };
}

/**
 * 验证文件
 *
 * 检查文件类型和大小是否符合要求。
 *
 * @param file - 要验证的文件
 * @param type - 上传类型
 * @returns 验证错误信息，验证通过返回 null
 *
 * @example
 * ```ts
 * const error = validateFile(file, '用户头像');
 * if (error) {
 *   toast.error(error);
 *   return;
 * }
 * ```
 */
export function validateFile(file: File, type: UploadType): string | null {
  const { allowedTypes, typeDescription } = getAllowedFileTypes(type);

  if (!allowedTypes.includes(file.type)) {
    return `只支持 ${typeDescription}`;
  }

  const maxSizeMB = getMaxFileSizeMB();
  const maxSize = maxSizeMB * 1024 * 1024;
  if (file.size > maxSize) {
    return `文件大小不能超过 ${maxSizeMB}MB`;
  }

  return null;
}

/**
 * 上传文件
 *
 * 将文件上传到服务器，返回文件 URL。
 *
 * @param file - 要上传的文件
 * @param type - 上传类型
 * @returns 上传响应数据，包含文件 URL
 * @throws 文件验证失败或上传失败时抛出错误
 *
 * @example
 * ```ts
 * try {
 *   const result = await uploadFile(file, '用户头像');
 *   console.log(result.file_url);
 * } catch (error) {
 *   toast.error(error.message);
 * }
 * ```
 */
export async function uploadFile(file: File, type: UploadType): Promise<UploadResponse> {
  const validationError = validateFile(file, type);
  if (validationError) {
    throw new Error(validationError);
  }

  const formData = new FormData();
  formData.append('file', file);
  formData.append('type', type);

  const config = useRuntimeConfig();
  const baseURL = config.public.apiUrl;

  const response = await $fetch<ApiResponse<UploadResponse>>('/upload', {
    baseURL,
    method: 'POST',
    body: formData,
  }).catch((error: unknown) => {
    const err = error as { data?: { message?: string }; message?: string };
    throw new Error(err?.data?.message || err?.message || '文件上传失败');
  });

  if (response.code !== 0) {
    throw new Error(response.message || '文件上传失败');
  }

  return response.data;
}
