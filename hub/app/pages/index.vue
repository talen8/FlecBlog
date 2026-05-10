<script setup lang="ts">
const { data: page } = await useAsyncData('index', () => queryCollection('index').first())

const title = page.value?.seo?.title || page.value?.title
const description = page.value?.seo?.description || page.value?.description

useSeoMeta({
  titleTemplate: '',
  title,
  ogTitle: title,
  description,
  ogDescription: description
})

/** 安装命令复制逻辑 */
const installCommand = computed(() => page.value?.hero?.command || '')
const copied = ref(false)

function copyCommand() {
  navigator.clipboard.writeText(installCommand.value)
  copied.value = true
  setTimeout(() => {
    copied.value = false
  }, 2000)
}
</script>

<template>
  <div v-if="page">
    <section class="relative min-h-screen flex items-center justify-center">
      <HeroBackground />

      <div class="relative text-center px-6 -translate-y-[80px]">
        <span class="mb-6 inline-flex items-center gap-1.5 text-sm font-semibold text-primary">
          {{ page.hero.headline }}
        </span>

        <h1 class="text-5xl sm:text-7xl font-bold tracking-tight text-pretty text-highlighted">
          <MDC
            :value="page.title"
            unwrap="p"
          />
        </h1>

        <p class="mt-8 text-lg sm:text-xl/8 text-muted text-balance">
          {{ page.description }}
        </p>

        <div class="mt-10 flex flex-wrap justify-center gap-x-6 gap-y-3">
          <UButton
            v-for="link in page.hero.links"
            :key="link.label"
            size="xl"
            v-bind="link"
          />
        </div>

        <div
          v-if="installCommand"
          class="mt-8 mx-auto w-fit inline-flex items-center gap-3 rounded-lg bg-slate-800 px-4 py-3 cursor-pointer transition-all hover:bg-slate-700 hover:shadow-lg dark:bg-slate-900 dark:hover:bg-slate-800"
          @click="copyCommand"
        >
          <code class="text-sm font-mono text-slate-200 dark:text-slate-100">{{ installCommand }}</code>
          <UIcon
            :name="copied ? 'i-lucide-check' : 'i-lucide-copy'"
            class="text-slate-400 transition-colors"
            :class="{ 'text-green-400': copied }"
          />
        </div>
      </div>
    </section>

    <UPageSection
      v-for="(section, index) in page.sections"
      :key="index"
      :title="section.title"
      :description="section.description"
      :orientation="section.orientation"
      :reverse="section.reverse"
      :features="section.features"
    >
      <NuxtImg
        :src="section.image"
        :alt="section.title"
        class="w-full rounded-lg border border-[var(--ui-border)] shadow-lg"
        loading="lazy"
      />
    </UPageSection>

    <UPageSection
      :title="page.features.title"
      :description="page.features.description"
    >
      <UPageGrid>
        <UPageCard
          v-for="(item, index) in page.features.items"
          :key="index"
          v-bind="item"
          spotlight
        />
      </UPageGrid>
    </UPageSection>

    <UPageSection
      id="testimonials"
      :headline="page.testimonials.headline"
      :title="page.testimonials.title"
      :description="page.testimonials.description"
    >
      <UPageColumns class="xl:columns-4">
        <UPageCard
          v-for="(testimonial, index) in page.testimonials.items"
          :key="index"
          variant="subtle"
          :description="testimonial.quote"
          :ui="{ description: 'before:content-[open-quote] after:content-[close-quote]' }"
        >
          <template #footer>
            <UUser
              v-bind="testimonial.user"
              size="lg"
            />
          </template>
        </UPageCard>
      </UPageColumns>
    </UPageSection>

    <USeparator />

    <UPageCTA
      v-bind="page.cta"
      variant="naked"
      class="overflow-hidden"
    >
      <LazyStarsBg />
    </UPageCTA>
  </div>
</template>
