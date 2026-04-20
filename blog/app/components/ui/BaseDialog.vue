<script setup lang="ts">
interface Props {
  modelValue: boolean;
  title?: string;
  confirmText?: string;
  loading?: boolean;
  closeOnClickOutside?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  title: '',
  confirmText: '确定',
  loading: false,
  closeOnClickOutside: true,
});

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  confirm: [];
}>();

const handleClose = () => {
  if (props.loading || !props.closeOnClickOutside) return;
  emit('update:modelValue', false);
};

const handleConfirm = () => {
  emit('confirm');
};
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="dialog-overlay" @click.self="handleClose">
        <div class="dialog-container">
          <!-- 头部 -->
          <div v-if="title" class="dialog-header">
            <h3 class="dialog-title">{{ title }}</h3>
            <button class="dialog-close" :disabled="loading" @click="handleClose">
              <i class="ri-close-line" />
            </button>
          </div>

          <!-- 内容 -->
          <div class="dialog-body">
            <slot />
          </div>

          <!-- 底部 -->
          <div v-if="confirmText" class="dialog-footer">
            <button class="btn btn-primary" :disabled="loading" @click="handleConfirm">
              <i v-if="loading" class="ri-loader-4-line loading" />
              <span>{{ loading ? '处理中..' : confirmText }}</span>
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 20px;
}

.dialog-container {
  background-color: var(--flec-card-bg);
  border-radius: 12px;
  width: 100%;
  max-width: var(--dialog-width, 480px);
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.dialog-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  border-bottom: 1px solid var(--flec-border);

  .dialog-title {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--font-color);
    flex: 1;
  }

  .dialog-close {
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: var(--theme-meta-color);
    cursor: pointer;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.2s;
    flex-shrink: 0;

    i {
      font-size: 20px;
    }

    &:hover:not(:disabled) {
      background: var(--flec-heavy-bg);
      color: var(--font-color);
    }

    &:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  }
}

.dialog-body {
  padding: 24px;
  overflow-y: auto;
  flex: 1;
}

.dialog-footer {
  padding: 16px 24px;
  border-top: 1px solid var(--flec-border);
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

// 动画
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;

  .dialog-container {
    transition: transform 0.3s ease;
  }
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;

  .dialog-container {
    transform: scale(0.95) translateY(-20px);
  }
}

@media (max-width: 768px) {
  .dialog-overlay {
    padding: 0;
  }

  .dialog-container {
    max-width: 100% !important;
    max-height: 100vh;
    border-radius: 0;
  }

  .dialog-header {
    padding: 16px 20px;
  }

  .dialog-body {
    padding: 20px;
  }

  .dialog-footer {
    padding: 12px 20px;
  }
}

// 按钮样式
.btn {
  padding: 10px 20px;
  border-radius: 8px;
  font-size: 0.95rem;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  border: none;

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.btn-primary {
  background: var(--theme-color);
  color: var(--font-light-color);

  &:hover:not(:disabled) {
    opacity: 0.9;
  }
}

.loading {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
