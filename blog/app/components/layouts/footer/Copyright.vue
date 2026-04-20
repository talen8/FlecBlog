<script lang="ts" setup>
import { parseJSON } from '@/utils/json';

const { basicConfig, blogConfig } = useSysConfig();
const startYear = 2024;
const copyrightYear = ref(`${startYear}`);

onMounted(() => {
  const currentYear = new Date().getFullYear();
  copyrightYear.value =
    startYear === currentYear ? `${currentYear}` : `${startYear} - ${currentYear}`;
});

/**
 * 页脚右侧链接列表
 * 从系统配置中读取 footer_links 字段
 */
const footerLinks = computed(() => {
  return parseJSON<Array<{ name: string; url: string }>>(blogConfig.value.footer_links, []).filter(
    item => item.name && item.url
  );
});

/**
 * 判断链接是否为外部链接
 * 以 / 开头的为内部链接，其他为外部链接
 * @param url - 链接地址
 * @returns 是否为外部链接
 */
const isExternalLink = (url: string) => {
  return !url.startsWith('/');
};
</script>

<template>
  <div class="footer-column">
    <div class="column-left">
      <div class="copyright">
        <span>©{{ copyrightYear }} By</span>
        <a
          :href="basicConfig.home_url || '#'"
          target="_blank"
          :aria-label="`作者 ${basicConfig.author}`"
          rel="noopener noreferrer"
          >{{ basicConfig.author }}</a
        >
      </div>
      <div v-if="basicConfig.icp || basicConfig.police_record" class="beian">
        <a
          v-if="basicConfig.icp"
          href="https://beian.miit.gov.cn/"
          target="_blank"
          :aria-label="`${basicConfig.icp} 备案信息`"
          rel="noopener noreferrer"
          >{{ basicConfig.icp }}</a
        >
        <a
          v-if="basicConfig.police_record"
          href="https://beian.mps.gov.cn/"
          target="_blank"
          :aria-label="`${basicConfig.police_record} 公安备案信息`"
          rel="noopener noreferrer"
          >{{ basicConfig.police_record }}</a
        >
      </div>
    </div>
    <div class="column-right">
      <!-- 可配置的页脚链接 -->
      <a
        v-for="link in footerLinks"
        :key="link.name"
        class="links"
        :href="link.url"
        :target="isExternalLink(link.url) ? '_blank' : '_self'"
        :rel="isExternalLink(link.url) ? 'noopener noreferrer' : undefined"
        :aria-label="link.name"
      >
        {{ link.name }}
      </a>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.footer-column {
  margin-top: 1rem;
  background: var(--flec-card-bg);
  display: flex;
  overflow: hidden;
  transition: 0.3s;
  width: 100%;
  justify-content: space-between;
  flex-wrap: wrap;
  align-items: center;
  line-height: 1;
  padding: 14px 5%;

  .column-left {
    display: flex;
    gap: 8px;
    flex-direction: column;

    .copyright {
      display: flex;
      flex-direction: row;
      align-items: center;

      a {
        margin: 0 4px;
        color: var(--flec-footer-font);
        font-weight: 700;
        white-space: nowrap;
        padding: 8px;
        border-radius: 32px;
        line-height: 1;
        display: flex;
        align-items: center;
        gap: 4px;

        &:hover {
          color: var(--flec-footer-font-hover);
          background: var(--flec-footer-font-bg-hover);
        }
      }
    }

    .beian {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;

      a {
        font-size: 0.9rem;
        font-weight: 400;
        color: var(--flec-footer-font);
        padding: 0;
        margin: 0;
        display: flex;
        align-items: center;
        gap: 3px;
      }
    }
  }

  .column-right {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    align-items: center;
    justify-content: center;

    .links {
      margin: 0 4px;
      color: var(--flec-footer-font);
      font-weight: 700;
      white-space: nowrap;
      padding: 8px;
      border-radius: 32px;
      line-height: 1;
      display: flex;
      align-items: center;
      gap: 4px;

      &:hover {
        color: var(--flec-footer-font-hover);
        background: var(--flec-footer-font-bg-hover);
      }
    }
  }
}

// 响应式设计
@media screen and (max-width: 768px) {
  .footer-column {
    flex-direction: column;
    text-align: center;
    padding: 18px 4%;

    .column-left {
      order: 2;
      align-items: center;

      .copyright {
        font-size: 0.95rem;
        justify-content: center;
      }

      .beian {
        justify-content: center;
        align-items: center;
        font-size: 0.82rem;
        flex-direction: column;
        gap: 4px;
      }
    }

    .column-right {
      order: 1;
      width: 100%;
      justify-content: center;
    }
  }
}
</style>
