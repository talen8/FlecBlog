/**
 * 数据库服务
 * 提供所有数据库操作方法
 */

import { Env, Version, Announcement, Setting } from '../types';

/**
 * 获取设置值
 */
export async function getSetting(env: Env, key: string): Promise<string | null> {
  const result = await env.DB.prepare('SELECT value FROM settings WHERE key = ?')
    .bind(key)
    .first<Setting>();

  return result?.value || null;
}

/**
 * 获取启用的版本列表
 */
export async function getEnabledVersions(env: Env, limit: number): Promise<Version[]> {
  const result = await env.DB.prepare(
    'SELECT id, version, date, changes FROM versions WHERE enabled = 1 ORDER BY date DESC LIMIT ?'
  )
    .bind(limit)
    .all<Version>();

  return result.results;
}

/**
 * 获取所有版本列表
 */
export async function getAllVersions(env: Env): Promise<Version[]> {
  const result = await env.DB.prepare(
    'SELECT * FROM versions ORDER BY date DESC'
  )
    .all<Version>();

  return result.results;
}

/**
 * 获取公告列表（公开接口，最近1天）
 */
export async function getPublicAnnouncements(env: Env): Promise<Announcement[]> {
  const result = await env.DB.prepare(
    `SELECT id, title, content, link FROM announcements
     WHERE created_at >= datetime('now', '-1 day')
     ORDER BY created_at DESC`
  ).all<Announcement>();

  return result.results;
}

/**
 * 获取所有公告列表（管理接口）
 */
export async function getAnnouncements(env: Env): Promise<Announcement[]> {
  const result = await env.DB.prepare(
    'SELECT id, title, content, link, created_at FROM announcements ORDER BY created_at DESC'
  ).all<Announcement>();

  return result.results;
}
