<template>
  <el-form label-width="100px" class="setting-form">
    <el-form-item label="文章数据">
      <el-button type="primary" :disabled="readonly" @click="articleImportVisible = true"
        >导入文章</el-button
      >
    </el-form-item>

    <el-form-item label="评论数据">
      <el-button type="primary" :disabled="readonly" @click="commentImportVisible = true"
        >导入评论</el-button
      >
    </el-form-item>
  </el-form>

  <!-- 文章导入对话框 -->
  <el-dialog
    v-model="articleImportVisible"
    title="导入文章"
    width="600px"
    :close-on-click-modal="false"
  >
    <el-form label-width="100px">
      <el-form-item label="数据来源">
        <el-select
          v-model="articleSourceType"
          placeholder="请选择数据来源"
          style="width: 100%"
          :disabled="readonly"
        >
          <el-option label="Hexo 格式" value="hexo" />
          <el-option label="Markdown 格式" value="markdown" />
        </el-select>
        <div class="form-tip">Hexo 格式需要包含 Front Matter，Markdown 格式 仅需 Markdown 内容</div>
      </el-form-item>

      <el-form-item label="上传文件">
        <el-upload
          :auto-upload="false"
          :file-list="articleFileList"
          :on-change="handleArticleFileChange"
          :on-remove="handleArticleFileRemove"
          accept=".md,.markdown"
          :limit="100"
          multiple
          drag
          :disabled="readonly"
        >
          <el-icon class="el-icon--upload"><upload-filled /></el-icon>
          <div class="el-upload__text">拖拽或点击选择文件</div>
          <template #tip>
            <div class="el-upload__tip">最多添加 100 个文件，如遇上传失败请减少数量后重试</div>
          </template>
        </el-upload>
      </el-form-item>

      <el-form-item label="图片处理">
        <el-switch v-model="articleUploadImages" :disabled="readonly" />
        <div class="form-tip" style="margin: 0 15px">开启后将自动下载并上传文章中的图片</div>
      </el-form-item>
    </el-form>

    <el-alert
      v-if="articleImportResult"
      :type="articleImportResult.failed > 0 ? 'warning' : 'success'"
      :closable="false"
      style="margin-top: 16px"
    >
      <div>成功 {{ articleImportResult.success }} 篇，失败 {{ articleImportResult.failed }} 篇</div>
      <div
        v-if="articleImportResult.errors?.length"
        style="margin-top: 8px; font-size: 12px; color: #909399"
      >
        <div v-for="(err, i) in articleImportResult.errors" :key="i">
          {{ err.filename }}: {{ err.error }}
        </div>
      </div>
    </el-alert>

    <template #footer>
      <el-button @click="articleImportVisible = false">取消</el-button>
      <el-button
        type="primary"
        :loading="articleUploading"
        :disabled="readonly || articleFileList.length === 0"
        @click="handleArticleImport"
      >
        {{ articleUploading ? '导入中...' : '开始导入' }}
      </el-button>
    </template>
  </el-dialog>

  <!-- 评论导入对话框 -->
  <el-dialog
    v-model="commentImportVisible"
    title="导入评论"
    width="600px"
    :close-on-click-modal="false"
  >
    <el-form label-width="100px">
      <el-form-item label="数据来源">
        <el-select
          v-model="commentSourceType"
          placeholder="请选择数据来源"
          style="width: 100%"
          :disabled="readonly"
        >
          <el-option label="Artalk" value="artalk" />
        </el-select>
        <div class="form-tip">选择评论数据的来源系统，目前支持 Artalk 评论系统</div>
      </el-form-item>

      <el-form-item label="上传文件">
        <el-upload
          :auto-upload="false"
          :file-list="commentFileList"
          :on-change="handleCommentFileChange"
          :on-remove="handleCommentFileRemove"
          accept=".json,.artrans"
          :limit="1"
          drag
          :disabled="readonly"
        >
          <el-icon class="el-icon--upload"><upload-filled /></el-icon>
          <div class="el-upload__text">拖拽或点击选择文件</div>
          <template #tip>
            <div class="el-upload__tip">支持 JSON 或 Artrans 格式文件，单个文件最大 10MB</div>
          </template>
        </el-upload>
      </el-form-item>
    </el-form>

    <el-alert
      v-if="commentImportResult"
      :type="commentImportResult.failed > 0 ? 'warning' : 'success'"
      :closable="false"
      style="margin-top: 16px"
    >
      <div>
        <strong>导入完成</strong>
      </div>
      <div style="margin-top: 8px">
        总计 {{ commentImportResult.total }} 条，成功 {{ commentImportResult.success }} 条，失败
        {{ commentImportResult.failed }}
        条
      </div>
      <div
        v-if="commentImportResult.user_created > 0"
        style="margin-top: 4px; font-size: 12px; color: #909399"
      >
        自动创建了 {{ commentImportResult.user_created }} 个游客用户账号
      </div>
      <div
        v-if="commentImportResult.errors?.length"
        style="
          margin-top: 12px;
          font-size: 12px;
          color: #909399;
          max-height: 200px;
          overflow-y: auto;
        "
      >
        <div><strong>失败详情：</strong></div>
        <div v-for="(err, i) in commentImportResult.errors" :key="i" style="margin-top: 4px">
          第 {{ err.index + 1 }} 条: {{ err.error }}
        </div>
      </div>
    </el-alert>

    <template #footer>
      <el-button @click="commentImportVisible = false">取消</el-button>
      <el-button
        type="primary"
        :loading="commentUploading"
        :disabled="readonly || commentFileList.length === 0 || !commentSourceType"
        @click="handleCommentImport"
      >
        {{ commentUploading ? '导入中...' : '开始导入' }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { UploadFilled } from '@element-plus/icons-vue';
import type { UploadUserFile, UploadFile } from 'element-plus';
import { importArticles } from '@/api/article';
import { importComments } from '@/api/comment';
import type { ImportArticlesResult } from '@/types/article';
import type { ImportCommentsResult } from '@/types/comment';

const props = withDefaults(defineProps<{ readonly?: boolean }>(), {
  readonly: false,
});

const emit = defineEmits<{
  'import-success': [];
}>();

// 文章导入相关
const articleImportVisible = ref(false);
const articleFileList = ref<UploadUserFile[]>([]);
const articleUploading = ref(false);
const articleImportResult = ref<ImportArticlesResult | undefined>();
const articleSourceType = ref<string>('hexo');
const articleUploadImages = ref(false);

const handleArticleFileChange = (file: UploadFile, files: UploadUserFile[]) => {
  articleFileList.value = files;
};

const handleArticleFileRemove = (file: UploadFile, files: UploadUserFile[]) => {
  articleFileList.value = files;
};

const handleArticleImport = async () => {
  if (props.readonly) {
    ElMessage.warning('仅超级管理员可执行导入操作');
    return;
  }

  if (articleFileList.value.length === 0) return;

  try {
    articleUploading.value = true;
    articleImportResult.value = undefined;

    const formData = new FormData();
    formData.append('source_type', articleSourceType.value);
    formData.append('upload_images', String(articleUploadImages.value));
    articleFileList.value.forEach(file => {
      if (file.raw) formData.append('files', file.raw);
    });

    const result = await importArticles(formData);
    articleImportResult.value = result;

    if (result.failed === 0) {
      ElMessage.success(`成功导入 ${result.success} 篇文章`);
      emit('import-success');
    } else if (result.success > 0) {
      ElMessage.warning(`导入完成：成功 ${result.success} 篇，失败 ${result.failed} 篇`);
      emit('import-success');
    } else {
      ElMessage.error('导入失败');
    }
  } catch (error: unknown) {
    ElMessage.error((error as Error)?.message || '导入失败');
  } finally {
    articleUploading.value = false;
  }
};

watch(articleImportVisible, val => {
  if (!val) {
    setTimeout(() => {
      articleFileList.value = [];
      articleImportResult.value = undefined;
      articleSourceType.value = 'hexo';
      articleUploadImages.value = false;
    }, 300);
  }
});

// 评论导入相关
const commentImportVisible = ref(false);
const commentFileList = ref<UploadUserFile[]>([]);
const commentUploading = ref(false);
const commentImportResult = ref<ImportCommentsResult | undefined>();
const commentSourceType = ref<string>('artalk');

const handleCommentFileChange = (file: UploadFile, files: UploadUserFile[]) => {
  commentFileList.value = files;
};

const handleCommentFileRemove = (file: UploadFile, files: UploadUserFile[]) => {
  commentFileList.value = files;
};

const handleCommentImport = async () => {
  if (props.readonly) {
    ElMessage.warning('仅超级管理员可执行导入操作');
    return;
  }

  if (commentFileList.value.length === 0) {
    ElMessage.warning('请选择要导入的文件');
    return;
  }

  if (!commentSourceType.value) {
    ElMessage.warning('请选择数据来源');
    return;
  }

  try {
    commentUploading.value = true;
    commentImportResult.value = undefined;

    const formData = new FormData();
    const rawFile = commentFileList.value[0]?.raw;
    if (!rawFile) {
      ElMessage.error('文件读取失败');
      return;
    }
    formData.append('file', rawFile);
    formData.append('source_type', commentSourceType.value);

    const result = await importComments(formData);
    commentImportResult.value = result;

    if (result.failed === 0) {
      ElMessage.success(`成功导入 ${result.success} 条评论`);
      emit('import-success');
    } else if (result.success > 0) {
      ElMessage.warning(`导入完成：成功 ${result.success} 条，失败 ${result.failed} 条`);
      emit('import-success');
    } else {
      ElMessage.error('导入失败，请检查文件格式');
    }
  } catch (error: unknown) {
    ElMessage.error((error as Error)?.message || '导入失败');
  } finally {
    commentUploading.value = false;
  }
};

watch(commentImportVisible, val => {
  if (!val) {
    setTimeout(() => {
      commentFileList.value = [];
      commentImportResult.value = undefined;
      commentSourceType.value = 'artalk';
    }, 300);
  }
});
</script>

<style lang="scss" scoped>
:deep(.el-icon--upload) {
  font-size: 40px;
  color: #409eff;
  margin-bottom: 12px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  margin-top: 8px;
}
</style>
