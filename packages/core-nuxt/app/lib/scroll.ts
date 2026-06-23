interface ScrollOptions {
  behavior?: ScrollBehavior;
  block?: ScrollLogicalPosition;
}

/**
 * 平滑滚动到顶部
 */
export function scrollToTop(): void {
  if (!import.meta.client) return;

  window.scrollTo({ top: 0, behavior: 'smooth' });
}

/**
 * 平滑滚动到指定元素
 * @param selector CSS 选择器
 * @param options 滚动选项
 */
export function scrollToElement(selector: string, options?: ScrollOptions): void {
  if (!import.meta.client) return;

  const { behavior = 'smooth', block = 'center' } = options || {};
  // # 开头用 getElementById 支持特殊字符，否则用 querySelector
  const element = selector.startsWith('#')
    ? document.getElementById(selector.slice(1))
    : document.querySelector(selector);

  if (element) {
    element.scrollIntoView({ behavior, block });
  }
}
