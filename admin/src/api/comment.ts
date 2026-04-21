import request from '@/utils/request';
import type {
  Comment,
  CommentListData,
  CommentListQuery,
  ImportCommentsResult,
} from '@/types/comment';

/**
 * 获取评论列表
 * @param params 查询参数
 * @returns Promise<CommentListData>
 */
export function getComments(params: CommentListQuery): Promise<CommentListData> {
  return request.get('/admin/comments', { params });
}

/**
 * 创建评论（用于回复）
 * @param data 评论数据
 * @returns Promise<Comment>
 */
export function createComment(data: {
  content: string;
  target_type: string;
  target_key: string;
  parent_id?: number;
}): Promise<Comment> {
  return request.post('/admin/comments', data);
}

/**
 * 切换评论状态
 * @param id 评论ID
 * @returns Promise<void>
 */
export function toggleCommentStatus(id: number): Promise<void> {
  return request.put(`/admin/comments/${id}/toggle-status`);
}

/**
 * 删除评论
 * @param id 评论ID
 * @returns Promise<void>
 */
export function deleteComment(id: number): Promise<void> {
  return request.delete(`/admin/comments/${id}`);
}

/**
 * 导入评论
 * @param formData 包含文件和参数的 FormData
 * @returns Promise<ImportCommentsResult>
 */
export function importComments(formData: FormData): Promise<ImportCommentsResult> {
  return request.post('/admin/comments/import', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}
