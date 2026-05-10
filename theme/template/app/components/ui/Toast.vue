<template>
  <Teleport to="body">
    <div class="toast-container">
      <div v-for="toast in toasts" :key="toast.id" class="toast-item" :class="toast.type">
        <i :class="getIcon(toast.type)" />
        <span>{{ toast.message }}</span>
      </div>
    </div>
  </Teleport>
</template>

<script lang="ts" setup>
const { toasts } = useToast();

const getIcon = (type: string) => {
  const icons: Record<string, string> = {
    success: 'ri-check-line',
    error: 'ri-error-warning-line',
    info: 'ri-information-line',
    warning: 'ri-alert-line',
  };
  return icons[type] || icons.info;
};
</script>

<style lang="scss" scoped>
.toast-container {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 2000;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.toast-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 20px;
  background: var(--flec-card-bg);
  border: 1px solid var(--flec-border);
  font-size: 0.9rem;

  &.success {
    color: #52c41a;
  }
  &.error {
    color: #ff4d4f;
  }
  &.info {
    color: var(--theme-color);
  }
  &.warning {
    color: #faad14;
  }
}
</style>
