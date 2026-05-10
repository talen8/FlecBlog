/**
 * 头像工具函数
 *
 * 提供用户头像 URL 的生成功能，支持自定义头像和 Gravatar 头像。
 *
 * @module utils/avatar
 */

/**
 * 获取用户头像 URL
 *
 * 根据用户信息生成头像 URL，优先级：
 * 1. 用户自定义头像
 * 2. 基于 email_hash 的 Cravatar 头像
 * 3. 默认 SVG 占位头像
 *
 * @param user - 用户对象
 * @param user.avatar - 用户自定义头像 URL
 * @param user.email_hash - 用户邮箱的 MD5 哈希值
 * @param size - 头像尺寸（像素），默认 48
 * @returns 头像 URL
 *
 * @example
 * ```ts
 * // 使用自定义头像
 * getAvatarUrl({ avatar: 'https://example.com/avatar.jpg' });
 *
 * // 使用 Gravatar 头像
 * getAvatarUrl({ email_hash: 'abc123...' });
 *
 * // 使用默认头像
 * getAvatarUrl({});
 * ```
 */
export function getAvatarUrl(user: { avatar?: string; email_hash?: string }, size = 48): string {
  if (user.avatar) return user.avatar;
  if (!user.email_hash) {
    return `data:image/svg+xml,%3Csvg width='${size}' height='${size}' viewBox='0 0 48 48' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='24' cy='24' r='24' fill='%23E5E7EB'/%3E%3Cpath d='M24 12C17.3726 12 12 17.3726 12 24C12 30.6274 17.3726 36 24 36C30.6274 36 36 30.6274 36 24C36 17.3726 30.6274 12 24 12Z' fill='%239CA3AF'/%3E%3Ccircle cx='24' cy='24' r='8' fill='%236B7280'/%3E%3C/svg%3E`;
  }
  return `https://cravatar.cn/avatar/${user.email_hash}?d=robohash&s=${size}`;
}
