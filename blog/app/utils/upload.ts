import type { ApiResponse } from '@@/types/request';

export type UploadType = '用户头像' | '评论贴图' | '反馈投诉';

export interface UploadResponse {
  original_name: string;
  file_url: string;
}

export function getMaxFileSizeMB(): number {
  try {
    const { uploadConfig } = useSysConfig();
    const configValue = uploadConfig.value['max_file_size'] || '5';
    const parsed = parseInt(configValue, 10);
    return isNaN(parsed) || parsed <= 0 ? 5 : parsed;
  } catch {
    return 5;
  }
}

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
