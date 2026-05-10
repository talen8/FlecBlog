<script lang="ts" setup>
import type { CommentTargetType, Comment } from '@@/types';

const props = defineProps<{
  targetType: CommentTargetType;
  targetKey: string | number;
}>();

const { comments, fetchComments } = useComments();
const replyTo = ref<Comment | null>(null);

onMounted(() => {
  fetchComments(props.targetType, props.targetKey);
});

const handleReply = (comment: Comment) => {
  replyTo.value = comment;
};

const cancelReply = () => {
  replyTo.value = null;
};

const handleSubmitted = () => {
  replyTo.value = null;
};
</script>

<template>
  <div class="comments">
    <h3 class="comments-title">
      <i class="ri-message-3-line" />
      评论 ({{ comments.length }})
    </h3>

    <FeaturesCommentInput
      :target-type="targetType"
      :target-key="targetKey"
      :reply-to="replyTo"
      @cancel-reply="cancelReply"
      @submitted="handleSubmitted"
    />

    <div v-if="comments.length > 0" class="comment-list">
      <FeaturesCommentItem
        v-for="comment in comments"
        :key="comment.id"
        :comment="comment"
        @reply="handleReply"
      />
    </div>

    <div v-else class="empty-state">
      <i class="ri-chat-3-line" />
      <p>暂无评论，快来抢沙发吧！</p>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.comments {
  margin-top: 40px;
  padding-top: 30px;
  border-top: 1px solid var(--flec-border);
}

.comments-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 1.2rem;
  font-weight: 600;
  margin-bottom: 24px;
}

.comment-list {
  display: flex;
  flex-direction: column;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--theme-meta-color);
  i {
    font-size: 3rem;
    margin-bottom: 16px;
    display: block;
  }
  p {
    margin: 0;
  }
}
</style>
