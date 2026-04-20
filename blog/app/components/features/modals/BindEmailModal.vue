<script setup lang="ts">
import { updateUserProfile } from '@/composables/api/user';

const props = defineProps<{
  modelValue: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
  (e: 'success'): void;
}>();

const { success: showSuccess, error: showError } = useToast();
const { onSkip } = useBindEmail();

const email = ref('');
const loading = ref(false);
const emailError = ref('');

// 邮箱验证
const validateEmail = (val: string) => {
  if (!val.trim()) return '请输入邮箱';
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(val)) return '请输入正确的邮箱格式';
  return '';
};

// 提交绑定
const handleSubmit = async () => {
  emailError.value = validateEmail(email.value);
  if (emailError.value) return;

  loading.value = true;
  try {
    await updateUserProfile({
      email: email.value.trim(),
    });
    showSuccess('邮箱绑定成功');
    emit('success');
    emit('update:modelValue', false);
  } catch (error: unknown) {
    showError((error as Error).message || '绑定失败');
  } finally {
    loading.value = false;
  }
};

// 稍后提醒（记录跳过时间）
const handleRemindLater = () => {
  onSkip();
  emit('update:modelValue', false);
};

// 重置表单
watch(
  () => props.modelValue,
  val => {
    if (val) {
      email.value = '';
      emailError.value = '';
    }
  }
);
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="modal-overlay" @click.self="handleRemindLater">
        <div class="modal-container">
          <button class="close-btn" :disabled="loading" @click="handleRemindLater">
            <i class="ri-close-line" />
          </button>

          <div class="modal-header">
            <i class="ri-mail-send-line header-icon" />
            <h3>绑定邮箱</h3>
          </div>

          <div class="modal-body">
            <p class="description">绑定真实邮箱后，您可以：</p>
            <ul class="benefits">
              <li><i class="ri-notification-3-line" />及时接收评论回复通知</li>
              <li><i class="ri-lock-password-line" />使用邮箱+密码登录</li>
              <li><i class="ri-key-2-line" />找回密码时使用</li>
            </ul>

            <form class="bind-form" @submit.prevent="handleSubmit">
              <div class="form-group">
                <input
                  v-model="email"
                  type="email"
                  class="form-input"
                  :class="{ error: emailError }"
                  placeholder="请输入您的邮箱"
                  :disabled="loading"
                />
                <p v-if="emailError" class="error-message">{{ emailError }}</p>
              </div>

              <div class="form-actions">
                <button
                  type="button"
                  class="btn-secondary"
                  :disabled="loading"
                  @click="handleRemindLater"
                >
                  稍后再说
                </button>
                <button type="submit" class="btn-primary" :disabled="loading">
                  {{ loading ? '绑定中...' : '确认绑定' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.modal-container {
  position: relative;
  width: 100%;
  max-width: 400px;
  background: var(--flec-card-bg);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.2);
  overflow: hidden;
}

.close-btn {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--theme-meta-color);
  transition: all 0.2s;

  &:hover:not(:disabled) {
    background: #8080801a;
    color: var(--font-color);
  }

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  i {
    font-size: 20px;
  }
}

.modal-header {
  padding: 24px 24px 16px;
  text-align: center;

  .header-icon {
    font-size: 48px;
    color: var(--theme-color);
    margin-bottom: 12px;
    display: block;
  }

  h3 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--font-color);
  }
}

.modal-body {
  padding: 0 24px 24px;

  .description {
    margin: 0 0 12px;
    font-size: 0.95rem;
    color: var(--theme-meta-color);
  }

  .benefits {
    margin: 0 0 20px;
    padding: 0;
    list-style: none;

    li {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 8px 0;
      font-size: 0.9rem;
      color: var(--font-color);

      i {
        font-size: 16px;
        color: var(--theme-color);
      }
    }
  }
}

.bind-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-input {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid var(--flec-border);
  border-radius: 8px;
  background: var(--flec-card-bg);
  color: var(--font-color);
  font-size: 0.95rem;
  transition: all 0.2s;

  &:focus {
    outline: none;
    border-color: var(--theme-color);
    box-shadow: 0 0 0 3px #49b1f526;
  }

  &.error {
    border-color: #e57373;
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  &::placeholder {
    color: var(--theme-meta-color);
  }
}

.error-message {
  margin: 0;
  font-size: 0.85rem;
  color: #e57373;
}

.form-actions {
  display: flex;
  gap: 12px;

  button {
    flex: 1;
    padding: 12px 16px;
    border-radius: 8px;
    font-size: 0.95rem;
    cursor: pointer;
    transition: all 0.2s;

    &:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
  }

  .btn-secondary {
    border: 1px solid var(--flec-border);
    background: transparent;
    color: var(--theme-meta-color);

    &:hover:not(:disabled) {
      background: #8080800d;
    }
  }

  .btn-primary {
    border: none;
    background: var(--theme-color);
    color: white;

    &:hover:not(:disabled) {
      opacity: 0.9;
    }
  }
}

// 响应式设计
@media screen and (max-width: 480px) {
  .modal-overlay {
    padding: 0.5rem;
    align-items: flex-start;
  }

  .modal-container {
    max-width: 100%;
    margin-top: 0.5rem;
  }

  .modal-header {
    padding: 20px 20px 14px;

    .header-icon {
      font-size: 40px;
      margin-bottom: 10px;
    }

    h3 {
      font-size: 1.15rem;
    }
  }

  .modal-body {
    padding: 0 20px 20px;

    .description {
      font-size: 0.9rem;
    }

    .benefits li {
      font-size: 0.85rem;
      padding: 6px 0;
    }

    .form-input {
      padding: 10px 12px;
      font-size: 0.9rem;
    }

    .form-actions button {
      padding: 10px 14px;
      font-size: 0.9rem;
    }
  }
}

// 过渡动画
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;

  .modal-container {
    transition: transform 0.2s ease;
  }
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;

  .modal-container {
    transform: scale(0.95) translateY(10px);
  }
}
</style>
