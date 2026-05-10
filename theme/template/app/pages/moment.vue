<script lang="ts" setup>
import { getMoments } from '@/composables/useApi';
import type { Moment } from '@@/types';

const moments = ref<Moment[]>([]);
const total = ref(0);
const currentPage = ref(1);
const pageSize = ref(30);

useSeoMeta({
  title: '动态',
  description: '浏览我的动态',
});

const { data: initialData } = await useAsyncData('moments-list', async () => {
  const { list, total: resTotal } = await getMoments({
    page: 1,
    page_size: pageSize.value,
  });
  return { list, total: resTotal };
});

if (initialData.value) {
  moments.value = initialData.value.list || [];
  total.value = initialData.value.total || 0;
  currentPage.value = 1;
}

const fetchData = async (page = 1) => {
  currentPage.value = page;
  const { list, total: resTotal } = await getMoments({
    page,
    page_size: pageSize.value,
  });
  moments.value = list;
  total.value = resTotal;
};

const handlePageChange = (page: number) => {
  fetchData(page);
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};
</script>

<template>
  <div id="page">
    <h1 class="page-title">动态</h1>
    <div v-if="moments.length > 0" class="moment-list">
      <div v-for="moment in moments" :key="moment.id" class="moment-item">
        <div class="moment-content">
          <p v-if="moment.content.text">{{ moment.content.text }}</p>

          <div v-if="moment.content.images?.length" class="moment-images">
            <img
              v-for="(img, index) in moment.content.images"
              :key="index"
              :src="img"
              alt="moment image"
              loading="lazy"
            />
          </div>

          <div v-if="moment.content.link" class="moment-link">
            <a :href="moment.content.link.url" target="_blank" rel="noopener noreferrer">
              <img
                v-if="moment.content.link.favicon"
                :src="moment.content.link.favicon"
                alt=""
                class="link-favicon"
              />
              <span>{{ moment.content.link.title || moment.content.link.url }}</span>
              <i class="ri-external-link-line" />
            </a>
          </div>

          <div v-if="moment.content.video" class="moment-video">
            <iframe
              v-if="moment.content.video.platform === 'bilibili' && moment.content.video.video_id"
              :src="`//player.bilibili.com/player.html?bvid=${moment.content.video.video_id}&autoplay=0`"
              allowfullscreen
            />
            <iframe
              v-else-if="
                moment.content.video.platform === 'youtube' && moment.content.video.video_id
              "
              :src="`https://www.youtube.com/embed/${moment.content.video.video_id}`"
              allowfullscreen
            />
            <video v-else-if="moment.content.video.url" :src="moment.content.video.url" controls />
          </div>

          <div v-if="moment.content.audio" class="moment-audio">
            <audio :src="moment.content.audio.url" controls />
          </div>

          <div v-if="moment.content.music" class="moment-music">
            <div class="music-player">
              <i class="ri-music-2-line" />
              <span>音乐播放</span>
            </div>
          </div>
        </div>
        <div class="moment-meta">
          <span>{{ formatMomentTime(moment.publish_time) }}</span>
          <span v-if="moment.content.location"
            ><i class="ri-map-pin-line" /> {{ moment.content.location }}</span
          >
          <span v-if="moment.content.tags" class="moment-tags">#{{ moment.content.tags }}</span>
        </div>
      </div>
    </div>
    <div v-else class="empty-state">暂无动态</div>
    <UiPagination
      v-if="total > pageSize"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
      @change="handlePageChange"
    />
  </div>
</template>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border: 1px solid var(--flec-border);
  padding: 40px;
}

.page-title {
  margin: 0 0 30px;
  font-size: 2rem;
  font-weight: 700;
}

.moment-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.moment-item {
  padding-bottom: 24px;
  border-bottom: 1px solid var(--flec-border);
  &:last-child {
    border-bottom: none;
    padding-bottom: 0;
  }
}

.moment-content {
  p {
    margin: 0 0 12px;
    line-height: 1.6;
    color: var(--font-color);
  }
}

.moment-images {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 8px;
  margin-bottom: 12px;
  img {
    width: 100%;
    aspect-ratio: 1;
    object-fit: cover;
  }
}

.moment-link {
  margin-bottom: 12px;
  a {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: var(--flec-heavy-bg);
    color: var(--font-color);
    text-decoration: none;
    font-size: 0.9rem;
  }
  .link-favicon {
    width: 16px;
    height: 16px;
  }
}

.moment-video {
  margin-bottom: 12px;
  iframe,
  video {
    width: 100%;
    aspect-ratio: 16 / 9;
    border: none;
  }
}

.moment-audio {
  margin-bottom: 12px;
  audio {
    width: 100%;
    max-width: 400px;
  }
}

.moment-music {
  margin-bottom: 12px;
  .music-player {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    background: var(--flec-heavy-bg);
    color: var(--theme-color);
    font-size: 0.9rem;
  }
}

.moment-meta {
  display: flex;
  gap: 16px;
  font-size: 0.85rem;
  color: var(--theme-meta-color);
  i {
    margin-right: 4px;
  }
  .moment-tags {
    color: var(--theme-color);
  }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--theme-meta-color);
}

@media screen and (max-width: 768px) {
  #page {
    padding: 20px;
  }
  .page-title {
    font-size: 1.5rem;
  }
  .moment-images {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
