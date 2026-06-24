import type { EmojiGroup } from '../../types/comment';

// 全局缓存
let groupsCache: EmojiGroup[] | null = null;
let groupsPromise: Promise<EmojiGroup[]> | null = null;
let emojiMapCache: Map<string, string> | null = null;

/**
 * 加载表情分组数据（带缓存）。
 * 选择器展示与 `:key:` 映射构建共用此唯一入口，避免重复请求。
 */
export async function loadEmojiGroups(emojisUrl: string): Promise<EmojiGroup[]> {
  if (groupsCache) return groupsCache;
  if (groupsPromise) return groupsPromise;

  groupsPromise = (async () => {
    try {
      const response = await fetch(emojisUrl);
      if (!response.ok) throw new Error('加载表情包失败');
      groupsCache = (await response.json()) as EmojiGroup[];
      return groupsCache;
    } catch (error) {
      console.error('加载表情包失败:', error);
      groupsPromise = null;
      return [];
    }
  })();

  return groupsPromise;
}

/**
 * 加载 `:key:` → url 映射（仅 image 类型），用于评论渲染。
 */
export async function loadEmojiMap(emojisUrl: string): Promise<Map<string, string>> {
  if (emojiMapCache) return emojiMapCache;

  const groups = await loadEmojiGroups(emojisUrl);
  const map = new Map<string, string>();
  for (const group of groups) {
    if (group.type === 'image') {
      for (const item of group.items) {
        map.set(item.key, item.val);
      }
    }
  }
  emojiMapCache = map;
  return map;
}

/**
 * 获取缓存的表情映射（同步）
 */
export function getEmojiMapSync(): Map<string, string> | null {
  return emojiMapCache;
}

/**
 * 清除表情缓存
 */
export function clearEmojiCache(): void {
  emojiMapCache = null;
  groupsCache = null;
  groupsPromise = null;
}

/**
 * 替换文本中的表情占位符为 img 标签
 */
export function replaceEmojisInText(text: string, emojiMap: Map<string, string>): string {
  return text.replace(/:([^:\s]+):/g, (match, key) => {
    const url = emojiMap.get(key);
    if (url) {
      return `<img src="${url}" alt="${key}" class="emoji-image" title="${key}" />`;
    }
    return match;
  });
}
