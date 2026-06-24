<template>
  <el-form :model="form" label-width="120px" class="setting-form">
    <el-divider content-position="left">基础信息</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('title') }">博客标题</span>
      </template>
      <el-input v-model="form.title" placeholder="博客标题" :disabled="loading" />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('subtitle') }">博客副标题</span>
      </template>
      <el-input v-model="form.subtitle" placeholder="博客副标题" :disabled="loading" />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('description') }">博客描述</span>
      </template>
      <el-input
        v-model="form.description"
        type="textarea"
        :rows="2"
        placeholder="用于 SEO 的博客描述"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('keywords') }">博客关键词</span>
      </template>
      <el-input v-model="form.keywords" placeholder="关键词，用逗号分隔" :disabled="loading" />
    </el-form-item>

    <div class="image-row">
      <el-form-item>
        <template #label>
          <span :class="{ 'field-modified': isFieldModified('favicon') }">网站 Favicon</span>
        </template>
        <ImageUploader
          ref="faviconUploaderRef"
          v-model="form.favicon"
          upload-type="博客图标"
          width="120px"
          height="120px"
          :disabled="loading"
        />
      </el-form-item>
    </div>

    <el-divider content-position="left">站长信息</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('author') }">站长姓名</span>
      </template>
      <el-input v-model="form.author" placeholder="站长姓名" :disabled="loading" />
    </el-form-item>

    <div class="image-row">
      <el-form-item>
        <template #label>
          <span :class="{ 'field-modified': isFieldModified('author_avatar') }">站长头像</span>
        </template>
        <ImageUploader
          ref="authorAvatarUploaderRef"
          v-model="form.author_avatar"
          upload-type="站长头像"
          width="120px"
          height="120px"
          :disabled="loading"
        />
      </el-form-item>
    </div>

    <el-divider content-position="left">备案信息</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('icp') }">ICP备案号</span>
      </template>
      <el-input v-model="form.icp" placeholder="ICP备案号" :disabled="loading" />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('police_record') }">公安备案号</span>
      </template>
      <el-input v-model="form.police_record" placeholder="公安备案号" :disabled="loading" />
    </el-form-item>

    <el-divider content-position="left">扩展功能</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('emojis') }">表情包配置</span>
      </template>
      <el-input v-model="form.emojis" placeholder="表情包 URL" :disabled="loading" />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('meting_api') }">Meting-API</span>
      </template>
      <el-input v-model="form.meting_api" placeholder="Meting-API 地址" :disabled="loading" />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('cravatar_url') }">Cravatar URL</span>
      </template>
      <el-input
        v-model="form.cravatar_url"
        placeholder="头像服务 URL（%s 为邮箱哈希）"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('ip_api_url') }">IP 归属地 API</span>
      </template>
      <el-input
        v-model="form.ip_api_url"
        placeholder="IP 归属地查询 URL（%s 为 IP）"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('cover_maker_api') }">封面制作 API</span>
      </template>
      <el-input
        v-model="form.cover_maker_api"
        placeholder="封面制作图片源 API"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">系统地址</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('admin_url') }">管理地址</span>
      </template>
      <el-input
        v-model="form.admin_url"
        placeholder="例如 https://admin.your-site.com"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('blog_url') }">博客地址</span>
      </template>
      <el-input
        v-model="form.blog_url"
        placeholder="例如 https://blog.your-site.com"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">自定义代码</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('custom_head') }">自定义 Head</span>
      </template>
      <el-input
        v-model="form.custom_head"
        type="textarea"
        :rows="4"
        placeholder="注入到 &lt;head&gt; 的自定义 HTML 代码"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('custom_body') }">自定义 Body</span>
      </template>
      <el-input
        v-model="form.custom_body"
        type="textarea"
        :rows="4"
        placeholder="注入到 &lt;body&gt; 的自定义 HTML 代码"
        :disabled="loading"
      />
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import ImageUploader from '@/components/common/ImageUploader.vue';

interface BasicForm {
  author: string;
  author_avatar: string;
  icp: string;
  police_record: string;
  admin_url: string;
  blog_url: string;
  title: string;
  subtitle: string;
  description: string;
  keywords: string;
  favicon: string;
  custom_head: string;
  custom_body: string;
  emojis: string;
  meting_api: string;
  cravatar_url: string;
  ip_api_url: string;
  cover_maker_api: string;
}

const form = defineModel<BasicForm>('form', { required: true });

defineProps<{
  loading?: boolean;
  isFieldModified: (key: string) => boolean;
}>();

// 图片上传器引用
const authorAvatarUploaderRef = ref<InstanceType<typeof ImageUploader>>();
const faviconUploaderRef = ref<InstanceType<typeof ImageUploader>>();

// 暴露给父组件使用
defineExpose({
  authorAvatarUploaderRef,
  faviconUploaderRef,
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
}

// 移动端适配
@media (max-width: 768px) {
  .setting-form {
    .image-row {
      flex-direction: column;
      gap: 0;
    }
  }

  :deep(.el-form-item__label) {
    width: 100px !important;
    font-size: 13px;
  }
}
</style>
