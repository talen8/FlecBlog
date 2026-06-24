<template>
  <div class="system-settings">
    <el-card shadow="never">
      <!-- 工具栏 -->
      <div class="toolbar">
        <h2>系统设置</h2>
        <div class="actions">
          <el-button
            type="primary"
            :loading="saving"
            :disabled="!canEditSettings"
            @click="handleSave"
          >
            保存配置
          </el-button>
          <el-button @click="loadAllConfigs">重置</el-button>
        </div>
      </div>

      <!-- 标签页 -->
      <el-tabs v-model="activeTab" class="setting-tabs">
        <!-- 基本配置标签页 -->
        <el-tab-pane label="基本配置" name="basic">
          <BasicSettingsTab
            ref="basicTabRef"
            v-model:form="basicForm"
            :is-field-modified="basicIsFieldModified"
            :loading="loading || !canEditSettings"
          />
        </el-tab-pane>

        <!-- 通知配置标签页 -->
        <el-tab-pane label="通知配置" name="notification">
          <NotificationSettingsTab
            v-model:form="notificationForm"
            :is-field-modified="notificationIsFieldModified"
            :loading="loading || !canEditSettings"
          />
        </el-tab-pane>

        <!-- 上传配置标签页 -->
        <el-tab-pane label="上传配置" name="upload">
          <UploadSettingsTab
            v-model:form="uploadForm"
            :is-field-modified="uploadIsFieldModified"
            :loading="loading || !canEditSettings"
          />
        </el-tab-pane>

        <!-- AI 配置标签页 -->
        <el-tab-pane label="AI 配置" name="ai">
          <AISettingsTab
            v-model:form="aiForm"
            :is-field-modified="aiIsFieldModified"
            :loading="loading || !canEditSettings"
          />
        </el-tab-pane>

        <!-- OAuth 配置标签页 -->
        <el-tab-pane label="OAuth 配置" name="oauth">
          <OAuthSettingsTab
            v-model:form="oauthForm"
            :is-field-modified="oauthIsFieldModified"
            :loading="loading || !canEditSettings"
          />
        </el-tab-pane>

        <!-- 导入导出标签页 -->
        <el-tab-pane label="导入导出" name="import-export">
          <ImportExportTab :readonly="!canEditSettings" @import-success="handleImportSuccess" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { ElMessage } from 'element-plus';
import { getSettingGroup, updateSettingGroup } from '@/api/sysconfig';
import { isSuperAdmin } from '@/utils/auth';
import BasicSettingsTab from './components/BasicSettingsTab.vue';
import NotificationSettingsTab from './components/NotificationSettingsTab.vue';
import UploadSettingsTab from './components/UploadSettingsTab.vue';
import AISettingsTab from './components/AISettingsTab.vue';
import OAuthSettingsTab from './components/OAuthSettingsTab.vue';
import ImportExportTab from './components/ImportExportTab.vue';
import type { SettingGroupType } from '@/types/sysconfig';
import type { NotificationForm } from './components/NotificationSettingsTab.vue';
import type { UploadForm } from './components/UploadSettingsTab.vue';

// 页面状态
const activeTab = ref('basic');
const route = useRoute();
const loading = ref(false);
const saving = ref(false);
const canEditSettings = computed(() => isSuperAdmin());

// 标签页引用
const basicTabRef = ref<InstanceType<typeof BasicSettingsTab>>();

// 原始配置快照（用于比较修改）
const originalBasicForm = ref<Record<string, string>>({});
const originalNotificationForm = ref<Record<string, string>>({});
const originalUploadForm = ref<Record<string, string | number | boolean>>({});
const originalAiForm = ref<Record<string, string>>({});
const originalOAuthForm = ref<Record<string, string>>({});

// 基本配置表单
const basicForm = ref({
  author: '',
  author_avatar: '',
  icp: '',
  police_record: '',
  admin_url: '',
  blog_url: '',
  title: '',
  subtitle: '',
  description: '',
  keywords: '',
  favicon: '',
  custom_head: '',
  custom_body: '',
  emojis: '',
  meting_api: '',
  cravatar_url: '',
  ip_api_url: '',
  cover_maker_api: '',
});

// 通知配置表单
const notificationForm = ref<NotificationForm>({
  email_host: '',
  email_port: '465',
  email_secure: 'ssl',
  email_username: '',
  email_from: '',
  email_password: '',
  feishu_app_id: '',
  feishu_secret: '',
  feishu_chat_id: '',
});

// 上传配置表单
const uploadForm = ref<UploadForm>({
  storage_type: 'local',
  max_file_size: 10,
  path_pattern: '{timestamp}_{random}{ext}',
  access_key: '',
  secret_key: '',
  region: '',
  bucket: '',
  endpoint: '',
  domain: '',
  use_ssl: true,
});

// AI 配置表单
const aiForm = ref({
  base_url: '',
  api_key: '',
  model: '',
  summary_prompt: '',
  ai_summary_prompt: '',
  title_prompt: '',
  mcp_secret: '',
});

// OAuth 配置表单
const oauthForm = ref({
  'github.enabled': 'false',
  'github.client_id': '',
  'github.client_secret': '',
  'github.redirect_url': '',
  'google.enabled': 'false',
  'google.client_id': '',
  'google.client_secret': '',
  'google.redirect_url': '',
  'qq.enabled': 'false',
  'qq.client_id': '',
  'qq.client_secret': '',
  'qq.redirect_url': '',
  'microsoft.enabled': 'false',
  'microsoft.client_id': '',
  'microsoft.client_secret': '',
  'microsoft.redirect_url': '',
  'oidc.enabled': 'false',
  'oidc.issuer_url': '',
  'oidc.client_id': '',
  'oidc.client_secret': '',
  'oidc.redirect_url': '',
  'wechat.enabled': 'false',
  'wechat.appid': '',
  'wechat.secret': '',
  worker_proxy: '',
});

// 通用配置加载函数
const loadConfigs = async (group: SettingGroupType) => {
  const data = await getSettingGroup(group);
  const configs: Record<string, string> = {};

  // 适配新的扁平化数据格式
  Object.entries(data).forEach(([key, value]) => {
    // 将键名中的分组前缀去掉，例如将 'basic.author' 转换为 'author'
    const shortKey = key.replace(`${group}.`, '');
    configs[shortKey] = value;
  });

  return configs;
};

// 加载基本配置
const loadBasicConfigs = async () => {
  try {
    const configs = await loadConfigs('basic');
    const data = {
      author: configs.author || '',
      author_avatar: configs.author_avatar || '',
      icp: configs.icp || '',
      police_record: configs.police_record || '',
      admin_url: configs.admin_url || '',
      blog_url: configs.blog_url || '',
      title: configs.title || '',
      subtitle: configs.subtitle || '',
      description: configs.description || '',
      keywords: configs.keywords || '',
      favicon: configs.favicon || '',
      custom_head: configs.custom_head || '',
      custom_body: configs.custom_body || '',
      emojis: configs.emojis || '',
      meting_api: configs.meting_api || '',
      cravatar_url: configs.cravatar_url || '',
      ip_api_url: configs.ip_api_url || '',
      cover_maker_api: configs.cover_maker_api || '',
    };
    Object.assign(basicForm.value, data);
    originalBasicForm.value = { ...data };
  } catch {
    ElMessage.error('获取基本配置失败');
  }
};

// 加载通知配置
const loadNotificationConfigs = async () => {
  try {
    const configs = await loadConfigs('notification');
    const data = {
      email_host: configs.email_host || '',
      email_port: configs.email_port || '465',
      email_secure: configs.email_secure || 'ssl',
      email_username: configs.email_username || '',
      email_from: configs.email_from || '',
      email_password: configs.email_password || '',
      feishu_app_id: configs.feishu_app_id || '',
      feishu_secret: configs.feishu_secret || '',
      feishu_chat_id: configs.feishu_chat_id || '',
    };
    Object.assign(notificationForm.value, data);
    originalNotificationForm.value = { ...data };
  } catch {
    ElMessage.error('获取通知配置失败');
  }
};

// 加载上传配置
const loadUploadConfigs = async () => {
  try {
    const configs = await loadConfigs('upload');
    const data = {
      storage_type: configs.storage_type || 'local',
      max_file_size: Number(configs.max_file_size || 10),
      path_pattern: configs.path_pattern || '{timestamp}_{random}{ext}',
      access_key: configs.access_key || '',
      secret_key: configs.secret_key || '',
      region: configs.region || '',
      bucket: configs.bucket || '',
      endpoint: configs.endpoint || '',
      domain: configs.domain || '',
      use_ssl: (configs.use_ssl || 'true') === 'true',
    };
    Object.assign(uploadForm.value, data);
    originalUploadForm.value = { ...data };
  } catch {
    ElMessage.error('获取上传配置失败');
  }
};

// 加载 AI 配置
const loadAIConfigs = async () => {
  try {
    const configs = await loadConfigs('ai');
    const data = {
      base_url: configs.base_url || '',
      api_key: configs.api_key || '',
      model: configs.model || '',
      summary_prompt: configs.summary_prompt || '',
      ai_summary_prompt: configs.ai_summary_prompt || '',
      title_prompt: configs.title_prompt || '',
      mcp_secret: configs.mcp_secret || '',
    };
    Object.assign(aiForm.value, data);
    originalAiForm.value = { ...data };
  } catch {
    ElMessage.error('获取 AI 配置失败');
  }
};

// 加载 OAuth 配置
const loadOAuthConfigs = async () => {
  try {
    const configs = await loadConfigs('oauth');
    const data = {
      'github.enabled': configs['github.enabled'] || 'false',
      'github.client_id': configs['github.client_id'] || '',
      'github.client_secret': configs['github.client_secret'] || '',
      'github.redirect_url': configs['github.redirect_url'] || '',
      'google.enabled': configs['google.enabled'] || 'false',
      'google.client_id': configs['google.client_id'] || '',
      'google.client_secret': configs['google.client_secret'] || '',
      'google.redirect_url': configs['google.redirect_url'] || '',
      'qq.enabled': configs['qq.enabled'] || 'false',
      'qq.client_id': configs['qq.client_id'] || '',
      'qq.client_secret': configs['qq.client_secret'] || '',
      'qq.redirect_url': configs['qq.redirect_url'] || '',
      'microsoft.enabled': configs['microsoft.enabled'] || 'false',
      'microsoft.client_id': configs['microsoft.client_id'] || '',
      'microsoft.client_secret': configs['microsoft.client_secret'] || '',
      'microsoft.redirect_url': configs['microsoft.redirect_url'] || '',
      'oidc.enabled': configs['oidc.enabled'] || 'false',
      'oidc.issuer_url': configs['oidc.issuer_url'] || '',
      'oidc.client_id': configs['oidc.client_id'] || '',
      'oidc.client_secret': configs['oidc.client_secret'] || '',
      'oidc.redirect_url': configs['oidc.redirect_url'] || '',
      'wechat.enabled': configs['wechat.enabled'] || 'false',
      'wechat.appid': configs['wechat.appid'] || '',
      'wechat.secret': configs['wechat.secret'] || '',
      worker_proxy: configs['worker_proxy'] || '',
    };
    Object.assign(oauthForm.value, data);
    originalOAuthForm.value = { ...data };
  } catch {
    ElMessage.error('获取 OAuth 配置失败');
  }
};

// 加载所有配置
const loadAllConfigs = async () => {
  loading.value = true;
  try {
    await Promise.all([
      loadBasicConfigs(),
      loadNotificationConfigs(),
      loadUploadConfigs(),
      loadAIConfigs(),
      loadOAuthConfigs(),
    ]);
  } finally {
    loading.value = false;
  }
};

const makeIsFieldModified = (form: Record<string, unknown>, original: Record<string, unknown>) => {
  return (key: string): boolean => {
    if (loading.value) return false;
    const current = form[key];
    const orig = original[key];
    if (current === orig) return false;
    if (current === undefined && orig === undefined) return false;
    if (current === null && orig === null) return false;
    return JSON.stringify(current) !== JSON.stringify(orig);
  };
};

const basicIsFieldModified = computed(() =>
  makeIsFieldModified(basicForm.value, originalBasicForm.value)
);
const notificationIsFieldModified = computed(() =>
  makeIsFieldModified(notificationForm.value, originalNotificationForm.value)
);
const uploadIsFieldModified = computed(() =>
  makeIsFieldModified(uploadForm.value, originalUploadForm.value)
);
const aiIsFieldModified = computed(() => makeIsFieldModified(aiForm.value, originalAiForm.value));
const oauthIsFieldModified = computed(() =>
  makeIsFieldModified(oauthForm.value, originalOAuthForm.value)
);

// 统一保存配置
const handleSave = async () => {
  if (!canEditSettings.value) {
    ElMessage.warning('仅超级管理员可修改系统配置');
    return;
  }

  saving.value = true;
  try {
    const uploadPromises: Promise<void>[] = [];

    // 收集所有待上传的图片（并行上传）
    const basicUploaders = basicTabRef.value;
    if (basicUploaders) {
      if (basicUploaders.authorAvatarUploaderRef?.getPendingCount()) {
        uploadPromises.push(
          basicUploaders.authorAvatarUploaderRef.uploadPendingFile().then(url => {
            if (url) basicForm.value.author_avatar = url;
          })
        );
      }
      if (basicUploaders.faviconUploaderRef?.getPendingCount()) {
        uploadPromises.push(
          basicUploaders.faviconUploaderRef.uploadPendingFile().then(url => {
            if (url) basicForm.value.favicon = url;
          })
        );
      }
    }

    // 等待所有上传完成（使用 allSettled 确保即使部分失败也继续）
    if (uploadPromises.length > 0) {
      const results = await Promise.allSettled(uploadPromises);
      const failedUploads = results.filter(r => r.status === 'rejected');
      if (failedUploads.length > 0) {
        saving.value = false;
        ElMessage.error(`${failedUploads.length} 个文件上传失败，请重试`);
        return;
      }
    }

    // 基本配置
    const basicPayload: Record<string, string> = {
      author: basicForm.value.author,
      author_avatar: basicForm.value.author_avatar,
      icp: basicForm.value.icp,
      police_record: basicForm.value.police_record,
      admin_url: basicForm.value.admin_url,
      blog_url: basicForm.value.blog_url,
      title: basicForm.value.title,
      subtitle: basicForm.value.subtitle,
      description: basicForm.value.description,
      keywords: basicForm.value.keywords,
      favicon: basicForm.value.favicon,
      custom_head: basicForm.value.custom_head,
      custom_body: basicForm.value.custom_body,
      emojis: basicForm.value.emojis,
      meting_api: basicForm.value.meting_api,
      cravatar_url: basicForm.value.cravatar_url,
      ip_api_url: basicForm.value.ip_api_url,
      cover_maker_api: basicForm.value.cover_maker_api,
    };

    // 通知配置
    const notificationPayload: Record<string, string> = {
      email_host: notificationForm.value.email_host,
      email_port: String(notificationForm.value.email_port),
      email_secure: notificationForm.value.email_secure,
      email_username: notificationForm.value.email_username,
      email_from: notificationForm.value.email_from,
      email_password: notificationForm.value.email_password,
      feishu_app_id: notificationForm.value.feishu_app_id,
      feishu_secret: notificationForm.value.feishu_secret,
      feishu_chat_id: notificationForm.value.feishu_chat_id,
    };

    // 上传配置
    const uploadPayload: Record<string, string> = {
      storage_type: uploadForm.value.storage_type,
      max_file_size: String(uploadForm.value.max_file_size),
      path_pattern: uploadForm.value.path_pattern,
      access_key: uploadForm.value.access_key,
      secret_key: uploadForm.value.secret_key,
      region: uploadForm.value.region,
      bucket: uploadForm.value.bucket,
      endpoint: uploadForm.value.endpoint,
      domain: uploadForm.value.domain,
      use_ssl: uploadForm.value.use_ssl ? 'true' : 'false',
    };

    // AI 配置
    const aiPayload: Record<string, string> = {
      base_url: aiForm.value.base_url,
      api_key: aiForm.value.api_key,
      model: aiForm.value.model,
      summary_prompt: aiForm.value.summary_prompt,
      ai_summary_prompt: aiForm.value.ai_summary_prompt,
      title_prompt: aiForm.value.title_prompt,
    };

    // OAuth 配置
    const oauthPayload: Record<string, string> = {
      'github.enabled': oauthForm.value['github.enabled'],
      'github.client_id': oauthForm.value['github.client_id'],
      'github.client_secret': oauthForm.value['github.client_secret'],
      'github.redirect_url': oauthForm.value['github.redirect_url'],
      'google.enabled': oauthForm.value['google.enabled'],
      'google.client_id': oauthForm.value['google.client_id'],
      'google.client_secret': oauthForm.value['google.client_secret'],
      'google.redirect_url': oauthForm.value['google.redirect_url'],
      'qq.enabled': oauthForm.value['qq.enabled'],
      'qq.client_id': oauthForm.value['qq.client_id'],
      'qq.client_secret': oauthForm.value['qq.client_secret'],
      'qq.redirect_url': oauthForm.value['qq.redirect_url'],
      'microsoft.enabled': oauthForm.value['microsoft.enabled'],
      'microsoft.client_id': oauthForm.value['microsoft.client_id'],
      'microsoft.client_secret': oauthForm.value['microsoft.client_secret'],
      'microsoft.redirect_url': oauthForm.value['microsoft.redirect_url'],
      'oidc.enabled': oauthForm.value['oidc.enabled'],
      'oidc.issuer_url': oauthForm.value['oidc.issuer_url'],
      'oidc.client_id': oauthForm.value['oidc.client_id'],
      'oidc.client_secret': oauthForm.value['oidc.client_secret'],
      'oidc.redirect_url': oauthForm.value['oidc.redirect_url'],
      'wechat.enabled': oauthForm.value['wechat.enabled'],
      'wechat.appid': oauthForm.value['wechat.appid'],
      'wechat.secret': oauthForm.value['wechat.secret'],
      worker_proxy: oauthForm.value['worker_proxy'],
    };

    // 构建需要保存的配置组列表
    const savePromises = [
      updateSettingGroup('basic', basicPayload),
      updateSettingGroup('notification', notificationPayload),
      updateSettingGroup('upload', uploadPayload),
      updateSettingGroup('ai', aiPayload),
      updateSettingGroup('oauth', oauthPayload),
    ];

    // 并行保存所有配置组
    await Promise.all(savePromises);

    // 更新原始值快照
    originalBasicForm.value = { ...basicForm.value };
    originalNotificationForm.value = { ...notificationForm.value };
    originalUploadForm.value = { ...uploadForm.value };
    originalAiForm.value = { ...aiForm.value };
    originalOAuthForm.value = { ...oauthForm.value };

    ElMessage.success('配置保存成功');
  } catch (e) {
    if (e instanceof Error) ElMessage.error(e.message);
    else ElMessage.error('保存失败');
  } finally {
    saving.value = false;
  }
};

const validTabs = new Set<SettingGroupType | 'import-export'>([
  'basic',
  'notification',
  'upload',
  'ai',
  'oauth',
  'import-export',
]);

watch(
  () => route.query.tab,
  tab => {
    if (typeof tab === 'string' && validTabs.has(tab as SettingGroupType | 'import-export')) {
      activeTab.value = tab;
    }
  },
  { immediate: true }
);

// 导入成功回调
const handleImportSuccess = () => {
  // 可以在这里添加导入成功后的逻辑
};

onMounted(() => {
  loadAllConfigs();
});
</script>

<style lang="scss" scoped>
.system-settings {
  height: 100%;

  :deep(.el-card) {
    height: 100%;
    display: flex;
    flex-direction: column;

    .el-card__body {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
    }
  }
}

.toolbar {
  margin-bottom: 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;

  h2 {
    margin: 0;
    font-size: 20px;
    font-weight: 500;
  }

  .actions {
    display: flex;
    gap: 12px;

    .el-button {
      margin: 0;
    }
  }
}

.setting-tabs {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;

  :deep(.el-tabs__header) {
    margin: 0 0 12px 0;
    flex-shrink: 0;
  }

  :deep(.el-tabs__nav-wrap) {
    justify-content: center;

    &::after {
      display: none;
    }
  }

  :deep(.el-tabs__nav) {
    float: none;
  }

  :deep(.el-tabs__content) {
    flex: 1;
    overflow: hidden;
  }

  :deep(.el-tab-pane) {
    height: 100%;
    overflow-y: auto;
    padding: 0 16px;

    .setting-form {
      max-width: 95%;
      margin: 0 auto;
    }
  }

  :deep(.field-modified) {
    color: #e6a23c;
    font-weight: 600;
  }
}

// 移动端适配
@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;

    h2 {
      font-size: 18px;
    }

    .actions {
      width: 100%;

      .el-button {
        flex: 1;
      }
    }
  }

  .setting-tabs {
    :deep(.el-tabs__nav-wrap) {
      justify-content: flex-start;
    }

    :deep(.el-tabs__nav-scroll) {
      overflow-x: auto;
      -webkit-overflow-scrolling: touch;
      scrollbar-width: none;

      &::-webkit-scrollbar {
        display: none;
      }
    }

    :deep(.el-tabs__nav-wrap.is-scrollable) {
      padding: 0;
    }

    :deep(.el-tab-pane) {
      padding: 0 8px;
      overflow-x: auto;
      -webkit-overflow-scrolling: touch;
      scrollbar-width: none;

      &::-webkit-scrollbar {
        display: none;
      }

      .setting-form {
        max-width: none;
        min-width: 800px;
      }
    }
  }

  :deep(.el-form-item__label) {
    width: 120px !important;
    flex-shrink: 0;
  }
}
</style>
