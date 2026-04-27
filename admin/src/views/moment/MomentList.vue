<template>
  <div class="moment-list-page">
    <!-- 筛选控制台 -->
    <transition name="filter-slide">
      <moment-filter
        v-if="showFilter"
        v-model="queryParams"
        @close="showFilter = false"
        @search="fetchMoments"
      />
    </transition>

    <common-list
      title="动态管理"
      :data="momentList"
      :loading="loading"
      :total="total"
      v-model:page="queryParams.page"
      v-model:page-size="queryParams.page_size"
      create-text="新增动态"
      :filter-active="showFilter"
      :filter-count="activeFilterCount"
      @create="handleCreate"
      @refresh="fetchMoments"
      @filter="toggleFilter"
      @update:page="fetchMoments"
      @update:pageSize="fetchMoments"
    >
      <!-- 快速筛选 -->
      <template #toolbar-before>
        <template v-if="!showFilter">
          <el-input
            v-model="quickFilters.keyword"
            placeholder="搜索关键词"
            clearable
            class="quick-filter-900"
            style="width: 180px"
            @keyup.enter="handleQuickFilterChange"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-select
            v-model="quickFilters.is_publish"
            placeholder="发布状态"
            clearable
            class="quick-filter-769"
            style="width: 110px"
            @change="handleQuickFilterChange"
          >
            <el-option label="已发布" :value="true" />
            <el-option label="草稿" :value="false" />
          </el-select>
        </template>
      </template>

      <el-table-column label="内容" min-width="400">
        <template #default="{ row }">
          <div class="moment-content">
            <!-- 文本内容 -->
            <div v-if="row.content.text" class="text-content">
              {{ row.content.text }}
            </div>

            <!-- 图片 -->
            <div v-if="row.content.images?.length" class="images-content">
              <el-image
                v-for="(image, index) in row.content.images.slice(0, 3)"
                :key="index"
                :src="image"
                fit="cover"
                style="width: 60px; height: 60px; border-radius: 4px; margin-right: 8px"
              />
              <span v-if="row.content.images.length > 3" class="more-images">
                +{{ row.content.images.length - 3 }}
              </span>
            </div>

            <!-- 所有标签（标签、视频、音频、音乐、链接、位置） -->
            <div
              v-if="
                row.content.tags ||
                row.content.video ||
                row.content.audio ||
                row.content.music ||
                row.content.link ||
                row.content.location
              "
              class="tags-container"
            >
              <!-- 标签 -->
              <el-tag v-if="row.content.tags" size="small" type="info">
                {{ row.content.tags }}
              </el-tag>

              <!-- 视频 -->
              <el-tag v-if="row.content.video" type="primary" size="small">
                <i class="ri-video-line"></i>
                {{ getVideoPlatformName(row.content.video.platform) }}
              </el-tag>

              <!-- 音频 -->
              <el-tag v-if="row.content.audio" type="success" size="small">
                <i class="ri-mic-line"></i>
                音频
              </el-tag>

              <!-- 音乐 -->
              <el-tag v-if="row.content.music" type="success" size="small">
                <i class="ri-music-line"></i>
                {{ getMusicLabel(row.content.music) }}
              </el-tag>

              <!-- 链接 -->
              <el-tag v-if="row.content.link" size="small" type="warning">
                <i class="ri-link"></i>
                {{ row.content.link.title || row.content.link.url }}
              </el-tag>

              <!-- 位置 -->
              <el-tag v-if="row.content.location" type="danger" size="small">
                <i class="ri-map-pin-line"></i>
                {{ row.content.location }}
              </el-tag>
            </div>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_publish ? 'success' : 'warning'" size="small">
            {{ row.is_publish ? '已发布' : '草稿' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="发布时间" width="180" align="center">
        <template #default="{ row }">
          <div v-if="row.publish_time">
            {{ formatDateTime(row.publish_time) }}
          </div>
          <span v-else style="color: #999">-</span>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="180" align="center" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="handleEdit(row.id)">编辑</el-button>
          <el-button type="danger" link size="small" @click="handleDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </common-list>
  </div>

  <!-- 动态表单弹窗 -->
  <moment-form-dialog
    v-model="momentDialogVisible"
    :edit-moment="editingMoment"
    @success="handleDialogSuccess"
  />
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Search } from '@element-plus/icons-vue';
import CommonList from '@/components/common/CommonList.vue';
import MomentFilter from './components/MomentFilter.vue';
import MomentFormDialog from './components/MomentFormDialog.vue';
import type { Moment, MomentListQuery } from '@/types/moment';
import { getMoments, deleteMoment } from '@/api/moment';
import { formatDateTime } from '@/utils/date';

const loading = ref(false);
const momentList = ref<Moment[]>([]);
const total = ref(0);
const showFilter = ref(false);
const queryParams = ref<MomentListQuery>({
  page: 1,
  page_size: 20,
});

// 快速筛选相关
const quickFilters = reactive({
  keyword: '',
  is_publish: undefined as boolean | undefined,
});

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
      fetchMoments();
    }, 500);
  }
);

const momentDialogVisible = ref(false);
const editingMoment = ref<Moment | null>(null);

// 音乐平台和类型映射
const MUSIC_LABELS = {
  type: {
    song: '单曲',
    album: '专辑',
    artist: '艺术家',
    playlist: '歌单',
  },
  server: {
    netease: '网易云',
    tencent: 'QQ音乐',
  },
};

/**
 * 计算当前激活的筛选项数量
 */
const activeFilterCount = computed(() => {
  let count = 0;
  if (queryParams.value.keyword) count++;
  if (queryParams.value.tags) count++;
  if (queryParams.value.location) count++;
  if (queryParams.value.is_publish !== undefined) count++;
  if (queryParams.value.has_images !== undefined) count++;
  if (queryParams.value.has_video !== undefined) count++;
  if (queryParams.value.has_music !== undefined) count++;
  if (queryParams.value.has_link !== undefined) count++;
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
  quickFilters.keyword = queryParams.value.keyword || '';
  quickFilters.is_publish = queryParams.value.is_publish;
};

/**
 * 处理快速筛选变化
 */
const handleQuickFilterChange = () => {
  // 将快速筛选条件同步到查询参数
  queryParams.value.keyword = quickFilters.keyword || undefined;
  queryParams.value.is_publish = quickFilters.is_publish;
  queryParams.value.page = 1;
  fetchMoments();
};

// 获取视频平台名称
const getVideoPlatformName = (platform?: string) => {
  if (!platform) return '本地视频';
  const platformMap: Record<string, string> = {
    bilibili: '哔哩哔哩',
    youtube: 'YouTube',
  };
  return platformMap[platform.toLowerCase()] || '本地视频';
};

// 获取音乐标签
const getMusicLabel = (music: { server: string; type: string }) => {
  const serverName =
    MUSIC_LABELS.server[music.server as keyof typeof MUSIC_LABELS.server] || music.server;
  const typeName = MUSIC_LABELS.type[music.type as keyof typeof MUSIC_LABELS.type] || music.type;
  return `${serverName} - ${typeName}`;
};

let errorMessageShown = false;

/**
 * 获取动态列表
 */
const fetchMoments = async () => {
  loading.value = true;
  try {
    const [result] = await Promise.all([
      getMoments(queryParams.value),
      new Promise(resolve => setTimeout(resolve, 300)),
    ]);
    momentList.value = result.list;
    total.value = result.total;
  } catch {
    if (!errorMessageShown) {
      errorMessageShown = true;
      ElMessage.error('获取动态列表失败');
      // 3秒后重置标记，允许再次提示
      setTimeout(() => {
        errorMessageShown = false;
      }, 3000);
    }
  } finally {
    loading.value = false;
  }
};

const handleCreate = () => {
  editingMoment.value = null;
  momentDialogVisible.value = true;
};

const handleEdit = (id: number) => {
  const moment = momentList.value.find(item => item.id === id);
  if (moment) {
    editingMoment.value = moment;
    momentDialogVisible.value = true;
  }
};

const handleDialogSuccess = () => {
  fetchMoments();
};

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这条动态吗？', '提示', {
      type: 'warning',
    });
    await deleteMoment(id);
    ElMessage.success('删除成功');
    fetchMoments();
  } catch (error) {
    if (error instanceof Error) ElMessage.error(error.message);
  }
};

onMounted(() => {
  // 初始化快速筛选值（从 queryParams）
  syncQuickFiltersFromQueryParams();
  fetchMoments();
});
</script>

<style scoped lang="scss">
.moment-list-page {
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

.moment-list-page > :deep(.filter-panel) {
  flex-shrink: 0;
}

.moment-list-page > :deep(.common-list) {
  flex: 1;
  min-height: 0;
}

.moment-content {
  .text-content {
    margin-bottom: 8px;
    line-height: 1.5;
    color: #333;
  }

  .images-content {
    display: flex;
    align-items: center;
    margin-bottom: 8px;

    .more-images {
      color: #666;
      font-size: 12px;
    }
  }

  .tags-container {
    display: flex;
    align-items: center;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 4px;
  }
}
</style>
