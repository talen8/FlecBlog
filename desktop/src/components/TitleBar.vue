<template>
  <div class="title-bar">
    <div class="title-bar-drag">
      <span class="app-name">FlecEditor</span>
    </div>
    <div class="title-bar-center">
      <span class="file-name" v-if="fileName">{{ fileName }}<span v-if="isModified" class="modified-dot"> •</span></span>
    </div>
    <div class="title-bar-drag"></div>
    <div class="title-bar-controls">
      <button class="title-btn" @click="minimize"><i class="ri-subtract-line"></i></button>
      <button class="title-btn" @click="maximize">
        <i :class="isMaximized ? 'ri-checkbox-blank-line' : 'ri-checkbox-multiple-blank-line'"></i>
      </button>
      <button class="title-btn close" @click="close"><i class="ri-close-line"></i></button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import { fileName, isModified } from '@/stores/editor'

const isMaximized = ref(false)

const minimize = () => window.electronAPI.minimize()
const maximize = () => window.electronAPI.maximize()
const close = () => window.electronAPI.close()

let removeListener: (() => void) | undefined
onMounted(async () => {
  isMaximized.value = await window.electronAPI.isMaximized()
  removeListener = window.electronAPI.onMaximizedChanged((v) => { isMaximized.value = v })
})
onUnmounted(() => removeListener?.())
</script>

<style scoped lang="scss">
.title-bar {
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background-color: var(--bg-surface);
  -webkit-app-region: drag;
  flex-shrink: 0;
  user-select: none;
  position: relative;
}

.title-bar-drag {
  flex: 1;
  padding-left: 12px;
  .app-name {
    font-size: 12px;
    color: var(--text-muted);
  }
}

.title-bar-center {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  .file-name {
    font-size: 12px;
    color: var(--text-secondary);
    white-space: nowrap;
    .modified-dot { color: var(--accent-gold); }
  }
}

.title-bar-controls {
  display: flex;
  -webkit-app-region: no-drag;
}

.title-btn {
  width: 46px;
  height: 32px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.1s;
  &:hover { background: rgba(255,255,255,0.1); color: #fff; }
  &.close:hover { background: #e81123; color: #fff; }
}
</style>
