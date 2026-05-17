export const isDark = ref(false);

export interface ThemeAutoSwitchConfig {
  lightStart: string;
  darkStart: string;
}

// 时间字符串 "HH:MM" 转为分钟数
const parseTimeToMinutes = (timeStr: string): number => {
  const parts = timeStr.split(':').map(Number);
  const hours = parts[0] ?? 0;
  const minutes = parts[1] ?? 0;
  return hours * 60 + minutes;
};

const getCurrentMinutes = (): number => {
  const now = new Date();
  return now.getHours() * 60 + now.getMinutes();
};

// 根据当前时间与配置判断是否应使用暗色主题
const shouldBeDarkByTime = (config: ThemeAutoSwitchConfig): boolean => {
  const currentMinutes = getCurrentMinutes();
  const lightStartMinutes = parseTimeToMinutes(config.lightStart);
  const darkStartMinutes = parseTimeToMinutes(config.darkStart);

  // 日间开始 < 夜间开始：在 [lightStart, darkStart) 区间内为日间，其余为夜间
  if (lightStartMinutes < darkStartMinutes) {
    if (currentMinutes >= lightStartMinutes && currentMinutes < darkStartMinutes) {
      return false;
    }
    return true;
    // 夜间开始 <= 日间开始（跨天）：在 [darkStart, lightStart) 区间内为夜间，其余为日间
  } else {
    if (currentMinutes >= darkStartMinutes || currentMinutes < lightStartMinutes) {
      return true;
    }
    return false;
  }
};

// 计算到下一次主题切换的毫秒数
const getMsToNextSwitch = (config: ThemeAutoSwitchConfig): number => {
  const currentMinutes = getCurrentMinutes();
  const lightStartMinutes = parseTimeToMinutes(config.lightStart);
  const darkStartMinutes = parseTimeToMinutes(config.darkStart);

  let nextSwitchMinutes: number;

  if (lightStartMinutes < darkStartMinutes) {
    if (currentMinutes < lightStartMinutes) {
      nextSwitchMinutes = lightStartMinutes;
    } else if (currentMinutes < darkStartMinutes) {
      nextSwitchMinutes = darkStartMinutes;
    } else {
      nextSwitchMinutes = lightStartMinutes + 24 * 60;
    }
  } else {
    if (currentMinutes < darkStartMinutes) {
      nextSwitchMinutes = darkStartMinutes;
    } else if (currentMinutes < lightStartMinutes) {
      nextSwitchMinutes = lightStartMinutes;
    } else {
      nextSwitchMinutes = darkStartMinutes + 24 * 60;
    }
  }

  const minutesUntilSwitch = nextSwitchMinutes - currentMinutes;
  return minutesUntilSwitch * 60 * 1000;
};

let autoSwitchTimer: ReturnType<typeof setTimeout> | null = null;

// 设置 setTimeout 定时器，到达切换点后切换主题并递归调度下一次
const setupAutoSwitchTimer = (config: ThemeAutoSwitchConfig): void => {
  if (autoSwitchTimer) {
    clearTimeout(autoSwitchTimer);
    autoSwitchTimer = null;
  }

  const msToNextSwitch = getMsToNextSwitch(config);
  autoSwitchTimer = setTimeout(() => {
    const shouldBeDark = shouldBeDarkByTime(config);
    if (isDark.value !== shouldBeDark) {
      isDark.value = shouldBeDark;
    }
    setupAutoSwitchTimer(config);
  }, msToNextSwitch);
};

// 由插件在 Nuxt 上下文内调用，启动定时主题切换
export const initAutoSwitch = (config: ThemeAutoSwitchConfig): void => {
  // 两时间相同时视为始终亮色，不启用自动切换
  if (config.lightStart !== config.darkStart) {
    isDark.value = shouldBeDarkByTime(config);
    setupAutoSwitchTimer(config);
  }
};

if (import.meta.client) {
  // 从服务端渲染的 data-theme 属性读取初始状态
  const currentTheme = document.documentElement.getAttribute('data-theme');
  isDark.value = currentTheme === 'dark';

  // 主题变化时同步 DOM 与 localStorage
  watch(isDark, dark => {
    document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
    localStorage.setItem('theme', dark ? 'dark' : 'light');
  });

  // 用户未手动选择主题时跟随系统主题
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
  mediaQuery.addEventListener('change', e => {
    if (!localStorage.getItem('theme')) {
      isDark.value = e.matches;
    }
  });
}

/**
 * 切换暗黑模式
 */
export const toggleTheme = (): void => {
  isDark.value = !isDark.value;
};
