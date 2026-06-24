import type {
  Comment,
  CommentTargetType,
  FlatComment,
  CreateCommentParams,
  GuestInfo,
} from '../../types/comment';
import { getComments, createComment, deleteComment } from './api/comment';

/**
 * 将嵌套评论扁平化为带深度的列表（用于列表渲染）
 * @param commentList - 嵌套评论数组
 * @param depth - 当前递归深度，默认 0
 * @returns 扁平化后的评论列表
 */
export function flattenComments(commentList: Comment[], depth = 0): FlatComment[] {
  const result: FlatComment[] = [];

  commentList.forEach(comment => {
    result.push({ comment, depth });
    if (comment.replies && comment.replies.length > 0) {
      result.push(...flattenComments(comment.replies, depth + 1));
    }
  });

  return result;
}

/** 评论上下文接口（provide/inject 模式，供评论子组件使用） */
export interface CommentContext {
  targetType: Ref<CommentTargetType>;
  targetKey: Ref<string | number>;
  addComment: (content: string, guestInfo?: GuestInfo) => Promise<void>;
  addReply: (commentId: number, content: string, guestInfo?: GuestInfo) => Promise<void>;
  deleteComment: (commentId: number) => Promise<void>;
  showLogin: () => void;
  replyState: {
    replyingToId: Ref<number | null>;
    replyingToNickname: Ref<string>;
    startReply: (commentId: number, nickname: string) => void;
    cancelReply: () => void;
  };
}

const CommentContextKey: InjectionKey<CommentContext> = Symbol('CommentContext');

/**
 * 提供评论上下文（父组件调用，子组件通过 useCommentContext 注入）
 * @param context - 评论上下文对象
 */
export function provideCommentContext(context: CommentContext) {
  provide(CommentContextKey, context);
}

/**
 * 注入评论上下文（子组件调用，需在 provideCommentContext 之后）
 * @returns 评论上下文对象
 */
export function useCommentContext() {
  const context = inject(CommentContextKey);
  if (!context) {
    throw new Error('useCommentContext must be used within a comment provider');
  }
  return context;
}

/**
 * 填充评论文本框并滚动聚焦（用于"引用评论"等场景）
 * @param content - 要填充的评论文本
 */
export async function fillComment(content: string) {
  const wrapper = document.querySelector('.comment-input');
  const textarea = wrapper?.querySelector('textarea') as HTMLTextAreaElement | null;

  if (!wrapper || !textarea) return;

  textarea.value = content;
  textarea.dispatchEvent(new Event('input', { bubbles: true }));

  await new Promise(resolve => {
    requestAnimationFrame(() => requestAnimationFrame(resolve));
  });

  scrollToElement('.comment-input');
  textarea.focus();
}

/**
 * 评论列表管理（获取/添加/删除，全局共享状态）
 * @returns comments - 评论列表
 * @returns fetchComments - 拉取评论
 * @returns addComment - 添加评论
 * @returns removeComment - 删除评论
 * @returns resetComments - 重置
 * @returns flattenComments - 扁平化工具
 */
export function useComments() {
  const comments = useState<Comment[]>('comments', () => []);
  const currentTargetType = useState<CommentTargetType | null>('comments-targetType', () => null);
  const currentTargetKey = useState<string | number | null>('comments-targetKey', () => null);

  const fetchComments = async (targetType: CommentTargetType, targetKey: string | number) => {
    if (!targetType || !targetKey) return;

    currentTargetType.value = targetType;
    currentTargetKey.value = targetKey;

    try {
      const data = await getComments({
        target_type: targetType,
        target_key: targetKey,
      });
      comments.value = data.list || [];
    } catch (error) {
      console.error('获取评论失败:', error);
      comments.value = [];
    }
  };

  const addComment = async (params: CreateCommentParams) => {
    const newComment = await createComment(params);

    if (!params.parent_id) {
      comments.value.unshift(newComment);
    } else {
      const addReplyToComment = (commentList: Comment[]): boolean => {
        for (const comment of commentList) {
          if (comment.id === params.parent_id) {
            if (!comment.replies) comment.replies = [];
            comment.replies.push(newComment);
            return true;
          }
          if (comment.replies?.length && addReplyToComment(comment.replies)) return true;
        }
        return false;
      };
      addReplyToComment(comments.value);
    }

    return newComment;
  };

  const removeComment = async (commentId: number) => {
    await deleteComment(commentId);

    const removeFromList = (commentList: Comment[]): boolean => {
      const index = commentList.findIndex(c => c.id === commentId);
      if (index !== -1) {
        commentList.splice(index, 1);
        return true;
      }
      for (const comment of commentList) {
        if (comment.replies?.length && removeFromList(comment.replies)) return true;
      }
      return false;
    };
    removeFromList(comments.value);
  };

  return {
    comments,
    fetchComments,
    addComment,
    removeComment,
    resetComments: () => {
      comments.value = [];
      currentTargetType.value = null;
      currentTargetKey.value = null;
    },
    flattenComments,
  };
}
