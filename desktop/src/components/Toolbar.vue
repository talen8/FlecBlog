<template>
  <div class="toolbar">
    <div class="toolbar-group">
      <button title="新建 (Ctrl+N)" @click="newFile"><i class="ri-file-add-line"></i></button>
      <button title="打开 (Ctrl+O)" @click="openFile"><i class="ri-folder-open-line"></i></button>
      <button title="保存 (Ctrl+S)" :class="{ disabled: isReadOnly }" @click="!isReadOnly && saveFile()"><i class="ri-save-line"></i></button>
    </div>
    <template v-if="!isReadOnly">
      <div class="separator"></div>
      <div class="toolbar-group">
        <button title="粗体 (Ctrl+B)" @click="insertFormat('bold')"><i class="ri-bold"></i></button>
        <button title="斜体 (Ctrl+I)" @click="insertFormat('italic')"><i class="ri-italic"></i></button>
        <button title="删除线" @click="insertFormat('strikethrough')"><i class="ri-strikethrough"></i></button>
        <button title="行内代码" @click="insertFormat('code')"><i class="ri-code-line"></i></button>
      </div>
      <div class="separator"></div>
      <div class="toolbar-group">
        <button title="标题" @click="insertFormat('heading')"><i class="ri-heading"></i></button>
        <button title="引用" @click="insertFormat('quote')"><i class="ri-double-quotes-l"></i></button>
        <button title="无序列表" @click="insertFormat('ul')"><i class="ri-list-unordered"></i></button>
        <button title="有序列表" @click="insertFormat('ol')"><i class="ri-list-ordered"></i></button>
        <button title="任务列表" @click="insertFormat('task')"><i class="ri-checkbox-line"></i></button>
      </div>
      <div class="separator"></div>
      <div class="toolbar-group">
        <button title="链接" @click="insertFormat('link')"><i class="ri-link"></i></button>
        <button title="图片" @click="insertFormat('image')"><i class="ri-image-line"></i></button>
        <button title="代码块" @click="insertFormat('codeblock')"><i class="ri-code-box-line"></i></button>
        <button title="表格" @click="insertFormat('table')"><i class="ri-table-line"></i></button>
        <button title="分割线" @click="insertFormat('hr')"><i class="ri-separator"></i></button>
      </div>
    </template>
    <div class="spacer"></div>
    <div class="toolbar-group">
      <button title="切换侧边栏" @click="toggleSidebar"><i class="ri-side-bar-line"></i></button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { toggleSidebar, isReadOnly } from '@/stores/editor'
import { useFile } from '@/composables/useFile'

const { newFile, openFile, saveFile } = useFile()

function insertFormat(type: string) {
  if (isReadOnly.value) return
  window.dispatchEvent(new CustomEvent('editor:insert', { detail: { type } }))
}
</script>

<style scoped lang="scss">
.toolbar {
  height: 36px;
  background: var(--bg-surface);
  border-bottom: 1px solid var(--border-secondary);
  display: flex;
  align-items: center;
  padding: 0 4px;
  flex-shrink: 0;
  gap: 2px;
}

.toolbar-group {
  display: flex;
  gap: 1px;
}

.separator {
  width: 1px;
  height: 20px;
  background: var(--border-primary);
  margin: 0 4px;
}

.spacer { flex: 1; }

button {
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  color: var(--text-secondary);
  font-size: 16px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.1s;
  &:hover { background: var(--bg-elevated); color: var(--text-primary); }
  &.disabled {
    opacity: 0.4;
    cursor: not-allowed;
    &:hover { background: transparent; color: var(--text-secondary); }
  }
}
</style>
