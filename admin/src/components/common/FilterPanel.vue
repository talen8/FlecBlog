<template>
  <div class="filter-panel">
    <el-card class="filter-card" shadow="never">
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
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
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

// 移动端适配
@media (max-width: 767px) {
  .filter-panel {
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
