<script lang="ts" setup>
import type { Comment } from '@@/types';
import { uploadFile } from '@/utils/upload';

const props = defineProps<{
  targetType: string;
  targetKey: string | number;
  replyTo?: Comment | null;
}>();

const emit = defineEmits<{
  cancelReply: [];
  submitted: [];
}>();

const { userInfo } = useUser();
const { addComment } = useComments();
const toast = useToast();
const showLoginModal = useState('showLoginModal', () => false);

const content = ref('');
const submitting = ref(false);
const uploading = ref(false);

const handleSubmit = async () => {
  if (!userInfo.value) {
    showLoginModal.value = true;
    return;
  }

  if (!content.value.trim()) {
    toast.error('请输入评论内容');
    return;
  }

  submitting.value = true;
  try {
    await addComment({
      // eslint-disable-next-line @typescript-eslint/no-explicit-any -- target_type 类型断言，由父组件保证类型正确
      target_type: props.targetType as any,
      target_key: props.targetKey,
      content: content.value.trim(),
      parent_id: props.replyTo?.id,
    });
    content.value = '';
    toast.success('评论成功');
    emit('submitted');
  } catch (error) {
    console.error('评论失败:', error);
    toast.error('评论失败');
  } finally {
    submitting.value = false;
  }
};

const handleImageUpload = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;

  uploading.value = true;
  try {
    const result = await uploadFile(file, '评论贴图');
    content.value += ` ![](${result.file_url})`;
    toast.success('图片上传成功');
  } catch (error) {
    console.error('上传失败:', error);
    toast.error('图片上传失败');
  } finally {
    uploading.value = false;
    input.value = '';
  }
};
</script>

<template>
  <div class="comment-input">
    <div v-if="replyTo" class="reply-info">
      <span
        >回复 <strong>{{ replyTo.user.nickname }}</strong></span
      >
      <button class="cancel-btn" @click="$emit('cancelReply')">取消</button>
    </div>
    <div class="input-wrapper">
      <img :src="getAvatarUrl(userInfo || {})" alt="avatar" class="user-avatar" />
      <textarea
        v-model="content"
        :placeholder="userInfo ? '写下你的评论...' : '请先登录后评论'"
        :disabled="submitting"
        rows="3"
      />
    </div>
    <div class="input-actions">
      <div class="left-actions">
        <label class="upload-btn" :class="{ disabled: uploading || !userInfo }">
          <i class="ri-image-line" />
          <input
            type="file"
            accept="image/*"
            :disabled="uploading || !userInfo"
            @change="handleImageUpload"
          />
        </label>
      </div>
      <button class="submit-btn" :disabled="!content.trim() || submitting" @click="handleSubmit">
        {{ submitting ? '提交中...' : '发表评论' }}
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.comment-input {
  margin-bottom: 24px;
}

.reply-info {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
  padding: 8px 12px;
  background: var(--flec-heavy-bg);
  font-size: 0.9rem;
}

.cancel-btn {
  margin-left: auto;
  color: var(--theme-meta-color);
  background: none;
  border: none;
  cursor: pointer;
  font-size: 0.85rem;
}

.input-wrapper {
  display: flex;
  gap: 12px;
}

.user-avatar {
  width: 40px;
  height: 40px;
  object-fit: cover;
  flex-shrink: 0;
}

textarea {
  flex: 1;
  padding: 12px;
  border: 1px solid var(--flec-border);
  background: var(--flec-card-bg);
  color: var(--font-color);
  font-size: 0.95rem;
  resize: vertical;
  font-family: inherit;

  &:focus {
    outline: none;
    border-color: var(--theme-color);
  }
  &::placeholder {
    color: var(--theme-meta-color);
  }
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 12px;
  padding-left: 52px;
}

.left-actions {
  display: flex;
  gap: 8px;
}

.upload-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  color: var(--theme-meta-color);
  cursor: pointer;

  &.disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  input {
    display: none;
  }
}

.submit-btn {
  padding: 10px 24px;
  background: var(--theme-color);
  color: #fff;
  border: none;
  font-size: 0.9rem;
  cursor: pointer;

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}
</style>
