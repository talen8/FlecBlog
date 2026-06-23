const isDark = ref(false);

/** 暗色模式自动切换时间配置 */
export interface DarkModeAutoSwitchConfig {
  lightStart: string;
  darkStart: string;
}

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

const shouldBeDarkByTime = (config: DarkModeAutoSwitchConfig): boolean => {
  const currentMinutes = getCurrentMinutes();
  const lightStartMinutes = parseTimeToMinutes(config.lightStart);
  const darkStartMinutes = parseTimeToMinutes(config.darkStart);

  if (lightStartMinutes < darkStartMinutes) {
    if (currentMinutes >= lightStartMinutes && currentMinutes < darkStartMinutes) {
      return false;
    }
    return true;
  } else {
    if (currentMinutes >= darkStartMinutes || currentMinutes < lightStartMinutes) {
      return true;
    }
    return false;
  }
};

const getMsToNextSwitch = (config: DarkModeAutoSwitchConfig): number => {
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

const setupAutoSwitchTimer = (config: DarkModeAutoSwitchConfig): void => {
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

const initAutoSwitch = (config: DarkModeAutoSwitchConfig): void => {
  if (config.lightStart !== config.darkStart) {
    isDark.value = shouldBeDarkByTime(config);
    setupAutoSwitchTimer(config);
  }
};

if (import.meta.client) {
  const currentTheme = document.documentElement.getAttribute('data-theme');
  isDark.value = currentTheme === 'dark';

  watch(isDark, dark => {
    document.documentElement.setAttribute('data-theme', dark ? 'dark' : 'light');
    localStorage.setItem('theme', dark ? 'dark' : 'light');
  });

  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
  mediaQuery.addEventListener('change', e => {
    if (!localStorage.getItem('theme')) {
      isDark.value = e.matches;
    }
  });
}

const toggleTheme = (): void => {
  isDark.value = !isDark.value;
};

/**
 * 暗色模式控制（与主题系统 Theme 区分）
 * - 读取 / 修改 data-theme 属性 + localStorage
 * - 支持跟随系统偏好和定时自动切换
 * @returns isDark - 当前是否为暗色模式
 * @returns toggleTheme - 手动切换
 * @returns initAutoSwitch(config) - 初始化定时自动切换
 */
export function useDarkMode() {
  return { isDark, toggleTheme, initAutoSwitch };
}
