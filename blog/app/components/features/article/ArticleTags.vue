<script lang="ts" setup>
defineProps<{ article: { tags: Array<{ id: number; name: string; url: string }> } }>();
const { blogConfig } = useSysConfig();
const showModal = ref(false);
const hasDonation = computed(() => {
  try {
    return JSON.parse(blogConfig.value.donation_methods || '[]').length > 0;
  } catch {
    return false;
  }
});
</script>

<template>
  <div v-if="article.tags?.length || hasDonation" class="wrapper">
    <div v-if="article.tags?.length" class="tags">
      <a v-for="tag in article.tags" :key="tag.id" :href="tag.url">{{ tag.name }}</a>
    </div>
    <button v-if="hasDonation" class="donate" @click="showModal = true">
      <i class="ri-heart-line" />
      <span>赞赏</span>
    </button>
    <FeaturesModalsDonationModal v-model="showModal" />
  </div>
</template>

<style lang="scss" scoped>
.wrapper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16px;
  margin-top: 20px;
}
.tags {
  display: flex;
  flex-wrap: nowrap;
  gap: 10px;
  flex: 1;
  min-width: 0;
  overflow-x: auto;
  overflow-y: hidden;
  scrollbar-width: none;
  -webkit-overflow-scrolling: touch;
  &::-webkit-scrollbar {
    display: none;
  }
  a {
    flex: 0 0 auto;
    padding: 1px 8px;
    border: 1px solid var(--theme-color);
    border-radius: 8px;
    color: var(--theme-color);
    text-decoration: none;
    font-size: 0.85rem;
    white-space: nowrap;
    &:hover {
      background: var(--theme-color);
      color: #fff;
    }
  }
}
.donate {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  border: 1px solid #e74c3c;
  border-radius: 8px;
  background: transparent;
  color: #e74c3c;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.3s;
  flex-shrink: 0;
  i {
    font-size: 1rem;
  }
  &:hover {
    background: #e74c3c;
    color: #fff;
  }
}
@media (max-width: 600px) {
  .wrapper {
    gap: 10px;
  }
  .donate {
    padding-inline: 10px;
  }
}
</style>
