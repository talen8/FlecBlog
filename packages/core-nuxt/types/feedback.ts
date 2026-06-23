// 反馈类型
export type ReportType = 'copyright' | 'inappropriate' | 'summary' | 'suggestion';

// 反馈状态
export type FeedbackStatus = 'pending' | 'resolved' | 'closed';

// 反馈对象
export interface Feedback {
  id: number;
  ticket_no: string;
  report_url: string;
  report_type: ReportType;
  form_content: Record<string, unknown>;
  email: string;
  status: FeedbackStatus;
  admin_reply: string;
  reply_time?: string;
  user_agent: string;
  ip: string;
  feedback_time: string;
}

// 提交反馈参数
export interface SubmitFeedbackParams {
  reportUrl: string;
  reportType: ReportType;
  email?: string;
  description: string;
  reason?: string;
  attachmentFiles?: string[];
}
