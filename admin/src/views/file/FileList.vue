<template>
  <div class="file-list-page">
    <transition name="filter-slide">
      <file-filter
        v-if="showFilter"
        v-model="query"
        @close="showFilter = false"
        @search="loadList"
      />
    </transition>

    <common-list
      title="文件管理"
      :data="fileList"
      :loading="loading"
      :total="total"
      :show-create="false"
      :filter-active="showFilter"
      :filter-count="activeFilterCount"
      v-model:page="query.page"
      v-model:page-size="query.page_size"
      @refresh="loadList"
      @filter="toggleFilter"
      @update:page="loadList"
      @update:pageSize="loadList"
    >
      <template #toolbar-before>
        <template v-if="!showFilter">
          <el-select
            v-model="quickFilters.file_type"
            placeholder="全部类型"
            clearable
            style="width: 120px; margin-right: 8px"
            @change="handleQuickFilterChange"
          >
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
            <el-option label="音频" value="audio" />
            <el-option label="文档" value="application" />
          </el-select>
          <el-select
            v-model="quickFilters.status"
            placeholder="使用状态"
            clearable
            style="width: 100px; margin-right: 12px"
            @change="handleQuickFilterChange"
          >
            <el-option label="使用中" :value="1" />
            <el-option label="未使用" :value="0" />
          </el-select>
        </template>
      </template>

      <el-table-column label="预览" width="80" align="center">
        <template #default="{ row }">
          <el-image
            v-if="isImage(row)"
            :src="row.file_url"
            fit="cover"
            style="width: 50px; height: 50px; border-radius: 4px"
          />
        </template>
      </el-table-column>

      <el-table-column label="文件名" min-width="180">
        <template #default="{ row }">
          <span style="margin-right: 8px; font-weight: 500">{{ row.file_name }}</span>
          <span style="font-size: 12px; color: #909399">{{ formatFileSize(row.file_size) }}</span>
        </template>
      </el-table-column>

      <el-table-column
        prop="original_name"
        label="原始文件名"
        min-width="200"
        show-overflow-tooltip
      />

      <el-table-column prop="file_type" label="类型" width="120" align="center" />

      <el-table-column label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getStatusTagType(row.status)" size="small" effect="light">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="upload_type" label="用途" width="100" align="center" />

      <el-table-column label="上传时间" width="180" align="center">
        <template #default="{ row }">
          {{ formatDateTime(row.upload_time) }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="180" align="center" fixed="right">
        <template #default="{ row }">
          <el-button link type="primary" size="small" @click="copyUrl(row)">复制链接</el-button>
          <el-button link type="danger" size="small" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </common-list>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import CommonList from '@/components/common/CommonList.vue';
import FileFilter from './components/FileFilter.vue';
import { getFileList, deleteFile } from '@/api/file';
import type { FileInfo, FileListQuery } from '@/types/file';
import { formatDateTime } from '@/utils/date';

const query = ref<FileListQuery>({ page: 1, page_size: 20 });
const fileList = ref<FileInfo[]>([]);
const total = ref(0);
const loading = ref(false);
const showFilter = ref(false);

const quickFilters = reactive({
  file_type: undefined as string | undefined,
  status: undefined as number | undefined,
});

/**
 * 计算当前激活的筛选项数量
 */
const activeFilterCount = computed(() => {
  let count = 0;
  if (query.value.keyword) count++;
  if (query.value.file_type !== undefined) count++;
  if (query.value.status !== undefined) count++;
  if (query.value.upload_type) count++;
  if (query.value.min_size || query.value.max_size) count++;
  if (query.value.start_time || query.value.end_time) count++;
  return count;
});

/**
 * 切换筛选面板显示状态
 */
const toggleFilter = () => {
  showFilter.value = !showFilter.value;
  if (!showFilter.value) {
    syncQuickFiltersFromQuery();
  }
};

/**
 * 从 query 同步筛选条件到快速筛选
 */
const syncQuickFiltersFromQuery = () => {
  quickFilters.file_type = query.value.file_type;
  quickFilters.status = query.value.status;
};

/**
 * 处理快速筛选变化
 */
const handleQuickFilterChange = () => {
  query.value.file_type = quickFilters.file_type;
  query.value.status = quickFilters.status;
  query.value.page = 1;
  loadList();
};

const loadList = async () => {
  loading.value = true;
  try {
    const [data] = await Promise.all([
      getFileList(query.value),
      new Promise(resolve => setTimeout(resolve, 300)),
    ]);
    fileList.value = data.list;
    total.value = data.total;
  } catch {
    ElMessage.error('加载失败');
  } finally {
    loading.value = false;
  }
};

const copyUrl = async (file: FileInfo) => {
  try {
    await navigator.clipboard.writeText(file.file_url);
    ElMessage.success('已复制');
  } catch {
    ElMessage.error('复制失败');
  }
};

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这个文件吗？', '提示', {
      type: 'warning',
    });
    await deleteFile(id);
    ElMessage.success('删除成功');
    loadList();
  } catch (error) {
    if (error !== 'cancel' && error instanceof Error) ElMessage.error(error.message);
  }
};

const isImage = (file: FileInfo) => file.file_type?.startsWith('image/');

const formatFileSize = (size: number) => {
  if (size < 1024) return size + ' B';
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB';
  return (size / (1024 * 1024)).toFixed(1) + ' MB';
};

const getStatusTagType = (status: number) => {
  return status === 1 ? 'success' : 'info';
};

const getStatusText = (status: number) => {
  return status === 1 ? '使用中' : '未使用';
};

onMounted(() => {
  syncQuickFiltersFromQuery();
  loadList();
});
</script>

<style scoped>
.file-list-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

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

.file-list-page > :deep(.filter-panel) {
  flex-shrink: 0;
}

.file-list-page > :deep(.common-list) {
  flex: 1;
  min-height: 0;
}
</style>
