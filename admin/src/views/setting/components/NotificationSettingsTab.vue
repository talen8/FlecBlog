<template>
  <el-form :model="form" label-width="120px" class="setting-form">
    <el-divider content-position="left">邮件通知配置</el-divider>

    <el-form-item label="邮箱服务器">
      <el-input v-model="form.email_host" placeholder="例如 smtp.163.com" :disabled="loading" />
    </el-form-item>

    <el-form-item label="邮箱端口">
      <el-input v-model="form.email_port" type="number" placeholder="例如 465" :disabled="loading" />
    </el-form-item>

    <el-form-item label="邮箱用户名">
      <el-input v-model="form.email_username" placeholder="邮箱账号，例如 user@domain.com" :disabled="loading" />
    </el-form-item>

    <el-form-item label="邮箱密码">
      <el-input v-model="form.email_password" type="password" show-password placeholder="邮箱授权码或密码"
        :disabled="loading" autocomplete="new-password" />
    </el-form-item>

    <el-divider content-position="left">飞书通知配置</el-divider>

    <el-form-item label="App ID">
      <el-input v-model="form.feishu_app_id" placeholder="飞书应用ID" :disabled="loading" />
    </el-form-item>

    <el-form-item label="App Secret">
      <el-input v-model="form.feishu_secret" type="password" show-password placeholder="飞书应用Secret"
        :disabled="loading" autocomplete="new-password" />
    </el-form-item>

    <el-form-item label="群聊 ID">
      <el-input v-model="form.feishu_chat_id" placeholder="接收通知的群聊ID" :disabled="loading" />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
export interface NotificationForm {
  email_host: string
  email_port: string
  email_username: string
  email_password: string
  feishu_app_id: string
  feishu_secret: string
  feishu_chat_id: string
}

const form = defineModel<NotificationForm>('form', { required: true })

defineProps<{
  loading?: boolean
}>()
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
