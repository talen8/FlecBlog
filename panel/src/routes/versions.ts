import { Hono } from 'hono';
import { Env } from '../types';
import * as db from '../services/db';
import { syncGitHubReleases } from '../services/github';

const routes = new Hono<{ Bindings: Env }>();

routes.get('/', async (c) => {
  const versions = await db.getAllVersions(c.env);
  return c.json({ versions });
});

routes.post('/sync', async (c) => {
  const result = await syncGitHubReleases(c.env);
  if (result.success) return c.json({ success: true, count: result.count });
  return c.json({ error: result.error }, 500);
});

routes.put('/:version{.+}', async (c) => {
  const version = c.req.param('version');
  const body = await c.req.json();
  const updates: string[] = [];
  const values: (string | number)[] = [];

  if (body.date !== undefined) { updates.push('date = ?'); values.push(body.date); }
  if (body.changes !== undefined) { updates.push('changes = ?'); values.push(body.changes); }
  if (body.enabled !== undefined) { updates.push('enabled = ?'); values.push(body.enabled ? 1 : 0); }

  if (updates.length === 0) return c.json({ error: 'No fields' }, 400);
  values.push(version);
  await c.env.DB.prepare(`UPDATE versions SET ${updates.join(', ')} WHERE version = ?`).bind(...values).run();
  return c.json({ success: true });
});

routes.delete('/:version{.+}', async (c) => {
  await c.env.DB.prepare('DELETE FROM versions WHERE version = ?').bind(c.req.param('version')).run();
  return c.json({ success: true });
});

export { routes as versionRoutes };
