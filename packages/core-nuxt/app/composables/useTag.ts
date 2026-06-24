import type { Tag } from '../../types/tag';
import { getTags, getTagBySlug } from './api/tag';

/**
 * 按 slug 获取标签详情（支持 SSR）
 * @param slug - 标签 slug，支持 ref / computed / getter
 * @returns data - 标签对象，refresh - 刷新方法
 */
export function useTag(slug: MaybeRefOrGetter<string>) {
  return useAsyncData(`tag:${toValue(slug)}`, () => getTagBySlug(toValue(slug)), {
    watch: [() => toValue(slug)],
  });
}

/**
 * 获取全部标签列表（全局缓存，仅首次请求）
 * @returns tags - 标签数组
 */
export function useTags() {
  const tags = useState<Tag[]>('tags', () => []);

  useAsyncData('tags-fetch', async () => {
    if (tags.value.length > 0) return tags.value;
    const { list } = await getTags();
    tags.value = list ?? [];
    return tags.value;
  });

  return { tags };
}
