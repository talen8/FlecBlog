<script lang="ts" setup>
import { parseJSON } from '@/utils/json';

const { blogConfig, basicConfig } = useSysConfig();

const currentYear = new Date().getFullYear();
const blogTitle = computed(() => blogConfig.value['blog.title'] || 'Blog');
const icp = computed(() => basicConfig.value?.['basic.icp'] || '');
const policeRecord = computed(() => basicConfig.value?.['basic.police_record'] || '');

const footerSocial = computed(() => {
  return parseJSON<Array<{ name: string; url: string; icon: string; position?: string }>>(
    blogConfig.value?.['blog.footer_social'],
    []
  ).filter(item => item.url && item.url.trim() !== '');
});

const footerLinks = computed(() => {
  const configured = parseJSON<Array<{ name: string; url: string }>>(
    blogConfig.value?.['blog.footer_links'],
    []
  ).filter(item => item.name && item.url);

  if (configured.length > 0) return configured;

  return [
    { name: '关于', url: '/about' },
    { name: '友链', url: '/friend' },
    { name: '动态', url: '/moment' },
    { name: '留言板', url: '/message' },
    { name: '反馈', url: '/feedback' },
  ];
});

const isExternalLink = (url: string) => {
  return !url.startsWith('/');
};
</script>

<template>
  <footer class="footer">
    <div class="footer-container">
      <div v-if="footerSocial.length > 0" class="footer-social">
        <a
          v-for="item in footerSocial"
          :key="item.name"
          :href="item.url"
          target="_blank"
          rel="noopener noreferrer"
          class="social-link"
          :aria-label="item.name"
        >
          <i :class="'ri-' + item.icon" />
        </a>
      </div>
      <div class="footer-links">
        <NuxtLink
          v-for="link in footerLinks"
          :key="link.name + link.url"
          :to="isExternalLink(link.url) ? undefined : link.url"
          :href="isExternalLink(link.url) ? link.url : undefined"
          :target="isExternalLink(link.url) ? '_blank' : undefined"
          :rel="isExternalLink(link.url) ? 'noopener noreferrer' : undefined"
          class="footer-link"
        >
          {{ link.name }}
        </NuxtLink>
      </div>
      <p class="copyright">© {{ currentYear }} {{ blogTitle }}</p>
      <div v-if="icp || policeRecord" class="records">
        <a v-if="icp" href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer">{{
          icp
        }}</a>
        <a
          v-if="policeRecord"
          href="https://beian.mps.gov.cn/"
          target="_blank"
          rel="noopener noreferrer"
          >{{ policeRecord }}</a
        >
      </div>
    </div>
  </footer>
</template>

<style lang="scss" scoped>
.footer {
  background: var(--flec-card-bg);
  border-top: 1px solid var(--flec-border);
  padding: 20px 15px;
  margin-top: auto;
}

.footer-container {
  max-width: 900px;
  margin: 0 auto;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.copyright {
  margin: 0;
  color: var(--theme-meta-color);
  font-size: 0.9rem;
}

.footer-social {
  display: flex;
  gap: 14px;
  align-items: center;
}

.social-link {
  color: var(--theme-meta-color);
  font-size: 1.2rem;
}

.footer-links {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
  justify-content: center;
}

.footer-link {
  color: var(--theme-meta-color);
  text-decoration: none;
  font-size: 0.9rem;
}

.records {
  display: flex;
  justify-content: center;
  gap: 16px;

  a {
    color: var(--theme-meta-color);
    text-decoration: none;
    font-size: 0.85rem;
  }
}
</style>
