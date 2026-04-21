<template>
  <filter-panel v-model="filterForm" title="筛选条件" @reset="handleReset" @close="$emit('close')">
    <!-- 第一行：关键词、标签、发布地点、发布状态 -->
    <el-col :span="6">
      <el-form-item label="关键词">
        <el-input v-model="filterForm.keyword" placeholder="搜索文本内容" clearable>
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="6">
      <el-form-item label="标签">
        <el-input v-model="filterForm.tags" placeholder="搜索标签" clearable>
          <template #prefix>
            <el-icon><PriceTag /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="6">
      <el-form-item label="发布地点">
        <el-input v-model="filterForm.location" placeholder="搜索发布地点" clearable>
          <template #prefix>
            <el-icon><Location /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="6">
      <el-form-item label="发布状态">
        <el-select
          v-model="filterForm.is_publish"
          placeholder="选择发布状态"
          clearable
          style="width: 100%"
        >
          <el-option label="已发布" :value="true" />
          <el-option label="草稿" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <!-- 第二行：内容类型筛选 -->
    <el-col :span="4">
      <el-form-item label="包含图片">
        <el-select v-model="filterForm.has_images" placeholder="全部" clearable style="width: 100%">
          <el-option label="有图片" :value="true" />
          <el-option label="无图片" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="4">
      <el-form-item label="包含视频">
        <el-select v-model="filterForm.has_video" placeholder="全部" clearable style="width: 100%">
          <el-option label="有视频" :value="true" />
          <el-option label="无视频" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="4">
      <el-form-item label="包含音乐">
        <el-select v-model="filterForm.has_music" placeholder="全部" clearable style="width: 100%">
          <el-option label="有音乐" :value="true" />
          <el-option label="无音乐" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="4">
      <el-form-item label="包含链接">
        <el-select v-model="filterForm.has_link" placeholder="全部" clearable style="width: 100%">
          <el-option label="有链接" :value="true" />
          <el-option label="无链接" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="8">
      <el-form-item label="发布时间">
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
import { Search, PriceTag, Location } from '@element-plus/icons-vue';
import FilterPanel from '@/components/common/FilterPanel.vue';
import type { MomentListQuery } from '@/types/moment';

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
  modelValue: MomentListQuery;
}>();

/**
 * 组件事件定义
 */
const emit = defineEmits<{
  'update:modelValue': [value: MomentListQuery];
  search: [];
  close: [];
}>();

const filterForm = ref<MomentListQuery>({ ...props.modelValue });
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
</script>
