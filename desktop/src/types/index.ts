export interface FileItem {
  name: string
  path: string
  isDirectory: boolean
  children?: FileItem[]
}

export interface OpenFile {
  path: string
  name: string
  content: string
  savedContent: string
  isModified: boolean
  source: 'local' | 'api'
  articleId?: number
  readOnly?: boolean
}

export type SidebarTab = 'files' | 'articles'
