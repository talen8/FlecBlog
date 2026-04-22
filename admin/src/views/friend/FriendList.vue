<template>
  <div class="friend-list-page">
    <!-- 筛选控制台 -->
    <transition name="filter-slide">
      <friend-filter
        v-if="showFilter"
        v-model="queryParams"
        @close="showFilter = false"
        @search="fetchFriends"
      />
    </transition>

    <common-list
      title="友链管理"
      :data="friendList"
      :loading="loading"
      :total="total"
      v-model:page="queryParams.page"
      v-model:page-size="queryParams.page_size"
      create-text="新增友链"
      :filter-active="showFilter"
      :filter-count="activeFilterCount"
      @create="handleCreate"
      @refresh="fetchFriends"
      @filter="toggleFilter"
      @update:page="fetchFriends"
      @update:pageSize="fetchFriends"
    >
      <!-- 快速筛选 -->
      <template #toolbar-before>
        <template v-if="!showFilter">
          <el-input
            v-model="quickFilters.keyword"
            placeholder="搜索关键词"
            clearable
            class="quick-filter-1080"
            style="width: 160px"
            @keyup.enter="handleQuickFilterChange"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select
            v-model="quickFilters.type_id"
            placeholder="友链类型"
            clearable
            class="quick-filter-1200"
            style="width: 130px"
            @change="handleQuickFilterChange"
          >
            <el-option
              v-for="type in friendTypes"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
          <el-select
            v-model="quickFilters.accessible_status"
            placeholder="访问状态"
            clearable
            class="quick-filter-800"
            style="width: 120px"
            @change="handleQuickFilterChange"
          >
            <el-option label="正常" value="normal" />
            <el-option label="异常" value="abnormal" />
            <el-option label="忽略检查" value="ignored" />
          </el-select>
        </template>
        <el-button class="icon-btn" @click="handleTypeManage">
          <el-icon><Folder /></el-icon><span class="btn-text">类型管理</span>
        </el-button>
      </template>

      <!-- 表格列 -->
      <el-table-column label="头像" width="80" align="center">
        <template #default="{ row }">
          <el-avatar v-if="row.avatar" :src="row.avatar" :size="40" />
          <el-avatar v-else :size="40">
            <el-icon>
              <Link />
            </el-icon>
          </el-avatar>
        </template>
      </el-table-column>

      <el-table-column label="友链名称" min-width="130">
        <template #default="{ row }">
          <span>{{ row.name }}</span>
          <el-tag v-if="row.is_invalid" type="warning" size="small" style="margin-left: 8px"
            >失效</el-tag
          >
          <el-tag v-else-if="row.accessible > 0" type="danger" size="small" style="margin-left: 8px"
            >异常({{ row.accessible }})</el-tag
          >
          <el-tag v-else-if="row.is_pending" type="info" size="small" style="margin-left: 8px"
            >待审核</el-tag
          >
        </template>
      </el-table-column>

      <el-table-column label="链接地址" min-width="160">
        <template #default="{ row }">
          <span>{{ row.url }}</span>
        </template>
      </el-table-column>

      <el-table-column label="描述" min-width="180">
        <template #default="{ row }">
          <span v-if="row.description">{{ row.description }}</span>
          <span v-else style="color: #999">-</span>
        </template>
      </el-table-column>

      <el-table-column label="类型" width="120" align="center">
        <template #default="{ row }">
          <span v-if="row.type_name">{{ row.type_name }}</span>
          <span v-else style="color: #999">-</span>
        </template>
      </el-table-column>

      <el-table-column prop="sort" label="排序" width="100" align="center" />

      <el-table-column label="最新文章" width="180" align="center">
        <template #default="{ row }">
          <span v-if="row.rss_latime" :class="getRSSTimeClass(row.rss_latime)">{{
            formatDateTime(row.rss_latime)
          }}</span>
          <span v-else style="color: #999">-</span>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="180" align="center" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button type="danger" link size="small" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </common-list>
  </div>

  <!-- 弹窗组件：懒挂载，首次打开时才渲染 -->
  <friend-form-dialog
    v-if="formMounted"
    v-model="dialogVisible"
    :edit-friend="currentFriend"
    @success="handleFriendSuccess"
  />
  <friend-type-manager
    v-if="typeManagerMounted"
    ref="typeManagerRef"
    v-model="typeManagerVisible"
    @success="fetchFriends"
  />
</template>

<script setup lang="ts">
import { ref, onMounted, reactive, computed, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Link, Search, Folder } from '@element-plus/icons-vue';
import CommonList from '@/components/common/CommonList.vue';
import FriendFilter from './components/FriendFilter.vue';
import type { Friend, FriendQuery, FriendType } from '@/types/friend';
import { getFriends, deleteFriend, getFriendTypes } from '@/api/friend';
import { formatDateTime } from '@/utils/date';
import FriendFormDialog from './components/FriendFormDialog.vue';
import FriendTypeManager from './components/FriendTypeManager.vue';

const loading = ref(false);
const friendList = ref<Friend[]>([]);
const total = ref(0);
const showFilter = ref(false);
const queryParams = ref<FriendQuery>({
  page: 1,
  page_size: 20,
});

// 快速筛选相关
const quickFilters = reactive({
  keyword: '',
  type_id: undefined as number | undefined,
  accessible_status: undefined as 'normal' | 'abnormal' | 'ignored' | undefined,
});

// 友链类型列表
const friendTypes = ref<FriendType[]>([]);

// 搜索防抖定时器
let searchTimer: ReturnType<typeof setTimeout> | null = null;

// 监听关键词变化，实时搜索
watch(
  () => quickFilters.keyword,
  newVal => {
    if (searchTimer) clearTimeout(searchTimer);
    searchTimer = setTimeout(() => {
      queryParams.value.keyword = newVal || undefined;
      queryParams.value.page = 1;
      fetchFriends();
    }, 500);
  }
);

// 加载友链类型
const loadFriendTypes = async () => {
  try {
    const result = await getFriendTypes();
    friendTypes.value = result.list || [];
  } catch {
    ElMessage.error('加载友链类型失败');
  }
};

// 对话框相关
const dialogVisible = ref(false);
const currentFriend = ref<Friend | null>(null);
const formMounted = ref(false);

// 类型管理对话框
const typeManagerVisible = ref(false);
const typeManagerMounted = ref(false);
const typeManagerRef = ref<InstanceType<typeof FriendTypeManager>>();

/**
 * 计算当前激活的筛选项数量
 */
const activeFilterCount = computed(() => {
  let count = 0;
  if (queryParams.value.keyword) count++;
  if (queryParams.value.type_id !== undefined) count++;
  if (queryParams.value.is_invalid !== undefined) count++;
  if (queryParams.value.is_pending !== undefined) count++;
  if (queryParams.value.accessible_status) count++;
  if (queryParams.value.rss_status) count++;
  if (queryParams.value.has_screenshot !== undefined) count++;
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
  quickFilters.keyword = queryParams.value.keyword || '';
  quickFilters.type_id = queryParams.value.type_id;
  quickFilters.accessible_status = queryParams.value.accessible_status;
};

/**
 * 处理快速筛选变化
 */
const handleQuickFilterChange = () => {
  // 将快速筛选条件同步到查询参数
  queryParams.value.keyword = quickFilters.keyword || undefined;
  queryParams.value.type_id = quickFilters.type_id;
  queryParams.value.accessible_status = quickFilters.accessible_status;
  // 重置到第一页并搜索
  queryParams.value.page = 1;
  fetchFriends();
};

const fetchFriends = async () => {
  loading.value = true;
  try {
    const [result] = await Promise.all([
      getFriends(queryParams.value),
      new Promise(resolve => setTimeout(resolve, 300)),
    ]);
    friendList.value = result.list;
    total.value = result.total;
  } catch {
    ElMessage.error('获取友链列表失败');
  } finally {
    loading.value = false;
  }
};

const handleCreate = () => {
  currentFriend.value = null;
  formMounted.value = true;
  dialogVisible.value = true;
};

const handleEdit = (friend: Friend) => {
  currentFriend.value = friend;
  formMounted.value = true;
  dialogVisible.value = true;
};

const handleTypeManage = () => {
  typeManagerMounted.value = true;
  typeManagerVisible.value = true;
};

const handleFriendSuccess = () => {
  fetchFriends();
  // 同时刷新类型管理器的数据（更新友链数量）
  typeManagerRef.value?.refreshData();
};

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个友链吗？', '提示', {
      type: 'warning',
    });
    await deleteFriend(id);
    ElMessage.success('删除成功');
    fetchFriends();
    // 同时刷新类型管理器的数据（更新友链数量）
    typeManagerRef.value?.refreshData();
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message);
  }
};

const getRSSTimeClass = (rssLatime?: string): string => {
  if (!rssLatime) return '';
  const months = (Date.now() - new Date(rssLatime).getTime()) / (1000 * 60 * 60 * 24 * 30);
  if (months > 6) return 'rss-danger';
  if (months > 3) return 'rss-warning';
  return '';
};

onMounted(() => {
  loadFriendTypes();
  // 初始化快速筛选值（从 queryParams）
  syncQuickFiltersFromQueryParams();
  fetchFriends();
});
</script>

<style scoped lang="scss">
.friend-list-page {
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

.friend-list-page > :deep(.filter-panel) {
  flex-shrink: 0;
}

.friend-list-page > :deep(.common-list) {
  flex: 1;
  min-height: 0;
}

.rss-warning {
  color: #e6a23c;
}

.rss-danger {
  color: #f56c6c;
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
