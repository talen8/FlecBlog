<template>
  <filter-panel v-model="filterForm" title="筛选条件" @reset="handleReset" @close="$emit('close')">
    <!-- 第一行：关键词、友链类型、失效状态、审核状态 -->
    <el-col :span="6">
      <el-form-item label="关键词">
        <el-input v-model="filterForm.keyword" placeholder="搜索名称、链接、描述" clearable>
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="6">
      <el-form-item label="友链类型">
        <el-select
          v-model="filterForm.type_id"
          placeholder="选择类型"
          clearable
          style="width: 100%"
        >
          <el-option
            v-for="type in friendTypes"
            :key="type.id"
            :label="type.name"
            :value="type.id"
          />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="6">
      <el-form-item label="失效状态">
        <el-select
          v-model="filterForm.is_invalid"
          placeholder="选择失效状态"
          clearable
          style="width: 100%"
        >
          <el-option label="已失效" :value="true" />
          <el-option label="正常" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="6">
      <el-form-item label="审核状态">
        <el-select
          v-model="filterForm.is_pending"
          placeholder="选择审核状态"
          clearable
          style="width: 100%"
        >
          <el-option label="待审核" :value="true" />
          <el-option label="已通过" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <!-- 第二行：访问状态、RSS订阅 -->
    <el-col :span="6">
      <el-form-item label="访问状态">
        <el-select
          v-model="filterForm.accessible_status"
          placeholder="选择访问状态"
          clearable
          style="width: 100%"
        >
          <el-option label="正常" value="normal" />
          <el-option label="异常" value="abnormal" />
          <el-option label="忽略检查" value="ignored" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="6">
      <el-form-item label="RSS状态">
        <el-select
          v-model="filterForm.rss_status"
          placeholder="选择RSS状态"
          clearable
          style="width: 100%"
        >
          <el-option label="无订阅" value="no_rss" />
          <el-option label="正常订阅" value="normal" />
          <el-option label="三个月未更新" value="warning" />
          <el-option label="六个月未更新" value="danger" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="6">
      <el-form-item label="包含截图">
        <el-select
          v-model="filterForm.has_screenshot"
          placeholder="是否添加网站截图"
          clearable
          style="width: 100%"
        >
          <el-option label="有截图" :value="true" />
          <el-option label="无截图" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>
  </filter-panel>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { Search } from '@element-plus/icons-vue';
import FilterPanel from '@/components/common/FilterPanel.vue';
import type { FriendQuery } from '@/types/friend';
import type { FriendType } from '@/types/friend';
import { getFriendTypes } from '@/api/friend';

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
  modelValue: FriendQuery;
}>();

/**
 * 组件事件定义
 */
const emit = defineEmits<{
  'update:modelValue': [value: FriendQuery];
  search: [];
  close: [];
}>();

const filterForm = ref<FriendQuery>({ ...props.modelValue });
const friendTypes = ref<FriendType[]>([]);

// 避免 watch 循环的标记
let isExternalUpdate = false;
let isResetting = false;

// 监听外部数据变化
watch(
  () => props.modelValue,
  newVal => {
    isExternalUpdate = true;
    filterForm.value = { ...newVal };
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
 * 处理重置
 */
const handleReset = () => {
  isResetting = true;

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
 * 加载友链类型列表
 */
const loadFriendTypes = async () => {
  try {
    const result = await getFriendTypes();
    friendTypes.value = result.list || [];
  } catch (error) {
    console.error('加载友链类型列表失败:', error);
  }
};

// 组件挂载时加载数据
onMounted(() => {
  loadFriendTypes();
});
</script>
