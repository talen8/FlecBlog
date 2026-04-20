<template>
  <el-dialog
    v-model="visible"
    :title="isEdit ? '编辑菜单' : '新增菜单'"
    width="600px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
      <!-- 提示信息 -->
      <div class="form-info">
        <div class="info-item">
          <span class="info-label">菜单类型</span>
          <span class="info-value">{{ getMenuTypeLabel(formData.type) }}</span>
        </div>
        <div v-if="parentMenu && !isEdit" class="info-item">
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
                <el-icon>
                  <Upload />
                </el-icon>
                上传
              </el-button>
            </template>
          </el-input>
          <div v-if="formData.icon" class="icon-preview">
            <!-- RemixIcon 图标预览 -->
            <i v-if="isRemixIcon(formData.icon)" :class="formData.icon"></i>
            <!-- 图片预览 -->
            <img v-else :src="formData.icon" alt="图标预览" @error="handleIconError" />
          </div>
        </div>
      </el-form-item>

      <!-- 编辑时显示父菜单选择器 -->
      <el-form-item v-if="isEdit" label="父菜单" prop="parent_id">
        <el-select
          v-model="formData.parent_id"
          :placeholder="hasChildren ? '包含子菜单，无法设置' : '请选择父菜单'"
          :disabled="hasChildren"
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
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="submitLoading" @click="handleSubmit"> 确定 </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { ElMessage } from 'element-plus';
import { Upload } from '@element-plus/icons-vue';
import type { FormInstance, FormRules } from 'element-plus';
import type { MenuTreeNode, CreateMenuRequest, UpdateMenuRequest, MenuType } from '@/types/menu';
import { createMenu, updateMenu, getMenuTree } from '@/api/menu';
import { uploadFile } from '@/api/file';

interface Props {
  modelValue: boolean;
  editMenu?: MenuTreeNode | null;
  parentMenu?: MenuTreeNode | null;
  currentType?: string;
}
const props = withDefaults(defineProps<Props>(), {
  editMenu: null,
  parentMenu: null,
  currentType: 'navigation',
});

const emit = defineEmits<{
  'update:modelValue': [value: boolean];
  success: [];
}>();

const visible = computed({
  get: () => props.modelValue,
  set: val => emit('update:modelValue', val),
});

const formRef = ref<FormInstance>();
const submitLoading = ref(false);
const isEdit = computed(() => !!props.editMenu);

// 判断当前菜单是否有子菜单（用于禁用父菜单选择器）
const hasChildren = computed(() => {
  return props.editMenu?.children && props.editMenu.children.length > 0;
});

// 父菜单列表
const parentMenuOptions = ref<MenuTreeNode[]>([]);

// 图标上传相关
interface IconItem {
  type: 'file' | 'url';
  file?: File;
  url: string;
}
const iconItem = ref<IconItem | null>(null);

// 表单数据
const formData = ref<CreateMenuRequest | UpdateMenuRequest>({
  title: '',
  type: 'navigation' as MenuType,
  url: '',
  icon: '',
  sort: 5,
  is_enabled: true,
  parent_id: null,
});

// 表单验证规则
const rules: FormRules = {
  title: [
    { message: '请输入菜单标题', trigger: 'blur' },
    { min: 1, max: 100, message: '长度在 1 到 100 个字符', trigger: 'blur' },
  ],
  type: [{ message: '请选择菜单类型', trigger: 'change' }],
  url: [{ max: 500, message: '链接地址不能超过 500 个字符', trigger: 'blur' }],
  icon: [{ max: 500, message: '图标不能超过 500 个字符', trigger: 'blur' }],
};

// 获取菜单类型标签
const getMenuTypeLabel = (type: MenuType) => {
  const typeMap: Record<MenuType, string> = {
    aggregate: '聚合菜单',
    navigation: '导航菜单',
    footer: '页脚菜单',
  };
  return typeMap[type] || type;
};

// 判断是否是 RemixIcon 图标类名
const isRemixIcon = (icon: string) => {
  return icon && icon.startsWith('ri-');
};

// 清理图标Blob URL
const cleanupIconBlob = () => {
  if (iconItem.value?.type === 'file' && iconItem.value.url.startsWith('blob:')) {
    URL.revokeObjectURL(iconItem.value.url);
  }
};

// 上传图标
const handleIconUpload = () => {
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = 'image/*';
  input.onchange = e => {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;

    // 清理旧的Blob URL
    cleanupIconBlob();

    // 创建临时预览URL
    const blobUrl = URL.createObjectURL(file);
    iconItem.value = {
      type: 'file',
      file,
      url: blobUrl,
    };
    formData.value.icon = blobUrl;
  };
  input.click();
};

// 图标加载错误处理
const handleIconError = (e: Event) => {
  const target = e.target as HTMLImageElement;
  target.style.display = 'none';
  ElMessage.warning('图标加载失败');
};

// 获取所有子菜单ID（用于防止循环引用）
const getAllChildrenIds = (menu: MenuTreeNode): number[] => {
  const ids: number[] = [menu.id];
  if (menu.children) {
    menu.children.forEach(child => {
      ids.push(...getAllChildrenIds(child));
    });
  }
  return ids;
};

// 获取父菜单选项列表
const fetchParentMenuOptions = async () => {
  try {
    const type = props.editMenu?.type || props.currentType;
    const allMenus = await getMenuTree(type);

    // 只显示主菜单（parent_id 为 null 的菜单）
    let options = allMenus.filter(menu => menu.parent_id === null);

    // 编辑时，过滤掉当前菜单及其所有子菜单（防止循环引用）
    if (props.editMenu) {
      const excludeIds = getAllChildrenIds(props.editMenu);
      options = options.filter(menu => !excludeIds.includes(menu.id));
    }

    parentMenuOptions.value = options;
  } catch (_error) {
    ElMessage.error('获取父菜单列表失败');
  }
};

// 初始化表单数据
const initFormData = () => {
  // 清理旧的Blob URL
  cleanupIconBlob();
  iconItem.value = null;

  if (props.editMenu) {
    const menu = props.editMenu;
    formData.value = {
      title: menu.title,
      type: menu.type,
      url: menu.url || '',
      icon: menu.icon || '',
      sort: menu.sort || 0,
      is_enabled: menu.is_enabled,
      parent_id: menu.parent_id,
    };

    // 如果有图标URL，标记为网络图标
    if (menu.icon) {
      iconItem.value = {
        type: 'url',
        url: menu.icon,
      };
    }
  } else {
    // 新增菜单时，重置为默认值
    formData.value = {
      title: '',
      type: props.parentMenu
        ? props.parentMenu.type
        : (props.currentType as MenuType) || 'navigation',
      url: '',
      icon: '',
      sort: 5,
      is_enabled: true,
      parent_id: props.parentMenu?.id || null,
    };
  }
};

// 监听对话框打开状态，打开时初始化表单
watch(
  () => props.modelValue,
  async newVal => {
    if (newVal) {
      // 对话框打开时初始化表单数据
      initFormData();
      // 编辑时加载父菜单选项
      if (isEdit.value) {
        await fetchParentMenuOptions();
      }
    }
  }
);

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    submitLoading.value = true;

    // 如果图标是本地文件，先上传
    if (iconItem.value?.type === 'file' && iconItem.value.file) {
      const result = await uploadFile(iconItem.value.file, '菜单图标');
      formData.value.icon = result.file_url;
    }

    if (isEdit.value && props.editMenu) {
      await updateMenu(props.editMenu.id, formData.value as UpdateMenuRequest);
      ElMessage.success('更新成功');
    } else {
      await createMenu(formData.value as CreateMenuRequest);
      ElMessage.success('创建成功');
    }

    emit('success');
    handleClose();
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error(error.message || '操作失败');
    }
  } finally {
    submitLoading.value = false;
  }
};

// 关闭对话框
const handleClose = () => {
  // 清理Blob URL
  cleanupIconBlob();
  iconItem.value = null;

  // 重置表单验证状态
  formRef.value?.clearValidate();

  emit('update:modelValue', false);
};
</script>

<style scoped lang="scss">
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
</style>
