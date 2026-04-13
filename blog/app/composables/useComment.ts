import type { Comment, CommentTargetType } from '@@/types/comment'

/**
 * 将树形评论结构递归转换为扁平数组
 * @param commentList 评论列表
 * @param depth 当前深度
 * @returns 扁平化的评论数组
 */
export function flattenComments(
  commentList: Comment[],
  depth = 0
): Array<{ comment: Comment; depth: number }> {
  const result: Array<{ comment: Comment; depth: number }> = []

  commentList.forEach(comment => {
    result.push({ comment, depth })
    if (comment.replies && comment.replies.length > 0) {
      result.push(...flattenComments(comment.replies, depth + 1))
    }
  })

  return result
}

/**
 * 游客信息类型
 */
export interface GuestInfo {
  nickname?: string
  email?: string
  website?: string
}

/**
 * 评论上下文类型定义
 */
export interface CommentContext {
  // 目标类型 (article/page)
  targetType: Ref<CommentTargetType>
  // 目标键值 (文章slug或页面key)
  targetKey: Ref<string | number>
  // 添加评论（顶层评论）
  addComment: (content: string, guestInfo?: GuestInfo) => Promise<void>
  // 添加回复
  addReply: (commentId: number, content: string, guestInfo?: GuestInfo) => Promise<void>
  // 删除评论
  deleteComment: (commentId: number) => Promise<void>
  // 显示登录模态框
  showLogin: () => void
  // 回复状态管理
  replyState: {
    replyingToId: Ref<number | null>
    replyingToNickname: Ref<string>
    startReply: (commentId: number, nickname: string) => void
    cancelReply: () => void
  }
}

// 注入键
const CommentContextKey: InjectionKey<CommentContext> = Symbol('CommentContext')

/**
 * 提供评论上下文
 */
export function provideCommentContext(context: CommentContext) {
  provide(CommentContextKey, context)
}

/**
 * 注入评论上下文
 */
export function useCommentContext() {
  const context = inject(CommentContextKey)
  if (!context) {
    throw new Error('useCommentContext must be used within a comment provider')
  }
  return context
}

/**
 * 填充评论内容并滚动到评论区
 * @param content 要填充的内容
 */
export async function fillComment(content: string) {
  const wrapper = document.querySelector('.comment-input')
  const textarea = wrapper?.querySelector('textarea') as HTMLTextAreaElement | null
  
  if (!wrapper || !textarea) return

  // 先填充并调整高度
  textarea.value = content
  textarea.dispatchEvent(new Event('input', { bubbles: true }))
  
  // 等待浏览器完成重排（双帧确保稳定）
  await new Promise(resolve => {
    requestAnimationFrame(() => requestAnimationFrame(resolve))
  })
  
  // 平滑滚动到评论区
  scrollToElement('.comment-input')
  textarea.focus()
}
