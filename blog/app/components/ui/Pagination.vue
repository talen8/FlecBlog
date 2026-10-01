<script setup lang="ts">
interface Props {
  total: number; // 总条目数
  currentPage: number; // 当前页码（从1开始）
  pageSize: number; // 每页条目数
  maxVisiblePages?: number; // 最多显示多少个页码按钮
}

const props = withDefaults(defineProps<Props>(), {
  maxVisiblePages: 5,
});

const emit = defineEmits<{
  change: [page: number];
}>();

// 计算总页数
const totalPages = computed(() => {
  return Math.ceil(props.total / props.pageSize);
});

// 是否有上一页
const hasPrev = computed(() => props.currentPage > 1);

// 是否有下一页
const hasNext = computed(() => props.currentPage < totalPages.value);

// 计算显示的页码数组
const visiblePages = computed(() => {
  const total = totalPages.value;
  const current = props.currentPage;
  const max = props.maxVisiblePages;

  if (total <= max) {
    // 总页数小于等于最大显示数，显示所有页码
    return Array.from({ length: total }, (_, i) => i + 1);
  }

  const pages: (number | string)[] = [];

  // 计算需要在中间显示的页码数量（除去首尾两页）
  const middleCount = max - 2;
  const half = Math.floor(middleCount / 2);

  // 始终显示第一页
  pages.push(1);

  let start: number;
  let end: number;

  // 判断当前页的位置，动态计算中间页码范围
  if (current <= half + 2) {
    // 当前页靠近开头，显示前面的页码
    start = 2;
    end = Math.min(middleCount + 1, total - 1);
  } else if (current >= total - half - 1) {
    // 当前页靠近结尾，显示后面的页码
    start = Math.max(2, total - middleCount);
    end = total - 1;
  } else {
    // 当前页在中间，以当前页为中心显示
    start = current - half;
    end = current + half;

    // 调整范围以保持显示的页码数量一致
    const diff = middleCount - (end - start);
    if (diff > 0) {
      if (end < total - 1) {
        end = Math.min(total - 1, end + diff);
      } else {
        start = Math.max(2, start - diff);
      }
    }
  }

  // 添加左侧省略号
  if (start > 2) {
    pages.push('...');
  }

  // 添加中间页码
  for (let i = start; i <= end; i++) {
    if (i > 1 && i < total) {
      pages.push(i);
    }
  }

  // 添加右侧省略号
  if (end < total - 1) {
    pages.push('...');
  }

  // 始终显示最后一页
  if (total > 1) {
    pages.push(total);
  }

  return pages;
});

// 跳转到指定页
const goToPage = (page: number) => {
  if (page === props.currentPage || page < 1 || page > totalPages.value) {
    return;
  }
  emit('change', page);
};

// 上一页
const prevPage = () => {
  if (hasPrev.value) {
    goToPage(props.currentPage - 1);
  }
};

// 下一页
const nextPage = () => {
  if (hasNext.value) {
    goToPage(props.currentPage + 1);
  }
};

// 判断是否为省略号
const isEllipsis = (page: number | string): page is string => {
  return typeof page === 'string';
};
</script>

<template>
  <div v-if="totalPages > 1" class="pagination">
    <div class="pagination-controls">
      <!-- 上一页 -->
      <button
        class="pagination-btn"
        :class="{ disabled: !hasPrev }"
        :disabled="!hasPrev"
        @click="prevPage"
      >
        <i class="ri-arrow-left-s-line" />
      </button>

      <!-- 页码 -->
      <template v-for="page in visiblePages" :key="page">
        <span v-if="isEllipsis(page)" class="pagination-ellipsis">...</span>
        <button
          v-else
          class="pagination-btn pagination-number"
          :class="{ active: page === currentPage }"
          @click="goToPage(page)"
        >
          {{ page }}
        </button>
      </template>

      <!-- 下一页 -->
      <button
        class="pagination-btn"
        :class="{ disabled: !hasNext }"
        :disabled="!hasNext"
        @click="nextPage"
      >
        <i class="ri-arrow-right-s-line" />
      </button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px 0;

  .pagination-controls {
    display: flex;
    align-items: center;
    gap: 10px;
    flex-wrap: wrap;
  }

  .pagination-btn {
    min-width: 40px;
    height: 40px;
    background: var(--flec-card-bg);
    color: var(--font-color);
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1.1rem;

    &:hover:not(.disabled):not(.active) {
      color: #fff;
      background: var(--flec-btn-hover);
    }

    &.active {
      background: var(--flec-btn);
      color: #fff;
      font-weight: 600;
      cursor: default;
    }

    &.disabled {
      opacity: 0.4;
      cursor: not-allowed;
    }

    i {
      font-size: 1.5rem;
    }
  }

  .pagination-ellipsis {
    color: var(--font-color);
    padding: 0 4px;
    user-select: none;
  }
}

// 响应式设计
@media screen and (max-width: 768px) {
  .pagination {
    padding: 14px 0;

    .pagination-controls {
      gap: 6px;
    }

    .pagination-btn {
      min-width: 34px;
      height: 34px;
      font-size: 0.95rem;

      i {
        font-size: 1.25rem;
      }
    }

    .pagination-ellipsis {
      padding: 0 3px;
    }
  }
}
</style>
