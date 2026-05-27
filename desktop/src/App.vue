<template>
  <div class="app">
    <TitleBar />
    <div class="app-body">
      <Sidebar v-show="sidebarVisible" @open-settings="showSettings = true" />
      <div class="editor-area">
        <Toolbar />
        <div class="editor-content">
          <MilkdownProvider v-if="editorMode === 'wysiwyg'">
            <WysiwygEditor />
          </MilkdownProvider>
          <SourceEditor v-else />
        </div>
        <StatusBar />
      </div>
    </div>
    <Settings v-if="showSettings" @close="showSettings = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import TitleBar from '@/components/TitleBar.vue'
import Sidebar from '@/components/Sidebar.vue'
import Toolbar from '@/components/Toolbar.vue'
import StatusBar from '@/components/StatusBar.vue'
import WysiwygEditor from '@/components/Editor/WysiwygEditor.vue'
import SourceEditor from '@/components/Editor/SourceEditor.vue'
import { MilkdownProvider } from '@milkdown/vue'
import Settings from '@/components/Settings/Settings.vue'
import { sidebarVisible, currentFile, editorMode, setEditorMode } from '@/stores/editor'
import { settings } from '@/composables/useStore'
import { useFile } from '@/composables/useFile'

const { openFile, saveFile, saveFileAs, newFile } = useFile()
const showSettings = ref(false)

let removeMenuListener: (() => void) | undefined
// 自动保存：内容变化后立即静默保存
watch(
  () => currentFile.value?.content,
  (val) => {
    if (val !== undefined && settings.autoSave && currentFile.value?.path) {
      saveFile(true)
    }
  }
)

onMounted(() => {
  document.documentElement.setAttribute('data-theme', 'dark')
  newFile()

  removeMenuListener = window.electronAPI.onMenuAction((action) => {
    switch (action) {
      case 'new-file': newFile(); break
      case 'open-file': openFile(); break
      case 'save-file': saveFile(); break
      case 'save-file-as': saveFileAs(); break
      case 'toggle-sidebar': sidebarVisible.value = !sidebarVisible.value; break
      case 'mode-source': setEditorMode('source'); break
      case 'mode-wysiwyg': setEditorMode('wysiwyg'); break
      case 'settings': showSettings.value = true; break
    }
  })
})

onUnmounted(() => {
  removeMenuListener?.()
})
</script>

<style scoped lang="scss">
.app {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-base);
  color: var(--text-regular);
  overflow: hidden;
}

.app-body {
  flex: 1;
  display: flex;
  overflow: hidden;
}

.editor-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.editor-content {
  flex: 1;
  overflow: hidden;
}
</style>
