<template>
  <filter-panel v-model="filterForm" title="筛选条件" @reset="handleReset" @close="$emit('close')">
    <!-- 第一行：关键词、是否启用、是否删除、注册时间 -->
    <el-col :span="6">
      <el-form-item label="关键词">
        <el-input v-model="filterForm.keyword" placeholder="搜索邮箱、昵称" clearable>
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </el-form-item>
    </el-col>

    <el-col :span="4">
      <el-form-item label="是否启用">
        <el-select
          v-model="filterForm.is_enabled"
          placeholder="选择状态"
          clearable
          style="width: 100%"
        >
          <el-option label="启用" :value="true" />
          <el-option label="禁用" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="4">
      <el-form-item label="是否删除">
        <el-select
          v-model="filterForm.is_deleted"
          placeholder="选择状态"
          clearable
          style="width: 100%"
        >
          <el-option label="已删除" :value="true" />
          <el-option label="未删除" :value="false" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="10">
      <el-form-item label="注册时间">
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

    <!-- 第二行：角色、登录方式、最后登录时间 -->
    <el-col :span="4">
      <el-form-item label="角色">
        <el-select v-model="filterForm.role" placeholder="选择角色" clearable style="width: 100%">
          <el-option label="超级管理员" value="super_admin" />
          <el-option label="管理员" value="admin" />
          <el-option label="普通用户" value="user" />
          <el-option label="访客" value="guest" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="4">
      <el-form-item label="登录方式">
        <el-select
          v-model="filterForm.login_method"
          placeholder="选择登录方式"
          clearable
          style="width: 100%"
        >
          <el-option label="密码登录" value="password" />
          <el-option label="GitHub" value="github" />
          <el-option label="Google" value="google" />
          <el-option label="QQ" value="qq" />
          <el-option label="Microsoft" value="microsoft" />
        </el-select>
      </el-form-item>
    </el-col>

    <el-col :span="10">
      <el-form-item label="最后登录时间">
        <el-date-picker
          v-model="lastLoginDateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 100%"
          @change="handleLastLoginDateChange"
        />
      </el-form-item>
    </el-col>
  </filter-panel>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { Search } from '@element-plus/icons-vue';
import FilterPanel from '@/components/common/FilterPanel.vue';
import type { UserListQuery } from '@/types/user';

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
  modelValue: UserListQuery;
}>();

/**
 * 组件事件定义
 */
const emit = defineEmits<{
  'update:modelValue': [value: UserListQuery];
  search: [];
  close: [];
}>();

const filterForm = ref<UserListQuery>({ ...props.modelValue });
const dateRange = ref<[string, string] | null>(null);
const lastLoginDateRange = ref<[string, string] | null>(null);

// 避免 watch 循环的标记
let isExternalUpdate = false;
let isResetting = false;

// 监听外部数据变化
watch(
  () => props.modelValue,
  newVal => {
    isExternalUpdate = true;
    filterForm.value = { ...newVal };
    // 初始化注册时间范围
    if (newVal.start_time && newVal.end_time) {
      dateRange.value = [newVal.start_time, newVal.end_time];
    } else {
      dateRange.value = null;
    }
    // 初始化最后登录时间范围
    if (newVal.last_login_start && newVal.last_login_end) {
      lastLoginDateRange.value = [newVal.last_login_start, newVal.last_login_end];
    } else {
      lastLoginDateRange.value = null;
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
 * 处理注册时间范围变化
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
 * 处理最后登录时间范围变化
 * @param val 日期范围值
 */
const handleLastLoginDateChange = (val: [string, string] | null) => {
  if (val) {
    filterForm.value.last_login_start = val[0];
    filterForm.value.last_login_end = val[1];
  } else {
    filterForm.value.last_login_start = undefined;
    filterForm.value.last_login_end = undefined;
  }
};

/**
 * 处理重置
 */
const handleReset = () => {
  isResetting = true;
  dateRange.value = null;
  lastLoginDateRange.value = null;

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
if (filterForm.value.last_login_start && filterForm.value.last_login_end) {
  lastLoginDateRange.value = [filterForm.value.last_login_start, filterForm.value.last_login_end];
}
</script>
