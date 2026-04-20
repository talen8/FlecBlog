<script lang="ts" setup>
import type { Friend } from '@@/types/friend';
import { getFriends } from '@/composables/api/friend';

const { footerMenus } = useMenus();

// 判断链接是否为外部链接
const isExternalLink = (url: string) => {
  return url.startsWith('http://') || url.startsWith('https://');
};

// 友链数据
const friendGroups = ref<Friend[]>([]);
const isLoadingFriends = ref(false);

// 使用SSR获取友链数据
const { data: initialFriendData } = await useAsyncData('footer-friends', async () => {
  try {
    const data = await getFriends();
    const allFriends: Friend[] = [];
    data.groups?.forEach(group => {
      group.friends.forEach(friend => {
        if (!friend.is_invalid) {
          allFriends.push(friend);
        }
      });
    });
    return allFriends;
  } catch (error) {
    console.error('获取友链失败:', error);
    return [];
  }
});

// 初始化数据
if (initialFriendData.value) {
  friendGroups.value = initialFriendData.value;
}

// 获取所有友链（客户端刷新使用）
const fetchFriends = async () => {
  try {
    isLoadingFriends.value = true;
    const data = await getFriends();
    const allFriends: Friend[] = [];
    data.groups?.forEach(group => {
      group.friends.forEach(friend => {
        if (!friend.is_invalid) {
          allFriends.push(friend);
        }
      });
    });
    friendGroups.value = allFriends;
  } catch (error) {
    console.error('获取友链失败:', error);
  } finally {
    isLoadingFriends.value = false;
  }
};

// 随机获取3个友链
const randomFriends = computed(() => {
  if (friendGroups.value.length <= 3) {
    return friendGroups.value;
  }
  const shuffled = [...friendGroups.value].sort(() => Math.random() - 0.5);
  return shuffled.slice(0, 3);
});

// 刷新友链
const refreshFriends = () => {
  fetchFriends();
};
</script>

<template>
  <div v-if="footerMenus.length > 0" class="footer-group">
    <div v-for="menu in footerMenus" :key="menu.id" class="group-item">
      <div class="item-title" role="heading" aria-level="2">{{ menu.title }}</div>
      <nav class="item-content" :aria-label="`${menu.title}导航`">
        <a
          v-for="child in menu.children"
          :key="child.id"
          class="content_link"
          :href="child.url"
          :target="isExternalLink(child.url) ? '_blank' : '_self'"
          :rel="isExternalLink(child.url) ? 'noopener noreferrer' : undefined"
          :aria-label="child.title"
        >
          {{ child.title }}
        </a>
      </nav>
    </div>

    <!-- 友链列 -->
    <div class="group-item">
      <div class="item-title friend-title" role="heading" aria-level="2">
        友链
        <i
          class="refresh-icon ri-refresh-line"
          :class="{ 'is-loading': isLoadingFriends }"
          :aria-label="isLoadingFriends ? '正在加载友链' : '刷新友链'"
          @click="refreshFriends"
        />
      </div>
      <nav class="item-content friend-content" aria-label="友情链接">
        <a
          v-for="friend in randomFriends"
          :key="friend.id"
          class="content_link"
          :href="friend.url"
          target="_blank"
          rel="noopener noreferrer"
          :aria-label="friend.name"
          :title="friend.description"
        >
          {{ friend.name }}
        </a>
        <a
          v-if="friendGroups.length > 3"
          href="/friend"
          class="content_link"
          aria-label="查看更多友链"
        >
          更多...
        </a>
      </nav>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.footer-group {
  display: flex;
  flex-direction: row;
  width: 100%;
  max-width: 1200px;
  justify-content: space-between;
  flex-wrap: wrap;
  padding: 0 1rem;
  gap: 16px;
  margin-top: 24px;

  .group-item {
    display: flex;
    flex-direction: column;
    gap: 16px;

    .item-title {
      color: var(--flec-footer-font);
      margin-left: 8px;
      margin-top: 0;
      margin-bottom: 0;
      width: fit-content;
      display: flex;
      align-items: center;
      gap: 8px;
    }

    .friend-title {
      .refresh-icon {
        cursor: pointer;
        transition: transform 0.3s ease;
        font-size: 1.1em;

        &:hover {
          color: var(--flec-footer-font-hover);
        }

        &.is-loading {
          animation: rotate 1s linear infinite;
        }
      }
    }

    @keyframes rotate {
      from {
        transform: rotate(0deg);
      }

      to {
        transform: rotate(360deg);
      }
    }

    .item-content {
      display: flex;
      flex-direction: column;
      gap: 8px;

      .content_link {
        color: var(--flec-footer-font);
        line-height: 0.6rem;
        margin-right: auto;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        max-width: 100px;
        cursor: pointer;
        padding: 8px;
        border-radius: 12px;

        &:hover {
          color: var(--flec-footer-font-hover);
          background: var(--flec-footer-font-bg-hover);
        }
      }
    }

    .friend-content {
      .content_link {
        width: 100%;
        box-sizing: border-box;
        min-width: 120px;
      }
    }
  }
}

// 响应式设计
@media screen and (max-width: 768px) {
  .footer-group {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    padding: 0 12px;
  }
}

@media screen and (max-width: 400px) {
  .footer-group {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
