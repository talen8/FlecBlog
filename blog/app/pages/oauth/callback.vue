<script setup lang="ts">
const router = useRouter();
const route = useRoute();
const { fetchUserInfo } = useUser();

onMounted(async () => {
  // 只在客户端执行
  if (!import.meta.client) return;

  const token = route.query.token as string;
  const redirect = route.query.redirect as string;

  if (token) {
    // 保存 Access Token
    const { setAccessToken } = await import('@/utils/auth');
    setAccessToken(token);

    // 获取用户信息
    await fetchUserInfo();

    // 跳转回原页面，如果有 redirect 参数则使用，否则回首页
    if (redirect) {
      router.push(decodeURIComponent(redirect));
    } else {
      router.push('/');
    }
  } else {
    // 失败处理：跳转首页
    router.push('/');
  }
});
</script>

<template>
  <div class="callback-page">
    <div class="loading-content">
      <i class="ri-loader-4-line spin" />
      <p>正在登录中，请稍候...</p>
    </div>
  </div>
</template>

<style scoped>
.callback-page {
  height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-content {
  text-align: center;
  color: var(--text-secondary);
}

.spin {
  font-size: 2rem;
  animation: spin 1s linear infinite;
  display: block;
  margin: 0 auto 1rem;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
