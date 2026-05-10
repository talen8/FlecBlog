<script lang="ts" setup>
import { login, register, forgotPassword, resetPassword } from '@/composables/useApi';

const showLoginModal = useState('showLoginModal', () => false);
const { fetchUserInfo } = useUser();
const { oauthConfig } = useSysConfig();
const toast = useToast();

const mode = ref<'login' | 'register' | 'forgot' | 'reset'>('login');
const loading = ref(false);

const loginForm = reactive({ email: '', password: '' });
const registerForm = reactive({ email: '', nickname: '', password: '', website: '' });
const forgotForm = reactive({ email: '' });
const resetForm = reactive({ email: '', code: '', password: '', confirmPassword: '' });

const githubEnabled = computed(() => oauthConfig.value['oauth.github.enabled'] === 'true');
const googleEnabled = computed(() => oauthConfig.value['oauth.google.enabled'] === 'true');
const qqEnabled = computed(() => oauthConfig.value['oauth.qq.enabled'] === 'true');
const microsoftEnabled = computed(() => oauthConfig.value['oauth.microsoft.enabled'] === 'true');

const hasOAuth = computed(
  () => githubEnabled.value || googleEnabled.value || qqEnabled.value || microsoftEnabled.value
);

const dialogTitle = computed(() => {
  switch (mode.value) {
    case 'login':
      return '登录';
    case 'register':
      return '注册';
    case 'forgot':
      return '忘记密码';
    case 'reset':
      return '重置密码';
    default:
      return '登录';
  }
});

const handleLogin = async () => {
  if (!loginForm.email || !loginForm.password) {
    toast.error('请输入邮箱和密码');
    return;
  }
  loading.value = true;
  try {
    const { access_token } = await login(loginForm);
    setAccessToken(access_token);
    await fetchUserInfo();
    showLoginModal.value = false;
    toast.success('登录成功');
    loginForm.email = '';
    loginForm.password = '';
  } catch (error) {
    console.error('登录失败:', error);
    toast.error('登录失败');
  } finally {
    loading.value = false;
  }
};

const handleRegister = async () => {
  if (!registerForm.email || !registerForm.nickname || !registerForm.password) {
    toast.error('请填写完整信息');
    return;
  }
  loading.value = true;
  try {
    const { access_token } = await register(registerForm);
    setAccessToken(access_token);
    await fetchUserInfo();
    showLoginModal.value = false;
    toast.success('注册成功');
  } catch (error) {
    console.error('注册失败:', error);
    toast.error('注册失败');
  } finally {
    loading.value = false;
  }
};

const handleForgot = async () => {
  if (!forgotForm.email || !forgotForm.email.includes('@')) {
    toast.error('请输入有效的邮箱地址');
    return;
  }
  loading.value = true;
  try {
    await forgotPassword({ email: forgotForm.email });
    resetForm.email = forgotForm.email;
    mode.value = 'reset';
    toast.success('验证码已发送');
  } catch (error) {
    console.error('发送失败:', error);
    toast.error('发送失败');
  } finally {
    loading.value = false;
  }
};

const handleReset = async () => {
  if (!resetForm.code) {
    toast.error('请输入验证码');
    return;
  }
  if (!resetForm.password || resetForm.password.length < 6) {
    toast.error('密码至少6位');
    return;
  }
  if (resetForm.password !== resetForm.confirmPassword) {
    toast.error('两次密码不一致');
    return;
  }
  loading.value = true;
  try {
    await resetPassword({
      email: resetForm.email,
      code: resetForm.code,
      password: resetForm.password,
    });
    toast.success('密码重置成功');
    mode.value = 'login';
    resetForm.code = '';
    resetForm.password = '';
    resetForm.confirmPassword = '';
  } catch (error) {
    console.error('重置失败:', error);
    toast.error('重置失败');
  } finally {
    loading.value = false;
  }
};

const handleOAuth = (provider: string) => {
  const config = useRuntimeConfig();
  const currentPath = encodeURIComponent(window.location.pathname);
  window.location.href = `${config.public.apiUrl}/auth/${provider}?redirect=${currentPath}`;
};

watch(showLoginModal, val => {
  if (val) mode.value = 'login';
});
</script>

<template>
  <UiBaseDialog
    v-model="showLoginModal"
    :title="dialogTitle"
    :show-footer="false"
    :close-on-click-outside="true"
  >
    <template v-if="mode === 'login' || mode === 'register'">
      <div v-if="hasOAuth" class="oauth-section">
        <div class="oauth-buttons">
          <button v-if="githubEnabled" class="oauth-btn" @click="handleOAuth('github')">
            <i class="ri-github-fill" /> GitHub
          </button>
          <button v-if="googleEnabled" class="oauth-btn" @click="handleOAuth('google')">
            <i class="ri-google-fill" /> Google
          </button>
          <button v-if="qqEnabled" class="oauth-btn" @click="handleOAuth('qq')">
            <i class="ri-qq-line" /> QQ
          </button>
          <button v-if="microsoftEnabled" class="oauth-btn" @click="handleOAuth('microsoft')">
            <i class="ri-microsoft-fill" /> Microsoft
          </button>
        </div>
        <div class="divider"><span>或</span></div>
      </div>
    </template>

    <form v-if="mode === 'login'" class="auth-form" @submit.prevent="handleLogin">
      <div class="form-group">
        <label>邮箱</label>
        <input
          v-model="loginForm.email"
          type="email"
          placeholder="请输入邮箱"
          :disabled="loading"
        />
      </div>
      <div class="form-group">
        <label>密码</label>
        <input
          v-model="loginForm.password"
          type="password"
          placeholder="请输入密码"
          :disabled="loading"
        />
      </div>
      <button class="submit-btn" type="submit" :disabled="loading">
        {{ loading ? '登录中...' : '登录' }}
      </button>
      <div class="form-links">
        <button type="button" class="link" @click="mode = 'forgot'">忘记密码？</button>
        <button type="button" class="link" @click="mode = 'register'">立即注册</button>
      </div>
    </form>

    <form v-else-if="mode === 'register'" class="auth-form" @submit.prevent="handleRegister">
      <div class="form-group">
        <label>邮箱</label>
        <input
          v-model="registerForm.email"
          type="email"
          placeholder="请输入邮箱"
          :disabled="loading"
        />
      </div>
      <div class="form-group">
        <label>昵称</label>
        <input
          v-model="registerForm.nickname"
          type="text"
          placeholder="请输入昵称"
          :disabled="loading"
        />
      </div>
      <div class="form-group">
        <label>密码</label>
        <input
          v-model="registerForm.password"
          type="password"
          placeholder="请输入密码（6-32位）"
          :disabled="loading"
        />
      </div>
      <div class="form-group">
        <label>网站（可选）</label>
        <input
          v-model="registerForm.website"
          type="url"
          placeholder="https://"
          :disabled="loading"
        />
      </div>
      <button class="submit-btn" type="submit" :disabled="loading">
        {{ loading ? '注册中...' : '注册' }}
      </button>
      <p class="switch-mode">
        已有账号？<button type="button" @click="mode = 'login'">立即登录</button>
      </p>
    </form>

    <form v-else-if="mode === 'forgot'" class="auth-form" @submit.prevent="handleForgot">
      <p class="hint">请输入注册时使用的邮箱地址，我们将发送验证码到您的邮箱。</p>
      <div class="form-group">
        <label>邮箱地址</label>
        <input
          v-model="forgotForm.email"
          type="email"
          placeholder="请输入邮箱"
          :disabled="loading"
        />
      </div>
      <button class="submit-btn" type="submit" :disabled="loading">
        {{ loading ? '发送中...' : '发送验证码' }}
      </button>
      <p class="switch-mode">
        <button type="button" @click="mode = 'login'">返回登录</button>
      </p>
    </form>

    <form v-else-if="mode === 'reset'" class="auth-form" @submit.prevent="handleReset">
      <p class="hint">验证码已发送到 {{ resetForm.email }}</p>
      <div class="form-group">
        <label>验证码</label>
        <input
          v-model="resetForm.code"
          type="text"
          placeholder="请输入验证码"
          :disabled="loading"
        />
      </div>
      <div class="form-group">
        <label>新密码</label>
        <input
          v-model="resetForm.password"
          type="password"
          placeholder="新密码（6-32位）"
          :disabled="loading"
        />
      </div>
      <div class="form-group">
        <label>确认新密码</label>
        <input
          v-model="resetForm.confirmPassword"
          type="password"
          placeholder="确认新密码"
          :disabled="loading"
        />
      </div>
      <button class="submit-btn" type="submit" :disabled="loading">
        {{ loading ? '重置中...' : '重置密码' }}
      </button>
      <p class="switch-mode">
        <button type="button" @click="mode = 'forgot'">重新发送验证码</button>
      </p>
    </form>
  </UiBaseDialog>
</template>

<style lang="scss" scoped>
.oauth-section {
  margin-bottom: 20px;
}
.oauth-buttons {
  display: flex;
  gap: 12px;
}

.oauth-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 10px;
  border: 1px solid var(--flec-border);
  background: var(--flec-card-bg);
  color: var(--font-color);
  font-size: 0.9rem;
  cursor: pointer;
}

.divider {
  display: flex;
  align-items: center;
  margin: 20px 0;
  gap: 16px;
  &::before,
  &::after {
    content: '';
    flex: 1;
    height: 1px;
    background: var(--flec-border);
  }
  span {
    color: var(--theme-meta-color);
    font-size: 0.85rem;
  }
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.hint {
  margin: 0;
  color: var(--theme-meta-color);
  font-size: 0.9rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  label {
    font-weight: 500;
    font-size: 0.9rem;
  }
  input {
    padding: 10px 12px;
    border: 1px solid var(--flec-border);
    background: var(--flec-card-bg);
    color: var(--font-color);
    font-size: 0.95rem;
    &:focus {
      outline: none;
      border-color: var(--theme-color);
    }
  }
}

.submit-btn {
  padding: 10px;
  background: var(--theme-color);
  color: #fff;
  border: none;
  font-size: 1rem;
  cursor: pointer;
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.form-links {
  display: flex;
  justify-content: space-between;
  font-size: 0.9rem;
  .link {
    color: var(--theme-color);
    background: none;
    border: none;
    cursor: pointer;
    font-size: 0.9rem;
    text-decoration: none;
  }
}

.switch-mode {
  text-align: center;
  font-size: 0.9rem;
  color: var(--theme-meta-color);
  margin: 0;
  button {
    color: var(--theme-color);
    background: none;
    border: none;
    cursor: pointer;
    font-size: 0.9rem;
  }
}
</style>
