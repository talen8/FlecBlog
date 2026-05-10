<script lang="ts" setup>
const { blogConfig } = useSysConfig();

const announcement = computed(() => {
  const content = blogConfig.value['blog.announcement']?.trim() || '';
  return content.replace(/<[^>]+>/g, '').trim();
});
</script>

<template>
  <div v-if="announcement" class="sidebar-card">
    <div class="headline"><i class="ri-megaphone-line" /> 公告</div>
    <p class="content">{{ announcement }}</p>
  </div>
</template>

<style lang="scss" scoped>
.headline {
  font-weight: 600;
  margin-bottom: 6px;
}
.headline i {
  color: #e03131;
  margin-right: 4px;
}
.content {
  margin: 0;
  font-size: 0.85rem;
  color: var(--theme-meta-color);
  white-space: pre-line;
}
</style>
