<template>
  <div class="rss-feed-list-page">
    <!-- 筛选面板 -->
    <transition name="filter-slide">
      <rss-feed-filter
        v-if="showFilter"
        v-model="queryParams"
        @search="fetchArticles"
        @close="showFilter = false"
      />
    </transition>

    <common-list
      title="RSS订阅"
      :data="articleList"
      :loading="loading"
      :total="total"
      v-model:page="queryParams.page"
      v-model:page-size="queryParams.page_size"
      :show-create="false"
      :filter-active="showFilter"
      :filter-count="filterCount"
      @refresh="fetchArticles"
      @filter="toggleFilter"
      @update:page="fetchArticles"
      @update:pageSize="fetchArticles"
    >
      <template #toolbar-before>
        <template v-if="!showFilter">
          <el-select
            v-model="quickFilters.friend_id"
            placeholder="全部友链"
            clearable
            class="quick-filter-960"
            style="width: 150px"
            @change="handleQuickFilterChange"
          >
            <el-option
              v-for="friend in friendList"
              :key="friend.id"
              :label="friend.name"
              :value="friend.id"
            />
          </el-select>
          <el-select
            v-model="quickFilters.is_read"
            placeholder="阅读状态"
            clearable
            class="quick-filter-800"
            style="width: 110px"
            @change="handleQuickFilterChange"
          >
            <el-option label="已读" :value="true" />
            <el-option label="未读" :value="false" />
          </el-select>
        </template>
        <el-button class="icon-btn" @click="openSubscriberDialog">
          <el-icon><Bell /></el-icon><span class="btn-text">本站订阅</span>
        </el-button>
        <el-badge :value="unreadCount" :hidden="unreadCount === 0" :max="99" class="unread-badge">
          <el-button
            class="icon-btn"
            :disabled="unreadCount === 0"
            @click="handleMarkAllRead"
            v-if="isSuperAdmin"
            type="primary"
          >
            <el-icon><Check /></el-icon><span class="btn-text">全部已读</span>
          </el-button>
        </el-badge>
      </template>

      <!-- 表格列 -->
      <el-table-column label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_read ? 'info' : 'danger'" size="small">
            {{ row.is_read ? '已读' : '未读' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="文章标题" min-width="300">
        <template #default="{ row }">
          <a :href="row.link" target="_blank" class="article-link" :class="{ read: row.is_read }">
            {{ row.title }}
          </a>
        </template>
      </el-table-column>

      <el-table-column label="来源" width="180">
        <template #default="{ row }">
          <a :href="row.friend_url" target="_blank" class="friend-link">
            {{ row.friend_name }}
          </a>
        </template>
      </el-table-column>

      <el-table-column label="发布时间" width="180" align="center">
        <template #default="{ row }">
          <span v-if="row.published_at">{{ formatDateTime(row.published_at) }}</span>
          <span v-else style="color: #999">-</span>
        </template>
      </el-table-column>

      <el-table-column label="抓取时间" width="180" align="center">
        <template #default="{ row }">
          {{ formatDateTime(row.created_at) }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="120" align="center" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="!row.is_read && isSuperAdmin"
            type="primary"
            link
            size="small"
            @click="handleMarkRead(row)"
          >
            标记已读
          </el-button>
          <span v-else style="color: #999">-</span>
        </template>
      </el-table-column>
    </common-list>
  </div>

  <!-- 本站订阅弹窗 -->
  <el-dialog
    v-model="subscriberDialogVisible"
    title="本站订阅者"
    width="90%"
    style="max-width: 700px"
    :align-center="true"
    destroy-on-close
  >
    <el-table :data="subscriberList" v-loading="subscriberLoading" style="width: 100%">
      <el-table-column label="邮箱地址" min-width="160">
        <template #default="{ row }">
          <div style="display: flex; align-items: center; gap: 8px">
            <el-icon size="16" color="#409eff">
              <Message />
            </el-icon>
            <span>{{ row.email }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.active ? 'success' : 'info'" size="small">
            {{ row.active ? '活跃' : '已退订' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="订阅时间" width="100" align="center">
        <template #default="{ row }">
          {{ formatDateTime(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="80" align="center" fixed="right">
        <template #default="{ row }">
          <el-button type="danger" link size="small" @click="handleDeleteSubscriber(row.id)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <div style="margin-top: 16px; display: flex; justify-content: flex-end">
      <el-pagination
        v-model:current-page="subscriberQuery.page"
        v-model:page-size="subscriberQuery.page_size"
        :total="subscriberTotal"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @current-change="fetchSubscribers"
        @size-change="fetchSubscribers"
      />
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Message, Bell, Check } from '@element-plus/icons-vue';
import CommonList from '@/components/common/CommonList.vue';
import RssFeedFilter from './components/RssFeedFilter.vue';
import type { RssArticle, RssArticleQuery } from '@/types/rssfeed';
import type { Subscriber } from '@/types/subscriber';
import type { Friend } from '@/types/friend';
import { getRssArticles, markRssArticleRead, markAllRssArticlesRead } from '@/api/rssfeed';
import { getSubscribers, deleteSubscriber } from '@/api/subscriber';
import { getFriends } from '@/api/friend';
import { formatDateTime } from '@/utils/date';
import { isSuperAdmin as checkSuperAdmin } from '@/utils/auth';

const isSuperAdmin = computed(() => checkSuperAdmin());

const loading = ref(false);
const articleList = ref<RssArticle[]>([]);
const total = ref(0);
const unreadCount = ref(0);
const showFilter = ref(false);
const queryParams = ref<RssArticleQuery>({ page: 1, page_size: 20 });

// 快速筛选相关
const friendList = ref<Friend[]>([]);
const quickFilters = reactive({
  friend_id: undefined as number | undefined,
  is_read: undefined as boolean | undefined,
});

/**
 * 计算当前应用的筛选条件数量
 */
const filterCount = computed(() => {
  let count = 0;
  if (queryParams.value.keyword) count++;
  if (queryParams.value.friend_id !== undefined) count++;
  if (queryParams.value.is_read !== undefined) count++;
  if (queryParams.value.start_time && queryParams.value.end_time) count++;
  return count;
});

/**
 * 切换筛选面板显示状态
 */
const toggleFilter = () => {
  showFilter.value = !showFilter.value;
  if (!showFilter.value) {
    syncQuickFiltersFromQueryParams();
  }
};

/**
 * 从 queryParams 同步筛选条件到快速筛选
 */
const syncQuickFiltersFromQueryParams = () => {
  quickFilters.friend_id = queryParams.value.friend_id;
  quickFilters.is_read = queryParams.value.is_read;
};

/**
 * 加载友链列表（用于快速筛选）
 */
const loadFriendList = async () => {
  try {
    const result = await getFriends({ page: 1, page_size: 1000 });
    friendList.value = result.list || [];
  } catch (error) {
    console.error('加载友链列表失败:', error);
  }
};

/**
 * 处理快速筛选变化
 */
const handleQuickFilterChange = () => {
  // 将快速筛选条件同步到查询参数
  queryParams.value.friend_id = quickFilters.friend_id;
  queryParams.value.is_read = quickFilters.is_read;
  // 重置到第一页并搜索
  queryParams.value.page = 1;
  fetchArticles();
};

/**
 * 获取RSS文章列表
 */
const fetchArticles = async () => {
  loading.value = true;
  try {
    const result = await getRssArticles(queryParams.value);
    articleList.value = result.list;
    total.value = result.total;
    unreadCount.value = result.unread_count;
  } catch {
    ElMessage.error('获取RSS文章列表失败');
  } finally {
    loading.value = false;
  }
};

/**
 * 标记单篇文章已读
 */
const handleMarkRead = async (article: RssArticle) => {
  try {
    await markRssArticleRead(article.id);
    article.is_read = true;
    unreadCount.value = Math.max(0, unreadCount.value - 1);
    ElMessage.success('已标记为已读');
  } catch {
    ElMessage.error('操作失败');
  }
};

/**
 * 全部标记已读
 */
const handleMarkAllRead = async () => {
  try {
    await ElMessageBox.confirm('确定要将所有未读文章标记为已读吗？', '提示', {
      type: 'warning',
    });
    const result = await markAllRssArticlesRead();
    ElMessage.success(`已标记 ${result.affected} 篇文章为已读`);
    fetchArticles();
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) {
      ElMessage.error(error.message);
    }
  }
};

// 订阅者相关
const subscriberDialogVisible = ref(false);
const subscriberLoading = ref(false);
const subscriberList = ref<Subscriber[]>([]);
const subscriberTotal = ref(0);
const subscriberQuery = reactive({ page: 1, page_size: 10 });

/**
 * 打开订阅者弹窗
 */
const openSubscriberDialog = () => {
  subscriberDialogVisible.value = true;
  subscriberQuery.page = 1;
  fetchSubscribers();
};

/**
 * 获取订阅者列表
 */
const fetchSubscribers = async () => {
  subscriberLoading.value = true;
  try {
    const result = await getSubscribers(subscriberQuery);
    subscriberList.value = result.list;
    subscriberTotal.value = result.total;
  } catch {
    ElMessage.error('获取订阅者列表失败');
  } finally {
    subscriberLoading.value = false;
  }
};

/**
 * 删除订阅者
 */
const handleDeleteSubscriber = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除该订阅者吗？此操作不可恢复。', '提示', {
      type: 'warning',
    });
    await deleteSubscriber(id);
    ElMessage.success('删除成功');
    fetchSubscribers();
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message);
  }
};

onMounted(() => {
  loadFriendList();
  // 初始化快速筛选值（从 queryParams）
  syncQuickFiltersFromQueryParams();
  fetchArticles();
});
</script>

<style scoped lang="scss">
.rss-feed-list-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 筛选面板滑入滑出动画 */
.filter-slide-enter-active,
.filter-slide-leave-active {
  transition: all 0.1s linear;
}

.filter-slide-enter-from,
.filter-slide-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

.filter-slide-enter-to,
.filter-slide-leave-from {
  opacity: 1;
  transform: translateY(0);
}

.rss-feed-list-page > :deep(.filter-panel) {
  flex-shrink: 0;
}

.rss-feed-list-page > :deep(.common-list) {
  flex: 1;
  min-height: 0;
}

.article-link {
  color: var(--el-color-primary);
  text-decoration: none;
  transition: color 0.2s;
}

.article-link:hover {
  text-decoration: underline;
}

.article-link.read {
  color: var(--el-text-color-regular);
}

.friend-link {
  color: var(--el-text-color-secondary);
  text-decoration: none;
}

.friend-link:hover {
  color: var(--el-color-primary);
  text-decoration: underline;
}

// 默认显示文字，移动端显示图标
.icon-btn {
  .el-icon {
    display: none;
  }
  // 覆盖 Element Plus 默认样式，消除图标隐藏后的左边距
  .btn-text {
    margin-left: 0;
  }
}

@media (max-width: 500px) {
  .icon-btn {
    .btn-text {
      display: none;
    }
    .el-icon {
      display: inline-flex;
    }
  }
}
</style>
