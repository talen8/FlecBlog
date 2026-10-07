import type {
  FriendGroupedResponse,
  FriendQueryParams,
  FriendApplyRequest,
  FriendCircleItem,
  FriendCircleResponse,
  FriendCircleStats,
  FriendCircleQueryParams,
} from '../../../types/friend';
import { createApi } from './createApi';

const friendApi = createApi<FriendGroupedResponse>('/friends');

/** 获取友链列表 */
export const getFriends = async (params?: FriendQueryParams) => {
  return friendApi.get('', params);
};

/** 获取友圈文章 */
export const getFriendCircle = async (params?: FriendCircleQueryParams) => {
  return friendApi.get<FriendCircleResponse>('/circle', params);
};

/** 随机获取一篇友圈文章 */
export const getFriendCircleRandom = async () => {
  return friendApi.get<FriendCircleItem>('/circle/random');
};

/** 获取友圈统计数据 */
export const getFriendCircleStats = async () => {
  return friendApi.get<FriendCircleStats>('/circle/stats');
};

/** 申请友链 */
export const applyFriend = async (data: FriendApplyRequest) => {
  return friendApi.post('/apply', data);
};
