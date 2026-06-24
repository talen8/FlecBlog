import { getSettingGroup } from './api/sysconfig';
import type { SysConfigData } from '../../types/sysconfig';

// 去掉 key 的指定前缀
const stripPrefix = (config: Record<string, string>, prefix: string) => {
  const p = `${prefix}.`;
  const result: Record<string, string> = {};
  for (const [key, value] of Object.entries(config)) {
    result[key.startsWith(p) ? key.slice(p.length) : key] = value;
  }
  return result;
};

const defaultBasic: Record<string, string> = {
  author: '',
  author_avatar: '',
  icp: '',
  police_record: '',
  admin_url: '',
  blog_url: '',
  title: 'FlecBLOG',
  subtitle: 'FlecBLOG',
  description: '',
  keywords: '',
  favicon: '',
  custom_head: '',
  custom_body: '',
  emojis: '',
  meting_api: '',
  cravatar_url: '',
  ip_api_url: '',
  cover_maker_api: '',
};

const defaultUpload: Record<string, string> = {
  storage_type: 'local',
  max_file_size: '5',
  path_pattern: '',
  access_key: '',
  secret_key: '',
  region: '',
  bucket: '',
  endpoint: '',
  domain: '',
  use_ssl: 'true',
};

const defaultOauth: Record<string, string> = {
  session_secret: '',
  worker_proxy: '',
  'github.enabled': 'false',
  'github.client_id': '',
  'github.client_secret': '',
  'github.redirect_url': '',
  'google.enabled': 'false',
  'google.client_id': '',
  'google.client_secret': '',
  'google.redirect_url': '',
  'qq.enabled': 'false',
  'qq.client_id': '',
  'qq.client_secret': '',
  'qq.redirect_url': '',
  'microsoft.enabled': 'false',
  'microsoft.client_id': '',
  'microsoft.client_secret': '',
  'microsoft.redirect_url': '',
  'oidc.enabled': 'false',
  'oidc.issuer_url': '',
  'oidc.client_id': '',
  'oidc.client_secret': '',
  'oidc.redirect_url': '',
  'wechat.enabled': 'false',
  'wechat.appid': '',
  'wechat.secret': '',
};

/**
 * 系统配置（支持 SSR）
 * @returns basicConfig - 基础配置（标题/描述/favicon 等）
 * @returns uploadConfig - 上传配置
 * @returns oauthConfig - OAuth 配置
 */
export function useSysConfig() {
  const { data } = useAsyncData('sysconfig-fetch', async () => {
    const [basic, upload, oauth] = await Promise.all([
      getSettingGroup('basic'),
      getSettingGroup('upload'),
      getSettingGroup('oauth'),
    ]);
    return {
      basic: stripPrefix(basic, 'basic'),
      upload: stripPrefix(upload, 'upload'),
      oauth: stripPrefix(oauth, 'oauth'),
    } satisfies SysConfigData;
  });

  return {
    basicConfig: computed(() => data.value?.basic ?? defaultBasic),
    uploadConfig: computed(() => data.value?.upload ?? defaultUpload),
    oauthConfig: computed(() => data.value?.oauth ?? defaultOauth),
  };
}
