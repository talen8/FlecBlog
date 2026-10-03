import type { Comment, CreateCommentParams, GetCommentsParams } from '../../../types/comment';
import { createApi } from './createApi';

const commentApi = createApi<Comment>('/comments', { stringifyTargetKey: true });

/** 获取评论列表 */
export const getComments = async (params: GetCommentsParams) => {
  return commentApi.getList(params);
};

/** 创建评论 */
export const createComment = async (params: CreateCommentParams) => {
  return commentApi.create(params);
};

/** 删除评论（仅可删除自己的评论） */
export const deleteComment = async (id: number) => {
  return commentApi.delete(id);
};
