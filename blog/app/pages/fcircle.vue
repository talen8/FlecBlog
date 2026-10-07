<script lang="ts" setup>
import { useDebounceFn } from '@vueuse/core';
import { getFriendCircleRandom, getFriendCircleStats } from '../composables/api/friend';

definePageMeta({
  showSidebar: false,
});

useSeoMeta({
  title: '友圈',
  description: '聚合友链站点的最新文章，看看朋友们最近更新了什么',
});

const { keyword, items, total, hasMore, loading, loadMore } = useFriendCircle(30);

const { data: stats } = useAsyncData('friend-circle-stats', () => getFriendCircleStats());

const {
  data: featured,
  pending: featuredPending,
  refresh: refreshFeatured,
} = useAsyncData('friend-circle-random', () => getFriendCircleRandom(), { default: () => null });

const searchText = ref('');
watch(
  searchText,
  useDebounceFn((text: string) => (keyword.value = text), 300)
);

const clearSearch = () => {
  searchText.value = '';
  keyword.value = '';
};

const failedAvatars = reactive(new Set<number>());
const authorInitial = (name?: string) => name?.trim().charAt(0) || '?';
</script>

<template>
  <div id="fcircle-page">
    <h1 class="page-title">友圈</h1>

    <div class="circle-overview">
      <template v-if="stats">
        <div class="stat-card">
          <i class="ri-article-line stat-icon" />
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">文章总数</div>
        </div>

        <div class="stat-card">
          <i class="ri-links-line stat-icon" />
          <div class="stat-value">{{ stats.site_count }}</div>
          <div class="stat-label">订阅站点</div>
        </div>

        <div class="stat-card">
          <i class="ri-sun-line stat-icon" />
          <div class="stat-value">{{ stats.today_count }}</div>
          <div class="stat-label">今日更新</div>
        </div>
      </template>

      <!-- 随机文章 -->
      <div v-if="featured" class="circle-featured">
        <div class="featured-bar">
          <span class="featured-label"><i class="ri-shuffle-line" />随机推荐</span>
          <button
            type="button"
            class="featured-shuffle"
            :disabled="featuredPending"
            @click="refreshFeatured()"
          >
            <i class="ri-refresh-line" />换一篇
          </button>
        </div>

        <a :href="featured.link" target="_blank" rel="noopener noreferrer" class="featured-main">
          <div class="featured-avatar">
            <NuxtImg
              v-if="featured.author?.avatar && !failedAvatars.has(featured.id)"
              :src="featured.author.avatar"
              :alt="featured.author.name"
              @error="failedAvatars.add(featured.id)"
            />
            <span v-else class="featured-avatar-fallback">
              {{ authorInitial(featured.author?.name) }}
            </span>
          </div>

          <div class="featured-body">
            <div class="featured-title">{{ featured.title }}</div>
            <div class="featured-meta">
              <span>{{ featured.author?.name ?? '未知站点' }}</span>
              <template v-if="featured.published_at">
                <span class="featured-dot">·</span>
                <span>{{ formatMomentTime(featured.published_at) }}</span>
              </template>
            </div>
          </div>

          <i class="ri-arrow-right-line featured-go" />
        </a>
      </div>
    </div>

    <!-- 搜索 -->
    <div v-if="items.length > 0 || keyword" class="circle-search">
      <i class="ri-search-line search-icon" />
      <input
        v-model="searchText"
        type="text"
        class="search-input"
        placeholder="搜索文章标题或站点名称"
      />
      <button v-if="searchText" type="button" class="search-clear" @click="clearSearch">
        <i class="ri-close-line" />
      </button>
    </div>

    <!-- 搜索结果计数 -->
    <div v-if="keyword" class="circle-search-result">
      找到 {{ total }} 篇与「{{ keyword }}」相关的文章
    </div>

    <!-- 友圈卡片列表 -->
    <div v-if="items.length > 0" class="circle-list">
      <div v-for="item in items" :key="item.id" class="circle-card">
        <a :href="item.link" target="_blank" rel="noopener noreferrer" class="circle-title">
          {{ item.title }}
        </a>

        <div class="circle-meta">
          <div class="circle-avatar">
            <NuxtImg
              v-if="item.author?.avatar && !failedAvatars.has(item.id)"
              :src="item.author.avatar"
              :alt="item.author.name"
              loading="lazy"
              @error="failedAvatars.add(item.id)"
            />
            <span v-else class="circle-avatar-fallback">
              {{ authorInitial(item.author?.name) }}
            </span>
          </div>

          <span class="circle-author">{{ item.author?.name ?? '未知站点' }}</span>
          <span v-if="item.published_at" class="circle-time">
            {{ formatMomentTime(item.published_at) }}
          </span>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <i :class="keyword ? 'ri-search-line' : 'ri-links-line'" />
      <p>{{ keyword ? '没有找到相关文章' : '暂无友圈内容' }}</p>
    </div>

    <!-- 加载更多 -->
    <div v-if="hasMore" class="circle-more">
      <button type="button" class="more-button" :disabled="loading" @click="loadMore">
        {{ loading ? '加载中…' : '加载更多' }}
      </button>
    </div>

    <!-- 底部提示 -->
    <div v-else-if="items.length > 0" class="circle-tip">
      <i class="ri-information-line" />
      <span>已展示全部 {{ total }} 篇文章</span>
    </div>

    <!-- 内容声明 -->
    <div class="circle-disclaimer">
      <i class="ri-shield-check-line disclaimer-icon" />
      <div class="disclaimer-body">
        <p>
          本页内容由本站自动聚合自友链站点的 RSS
          源，标题与链接均指向原站原文，版权归原作者所有；本站仅作分享与导航，不代表赞同其观点，也不对其真实性与合法性负责。
        </p>
        <p>
          若你是相关内容的权利人，或发现其中存在侵权、违规内容，请通过
          <NuxtLink to="/feedback" class="disclaimer-link">意见反馈</NuxtLink>
          提交举报，核实后将即时从本页移除并停止收录该站点。
        </p>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;
#fcircle-page {
  @extend .cardHover;
  width: 100%;
  padding: 40px;

  .page-title {
    margin: 0 0 10px;
    font-weight: bold;
    font-size: 2rem;
  }

  .circle-overview {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;
    margin-bottom: 24px;

    .stat-card {
      @extend .cardHover;
      flex: 1 1 0;
      display: flex;
      flex-direction: column;
      justify-content: center;
      min-width: 0;
      gap: 4px;
      padding: 16px 18px;

      .stat-icon {
        font-size: 1.35rem;
        color: var(--theme-color);
      }

      .stat-value {
        font-size: 1.5rem;
        font-weight: 700;
        line-height: 1.2;
      }

      .stat-label {
        font-size: 13px;
        color: var(--theme-meta-color);
      }
    }
  }

  .circle-disclaimer {
    display: flex;
    gap: 12px;
    margin-top: 24px;
    padding: 16px 18px;
    border-radius: 8px;
    background: var(--flec-heavy-bg);
    font-size: 13px;
    line-height: 1.8;
    color: var(--theme-meta-color);

    .disclaimer-icon {
      flex-shrink: 0;
      font-size: 1.1rem;
      line-height: 1.8;
      color: var(--theme-color);
    }

    .disclaimer-body {
      min-width: 0;

      p {
        margin: 0 0 6px;
      }

      p:last-child {
        margin-bottom: 0;
      }
    }

    .disclaimer-link {
      color: var(--theme-color);

      &:hover {
        text-decoration: underline;
      }
    }
  }

  .circle-featured {
    @extend .cardHover;
    flex: 3 1 0;
    display: flex;
    flex-direction: column;
    justify-content: center;
    min-width: 0;
    gap: 12px;
    padding: 16px 20px;

    .featured-bar {
      display: flex;
      align-items: center;
      justify-content: space-between;
      font-size: 13px;
      color: var(--theme-meta-color);
    }

    .featured-label {
      display: flex;
      align-items: center;
      gap: 6px;

      i {
        font-size: 1rem;
        color: var(--theme-color);
      }
    }

    .featured-shuffle {
      display: flex;
      align-items: center;
      gap: 4px;
      padding: 2px 8px;
      font-size: 13px;
      color: var(--theme-meta-color);
      background: transparent;
      border: none;
      border-radius: 12px;
      cursor: pointer;
      transition: all 0.3s ease;

      i {
        font-size: 0.95rem;
        transition: transform 0.4s ease;
      }

      &:hover {
        color: var(--theme-color);
        background: var(--flec-heavy-bg);

        i {
          transform: rotate(180deg);
        }
      }
    }

    .featured-main {
      display: flex;
      align-items: center;
      gap: 14px;
      text-decoration: none;
      color: inherit;

      &:hover {
        .featured-title {
          color: var(--theme-color);
        }

        .featured-go {
          color: var(--theme-color);
          opacity: 1;
          transform: translateX(3px);
        }
      }
    }

    .featured-avatar {
      flex-shrink: 0;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 44px;
      height: 44px;
      border-radius: 50%;
      overflow: hidden;
      background: var(--flec-heavy-bg);

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }

      .featured-avatar-fallback {
        font-size: 18px;
        font-weight: 600;
        color: var(--theme-meta-color);
      }
    }

    .featured-body {
      flex: 1;
      min-width: 0;
      display: flex;
      flex-direction: column;
      gap: 4px;
    }

    .featured-title {
      font-size: 1.05rem;
      font-weight: 600;
      line-height: 1.5;
      overflow: hidden;
      text-overflow: ellipsis;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      line-clamp: 2;
      -webkit-box-orient: vertical;
      transition: color 0.3s ease;
    }

    .featured-meta {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 13px;
      color: var(--theme-meta-color);

      span {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }

    .featured-go {
      flex-shrink: 0;
      font-size: 1.2rem;
      color: var(--theme-meta-color);
      opacity: 0.4;
      transition: all 0.3s ease;
    }
  }

  .circle-search {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 20px;
    padding: 0 14px;
    border-radius: 24px;
    background: var(--flec-heavy-bg);
    transition: box-shadow 0.3s ease;

    &:focus-within {
      box-shadow: inset 0 0 0 1px var(--theme-color);
    }

    .search-icon {
      flex-shrink: 0;
      font-size: 1rem;
      color: var(--theme-meta-color);
    }

    .search-input {
      flex: 1;
      min-width: 0;
      padding: 10px 0;
      font-size: 14px;
      font-family: inherit;
      color: var(--font-color);
      background: transparent;
      border: none;
      outline: none;

      &::placeholder {
        color: var(--theme-meta-color);
      }
    }

    .search-clear {
      flex-shrink: 0;
      display: flex;
      align-items: center;
      padding: 4px;
      font-size: 1rem;
      color: var(--theme-meta-color);
      background: transparent;
      border: none;
      border-radius: 50%;
      cursor: pointer;
      transition: color 0.3s ease;

      &:hover {
        color: var(--theme-color);
      }
    }
  }

  .circle-search-result {
    margin-bottom: 16px;
    font-size: 13px;
    color: var(--theme-meta-color);
  }

  .circle-list {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 24px;

    .circle-card {
      @extend .cardHover;
      display: flex;
      flex-direction: column;
      gap: 12px;
      padding: 16px;
      overflow: hidden;

      .circle-meta {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 13px;
        color: var(--theme-meta-color);

        .circle-author {
          max-width: 45%;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .circle-time {
          flex-shrink: 0;
          margin-left: auto;
          opacity: 0.85;
        }
      }

      .circle-avatar {
        flex-shrink: 0;
        width: 22px;
        height: 22px;
        border-radius: 50%;
        overflow: hidden;
        background: var(--flec-heavy-bg);
        display: flex;
        align-items: center;
        justify-content: center;

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
        }

        .circle-avatar-fallback {
          font-size: 11px;
          font-weight: 600;
          color: var(--theme-meta-color);
        }
      }

      .circle-title {
        font-size: 16px;
        font-weight: 600;
        line-height: 1.5;
        min-height: 3em;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        color: inherit;
        text-decoration: none;
        transition: color 0.3s ease;

        &:hover {
          color: var(--theme-color);
        }
      }
    }
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 80px 20px;
    color: var(--theme-meta-color);

    i {
      font-size: 4rem;
      margin-bottom: 15px;
      opacity: 0.5;
    }

    p {
      font-size: 1.1rem;
      margin: 0;
    }
  }

  .circle-more {
    display: flex;
    justify-content: center;
    margin-top: 24px;

    .more-button {
      padding: 8px 24px;
      font-size: 14px;
      color: var(--font-color);
      background: var(--flec-heavy-bg);
      border: none;
      border-radius: 20px;
      cursor: pointer;
      transition: all 0.3s ease;

      &:hover:not(:disabled) {
        color: var(--theme-color);
      }

      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }
    }
  }

  .circle-tip {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    margin-top: 24px;
    padding: 15px 20px;
    font-size: 0.9rem;
    color: var(--theme-meta-color);

    i {
      font-size: 1.1rem;
    }
  }
}

@media screen and (max-width: 1024px) {
  #fcircle-page {
    padding: 30px;

    .page-title {
      font-size: 1.75rem;
    }

    .circle-overview {
      gap: 12px;

      .stat-card {
        padding: 14px 16px;
      }

      .stat-value {
        font-size: 1.35rem;
      }
    }

    .circle-list {
      grid-template-columns: repeat(2, 1fr);
      gap: 20px;
    }
  }
}

@media screen and (max-width: 768px) {
  #fcircle-page {
    padding: 18px;

    .page-title {
      font-size: 1.4rem;
    }

    .circle-disclaimer {
      gap: 10px;
      padding: 14px 16px;
      font-size: 0.8rem;
    }

    .circle-overview {
      gap: 10px;
      margin-bottom: 20px;

      .stat-label {
        font-size: 0.8rem;
      }
    }

    .circle-featured {
      flex-basis: 100%;
      padding: 14px 16px;

      .featured-main {
        gap: 10px;
      }

      .featured-avatar {
        width: 36px;
        height: 36px;
      }

      .featured-title {
        font-size: 0.95rem;
      }

      .featured-meta {
        font-size: 0.8rem;
      }
    }

    .circle-list {
      grid-template-columns: 1fr;
      gap: 16px;

      .circle-card {
        gap: 10px;
        padding: 14px;

        .circle-meta {
          gap: 6px;
          font-size: 0.8rem;
        }

        .circle-avatar {
          width: 20px;
          height: 20px;
        }

        .circle-title {
          font-size: 0.95rem;
        }
      }
    }

    .circle-tip {
      margin-top: 20px;
      padding: 12px 16px;
      font-size: 0.875rem;

      i {
        font-size: 1rem;
      }
    }
  }
}
</style>
