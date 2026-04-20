<script setup lang="ts">
import type { CommentTargetType } from '@@/types/comment';
import CommentInput from './CommentInput.vue';
import CommentList from './CommentList.vue';
import CommentEmpty from './CommentEmpty.vue';
import { loadEmojiMap } from '@/composables/useEmojis';

// 组件属性
const props = defineProps<{
  targetType: CommentTargetType; // 目标类型 (article/page)
  targetKey: string | number; // 目标键值 (文章slug或页面key)
}>();

// 使用评论 store
const { comments, fetchComments, addComment, removeComment, resetComments, flattenComments } =
  useComments();

// 加载表情映射
const { blogConfig } = useSysConfig();
const emojiMapLoaded = ref(false);

const initEmojiMap = async () => {
  const emojisUrl = blogConfig.value.emojis;
  if (emojisUrl) {
    await loadEmojiMap(emojisUrl);
  }
  emojiMapLoaded.value = true;
};

onMounted(initEmojiMap);

// 全局登录弹窗
const { open: openLogin } = useLoginModal();

// 回复状态管理
const replyingToId = ref<number | null>(null);
const replyingToNickname = ref<string>('');

const route = useRoute();

let hasHandledInitialHash = false;

const scrollToComment = (hash?: string | null) => {
  if (!hash || !hash.startsWith('#comment-')) return;

  nextTick(() => {
    const elementId = hash.slice(1);
    if (!document.getElementById(elementId)) return;

    scrollToElement(hash);
    if (window.history?.replaceState) {
      const { pathname, search } = window.location;
      window.history.replaceState(null, '', pathname + search);
    }
  });
};

// 监听目标变化，自动加载评论
watch(
  () => [props.targetType, props.targetKey],
  ([type, key]) => {
    // 只在 targetKey 有效时才加载评论，避免传递 undefined
    if (key) {
      fetchComments(type as CommentTargetType, key as string | number);
    }
  },
  { immediate: true }
);

// 评论加载完成后滚动一次
watch(
  comments,
  newComments => {
    if (!hasHandledInitialHash && newComments.length > 0) {
      scrollToComment(route.hash);
      hasHandledInitialHash = true;
    }
  },
  { flush: 'post' }
);

// 监听 hash 变化
watch(
  () => route.hash,
  hash => {
    scrollToComment(hash);
  }
);

// 组件卸载时重置评论
onUnmounted(resetComments);

// 是否有评论
const hasComments = computed(() => comments.value.length > 0);

// 扁平化的评论列表（用于渲染）
const flatComments = computed(() => flattenComments(comments.value));

// 计算所有评论总数（包括回复）
const totalCommentsCount = computed(() => {
  const count = (list: typeof comments.value): number =>
    list.reduce((total, c) => total + 1 + (c.replies?.length ? count(c.replies) : 0), 0);
  return count(comments.value);
});

// 处理评论提交（顶层评论）
const handleAddComment = async (
  content: string,
  guestInfo?: { nickname?: string; email?: string; website?: string }
) => {
  await addComment({
    target_type: props.targetType,
    target_key: props.targetKey,
    content,
    ...guestInfo,
  });
};

// 处理回复提交
const handleAddReply = async (
  commentId: number,
  content: string,
  guestInfo?: { nickname?: string; email?: string; website?: string }
) => {
  await addComment({
    target_type: props.targetType,
    target_key: props.targetKey,
    content,
    parent_id: commentId,
    ...guestInfo,
  });
};

// 开始回复
const startReply = (id: number, nickname: string) => {
  if (replyingToId.value === id) {
    // 如果点击的是同一个评论，则取消回复
    replyingToId.value = null;
    replyingToNickname.value = '';
  } else {
    replyingToId.value = id;
    replyingToNickname.value = nickname;
  }
};

// 取消回复
const cancelReply = () => {
  replyingToId.value = null;
  replyingToNickname.value = '';
};

// 提供评论上下文给所有子组件
provideCommentContext({
  targetType: computed(() => props.targetType),
  targetKey: computed(() => props.targetKey),
  addComment: handleAddComment,
  addReply: handleAddReply,
  deleteComment: removeComment,
  showLogin: openLogin,
  replyState: {
    replyingToId,
    replyingToNickname,
    startReply,
    cancelReply,
  },
});
</script>

<template>
  <div class="comments-section">
    <div class="comments-header">
      <h3 class="comments-title">
        <i class="ri-chat-3-line" />
        评论
        <span v-if="totalCommentsCount > 0" class="comments-count">({{ totalCommentsCount }})</span>
      </h3>
    </div>

    <!-- 评论输入框 -->
    <CommentInput />

    <!-- 评论列表 -->
    <CommentList v-if="emojiMapLoaded && hasComments" :comments="flatComments" />

    <!-- 无评论状态 -->
    <CommentEmpty v-else-if="emojiMapLoaded && !hasComments" />
  </div>
</template>

<style lang="scss" scoped>
.comments-section {
  margin-top: 40px;
  padding-top: 30px;
  border-top: 1px solid var(--flec-border-color);
}

.comments-header {
  margin-bottom: 30px;

  .comments-title {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--font-color);
    display: flex;
    align-items: center;
    gap: 10px;
    margin: 0;

    i {
      font-size: 1.6rem;
    }

    .comments-count {
      font-size: 1rem;
      color: var(--flec-lightText);
      font-weight: normal;
    }
  }
}

// 响应式设计
@media screen and (max-width: 768px) {
  .comments-section {
    margin-top: 30px;
    padding-top: 20px;
  }

  .comments-header {
    margin-bottom: 20px;

    .comments-title {
      font-size: 1.25rem;
    }
  }
}
</style>

<style lang="scss">
// Markdown 样式 - 评论模块共用（不使用 scoped，让子组件可以使用）
.comments-section .markdown-body {
  font-family: inherit;

  p {
    margin: 0.5em 0;
  }

  strong {
    font-weight: 600;
  }

  em {
    font-style: italic;
  }

  code {
    padding: 2px 6px;
    background-color: rgba(0, 0, 0, 0.06);
    border-radius: 3px;
    font-size: 0.9em;
    font-family: 'Consolas', 'Monaco', monospace;
  }

  pre {
    padding: 10px;
    background-color: rgba(0, 0, 0, 0.04);
    border-radius: 4px;
    overflow-x: auto;
    font-family: 'Consolas', 'Monaco', monospace;

    code {
      padding: 0;
      background: transparent;
      font-size: inherit;
    }
  }

  blockquote {
    margin: 0.5em 0;
    padding-left: 12px;
    border-left: 3px solid var(--theme-color);
    opacity: 0.8;
  }

  ul,
  ol {
    margin: 0.5em 0;
    padding-left: 1.8em;
  }

  a {
    color: var(--theme-color);
    text-decoration: none;

    &:hover {
      text-decoration: underline;
    }
  }

  img {
    max-width: 100%;
  }

  // 表情图片样式
  .emoji-image {
    display: inline-block;
    width: 36px;
    height: 36px;
    margin: 0 2px;
    object-fit: contain;
  }
}
</style>
