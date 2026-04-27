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
            :loading="loading || !canEditSettings"
          />
        </el-tab-pane>

        <!-- 博客配置标签页 -->
        <el-tab-pane label="博客配置" name="blog">
          <BlogSettingsTab
            ref="blogTabRef"
            v-model:form="blogForm"
            :loading="loading || !canEditSettings"
          />
        </el-tab-pane>

        <!-- 通知配置标签页 -->
        <el-tab-pane label="通知配置" name="notification">
          <NotificationSettingsTab
            v-model:form="notificationForm"
            :loading="loading || !canEditSettings"
          />
        </el-tab-pane>

        <!-- 上传配置标签页 -->
        <el-tab-pane label="上传配置" name="upload">
          <UploadSettingsTab v-model:form="uploadForm" :loading="loading || !canEditSettings" />
        </el-tab-pane>

        <!-- AI 配置标签页 -->
        <el-tab-pane label="AI 配置" name="ai">
          <AISettingsTab v-model:form="aiForm" :loading="loading || !canEditSettings" />
        </el-tab-pane>

        <!-- OAuth 配置标签页 -->
        <el-tab-pane label="OAuth 配置" name="oauth">
          <OAuthSettingsTab v-model:form="oauthForm" :loading="loading || !canEditSettings" />
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
import BlogSettingsTab from './components/BlogSettingsTab.vue';
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
const blogTabRef = ref<InstanceType<typeof BlogSettingsTab>>();
const basicTabRef = ref<InstanceType<typeof BasicSettingsTab>>();

// 基本配置表单
const basicForm = ref({
  author: '',
  author_email: '',
  author_desc: '',
  author_avatar: '',
  author_photo: '',
  icp: '',
  police_record: '',
  admin_url: '',
  blog_url: '',
  home_url: '',
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

// 博客配置表单
const blogForm = ref({
  // 博客网站信息
  title: '',
  subtitle: '',
  slogan: '',
  description: '',
  keywords: '',
  established: '',

  // 全局样式
  favicon: '',
  background_image: '',
  screenshot: '',
  announcement: '',
  typingTextsList: [] as Array<{ value: string }>,

  // 社交媒体
  sidebarSocialList: [] as Array<{ name: string; url: string; icon: string }>,
  footerSocialList: [] as Array<{
    name: string;
    url: string;
    icon: string;
    position: string;
  }>,

  // 页脚链接
  footerLinksList: [] as Array<{ name: string; url: string }>,

  // 页面配置
  moments_size: 30,
  message_content: '',
  home_layout: 'waterfall',

  // 关于页面配置
  about_describe: '',
  about_describe_tips: '',
  about_exhibition: '',
  profileList: [] as Array<{ label: string; value: string; color: string }>,
  about_personality: '',
  mottoMainList: [] as string[],
  about_motto_sub: '',
  socializeList: [] as Array<{ name: string; url: string }>,
  creationList: [] as Array<{ name: string; url: string }>,
  versionsList: [] as Array<{ name: string; version: string }>,
  unionsList: [] as Array<{ name: string; url: string }>,
  about_story: '',
  custom_head: '',
  custom_body: '',
  emojis: '',
  font: '',
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
    Object.assign(basicForm.value, {
      author: configs.author || '',
      author_email: configs.author_email || '',
      author_desc: configs.author_desc || '',
      author_avatar: configs.author_avatar || '',
      author_photo: configs.author_photo || '',
      icp: configs.icp || '',
      police_record: configs.police_record || '',
      admin_url: configs.admin_url || '',
      blog_url: configs.blog_url || '',
      home_url: configs.home_url || '',
    });
  } catch {
    ElMessage.error('获取基本配置失败');
  }
};

// 加载博客配置
const loadBlogConfigs = async () => {
  try {
    const configs = await loadConfigs('blog');

    // 博客网站信息
    Object.assign(blogForm.value, {
      title: configs.title || '',
      subtitle: configs.subtitle || '',
      slogan: configs.slogan || '',
      description: configs.description || '',
      keywords: configs.keywords || '',
      established: configs.established || '',

      // 全局样式
      favicon: configs.favicon || '',
      background_image: configs.background_image || '',
      screenshot: configs.screenshot || '',
      announcement: configs.announcement || '',

      // 关于页面配置
      about_describe: configs.about_describe || '',
      about_describe_tips: configs.about_describe_tips || '',
      about_exhibition: configs.about_exhibition || '',
      about_personality: configs.about_personality || '',
      about_motto_sub: configs.about_motto_sub || '',
      about_story: configs.about_story || '',
      moments_size: Number(configs.moments_size) || 30,
      message_content: configs.message_content || '',
      home_layout: configs.home_layout || 'waterfall',
    });

    // 解析 JSON 字段
    const parsed = parseJSON<(string | { value: string } | { label: string; value: string })[]>(
      configs.typing_texts || '',
      []
    );
    blogForm.value.typingTextsList = parsed.map(item =>
      typeof item === 'string' ? { value: item } : item
    );

    blogForm.value.sidebarSocialList = parseJSON(configs.sidebar_social || '', []);
    blogForm.value.footerSocialList = parseJSON(configs.footer_social || '', []);
    blogForm.value.footerLinksList = parseJSON(configs.footer_links || '', []);

    blogForm.value.profileList = Array(6)
      .fill(null)
      .map(
        (_, i) =>
          parseJSON(configs.about_profile || '', [])[i] || {
            label: '',
            value: '',
            color: '#43a6c6',
          }
      );

    blogForm.value.mottoMainList = Array(2)
      .fill(null)
      .map((_, i) => parseJSON(configs.about_motto_main || '', [])[i] || '');

    blogForm.value.socializeList = parseJSON(configs.about_socialize || '', []);
    blogForm.value.creationList = parseJSON(configs.about_creation || '', []);

    blogForm.value.versionsList = Array(3)
      .fill(null)
      .map(
        (_, i) =>
          parseJSON(configs.about_versions || '', [])[i] || {
            name: '',
            version: '',
          }
      );

    blogForm.value.unionsList = parseJSON(configs.about_unions || '', []);
    blogForm.value.custom_head = configs.custom_head || '';
    blogForm.value.custom_body = configs.custom_body || '';
    blogForm.value.emojis = configs.emojis || '';
    blogForm.value.font = configs.font || '';
  } catch {
    ElMessage.error('获取博客配置失败');
  }
};

// 加载通知配置
const loadNotificationConfigs = async () => {
  try {
    const configs = await loadConfigs('notification');
    Object.assign(notificationForm.value, {
      email_host: configs.email_host || '',
      email_port: configs.email_port || '465',
      email_secure: configs.email_secure || 'ssl',
      email_username: configs.email_username || '',
      email_from: configs.email_from || '',
      email_password: configs.email_password || '',
      feishu_app_id: configs.feishu_app_id || '',
      feishu_secret: configs.feishu_secret || '',
      feishu_chat_id: configs.feishu_chat_id || '',
    });
  } catch {
    ElMessage.error('获取通知配置失败');
  }
};

// 加载上传配置
const loadUploadConfigs = async () => {
  try {
    const configs = await loadConfigs('upload');
    Object.assign(uploadForm.value, {
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
    });
  } catch {
    ElMessage.error('获取上传配置失败');
  }
};

// JSON 解析辅助函数
const parseJSON = <T,>(jsonStr: string, fallback: T): T => {
  try {
    return jsonStr ? JSON.parse(jsonStr) : fallback;
  } catch {
    return fallback;
  }
};

// 加载 AI 配置
const loadAIConfigs = async () => {
  try {
    const configs = await loadConfigs('ai');
    Object.assign(aiForm.value, {
      base_url: configs.base_url || '',
      api_key: configs.api_key || '',
      model: configs.model || '',
      summary_prompt: configs.summary_prompt || '',
      ai_summary_prompt: configs.ai_summary_prompt || '',
      title_prompt: configs.title_prompt || '',
      mcp_secret: configs.mcp_secret || '',
    });
  } catch {
    ElMessage.error('获取 AI 配置失败');
  }
};

// 加载 OAuth 配置
const loadOAuthConfigs = async () => {
  try {
    const configs = await loadConfigs('oauth');
    Object.assign(oauthForm.value, {
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
    });
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
      loadBlogConfigs(),
      loadNotificationConfigs(),
      loadUploadConfigs(),
      loadAIConfigs(),
      loadOAuthConfigs(),
    ]);
  } finally {
    loading.value = false;
  }
};

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
      if (basicUploaders.authorPhotoUploaderRef?.getPendingCount()) {
        uploadPromises.push(
          basicUploaders.authorPhotoUploaderRef.uploadPendingFile().then(url => {
            if (url) basicForm.value.author_photo = url;
          })
        );
      }
    }

    const blogUploaders = blogTabRef.value;
    if (blogUploaders) {
      if (blogUploaders.faviconUploaderRef?.getPendingCount()) {
        uploadPromises.push(
          blogUploaders.faviconUploaderRef.uploadPendingFile().then(url => {
            if (url) blogForm.value.favicon = url;
          })
        );
      }
      if (blogUploaders.backgroundUploaderRef?.getPendingCount()) {
        uploadPromises.push(
          blogUploaders.backgroundUploaderRef.uploadPendingFile().then(url => {
            if (url) blogForm.value.background_image = url;
          })
        );
      }
      if (blogUploaders.screenshotUploaderRef?.getPendingCount()) {
        uploadPromises.push(
          blogUploaders.screenshotUploaderRef.uploadPendingFile().then(url => {
            if (url) blogForm.value.screenshot = url;
          })
        );
      }
      if (blogUploaders.aboutExhibitionUploaderRef?.getPendingCount()) {
        uploadPromises.push(
          blogUploaders.aboutExhibitionUploaderRef.uploadPendingFile().then(url => {
            if (url) blogForm.value.about_exhibition = url;
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
      'basic.author': basicForm.value.author,
      'basic.author_email': basicForm.value.author_email,
      'basic.author_desc': basicForm.value.author_desc,
      'basic.author_avatar': basicForm.value.author_avatar,
      'basic.author_photo': basicForm.value.author_photo,
      'basic.icp': basicForm.value.icp,
      'basic.police_record': basicForm.value.police_record,
      'basic.admin_url': basicForm.value.admin_url,
      'basic.blog_url': basicForm.value.blog_url,
      'basic.home_url': basicForm.value.home_url,
    };

    // 博客配置
    const blogPayload: Record<string, string> = {
      'blog.title': blogForm.value.title,
      'blog.subtitle': blogForm.value.subtitle,
      'blog.slogan': blogForm.value.slogan,
      'blog.description': blogForm.value.description,
      'blog.keywords': blogForm.value.keywords,
      'blog.established': blogForm.value.established,
      'blog.favicon': blogForm.value.favicon,
      'blog.background_image': blogForm.value.background_image,
      'blog.screenshot': blogForm.value.screenshot,
      'blog.announcement': blogForm.value.announcement,
      'blog.typing_texts': JSON.stringify(blogForm.value.typingTextsList.map(item => item.value)),
      'blog.sidebar_social': JSON.stringify(blogForm.value.sidebarSocialList),
      'blog.footer_social': JSON.stringify(blogForm.value.footerSocialList),
      'blog.footer_links': JSON.stringify(blogForm.value.footerLinksList),
      'blog.about_describe': blogForm.value.about_describe,
      'blog.about_describe_tips': blogForm.value.about_describe_tips,
      'blog.about_exhibition': blogForm.value.about_exhibition,
      'blog.about_profile': JSON.stringify(blogForm.value.profileList),
      'blog.about_personality': blogForm.value.about_personality,
      'blog.about_motto_main': JSON.stringify(blogForm.value.mottoMainList),
      'blog.about_motto_sub': blogForm.value.about_motto_sub,
      'blog.about_socialize': JSON.stringify(blogForm.value.socializeList),
      'blog.about_creation': JSON.stringify(blogForm.value.creationList),
      'blog.about_versions': JSON.stringify(blogForm.value.versionsList),
      'blog.about_unions': JSON.stringify(blogForm.value.unionsList),
      'blog.about_story': blogForm.value.about_story,
      'blog.custom_head': blogForm.value.custom_head,
      'blog.custom_body': blogForm.value.custom_body,
      'blog.emojis': blogForm.value.emojis,
      'blog.font': blogForm.value.font,
      'blog.moments_size': String(blogForm.value.moments_size),
      'blog.message_content': blogForm.value.message_content,
      'blog.home_layout': blogForm.value.home_layout,
    };

    // 通知配置
    const notificationPayload: Record<string, string> = {
      'notification.email_host': notificationForm.value.email_host,
      'notification.email_port': String(notificationForm.value.email_port),
      'notification.email_secure': notificationForm.value.email_secure,
      'notification.email_username': notificationForm.value.email_username,
      'notification.email_from': notificationForm.value.email_from,
      'notification.email_password': notificationForm.value.email_password,
      'notification.feishu_app_id': notificationForm.value.feishu_app_id,
      'notification.feishu_secret': notificationForm.value.feishu_secret,
      'notification.feishu_chat_id': notificationForm.value.feishu_chat_id,
    };

    // 上传配置
    const uploadPayload: Record<string, string> = {
      'upload.storage_type': uploadForm.value.storage_type,
      'upload.max_file_size': String(uploadForm.value.max_file_size),
      'upload.path_pattern': uploadForm.value.path_pattern,
      'upload.access_key': uploadForm.value.access_key,
      'upload.secret_key': uploadForm.value.secret_key,
      'upload.region': uploadForm.value.region,
      'upload.bucket': uploadForm.value.bucket,
      'upload.endpoint': uploadForm.value.endpoint,
      'upload.domain': uploadForm.value.domain,
      'upload.use_ssl': uploadForm.value.use_ssl ? 'true' : 'false',
    };

    // AI 配置
    const aiPayload: Record<string, string> = {
      'ai.base_url': aiForm.value.base_url,
      'ai.api_key': aiForm.value.api_key,
      'ai.model': aiForm.value.model,
      'ai.summary_prompt': aiForm.value.summary_prompt,
      'ai.ai_summary_prompt': aiForm.value.ai_summary_prompt,
      'ai.title_prompt': aiForm.value.title_prompt,
    };

    // OAuth 配置
    const oauthPayload: Record<string, string> = {
      'oauth.github.enabled': oauthForm.value['github.enabled'],
      'oauth.github.client_id': oauthForm.value['github.client_id'],
      'oauth.github.client_secret': oauthForm.value['github.client_secret'],
      'oauth.github.redirect_url': oauthForm.value['github.redirect_url'],
      'oauth.google.enabled': oauthForm.value['google.enabled'],
      'oauth.google.client_id': oauthForm.value['google.client_id'],
      'oauth.google.client_secret': oauthForm.value['google.client_secret'],
      'oauth.google.redirect_url': oauthForm.value['google.redirect_url'],
      'oauth.qq.enabled': oauthForm.value['qq.enabled'],
      'oauth.qq.client_id': oauthForm.value['qq.client_id'],
      'oauth.qq.client_secret': oauthForm.value['qq.client_secret'],
      'oauth.qq.redirect_url': oauthForm.value['qq.redirect_url'],
      'oauth.microsoft.enabled': oauthForm.value['microsoft.enabled'],
      'oauth.microsoft.client_id': oauthForm.value['microsoft.client_id'],
      'oauth.microsoft.client_secret': oauthForm.value['microsoft.client_secret'],
      'oauth.microsoft.redirect_url': oauthForm.value['microsoft.redirect_url'],
    };

    // 构建需要保存的配置组列表
    const savePromises = [
      updateSettingGroup('basic', basicPayload),
      updateSettingGroup('blog', blogPayload),
      updateSettingGroup('notification', notificationPayload),
      updateSettingGroup('upload', uploadPayload),
      updateSettingGroup('ai', aiPayload),
      updateSettingGroup('oauth', oauthPayload),
    ];

    // 并行保存所有配置组
    await Promise.all(savePromises);

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
  'blog',
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
