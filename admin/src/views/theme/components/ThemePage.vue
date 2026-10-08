<template>
  <div class="theme-page-panel">
    <el-table :data="rows" border class="page-table">
      <el-table-column label="路径" min-width="200">
        <template #default="{ row }">
          <span>{{ (row as PageRow).path }}</span>
        </template>
      </el-table-column>

      <el-table-column label="页面标题" min-width="200">
        <template #default="{ row }">
          <span v-if="(row as PageRow).item?.title">{{ (row as PageRow).item?.title }}</span>
          <span v-else class="empty-value"></span>
        </template>
      </el-table-column>

      <el-table-column label="页面描述" min-width="280">
        <template #default="{ row }">
          <span v-if="(row as PageRow).item?.description" class="desc-text">
            {{ (row as PageRow).item?.description }}
          </span>
          <span v-else class="empty-value"></span>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="180" align="center">
        <template #default="{ row }">
          <el-button
            type="primary"
            link
            size="small"
            :disabled="disabled"
            @click="handleEdit(row as PageRow)"
          >
            编辑
          </el-button>
          <el-button
            type="danger"
            link
            size="small"
            :disabled="disabled || !(row as PageRow).customized"
            @click="handleReset(row as PageRow)"
          >
            恢复默认
          </el-button>
        </template>
      </el-table-column>

      <template #empty>
        <el-empty description="暂无可管理的页面" />
      </template>
    </el-table>

    <el-dialog
      v-model="dialogVisible"
      title="编辑页面文案"
      width="90%"
      style="max-width: 600px"
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="90px">
        <div class="form-info">
          <div class="info-item">
            <span class="info-label">页面路径</span>
            <span class="info-value">{{ editingPath }}</span>
          </div>
        </div>

        <el-form-item label="页面标题" prop="title">
          <el-input
            v-model="formData.title"
            placeholder="留空则使用默认标题"
            maxlength="100"
            show-word-limit
            clearable
          />
        </el-form-item>

        <el-form-item label="页面描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="留空则使用默认描述，用于页面 SEO 描述"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { FormInstance, FormRules } from 'element-plus';
import type { ThemePageItem, ThemeSchema } from '@/types/theme';
import { updateThemePages } from '@/api/theme';

interface PageRow {
  path: string;
  item?: ThemePageItem;
  customized: boolean;
}

const props = withDefaults(
  defineProps<{
    schema?: ThemeSchema | Record<string, unknown>;
    pages?: Record<string, ThemePageItem>;
    disabled?: boolean;
  }>(),
  {
    schema: () => ({}),
    pages: () => ({}),
    disabled: false,
  }
);

const emit = defineEmits<{
  refresh: [];
}>();

const pagePaths = computed(() => (props.schema as ThemeSchema)?.$pages || []);

const dialogVisible = ref(false);
const submitLoading = ref(false);
const formRef = ref<FormInstance>();
const editingPath = ref('');
const formData = ref({ title: '', description: '' });

const rules: FormRules = {
  title: [{ max: 100, message: '标题不能超过 100 个字符', trigger: 'blur' }],
  description: [{ max: 500, message: '描述不能超过 500 个字符', trigger: 'blur' }],
};

const rows = computed<PageRow[]>(() =>
  pagePaths.value.map(path => {
    const item = props.pages?.[path];
    return {
      path,
      item,
      customized: Boolean(item?.title || item?.description),
    };
  })
);

const savePages = async (path: string, item: ThemePageItem | null) => {
  const nextPages: Record<string, ThemePageItem> = { ...props.pages };
  if (item) {
    nextPages[path] = item;
  } else {
    delete nextPages[path];
  }
  await updateThemePages(nextPages);
  emit('refresh');
};

const handleEdit = (row: PageRow) => {
  editingPath.value = row.path;
  formData.value = {
    title: row.item?.title ?? '',
    description: row.item?.description ?? '',
  };
  dialogVisible.value = true;
};

const handleReset = async (row: PageRow) => {
  try {
    await ElMessageBox.confirm(`确定将「${row.path}」恢复为默认文案吗？`, '提示', {
      type: 'warning',
    });
    await savePages(row.path, null);
    ElMessage.success('已恢复默认');
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error(error.message || '操作失败');
    }
  }
};

const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    submitLoading.value = true;

    const title = formData.value.title.trim();
    const description = formData.value.description.trim();
    if (!title && !description) {
      await savePages(editingPath.value, null);
      ElMessage.success('已恢复默认');
    } else {
      await savePages(editingPath.value, { title, description });
      ElMessage.success('保存成功');
    }

    dialogVisible.value = false;
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error(error.message || '操作失败');
    }
  } finally {
    submitLoading.value = false;
  }
};
</script>

<style scoped lang="scss">
.theme-page-panel {
  .page-table {
    width: 100%;
  }

  .desc-text {
    display: -webkit-box;
    -webkit-box-orient: vertical;
    line-clamp: 2;
    overflow: hidden;
  }

  .empty-value {
    color: #909399;
  }
}

.form-info {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 16px;
  margin-bottom: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
  border: 1px solid #e4e7ed;

  .info-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    flex: 1;
    text-align: center;

    .info-label {
      font-size: 12px;
      color: #909399;
    }

    .info-value {
      font-size: 14px;
      color: #303133;
      font-weight: 500;
    }
  }
}

@media (max-width: 768px) {
  .theme-page-panel {
    .search-input {
      width: 100%;
    }
  }
}
</style>
