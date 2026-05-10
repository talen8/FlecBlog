import { Hono } from 'hono';
import { cors } from 'hono/cors';
import { logger } from 'hono/logger';
import { Env } from './types';

import { publicRoutes } from './routes/public';
import { versionRoutes } from './routes/versions';
import { announcementRoutes } from './routes/announcements';
import { settingsRoutes } from './routes/settings';
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

app.route('/api', publicRoutes);
app.route('/admin/versions', versionRoutes);
app.route('/admin/announcements', announcementRoutes);
app.route('/admin', settingsRoutes);

app.get('/admin', async (c) => c.env.ASSETS.fetch(new URL('admin.html', c.req.url)));

export default {
  fetch: app.fetch,
  async scheduled(event: ScheduledEvent, env: Env) {
    await syncGitHubReleases(env);
    await autoEnablePendingVersions(env);
  },
};
