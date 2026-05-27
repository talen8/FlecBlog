import { app, BrowserWindow, shell, Menu, Tray, nativeImage } from 'electron'
import { join } from 'path'
import { is } from '@electron-toolkit/utils'
import { registerIpcHandlers } from './ipc'

let mainWindow: BrowserWindow | null = null
let tray: Tray | null = null

function createWindow(): void {
  mainWindow = new BrowserWindow({
    width: 1280,
    height: 800,
    minWidth: 960,
    minHeight: 600,
    show: false,
    frame: false,
    titleBarStyle: 'hidden',
    backgroundColor: '#1e1e1e',
    webPreferences: {
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false,
      contextIsolation: true,
      nodeIntegration: false,
    },
  })

  mainWindow.on('ready-to-show', () => mainWindow?.show())

  mainWindow.on('close', (event) => {
    if (tray && !is.dev) {
      event.preventDefault()
      mainWindow?.hide()
    }
  })

  mainWindow.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url)
    return { action: 'deny' }
  })

  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    mainWindow.loadURL(process.env['ELECTRON_RENDERER_URL'])
  } else {
    mainWindow.loadFile(join(__dirname, '../renderer/index.html'))
  }

  registerIpcHandlers(mainWindow)
}

function createTray(): void {
  const iconPath = join(__dirname, '../resources/icon.ico')
  const icon = nativeImage.createFromPath(iconPath)
  tray = new Tray(icon.resize({ width: 16, height: 16 }))

  const contextMenu = Menu.buildFromTemplate([
    { label: '显示主窗口', click: () => { mainWindow?.show(); mainWindow?.focus() } },
    { type: 'separator' },
    { label: '退出', click: () => { tray?.destroy(); app.quit() } },
  ])

  tray.setToolTip('FlecEditor')
  tray.setContextMenu(contextMenu)
  tray.on('double-click', () => { mainWindow?.show(); mainWindow?.focus() })
}

function sendMenuAction(action: string) {
  mainWindow?.webContents.send('menu-action', action)
}

function createMenu(): void {
  const template: Electron.MenuItemConstructorOptions[] = [
    {
      label: '文件',
      submenu: [
        { label: '新建', accelerator: 'CmdOrCtrl+N', click: () => sendMenuAction('new-file') },
        { label: '打开文件', accelerator: 'CmdOrCtrl+O', click: () => sendMenuAction('open-file') },
        { label: '打开文件夹', accelerator: 'CmdOrCtrl+Shift+O', click: () => sendMenuAction('open-folder') },
        { type: 'separator' },
        { label: '保存', accelerator: 'CmdOrCtrl+S', click: () => sendMenuAction('save-file') },
        { label: '另存为', accelerator: 'CmdOrCtrl+Shift+S', click: () => sendMenuAction('save-file-as') },
        { type: 'separator' },
        { label: '设置', accelerator: 'CmdOrCtrl+,', click: () => sendMenuAction('settings') },
        { type: 'separator' },
        { role: 'quit', label: '退出' },
      ],
    },
    {
      label: '编辑',
      submenu: [
        { role: 'undo', label: '撤销' },
        { role: 'redo', label: '重做' },
        { type: 'separator' },
        { role: 'cut', label: '剪切' },
        { role: 'copy', label: '复制' },
        { role: 'paste', label: '粘贴' },
        { role: 'selectAll', label: '全选' },
        { type: 'separator' },
        { role: 'find', label: '查找' },
      ],
    },
    {
      label: '视图',
      submenu: [
        { label: '所见即所得', accelerator: 'CmdOrCtrl+Shift+1', click: () => sendMenuAction('mode-wysiwyg') },
        { label: '源码模式', accelerator: 'CmdOrCtrl+Shift+2', click: () => sendMenuAction('mode-source') },
        { type: 'separator' },
        { label: '切换侧边栏', accelerator: 'CmdOrCtrl+B', click: () => sendMenuAction('toggle-sidebar') },
        { type: 'separator' },
        { role: 'resetZoom', label: '重置缩放' },
        { role: 'zoomIn', label: '放大' },
        { role: 'zoomOut', label: '缩小' },
        { type: 'separator' },
        { role: 'togglefullscreen', label: '全屏' },
        { role: 'toggleDevTools', label: '开发者工具' },
      ],
    },
    {
      label: '帮助',
      submenu: [
        { label: '关于', click: () => sendMenuAction('about') },
      ],
    },
  ]

  Menu.setApplicationMenu(Menu.buildFromTemplate(template))
}

app.whenReady().then(() => {
  createMenu()
  createWindow()
  createTray()

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') app.quit()
})
