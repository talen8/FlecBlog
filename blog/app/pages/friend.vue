<script lang="ts" setup>
import type { FriendGroup } from '@@/types/friend';
import { getFriends } from '@/composables/api/friend';

definePageMeta({
  showSidebar: false,
});

useSeoMeta({
  title: '友链',
  description: '浏览我的友情链接，发现更多优秀的博客和网站',
});

const allGroups = ref<FriendGroup[]>([]);

// 使用SSR获取友链数据
const { data: initialData } = await useAsyncData('friends-list', async () => {
  try {
    const data = await getFriends();
    return data.groups || [];
  } catch (error) {
    console.error('获取友链失败:', error);
    return [];
  }
});

// 初始化数据
if (initialData.value) {
  allGroups.value = initialData.value;
}

// 正常友链分组（is_invalid = false）
const friendGroups = computed(() => {
  return allGroups.value
    .map((group: FriendGroup) => ({
      ...group,
      friends: group.friends.filter(f => !f.is_invalid),
    }))
    .filter((group: FriendGroup) => group.friends.length > 0);
});

// 失效友链分组（is_invalid = true）
const invalidFriendGroups = computed(() => {
  return allGroups.value
    .map((group: FriendGroup) => ({
      ...group,
      friends: group.friends.filter(f => f.is_invalid),
    }))
    .filter((group: FriendGroup) => group.friends.length > 0);
});

// 空状态判断
const isEmpty = computed(() => {
  return (
    allGroups.value.length === 0 ||
    allGroups.value.every((g: FriendGroup) => g.friends.length === 0)
  );
});
</script>

<template>
  <div id="friend-page">
    <h1 class="page-title">友链</h1>

    <div class="friend-sections">
      <!-- 友链分组 -->
      <section
        v-for="group in friendGroups"
        :key="group.type_id ?? 'uncategorized'"
        class="friend-section"
      >
        <h2 class="section-title">
          <i class="ri-links-line" />
          {{ group.type_name }}
        </h2>

        <div class="friend-list">
          <a
            v-for="friend in group.friends"
            :key="friend.id"
            :href="friend.url"
            target="_blank"
            class="friend-card"
            rel="noopener noreferrer"
            :title="friend.description"
          >
            <!-- 网站截图 -->
            <div class="friend-screenshot">
              <NuxtImg :src="friend.screenshot" :alt="friend.name" loading="lazy" />
            </div>

            <!-- 网站信息 -->
            <div class="friend-content">
              <NuxtImg :src="friend.avatar" :alt="friend.name" loading="lazy" />
              <div class="friend-info">
                <div class="friend-name">{{ friend.name }}</div>
                <div class="friend-description">{{ friend.description }}</div>
              </div>
            </div>
          </a>
        </div>
      </section>

      <!-- 失效友链 -->
      <section
        v-for="group in invalidFriendGroups"
        :key="'invalid-' + (group.type_id ?? 'uncategorized')"
        class="friend-section inactive-section"
      >
        <h2 class="section-title">
          <i class="ri-links-line" />
          失联友链
        </h2>

        <div class="inactive-list">
          <span
            v-for="friend in group.friends"
            :key="friend.id"
            class="inactive-tag"
            :title="friend.name + ' - ' + friend.description"
          >
            {{ friend.name }}
          </span>
        </div>
      </section>

      <!-- 空状态-->
      <div v-if="isEmpty" class="empty-state">
        <i class="ri-links-line" />
        <p>暂无友链数据</p>
      </div>
    </div>

    <!-- 申请友链 -->
    <h2 id="apply" class="section-title">申请友链</h2>
    <FeaturesFriendGuide />

    <!-- 评论区域 -->
    <LazyFeaturesCommentComments target-type="page" target-key="friend" />
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;
#friend-page {
  @extend .cardHover;
  width: 100%;
  padding: 40px;
}

.page-title {
  margin: 0 0 10px;
  font-weight: bold;
  font-size: 2rem;
}

.friend-section {
  margin-bottom: 40px;
}

.section-title {
  margin-bottom: 20px;
}

// 友链列表
.friend-list {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;

  .friend-card {
    @extend .cardHover;
    display: flex;
    flex-direction: column;
    text-decoration: none;
    color: inherit;
    overflow: hidden;
    cursor: pointer;

    &:hover {
      .friend-screenshot img {
        transform: scale(1.05);
      }

      .friend-content img {
        transform: scale(1.05) rotate(8deg);
      }

      .friend-name {
        color: var(--theme-color);
      }
    }

    .friend-screenshot {
      width: 100%;
      aspect-ratio: 16 / 9;
      overflow: hidden;
      background: transparent;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
      }
    }

    .friend-content {
      display: flex;
      gap: 12px;
      padding: 12px;

      img {
        width: 45px;
        height: 45px;
        border-radius: 50%;
        flex-shrink: 0;
        object-fit: cover;
        transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      }

      .friend-info {
        flex: 1;
        min-width: 0;
        display: flex;
        flex-direction: column;
        justify-content: center;
        gap: 6px;

        .friend-name {
          font-weight: 600;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
          transition: color 0.3s ease;
          line-height: 1.3;
        }

        .friend-description {
          font-size: 13px;
          line-height: 1.5;
          color: var(--theme-meta-color);
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }
      }
    }
  }
}

// 失效友链
.inactive-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;

  .inactive-tag {
    padding: 4px 10px;
    font-size: 14px;
    color: var(--theme-meta-color);
    text-decoration: line-through;
    background: var(--flec-card-bg);
    border-radius: 4px;
  }
}

// 空状态
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80px 20px;
}

// 响应式设计
@media screen and (max-width: 1024px) {
  #friend-page {
    padding: 30px;

    .page-title {
      font-size: 1.75rem;
    }

    .section-title {
      font-size: 1.35rem;
    }
  }

  .friend-list {
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
  }
}

@media screen and (max-width: 768px) {
  #friend-page {
    padding: 18px;

    .page-title {
      font-size: 1.4rem;
    }

    .section-title {
      font-size: 1.15rem;
      margin-bottom: 16px;
    }
  }

  .friend-list {
    grid-template-columns: 1fr;
    gap: 16px;

    .friend-card {
      .friend-screenshot {
        aspect-ratio: 16 / 9;
      }

      .friend-content {
        padding: 10px;
        gap: 10px;

        img {
          width: 40px;
          height: 40px;
        }

        .friend-info {
          gap: 4px;

          .friend-name {
            font-size: 0.9rem;
          }

          .friend-description {
            font-size: 0.8rem;
          }
        }
      }
    }
  }

  .inactive-list {
    .inactive-tag {
      font-size: 0.8rem;
      padding: 3px 8px;
    }
  }
}
</style>
