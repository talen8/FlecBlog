<script lang="ts" setup>
import { subscribe } from '@/composables/useApi';

const config = useRuntimeConfig();
const toast = useToast();

const email = ref('');
const submitting = ref(false);

const handleSubscribe = async () => {
  if (!email.value || !email.value.includes('@')) {
    toast.error('请输入有效的邮箱地址');
    return;
  }

  submitting.value = true;
  try {
    await subscribe(email.value);
    toast.success('订阅成功！');
    email.value = '';
  } catch (error) {
    console.error('订阅失败:', error);
    toast.error('订阅失败，请稍后重试');
  } finally {
    submitting.value = false;
  }
};

useSeoMeta({
  title: '订阅',
  description: '订阅本站获取最新文章更新',
});
</script>

<template>
  <div id="page">
    <h1 class="page-title">订阅</h1>
    <div class="subscribe-content">
      <p>订阅本站，获取最新文章更新通知。</p>
      <div class="subscribe-form">
        <input v-model="email" type="email" placeholder="请输入邮箱地址" :disabled="submitting" />
        <button class="subscribe-btn" :disabled="submitting" @click="handleSubscribe">
          {{ submitting ? '提交中...' : '订阅' }}
        </button>
      </div>
      <div class="subscribe-links">
        <a :href="`${config.public.apiUrl}/rss.xml`" target="_blank" class="subscribe-link">
          <i class="ri-rss-line" /> RSS 订阅
        </a>
        <a :href="`${config.public.apiUrl}/atom.xml`" target="_blank" class="subscribe-link">
          <i class="ri-rss-line" /> Atom 订阅
        </a>
      </div>
    </div>
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

.subscribe-content {
  text-align: center;

  p {
    margin-bottom: 24px;
    color: var(--theme-meta-color);
  }
}

.subscribe-form {
  display: flex;
  gap: 12px;
  margin-bottom: 30px;

  input {
    flex: 1;
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

.subscribe-btn {
  padding: 12px 24px;
  background: var(--theme-color);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
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

.subscribe-links {
  display: flex;
  justify-content: center;
  gap: 20px;
}

.subscribe-link {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  border: 1px solid var(--flec-border);
  border-radius: 8px;
  color: var(--font-color);
  text-decoration: none;
  font-size: 0.95rem;
  transition: all 0.2s;

  &:hover {
    border-color: var(--theme-color);
    color: var(--theme-color);
  }
}

@media screen and (max-width: 768px) {
  #page {
    padding: 20px;
  }

  .page-title {
    font-size: 1.5rem;
  }

  .subscribe-form {
    flex-direction: column;
  }

  .subscribe-links {
    flex-direction: column;
  }
}
</style>
