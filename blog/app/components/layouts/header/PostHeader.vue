<script lang="ts" setup>
import { getArticleBySlug } from '@/composables/api/article';

const router = useRouter();
const route = useRoute();

// 使用 useAsyncData 确保服务端渲染时数据已加载
const { data: article } = await useAsyncData('post-header-article', async () => {
  const slug = route.params.slug as string;
  return await getArticleBySlug(slug);
});

// 计算文章字数（去除 Markdown 标记后的准确字数）
const wordCount = computed(() => {
  if (!article.value?.content) return 0;
  return countWords(article.value.content);
});

// 计算阅读时长（按每分钟300字计算）
const readingTime = computed(() => {
  if (!article.value?.content) return 0;
  return estimateReadingTime(article.value.content, 300);
});

// 文章评论总数
const commentCount = computed(() => {
  return article.value?.comment_count || 0;
});

// 跳转到分类详情页
const goToCategory = () => {
  if (article.value?.category?.url) {
    router.push(article.value.category.url);
  }
};
</script>

<template>
  <header
    class="post-header"
    :style="{ backgroundImage: article?.cover ? `url(${article.cover})` : 'none' }"
  >
    <div v-if="article" class="post-info">
      <h1 class="post-title">{{ article.title }}</h1>

      <!-- 移动端：合并为一行 -->
      <div class="post-meta post-meta-mobile">
        <span class="post-meta-item">
          <i class="ri-calendar-line" />
          <span>发表于 {{ formatFriendly(article.publish_time) }}</span>
        </span>
        <span v-if="article.update_time" class="post-meta-item">
          <i class="ri-refresh-line" />
          <span>更新于 {{ formatFriendly(article.update_time) }}</span>
        </span>
        <span v-if="article.location" class="post-meta-item">
          <i class="ri-map-pin-line" />
          <span>{{ article.location }}</span>
        </span>
        <span v-if="article.category" class="post-meta-item clickable" @click="goToCategory">
          <i class="ri-folder-line" />
          <span>{{ article.category.name }}</span>
        </span>
        <span class="post-meta-item">
          <i class="ri-file-word-line" />
          <span>总字数: {{ wordCount }}</span>
        </span>
        <span class="post-meta-item">
          <i class="ri-time-line" />
          <span>阅读时长: {{ readingTime }}分钟</span>
        </span>
        <span class="post-meta-item">
          <i class="ri-eye-line" />
          <span>浏览量: {{ article.view_count }}</span>
        </span>
        <span class="post-meta-item clickable" @click="scrollToElement('.comment-input')">
          <i class="ri-message-3-line" />
          <span>评论数: {{ commentCount }}</span>
        </span>
      </div>

      <!-- 桌面端：分两行显示 -->
      <div class="post-meta-desktop">
        <div class="post-meta">
          <span class="post-meta-item">
            <i class="ri-calendar-line" />
            <span>发表于 {{ formatFriendly(article.publish_time) }}</span>
          </span>
          <span v-if="article.update_time" class="post-meta-item">
            <i class="ri-refresh-line" />
            <span>更新于 {{ formatFriendly(article.update_time) }}</span>
          </span>
          <span v-if="article.location" class="post-meta-item">
            <i class="ri-map-pin-line" />
            <span>{{ article.location }}</span>
          </span>
          <span v-if="article.category" class="post-meta-item clickable" @click="goToCategory">
            <i class="ri-folder-line" />
            <span>{{ article.category.name }}</span>
          </span>
        </div>
        <div class="post-meta">
          <span class="post-meta-item">
            <i class="ri-file-word-line" />
            <span>总字数: {{ wordCount }}</span>
          </span>
          <span class="post-meta-item">
            <i class="ri-time-line" />
            <span>阅读时长: {{ readingTime }}分钟</span>
          </span>
          <span class="post-meta-item">
            <i class="ri-eye-line" />
            <span>浏览量: {{ article.view_count }}</span>
          </span>
          <span class="post-meta-item clickable" @click="scrollToElement('.comment-input')">
            <i class="ri-message-3-line" />
            <span>评论数: {{ commentCount }}</span>
          </span>
        </div>
      </div>
    </div>
  </header>
</template>

<style lang="scss" scoped>
.post-header {
  position: relative;
  height: 400px;
  width: 100%;
  margin-top: -4rem; // 向上延伸 4rem（NavBar 的高度）
  padding-top: 4rem; // 补偿内容位置
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  z-index: 0; // 确保在 NavBar (z-index: 50) 下方

  // 添加遮罩层，确保文字清晰可见
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(8px);
    z-index: 1;
  }

  .post-info {
    position: absolute;
    top: 35%;
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-direction: column;
    padding: 0 2rem;
    z-index: 2;

    .post-title {
      font-size: 2.2rem;
      color: #fff;
      text-align: center;
    }

    // 默认显示桌面端，隐藏移动端
    .post-meta-mobile {
      display: none !important;
    }

    .post-meta-desktop {
      display: block !important;
    }

    .post-meta {
      display: flex;
      align-items: center;
      flex-wrap: wrap;
      justify-content: center;
      color: rgba(255, 255, 255, 0.9);

      .post-meta-item {
        display: flex;
        align-items: center;

        &:not(:last-child)::after {
          content: '|';
          color: rgba(255, 255, 255, 0.6);
          margin: 0 0.5rem;
        }

        i {
          font-size: 1rem;
          margin-right: 0.3rem;
        }

        &.clickable {
          cursor: pointer;
          transition: all 0.3s ease;

          &:hover {
            color: rgba(255, 255, 255, 1);
          }
        }
      }
    }
  }
}

// 响应式设计
@media screen and (max-width: 768px) {
  .post-header {
    height: 350px;

    .post-info {
      padding: 0 1.25rem;

      .post-meta {
        font-size: 0.8rem;
        line-height: 1.6;
      }
    }
  }
}

@media screen and (max-width: 500px) {
  .post-header {
    .post-info {
      // 移动端显示移动版，隐藏桌面版
      .post-meta-mobile {
        display: flex !important;
      }

      .post-meta-desktop {
        display: none !important;
      }
    }
  }
}
</style>
