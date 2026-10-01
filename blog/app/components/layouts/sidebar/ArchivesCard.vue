<script setup lang="ts">
const { isExpanded, toggleExpand } = useExpandable();
const { data: archives } = await useArchives();

const displayArchives = computed(() => {
  const arr = archives.value ?? [];
  const list = arr
    .slice(0, 6)
    .map(a => ({ ...a, displayText: `${a.year} ${a.month}`, isEarlier: false }));
  if (arr.length > 6) {
    const earlierCount = arr.slice(6).reduce((s, a) => s + a.count, 0);
    list.push({
      year: '',
      month: '',
      displayText: '在此之前',
      count: earlierCount,
      isEarlier: true,
    });
  }
  return list;
});
</script>

<template>
  <div class="card-widget card-archives">
    <div class="item-headline" :class="{ 'is-expanded': isExpanded }">
      <i class="ri-archive-fill" />
      <span>归档</span>
      <i
        class="collapse-icon ri-arrow-left-s-fill"
        :class="{ 'is-expanded': isExpanded }"
        @click="toggleExpand"
      />
    </div>
    <ul class="card-list" :class="{ 'is-expanded': isExpanded }">
      <li v-for="(archive, i) in displayArchives" :key="i" class="card-list-item">
        <router-link
          class="card-list-link"
          :to="archive.isEarlier ? '/archive' : `/archive/${archive.year}/${archive.month}`"
        >
          <span class="card-list-name">{{ archive.displayText }}</span>
          <span class="card-list-count">{{ archive.count }}</span>
        </router-link>
      </li>
    </ul>
  </div>
</template>
