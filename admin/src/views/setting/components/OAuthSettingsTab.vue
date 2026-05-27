<template>
  <el-form :model="form" label-width="140px" class="setting-form">
    <el-divider content-position="left">GitHub 登录</el-divider>
    <el-form-item label="启用 GitHub">
      <el-switch
        v-model="form['github.enabled']"
        active-value="true"
        inactive-value="false"
        :disabled="loading"
      />
    </el-form-item>
    <el-form-item label="Client ID">
      <el-input
        v-model="form['github.client_id']"
        placeholder="GitHub Client ID"
        :disabled="loading"
      />
    </el-form-item>
    <el-form-item label="Client Secret">
      <el-input
        v-model="form['github.client_secret']"
        type="password"
        show-password
        placeholder="GitHub Client Secret"
        :disabled="loading"
        autocomplete="off"
      />
    </el-form-item>
    <el-form-item label="回调地址">
      <el-input
        v-model="form['github.redirect_url']"
        placeholder="例如 https://api.your-domain.com/v1/auth/github/callback"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">Google 登录</el-divider>
    <el-form-item label="启用 Google">
      <el-switch
        v-model="form['google.enabled']"
        active-value="true"
        inactive-value="false"
        :disabled="loading"
      />
    </el-form-item>
    <el-form-item label="Client ID">
      <el-input
        v-model="form['google.client_id']"
        placeholder="Google Client ID"
        :disabled="loading"
      />
    </el-form-item>
    <el-form-item label="Client Secret">
      <el-input
        v-model="form['google.client_secret']"
        type="password"
        show-password
        placeholder="Google Client Secret"
        :disabled="loading"
        autocomplete="off"
      />
    </el-form-item>
    <el-form-item label="回调地址">
      <el-input
        v-model="form['google.redirect_url']"
        placeholder="例如 https://api.your-domain.com/v1/auth/google/callback"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">QQ 登录</el-divider>
    <el-form-item label="启用 QQ">
      <el-switch
        v-model="form['qq.enabled']"
        active-value="true"
        inactive-value="false"
        :disabled="loading"
      />
    </el-form-item>
    <el-form-item label="AppID">
      <el-input v-model="form['qq.client_id']" placeholder="QQ 互联 AppID" :disabled="loading" />
    </el-form-item>
    <el-form-item label="AppKey">
      <el-input
        v-model="form['qq.client_secret']"
        type="password"
        show-password
        placeholder="QQ 互联 AppKey"
        :disabled="loading"
        autocomplete="off"
      />
    </el-form-item>
    <el-form-item label="回调地址">
      <el-input
        v-model="form['qq.redirect_url']"
        placeholder="例如 https://api.your-domain.com/v1/auth/qq/callback"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">Microsoft 登录</el-divider>
    <el-form-item label="启用 Microsoft">
      <el-switch
        v-model="form['microsoft.enabled']"
        active-value="true"
        inactive-value="false"
        :disabled="loading"
      />
    </el-form-item>
    <el-form-item label="Client ID">
      <el-input
        v-model="form['microsoft.client_id']"
        placeholder="Microsoft Application (client) ID"
        :disabled="loading"
      />
    </el-form-item>
    <el-form-item label="Client Secret">
      <el-input
        v-model="form['microsoft.client_secret']"
        type="password"
        show-password
        placeholder="Microsoft Client Secret"
        :disabled="loading"
        autocomplete="off"
      />
    </el-form-item>
    <el-form-item label="回调地址">
      <el-input
        v-model="form['microsoft.redirect_url']"
        placeholder="例如 https://api.your-domain.com/v1/auth/microsoft/callback"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">微信小程序登录</el-divider>
    <el-form-item label="启用微信">
      <el-switch
        v-model="form['wechat.enabled']"
        active-value="true"
        inactive-value="false"
        :disabled="loading"
      />
    </el-form-item>
    <el-form-item label="AppID">
      <el-input v-model="form['wechat.appid']" placeholder="微信小程序 AppID" :disabled="loading" />
    </el-form-item>
    <el-form-item label="AppSecret">
      <el-input
        v-model="form['wechat.secret']"
        type="password"
        show-password
        placeholder="微信小程序 AppSecret"
        :disabled="loading"
        autocomplete="off"
      />
    </el-form-item>

    <el-divider content-position="left">OIDC 登录</el-divider>
    <el-form-item label="启用 OIDC">
      <el-switch
        v-model="form['oidc.enabled']"
        active-value="true"
        inactive-value="false"
        :disabled="loading"
      />
    </el-form-item>
    <el-form-item label="Issuer URL">
      <el-input
        v-model="form['oidc.issuer_url']"
        placeholder="https://your-casdoor.example.com/.well-known/openid-configuration"
        :disabled="loading"
      />
    </el-form-item>
    <el-form-item label="Client ID">
      <el-input v-model="form['oidc.client_id']" placeholder="OIDC Client ID" :disabled="loading" />
    </el-form-item>
    <el-form-item label="Client Secret">
      <el-input
        v-model="form['oidc.client_secret']"
        type="password"
        show-password
        placeholder="OIDC Client Secret"
        :disabled="loading"
        autocomplete="off"
      />
    </el-form-item>
    <el-form-item label="回调地址">
      <el-input
        v-model="form['oidc.redirect_url']"
        placeholder="例如 https://api.your-domain.com/v1/auth/oidc/callback"
        :disabled="loading"
      />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
interface OAuthForm {
  'github.enabled': string;
  'github.client_id': string;
  'github.client_secret': string;
  'github.redirect_url': string;
  'google.enabled': string;
  'google.client_id': string;
  'google.client_secret': string;
  'google.redirect_url': string;
  'qq.enabled': string;
  'qq.client_id': string;
  'qq.client_secret': string;
  'qq.redirect_url': string;
  'microsoft.enabled': string;
  'microsoft.client_id': string;
  'microsoft.client_secret': string;
  'microsoft.redirect_url': string;
  'wechat.enabled': string;
  'wechat.appid': string;
  'wechat.secret': string;
  'oidc.enabled': string;
  'oidc.issuer_url': string;
  'oidc.client_id': string;
  'oidc.client_secret': string;
  'oidc.redirect_url': string;
}

const form = defineModel<OAuthForm>('form', { required: true });

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
