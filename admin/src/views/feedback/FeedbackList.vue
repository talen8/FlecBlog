<template>
  <div class="feedback-list-page">
    <!-- 筛选面板 -->
    <transition name="filter-slide">
      <feedback-filter
        v-if="showFilter"
        v-model="queryParams"
        @search="fetchList"
        @close="showFilter = false"
      />
    </transition>

    <common-list
      title="反馈投诉"
      :data="list"
      :loading="loading"
      :total="total"
      :show-create="false"
      :filter-active="showFilter"
      :filter-count="filterCount"
      v-model:page="queryParams.page"
      v-model:page-size="queryParams.page_size"
      @refresh="fetchList"
      @filter="toggleFilter"
      @update:page="fetchList"
      @update:pageSize="fetchList"
    >
      <!-- 快速筛选 -->
      <template #toolbar-before>
        <template v-if="!showFilter">
          <el-select
            v-model="quickFilters.report_type"
            placeholder="反馈类型"
            clearable
            style="width: 150px; margin-right: 8px"
            @change="handleQuickFilterChange"
          >
            <el-option label="版权侵权" value="copyright" />
            <el-option label="不当内容" value="inappropriate" />
            <el-option label="摘要问题" value="summary" />
            <el-option label="功能建议" value="suggestion" />
          </el-select>
          <el-select
            v-model="quickFilters.status"
            placeholder="状态"
            clearable
            style="width: 110px; margin-right: 12px"
            @change="handleQuickFilterChange"
          >
            <el-option label="待处理" value="pending" />
            <el-option label="已解决" value="resolved" />
            <el-option label="已关闭" value="closed" />
          </el-select>
        </template>
      </template>

      <el-table-column label="工单号" width="150" align="center">
        <template #default="{ row }">
          <span>{{ row.ticket_no }}</span>
        </template>
      </el-table-column>

      <el-table-column label="类型" width="200" align="center">
        <template #default="{ row }">
          <el-tag :type="getReportTypeTagType(row.report_type)">
            {{ getReportTypeLabel(row.report_type) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="投诉地址" min-width="200" align="center">
        <template #default="{ row }">
          <span :title="row.report_url">
            {{ truncateUrl(row.report_url) }}
          </span>
        </template>
      </el-table-column>

      <el-table-column label="联系方式" width="200" align="center">
        <template #default="{ row }">
          <span v-if="row.email">{{ row.email }}</span>
        </template>
      </el-table-column>

      <el-table-column label="状态" width="120" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusTagType(row.status)">
            {{ getStatusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="反馈时间" width="180" align="center">
        <template #default="{ row }">
          {{ formatDateTime(row.feedback_time) }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="180" align="center" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" size="small" text @click="handleView(row.id)">
            查看详情
          </el-button>
          <el-button type="danger" size="small" text @click="handleDelete(row.id)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </common-list>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import CommonList from '@/components/common/CommonList.vue';
import FeedbackFilter from './components/FeedbackFilter.vue';
import { getFeedbackList, deleteFeedback } from '@/api/feedback';
import type { Feedback, FeedbackStatus, FeedbackListQuery, ReportType } from '@/types/feedback';
import { formatDateTime } from '@/utils/date';

const router = useRouter();
const loading = ref(false);
const list = ref<Feedback[]>([]);
const total = ref(0);
const showFilter = ref(false);
const queryParams = ref<FeedbackListQuery>({ page: 1, page_size: 10 });

// 快速筛选相关
const quickFilters = reactive({
  report_type: undefined as ReportType | undefined,
  status: undefined as FeedbackStatus | undefined,
});

/**
 * 计算当前应用的筛选条件数量
 */
const filterCount = computed(() => {
  let count = 0;
  if (queryParams.value.keyword) count++;
  if (queryParams.value.report_type) count++;
  if (queryParams.value.status) count++;
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
  quickFilters.report_type = queryParams.value.report_type;
  quickFilters.status = queryParams.value.status;
};

/**
 * 处理快速筛选变化
 */
const handleQuickFilterChange = () => {
  // 将快速筛选条件同步到查询参数
  queryParams.value.report_type = quickFilters.report_type;
  queryParams.value.status = quickFilters.status;
  // 重置到第一页并搜索
  queryParams.value.page = 1;
  fetchList();
};

// 获取反馈列表
const fetchList = async () => {
  loading.value = true;
  try {
    const res = await getFeedbackList(queryParams.value);
    list.value = res.list || [];
    total.value = res.total || 0;
  } catch (_error) {
    ElMessage.error('获取反馈列表失败');
    list.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
};

const handleView = (id: number) => {
  router.push(`/feedback/${id}`);
};

// 删除反馈
const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除此反馈吗？', '提示', {
      type: 'warning',
    });
    await deleteFeedback(id);
    ElMessage.success('删除成功');
    fetchList();
  } catch (error: unknown) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败');
    }
  }
};

const getReportTypeLabel = (reportType: string) => {
  const labels: Record<string, string> = {
    copyright: '版权侵权内容投诉',
    inappropriate: '不当内容举报投诉',
    summary: '文章摘要问题反馈',
    suggestion: '功能建议优化反馈',
  };
  return labels[reportType] || reportType;
};

// Element Plus 标签类型
type TagType = 'success' | 'warning' | 'danger' | 'info';

const getReportTypeTagType = (reportType: string): TagType => {
  const types: Record<string, TagType> = {
    copyright: 'warning',
    inappropriate: 'danger',
    summary: 'info',
    suggestion: 'success',
  };
  return types[reportType] || 'info';
};

const getStatusLabel = (status: FeedbackStatus) => {
  const labels: Record<FeedbackStatus, string> = {
    pending: '待处理',
    resolved: '已解决',
    closed: '已关闭',
  };
  return labels[status] || status;
};

const getStatusTagType = (status: FeedbackStatus): TagType => {
  const types: Record<FeedbackStatus, TagType> = {
    pending: 'warning',
    resolved: 'success',
    closed: 'info',
  };
  return types[status] || 'info';
};

const truncateUrl = (url: string) => {
  if (!url) return '无地址';
  if (url.length <= 60) return url;
  return url.substring(0, 60) + '...';
};

onMounted(() => {
  // 初始化快速筛选值（从 queryParams）
  syncQuickFiltersFromQueryParams();
  fetchList();
});
</script>

<style scoped lang="scss">
.feedback-list-page {
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

.feedback-list-page > :deep(.filter-panel) {
  flex-shrink: 0;
}

.feedback-list-page > :deep(.common-list) {
  flex: 1;
  min-height: 0;
}
</style>
