<script setup lang="ts">
import { parseJSON } from '@/utils/json';

const { basicConfig, blogConfig } = useSysConfig();
const { total: articlesTotal, fetchArticles } = useArticles();
const avatarUrl = computed(() => basicConfig.value.author_avatar || '/avatar.webp');
const { total: categoriesTotal } = useCategories();
const { total: tagsTotal } = useTags();

const contacts = computed(() => {
  const socialList = parseJSON<Array<{ name: string; url: string; icon: string }>>(
    blogConfig.value.sidebar_social,
    []
  );
  return socialList.filter(item => item.url && item.url.trim() !== '');
});

onMounted(async () => {
  if (articlesTotal.value === 0) {
    await fetchArticles({ page: 1, page_size: 1 });
  }
});
</script>

<template>
  <div class="card-widget card-info">
    <div class="author-avatar">
      <NuxtImg :src="avatarUrl" alt="头像" loading="lazy" />
    </div>
    <div class="author-name">{{ basicConfig.author }}</div>
    <div class="author-desc">{{ basicConfig.author_desc }}</div>
    <div class="site-data">
      <router-link to="/archive" :aria-label="`查看全部 ${articlesTotal} 篇文章`">
        <div class="headline">文章</div>
        <div class="num">{{ articlesTotal }}</div>
      </router-link>
      <router-link to="/categories" :aria-label="`查看全部 ${categoriesTotal} 个分类`">
        <div class="headline">分类</div>
        <div class="num">{{ categoriesTotal }}</div>
      </router-link>
      <router-link to="/tags" :aria-label="`查看全部 ${tagsTotal} 个标签`">
        <div class="headline">标签</div>
        <div class="num">{{ tagsTotal }}</div>
      </router-link>
    </div>
    <a href="/subscribe" class="card-info-btn">订阅本站</a>
    <div class="card-info-icons">
      <a
        v-for="contact in contacts"
        :key="contact.name"
        :href="contact.url"
        class="icon"
        target="_blank"
        :aria-label="`访问 ${contact.name}`"
        rel="noopener noreferrer"
      >
        <i :class="'ri-' + contact.icon" aria-hidden="true" />
      </a>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.card-info {
  text-align: center;

  .author-avatar {
    overflow: hidden;
    margin: 0 auto;
    width: 110px;
    height: 110px;
    border-radius: 50%;

    img {
      object-fit: cover;
      width: 100%;
      height: 100%;
    }
  }

  .author-name {
    font-weight: 500;
    font-size: 1.57em;
  }

  .author-desc {
    margin-top: -0.42em;
  }

  .site-data {
    margin: 14px 0 4px;
    display: table;
    width: 100%;
    table-layout: fixed;

    a {
      display: table-cell;
      transition: all 0.2s;

      .headline {
        font-size: 0.95rem;
      }

      .num {
        margin-top: -0.45rem;
        font-size: 1.2rem;
      }

      &:hover {
        color: var(--flec-btn-hover);
      }
    }
  }

  .card-info-btn {
    display: block;
    margin-top: 14px;
    background-color: var(--flec-btn);
    color: #fff;
    line-height: 2.4;
    border-radius: 7px;

    &:hover {
      background-color: var(--flec-btn-hover);
    }
  }

  .card-info-icons {
    margin: 6px 0 -6px;

    .icon {
      margin: 0 10px;
      color: var(--font-color);
      font-size: 1.4em;
    }
  }
}

@media screen and (max-width: 900px) {
  .card-info {
    .author-avatar {
      width: 90px;
      height: 90px;
    }

    .author-name {
      font-size: 1.4em;
    }

    .site-data {
      margin: 10px 0 2px;

      a {
        .headline {
          font-size: 0.9rem;
        }

        .num {
          font-size: 1.1rem;
        }
      }
    }

    .card-info-btn {
      margin-top: 10px;
      line-height: 2.2;
    }

    .card-info-icons {
      margin: 4px 0 -4px;

      .icon {
        margin: 0 8px;
        font-size: 1.3em;
      }
    }
  }
}
</style>
