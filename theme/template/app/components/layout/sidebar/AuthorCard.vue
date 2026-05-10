<script lang="ts" setup>
import { parseJSON } from '@/utils/json';

const { basicConfig, blogConfig } = useSysConfig();
const { total: articlesTotal } = useArticles();
const { total: categoriesTotal } = useCategories();
const { total: tagsTotal } = useTags();

const avatarUrl = computed(() => basicConfig.value?.['basic.author_avatar'] || '/avatar.webp');
const authorName = computed(() => basicConfig.value?.['basic.author'] || '');
const authorDesc = computed(() => basicConfig.value?.['basic.author_desc'] || '');

const contacts = computed(() => {
  return parseJSON<Array<{ name: string; url: string; icon: string }>>(
    blogConfig.value?.['blog.sidebar_social'],
    []
  ).filter(item => item.url && item.url.trim() !== '');
});
</script>

<template>
  <div class="sidebar-card author-card">
    <img class="avatar" :src="avatarUrl" alt="头像" loading="lazy" />
    <div class="name">{{ authorName }}</div>
    <div v-if="authorDesc" class="desc">{{ authorDesc }}</div>
    <div class="stats">
      <NuxtLink to="/archive"
        ><span class="num">{{ articlesTotal }}</span
        ><span class="label">文章</span></NuxtLink
      >
      <NuxtLink to="/categories"
        ><span class="num">{{ categoriesTotal }}</span
        ><span class="label">分类</span></NuxtLink
      >
      <NuxtLink to="/tags"
        ><span class="num">{{ tagsTotal }}</span
        ><span class="label">标签</span></NuxtLink
      >
    </div>
    <NuxtLink to="/subscribe" class="subscribe-btn">订阅本站</NuxtLink>
    <div v-if="contacts.length" class="social">
      <a
        v-for="c in contacts"
        :key="c.name"
        :href="c.url"
        target="_blank"
        rel="noopener noreferrer"
        :aria-label="c.name"
      >
        <i :class="'ri-' + c.icon" />
      </a>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.author-card {
  text-align: center;
}
.avatar {
  width: 96px;
  height: 96px;
  object-fit: cover;
}
.name {
  margin-top: 8px;
  font-weight: 600;
}
.desc {
  font-size: 0.85rem;
  color: var(--theme-meta-color);
}
.stats {
  display: flex;
  justify-content: center;
  gap: 24px;
  margin-top: 12px;
}
.stats a {
  color: inherit;
  text-decoration: none;
  text-align: center;
}
.num {
  display: block;
  font-weight: 600;
}
.label {
  font-size: 0.8rem;
  color: var(--theme-meta-color);
}
.subscribe-btn {
  display: block;
  margin-top: 12px;
  padding: 6px 0;
  background: var(--theme-color);
  color: #fff;
  text-decoration: none;
  font-size: 0.9rem;
}
.social {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 10px;
}
.social a {
  color: var(--theme-meta-color);
  font-size: 1.2rem;
}
</style>
