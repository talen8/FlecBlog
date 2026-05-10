<template>
  <div id="page">
    <div class="callback-content">
      <p v-if="loading">正在处理登录...</p>
      <p v-else-if="error">{{ error }}</p>
      <p v-else>登录成功，正在跳转...</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { setAccessToken } from '@/utils/auth';

const route = useRoute();
const router = useRouter();
const { fetchUserInfo } = useUser();
const toast = useToast();

const loading = ref(true);
const error = ref('');

onMounted(async () => {
  const token = route.query.token as string;
  const redirect = route.query.redirect as string;

  if (!token) {
    error.value = '登录参数错误';
    loading.value = false;
    return;
  }

  try {
    setAccessToken(token);
    await fetchUserInfo();

    toast.success('登录成功');

    if (redirect) {
      router.replace(decodeURIComponent(redirect));
    } else {
      router.replace('/');
    }
  } catch (err) {
    console.error('OAuth 登录失败:', err);
    error.value = '登录失败，请重试';
    loading.value = false;
  }
});

useSeoMeta({
  title: '登录中...',
});
</script>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border-radius: 12px;
  border: 1px solid var(--flec-border);
  padding: 40px;
  max-width: 400px;
  margin: 0 auto;
  text-align: center;
}

.callback-content {
  p {
    margin: 0;
    color: var(--font-color);
    font-size: 1rem;
  }
}
</style>
