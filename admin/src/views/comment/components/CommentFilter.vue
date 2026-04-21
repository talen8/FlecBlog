<template>
  <filter-panel v-model="filterForm" title="筛选条件" @reset="handleReset" @close="$emit('close')">
    <!-- 第一行：关键词、显示状态、删除状态、评论类型、评论时间 -->
    <el-col :span="5">
      <el-form-item label="关键词">
        <el-input v-model="filterForm.keyword" placeholder="搜索评论内容" clearable>
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="4">
      <el-form-item label="显示状态">
        <el-select
          v-model="filterForm.status"
          placeholder="选择显示状态"
          clearable
          style="width: 100%"
        >
          <el-option label="显示" :value="1" />
          <el-option label="隐藏" :value="0" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="4">
      <el-form-item label="删除状态">
        <el-select
          v-model="filterForm.is_deleted"
          placeholder="选择删除状态"
          clearable
          style="width: 100%"
        >
          <el-option label="已删除" :value="true" />
          <el-option label="未删除" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="4">
      <el-form-item label="评论类型">
        <el-select
          v-model="filterForm.is_sub"
          placeholder="选择评论类型"
          clearable
          style="width: 100%"
        >
          <el-option label="子评论" :value="true" />
          <el-option label="父评论" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="7">
      <el-form-item label="评论时间">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 100%"
          @change="handleDateChange"
        />
      </el-form-item>
    </el-col>
  </filter-panel>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { Search } from '@element-plus/icons-vue';
import FilterPanel from '@/components/common/FilterPanel.vue';
import type { CommentListQuery } from '@/types/comment';

/**
 * 防抖函数
 * @param fn 要执行的函数
 * @param delay 延迟时间（毫秒）
 * @returns 防抖后的函数
 */
function debounce<T extends (...args: unknown[]) => unknown>(fn: T, delay: number) {
  let timer: ReturnType<typeof setTimeout> | null = null;
  return function (this: ThisParameterType<T>, ...args: Parameters<T>) {
    if (timer) clearTimeout(timer);
    timer = setTimeout(() => {
      fn.apply(this, args);
    }, delay);
  };
}

/**
 * 组件属性定义
 */
const props = defineProps<{
  modelValue: CommentListQuery;
}>();

/**
 * 组件事件定义
 */
const emit = defineEmits<{
  'update:modelValue': [value: CommentListQuery];
  search: [];
  close: [];
}>();

const filterForm = ref<CommentListQuery>({ ...props.modelValue });
const dateRange = ref<[string, string] | null>(null);

// 避免 watch 循环的标记
let isExternalUpdate = false;
let isResetting = false;

// 监听外部数据变化
watch(
  () => props.modelValue,
  newVal => {
    isExternalUpdate = true;
    filterForm.value = { ...newVal };
    if (newVal.start_time && newVal.end_time) {
      dateRange.value = [newVal.start_time, newVal.end_time];
    } else {
      dateRange.value = null;
    }
    setTimeout(() => {
      isExternalUpdate = false;
    }, 0);
  },
  { deep: true }
);

// 防抖的实时搜索
const debouncedSearch = debounce(() => {
  emit('update:modelValue', { ...filterForm.value });
  emit('search');
}, 500);

// 监听表单变化，实时触发搜索
watch(
  filterForm,
  () => {
    if (!isExternalUpdate && !isResetting) {
      debouncedSearch();
    }
  },
  { deep: true }
);

/**
 * 处理日期范围变化
 * @param val 日期范围值
 */
const handleDateChange = (val: [string, string] | null) => {
  if (val) {
    filterForm.value.start_time = val[0];
    filterForm.value.end_time = val[1];
  } else {
    filterForm.value.start_time = undefined;
    filterForm.value.end_time = undefined;
  }
};

/**
 * 处理重置
 */
const handleReset = () => {
  isResetting = true;
  dateRange.value = null;

  const page = filterForm.value.page;
  const pageSize = filterForm.value.page_size;
  filterForm.value = { page, page_size: pageSize };

  emit('update:modelValue', { ...filterForm.value });
  emit('search');

  setTimeout(() => {
    isResetting = false;
  }, 100);
};

// 组件挂载时初始化日期范围
if (filterForm.value.start_time && filterForm.value.end_time) {
  dateRange.value = [filterForm.value.start_time, filterForm.value.end_time];
}
</script>
