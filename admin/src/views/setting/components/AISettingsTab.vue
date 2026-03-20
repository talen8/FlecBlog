<template>
  <el-form :model="form" label-width="120px" class="setting-form">
    <el-form-item label="API 端点">
      <el-input v-model="form.base_url" placeholder="例如 https://api.deepseek.com" :disabled="loading" />
    </el-form-item>

    <el-form-item label="API 密钥">
      <el-input v-model="form.api_key" type="password" show-password placeholder="输入 API Key" :disabled="loading" autocomplete="off" />
    </el-form-item>

    <el-form-item label="模型名称">
      <el-input v-model="form.model" placeholder="例如 deepseek-chat" :disabled="loading" />
    </el-form-item>

    <el-form-item label=" ">
      <el-button :loading="testing" @click="handleTest">测试连接</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { testAIConfig } from '@/api/ai'

interface AIForm {
  base_url: string
  api_key: string
  model: string
}

const form = defineModel<AIForm>('form', { required: true })

defineProps<{
  loading?: boolean
}>()

const testing = ref(false)

async function handleTest() {
  if (!form.value.base_url || !form.value.api_key || !form.value.model) {
    ElMessage.warning('请先填写完整的 API 端点、密钥和模型名称')
    return
  }
  testing.value = true
  try {
    await testAIConfig({
      base_url: form.value.base_url,
      api_key: form.value.api_key,
      model: form.value.model,
    })
    ElMessage.success('连接成功，配置可用')
  } catch (e: any) {
    ElMessage.error(e?.message || '连接失败，请检查配置')
  } finally {
    testing.value = false
  }
}
</script>

<style lang="scss" scoped>
// 移动端适配
@media (max-width: 768px) {
  :deep(.el-form-item__label) {
    width: 100px !important;
    font-size: 13px;
  }
}
</style>
