<template>
  <div id="page">
    <div class="page-header">
      <h1 class="page-title">通知中心</h1>
      <button
        v-if="notifications.length > 0 && unreadCount > 0"
        class="mark-all-btn"
        :disabled="markingAll"
        @click="handleMarkAllRead"
      >
        {{ markingAll ? '处理中...' : '全部已读' }}
      </button>
    </div>

    <div v-if="notifications.length > 0" class="notification-list">
      <div
        v-for="notification in notifications"
        :key="notification.id"
        class="notification-item"
        :class="{ unread: !notification.is_read }"
      >
        <NuxtLink
          :to="notification.link"
          class="notification-link"
          @click="handleClick(notification)"
        >
          <div class="notification-content">
            <h3>{{ notification.title }}</h3>
            <p>{{ notification.content }}</p>
            <span class="notification-time">{{ formatMomentTime(notification.created_at) }}</span>
          </div>
        </NuxtLink>
        <button
          v-if="!notification.is_read"
          class="mark-read-btn"
          title="标记已读"
          @click.stop="handleMarkRead(notification.id)"
        >
          <i class="ri-check-line" />
        </button>
      </div>
    </div>
    <div v-else class="empty-state">暂无通知</div>
    <UiPagination
      v-if="total > pageSize"
      :total="total"
      :current-page="currentPage"
      :page-size="pageSize"
      @change="handlePageChange"
    />
  </div>
</template>

<script lang="ts" setup>
const isLoggedIn = useAuth();
const router = useRouter();
const {
  notifications,
  total,
  currentPage,
  pageSize,
  unreadCount,
  fetchNotifications,
  markNotificationAsRead,
  markAllNotificationsAsRead,
} = useNotifications();

const markingAll = ref(false);

onMounted(() => {
  if (!isLoggedIn.value) {
    router.push('/');
    return;
  }
  fetchNotifications();
});

const handlePageChange = (page: number) => {
  fetchNotifications({ page });
  if (import.meta.client) {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
};

const handleClick = async (notification: { id: number; is_read: boolean }) => {
  if (!notification.is_read) {
    await markNotificationAsRead(notification.id);
  }
};

const handleMarkRead = async (id: number) => {
  try {
    await markNotificationAsRead(id);
  } catch (error) {
    console.error('标记已读失败:', error);
  }
};

const handleMarkAllRead = async () => {
  markingAll.value = true;
  try {
    await markAllNotificationsAsRead();
  } catch (error) {
    console.error('全部标记已读失败:', error);
  } finally {
    markingAll.value = false;
  }
};

useSeoMeta({
  title: '通知中心',
});
</script>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border-radius: 12px;
  border: 1px solid var(--flec-border);
  padding: 40px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-title {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
}

.mark-all-btn {
  padding: 8px 16px;
  border: 1px solid var(--flec-border);
  border-radius: 6px;
  background: var(--flec-card-bg);
  color: var(--font-color);
  font-size: 0.85rem;
  cursor: pointer;

  &:hover:not(:disabled) {
    border-color: var(--theme-color);
    color: var(--theme-color);
  }

  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.notification-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.notification-item {
  display: flex;
  align-items: center;
  border-radius: 8px;
  background: var(--flec-heavy-bg);

  &.unread {
    background: rgba(73, 177, 245, 0.1);
  }
}

.notification-link {
  flex: 1;
  display: block;
  padding: 16px;
  text-decoration: none;
}

.notification-content {
  h3 {
    margin: 0 0 8px;
    font-size: 1rem;
    font-weight: 500;
    color: var(--font-color);
  }

  p {
    margin: 0 0 8px;
    font-size: 0.9rem;
    color: var(--theme-meta-color);
  }
}

.notification-time {
  font-size: 0.85rem;
  color: var(--theme-meta-color);
}

.mark-read-btn {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  margin-right: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: var(--theme-meta-color);
  font-size: 1rem;
  cursor: pointer;

  &:hover {
    background: var(--flec-hover-bg);
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
}
</style>
