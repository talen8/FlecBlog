/**
 * Markdown 分段解析器
 * 将 Markdown 内容解析为多个段落，支持自定义块
 */

/**
 * 内容段落类型
 */
export type BlockType =
  | 'markdown'
  | 'note'
  | 'fold'
  | 'tabs'
  | 'link'
  | 'video'
  | 'audio'
  | 'music'
  | 'photo';

/**
 * 基础段落接口
 */
export interface BaseBlock {
  type: BlockType;
}

/**
 * Markdown 段落
 */
export interface MarkdownBlock extends BaseBlock {
  type: 'markdown';
  content: string;
}

/**
 * 提示框段落
 */
export interface NoteBlock extends BaseBlock {
  type: 'note';
  noteType: string;
  title: string;
  content: string;
}

/**
 * 折叠面板段落
 */
export interface FoldBlock extends BaseBlock {
  type: 'fold';
  title: string;
  open: boolean;
  content: string;
}

/**
 * 标签页段落
 */
export interface TabsBlock extends BaseBlock {
  type: 'tabs';
  tabs: Array<{ name: string; content: string }>;
  activeTab: string;
}

/**
 * 链接卡片段落
 */
export interface LinkBlock extends BaseBlock {
  type: 'link';
  title: string;
  link: string;
  description: string;
}

/**
 * 视频段落
 */
export interface VideoBlock extends BaseBlock {
  type: 'video';
  platform: string;
  videoId: string;
}

/**
 * 音频段落
 */
export interface AudioBlock extends BaseBlock {
  type: 'audio';
  title: string;
  audioUrl: string;
}

/**
 * 在线音乐段落
 */
export interface MusicBlock extends BaseBlock {
  type: 'music';
  server: string;
  musicId: string;
}

/**
 * 照片墙段落
 */
export interface PhotoBlock extends BaseBlock {
  type: 'photo';
  rows: string[][];
}

/**
 * 所有段落类型联合
 */
export type ContentBlock =
  | MarkdownBlock
  | NoteBlock
  | FoldBlock
  | TabsBlock
  | LinkBlock
  | VideoBlock
  | AudioBlock
  | MusicBlock
  | PhotoBlock;

/**
 * 提取标签名和参数
 * @param line - 完整的标签行
 * @returns 标签名和参数数组
 */
function extractTagAndParams(line: string): { tag: string; params: string[] } {
  const match = line.match(/^:::(\w+)(.*)$/);
  if (!match) return { tag: '', params: [] };
  const tag = match[1] || '';
  const paramsString = match[2]?.trim() || '';
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

/**
 * 解析 Markdown 内容为段落数组
 * @param markdown - 原始 Markdown 内容
 * @returns 段落数组
 */
export function parseMarkdownToBlocks(markdown: string): ContentBlock[] {
  if (!markdown) return [];

  const lines = markdown.split('\n');
  const blocks: ContentBlock[] = [];
  let currentMarkdown: string[] = [];
  let i = 0;

  /**
   * 保存当前累积的 Markdown 内容
   */
  function flushMarkdown() {
    const content = currentMarkdown.join('\n').trim();
    if (content) {
      blocks.push({ type: 'markdown', content });
    }
    currentMarkdown = [];
  }

  while (i < lines.length) {
    const line = lines[i] || '';
    const trimmedLine = line.trim();

    if (!trimmedLine.startsWith(':::')) {
      currentMarkdown.push(line);
      i++;
      continue;
    }

    flushMarkdown();

    if (isSelfClosing(trimmedLine)) {
      const { tag, params } = extractTagAndParams(trimmedLine);

      if (tag === 'link') {
        blocks.push({
          type: 'link',
          title: params[0] || '',
          link: params[1] || '',
          description: params.slice(2).join(' '),
        });
      } else if (tag === 'video') {
        blocks.push({
          type: 'video',
          platform: params[0] || '',
          videoId: params[1] || '',
        });
      } else if (tag === 'audio') {
        blocks.push({
          type: 'audio',
          title: params[0] || '',
          audioUrl: params[1] || '',
        });
      } else if (tag === 'music') {
        blocks.push({
          type: 'music',
          server: params[0] || '',
          musicId: params[1] || '',
        });
      }
      i++;
      continue;
    }

    const { tag, params } = extractTagAndParams(trimmedLine);
    if (!tag) {
      currentMarkdown.push(line);
      i++;
      continue;
    }

    if (tag === 'tabs') {
      const tabs: Array<{ name: string; content: string }> = [];
      let currentTab: { name: string; content: string } | null = null;
      i++;

      while (i < lines.length) {
        const tabLine = (lines[i] || '').trim();
        if (tabLine.startsWith(':::endtabs')) {
          if (currentTab) {
            tabs.push(currentTab);
          }
          i++;
          break;
        }
        if (tabLine.startsWith(':::tab')) {
          if (currentTab) {
            tabs.push(currentTab);
          }
          const tabParams = extractTagAndParams(tabLine).params;
          currentTab = { name: tabParams[0] || `Tab ${tabs.length + 1}`, content: '' };
        } else if (!tabLine.startsWith(':::endtab')) {
          if (currentTab) {
            currentTab.content += lines[i] + '\n';
          }
        }
        i++;
      }

      blocks.push({
        type: 'tabs',
        tabs,
        activeTab: params[0] || (tabs[0] && tabs[0].name) || '',
      });
      continue;
    }

    if (tag === 'photo') {
      const rows: string[][] = [];
      let currentRow: string[] = [];
      i++;

      while (i < lines.length) {
        const photoLine = (lines[i] || '').trim();
        if (photoLine === ':::endphoto') {
          if (currentRow.length > 0) {
            rows.push(currentRow);
          }
          i++;
          break;
        }
        if (photoLine === ':::n') {
          if (currentRow.length > 0) {
            rows.push(currentRow);
            currentRow = [];
          }
        } else {
          const images = photoLine.split(/\s+/).filter(img => img.trim());
          currentRow.push(...images);
        }
        i++;
      }

      blocks.push({ type: 'photo', rows });
      continue;
    }

    const endTag = `:::end${tag}`;
    const contentLines: string[] = [];
    i++;

    while (i < lines.length) {
      const contentLine = lines[i] || '';
      if (contentLine.trim() === endTag) {
        i++;
        break;
      }
      contentLines.push(contentLine);
      i++;
    }

    const content = contentLines.join('\n');

    if (tag === 'note') {
      blocks.push({
        type: 'note',
        noteType: params[0] || 'info',
        title: params[1] || '',
        content,
      });
    } else if (tag === 'fold') {
      blocks.push({
        type: 'fold',
        title: params[0] || '点击展开',
        open: params[1] === 'true' || params[1] === 'open',
        content,
      });
    } else {
      blocks.push({ type: 'markdown', content });
    }
  }

  flushMarkdown();

  return blocks;
}

