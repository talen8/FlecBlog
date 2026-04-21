// RSS文章实体
export interface RssArticle {
  id: number;
  friend_id: number;
  friend_name: string;
  friend_url: string;
  title: string;
  link: string;
  description: string;
  published_at?: string;
  is_read: boolean;
  created_at: string;
}

// RSS文章列表查询参数
export interface RssArticleQuery {
  page?: number;
  page_size?: number;
  keyword?: string; // 搜索关键词
  friend_id?: number; // 友链ID筛选
  is_read?: boolean; // 已读状态筛选
  start_time?: string; // 发布开始时间
  end_time?: string; // 发布结束时间
}

// RSS文章列表响应数据
export interface RssArticleListData {
  list: RssArticle[];
  total: number;
  page: number;
  page_size: number;
  unread_count: number;
}
