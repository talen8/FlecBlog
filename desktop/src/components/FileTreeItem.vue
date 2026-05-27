<template>
  <div>
    <div
      class="tree-item"
      :style="{ paddingLeft: level * 12 + 8 + 'px' }"
      :class="{ active: isActive }"
      @click="handleClick"
      @contextmenu.prevent="$emit('contextmenu', $event)"
    >
      <i v-if="item.isDirectory" :class="expanded ? 'ri-folder-open-line' : 'ri-folder-3-line'" class="icon folder"></i>
      <i v-else class="ri-file-text-line icon file"></i>
      <span class="name">{{ item.name }}</span>
    </div>
    <div v-if="item.isDirectory && expanded && children.length">
      <FileTreeItem
        v-for="child in children"
        :key="child.path"
        :item="child"
        :level="level + 1"
        @open="(e) => $emit('open', e)"
        @contextmenu="(e) => $emit('contextmenu', e)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { currentFile } from '@/stores/editor'
import { useFile } from '@/composables/useFile'
import type { FileItem } from '@/types'

const props = defineProps<{ item: FileItem; level: number }>()
const emit = defineEmits<{ open: [item: FileItem]; contextmenu: [event: MouseEvent] }>()

const { readDir } = useFile()
const expanded = ref(false)
const children = ref<FileItem[]>([])

const isActive = computed(() => currentFile.value?.path === props.item.path)

async function handleClick() {
  if (props.item.isDirectory) {
    if (!expanded.value && children.value.length === 0) {
      children.value = await readDir(props.item.path)
    }
    expanded.value = !expanded.value
  } else {
    emit('open', props.item)
  }
}
</script>

<style scoped lang="scss">
.tree-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 3px 0;
  cursor: pointer;
  font-size: 13px;
  color: var(--text-regular);
  transition: background 0.1s;
  &:hover { background: var(--bg-hover); }
  &.active { background: var(--bg-active); }
}

.icon {
  font-size: 14px;
  &.folder { color: var(--accent-gold); }
  &.file { color: var(--text-secondary); }
}

.name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
