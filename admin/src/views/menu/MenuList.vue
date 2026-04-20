<template>
  <common-list
    title="菜单管理"
    :data="menuTree"
    :loading="loading"
    :show-pagination="false"
    create-text="新增菜单"
    @create="handleCreate"
    @refresh="fetchMenuTree"
    row-key="id"
    :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
    default-expand-all
  >
    <!-- 菜单类型切换器 -->
    <template #toolbar-before>
      <el-segmented v-model="selectedType" :options="menuTypeOptions" @change="handleTypeChange" />
    </template>

    <!-- 表格列 -->
    <el-table-column label="菜单标题" min-width="200">
      <template #default="{ row }">
        <div class="menu-title">
          <!-- RemixIcon 图标 -->
          <i v-if="row.icon && isRemixIcon(row.icon)" :class="row.icon" class="menu-icon"></i>
          <!-- 图片图标 -->
          <img v-else-if="row.icon" :src="row.icon" class="menu-icon-img" alt="icon" />
          <span>{{ row.title }}</span>
        </div>
      </template>
    </el-table-column>

    <el-table-column label="链接地址" min-width="250">
      <template #default="{ row }">
        <span v-if="row.url">{{ row.url }}</span>
        <span v-else style="color: #999">-</span>
      </template>
    </el-table-column>

    <el-table-column prop="sort" label="排序" width="100" align="center" />

    <el-table-column label="操作" width="180" align="center" fixed="right">
      <template #default="{ row }">
        <el-button
          v-if="row.parent_id === null"
          type="primary"
          link
          size="small"
          @click="handleAddChild(row)"
        >
          新增子菜单
        </el-button>
        <el-button type="primary" link size="small" @click="handleEdit(row)"> 编辑 </el-button>
        <el-button type="danger" link size="small" @click="handleDelete(row.id)"> 删除 </el-button>
      </template>
    </el-table-column>

    <!-- 额外内容 -->
    <template #extra>
      <menu-form-dialog
        v-model="dialogVisible"
        :edit-menu="currentMenu"
        :parent-menu="parentMenu"
        :current-type="selectedType"
        @success="fetchMenuTree"
      />
    </template>
  </common-list>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import CommonList from '@/components/common/CommonList.vue';
import type { MenuTreeNode } from '@/types/menu';
import { getMenuTree, deleteMenu } from '@/api/menu';
import MenuFormDialog from './components/MenuFormDialog.vue';

const loading = ref(false);
const menuTree = ref<MenuTreeNode[]>([]);
const selectedType = ref<string>('aggregate');
const dialogVisible = ref(false);
const currentMenu = ref<MenuTreeNode | null>(null);
const parentMenu = ref<MenuTreeNode | null>(null);

// 菜单类型选项
const menuTypeOptions = [
  { label: '聚合菜单', value: 'aggregate' },
  { label: '导航菜单', value: 'navigation' },
  { label: '页脚菜单', value: 'footer' },
];

// 判断是否是 RemixIcon 图标类名
const isRemixIcon = (icon: string) => {
  return icon && icon.startsWith('ri-');
};

// 获取菜单树
const fetchMenuTree = async () => {
  loading.value = true;
  try {
    const data = await getMenuTree(selectedType.value);
    menuTree.value = data;
  } catch (_error) {
    ElMessage.error('获取菜单列表失败');
  } finally {
    loading.value = false;
  }
};

// 菜单类型切换
const handleTypeChange = () => {
  fetchMenuTree();
};

// 新增菜单
const handleCreate = () => {
  currentMenu.value = null;
  parentMenu.value = null;
  dialogVisible.value = true;
};

// 新增子菜单
const handleAddChild = (menu: MenuTreeNode) => {
  currentMenu.value = null;
  parentMenu.value = menu;
  dialogVisible.value = true;
};

// 查找父菜单
const findParentMenu = (parentId: number | null): MenuTreeNode | null => {
  if (!parentId) return null;

  for (const menu of menuTree.value) {
    if (menu.id === parentId) {
      return menu;
    }
  }
  return null;
};

// 编辑菜单
const handleEdit = (menu: MenuTreeNode) => {
  currentMenu.value = menu;
  // 如果是子菜单，查找并设置父菜单
  parentMenu.value = menu.parent_id ? findParentMenu(menu.parent_id) : null;
  dialogVisible.value = true;
};

// 查找菜单节点
const findMenuNode = (id: number, nodes: MenuTreeNode[] = menuTree.value): MenuTreeNode | null => {
  for (const node of nodes) {
    if (node.id === id) {
      return node;
    }
    if (node.children && node.children.length > 0) {
      const found = findMenuNode(id, node.children);
      if (found) return found;
    }
  }
  return null;
};

// 删除菜单
const handleDelete = async (id: number) => {
  try {
    const menuNode = findMenuNode(id);
    const hasChildren = menuNode?.children && menuNode.children.length > 0;

    if (hasChildren) {
      // 有子菜单，询问如何处理
      const result = await ElMessageBox.confirm(
        `包含 ${menuNode.children?.length} 个子菜单。`,
        '提示',
        {
          distinguishCancelAndClose: true,
          confirmButtonText: '保留子菜单',
          cancelButtonText: '全部删除',
          type: 'warning',
        }
      )
        .then(() => 'upgrade' as const)
        .catch(action => {
          if (action === 'cancel') return 'delete' as const;
          throw 'close';
        });

      await deleteMenu(id, { children_action: result });
      ElMessage.success('删除成功');
      fetchMenuTree();
    } else {
      // 无子菜单，直接删除
      await ElMessageBox.confirm('确定要删除此菜单吗？', '提示', {
        type: 'warning',
      });
      await deleteMenu(id);
      ElMessage.success('删除成功');
      fetchMenuTree();
    }
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      if (error instanceof Error) {
        ElMessage.error(error.message);
      }
    }
  }
};

onMounted(() => {
  fetchMenuTree();
});
</script>

<style scoped lang="scss">
.menu-title {
  display: flex;
  align-items: center;

  .menu-icon {
    margin-right: 8px;
    font-size: 16px;
    color: #606266;
  }

  .menu-icon-img {
    width: 16px;
    height: 16px;
    margin-right: 8px;
    object-fit: contain;
    vertical-align: middle;
  }
}

// 树形表格样式优化
:deep(.el-table) {
  .el-table__expand-icon {
    font-size: 14px;
    color: #606266;

    &.el-table__expand-icon--expanded {
      transform: rotate(90deg);
    }
  }

  .el-table__indent {
    padding-left: 20px;
  }

  .el-table__placeholder {
    display: inline-block;
    width: 20px;
  }

  // 确保树形结构有正确的缩进
  .el-table__body {
    .el-table__row {
      .el-table__cell {
        &:first-child {
          .cell {
            display: flex;
            align-items: center;
          }
        }
      }
    }
  }
}
</style>
