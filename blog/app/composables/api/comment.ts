import type { Comment, CommentTargetType, CreateCommentParams } from '@@/types/comment'
import type { PaginationData, PaginationQuery } from '@@/types/request'
import { createApi } from './createApi'

interface GetCommentsParams extends PaginationQuery {
  target_type: CommentTargetType
  target_key: string | number
}

const commentApi = createApi<Comment>('/comments', { stringifyTargetKey: true })

/** 获取评论列表 */
export const getComments = async (params: GetCommentsParams) => {
  return commentApi.getList(params)
}

/** 创建评论 */
export const createComment = async (params: CreateCommentParams) => {
  return commentApi.create(params)
}

/** 删除评论（仅可删除自己的评论） */
export const deleteComment = async (id: number) => {
  return commentApi.delete(id)
}
