<template>
  <filter-panel v-model="filterForm" title="筛选条件" @reset="handleReset" @close="$emit('close')">
    <el-col :span="5">
      <el-form-item label="关键词">
        <el-input v-model="filterForm.keyword" placeholder="搜索文件名" clearable>
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="文件类型">
        <el-select
          v-model="filterForm.file_type"
          placeholder="选择文件类型"
          clearable
          style="width: 100%"
        >
          <el-option label="图片" value="image" />
          <el-option label="视频" value="video" />
          <el-option label="音频" value="audio" />
          <el-option label="文档" value="application" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="使用状态">
        <el-select
          v-model="filterForm.status"
          placeholder="选择使用状态"
          clearable
          style="width: 100%"
        >
          <el-option label="使用中" :value="1" />
          <el-option label="未使用" :value="0" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="上传用途">
        <el-select
          v-model="filterForm.upload_type"
          placeholder="选择上传用途"
          clearable
          style="width: 100%"
        >
          <el-option-group label="文章板块">
            <el-option label="文章封面" value="文章封面" />
            <el-option label="文章图片" value="文章图片" />
            <el-option label="文章视频" value="文章视频" />
            <el-option label="文章音频" value="文章音频" />
          </el-option-group>
          <el-option-group label="动态板块">
            <el-option label="动态配图" value="动态配图" />
            <el-option label="动态视频" value="动态视频" />
          </el-option-group>
          <el-option-group label="互动板块">
            <el-option label="评论贴图" value="评论贴图" />
            <el-option label="反馈投诉" value="反馈投诉" />
            <el-option label="友情链接A" value="友情链接A" />
            <el-option label="友情链接S" value="友情链接S" />
          </el-option-group>
          <el-option-group label="系统板块">
            <el-option label="菜单图标" value="菜单图标" />
            <el-option label="用户头像" value="用户头像" />
          </el-option-group>
          <el-option-group label="设置板块">
            <el-option label="博客图标" value="博客图标" />
            <el-option label="博客背景" value="博客背景" />
            <el-option label="博客截图" value="博客截图" />
            <el-option label="展览图片" value="展览图片" />
            <el-option label="站长头像" value="站长头像" />
            <el-option label="站长形象" value="站长形象" />
          </el-option-group>
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="3">
      <el-form-item label="文件大小">
        <el-select
          v-model="sizeRange"
          placeholder="选择文件大小"
          clearable
          style="width: 100%"
          @change="handleSizeChange"
        >
          <el-option label="< 100KB" value="0-102400" />
          <el-option label="100KB - 1MB" value="102400-1048576" />
          <el-option label="1MB - 10MB" value="1048576-10485760" />
          <el-option label="> 10MB" value="10485760-0" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="7">
      <el-form-item label="上传时间">
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
import { Search } from '@element-plus/icons-vue';
import FilterPanel from '@/components/common/FilterPanel.vue';
import type { FileListQuery } from '@/types/file';

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
  modelValue: FileListQuery;
}>();

/**
 * 组件事件定义
 */
const emit = defineEmits<{
  'update:modelValue': [value: FileListQuery];
  search: [];
  close: [];
}>();

const filterForm = ref<FileListQuery>({ ...props.modelValue });
const dateRange = ref<[string, string] | null>(null);
const sizeRange = ref<string | null>(null);

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
 * 处理文件大小范围变化
 * @param val 大小范围值
 */
const handleSizeChange = (val: string | null) => {
  if (val) {
    const [min, max] = val.split('-').map(Number);
    filterForm.value.min_size = min;
    filterForm.value.max_size = max || undefined;
  } else {
    filterForm.value.min_size = undefined;
    filterForm.value.max_size = undefined;
  }
};

/**
 * 处理重置
 */
const handleReset = () => {
  isResetting = true;
  dateRange.value = null;
  sizeRange.value = null;

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
