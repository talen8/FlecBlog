import { Hono } from 'hono';
import { Env } from '../types';
import { getSetting } from '../services/db';

const routes = new Hono<{ Bindings: Env }>();

routes.get('/', async (c) => {
  const [repo, token] = await Promise.all([
    getSetting(c.env, 'github_repo'),
    getSetting(c.env, 'github_token'),
  ]);
  return c.json({
    settings: {
      github_repo: repo || '',
      github_token: token || '',
    },
  });
});

routes.put('/', async (c) => {
  const body = await c.req.json();
  const settings = [
    { key: 'github_repo', value: body.github_repo },
    { key: 'github_token', value: body.github_token },
  ];
  for (const { key, value } of settings) {
    if (value !== undefined) {
      await c.env.DB.prepare('INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)')
        .bind(key, value || '').run();
    }
  }
  return c.json({ success: true });
});

export { routes as settingsRoutes };
