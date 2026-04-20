// 动态生成 sitemap 路由

interface SitemapArticle {
  url: string;
  update_time?: string;
  publish_time?: string;
}

interface SitemapCategory {
  url: string;
}

interface SitemapTag {
  url: string;
}

interface ArticlesResponse {
  data?: {
    list?: SitemapArticle[];
  };
}

interface CategoriesResponse {
  data?: {
    list?: SitemapCategory[];
  };
}

interface TagsResponse {
  data?: {
    list?: SitemapTag[];
  };
}

interface SitemapUrl {
  loc: string;
  lastmod?: string;
  changefreq: 'weekly';
  priority: number;
}

export default defineNitroPlugin(async nitroApp => {
  nitroApp.hooks.hook('sitemap:resolved', async ctx => {
    try {
      const config = useRuntimeConfig();
      const apiUrl = config.public.apiUrl;

      // 从后端 API 获取文章列表
      const articlesRes = await $fetch<ArticlesResponse>(`${apiUrl}/articles`).catch(() => null);

      // 添加文章路由到 sitemap
      const articles = articlesRes?.data?.list || [];
      articles.forEach((article: SitemapArticle) => {
        ctx.urls.push({
          loc: article.url,
          lastmod: article.update_time || article.publish_time,
          changefreq: 'weekly' as const,
          priority: 0.8,
        } as SitemapUrl);
      });

      // 添加分类路由
      const categoriesRes = await $fetch<CategoriesResponse>(`${apiUrl}/categories`).catch(
        () => null
      );
      const categories = categoriesRes?.data?.list || [];
      categories.forEach((category: SitemapCategory) => {
        ctx.urls.push({
          loc: category.url,
          changefreq: 'weekly' as const,
          priority: 0.6,
        } as SitemapUrl);
      });

      // 添加标签路由
      const tagsRes = await $fetch<TagsResponse>(`${apiUrl}/tags`).catch(() => null);
      const tags = tagsRes?.data?.list || [];
      tags.forEach((tag: SitemapTag) => {
        ctx.urls.push({
          loc: tag.url,
          changefreq: 'weekly' as const,
          priority: 0.5,
        } as SitemapUrl);
      });
    } catch {
      // sitemap 生成失败不影响应用运行，静默忽略
    }
  });
});
