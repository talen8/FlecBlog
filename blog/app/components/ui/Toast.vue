<script setup lang="ts">
const { toasts } = useToast();

// 语义图标：与全局提示的语义映射保持一致
const ICONS: Record<ToastType, string> = {
  success: 'ri-checkbox-circle-line',
  error: 'ri-error-warning-line',
  warning: 'ri-alarm-warning-line',
  info: 'ri-information-line',
};
</script>

<template>
  <Teleport to="body">
    <TransitionGroup name="toast" tag="div" class="toast" appear>
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="toast-item"
        :class="`toast-${toast.type}`"
      >
        <i :class="ICONS[toast.type]" />
        <span>{{ toast.message }}</span>
      </div>
    </TransitionGroup>
  </Teleport>
</template>

<style lang="scss" scoped>
.toast {
  position: fixed;
  top: 80px;
  left: 0;
  right: 0;
  z-index: 10000;
  pointer-events: none;
}

.toast-item {
  --toast-tint: var(--toast-info);

  display: flex;
  align-items: center;
  gap: 8px;
  width: max-content;
  max-width: min(90vw, 420px);
  margin-inline: auto;
  margin-bottom: 8px;
  padding: 10px 16px;
  border: 1px solid transparent;
  border-radius: 6px;
  font-size: 0.9rem;
  pointer-events: auto;
  background: color-mix(in srgb, var(--toast-tint) 15%, var(--flec-heavy-bg, #f5f5f5));
  border-color: color-mix(in srgb, var(--toast-tint) 18%, var(--flec-border, #e0e0e0));
  color: var(--toast-tint);

  &.toast-success {
    --toast-tint: var(--toast-success);
  }

  &.toast-error {
    --toast-tint: var(--toast-error);
  }

  &.toast-warning {
    --toast-tint: var(--toast-warning);
  }

  &.toast-info {
    --toast-tint: var(--toast-info);
  }
}

.toast-enter-active,
.toast-leave-active {
  transition:
    opacity 0.3s ease,
    transform 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

.toast-leave-active {
  pointer-events: none;
}

.toast-move {
  transition: transform 0.3s ease;
}

@media screen and (max-width: 768px) {
  .toast {
    left: 16px;
    right: 16px;
  }

  .toast-item {
    max-width: 100%;
  }
}
</style>
