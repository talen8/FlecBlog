/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<Record<string, unknown>, Record<string, unknown>, unknown>
  export default component
}

interface FileItem {
  name: string
  path: string
  isDirectory: boolean
}

interface FileOpenResult {
  path: string
  name: string
  content: string
}

interface Window {
  electronAPI: {
    minimize: () => void
    maximize: () => void
    close: () => void
    getVersion: () => Promise<string>
    isMaximized: () => Promise<boolean>
    onMaximizedChanged: (callback: (maximized: boolean) => void) => () => void
    onNavigate: (callback: (path: string) => void) => () => void
    onMenuAction: (callback: (action: string) => void) => () => void
  }
  fileAPI: {
    openFile: () => Promise<FileOpenResult | null>
    saveFile: (path: string, content: string) => Promise<boolean>
    saveFileDialog: (content: string) => Promise<string | null>
    openFolder: () => Promise<string | null>
    readDir: (dirPath: string) => Promise<FileItem[]>
    readFile: (filePath: string) => Promise<string>
    writeFile: (path: string, content: string) => Promise<boolean>
    createFile: (dirPath: string, name: string) => Promise<string | null>
    createDir: (dirPath: string, name: string) => Promise<string | null>
    deleteFile: (path: string) => Promise<boolean>
  }
  storeAPI: {
    get: <T = unknown>(key: string) => Promise<T>
    set: (key: string, value: unknown) => Promise<void>
    getAll: () => Promise<Record<string, unknown>>
  }
}

interface ImportMetaEnv {
  readonly VITE_API_URL: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
