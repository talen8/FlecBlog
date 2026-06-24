import type { ApiResponse } from '../../../types/request';
import type { UploadType, UploadResponse } from '../../../types/upload';
import { post } from './http';

/** 上传文件 */
export const uploadFile = async (file: File, type: UploadType): Promise<UploadResponse> => {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('type', type);

  const response = await post<ApiResponse<UploadResponse>>('/upload', formData);
  return response.data;
};
