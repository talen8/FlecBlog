/**
 * 全局配置文件
 */

import type { SiteBasicConfig, BlogConfig } from './types';

// ==================== API 配置 ====================

export const API_CONFIG = {
  BASE_URL: 'http://localhost:8080/api/v1',
  TIMEOUT: 30000,
  DEFAULT_PAGE_SIZE: 10,
  MAX_PAGE_SIZE: 100,
} as const;

// ==================== 页面配置 ====================

export const PAGES = {
  INDEX: '/pages/index/index',
} as const;

// ==================== 存储键名配置 ====================

export const STORAGE_KEYS = {
  SITE_CONFIG: 'site_config',
  BLOG_CONFIG: 'blog_config',
} as const;

// ==================== 默认配置 ====================

export const DEFAULT_SITE_CONFIG: Partial<SiteBasicConfig> = {
  site_name: 'FlecBlog',
  site_description: '一个简洁优雅的博客系统',
};

export const DEFAULT_BLOG_CONFIG: Partial<BlogConfig> = {
  blog_title: 'FlecBlog',
  blog_subtitle: '记录生活，分享技术',
  page_size: 10,
  enable_comment: true,
};

// ==================== 错误提示信息 ====================

export const ERROR_MESSAGES = {
  NETWORK_ERROR: '网络连接失败，请检查网络设置',
  TIMEOUT_ERROR: '请求超时，请稍后重试',
  SERVER_ERROR: '服务器繁忙，请稍后重试',
  UNAUTHORIZED: '登录已过期，请重新登录',
  REQUEST_FAILED: '请求失败，请稍后重试',
  NO_DATA: '暂无数据',
  LOAD_FAILED: '加载失败，点击重试',
} as const;
