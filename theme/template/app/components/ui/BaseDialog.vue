<template>
  <Teleport to="body">
    <div v-if="modelValue" class="dialog-overlay" @click.self="handleClose">
      <div class="dialog-container" :style="{ '--dialog-width': dialogWidth }">
        <div class="dialog-header">
          <h3 class="dialog-title">{{ title }}</h3>
          <button class="dialog-close" @click="handleClose">
            <i class="ri-close-line" />
          </button>
        </div>
        <div class="dialog-content">
          <slot />
        </div>
        <div v-if="showFooter" class="dialog-footer">
          <button class="btn-secondary" @click="handleClose">取消</button>
          <button class="btn-primary" :disabled="loading" @click="handleConfirm">
            {{ loading ? '处理中...' : confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script lang="ts" setup>
const props = withDefaults(
  defineProps<{
    modelValue: boolean;
    title: string;
    confirmText?: string;
    loading?: boolean;
    showFooter?: boolean;
    closeOnClickOutside?: boolean;
  }>(),
  {
    confirmText: '确认',
    loading: false,
    showFooter: true,
    closeOnClickOutside: true,
  }
);

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  confirm: [];
}>();

const dialogWidth = computed(() => (props.showFooter ? '480px' : '400px'));

const handleClose = () => {
  if (props.closeOnClickOutside && !props.loading) {
    emit('update:modelValue', false);
  }
};

const handleConfirm = () => {
  emit('confirm');
};
</script>

<style lang="scss" scoped>
.dialog-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
}

.dialog-container {
  width: var(--dialog-width, 480px);
  max-width: 90vw;
  max-height: 90vh;
  background: var(--flec-card-bg);
  border: 1px solid var(--flec-border);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--flec-border);
}

.dialog-title {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
}

.dialog-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  color: var(--theme-meta-color);
  background: none;
  border: none;
  cursor: pointer;
}

.dialog-content {
  padding: 20px;
  overflow-y: auto;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 20px;
  border-top: 1px solid var(--flec-border);
}

.btn-primary,
.btn-secondary {
  padding: 8px 20px;
  font-size: 0.9rem;
  cursor: pointer;
}

.btn-primary {
  background: var(--theme-color);
  color: #fff;
  border: none;
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.btn-secondary {
  background: var(--flec-card-bg);
  color: var(--font-color);
  border: 1px solid var(--flec-border);
}
</style>
