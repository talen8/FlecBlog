<script setup lang="ts">
import {
  getUserProfile,
  updateUserProfile,
  changePassword,
  setPassword,
  deactivateAccount,
  unbindOAuth,
} from '@/composables/api/user';
import type { UserInfo, UserRole } from '@@/types/user';

definePageMeta({});

useSeoMeta({
  title: '个人信息',
  description: '管理和编辑您的个人资料、账户设置和登录方式',
});

const router = useRouter();
const route = useRoute();
const config = useRuntimeConfig();
const { success: showSuccess, error: showError } = useToast();

const userInfo = ref<UserInfo | null>(null);
const showEditDialog = ref(false);
const showBadgeDialog = ref(false);
const showPasswordDialog = ref(false);
const showSetPasswordDialog = ref(false);
const showDeactivateDialog = ref(false);
const showUnbindDialog = ref(false);
const unbindProvider = ref('');

// ===== 通用验证函数 =====
const validators = {
  nickname: (val: string) => {
    if (!val.trim()) return '昵称不能为空';
    if (val.length < 2 || val.length > 32) return '昵称长度需要在2-32个字符之间';
    return '';
  },
  email: (val: string) => {
    if (!val.trim()) return '邮箱不能为空';
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(val)) return '请输入正确的邮箱格式';
    return '';
  },
  website: (val: string) => {
    if (val && !/^https?:\/\/.+/.test(val))
      return '网站地址格式不正确，请以 http:// 或 https:// 开头';
    return '';
  },
  badge: (val: string) => {
    const forbidden = ['站长', '博主', '管理员', 'admin', 'root', 'super_admin'];
    if (forbidden.includes(val)) return '禁止使用该铭牌';
    return '';
  },
  password: (val: string, isNew = false) => {
    if (!val) return isNew ? '请输入新密码' : '请输入密码';
    if (val.length < 6 || val.length > 32) return '密码长度需要在6-32个字符之间';
    return '';
  },
};

// ===== 基础功能 =====
const getRoleName = (role: UserRole): string => {
  const roleMap: Record<UserRole, string> = {
    super_admin: '超级管理员',
    admin: '管理员',
    user: '普通用户',
  };
  return roleMap[role] || role;
};

const fetchProfile = async () => {
  const data = await getUserProfile();
  userInfo.value = data;
};

const handleLogout = () => {
  logout();
  router.push('/');
};

// ===== 编辑资料对话框 =====
const editForm = ref({ nickname: '', email: '', website: '', avatar: '' });
const editLoading = ref(false);
const uploading = ref(false);
const editErrors = ref<Record<string, string>>({});

watch(showEditDialog, val => {
  if (val && userInfo.value) {
    editForm.value = {
      nickname: userInfo.value.nickname || '',
      email: userInfo.value.email || '',
      website: userInfo.value.website || '',
      avatar: userInfo.value.avatar || '',
    };
  }
  editErrors.value = {};
});

const handleAvatarUpload = async () => {
  if (uploading.value) return;
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = 'image/*';
  input.onchange = async e => {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;
    try {
      uploading.value = true;
      editErrors.value = {};
      const result = await uploadFile(file, '用户头像');
      editForm.value.avatar = result.file_url;
    } catch (error: unknown) {
      const err = error as Error;
      editErrors.value.avatar = err.message || '头像上传失败';
    } finally {
      uploading.value = false;
    }
  };
  input.click();
};

const handleEditSubmit = async () => {
  editErrors.value = {};
  const { nickname, email, website } = editForm.value;

  const nicknameError = validators.nickname(nickname);
  if (nicknameError) return (editErrors.value.nickname = nicknameError);

  const emailError = validators.email(email);
  if (emailError) return (editErrors.value.email = emailError);

  const websiteError = validators.website(website.trim());
  if (websiteError) return (editErrors.value.website = websiteError);

  editLoading.value = true;
  try {
    const data = await updateUserProfile({
      nickname: nickname.trim(),
      email: email.trim(),
      website: website.trim() || undefined,
      avatar: editForm.value.avatar || undefined,
    });
    showSuccess('保存成功');
    setTimeout(() => {
      userInfo.value = data;
      showEditDialog.value = false;
    }, 1000);
  } catch (error: unknown) {
    const err = error as Error;
    showError(err.message || '保存失败');
  } finally {
    editLoading.value = false;
  }
};

// ===== 铭牌设置对话框 =====
const badge = ref('');
const badgeLoading = ref(false);
const badgeError = ref('');

watch(showBadgeDialog, val => {
  if (val) {
    badge.value = userInfo.value?.badge || '';
    badgeError.value = '';
  }
});

const handleBadgeSubmit = async () => {
  badgeError.value = validators.badge(badge.value);
  if (badgeError.value) return;

  badgeLoading.value = true;
  try {
    const data = await updateUserProfile({
      nickname: userInfo.value!.nickname,
      email: userInfo.value!.email,
      badge: badge.value.trim() || undefined,
    });
    showSuccess('铭牌设置成功');
    setTimeout(() => {
      userInfo.value = data;
      showBadgeDialog.value = false;
    }, 1000);
  } catch (error: unknown) {
    const err = error as Error;
    showError(err.message || '设置失败');
  } finally {
    badgeLoading.value = false;
  }
};

// ===== 修改密码对话框 =====
const passwordForm = ref({ old_password: '', new_password: '', confirm_password: '' });
const passwordLoading = ref(false);
const passwordErrors = ref<Record<string, string>>({});

watch(showPasswordDialog, val => {
  if (val) passwordForm.value = { old_password: '', new_password: '', confirm_password: '' };
  passwordErrors.value = {};
});

const handlePasswordSubmit = async () => {
  passwordErrors.value = {};
  const { old_password, new_password, confirm_password } = passwordForm.value;

  if (!old_password) return (passwordErrors.value.old_password = '请输入旧密码');

  const pwdError = validators.password(new_password, true);
  if (pwdError) return (passwordErrors.value.new_password = pwdError);

  if (!confirm_password) return (passwordErrors.value.confirm_password = '请确认新密码');
  if (new_password !== confirm_password)
    return (passwordErrors.value.confirm_password = '两次输入的密码不一致');
  if (old_password === new_password)
    return (passwordErrors.value.new_password = '新密码不能与旧密码相同');

  passwordLoading.value = true;
  try {
    await changePassword({ old_password, new_password });
    showSuccess('密码修改成功，请重新登录');
    showPasswordDialog.value = false;
    setTimeout(() => {
      logout();
      router.push('/');
    }, 1500);
  } catch (error: unknown) {
    const err = error as Error & { message?: string };
    const errorMsg = err?.message || '密码修改失败';
    if (errorMsg.includes('旧密码')) {
      passwordErrors.value.old_password = '原密码错误';
    } else {
      showError(errorMsg);
    }
  } finally {
    passwordLoading.value = false;
  }
};

// ===== 设置密码对话框（OAuth 用户首次设置密码）=====
const setPasswordForm = ref({ password: '', confirm_password: '' });
const setPasswordLoading = ref(false);
const setPasswordErrors = ref<Record<string, string>>({});

watch(showSetPasswordDialog, val => {
  if (val) setPasswordForm.value = { password: '', confirm_password: '' };
  setPasswordErrors.value = {};
});

const handleSetPasswordSubmit = async () => {
  setPasswordErrors.value = {};
  const { password, confirm_password } = setPasswordForm.value;

  const pwdError = validators.password(password, true);
  if (pwdError) return (setPasswordErrors.value.password = pwdError);

  if (!confirm_password) return (setPasswordErrors.value.confirm_password = '请确认密码');
  if (password !== confirm_password)
    return (setPasswordErrors.value.confirm_password = '两次输入的密码不一致');

  setPasswordLoading.value = true;
  try {
    await setPassword({ password, confirm_password });
    showSuccess('密码设置成功');
    showSetPasswordDialog.value = false;
    await fetchProfile(); // 刷新用户信息
  } catch (error: unknown) {
    const err = error as Error;
    showError(err?.message || '密码设置失败');
  } finally {
    setPasswordLoading.value = false;
  }
};

// ===== 注销账户对话框 =====
const deactivatePassword = ref('');
const deactivateConfirmed = ref(false);
const deactivateLoading = ref(false);
const deactivateErrors = ref<Record<string, string>>({});

watch(showDeactivateDialog, val => {
  if (val) {
    deactivatePassword.value = '';
    deactivateConfirmed.value = false;
  }
  deactivateErrors.value = {};
});

const { oauthConfig } = useSysConfig();

// 处理登录方式点击
const handleLoginMethodClick = (provider: string, enabled: boolean) => {
  // 密码登录方式特殊处理
  if (provider === 'password') {
    if (!enabled) {
      // 未设置密码，打开设置密码对话框
      showSetPasswordDialog.value = true;
    }
    // 已设置密码时不做任何操作（密码一旦设置不可撤销）
    return;
  }

  // OAuth 方式处理
  if (enabled) {
    // 已绑定，检查是否可以解绑
    const loginCount = getLoginMethodCount();
    if (loginCount <= 1) {
      showError('至少保留一种登录方式，无法解绑');
      return;
    }
    // 显示解绑对话框
    unbindProvider.value = provider;
    showUnbindDialog.value = true;
  } else {
    // 未绑定，跳转绑定
    bindOAuth(provider);
  }
};

// 计算当前登录方式数量
const getLoginMethodCount = () => {
  if (!userInfo.value) return 0;
  let count = 0;
  if (userInfo.value.has_password) count++;
  if (userInfo.value.linked_oauths?.includes('github')) count++;
  if (userInfo.value.linked_oauths?.includes('google')) count++;
  if (userInfo.value.linked_oauths?.includes('qq')) count++;
  if (userInfo.value.linked_oauths?.includes('microsoft')) count++;
  return count;
};

// OAuth 绑定 loading 状态
const oauthBindLoading = ref<string | null>(null);

// 绑定 OAuth
const bindOAuth = (provider: string) => {
  oauthBindLoading.value = provider;
  const apiUrl = config.public.apiUrl;
  const token = accessToken.value;
  const redirect = encodeURIComponent(route.fullPath);
  window.location.href = `${apiUrl}/auth/${provider}?action=bind&token=${token}&redirect=${redirect}`;
};

// 解绑 OAuth 相关
const unbindLoading = ref(false);

const getProviderName = (provider: string) => {
  const names: Record<string, string> = {
    github: 'GitHub',
    google: 'Google',
    qq: 'QQ',
    microsoft: 'Microsoft',
  };
  return names[provider] || provider;
};

const handleUnbindSubmit = async () => {
  unbindLoading.value = true;
  try {
    await unbindOAuth(unbindProvider.value);
    showSuccess(`已解绑 ${getProviderName(unbindProvider.value)}`);
    showUnbindDialog.value = false;
    await fetchProfile();
  } catch (error: unknown) {
    const err = error as Error;
    showError(err?.message || '解绑失败');
  } finally {
    unbindLoading.value = false;
  }
};

const handleDeactivateSubmit = async () => {
  deactivateErrors.value = {};
  if (!deactivateConfirmed.value)
    return (deactivateErrors.value.confirmed = '请确认您已了解注销账户的后果');
  if (!deactivatePassword.value) return (deactivateErrors.value.password = '请输入密码以确认身份');

  deactivateLoading.value = true;
  try {
    await deactivateAccount({ password: deactivatePassword.value });
    showSuccess('账户注销成功');
    showDeactivateDialog.value = false;
    setTimeout(() => {
      logout();
      router.push('/');
    }, 1500);
  } catch (error: unknown) {
    const err = error as Error;
    const errorMsg = err?.message || '账户注销失败';
    if (errorMsg.includes('密码') || errorMsg.includes('password')) {
      deactivateErrors.value.password = '密码错误，请重新输入';
    } else {
      showError(errorMsg);
    }
  } finally {
    deactivateLoading.value = false;
  }
};

onMounted(async () => {
  if (!accessToken.value) {
    router.push('/');
    return;
  }

  // 只在客户端获取需要认证的数据
  if (import.meta.client) {
    await fetchProfile();

    // 处理 OAuth 绑定回调消息
    const bindStatus = route.query.bind as string;
    const provider = route.query.provider as string;

    if (bindStatus === 'success' && provider) {
      showSuccess(`已成功绑定 ${getProviderName(provider)}`);
      // 清除 URL 参数
      router.replace({ path: '/profile' });
    } else if (bindStatus === 'error') {
      const message = route.query.message as string;
      showError(message || '绑定失败');
      router.replace({ path: '/profile' });
    }
  }
});
</script>

<template>
  <div id="page">
    <h1 class="page-title">个人信息</h1>

    <div v-if="userInfo" class="profile-content">
      <!-- 基础信息 -->
      <div class="info-card">
        <div class="card-header">
          <h3 class="card-title">基础信息</h3>
          <button class="btn-primary" @click="showEditDialog = true">
            <i class="ri-edit-line" /> 编辑资料
          </button>
        </div>
        <div class="info-list">
          <div class="info-item">
            <span class="label">头像</span>
            <span class="value">
              <NuxtImg
                :src="getAvatarUrl(userInfo)"
                alt="头像"
                class="avatar-preview"
                loading="lazy"
              />
            </span>
          </div>
          <div class="info-item">
            <span class="label">昵称</span>
            <span class="value">{{ userInfo.nickname || '未设置' }}</span>
          </div>
          <div class="info-item">
            <span class="label">邮箱</span>
            <span class="value">{{ userInfo.email }}</span>
          </div>
          <div class="info-item">
            <span class="label">个人网站</span>
            <span class="value">
              <span v-if="userInfo.website">{{ userInfo.website }}</span>
              <span v-else class="empty">未设置</span>
            </span>
          </div>
        </div>
      </div>

      <!-- 账户信息 -->
      <div class="info-card">
        <h3 class="card-title">账户信息</h3>
        <div class="info-list">
          <div class="info-item">
            <span class="label">铭牌标识</span>
            <span class="value">
              <span v-if="userInfo.badge" class="badge-text">{{ userInfo.badge }}</span>
              <span v-else class="empty">未设置</span>
              <button class="btn-icon" title="设置铭牌" @click="showBadgeDialog = true">
                <i class="ri-settings-4-line" />
              </button>
            </span>
          </div>
          <div class="info-item">
            <span class="label">登录方式</span>
            <span class="value">
              <div class="login-methods">
                <!-- 密码登录方式 -->
                <div
                  class="method-icon"
                  :class="{
                    enabled: userInfo?.has_password,
                    disabled: !userInfo?.has_password,
                    clickable: !userInfo?.has_password,
                  }"
                  :title="userInfo?.has_password ? '密码登录' : '设置密码'"
                  @click="handleLoginMethodClick('password', userInfo?.has_password ?? false)"
                >
                  <i class="ri-lock-password-line" />
                </div>

                <!-- GitHub登录方式 -->
                <div
                  v-if="oauthConfig['github.enabled'] === 'true'"
                  class="method-icon clickable"
                  :class="{
                    enabled: userInfo?.linked_oauths?.includes('github'),
                    disabled: !userInfo?.linked_oauths?.includes('github'),
                    loading: oauthBindLoading === 'github',
                  }"
                  :title="
                    oauthBindLoading === 'github'
                      ? '绑定中...'
                      : userInfo?.linked_oauths?.includes('github')
                        ? 'GitHub'
                        : '绑定GitHub'
                  "
                  @click="
                    oauthBindLoading
                      ? null
                      : handleLoginMethodClick(
                          'github',
                          userInfo?.linked_oauths?.includes('github') ?? false
                        )
                  "
                >
                  <i v-if="oauthBindLoading === 'github'" class="ri-loader-4-line spin" />
                  <i v-else class="ri-github-fill" />
                </div>

                <!-- Google登录方式 -->
                <div
                  v-if="oauthConfig['google.enabled'] === 'true'"
                  class="method-icon clickable"
                  :class="{
                    enabled: userInfo?.linked_oauths?.includes('google'),
                    disabled: !userInfo?.linked_oauths?.includes('google'),
                    loading: oauthBindLoading === 'google',
                  }"
                  :title="
                    oauthBindLoading === 'google'
                      ? '绑定中...'
                      : userInfo?.linked_oauths?.includes('google')
                        ? 'Google'
                        : '绑定Google'
                  "
                  @click="
                    oauthBindLoading
                      ? null
                      : handleLoginMethodClick(
                          'google',
                          userInfo?.linked_oauths?.includes('google') ?? false
                        )
                  "
                >
                  <i v-if="oauthBindLoading === 'google'" class="ri-loader-4-line spin" />
                  <i v-else class="ri-google-fill" />
                </div>

                <!-- QQ登录方式 -->
                <div
                  v-if="oauthConfig['qq.enabled'] === 'true'"
                  class="method-icon clickable"
                  :class="{
                    enabled: userInfo?.linked_oauths?.includes('qq'),
                    disabled: !userInfo?.linked_oauths?.includes('qq'),
                    loading: oauthBindLoading === 'qq',
                  }"
                  :title="
                    oauthBindLoading === 'qq'
                      ? '绑定中...'
                      : userInfo?.linked_oauths?.includes('qq')
                        ? 'QQ'
                        : '绑定QQ'
                  "
                  @click="
                    oauthBindLoading
                      ? null
                      : handleLoginMethodClick(
                          'qq',
                          userInfo?.linked_oauths?.includes('qq') ?? false
                        )
                  "
                >
                  <i v-if="oauthBindLoading === 'qq'" class="ri-loader-4-line spin" />
                  <i v-else class="ri-qq-fill" />
                </div>

                <!-- Microsoft登录方式 -->
                <div
                  v-if="oauthConfig['microsoft.enabled'] === 'true'"
                  class="method-icon clickable"
                  :class="{
                    enabled: userInfo?.linked_oauths?.includes('microsoft'),
                    disabled: !userInfo?.linked_oauths?.includes('microsoft'),
                    loading: oauthBindLoading === 'microsoft',
                  }"
                  :title="
                    oauthBindLoading === 'microsoft'
                      ? '绑定中...'
                      : userInfo?.linked_oauths?.includes('microsoft')
                        ? 'Microsoft'
                        : '绑定Microsoft'
                  "
                  @click="
                    oauthBindLoading
                      ? null
                      : handleLoginMethodClick(
                          'microsoft',
                          userInfo?.linked_oauths?.includes('microsoft') ?? false
                        )
                  "
                >
                  <i v-if="oauthBindLoading === 'microsoft'" class="ri-loader-4-line spin" />
                  <i v-else class="ri-microsoft-fill" />
                </div>
              </div>
            </span>
          </div>
          <div class="info-item">
            <span class="label">角色权限</span>
            <span class="value">{{ getRoleName(userInfo.role) }}</span>
          </div>
          <div class="info-item">
            <span class="label">注册时间</span>
            <span class="value">{{ formatFriendly(userInfo.created_at) }}</span>
          </div>
          <div class="info-item">
            <span class="label">最后登录</span>
            <span class="value">{{ formatFriendly(userInfo.last_login) }}</span>
          </div>
        </div>
      </div>

      <!-- 账户管理 -->
      <div class="info-card danger-zone">
        <h3 class="card-title">账户管理</h3>
        <div class="info-list">
          <div class="info-item">
            <div class="action-item">
              <div class="action-info">
                <span class="action-title">退出登录</span>
                <span class="action-desc">退出当前账户，需要重新登录</span>
              </div>
              <button class="btn-secondary" @click="handleLogout">
                <i class="ri-logout-box-line" /> 退出
              </button>
            </div>
          </div>
          <div class="info-item">
            <div class="action-item">
              <div class="action-info">
                <span class="action-title">{{
                  userInfo?.has_password ? '修改密码' : '设置密码'
                }}</span>
                <span class="action-desc">{{
                  userInfo?.has_password ? '修改账户登录密码' : '设置账户登录密码'
                }}</span>
              </div>
              <button
                v-if="userInfo?.has_password"
                class="btn-secondary"
                @click="showPasswordDialog = true"
              >
                <i class="ri-lock-password-line" /> 修改密码
              </button>
              <button v-else class="btn-secondary" @click="showSetPasswordDialog = true">
                <i class="ri-lock-password-line" /> 设置密码
              </button>
            </div>
          </div>
          <div class="info-item">
            <div class="action-item">
              <div class="action-info">
                <span class="action-title">注销账户</span>
                <span class="action-desc">永久删除账户及所有数据，此操作不可恢复</span>
              </div>
              <button class="btn-danger" @click="showDeactivateDialog = true">
                <i class="ri-delete-bin-line" /> 注销账户
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑资料对话框 -->
    <UiBaseDialog
      v-model="showEditDialog"
      title="编辑个人信息"
      style="--dialog-width: 500px"
      :loading="editLoading"
      confirm-text="保存"
      @confirm="handleEditSubmit"
    >
      <form
        style="display: flex; flex-direction: column; gap: 20px"
        @submit.prevent="handleEditSubmit"
      >
        <div class="form-group">
          <label>头像</label>
          <div class="avatar-section">
            <NuxtImg
              :src="editForm.avatar || (userInfo ? getAvatarUrl(userInfo) : '')"
              alt="头像"
              loading="lazy"
            />
            <button type="button" :disabled="uploading || editLoading" @click="handleAvatarUpload">
              {{ uploading ? '上传中...' : '更换头像' }}
            </button>
          </div>
          <p v-if="editErrors.avatar" class="error-message">{{ editErrors.avatar }}</p>
        </div>

        <div class="form-group">
          <label class="form-label">昵称</label>
          <input
            v-model="editForm.nickname"
            type="text"
            class="form-input"
            :class="{ error: editErrors.nickname }"
            placeholder="请输入昵称（2-32个字符）"
            :disabled="editLoading"
          />
          <p v-if="editErrors.nickname" class="error-message">{{ editErrors.nickname }}</p>
        </div>

        <div class="form-group">
          <label class="form-label">邮箱</label>
          <input
            v-model="editForm.email"
            type="email"
            class="form-input"
            :class="{ error: editErrors.email }"
            placeholder="请输入邮箱"
            :disabled="editLoading"
          />
          <p v-if="editErrors.email" class="error-message">{{ editErrors.email }}</p>
        </div>

        <div class="form-group">
          <label class="form-label">网站</label>
          <input
            v-model="editForm.website"
            type="url"
            class="form-input"
            :class="{ error: editErrors.website }"
            placeholder="https://example.com（选填）"
            :disabled="editLoading"
          />
          <p v-if="editErrors.website" class="error-message">{{ editErrors.website }}</p>
        </div>
      </form>
    </UiBaseDialog>

    <!-- 铭牌设置对话框 -->
    <UiBaseDialog
      v-model="showBadgeDialog"
      title="设置铭牌标识"
      :loading="badgeLoading"
      confirm-text="保存"
      style="--dialog-width: 400px"
      @confirm="handleBadgeSubmit"
    >
      <form
        style="display: flex; flex-direction: column; gap: 20px"
        @submit.prevent="handleBadgeSubmit"
      >
        <div class="form-group">
          <input
            v-model="badge"
            type="text"
            class="form-input"
            :class="{ error: badgeError }"
            placeholder="请输入铭牌内容"
            :disabled="badgeLoading"
          />
          <p v-if="badgeError" class="error-message">{{ badgeError }}</p>
        </div>
      </form>
    </UiBaseDialog>

    <!-- 修改密码对话框 -->
    <UiBaseDialog
      v-model="showPasswordDialog"
      title="修改密码"
      :loading="passwordLoading"
      confirm-text="确认修改"
      @confirm="handlePasswordSubmit"
    >
      <form
        style="display: flex; flex-direction: column; gap: 20px"
        @submit.prevent="handlePasswordSubmit"
      >
        <div class="form-group">
          <label class="form-label">旧密码</label>
          <input
            v-model="passwordForm.old_password"
            type="password"
            class="form-input"
            :class="{ error: passwordErrors.old_password }"
            placeholder="请输入旧密码"
            :disabled="passwordLoading"
            autocomplete="current-password"
          />
          <p v-if="passwordErrors.old_password" class="error-message">
            {{ passwordErrors.old_password }}
          </p>
        </div>

        <div class="form-group">
          <label class="form-label">新密码</label>
          <input
            v-model="passwordForm.new_password"
            type="password"
            class="form-input"
            :class="{ error: passwordErrors.new_password }"
            placeholder="请输入新密码（6-32个字符）"
            :disabled="passwordLoading"
            autocomplete="new-password"
          />
          <p v-if="passwordErrors.new_password" class="error-message">
            {{ passwordErrors.new_password }}
          </p>
        </div>

        <div class="form-group">
          <label class="form-label">确认新密码</label>
          <input
            v-model="passwordForm.confirm_password"
            type="password"
            class="form-input"
            :class="{ error: passwordErrors.confirm_password }"
            placeholder="请再次输入新密码"
            :disabled="passwordLoading"
            autocomplete="new-password"
          />
          <p v-if="passwordErrors.confirm_password" class="error-message">
            {{ passwordErrors.confirm_password }}
          </p>
        </div>

        <div class="tip">
          <i class="ri-information-line" />
          <span>密码修改成功后，系统将自动退出登录，请使用新密码重新登录</span>
        </div>
      </form>
    </UiBaseDialog>

    <!-- 注销账户对话框 -->
    <UiBaseDialog
      v-model="showDeactivateDialog"
      title="注销账户"
      style="--dialog-width: 520px"
      :loading="deactivateLoading"
      confirm-text="确认注销"
      @confirm="handleDeactivateSubmit"
    >
      <div style="display: flex; flex-direction: column; gap: 24px">
        <div class="warning">
          <p class="warning-title">⚠️ 请谨慎操作，此操作不可恢复！</p>
          <p>注销账户后，您将无法再登录此账户，您的个人信息将被永久删除。</p>
        </div>

        <form
          style="display: flex; flex-direction: column; gap: 20px"
          @submit.prevent="handleDeactivateSubmit"
        >
          <div class="form-group">
            <label class="checkbox-label">
              <input v-model="deactivateConfirmed" type="checkbox" :disabled="deactivateLoading" />
              <span>我已充分了解注销账户的后果，并确认要注销我的账户</span>
            </label>
            <p v-if="deactivateErrors.confirmed" class="error-message">
              {{ deactivateErrors.confirmed }}
            </p>
          </div>

          <div class="form-group">
            <label class="form-label">输入密码以确认</label>
            <input
              v-model="deactivatePassword"
              type="password"
              class="form-input"
              :class="{ error: deactivateErrors.password }"
              placeholder="请输入您的账户密码"
              :disabled="deactivateLoading"
              autocomplete="current-password"
            />
            <p v-if="deactivateErrors.password" class="error-message">
              {{ deactivateErrors.password }}
            </p>
          </div>
        </form>
      </div>
    </UiBaseDialog>

    <!-- 解绑 OAuth 对话框 -->
    <UiBaseDialog
      v-model="showUnbindDialog"
      :title="`解绑 ${getProviderName(unbindProvider)}`"
      style="--dialog-width: 400px"
      :loading="unbindLoading"
      confirm-text="确认解绑"
      @confirm="handleUnbindSubmit"
    >
      <div style="display: flex; flex-direction: column; gap: 16px">
        <p>确定要解绑 {{ getProviderName(unbindProvider) }} 登录方式吗？</p>
        <p class="dialog-hint">解绑后，您将无法使用该方式登录。需至少保留一种登录方式。</p>
      </div>
    </UiBaseDialog>

    <!-- 设置密码对话框（OAuth 用户首次设置密码）-->
    <UiBaseDialog
      v-model="showSetPasswordDialog"
      title="设置密码"
      :loading="setPasswordLoading"
      confirm-text="确认设置"
      @confirm="handleSetPasswordSubmit"
    >
      <form
        style="display: flex; flex-direction: column; gap: 20px"
        @submit.prevent="handleSetPasswordSubmit"
      >
        <div class="tip">
          <i class="ri-information-line" />
          <span>设置密码后可使用邮箱+密码登录</span>
        </div>

        <div class="form-group">
          <label class="form-label">设置密码</label>
          <input
            v-model="setPasswordForm.password"
            type="password"
            class="form-input"
            :class="{ error: setPasswordErrors.password }"
            placeholder="请输入密码（6-32个字符）"
            :disabled="setPasswordLoading"
            autocomplete="new-password"
          />
          <p v-if="setPasswordErrors.password" class="error-message">
            {{ setPasswordErrors.password }}
          </p>
        </div>

        <div class="form-group">
          <label class="form-label">确认密码</label>
          <input
            v-model="setPasswordForm.confirm_password"
            type="password"
            class="form-input"
            :class="{ error: setPasswordErrors.confirm_password }"
            placeholder="请再次输入密码"
            :disabled="setPasswordLoading"
            autocomplete="new-password"
          />
          <p v-if="setPasswordErrors.confirm_password" class="error-message">
            {{ setPasswordErrors.confirm_password }}
          </p>
        </div>
      </form>
    </UiBaseDialog>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

// 按钮基础样式
@mixin btn-base {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

// 页面容器
#page {
  @extend .cardHover;
  align-self: flex-start;
  padding: 40px;

  .page-title {
    margin: 0 0 30px;
    font-weight: bold;
    font-size: 2rem;
  }
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

// 信息卡片
.info-card {
  padding: 25px;
  background: var(--flec-card-bg);
  border-radius: 8px;

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 10px;
    border-bottom: 2px solid var(--flec-border);

    .card-title {
      margin: 0;
      padding-bottom: 0;
      border-bottom: none;
    }
  }

  .card-title {
    margin: 0 0 20px;
    padding-bottom: 10px;
    border-bottom: 2px solid var(--flec-border);
    font-weight: 600;
    color: var(--font-color);
  }

  .info-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .info-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    min-height: 48px;
    padding: 8px 0;
    border-bottom: 1px solid #8080800d;

    &:last-child {
      border-bottom: none;
    }

    .label {
      font-weight: 500;
      color: var(--theme-meta-color);
    }

    .value {
      display: flex;
      align-items: center;
      gap: 10px;
      color: var(--font-color);
      text-align: right;

      .avatar-preview {
        width: 50px;
        height: 50px;
        border-radius: 50%;
        object-fit: cover;
        border: 2px solid var(--flec-border);
      }

      .badge-text {
        background: var(--flec-card-bg);
        border: 1px solid var(--theme-color);
        color: var(--theme-color);
        padding: 2px 8px;
        border-radius: 4px;
        font-size: 0.9em;
      }

      .btn-icon {
        padding: 6px;
        border-radius: 6px;
        color: var(--theme-meta-color);
        background: transparent;
        border: none;
        cursor: pointer;
        transition: all 0.2s;
        display: flex;
        align-items: center;

        &:hover {
          background: var(--flec-border);
          color: var(--theme-color);
        }
      }

      .empty {
        color: var(--theme-meta-color);
        font-size: 0.9rem;
      }
    }

    .action-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      width: 100%;
      gap: 20px;

      .action-info {
        display: flex;
        flex-direction: column;
        gap: 4px;

        .action-title {
          font-weight: 500;
          color: var(--font-color);
          font-size: 0.95rem;
        }

        .action-desc {
          color: var(--theme-meta-color);
          font-size: 0.85rem;
          line-height: 1.4;
        }
      }
    }
  }

  .login-methods {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;

    .method-icon {
      width: 36px;
      height: 36px;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 8px;
      background: #8080800a;
      cursor: default;
      transition: all 0.2s;

      i {
        font-size: 20px;
        color: var(--theme-meta-color);
      }

      &.enabled i {
        color: var(--theme-color);
      }

      &.loading {
        cursor: wait;
        pointer-events: none;

        i {
          color: var(--theme-color);
        }
      }

      &.clickable {
        cursor: pointer;

        &:hover {
          background: #80808014;
        }
      }

      &.enabled.clickable:hover {
        i {
          color: var(--theme-color);
        }
      }
    }

    .spin {
      animation: spin 1s linear infinite;
    }

    @keyframes spin {
      from {
        transform: rotate(0deg);
      }

      to {
        transform: rotate(360deg);
      }
    }
  }
}

// 按钮
.btn-primary {
  @include btn-base;
  border: none;
  background: var(--theme-color);
  color: white;

  &:hover {
    opacity: 0.9;
  }
}

.btn-secondary {
  @include btn-base;
  border: 1px solid var(--flec-border);
  background: var(--flec-card-bg);
  color: var(--font-color);

  &:hover {
    background: #80808014;
  }
}

.btn-danger {
  @include btn-base;
  border: 1px solid #e57373;
  background: transparent;
  color: #e57373;

  &:hover {
    background: #e57373;
    color: white;
  }
}

.btn-link {
  color: var(--theme-color);
  font-size: 1rem;
  transition: opacity 0.2s;

  &:hover {
    opacity: 0.8;
    text-decoration: underline;
  }
}

// ===== 对话框表单样式 =====
.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-size: 0.95rem;
  font-weight: 500;
  color: var(--font-color);
}

.form-input {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid var(--flec-border);
  border-radius: 8px;
  background: var(--flec-card-bg);
  color: var(--font-color);
  font-size: 0.95rem;
  transition: all 0.2s;

  &:focus {
    outline: none;
    border-color: var(--theme-color);
    box-shadow: 0 0 0 3px #49b1f526;
  }

  &.error {
    border-color: #e57373;

    &:focus {
      box-shadow: 0 0 0 3px #e5737326;
    }
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  &::placeholder {
    color: var(--theme-meta-color);
  }
}

.error-message {
  margin: 0;
  font-size: 0.85rem;
  color: #e57373;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 15px;

  img {
    width: 70px;
    height: 70px;
    border-radius: 50%;
    object-fit: cover;
    border: 2px solid var(--flec-border);
  }

  button {
    padding: 8px 16px;
    border: 1px solid var(--flec-border);
    border-radius: 8px;
    background: var(--flec-card-bg);
    color: var(--font-color);
    cursor: pointer;
    transition: all 0.2s;

    &:hover:not(:disabled) {
      border-color: var(--theme-color);
      color: var(--theme-color);
    }

    &:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
  }
}

.tip {
  display: flex;
  gap: 8px;
  padding: 12px;
  background: #49b1f50d;
  border-radius: 8px;
  font-size: 0.9rem;
  color: var(--theme-meta-color);
  line-height: 1.5;

  i {
    font-size: 16px;
    color: var(--theme-color);
    flex-shrink: 0;
    margin-top: 2px;
  }
}

.warning {
  padding: 16px;
  background: #8080800d;
  border: 1px solid var(--flec-border);
  border-radius: 8px;

  .warning-title {
    margin: 0 0 12px;
    font-size: 1rem;
    font-weight: 600;
    color: var(--font-color);
  }

  p {
    margin: 0 0 8px;
    font-size: 0.9rem;
    color: var(--font-color);
  }

  ul {
    margin: 8px 0 0;
    padding-left: 20px;
    font-size: 0.9rem;
    color: var(--theme-meta-color);
    line-height: 1.8;

    li {
      margin: 4px 0;
    }
  }
}

.checkbox-label {
  display: flex;
  gap: 10px;
  cursor: pointer;
  padding: 12px;
  border-radius: 8px;
  transition: background 0.2s;

  &:hover {
    background: #8080800d;
  }

  input[type='checkbox'] {
    margin-top: 2px;
    width: 18px;
    height: 18px;
    cursor: pointer;
    flex-shrink: 0;
    accent-color: var(--theme-color);
  }

  span {
    font-size: 0.95rem;
    color: var(--font-color);
    line-height: 1.5;
  }
}

.dialog-hint {
  margin: 0;
  font-size: 0.9rem;
  color: var(--theme-meta-color);
  line-height: 1.5;
}

// 响应式
@media (max-width: 768px) {
  #page {
    padding: 20px;

    .page-title {
      font-size: 1.5rem;
      margin-bottom: 20px;
    }
  }

  .info-card {
    padding: 20px;

    .card-header {
      flex-direction: column;
      align-items: stretch;
      gap: 12px;

      .card-title {
        padding-bottom: 0;
        border-bottom: none;
      }
    }

    .info-item {
      flex-direction: column;
      gap: 8px;
      align-items: flex-start;
      min-height: auto;

      .value {
        width: 100%;
        text-align: left;
        flex-wrap: wrap;
      }

      .action-item {
        flex-direction: column;
        align-items: stretch;
        gap: 12px;

        .action-info .action-title {
          font-size: 1rem;
        }

        .action-info .action-desc {
          font-size: 0.8rem;
        }

        button {
          width: 100%;
          justify-content: center;
        }
      }
    }

    .login-methods {
      justify-content: flex-start;
      gap: 10px;

      .method-icon {
        width: 32px;
        height: 32px;

        i {
          font-size: 18px;
        }
      }
    }
  }

  .btn-primary,
  .btn-secondary {
    width: 100%;
    justify-content: center;
  }

  .btn-link {
    padding: 6px 12px;
  }
}
</style>
