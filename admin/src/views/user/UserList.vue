<template>
  <div class="user-list-page">
    <!-- 筛选面板 -->
    <transition name="filter-slide">
      <user-filter
        v-if="showFilter"
        v-model="queryParams"
        @search="fetchUsers"
        @close="showFilter = false"
      />
    </transition>

    <common-list
      title="用户管理"
      :data="userList"
      :loading="loading"
      :total="total"
      :filter-active="showFilter"
      :filter-count="filterCount"
      v-model:page="queryParams.page"
      v-model:page-size="queryParams.page_size"
      create-text="新增用户"
      @create="handleCreate"
      @refresh="fetchUsers"
      @filter="toggleFilter"
      @update:page="fetchUsers"
      @update:pageSize="fetchUsers"
    >
      <!-- 快速筛选 -->
      <template #toolbar-before>
        <template v-if="!showFilter">
          <el-select
            v-model="quickFilters.role"
            placeholder="全部角色"
            clearable
            class="quick-filter-769"
            style="width: 130px"
            @change="handleQuickFilterChange"
          >
            <el-option label="超级管理员" value="super_admin" />
            <el-option label="管理员" value="admin" />
            <el-option label="普通用户" value="user" />
            <el-option label="访客" value="guest" />
          </el-select>
          <el-select
            v-model="quickFilters.is_enabled"
            placeholder="启用状态"
            clearable
            class="quick-filter-840"
            style="width: 110px"
            @change="handleQuickFilterChange"
          >
            <el-option label="启用" :value="true" />
            <el-option label="禁用" :value="false" />
          </el-select>
        </template>
      </template>

      <!-- 表格列 -->
      <el-table-column label="头像" width="100" align="center">
        <template #default="{ row }">
          <el-avatar v-if="row.avatar" :src="row.avatar" :size="40" />
          <el-avatar v-else :size="40">
            <el-icon>
              <User />
            </el-icon>
          </el-avatar>
        </template>
      </el-table-column>

      <el-table-column label="昵称" min-width="130" align="center">
        <template #default="{ row }">
          <div style="display: flex; align-items: center; justify-content: center; gap: 8px">
            <span>{{ row.nickname }}</span>
            <el-tag v-if="row.badge" type="info" effect="plain" size="small">{{
              row.badge
            }}</el-tag>
            <el-tag v-if="row.deleted_at" type="danger" size="small">已删除</el-tag>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="邮箱" min-width="150" align="center">
        <template #default="{ row }">
          <div style="display: flex; align-items: center; justify-content: center; gap: 8px">
            <span v-if="row.email">{{ row.email }}</span>
            <span v-else style="color: #999">-</span>
            <el-tag v-if="!row.is_enabled" type="danger" size="small">禁用</el-tag>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="网站地址" min-width="150" align="center">
        <template #default="{ row }">
          <span v-if="row.website">{{ row.website }}</span>
          <span v-else style="color: #999">-</span>
        </template>
      </el-table-column>

      <el-table-column label="角色" width="120" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.role === 'super_admin'" type="danger" size="small">超级管理员</el-tag>
          <el-tag v-else-if="row.role === 'admin'" type="warning" size="small">管理员</el-tag>
          <el-tag v-else-if="row.role === 'user'" type="success" size="small">普通用户</el-tag>
          <el-tag v-else type="info" size="small">访客</el-tag>
        </template>
      </el-table-column>

      <el-table-column label="登录方式" width="150" align="center">
        <template #default="{ row }">
          <div class="login-methods">
            <template v-if="row.has_password">
              <i class="ri-lock-password-fill"></i>
            </template>

            <template v-if="row.github_id">
              <i class="ri-github-fill"></i>
            </template>

            <template v-if="row.google_id">
              <i class="ri-google-fill"></i>
            </template>

            <template v-if="row.qq_id">
              <i class="ri-qq-fill"></i>
            </template>

            <template v-if="row.microsoft_id">
              <i class="ri-microsoft-fill"></i>
            </template>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="最后登录" width="180" align="center">
        <template #default="{ row }">
          {{ formatDateTime(row.last_login) }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="180" align="center" fixed="right">
        <template #default="{ row }">
          <template v-if="!row.deleted_at && canOperateUser(row)">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row.id)"
              >删除</el-button
            >
          </template>
        </template>
      </el-table-column>

      <!-- 额外内容 -->
      <template #extra>
        <!-- 用户表单对话框 -->
        <user-form-dialog v-model="dialogVisible" :edit-user="currentUser" @success="fetchUsers" />
      </template>
    </common-list>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, reactive } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { User } from '@element-plus/icons-vue';
import CommonList from '@/components/common/CommonList.vue';
import UserFilter from './components/UserFilter.vue';
import type { User as UserType, UserListQuery } from '@/types/user';
import { getUsers, deleteUser } from '@/api/user';
import UserFormDialog from './components/UserFormDialog.vue';
import { formatDateTime } from '@/utils/date';
import { getCurrentUserRole } from '@/utils/auth';

const loading = ref(false);
const userList = ref<UserType[]>([]);
const total = ref(0);
const showFilter = ref(false);
const queryParams = ref<UserListQuery>({ page: 1, page_size: 20 });
const currentRole = computed(() => getCurrentUserRole());

// 快速筛选相关
const quickFilters = reactive({
  role: undefined as string | undefined,
  is_enabled: undefined as boolean | undefined,
});

/**
 * 计算当前应用的筛选条件数量
 */
const filterCount = computed(() => {
  let count = 0;
  if (queryParams.value.keyword) count++;
  if (queryParams.value.role) count++;
  if (queryParams.value.is_enabled !== undefined) count++;
  if (queryParams.value.is_deleted !== undefined) count++;
  if (queryParams.value.login_method) count++;
  if (queryParams.value.last_login_start && queryParams.value.last_login_end) count++;
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
  quickFilters.role = queryParams.value.role;
  quickFilters.is_enabled = queryParams.value.is_enabled;
};

/**
 * 处理快速筛选变化
 */
const handleQuickFilterChange = () => {
  // 将快速筛选条件同步到查询参数
  queryParams.value.role = quickFilters.role;
  queryParams.value.is_enabled = quickFilters.is_enabled;
  // 重置到第一页并搜索
  queryParams.value.page = 1;
  fetchUsers();
};

// 对话框相关
const dialogVisible = ref(false);
const currentUser = ref<UserType | null>(null);

const isManagedRole = (role: string) => role === 'admin' || role === 'super_admin';
const canOperateUser = (user: UserType) =>
  currentRole.value === 'super_admin' || !isManagedRole(user.role);

const fetchUsers = async () => {
  loading.value = true;
  try {
    const [result] = await Promise.all([
      getUsers(queryParams.value),
      new Promise(resolve => setTimeout(resolve, 300)),
    ]);
    userList.value = result.list;
    total.value = result.total;
  } catch {
    ElMessage.error('获取用户列表失败');
  } finally {
    loading.value = false;
  }
};

const handleCreate = () => {
  currentUser.value = null;
  dialogVisible.value = true;
};

const handleEdit = (user: UserType) => {
  if (!canOperateUser(user)) return;
  currentUser.value = user;
  dialogVisible.value = true;
};

const handleDelete = async (id: number) => {
  const target = userList.value.find(user => user.id === id);
  if (target && !canOperateUser(target)) return;

  try {
    await ElMessageBox.confirm('确定要删除这个用户吗？', '提示', {
      type: 'warning',
    });
    await deleteUser(id);
    ElMessage.success('删除成功');
    fetchUsers();
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message);
  }
};

onMounted(() => {
  // 初始化快速筛选值（从 queryParams）
  syncQuickFiltersFromQueryParams();
  fetchUsers();
});
</script>

<style scoped lang="scss">
.user-list-page {
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

.user-list-page > :deep(.filter-panel) {
  flex-shrink: 0;
}

.user-list-page > :deep(.common-list) {
  flex: 1;
  min-height: 0;
}

.login-methods {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
}

.login-methods i {
  font-size: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
