import { Hono } from 'hono';
import { Env } from '../types';
import { getSetting } from '../services/db';

// --- 代理接口 ---
// 托管存储服务代理，地址与 Token 从系统设置读取（hosted_storage_url / hosted_storage_token）
const storageApi = new Hono<{ Bindings: Env }>();

async function getHostedStorageConfig(env: Env) {
  const baseUrl = await getSetting(env, 'hosted_storage_url');
  const token = await getSetting(env, 'hosted_storage_token');
  if (!baseUrl || !token) {
    return null;
  }
  return { baseUrl: baseUrl.replace(/\/$/, ''), token };
}

storageApi.post('/upload', async (c) => {
  const config = await getHostedStorageConfig(c.env);
  if (!config) {
    return c.json({ error: '未配置托管存储服务（hosted_storage_url / hosted_storage_token）' }, 503);
  }

  const formData = await c.req.formData();
  const file = formData.get('file');
  if (typeof file !== 'object' || file === null) {
    return c.json({ error: '缺少文件字段 file' }, 400);
  }

  try {
    const url = await uploadToHostedStorage(config, file);
    return c.json({ url });
  } catch (err) {
    const message = err instanceof Error ? err.message : String(err);
    return c.json({ error: message }, 502);
  }
});

// --- 调用托管存储服务 ---

// 上传到托管存储服务
async function uploadToHostedStorage(config: { baseUrl: string; token: string }, file: File): Promise<string> {
  const formData = new FormData();
  formData.append('file', file);

  const resp = await fetch(`${config.baseUrl}/upload?returnFormat=full`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${config.token}` },
    body: formData,
  });

  const result = await resp.json<any>();
  if (!Array.isArray(result) || result.length === 0 || !result[0].src) {
    throw new Error(resp.ok ? '上传失败：响应格式异常' : `托管存储服务返回 ${resp.status}`);
  }

  const first = result[0];
  return first.publicUrl || first.src;
}

export { storageApi };
