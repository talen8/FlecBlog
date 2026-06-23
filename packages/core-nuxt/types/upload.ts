/** 上传类型（决定后端归类与允许的文件格式） */
export type UploadType = '用户头像' | '评论贴图' | '反馈投诉';

/** 上传响应 */
export interface UploadResponse {
  original_name: string;
  file_url: string;
}
