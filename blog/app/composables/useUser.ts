import type { UserInfo } from '@@/types/user';
import { getUserProfile } from '@/composables/api/user';

// 用户信息状态
const userInfo = ref<UserInfo | null>(null);

/**
 * 用户信息 Composable
 */
export const useUser = () => {
  const isLoggedIn = useAuth();

  // 获取用户信息
  const fetchUserInfo = async () => {
    if (!isLoggedIn.value) {
      userInfo.value = null;
      return;
    }

    try {
      const data = await getUserProfile();
      userInfo.value = data;
    } catch (error) {
      console.error('获取用户信息失败:', error);
      userInfo.value = null;
    }
  };

  // 清除用户信息
  const clearUserInfo = () => {
    userInfo.value = null;
  };

  return {
    userInfo,
    userAvatar: computed(() => getAvatarUrl(userInfo.value || {})),
    userNickname: computed(() => userInfo.value?.nickname || '用户'),
    userEmail: computed(() => userInfo.value?.email || ''),
    fetchUserInfo,
    clearUserInfo,
  };
};
