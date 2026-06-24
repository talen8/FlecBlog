import { useSysConfig } from './useSysConfig';

const DEFAULT_CRAVATAR_URL = 'https://cravatar.cn/avatar/%s?s=200&d=robohash';

/**
 * 头像工具（支持自定义头像、Gravatar/Cravatar 回退、默认 SVG 占位）
 * @returns getAvatarUrl(user, size?) - 根据用户信息生成头像 URL
 */
export function useAvatar() {
  const { basicConfig } = useSysConfig();

  function getAvatarUrl(user: { avatar?: string; email_hash?: string }, size = 48): string {
    if (user.avatar) return user.avatar;
    if (!user.email_hash) {
      return `data:image/svg+xml,%3Csvg width='${size}' height='${size}' viewBox='0 0 48 48' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Ccircle cx='24' cy='24' r='24' fill='%23E5E7EB'/%3E%3Cpath d='M24 12C17.3726 12 12 17.3726 12 24C12 30.6274 17.3726 36 24 36C30.6274 36 36 30.6274 36 24C36 17.3726 30.6274 12 24 12Z' fill='%239CA3AF'/%3E%3Ccircle cx='24' cy='24' r='8' fill='%236B7280'/%3E%3C/svg%3E`;
    }
    const cravatarUrl = basicConfig.value.cravatar_url || DEFAULT_CRAVATAR_URL;
    const url = cravatarUrl.replace('%s', user.email_hash);
    return url.replace(/s=\d+/, `s=${size}`);
  }

  return { getAvatarUrl };
}
