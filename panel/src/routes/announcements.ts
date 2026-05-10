import { Hono } from 'hono';
import { Env, Announcement } from '../types';
import * as db from '../services/db';

const routes = new Hono<{ Bindings: Env }>();

routes.get('/', async (c) => {
  return c.json({ announcements: await db.getAnnouncements(c.env) });
});

routes.post('/', async (c) => {
  const body = await c.req.json<Announcement>();
  if (!body.title || !body.content) return c.json({ error: 'Missing fields' }, 400);
  await c.env.DB.prepare('INSERT INTO announcements (title, content, link) VALUES (?, ?, ?)')
    .bind(body.title, body.content, body.link || null).run();
  return c.json({ success: true });
});

routes.delete('/:id', async (c) => {
  await c.env.DB.prepare('DELETE FROM announcements WHERE id = ?').bind(c.req.param('id')).run();
  return c.json({ success: true });
});

export { routes as announcementRoutes };
