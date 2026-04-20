import request from '@/utils/request';
import type { FileListData, FileListQuery } from '@/types/file';

/**
 * 上传文件响应接口
 */
export interface UploadResponse {
  file_url: string;
  file_name: string;
  file_size: number;
}

/**
 * 上传文件
 * @param file 要上传的文件
 * @param type 文件类型（默认为'image'）
 * @returns Promise<UploadResponse>
 */
export async function uploadFile(file: File, type = 'image'): Promise<UploadResponse> {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('type', type);
  try {
    return await request.post('/admin/files', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
  } catch (error: unknown) {
    const err = error as { response?: { data?: { message?: string } }; message?: string };
    const serverMessage = err.response?.data?.message;
    const errorMessage = serverMessage || err.message || '上传失败';
    throw new Error(errorMessage);
  }
}

/**
 * 获取文件列表
 * @param params 查询参数
 * @returns Promise<FileListData>
 */
export function getFileList(params: FileListQuery): Promise<FileListData> {
  return request.get('/admin/files', { params });
}

/**
 * 删除文件
 * @param id 文件ID
 * @returns Promise<void>
 */
export function deleteFile(id: number): Promise<void> {
  return request.delete(`/admin/files/${id}`);
}
