import type { SubmitFeedbackParams } from '../../types/feedback';
import { submitFeedback, getFeedbackByTicketNo } from './api/feedback';

/**
 * 反馈管理（提交反馈 + 按工单号查询反馈状态）
 * @returns submit - 提交反馈，queryByTicket - 按工单号查询
 */
export function useFeedback() {
  const submit = (data: SubmitFeedbackParams) => submitFeedback(data);
  const queryByTicket = (ticketNo: string) => getFeedbackByTicketNo(ticketNo);
  return { submit, queryByTicket };
}
