import { contextBridge, ipcRenderer } from 'electron'

const electronAPI = {
  minimize: () => ipcRenderer.send('window:minimize'),
  maximize: () => ipcRenderer.send('window:maximize'),
  close: () => ipcRenderer.send('window:close'),
  getVersion: () => ipcRenderer.invoke('app:version') as Promise<string>,
  isMaximized: () => ipcRenderer.invoke('window:isMaximized') as Promise<boolean>,
  onMaximizedChanged: (callback: (maximized: boolean) => void) => {
    const handler = (_event: Electron.IpcRendererEvent, maximized: boolean): void => { callback(maximized) }
    ipcRenderer.on('window:maximized-changed', handler)
    return () => { ipcRenderer.removeListener('window:maximized-changed', handler) }
  },
  onNavigate: (callback: (path: string) => void) => {
    const handler = (_event: Electron.IpcRendererEvent, path: string): void => { callback(path) }
    ipcRenderer.on('navigate', handler)
    return () => { ipcRenderer.removeListener('navigate', handler) }
  },
  onMenuAction: (callback: (action: string) => void) => {
    const handler = (_event: Electron.IpcRendererEvent, action: string): void => { callback(action) }
    ipcRenderer.on('menu-action', handler)
    return () => { ipcRenderer.removeListener('menu-action', handler) }
  },
}

const fileAPI = {
  openFile: () => ipcRenderer.invoke('file:open') as Promise<{ path: string; name: string; content: string } | null>,
  saveFile: (path: string, content: string) => ipcRenderer.invoke('file:save', path, content) as Promise<boolean>,
  saveFileDialog: (content: string) => ipcRenderer.invoke('file:save-dialog', content) as Promise<string | null>,
  openFolder: () => ipcRenderer.invoke('file:open-folder') as Promise<string | null>,
  readDir: (dirPath: string) => ipcRenderer.invoke('file:read-dir', dirPath) as Promise<{ name: string; path: string; isDirectory: boolean }[]>,
  readFile: (filePath: string) => ipcRenderer.invoke('file:read', filePath) as Promise<string>,
  writeFile: (filePath: string, content: string) => ipcRenderer.invoke('file:write', filePath, content) as Promise<boolean>,
  createFile: (dirPath: string, name: string) => ipcRenderer.invoke('file:create', dirPath, name) as Promise<string | null>,
  createDir: (dirPath: string, name: string) => ipcRenderer.invoke('file:mkdir', dirPath, name) as Promise<string | null>,
  deleteFile: (filePath: string) => ipcRenderer.invoke('file:delete', filePath) as Promise<boolean>,
}

const storeAPI = {
  get: <T = unknown>(key: string) => ipcRenderer.invoke('store:get', key) as Promise<T>,
  set: (key: string, value: unknown) => ipcRenderer.invoke('store:set', key, value) as Promise<void>,
  getAll: () => ipcRenderer.invoke('store:get-all') as Promise<Record<string, unknown>>,
}

contextBridge.exposeInMainWorld('electronAPI', electronAPI)
contextBridge.exposeInMainWorld('fileAPI', fileAPI)
contextBridge.exposeInMainWorld('storeAPI', storeAPI)
