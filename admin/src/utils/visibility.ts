import { reactive } from 'vue';
import { getTheme, getThemes } from '@/api/theme';

const state = reactive({
  loaded: false,
  loading: null as Promise<void> | null,
  features: {} as Record<string, boolean>,
});

const readFeatureMap = (schema: Record<string, unknown> | undefined) => {
  const raw = schema?.$features;
  const next: Record<string, boolean> = {};

  if (raw && typeof raw === 'object' && !Array.isArray(raw)) {
    Object.entries(raw as Record<string, unknown>).forEach(([key, value]) => {
      next[key] = value !== false;
    });
  }

  return next;
};

export const loadThemeFeatures = async (force = false) => {
  if (state.loading) return state.loading;
  if (state.loaded && !force) return;

  state.loading = (async () => {
    try {
      const themes = await getThemes();
      const activeTheme = themes.find(theme => theme.is_active) || themes[0];
      if (!activeTheme) {
        state.features = {};
        return;
      }

      const theme = await getTheme(activeTheme.slug);
      state.features = readFeatureMap(theme.schema);
    } finally {
      state.loaded = true;
      state.loading = null;
    }
  })();

  return state.loading;
};

export const useThemeFeatures = () => {
  if (!state.loaded && !state.loading) {
    void loadThemeFeatures();
  }

  return {
    isFeatureEnabled: (key: string) => state.features[key] ?? true,
    loadThemeFeatures,
  };
};
