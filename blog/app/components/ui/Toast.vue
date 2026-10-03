<script setup lang="ts">
export interface ToastProps {
  message: string;
  type?: 'success' | 'error' | 'warning' | 'info';
  show?: boolean;
}

withDefaults(defineProps<ToastProps>(), {
  type: 'info',
  show: false,
});
</script>

<template>
  <Teleport to="body">
    <Transition name="toast">
      <div v-if="show" :class="['toast', `toast-${type}`]">
        {{ message }}
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
.toast {
  position: fixed;
  top: 80px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10000;
  padding: 12px 20px;
  border-radius: 6px;
  font-size: 0.95rem;
  color: white;
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

  &.toast-success {
    background: rgba(76, 175, 80, 0.8);
  }

  &.toast-error {
    background: rgba(244, 67, 54, 0.8);
  }

  &.toast-warning {
    background: rgba(255, 152, 0, 0.8);
  }

  &.toast-info {
    background: rgba(59, 130, 246, 0.8);
  }
}

// 动画
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(-50%) translateY(-20px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-10px);
}

@media screen and (max-width: 768px) {
  .toast {
    left: 16px;
    right: 16px;
    transform: none;
  }

  .toast-enter-from {
    transform: translateY(-20px);
  }

  .toast-leave-to {
    transform: translateY(-10px);
  }
}
</style>
