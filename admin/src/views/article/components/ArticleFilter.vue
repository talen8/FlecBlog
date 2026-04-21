<template>
  <filter-panel v-model="filterForm" title="筛选条件" @reset="handleReset" @close="$emit('close')">
    <!-- 第一行：关键字、分类、标签、发布状态、过时状态 -->
    <el-col :span="5">
      <el-form-item label="关键词">
        <el-input v-model="filterForm.keyword" placeholder="搜索标题或内容" clearable>
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="5">
      <el-form-item label="分类">
        <el-select
          v-model="filterForm.category_id"
          placeholder="选择分类"
          clearable
          style="width: 100%"
        >
          <el-option
            v-for="category in categoryList"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="5">
      <el-form-item label="标签">
        <el-select
          v-model="filterForm.tag_ids"
          placeholder="选择标签"
          multiple
          collapse-tags
          collapse-tags-tooltip
          clearable
          style="width: 100%"
        >
          <el-option v-for="tag in tagList" :key="tag.id" :label="tag.name" :value="tag.id" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="5">
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

    <el-col :span="4">
      <el-form-item label="过时状态">
        <el-select
          v-model="filterForm.is_outdated"
          placeholder="选择过时状态"
          clearable
          style="width: 100%"
        >
          <el-option label="已过时" :value="true" />
          <el-option label="未过时" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <!-- 第二行：置顶状态、精选状态、发布地点、发布时间 -->
    <el-col :span="5">
      <el-form-item label="置顶状态">
        <el-select
          v-model="filterForm.is_top"
          placeholder="选择置顶状态"
          clearable
          style="width: 100%"
        >
          <el-option label="已置顶" :value="true" />
          <el-option label="未置顶" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="5">
      <el-form-item label="精选状态">
        <el-select
          v-model="filterForm.is_essence"
          placeholder="选择精选状态"
          clearable
          style="width: 100%"
        >
          <el-option label="已精选" :value="true" />
          <el-option label="未精选" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="5">
      <el-form-item label="发布地点">
        <el-input v-model="filterForm.location" placeholder="搜索发布地点" clearable>
          <template #prefix>
            <el-icon><Location /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="9">
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
import { ref, onMounted, watch } from 'vue';
import { Search, Location } from '@element-plus/icons-vue';
import FilterPanel from '@/components/common/FilterPanel.vue';
import type { Category } from '@/types/category';
import type { Tag } from '@/types/tag';
import type { ArticleListQuery } from '@/types/article';
import { getCategories } from '@/api/category';
import { getTags } from '@/api/tag';

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
  modelValue: ArticleListQuery;
}>();

/**
 * 组件事件定义
 */
const emit = defineEmits<{
  'update:modelValue': [value: ArticleListQuery];
  search: [];
  close: [];
}>();

const filterForm = ref<ArticleListQuery>({ ...props.modelValue });
const dateRange = ref<[string, string] | null>(null);
const categoryList = ref<Category[]>([]);
const tagList = ref<Tag[]>([]);

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

/**
 * 加载分类列表
 */
const loadCategories = async () => {
  try {
    const result = await getCategories({ page: 1, page_size: 1000 });
    categoryList.value = result.list;
  } catch (error) {
    console.error('加载分类列表失败:', error);
  }
};

/**
 * 加载标签列表
 */
const loadTags = async () => {
  try {
    const result = await getTags({ page: 1, page_size: 1000 });
    tagList.value = result.list;
  } catch (error) {
    console.error('加载标签列表失败:', error);
  }
};

// 组件挂载时加载数据
onMounted(() => {
  loadCategories();
  loadTags();

  if (filterForm.value.start_time && filterForm.value.end_time) {
    dateRange.value = [filterForm.value.start_time, filterForm.value.end_time];
  }
});
</script>
