<script setup lang="ts">
const { data: page } = await useAsyncData('themes', () => queryCollection('themes').first())

const title = page.value?.seo?.title || page.value?.title
const description = page.value?.seo?.description || page.value?.description

useSeoMeta({
  title,
  ogTitle: title,
  description,
  ogDescription: description
})

/** 复制主题名称逻辑 */
const copiedTheme = ref('')

function copyThemeName(name: string) {
  navigator.clipboard.writeText(name)
  copiedTheme.value = name
  setTimeout(() => {
    copiedTheme.value = ''
  }, 2000)
}
</script>

<template>
  <div v-if="page">
    <UContainer>
      <UPageHeader
        v-bind="page"
        class="py-[50px]"
      />
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 pt-8">
        <div
          v-for="theme in page.items"
          :key="theme.name"
          class="group rounded-lg border border-[var(--ui-border)] overflow-hidden bg-[var(--ui-bg)] transition-all hover:shadow-lg hover:border-[var(--ui-primary)]/30"
        >
          <!-- 主题预览图 -->
          <div class="aspect-video overflow-hidden bg-[var(--ui-bg-elevated)]">
            <NuxtImg
              :src="theme.image"
              :alt="theme.display_name"
              class="w-full h-full object-cover"
              loading="lazy"
            />
          </div>

          <!-- 主题信息 -->
          <div class="p-5">
            <div class="flex items-center justify-between mb-2">
              <h3 class="text-base font-semibold text-highlighted">
                {{ theme.display_name }}
              </h3>
            </div>

            <p class="text-sm text-muted leading-relaxed mb-3 line-clamp-2">
              {{ theme.description || '暂无描述' }}
            </p>

            <div class="flex items-center gap-1.5 text-xs text-muted mb-4">
              <UIcon
                name="i-lucide-user"
                class="size-3.5"
              />
              <span>{{ theme.author }}</span>
            </div>

            <!-- 操作区 -->
            <div class="flex flex-wrap gap-2">
              <UButton
                size="xs"
                variant="subtle"
                icon="i-lucide-copy"
                class="cursor-pointer"
                @click="copyThemeName(theme.name)"
              >
                {{ copiedTheme === theme.name ? '已复制' : '复制名称' }}
              </UButton>
              <UButton
                v-if="theme.links?.preview"
                size="xs"
                variant="ghost"
                icon="i-lucide-external-link"
                :to="theme.links.preview"
                target="_blank"
              >
                预览
              </UButton>
              <UButton
                v-if="theme.links?.source"
                size="xs"
                variant="ghost"
                icon="i-simple-icons-github"
                :to="theme.links.source"
                target="_blank"
              >
                源码
              </UButton>
            </div>
          </div>
        </div>
      </div>
    </UContainer>
  </div>
</template>
