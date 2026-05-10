import { Hono } from 'hono';
import { Env } from '../types';
import * as db from '../services/db';

const routes = new Hono<{ Bindings: Env }>();

routes.get('/versions', async (c) => {
  const versions = await db.getEnabledVersions(c.env, 10);
  return c.json(versions);
});

routes.get('/announcements', async (c) => {
  const announcements = await db.getPublicAnnouncements(c.env);
  return c.json(announcements);
});

export { routes as publicRoutes };
