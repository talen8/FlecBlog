<template>
  <div class="visit-list-page">
    <!-- 筛选控制台 -->
    <transition name="filter-slide">
      <visit-filter
        v-if="showFilter"
        v-model="queryParams"
        @close="showFilter = false"
        @search="fetchVisits"
      />
    </transition>

    <common-list
      title="访问日志"
      :data="visitList"
      :loading="loading"
      :total="total"
      :show-create="false"
      :filter-active="showFilter"
      :filter-count="activeFilterCount"
      v-model:page="queryParams.page"
      v-model:page-size="queryParams.page_size"
      @refresh="fetchVisits"
      @filter="toggleFilter"
      @update:page="fetchVisits"
      @update:pageSize="fetchVisits"
    >
      <!-- 快速筛选 -->
      <template #toolbar-before>
        <template v-if="!showFilter">
          <el-input
            v-model="quickFilters.ip"
            placeholder="筛选 IP 地址"
            clearable
            class="quick-filter-769"
            style="width: 180px; margin-right: 12px"
            @keyup.enter="handleQuickFilterChange"
          >
            <template #prefix>
              <el-icon><Connection /></el-icon>
            </template>
          </el-input>
        </template>
      </template>

      <el-table-column label="访客ID" width="150" align="center">
        <template #default="{ row }">
          <el-tooltip :content="row.visitor_id" placement="top">
            <div style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap">
              {{ row.visitor_id.substring(0, 12) }}...
            </div>
          </el-tooltip>
        </template>
      </el-table-column>

      <el-table-column label="IP地址" width="140" align="center" prop="ip" />

      <el-table-column label="访问页面" min-width="250">
        <template #default="{ row }">
          <div style="display: flex; align-items: center; gap: 8px; width: 100%">
            <el-tooltip :content="row.page_url" placement="top">
              <div style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap; flex: 1">
                {{ row.page_url }}
              </div>
            </el-tooltip>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="地理位置" width="150" align="center" prop="location" />

      <el-table-column label="浏览器" width="140" align="center" prop="browser" />

      <el-table-column label="操作系统" width="120" align="center" prop="os" />

      <el-table-column label="来源" width="250" align="center">
        <template #default="{ row }">
          <el-tooltip v-if="row.referer" :content="row.referer" placement="top">
            <div style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap">
              {{ row.referer }}
            </div>
          </el-tooltip>
          <span v-else style="color: #999">-</span>
        </template>
      </el-table-column>

      <el-table-column label="访问时间" width="180" align="center">
        <template #default="{ row }">
          {{ formatDateTime(row.created_at) }}
        </template>
      </el-table-column>
    </common-list>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { Connection } from '@element-plus/icons-vue';
import CommonList from '@/components/common/CommonList.vue';
import VisitFilter from './components/VisitFilter.vue';
import type { Visit, VisitListQuery } from '@/types/stats';
import { getVisits } from '@/api/stats';
import { formatDateTime } from '@/utils/date';

const loading = ref(false);
const visitList = ref<Visit[]>([]);
const total = ref(0);
const showFilter = ref(false);
const queryParams = ref<VisitListQuery>({
  page: 1,
  page_size: 20,
});

// 快速筛选相关
const quickFilters = reactive({
  ip: '',
});

// 搜索防抖定时器
let searchTimer: ReturnType<typeof setTimeout> | null = null;

// 监听 IP 变化，实时搜索
watch(
  () => quickFilters.ip,
  newVal => {
    if (searchTimer) clearTimeout(searchTimer);
    searchTimer = setTimeout(() => {
      queryParams.value.ip = newVal || undefined;
      queryParams.value.page = 1;
      fetchVisits();
    }, 500);
  }
);

/**
 * 计算当前激活的筛选项数量
 */
const activeFilterCount = computed(() => {
  let count = 0;
  if (queryParams.value.keyword) count++;
  if (queryParams.value.visitor_id) count++;
  if (queryParams.value.ip) count++;
  if (queryParams.value.exclude_ips) count++;
  if (queryParams.value.location) count++;
  if (queryParams.value.browser) count++;
  if (queryParams.value.os) count++;
  if (queryParams.value.start_time || queryParams.value.end_time) count++;
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
  quickFilters.ip = queryParams.value.ip || '';
};

/**
 * 处理快速筛选变化
 */
const handleQuickFilterChange = () => {
  // 将快速筛选条件同步到查询参数
  queryParams.value.ip = quickFilters.ip || undefined;
  queryParams.value.page = 1;
  fetchVisits();
};

let errorMessageShown = false;

/**
 * 获取访问日志列表
 */
const fetchVisits = async () => {
  loading.value = true;
  try {
    const [result] = await Promise.all([
      getVisits(queryParams.value),
      new Promise(resolve => setTimeout(resolve, 300)),
    ]);
    visitList.value = result.list;
    total.value = result.total;
  } catch {
    if (!errorMessageShown) {
      errorMessageShown = true;
      ElMessage.error('获取访问日志失败');
      // 3秒后重置标记，允许再次提示
      setTimeout(() => {
        errorMessageShown = false;
      }, 3000);
    }
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  // 初始化快速筛选值（从 queryParams）
  syncQuickFiltersFromQueryParams();
  fetchVisits();
});
</script>

<style scoped lang="scss">
.visit-list-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* 筛选控制台滑入滑出动画 */
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

.visit-list-page > :deep(.filter-panel) {
  flex-shrink: 0;
}

.visit-list-page > :deep(.common-list) {
  flex: 1;
  min-height: 0;
}
</style>
