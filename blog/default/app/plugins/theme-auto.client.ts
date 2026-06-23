export default defineNuxtPlugin(() => {
  const { themeConfig } = useTheme();
  const { initAutoSwitch } = useDarkMode();

  const tryInit = () => {
    initAutoSwitch({
      lightStart: themeConfig.value.theme_light_start || '06:00',
      darkStart: themeConfig.value.theme_dark_start || '18:00',
    });
  };

  tryInit();
  watch(themeConfig, tryInit);
});
