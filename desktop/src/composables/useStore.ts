import { reactive, toRaw } from 'vue'

export interface AppSettings {
  launchBehavior: 'empty' | 'last' | 'folder'
  autoSave: boolean
  confirmOnExit: boolean
  fontSize: number
  fontFamily: string
  lineHeight: number
  tabSize: number
  wordWrap: boolean
  lineNumbers: boolean
  spellCheck: boolean
  theme: 'dark' | 'light' | 'system'
  previewFontSize: number
  showSidebar: boolean
  showToolbar: boolean
  showStatusBar: boolean
  defaultFormat: string
  encoding: string
  lineEnding: string
  insertFinalNewline: boolean
  trimTrailingWhitespace: boolean
  codeTheme: string
  enableMath: boolean
  enableMermaid: boolean
}

const defaultSettings: AppSettings = {
  launchBehavior: 'empty',
  autoSave: false,
  confirmOnExit: true,
  fontSize: 14,
  fontFamily: 'default',
  lineHeight: 1.5,
  tabSize: 2,
  wordWrap: true,
  lineNumbers: true,
  spellCheck: false,
  theme: 'dark',
  previewFontSize: 15,
  showSidebar: true,
  showToolbar: true,
  showStatusBar: true,
  defaultFormat: '.md',
  encoding: 'utf-8',
  lineEnding: 'lf',
  insertFinalNewline: true,
  trimTrailingWhitespace: false,
  codeTheme: 'one-dark',
  enableMath: true,
  enableMermaid: true,
}

export const settings = reactive<AppSettings>({ ...defaultSettings })

export const apiConfig = reactive({ url: '', token: '' })

export const lastFolder = { value: null as string | null }

export async function loadPersistedData(): Promise<void> {
  const all = await window.storeAPI.getAll()
  if (all.settings) Object.assign(settings, all.settings)
  if (all.api) Object.assign(apiConfig, all.api)
  if (all.lastFolder !== undefined) lastFolder.value = all.lastFolder as string | null
}

export async function saveSettings(): Promise<void> {
  await window.storeAPI.set('settings', toRaw(settings))
}

export async function saveApiUrl(url: string): Promise<void> {
  apiConfig.url = url.replace(/\/+$/, '')
  await window.storeAPI.set('api', { ...toRaw(apiConfig) })
}

export async function saveApiToken(token: string): Promise<void> {
  apiConfig.token = token
  await window.storeAPI.set('api', { ...toRaw(apiConfig) })
}

export async function saveLastFolder(folder: string | null): Promise<void> {
  lastFolder.value = folder
  await window.storeAPI.set('lastFolder', folder)
}

export function getApiUrl(): string {
  return apiConfig.url
}

export function hasApiUrl(): boolean {
  return apiConfig.url !== ''
}

export function getApiToken(): string {
  return apiConfig.token
}
