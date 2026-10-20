/**
 * Cloudflare Worker - 外站 API 反向代理
 * 用于代理 GitHub、Google 和 Microsoft 的 OAuth API 请求
 *
 * 部署地址: https://proxy.flec.top
 *
 * 路由映射:
 *   /github-api/*  -> https://api.github.com/*
 *   /github/*      -> https://github.com/*
 *   /google-oauth2/* -> https://oauth2.googleapis.com/*
 *   /google-api/*  -> https://www.googleapis.com/*
 *   /google/*      -> https://accounts.google.com/*
 *   /microsoft-graph/* -> https://graph.microsoft.com/*
 *   /microsoft/*   -> https://login.microsoftonline.com/*
 *   /telegram/*    -> https://api.telegram.org/*
 */
export default {
  async fetch(request) {
    const url = new URL(request.url);
    const path = url.pathname;

    if (request.method === 'OPTIONS') {
      return new Response(null, {
        headers: {
          'Access-Control-Allow-Origin': '*',
          'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
          'Access-Control-Allow-Headers': '*',
        }
      });
    }

    // 路由表（长前缀在前）
    const targets = [
      ['/github-api', 'https://api.github.com'],
      ['/github', 'https://github.com'],
      ['/google-oauth2', 'https://oauth2.googleapis.com'],
      ['/google-api', 'https://www.googleapis.com'],
      ['/google', 'https://accounts.google.com'],
      ['/microsoft-graph', 'https://graph.microsoft.com'],
      ['/microsoft', 'https://login.microsoftonline.com'],
      ['/telegram', 'https://api.telegram.org'],
    ];

    let targetBase = null, matchedPrefix = null;
    for (const [prefix, target] of targets) {
      if (path.startsWith(prefix)) {
        targetBase = target;
        matchedPrefix = prefix;
        break;
      }
    }

    if (!targetBase) {
      return new Response(JSON.stringify({ error: 'Invalid path' }), {
        status: 400,
        headers: { 'Content-Type': 'application/json' }
      });
    }

    // 构建目标 URL
    const targetPath = path.slice(matchedPrefix.length) || '/';
    const targetUrl = new URL(targetPath, targetBase);
    targetUrl.search = url.search;

    // 构建请求头
    const headers = new Headers();
    const skipHeaders = ['cf-connecting-ip', 'cf-ipcountry', 'cf-ray', 'cf-visitor', 'x-forwarded-proto', 'x-real-ip', 'host', 'connection'];
    for (const [key, value] of request.headers.entries()) {
      if (!skipHeaders.includes(key.toLowerCase()) && !key.toLowerCase().startsWith('cf-')) {
        headers.set(key, value);
      }
    }
    headers.set('Host', targetUrl.host);

    // 构建请求
    const fetchOptions = { method: request.method, headers, redirect: 'follow' };
    if (request.method !== 'GET' && request.method !== 'HEAD') {
      fetchOptions.body = await request.arrayBuffer();
    }

    try {
      const response = await fetch(targetUrl.toString(), fetchOptions);
      const responseHeaders = new Headers(response.headers);
      responseHeaders.set('Access-Control-Allow-Origin', '*');
      responseHeaders.delete('content-encoding');

      return new Response(response.body, {
        status: response.status,
        statusText: response.statusText,
        headers: responseHeaders,
      });
    } catch (error) {
      return new Response(JSON.stringify({ error: error.message }), {
        status: 502,
        headers: { 'Content-Type': 'application/json' }
      });
    }
  },
};
