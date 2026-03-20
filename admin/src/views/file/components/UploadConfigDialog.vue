<template>
  <el-dialog v-model="visible" :title="'上传配置'" width="600px" :close-on-click-modal="false" @close="handleClose">
    <el-form ref="formRef" :model="form" label-width="100px">
      <el-form-item label="存储类型">
        <el-select v-model="form.storage_type" placeholder="选择存储类型" style="width: 220px">
          <el-option label="本地存储" value="local" />
          <el-option label="亚马逊 S3" value="s3" />
          <el-option label="阿里云 OSS" value="oss" />
          <el-option label="腾讯云 COS" value="cos" />
          <el-option label="七牛云 Kodo" value="kodo" />
          <el-option label="Cloudflare R2" value="r2" />
          <el-option label="MinIO" value="minio" />
        </el-select>
      </el-form-item>
      <el-form-item label="最大文件大小">
        <el-input-number v-model="form.max_file_size" :min="0" :step="1" />
        <span style="margin-left:8px;color:#909399">MB</span>
      </el-form-item>
      <el-form-item label="文件命名">
        <el-input v-model="form.path_pattern" placeholder="{timestamp}_{random}{ext}" />
      </el-form-item>
      <template v-if="form.storage_type !== 'local'">
        <el-form-item :label="accessLabel">
          <el-input v-model="form.access_key" :placeholder="accessPlaceholder" clearable />
        </el-form-item>
        <el-form-item :label="secretLabel">
          <el-input v-model="form.secret_key" type="password" show-password :placeholder="secretPlaceholder"
            clearable />
        </el-form-item>
        <el-form-item v-if="showRegion" label="地域">
          <el-input v-model="form.region" :placeholder="regionPlaceholder" clearable />
        </el-form-item>
        <el-form-item label="存储桶">
          <el-input v-model="form.bucket" :placeholder="bucketPlaceholder" clearable />
        </el-form-item>
        <el-form-item v-if="showEndpoint" label="服务端点">
          <el-input v-model="form.endpoint" :placeholder="endpointPlaceholder" clearable />
        </el-form-item>
        <el-form-item label="自定义域名">
          <el-input v-model="form.domain" :placeholder="domainPlaceholder" clearable />
        </el-form-item>
        <el-form-item v-if="showUseSSL" label="启用 HTTPS">
          <el-switch v-model="form.use_ssl" :active-value="true" :inactive-value="false" />
        </el-form-item>
      </template>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSubmit">保存</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { getSettingGroup, updateSettingGroup } from '@/api/sysconfig'
const props = defineProps<{ modelValue: boolean }>()
const emit = defineEmits<{ 'update:modelValue': [value: boolean], 'success': [] }>()

const visible = computed({
  get: () => props.modelValue,
  set: (val: boolean) => emit('update:modelValue', val)
})

const formRef = ref()
const saving = ref(false)

const form = ref({
  storage_type: 'local',
  max_file_size: 10,
  path_pattern: '{timestamp}_{random}{ext}',
  access_key: '',
  secret_key: '',
  region: '',
  bucket: '',
  endpoint: '',
  domain: '',
  use_ssl: true
})

const loadConfigs = async () => {
  try {
    const data = await getSettingGroup('upload')
    // 适配新的扁平化数据格式
    const configs: Record<string, string> = {}
    Object.entries(data).forEach(([key, value]) => {
      const key2 = key.replace('upload.', '') // 移除前缀
      configs[key2] = value
    })

    form.value.storage_type = configs.storage_type || 'local'
    form.value.max_file_size = Number(configs.max_file_size || 10)
    form.value.path_pattern = configs.path_pattern || '{timestamp}_{random}{ext}'
    form.value.access_key = configs.access_key || ''
    form.value.secret_key = configs.secret_key || ''
    form.value.region = configs.region || ''
    form.value.bucket = configs.bucket || ''
    form.value.endpoint = configs.endpoint || ''
    form.value.domain = configs.domain || ''
    form.value.use_ssl = (configs.use_ssl || 'true') === 'true'
  } catch {
    ElMessage.error('获取配置失败')
  }
}

watch(() => visible.value, (val) => { if (val) loadConfigs() })

const handleSubmit = async () => {
  saving.value = true
  try {
    const payload: Record<string, string> = {
      'upload.storage_type': form.value.storage_type,
      'upload.max_file_size': String(form.value.max_file_size),
      'upload.path_pattern': form.value.path_pattern,
      'upload.access_key': form.value.access_key,
      'upload.secret_key': form.value.secret_key,
      'upload.region': form.value.region,
      'upload.bucket': form.value.bucket,
      'upload.endpoint': form.value.endpoint,
      'upload.domain': form.value.domain,
      'upload.use_ssl': form.value.use_ssl ? 'true' : 'false'
    }
    await updateSettingGroup('upload', payload)
    ElMessage.success('保存成功（配置已自动热重载）')
    emit('success')
    visible.value = false
  } catch (e) {
    if (e instanceof Error) ElMessage.error(e.message)
    else ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const handleClose = () => { }

const accessLabel = computed(() => {
  switch (form.value.storage_type) {
    case 'cos': return 'SecretId'
    case 'oss': return 'AccessKeyId'
    case 'kodo': return 'AccessKey'
    case 'r2': return 'Access Key'
    case 'minio': return 'Access Key'
    default: return 'Access Key'
  }
})

const secretLabel = computed(() => {
  switch (form.value.storage_type) {
    case 'cos': return 'SecretKey'
    case 'oss': return 'AccessKeySecret'
    case 'kodo': return 'SecretKey'
    case 'r2': return 'Secret Key'
    case 'minio': return 'Secret Key'
    default: return 'Secret Key'
  }
})

const accessPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 'cos': return '例如 AKIDxxxxxxxxxxxxxxxxxxxx'
    case 'oss': return '例如 LTAIxxxxxxxxxxxxxxxx'
    default: return ''
  }
})

const secretPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 'cos': return 'COS 的 SecretKey'
    case 'oss': return 'OSS 的 AccessKeySecret'
    default: return ''
  }
})

const regionPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 's3': return '例如 us-east-1, ap-southeast-1'
    case 'cos': return '例如 ap-guangzhou, ap-beijing'
    case 'oss': return '例如 oss-cn-hangzhou, oss-cn-beijing'
    case 'kodo': return '例如 cn-east-1, cn-north-1, cn-south-1'
    default: return ''
  }
})

const endpointPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 's3': return '可选，例如 s3.us-east-1.amazonaws.com'
    case 'r2': return '例如 <account-id>.r2.cloudflarestorage.com'
    case 'minio': return '例如 localhost:9000 或 minio.example.com'
    default: return ''
  }
})

const showRegion = computed(() => {
  const t = form.value.storage_type
  return t === 's3' || t === 'cos' || t === 'oss' || t === 'kodo'
})

const showEndpoint = computed(() => {
  const t = form.value.storage_type
  return t === 's3' || t === 'r2' || t === 'minio'
})

const showUseSSL = computed(() => {
  const t = form.value.storage_type
  return t === 'r2' || t === 'minio'
})

const bucketPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 's3': return '例如 my-bucket'
    case 'cos': return '例如 my-bucket-1234567890'
    case 'oss': return '例如 my-bucket'
    case 'kodo': return '例如 my-bucket'
    case 'r2': return '例如 my-bucket'
    case 'minio': return '例如 my-bucket'
    default: return '例如 my-bucket'
  }
})

const domainPlaceholder = computed(() => {
  switch (form.value.storage_type) {
    case 'kodo': return '必需，例如 https://cdn.example.com (七牛云CDN域名)'
    case 'cos': return '可选，例如 https://cdn.example.com (腾讯云CDN域名)'
    case 'oss': return '可选，例如 https://cdn.example.com (阿里云CDN域名)'
    default: return '可选，例如 https://cdn.example.com'
  }
})
</script>

<style scoped>
.dialog-footer {
  display: inline-flex;
  gap: 12px;
}
</style>