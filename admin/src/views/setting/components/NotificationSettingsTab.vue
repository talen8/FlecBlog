<template>
  <el-form :model="form" label-width="120px" class="setting-form">
    <el-divider content-position="left">邮件通知</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('email_host') }">主机名</span>
      </template>
      <el-input v-model="form.email_host" placeholder="例如 smtp.163.com" :disabled="loading" />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('email_port') }">端口</span>
      </template>
      <el-input
        v-model="form.email_port"
        type="number"
        placeholder="例如 465"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('email_secure') }">安全性</span>
      </template>
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

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('email_username') }">用户名</span>
      </template>
      <el-input v-model="form.email_username" placeholder="SMTP登录账号" :disabled="loading" />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('email_password') }">密码</span>
      </template>
      <el-input
        v-model="form.email_password"
        type="password"
        show-password
        placeholder="邮箱授权码或密码"
        :disabled="loading"
        autocomplete="new-password"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('email_from') }">发信人</span>
      </template>
      <el-input
        v-model="form.email_from"
        placeholder="收件人看到的发件人地址，为空则使用用户名"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">飞书通知</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('feishu_webhook_url') }"
          >Webhook 地址</span
        >
      </template>
      <el-input
        v-model="form.feishu_webhook_url"
        placeholder="https://open.feishu.cn/open-apis/bot/v2/hook/xxxxx"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('feishu_secret') }">签名密钥</span>
      </template>
      <el-input
        v-model="form.feishu_secret"
        type="password"
        show-password
        placeholder="留空表示机器人未开启签名校验"
        :disabled="loading"
        autocomplete="new-password"
      />
    </el-form-item>

    <el-divider content-position="left">Telegram</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('telegram_bot_token') }">Bot Token</span>
      </template>
      <el-input
        v-model="form.telegram_bot_token"
        type="password"
        show-password
        placeholder="123456:ABC-DEF..."
        :disabled="loading"
        autocomplete="new-password"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('telegram_chat_id') }">Chat ID</span>
      </template>
      <el-input v-model="form.telegram_chat_id" placeholder="-1001234567890" :disabled="loading" />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import type { NotificationForm } from '@/types/setting';

const form = defineModel<NotificationForm>('form', { required: true });

defineProps<{
  loading?: boolean;
  isFieldModified: (key: string) => boolean;
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
