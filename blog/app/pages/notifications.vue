<script setup lang="ts">
import type { Notification } from '@@/types/notification';

const {
  notifications,
  total,
  currentPage,
  pageSize,
  unreadCount,
  loading,
  fetchNotifications,
  markNotificationAsRead,
  markAllNotificationsAsRead,
  resetPage,
} = useNotifications();
const router = useRouter();

definePageMeta({});

useSeoMeta({
  title: '通知中心',
  description: '查看您的所有通知，及时了解最新动态',
});

const loadNotifications = async () => {
  await fetchNotifications({
    page: currentPage.value,
    page_size: pageSize.value,
  });
};

const handleMarkAllAsRead = async () => {
  await markAllNotificationsAsRead();
  await loadNotifications();
};

const handlePageChange = (page: number) => {
  currentPage.value = page;
  loadNotifications();
};

const handleNotificationClick = (notification: Notification) => {
  if (!notification.is_read) {
    markNotificationAsRead(notification.id).catch(console.error);
  }
  if (notification.link) {
    router.push(notification.link);
  }
};

onMounted(() => {
  // 只在客户端获取需要认证的数据
  if (import.meta.client) {
    resetPage();
    loadNotifications();
  }
});
</script>

<template>
  <div id="page">
    <h1 class="page-title">通知中心</h1>

    <div class="notification-list">
      <div class="list-header">
        <h3 class="title">
          通知
          <span v-if="unreadCount > 0" class="count">({{ unreadCount }})</span>
        </h3>
        <button v-if="unreadCount > 0" class="mark-all-btn" @click="handleMarkAllAsRead">
          全部已读
        </button>
      </div>

      <div v-if="loading" class="loading">加载中...</div>

      <div v-else-if="notifications.length > 0" class="list-content">
        <div
          v-for="notification in notifications"
          :key="notification.id"
          class="notification-item"
          :class="{ unread: !notification.is_read }"
          @click="handleNotificationClick(notification)"
        >
          <div class="content">
            <div class="header">
              <h4 class="title">{{ notification.title }}</h4>
              <span class="time">{{ formatMomentTime(notification.created_at) }}</span>
            </div>
            <p v-if="notification.content" class="content-text">
              {{ notification.content }}
            </p>
          </div>
          <div v-if="!notification.is_read" class="unread-indicator" />
        </div>
      </div>

      <div v-else class="empty-state">暂无通知</div>

      <UiPagination
        v-if="total > pageSize"
        :current-page="currentPage"
        :total="total"
        :page-size="pageSize"
        @page-change="handlePageChange"
      />
    </div>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#page {
  @extend .cardHover;
  align-self: flex-start;
  padding: 40px;
  min-height: 500px;

  .page-title {
    margin: 0 0 30px;
    font-weight: bold;
    font-size: 2rem;
    color: var(--font-color);
  }
}

.notification-list {
  .list-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1rem;
    padding-bottom: 0.75rem;
    border-bottom: 1px solid var(--flec-border, #e5e7eb);

    .title {
      margin: 0;
      font-size: 1rem;
      font-weight: 500;
      color: var(--flec-font, #333);

      .count {
        font-size: 0.9rem;
        color: var(--flec-secondary-text, #6b7280);
        font-weight: 400;
      }
    }

    .mark-all-btn {
      padding: 0.4rem 0.8rem;
      background: transparent;
      color: var(--flec-secondary-text, #6b7280);
      border: 1px solid var(--flec-border, #e5e7eb);
      border-radius: 0.25rem;
      font-size: 0.8rem;
      cursor: pointer;
      transition: all 0.15s;

      &:hover {
        border-color: var(--flec-primary, #3b82f6);
        color: var(--flec-primary, #3b82f6);
      }
    }
  }

  .loading {
    padding: 2rem;
    text-align: center;
    font-size: 0.85rem;
    color: var(--flec-secondary-text, #9ca3af);
  }

  .list-content {
    border: 1px solid var(--flec-border, #e5e7eb);
    border-radius: 0.25rem;
    background-color: var(--flec-card-bg, #fff);
    overflow: hidden;
  }

  .empty-state {
    padding: 3rem;
    text-align: center;
    font-size: 0.85rem;
    color: var(--flec-secondary-text, #9ca3af);
  }
}

.notification-item {
  position: relative;
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--flec-border, #e5e7eb);
  cursor: pointer;
  transition: background-color 0.15s ease;

  &:last-child {
    border-bottom: none;
  }

  &:hover {
    background-color: var(--flec-hover-bg, rgba(0, 0, 0, 0.01));
  }

  .unread-indicator {
    position: absolute;
    top: 0.75rem;
    right: 0.75rem;
    width: 6px;
    height: 6px;
    background-color: var(--flec-primary, #3b82f6);
    border-radius: 50%;
  }

  .content {
    flex: 1;
    min-width: 0;

    .header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 0.5rem;
      gap: 0.5rem;

      .title {
        flex: 1;
        font-size: 0.9rem;
        font-weight: 500;
        color: var(--flec-font, #333);
        margin: 0;
      }

      .time {
        font-size: 0.75rem;
        color: var(--flec-secondary-text, #9ca3af);
        white-space: nowrap;
      }
    }

    .content-text {
      font-size: 0.85rem;
      color: var(--flec-secondary-text, #6b7280);
      margin: 0;
      line-height: 1.5;
    }
  }
}

@media (max-width: 768px) {
  #page {
    padding: 20px;

    .page-title {
      font-size: 1.5rem;
      margin-bottom: 20px;
    }
  }
}
</style>
