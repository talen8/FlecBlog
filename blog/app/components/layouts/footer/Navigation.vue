<script lang="ts" setup>
const { getMenus } = useTheme();
const footerMenus = getMenus('footer');
const { data: friendGroups, refresh, status } = await useFriends();

// 刷新友链（包装成无参，供模板 @click 调用，避免把事件对象传进 refresh）
const refreshFriends = () => refresh();

// 判断链接是否为外部链接
const isExternalLink = (url: string) => {
  return url.startsWith('http://') || url.startsWith('https://');
};

const isLoadingFriends = computed(() => status.value === 'pending');

// 所有有效友链（扁平化，过滤失效）
const validFriends = computed(() =>
  (friendGroups.value ?? []).flatMap(group => group.friends).filter(friend => !friend.is_invalid)
);

// 随机获取 3 个友链
const randomFriends = computed(() => {
  const all = validFriends.value;
  if (all.length <= 3) return all;
  return [...all].sort(() => Math.random() - 0.5).slice(0, 3);
});
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
    <ClientOnly>
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
            v-if="validFriends.length > 3"
            href="/friend"
            class="content_link"
            aria-label="查看更多友链"
          >
            更多...
          </a>
        </nav>
      </div>
    </ClientOnly>
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
