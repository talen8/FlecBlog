<script setup lang="ts">
import { h } from 'vue'

const IconLogo = () => h(
  'svg',
  {
    'xmlns': 'http://www.w3.org/2000/svg',
    'xml:space': 'preserve',
    'viewBox': '0 0 512 512',
    'width': '1em',
    'height': '1em'
  },
  [
    h('path', {
      fill: 'currentColor',
      d: 'M504 196c-10-7-29-44-34-48-4-4-28-14-28-14l-14-43c-27 3-43 45-44 48 0 0-2 74-120 45-85-20-135-40-178 0q-4-6-6-14-3-14-3-31v-19q1-20-4-44-4-24-22-49a16 16 0 1 0-26 19q15 19 17 37 4 18 4 37l-1 19q0 18 4 39 4 16 16 32c-33 51-15 107-15 107l-34 54-1 3-15 80q0 6 4 9l41 28h3l2-2 9-17-1-4-11-12q-3-3-2-9l18-61 3-4 24-17a9 9 0 0 1 14 10l-11 27q-1 5 2 9l60 73h45q4 0 4-3v-23q0-3-4-3h-14q-5 0-8-5l-24-41q-3-5 0-10l22-29 7-4h62a9 9 0 0 1 8 12l-19 62v6l12 33h40v-9l-7-20v-6l27-72q2-5 8-6h27q4 0 7 3l75 104h51c12 0 7-26-5-27-23-2-52-66-64-88-13-22 24-91 53-94q41-6 47-31c9-10 10-31-1-37'
    })
  ]
)

const columns = [
  {
    label: '产品',
    children: [
      { label: '文档', to: '/docs' },
      { label: '公告', to: '/changelog' },
      { label: '方案', to: '/pricing' }
    ]
  },
  {
    label: '资源',
    children: [
      { label: 'GitHub', to: 'https://github.com/talen8/FlecBlog', target: '_blank' },
      { label: '主题', to: '/themes' }
    ]
  },
  {
    label: '社区',
    children: [
      { label: '社群', to: 'https://qm.qq.com/q/Zzm9XN6lOi', target: '_blank' },
      { label: '作者博客', to: 'https://blog.talen.top/', target: '_blank' }
    ]
  }
]
</script>

<template>
  <USeparator
    :icon="IconLogo"
    class="h-px"
  />

  <UFooter :ui="{ top: 'border-b border-default' }">
    <template #top>
      <UContainer>
        <UFooterColumns :columns="columns">
          <template #right>
            <hgroup class="h-full flex flex-col justify-center items-center text-center">
              <h3 class="text-lg font-semibold">
                FlecBlog
              </h3>
              <p class="text-sm text-muted">
                一款有边界感的博客系统
              </p>
            </hgroup>
          </template>
        </UFooterColumns>
      </UContainer>
    </template>

    <template #left>
      <p class="text-muted text-sm">
        基于 Nuxt UI 构建 • © {{ new Date().getFullYear() }} FlecBlog
      </p>
    </template>

    <template #right>
      <UButton
        to="https://github.com/talen8/FlecBlog"
        target="_blank"
        icon="i-simple-icons-github"
        aria-label="FlecBlog on GitHub"
        color="neutral"
        variant="ghost"
      />
    </template>
  </UFooter>
</template>
