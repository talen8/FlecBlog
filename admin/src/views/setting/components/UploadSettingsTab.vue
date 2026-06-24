<template>
  <el-form :model="form" label-width="120px" class="setting-form">
    <el-divider content-position="left">基础配置</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('storage_type') }">存储类型</span>
      </template>
      <el-select
        v-model="form.storage_type"
        placeholder="选择存储类型"
        style="width: 220px"
        :disabled="loading"
      >
        <el-option label="本地存储" value="local" />
        <el-option label="亚马逊 S3" value="s3" />
        <el-option label="阿里云 OSS" value="oss" />
        <el-option label="腾讯云 COS" value="cos" />
        <el-option label="七牛云 Kodo" value="kodo" />
        <el-option label="Cloudflare R2" value="r2" />
        <el-option label="MinIO" value="minio" />
      </el-select>
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('max_file_size') }">最大文件大小</span>
      </template>
      <el-input-number v-model="form.max_file_size" :min="0" :step="1" :disabled="loading" />
      <span class="unit-tip">MB</span>
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('path_pattern') }">文件命名</span>
      </template>
      <el-input
        v-model="form.path_pattern"
        placeholder="{timestamp}_{random}{ext}"
        :disabled="loading"
      />
    </el-form-item>

    <template v-if="form.storage_type !== 'local'">
      <el-form-item>
        <template #label>
          <span :class="{ 'field-modified': isFieldModified('access_key') }">{{
            accessLabel
          }}</span>
        </template>
        <el-input
          v-model="form.access_key"
          :placeholder="accessPlaceholder"
          clearable
          :disabled="loading"
        />
      </el-form-item>

      <el-form-item>
        <template #label>
          <span :class="{ 'field-modified': isFieldModified('secret_key') }">{{
            secretLabel
          }}</span>
        </template>
        <el-input
          v-model="form.secret_key"
          type="password"
          show-password
          :placeholder="secretPlaceholder"
          clearable
          :disabled="loading"
          autocomplete="new-password"
        />
      </el-form-item>

      <el-form-item v-if="showRegion">
        <template #label>
          <span :class="{ 'field-modified': isFieldModified('region') }">地域</span>
        </template>
        <el-input
          v-model="form.region"
          :placeholder="regionPlaceholder"
          clearable
          :disabled="loading"
        />
      </el-form-item>

      <el-form-item>
        <template #label>
          <span :class="{ 'field-modified': isFieldModified('bucket') }">存储桶</span>
        </template>
        <el-input
          v-model="form.bucket"
          :placeholder="bucketPlaceholder"
          clearable
          :disabled="loading"
        />
      </el-form-item>

      <el-form-item v-if="showEndpoint">
        <template #label>
          <span :class="{ 'field-modified': isFieldModified('endpoint') }">服务端点</span>
        </template>
        <el-input
          v-model="form.endpoint"
          :placeholder="endpointPlaceholder"
          clearable
          :disabled="loading"
        />
      </el-form-item>

      <el-form-item>
        <template #label>
          <span :class="{ 'field-modified': isFieldModified('domain') }">自定义域名</span>
        </template>
        <el-input
          v-model="form.domain"
          :placeholder="domainPlaceholder"
          clearable
          :disabled="loading"
        />
      </el-form-item>

      <el-form-item v-if="showUseSSL">
        <template #label>
          <span :class="{ 'field-modified': isFieldModified('use_ssl') }">启用 HTTPS</span>
        </template>
        <el-switch
          v-model="form.use_ssl"
          :active-value="true"
          :inactive-value="false"
          :disabled="loading"
        />
      </el-form-item>
    </template>
  </el-form>
</template>

<script setup lang="ts">
import { computed } from 'vue';

export interface UploadForm {
  storage_type: string;
  max_file_size: number;
  path_pattern: string;
  access_key: string;
  secret_key: string;
  region: string;
  bucket: string;
  endpoint: string;
  domain: string;
  use_ssl: boolean;
}

const form = defineModel<UploadForm>('form', { required: true });

defineProps<{
  loading?: boolean;
  isFieldModified: (key: string) => boolean;
}>();

const accessLabel = computed(() => {
  switch (form.value.storage_type) {
    case 'cos':
      return 'SecretId';
    case 'oss':
      return 'AccessKeyId';
    case 'kodo':
      return 'AccessKey';
    case 'r2':
    case 'minio':
      return 'Access Key';
    default:
      return 'Access Key';
  }
});

const secretLabel = computed(() => {
  switch (form.value.storage_type) {
    case 'cos':
      return 'SecretKey';
    case 'oss':
      return 'AccessKeySecret';
    case 'kodo':
      return 'SecretKey';
    case 'r2':
    case 'minio':
      return 'Secret Key';
    default:
      return 'Secret Key';
  }
});

const accessPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 'cos':
      return '例如 AKIDxxxxxxxxxxxxxxxxxxxx';
    case 'oss':
      return '例如 LTAIxxxxxxxxxxxxxxxx';
    default:
      return '';
  }
});

const secretPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 'cos':
      return 'COS 的 SecretKey';
    case 'oss':
      return 'OSS 的 AccessKeySecret';
    default:
      return '';
  }
});

const regionPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 's3':
      return '例如 us-east-1, ap-southeast-1';
    case 'cos':
      return '例如 ap-guangzhou, ap-beijing';
    case 'oss':
      return '例如 oss-cn-hangzhou, oss-cn-beijing';
    case 'kodo':
      return '例如 cn-east-1, cn-north-1, cn-south-1';
    case 'minio':
      return '例如 us-east-1, cn-east-1';
    default:
      return '';
  }
});

const endpointPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 's3':
      return '可选，例如 s3.us-east-1.amazonaws.com';
    case 'r2':
      return '例如 <account-id>.r2.cloudflarestorage.com';
    case 'minio':
      return '例如 localhost:9000 或 minio.example.com';
    default:
      return '';
  }
});

const showRegion = computed(() => {
  const type = form.value.storage_type;
  return type === 's3' || type === 'cos' || type === 'oss' || type === 'kodo' || type === 'minio';
});

const showEndpoint = computed(() => {
  const type = form.value.storage_type;
  return type === 's3' || type === 'r2' || type === 'minio';
});

const showUseSSL = computed(() => {
  const type = form.value.storage_type;
  return type === 'r2' || type === 'minio';
});

const bucketPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 'cos':
      return '例如 my-bucket-1234567890';
    default:
      return '例如 my-bucket';
  }
});

const domainPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 'kodo':
      return '必需，例如 https://cdn.example.com (七牛云CDN域名)';
    case 'cos':
      return '可选，例如 https://cdn.example.com (腾讯云CDN域名)';
    case 'oss':
      return '可选，例如 https://cdn.example.com (阿里云CDN域名)';
    default:
      return '可选，例如 https://cdn.example.com';
  }
});
</script>

<style lang="scss" scoped>
.unit-tip {
  margin-left: 8px;
  color: #909399;
}

@media (max-width: 768px) {
  :deep(.el-form-item__label) {
    width: 100px !important;
    font-size: 13px;
  }
}
</style>
