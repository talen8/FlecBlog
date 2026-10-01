import type { Category } from '../../types/category';
import { getCategories, getCategoryBySlug } from './api/category';

/**
 * 按 slug 获取分类详情（支持 SSR）
 * @param slug - 分类 slug，支持 ref / computed / getter
 * @returns data - 分类对象，refresh - 刷新方法
 */
export function useCategory(slug: MaybeRefOrGetter<string>) {
  return useAsyncData(`category:${toValue(slug)}`, () => getCategoryBySlug(toValue(slug)), {
    watch: [() => toValue(slug)],
  });
}

/**
 * 获取全部分类列表（全局缓存，仅首次请求）
 * @returns categories - 分类数组
 */
export function useCategories() {
  const categories = useState<Category[]>('categories', () => []);

  useAsyncData('categories-fetch', async () => {
    if (categories.value.length > 0) return categories.value;
    const { list } = await getCategories();
    categories.value = list ?? [];
    return categories.value;
  });

  return { categories };
}
