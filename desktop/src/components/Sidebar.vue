<template>
  <div class="sidebar">
    <div class="sidebar-tabs">
      <button
        :class="{ active: tab === 'files' }"
        @click="tab = 'files'"
        title="本地文件"
      >
        <i class="ri-folder-3-line"></i>
      </button>
      <button
        :class="{ active: tab === 'articles' }"
        @click="tab = 'articles'; loadArticles()"
        title="博客文章"
      >
        <i class="ri-article-line"></i>
      </button>
    </div>

    <div class="sidebar-content">
      <!-- 本地文件 -->
      <div v-if="tab === 'files'" class="file-panel">
        <div class="panel-header">
          <span>文件浏览器</span>
          <div class="header-actions">
            <button @click="startCreate('file')" title="新建文件"><i class="ri-file-add-line"></i></button>
            <button @click="startCreate('dir')" title="新建文件夹"><i class="ri-folder-add-line"></i></button>
            <button @click="openFolder" title="打开文件夹"><i class="ri-folder-open-line"></i></button>
          </div>
        </div>
        <div v-if="!currentFolder" class="empty-hint">
          <p>打开文件夹开始</p>
          <button @click="openFolder">打开文件夹</button>
        </div>
        <div v-else class="file-tree">
          <!-- 内联新建输入框 -->
          <div v-if="creating" class="inline-create" :style="{ paddingLeft: '8px' }">
            <i :class="creating.type === 'dir' ? 'ri-folder-3-line' : 'ri-file-text-line'" class="icon"></i>
            <input
              ref="createInput"
              v-model="creating.name"
              class="create-input"
              :placeholder="creating.type === 'dir' ? '文件夹名' : '文件名'"
              @keyup.enter="confirmCreate"
              @keyup.escape="cancelCreate"
              @blur="cancelCreate"
            />
          </div>
          <FileTreeItem
            v-for="item in fileTree"
            :key="item.path"
            :item="item"
            :level="0"
            @open="openFileFromTree"
            @contextmenu="onContextMenu($event, item)"
          />
        </div>
      </div>

      <!-- 博客文章 -->
      <div v-if="tab === 'articles'" class="article-panel">
        <div class="panel-header">
          <span>博客文章</span>
        </div>
        <div v-if="!isLoggedIn" class="empty-hint">
          <p>请先配置后端</p>
          <button @click="showLogin = true">登录</button>
        </div>
        <div v-else class="article-list">
          <div
            v-for="article in articles"
            :key="article.id"
            class="article-item"
            :class="{ active: currentFile?.articleId === article.id }"
            @click="openArticle(article)"
          >
            <span class="article-title">{{ article.title }}</span>
            <span class="article-status" :class="{ published: article.is_publish }">
              {{ article.is_publish ? '已发布' : '草稿' }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部设置按钮 -->
    <div class="sidebar-footer">
      <button class="settings-btn" @click="$emit('open-settings')" title="设置">
        <i class="ri-settings-3-line"></i>
        <span>设置</span>
      </button>
    </div>

    <!-- 右键菜单 -->
    <div
      v-if="contextMenu.visible"
      class="context-menu"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
    >
      <button @click="startCreate('file')"><i class="ri-file-add-line"></i> 新建文件</button>
      <button @click="startCreate('dir')"><i class="ri-folder-add-line"></i> 新建文件夹</button>
      <div class="context-divider"></div>
      <button @click="handleDelete" class="danger"><i class="ri-delete-bin-line"></i> 删除</button>
    </div>

    <!-- 登录对话框 -->
    <div v-if="showLogin" class="login-overlay" @click.self="showLogin = false">
      <div class="login-dialog">
        <h3>连接博客</h3>
        <div v-if="!hasApiUrl()" class="input-group">
          <input v-model="serverUrl" placeholder="后端地址 http://..." @keyup.enter="doLogin" />
        </div>
        <div class="input-group">
          <input v-model="loginForm.email" placeholder="邮箱" @keyup.enter="doLogin" />
        </div>
        <div class="input-group">
          <input v-model="loginForm.password" type="password" placeholder="密码" @keyup.enter="doLogin" />
        </div>
        <button class="btn-primary" @click="doLogin">登录</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, nextTick, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { currentFile, setContent } from '@/stores/editor'
import { hasApiUrl, saveApiUrl, lastFolder, saveLastFolder } from '@/composables/useStore'
import { useFile } from '@/composables/useFile'
import { getArticles, getArticle } from '@/api/article'
import { login, getProfile, setAccessToken } from '@/api/user'
import type { Article } from '@/types/article'
import type { OpenFile, SidebarTab, FileItem } from '@/types'
import FileTreeItem from './FileTreeItem.vue'

const { openFolder: openFolderDialog, readDir, readFile, createFile, createDir, deleteFile } = useFile()

defineEmits<{ 'open-settings': [] }>()

const tab = ref<SidebarTab>('files')
const currentFolder = ref<string | null>(null)
const fileTree = ref<FileItem[]>([])
const articles = ref<Article[]>([])
const isLoggedIn = ref(false)
const showLogin = ref(false)
const serverUrl = ref('')
const loginForm = reactive({ email: '', password: '' })

// 右键菜单
const contextMenu = reactive({ visible: false, x: 0, y: 0, item: null as FileItem | null })

// 内联新建
const creating = ref<{ type: 'file' | 'dir'; parentDir: string; name: string } | null>(null)
const createInput = ref<HTMLInputElement | null>(null)

function onContextMenu(e: MouseEvent, item: FileItem) {
  e.preventDefault()
  e.stopPropagation()
  contextMenu.visible = true
  contextMenu.x = e.clientX
  contextMenu.y = e.clientY
  contextMenu.item = item
}

function closeContextMenu() {
  contextMenu.visible = false
}

function startCreate(type: 'file' | 'dir') {
  closeContextMenu()
  let parentDir = currentFolder.value || ''
  if (contextMenu.item) {
    parentDir = contextMenu.item.isDirectory
      ? contextMenu.item.path
      : contextMenu.item.path.replace(/[/\\][^/\\]+$/, '')
  }
  creating.value = { type, parentDir, name: '' }
  nextTick(() => createInput.value?.focus())
}

async function confirmCreate() {
  if (!creating.value || !creating.value.name.trim()) {
    creating.value = null
    return
  }
  const { type, parentDir, name } = creating.value
  const result = type === 'file'
    ? await createFile(parentDir, name.trim())
    : await createDir(parentDir, name.trim())
  if (result) {
    ElMessage.success(type === 'file' ? '文件已创建' : '文件夹已创建')
    await refreshTree()
  } else {
    ElMessage.error('创建失败')
  }
  creating.value = null
}

function cancelCreate() {
  creating.value = null
}

async function handleDelete() {
  if (!contextMenu.item) return
  const item = contextMenu.item
  closeContextMenu()
  const ok = await deleteFile(item.path)
  if (ok) {
    ElMessage.success('已删除')
    if (currentFile.value?.path === item.path) {
      currentFile.value = null
    }
    await refreshTree()
  } else {
    ElMessage.error('删除失败')
  }
}

async function refreshTree() {
  if (currentFolder.value) {
    fileTree.value = await readDir(currentFolder.value)
  }
}

async function openFolder() {
  const folder = await openFolderDialog()
  if (!folder) return
  currentFolder.value = folder
  fileTree.value = await readDir(folder)
  saveLastFolder(folder)
}

async function openFileFromTree(item: FileItem) {
  if (item.isDirectory) return
  const content = await readFile(item.path)
  const file: OpenFile = {
    path: item.path,
    name: item.name,
    content,
    savedContent: content,
    isModified: false,
    source: 'local',
  }
  currentFile.value = file
  setContent(content)
}

async function loadArticles() {
  if (!hasApiUrl() || isLoggedIn.value) return
  try {
    await getProfile()
    isLoggedIn.value = true
    const data = await getArticles({ page: 1, page_size: 100 })
    articles.value = data.list
  } catch {
    isLoggedIn.value = false
  }
}

async function openArticle(article: Article) {
  try {
    const detail = await getArticle(article.id)
    const file: OpenFile = {
      path: '',
      name: detail.title + '.md',
      content: detail.content,
      savedContent: detail.content,
      isModified: false,
      source: 'api',
      articleId: detail.id,
      readOnly: true,
    }
    currentFile.value = file
    setContent(detail.content)
  } catch {
    ElMessage.error('加载文章失败')
  }
}

async function doLogin() {
  if (!hasApiUrl() && serverUrl.value) {
    try {
      new URL(serverUrl.value)
      await saveApiUrl(serverUrl.value)
    } catch {
      ElMessage.error('地址格式不正确')
      return
    }
  }
  try {
    const { access_token } = await login(loginForm)
    setAccessToken(access_token)
    isLoggedIn.value = true
    showLogin.value = false
    ElMessage.success('登录成功')
    const data = await getArticles({ page: 1, page_size: 100 })
    articles.value = data.list
  } catch {
    ElMessage.error('登录失败')
  }
}

onMounted(() => {
  document.addEventListener('click', closeContextMenu)
  // 恢复上次打开的文件夹
  if (lastFolder.value) {
    currentFolder.value = lastFolder.value
    readDir(lastFolder.value).then(tree => { fileTree.value = tree })
  }
})

onUnmounted(() => {
  document.removeEventListener('click', closeContextMenu)
})
</script>

<style scoped lang="scss">
.sidebar {
  width: 260px;
  background: var(--bg-surface);
  border-right: 1px solid var(--border-secondary);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.sidebar-tabs {
  display: flex;
  border-bottom: 1px solid var(--border-secondary);
  button {
    flex: 1;
    height: 36px;
    border: none;
    background: transparent;
    color: var(--text-muted);
    font-size: 18px;
    cursor: pointer;
    transition: all 0.15s;
    &:hover { color: var(--text-regular); background: var(--bg-elevated); }
    &.active { color: #fff; border-bottom: 2px solid var(--accent-blue); }
  }
}

.sidebar-content {
  flex: 1;
  overflow: auto;
}

.sidebar-footer {
  border-top: 1px solid var(--border-secondary);
  padding: 4px;
  flex-shrink: 0;

  .settings-btn {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 10px;
    border: none;
    background: transparent;
    color: var(--text-secondary);
    font-size: 12px;
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.15s;
    &:hover { background: var(--bg-elevated); color: var(--text-regular); }
    i { font-size: 14px; }
  }
}

.panel-header {
  padding: 8px 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 11px;
  text-transform: uppercase;
  color: var(--text-secondary);
  letter-spacing: 0.5px;

  .header-actions {
    display: flex;
    gap: 2px;
  }

  button {
    border: none;
    background: transparent;
    color: var(--text-secondary);
    cursor: pointer;
    font-size: 14px;
    padding: 2px;
    border-radius: 3px;
    &:hover { color: var(--text-regular); background: var(--bg-elevated); }
  }
}

.empty-hint {
  padding: 24px 16px;
  text-align: center;
  p { color: var(--text-muted); font-size: 13px; margin-bottom: 12px; }
  button {
    padding: 6px 16px;
    background: var(--accent-blue);
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 13px;
    &:hover { background: var(--accent-blue-hover); }
  }
}

.file-tree {
  padding: 4px 0;
}

.article-list {
  padding: 4px 0;
}

.article-item {
  padding: 6px 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
  transition: background 0.1s;
  &:hover { background: var(--bg-hover); }
  &.active { background: var(--bg-active); }
  .article-title {
    font-size: 13px;
    color: var(--text-regular);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
  }
  .article-status {
    font-size: 11px;
    color: var(--text-secondary);
    &.published { color: var(--accent-green); }
  }
}

.context-menu {
  position: fixed;
  background: var(--bg-elevated);
  border: 1px solid var(--border-primary);
  border-radius: 6px;
  padding: 4px;
  min-width: 160px;
  z-index: 200;
  box-shadow: 0 4px 12px rgba(0,0,0,0.4);

  button {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 10px;
    border: none;
    background: transparent;
    color: var(--text-regular);
    font-size: 13px;
    cursor: pointer;
    border-radius: 4px;
    text-align: left;
    &:hover { background: var(--border-primary); }
    &.danger { color: var(--accent-red); }
    i { font-size: 14px; }
  }

  .context-divider {
    height: 1px;
    background: var(--border-primary);
    margin: 4px 0;
  }
}

.inline-create {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 3px 0;
  font-size: 13px;

  .icon {
    font-size: 14px;
    color: var(--accent-gold);
    flex-shrink: 0;
  }

  .create-input {
    flex: 1;
    height: 22px;
    padding: 0 4px;
    background: var(--bg-base);
    border: 1px solid var(--accent-blue);
    border-radius: 3px;
    color: var(--text-primary);
    font-size: 13px;
    outline: none;
  }
}

.login-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.login-dialog {
  background: var(--bg-elevated);
  border: 1px solid var(--border-primary);
  border-radius: 8px;
  padding: 24px;
  width: 300px;
  h3 { color: var(--text-primary); margin-bottom: 16px; font-size: 16px; }
  .input-group {
    margin-bottom: 10px;
    input {
      width: 100%;
      height: 36px;
      padding: 0 10px;
      background: var(--bg-base);
      border: 1px solid var(--border-primary);
      border-radius: 4px;
      color: var(--text-primary);
      font-size: 13px;
      box-sizing: border-box;
      &:focus { border-color: var(--accent-blue); outline: none; }
    }
  }
  .btn-primary {
    width: 100%;
    height: 36px;
    background: var(--accent-blue);
    color: #fff;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    margin-top: 4px;
    &:hover { background: var(--accent-blue-hover); }
  }
}
</style>
