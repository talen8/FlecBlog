import type { UserInfo } from '../../types/user';
import { ref, computed } from 'vue';
import {
  login,
  register,
  forgotPassword,
  resetPassword,
  getWechatQrcode,
  getWechatScene,
  logout,
  getOAuthUrl,
  getOAuthBindUrl,
} from './api/auth';
import { getUserProfile } from './api/user';

// ============ Token 管理 ============

const ACCESS_TOKEN_KEY = 'access_token';

const getStoredToken = (): string | null => {
  if (import.meta.server) return null;
  return localStorage.getItem(ACCESS_TOKEN_KEY);
};

/** 当前访问令牌（响应式，自动从 localStorage 读取） */
export const accessToken = ref<string | null>(getStoredToken());
/** 是否已登录（响应式计算属性） */
export const isLoggedIn = computed(() => !!accessToken.value && accessToken.value !== '');

/**
 * 获取响应式登录状态
 * @returns 是否已登录（computed ref）
 */
export const useAuth = () => isLoggedIn;

/**
 * 设置访问令牌（写入 localStorage 和响应式状态）
 * @param access - 访问令牌字符串
 */
export const setAccessToken = (access: string): void => {
  if (import.meta.client) {
    localStorage.setItem(ACCESS_TOKEN_KEY, access);
  }
  accessToken.value = access;
};

/**
 * 读取访问令牌（优先从 localStorage 读取）
 * @returns 访问令牌字符串或 null
 */
export const getAccessToken = (): string | null => {
  if (import.meta.client) {
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  }
  return accessToken.value;
};

/** 清除访问令牌（删除 localStorage 并重置响应式状态） */
export const clearAccessToken = (): void => {
  if (import.meta.client) {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
  }
  accessToken.value = null;
};

/**
 * 登出操作
 * 注：登出时 accessToken 已清除，这里通知后端清除 session。
 */
export const logoutUser = (): void => {
  clearAccessToken();
  logout().catch(() => {});
};

// ============ 认证流程 ============

/**
 * 认证操作集合
 * @returns loginUser - 邮箱登录
 * @returns registerUser - 邮箱注册
 * @returns loginWithToken - 凭 token 登录
 * @returns bindOAuthUrl - 获取 OAuth 绑定地址
 * @returns forgotPassword / resetPassword / getWechatQrcode / getWechatScene / getOAuthUrl - 直接暴露的 API 函数
 */
export const useAuthActions = () => {
  const loginUser = async (data: { email: string; password: string }) => {
    const res = await login(data);
    setAccessToken(res.access_token);
    return res;
  };

  const registerUser = async (data: {
    email: string;
    nickname: string;
    password: string;
    website?: string;
  }) => {
    const res = await register(data);
    setAccessToken(res.access_token);
    return res;
  };

  /** 凭 token 登录（微信扫码 / OAuth 回调）：写 token + 拉取用户信息 */
  const loginWithToken = async (token: string) => {
    setAccessToken(token);
    await useUser().fetchUserInfo();
  };

  const bindOAuthUrl = (provider: string, redirect: string) =>
    getOAuthBindUrl(provider, redirect, accessToken.value ?? '');

  return {
    loginUser,
    registerUser,
    loginWithToken,
    forgotPassword,
    resetPassword,
    getWechatQrcode,
    getWechatScene,
    getOAuthUrl,
    bindOAuthUrl,
  };
};

// ============ 邮箱绑定提示 ============

// 触发间隔配置
const GLOBAL_REMIND_INTERVAL = 12 * 60 * 60 * 1000; // 全局触发：12小时
const COMMENT_REMIND_INTERVAL = 10 * 60 * 1000; // 评论触发：10分钟

// 存储 key
const SKIP_TIME_KEY = 'bindEmailSkipTime';

// 全局状态
const showBindEmailModal = ref(false);

// 触发类型
type TriggerType = 'global' | 'comment';

/**
 * 邮箱绑定提示管理
 * - 全局触发（页面访问/刷新/路由切换）：间隔 12 小时
 * - 评论触发：间隔 10 分钟
 * - 关闭弹窗会重置计时器
 * @returns showBindEmailModal - 弹窗显隐状态
 * @returns triggerGlobal - 全局触发检查
 * @returns triggerOnComment - 评论触发检查
 * @returns onBindSuccess - 绑定成功回调
 * @returns onSkip - 跳过回调
 */
export function useBindEmail() {
  const shouldShowPrompt = async (
    trigger: TriggerType,
    userInfo?: UserInfo | null
  ): Promise<boolean> => {
    if (!isLoggedIn.value) return false;

    let user = userInfo;
    if (!user) {
      try {
        user = await getUserProfile();
      } catch {
        return false;
      }
    }

    if (!user?.is_virtual_email) return false;

    const skipTime = localStorage.getItem(SKIP_TIME_KEY);
    if (skipTime) {
      const elapsed = Date.now() - parseInt(skipTime, 10);
      const interval = trigger === 'comment' ? COMMENT_REMIND_INTERVAL : GLOBAL_REMIND_INTERVAL;
      if (elapsed < interval) return false;
    }

    return true;
  };

  const triggerGlobal = async (userInfo?: UserInfo | null) => {
    if (await shouldShowPrompt('global', userInfo)) {
      showBindEmailModal.value = true;
    }
  };

  const triggerOnComment = async () => {
    if (await shouldShowPrompt('comment')) {
      showBindEmailModal.value = true;
    }
  };

  const onBindSuccess = () => {
    localStorage.removeItem(SKIP_TIME_KEY);
  };

  const onSkip = () => {
    localStorage.setItem(SKIP_TIME_KEY, String(Date.now()));
  };

  return {
    showBindEmailModal,
    triggerGlobal,
    triggerOnComment,
    onBindSuccess,
    onSkip,
  };
}
