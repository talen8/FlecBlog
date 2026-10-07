import type {
  FriendApplyRequest,
  FriendCircleItem,
  FriendCircleResponse,
} from '../../types/friend';
import { getFriends, applyFriend, getFriendCircle } from './api/friend';

/**
 * 获取友链分组列表（支持 SSR）
 * @returns data - 友链分组数组，refresh - 刷新方法
 */
export function useFriends() {
  return useAsyncData('friends', async () => {
    const { groups } = await getFriends();
    return groups ?? [];
  });
}

/**
 * 获取友圈文章（SSR 首屏 + 客户端加载更多 + 关键词搜索）
 * @param pageSize - 每页数量，默认 20
 * @returns keyword - 搜索关键词，items - 已加载的文章列表，total - 文章总数，hasMore - 是否还有更多，loading - 加载中状态，loadMore() - 加载下一页
 */
export function useFriendCircle(pageSize = 20) {
  const keyword = ref('');
  const page = ref(1);
  const loading = ref(false);
  const extraItems = ref<FriendCircleItem[]>([]);

  const { data, refresh } = useAsyncData(
    'friend-circle',
    () =>
      getFriendCircle({
        page: 1,
        page_size: pageSize,
        keyword: keyword.value.trim() || undefined,
      }),
    {
      default: (): FriendCircleResponse => ({ list: [], total: 0, page: 1, page_size: pageSize }),
    }
  );

  // 首屏数据与后续追加数据的合并结果
  const items = computed<FriendCircleItem[]>(() => [
    ...(data.value?.list ?? []),
    ...extraItems.value,
  ]);

  const total = computed(() => data.value?.total ?? 0);
  const hasMore = computed(() => items.value.length < total.value);

  /** 加载下一页并追加到列表末尾 */
  const loadMore = async () => {
    if (loading.value || !hasMore.value) return;
    loading.value = true;
    try {
      const nextPage = page.value + 1;
      const res = await getFriendCircle({
        page: nextPage,
        page_size: pageSize,
        keyword: keyword.value.trim() || undefined,
      });
      extraItems.value = [...extraItems.value, ...(res?.list ?? [])];
      page.value = nextPage;
    } finally {
      loading.value = false;
    }
  };

  /** 关键词变化后回到第一页重新搜索 */
  watch(keyword, async () => {
    page.value = 1;
    extraItems.value = [];
    await refresh();
  });

  return { keyword, items, total, hasMore, loading, loadMore };
}

/**
 * 提交友链申请
 * @returns submitting - 提交中状态，apply(data) - 执行申请
 */
export function useFriendApply() {
  const submitting = ref(false);

  const apply = async (data: FriendApplyRequest) => {
    submitting.value = true;
    try {
      await applyFriend(data);
    } finally {
      submitting.value = false;
    }
  };

  return { submitting, apply };
}
