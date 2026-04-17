<script setup lang="ts">
import { uploadFile } from '~/utils/upload';

interface Props {
  // 评论ID，如果提供则为回复模式
  commentId?: number;
  // 回复的目标昵称
  replyTo?: string;
}

const props = defineProps<Props>();
const { success, info } = useToast();
const { triggerOnComment } = useBindEmail();

// 获取评论上下文
const context = useCommentContext();

// 状态（延迟从本地存储初始化，避免 SSR 时 localStorage 不存在）
const nickname = ref('');
const email = ref('');
const website = ref('');
const commentContent = ref('');

// 在客户端加载本地存储数据
onMounted(() => {
  // 加载游客信息
  const stored = localStorage.getItem('guest_info');
  if (stored) {
    const saved = JSON.parse(stored);
    nickname.value = saved.nickname || '';
    email.value = saved.email || '';
    website.value = saved.website || '';
  }
  // 加载评论草稿
  commentContent.value = localStorage.getItem('comment_draft') || '';
  // 根据草稿内容调整输入框高度
  nextTick(() => {
    resetTextareaHeight();
  });
});

const isSubmitting = ref(false);
const showPreview = ref(false);
const showExpandedBtn = ref(false);
const textareaRef = ref<HTMLTextAreaElement | null>(null);
const buttonGroupRef = ref<HTMLElement | null>(null);
const fileInputRef = ref<HTMLInputElement | null>(null);

// 错误提示
const errors = ref({
  nickname: '',
  email: '',
  website: '',
  content: '',
});

// 计算属性
const isLoggedIn = useAuth();
const isReplyMode = computed(() => !!props.replyTo);
const isUserInfoFilled = computed(() => nickname.value.trim() && email.value.trim());
const shouldShowSend = computed(() => isLoggedIn.value || isUserInfoFilled.value);
const mainBtn = computed(() => {
  if (isSubmitting.value) return { text: '发送中...', icon: 'ri-loader-4-line rotating' };
  return shouldShowSend.value
    ? { text: '发送', icon: 'ri-send-plane-fill' }
    : { text: '登录', icon: 'ri-login-box-line' };
});

const secondaryBtn = computed(() =>
  isUserInfoFilled.value
    ? { text: '登录', icon: 'ri-login-box-line' }
    : { text: '发送', icon: 'ri-send-plane-fill' }
);
const renderedMarkdown = computed(() => renderSimpleMarkdown(commentContent.value));
const guestPrivacyNotice = [
  '游客无需注册即可评论。',
  '你提交的昵称、邮箱、网址和评论内容会保存在服务端，用于展示评论身份、接收回复及必要的安全审计。',
  '浏览器会本地保存已填游客信息和评论草稿，方便下次免填。',
  '回复提醒会通过站内消息和邮件通知。',
];

// 工具函数
const resetTextareaHeight = () => {
  if (textareaRef.value) {
    textareaRef.value.style.height = 'auto';
    textareaRef.value.style.height = textareaRef.value.scrollHeight + 'px';
  }
};

// 重置输入框为默认高度
const resetToDefaultHeight = () => {
  if (textareaRef.value) {
    textareaRef.value.style.height = 'auto';
  }
};

const validateForm = () => {
  errors.value = { nickname: '', email: '', website: '', content: '' };

  if (!isLoggedIn.value) {
    const nick = nickname.value.trim();
    if (!nick) {
      errors.value.nickname = '请输入昵称';
      return false;
    }
    if (nick.length < 2 || nick.length > 32) {
      errors.value.nickname = nick.length < 2 ? '昵称至少需要2个字符' : '昵称不能超过32个字符';
      return false;
    }

    if (!email.value.trim()) {
      errors.value.email = '请输入邮箱';
      return false;
    }
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
      errors.value.email = '请输入正确的邮箱格式';
      return false;
    }

    const site = website.value.trim();
    if (site && !/^https?:\/\/.+/.test(site)) {
      errors.value.website = '网站地址格式不正确';
      return false;
    }
  }

  if (!commentContent.value.trim()) {
    errors.value.content = '请输入评论内容';
    return false;
  }

  return true;
};

// 清除特定字段的错误
const clearError = (field: 'nickname' | 'email' | 'website' | 'content') => {
  errors.value[field] = '';
};

/**
 * 处理文本框输入事件
 */
const handleTextareaInput = () => {
  clearError('content');
  resetTextareaHeight();
};

// 事件处理
const handleSubmitComment = async () => {
  if (!validateForm()) return;

  isSubmitting.value = true;
  try {
    if (pendingImages.value.size > 0) {
      isUploading.value = true;
      try {
        await uploadPendingImages();
      } catch (error: any) {
        info(error.message || '图片上传失败');
        return;
      } finally {
        isUploading.value = false;
      }
    }

    const content = commentContent.value.trim();
    const guestInfo = !isLoggedIn.value
      ? {
          nickname: nickname.value.trim(),
          email: email.value.trim(),
          website: website.value.trim() || undefined,
        }
      : undefined;

    if (isReplyMode.value && props.commentId) {
      await context.addReply(props.commentId, content, guestInfo);
      // 回复成功后关闭回复输入框
      context.replyState.cancelReply();
    } else {
      await context.addComment(content, guestInfo);
    }

    if (!isLoggedIn.value) {
      localStorage.setItem(
        'guest_info',
        JSON.stringify({
          nickname: nickname.value.trim(),
          email: email.value.trim(),
          website: website.value.trim(),
        })
      );
    }

    commentContent.value = '';
    resetToDefaultHeight();
    success('评论发表成功');

    if (isLoggedIn.value) triggerOnComment();
  } catch (error: any) {
    errors.value.email = error.message || error.response?.data?.message || '评论发表失败';
  } finally {
    isSubmitting.value = false;
  }
};

const handleLogin = () => context.showLogin();
const handleMainAction = () => (shouldShowSend.value ? handleSubmitComment() : handleLogin());
const handleSecondaryAction = (event: Event) => {
  event.stopPropagation();
  showExpandedBtn.value = false;
  isUserInfoFilled.value ? handleLogin() : handleSubmitComment();
};
const toggleExpandedBtn = (event: Event) => {
  event.stopPropagation();
  showExpandedBtn.value = !showExpandedBtn.value;
};
const togglePreview = () => (showPreview.value = !showPreview.value);

const handleCancelReply = () => {
  context.replyState.cancelReply();
  cleanupPendingImages();
  commentContent.value = '';
  resetTextareaHeight();
};

// 图片上传相关
const isUploading = ref(false);
const pendingImages = ref<Map<string, File>>(new Map());

// 表情选择器相关
const showEmojiPicker = ref(false);
const emojiButtonRef = ref<HTMLElement | null>(null);
const emojiPickerRef = ref<HTMLElement | null>(null);

const handleImageUpload = () => fileInputRef.value?.click();
const handleFileSelect = (event: Event) => {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (file) {
    insertImagePlaceholder(file);
    (event.target as HTMLInputElement).value = '';
  }
};

const insertImagePlaceholder = (file: File) => {
  const error = validateFile(file, '评论贴图');
  if (error) {
    info(error);
    return;
  }

  const blobUrl = URL.createObjectURL(file);
  pendingImages.value.set(blobUrl, file);

  const textarea = textareaRef.value;
  if (!textarea) return;

  const { selectionStart, selectionEnd } = textarea;
  const imageMarkdown = `![图片](${blobUrl})`;

  commentContent.value =
    commentContent.value.substring(0, selectionStart) +
    imageMarkdown +
    commentContent.value.substring(selectionEnd);

  nextTick(() => {
    textarea.focus();
    const newPosition = selectionStart + imageMarkdown.length;
    textarea.setSelectionRange(newPosition, newPosition);
    resetTextareaHeight();
  });
};

const handlePaste = (event: ClipboardEvent) => {
  const file = Array.from(event.clipboardData?.items || [])
    .find(item => item.type.startsWith('image/'))
    ?.getAsFile();

  if (file) {
    event.preventDefault();
    insertImagePlaceholder(file);
  }
};

// 表情选择
const toggleEmojiPicker = () => {
  showEmojiPicker.value = !showEmojiPicker.value;
};

const handleEmojiSelect = (emoji: string) => {
  const textarea = textareaRef.value;
  if (!textarea) return;

  const { selectionStart, selectionEnd } = textarea;

  commentContent.value =
    commentContent.value.substring(0, selectionStart) +
    emoji +
    commentContent.value.substring(selectionEnd);

  nextTick(() => {
    textarea.focus();
    const newPosition = selectionStart + emoji.length;
    textarea.setSelectionRange(newPosition, newPosition);
    resetTextareaHeight();
  });

  showEmojiPicker.value = false;
};

const uploadPendingImages = async () => {
  if (pendingImages.value.size === 0) return;

  const uploads = Array.from(pendingImages.value.entries()).map(async ([blobUrl, file]) => {
    const result = await uploadFile(file, '评论贴图');
    return { blobUrl, realUrl: result.file_url };
  });

  const results = await Promise.all(uploads);

  results.forEach(({ blobUrl, realUrl }) => {
    commentContent.value = commentContent.value.replace(blobUrl, realUrl);
    URL.revokeObjectURL(blobUrl);
  });

  pendingImages.value.clear();
};

const cleanupPendingImages = () => {
  pendingImages.value.forEach((_, blobUrl) => URL.revokeObjectURL(blobUrl));
  pendingImages.value.clear();
};

// 实时保存评论内容到本地存储
watch(commentContent, newContent => {
  if (newContent) {
    localStorage.setItem('comment_draft', newContent);
  } else {
    localStorage.removeItem('comment_draft');
  }
});

// 点击外部关闭扩展按钮
onClickOutside(buttonGroupRef, () => {
  showExpandedBtn.value = false;
});

// 点击外部关闭表情选择器
onClickOutside(
  emojiButtonRef,
  () => {
    showEmojiPicker.value = false;
  },
  {
    ignore: [emojiPickerRef],
  }
);

// 组件卸载时清理 blob URL
onUnmounted(() => {
  cleanupPendingImages();
});
</script>

<template>
  <div class="comment-input" :class="{ 'reply-mode': isReplyMode }">
    <template v-if="!isLoggedIn">
      <div class="user-info-row">
        <div class="input-wrapper">
          <input
            v-model="nickname"
            type="text"
            placeholder="昵称 *"
            :disabled="isSubmitting"
            :class="{ error: errors.nickname }"
            @input="clearError('nickname')"
          />
          <transition name="fade">
            <div v-if="errors.nickname" class="error-tooltip">{{ errors.nickname }}</div>
          </transition>
        </div>
        <div class="input-wrapper">
          <input
            v-model="email"
            type="email"
            placeholder="邮箱 *"
            :disabled="isSubmitting"
            :class="{ error: errors.email }"
            @input="clearError('email')"
          />
          <transition name="fade">
            <div v-if="errors.email" class="error-tooltip">{{ errors.email }}</div>
          </transition>
        </div>
        <div class="input-wrapper policy-input">
          <input
            v-model="website"
            type="url"
            placeholder="网址"
            :disabled="isSubmitting"
            :class="{ error: errors.website }"
            @input="clearError('website')"
          />
          <div class="guest-policy-tip">
            <button
              type="button"
              class="guest-policy-trigger"
              aria-label="游客评论信息说明"
              title="游客评论信息说明"
            >
              <i class="ri-information-line"></i>
            </button>
            <div class="guest-policy-tooltip" role="note">
              <p v-for="line in guestPrivacyNotice" :key="line">{{ line }}</p>
            </div>
          </div>
          <transition name="fade">
            <div v-if="errors.website" class="error-tooltip">{{ errors.website }}</div>
          </transition>
        </div>
      </div>
    </template>

    <div class="editor-container">
      <textarea
        ref="textareaRef"
        v-model="commentContent"
        placeholder="写下你的评论...支持 Markdown 语法"
        rows="3"
        maxlength="500"
        :disabled="isSubmitting"
        :class="{ error: errors.content }"
        data-lenis-prevent
        @input="handleTextareaInput"
        @paste="handlePaste"
      />
      <transition name="fade">
        <div v-if="errors.content" class="error-tooltip content-error">{{ errors.content }}</div>
      </transition>
      <!-- 字符数提示，超过450字后显示 -->
      <div
        v-if="commentContent.length > 450"
        class="char-count"
        :class="{ 'near-limit': commentContent.length > 480 }"
      >
        {{ commentContent.length }}/500
      </div>
      <transition name="expand">
        <div
          v-if="showPreview"
          class="preview-area markdown-body"
          v-html="renderedMarkdown || '<p class=\'empty-hint\'>暂无内容</p>'"
        ></div>
      </transition>
    </div>

    <div class="toolbar">
      <div class="toolbar-left">
        <div v-if="isReplyMode" class="reply-tag">
          <span class="reply-tag-text">回复 {{ replyTo }}</span>
          <button
            class="reply-tag-close"
            @click="handleCancelReply"
            :disabled="isSubmitting"
            aria-label="取消回复"
          >
            <i class="ri-close-line"></i>
          </button>
        </div>
        <div class="emoji-wrapper">
          <button
            ref="emojiButtonRef"
            class="tool-btn"
            @click="toggleEmojiPicker"
            title="表情"
            aria-label="插入表情"
            :disabled="isSubmitting || isUploading"
            :class="{ active: showEmojiPicker }"
          >
            <i class="ri-emotion-line"></i>
          </button>
        </div>
        <transition name="fade-scale">
          <FeaturesCommentEmojiPicker
            v-if="showEmojiPicker"
            ref="emojiPickerRef"
            class="emoji-picker-portal"
            @select="handleEmojiSelect"
          />
        </transition>
        <button
          class="tool-btn"
          @click="handleImageUpload"
          title="图片"
          aria-label="上传图片"
          :disabled="isSubmitting || isUploading"
          :class="{ uploading: isUploading }"
        >
          <i :class="isUploading ? 'ri-loader-4-line rotating' : 'ri-image-line'"></i>
        </button>
        <input
          ref="fileInputRef"
          type="file"
          accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
          style="display: none"
          @change="handleFileSelect"
        />
        <button
          class="tool-btn"
          @click="togglePreview"
          :title="showPreview ? '编辑' : 'Markdown预览'"
          :aria-label="showPreview ? '切换到编辑模式' : '切换到预览模式'"
          :class="{ active: showPreview }"
          :disabled="isSubmitting || isUploading"
        >
          <i :class="showPreview ? 'ri-edit-line' : 'ri-eye-line'"></i>
        </button>
      </div>
      <div ref="buttonGroupRef" class="button-group">
        <button
          class="submit-btn main-btn"
          @click="handleMainAction"
          :disabled="isSubmitting"
          :aria-label="mainBtn.text"
        >
          <i :class="mainBtn.icon"></i>{{ mainBtn.text }}
        </button>
        <template v-if="!isLoggedIn">
          <button
            class="submit-btn expand-btn"
            @click="toggleExpandedBtn"
            :disabled="isSubmitting"
            aria-label="更多选项"
          >
            <i class="ri-more-2-fill"></i>
          </button>
          <transition name="slide-fade">
            <button
              v-if="showExpandedBtn"
              class="submit-btn secondary-btn"
              @click="handleSecondaryAction"
              :disabled="isSubmitting"
              :aria-label="secondaryBtn.text"
            >
              <i :class="secondaryBtn.icon"></i>{{ secondaryBtn.text }}
            </button>
          </transition>
        </template>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@keyframes rotating {
  to {
    transform: rotate(360deg);
  }
}

.rotating {
  animation: rotating 1s linear infinite;
}

.comment-input {
  margin-bottom: 30px;
  background: var(--flec-card-bg);
  border-radius: 8px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;

  &.reply-mode {
    background: var(--flec-heavy-bg);
    margin-bottom: 0;
  }
}

input,
textarea {
  border: 1px solid transparent;
  background-color: transparent;
  border-radius: 6px;
  color: var(--font-color);
  outline: none;
  transition: border-color 0.2s;
  padding: 8px 12px;
  font-size: 0.9rem;

  &:focus {
    border-color: var(--theme-color);
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  &.error {
    border-color: #ef4444;

    &:focus {
      border-color: #ef4444;
    }
  }
}

.user-info-row {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.guest-policy-tip {
  position: relative;
  display: inline-flex;
}

.policy-input {
  input {
    padding-right: 40px;
  }

  .guest-policy-tip {
    position: absolute;
    top: 50%;
    right: 10px;
    transform: translateY(-50%);
    z-index: 2;
  }
}

.guest-policy-trigger {
  width: 24px;
  height: 24px;
  padding: 0;
  border: 1px solid var(--flec-border-color);
  border-radius: 999px;
  background: transparent;
  color: var(--theme-meta-color);
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: help;
  transition:
    color 0.2s ease,
    border-color 0.2s ease,
    background-color 0.2s ease;

  i {
    font-size: 0.9rem;
  }

  &:hover,
  &:focus-visible {
    color: var(--theme-color);
    border-color: rgba(73, 177, 245, 0.35);
    background: rgba(73, 177, 245, 0.08);
  }
}

.guest-policy-tooltip {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  width: min(360px, calc(100vw - 48px));
  padding: 10px 12px;
  border-radius: 10px;
  background: rgba(20, 24, 31, 0.96);
  color: rgba(255, 255, 255, 0.92);
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.18);
  font-size: 0.76rem;
  line-height: 1.65;
  opacity: 0;
  visibility: hidden;
  transform: translateY(-4px);
  transition:
    opacity 0.2s ease,
    transform 0.2s ease,
    visibility 0.2s ease;
  z-index: 20;

  &::before {
    content: '';
    position: absolute;
    right: 8px;
    bottom: 100%;
    border: 6px solid transparent;
    border-bottom-color: rgba(20, 24, 31, 0.96);
  }

  p {
    margin: 0;
  }

  p + p {
    margin-top: 6px;
  }
}

.guest-policy-tip:hover .guest-policy-tooltip,
.guest-policy-tip:focus-within .guest-policy-tooltip {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
}

.input-wrapper {
  position: relative;
}

.error-tooltip {
  position: absolute;
  bottom: calc(100% + 8px);
  left: 50%;
  transform: translateX(-50%);
  padding: 6px 12px;
  background: rgba(0, 0, 0, 0.85);
  color: white;
  font-size: 0.75rem;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  z-index: 10;
  white-space: nowrap;

  &::after {
    content: '';
    position: absolute;
    top: 100%;
    left: 50%;
    transform: translateX(-50%);
    border: 4px solid transparent;
    border-top-color: rgba(0, 0, 0, 0.85);
  }

  &.content-error {
    left: 20px;
    transform: translateX(0);

    &::after {
      left: 20px;
      transform: translateX(0);
    }
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(4px);
}

.editor-container {
  position: relative;
  width: 100%;
}

// 字符数提示样式
.char-count {
  position: absolute;
  bottom: 12px;
  right: 12px;
  font-size: 0.75rem;
  color: var(--theme-meta-color);
  padding: 6px;
  border-radius: 4px;
  pointer-events: none;

  &.near-limit {
    color: #ef4444;
    font-weight: 500;
  }
}

textarea {
  width: 100%;
  padding: 10px 12px;
  font-size: 0.95rem;
  line-height: 1.6;
  resize: none;
  min-height: 60px;
  max-height: 300px;
  overflow-y: auto;
  // 隐藏滚动条但保持可滚动
  scrollbar-width: none;
  -ms-overflow-style: none;

  &::-webkit-scrollbar {
    display: none;
  }

  &.error {
    border-color: #ef4444;

    &:focus {
      border-color: #ef4444;
    }
  }
}

.preview-area {
  width: 100%;
  padding: 10px 12px;
  font-size: 0.95rem;
  line-height: 1.6;
  min-height: 80px;
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 6px;
  background-color: rgba(0, 0, 0, 0.02);
  font-family: inherit;

  .empty-hint {
    color: var(--theme-meta-color);
    font-style: italic;
    margin: 0;
  }
}

// 展开动画
.expand-enter-active,
.expand-leave-active {
  transition: all 0.3s ease;
  overflow: hidden;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
  margin-top: 0;
  padding-top: 0;
  padding-bottom: 0;
}

.expand-enter-to,
.expand-leave-from {
  opacity: 1;
  max-height: 300px;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;

  &-left {
    display: flex;
    align-items: center;
    gap: 6px;
    position: relative;
  }
}

.emoji-wrapper {
  position: relative;
}

.emoji-picker-portal {
  position: absolute;
  bottom: calc(100% + 8px);
  left: 0;
  z-index: 100;
}

.fade-scale {
  &-enter-active,
  &-leave-active {
    transition: all 0.2s ease;
  }

  &-enter-from,
  &-leave-to {
    opacity: 0;
    transform: translateY(-8px) scale(0.95);
  }
}

:deep(.emoji-picker) {
  position: absolute;
  top: calc(100% + 8px);
  left: 0;
  z-index: 100;
}

.reply-tag {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 8px 4px 12px;
  background: var(--flec-card-bg);
  border-radius: 6px;
  margin-right: 8px;
  border: 1px solid var(--flec-border-color);

  &-text {
    font-size: 0.85rem;
    color: var(--theme-color);
  }

  &-close {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.2s;
    padding: 0;
    border: none;
    background: transparent;
    color: var(--theme-meta-color);

    i {
      font-size: 1rem;
    }

    &:hover:not(:disabled) {
      background: var(--flec-card-bg);
      color: var(--font-color);
    }

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }
}

.tool-btn {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
  color: var(--font-color);

  i {
    font-size: 1.1rem;
  }

  &:hover:not(:disabled) {
    color: var(--theme-color);
  }

  &.active {
    color: var(--theme-color);
    background: rgba(73, 177, 245, 0.1);
  }

  &.uploading {
    color: var(--theme-color);
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

.button-group {
  display: flex;
  gap: 4px;
  position: relative;
}

.submit-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  border: none;
  border-radius: 6px;
  background: var(--flec-btn);
  color: var(--font-light-color);
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;

  i {
    font-size: 1rem;
  }

  &:hover:not(:disabled) {
    background: var(--flec-btn-hover);
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  &.main-btn {
    padding: 6px 20px;
    min-width: 90px;
  }

  &.expand-btn {
    padding: 6px 8px;

    i {
      font-size: 1.1rem;
    }
  }

  &.secondary-btn {
    padding: 6px 16px;
    position: absolute;
    right: 0;
    top: calc(100% + 4px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 10;
  }
}

.slide-fade {
  &-enter-active,
  &-leave-active {
    transition: all 0.2s;
  }

  &-enter-from {
    opacity: 0;
    transform: translateY(-8px);
  }

  &-leave-to {
    opacity: 0;
    transform: translateY(-4px);
  }
}

@media screen and (max-width: 768px) {
  .comment-input {
    padding: 12px;
    gap: 8px;
  }

  .user-info-row {
    grid-template-columns: 1fr;
  }

  .guest-policy-tooltip {
    right: -4px;
    width: min(320px, calc(100vw - 32px));
    font-size: 0.74rem;
  }

  .error-tooltip {
    font-size: 0.7rem;
    padding: 5px 10px;
    max-width: 90%;
  }

  input {
    font-size: 0.875rem;
  }

  textarea {
    font-size: 0.9rem;
  }

  .preview-area {
    font-size: 0.9rem;
  }

  .toolbar {
    flex-direction: column;
    gap: 10px;

    &-left {
      width: 100%;
      flex-wrap: wrap;
    }
  }

  .reply-tag {
    &-text {
      font-size: 0.8rem;
    }
  }

  .button-group {
    width: 100%;
  }

  .submit-btn {
    &.main-btn {
      flex: 1;
      justify-content: center;
      min-width: auto;
    }

    &.expand-btn {
      flex-shrink: 0;
    }

    &.secondary-btn {
      width: calc(100% - 4px);
    }
  }
}
</style>
