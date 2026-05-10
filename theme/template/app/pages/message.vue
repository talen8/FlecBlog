<template>
  <div id="page">
    <h1 class="page-title">留言板</h1>
    <div class="message-content">
      <p>{{ blogConfig['blog.message_content'] || '欢迎留言交流~' }}</p>
    </div>
    <FeaturesCommentComments target-type="page" target-key="message" />
  </div>
</template>

<script lang="ts" setup>
const { blogConfig } = useSysConfig();

useSeoMeta({
  title: '留言板',
  description: '欢迎留言交流',
});
</script>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border-radius: 12px;
  border: 1px solid var(--flec-border);
  padding: 40px;
}

.page-title {
  margin: 0 0 30px;
  font-size: 2rem;
  font-weight: 700;
}

.message-content {
  line-height: 1.8;
  color: var(--font-color);
  margin-bottom: 40px;
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
