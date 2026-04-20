<template>
  <el-dialog v-model="visible" :title="dialogTitle" width="600px" :close-on-click-modal="false">
    <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
      <el-form-item label="友链名称" prop="name">
        <el-input v-model="formData.name" placeholder="请输入友链名称" clearable />
      </el-form-item>

      <el-form-item label="链接地址" prop="url">
        <el-input
          v-model="formData.url"
          placeholder="请输入链接地址，如：https://example.com"
          clearable
        >
          <template #append>
            <el-button
              type="primary"
              @click="handleParseLink"
              :disabled="!formData.url || parseLoading"
            >
              {{ parseLoading ? '解析中...' : '解析' }}
            </el-button>
          </template>
        </el-input>
      </el-form-item>

      <el-form-item label="友链描述" prop="description">
        <el-input
          v-model="formData.description"
          type="textarea"
          :rows="3"
          placeholder="请输入友链描述"
          clearable
        />
      </el-form-item>

      <el-form-item label="RSS地址" prop="rss_url">
        <el-input
          v-model="formData.rss_url"
          placeholder="请输入RSS订阅地址，如：https://example.com/feed"
          clearable
        />
      </el-form-item>

      <el-row :gutter="20">
        <el-col :span="9">
          <el-form-item label="友链头像" prop="avatar">
            <ImageUploader
              ref="avatarUploaderRef"
              v-model="formData.avatar"
              upload-type="友情链接A"
              width="120px"
              height="120px"
            />
          </el-form-item>
        </el-col>

        <el-col :span="15">
          <el-form-item label="网站截图" prop="screenshot">
            <ImageUploader
              ref="screenshotUploaderRef"
              v-model="formData.screenshot"
              upload-type="友情链接S"
              width="213px"
              height="120px"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <el-row :gutter="20">
        <el-col :span="12">
          <el-form-item label="友链类型" prop="type_id">
            <el-select v-model="formData.type_id" placeholder="请选择友链类型" style="width: 100%">
              <el-option
                v-for="type in friendTypeOptions"
                :key="type.id"
                :label="type.name"
                :value="type.id"
              >
                <span>{{ type.name }}</span>
                <el-tag v-if="!type.is_visible" size="small" type="info" style="margin-left: 8px"
                  >已隐藏</el-tag
                >
              </el-option>
            </el-select>
          </el-form-item>
        </el-col>

        <el-col :span="12">
          <el-form-item label="排序值" prop="sort">
            <el-input-number
              v-model="formData.sort"
              :min="1"
              :max="10"
              placeholder="排序值，范围1-10，数值越大排序越靠前"
              style="width: 100%"
            />
          </el-form-item>
        </el-col>
      </el-row>

      <!-- 只有编辑时才显示状态控制 -->
      <el-row v-if="isEdit" :gutter="20">
        <el-col :span="8">
          <el-form-item label="已失效">
            <el-switch v-model="formData.is_invalid" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="待审核">
            <el-switch v-model="formData.is_pending" />
          </el-form-item>
        </el-col>
        <el-col :span="8">
          <el-form-item label="忽略检查">
            <el-switch v-model="formData.ignoreCheck" />
          </el-form-item>
        </el-col>
      </el-row>
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
import { ref, computed, watch } from 'vue';
import { ElMessage, type FormInstance, type FormRules } from 'element-plus';
import type { Friend, FriendType, CreateFriendRequest, UpdateFriendRequest } from '@/types/friend';
import { createFriend, updateFriend, getFriendTypes } from '@/api/friend';
import { fetchLinkInfo } from '@/api/tools';
import request from '@/utils/request';
import ImageUploader from '@/components/common/ImageUploader.vue';
const props = defineProps<{
  modelValue: boolean;
  editFriend?: Friend | null;
}>();

const emit = defineEmits(['update:modelValue', 'success']);

const visible = computed({
  get: () => props.modelValue,
  set: val => emit('update:modelValue', val),
});

const isEdit = computed(() => !!props.editFriend);
const dialogTitle = computed(() => (isEdit.value ? '编辑友链' : '新增友链'));

const submitLoading = ref(false);
const parseLoading = ref(false);
const formRef = ref<FormInstance>();
const avatarUploaderRef = ref<InstanceType<typeof ImageUploader>>();
const screenshotUploaderRef = ref<InstanceType<typeof ImageUploader>>();

// 友链类型选项
const friendTypeOptions = ref<FriendType[]>([]);

// 表单数据类型（编辑时使用）
interface FriendFormData {
  name: string;
  url: string;
  description?: string;
  avatar?: string;
  screenshot?: string;
  sort?: number;
  type_id: number | null;
  is_invalid?: boolean;
  is_pending?: boolean;
  rss_url?: string;
  ignoreCheck?: boolean;
}

// 表单数据
const formData = ref<FriendFormData>({
  name: '',
  url: '',
  description: '',
  avatar: '',
  screenshot: '',
  sort: 5,
  type_id: null,
  is_invalid: false,
  is_pending: false,
  rss_url: '',
  ignoreCheck: false,
});

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入友链名称', trigger: 'blur' },
    { min: 1, max: 50, message: '友链名称长度为1-50个字符', trigger: 'blur' },
  ],
  url: [
    { required: true, message: '请输入链接地址', trigger: 'blur' },
    {
      pattern: /^https?:\/\/.+/,
      message: '请输入正确的链接地址，必须以http://或https://开头',
      trigger: 'blur',
    },
    { max: 255, message: '链接地址长度不能超过255个字符', trigger: 'blur' },
  ],
  description: [{ max: 500, message: '描述长度不能超过500个字符', trigger: 'blur' }],
  type_id: [{ required: true, message: '请选择友链类型', trigger: 'change' }],
  sort: [
    { required: true, message: '请输入排序值', trigger: 'blur' },
    {
      type: 'number',
      min: 1,
      max: 10,
      message: '排序值必须在 1-10 之间',
      trigger: 'blur',
    },
  ],
};

// 加载友链类型列表
const loadFriendTypes = async () => {
  try {
    const res = await getFriendTypes();
    friendTypeOptions.value = res.list;
  } catch (error) {
    console.error('加载友链类型失败:', error);
  }
};

// 重置表单数据
const resetFormData = () => {
  formData.value = {
    name: '',
    url: '',
    description: '',
    avatar: '',
    screenshot: '',
    sort: 5,
    type_id: null,
    is_invalid: false,
    is_pending: false,
    rss_url: '',
    ignoreCheck: false,
  };
  // 清除表单验证状态
  setTimeout(() => {
    formRef.value?.clearValidate();
  }, 0);
};

// 弹窗打开时加载友链类型（immediate 确保懒挂载组件首次打开时也能加载）
watch(
  visible,
  val => {
    if (val) loadFriendTypes();
  },
  { immediate: true }
);

// 监听编辑友链变化
watch(
  () => props.editFriend,
  friend => {
    if (friend) {
      formData.value = {
        name: friend.name,
        url: friend.url,
        description: friend.description || '',
        avatar: friend.avatar || '',
        screenshot: friend.screenshot || '',
        sort: friend.sort || 5,
        type_id: friend.type_id ?? null,
        is_invalid: friend.is_invalid ?? false,
        is_pending: friend.is_pending ?? false,
        rss_url: friend.rss_url || '',
        ignoreCheck: friend.accessible === -1,
      };
      // 清除表单验证
      setTimeout(() => {
        formRef.value?.clearValidate();
      }, 0);
    } else {
      resetFormData();
    }
  },
  { immediate: true }
);

// 取消
const handleCancel = () => {
  resetFormData();
  visible.value = false;
};

// 下载预览图片并返回完整信息
interface PreviewImageInfo {
  blobUrl: string;
  file: File;
}

const downloadPreviewImage = async (
  url: string,
  filename: string
): Promise<PreviewImageInfo | null> => {
  try {
    // 使用更长的超时时间（60秒）下载图片
    const response = await request.post('/admin/tools/download-image', { url }, { timeout: 60000 });

    // 简化：直接使用base64创建Blob
    const blob = await fetch(`data:image/png;base64,${response.data}`).then(res => res.blob());

    // 获取content-type
    const contentType = response.headers?.['content-type'] || 'image/png';

    // 创建File对象
    const file = new File([blob], filename, { type: contentType });

    // 创建Blob URL用于显示
    const blobUrl = URL.createObjectURL(blob);

    return { blobUrl, file };
  } catch (error) {
    console.error('下载图片失败:', error);
    return null;
  }
};

// 解析链接
const handleParseLink = async () => {
  if (!formData.value.url) {
    ElMessage.warning('请输入链接地址');
    return;
  }

  try {
    parseLoading.value = true;
    const result = await fetchLinkInfo({ url: formData.value.url });

    // 更新表单数据
    formData.value = {
      ...formData.value,
      name: result.title || formData.value.name,
      description: result.description || formData.value.description,
    };

    // 下载并设置favicon（如果存在）
    if (result.favicon && avatarUploaderRef.value) {
      const previewInfo = await downloadPreviewImage(result.favicon, 'avatar.png');
      if (previewInfo) {
        formData.value.avatar = previewInfo.blobUrl;
        avatarUploaderRef.value.setPendingFile?.(previewInfo.file);
      }
    }

    ElMessage.success('解析成功');
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : '解析失败');
  } finally {
    parseLoading.value = false;
  }
};

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    submitLoading.value = true;

    // 收集待上传的图片（并行上传）
    const uploadPromises: Promise<void>[] = [];
    const uploadFields: Array<{
      uploader: typeof avatarUploaderRef.value;
      field: 'avatar' | 'screenshot';
    }> = [
      { uploader: avatarUploaderRef.value, field: 'avatar' },
      { uploader: screenshotUploaderRef.value, field: 'screenshot' },
    ];

    for (const { uploader, field } of uploadFields) {
      if (!uploader) continue;
      const hasPending = uploader.getPendingCount() > 0;
      const hasBlobUrl = formData.value[field]?.startsWith('blob:');
      if (!hasPending && !hasBlobUrl) continue;

      uploadPromises.push(
        uploader
          .uploadPendingFile()
          .then(uploadedUrl => {
            if (uploadedUrl) formData.value[field] = uploadedUrl;
          })
          .catch((error: unknown) => {
            ElMessage.error(
              (error as Error)?.message || `${field === 'avatar' ? '头像' : '截图'}上传失败`
            );
          })
      );
    }

    // 等待所有上传完成（使用 allSettled 确保即使部分失败也继续）
    if (uploadPromises.length > 0) {
      const results = await Promise.allSettled(uploadPromises);
      const failedCount = results.filter(r => r.status === 'rejected').length;
      if (failedCount > 0) {
        submitLoading.value = false;
        ElMessage.error(`${failedCount} 个文件上传失败，请重试`);
        return;
      }
    }

    if (isEdit.value && props.editFriend) {
      // 编辑友链
      const updateData: UpdateFriendRequest = {
        name: formData.value.name,
        url: formData.value.url,
        description: formData.value.description,
        avatar: formData.value.avatar,
        screenshot: formData.value.screenshot,
        sort: formData.value.sort,
        type_id: formData.value.type_id ?? undefined,
        is_invalid: formData.value.is_invalid,
        is_pending: formData.value.is_pending,
        rss_url: formData.value.rss_url,
        accessible: formData.value.ignoreCheck ? -1 : 0,
      };
      await updateFriend(props.editFriend.id, updateData);
      ElMessage.success('更新友链成功');
    } else {
      // 新增友链
      const createData: CreateFriendRequest = {
        name: formData.value.name,
        url: formData.value.url,
        description: formData.value.description,
        avatar: formData.value.avatar,
        screenshot: formData.value.screenshot,
        sort: formData.value.sort,
        type_id: formData.value.type_id!,
        rss_url: formData.value.rss_url,
      };
      await createFriend(createData);
      ElMessage.success('创建友链成功');
    }

    resetFormData();
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
