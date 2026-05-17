export default defineNuxtPlugin(() => {
  const { blogConfig } = useSysConfig();

  const tryInit = () => {
    initAutoSwitch({
      lightStart: blogConfig.value.theme_light_start || '06:00',
      darkStart: blogConfig.value.theme_dark_start || '18:00',
    });
  };

  tryInit();
  watch(blogConfig, tryInit);
});
