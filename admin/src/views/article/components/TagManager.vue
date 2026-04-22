<template>
  <el-dialog
    v-model="visible"
    title="标签管理"
    width="90%"
    style="max-width: 600px"
    :align-center="true"
  >
    <el-table v-loading="loading" :data="list" style="margin: 20px 0" max-height="350">
      <el-table-column prop="name" label="标签名称" min-width="100" show-overflow-tooltip />
      <el-table-column prop="description" label="描述" min-width="120" show-overflow-tooltip />
      <el-table-column prop="count" label="文章数" width="80" align="center" />
      <el-table-column label="操作" width="100" align="center" fixed="right">
        <template #header>
          <el-button type="primary" plain size="small" @click="openForm()">新增</el-button>
        </template>
        <template #default="{ row }">
          <el-button type="primary" link size="small" @click="openForm(row)">编辑</el-button>
          <el-button type="danger" link size="small" @click="remove(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="formVisible"
      :title="current.id ? '编辑' : '新增'"
      width="90%"
      style="max-width: 400px"
      append-to-body
    >
      <el-form :model="current" label-width="80px">
        <el-form-item label="名称" required>
          <el-input v-model="current.name" placeholder="请输入标签名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="current.description" type="textarea" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="formVisible = false">取消</el-button>
        <el-button type="primary" @click="save">确定</el-button>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { getTags, createTag, updateTag, deleteTag } from '@/api/tag';
import type { Tag } from '@/types/tag';

const props = defineProps<{ modelValue: boolean }>();

const emit = defineEmits(['update:modelValue']);

const visible = computed({
  get: () => props.modelValue,
  set: val => emit('update:modelValue', val),
});

const loading = ref(false);
const list = ref<Tag[]>([]);

const formVisible = ref(false);
const current = ref<Partial<Tag>>({ id: 0, name: '', description: '' });

// 弹窗打开时加载数据（immediate 确保懒挂载组件首次打开时也能加载）
watch(
  visible,
  val => {
    if (val) loadData();
  },
  { immediate: true }
);

// 加载标签列表
async function loadData() {
  loading.value = true;
  try {
    const res = await getTags();
    list.value = res.list;
  } catch (_error) {
    ElMessage.error('加载标签列表失败');
  } finally {
    loading.value = false;
  }
}

// 打开表单
function openForm(row?: Tag) {
  if (row) {
    current.value = { ...row };
  } else {
    current.value = { id: 0, name: '', description: '' };
  }
  formVisible.value = true;
}

async function remove(row: Tag) {
  try {
    await ElMessageBox.confirm('确定要删除这个标签吗？');
    await deleteTag(row.id);
    await loadData();
    ElMessage.success('删除成功');
  } catch {}
}

async function save() {
  if (!current.value.name?.trim()) {
    return ElMessage.warning('请输入标签名称');
  }

  loading.value = true;
  try {
    if (current.value.id) {
      await updateTag(current.value.id, current.value);
    } else {
      await createTag(current.value);
    }
    await loadData();
    formVisible.value = false;
    ElMessage.success('保存成功');
  } catch (_error) {
    ElMessage.error('保存失败');
  } finally {
    loading.value = false;
  }
}
</script>
