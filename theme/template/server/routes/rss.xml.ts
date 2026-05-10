export default defineEventHandler(async event => {
  const config = useRuntimeConfig();
  const baseURL = config.public.apiUrl;

  try {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- API 响应类型，结构由后端确定
    const response: any = await $fetch(`${baseURL}/articles`, {
      params: {
        page: 1,
        page_size: 100,
      },
    });

    const articles = response.data?.list || [];

    const rssItems = articles
      // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 文章数据类型，结构由后端确定
      .map((article: any) => {
        const pubDate = new Date(article.publish_time).toUTCString();
        return `
    <item>
      <title><![CDATA[${article.title}]]></title>
      <link>${article.url}</link>
      <description><![CDATA[${article.summary}]]></description>
      <pubDate>${pubDate}</pubDate>
      <guid>${article.url}</guid>
    </item>`;
      })
      .join('\n');

    const rss = `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
  <channel>
    <title>Blog RSS Feed</title>
    <link>${baseURL}</link>
    <description>Blog RSS Feed</description>
    <language>zh-CN</language>
    ${rssItems}
  </channel>
</rss>`;

    setHeader(event, 'Content-Type', 'application/xml; charset=utf-8');
    return rss;
  } catch {
    throw createError({
      statusCode: 500,
      message: 'Failed to generate RSS feed',
    });
  }
});
