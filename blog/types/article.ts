export interface ArticleNav {
  title: string;
  url: string;
}

export interface Article {
  id: number;
  title: string;
  slug?: string;
  url: string;
  content?: string;
  summary: string;
  ai_summary?: string;
  excerpt?: string;
  cover?: string;
  is_top: boolean;
  is_essence: boolean;
  is_outdated?: boolean;
  view_count?: number;
  comment_count: number;
  publish_time: string;
  update_time: string;
  location?: string;
  category: {
    id: number;
    name: string;
    url: string;
  };
  tags: Array<{
    id: number;
    name: string;
    url: string;
  }>;
  prev?: ArticleNav;
  next?: ArticleNav;
}

/**
 * 文章查询参数（所有参数可选，支持灵活组合）
 */
export interface ArticleQuery {
  page?: number; // 页码；不传返回全部数据
  page_size?: number; // 每页数量；不传返回全部数据
  year?: string; // 按年份筛选
  month?: string; // 按月份筛选
  category?: string; // 按分类筛选（slug）
  tag?: string; // 按标签筛选（slug）
}
