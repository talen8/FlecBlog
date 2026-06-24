import type { FriendGroupedResponse, FriendQueryParams, FriendApplyRequest } from '../../../types/friend';
import { createApi } from './createApi';

const friendApi = createApi<FriendGroupedResponse>('/friends');

/** 获取友链列表 */
export const getFriends = async (params?: FriendQueryParams) => {
  return friendApi.get('', params);
};

/** 申请友链 */
export const applyFriend = async (data: FriendApplyRequest) => {
  return friendApi.post('/apply', data);
};
