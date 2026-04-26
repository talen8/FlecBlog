<template>
  <div class="admin-login">
    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <h1>后台管理系统</h1>
          <p>Admin Management System</p>
        </div>
        <div class="login-form">
          <div class="form-item">
            <i class="ri-user-line"></i>
            <input
              type="text"
              v-model="formState.email"
              placeholder="请输入邮箱"
              @keyup.enter="handleLogin"
            />
          </div>
          <div class="form-item">
            <i class="ri-lock-line"></i>
            <input
              type="password"
              v-model="formState.password"
              placeholder="请输入密码"
              @keyup.enter="handleLogin"
            />
          </div>
          <button class="submit-btn" :disabled="loading" @click="handleLogin">
            <span v-if="!loading">登 录</span>
            <span v-else class="loading-text">
              <i class="ri-loader-4-line"></i>
              登录中...
            </span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { login } from '@/api/user';
import { fetchUserInfo, setAccessToken, clearAuthState } from '@/utils/auth';
import type { LoginParams } from '@/types/user';

const router = useRouter();
const loading = ref(false);
const formState = reactive<LoginParams>({
  email: '',
  password: '',
});

const handleLogin = async () => {
  if (!formState.email || !formState.password) {
    ElMessage.warning('请输入邮箱和密码');
    return;
  }

  loading.value = true;
  try {
    const { access_token } = await login(formState);
    setAccessToken(access_token);
    await fetchUserInfo();
    ElMessage.success('登录成功');
    router.push('/');
  } catch (error) {
    clearAuthState();
    ElMessage.error(error instanceof Error ? error.message : '登录失败');
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped lang="scss">
.admin-login {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: #f8fafc;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: -50%;
    right: -20%;
    width: 600px;
    height: 600px;
    background: #e2e8f0;
    border-radius: 50%;
    opacity: 0.5;
  }

  &::after {
    content: '';
    position: absolute;
    bottom: -30%;
    left: -10%;
    width: 400px;
    height: 400px;
    background: #e2e8f0;
    border-radius: 50%;
    opacity: 0.5;
  }
}

.login-container {
  width: 100%;
  max-width: 400px;
  padding: 20px;
  position: relative;
  z-index: 1;
}

.login-card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06);
  padding: 48px 40px;
  border: 1px solid #e2e8f0;
}

.login-header {
  text-align: center;
  margin-bottom: 36px;

  h1 {
    font-size: 24px;
    color: #1e293b;
    margin-bottom: 8px;
    font-weight: 600;
  }

  p {
    color: #64748b;
    font-size: 14px;
  }
}

.login-form {
  .form-item {
    position: relative;
    margin-bottom: 20px;

    i {
      position: absolute;
      left: 14px;
      top: 50%;
      transform: translateY(-50%);
      color: #94a3b8;
      font-size: 18px;
    }

    input {
      width: 100%;
      height: 48px;
      line-height: 48px;
      padding: 0 16px 0 44px;
      border: 1px solid #e2e8f0;
      border-radius: 10px;
      color: #1e293b;
      font-size: 15px;
      background: #f8fafc;
      transition: all 0.2s ease;

      &:focus {
        border-color: #1e293b;
        background: #fff;
        box-shadow: 0 0 0 3px rgba(30, 41, 59, 0.1);
      }

      &::placeholder {
        color: #94a3b8;
      }
    }
  }

  .submit-btn {
    width: 100%;
    height: 48px;
    background: #1e293b;
    border: none;
    border-radius: 10px;
    color: #fff;
    font-size: 16px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    margin-top: 8px;

    &:hover:not(:disabled) {
      background: #334155;
    }

    &:active:not(:disabled) {
      transform: scale(0.98);
    }

    &:disabled {
      opacity: 0.7;
      cursor: not-allowed;
    }

    .loading-text {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 8px;

      i {
        animation: spin 1s linear infinite;
        position: static;
        transform: none;
        color: #fff;
        font-size: 16px;
      }
    }
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }

  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 768px) {
  .admin-login {
    &::before,
    &::after {
      display: none;
    }
  }

  .login-container {
    max-width: 100%;
    padding: 16px;
  }

  .login-card {
    padding: 36px 24px;
    border-radius: 12px;
  }

  .login-header {
    margin-bottom: 28px;

    h1 {
      font-size: 22px;
    }

    p {
      font-size: 13px;
    }
  }

  .login-form {
    .form-item {
      margin-bottom: 16px;

      input {
        height: 44px;
        font-size: 14px;
      }
    }

    .submit-btn {
      height: 44px;
      font-size: 15px;
    }
  }
}
</style>
