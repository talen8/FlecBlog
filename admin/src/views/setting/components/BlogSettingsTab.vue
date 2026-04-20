<template>
  <el-form :model="form" label-width="120px" class="setting-form">
    <el-divider content-position="left">网站信息</el-divider>

    <el-form-item label="网站标题">
      <el-input
        v-model="form.title"
        placeholder="用于RSS订阅和邮件显示的站点标题"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="网站副标题">
      <el-input v-model="form.subtitle" placeholder="网站的副标题或简短描述" :disabled="loading" />
    </el-form-item>

    <el-form-item label="网站标语">
      <el-input v-model="form.slogan" placeholder="网站标语或座右铭" :disabled="loading" />
    </el-form-item>

    <el-form-item label="网站描述">
      <el-input
        v-model="form.description"
        type="textarea"
        :rows="3"
        placeholder="网站描述，用于SEO"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="关键词">
      <el-input
        v-model="form.keywords"
        placeholder="网站关键词，多个用逗号分隔"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="建站日期">
      <el-date-picker
        v-model="form.established"
        type="date"
        placeholder="选择建站日期"
        format="YYYY-MM-DD"
        value-format="YYYY-MM-DD"
        :disabled="loading"
        style="width: 100%"
      />
    </el-form-item>

    <el-divider content-position="left">全局样式</el-divider>

    <div class="image-row">
      <el-form-item label="网站Favicon">
        <ImageUploader
          ref="faviconUploaderRef"
          v-model="form.favicon"
          upload-type="博客图标"
          width="120px"
          height="120px"
          :disabled="loading"
        />
      </el-form-item>

      <el-form-item>
        <template #label>
          <span class="clickable-label" @click="toggleBackgroundMode"> 背景图片 </span>
        </template>
        <div class="background-image-input">
          <ImageUploader
            v-if="backgroundInputMode === 'upload'"
            ref="backgroundUploaderRef"
            v-model="form.background_image"
            upload-type="博客背景"
            width="213px"
            height="120px"
            :disabled="loading"
          />
          <div v-else class="url-input-wrapper">
            <el-input
              v-model="form.background_image"
              placeholder="输入图片链接"
              :disabled="loading"
              clearable
            >
              <template #prefix>
                <i class="ri-link"></i>
              </template>
            </el-input>
          </div>
        </div>
      </el-form-item>

      <el-form-item label="站点截图">
        <ImageUploader
          ref="screenshotUploaderRef"
          v-model="form.screenshot"
          upload-type="博客截图"
          width="213px"
          height="120px"
          :disabled="loading"
        />
      </el-form-item>
    </div>

    <el-form-item label="侧边栏公告">
      <el-input
        v-model="form.announcement"
        type="textarea"
        :rows="4"
        placeholder="支持多行文案和简单 HTML"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="打字机文本">
      <JsonListEditor
        v-model="form.typingTextsList"
        :fields="typingTextsFields"
        :default-item="{ value: '' }"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">社交媒体</el-divider>

    <el-form-item label="侧边栏社交">
      <JsonListEditor
        v-model="form.sidebarSocialList"
        :fields="sidebarSocialFields"
        :default-item="{ name: '', url: '', icon: '' }"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="页脚社交">
      <JsonListEditor
        v-model="form.footerSocialList"
        :fields="footerSocialFields"
        :default-item="{ name: '', url: '', icon: '', position: 'left' }"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="页脚链接">
      <JsonListEditor
        v-model="form.footerLinksList"
        :fields="footerLinksFields"
        :default-item="{ name: '', url: '' }"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">页面配置</el-divider>

    <el-form-item label="动态数量">
      <el-input-number
        v-model="form.moments_size"
        :min="1"
        :disabled="loading"
        placeholder="动态列表每页显示数量"
        style="width: 200px"
      />
    </el-form-item>

    <el-form-item label="首页布局">
      <el-select
        v-model="form.home_layout"
        :disabled="loading"
        placeholder="选择首页布局"
        style="width: 200px"
      >
        <el-option label="瀑布流" value="waterfall" />
        <el-option label="单列布局" value="single_column" />
      </el-select>
    </el-form-item>

    <el-form-item label="留言信">
      <el-input
        v-model="form.message_content"
        type="textarea"
        :rows="8"
        placeholder="每行一段文字，显示在留言页面的信封中"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">关于页面配置</el-divider>

    <el-form-item label="个人描述">
      <el-input
        v-model="form.about_describe"
        type="textarea"
        :rows="3"
        placeholder="关于页面的个人描述"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="描述提示">
      <el-input
        v-model="form.about_describe_tips"
        placeholder="例如：前端工程师 · 业余 · 专注 · 享受生活"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="展览图片">
      <ImageUploader
        ref="aboutExhibitionUploaderRef"
        v-model="form.about_exhibition"
        upload-type="展览图片"
        width="213px"
        height="120px"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="个人资料">
      <JsonListEditor
        v-model="form.profileList"
        :fields="profileFields"
        :disabled="loading"
        hide-controls
      />
    </el-form-item>

    <el-form-item label="性格类型">
      <el-select
        v-model="form.about_personality"
        placeholder="请选择性格类型"
        clearable
        :disabled="loading"
        style="width: 100%"
      >
        <el-option-group label="分析家">
          <el-option label="INTJ-A (建筑师-自信型)" value="INTJ-A" />
          <el-option label="INTJ-T (建筑师-波动型)" value="INTJ-T" />
          <el-option label="INTP-A (逻辑学家-自信型)" value="INTP-A" />
          <el-option label="INTP-T (逻辑学家-波动型)" value="INTP-T" />
          <el-option label="ENTJ-A (指挥官-自信型)" value="ENTJ-A" />
          <el-option label="ENTJ-T (指挥官-波动型)" value="ENTJ-T" />
          <el-option label="ENTP-A (辩论家-自信型)" value="ENTP-A" />
          <el-option label="ENTP-T (辩论家-波动型)" value="ENTP-T" />
        </el-option-group>
        <el-option-group label="外交家">
          <el-option label="INFJ-A (提倡者-自信型)" value="INFJ-A" />
          <el-option label="INFJ-T (提倡者-波动型)" value="INFJ-T" />
          <el-option label="INFP-A (调停者-自信型)" value="INFP-A" />
          <el-option label="INFP-T (调停者-波动型)" value="INFP-T" />
          <el-option label="ENFJ-A (主人公-自信型)" value="ENFJ-A" />
          <el-option label="ENFJ-T (主人公-波动型)" value="ENFJ-T" />
          <el-option label="ENFP-A (竞选者-自信型)" value="ENFP-A" />
          <el-option label="ENFP-T (竞选者-波动型)" value="ENFP-T" />
        </el-option-group>
        <el-option-group label="守护者">
          <el-option label="ISTJ-A (物流师-自信型)" value="ISTJ-A" />
          <el-option label="ISTJ-T (物流师-波动型)" value="ISTJ-T" />
          <el-option label="ISFJ-A (守卫者-自信型)" value="ISFJ-A" />
          <el-option label="ISFJ-T (守卫者-波动型)" value="ISFJ-T" />
          <el-option label="ESTJ-A (总经理-自信型)" value="ESTJ-A" />
          <el-option label="ESTJ-T (总经理-波动型)" value="ESTJ-T" />
          <el-option label="ESFJ-A (执政官-自信型)" value="ESFJ-A" />
          <el-option label="ESFJ-T (执政官-波动型)" value="ESFJ-T" />
        </el-option-group>
        <el-option-group label="探险家">
          <el-option label="ISTP-A (鉴赏家-自信型)" value="ISTP-A" />
          <el-option label="ISTP-T (鉴赏家-波动型)" value="ISTP-T" />
          <el-option label="ISFP-A (探险家-自信型)" value="ISFP-A" />
          <el-option label="ISFP-T (探险家-波动型)" value="ISFP-T" />
          <el-option label="ESTP-A (企业家-自信型)" value="ESTP-A" />
          <el-option label="ESTP-T (企业家-波动型)" value="ESTP-T" />
          <el-option label="ESFP-A (表演者-自信型)" value="ESFP-A" />
          <el-option label="ESFP-T (表演者-波动型)" value="ESFP-T" />
        </el-option-group>
      </el-select>
    </el-form-item>

    <el-form-item label="座右铭">
      <div class="motto-inputs">
        <el-input v-model="form.mottoMainList[0]" placeholder="第 1 行" :disabled="loading" />
        <el-input v-model="form.mottoMainList[1]" placeholder="第 2 行" :disabled="loading" />
      </div>
    </el-form-item>

    <el-form-item label="一言">
      <el-input v-model="form.about_motto_sub" placeholder="一句话介绍" :disabled="loading" />
    </el-form-item>

    <el-form-item label="联系方式">
      <JsonListEditor
        v-model="form.socializeList"
        :fields="socializeFields"
        :default-item="{ name: '', url: '' }"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="创作平台">
      <JsonListEditor
        v-model="form.creationList"
        :fields="creationFields"
        :default-item="{ name: '', url: '' }"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="版本信息">
      <JsonListEditor
        v-model="form.versionsList"
        :fields="versionsFields"
        :disabled="loading"
        hide-controls
      />
    </el-form-item>

    <el-form-item label="站长联盟">
      <JsonListEditor
        v-model="form.unionsList"
        :fields="unionsFields"
        :default-item="{ name: '', url: '' }"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="心路历程">
      <el-input
        v-model="form.about_story"
        type="textarea"
        :rows="6"
        placeholder="关于本站的介绍和心路历程"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">自定义代码</el-divider>

    <el-form-item label="字体配置">
      <el-input
        v-model="form.font"
        placeholder="字体文件URL|字体名称（https://hanzi.bluu.pl/fonts/986/result.css|LXGW WenKai）"
        :disabled="loading"
      >
        <template #append>
          <el-dropdown trigger="click" @command="handleFontSiteCommand">
            <el-button> 查找字体 </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="https://hanzi.bluu.pl/">
                  <i class="ri-global-line" style="margin-right: 8px"></i>
                  千字网
                </el-dropdown-item>
                <el-dropdown-item command="https://fonts.zeoseven.com/">
                  <i class="ri-global-line" style="margin-right: 8px"></i>
                  ZSFT
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-input>
    </el-form-item>

    <el-form-item label="自定义 Head">
      <el-input
        v-model="form.custom_head"
        type="textarea"
        :rows="10"
        placeholder="输入要在 <head> 中插入的代码"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label="自定义 Body">
      <el-input
        v-model="form.custom_body"
        type="textarea"
        :rows="10"
        placeholder="输入要在 <body> 中插入的代码"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">表情包配置</el-divider>

    <el-form-item label="表情包">
      <el-input
        v-model="form.emojis"
        placeholder="输入表情包配置（JSON文件）"
        :disabled="loading"
      />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import ImageUploader from '@/components/common/ImageUploader.vue';
import JsonListEditor from '@/components/common/JsonListEditor.vue';
import type { FieldConfig } from '@/components/common/JsonListEditor.vue';

interface BlogFormData {
  // 博客网站信息
  title: string;
  subtitle: string;
  slogan: string;
  description: string;
  keywords: string;
  established: string;

  // 全局样式
  favicon: string;
  background_image: string;
  screenshot: string;
  announcement: string;
  typingTextsList: Array<{ value: string }>;

  // 社交媒体
  sidebarSocialList: Array<{ name: string; url: string; icon: string }>;
  footerSocialList: Array<{
    name: string;
    url: string;
    icon: string;
    position: string;
  }>;

  // 页脚链接
  footerLinksList: Array<{ name: string; url: string }>;

  // 页面配置
  moments_size: number;
  message_content: string;
  home_layout: string;

  // 关于页面配置
  about_describe: string;
  about_describe_tips: string;
  about_exhibition: string;
  profileList: Array<{ label: string; value: string; color: string }>;
  about_personality: string;
  mottoMainList: string[];
  about_motto_sub: string;
  socializeList: Array<{ name: string; url: string }>;
  creationList: Array<{ name: string; url: string }>;
  versionsList: Array<{ name: string; version: string }>;
  unionsList: Array<{ name: string; url: string }>;
  about_story: string;
  custom_head: string;
  custom_body: string;
  emojis: string;
  font: string;
}

const form = defineModel<BlogFormData>('form', { required: true });

defineProps<{
  loading?: boolean;
}>();

// 图片上传器引用
const faviconUploaderRef = ref<InstanceType<typeof ImageUploader>>();
const backgroundUploaderRef = ref<InstanceType<typeof ImageUploader>>();
const screenshotUploaderRef = ref<InstanceType<typeof ImageUploader>>();
const aboutExhibitionUploaderRef = ref<InstanceType<typeof ImageUploader>>();

// 背景图片输入模式
const backgroundInputMode = ref<'upload' | 'url'>('upload');

const toggleBackgroundMode = () => {
  backgroundInputMode.value = backgroundInputMode.value === 'upload' ? 'url' : 'upload';
};

// 预设的常用社交平台图标
const commonIcons = [
  'github-line',
  'mail-line',
  'twitter-x-line',
  'bilibili-line',
  'wechat-line',
  'qq-line',
  'weibo-line',
  'zhihu-line',
  'douban-line',
  'linkedin-line',
  'facebook-line',
  'instagram-line',
  'youtube-line',
  'tiktok-line',
  'discord-line',
  'telegram-line',
  'slack-line',
  'rss-line',
  'links-line',
];

const iconOptions = commonIcons.map(icon => ({
  label: icon,
  value: icon,
  icon: 'ri-' + icon,
}));
const iconField = {
  key: 'icon',
  type: 'select' as const,
  placeholder: '图标',
  style: 'width: 200px; margin-right: 8px',
  prefix: 'ri-',
  filterable: true,
  allowCreate: true,
  options: iconOptions,
};

// 字段配置
const typingTextsFields: FieldConfig[] = [
  { key: 'value', type: 'text', placeholder: '打字机文本', style: 'flex: 1' },
];

const sidebarSocialFields: FieldConfig[] = [
  { key: 'name', type: 'text', placeholder: '平台名称', style: 'width: 120px' },
  {
    key: 'url',
    type: 'text',
    placeholder: '链接地址',
    style: 'flex: 1; margin: 0 8px',
  },
  iconField,
];

const footerSocialFields: FieldConfig[] = [
  { key: 'name', type: 'text', placeholder: '平台名称', style: 'width: 100px' },
  {
    key: 'url',
    type: 'text',
    placeholder: '链接地址',
    style: 'flex: 1; margin: 0 8px',
  },
  iconField,
  {
    key: 'position',
    type: 'select',
    placeholder: '位置',
    style: 'width: 80px; margin-right: 8px',
    options: [
      { label: '左', value: 'left' },
      { label: '右', value: 'right' },
    ],
  },
];

const footerLinksFields: FieldConfig[] = [
  { key: 'name', type: 'text', placeholder: '链接名称', style: 'width: 120px' },
  {
    key: 'url',
    type: 'text',
    placeholder: '链接地址 (/开头为内链)',
    style: 'flex: 1; margin: 0 8px',
  },
];

const nameUrlFields: FieldConfig[] = [
  { key: 'name', type: 'text', placeholder: '平台名称', style: 'width: 120px' },
  {
    key: 'url',
    type: 'text',
    placeholder: '链接地址',
    style: 'flex: 1; margin: 0 8px',
  },
];

const socializeFields = nameUrlFields;
const creationFields = nameUrlFields;

const unionsFields: FieldConfig[] = [
  { key: 'name', type: 'text', placeholder: '联盟名称', style: 'width: 150px' },
  {
    key: 'url',
    type: 'text',
    placeholder: '链接地址',
    style: 'flex: 1; margin: 0 8px',
  },
];

const profileFields: FieldConfig[] = [
  { key: 'label', type: 'text', placeholder: '标签', style: 'width: 100px' },
  {
    key: 'value',
    type: 'text',
    placeholder: '值',
    style: 'flex: 1; margin: 0 8px',
  },
  { key: 'color', type: 'color' },
];

const versionsFields: FieldConfig[] = [
  { key: 'name', type: 'text', placeholder: '技术名称', style: 'width: 150px' },
  {
    key: 'version',
    type: 'text',
    placeholder: '版本号',
    style: 'flex: 1; margin: 0 8px',
  },
];

// 处理字体网站跳转
const handleFontSiteCommand = (url: string) => {
  window.open(url, '_blank');
};

// 暴露上传器引用给父组件
defineExpose({
  faviconUploaderRef,
  backgroundUploaderRef,
  screenshotUploaderRef,
  aboutExhibitionUploaderRef,
});
</script>

<style lang="scss" scoped>
.setting-form {
  .image-row {
    display: flex;
    gap: 40px;

    .el-form-item {
      margin-bottom: 22px;
    }
  }

  :deep(.el-form-item__content) {
    align-items: flex-start !important;
  }

  .background-image-input {
    display: flex;
    flex-direction: column;
    gap: 8px;

    .url-input-wrapper {
      width: 213px;

      .ri-link {
        font-size: 14px;
        color: #909399;
      }
    }
  }

  .clickable-label {
    cursor: pointer;
    display: inline-flex;
    align-items: center;
    gap: 8px;
    transition: color 0.2s;

    &:hover {
      color: var(--el-color-primary);
    }
  }
}

.motto-inputs {
  display: flex;
  gap: 8px;
}

// 移动端适配
@media (max-width: 768px) {
  .setting-form {
    .image-row {
      flex-direction: column;
      gap: 0;
    }

    .motto-inputs {
      flex-direction: column;
    }

    .background-image-input {
      .url-input-wrapper {
        width: 100%;
      }
    }
  }

  :deep(.el-form-item__label) {
    width: 100px !important;
    font-size: 13px;
  }
}
</style>
