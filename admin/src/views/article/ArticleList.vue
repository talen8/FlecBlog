<template>
  <common-list
    title="文章列表"
    :data="articleList"
    :loading="loading"
    :total="total"
    v-model:page="queryParams.page"
    v-model:page-size="queryParams.page_size"
    create-text="新增文章"
    @create="handleCreate"
    @refresh="fetchArticles"
    @update:page="fetchArticles"
    @update:pageSize="fetchArticles"
  >
    <!-- 额外按钮 -->
    <template #toolbar-after>
      <el-button @click="openCategoryManager"> 分类管理 </el-button>
      <el-button @click="openTagManager"> 标签管理 </el-button>
    </template>

    <!-- 表格列 - 直接使用 el-table-column -->
    <el-table-column label="封面" width="120" align="center">
      <template #default="{ row }">
        <el-image
          :src="row.cover"
          fit="cover"
          style="width: 80px; height: 50px; border-radius: 4px"
        />
      </template>
    </el-table-column>

    <el-table-column label="标题" min-width="300">
      <template #default="{ row }">
        <span>{{ row.title }}</span>
        <el-tag v-if="row.is_top" type="primary" size="small" style="margin-left: 8px">置顶</el-tag>
        <el-tag v-if="row.is_essence" type="success" size="small" style="margin-left: 8px"
          >精选</el-tag
        >
        <el-tag v-if="!row.is_publish" type="warning" size="small" style="margin-left: 8px"
          >草稿</el-tag
        >
      </template>
    </el-table-column>

    <el-table-column label="分类" width="120" align="center">
      <template #default="{ row }">
        <span v-if="row.category">{{ row.category.name }}</span>
        <span v-else style="color: #999">-</span>
      </template>
    </el-table-column>

    <el-table-column label="标签" width="200" align="center">
      <template #default="{ row }">
        <el-tag v-for="tag in row.tags" :key="tag.id" size="small" type="info" style="margin: 2px">
          {{ tag.name }}
        </el-tag>
        <span v-if="!row.tags?.length" style="color: #999">-</span>
      </template>
    </el-table-column>

    <el-table-column label="发布地点" width="120" align="center">
      <template #default="{ row }">
        <span v-if="row.location">{{ row.location }}</span>
        <span v-else style="color: #999">-</span>
      </template>
    </el-table-column>

    <el-table-column label="统计" width="140" align="center">
      <template #default="{ row }">
        <div
          style="
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 12px;
            font-size: 13px;
          "
        >
          <div style="display: flex; align-items: center; gap: 4px">
            <el-icon size="14" style="color: #409eff">
              <View />
            </el-icon>
            <span>{{ row.view_count || 0 }}</span>
          </div>
          <div style="display: flex; align-items: center; gap: 4px">
            <el-icon size="14" style="color: #67c23a">
              <ChatDotRound />
            </el-icon>
            <span>{{ row.comment_count || 0 }}</span>
          </div>
        </div>
      </template>
    </el-table-column>

    <el-table-column label="发布时间" width="180" align="center">
      <template #default="{ row }">
        <div v-if="row.publish_time" style="font-size: 13px; line-height: 1.8">
          <div style="display: flex; align-items: center; justify-content: center; gap: 4px">
            <el-icon size="13" style="color: #67c23a">
              <Upload />
            </el-icon>
            <span>{{ formatDateTime(row.publish_time) }}</span>
          </div>
          <div
            v-if="row.update_time && row.update_time !== row.publish_time"
            style="display: flex; align-items: center; justify-content: center; gap: 4px"
          >
            <el-icon size="13" style="color: #409eff">
              <EditPen />
            </el-icon>
            <span>{{ formatDateTime(row.update_time) }}</span>
          </div>
        </div>
        <span v-else style="color: #999">未发布</span>
      </template>
    </el-table-column>

    <el-table-column label="操作" width="180" align="center" fixed="right">
      <template #default="{ row }">
        <el-button type="primary" link size="small" @click="handleEdit(row.id)">编辑</el-button>
        <el-button type="success" link size="small" @click="openExportDialog(row.id)"
          >导出</el-button
        >
        <el-button type="danger" link size="small" @click="handleDelete(row.id)">删除</el-button>
      </template>
    </el-table-column>
  </common-list>

  <!-- 弹窗组件：懒挂载，首次打开时才渲染 -->
  <category-manager v-if="categoryMounted" v-model="categoryDialogVisible" />
  <tag-manager v-if="tagMounted" v-model="tagDialogVisible" />

  <!-- 导出弹窗 -->
  <el-dialog
    v-model="exportDialogVisible"
    title="导出文章"
    width="480px"
    :close-on-click-modal="false"
  >
    <div class="export-options">
      <div
        v-for="option in exportOptions"
        :key="option.key"
        class="export-option"
        :class="{ disabled: option.loading }"
        @click="handleExport(option.key)"
      >
        <div class="option-icon">
          <i :class="option.icon"></i>
        </div>
        <div class="option-content">
          <div class="option-title">{{ option.title }}</div>
          <div class="option-desc">{{ option.desc }}</div>
        </div>
        <el-icon v-if="option.loading" class="is-loading">
          <Loading />
        </el-icon>
      </div>
    </div>
  </el-dialog>
</template>

<style scoped>
.export-options {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.export-option {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.export-option:hover {
  border-color: #409eff;
  background: #f5f7fa;
}

.export-option.disabled {
  opacity: 0.6;
  pointer-events: none;
}

.option-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f9eb;
  border-radius: 8px;
  font-size: 20px;
  color: #67c23a;
}

.export-option:nth-child(1) .option-icon {
  background: #e6f7e6;
  color: #07c160;
}

.export-option:nth-child(2) .option-icon {
  background: #f4f4f5;
  color: #909399;
}

.option-content {
  flex: 1;
}

.option-title {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 4px;
}

.option-desc {
  font-size: 12px;
  color: #909399;
}
</style>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import { View, ChatDotRound, Upload, EditPen, Loading } from '@element-plus/icons-vue';
import CommonList from '@/components/common/CommonList.vue';
import type { Article } from '@/types/article';
import type { PaginationQuery } from '@/types/request';
import { getArticles, deleteArticle, exportToWeChat, downloadArticleZip } from '@/api/article';
import CategoryManager from './components/CategoryManager.vue';
import TagManager from './components/TagManager.vue';
import { formatDateTime } from '@/utils/date';

const router = useRouter();
const loading = ref(false);
const categoryDialogVisible = ref(false);
const tagDialogVisible = ref(false);
const categoryMounted = ref(false);
const tagMounted = ref(false);
const articleList = ref<Article[]>([]);
const total = ref(0);
const queryParams = ref<PaginationQuery>({ page: 1, page_size: 20 });

const openCategoryManager = () => {
  categoryMounted.value = true;
  categoryDialogVisible.value = true;
};

const openTagManager = () => {
  tagMounted.value = true;
  tagDialogVisible.value = true;
};

const fetchArticles = async () => {
  loading.value = true;
  try {
    const [result] = await Promise.all([
      getArticles(queryParams.value),
      new Promise(resolve => setTimeout(resolve, 300)),
    ]);
    articleList.value = result.list;
    total.value = result.total;
  } catch {
    ElMessage.error('获取文章列表失败');
  } finally {
    loading.value = false;
  }
};

const handleCreate = () => router.push('/articles/create');
const handleEdit = (id: number) => router.push(`/articles/edit/${id}`);

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这篇文章吗？', '提示', {
      type: 'warning',
    });
    await deleteArticle(id);
    ElMessage.success('删除成功');
    fetchArticles();
  } catch (error) {
    if (error instanceof Error) ElMessage.error(error.message);
  }
};

// ==================== 导出功能 ====================

const exportDialogVisible = ref(false);
const exportArticleId = ref<number>(0);

const exportOptions = reactive([
  {
    key: 'wechat',
    title: '复制微信公众号格式',
    desc: '转换为公众号 HTML 并复制到剪贴板',
    icon: 'ri-wechat-line',
    loading: false,
  },
  {
    key: 'markdown',
    title: '下载为 Markdown',
    desc: '下载含图片资源的完整文章',
    icon: 'ri-markdown-line',
    loading: false,
  },
]);

const openExportDialog = (id: number) => {
  exportArticleId.value = id;
  exportDialogVisible.value = true;
};

const handleExport = async (key: string) => {
  const option = exportOptions.find(o => o.key === key);
  if (!option || option.loading) return;

  option.loading = true;

  try {
    switch (key) {
      case 'wechat':
        await handleExportToWeChat();
        break;
      case 'markdown':
        await handleDownloadMarkdown();
        break;
    }
  } finally {
    option.loading = false;
  }
};

// 复制微信公众号格式到剪贴板
const handleExportToWeChat = async () => {
  const result = await exportToWeChat(exportArticleId.value);
  if (result.html) {
    await copyRichText(result.html);
    ElMessage.success('已复制到剪贴板，请粘贴到微信公众平台编辑器');
  } else {
    ElMessage.error('生成失败');
  }
  exportDialogVisible.value = false;
};

// 下载为 Markdown
const handleDownloadMarkdown = async () => {
  let waitingMessage: ReturnType<typeof ElMessage> | undefined;
  const waitingTimer = setTimeout(() => {
    waitingMessage = ElMessage({
      message: '网络较慢或文件资源较大，请耐心等待...',
      type: 'info',
      duration: 0,
    });
  }, 10000);

  try {
    const blob = await downloadArticleZip(exportArticleId.value);
    clearTimeout(waitingTimer);
    waitingMessage?.close();

    const article = articleList.value.find(a => a.id === exportArticleId.value);
    const filename = article ? `${article.title}.zip` : `article-${exportArticleId.value}.zip`;
    downloadBlob(blob, filename);
    ElMessage.success('下载完成');
    exportDialogVisible.value = false;
  } catch (error) {
    clearTimeout(waitingTimer);
    waitingMessage?.close();
    const errorMsg = error instanceof Error ? error.message : '下载失败，请稍后重试';
    ElMessage.error(errorMsg);
  }
};

// 下载 Blob 文件
const downloadBlob = (blob: Blob, filename: string) => {
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
};

// 复制富文本到剪贴板（HTML 格式）
const copyRichText = async (html: string) => {
  try {
    // 使用 Clipboard API 写入富文本
    const blob = new Blob([html], { type: 'text/html' });
    const clipboardItem = new ClipboardItem({
      'text/html': blob,
      'text/plain': new Blob([html], { type: 'text/plain' }),
    });
    await navigator.clipboard.write([clipboardItem]);
  } catch {
    // 降级方案：通过临时元素复制
    const container = document.createElement('div');
    container.innerHTML = html;
    container.style.position = 'fixed';
    container.style.left = '-9999px';
    container.style.whiteSpace = 'pre-wrap';
    document.body.appendChild(container);

    const range = document.createRange();
    range.selectNodeContents(container);
    const selection = window.getSelection();
    selection?.removeAllRanges();
    selection?.addRange(range);

    document.execCommand('copy');
    selection?.removeAllRanges();
    document.body.removeChild(container);
  }
};

onMounted(fetchArticles);
</script>
