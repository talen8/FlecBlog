/**
 * 文章详情页
 * 使用分段渲染支持自定义块
 */
import type { ArticleDetail } from '../../../types';
import type { IAppOption } from '../../../app';
import { getArticleBySlug } from '../../../api/article';
import { formatRelativeTime } from '../../../utils/format';
import { parseMarkdownToBlocks, type ContentBlock } from '../../../utils/markdown';

/**
 * 渲染后的内容块
 */
interface RenderedBlock {
  type: string;
  [key: string]: unknown;
}

Page({
  data: {
    article: {} as ArticleDetail,
    loading: true,
    error: '',
    slug: '',
    formattedTime: '',
    contentBlocks: [] as RenderedBlock[],
  },

  /**
   * 页面加载时触发
   * @param options - 页面参数
   */
  onLoad(options: { slug?: string }) {
    const slug = options.slug || '';
    if (!slug) {
      this.setData({ loading: false, error: '参数错误' });
      return;
    }
    this.setData({ slug });
    this.loadArticle(slug);
  },

  /**
   * 加载文章详情
   * @param slug - 文章 slug
   */
  async loadArticle(slug: string) {
    this.setData({ loading: true, error: '' });
    try {
      const article = await getArticleBySlug(slug);
      wx.setNavigationBarTitle({ title: article.title });

      const app = getApp<IAppOption>();
      const blocks = parseMarkdownToBlocks(article.content);
      const contentBlocks = this.renderBlocks(blocks, app);

      this.setData({
        article,
        loading: false,
        formattedTime: formatRelativeTime(article.publish_time),
        contentBlocks,
      });
    } catch {
      this.setData({ loading: false, error: '加载失败' });
    }
  },

  /**
 * 渲染内容块
 * @param blocks - 解析后的内容块
 * @param app - 应用实例
 * @returns 渲染后的内容块
 */
  renderBlocks(blocks: ContentBlock[], app: IAppOption): RenderedBlock[] {
    return blocks.map(block => {
      if (block.type === 'markdown') {
        return {
          type: 'markdown',
          content: app.towxml(block.content, 'markdown', {
            theme: 'light',
            events: {
              tap: (e: { target: { dataset: { href?: string } } }) => {
                const href = e.target?.dataset?.href;
                if (href) {
                  if (href.startsWith('http://') || href.startsWith('https://')) {
                    wx.setClipboardData({
                      data: href,
                      success: () => {
                        wx.showToast({ title: '链接已复制', icon: 'success' });
                      },
                    });
                  } else {
                    wx.navigateTo({ url: href });
                  }
                }
              },
            },
          }),
        };
      }

      if (block.type === 'note') {
        const innerBlocks = parseMarkdownToBlocks(block.content);
        const contentBlocks = this.renderBlocks(innerBlocks, app);
        return {
          type: 'note',
          noteType: block.noteType,
          title: block.title,
          contentBlocks,
        };
      }

      if (block.type === 'fold') {
        const innerBlocks = parseMarkdownToBlocks(block.content);
        const contentBlocks = this.renderBlocks(innerBlocks, app);
        return {
          type: 'fold',
          title: block.title,
          open: block.open,
          contentBlocks,
        };
      }

      if (block.type === 'tabs') {
        return {
          type: 'tabs',
          activeTab: block.activeTab,
          tabs: block.tabs.map(tab => {
            const innerBlocks = parseMarkdownToBlocks(tab.content);
            const contentBlocks = this.renderBlocks(innerBlocks, app);
            return {
              name: tab.name,
              contentBlocks,
            };
          }),
        };
      }

      return block as unknown as RenderedBlock;
    });
  },

  onRetry() {
    if (this.data.slug) this.loadArticle(this.data.slug);
  },

  onCopyOriginal() {
    const { url } = this.data.article;
    if (!url) return;
    const app = getApp<IAppOption>();
    const blogUrl = app.globalData.siteConfig.blog_url || '';
    const fullUrl = blogUrl ? `${blogUrl.replace(/\/$/, '')}${url}` : url;
    wx.setClipboardData({
      data: fullUrl,
      success: () => {
        wx.showToast({ title: '链接已复制', icon: 'success' });
      },
    });
  },

  onShareAppMessage() {
    const { title } = this.data.article;
    return {
      title: title || '文章分享',
      path: `/pages/article/detail/index?slug=${this.data.slug}`,
      imageUrl: this.data.article.cover,
    };
  },

  onShareTimeline() {
    const { title } = this.data.article;
    return {
      title: title || '文章分享',
      imageUrl: this.data.article.cover,
    };
  },
});
