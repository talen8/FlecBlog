<script lang="ts" setup>
import mediumZoom from 'medium-zoom';
import type { Moment } from '@@/types/moment';
import { getMoments } from '@/composables/api/moment';

const { basicConfig, blogConfig } = useSysConfig();
const avatarUrl = computed(() => basicConfig.value.author_avatar || '/avatar.webp');
const momentsPageSize = computed(() => {
  const configSize = parseInt(blogConfig.value['moments_size'] || '30');
  return configSize > 0 ? configSize : 30;
});

definePageMeta({
  showSidebar: false,
});

useSeoMeta({
  title: '动态',
  description: '查看我的最新动态，分享生活点滴和即时想法',
});

const { moments } = useMoments();

// 使用SSR获取动态列表
const { data: initialData } = await useAsyncData('moments-list', async () => {
  const response = await getMoments({
    page: 1,
    page_size: momentsPageSize.value,
  });
  return response;
});

// 初始化数据
if (initialData.value) {
  moments.value = initialData.value.list || [];
}

const { waterfall, isLayoutReady } = useWaterfall({
  containerSelector: '#moment-list',
  columns: 3,
  gap: 15,
  debounceDelay: 150,
  waitForImages: true,
  breakpoints: { mobile: 768, tablet: 1200 },
});

// 图片缩放实例
let zoom: ReturnType<typeof mediumZoom> | null = null;

// 初始化图片缩放
const initZoom = () => {
  const contentEl = document.querySelector('#moment-list');
  if (!contentEl) return;

  const images = contentEl.querySelectorAll('.moment-images img');
  if (images.length === 0) return;

  // 如果已有实例，先销毁
  if (zoom) {
    zoom.detach();
  }

  // 初始化新的缩放实例
  zoom = mediumZoom(images, {
    margin: 24,
    background: 'rgba(0, 0, 0, 0.9)',
    scrollOffset: 48,
  });
};

onMounted(async () => {
  await nextTick();
  await waterfall();
  initZoom();
});

watch(
  () => moments.value.length,
  async () => {
    await nextTick();
    await waterfall();
    initZoom();
  }
);

// 组件卸载时清理
onUnmounted(() => {
  if (zoom) {
    zoom.detach();
    zoom = null;
  }
});

const getMomentContentType = (moment: Moment) => {
  if (moment.content.images?.length) return '图片动态';
  if (moment.content.video) return '视频动态';
  if (moment.content.audio) return '音频动态';
  if (moment.content.music) return '音乐动态';
  if (moment.content.link) return '链接分享';
  return '动态';
};

const handleCommentClick = (moment: Moment) => {
  const text = moment.content.text;
  const quote = text
    ? `> ${text.length > 100 ? text.substring(0, 100) + '...' : text}\n\n`
    : `> [${getMomentContentType(moment)}]\n\n`;
  fillComment(quote);
};
</script>

<template>
  <div id="moment-page">
    <h1 class="page-title">动态</h1>

    <div v-if="moments.length === 0" class="empty-state">
      <i class="ri-chat-3-line" />
      <p>暂无动态</p>
    </div>

    <div v-else id="moment-list" class="moment-list">
      <div
        v-for="moment in moments"
        :key="moment.id"
        class="moment-item"
        :class="{ 'layout-ready': isLayoutReady }"
      >
        <!-- 上部分：头像、作者、时间 -->
        <div class="moment-header">
          <div class="moment-avatar">
            <NuxtImg :src="avatarUrl" alt="avatar" loading="lazy" />
          </div>
          <div class="moment-meta">
            <div class="moment-author">{{ basicConfig.author }}</div>
            <div class="moment-time">{{ formatMomentTime(moment.publish_time) }}</div>
          </div>
        </div>

        <!-- 中部分：内容 -->
        <div class="moment-content">
          <!-- 文本内容 -->
          <div v-if="moment.content.text" class="moment-text">
            {{ moment.content.text }}
          </div>

          <!-- 图片内容 -->
          <div
            v-if="moment.content.images?.length"
            class="moment-images"
            :class="`images-${Math.min(moment.content.images.length, 6)}`"
          >
            <div
              v-for="(image, index) in moment.content.images.slice(0, 6)"
              :key="index"
              class="image-item"
            >
              <NuxtImg :src="image" :alt="`图片 ${index + 1}`" loading="lazy" />
              <div
                v-if="index === 5 && moment.content.images.length > 6"
                class="more-images-overlay"
              >
                <i class="ri-image-line" />
                <span>+{{ moment.content.images.length - 6 }}</span>
              </div>
            </div>
          </div>

          <!-- 视频内容 -->
          <div v-if="moment.content.video" class="moment-video">
            <video
              v-if="!moment.content.video.platform || moment.content.video.platform === 'local'"
              :src="moment.content.video.url"
              controls
              preload="metadata"
            />

            <iframe
              v-else-if="moment.content.video.platform === 'bilibili'"
              :src="`//player.bilibili.com/player.html?bvid=${moment.content.video.video_id}&autoplay=0`"
              scrolling="no"
              border="0"
              frameborder="no"
              framespacing="0"
              allowfullscreen="true"
            />

            <iframe
              v-else-if="moment.content.video.platform === 'youtube'"
              :src="`https://www.youtube.com/embed/${moment.content.video.video_id}`"
              frameborder="0"
              allow="
                accelerometer;
                autoplay;
                clipboard-write;
                encrypted-media;
                gyroscope;
                picture-in-picture;
              "
              allowfullscreen
            />
          </div>

          <!-- 音频内容 -->
          <div v-if="moment.content.audio" class="moment-audio">
            <audio :src="moment.content.audio.url" controls preload="metadata" />
          </div>

          <!-- 音乐内容 -->
          <div v-if="moment.content.music" class="moment-music">
            <FeaturesMomentMusicPlayer :music="moment.content.music" />
          </div>

          <!-- 链接内容 -->
          <a
            v-if="moment.content.link"
            :href="moment.content.link.url"
            target="_blank"
            rel="noopener noreferrer"
            class="moment-link"
          >
            <NuxtImg
              v-if="moment.content.link.favicon"
              :src="moment.content.link.favicon"
              alt="favicon"
              loading="lazy"
              class="link-favicon"
            />
            <div class="link-info">
              <div class="link-title">{{ moment.content.link.title }}</div>
              <div class="link-url">{{ moment.content.link.url }}</div>
            </div>
            <i class="ri-external-link-line" />
          </a>
        </div>

        <!-- 下部分：位置、分类标签、评论按钮 -->
        <div class="moment-footer">
          <div class="moment-info">
            <span v-if="moment.content.location" class="location">
              <i class="ri-map-pin-line" />
              {{ moment.content.location }}
            </span>
            <span v-if="moment.content.tags" class="tags">
              <i class="ri-price-tag-3-line" />
              {{ moment.content.tags }}
            </span>
          </div>
          <div class="moment-actions">
            <button class="comment-btn" title="评论此动态" @click="handleCommentClick(moment)">
              <i class="ri-chat-3-line" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部提示 -->
    <div v-if="moments.length > 0" class="moment-tip">
      <i class="ri-information-line" />
      <span>只显示最近{{ momentsPageSize }}条动态</span>
    </div>

    <!-- 评论区域 -->
    <LazyFeaturesCommentComments target-type="page" target-key="moment" />
  </div>
</template>

<style>
/* medium-zoom 样式覆盖 */
.medium-zoom-overlay {
  z-index: 9999 !important;
}

.medium-zoom-image {
  z-index: 10000 !important;
}
</style>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#moment-page {
  @extend .cardHover;
  align-self: flex-start;
  padding: 40px;

  .page-title {
    margin: 0 0 10px;
    font-weight: bold;
    font-size: 2rem;
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
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

  .moment-list {
    position: relative;
    width: 100%;
  }

  .moment-item {
    @extend .cardHover;
    position: absolute;
    padding: 0;
    overflow: hidden;
    opacity: 0;

    &.layout-ready {
      opacity: 1;
    }
  }

  // 上部分：头像、作者、时间
  .moment-header {
    display: flex;
    align-items: center;
    padding: 0.5rem 1rem;
    border-bottom: 1px solid var(--flec-moment-divider);

    .moment-avatar {
      width: 50px;
      height: 50px;
      border-radius: 10px;
      overflow: hidden;
      margin-right: 12px;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .moment-meta {
      flex: 1;
      min-width: 0;

      .moment-author {
        font-weight: 600;
        color: var(--flec-moment-title);
      }

      .moment-time {
        font-size: 0.875rem;
        color: var(--flec-moment-date);
      }
    }
  }

  // 中部分：内容
  .moment-content {
    padding: 15px 20px;
    border-bottom: 1px solid var(--flec-moment-divider);

    .moment-text {
      line-height: 1.7;
      color: var(--flec-moment-font);
      margin-bottom: 12px;
      white-space: pre-wrap;
      word-break: break-word;

      &:last-child {
        margin-bottom: 0;
      }
    }

    .moment-images {
      display: grid;
      gap: 6px;
      margin-top: 12px;

      // 1张图片：100%宽，高度自动
      &.images-1 {
        grid-template-columns: 1fr;

        .image-item {
          padding-bottom: 0;
          height: auto;

          img {
            position: relative;
            height: auto;
            max-height: 500px;
          }
        }
      }

      // 2张图片：一行2个
      &.images-2 {
        grid-template-columns: repeat(2, 1fr);
      }

      // 3张图片：一行3个
      &.images-3 {
        grid-template-columns: repeat(3, 1fr);
      }

      // 4张图片：2+2结构
      &.images-4 {
        grid-template-columns: repeat(2, 1fr);
      }

      // 5张图片：2+3结构
      &.images-5 {
        grid-template-columns: repeat(6, 1fr);

        .image-item:nth-child(1),
        .image-item:nth-child(2) {
          grid-column: span 3;
        }

        .image-item:nth-child(n + 3) {
          grid-column: span 2;
        }
      }

      // 6张图片：3+3结构
      &.images-6 {
        grid-template-columns: repeat(3, 1fr);
      }

      .image-item {
        position: relative;
        width: 100%;
        padding-bottom: 100%;
        overflow: hidden;
        border-radius: 6px;
        cursor: zoom-in;
        background: #f5f5f5;
        transition: transform 0.3s ease;

        &:hover {
          transform: translate3d(0, -2px, 0) scale(1.02);
        }

        img {
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          object-fit: cover;
          transition: transform 0.2s ease;

          &:hover {
            transform: scale(1.02);
          }
        }

        // 剩余图片数量覆盖层
        .more-images-overlay {
          position: absolute;
          bottom: 0;
          right: 0;
          left: 0;
          top: 0;
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          background: rgba(0, 0, 0, 0.6);
          color: #fff;
          font-weight: 600;
          border-radius: 6px;
          backdrop-filter: blur(2px);
          transition: background 0.3s ease;
        }

        &:hover .more-images-overlay {
          background: rgba(0, 0, 0, 0.7);
        }
      }
    }

    .moment-video {
      margin-top: 12px;
      border-radius: 6px;
      overflow: hidden;
      background: #000;
      transition: transform 0.3s ease;

      video,
      iframe {
        width: 100%;
        height: auto;
        aspect-ratio: 16 / 9;
        border: none;
        display: block;
      }

      &:hover {
        transform: translate3d(0, -2px, 0) scale(1.02);
      }
    }

    .moment-music {
      margin-top: 12px;
      transition: transform 0.3s ease;

      &:hover {
        transform: translate3d(0, -2px, 0) scale(1.02);
      }
    }

    .moment-audio {
      margin-top: 12px;
      border-radius: 6px;
      overflow: hidden;
      background: var(--flec-moment-card-bg);
      transition: transform 0.3s ease;

      audio {
        width: 100%;
        height: 50px;
        display: block;
      }

      &:hover {
        transform: translate3d(0, -2px, 0) scale(1.02);
      }
    }

    .moment-link {
      display: flex;
      align-items: center;
      margin-top: 12px;
      padding: 12px;
      background: var(--flec-moment-card-bg);
      border-radius: 6px;
      text-decoration: none;
      color: var(--flec-moment-font);
      transition: transform 0.3s ease;

      .link-favicon {
        flex-shrink: 0;
        width: 50px;
        height: 50px;
        margin-right: 12px;
        border-radius: 4px;
      }

      .link-info {
        flex: 1;
        min-width: 0;

        .link-title {
          font-size: 0.9rem;
          font-weight: 500;
          margin-bottom: 3px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .link-url {
          font-size: 0.75rem;
          color: var(--flec-moment-date);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }

      i {
        margin-left: 10px;
        font-size: 1.1rem;
        color: var(--flec-moment-date);
      }

      &:hover {
        transform: translate3d(0, -2px, 0) scale(1.02);
      }
    }
  }

  // 下部分：位置、分类标签、评论按钮
  .moment-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 10px 20px;

    .moment-info {
      display: flex;
      align-items: center;
      gap: 12px;
      font-size: 0.85rem;
      color: var(--flec-moment-date);

      .location {
        display: flex;
        align-items: center;
        gap: 4px;

        i {
          font-size: 0.9rem;
        }
      }
    }

    .moment-actions {
      display: flex;
      align-items: center;
      gap: 10px;

      .comment-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
        border: none;
        background: transparent;
        color: var(--flec-moment-date);
        cursor: pointer;
        border-radius: 4px;
        transition: all 0.3s ease;

        i {
          font-size: 1.1rem;
        }

        &:hover {
          background: var(--flec-moment-card-bg);
          color: #49b1f5;
        }
      }
    }
  }

  .moment-tip {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    margin-top: 30px;
    padding: 15px 20px;
    text-align: center;

    i {
      font-size: 1.1rem;
    }
  }
}

// 响应式设计
@media screen and (max-width: 1024px) {
  #moment-page {
    padding: 30px;

    .page-title {
      font-size: 1.75rem;
    }

    .moment-header {
      .moment-avatar {
        width: 45px;
        height: 45px;
      }

      .moment-meta {
        .moment-author {
          font-size: 0.95rem;
        }

        .moment-time {
          font-size: 0.8rem;
        }
      }
    }

    .moment-content {
      padding: 12px 18px;

      .moment-text {
        font-size: 0.95rem;
      }
    }

    .moment-footer {
      padding: 8px 18px;

      .moment-info {
        font-size: 0.8rem;
      }
    }
  }
}

@media screen and (max-width: 768px) {
  #moment-page {
    padding: 18px;

    .page-title {
      font-size: 1.4rem;
      margin-bottom: 8px;
    }

    .moment-header {
      padding: 0.4rem 0.8rem;

      .moment-avatar {
        width: 40px;
        height: 40px;
        margin-right: 10px;
      }

      .moment-meta {
        .moment-author {
          font-size: 0.9rem;
        }

        .moment-time {
          font-size: 0.75rem;
        }
      }
    }

    .moment-content {
      padding: 10px 14px;

      .moment-text {
        font-size: 0.9rem;
        line-height: 1.6;
      }

      .moment-images {
        gap: 4px;
        margin-top: 10px;

        // 在移动端，5张图片改为 2+3 结构
        &.images-5 {
          grid-template-columns: repeat(4, 1fr);

          .image-item:nth-child(1),
          .image-item:nth-child(2) {
            grid-column: span 2;
          }

          .image-item:nth-child(n + 3) {
            grid-column: span 1;
          }
        }

        // 在移动端，6张图片改为 2x3 结构
        &.images-6 {
          grid-template-columns: repeat(2, 1fr);
        }
      }

      .moment-link {
        padding: 10px;

        .link-favicon {
          width: 40px;
          height: 40px;
          margin-right: 10px;
        }

        .link-info {
          .link-title {
            font-size: 0.85rem;
          }

          .link-url {
            font-size: 0.7rem;
          }
        }

        i {
          font-size: 1rem;
        }
      }
    }

    .moment-footer {
      padding: 8px 14px;

      .moment-info {
        font-size: 0.75rem;
        gap: 10px;
        flex-wrap: wrap;

        .location,
        .tags {
          i {
            font-size: 0.85rem;
          }
        }
      }

      .moment-actions {
        .comment-btn {
          width: 28px;
          height: 28px;

          i {
            font-size: 1rem;
          }
        }
      }
    }

    .moment-tip {
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
