<script setup lang="ts">
const route = useRoute()

const items = computed(() => [{
  label: '文档',
  to: '/docs',
  active: route.path.startsWith('/docs')
}, {
  label: '主题',
  to: '/themes'
}, {
  label: '公告',
  to: '/changelog'
}, {
  label: '方案',
  to: '/pricing'
}, {
  label: '社区',
  to: '/community'
}])
</script>

<template>
  <UHeader>
    <template #left>
      <NuxtLink
        to="/"
        class="text-lg font-bold"
      >
        FlecBlog
      </NuxtLink>
    </template>

    <UNavigationMenu
      :items="items"
      variant="link"
    />

    <template #right>
      <UContentSearchButton />
      <UColorModeButton />
      <UButton
        icon="i-simple-icons-github"
        color="neutral"
        variant="ghost"
        to="https://github.com/talen8/FlecBlog"
        target="_blank"
        aria-label="GitHub"
      />
    </template>

    <template #body>
      <UNavigationMenu
        :items="items"
        orientation="vertical"
        class="-mx-2.5"
      />
    </template>
  </UHeader>
</template>
