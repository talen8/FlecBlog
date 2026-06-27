<template>
  <div class="theme-menu-panel">
    <div class="menu-toolbar">
      <el-segmented
        v-model="selectedType"
        :options="menuTypeOptions"
        class="menu-type-segmented"
        @change="fetchMenuTree"
      >
        <template #default="{ item }">
          <div class="segmented-item">
            <el-icon class="segmented-icon"><component :is="item.icon" /></el-icon>
            <span class="segmented-text">{{ item.label }}</span>
          </div>
        </template>
      </el-segmented>

      <div class="menu-actions">
        <el-button type="primary" :disabled="disabled || !themeSlug" @click="handleCreate">
          新增菜单
        </el-button>
        <el-button :disabled="loading || !themeSlug" @click="fetchMenuTree">刷新</el-button>
      </div>
    </div>

    <el-table
      v-loading="loading"
      :data="menuTree"
      border
      row-key="id"
      :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
      default-expand-all
      class="menu-table"
    >
      <el-table-column label="菜单标题" min-width="200">
        <template #default="{ row }">
          <div class="menu-title">
            <i v-if="row.icon && isRemixIcon(row.icon)" :class="row.icon" class="menu-icon"></i>
            <img v-else-if="row.icon" :src="row.icon" class="menu-icon-img" alt="icon" />
            <span>{{ row.title }}</span>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="链接地址" min-width="250">
        <template #default="{ row }">
          <span v-if="row.url">{{ row.url }}</span>
          <span v-else class="empty-value">-</span>
        </template>
      </el-table-column>

      <el-table-column prop="sort" label="排序" width="100" align="center" />

      <el-table-column label="操作" width="180" align="center" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="(row as any)._depth < currentMaxDepth - 1"
            type="primary"
            link
            size="small"
            :disabled="disabled"
            @click="handleAddChild(row)"
          >
            新增子菜单
          </el-button>
          <el-button type="primary" link size="small" :disabled="disabled" @click="handleEdit(row)">
            编辑
          </el-button>
          <el-button
            type="danger"
            link
            size="small"
            :disabled="disabled"
            @click="handleDelete(row.id)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="dialogVisible"
      :title="editingMenu ? '编辑菜单' : '新增菜单'"
      width="90%"
      style="max-width: 600px"
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <div class="form-info">
          <div class="info-item">
            <span class="info-label">菜单类型</span>
            <span class="info-value">{{ currentTypeLabel }}</span>
          </div>
          <div v-if="parentMenu && !editingMenu" class="info-item">
            <span class="info-label">父菜单</span>
            <span class="info-value">{{ parentMenu.title }}</span>
          </div>
        </div>

        <el-form-item label="菜单标题" prop="title">
          <el-input
            v-model="formData.title"
            placeholder="请输入菜单标题"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="链接地址" prop="url">
          <el-input
            v-model="formData.url"
            placeholder="请输入链接地址"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="图标" prop="icon">
          <div class="icon-input-wrapper">
            <el-input
              v-model="formData.icon"
              placeholder="请输入图标类名(ri-home-line)或上传图片"
              maxlength="500"
            >
              <template #append>
                <el-button @click="handleIconUpload">
                  <el-icon><Upload /></el-icon>
                  上传
                </el-button>
              </template>
            </el-input>
            <div v-if="formData.icon" class="icon-preview">
              <i v-if="isRemixIcon(formData.icon)" :class="formData.icon"></i>
              <img v-else :src="formData.icon" alt="图标预览" @error="handleIconError" />
            </div>
          </div>
        </el-form-item>

        <el-form-item v-if="editingMenu" label="父菜单">
          <el-select
            v-model="selectedParentId"
            :placeholder="editingMenu.children?.length ? '包含子菜单，无法设置' : '请选择父菜单'"
            :disabled="editingMenu.children?.length! > 0"
            clearable
            style="width: 100%"
          >
            <el-option
              v-for="menu in parentMenuOptions"
              :key="menu.id"
              :label="menu.title"
              :value="menu.id"
            >
              <div style="display: flex; align-items: center">
                <i
                  v-if="menu.icon && isRemixIcon(menu.icon)"
                  :class="menu.icon"
                  style="margin-right: 8px"
                ></i>
                <img
                  v-else-if="menu.icon"
                  :src="menu.icon"
                  style="width: 16px; height: 16px; margin-right: 8px; object-fit: contain"
                />
                <span>{{ menu.title }}</span>
              </div>
            </el-option>
          </el-select>
        </el-form-item>

        <el-form-item label="排序" prop="sort">
          <el-input-number v-model="formData.sort" :min="1" :max="10" />
        </el-form-item>

        <el-form-item label="是否启用" prop="is_enabled">
          <el-switch v-model="formData.is_enabled" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Menu, Upload } from '@element-plus/icons-vue';
import type { Component } from 'vue';
import type { FormInstance, FormRules } from 'element-plus';
import type { ThemeMenuItem, MenuSlot, ThemeSchema } from '@/types/theme';
import { updateThemeMenus } from '@/api/theme';
import { uploadFile } from '@/api/file';

const findMenu = (items: ThemeMenuItem[], id: number): ThemeMenuItem | null => {
  for (const item of items) {
    if (item.id === id) return item;
    const found = findMenu(item.children || [], id);
    if (found) return found;
  }
  return null;
};

const findParentOf = (items: ThemeMenuItem[], id: number): ThemeMenuItem | null => {
  for (const item of items) {
    if (item.children?.some(child => child.id === id)) return item;
    const found = findParentOf(item.children || [], id);
    if (found) return found;
  }
  return null;
};

const removeMenu = (
  items: ThemeMenuItem[],
  id: number,
  action: 'delete' | 'upgrade' = 'delete'
): ThemeMenuItem | null => {
  for (const [index, item] of items.entries()) {
    if (item.id === id) {
      const [removed] = items.splice(index, 1);
      if (!removed) return null;
      if (action === 'upgrade' && removed.children?.length) {
        items.splice(index, 0, ...removed.children);
      }
      return removed;
    }
    const found = removeMenu(item.children || [], id, action);
    if (found) return found;
  }
  return null;
};

const appendMenu = (
  items: ThemeMenuItem[],
  parentId: number | null,
  menu: ThemeMenuItem
): boolean => {
  if (!parentId) {
    items.push(menu);
    return true;
  }
  for (const item of items) {
    if (item.id === parentId) {
      item.children = [...(item.children || []), menu];
      return true;
    }
    if (appendMenu(item.children || [], parentId, menu)) return true;
  }
  return false;
};

interface MenuTypeOption {
  label: string;
  value: string;
  icon: Component;
}

const props = withDefaults(
  defineProps<{
    themeSlug: string;
    schema?: ThemeSchema | Record<string, unknown>;
    menus?: Record<string, ThemeMenuItem[]>;
    disabled?: boolean;
  }>(),
  {
    schema: () => ({}),
    menus: () => ({}),
    disabled: false,
  }
);

const emit = defineEmits<{
  refresh: [];
}>();

const menuTree = ref<ThemeMenuItem[]>([]);
const selectedType = ref('aggregate');
const loading = ref(false);
const dialogVisible = ref(false);
const editingMenu = ref<ThemeMenuItem | null>(null);
const parentMenu = ref<ThemeMenuItem | null>(null);

const menuTypeOptions = computed<MenuTypeOption[]>(() => {
  const slots = (props.schema as ThemeSchema)?.$menus || {};
  return Object.entries(slots).map(([value, slot]) => ({
    label: slot.label || slot.title || value,
    value,
    icon: Menu,
  }));
});

const currentTypeLabel = computed(() => {
  const slots = (props.schema as ThemeSchema)?.$menus || {};
  const slot = slots[selectedType.value];
  return slot?.label || slot?.title || selectedType.value;
});

const currentMaxDepth = computed(() => {
  const slots = (props.schema as ThemeSchema)?.$menus || {};
  const slot = slots[selectedType.value];
  return slot?.maxDepth ?? 2;
});

const isRemixIcon = (icon: string) => icon && icon.startsWith('ri-');

const sortMenuTree = (items: ThemeMenuItem[], depth = 0): ThemeMenuItem[] =>
  items
    .slice()
    .sort((a, b) => a.sort - b.sort || a.id - b.id)
    .map(item => ({
      ...item,
      _depth: depth,
      children: sortMenuTree(item.children || [], depth + 1),
    })) as ThemeMenuItem[];

let defaultIdCounter = 0;

const buildDefaults = (items: Partial<ThemeMenuItem>[]): ThemeMenuItem[] =>
  items.map(item => ({
    id: --defaultIdCounter,
    title: item.title || '',
    url: item.url || '',
    icon: item.icon || '',
    sort: item.sort ?? 5,
    is_enabled: item.is_enabled ?? true,
    children: buildDefaults(item.children || []),
  }));

const fetchMenuTree = () => {
  if (!props.themeSlug) return;

  const menus = props.menus || {};
  const type = selectedType.value;
  const saved = menus[type];
  const tree = sortMenuTree(
    saved?.length ? saved : buildDefaults(getCurrentSlot()?.defaults || [])
  );
  menuTree.value = tree;
};

const getCurrentSlot = () =>
  ((props.schema as ThemeSchema)?.$menus || {})[selectedType.value] as MenuSlot | undefined;

const handleCreate = () => {
  editingMenu.value = null;
  parentMenu.value = null;
  dialogVisible.value = true;
};

const handleAddChild = (menu: ThemeMenuItem) => {
  editingMenu.value = null;
  parentMenu.value = menu;
  dialogVisible.value = true;
};

const handleEdit = (menu: ThemeMenuItem) => {
  editingMenu.value = menu;
  parentMenu.value = findParentOf(menuTree.value, menu.id);
  dialogVisible.value = true;
};

const flushMenus = (): Record<string, ThemeMenuItem[]> => {
  return { ...props.menus, [selectedType.value]: menuTree.value };
};

const handleDelete = async (id: number) => {
  if (!props.themeSlug) return;

  try {
    const menuNode = findMenu(menuTree.value, id);
    const hasChildren = menuNode?.children && menuNode.children.length > 0;

    const action = hasChildren
      ? await ElMessageBox.confirm(`包含 ${menuNode.children?.length} 个子菜单。`, '提示', {
          distinguishCancelAndClose: true,
          confirmButtonText: '保留子菜单',
          cancelButtonText: '全部删除',
          type: 'warning',
        })
          .then(() => 'upgrade' as const)
          .catch(action => {
            if (action === 'cancel') return 'delete' as const;
            throw 'close';
          })
      : await ElMessageBox.confirm('确定要删除此菜单吗？', '提示', { type: 'warning' }).then(
          () => 'delete' as const
        );

    const tree = [...menuTree.value];
    removeMenu(tree, id, action);
    menuTree.value = tree;
    await updateThemeMenus(props.themeSlug, flushMenus());

    ElMessage.success('删除成功');
    emit('refresh');
  } catch (error) {
    if (error !== 'cancel' && error !== 'close' && error instanceof Error) {
      ElMessage.error(error.message);
    }
  }
};

const formRef = ref<FormInstance>();
const submitLoading = ref(false);
const parentMenuOptions = ref<ThemeMenuItem[]>([]);
const selectedParentId = ref<number | null>(null);
const pendingFile = ref<File | null>(null);

const formData = ref<ThemeMenuItem>({
  id: 0,
  title: '',
  url: '',
  icon: '',
  sort: 5,
  is_enabled: true,
  children: [],
});

const rules: FormRules = {
  title: [
    { message: '请输入菜单标题', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' },
  ],
  url: [{ max: 500, message: '链接地址不能超过 500 个字符', trigger: 'blur' }],
  icon: [{ max: 500, message: '图标不能超过 500 个字符', trigger: 'blur' }],
};

const cleanupIconBlob = () => {
  if (formData.value.icon.startsWith('blob:')) {
    URL.revokeObjectURL(formData.value.icon);
  }
};

const handleIconUpload = () => {
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = 'image/*';
  input.onchange = e => {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;

    cleanupIconBlob();
    const blobUrl = URL.createObjectURL(file);
    pendingFile.value = file;
    formData.value.icon = blobUrl;
  };
  input.click();
};

const handleIconError = (e: Event) => {
  const target = e.target as HTMLImageElement;
  target.style.display = 'none';
  ElMessage.warning('图标加载失败');
};

const fetchParentMenuOptions = () => {
  const allMenus = menuTree.value;
  let options = [...allMenus];
  if (editingMenu.value) {
    const excludeIds = [editingMenu.value.id];
    const stack = [...(editingMenu.value.children || [])];
    while (stack.length) {
      const node = stack.pop()!;
      excludeIds.push(node.id);
      stack.push(...(node.children || []));
    }
    options = options.filter(menu => !excludeIds.includes(menu.id));
  }
  parentMenuOptions.value = options;
};

const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    submitLoading.value = true;

    if (pendingFile.value) {
      const result = await uploadFile(pendingFile.value, '菜单图标');
      formData.value.icon = result.file_url;
    }

    const tree = menuTree.value;

    if (editingMenu.value) {
      const existing = findMenu(tree, editingMenu.value.id);
      if (!existing) throw new Error('菜单不存在');

      Object.assign(existing, {
        title: formData.value.title,
        url: formData.value.url,
        icon: formData.value.icon,
        sort: formData.value.sort,
        is_enabled: formData.value.is_enabled,
      });

      const currentParent = findParentOf(tree, editingMenu.value.id);
      if (currentParent?.id !== selectedParentId.value) {
        removeMenu(tree, editingMenu.value.id);
        if (!appendMenu(tree, selectedParentId.value, existing)) tree.push(existing);
      }

      ElMessage.success('更新成功');
    } else {
      const menu: ThemeMenuItem = {
        id: 0,
        title: formData.value.title,
        url: formData.value.url,
        icon: formData.value.icon,
        sort: formData.value.sort,
        is_enabled: formData.value.is_enabled,
        children: [],
      };
      if (!appendMenu(tree, selectedParentId.value, menu)) tree.push(menu);
      ElMessage.success('创建成功');
    }

    await updateThemeMenus(props.themeSlug, flushMenus());
    dialogVisible.value = false;
    emit('refresh');
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error(error.message || '操作失败');
    }
  } finally {
    submitLoading.value = false;
  }
};

watch(dialogVisible, val => {
  if (val) {
    cleanupIconBlob();
    pendingFile.value = null;
    if (editingMenu.value) {
      formData.value = { ...editingMenu.value };
      selectedParentId.value = parentMenu.value?.id ?? null;
    } else {
      formData.value = {
        id: 0,
        title: '',
        url: '',
        icon: '',
        sort: 5,
        is_enabled: true,
        children: [],
      };
      selectedParentId.value = parentMenu.value?.id ?? null;
    }
    if (editingMenu.value) {
      fetchParentMenuOptions();
    }
  }
});

watch(
  () => [props.themeSlug, props.menus, menuTypeOptions.value.map(option => option.value).join(',')],
  () => {
    const firstType = menuTypeOptions.value[0]?.value || 'aggregate';
    if (!menuTypeOptions.value.some(option => option.value === selectedType.value)) {
      selectedType.value = firstType;
    }
    fetchMenuTree();
  },
  { immediate: true }
);
</script>

<style scoped lang="scss">
.theme-menu-panel {
  .menu-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
    flex-wrap: wrap;
  }

  .menu-actions {
    display: flex;
    align-items: center;
    gap: 12px;

    :deep(.el-button + .el-button) {
      margin-left: 0;
    }
  }

  .menu-table {
    width: 100%;
  }

  .menu-title {
    display: flex;
    align-items: center;
  }

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

  .empty-value {
    color: #909399;
  }

  .menu-type-segmented {
    .segmented-item {
      display: flex;
      align-items: center;
      gap: 4px;
    }

    .segmented-icon {
      display: none;
    }
  }

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
}

.form-info {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 16px;
  margin-bottom: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
  border: 1px solid #e4e7ed;

  .info-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    flex: 1;
    text-align: center;

    .info-label {
      font-size: 12px;
      color: #909399;
    }

    .info-value {
      font-size: 14px;
      color: #303133;
      font-weight: 500;
    }
  }
}

.icon-input-wrapper {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12px;

  .el-input {
    flex: 1;
  }

  .icon-preview {
    width: 40px;
    height: 40px;
    border: 1px solid #e4e7ed;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    background-color: #f5f7fa;
    flex-shrink: 0;

    img {
      max-width: 100%;
      max-height: 100%;
      object-fit: contain;
    }

    i {
      font-size: 24px;
      color: #606266;
    }
  }
}

@media (max-width: 500px) {
  .theme-menu-panel {
    .menu-type-segmented {
      .segmented-icon {
        display: inline-flex;
      }

      .segmented-text {
        display: none;
      }
    }
  }
}
</style>
