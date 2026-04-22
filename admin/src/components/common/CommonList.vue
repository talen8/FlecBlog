<template>
  <div class="common-list">
    <el-card shadow="never">
      <!-- 工具栏 -->
      <div class="toolbar">
        <h2>{{ title }}</h2>
        <div class="actions">
          <!-- 前工具栏 -->
          <slot name="toolbar-before" />
          <el-button v-if="showCreate" type="primary" @click="$emit('create')">
            {{ createText }}
          </el-button>
          <!-- 后工具栏 -->
          <slot name="toolbar-after" />
          <el-badge
            v-if="showFilter"
            :value="filterCount"
            :hidden="filterCount === 0"
            class="filter-badge"
          >
            <el-button :type="filterActive ? 'success' : 'default'" @click="$emit('filter')">
              <el-icon>
                <Filter />
              </el-icon>
            </el-button>
          </el-badge>
          <el-button class="refresh-btn" @click="$emit('refresh')">
            <el-icon>
              <Refresh />
            </el-icon>
          </el-button>
        </div>
      </div>

      <!-- 额外内容 -->
      <slot name="extra" />

      <!-- 表格区域 -->
      <div class="table-wrapper">
        <!-- 加载状态 -->
        <div v-if="loading" class="common-list-loading">
          <el-skeleton :rows="5" animated />
        </div>

        <!-- 表格 - 完全由外部控制 -->
        <el-table v-else :data="data" border style="width: 100%; height: 100%" v-bind="$attrs">
          <slot />
        </el-table>
      </div>

      <!-- 分页 -->
      <div v-if="showPagination" class="pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next"
          @current-change="$emit('update:page', $event)"
          @size-change="$emit('update:pageSize', $event)"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { Refresh, Filter } from '@element-plus/icons-vue';

withDefaults(
  defineProps<{
    title: string;
    data: unknown[];
    loading?: boolean;
    total?: number;
    page?: number;
    pageSize?: number;
    showPagination?: boolean;
    showCreate?: boolean;
    createText?: string;
    showFilter?: boolean;
    filterActive?: boolean;
    filterCount?: number;
  }>(),
  {
    loading: false,
    total: 0,
    page: 1,
    pageSize: 10,
    showPagination: true,
    showCreate: true,
    createText: '新增',
    showFilter: true,
    filterActive: false,
    filterCount: 0,
  }
);

defineEmits<{
  create: [];
  refresh: [];
  filter: [];
  'update:page': [page: number];
  'update:pageSize': [size: number];
}>();
</script>

<style scoped lang="scss">
.common-list {
  height: 100%;

  :deep(.el-card) {
    height: 100%;
    display: flex;
    flex-direction: column;

    .el-card__body {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
    }
  }

  .toolbar {
    margin-bottom: 12px;
    display: flex;
    justify-content: space-between;
    align-items: center;

    h2 {
      margin: 0;
      font-size: 20px;
      font-weight: 500;
    }

    .actions {
      display: flex;
      gap: 12px;

      :deep(.el-button + .el-button) {
        margin-left: 0;
      }

      // 筛选按钮hover时显示绿色，与active状态保持一致
      :deep(.el-button--default:hover) {
        color: var(--el-color-success);
        border-color: var(--el-color-success-light-5);
        background-color: var(--el-color-success-light-9);
      }
    }

    @media (max-width: 767px) {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;

      h2 {
        font-size: 18px;
      }

      .actions {
        width: 100%;
        flex-wrap: wrap;

        .refresh-btn {
          display: none;
        }
      }
    }
  }

  .table-wrapper {
    flex: 1;
    overflow: auto;

    :deep(.el-table__header th .cell) {
      text-align: center;
    }
  }

  .pagination {
    display: flex;
    justify-content: flex-end;
    padding-top: 12px;

    @media (max-width: 767px) {
      justify-content: center;
    }
  }
}
</style>
