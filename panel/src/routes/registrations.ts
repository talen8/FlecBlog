import { Hono } from 'hono';
import { Env, RegistrationUpdate } from '../types';

// ── Public API ──────────────────────────────────────────────────────────
const api = new Hono<{ Bindings: Env }>();

api.post('/', async (c) => {
  const body = await c.req.json<{
    client_key?: string;
    api_key?: string;
    site_url: string;
    version?: string;
  }>();

  const { client_key, api_key, site_url, version } = body;

  if (!site_url) {
    return c.json({ error: 'site_url is required' }, 400);
  }

  if (site_url.includes('localhost')) {
    return c.json({ error: '请确保后端服务外网可访问' }, 400);
  }

  // Heartbeat: existing client_key
  if (client_key) {
    const existing = await c.env.DB.prepare(
      'SELECT * FROM registrations WHERE client_key = ?'
    ).bind(client_key).first<any>();

    if (!existing) {
      return c.json({ error: 'client_key not found' }, 401);
    }

    // 校验 api_key，防止自编译用户用旧 client_key 白嫖
    const panelApiKey = c.env.PANEL_API_KEY;
    if (panelApiKey) {
      if (!api_key || api_key !== panelApiKey) {
        return c.json({ error: 'Unauthorized' }, 401);
      }
    }

    const siteUrlChanged = site_url !== existing.site_url;
    const history = JSON.parse(existing.site_url_history || '[]');

    if (siteUrlChanged) {
      const now = new Date().toISOString().replace('T', ' ').slice(0, 19);
      if (existing.site_url) {
        const existingEntry = history.find((h: any) => h.url === existing.site_url);
        if (existingEntry) {
          existingEntry.last_seen = now;
        } else {
          history.push({ url: existing.site_url, first_seen: existing.last_seen_at || now, last_seen: now });
        }
      }
      const newEntry = history.find((h: any) => h.url === site_url);
      if (newEntry) {
        newEntry.first_seen = now;
      } else {
        history.push({ url: site_url, first_seen: now, last_seen: now });
      }
    }

    const updateFields: string[] = ["last_seen_at = datetime('now')"];
    const params: any[] = [];

    if (siteUrlChanged) {
      updateFields.push('site_url = ?');
      params.push(site_url);
      updateFields.push('site_url_history = ?');
      params.push(JSON.stringify(history));
    }
    if (version !== undefined && version !== existing.version) {
      updateFields.push('version = ?');
      params.push(version);
    }

    params.push(client_key);
    await c.env.DB.prepare(
      'UPDATE registrations SET ' + updateFields.join(', ') + ' WHERE client_key = ?'
    ).bind(...params).run();

    return c.json({ status: 'ok' });
  }

  // First-time registration
  if (api_key) {
    const panelApiKey = c.env.PANEL_API_KEY;
    if (panelApiKey && api_key !== panelApiKey) {
      return c.json({ error: 'Unauthorized' }, 401);
    }
  } else {
    return c.json({ error: 'api_key or client_key is required' }, 400);
  }

  // Check if site_url already exists
  const existing = await c.env.DB.prepare(
    'SELECT * FROM registrations WHERE site_url = ?'
  ).bind(site_url).first<any>();

  if (existing) {
    await c.env.DB.prepare(
      "UPDATE registrations SET version = ?, last_seen_at = datetime('now') WHERE id = ?"
    ).bind(version || '', existing.id).run();
    return c.json({ client_key: existing.client_key });
  }

  // Check site_url_history
  const historyRows = await c.env.DB.prepare(
    'SELECT id, client_key, site_url_history, last_seen_at FROM registrations WHERE site_url_history LIKE ?'
  ).bind('%"url":"' + site_url + '"%').all<any>();

  if (historyRows.results && historyRows.results.length > 0) {
    const match = historyRows.results[0];
    const now = new Date().toISOString().replace('T', ' ').slice(0, 19);
    const history = JSON.parse(match.site_url_history || '[]');

    if (match.site_url) {
      const oldEntry = history.find((h: any) => h.url === match.site_url);
      if (oldEntry) {
        oldEntry.last_seen = now;
      } else {
        history.push({ url: match.site_url, first_seen: match.last_seen_at || now, last_seen: now });
      }
    }
    const newEntry = history.find((h: any) => h.url === site_url);
    if (newEntry) {
      newEntry.first_seen = now;
    } else {
      history.push({ url: site_url, first_seen: now, last_seen: now });
    }

    await c.env.DB.prepare(
      'UPDATE registrations SET site_url = ?, site_url_history = ?, version = ?, last_seen_at = datetime(\'now\') WHERE id = ?'
    ).bind(site_url, JSON.stringify(history), version || '', match.id).run();

    return c.json({ client_key: match.client_key });
  }

  // Fresh registration
  const clientKey = 'panel_' + crypto.randomUUID().replace(/-/g, '').slice(0, 24);
  await c.env.DB.prepare(
    'INSERT INTO registrations (client_key, site_url, version) VALUES (?, ?, ?)'
  ).bind(clientKey, site_url, version || '').run();

  return c.json({ client_key: clientKey });
});

// ── Admin API ───────────────────────────────────────────────────────────
const routes = new Hono<{ Bindings: Env }>();

routes.get('/', async (c) => {
  const rows = await c.env.DB.prepare(
    'SELECT * FROM registrations ORDER BY last_seen_at DESC'
  ).all<any>();
  return c.json(rows.results || []);
});

routes.get('/:id', async (c) => {
  const id = c.req.param('id');
  const row = await c.env.DB.prepare(
    'SELECT * FROM registrations WHERE id = ?'
  ).bind(id).first<any>();
  if (!row) return c.json({ error: 'Not found' }, 404);
  return c.json(row);
});

routes.put('/:id', async (c) => {
  const id = c.req.param('id');
  const body = await c.req.json<RegistrationUpdate>();

  const updateFields: string[] = [];
  const params: any[] = [];

  if (body.banned !== undefined) {
    updateFields.push('banned = ?');
    params.push(body.banned);
  }
  if (body.remark !== undefined) {
    updateFields.push('remark = ?');
    params.push(body.remark);
  }

  if (updateFields.length === 0) {
    return c.json({ error: 'No fields to update' }, 400);
  }

  params.push(id);
  await c.env.DB.prepare(
    'UPDATE registrations SET ' + updateFields.join(', ') + ' WHERE id = ?'
  ).bind(...params).run();

  return c.json({ status: 'ok' });
});

export { api as registerApi, routes as registrationRoutes };
