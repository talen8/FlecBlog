<script setup lang="ts">
const { data: page } = await useAsyncData('community', () => queryCollection('community').first())

const title = page.value?.seo?.title || page.value?.title
const description = page.value?.seo?.description || page.value?.description

useSeoMeta({
  title,
  ogTitle: title,
  description,
  ogDescription: description
})
</script>

<template>
  <div v-if="page">
    <UContainer>
      <UPageHeader
        v-bind="page"
        class="py-[50px]"
      />

      <div class="flex flex-col gap-4 pt-8">
        <NuxtLink
          v-for="item in page.items"
          :key="item.title"
          :to="item.link"
          target="_blank"
          class="group flex items-center gap-6 rounded-xl border border-[var(--ui-border)] p-6 bg-[var(--ui-bg)] transition-all hover:shadow-lg hover:border-[var(--ui-primary)]/30"
        >
          <!-- 图标区域 -->
          <div class="flex-shrink-0 flex items-center justify-center size-14 rounded-xl bg-[var(--ui-primary)]/10 text-[var(--ui-primary)] group-hover:bg-[var(--ui-primary)]/15 transition-colors">
            <UIcon
              :name="item.icon"
              class="size-7"
            />
          </div>

          <!-- 文字区域 -->
          <div class="flex-1 min-w-0">
            <h3 class="text-lg font-semibold text-highlighted mb-1">
              {{ item.title }}
            </h3>
            <p class="text-sm text-muted leading-relaxed">
              {{ item.description }}
            </p>
          </div>

          <!-- 箭头 -->
          <UIcon
            name="i-lucide-arrow-right"
            class="size-5 text-muted group-hover:text-[var(--ui-primary)] transition-colors flex-shrink-0"
          />
        </NuxtLink>
      </div>
    </UContainer>
  </div>
</template>
