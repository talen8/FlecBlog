<template>
  <div class="image-uploader">
    <div class="uploader-container" :style="{ width, height }">
      <el-upload
        class="uploader-box"
        :show-file-list="false"
        :http-request="handleUpload"
        accept="image/*"
        :disabled="disabled"
      >
        <img v-if="imageUrl" :src="imageUrl" class="preview-image" />
        <div v-else class="upload-placeholder">
          <el-icon :size="40">
            <Plus />
          </el-icon>
        </div>
      </el-upload>

      <div v-if="imageUrl && !disabled" class="delete-btn" @click.stop="handleDelete" title="删除">
        <el-icon>
          <Delete />
        </el-icon>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { ElMessage, type UploadRequestOptions } from 'element-plus';
import { Plus, Delete } from '@element-plus/icons-vue';
import { uploadFile } from '@/api/file';

export interface ImageUploaderProps {
  modelValue?: string; // 图片 URL
  uploadType?: string; // 上传用途（如：用户头像、文章封面）
  width?: string; // 宽度
  height?: string; // 高度
  disabled?: boolean;
}

const props = withDefaults(defineProps<ImageUploaderProps>(), {
  uploadType: '图片',
  width: '120px',
  height: '120px',
  disabled: false,
});

const emit = defineEmits<{
  'update:modelValue': [value: string];
}>();

const pendingFile = ref<File | null>(null); // 待上传的文件
const previewUrl = ref<string>(''); // 本地预览 URL

// 图片 URL（本地预览或已上传）
const imageUrl = computed(() => {
  // 如果有本地预览，优先显示本地预览
  if (previewUrl.value) return previewUrl.value;
  return props.modelValue || '';
});

// 上传处理（延迟上传：只做本地预览）
const handleUpload = async (options: UploadRequestOptions): Promise<void> => {
  if (props.disabled) {
    return Promise.resolve();
  }

  const file = options.file as File;

  // 验证文件类型
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件');
    return Promise.reject();
  }

  // 清理旧的预览 URL
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value);
  }

  // 保存文件和创建本地预览
  pendingFile.value = file;
  previewUrl.value = URL.createObjectURL(file);

  return Promise.resolve();
};

// 删除文件
const handleDelete = () => {
  if (props.disabled) return;

  // 清理本地预览
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value);
    previewUrl.value = '';
  }
  pendingFile.value = null;
  emit('update:modelValue', '');
};

// 暴露上传方法供父组件调用
const uploadPendingFile = async (): Promise<string | null> => {
  if (props.disabled || !pendingFile.value) return null;

  const loading = ElMessage.info({ message: '正在上传...', duration: 0 });
  try {
    const result = await uploadFile(pendingFile.value, props.uploadType);

    // 清理本地预览
    if (previewUrl.value) {
      URL.revokeObjectURL(previewUrl.value);
      previewUrl.value = '';
    }
    pendingFile.value = null;

    // 更新值
    emit('update:modelValue', result.file_url);
    loading.close();
    return result.file_url;
  } catch (error: unknown) {
    loading.close();
    ElMessage.error((error as Error)?.message || '上传失败');
    throw error;
  }
};

// 获取待上传文件数量
const getPendingCount = () => {
  return pendingFile.value ? 1 : 0;
};

// 设置待上传文件
const setPendingFile = (file: File) => {
  pendingFile.value = file;
};

// 暴露方法给父组件
defineExpose({
  uploadPendingFile,
  getPendingCount,
  setPendingFile,
});
</script>

<style scoped lang="scss">
.image-uploader {
  display: inline-block;

  .uploader-container {
    position: relative;
    display: inline-block;

    .delete-btn {
      position: absolute;
      top: 8px;
      right: 8px;
      z-index: 10;
      width: 28px;
      height: 28px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      background: rgba(0, 0, 0, 0.5);
      border-radius: 4px;
      transition: all 0.2s;

      .el-icon {
        color: #fff;
        font-size: 16px;
      }

      &:hover {
        background: rgba(245, 108, 108, 0.9);
        transform: scale(1.1);
      }
    }
  }

  .uploader-box {
    width: 100%;
    height: 100%;
    box-sizing: border-box;
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    overflow: hidden;
    transition: var(--el-transition-duration-fast);

    &:hover {
      border-color: var(--el-color-primary);
    }

    :deep(.el-upload) {
      width: 100%;
      height: 100%;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
    }

    .preview-image {
      width: 100%;
      height: 100%;
      object-fit: cover;
      display: block;
    }

    .upload-placeholder {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      color: #8c939d;
      text-align: center;
    }
  }
}
</style>
