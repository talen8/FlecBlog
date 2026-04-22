<template>
  <div class="filter-panel">
    <!-- 桌面端：使用卡片形式 -->
    <el-card v-if="!isMobile" class="filter-card" shadow="never">
      <div class="filter-header">
        <span class="filter-title">
          <el-icon><Filter /></el-icon>
          {{ title }}
        </span>
        <div class="filter-actions">
          <el-button link type="primary" @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
          <el-button link @click="handleClose">
            <el-icon><Close /></el-icon>
            收起
          </el-button>
        </div>
      </div>

      <el-form :model="formData" label-position="top">
        <el-row :gutter="12">
          <slot />
        </el-row>
      </el-form>
    </el-card>

    <!-- 移动端：使用 Dialog 形式 -->
    <el-dialog
      v-else
      :model-value="true"
      :title="title"
      width="90%"
      style="max-width: 350px"
      :close-on-click-modal="true"
      :show-close="false"
      class="filter-dialog"
      append-to-body
      destroy-on-close
      align-center
      @close="handleClose"
    >
      <div class="dialog-content" @click.stop>
        <el-form :model="formData" label-position="top" size="small">
          <div class="mobile-form-items">
            <slot />
          </div>
        </el-form>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="handleReset">重置</el-button>
          <el-button type="primary" size="small" @click="handleClose">关闭</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue';
import { Filter, Refresh, Close } from '@element-plus/icons-vue';

/**
 * 组件属性定义
 */
const props = withDefaults(
  defineProps<{
    /** 筛选面板标题 */
    title?: string;
    /** 表单数据 */
    modelValue: Record<string, unknown>;
  }>(),
  {
    title: '筛选条件',
  }
);

/**
 * 组件事件定义
 */
const emit = defineEmits<{
  'update:modelValue': [value: Record<string, unknown>];
  reset: [];
  close: [];
}>();

const formData = ref<Record<string, unknown>>({ ...props.modelValue });
const isMobile = ref(false);

const checkMobile = () => {
  isMobile.value = window.innerWidth < 600;
};

onMounted(() => {
  checkMobile();
  window.addEventListener('resize', checkMobile);
});

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile);
});

watch(
  () => props.modelValue,
  newVal => {
    formData.value = { ...newVal };
  },
  { deep: true }
);

/**
 * 处理重置
 */
const handleReset = () => {
  const pageSize = formData.value.page_size;
  formData.value = {
    page: 1,
    ...(pageSize !== undefined && { page_size: pageSize }),
  };
  emit('update:modelValue', { ...formData.value });
  emit('reset');
};

/**
 * 处理关闭
 */
const handleClose = () => {
  emit('close');
};
</script>

<style scoped lang="scss">
.filter-panel {
  margin-bottom: 16px;

  .filter-card {
    :deep(.el-card__body) {
      padding: 14px 20px;
    }
  }

  .filter-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid var(--el-border-color-lighter);

    .filter-title {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 15px;
      font-weight: 500;
      color: var(--el-text-color-primary);
    }

    .filter-actions {
      display: flex;
      gap: 8px;
    }
  }
}

.filter-dialog {
  :deep(.el-dialog) {
    max-height: 80vh;
    display: flex;
    flex-direction: column;
  }

  :deep(.el-dialog__body) {
    padding: 12px 16px;
    overflow-y: auto;
    flex: 1;
  }

  .dialog-content {
    max-height: calc(80vh - 140px);
    overflow-y: auto;

    :deep(.el-form-item) {
      margin-bottom: 12px;
    }

    :deep(.el-form-item__label) {
      font-size: 12px;
      padding-bottom: 4px;
    }

    :deep(.el-col) {
      width: 100% !important;
      max-width: 100% !important;
      flex: 0 0 100% !important;
    }
  }

  .mobile-form-items {
    :deep(.el-col) {
      width: 100% !important;
      max-width: 100% !important;
      flex: 0 0 100% !important;
    }

    :deep(.el-select) {
      width: 100% !important;
    }

    :deep(.el-date-editor) {
      width: 100% !important;
    }
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    flex-shrink: 0;

    .el-button {
      margin-left: 0;
    }
  }
}

// 移动端适配
@media (max-width: 767px) {
  .filter-panel {
    margin-bottom: 0;

    .filter-card {
      :deep(.el-card__body) {
        padding: 12px;
      }
    }

    .filter-header {
      .filter-title {
        font-size: 14px;
      }
    }
  }
}
</style>
