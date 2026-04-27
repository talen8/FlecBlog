<template>
  <el-form :model="form" label-width="120px" class="setting-form">
    <el-divider content-position="left">邮件通知配置</el-divider>

    <el-form-item label="主机名">
      <el-input v-model="form.email_host" placeholder="例如 smtp.163.com" :disabled="loading" />
    </el-form-item>

    <el-form-item label="端口">
      <el-input
        v-model="form.email_port"
        type="number"
        placeholder="例如 465"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="安全性">
      <el-select
        v-model="form.email_secure"
        placeholder="选择加密方式"
        :disabled="loading"
        style="width: 100%"
      >
        <el-option label="SSL/TLS" value="ssl" />
        <el-option label="STARTTLS" value="starttls" />
        <el-option label="无加密" value="none" />
      </el-select>
    </el-form-item>

    <el-form-item label="用户名">
      <el-input v-model="form.email_username" placeholder="SMTP登录账号" :disabled="loading" />
    </el-form-item>

    <el-form-item label="密码">
      <el-input
        v-model="form.email_password"
        type="password"
        show-password
        placeholder="邮箱授权码或密码"
        :disabled="loading"
        autocomplete="new-password"
      />
    </el-form-item>

    <el-form-item label="发信人">
      <el-input
        v-model="form.email_from"
        placeholder="收件人看到的发件人地址，为空则使用用户名"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">飞书通知配置</el-divider>

    <el-form-item label="App ID">
      <el-input v-model="form.feishu_app_id" placeholder="飞书应用ID" :disabled="loading" />
    </el-form-item>

    <el-form-item label="App Secret">
      <el-input
        v-model="form.feishu_secret"
        type="password"
        show-password
        placeholder="飞书应用Secret"
        :disabled="loading"
        autocomplete="new-password"
      />
    </el-form-item>

    <el-form-item label="群聊 ID">
      <el-input v-model="form.feishu_chat_id" placeholder="接收通知的群聊ID" :disabled="loading" />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
export interface NotificationForm {
  email_host: string;
  email_port: string;
  email_secure: string;
  email_username: string;
  email_from: string;
  email_password: string;
  feishu_app_id: string;
  feishu_secret: string;
  feishu_chat_id: string;
}

const form = defineModel<NotificationForm>('form', { required: true });

defineProps<{
  loading?: boolean;
}>();
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
