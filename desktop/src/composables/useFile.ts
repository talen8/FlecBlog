import { ElMessage } from 'element-plus'
import { currentFile, content, markSaved, setContent } from '@/stores/editor'
import type { OpenFile } from '@/types'

export function useFile() {
  async function openFile(): Promise<void> {
    const result = await window.fileAPI.openFile()
    if (!result) return

    const file: OpenFile = {
      path: result.path,
      name: result.name,
      content: result.content,
      savedContent: result.content,
      isModified: false,
      source: 'local',
    }
    currentFile.value = file
    content.value = result.content
  }

  async function saveFile(silent = false): Promise<void> {
    if (!currentFile.value) return

    if (!currentFile.value.path) {
      if (!silent) await saveFileAs()
      return
    }

    const ok = await window.fileAPI.saveFile(currentFile.value.path, currentFile.value.content)
    if (ok) {
      markSaved()
      if (!silent) ElMessage.success('已保存')
    }
  }

  async function saveFileAs(): Promise<void> {
    if (!currentFile.value) return

    const result = await window.fileAPI.saveFileDialog(currentFile.value.content)
    if (result) {
      currentFile.value.path = result
      currentFile.value.name = result.split(/[/\\]/).pop() || ''
      markSaved()
      ElMessage.success('已保存')
    }
  }

  function newFile() {
    const file: OpenFile = {
      path: '',
      name: '未命名.md',
      content: '',
      savedContent: '',
      isModified: false,
      source: 'local',
    }
    currentFile.value = file
    content.value = ''
  }

  async function openFolder(): Promise<string | null> {
    return window.fileAPI.openFolder()
  }

  async function readDir(dirPath: string) {
    return window.fileAPI.readDir(dirPath)
  }

  async function readFile(filePath: string) {
    return window.fileAPI.readFile(filePath)
  }

  async function createFile(dirPath: string, name: string): Promise<string | null> {
    return window.fileAPI.createFile(dirPath, name)
  }

  async function createDir(dirPath: string, name: string): Promise<string | null> {
    return window.fileAPI.createDir(dirPath, name)
  }

  async function deleteFile(filePath: string): Promise<boolean> {
    return window.fileAPI.deleteFile(filePath)
  }

  return { openFile, saveFile, saveFileAs, newFile, openFolder, readDir, readFile, createFile, createDir, deleteFile }
}
