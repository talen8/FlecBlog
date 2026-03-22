<template>
  <common-list title="友链管理" :data="friendList" :loading="loading" :total="total" v-model:page="queryParams.page"
    v-model:page-size="queryParams.page_size" create-text="新增友链" @create="handleCreate" @refresh="fetchFriends"
    @update:page="fetchFriends" @update:pageSize="fetchFriends">
    <!-- 额外按钮 -->
    <template #toolbar-after>
      <el-button @click="handleTypeManage">
        类型管理
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
        <el-tag v-if="row.is_invalid" type="warning" size="small" style="margin-left: 8px">失效</el-tag>
        <el-tag v-else-if="row.accessible > 0" type="danger" size="small" style="margin-left: 8px">异常({{ row.accessible }})</el-tag>
        <el-tag v-else-if="row.is_pending" type="info" size="small" style="margin-left: 8px">待审核</el-tag>
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
        <span v-if="row.rss_latime" :class="getRSSTimeClass(row.rss_latime)">{{ formatDateTime(row.rss_latime) }}</span>
        <span v-else style="color: #999">-</span>
      </template>
    </el-table-column>

    <el-table-column label="操作" width="180" align="center" fixed="right">
      <template #default="{ row }">
        <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
        <el-button type="danger" link size="small" @click="handleDelete(row.id)">删除</el-button>
      </template>
    </el-table-column>

    <!-- 额外内容 -->
    <template #extra>
      <friend-form-dialog v-model="dialogVisible" :edit-friend="currentFriend" @success="handleFriendSuccess" />
      <friend-type-manager ref="typeManagerRef" v-model="typeManagerVisible" @success="fetchFriends" />
    </template>
  </common-list>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Link } from '@element-plus/icons-vue'
import CommonList from '@/components/common/CommonList.vue'
import type { Friend } from '@/types/friend'
import type { PaginationQuery } from '@/types/request'
import { getFriends, deleteFriend } from '@/api/friend'
import { formatDateTime } from '@/utils/date'
import FriendFormDialog from './components/FriendFormDialog.vue'
import FriendTypeManager from './components/FriendTypeManager.vue'

const loading = ref(false)
const friendList = ref<Friend[]>([])
const total = ref(0)
const queryParams = ref<PaginationQuery>({ page: 1, page_size: 20 })

// 对话框相关
const dialogVisible = ref(false)
const currentFriend = ref<Friend | null>(null)

// 类型管理对话框
const typeManagerVisible = ref(false)
const typeManagerRef = ref<InstanceType<typeof FriendTypeManager>>()

const fetchFriends = async () => {
  loading.value = true
  try {
    const [result] = await Promise.all([
      getFriends(queryParams.value),
      new Promise(resolve => setTimeout(resolve, 300))
    ])
    friendList.value = result.list
    total.value = result.total
  } catch {
    ElMessage.error('获取友链列表失败')
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  currentFriend.value = null
  dialogVisible.value = true
}

const handleEdit = (friend: Friend) => {
  currentFriend.value = friend
  dialogVisible.value = true
}

const handleTypeManage = () => {
  typeManagerVisible.value = true
}

const handleFriendSuccess = () => {
  fetchFriends()
  // 同时刷新类型管理器的数据（更新友链数量）
  typeManagerRef.value?.refreshData()
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个友链吗？', '提示', { type: 'warning' })
    await deleteFriend(id)
    ElMessage.success('删除成功')
    fetchFriends()
    // 同时刷新类型管理器的数据（更新友链数量）
    typeManagerRef.value?.refreshData()
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message)
  }
}

const getRSSTimeClass = (rssLatime?: string): string => {
  if (!rssLatime) return ''
  const months = (Date.now() - new Date(rssLatime).getTime()) / (1000 * 60 * 60 * 24 * 30)
  if (months > 6) return 'rss-danger'
  if (months > 3) return 'rss-warning'
  return ''
}

onMounted(fetchFriends)
</script>

<style scoped>
.rss-warning {
  color: #e6a23c;
}

.rss-danger {
  color: #f56c6c;
}
</style>
