import { ref, computed } from 'vue'
import type { OpenFile } from '@/types'

export type EditorMode = 'wysiwyg' | 'source'

export const editorMode = ref<EditorMode>('wysiwyg')
export const currentFile = ref<OpenFile | null>(null)
export const sidebarVisible = ref(true)
export const charCount = ref(0)

export const isModified = computed(() => currentFile.value?.isModified ?? false)
export const fileName = computed(() => currentFile.value?.name ?? '')
export const filePath = computed(() => currentFile.value?.path ?? '')
export const fileSource = computed(() => currentFile.value?.source ?? 'local')
export const isReadOnly = computed(() => currentFile.value?.readOnly ?? false)

export const content = ref('')

export function setContent(value: string) {
  if (currentFile.value) {
    content.value = value
    currentFile.value.content = value
    currentFile.value.isModified = value !== currentFile.value.savedContent
  }
}

export function markSaved() {
  if (currentFile.value) {
    currentFile.value.savedContent = currentFile.value.content
    currentFile.value.isModified = false
  }
}

export function setEditorMode(mode: EditorMode) {
  editorMode.value = mode
}

export function toggleSidebar() {
  sidebarVisible.value = !sidebarVisible.value
}
