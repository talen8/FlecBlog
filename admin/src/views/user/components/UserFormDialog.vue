<template>
  <el-dialog
    v-model="visible"
    :title="dialogTitle"
    width="90%"
    style="max-width: 600px"
    :close-on-click-modal="false"
  >
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
      <el-form-item label="邮箱" prop="email">
        <el-input v-model="formData.email" placeholder="请输入用户邮箱" clearable />
      </el-form-item>

      <el-form-item label="昵称" prop="nickname">
        <el-input v-model="formData.nickname" placeholder="请输入用户昵称" clearable />
      </el-form-item>

      <el-form-item label="网站" prop="website">
        <el-input v-model="formData.website" placeholder="请输入用户网站（可选）" clearable />
      </el-form-item>

      <el-form-item label="徽章" prop="badge">
        <el-input v-model="formData.badge" placeholder="请输入用户徽章（可选）" clearable />
      </el-form-item>

      <el-form-item label="密码" prop="password">
        <el-input
          v-model="formData.password"
          type="password"
          :placeholder="isEdit ? '留空则不修改密码' : '请输入密码'"
          show-password
          clearable
          autocomplete="new-password"
        />
      </el-form-item>

      <el-form-item label="头像" prop="avatar">
        <ImageUploader
          ref="avatarUploaderRef"
          v-model="formData.avatar"
          upload-type="用户头像"
          width="120px"
          height="120px"
        />
      </el-form-item>

      <el-form-item label="角色" prop="role">
        <el-select v-model="formData.role" placeholder="请选择角色" style="width: 100%">
          <el-option
            v-for="option in roleOptions"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          />
        </el-select>
      </el-form-item>

      <el-form-item v-if="isEdit" label="状态" prop="is_enabled">
        <el-switch v-model="formData.is_enabled" active-text="启用" inactive-text="禁用" />
      </el-form-item>
    </el-form>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading"> 确定 </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import type { User, CreateUserRequest, UpdateUserRequest } from '@/types/user';
import { createUser, updateUser } from '@/api/user';
import ImageUploader from '@/components/common/ImageUploader.vue';
import { isSuperAdmin } from '@/utils/auth';

type UserRoleOption = CreateUserRequest['role'];

interface UserFormData {
  password: string;
  email: string;
  nickname: string;
  avatar?: string;
  badge?: string;
  website?: string;
  role: UserRoleOption;
  is_enabled?: boolean;
}

const props = defineProps<{
  modelValue: boolean;
  editUser?: User | null;
}>();

const emit = defineEmits(['update:modelValue', 'success']);

const visible = computed({
  get: () => props.modelValue,
  set: val => emit('update:modelValue', val),
});

const isEdit = computed(() => !!props.editUser);
const dialogTitle = computed(() => (isEdit.value ? '编辑用户' : '新增用户'));
const canAssignAdminRole = computed(() => isSuperAdmin());
const roleOptions = computed<{ label: string; value: UserRoleOption }[]>(() => {
  const options: { label: string; value: UserRoleOption }[] = [
    { label: '普通用户', value: 'user' },
    { label: '访客', value: 'guest' },
  ];

  if (canAssignAdminRole.value) {
    return [
      { label: '超级管理员', value: 'super_admin' },
      { label: '管理员', value: 'admin' },
      ...options,
    ];
  }

  return options;
});

const submitLoading = ref(false);
const formRef = ref<FormInstance>();
const avatarUploaderRef = ref<InstanceType<typeof ImageUploader>>();

const getDefaultRole = (): UserRoleOption => roleOptions.value[0]?.value || 'user';

// 表单数据
const formData = ref<UserFormData>({
  password: '',
  email: '',
  nickname: '',
  avatar: '',
  badge: '',
  website: '',
  role: getDefaultRole(),
  is_enabled: true,
});

// 表单验证规则
const formRules = computed<FormRules>(() => ({
  password: [
    {
      required: !isEdit.value,
      message: '请输入密码',
      trigger: 'blur',
    },
    {
      min: 6,
      max: 32,
      message: '密码长度为6-32个字符',
      trigger: 'blur',
    },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' },
  ],
  website: [{ type: 'url', message: '请输入正确的网址格式', trigger: 'blur' }],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 32, message: '昵称长度为2-32个字符', trigger: 'blur' },
  ],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
}));

// 重置表单数据
const resetFormData = () => {
  formData.value = {
    password: '',
    email: '',
    nickname: '',
    avatar: '',
    badge: '',
    website: '',
    role: getDefaultRole(),
    is_enabled: true,
  };
};

// 监听编辑用户变化
watch(
  () => props.editUser,
  user => {
    if (user) {
      formData.value = {
        password: '',
        email: user.email || '',
        nickname: user.nickname || '',
        avatar: user.avatar || '',
        badge: user.badge || '',
        website: user.website || '',
        role: user.role as UserRoleOption,
        is_enabled: user.is_enabled,
      };
    } else {
      resetFormData();
    }

    // 清除表单验证
    setTimeout(() => {
      formRef.value?.clearValidate();
    }, 0);
  },
  { immediate: true }
);

watch(roleOptions, options => {
  if (!options.some(option => option.value === formData.value.role)) {
    formData.value.role = getDefaultRole();
  }
});

// 取消
const handleCancel = () => {
  visible.value = false;
};

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    submitLoading.value = true;

    // 如果有新选择的头像文件，先上传
    if (avatarUploaderRef.value?.getPendingCount()) {
      try {
        const uploadedUrl = await avatarUploaderRef.value.uploadPendingFile();
        if (uploadedUrl) {
          formData.value.avatar = uploadedUrl;
        }
      } catch (error: unknown) {
        submitLoading.value = false;
        ElMessage.error((error as Error)?.message || '头像上传失败');
        return;
      }
    }

    if (isEdit.value && props.editUser) {
      // 编辑用户
      const updateData: UpdateUserRequest = {
        email: formData.value.email,
        nickname: formData.value.nickname,
        avatar: formData.value.avatar,
        badge: formData.value.badge,
        website: formData.value.website,
        role: formData.value.role,
        is_enabled: formData.value.is_enabled,
      };

      // 如果填写了新密码，则包含密码字段
      if (formData.value.password) {
        updateData.password = formData.value.password;
      }

      await updateUser(props.editUser.id, updateData);
      ElMessage.success('更新用户成功');
    } else {
      // 新增用户
      const createData: CreateUserRequest = {
        password: formData.value.password,
        email: formData.value.email,
        nickname: formData.value.nickname,
        avatar: formData.value.avatar,
        badge: formData.value.badge,
        website: formData.value.website,
        role: formData.value.role,
      };
      await createUser(createData);
      ElMessage.success('创建用户成功');
    }

    visible.value = false;
    emit('success');
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error(error.message);
    }
  } finally {
    submitLoading.value = false;
  }
};
</script>
