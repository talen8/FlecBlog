import type { Feedback, SubmitFeedbackParams } from '../../../types/feedback';
import { createApi } from './createApi';

const feedbackApi = createApi<Feedback>('/feedback');

/** 提交反馈 */
export const submitFeedback = async (data: SubmitFeedbackParams) => {
  return feedbackApi.post('', data);
};

/** 查询反馈状态 */
export const getFeedbackByTicketNo = async (ticketNo: string) => {
  return feedbackApi.get(`/ticket/${ticketNo}`);
};
