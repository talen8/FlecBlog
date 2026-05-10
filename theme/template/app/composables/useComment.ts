/**
 * 评论交互状态管理
 *
 * 管理评论的回复状态和内容，提供评论操作的辅助方法。
 * 用于评论组件之间的状态共享。
 *
 * @module composables/useComment
 */

import type { Comment } from '@@/types';

/**
 * 将嵌套评论扁平化
 *
 * 将树形结构的评论列表转换为扁平数组，
 * 同时保留 parent_id 信息用于判断层级关系。
 *
 * @param comments - 嵌套的评论列表
 * @returns 扁平化的评论数组
 *
 * @example
 * ```ts
 * const flatList = flattenComments(nestedComments);
 * // 返回所有评论的一维数组，包含 parent_id 信息
 * ```
 */
export const flattenComments = (comments: Comment[]): Comment[] => {
  const result: Comment[] = [];

  const flatten = (commentList: Comment[], parentId: number | null = null) => {
    for (const comment of commentList) {
      result.push({
        ...comment,
        parent_id: parentId,
      });
      if (comment.replies?.length) {
        flatten(comment.replies, comment.id);
      }
    }
  };

  flatten(comments);
  return result;
};

/**
 * 评论交互状态 Composable
 *
 * 管理评论回复的相关状态，包括回复目标、评论内容等。
 *
 * @returns {Object} 评论状态和方法
 * @returns {Ref<Comment | null>} replyTo - 回复目标的评论
 * @returns {Ref<string>} commentContent - 评论内容
 * @returns {Function} setReplyTo - 设置回复目标
 * @returns {Function} clearReply - 清除回复状态
 * @returns {Function} fillComment - 填充评论内容
 * @returns {Function} flattenComments - 扁平化评论列表
 *
 * @example
 * ```ts
 * const { replyTo, setReplyTo, clearReply } = useComment();
 *
 * // 设置回复目标
 * setReplyTo(comment);
 *
 * // 清除回复状态
 * clearReply();
 * ```
 */
export const useComment = () => {
  const replyTo = ref<Comment | null>(null);
  const commentContent = ref('');

  /**
   * 设置回复目标
   *
   * @param comment - 要回复的评论，设为 null 取消回复
   */
  const setReplyTo = (comment: Comment | null) => {
    replyTo.value = comment;
  };

  /**
   * 清除回复状态
   *
   * 重置回复目标和评论内容。
   */
  const clearReply = () => {
    replyTo.value = null;
    commentContent.value = '';
  };

  /**
   * 填充评论内容
   *
   * @param content - 评论内容
   */
  const fillComment = (content: string) => {
    commentContent.value = content;
  };

  return {
    replyTo,
    commentContent,
    setReplyTo,
    clearReply,
    fillComment,
    flattenComments,
  };
};
