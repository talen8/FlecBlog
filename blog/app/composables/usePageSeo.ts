import { computed, type ComputedRef } from 'vue';
import { useTheme } from './useTheme';

interface PageSeoOptions {
  path: string;
  title: string;
  description: string | (() => string);
}

interface PageSeoResult {
  title: ComputedRef<string>;
  description: ComputedRef<string>;
}

/**
 * 设置页面标题与描述，并同步到 SEO 元信息
 * @param options 页面路径与内置文案
 * @returns title - 生效的页面标题；description - 生效的页面描述，可直接用于页面渲染
 */
export function usePageSeo(options: PageSeoOptions): PageSeoResult {
  const { getPage } = useTheme();

  const override = computed(() => getPage(options.path));
  const fallbackDescription = computed(() =>
    typeof options.description === 'function' ? options.description() : options.description
  );

  const title = computed(() => override.value?.title || options.title);
  const description = computed(() => override.value?.description || fallbackDescription.value);

  useSeoMeta({
    title: () => title.value,
    description: () => description.value,
  });

  return { title, description };
}
