<template>
  <filter-panel v-model="filterForm" title="筛选条件" @reset="handleReset" @close="$emit('close')">
    <el-col :span="3">
      <el-form-item label="页面URL">
        <el-input v-model="filterForm.keyword" placeholder="搜索页面URL" clearable>
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="访客ID">
        <el-input v-model="filterForm.visitor_id" placeholder="输入访客ID" clearable>
          <template #prefix>
            <el-icon><User /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="IP地址">
        <el-input v-model="filterForm.ip" placeholder="输入IP地址" clearable>
          <template #prefix>
            <el-icon><Connection /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="地理位置">
        <el-input v-model="filterForm.location" placeholder="搜索地理位置" clearable>
          <template #prefix>
            <el-icon><Location /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="浏览器">
        <el-input v-model="filterForm.browser" placeholder="搜索浏览器" clearable>
          <template #prefix>
            <el-icon><Monitor /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="操作系统">
        <el-input v-model="filterForm.os" placeholder="搜索操作系统" clearable>
          <template #prefix>
            <el-icon><Platform /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="6">
      <el-form-item label="访问时间">
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
import { ref, watch, onMounted } from 'vue';
import { Search, User, Connection, Location, Monitor, Platform } from '@element-plus/icons-vue';
import FilterPanel from '@/components/common/FilterPanel.vue';
import type { VisitListQuery } from '@/types/stats';

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
  modelValue: VisitListQuery;
}>();

/**
 * 组件事件定义
 */
const emit = defineEmits<{
  'update:modelValue': [value: VisitListQuery];
  search: [];
  close: [];
}>();

const filterForm = ref<VisitListQuery>({ ...props.modelValue });
const dateRange = ref<[string, string] | null>(null);

let isExternalUpdate = false;
let isResetting = false;

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

const debouncedSearch = debounce(() => {
  emit('update:modelValue', { ...filterForm.value });
  emit('search');
}, 500);

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

onMounted(() => {
  if (filterForm.value.start_time && filterForm.value.end_time) {
    dateRange.value = [filterForm.value.start_time, filterForm.value.end_time];
  }
});
</script>
