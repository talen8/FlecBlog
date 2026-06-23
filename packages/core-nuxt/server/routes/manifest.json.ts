import { getSettingGroup } from '../../app/composables/api/sysconfig';

export default defineEventHandler(async event => {
  try {
    const basic = await getSettingGroup('basic');

    const manifest = {
      name: basic.title || 'FlecBlog',
      short_name: basic.title?.substring(0, 12) || 'Flec',
      description: basic.description || 'Flec个人博客',
      theme_color: '#f7f7f7',
      background_color: '#ffffff',
      display: 'standalone',
      start_url: '/',
      icons: [
        {
          src: basic.favicon || '/favicon.ico',
          sizes: '192x192',
          type: 'image/png',
        },
        {
          src: basic.favicon || '/favicon.ico',
          sizes: '512x512',
          type: 'image/png',
        },
      ],
    };

    setHeader(event, 'Content-Type', 'application/manifest+json');
    return manifest;
  } catch {
    return {
      name: 'FlecBlog',
      short_name: 'Flec',
      theme_color: '#f7f7f7',
      background_color: '#ffffff',
      display: 'standalone',
      start_url: '/',
      icons: [{ src: '/favicon.ico', sizes: '192x192', type: 'image/png' }],
    };
  }
});
