import { Hono } from 'hono';
import { Env, AIProvider } from '../types';

// --- 管理接口 ---
const aiRoutes = new Hono<{ Bindings: Env }>();

aiRoutes.get('/providers', async (c) => {
  const result = await c.env.DB.prepare(
    'SELECT * FROM ai_providers ORDER BY priority_immediate ASC, priority_deferred ASC'
  ).all<AIProvider>();
  return c.json({ providers: result.results });
});

aiRoutes.post('/providers', async (c) => {
  const body = await c.req.json();
  await c.env.DB.prepare(
    `INSERT INTO ai_providers (name, base_url, api_key, model, priority_immediate, priority_deferred)
     VALUES (?, ?, ?, ?, ?, ?)`
  ).bind(body.name, body.base_url, body.api_key, body.model,
    body.priority_immediate ?? 999, body.priority_deferred ?? 999).run();
  return c.json({ success: true });
});

aiRoutes.put('/providers/:id', async (c) => {
  const id = c.req.param('id');
  const body = await c.req.json();
  const updates: string[] = [];
  const values: (string | number)[] = [];

  if (body.name !== undefined) { updates.push('name = ?'); values.push(body.name); }
  if (body.base_url !== undefined) { updates.push('base_url = ?'); values.push(body.base_url); }
  if (body.api_key !== undefined) { updates.push('api_key = ?'); values.push(body.api_key); }
  if (body.model !== undefined) { updates.push('model = ?'); values.push(body.model); }
  if (body.priority_immediate !== undefined) { updates.push('priority_immediate = ?'); values.push(body.priority_immediate); }
  if (body.priority_deferred !== undefined) { updates.push('priority_deferred = ?'); values.push(body.priority_deferred); }

  if (updates.length === 0) return c.json({ error: 'No fields' }, 400);
  updates.push("updated_at = datetime('now')");
  values.push(id);
  await c.env.DB.prepare(
    `UPDATE ai_providers SET ${updates.join(', ')} WHERE id = ?`
  ).bind(...values).run();
  return c.json({ success: true });
});

aiRoutes.delete('/providers/:id', async (c) => {
  const id = c.req.param('id');
  await c.env.DB.prepare('DELETE FROM ai_providers WHERE id = ?').bind(id).run();
  return c.json({ success: true });
});

// --- 代理接口 ---
const aiApi = new Hono<{ Bindings: Env }>();
aiApi.post('/', async (c) => {
  const mode = c.req.header('X-AI-Mode') || 'immediate';

  const orderBy = mode === 'deferred'
    ? 'ORDER BY priority_deferred ASC, priority_immediate ASC'
    : 'ORDER BY priority_immediate ASC, priority_deferred ASC';

  const result = await c.env.DB.prepare(
    `SELECT * FROM ai_providers ${orderBy}`
  ).all<AIProvider>();

  if (result.results.length === 0) {
    return c.json({ error: '未配置 AI 提供商' }, 503);
  }

  const body = await c.req.json();
  const messages = body.messages;
  if (!messages || !Array.isArray(messages)) {
    return c.json({ error: '缺少 messages 字段' }, 400);
  }

  let lastError: string | null = null;

  for (const provider of result.results) {
    const startTime = Date.now();
    try {
      const upstream = await forwardToProvider(provider, messages, body);
      const latency = Date.now() - startTime;

      c.executionCtx.waitUntil(updateStats(c.env, provider.id, latency, upstream.usage));

      return new Response(JSON.stringify(upstream.body), {
        status: upstream.status,
        headers: { 'Content-Type': 'application/json' },
      });
    } catch (err) {
      const latency = Date.now() - startTime;
      lastError = err instanceof Error ? err.message : String(err);
      c.executionCtx.waitUntil(recordFailure(c.env, provider.id, lastError));
    }
  }

  return c.json({ error: lastError || '所有 AI 提供商均失败' }, 502);
});

async function forwardToProvider(
  provider: AIProvider,
  messages: unknown[],
  originalBody: any,
): Promise<{ body: any; status: number; usage: any }> {
  const reqBody: Record<string, any> = {
    model: provider.model,
    messages,
  };

  if (originalBody.stream) {
    reqBody.stream = true;
  }
  if (originalBody.temperature !== undefined) {
    reqBody.temperature = originalBody.temperature;
  }
  if (originalBody.max_tokens !== undefined) {
    reqBody.max_tokens = originalBody.max_tokens;
  }

  const resp = await fetch(provider.base_url.replace(/\/$/, '') + '/chat/completions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${provider.api_key}`,
    },
    body: JSON.stringify(reqBody),
    signal: AbortSignal.timeout(60000),
  });

  const respBody = await resp.json<any>();

  if (resp.status >= 500) {
    throw new Error(`上游返回 ${resp.status}: ${JSON.stringify(respBody)}`);
  }

  if (!resp.ok) {
    return { body: respBody, status: resp.status, usage: null };
  }

  return { body: respBody, status: resp.status, usage: respBody.usage || null };
}

async function updateStats(env: Env, providerId: number, latency: number, usage: any) {
  const updates = [
    'call_count = call_count + 1',
    'total_latency_ms = total_latency_ms + ?',
    "last_call_at = datetime('now')",
    "updated_at = datetime('now')",
  ];
  const values: (string | number)[] = [latency];

  if (usage) {
    updates.push('total_prompt_tokens = total_prompt_tokens + ?');
    values.push(usage.prompt_tokens || 0);
    updates.push('total_completion_tokens = total_completion_tokens + ?');
    values.push(usage.completion_tokens || 0);
  }

  values.push(providerId);
  await env.DB.prepare(
    `UPDATE ai_providers SET ${updates.join(', ')} WHERE id = ?`
  ).bind(...values).run();
}

async function recordFailure(env: Env, providerId: number, error: string) {
  await env.DB.prepare(
    `UPDATE ai_providers SET
      fail_count = fail_count + 1,
      last_error_message = ?,
      updated_at = datetime('now')
     WHERE id = ?`
  ).bind(error.slice(0, 500), providerId).run();
}

export { aiRoutes, aiApi };
