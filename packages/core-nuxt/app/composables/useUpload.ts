import type { UploadType, UploadResponse } from '../../types/upload';
import { uploadFile } from './api/upload';

export type { UploadType, UploadResponse };

/**
 * 文件上传管理（含文件类型校验和大小限制）
 * @returns maxSizeMB - 最大文件大小（MB）
 * @returns allowedFileTypes(type) - 允许的文件类型
 * @returns validate(file, type) - 校验文件
 * @returns upload(file, type) - 执行上传
 */
export const useUpload = () => {
  const { uploadConfig } = useSysConfig();

  const maxSizeMB = computed(() => {
    const configValue = uploadConfig.value['max_file_size'] || '5';
    const parsed = parseInt(configValue, 10);
    return isNaN(parsed) || parsed <= 0 ? 5 : parsed;
  });

  const allowedFileTypes = (type: UploadType) => {
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
  };

  const validate = (file: File, type: UploadType): string | null => {
    const { allowedTypes, typeDescription } = allowedFileTypes(type);
    if (!allowedTypes.includes(file.type)) {
      return `只支持 ${typeDescription}`;
    }
    const maxSize = maxSizeMB.value * 1024 * 1024;
    if (file.size > maxSize) {
      return `文件大小不能超过 ${maxSizeMB.value}MB`;
    }
    return null;
  };

  const upload = async (file: File, type: UploadType): Promise<UploadResponse> => {
    const validationError = validate(file, type);
    if (validationError) {
      throw new Error(validationError);
    }

    try {
      return await uploadFile(file, type);
    } catch (error: unknown) {
      const err = error as { message?: string; response?: { data?: { message?: string } } };
      throw new Error(err?.message || err?.response?.data?.message || '文件上传失败');
    }
  };

  return {
    maxSizeMB,
    allowedFileTypes,
    validate,
    upload,
  };
};
