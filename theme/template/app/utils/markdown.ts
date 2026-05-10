/**
 * Markdown 渲染工具函数
 *
 * 基于 markdown-it 提供 Markdown 渲染功能，
 * 并通过 DOMPurify 进行 XSS 防护。
 *
 * @module utils/markdown
 */

import MarkdownIt from 'markdown-it';
import anchor from 'markdown-it-anchor';
import DOMPurify from 'isomorphic-dompurify';

const md = new MarkdownIt({
  html: false,
  breaks: true,
  linkify: true,
});

md.use(anchor, {
  slugify: (s: string) =>
    s
      .toLowerCase()
      .replace(/[^\u4e00-\u9fa5a-z0-9]+/g, '-')
      .replace(/^-+|-+$/g, ''),
  permalink: false,
  level: [1, 2, 3, 4, 5, 6],
});

/**
 * 渲染完整 Markdown 内容
 *
 * 将 Markdown 文本渲染为 HTML，支持：
 * - 标题锚点（用于目录跳转）
 * - 自动链接识别
 * - 换行转换
 * - XSS 防护（通过 DOMPurify）
 *
 * @param markdown - Markdown 文本
 * @returns 渲染后的安全 HTML 字符串
 *
 * @example
 * ```ts
 * const html = renderMarkdown('# 标题\n\n这是**粗体**文本');
 * // '<h1 id="标题">标题</h1><p>这是<strong>粗体</strong>文本</p>'
 * ```
 */
export function renderMarkdown(markdown: string): string {
  if (!markdown) return '';
  const rawHtml = md.render(markdown);
  return DOMPurify.sanitize(rawHtml, {
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
      'input',
      'iframe',
      'video',
      'audio',
      'source',
      'picture',
      'details',
      'summary',
      'figure',
      'figcaption',
      'aside',
      'section',
      'article',
      'header',
      'footer',
      'nav',
      'button',
      'label',
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
      'name',
      'value',
      'checked',
      'disabled',
      'controls',
      'autoplay',
      'loop',
      'muted',
      'poster',
      'preload',
      'allow',
      'allowfullscreen',
      'frameborder',
      'scrolling',
      'style',
      'role',
      'aria-label',
      'aria-expanded',
      'data-*',
    ],
    ALLOW_DATA_ATTR: true,
    ADD_ATTR: ['target', 'allowfullscreen'],
  });
}

/**
 * 渲染简单 Markdown 内容
 *
 * 将 Markdown 文本渲染为 HTML，仅支持基础格式：
 * - 段落、换行
 * - 粗体、斜体
 * - 代码块、行内代码
 * - 列表
 * - 引用
 * - 链接、图片
 *
 * 用于渲染用户输入的简单内容，如评论、动态等。
 *
 * @param markdown - Markdown 文本
 * @returns 渲染后的安全 HTML 字符串
 *
 * @example
 * ```ts
 * const html = renderSimpleMarkdown('这是**粗体**和`代码`');
 * // '<p>这是<strong>粗体</strong>和<code>代码</code></p>'
 * ```
 */
export function renderSimpleMarkdown(markdown: string): string {
  if (!markdown) return '';
  const simpleMd = new MarkdownIt({
    html: false,
    breaks: true,
    linkify: true,
  });
  const simpleHtml = simpleMd.render(markdown);
  return DOMPurify.sanitize(simpleHtml, {
    ALLOWED_TAGS: [
      'p',
      'br',
      'strong',
      'em',
      'code',
      'pre',
      'ul',
      'ol',
      'li',
      'blockquote',
      'a',
      'img',
    ],
    ALLOWED_ATTR: ['href', 'title', 'src', 'alt', 'width', 'height', 'class'],
  });
}
