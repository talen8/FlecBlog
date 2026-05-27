import Store from 'electron-store'

interface StoreSchema {
  settings: {
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
  api: {
    url: string
    token: string
  }
  lastFolder: string | null
}

const store = new Store<StoreSchema>({
  defaults: {
    settings: {
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
    },
    api: {
      url: '',
      token: '',
    },
    lastFolder: null,
  },
})

export default store
