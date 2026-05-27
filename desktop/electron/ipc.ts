import { BrowserWindow, ipcMain, app, dialog } from 'electron'
import { readFile, writeFile, readdir, stat, mkdir } from 'fs/promises'
import { join, basename } from 'path'
import type { FileItem } from '../src/types'
import store from './store'

export function registerIpcHandlers(win: BrowserWindow): void {
  // 窗口控制
  ipcMain.on('window:minimize', () => win.minimize())
  ipcMain.on('window:maximize', () => {
    win.isMaximized() ? win.unmaximize() : win.maximize()
  })
  ipcMain.on('window:close', () => win.close())
  ipcMain.handle('app:version', () => app.getVersion())
  ipcMain.handle('window:isMaximized', () => win.isMaximized())

  win.on('maximize', () => win.webContents.send('window:maximized-changed', true))
  win.on('unmaximize', () => win.webContents.send('window:maximized-changed', false))

  // 文件操作
  ipcMain.handle('file:open', async () => {
    const result = await dialog.showOpenDialog(win, {
      properties: ['openFile'],
      filters: [{ name: 'Markdown', extensions: ['md', 'markdown', 'txt'] }],
    })
    if (result.canceled || !result.filePaths[0]) return null
    const filePath = result.filePaths[0]
    const content = await readFile(filePath, 'utf-8')
    return { path: filePath, name: basename(filePath), content }
  })

  ipcMain.handle('file:save', async (_e, filePath: string, content: string) => {
    try {
      await writeFile(filePath, content, 'utf-8')
      return true
    } catch {
      return false
    }
  })

  ipcMain.handle('file:save-dialog', async (_e, content: string) => {
    const result = await dialog.showSaveDialog(win, {
      filters: [{ name: 'Markdown', extensions: ['md'] }],
    })
    if (result.canceled || !result.filePath) return null
    await writeFile(result.filePath, content, 'utf-8')
    return result.filePath
  })

  ipcMain.handle('file:open-folder', async () => {
    const result = await dialog.showOpenDialog(win, {
      properties: ['openDirectory'],
    })
    if (result.canceled || !result.filePaths[0]) return null
    return result.filePaths[0]
  })

  ipcMain.handle('file:read-dir', async (_e, dirPath: string): Promise<FileItem[]> => {
    try {
      const entries = await readdir(dirPath, { withFileTypes: true })
      const items: FileItem[] = []
      for (const entry of entries) {
        if (entry.name.startsWith('.')) continue
        items.push({
          name: entry.name,
          path: join(dirPath, entry.name),
          isDirectory: entry.isDirectory(),
        })
      }
      // 目录在前，文件在后，按名称排序
      items.sort((a, b) => {
        if (a.isDirectory !== b.isDirectory) return a.isDirectory ? -1 : 1
        return a.name.localeCompare(b.name)
      })
      return items
    } catch {
      return []
    }
  })

  ipcMain.handle('file:read', async (_e, filePath: string) => {
    try {
      return await readFile(filePath, 'utf-8')
    } catch {
      return ''
    }
  })

  ipcMain.handle('file:write', async (_e, filePath: string, content: string) => {
    try {
      await writeFile(filePath, content, 'utf-8')
      return true
    } catch {
      return false
    }
  })

  ipcMain.handle('file:create', async (_e, dirPath: string, name: string) => {
    try {
      const filePath = join(dirPath, name)
      await writeFile(filePath, '', 'utf-8')
      return filePath
    } catch {
      return null
    }
  })

  ipcMain.handle('file:mkdir', async (_e, dirPath: string, name: string) => {
    try {
      const newDir = join(dirPath, name)
      await mkdir(newDir, { recursive: true })
      return newDir
    } catch {
      return null
    }
  })

  ipcMain.handle('file:delete', async (_e, filePath: string) => {
    try {
      const { unlink, rmdir } = await import('fs/promises')
      const s = await stat(filePath)
      if (s.isDirectory()) {
        await rmdir(filePath, { recursive: true })
      } else {
        await unlink(filePath)
      }
      return true
    } catch {
      return false
    }
  })

  // 持久化存储
  ipcMain.handle('store:get', (_e, key: string) => {
    return store.get(key)
  })

  ipcMain.handle('store:set', (_e, key: string, value: unknown) => {
    store.set(key, value)
  })

  ipcMain.handle('store:get-all', () => {
    return store.store
  })
}
