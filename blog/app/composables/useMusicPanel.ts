/**
 * 全局音乐面板的显示形态
 * @returns panel - 面板状态
 * @returns hasHover - 设备是否支持悬浮指针
 * @returns setHovered() - 更新停留状态；expand() - 展开完整面板
 */
export function useMusicPanel() {
  const panel = useState('flec-music-panel', () => ({
    hovered: false,
  }));

  const hasHover = useState('flec-music-panel-hover', () => true);

  let leaveTimer: ReturnType<typeof setTimeout> | null = null;

  onMounted(() => {
    hasHover.value = window.matchMedia('(hover: hover)').matches;
  });

  onUnmounted(() => {
    if (leaveTimer) clearTimeout(leaveTimer);
    leaveTimer = null;
  });

  const setHovered = (value: boolean) => {
    if (leaveTimer) clearTimeout(leaveTimer);
    leaveTimer = null;

    if (value) {
      panel.value.hovered = true;
      return;
    }

    leaveTimer = setTimeout(() => {
      leaveTimer = null;
      panel.value.hovered = false;
    }, 150);
  };

  const expand = () => {
    setHovered(true);
  };

  return { panel, hasHover, setHovered, expand };
}
