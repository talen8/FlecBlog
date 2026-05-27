<template>
  <div class="settings-page">
    <div class="settings-sidebar">
      <div class="settings-title">
        <button class="back-btn" @click="$emit('close')" title="返回">
          <i class="ri-arrow-left-line"></i>
        </button>
        <span>设置</span>
      </div>
      <nav class="settings-nav">
        <button
          v-for="section in sections"
          :key="section.key"
          :class="{ active: activeSection === section.key }"
          @click="activeSection = section.key"
        >
          <i :class="section.icon"></i>
          <span>{{ section.label }}</span>
        </button>
      </nav>
    </div>

    <div class="settings-content">
      <!-- 通用 -->
      <div v-if="activeSection === 'general'" class="section">
        <h2>通用</h2>
        <div class="setting-group">
          <div class="setting-row">
            <div class="setting-label">
              <span>启动时</span>
              <span class="desc">应用启动时的行为</span>
            </div>
            <select v-model="settings.launchBehavior">
              <option value="empty">打开空白编辑器</option>
              <option value="last">恢复上次打开的文件</option>
              <option value="folder">打开上次的文件夹</option>
            </select>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>自动保存</span>
              <span class="desc">文件修改后实时保存到本地</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.autoSave" />
              <span class="slider"></span>
            </label>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>退出时确认</span>
              <span class="desc">有未保存修改时提示确认</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.confirmOnExit" />
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <!-- 编辑器 -->
      <div v-if="activeSection === 'editor'" class="section">
        <h2>编辑器</h2>
        <div class="setting-group">
          <div class="setting-row">
            <div class="setting-label">
              <span>字体大小</span>
              <span class="desc">编辑器字体大小（像素）</span>
            </div>
            <input type="number" v-model.number="settings.fontSize" min="10" max="32" class="number-input" />
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>字体</span>
              <span class="desc">编辑器使用的等宽字体</span>
            </div>
            <select v-model="settings.fontFamily">
              <option value="default">默认</option>
              <option value="'JetBrains Mono', monospace">JetBrains Mono</option>
              <option value="'Fira Code', monospace">Fira Code</option>
              <option value="'Cascadia Code', monospace">Cascadia Code</option>
              <option value="'Source Code Pro', monospace">Source Code Pro</option>
              <option value="Consolas, monospace">Consolas</option>
            </select>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>行高</span>
              <span class="desc">行间距倍数</span>
            </div>
            <select v-model="settings.lineHeight">
              <option :value="1.2">紧凑 (1.2)</option>
              <option :value="1.5">标准 (1.5)</option>
              <option :value="1.8">宽松 (1.8)</option>
              <option :value="2.0">双倍 (2.0)</option>
            </select>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>Tab 大小</span>
              <span class="desc">一个 Tab 对应的空格数</span>
            </div>
            <select v-model.number="settings.tabSize">
              <option :value="2">2</option>
              <option :value="4">4</option>
              <option :value="8">8</option>
            </select>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>自动换行</span>
              <span class="desc">超出编辑器宽度时自动换行</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.wordWrap" />
              <span class="slider"></span>
            </label>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>显示行号</span>
              <span class="desc">编辑器左侧显示行号</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.lineNumbers" />
              <span class="slider"></span>
            </label>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>拼写检查</span>
              <span class="desc">启用浏览器拼写检查</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.spellCheck" />
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <!-- 外观 -->
      <div v-if="activeSection === 'appearance'" class="section">
        <h2>外观</h2>
        <div class="setting-group">
          <div class="setting-row">
            <div class="setting-label">
              <span>主题</span>
              <span class="desc">编辑器颜色主题</span>
            </div>
            <select v-model="settings.theme">
              <option value="dark">深色</option>
              <option value="light">浅色</option>
              <option value="system">跟随系统</option>
            </select>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>预览区字体大小</span>
              <span class="desc">Markdown 预览的字体大小</span>
            </div>
            <input type="number" v-model.number="settings.previewFontSize" min="12" max="24" class="number-input" />
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>显示侧边栏</span>
              <span class="desc">启动时是否显示侧边栏</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.showSidebar" />
              <span class="slider"></span>
            </label>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>显示工具栏</span>
              <span class="desc">是否显示格式工具栏</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.showToolbar" />
              <span class="slider"></span>
            </label>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>显示状态栏</span>
              <span class="desc">是否显示底部状态栏</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.showStatusBar" />
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <!-- 文件 -->
      <div v-if="activeSection === 'file'" class="section">
        <h2>文件</h2>
        <div class="setting-group">
          <div class="setting-row">
            <div class="setting-label">
              <span>默认保存格式</span>
              <span class="desc">新建文件的默认格式</span>
            </div>
            <select v-model="settings.defaultFormat">
              <option value=".md">Markdown (.md)</option>
              <option value=".txt">纯文本 (.txt)</option>
            </select>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>文件编码</span>
              <span class="desc">读写文件使用的编码</span>
            </div>
            <select v-model="settings.encoding">
              <option value="utf-8">UTF-8</option>
              <option value="gbk">GBK</option>
              <option value="utf-16">UTF-16</option>
            </select>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>行尾符号</span>
              <span class="desc">文件换行符类型</span>
            </div>
            <select v-model="settings.lineEnding">
              <option value="lf">LF (Unix/macOS)</option>
              <option value="crlf">CRLF (Windows)</option>
              <option value="auto">自动检测</option>
            </select>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>插入尾部换行</span>
              <span class="desc">保存时确保文件末尾有换行</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.insertFinalNewline" />
              <span class="slider"></span>
            </label>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>去除尾部空格</span>
              <span class="desc">保存时去除行尾多余空格</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.trimTrailingWhitespace" />
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <!-- Markdown -->
      <div v-if="activeSection === 'markdown'" class="section">
        <h2>Markdown</h2>
        <div class="setting-group">
          <div class="setting-row">
            <div class="setting-label">
              <span>代码高亮主题</span>
              <span class="desc">预览区代码块的高亮配色</span>
            </div>
            <select v-model="settings.codeTheme">
              <option value="one-dark">One Dark</option>
              <option value="github">GitHub</option>
              <option value="monokai">Monokai</option>
            </select>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>启用数学公式</span>
              <span class="desc">渲染 LaTeX 数学公式</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.enableMath" />
              <span class="slider"></span>
            </label>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>启用 Mermaid 图表</span>
              <span class="desc">渲染 Mermaid 流程图等</span>
            </div>
            <label class="switch">
              <input type="checkbox" v-model="settings.enableMermaid" />
              <span class="slider"></span>
            </label>
          </div>
        </div>
      </div>

      <!-- 连接 -->
      <div v-if="activeSection === 'connection'" class="section">
        <h2>连接</h2>
        <div class="setting-group">
          <div class="setting-row">
            <div class="setting-label">
              <span>后端 API 地址</span>
              <span class="desc">FlecBlog 后端地址，用于同步博客文章</span>
            </div>
            <div class="input-with-btn">
              <input type="text" v-model="apiConfig.url" placeholder="http://localhost:8080/api/v1" />
              <button class="btn-small" @click="handleSaveApiUrl">保存</button>
            </div>
          </div>
          <div class="setting-row">
            <div class="setting-label">
              <span>连接状态</span>
            </div>
            <span class="status" :class="hasApiUrl() ? 'ok' : 'fail'">
              {{ hasApiUrl() ? '已配置' : '未配置' }}
            </span>
          </div>
        </div>
      </div>

      <!-- 快捷键 -->
      <div v-if="activeSection === 'keybindings'" class="section">
        <h2>快捷键</h2>
        <div class="setting-group">
          <div class="shortcut-list">
            <div class="shortcut-item" v-for="s in shortcuts" :key="s.action">
              <span class="action">{{ s.label }}</span>
              <kbd>{{ s.key }}</kbd>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { settings, apiConfig, hasApiUrl, saveApiUrl as storeSaveApiUrl, saveSettings } from '@/composables/useStore'

defineEmits<{ close: [] }>()

const activeSection = ref('general')

const sections = [
  { key: 'general', label: '通用', icon: 'ri-settings-3-line' },
  { key: 'editor', label: '编辑器', icon: 'ri-code-line' },
  { key: 'appearance', label: '外观', icon: 'ri-palette-line' },
  { key: 'file', label: '文件', icon: 'ri-file-text-line' },
  { key: 'markdown', label: 'Markdown', icon: 'ri-markdown-line' },
  { key: 'connection', label: '连接', icon: 'ri-server-line' },
  { key: 'keybindings', label: '快捷键', icon: 'ri-keyboard-line' },
]

const shortcuts = [
  { action: 'new', label: '新建文件', key: 'Ctrl+N' },
  { action: 'open', label: '打开文件', key: 'Ctrl+O' },
  { action: 'save', label: '保存', key: 'Ctrl+S' },
  { action: 'saveAs', label: '另存为', key: 'Ctrl+Shift+S' },
  { action: 'bold', label: '粗体', key: 'Ctrl+B' },
  { action: 'italic', label: '斜体', key: 'Ctrl+I' },
  { action: 'find', label: '查找', key: 'Ctrl+F' },
  { action: 'sidebar', label: '切换侧边栏', key: 'Ctrl+B' },
  { action: 'settings', label: '打开设置', key: 'Ctrl+,' },
  { action: 'fullscreen', label: '全屏', key: 'F11' },
  { action: 'zoomIn', label: '放大', key: 'Ctrl+=' },
  { action: 'zoomOut', label: '缩小', key: 'Ctrl+-' },
]

// 设置变更时自动持久化
watch(settings, () => saveSettings(), { deep: true })

async function handleSaveApiUrl() {
  const url = apiConfig.url.trim()
  if (!url) {
    ElMessage.warning('请输入 API 地址')
    return
  }
  try {
    new URL(url)
  } catch {
    ElMessage.error('地址格式不正确')
    return
  }
  await storeSaveApiUrl(url)
  ElMessage.success('已保存')
}
</script>

<style scoped lang="scss">
.settings-page {
  display: flex;
  height: 100%;
  background: #1e1e1e;
}

.settings-sidebar {
  width: 220px;
  background: #252526;
  border-right: 1px solid #333;
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
}

.settings-title {
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 8px;
  border-bottom: 1px solid #333;
  font-size: 14px;
  font-weight: 500;
  color: #e0e0e0;

  .back-btn {
    border: none;
    background: transparent;
    color: #888;
    cursor: pointer;
    padding: 4px;
    border-radius: 4px;
    font-size: 16px;
    display: flex;
    align-items: center;
    &:hover { background: var(--border-primary); color: var(--text-regular); }
  }
}

.settings-nav {
  flex: 1;
  padding: 8px;
  overflow-y: auto;

  button {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    border: none;
    background: transparent;
    color: var(--text-regular);
    font-size: 13px;
    cursor: pointer;
    border-radius: 4px;
    transition: all 0.1s;
    text-align: left;

    i { font-size: 16px; color: var(--text-secondary); }

    &:hover { background: var(--bg-elevated); }

    &.active {
      background: var(--bg-active);
      color: #fff;
      i { color: #fff; }
    }
  }
}

.settings-content {
  flex: 1;
  overflow-y: auto;
  padding: 0;
}

.section {
  padding: 24px 32px;
  max-width: 700px;

  h2 {
    font-size: 20px;
    font-weight: 500;
    color: var(--text-primary);
    margin-bottom: 24px;
  }
}

.setting-group {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 0;
  border-bottom: 1px solid var(--bg-elevated);
  gap: 24px;

  &:last-child { border-bottom: none; }
}

.setting-label {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;

  span {
    font-size: 13px;
    color: var(--text-primary);
  }

  .desc {
    font-size: 12px;
    color: var(--text-muted);
  }
}

select {
  height: 32px;
  padding: 0 10px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-primary);
  border-radius: 4px;
  color: var(--text-primary);
  font-size: 13px;
  min-width: 160px;
  &:focus { border-color: var(--accent-blue); outline: none; }
}

.number-input {
  width: 80px;
  height: 32px;
  padding: 0 10px;
  background: var(--bg-elevated);
  border: 1px solid var(--border-primary);
  border-radius: 4px;
  color: var(--text-primary);
  font-size: 13px;
  text-align: center;
  &:focus { border-color: var(--accent-blue); outline: none; }
}

.input-with-btn {
  display: flex;
  gap: 8px;

  input {
    flex: 1;
    height: 32px;
    padding: 0 10px;
    background: var(--bg-elevated);
    border: 1px solid var(--border-primary);
    border-radius: 4px;
    color: var(--text-primary);
    font-size: 13px;
    min-width: 240px;
    &:focus { border-color: var(--accent-blue); outline: none; }
  }
}

.btn-small {
  height: 32px;
  padding: 0 14px;
  background: var(--accent-blue);
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 13px;
  cursor: pointer;
  white-space: nowrap;
  &:hover { background: var(--accent-blue-hover); }
}

/* Toggle switch */
.switch {
  position: relative;
  display: inline-block;
  width: 40px;
  height: 22px;
  flex-shrink: 0;

  input {
    opacity: 0;
    width: 0;
    height: 0;
  }

  .slider {
    position: absolute;
    cursor: pointer;
    inset: 0;
    background: var(--border-primary);
    border-radius: 22px;
    transition: 0.2s;

    &::before {
      content: '';
      position: absolute;
      height: 16px;
      width: 16px;
      left: 3px;
      bottom: 3px;
      background: var(--text-regular);
      border-radius: 50%;
      transition: 0.2s;
    }
  }

  input:checked + .slider {
    background: var(--accent-blue);
  }

  input:checked + .slider::before {
    transform: translateX(18px);
    background: #fff;
  }
}

.status {
  font-size: 13px;
  &.ok { color: var(--accent-green); }
  &.fail { color: var(--accent-red); }
}

.shortcut-list {
  display: flex;
  flex-direction: column;
  width: 100%;
}

.shortcut-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid var(--bg-elevated);

  .action {
    font-size: 13px;
    color: var(--text-regular);
  }

  kbd {
    background: var(--bg-elevated);
    border: 1px solid var(--border-primary);
    border-radius: 4px;
    padding: 2px 8px;
    font-size: 12px;
    color: var(--text-regular);
    font-family: monospace;
  }
}
</style>
