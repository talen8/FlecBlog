<script setup lang="ts">
const props = defineProps<{ modelValue: boolean }>();
const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>();
const { themeConfig } = useTheme();
const methods = computed(() => {
  try {
    return JSON.parse(themeConfig.value.donation_methods || '[]');
  } catch {
    return [];
  }
});
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="props.modelValue" class="overlay" @click="emit('update:modelValue', false)">
        <div class="box" @click.stop>
          <div v-for="(m, i) in methods" :key="i" class="item">
            <img :src="m.qrcode" :alt="m.name" class="qr" />
            <div class="name">{{ m.name }}</div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
.overlay {
  position: fixed;
  inset: 0;
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}
.box {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 1.5rem;
  padding: 1.5rem 1.5rem 0.9rem;
  background: var(--flec-card-bg);
  border-radius: 1rem;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  animation: up 0.3s;
}
.item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
}
.name {
  line-height: 1rem;
  color: var(--font-color);
}
.qr {
  width: 160px;
  height: 160px;
  object-fit: contain;
}
@keyframes up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s;
  .box {
    transition: all 0.3s;
  }
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
  .box {
    transform: scale(0.9) translateY(20px);
  }
}
@media (max-width: 768px) {
  .box {
    padding: 1.5rem 1.2rem;
    max-width: 360px;
  }
  .qr {
    width: 150px;
    height: 150px;
  }
}
</style>
