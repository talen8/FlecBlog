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

    const atomEntries = articles
      // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 文章数据类型，结构由后端确定
      .map((article: any) => {
        const pubDate = new Date(article.publish_time).toISOString();
        return `
  <entry>
    <title><![CDATA[${article.title}]]></title>
    <link href="${article.url}"/>
    <id>${article.url}</id>
    <updated>${pubDate}</updated>
    <summary><![CDATA[${article.summary}]]></summary>
  </entry>`;
      })
      .join('\n');

    const atom = `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>Blog Atom Feed</title>
  <link href="${baseURL}"/>
  <updated>${new Date().toISOString()}</updated>
  <id>${baseURL}</id>
  ${atomEntries}
</feed>`;

    setHeader(event, 'Content-Type', 'application/atom+xml; charset=utf-8');
    return atom;
  } catch {
    throw createError({
      statusCode: 500,
      message: 'Failed to generate Atom feed',
    });
  }
});
