import MarkdownIt from 'markdown-it'
import anchor from 'markdown-it-anchor'
// @ts-ignore - 没有类型定义
import taskLists from 'markdown-it-task-lists'
// @ts-ignore - 没有类型定义
import mark from 'markdown-it-mark'
// @ts-ignore - 没有类型定义
import linkAttributes from 'markdown-it-link-attributes'
// @ts-ignore - 没有类型定义
import kbd from 'markdown-it-kbd'
// @ts-ignore - 没有类型定义
import sub from 'markdown-it-sub'
// @ts-ignore - 没有类型定义
import sup from 'markdown-it-sup'
// @ts-ignore - 没有类型定义
import underline from 'markdown-it-plugin-underline'

import DOMPurify from 'isomorphic-dompurify'
import { getEmojiMapSync, replaceEmojisInText } from '@/composables/useEmojis'

// highlight.js 按需加载
import hljs from 'highlight.js/lib/core'
// 核心语言
import javascript from 'highlight.js/lib/languages/javascript'
import typescript from 'highlight.js/lib/languages/typescript'
import python from 'highlight.js/lib/languages/python'
import go from 'highlight.js/lib/languages/go'
import java from 'highlight.js/lib/languages/java'
import sql from 'highlight.js/lib/languages/sql'
// Web 相关
import xml from 'highlight.js/lib/languages/xml' // 包含 HTML
import css from 'highlight.js/lib/languages/css'
import json from 'highlight.js/lib/languages/json'
import yaml from 'highlight.js/lib/languages/yaml'
import markdown from 'highlight.js/lib/languages/markdown'
// DevOps / Shell
import bash from 'highlight.js/lib/languages/bash'
import shell from 'highlight.js/lib/languages/shell'
import dockerfile from 'highlight.js/lib/languages/dockerfile'
// 其他
import diff from 'highlight.js/lib/languages/diff'

// 注册语言
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('js', javascript)
hljs.registerLanguage('typescript', typescript)
hljs.registerLanguage('ts', typescript)
hljs.registerLanguage('python', python)
hljs.registerLanguage('py', python)
hljs.registerLanguage('go', go)
hljs.registerLanguage('java', java)
hljs.registerLanguage('sql', sql)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('html', xml)
hljs.registerLanguage('css', css)
hljs.registerLanguage('json', json)
hljs.registerLanguage('yaml', yaml)
hljs.registerLanguage('yml', yaml)
hljs.registerLanguage('markdown', markdown)
hljs.registerLanguage('md', markdown)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('sh', bash)
hljs.registerLanguage('shell', shell)
hljs.registerLanguage('dockerfile', dockerfile)
hljs.registerLanguage('docker', dockerfile)
hljs.registerLanguage('diff', diff)

// ========== 属性解析函数 ==========

/**
 * 提取标签名和参数
 * @param line - 完整的标签行，格式：`:::标签名 参数1 参数2 ...`
 * @returns 标签名和参数数组
 */
function extractTagAndParams(line: string): { tag: string; params: string[] } {
  const match = line.match(/^:::(\w+)(.*)$/)
  if (!match) return { tag: '', params: [] }
  const tag = match[1] || ''
  const paramsString = match[2]?.trim() || ''

  // 简单按空格分割参数
  const params = paramsString ? paramsString.split(/\s+/).filter(p => p && p !== ':::') : []

  return { tag, params }
}

/**
 * 检查是否为自闭合标签
 * @param line - 标签行
 * @returns 是否为自闭合标签
 */
function isSelfClosing(line: string): boolean {
  return /:::$/.test(line.trim())
}

// 简单哈希函数（确定性）
function simpleHash(str: string): string {
  let hash = 0
  for (let i = 0; i < str.length; i++) {
    hash = ((hash << 5) - hash) + str.charCodeAt(i)
    hash |= 0
  }
  return Math.abs(hash).toString(36)
}

// 生成标题 ID（支持中文）
function generateHeadingId(text: string): string {
  const id = text.toLowerCase()
    .replace(/[^\u4e00-\u9fa5a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
  return id || `heading-${simpleHash(text)}`
}

// ========== 自定义块渲染函数 ==========

/**
 * 渲染提示框
 * @param content - 内容
 * @param params - [类型, 标题(可选)]
 */
function renderNote(content: string, params: string[]): string {
  const type = params[0] || 'info'
  const title = params[1] || ''

  const titleHtml = title ? `<div class="custom-note-title">${title}</div>` : ''

  return `<div class="custom-note custom-note-${type}">${titleHtml}<div class="custom-note-content">${content}</div></div>`
}

/**
 * 渲染标签页
 * @param tabsData - 标签数据
 * @param params - [默认标签名(可选)]
 */
function renderTabs(tabsData: Array<{ name: string; content: string }>, params: string[]): string {
  if (tabsData.length === 0) return ''

  const tabsId = `tabs-${simpleHash(tabsData.map(t => t.name).join('-'))}`
  const activeTab = params[0] || tabsData[0]?.name || ''

  // 生成标签头
  const tabHeaders = tabsData.map(tab => {
    const isActive = tab.name === activeTab ? 'active' : ''
    return `<button class="custom-tab-btn ${isActive}" onclick="switchTab('${tabsId}', '${tab.name}')">${tab.name}</button>`
  }).join('')

  // 生成标签内容
  const tabContents = tabsData.map(tab => {
    const isActive = tab.name === activeTab ? 'active' : ''
    return `<div class="custom-tab-panel ${isActive}" data-tab="${tab.name}">${tab.content}</div>`
  }).join('')

  return `<div class="custom-tabs" id="${tabsId}"><div class="custom-tabs-header">${tabHeaders}</div><div class="custom-tabs-content">${tabContents}</div></div>`
}

/**
 * 渲染折叠面板
 * @param content - 内容
 * @param params - [标题, open(可选)]
 */
function renderFold(content: string, params: string[]): string {
  const title = params[0] || '点击展开'
  const open = params[1] === 'true' || params[1] === 'open'
  const foldId = `fold-${simpleHash(title + content.slice(0, 50))}`
  const openClass = open ? 'open' : ''

  return `<div class="custom-fold ${openClass}" id="${foldId}"><div class="custom-fold-header" onclick="toggleFold('${foldId}')"><i class="ri-arrow-right-s-line"></i><span>${title}</span></div><div class="custom-fold-content"><div>${content}</div></div></div>`
}

/**
 * 渲染链接卡片
 * @param params - [标题, 链接, 描述(可包含空格)]
 */
function renderLinkCard(params: string[]): string {
  const title = params[0] || ''
  const link = params[1] || ''
  const description = params.slice(2).join(' ')

  if (!link) return ''

  // 判断是否为外部链接
  const isExternal = link.startsWith('http://') || link.startsWith('https://')
  const linkType = isExternal ? '引用站外链接' : '站内链接'
  const linkTypeClass = isExternal ? 'external' : 'internal'

  return `<div class="custom-link-card ${linkTypeClass}">
    <div class="custom-link-type">${linkType}</div>
    <a href="${link}" class="custom-link-main" target="${isExternal ? '_blank' : '_self'}" rel="${isExternal ? 'noopener noreferrer' : ''}">
      <div class="custom-link-icon">
        <i class="ri-global-line"></i>
      </div>
      <div class="custom-link-info">
        <div class="custom-link-title">${title}</div>
        <div class="custom-link-desc">${description || link}</div>
      </div>
      <div class="custom-link-arrow">
        <i class="ri-arrow-right-up-line"></i>
      </div>
    </a>
  </div>`
}

// 创建 markdown-it 实例
const md = new MarkdownIt({
  html: false,
  breaks: true,
  linkify: true
})

// 自定义代码块渲染规则
md.renderer.rules.fence = (tokens, idx) => {
  const token = tokens[idx]
  if (!token) return ''

  const code = token.content
  const lang = token.info.trim()

  // 特殊处理 Mermaid 代码块（不进行代码高亮）
  if (lang === 'mermaid') {
    return `<pre class="mermaid"><code>${md.utils.escapeHtml(code)}</code></pre>`
  }

  // 高亮代码
  let highlightedCode = ''
  const displayLang = (lang || 'text').toUpperCase()

  if (lang && hljs.getLanguage(lang)) {
    try {
      highlightedCode = hljs.highlight(code, { language: lang, ignoreIllegals: true }).value
    } catch {
      highlightedCode = md.utils.escapeHtml(code)
    }
  } else {
    highlightedCode = md.utils.escapeHtml(code)
  }

  // 添加行号（移除末尾换行符避免空行）
  const numberedLines = highlightedCode
    .replace(/\n$/, '')
    .split('\n')
    .map((line, index) => `<span class="line-number" data-line="${index + 1}"></span><span class="line-content">${line}</span>`)
    .join('\n')

  // 返回完整结构
  return `<div class="code-block-container"><div class="code-toolbar"><button class="code-fold-btn" onclick="this.closest('.code-block-container').classList.toggle('collapsed')" title="折叠/展开"><i class="ri-arrow-down-s-line"></i></button><span class="code-lang">${displayLang}</span><button class="code-copy-btn" onclick="copyCodeBlock(this)" title="复制代码"><i class="ri-file-copy-fill"></i></button></div><pre><code>${numberedLines}</code></pre></div>`
}

// 使用 anchor 插件生成标题 ID
md.use(anchor, {
  slugify: generateHeadingId,
  permalink: false,
  level: [1, 2, 3, 4, 5, 6]
})

// 使用任务列表插件
md.use(taskLists, {
  enabled: true,
  label: true,
  labelAfter: false
})

// 使用高亮文本插件
md.use(mark)

// 使用链接属性插件（外部链接在新窗口打开）
md.use(linkAttributes, {
  matcher(href: string) {
    return href.startsWith('http://') || href.startsWith('https://')
  },
  attrs: {
    target: '_blank',
    rel: 'noopener noreferrer'
  }
})

// 使用键盘按键插件（支持 [[Ctrl]] 语法）
md.use(kbd)

// 使用上标插件（支持 ^上标^ 语法）
md.use(sup)

// 使用下标插件（支持 ~下标~ 语法）
md.use(sub)

// 使用下划线插件（支持 ++下划线++ 语法）
md.use(underline)

// ========== 自定义块插件 ==========

/**
 * 自定义块插件
 */
function customBlocksPlugin(md: MarkdownIt) {
  // 块级规则
  md.block.ruler.before('fence', 'custom_blocks', (state, startLine, endLine, silent) => {
    const pos = (state.bMarks[startLine] ?? 0) + (state.tShift[startLine] ?? 0)
    const max = state.eMarks[startLine] ?? 0
    const lineText = state.src.slice(pos, max).trim()

    // 检查是否为自定义块起始标签
    if (!lineText.startsWith(':::')) {
      return false
    }

    // 检查是否为自闭合标签
    if (isSelfClosing(lineText)) {
      if (silent) return true

      const { tag, params } = extractTagAndParams(lineText)

      // 处理自闭合标签
      let html = ''
      if (tag === 'link') {
        html = renderLinkCard(params)
      }

      if (html) {
        const token = state.push('html_block', '', 0)
        token.content = html
        token.map = [startLine, startLine + 1]
        state.line = startLine + 1
        return true
      }

      return false
    }

    // 处理块级标签
    const { tag, params } = extractTagAndParams(lineText)
    if (!tag) return false

    // 查找结束标签
    const endTagFull = `end${tag}`
    let nextLine = startLine + 1
    let foundEnd = false
    let contentLines: string[] = []

    // 特殊处理 tabs
    if (tag === 'tabs') {
      const tabsData: Array<{ name: string; content: string }> = []
      let currentTab: { name: string; content: string } | null = null

      while (nextLine < endLine) {
        const linePos = state.bMarks[nextLine] ?? 0
        const lineMax = state.eMarks[nextLine] ?? 0
        const line = state.src.slice(linePos, lineMax).trim()

        if (line.startsWith(':::endtabs')) {
          foundEnd = true
          break
        }

        if (line.startsWith(':::tab')) {
          // 保存上一个 tab
          if (currentTab) {
            tabsData.push(currentTab)
          }
          // 开始新 tab
          const tabParams = extractTagAndParams(line).params
          currentTab = { name: tabParams[0] || `Tab ${tabsData.length + 1}`, content: '' }
        } else if (line.startsWith(':::endtab')) {
          // tab 结束，不做处理
        } else {
          // tab 内容
          if (currentTab) {
            currentTab.content += state.src.slice(linePos, lineMax) + '\n'
          }
        }

        nextLine++
      }

      // 保存最后一个 tab
      if (currentTab) {
        tabsData.push(currentTab)
      }

      if (foundEnd && tabsData.length > 0) {
        if (silent) return true

        // 渲染每个 tab 的内容
        const renderedTabs = tabsData.map(tab => ({
          name: tab.name,
          content: md.render(tab.content)
        }))

        const html = renderTabs(renderedTabs, params)

        const token = state.push('html_block', '', 0)
        token.content = html
        token.map = [startLine, nextLine + 1]
        state.line = nextLine + 1
        return true
      }

      return false
    }

    // 特殊处理 photo
    if (tag === 'photo') {
      const rows: string[][] = []
      let currentRow: string[] = []

      while (nextLine < endLine) {
        const linePos = (state.bMarks[nextLine] ?? 0) + (state.tShift[nextLine] ?? 0)
        const lineMax = state.eMarks[nextLine] ?? 0
        const line = state.src.slice(linePos, lineMax).trim()

        if (line === ':::endphoto') {
          foundEnd = true
          break
        }

        // 检查是否为换行标记 :::n
        if (line === ':::n') {
          // 保存当前行并开始新行
          if (currentRow.length > 0) {
            rows.push(currentRow)
            currentRow = []
          }
        } else {
          // 解析图片（支持多个图片用空格分隔）
          const images = line.split(/\s+/).filter(img => img.trim())
          currentRow.push(...images)
        }

        nextLine++
      }

      // 保存最后一行
      if (currentRow.length > 0) {
        rows.push(currentRow)
      }

      if (foundEnd && rows.length > 0) {
        if (silent) return true

        const html = renderPhotoWall(rows, startLine)

        const token = state.push('html_block', '', 0)
        token.content = html
        token.map = [startLine, nextLine + 1]
        state.line = nextLine + 1
        return true
      }

      return false
    }

    // 处理其他块级标签（note, fold）
    while (nextLine < endLine) {
      const linePos = state.bMarks[nextLine] ?? 0
      const lineMax = state.eMarks[nextLine] ?? 0
      const line = state.src.slice(linePos, lineMax).trim()

      if (line === `:::${endTagFull}`) {
        foundEnd = true
        break
      }

      contentLines.push(state.src.slice(linePos, lineMax))
      nextLine++
    }

    if (!foundEnd) return false
    if (silent) return true

    // 渲染内容
    const content = md.render(contentLines.join('\n'))

    let html = ''
    if (tag === 'note') {
      html = renderNote(content, params)
    } else if (tag === 'fold') {
      html = renderFold(content, params)
    }

    if (html) {
      const token = state.push('html_block', '', 0)
      token.content = html
      token.map = [startLine, nextLine + 1]
      state.line = nextLine + 1
      return true
    }

    return false
  })
}

// 使用自定义块插件
md.use(customBlocksPlugin)

// 渲染 Markdown 为 HTML
export function renderMarkdown(markdown: string): string {
  if (!markdown) return ''

  const rawHtml = md.render(markdown)

  // 替换表情占位符为 img 标签
  const emojiMap = getEmojiMapSync()
  let processedHtml = rawHtml
  if (emojiMap && emojiMap.size > 0) {
    processedHtml = replaceEmojisInText(rawHtml, emojiMap)
  }

  return DOMPurify.sanitize(processedHtml, {
    ALLOWED_TAGS: [
      'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'p', 'br', 'hr',
      'strong', 'em', 'u', 's', 'del', 'ins', 'mark', 'code', 'pre',
      'ul', 'ol', 'li', 'blockquote', 'cite', 'footer',
      'a', 'img', 'table', 'thead', 'tbody', 'tr', 'th', 'td',
      'div', 'span', 'sup', 'sub', 'kbd', 'abbr',
      'input', 'label', 'button', 'i', 'section',
      'svg', 'path', 'g', 'rect', 'circle', 'ellipse', 'line', 'polygon', 'polyline', 'text', 'foreignObject'
    ],
    ALLOWED_ATTR: [
      'href', 'title', 'target', 'rel', 'src', 'alt', 'width', 'height',
      'class', 'id', 'colspan', 'rowspan', 'align',
      'type', 'checked', 'disabled', 'for', 'onclick', 'start',
      'd', 'fill', 'stroke', 'stroke-width', 'x', 'y', 'cx', 'cy', 'r', 'rx', 'ry',
      'x1', 'y1', 'x2', 'y2', 'points', 'transform', 'viewBox', 'xmlns',
      'text-anchor', 'font-size', 'font-family', 'dominant-baseline', 'data-processed'
    ],
    ALLOW_DATA_ATTR: true,
    ADD_ATTR: ['target', 'onclick']
  })
}

// 计算字数
export function countWords(markdown: string): number {
  if (!markdown) return 0

  // 服务端渲染时，直接从 markdown 统计
  if (!process.client) {
    const text = markdown
      .replace(/```[\s\S]*?```/g, '') // 移除代码块
      .replace(/`[^`]+`/g, '') // 移除行内代码
      .replace(/!\[.*?\]\(.*?\)/g, '') // 移除图片
      .replace(/\[.*?\]\(.*?\)/g, '') // 移除链接
      .replace(/[#*_~`|]/g, '') // 移除 markdown 符号
      .trim()

    const chineseChars = text.match(/[\u4e00-\u9fa5]/g) || []
    const englishWords = text.match(/[a-zA-Z]+/g) || []
    return chineseChars.length + englishWords.length
  }

  // 客户端：使用 DOM 提取文本
  const html = md.render(markdown)
  const temp = document.createElement('div')
  temp.innerHTML = html

  // 移除代码块（不统计代码）
  temp.querySelectorAll('pre, code').forEach(el => el.remove())

  // 提取纯文本
  const text = temp.textContent?.trim() || ''

  // 统计中英文字数
  const chineseChars = text.match(/[\u4e00-\u9fa5]/g) || []
  const englishWords = text.match(/[a-zA-Z]+/g) || []
  return chineseChars.length + englishWords.length
}

// 计算阅读时长（分钟）
export function estimateReadingTime(markdown: string, wordsPerMinute = 300): number {
  return Math.ceil(countWords(markdown) / wordsPerMinute)
}

// 目录项接口
export interface TocItem {
  id: string
  level: number
  text: string
  children?: TocItem[]
}

// 提取目录
export function extractToc(markdown: string): TocItem[] {
  if (!markdown) return []

  // 移除代码块
  let cleanedMarkdown = markdown
    .replace(/```[\s\S]*?```/g, '')
    .replace(/~~~[\s\S]*?~~~\s*/g, '')
    .replace(/^(    |\t).+$/gm, '')

  // 处理单行自定义块
  cleanedMarkdown = cleanedMarkdown.replace(/^:::link\s+.*?:::$/gm, '')
  // 处理多行自定义块
  cleanedMarkdown = cleanedMarkdown.replace(/^:::note[\s\S]*?^:::endnote$/gm, '')
  cleanedMarkdown = cleanedMarkdown.replace(/^:::tabs[\s\S]*?^:::endtabs$/gm, '')
  cleanedMarkdown = cleanedMarkdown.replace(/^:::fold[\s\S]*?^:::endfold$/gm, '')
  cleanedMarkdown = cleanedMarkdown.replace(/^:::photo[\s\S]*?^:::endphoto$/gm, '')

  const headings: TocItem[] = []

  for (const line of cleanedMarkdown.split('\n')) {
    const match = line.match(/^(#{1,6})\s+(.+)$/)
    if (match?.[1] && match[2]) {
      headings.push({
        id: generateHeadingId(match[2].trim()),
        level: match[1].length,
        text: match[2].trim()
      })
    }
  }

  return headings
}

/**
 * 渲染照片展示墙
 * @param rows - 每行的图片数组
 * @param lineNum - 源码行号（可选，用于滚动同步）
 */
function renderPhotoWall(rows: string[][], lineNum?: number): string {
  if (rows.length === 0) return ''

  const lineAttr = lineNum !== undefined ? ` data-source-line="${lineNum}"` : ''

  // 生成每一行的图片
  const rowsHtml = rows.map(row => {
    const imagesHtml = row.map(img => {
      // 处理图片语法：支持 markdown 图片语法和直接 URL
      let imgSrc = img
      let imgAlt = ''

      // 检查是否为 markdown 图片语法 ![alt](url)
      const imgMatch = img.match(/^!\[(.*?)\]\((.*?)\)$/)
      if (imgMatch) {
        imgAlt = imgMatch[1] || ''
        imgSrc = imgMatch[2] || img
      }

      return `<div class="custom-photo-wall-item"><img src="${imgSrc}" alt="${imgAlt || '图片'}" loading="lazy" /></div>`
    }).join('')

    return `<div class="custom-photo-wall-row">${imagesHtml}</div>`
  }).join('')

  return `<div class="custom-photo-wall"${lineAttr}><div class="custom-photo-wall-container">${rowsHtml}</div></div>`
}

// 简单 Markdown 渲染（用于评论）
export function renderSimpleMarkdown(markdown: string): string {
  if (!markdown) return ''

  const simpleMd = new MarkdownIt({
    html: false,
    breaks: true,
    linkify: true
  })

  // 先渲染 Markdown
  let simpleHtml = simpleMd.render(markdown)

  // 然后替换表情占位符（在 HTML 中替换）
  const emojiMap = getEmojiMapSync()
  if (emojiMap && emojiMap.size > 0) {
    simpleHtml = replaceEmojisInText(simpleHtml, emojiMap)
  }

  return DOMPurify.sanitize(simpleHtml, {
    ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'code', 'pre', 'ul', 'ol', 'li', 'blockquote', 'a', 'img'],
    ALLOWED_ATTR: ['href', 'title', 'src', 'alt', 'width', 'height', 'class'],
    ALLOWED_URI_REGEXP: /^(?:(?:(?:f|ht)tps?|mailto|tel|callto|sms|cid|xmpp|blob):|[^a-z]|[a-z+.\-]+(?:[^a-z+.\-:]|$))/i,
    ALLOW_DATA_ATTR: false
  })
}

// 复制代码块功能
export function copyCodeBlock(button: HTMLElement): void {
  const container = button.closest('.code-block-container')
  if (!container) return

  const code = container.querySelector('code')
  if (!code) return

  // 只提取代码内容，不包含行号
  const codeLines = Array.from(code.querySelectorAll('.line-content'))
  const codeText = codeLines.map(line => line.textContent || '').join('\n')

  // 复制到剪贴板
  navigator.clipboard.writeText(codeText).then(() => {
    // 更新按钮状态
    const icon = button.querySelector('i')
    if (icon) {
      icon.className = 'ri-check-line'
      button.classList.add('copied')
    }

    // 2秒后恢复
    setTimeout(() => {
      if (icon) {
        icon.className = 'ri-file-copy-fill'
        button.classList.remove('copied')
      }
    }, 2000)
  }).catch(err => {
    console.error('复制失败:', err)
  })
}

// 标签页切换功能
export function switchTab(tabsId: string, tabName: string): void {
  const tabsContainer = document.getElementById(tabsId)
  if (!tabsContainer) return

  // 更新标签按钮状态
  const buttons = tabsContainer.querySelectorAll('.custom-tab-btn')
  buttons.forEach(btn => {
    if (btn.textContent === tabName) {
      btn.classList.add('active')
    } else {
      btn.classList.remove('active')
    }
  })

  // 更新内容面板状态
  const panels = tabsContainer.querySelectorAll('.custom-tab-panel')
  panels.forEach(panel => {
    const panelElement = panel as HTMLElement
    if (panelElement.dataset.tab === tabName) {
      panel.classList.add('active')
    } else {
      panel.classList.remove('active')
    }
  })
}

// 折叠面板切换功能
export function toggleFold(foldId: string): void {
  const foldContainer = document.getElementById(foldId)
  if (!foldContainer) return

  foldContainer.classList.toggle('open')
}

// 挂载全局函数供内联 onclick 使用
if (typeof window !== 'undefined') {
  (window as any).copyCodeBlock = copyCodeBlock;
  (window as any).switchTab = switchTab;
  (window as any).toggleFold = toggleFold
}

export default {
  render: renderMarkdown,
  renderSimple: renderSimpleMarkdown,
  countWords,
  estimateReadingTime,
  extractToc,
  copyCodeBlock
}
