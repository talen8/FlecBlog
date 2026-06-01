<script setup lang="ts">
import type { ContentNavigationItem } from '@nuxt/content'

const navigation = inject<Ref<ContentNavigationItem[]>>('navigation')
const isMobileNavOpen = ref(false)
</script>

<template>
  <div>
    <AppHeader />

    <UMain>
      <UContainer>
        <UPage>
          <template #left>
            <UPageAside>
              <template #top>
                <UContentSearchButton
                  :collapsed="false"
                  label="搜索"
                />
              </template>

              <UContentNavigation
                :navigation="navigation"
                highlight
              />
            </UPageAside>
          </template>

          <slot />
        </UPage>
      </UContainer>
    </UMain>

    <UButton
      icon="i-lucide-menu"
      color="primary"
      size="xl"
      class="lg:hidden fixed bottom-6 right-6 z-40 rounded-full shadow-lg"
      aria-label="文档目录"
      @click="isMobileNavOpen = true"
    />

    <USlideover
      v-model:open="isMobileNavOpen"
      side="left"
      title="文档导航"
      :ui="{ content: 'max-w-72 sm:max-w-sm' }"
    >
      <template #body>
        <UContentSearchButton
          :collapsed="false"
          label="搜索"
          class="mb-4 w-full"
        />
        <UContentNavigation
          :navigation="navigation"
          highlight
          @click="isMobileNavOpen = false"
        />
      </template>
    </USlideover>

    <AppFooter />
  </div>
</template>
