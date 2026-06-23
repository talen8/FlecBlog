/**
 * 展开/折叠状态管理
 * @param initialState - 初始展开状态，默认 false
 * @returns isExpanded - 当前展开状态，toggleExpand - 切换展开/折叠
 */
export function useExpandable(initialState = false) {
  const isExpanded = ref(initialState);
  const toggleExpand = () => (isExpanded.value = !isExpanded.value);
  return { isExpanded, toggleExpand };
}
