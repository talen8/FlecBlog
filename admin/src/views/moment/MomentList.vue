<template>
  <common-list
    title="动态列表"
    :data="momentList"
    :loading="loading"
    :total="total"
    v-model:page="queryParams.page"
    v-model:page-size="queryParams.page_size"
    create-text="新增动态"
    @create="handleCreate"
    @refresh="fetchMoments"
    @update:page="fetchMoments"
    @update:pageSize="fetchMoments"
  >
    <!-- 表格列 -->
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

          <!-- 所有标签（标签、视频、音乐、链接、位置） -->
          <div
            v-if="
              row.content.tags ||
              row.content.video ||
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

  <!-- 动态表单弹窗 -->
  <moment-form-dialog
    v-model="momentDialogVisible"
    :edit-moment="editingMoment"
    @success="handleDialogSuccess"
  />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import CommonList from '@/components/common/CommonList.vue';
import MomentFormDialog from './components/MomentFormDialog.vue';
import type { Moment } from '@/types/moment';
import type { PaginationQuery } from '@/types/request';
import { getMoments, deleteMoment } from '@/api/moment';
import { formatDateTime } from '@/utils/date';

const loading = ref(false);
const momentList = ref<Moment[]>([]);
const total = ref(0);
const queryParams = ref<PaginationQuery>({ page: 1, page_size: 20 });
const momentDialogVisible = ref(false);
const editingMoment = ref<Moment | null>(null);

// 音乐平台和类型映射
const MUSIC_LABELS = {
  type: {
    search: '搜索',
    song: '单曲',
    album: '专辑',
    artist: '艺术家',
    playlist: '歌单',
  },
  server: {
    netease: '网易云',
    tencent: 'QQ音乐',
    kugou: '酷狗',
    xiami: '虾米',
    baidu: '百度',
    kuwo: '酷我',
  },
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
    ElMessage.error('获取动态列表失败');
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

onMounted(fetchMoments);
</script>

<style scoped lang="scss">
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
