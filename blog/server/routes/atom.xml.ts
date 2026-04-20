export default defineEventHandler(async event => {
  const config = useRuntimeConfig();
  const backendUrl = config.public.apiUrl.replace(/\/+$/, '').replace(/\/api\/v\d+$/, '');

  try {
    const response = await fetch(`${backendUrl}/atom.xml`);

    if (!response.ok) {
      throw new Error(`Backend returned ${response.status}`);
    }

    const atomFeed = await response.text();
    setResponseHeader(event, 'Content-Type', 'application/atom+xml; charset=utf-8');
    return atomFeed;
  } catch {
    throw createError({
      statusCode: 500,
      statusMessage: 'Failed to fetch Atom feed',
    });
  }
});
