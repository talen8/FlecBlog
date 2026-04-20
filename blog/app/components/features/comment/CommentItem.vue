<script setup lang="ts">
import type { Comment } from '@@/types/comment';
import CommentInput from './CommentInput.vue';

interface Props {
  comment: Comment;
  depth: number;
}

const props = defineProps<Props>();
const context = useCommentContext();
const { userInfo } = useUser();
const isLoggedIn = useAuth();

/** 判断当前登录用户是否为该评论的作者 */
const isCommentOwner = computed(
  () =>
    isLoggedIn.value && !props.comment.is_deleted && userInfo.value?.id === props.comment.user.id
);

/** 是否正在删除中 */
const isDeleting = ref(false);

/** 是否显示删除确认弹窗 */
const showDeleteDialog = ref(false);

/** 点击删除按钮，弹出确认框 */
const handleDeleteClick = () => {
  showDeleteDialog.value = true;
};

/** 确认删除评论 */
const handleDeleteConfirm = async () => {
  if (isDeleting.value) return;
  isDeleting.value = true;
  try {
    await context.deleteComment(props.comment.id);
  } catch (error) {
    console.error('删除评论失败:', error);
  } finally {
    isDeleting.value = false;
    showDeleteDialog.value = false;
  }
};

// 是否正在回复此评论
const isReplying = computed(() => context.replyState.replyingToId.value === props.comment.id);

// 长评论折叠控制（超过 300 字符）
const MAX_CONTENT_LENGTH = 300;
const isExpanded = ref(false);

const isLongComment = computed(() => props.comment.content.length > MAX_CONTENT_LENGTH);

const displayContent = computed(() => {
  if (!isLongComment.value || isExpanded.value) {
    return props.comment.content;
  }
  return props.comment.content.substring(0, MAX_CONTENT_LENGTH) + '...';
});

// 渲染评论内容为markdown
const renderedContent = computed(() => renderSimpleMarkdown(displayContent.value));

const toggleExpand = () => {
  isExpanded.value = !isExpanded.value;
};

// 处理回复按钮点击
const handleReplyClick = () => {
  const wasReplying = isReplying.value;
  context.replyState.startReply(props.comment.id, props.comment.user.nickname);

  // 如果是展开回复（而不是取消回复），则滚动到输入框
  if (!wasReplying) {
    nextTick(() => {
      scrollToElement(`#reply-input-${props.comment.id}`);
    });
  }
};
</script>

<template>
  <div :id="`comment-${comment.id}`" class="comment-item" :class="{ 'is-reply': depth > 0 }">
    <!-- 左：头像 -->
    <div class="comment-left">
      <a
        v-if="comment.user.website"
        :href="comment.user.website"
        target="_blank"
        rel="noopener noreferrer"
        class="comment-avatar-link"
      >
        <NuxtImg
          :src="getAvatarUrl(comment.user)"
          :alt="comment.user.nickname"
          class="comment-avatar"
          loading="lazy"
        />
      </a>
      <NuxtImg
        v-else
        :src="getAvatarUrl(comment.user)"
        :alt="comment.user.nickname"
        class="comment-avatar"
        loading="lazy"
      />
    </div>

    <!-- 右：上中下 -->
    <div class="comment-right">
      <!-- 上：昵称、时间、系统信息 -->
      <div class="comment-header">
        <span class="comment-author">
          {{ comment.user.nickname }}
          <span v-if="comment.user.badge" class="comment-badge" :class="comment.user.role">{{
            comment.user.badge
          }}</span>
          <span v-if="comment.reply_user" class="reply-arrow">
            <i class="ri-arrow-right-s-fill" />
            {{ comment.reply_user.nickname }}
            <span
              v-if="comment.reply_user.badge"
              class="comment-badge"
              :class="comment.reply_user.role"
              >{{ comment.reply_user.badge }}</span
            >
          </span>
        </span>
        <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
        <span v-if="comment.location" class="comment-meta">{{ comment.location }}</span>
        <span v-if="comment.os" class="comment-meta">{{ comment.os }}</span>
        <span v-if="comment.browser" class="comment-meta">{{ comment.browser }}</span>
      </div>

      <!-- 中：内容 -->
      <!-- eslint-disable-next-line vue/no-v-html -- 评论内容经过 DOMPurify 净化处理 -->
      <div class="comment-body markdown-body" v-html="renderedContent" />

      <!-- 长评论展开/折叠按钮 -->
      <button
        v-if="isLongComment"
        class="expand-btn"
        aria-label="展开/折叠评论"
        @click="toggleExpand"
      >
        {{ isExpanded ? '收起' : '展开全文' }}
        <i :class="isExpanded ? 'ri-arrow-up-s-line' : 'ri-arrow-down-s-line'" />
      </button>

      <!-- 下：功能按钮 -->
      <div class="comment-actions">
        <button class="action-btn" aria-label="回复评论" @click="handleReplyClick">回复</button>
        <button
          v-if="isCommentOwner"
          class="action-btn delete-btn"
          :disabled="isDeleting"
          aria-label="删除评论"
          @click="handleDeleteClick"
        >
          {{ isDeleting ? '删除中...' : '删除' }}
        </button>
      </div>

      <!-- 回复输入框 -->
      <div v-if="isReplying" :id="`reply-input-${comment.id}`" class="reply-input-wrapper">
        <CommentInput
          :comment-id="comment.id"
          :reply-to="context.replyState.replyingToNickname.value"
        />
      </div>

      <!-- 子评论插槽 -->
      <slot name="replies" />
    </div>

    <!-- 删除确认弹窗 -->
    <UiBaseDialog
      v-model="showDeleteDialog"
      title="删除评论"
      confirm-text="确认删除"
      :loading="isDeleting"
      style="--dialog-width: 400px"
      @confirm="handleDeleteConfirm"
    >
      <p>确定要删除这条评论吗？删除后将无法恢复。</p>
    </UiBaseDialog>
  </div>
</template>

<style lang="scss" scoped>
.comment-item {
  display: flex;
  gap: 10px;

  &.is-reply {
    padding: 10px 0;

    &:not(:last-child) {
      border-bottom: 1px solid var(--flec-border-color);
    }
  }
}

.comment-left {
  flex-shrink: 0;
}

.comment-avatar {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  object-fit: cover;
}

.comment-avatar-link {
  display: block;
  transition: transform 0.2s;

  &:hover {
    transform: scale(1.05);
  }
}

.comment-right {
  flex: 1;
  min-width: 0;
  margin-left: 5px;
}

.comment-header {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  margin-bottom: 10px;
  margin-top: 5px;
  gap: 8px;
}

.comment-author {
  color: var(--theme-color);
  font-size: 0.85rem;
  display: inline-flex;
  align-items: center;
  line-height: 1;

  .reply-arrow {
    display: inline-flex;
    align-items: center;

    i {
      color: var(--theme-meta-color);
      font-size: 1rem;
      line-height: 1;
      display: inline-block;
      margin: 0 2px;
    }
  }

  .comment-badge {
    display: inline-block;
    font-size: 12px;
    line-height: 1;
    padding: 2px 6px;
    border-radius: 4px;
    margin-left: 6px;

    &.super_admin,
    &.admin {
      background: rgba(255, 193, 7, 0.15);
      color: #e3ae0d;
      border: 1px solid rgba(255, 193, 7, 0.3);
    }

    &.user,
    &.guest {
      background: rgba(64, 158, 255, 0.15);
      color: #409eff;
      border: 1px solid rgba(64, 158, 255, 0.3);
    }
  }
}

.comment-time {
  color: var(--theme-meta-color);
  font-size: 0.85rem;
  white-space: nowrap;
  line-height: 1;
}

.comment-meta {
  color: var(--theme-meta-color);
  font-size: 0.85rem;
  white-space: nowrap;
  line-height: 1;
  position: relative;
  padding-left: 12px;

  &::before {
    content: '·';
    position: absolute;
    left: 0;
    color: var(--font-color);
    opacity: 0.5;
  }
}

.comment-body {
  color: var(--font-color);
  line-height: 1.5;
  word-wrap: break-word;
}

.expand-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 8px;
  padding: 4px 8px;
  border: none;
  background: transparent;
  color: var(--theme-color);
  font-size: 0.85rem;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;

  &:hover {
    background: var(--flec-hover-bg);
  }

  i {
    font-size: 0.95rem;
  }
}

.comment-actions {
  display: flex;
  gap: 8px;
  margin-top: 10px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  border: none;
  background: transparent;
  color: var(--theme-meta-color);
  font-size: 0.85rem;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;

  i {
    font-size: 0.95rem;
  }

  &.delete-btn {
    &:hover:not(:disabled) {
      color: var(--flec-danger, #f56c6c);
    }

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }
}

.reply-input-wrapper {
  margin-top: 12px;
}

.replies-list {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--flec-border-color);
}

@media screen and (max-width: 768px) {
  .comment-avatar {
    width: 40px;
    height: 40px;
    border-radius: 6px;
  }

  .comment-header {
    gap: 6px;
  }

  .comment-author {
    font-size: 0.85rem;

    .reply-arrow {
      display: block;
      margin: 2px 0;
    }
  }

  .comment-time {
    font-size: 0.75rem;
  }

  .comment-meta {
    font-size: 0.75rem;
    white-space: normal;
    word-break: break-all;
    padding-left: 0;
    max-width: 100%;

    &::before {
      display: none;
    }
  }

  .comment-body {
    font-size: 0.85rem;
  }

  .action-btn {
    font-size: 0.75rem;
    padding: 3px 8px;
  }
}
</style>
