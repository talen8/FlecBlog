import MarkdownIt from 'markdown-it';
import anchor from 'markdown-it-anchor';
// @ts-expect-error - 没有类型定义
import taskLists from 'markdown-it-task-lists';
// @ts-expect-error - 没有类型定义
import mark from 'markdown-it-mark';
// @ts-expect-error - 没有类型定义
import linkAttributes from 'markdown-it-link-attributes';
import kbd from 'markdown-it-kbd';
// @ts-expect-error - 没有类型定义
import sub from 'markdown-it-sub';
// @ts-expect-error - 没有类型定义
import sup from 'markdown-it-sup';
// @ts-expect-error - 没有类型定义
import underline from 'markdown-it-plugin-underline';
import katex from '@traptitech/markdown-it-katex';
import DOMPurify from 'dompurify';
import hljs from 'highlight.js';

// ========== 属性解析函数 ==========

/**
 * 提取标签名和参数
 * @param line - 完整的标签行，格式：`:::标签名 参数1 参数2 ...`
 * @returns 标签名和参数数组
 */
function extractTagAndParams(line: string): { tag: string; params: string[] } {
  const match = line.match(/^:::(\w+)(.*)$/);

  if (!match) {
    return { tag: '', params: [] };
  }

  const tag = match[1] || '';
  const paramsString = match[2]?.trim() || '';

  // 简单按空格分割参数
  const params = paramsString ? paramsString.split(/\s+/).filter(p => p && p !== ':::') : [];

  return { tag, params };
}

/**
 * 检查是否为自闭合标签
 * @param line - 标签行
 * @returns 是否为自闭合标签
 */
function isSelfClosing(line: string): boolean {
  return /:::$/.test(line.trim());
}

// 生成标题 ID（支持中文）
function generateHeadingId(text: string): string {
  const id = text
    .toLowerCase()
    .replace(/[^\u4e00-\u9fa5a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
  return id || `heading-${Math.random().toString(36).slice(2, 9)}`;
}

const INLINE_ANCHOR_CHUNK_SIZE = 24;

function escapeHtmlContent(value: string): string {
  return value
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

function getLineStartOffsets(source: string): number[] {
  const offsets = [0];
  for (let index = 0; index < source.length; index++) {
    if (source[index] === '\n') offsets.push(index + 1);
  }
  return offsets;
}

function getOffsetForLine(lineStarts: number[], line: number, sourceLength: number): number {
  if (line <= 0) return 0;
  if (line >= lineStarts.length) return sourceLength;
  return lineStarts[line] ?? sourceLength;
}

function setTokenSourceMeta(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  token: any,
  sourceStartLine: number,
  sourceEndLine: number,
  sourceStartOffset: number,
  sourceEndOffset: number
) {
  token.meta = {
    ...(token.meta || {}),
    sourceStartLine,
    sourceEndLine,
    sourceStartOffset,
    sourceEndOffset,
  };
}

function buildSourceAttrs(
  meta?: {
    sourceStartLine?: number;
    sourceEndLine?: number;
    sourceStartOffset?: number;
    sourceEndOffset?: number;
  },
  kind: 'block' | 'text' = 'block'
): string {
  if (!meta) return '';

  const attrs: string[] = [`data-sync-kind="${kind}"`];
  if (meta.sourceStartLine !== undefined) {
    attrs.push(`data-source-line="${meta.sourceStartLine}"`);
    attrs.push(`data-source-start-line="${meta.sourceStartLine}"`);
  }
  if (meta.sourceEndLine !== undefined) attrs.push(`data-source-end-line="${meta.sourceEndLine}"`);
  if (meta.sourceStartOffset !== undefined)
    attrs.push(`data-source-start-offset="${meta.sourceStartOffset}"`);
  if (meta.sourceEndOffset !== undefined)
    attrs.push(`data-source-end-offset="${meta.sourceEndOffset}"`);
  return attrs.length ? ` ${attrs.join(' ')}` : '';
}

function splitTextIntoChunks(text: string): Array<{ text: string; start: number; end: number }> {
  const chunks: Array<{ text: string; start: number; end: number }> = [];
  let start = 0;

  while (start < text.length) {
    let end = Math.min(start + INLINE_ANCHOR_CHUNK_SIZE, text.length);
    if (end < text.length) {
      const boundary = text.lastIndexOf(' ', end);
      if (boundary > start + 8) end = boundary + 1;
    }
    if (end <= start) end = Math.min(start + INLINE_ANCHOR_CHUNK_SIZE, text.length);
    chunks.push({ text: text.slice(start, end), start, end });
    start = end;
  }

  return chunks;
}

function renderAnchoredText(
  text: string,
  sourceStartOffset?: number,
  sourceEndOffset?: number,
  escapeHtml: (value: string) => string = escapeHtmlContent
): string {
  const escapedText = escapeHtml(text);
  if (
    sourceStartOffset === undefined ||
    sourceEndOffset === undefined ||
    sourceEndOffset <= sourceStartOffset
  ) {
    return escapedText;
  }
  if (!text.trim()) return escapedText;

  return splitTextIntoChunks(text)
    .map(chunk => {
      const attrs = buildSourceAttrs(
        {
          sourceStartOffset: sourceStartOffset + chunk.start,
          sourceEndOffset: sourceStartOffset + chunk.end,
        },
        'text'
      );
      return `<span class="sync-text-anchor"${attrs}>${escapeHtml(chunk.text)}</span>`;
    })
    .join('');
}

// ========== 自定义块渲染函数 ==========

/**
 * 渲染提示框
 * @param content - 内容
 * @param params - [类型, 标题(可选)]
 * @param lineNum - 源码行号（可选，用于滚动同步）
 */
function renderNote(content: string, params: string[], sourceAttrs = ''): string {
  const type = params[0] || 'info';
  const title = params[1] || '';

  const titleHtml = title ? `<div class="custom-note-title">${title}</div>` : '';

  return `<div class="custom-note custom-note-${type}"${sourceAttrs}>${titleHtml}<div class="custom-note-content">${content}</div></div>`;
}

/**
 * 渲染标签页
 * @param tabsData - 标签数据
 * @param params - [默认标签名(可选)]
 * @param lineNum - 源码行号（可选，用于滚动同步）
 */
function renderTabs(
  tabsData: Array<{ name: string; content: string }>,
  params: string[],
  sourceAttrs = ''
): string {
  if (tabsData.length === 0) return '';

  const tabsId = `tabs-${Math.random().toString(36).slice(2, 9)}`;
  const activeTab = params[0] || tabsData[0]?.name || '';

  // 生成标签头
  const tabHeaders = tabsData
    .map(tab => {
      const isActive = tab.name === activeTab ? 'active' : '';
      return `<button class="custom-tab-btn ${isActive}" onclick="switchTab('${tabsId}', '${tab.name}')">${tab.name}</button>`;
    })
    .join('');

  // 生成标签内容
  const tabContents = tabsData
    .map(tab => {
      const isActive = tab.name === activeTab ? 'active' : '';
      return `<div class="custom-tab-panel ${isActive}" data-tab="${tab.name}">${tab.content}</div>`;
    })
    .join('');

  return `<div class="custom-tabs" id="${tabsId}"${sourceAttrs}><div class="custom-tabs-header">${tabHeaders}</div><div class="custom-tabs-content">${tabContents}</div></div>`;
}

/**
 * 渲染折叠面板
 * @param content - 内容
 * @param params - [标题, open(可选)]
 * @param lineNum - 源码行号（可选，用于滚动同步）
 */
function renderFold(content: string, params: string[], sourceAttrs = ''): string {
  const title = params[0] || '点击展开';
  const open = params[1] === 'true' || params[1] === 'open';
  const foldId = `fold-${Math.random().toString(36).slice(2, 9)}`;
  const openClass = open ? 'open' : '';

  return `<div class="custom-fold ${openClass}" id="${foldId}"${sourceAttrs}><div class="custom-fold-header" onclick="toggleFold('${foldId}')"><i class="ri-arrow-right-s-line"></i><span>${title}</span></div><div class="custom-fold-content"><div>${content}</div></div></div>`;
}

/**
 * 渲染链接卡片
 * @param params - [标题, 链接, 描述(可包含空格)]
 * @param lineNum - 源码行号（可选，用于滚动同步）
 */
function renderLinkCard(params: string[], sourceAttrs = ''): string {
  const title = params[0] || '';
  const link = params[1] || '';
  const description = params.slice(2).join(' ');

  if (!link) return '';

  // 判断是否为外部链接
  const isExternal = link.startsWith('http://') || link.startsWith('https://');
  const linkType = isExternal ? '引用站外链接' : '站内链接';
  const linkTypeClass = isExternal ? 'external' : 'internal';

  return `<div class="custom-link-card ${linkTypeClass}"${sourceAttrs}>
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
  </div>`;
}

/**
 * 渲染在线视频
 * @param params - [平台或URL, 视频ID(可选)]
 * @param lineNum - 源码行号（可选，用于滚动同步）
 */
function renderVideo(params: string[], sourceAttrs = ''): string {
  if (params.length === 0) return '';
  const platformOrUrl = params[0] || '';
  const videoId = params[1] || '';

  // B站视频
  if (platformOrUrl === 'bilibili' && videoId) {
    return `<div class="custom-video"${sourceAttrs}><iframe src="//player.bilibili.com/player.html?bvid=${videoId}&autoplay=0" scrolling="no" border="0" frameborder="no" framespacing="0" allowfullscreen="true" sandbox="allow-scripts allow-same-origin allow-popups" referrerpolicy="strict-origin-when-cross-origin"></iframe></div>`;
  }

  // YouTube视频
  if (platformOrUrl === 'youtube' && videoId) {
    return `<div class="custom-video"${sourceAttrs}><iframe src="https://www.youtube.com/embed/${videoId}" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" allowfullscreen sandbox="allow-scripts allow-same-origin allow-popups" referrerpolicy="strict-origin-when-cross-origin"></iframe></div>`;
  }

  // 本地/在线视频URL
  if (
    platformOrUrl.startsWith('http://') ||
    platformOrUrl.startsWith('https://') ||
    platformOrUrl.startsWith('/')
  ) {
    return `<div class="custom-video"${sourceAttrs}><video src="${platformOrUrl}" controls preload="metadata"></video></div>`;
  }

  return '';
}

/**
 * 渲染在线音乐/音频
 * @param params - [标题, 音频URL]
 * @param sourceAttrs - 源码属性（可选，用于滚动同步）
 */
function renderAudio(params: string[], sourceAttrs = ''): string {
  const title = params[0] || '';
  const audioUrl = params[1] || '';

  if (!title || (!audioUrl.startsWith('http://') && !audioUrl.startsWith('https://'))) {
    return '';
  }

  const audioId = `audio-${Math.random().toString(36).slice(2, 9)}`;

  return `<div class="custom-audio" data-audio-id="${audioId}"${sourceAttrs}>
  <div class="custom-audio-type">播放音频</div>
  <div class="custom-audio-main">
    <div class="custom-audio-icon">
      <i class="ri-music-line"></i>
      <button class="custom-audio-btn" onclick="toggleAudioPlay('${audioId}')">
        <i class="ri-play-fill"></i>
      </button>
    </div>
    <div class="custom-audio-content">
      <div class="custom-audio-info">${title}</div>
      <div class="custom-audio-controls">
        <div class="custom-audio-progress" onclick="seekAudio('${audioId}', event)">
          <div class="custom-audio-progress-bar"></div>
        </div>
        <div class="custom-audio-time"><span class="custom-audio-current">0:00</span> / <span class="custom-audio-duration">0:00</span></div>
      </div>
    </div>
  </div>
  <audio src="${audioUrl}" preload="auto" style="display:none;"></audio>
</div>`;
}

/**
 * 渲染在线音乐
 * @param params - [平台, 音乐ID]
 * @param sourceAttrs - 源码属性（可选）
 */
function renderMusic(params: string[], sourceAttrs = ''): string {
  if (params.length < 2) return '';

  const server = params[0] || '';
  const musicId = params[1] || '';

  if (!server || !musicId) return '';

  const audioId = `audio-${Math.random().toString(36).slice(2, 9)}`;

  const embedUrl = `https://meting.flec.top/api?server=${server}&type=song&id=${musicId}`;

  return `<div class="custom-audio" data-audio-id="${audioId}" data-music-id="${musicId}"${sourceAttrs}>
  <div class="custom-audio-type">播放在线音乐</div>
  <div class="custom-audio-main">
    <div class="custom-audio-icon">
      <i class="ri-music-line"></i>
      <button class="custom-audio-btn" onclick="toggleMusicPlay('${audioId}', '${server}', '${musicId}')">
        <i class="ri-play-fill"></i>
      </button>
    </div>
    <div class="custom-audio-content">
      <div class="custom-audio-info" data-music-info="${audioId}">加载中...</div>
      <div class="custom-audio-controls">
        <div class="custom-audio-progress" onclick="seekMusic('${audioId}', event)">
          <div class="custom-audio-progress-bar"></div>
        </div>
        <div class="custom-audio-time"><span class="custom-audio-current">0:00</span> / <span class="custom-audio-duration">0:00</span></div>
      </div>
    </div>
  </div>
  <div class="custom-audio-source" data-embed-url="${embedUrl}" style="display:none;"></div>
</div>`;
}

/**
 * 渲染照片展示墙
 * @param rows - 每行的图片数组
 * @param lineNum - 源码行号（可选，用于滚动同步）
 */
function renderPhotoWall(rows: string[][], sourceAttrs = ''): string {
  if (rows.length === 0) return '';

  // 生成每一行的图片
  const rowsHtml = rows
    .map(row => {
      const imagesHtml = row
        .map(img => {
          // 处理图片语法：支持 markdown 图片语法和直接 URL
          let imgSrc = img;
          let imgAlt = '';

          // 检查是否为 markdown 图片语法 ![alt](url)
          const imgMatch = img.match(/^!\[(.*?)\]\((.*?)\)$/);
          if (imgMatch) {
            imgAlt = imgMatch[1] || '';
            imgSrc = imgMatch[2] || img;
          }

          return `<div class="custom-photo-wall-item"><img src="${imgSrc}" alt="${imgAlt || '图片'}" loading="lazy" /></div>`;
        })
        .join('');

      return `<div class="custom-photo-wall-row">${imagesHtml}</div>`;
    })
    .join('');

  return `<div class="custom-photo-wall"${sourceAttrs}><div class="custom-photo-wall-container">${rowsHtml}</div></div>`;
}

function stripNestedSyncAttrs(content: string): string {
  return content.replace(/\s*data-(?:source|sync)-[\w-]+="[^"]*"/g, '');
}

function createMarkdownRenderer(): MarkdownIt {
  const instance = new MarkdownIt({
    html: false,
    breaks: true,
    linkify: true,
  });

  instance.use(anchor, {
    slugify: generateHeadingId,
    permalink: false,
    level: [1, 2, 3, 4, 5, 6],
  });

  instance.use(taskLists, {
    enabled: true,
    label: true,
    labelAfter: false,
  });

  instance.use(mark);
  instance.use(linkAttributes, {
    matcher(href: string) {
      return href.startsWith('http://') || href.startsWith('https://');
    },
    attrs: {
      target: '_blank',
      rel: 'noopener noreferrer',
    },
  });
  instance.use(kbd);
  instance.use(sup);
  instance.use(sub);
  instance.use(underline);
  instance.use(katex, { throwOnError: false, errorColor: '#cc0000' });
  instance.use(customBlocksPlugin);

  return instance;
}

function renderFence(
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  token: any,
  escapeHtml: (value: string) => string,
  sourceAttrs = ''
): string {
  const code = token.content;
  const lang = token.info.trim();

  if (lang === 'mermaid') {
    return `<pre class="mermaid"${sourceAttrs}><code>${escapeHtml(code)}</code></pre>`;
  }

  let highlightedCode = '';
  const displayLang = (lang || 'text').toUpperCase();

  if (lang && hljs.getLanguage(lang)) {
    try {
      highlightedCode = hljs.highlight(code, {
        language: lang,
        ignoreIllegals: true,
      }).value;
    } catch {
      highlightedCode = escapeHtml(code);
    }
  } else {
    highlightedCode = escapeHtml(code);
  }

  const numberedLines = highlightedCode
    .replace(/\n$/, '')
    .split('\n')
    .map(
      (line, index) =>
        `<span class="line-number" data-line="${index + 1}"></span><span class="line-content">${line}</span>`
    )
    .join('\n');

  return `<div class="code-block-container"${sourceAttrs}><div class="code-toolbar"><button class="code-fold-btn" onclick="this.closest('.code-block-container').classList.toggle('collapsed')" title="折叠/展开"><i class="ri-arrow-down-s-line"></i></button><span class="code-lang">${displayLang}</span><button class="code-copy-btn" onclick="copyCodeBlock(this)" title="复制代码"><i class="ri-file-copy-fill"></i></button></div><pre><code>${numberedLines}</code></pre></div>`;
}

// 创建 markdown-it 实例
const md = createMarkdownRenderer();
md.renderer.rules.fence = (tokens, idx) => {
  const token = tokens[idx];
  if (!token) return '';
  return renderFence(token, md.utils.escapeHtml);
};
function customBlocksPlugin(md: MarkdownIt) {
  // 块级规则
  md.block.ruler.before('fence', 'custom_blocks', (state, startLine, endLine, silent) => {
    const buildBlockSourceAttrs = (fromLine: number, toLine: number) => {
      const sourceStartOffset = state.bMarks[fromLine] ?? 0;
      const sourceEndOffset =
        toLine < state.bMarks.length
          ? (state.bMarks[toLine] ?? state.src.length)
          : state.src.length;
      return buildSourceAttrs(
        {
          sourceStartLine: fromLine,
          sourceEndLine: toLine,
          sourceStartOffset,
          sourceEndOffset,
        },
        'block'
      );
    };

    const pos = (state.bMarks[startLine] ?? 0) + (state.tShift[startLine] ?? 0);
    const max = state.eMarks[startLine] ?? 0;
    const lineText = state.src.slice(pos, max).trim();

    // 检查是否为自定义块起始标签
    if (!lineText.startsWith(':::')) {
      return false;
    }

    // 检查是否为自闭合标签
    if (isSelfClosing(lineText)) {
      if (silent) return true;

      const { tag, params } = extractTagAndParams(lineText);

      // 处理自闭合标签
      let html = '';
      if (tag === 'link') {
        html = renderLinkCard(params, buildBlockSourceAttrs(startLine, startLine + 1));
      } else if (tag === 'video') {
        html = renderVideo(params, buildBlockSourceAttrs(startLine, startLine + 1));
      } else if (tag === 'audio') {
        html = renderAudio(params, buildBlockSourceAttrs(startLine, startLine + 1));
      } else if (tag === 'music') {
        html = renderMusic(params, buildBlockSourceAttrs(startLine, startLine + 1));
      }

      if (html) {
        const token = state.push('html_block', '', 0);
        token.content = html;
        token.map = [startLine, startLine + 1];
        state.line = startLine + 1;
        return true;
      }

      return false;
    }

    // 处理块级标签
    const { tag, params } = extractTagAndParams(lineText);
    if (!tag) return false;

    // 查找结束标签
    const endTagFull = `end${tag}`;
    let nextLine = startLine + 1;
    let foundEnd = false;
    const contentLines: string[] = [];

    // 特殊处理 tabs
    if (tag === 'tabs') {
      const tabsData: Array<{ name: string; content: string }> = [];
      let currentTab: { name: string; content: string } | null = null;

      while (nextLine < endLine) {
        const linePos = state.bMarks[nextLine] ?? 0;
        const lineMax = state.eMarks[nextLine] ?? 0;
        const line = state.src.slice(linePos, lineMax).trim();

        if (line.startsWith(':::endtabs')) {
          foundEnd = true;
          break;
        }

        if (line.startsWith(':::tab')) {
          // 保存上一个 tab
          if (currentTab) {
            tabsData.push(currentTab);
          }
          // 开始新 tab
          const tabParams = extractTagAndParams(line).params;
          currentTab = {
            name: tabParams[0] || `Tab ${tabsData.length + 1}`,
            content: '',
          };
        } else if (line.startsWith(':::endtab')) {
          // tab 结束，不做处理
        } else {
          // tab 内容
          if (currentTab) {
            currentTab.content += state.src.slice(linePos, lineMax) + '\n';
          }
        }
        nextLine++;
      }

      // 保存最后一个 tab
      if (currentTab) {
        tabsData.push(currentTab);
      }

      if (foundEnd && tabsData.length > 0) {
        if (silent) return true;

        // 渲染每个 tab 的内容
        // 注意：嵌套内容会产生错误的行号（从0开始），需要移除
        const renderedTabs = tabsData.map(tab => {
          let content = md.render(tab.content);
          // 移除嵌套块的同步属性，避免行号/偏移冲突
          content = stripNestedSyncAttrs(content);
          return { name: tab.name, content };
        });

        const html = renderTabs(
          renderedTabs,
          params,
          buildBlockSourceAttrs(startLine, nextLine + 1)
        );

        const token = state.push('html_block', '', 0);
        token.content = html;
        token.map = [startLine, nextLine + 1];
        state.line = nextLine + 1;
        return true;
      }

      return false;
    }

    // 特殊处理 photo
    if (tag === 'photo') {
      const rows: string[][] = [];
      let currentRow: string[] = [];

      while (nextLine < endLine) {
        const linePos = (state.bMarks[nextLine] ?? 0) + (state.tShift[nextLine] ?? 0);
        const lineMax = state.eMarks[nextLine] ?? 0;
        const line = state.src.slice(linePos, lineMax).trim();

        if (line === ':::endphoto') {
          foundEnd = true;
          break;
        }

        // 检查是否为换行标记 :::n
        if (line === ':::n') {
          // 保存当前行并开始新行
          if (currentRow.length > 0) {
            rows.push(currentRow);
            currentRow = [];
          }
        } else {
          // 解析图片（支持多个图片用空格分隔）
          const images = line.split(/\s+/).filter(img => img.trim());
          currentRow.push(...images);
        }

        nextLine++;
      }

      // 保存最后一行
      if (currentRow.length > 0) {
        rows.push(currentRow);
      }

      if (foundEnd && rows.length > 0) {
        if (silent) return true;

        const html = renderPhotoWall(rows, buildBlockSourceAttrs(startLine, nextLine + 1));

        const token = state.push('html_block', '', 0);
        token.content = html;
        token.map = [startLine, nextLine + 1];
        state.line = nextLine + 1;
        return true;
      }

      return false;
    }

    // 处理其他块级标签（note, fold）
    while (nextLine < endLine) {
      const linePos = state.bMarks[nextLine] ?? 0;
      const lineMax = state.eMarks[nextLine] ?? 0;
      const line = state.src.slice(linePos, lineMax).trim();

      if (line === `:::${endTagFull}`) {
        foundEnd = true;
        break;
      }

      contentLines.push(state.src.slice(linePos, lineMax));
      nextLine++;
    }

    if (!foundEnd) return false;
    if (silent) return true;

    // 渲染内容
    // 注意：嵌套内容会产生错误的行号（从0开始），需要移除
    let content = md.render(contentLines.join('\n'));
    content = stripNestedSyncAttrs(content);

    let html = '';
    if (tag === 'note') {
      html = renderNote(content, params, buildBlockSourceAttrs(startLine, nextLine + 1));
    } else if (tag === 'fold') {
      html = renderFold(content, params, buildBlockSourceAttrs(startLine, nextLine + 1));
    }

    if (html) {
      const token = state.push('html_block', '', 0);
      token.content = html;
      token.map = [startLine, nextLine + 1];
      state.line = nextLine + 1;
      return true;
    }

    return false;
  });
}

// DOMPurify 配置
const SANITIZE_CONFIG = {
  ALLOWED_TAGS: [
    'h1',
    'h2',
    'h3',
    'h4',
    'h5',
    'h6',
    'p',
    'br',
    'hr',
    'strong',
    'em',
    'u',
    's',
    'del',
    'ins',
    'mark',
    'code',
    'pre',
    'ul',
    'ol',
    'li',
    'blockquote',
    'cite',
    'footer',
    'a',
    'img',
    'table',
    'thead',
    'tbody',
    'tr',
    'th',
    'td',
    'div',
    'span',
    'sup',
    'sub',
    'kbd',
    'abbr',
    'input',
    'label',
    'button',
    'i',
    'section',
    'svg',
    'path',
    'g',
    'rect',
    'circle',
    'ellipse',
    'line',
    'polygon',
    'polyline',
    'text',
    'foreignObject',
    'video',
    'iframe',
    'audio',
    'source',
    // KaTeX / MathML 标签
    'math',
    'mrow',
    'mi',
    'mo',
    'mn',
    'msup',
    'msub',
    'msubsup',
    'mfrac',
    'msqrt',
    'mroot',
    'mover',
    'munder',
    'munderover',
    'mtable',
    'mtr',
    'mtd',
    'mtext',
    'mspace',
    'mpadded',
    'menclose',
    'mstyle',
    'merror',
    'mfenced',
    'mphantom',
    'annotation',
    'semantics',
  ],
  ALLOWED_ATTR: [
    'href',
    'title',
    'target',
    'rel',
    'src',
    'alt',
    'width',
    'height',
    'class',
    'id',
    'colspan',
    'rowspan',
    'align',
    'type',
    'checked',
    'disabled',
    'for',
    'onclick',
    'start',
    'data-source-line',
    'data-source-start-line',
    'data-source-end-line',
    'data-source-start-offset',
    'data-source-end-offset',
    'data-sync-kind',
    'd',
    'fill',
    'stroke',
    'stroke-width',
    'x',
    'y',
    'cx',
    'cy',
    'r',
    'rx',
    'ry',
    'x1',
    'y1',
    'x2',
    'y2',
    'points',
    'transform',
    'viewBox',
    'xmlns',
    'text-anchor',
    'font-size',
    'font-family',
    'dominant-baseline',
    'data-processed',
    'controls',
    'preload',
    'autoplay',
    'loop',
    'muted',
    'poster',
    'allowfullscreen',
    'scrolling',
    'border',
    'frameborder',
    'framespacing',
    'allow',
    'sandbox',
    'referrerpolicy',
    'data-server',
    'data-type',
    'data-id',
    // KaTeX / MathML 属性
    'style',
    'mathvariant',
    'mathcolor',
    'mathbackground',
    'mathsize',
    'displaystyle',
    'scriptlevel',
    'linethickness',
    'lspace',
    'rspace',
    'stretchy',
    'symmetric',
    'largeop',
    'movablelimits',
    'accent',
    'minsize',
    'maxsize',
    'open',
    'close',
    'separators',
    'notation',
    'encoding',
    'definitionurl',
    'display',
    'xmlns:xlink',
    'height',
    'depth',
    'voffset',
    'width',
    'lspace',
    'width',
    'columnalign',
    'rowalign',
    'columnspacing',
    'rowspacing',
  ],
  ALLOW_DATA_ATTR: true,
  ADD_ATTR: ['target', 'onclick', 'allowfullscreen'],
};

// 渲染 Markdown 为 HTML
export function renderMarkdown(markdown: string): string {
  if (!markdown) return '';

  const rawHtml = md.render(markdown);

  return DOMPurify.sanitize(rawHtml, SANITIZE_CONFIG);
}

// 创建带行号映射的 markdown-it 实例
function createLineNumberMd(): MarkdownIt {
  const lineMd = createMarkdownRenderer();

  lineMd.core.ruler.push('sync_source_meta', state => {
    const lineStarts = getLineStartOffsets(state.src);

    state.tokens.forEach(token => {
      if (token.map?.[0] !== undefined) {
        const sourceStartLine = token.map[0];
        const sourceEndLine = token.map[1] ?? sourceStartLine + 1;
        setTokenSourceMeta(
          token,
          sourceStartLine,
          sourceEndLine,
          getOffsetForLine(lineStarts, sourceStartLine, state.src.length),
          getOffsetForLine(lineStarts, sourceEndLine, state.src.length)
        );
      }

      if (token.type !== 'inline' || !token.children?.length || token.map?.[0] === undefined)
        return;

      const sourceStartLine = token.map[0];
      const sourceEndLine = token.map[1] ?? sourceStartLine + 1;
      const blockStartOffset = getOffsetForLine(lineStarts, sourceStartLine, state.src.length);
      const blockEndOffset = getOffsetForLine(lineStarts, sourceEndLine, state.src.length);
      let cursor = blockStartOffset;

      token.children.forEach(child => {
        if (child.type === 'image') {
          setTokenSourceMeta(
            child,
            sourceStartLine,
            sourceEndLine,
            blockStartOffset,
            blockEndOffset
          );
          return;
        }

        if (child.type !== 'text' && child.type !== 'code_inline') return;
        if (!child.content) return;

        let matchIndex = state.src.indexOf(child.content, cursor);
        if (matchIndex === -1 || matchIndex >= blockEndOffset) {
          matchIndex = state.src.indexOf(child.content, blockStartOffset);
        }
        if (matchIndex === -1 || matchIndex >= blockEndOffset) return;

        const matchEnd = Math.min(matchIndex + child.content.length, blockEndOffset);
        setTokenSourceMeta(child, sourceStartLine, sourceEndLine, matchIndex, matchEnd);
        cursor = matchEnd;
      });
    });
  });

  const applySourceAttrsToToken = (
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    token: any,
    kind: 'block' | 'text' = 'block'
  ) => {
    const meta = token?.meta;
    if (!meta) return;

    if (meta.sourceStartLine !== undefined) {
      token.attrSet('data-source-line', String(meta.sourceStartLine));
      token.attrSet('data-source-start-line', String(meta.sourceStartLine));
    }
    if (meta.sourceEndLine !== undefined)
      token.attrSet('data-source-end-line', String(meta.sourceEndLine));
    if (meta.sourceStartOffset !== undefined)
      token.attrSet('data-source-start-offset', String(meta.sourceStartOffset));
    if (meta.sourceEndOffset !== undefined)
      token.attrSet('data-source-end-offset', String(meta.sourceEndOffset));
    token.attrSet('data-sync-kind', kind);
  };

  lineMd.renderer.rules.fence = (tokens, idx) => {
    const token = tokens[idx];
    if (!token) return '';
    return renderFence(token, lineMd.utils.escapeHtml, buildSourceAttrs(token.meta, 'block'));
  };

  const blockTags = [
    'heading_open',
    'blockquote_open',
    'bullet_list_open',
    'ordered_list_open',
    'list_item_open',
    'table_open',
    'hr',
  ];

  blockTags.forEach(tag => {
    const originalRule = lineMd.renderer.rules[tag];
    lineMd.renderer.rules[tag] = (tokens, idx, options, env, self) => {
      const token = tokens[idx];
      if (token) applySourceAttrsToToken(token, 'block');
      return originalRule
        ? originalRule(tokens, idx, options, env, self)
        : self.renderToken(tokens, idx, options);
    };
  });

  const originalImageRule = lineMd.renderer.rules.image;
  lineMd.renderer.rules.image = (tokens, idx, options, env, self) => {
    const token = tokens[idx];
    if (token) {
      applySourceAttrsToToken(token, 'block');
      token.attrJoin('class', 'preview-collapsible-image');
    }
    return originalImageRule
      ? originalImageRule(tokens, idx, options, env, self)
      : self.renderToken(tokens, idx, options);
  };

  const originalParagraphOpen = lineMd.renderer.rules.paragraph_open;
  lineMd.renderer.rules.paragraph_open = (tokens, idx, options, env, self) => {
    const token = tokens[idx];
    if (token) applySourceAttrsToToken(token, 'block');
    return originalParagraphOpen
      ? originalParagraphOpen(tokens, idx, options, env, self)
      : self.renderToken(tokens, idx, options);
  };

  lineMd.renderer.rules.text = (tokens, idx) => {
    const token = tokens[idx];
    if (!token) return '';
    return renderAnchoredText(
      token.content,
      token.meta?.sourceStartOffset,
      token.meta?.sourceEndOffset,
      lineMd.utils.escapeHtml
    );
  };

  lineMd.renderer.rules.code_inline = (tokens, idx) => {
    const token = tokens[idx];
    if (!token) return '';
    return `<code${buildSourceAttrs(token.meta, 'text')}>${lineMd.utils.escapeHtml(token.content)}</code>`;
  };

  return lineMd;
}

const lineMd = createLineNumberMd();

// 渲染 Markdown 为带行号映射的 HTML（用于滚动同步）
export function renderMarkdownWithSourceMap(markdown: string): string {
  if (!markdown) return '';

  const rawHtml = lineMd.render(markdown);

  return DOMPurify.sanitize(rawHtml, SANITIZE_CONFIG);
}

// Markdown 内容样式（从 _prose.scss 提取的核心样式）
const MARKDOWN_STYLES = `
.markdown-content { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', sans-serif; font-size: 1rem; line-height: 1.8; color: #4c4948; max-width: 800px; margin: 0 auto; word-wrap: break-word; }
.markdown-content h1, .markdown-content h2, .markdown-content h3, .markdown-content h4, .markdown-content h5, .markdown-content h6 { margin: 20px 0; font-weight: 600; line-height: 1.3; color: #2c3e50; }
.markdown-content h1 { font-size: 2.25em; margin-top: 0; } .markdown-content h2 { font-size: 1.875em; } .markdown-content h3 { font-size: 1.5em; } .markdown-content h4 { font-size: 1.25em; } .markdown-content h5 { font-size: 1.125em; } .markdown-content h6 { font-size: 1em; }
.markdown-content p { margin: 1.25em 0; text-align: justify; }
.markdown-content a { color: #49b1f5; text-decoration: none; font-weight: 500; }
.markdown-content ul, .markdown-content ol { padding-left: 2em; margin: 1.25em 0; } .markdown-content li { margin: 0.5em 0; line-height: 1.75; }
.markdown-content blockquote { margin: 1.5em 0; padding: 0.75em 1.25em; background: linear-gradient(to right, rgba(73, 177, 245, 0.05), transparent); border-left: 4px solid #49b1f5; border-radius: 0 6px 6px 0; color: #5a6c7d; font-size: 0.95em; } .markdown-content blockquote p { margin: 0.75em 0; }
.markdown-content code { padding: 0.2em 0.5em; margin: 0; font-size: 0.875em; background-color: rgba(73, 177, 245, 0.08); border: 1px solid rgba(73, 177, 245, 0.15); border-radius: 4px; font-family: 'Consolas', 'Monaco', 'Courier New', monospace; color: #e74c3c; font-weight: 500; }
.markdown-content .code-block-container { margin: 1.5em 0; border: 1px solid #e1e4e8; border-radius: 6px; background: #f6f8fa; overflow: hidden; }
.markdown-content .code-block-container .code-toolbar { display: flex; align-items: center; gap: 0.75em; padding: 0.3em 1em; border-bottom: 1px solid #e1e4e8; }
.markdown-content .code-block-container .code-toolbar .code-fold-btn { display: flex; align-items: center; justify-content: center; padding: 0; background: transparent; border: none; cursor: pointer; transition: color 0.2s ease; }
.markdown-content .code-block-container .code-toolbar .code-fold-btn i { font-size: 24px; font-style: normal; color: #586069; transition: color 0.2s ease, transform 0.3s ease; }
.markdown-content .code-block-container .code-toolbar .code-fold-btn:hover i { color: #49b1f5; }
.markdown-content .code-block-container .code-toolbar .code-lang { font-weight: 600; color: #586069; line-height: 1; user-select: none; }
.markdown-content .code-block-container .code-toolbar .code-copy-btn { margin-left: auto; display: flex; align-items: center; justify-content: center; padding: 0; background: transparent; border: none; cursor: pointer; transition: color 0.2s ease; }
.markdown-content .code-block-container .code-toolbar .code-copy-btn i { font-size: 22px; font-style: normal; color: #586069; transition: color 0.2s ease; }
.markdown-content .code-block-container .code-toolbar .code-copy-btn:hover i { color: #49b1f5; }
.markdown-content .code-block-container .code-toolbar .code-copy-btn.copied i { color: #4caf50; }
.markdown-content .code-block-container pre { margin: 0; padding: 0; background: transparent; border: none; border-radius: 0; }
.markdown-content .code-block-container pre code { display: block; padding: 1em; overflow-x: auto; background: none; border: none; color: #24292e; line-height: 1.45; font-size: 0.875em; font-family: 'Consolas', 'Monaco', 'Courier New', monospace; }
.markdown-content .code-block-container pre code .line-number { display: inline-block; width: 2em; padding-right: 1em; margin-right: 1em; text-align: right; color: #858585; border-right: 1px solid #d1d5da; user-select: none; }
.markdown-content .code-block-container pre code .line-number::before { content: attr(data-line); }
.markdown-content .code-block-container pre code .line-content { display: inline; white-space: pre; }
.markdown-content .code-block-container.collapsed .code-fold-btn i { transform: rotate(-90deg); }
.markdown-content .code-block-container.collapsed pre { display: none; }
.markdown-content table { border-collapse: collapse; width: 100%; margin: 1.5em 0; font-size: 0.9em; } .markdown-content th, .markdown-content td { padding: 0.75em 1em; border: 1px solid rgba(128, 128, 128, 0.15); text-align: left; } .markdown-content th { font-weight: 600; color: #2c3e50; background: linear-gradient(135deg, rgba(73, 177, 245, 0.1), rgba(73, 177, 245, 0.05)); }
.markdown-content hr { height: 2px; margin: 2.5em 0; background: linear-gradient(to right, transparent, rgba(73, 177, 245, 0.2), rgba(73, 177, 245, 0.6), rgba(73, 177, 245, 0.2), transparent); border: 0; }
.markdown-content img { max-width: 100%; height: auto; border-radius: 8px; margin: 1.5em auto; display: block; }
.markdown-content strong, .markdown-content b { font-weight: 600; color: #2c3e50; } .markdown-content em, .markdown-content i { font-style: italic; color: #5a6c7d; } .markdown-content del, .markdown-content s { text-decoration: line-through; opacity: 0.7; color: #858585; }
.markdown-content u { text-decoration: underline; text-decoration-color: #49b1f5; } .markdown-content mark { background: linear-gradient(to bottom, transparent 50%, rgba(255, 235, 59, 0.5) 50%); padding: 0.1em 0.2em; border-radius: 2px; }
.markdown-content sup { font-size: 0.75em; vertical-align: super; } .markdown-content sub { font-size: 0.75em; vertical-align: sub; }
.markdown-content kbd { display: inline-block; padding: 0.2em 0.5em; font-size: 0.875em; font-family: 'Consolas', 'Monaco', monospace; background: linear-gradient(to bottom, #f7f7f7, #e8e8e8); border: 1px solid #ccc; border-radius: 4px; box-shadow: 0 2px 0 rgba(0, 0, 0, 0.1); color: #333; }
.markdown-content .custom-note { margin: 1.5em 0; padding: 1em 1.2em; border-radius: 8px; background: rgba(249, 250, 251, 0.5); border-left: 4px solid; } .markdown-content .custom-note-title { font-weight: 700; font-size: 1.05em; margin-bottom: 0.6em; }
.markdown-content .custom-note-info { border-left-color: #2196f3; } .markdown-content .custom-note-info .custom-note-title { color: #2196f3; } .markdown-content .custom-note-warning { border-left-color: #ff9800; } .markdown-content .custom-note-warning .custom-note-title { color: #ff9800; } .markdown-content .custom-note-success { border-left-color: #4caf50; } .markdown-content .custom-note-success .custom-note-title { color: #4caf50; } .markdown-content .custom-note-error { border-left-color: #f44336; } .markdown-content .custom-note-error .custom-note-title { color: #f44336; }
.markdown-content .custom-tabs { margin: 1.5em 0; padding: 12px; border-radius: 8px; background: rgba(249, 250, 251, 0.5); } .markdown-content .custom-tabs-header { display: flex; gap: 8px; margin-bottom: 12px; padding-bottom: 12px; border-bottom: 1px solid rgba(128, 128, 128, 0.2); }
.markdown-content .custom-tab-btn { padding: 6px 14px; background: transparent; border: none; border-radius: 4px; color: #858585; font-size: 0.9rem; font-weight: 500; cursor: pointer; } .markdown-content .custom-tab-btn.active { color: #fff; background: #49b1f5; font-weight: 600; }
.markdown-content .custom-tab-panel { display: none; } .markdown-content .custom-tab-panel.active { display: block; }
.markdown-content .custom-fold { margin: 1.5em 0; border-radius: 8px; background: rgba(249, 250, 251, 0.5); overflow: hidden; } .markdown-content .custom-fold-header { padding: 8px 16px; font-weight: 700; cursor: pointer; display: flex; align-items: center; gap: 8px; }
.markdown-content .custom-fold-content { max-height: 0; overflow: hidden; } .markdown-content .custom-fold.open .custom-fold-content { max-height: 800px; } .markdown-content .custom-fold-content > div { padding: 0 16px 16px; }
.markdown-content .custom-link-card { margin: 1.5em 0; border-radius: 8px; background: rgba(249, 250, 251, 0.5); overflow: hidden; } .markdown-content .custom-link-type { padding: 6px 16px; font-size: 0.75em; font-weight: 600; color: #858585; background: rgba(73, 177, 245, 0.05); border-bottom: 1px solid rgba(128, 128, 128, 0.2); }
.markdown-content .custom-link-main { display: flex; align-items: center; gap: 12px; padding: 12px 16px; text-decoration: none; color: inherit; } .markdown-content .custom-link-icon { flex-shrink: 0; width: 64px; height: 64px; display: flex; align-items: center; justify-content: center; border-radius: 8px; background: rgba(249, 250, 251, 0.5); border: 1px solid rgba(128, 128, 128, 0.2); }
.markdown-content .custom-link-info { flex: 1; min-width: 0; } .markdown-content .custom-link-title { font-weight: 600; font-size: 1.1em; color: #2c3e50; } .markdown-content .custom-link-desc { font-size: 0.875em; color: #5a6c7d; line-height: 1.5; }
.markdown-content .custom-photo-wall{margin:1em 0;border-radius:8px;overflow-x:auto}.markdown-content .custom-photo-wall-container{display:flex;flex-direction:column}.markdown-content .custom-photo-wall-row{display:flex;align-items:stretch;flex-wrap:nowrap}.markdown-content .custom-photo-wall-item{flex:1 1 0;min-width:0;display:flex;align-items:center;justify-content:center;padding:5px;height:100%}.markdown-content .custom-photo-wall-item img{margin:0;display:block;max-width:100%;max-height:100%;object-fit:contain;border-radius:6px;background:rgba(249,250,251,.8)}
.markdown-content .katex-inline{display:inline}.markdown-content .katex-block{display:block;margin:1.5rem 0;text-align:center;overflow-x:auto;padding:0.5rem 0}.markdown-content .katex-block .katex{font-size:1.15em}.markdown-content .katex{font-size:1em;line-height:1.6}.markdown-content .katex .base{color:inherit}.markdown-content .katex .katex-mathml{position:absolute;clip:rect(1px,1px,1px,1px);padding:0;border:0;height:1px;width:1px;overflow:hidden}.markdown-content .katex-error{color:#cc0000;font-style:italic}
`;

// 渲染 Markdown 为带样式的完整 HTML（用于复制）
export function renderMarkdownWithStyles(markdown: string): string {
  if (!markdown) return '';

  const html = renderMarkdown(markdown);
  const script = `;(function(){function f(t){var e=t.closest('.code-block-container');if(!e)return'';var n=e.querySelector('code');if(!n)return'';var r=Array.from(n.querySelectorAll('.line-content'));return r.map(function(o){return o.textContent||''}).join('\\n')}function c(t,e){try{if(navigator.clipboard&&navigator.clipboard.writeText){return navigator.clipboard.writeText(t).then(e)} }catch(o){}var n=document.createElement('textarea');n.value=t;n.style.position='fixed';n.style.opacity='0';document.body.appendChild(n);n.select();try{document.execCommand('copy')}catch(o){}document.body.removeChild(n);e&&e()}function copyCodeBlock(t){var e=f(t);if(!e)return;c(e,function(){var n=t.querySelector('i');if(n){n.className='ri-check-line';t.classList.add('copied')}setTimeout(function(){if(n){n.className='ri-file-copy-fill';t.classList.remove('copied')}},2000)})}function switchTab(t,e){var n=document.getElementById(t);if(!n)return;Array.from(n.querySelectorAll('.custom-tab-btn')).forEach(function(r){r.textContent===e?r.classList.add('active'):r.classList.remove('active')});Array.from(n.querySelectorAll('.custom-tab-panel')).forEach(function(r){var o=r; o.dataset&&o.dataset.tab===e?r.classList.add('active'):r.classList.remove('active')})}function toggleFold(t){var e=document.getElementById(t);if(!e)return;e.classList.toggle('open')}window.copyCodeBlock=copyCodeBlock;window.switchTab=switchTab;window.toggleFold=toggleFold})();`;

  return `<style>${MARKDOWN_STYLES}</style><div class="markdown-content">${html}</div>\n<script>${script}</script>\n`;
}

// 计算字数
export function countWords(markdown: string): number {
  if (!markdown) return 0;

  const text = markdown
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/~~~[\s\S]*?~~~/g, ' ')
    .replace(/^:::note[\s\S]*?^:::endnote$/gm, ' ')
    .replace(/^:::tabs[\s\S]*?^:::endtabs$/gm, ' ')
    .replace(/^:::fold[\s\S]*?^:::endfold$/gm, ' ')
    .replace(/^:::photo[\s\S]*?^:::endphoto$/gm, ' ')
    .replace(/^:::link\s+.*?:::\s*$/gm, ' ')
    .replace(/^:::video\s+.*?:::\s*$/gm, ' ')
    .replace(/^:::audio\s+.*?:::\s*$/gm, ' ')
    .replace(/!\[[^\]]*\]\([^)]*\)/g, ' ')
    .replace(/\[[^\]]+\]\([^)]*\)/g, '$1')
    .replace(/<[^>]+>/g, ' ')
    .replace(/[`*_~>#\-|]/g, ' ')
    .replace(/\s+/g, ' ')
    .trim();

  const chineseChars = text.match(/[\u4e00-\u9fa5]/g) || [];
  const englishWords = text.match(/[a-zA-Z]+/g) || [];
  return chineseChars.length + englishWords.length;
}

// 计算阅读时长（分钟）
export function estimateReadingTime(markdown: string, wordsPerMinute = 300): number {
  return Math.ceil(countWords(markdown) / wordsPerMinute);
}

// 目录项接口
export interface TocItem {
  id: string;
  level: number;
  text: string;
  children?: TocItem[];
}

// 提取目录
export function extractToc(markdown: string): TocItem[] {
  if (!markdown) return [];

  // 移除代码块
  let cleanedMarkdown = markdown
    .replace(/```[\s\S]*?```/g, '')
    .replace(/~~~[\s\S]*?~~~\s*/g, '')
    .replace(/^( {4}|\t).+$/gm, '');

  // 处理单行自定义块
  cleanedMarkdown = cleanedMarkdown.replace(/^:::link\s+.*?:::$/gm, '');
  cleanedMarkdown = cleanedMarkdown.replace(/^:::video\s+.*?:::$/gm, '');
  cleanedMarkdown = cleanedMarkdown.replace(/^:::audio\s+.*?:::$/gm, '');
  // 处理多行自定义块
  cleanedMarkdown = cleanedMarkdown.replace(/^:::note[\s\S]*?^:::endnote$/gm, '');
  cleanedMarkdown = cleanedMarkdown.replace(/^:::tabs[\s\S]*?^:::endtabs$/gm, '');
  cleanedMarkdown = cleanedMarkdown.replace(/^:::fold[\s\S]*?^:::endfold$/gm, '');
  cleanedMarkdown = cleanedMarkdown.replace(/^:::photo[\s\S]*?^:::endphoto$/gm, '');

  const headings: TocItem[] = [];

  for (const line of cleanedMarkdown.split('\n')) {
    const match = line.match(/^(#{1,6})\s+(.+)$/);
    if (match?.[1] && match[2]) {
      headings.push({
        id: generateHeadingId(match[2].trim()),
        level: match[1].length,
        text: match[2].trim(),
      });
    }
  }

  return headings;
}

// 简单 Markdown 渲染（用于评论）
export function renderSimpleMarkdown(markdown: string): string {
  if (!markdown) return '';

  const simpleMd = new MarkdownIt({
    html: false,
    breaks: true,
    linkify: true,
  });

  const simpleHtml = simpleMd.render(markdown);

  return DOMPurify.sanitize(simpleHtml, {
    ALLOWED_TAGS: ['p', 'br', 'strong', 'em', 'code', 'pre', 'ul', 'ol', 'li', 'blockquote', 'a'],
    ALLOWED_ATTR: ['href', 'title'],
    ALLOW_DATA_ATTR: false,
  });
}

// 复制代码块功能
export function copyCodeBlock(button: HTMLElement): void {
  const container = button.closest('.code-block-container');
  if (!container) return;

  const code = container.querySelector('code');
  if (!code) return;

  // 只提取代码内容，不包含行号
  const codeLines = Array.from(code.querySelectorAll('.line-content'));
  const codeText = codeLines.map(line => line.textContent || '').join('\n');

  // 复制到剪贴板
  navigator.clipboard
    .writeText(codeText)
    .then(() => {
      // 更新按钮状态
      const icon = button.querySelector('i');
      if (icon) {
        icon.className = 'ri-check-line';
        button.classList.add('copied');
      }

      // 2秒后恢复
      setTimeout(() => {
        if (icon) {
          icon.className = 'ri-file-copy-fill';
          button.classList.remove('copied');
        }
      }, 2000);
    })
    .catch(err => {
      console.error('复制失败:', err);
    });
}

// 标签页切换功能
export function switchTab(tabsId: string, tabName: string): void {
  const tabsContainer = document.getElementById(tabsId);
  if (!tabsContainer) return;

  // 更新标签按钮状态
  const buttons = tabsContainer.querySelectorAll('.custom-tab-btn');
  buttons.forEach(btn => {
    if (btn.textContent === tabName) {
      btn.classList.add('active');
    } else {
      btn.classList.remove('active');
    }
  });

  // 更新内容面板状态
  const panels = tabsContainer.querySelectorAll('.custom-tab-panel');
  panels.forEach(panel => {
    const panelElement = panel as HTMLElement;
    if (panelElement.dataset.tab === tabName) {
      panel.classList.add('active');
    } else {
      panel.classList.remove('active');
    }
  });
}

// 折叠面板切换功能
export function toggleFold(foldId: string): void {
  const foldContainer = document.getElementById(foldId);
  if (!foldContainer) return;

  const isOpening = !foldContainer.classList.contains('open');
  const contentDiv = foldContainer.querySelector('.custom-fold-content > div') as HTMLElement;

  if (isOpening && contentDiv) {
    // 展开时：先获取内容的实际高度
    const contentHeight = contentDiv.scrollHeight;
    // 设置 max-height 为内容的实际高度
    const contentContainer = foldContainer.querySelector('.custom-fold-content') as HTMLElement;
    if (contentContainer) {
      contentContainer.style.maxHeight = `${contentHeight}px`;
    }
  } else {
    // 折叠时：重置 max-height 为 0
    const contentContainer = foldContainer.querySelector('.custom-fold-content') as HTMLElement;
    if (contentContainer) {
      contentContainer.style.maxHeight = '0px';
    }
  }

  foldContainer.classList.toggle('open');
}

function formatTime(seconds: number): string {
  if (isNaN(seconds) || !isFinite(seconds)) return '0:00';
  const mins = Math.floor(seconds / 60);
  const secs = Math.floor(seconds % 60);
  return `${mins}:${secs.toString().padStart(2, '0')}`;
}

export function toggleAudioPlay(audioId: string): void {
  const container = document.querySelector(`[data-audio-id="${audioId}"]`);
  if (!container) return;

  const audio = container.querySelector('audio') as HTMLAudioElement;
  const btn = container.querySelector('.custom-audio-btn i');
  const durationTime = container.querySelector('.custom-audio-duration');
  if (!audio || !btn) return;

  if (!(container as HTMLElement).dataset.audioInitialized) {
    initAudioEvents(container);
    (container as HTMLElement).dataset.audioInitialized = 'true';
  }

  // 手动更新时长（以防错过 loadedmetadata 事件）
  if (durationTime && audio.duration && isFinite(audio.duration)) {
    const mins = Math.floor(audio.duration / 60);
    const secs = Math.floor(audio.duration % 60);
    durationTime.textContent = `${mins}:${secs.toString().padStart(2, '0')}`;
  }

  if (audio.paused) {
    audio.play();
    btn.className = 'ri-pause-fill';
  } else {
    audio.pause();
    btn.className = 'ri-play-fill';
  }
}

export function seekAudio(audioId: string, event: MouseEvent): void {
  const container = document.querySelector(`[data-audio-id="${audioId}"]`);
  if (!container) return;

  const audio = container.querySelector('audio') as HTMLAudioElement;
  const progressBar = container.querySelector('.custom-audio-progress');
  const durationTime = container.querySelector('.custom-audio-duration');
  if (!audio || !progressBar) return;

  if (!(container as HTMLElement).dataset.audioInitialized) {
    initAudioEvents(container);
    (container as HTMLElement).dataset.audioInitialized = 'true';
  }

  // 预加载元数据并更新时长显示
  if (durationTime && (!audio.duration || !isFinite(audio.duration))) {
    audio.preload = 'metadata';
    audio.addEventListener(
      'loadedmetadata',
      () => {
        const mins = Math.floor(audio.duration / 60);
        const secs = Math.floor(audio.duration % 60);
        durationTime.textContent = `${mins}:${secs.toString().padStart(2, '0')}`;
      },
      { once: true }
    );
  }

  const rect = progressBar.getBoundingClientRect();
  const percent = (event.clientX - rect.left) / rect.width;
  audio.currentTime = percent * audio.duration;
}

export function toggleMusicPlay(audioId: string, _server: string, _musicId: string): void {
  const container = document.querySelector(`[data-audio-id="${audioId}"]`);
  if (!container) return;

  const btn = container.querySelector('.custom-audio-btn i');
  const musicInfoEl = container.querySelector(`[data-music-info="${audioId}"]`);
  const durationTime = container.querySelector('.custom-audio-duration');
  if (!btn) return;

  let audio = container.querySelector('audio') as HTMLAudioElement;
  const musicSource = container.querySelector('.custom-audio-source') as HTMLElement;

  if (!audio && musicSource) {
    const embedUrl = musicSource.dataset.embedUrl || '';

    fetch(embedUrl)
      .then(res => res.json())
      .then(data => {
        if (data && data.length > 0) {
          const info = data[0];
          if (musicInfoEl) {
            musicInfoEl.textContent = `${info.name || '未知歌曲'} - ${info.artist || info.author || '未知艺术家'}`;
          }

          audio = container.querySelector('audio') as HTMLAudioElement;
          if (!audio && info.url) {
            const newAudio = document.createElement('audio');
            newAudio.src = info.url;
            newAudio.preload = 'auto';
            container.appendChild(newAudio);
            audio = newAudio;
          } else if (audio) {
            audio.src = info.url;
          }

          if (audio && durationTime) {
            audio.addEventListener(
              'loadedmetadata',
              () => {
                if (durationTime && audio.duration && isFinite(audio.duration)) {
                  const mins = Math.floor(audio.duration / 60);
                  const secs = Math.floor(audio.duration % 60);
                  durationTime.textContent = `${mins}:${secs.toString().padStart(2, '0')}`;
                }
              },
              { once: true }
            );
          }

          initMusicEvents(container);
          audio.play();
          btn.className = 'ri-pause-fill';
        }
      })
      .catch(err => {
        console.error('Failed to load music:', err);
        if (musicInfoEl) {
          musicInfoEl.textContent = '加载失败';
        }
      });
    return;
  }

  if (audio) {
    if (audio.paused) {
      audio.play();
      btn.className = 'ri-pause-fill';
    } else {
      audio.pause();
      btn.className = 'ri-play-fill';
    }
  }
}

export function seekMusic(audioId: string, event: MouseEvent): void {
  const container = document.querySelector(`[data-audio-id="${audioId}"]`);
  if (!container) return;

  const audio = container.querySelector('audio') as HTMLAudioElement;
  const progressBar = container.querySelector('.custom-audio-progress');
  const durationTime = container.querySelector('.custom-audio-duration');
  if (!audio || !progressBar) return;

  initMusicEvents(container);

  if (durationTime && (!audio.duration || !isFinite(audio.duration))) {
    audio.preload = 'metadata';
    audio.addEventListener(
      'loadedmetadata',
      () => {
        if (durationTime && audio.duration && isFinite(audio.duration)) {
          const mins = Math.floor(audio.duration / 60);
          const secs = Math.floor(audio.duration % 60);
          durationTime.textContent = `${mins}:${secs.toString().padStart(2, '0')}`;
        }
      },
      { once: true }
    );
  }

  const rect = progressBar.getBoundingClientRect();
  const percent = (event.clientX - rect.left) / rect.width;
  audio.currentTime = percent * audio.duration;
}

function initMusicEvents(container: Element): void {
  if ((container as HTMLElement).dataset.musicInitialized) return;
  (container as HTMLElement).dataset.musicInitialized = 'true';

  const audio = container.querySelector('audio') as HTMLAudioElement;
  const progressBar = container.querySelector('.custom-audio-progress-bar') as HTMLElement;
  const currentTimeEl = container.querySelector('.custom-audio-current') as HTMLElement;
  const durationTime = container.querySelector('.custom-audio-duration') as HTMLElement;
  const btn = container.querySelector('.custom-audio-btn i');

  if (!audio) return;

  audio.load();

  const updateDuration = () => {
    if (durationTime && audio.duration && isFinite(audio.duration)) {
      durationTime.textContent = formatTime(audio.duration);
    }
  };

  audio.addEventListener('loadedmetadata', updateDuration);
  audio.addEventListener('canplay', updateDuration);

  audio.addEventListener('timeupdate', () => {
    if (currentTimeEl) {
      currentTimeEl.textContent = formatTime(audio.currentTime);
    }
    if (progressBar && audio.duration && isFinite(audio.duration)) {
      progressBar.style.width = `${(audio.currentTime / audio.duration) * 100}%`;
    }
    if (durationTime && audio.duration && isFinite(audio.duration)) {
      durationTime.textContent = formatTime(audio.duration);
    }
  });

  audio.addEventListener('ended', () => {
    if (btn) btn.className = 'ri-play-fill';
    if (progressBar) progressBar.style.width = '0%';
    if (currentTimeEl) currentTimeEl.textContent = '0:00';
  });
}

function initAudioEvents(container: Element): void {
  const audio = container.querySelector('audio') as HTMLAudioElement;
  const progressBar = container.querySelector('.custom-audio-progress-bar') as HTMLElement;
  const currentTimeEl = container.querySelector('.custom-audio-current') as HTMLElement;
  const durationTime = container.querySelector('.custom-audio-duration') as HTMLElement;

  if (!audio) return;
  if ((container as HTMLElement).dataset.audioInitialized) return;
  (container as HTMLElement).dataset.audioInitialized = 'true';

  audio.load();

  const updateDuration = () => {
    if (durationTime && audio.duration && isFinite(audio.duration)) {
      durationTime.textContent = formatTime(audio.duration);
    }
  };

  audio.addEventListener('loadedmetadata', updateDuration);
  audio.addEventListener('canplay', updateDuration);

  audio.addEventListener('timeupdate', () => {
    if (progressBar && audio.duration && isFinite(audio.duration)) {
      progressBar.style.width = `${(audio.currentTime / audio.duration) * 100}%`;
    }
    if (durationTime && audio.duration && isFinite(audio.duration)) {
      durationTime.textContent = formatTime(audio.duration);
    }
    if (currentTimeEl) currentTimeEl.textContent = formatTime(audio.currentTime);
  });

  audio.addEventListener('ended', () => {
    const btn = container.querySelector('.custom-audio-btn i') as HTMLElement;
    if (btn) btn.className = 'ri-play-fill';
    if (progressBar) progressBar.style.width = '0%';
    if (currentTimeEl) currentTimeEl.textContent = '0:00';
  });
}

// 使用 MutationObserver 自动初始化动态添加的音频播放器
function observeAudioPlayers(): void {
  const observer = new MutationObserver(mutations => {
    for (const mutation of mutations) {
      for (const node of mutation.addedNodes) {
        if (node instanceof Element) {
          // 初始化音频卡片
          const audioContainers = node.matches?.('[data-audio-id]')
            ? [node]
            : node.querySelectorAll?.('[data-audio-id]');
          audioContainers?.forEach(container => {
            initAudioEvents(container);
            initMusicCard(container);
          });
        }
      }
    }
  });

  observer.observe(document.body, { childList: true, subtree: true });

  // 初始化已存在的音频播放器
  document.querySelectorAll('.custom-audio[data-audio-id]').forEach(container => {
    initAudioEvents(container);
    initMusicCard(container);
  });
}

// 初始化音乐卡片 - 预加载音乐信息
function initMusicCard(container: Element): void {
  const musicSource = container.querySelector('.custom-audio-source') as HTMLElement;
  const musicInfoEl = container.querySelector('.custom-audio-info') as HTMLElement;
  if (!musicSource || musicInfoEl?.dataset.initialized) return;

  const embedUrl = musicSource.dataset.embedUrl;
  if (!embedUrl) return;

  musicInfoEl.dataset.initialized = 'true';

  fetch(embedUrl)
    .then(res => res.json())
    .then(data => {
      if (data && data.length > 0) {
        const info = data[0];
        if (musicInfoEl) {
          musicInfoEl.textContent = `${info.name || '未知歌曲'} - ${info.artist || info.author || '未知艺术家'}`;
        }
        // 创建 audio 元素但不播放
        const existingAudio = container.querySelector('audio');
        if (!existingAudio && info.url) {
          const audio = document.createElement('audio');
          audio.src = info.url;
          audio.preload = 'auto';
          audio.style.display = 'none';
          container.appendChild(audio);
          initMusicEvents(container);
        }
      }
    })
    .catch(() => {
      if (musicInfoEl) {
        musicInfoEl.textContent = '加载失败';
      }
    });
}

// 挂载全局函数供内联 onclick 使用
if (typeof window !== 'undefined') {
  const win = window as unknown as Window & {
    copyCodeBlock: typeof copyCodeBlock;
    switchTab: typeof switchTab;
    toggleFold: typeof toggleFold;
    toggleAudioPlay: typeof toggleAudioPlay;
    seekAudio: typeof seekAudio;
    toggleMusicPlay: typeof toggleMusicPlay;
    seekMusic: typeof seekMusic;
  };
  win.copyCodeBlock = copyCodeBlock;
  win.switchTab = switchTab;
  win.toggleFold = toggleFold;
  win.toggleAudioPlay = toggleAudioPlay;
  win.seekAudio = seekAudio;
  win.toggleMusicPlay = toggleMusicPlay;
  win.seekMusic = seekMusic;

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', observeAudioPlayers);
  } else {
    observeAudioPlayers();
  }
}

export default {
  render: renderMarkdown,
  renderSimple: renderSimpleMarkdown,
  countWords,
  estimateReadingTime,
  extractToc,
  copyCodeBlock,
};
