/**
 * 用户状态管理
 *
 * 管理用户登录状态和用户信息，提供用户信息获取和清除方法。
 * 用户信息在全局共享，使用 ref 而非 useState 以避免 SSR 问题。
 *
 * @module composables/useUser
 */

import type { UserInfo } from '@@/types';
import { getUserProfile } from '@/composables/useApi';

const userInfo = ref<UserInfo | null>(null);

/**
 * 用户状态管理 Composable
 *
 * 提供用户信息的获取、清除等功能。
 *
 * @returns {Object} 用户状态和方法
 * @returns {Ref<UserInfo | null>} userInfo - 用户信息
 * @returns {ComputedRef<string>} userAvatar - 用户头像 URL
 * @returns {ComputedRef<string>} userNickname - 用户昵称
 * @returns {ComputedRef<string>} userEmail - 用户邮箱
 * @returns {Function} fetchUserInfo - 获取用户信息
 * @returns {Function} clearUserInfo - 清除用户信息
 *
 * @example
 * ```ts
 * const { userInfo, fetchUserInfo, clearUserInfo } = useUser();
 *
 * // 获取用户信息
 * await fetchUserInfo();
 *
 * // 访问用户数据
 * console.log(userInfo.value?.nickname);
 * ```
 */
export const useUser = () => {
  const isLoggedIn = useAuth();

  /**
   * 获取当前登录用户信息
   *
   * 从后端 API 获取用户详细信息并更新本地状态。
   * 如果用户未登录，会清除本地用户信息。
   */
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

  /**
   * 清除本地用户信息
   *
   * 用于用户登出时清除本地缓存的用户数据。
   */
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
