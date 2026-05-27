<template>
  <div class="status-bar">
    <div class="left">
      <span class="item" :class="connectionClass">
        <i class="ri-checkbox-blank-circle-fill"></i>
        {{ connectionText }}
      </span>
      <span class="item" v-if="fileName">
        <template v-if="isReadOnly">
          <i class="ri-lock-line"></i> 只读
        </template>
        <template v-else>
          <i :class="isModified ? 'ri-edit-2-fill' : 'ri-check-line'"></i>
          {{ isModified ? '未保存' : '已保存' }}
        </template>
      </span>
      <span class="item">{{ charCount }} 字</span>
    </div>
    <div class="right">
      <button class="mode-btn" :class="{ active: editorMode === 'wysiwyg' }" @click="setEditorMode('wysiwyg')">
        <i class="ri-edit-2-line"></i> WYSIWYG
      </button>
      <button class="mode-btn" :class="{ active: editorMode === 'source' }" @click="setEditorMode('source')">
        <i class="ri-code-line"></i> 源码
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { charCount, fileName, isModified, isReadOnly, editorMode, setEditorMode } from '@/stores/editor'
import { hasApiUrl } from '@/composables/useStore'

const connectionClass = computed(() => hasApiUrl() ? 'connected' : 'disconnected')
const connectionText = computed(() => hasApiUrl() ? '已连接' : '本地模式')
</script>

<style scoped lang="scss">
.status-bar {
  height: 24px;
  background: var(--accent-blue);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 8px;
  font-size: 12px;
  flex-shrink: 0;
  user-select: none;
}

.left, .right { display: flex; align-items: center; gap: 12px; }

.item {
  display: flex; align-items: center; gap: 4px;
  i { font-size: 10px; }
  &.connected i { color: var(--accent-green); }
  &.disconnected i { color: var(--accent-red); }
}

.mode-btn {
  border: none;
  background: transparent;
  color: rgba(255,255,255,0.7);
  font-size: 11px;
  cursor: pointer;
  padding: 2px 8px;
  border-radius: 3px;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: all 0.15s;
  &:hover { background: rgba(255,255,255,0.15); color: #fff; }
  &.active { background: rgba(255,255,255,0.2); color: #fff; }
}
</style>
