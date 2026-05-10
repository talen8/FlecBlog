<script lang="ts" setup>
import type { Comment } from '@@/types';

const props = defineProps<{ comment: Comment }>();
const emit = defineEmits<{ reply: [comment: Comment] }>();

const { userInfo } = useUser();
const { removeComment } = useComments();

const canDelete = computed(() => {
  if (!userInfo.value) return false;
  return (
    userInfo.value.id === props.comment.user.id ||
    userInfo.value.role === 'admin' ||
    userInfo.value.role === 'super_admin'
  );
});

const handleReply = () => {
  emit('reply', props.comment);
};

const handleDelete = async () => {
  if (!confirm('确定要删除这条评论吗？')) return;
  await removeComment(props.comment.id);
};

const showReplies = ref(false);
</script>

<template>
  <div class="comment-item" :class="{ 'is-deleted': comment.is_deleted }">
    <div class="comment-avatar">
      <img :src="getAvatarUrl(comment.user)" alt="avatar" />
    </div>
    <div class="comment-body">
      <div class="comment-header">
        <div class="user-info">
          <span class="nickname">{{ comment.user.nickname }}</span>
          <span v-if="comment.user.badge" class="badge">{{ comment.user.badge }}</span>
        </div>
        <span class="time">{{ formatMomentTime(comment.created_at) }}</span>
      </div>
      <div class="comment-content">
        <template v-if="comment.is_deleted">
          <span class="deleted-text">该评论已被删除</span>
        </template>
        <template v-else>
          <span v-if="comment.reply_user" class="reply-to">
            回复
            <a :href="'#comment-' + comment.reply_user.id">@{{ comment.reply_user.nickname }}</a
            >：
          </span>
          {{ comment.content }}
        </template>
      </div>
      <div class="comment-meta">
        <span v-if="comment.location"><i class="ri-map-pin-line" /> {{ comment.location }}</span>
        <span v-if="comment.browser"><i class="ri-chrome-line" /> {{ comment.browser }}</span>
        <span v-if="comment.os"><i class="ri-computer-line" /> {{ comment.os }}</span>
      </div>
      <div v-if="!comment.is_deleted" class="comment-actions">
        <button class="action-btn" @click="handleReply"><i class="ri-reply-line" /> 回复</button>
        <button v-if="canDelete" class="action-btn delete" @click="handleDelete">
          <i class="ri-delete-bin-line" /> 删除
        </button>
      </div>
      <div v-if="comment.replies?.length" class="comment-replies">
        <button class="toggle-replies" @click="showReplies = !showReplies">
          {{ showReplies ? '收起回复' : `展开 ${comment.replies.length} 条回复` }}
        </button>
        <div v-show="showReplies" class="replies-list">
          <FeaturesCommentItem
            v-for="reply in comment.replies"
            :key="reply.id"
            :comment="reply"
            @reply="$emit('reply', $event)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.comment-item {
  display: flex;
  gap: 12px;
  padding: 16px 0;
  &.is-deleted {
    opacity: 0.6;
  }
}

.comment-avatar {
  flex-shrink: 0;
  img {
    width: 40px;
    height: 40px;
    object-fit: cover;
  }
}

.comment-body {
  flex: 1;
  min-width: 0;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nickname {
  font-weight: 500;
  color: var(--font-color);
}

.badge {
  padding: 2px 6px;
  background: var(--theme-color);
  color: #fff;
  font-size: 0.75rem;
}

.time {
  color: var(--theme-meta-color);
  font-size: 0.85rem;
}

.comment-content {
  line-height: 1.6;
  color: var(--font-color);
}

.reply-to {
  color: var(--theme-meta-color);
  a {
    color: var(--theme-color);
    text-decoration: none;
  }
}

.deleted-text {
  color: var(--theme-meta-color);
  font-style: italic;
}

.comment-meta {
  display: flex;
  gap: 12px;
  margin-top: 6px;
  font-size: 0.8rem;
  color: var(--theme-meta-color);
  i {
    margin-right: 2px;
  }
}

.comment-actions {
  display: flex;
  gap: 16px;
  margin-top: 8px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  color: var(--theme-meta-color);
  font-size: 0.85rem;
  background: none;
  border: none;
  cursor: pointer;
  &.delete {
    color: #ff4d4f;
  }
}

.comment-replies {
  margin-top: 12px;
  padding-left: 20px;
  border-left: 2px solid var(--flec-border);
}

.toggle-replies {
  padding: 4px 0;
  color: var(--theme-color);
  font-size: 0.85rem;
  background: none;
  border: none;
  cursor: pointer;
}

.replies-list {
  margin-top: 8px;
}
</style>
