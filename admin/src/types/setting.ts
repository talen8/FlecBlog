// 系统设置页各配置分组的表单类型

// BasicForm 基础配置表单
export interface BasicForm {
  author: string;
  author_email: string;
  author_desc: string;
  home_url: string;
  author_avatar: string;
  icp: string;
  police_record: string;
  admin_url: string;
  blog_url: string;
  title: string;
  subtitle: string;
  description: string;
  keywords: string;
  favicon: string;
  established: string;
  custom_head: string;
  custom_body: string;
  emojis: string;
  meting_api: string;
  cravatar_url: string;
  ip_api_url: string;
  cover_maker_api: string;
  worker_proxy: string;
}

// NotificationForm 通知配置表单
export interface NotificationForm {
  email_host: string;
  email_port: string;
  email_secure: string;
  email_username: string;
  email_from: string;
  email_password: string;
  feishu_webhook_url: string;
  feishu_secret: string;
  telegram_bot_token: string;
  telegram_chat_id: string;
}

// UploadForm 上传配置表单
export interface UploadForm {
  storage_type: string;
  max_file_size: number;
  path_pattern: string;
  access_key: string;
  secret_key: string;
  region: string;
  bucket: string;
  endpoint: string;
  domain: string;
  use_ssl: boolean;
}

// AIForm AI 配置表单
export interface AIForm {
  base_url: string;
  api_key: string;
  model: string;
  summary_prompt: string;
  ai_summary_prompt: string;
  title_prompt: string;
  mcp_secret: string;
}

// OAuthForm OAuth 配置表单
export interface OAuthForm {
  'github.enabled': string;
  'github.client_id': string;
  'github.client_secret': string;
  'github.redirect_url': string;
  'google.enabled': string;
  'google.client_id': string;
  'google.client_secret': string;
  'google.redirect_url': string;
  'qq.enabled': string;
  'qq.client_id': string;
  'qq.client_secret': string;
  'qq.redirect_url': string;
  'microsoft.enabled': string;
  'microsoft.client_id': string;
  'microsoft.client_secret': string;
  'microsoft.redirect_url': string;
  'wechat.enabled': string;
  'wechat.appid': string;
  'wechat.secret': string;
  'oidc.enabled': string;
  'oidc.issuer_url': string;
  'oidc.client_id': string;
  'oidc.client_secret': string;
  'oidc.redirect_url': string;
}
