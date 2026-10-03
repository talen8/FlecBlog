/**
 * 主题配色注入插件
 */

const CSS_COLOR_RE = /^(?:#(?:[0-9a-f]{3}|[0-9a-f]{6}|[0-9a-f]{8})|rgba?\([0-9.,%\s/]+\))$/i;

const parseColor = (value: unknown): string | null => {
  if (typeof value !== 'string') return null;
  const trimmed = value.trim();
  return CSS_COLOR_RE.test(trimmed) ? trimmed : null;
};

export default defineNuxtPlugin({
  name: 'theme-color',
  setup() {
    const { themeConfig } = useTheme();

    const cssText = computed(() => {
      const themeColor = parseColor(themeConfig.value.theme_color);
      return themeColor ? `:root{--theme-color-custom:${themeColor};}` : '';
    });

    useHead(
      computed(() => ({
        style: cssText.value ? [{ key: 'theme-color-vars', innerHTML: cssText.value }] : [],
      }))
    );
  },
});
