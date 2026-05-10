<script lang="ts" setup>
import { renderSimpleMarkdown } from '@/utils/markdown';
import { parseJSON } from '@/utils/json';

const { blogConfig } = useSysConfig();

const aboutDescribe = computed(() => blogConfig.value?.['blog.about_describe'] || '');
const aboutStory = computed(() => blogConfig.value?.['blog.about_story'] || '');
const aboutMottoMain = computed(() => blogConfig.value?.['blog.about_motto_main'] || '');
const aboutMottoSub = computed(() => blogConfig.value?.['blog.about_motto_sub'] || '');
const aboutPersonality = computed(() => blogConfig.value?.['blog.about_personality'] || '');
const aboutExhibition = computed(() => blogConfig.value?.['blog.about_exhibition'] || '');

const aboutProfile = computed(() =>
  parseJSON<Array<{ label: string; value: string }>>(blogConfig.value?.['blog.about_profile'], [])
);
const aboutSocialize = computed(() =>
  parseJSON<Array<{ name: string; url: string; icon?: string }>>(
    blogConfig.value?.['blog.about_socialize'],
    []
  )
);
const aboutCreation = computed(() =>
  parseJSON<Array<{ name: string; url: string; icon?: string }>>(
    blogConfig.value?.['blog.about_creation'],
    []
  )
);
const aboutVersions = computed(() =>
  parseJSON<Array<{ label: string; value: string }>>(blogConfig.value?.['blog.about_versions'], [])
);
const aboutUnions = computed(() =>
  parseJSON<Array<{ name: string; url: string; icon?: string }>>(
    blogConfig.value?.['blog.about_unions'],
    []
  )
);

useSeoMeta({
  title: '关于',
  description: () => aboutDescribe.value || '关于我',
});
</script>

<template>
  <div id="page">
    <h1 class="page-title">关于</h1>

    <div v-if="aboutDescribe" class="about-section">
      <p>{{ aboutDescribe }}</p>
    </div>

    <div v-if="aboutMottoMain || aboutMottoSub" class="about-section motto">
      <h2 class="section-title">座右铭</h2>
      <blockquote v-if="aboutMottoMain" class="motto-main">{{ aboutMottoMain }}</blockquote>
      <p v-if="aboutMottoSub" class="motto-sub">{{ aboutMottoSub }}</p>
    </div>

    <div v-if="aboutPersonality" class="about-section">
      <h2 class="section-title">性格</h2>
      <p class="personality">{{ aboutPersonality }}</p>
    </div>

    <div v-if="aboutProfile.length" class="about-section">
      <h2 class="section-title">个人资料</h2>
      <div class="info-list">
        <div v-for="(item, index) in aboutProfile" :key="index" class="info-item">
          <span class="label">{{ item.label }}</span>
          <span class="value">{{ item.value }}</span>
        </div>
      </div>
    </div>

    <div v-if="aboutSocialize.length" class="about-section">
      <h2 class="section-title">联系方式</h2>
      <div class="link-list">
        <a
          v-for="(item, index) in aboutSocialize"
          :key="index"
          :href="item.url"
          target="_blank"
          rel="noopener noreferrer"
          class="link-item"
        >
          <i v-if="item.icon" :class="item.icon" />
          <span>{{ item.name }}</span>
        </a>
      </div>
    </div>

    <div v-if="aboutCreation.length" class="about-section">
      <h2 class="section-title">创作平台</h2>
      <div class="link-list">
        <a
          v-for="(item, index) in aboutCreation"
          :key="index"
          :href="item.url"
          target="_blank"
          rel="noopener noreferrer"
          class="link-item"
        >
          <i v-if="item.icon" :class="item.icon" />
          <span>{{ item.name }}</span>
        </a>
      </div>
    </div>

    <div v-if="aboutVersions.length" class="about-section">
      <h2 class="section-title">版本信息</h2>
      <div class="info-list">
        <div v-for="(item, index) in aboutVersions" :key="index" class="info-item">
          <span class="label">{{ item.label }}</span>
          <span class="value">{{ item.value }}</span>
        </div>
      </div>
    </div>

    <div v-if="aboutUnions.length" class="about-section">
      <h2 class="section-title">站长联盟</h2>
      <div class="link-list">
        <a
          v-for="(item, index) in aboutUnions"
          :key="index"
          :href="item.url"
          target="_blank"
          rel="noopener noreferrer"
          class="link-item"
        >
          <i v-if="item.icon" :class="item.icon" />
          <span>{{ item.name }}</span>
        </a>
      </div>
    </div>

    <div v-if="aboutExhibition" class="about-section">
      <h2 class="section-title">展览</h2>
      <img :src="aboutExhibition" alt="exhibition" class="exhibition-img" />
    </div>

    <div v-if="aboutStory" class="about-section">
      <h2 class="section-title">心路历程</h2>
      <!-- eslint-disable-next-line vue/no-v-html -- Markdown 渲染内容，已通过 DOMPurify 净化 -->
      <div class="story-content" v-html="renderSimpleMarkdown(aboutStory)" />
    </div>

    <FeaturesCommentComments target-type="page" target-key="about" />
  </div>
</template>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border: 1px solid var(--flec-border);
  padding: 40px;
}

.page-title {
  margin: 0 0 30px;
  font-size: 2rem;
  font-weight: 700;
}

.about-section {
  margin-bottom: 30px;
  p {
    line-height: 1.8;
    color: var(--font-color);
  }
}

.section-title {
  font-size: 1.2rem;
  font-weight: 600;
  margin: 0 0 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--flec-border);
}

.motto-main {
  margin: 0 0 8px;
  padding: 16px;
  background: var(--flec-heavy-bg);
  border-left: 4px solid var(--theme-color);
  font-size: 1.1rem;
  font-style: italic;
  color: var(--font-color);
}

.motto-sub {
  color: var(--theme-meta-color);
  font-size: 0.9rem;
}

.personality {
  display: inline-block;
  padding: 8px 16px;
  background: var(--flec-heavy-bg);
  font-weight: 500;
}

.info-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.info-item {
  display: flex;
  gap: 16px;
  padding: 8px 0;
  border-bottom: 1px dashed var(--flec-border);
  .label {
    color: var(--theme-meta-color);
    min-width: 80px;
    flex-shrink: 0;
  }
  .value {
    color: var(--font-color);
  }
}

.link-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.link-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background: var(--flec-heavy-bg);
  color: var(--font-color);
  text-decoration: none;
  font-size: 0.9rem;
}

.exhibition-img {
  max-width: 100%;
}

.story-content {
  line-height: 1.8;
  color: var(--font-color);
}

@media screen and (max-width: 768px) {
  #page {
    padding: 20px;
  }
  .page-title {
    font-size: 1.5rem;
  }
}
</style>
