<script setup lang="ts">
const { tags } = useTags();
const { isExpanded, toggleExpand } = useExpandable();

const tagCloudRef = ref<HTMLElement | null>(null);
const needsExpand = ref(false);

const updateNeedsExpand = () => {
  needsExpand.value = tagCloudRef.value ? tagCloudRef.value.scrollHeight > 150 : false;
};

onMounted(updateNeedsExpand);
watch(tags, () => nextTick(updateNeedsExpand), { deep: true });

const getTagSize = (count: number) => {
  const maxCount = Math.max(...tags.value.map(t => t.count), 1);
  return `${0.9 + 0.6 * (count / maxCount)}em`;
};
</script>

<template>
  <div class="card-widget card-tags">
    <div class="item-headline" :class="{ 'is-expanded': isExpanded }">
      <i class="ri-price-tag-3-fill" />
      <span>标签</span>
      <i
        v-if="needsExpand"
        class="collapse-icon ri-arrow-left-s-fill"
        :class="{ 'is-expanded': isExpanded }"
        @click="toggleExpand"
      />
    </div>
    <div
      ref="tagCloudRef"
      class="card-tag-cloud"
      :class="{ 'is-expanded': isExpanded, 'can-expand': needsExpand }"
    >
      <router-link
        v-for="tag in tags"
        :key="tag.id"
        :to="tag.url"
        :style="{ fontSize: getTagSize(tag.count) }"
        :aria-label="`查看标签：${tag.name}，共 ${tag.count} 篇文章`"
      >
        {{ tag.name }}
      </router-link>
    </div>
    <div v-if="needsExpand" class="expand-toggle" @click="toggleExpand">
      <span>{{ isExpanded ? '收起' : '展开更多' }}</span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.card-tag-cloud {
  a {
    font-size: 1.1em;
    color: #666;
    display: inline-block;
    padding: 0 4px;
    line-height: 1.8;

    &:hover {
      color: var(--flec-btn-hover);
    }
  }

  &.can-expand:not(.is-expanded) {
    max-height: 200px;
    overflow: hidden;
    position: relative;

    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      height: 30px;
      pointer-events: none;
    }
  }

  &.is-expanded {
    max-height: none;
  }
}

.expand-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  margin-top: 8px;
  color: #999;
  font-size: 0.9em;
  cursor: pointer;
  transition: color 0.3s;

  &:hover {
    color: var(--flec-btn-hover);
  }
}

@media (max-width: 900px) {
  .expand-toggle {
    display: none;
  }

  .card-tag-cloud.can-expand:not(.is-expanded) {
    max-height: none;
    overflow: visible;

    &::after {
      display: none;
    }
  }
}
</style>
