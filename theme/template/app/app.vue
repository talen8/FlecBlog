<script lang="ts" setup>
import { getMenus, getSettingGroup } from '@/composables/useApi';

const { menus } = useMenus();
const { basicConfig, blogConfig, oauthConfig } = useSysConfig();
const { userInfo, fetchUserInfo } = useUser();
const isLoggedIn = useAuth();

const { data: menusData } = await useAsyncData('menus', () => getMenus());
if (menusData.value) {
  menus.value = menusData.value;
}

const { data: basicData } = await useAsyncData('settings-basic', () => getSettingGroup('basic'));
if (basicData.value) {
  basicConfig.value = basicData.value;
}

const { data: blogData } = await useAsyncData('settings-blog', () => getSettingGroup('blog'));
if (blogData.value) {
  blogConfig.value = blogData.value;
}

const { data: oauthData } = await useAsyncData('settings-oauth', () => getSettingGroup('oauth'));
if (oauthData.value) {
  oauthConfig.value = oauthData.value;
}

onMounted(() => {
  if (isLoggedIn.value && !userInfo.value) {
    fetchUserInfo();
  }
});

const favicon = computed(() => blogConfig.value['blog.favicon'] || '/favicon.ico');
const keywords = computed(() => blogConfig.value['blog.keywords'] || '');
const customFont = computed(() => blogConfig.value['blog.font'] || '');
const customHead = computed(() => blogConfig.value['blog.custom_head'] || '');
const customBody = computed(() => blogConfig.value['blog.custom_body'] || '');

useHead({
  titleTemplate: title =>
    title
      ? `${title} - ${blogConfig.value['blog.title'] || 'Blog'}`
      : blogConfig.value['blog.title'] || 'Blog',
  meta: [{ name: 'keywords', content: keywords }],
  link: [
    { rel: 'icon', type: 'image/x-icon', href: favicon },
    {
      rel: 'stylesheet',
      href: 'https://cdn.jsdelivr.net/npm/remixicon@4.2.0/fonts/remixicon.min.css',
    },
    ...(customFont.value ? [{ rel: 'stylesheet', href: customFont.value.split('|')[0] }] : []),
  ],
});

const fontFamily = computed(() => {
  if (!customFont.value) return '';
  const parts = customFont.value.split('|');
  return parts[1] ? parts[1] : '';
});

watchEffect(() => {
  if (import.meta.client && fontFamily.value) {
    document.documentElement.style.setProperty('--font-family', fontFamily.value);
    document.body.style.fontFamily = `${fontFamily.value}, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif`;
  }
});
</script>

<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
  <UiToast />
  <ClientOnly>
    <!-- eslint-disable-next-line vue/no-v-html -- 自定义头部代码，由管理员配置 -->
    <div v-if="customHead" v-html="customHead" />
    <!-- eslint-disable-next-line vue/no-v-html -- 自定义底部代码，由管理员配置 -->
    <div v-if="customBody" v-html="customBody" />
  </ClientOnly>
</template>
