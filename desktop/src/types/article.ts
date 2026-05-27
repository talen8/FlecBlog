// 文章实体
export interface Article {
  id: number;
  title: string;
  slug: string;
  content: string;
  summary?: string;
  ai_summary?: string;
  cover?: string;
  is_publish: boolean; // 是否已发布
  is_top: boolean;
  is_essence: boolean;
  is_outdated: boolean; // 是否过时
  view_count: number;
  comment_count: number;
  location?: string;
  publish_time?: string;
  update_time?: string;
  category?: {
    id: number;
    name: string;
  };
  tags?: Array<{
    id: number;
    name: string;
  }>;
}

// 创建文章请求
export interface CreateArticleRequest {
  title: string;
  slug?: string;
  content: string;
  summary?: string;
  cover?: string;
  category_id?: number;
  tag_ids?: number[];
  location?: string;
  is_top?: boolean;
  is_essence?: boolean;
  is_outdated?: boolean; // 是否过时
}

// 更新文章请求
export interface UpdateArticleRequest {
  title?: string;
  slug?: string;
  content?: string;
  summary?: string;
  cover?: string;
  category_id?: number;
  tag_ids?: number[];
  location?: string;
  is_top?: boolean;
  is_essence?: boolean;
  is_outdated?: boolean; // 是否过时
  publish_time?: string;
  update_time?: string;
}

// 分页数据
export interface ArticleListData {
  list: Article[];
  total: number;
  page: number;
  page_size: number;
}

// 文章导入相关类型
export interface ImportArticleError {
  filename: string;
  title: string;
  error: string;
}

export interface ImportArticlesResult {
  total: number;
  success: number;
  failed: number;
  categories_added: number;
  tags_added: number;
  errors?: ImportArticleError[];
}

// 微信公众号导出结果
export interface WeChatExportResult {
  html: string; // 公众号 HTML，用于复制粘贴到微信公众平台
}

// 文章列表查询参数
export interface ArticleListQuery {
  page?: number;
  page_size?: number;
  keyword?: string;
  category_id?: number;
  tag_ids?: number[];
  location?: string;
  is_publish?: boolean;
  is_top?: boolean;
  is_essence?: boolean;
  is_outdated?: boolean;
  start_time?: string;
  end_time?: string;
}
