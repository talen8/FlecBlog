import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { logger } from 'hono/logger';
import { Env } from './types';

import { versionRoutes, versionsApi } from './routes/versions';
import { announcementRoutes, announcementsApi } from './routes/announcements';
import { settingsRoutes } from './routes/settings';
import { aiRoutes, aiApi } from './routes/ai';
import { registerApi, registrationRoutes } from './routes/registrations';
import { syncGitHubReleases, autoEnablePendingVersions } from './services/github';

type AppEnv = Env & { ASSETS: Fetcher };
const app = new Hono<{ Bindings: AppEnv }>();

app.use(
  '*',
  cors({
    origin: (origin) => origin,
    credentials: true,
  })
);
app.use('*', logger());

// 校验注册密钥，保护 /api 路由
app.use('/api/*', async (c, next) => {
  if (c.req.path === '/api/announcements' || c.req.path === '/api/register') {
    return next();
  }
  const clientKey = c.req.header('X-Client-Key');
  if (!clientKey) {
    return c.json({ error: 'Unauthorized' }, 401);
  }
  const reg = await c.env.DB.prepare(
    'SELECT * FROM registrations WHERE client_key = ?'
  ).bind(clientKey).first<any>();
  if (!reg) {
    return c.json({ error: '该功能需要官方版本' }, 401);
  }
  if (reg.banned === 1) {
    return c.json({ error: '该部署已被封禁' }, 403);
  }

  // Version mismatch detection
  const reqVersion = c.req.header('X-Version');
  if (reqVersion && reg.version && reqVersion !== reg.version) {
    return c.json({ error: '疑似非官方版本，请重启容器后重试' }, 409);
  }

  // Site URL tracking + call counters + last_seen_at（单条 UPDATE）
  const updateFields: string[] = ["last_seen_at = datetime('now')"];
  const updateParams: any[] = [];

  const reqSiteUrl = c.req.header('X-Site-Url');
  if (reqSiteUrl && reqSiteUrl !== reg.site_url) {
    const now = new Date().toISOString().replace('T', ' ').slice(0, 19);
    const history = JSON.parse(reg.site_url_history || '[]');
    if (reg.site_url) {
      const entry = history.find((h: any) => h.url === reg.site_url);
      if (entry) entry.last_seen = now;
      else history.push({ url: reg.site_url, first_seen: reg.last_seen_at || now, last_seen: now });
    }
    const newEntry = history.find((h: any) => h.url === reqSiteUrl);
    if (newEntry) newEntry.first_seen = now;
    else history.push({ url: reqSiteUrl, first_seen: now, last_seen: now });
    updateFields.push('site_url = ?', 'site_url_history = ?');
    updateParams.push(reqSiteUrl, JSON.stringify(history));
  }

  const today = new Date().toISOString().slice(0, 10);
  const counts = JSON.parse(reg.call_counts || '{}');
  if (!counts[today]) counts[today] = {};
  const path = c.req.path;
  if (path === '/api/ai') counts[today].ai = (counts[today].ai || 0) + 1;
  else if (path === '/api/versions') counts[today].versions = (counts[today].versions || 0) + 1;
  updateFields.push('call_counts = ?');
  updateParams.push(JSON.stringify(counts));

  updateParams.push(reg.id);
  await c.env.DB.prepare(
    'UPDATE registrations SET ' + updateFields.join(', ') + ' WHERE id = ?'
  ).bind(...updateParams).run();

  await next();
});

app.route('/api/versions', versionsApi);
app.route('/api/announcements', announcementsApi);
app.route('/api/ai', aiApi);
app.route('/api/register', registerApi);
app.route('/admin/versions', versionRoutes);
app.route('/admin/announcements', announcementRoutes);
app.route('/admin/ai', aiRoutes);
app.route('/admin/settings', settingsRoutes);
app.route('/admin/registrations', registrationRoutes);

app.get('/admin', async (c) => c.env.ASSETS.fetch(new URL('admin.html', c.req.url)));

export default {
  fetch: app.fetch,
  async scheduled(event: ScheduledEvent, env: Env) {
    await syncGitHubReleases(env);
    await autoEnablePendingVersions(env);
  },
};
