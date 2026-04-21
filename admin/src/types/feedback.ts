// 举报类型
export type ReportType = 'copyright' | 'inappropriate' | 'summary' | 'suggestion';

// 反馈状态
export type FeedbackStatus = 'pending' | 'resolved' | 'closed';

// 反馈内容结构
export interface FeedbackContent {
  description: string;
  reason?: string;
  attachmentFiles?: string[];
}

// 反馈对象
export interface Feedback {
  id: number;
  ticket_no: string;
  report_url: string;
  report_type: ReportType;
  form_content: FeedbackContent;
  email: string;
  status: FeedbackStatus;
  admin_reply: string;
  reply_time?: string;
  user_agent: string;
  ip: string;
  feedback_time: string;
}

// 反馈列表查询参数
export interface FeedbackListQuery {
  page: number;
  page_size: number;
  keyword?: string; // 搜索关键词（工单号、投诉地址）
  report_type?: ReportType; // 反馈类型筛选
  status?: FeedbackStatus; // 状态筛选
  start_time?: string; // 反馈开始时间
  end_time?: string; // 反馈结束时间
}

// 反馈列表数据
export interface FeedbackListData {
  list: Feedback[];
  total: number;
  page: number;
  page_size: number;
}

// 更新反馈请求
export interface FeedbackUpdateRequest {
  status: FeedbackStatus;
  admin_reply?: string;
}
