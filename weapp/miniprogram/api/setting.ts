/**
 * 配置相关 API 接口
 * 对接后端公开配置接口
 */

import { get } from '../utils/request';

/**
 * 获取站点基础配置（公开接口）
 * @returns 站点基础配置
 */
export function getSiteBasicConfig(): Promise<Record<string, string>> {
  return get<Record<string, string>>('/settings/basic');
}

/**
 * 获取博客配置（公开接口）
 * @returns 博客配置
 */
export function getBlogConfig(): Promise<Record<string, string>> {
  return get<Record<string, string>>('/settings/blog');
}

/**
 * 获取 OAuth 配置（公开接口）
 * 只返回各提供商的启用状态
 * @returns OAuth 配置
 */
export function getOAuthConfig(): Promise<Record<string, string>> {
  return get<Record<string, string>>('/settings/oauth');
}
