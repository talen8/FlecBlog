import type { FriendApplyRequest } from '../../types/friend';
import { getFriends, applyFriend } from './api/friend';

/**
 * 获取友链分组列表（支持 SSR）
 * @returns data - 友链分组数组，refresh - 刷新方法
 */
export function useFriends() {
  return useAsyncData('friends', async () => {
    const { groups } = await getFriends();
    return groups ?? [];
  });
}

/**
 * 提交友链申请
 * @returns submitting - 提交中状态，apply(data) - 执行申请
 */
export function useFriendApply() {
  const submitting = ref(false);

  const apply = async (data: FriendApplyRequest) => {
    submitting.value = true;
    try {
      await applyFriend(data);
    } finally {
      submitting.value = false;
    }
  };

  return { submitting, apply };
}
