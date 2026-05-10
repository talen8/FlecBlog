<template>
  <div class="pagination">
    <button
      class="pagination-btn"
      :disabled="currentPage <= 1"
      @click="handlePageChange(currentPage - 1)"
    >
      上一页
    </button>
    <span class="pagination-info">{{ currentPage }} / {{ totalPages }}</span>
    <button
      class="pagination-btn"
      :disabled="currentPage >= totalPages"
      @click="handlePageChange(currentPage + 1)"
    >
      下一页
    </button>
  </div>
</template>

<script lang="ts" setup>
const props = defineProps<{
  total: number;
  currentPage: number;
  pageSize: number;
}>();

const emit = defineEmits<{
  change: [page: number];
}>();

const totalPages = computed(() => Math.ceil(props.total / props.pageSize));

const handlePageChange = (page: number) => {
  if (page >= 1 && page <= totalPages.value) {
    emit('change', page);
  }
};
</script>

<style lang="scss" scoped>
.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 16px;
  margin-top: 30px;
}

.pagination-btn {
  padding: 8px 16px;
  background: var(--flec-card-bg);
  border: 1px solid var(--flec-border);
  color: var(--font-color);
  font-size: 0.9rem;
  cursor: pointer;

  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
}

.pagination-info {
  color: var(--theme-meta-color);
  font-size: 0.9rem;
}
</style>
