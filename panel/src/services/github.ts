import { Env } from '../types';
import { getSetting } from './db';

interface GitHubApiRelease {
  tag_name: string;
  published_at: string;
  body: string;
}

export async function syncGitHubReleases(env: Env): Promise<{ success: boolean; count: number; error?: string }> {
  const githubRepo = await getSetting(env, 'github_repo');
  if (!githubRepo) {
    return { success: false, count: 0, error: '请先配置 GitHub 仓库' };
  }

  const githubToken = await getSetting(env, 'github_token');

  try {
    const headers: Record<string, string> = {
      'Accept': 'application/vnd.github.v3+json',
      'User-Agent': 'FlecPanel/1.0',
    };

    if (githubToken) {
      headers['Authorization'] = `Bearer ${githubToken}`;
    }

    const res = await fetch(`https://api.github.com/repos/${githubRepo}/releases?per_page=10`, {
      headers,
    });

    if (!res.ok) {
      if (res.status === 403) {
        return { success: false, count: 0, error: 'API 速率限制，请配置 GitHub Token' };
      }
      return { success: false, count: 0, error: `GitHub API 错误: ${res.status}` };
    }

    const releases = await res.json() as GitHubApiRelease[];
    let count = 0;

    for (const release of releases) {
      let version = release.tag_name;
      if (!version.startsWith('v')) continue;

      version = version.substring(1);
      const date = release.published_at.split('T')[0];
      const changes = parseReleaseBody(release.body || '');

      const result = await env.DB.prepare(
        'INSERT OR IGNORE INTO versions (version, date, changes, enabled) VALUES (?, ?, ?, 0)'
      )
        .bind(version, date, changes)
        .run();

      if (result.meta?.changes && result.meta.changes > 0) {
        count++;
      }
    }

    return { success: true, count };
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unknown error';
    return { success: false, count: 0, error: message };
  }
}

function parseReleaseBody(body: string): string {
  const lines: string[] = [];

  const patterns: Record<string, RegExp> = {
    feat: /### ?(?:✨|New Features?|Features?)\n([\s\S]*?)(?=\n###|$)/i,
    fix: /### ?(?:🐛|Bug Fixes?|Fixes?)\n([\s\S]*?)(?=\n###|$)/i,
    refactor: /### ?(?:♻️|Refactors?|Code Changes?)\n([\s\S]*?)(?=\n###|$)/i,
    docs: /### ?(?:📝|Documentation|Docs?)\n([\s\S]*?)(?=\n###|$)/i,
    style: /### ?(?:💄|Styles?)\n([\s\S]*?)(?=\n###|$)/i,
    test: /### ?(?:✅|Tests?)\n([\s\S]*?)(?=\n###|$)/i,
    chore: /### ?(?:🔧|Chores?|Maintenance)\n([\s\S]*?)(?=\n###|$)/i,
  };

  const labels: Record<string, string> = {
    feat: '新功能',
    fix: '修复',
    refactor: '重构',
    docs: '文档',
    style: '样式',
    test: '测试',
    chore: '杂项',
  };

  for (const [key, pattern] of Object.entries(patterns)) {
    const match = body.match(pattern);
    if (match) {
      const items = match[1]
        .split('\n')
        .map(line => line.replace(/^[-*]\s*/, '').trim())
        .filter(line => line.length > 0);
      for (const item of items) {
        lines.push(`[${labels[key]}] ${item}`);
      }
    }
  }

  return lines.join('\n');
}
