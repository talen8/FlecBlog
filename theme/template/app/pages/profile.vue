<script lang="ts" setup>
import {
  updateUserProfile,
  changePassword,
  deactivateAccount,
  unbindOAuth,
} from '@/composables/useApi';
import { logout } from '@/utils/auth';

const isLoggedIn = useAuth();
const router = useRouter();
const { userInfo, fetchUserInfo, clearUserInfo } = useUser();
const showLoginModal = useState('showLoginModal', () => false);
const toast = useToast();

const form = reactive({
  nickname: '',
  website: '',
});
const saving = ref(false);

const showPasswordForm = ref(false);
const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: '',
});
const changingPassword = ref(false);

const showDeactivateConfirm = ref(false);
const deactivatePassword = ref('');
const deactivating = ref(false);

watch(
  userInfo,
  newVal => {
    if (newVal) {
      form.nickname = newVal.nickname || '';
      form.website = newVal.website || '';
    }
  },
  { immediate: true }
);

onMounted(() => {
  if (isLoggedIn.value && !userInfo.value) {
    fetchUserInfo();
  }
});

const handleSave = async () => {
  if (!form.nickname.trim()) {
    toast.error('请输入昵称');
    return;
  }
  saving.value = true;
  try {
    await updateUserProfile({
      nickname: form.nickname.trim(),
      website: form.website.trim() || undefined,
    });
    await fetchUserInfo();
    toast.success('保存成功');
  } catch (error) {
    console.error('保存失败:', error);
    toast.error('保存失败');
  } finally {
    saving.value = false;
  }
};

const handleChangePassword = async () => {
  if (!passwordForm.old_password || !passwordForm.new_password) {
    toast.error('请填写密码');
    return;
  }
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    toast.error('两次密码不一致');
    return;
  }
  if (passwordForm.new_password.length < 6) {
    toast.error('新密码至少6位');
    return;
  }
  changingPassword.value = true;
  try {
    await changePassword({
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password,
    });
    toast.success('密码修改成功');
    showPasswordForm.value = false;
    passwordForm.old_password = '';
    passwordForm.new_password = '';
    passwordForm.confirm_password = '';
  } catch (error) {
    console.error('修改密码失败:', error);
    toast.error('修改密码失败');
  } finally {
    changingPassword.value = false;
  }
};

const handleDeactivate = async () => {
  if (!deactivatePassword.value) {
    toast.error('请输入密码确认');
    return;
  }
  deactivating.value = true;
  try {
    await deactivateAccount({ password: deactivatePassword.value });
    toast.success('账号已注销');
    clearUserInfo();
    logout();
    router.push('/');
  } catch (error) {
    console.error('注销失败:', error);
    toast.error('注销失败');
  } finally {
    deactivating.value = false;
  }
};

const handleUnbindOAuth = async (provider: string) => {
  if (!confirm(`确定要解绑 ${provider} 吗？`)) return;
  try {
    await unbindOAuth(provider);
    await fetchUserInfo();
    toast.success(`已解绑 ${provider}`);
  } catch (error) {
    console.error('解绑失败:', error);
    toast.error('解绑失败');
  }
};

useSeoMeta({
  title: '个人中心',
});
</script>

<template>
  <div v-if="userInfo" id="page">
    <h1 class="page-title">个人中心</h1>

    <div class="profile-header">
      <img :src="getAvatarUrl(userInfo)" alt="avatar" class="avatar" />
      <div class="profile-info">
        <h2>{{ userInfo.nickname }}</h2>
        <p>{{ userInfo.email }}</p>
        <span class="role-badge">{{ userInfo.role }}</span>
      </div>
    </div>

    <div class="section">
      <h2 class="section-title">基本信息</h2>
      <div class="form-group">
        <label>昵称</label>
        <input v-model="form.nickname" type="text" placeholder="昵称" />
      </div>
      <div class="form-group">
        <label>网站</label>
        <input v-model="form.website" type="url" placeholder="网站" />
      </div>
      <button class="save-btn" :disabled="saving" @click="handleSave">
        {{ saving ? '保存中...' : '保存' }}
      </button>
    </div>

    <div v-if="userInfo.linked_oauths?.length" class="section">
      <h2 class="section-title">已绑定账号</h2>
      <div class="oauth-list">
        <div v-for="provider in userInfo.linked_oauths" :key="provider" class="oauth-item">
          <span
            ><i
              :class="
                provider === 'github'
                  ? 'ri-github-fill'
                  : provider === 'google'
                    ? 'ri-google-fill'
                    : 'ri-qq-line'
              "
            />
            {{ provider }}</span
          >
          <button class="unbind-btn" @click="handleUnbindOAuth(provider)">解绑</button>
        </div>
      </div>
    </div>

    <div class="section">
      <h2 class="section-title">安全设置</h2>
      <button class="action-btn" @click="showPasswordForm = !showPasswordForm">
        {{ showPasswordForm ? '取消' : '修改密码' }}
      </button>
      <div v-if="showPasswordForm" class="password-form">
        <div class="form-group">
          <label>当前密码</label>
          <input v-model="passwordForm.old_password" type="password" placeholder="当前密码" />
        </div>
        <div class="form-group">
          <label>新密码</label>
          <input
            v-model="passwordForm.new_password"
            type="password"
            placeholder="新密码（6-32位）"
          />
        </div>
        <div class="form-group">
          <label>确认新密码</label>
          <input v-model="passwordForm.confirm_password" type="password" placeholder="确认新密码" />
        </div>
        <button class="save-btn" :disabled="changingPassword" @click="handleChangePassword">
          {{ changingPassword ? '修改中...' : '确认修改' }}
        </button>
      </div>
    </div>

    <div class="section danger-zone">
      <h2 class="section-title">危险操作</h2>
      <button class="danger-btn" @click="showDeactivateConfirm = !showDeactivateConfirm">
        注销账号
      </button>
      <div v-if="showDeactivateConfirm" class="deactivate-form">
        <p class="warning-text">注销后账号数据将无法恢复，请输入密码确认。</p>
        <div class="form-group">
          <label>密码</label>
          <input v-model="deactivatePassword" type="password" placeholder="请输入密码确认注销" />
        </div>
        <button class="danger-btn" :disabled="deactivating" @click="handleDeactivate">
          {{ deactivating ? '注销中...' : '确认注销' }}
        </button>
      </div>
    </div>
  </div>

  <div v-else class="login-prompt">
    <p>请先登录</p>
    <button class="login-btn" @click="showLoginModal = true">登录</button>
  </div>
</template>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border-radius: 12px;
  border: 1px solid var(--flec-border);
  padding: 40px;
  max-width: 600px;
  margin: 0 auto;
}

.page-title {
  margin: 0 0 30px;
  font-size: 2rem;
  font-weight: 700;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  object-fit: cover;
}

.profile-info {
  h2 {
    margin: 0 0 8px;
    font-size: 1.3rem;
  }
  p {
    margin: 0 0 8px;
    color: var(--theme-meta-color);
  }
}

.role-badge {
  display: inline-block;
  padding: 2px 8px;
  background: var(--theme-color);
  color: #fff;
  font-size: 0.75rem;
  border-radius: 4px;
}

.section {
  margin-bottom: 30px;
  padding-top: 20px;
  border-top: 1px solid var(--flec-border);
}

.section-title {
  font-size: 1.1rem;
  font-weight: 600;
  margin: 0 0 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 16px;

  label {
    font-weight: 500;
    font-size: 0.9rem;
  }

  input {
    padding: 12px;
    border: 1px solid var(--flec-border);
    border-radius: 8px;
    background: var(--flec-card-bg);
    color: var(--font-color);
    font-size: 0.95rem;

    &:focus {
      outline: none;
      border-color: var(--theme-color);
    }
  }
}

.save-btn {
  padding: 12px 24px;
  background: var(--theme-color);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
  transition: opacity 0.2s;

  &:hover:not(:disabled) {
    opacity: 0.9;
  }
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.action-btn {
  padding: 10px 20px;
  border: 1px solid var(--flec-border);
  border-radius: 8px;
  background: var(--flec-card-bg);
  color: var(--font-color);
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    border-color: var(--theme-color);
    color: var(--theme-color);
  }
}

.oauth-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.oauth-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 16px;
  background: var(--flec-heavy-bg);
  border-radius: 8px;

  i {
    margin-right: 8px;
  }
}

.unbind-btn {
  padding: 4px 12px;
  border: 1px solid var(--flec-border);
  border-radius: 6px;
  background: none;
  color: var(--theme-meta-color);
  font-size: 0.85rem;
  cursor: pointer;

  &:hover {
    border-color: #ff4d4f;
    color: #ff4d4f;
  }
}

.danger-zone {
  .danger-btn {
    padding: 10px 20px;
    background: none;
    border: 1px solid #ff4d4f;
    border-radius: 8px;
    color: #ff4d4f;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      background: #ff4d4f;
      color: #fff;
    }
    &:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
  }
}

.warning-text {
  color: #ff4d4f;
  font-size: 0.9rem;
  margin: 0 0 16px;
}

.login-prompt {
  text-align: center;
  padding: 60px 20px;

  p {
    margin-bottom: 20px;
    color: var(--theme-meta-color);
  }
}

.login-btn {
  padding: 12px 32px;
  background: var(--theme-color);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;
}

@media screen and (max-width: 768px) {
  #page {
    padding: 20px;
  }
  .page-title {
    font-size: 1.5rem;
  }
}
</style>
