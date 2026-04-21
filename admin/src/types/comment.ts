// 评论实体
export interface Comment {
  id: number;
  content: string;
  status: number; // 0: 隐藏, 1: 显示
  parent_id?: number; // 父评论ID，用于回复
  created_at: string;
  deleted_at?: string;
  target: {
    type: string; // 评论目标类型：article(文章)、page(页面)等
    key: string; // 目标标识符
    title: string; // 目标标题
  };
  user: {
    id: number;
    nickname: string;
    email: string;
    avatar: string;
  };
}

// 分页数据
export interface CommentListData {
  list: Comment[];
  total: number;
  page: number;
  page_size: number;
}

// 评论列表查询参数
export interface CommentListQuery {
  page: number;
  page_size: number;
  keyword?: string; // 搜索关键词（评论内容）
  status?: number; // 状态筛选 0:隐藏 1:显示
  is_deleted?: boolean; // 是否已删除
  is_sub?: boolean; // 是否子评论（true:子评论, false:父评论）
  start_time?: string; // 评论开始时间
  end_time?: string; // 评论结束时间
}

// 评论导入相关类型
export interface ImportCommentError {
  index: number;
  content: string;
  error: string;
}

export interface ImportCommentsResult {
  total: number;
  success: number;
  failed: number;
  user_created: number;
  errors?: ImportCommentError[];
}
