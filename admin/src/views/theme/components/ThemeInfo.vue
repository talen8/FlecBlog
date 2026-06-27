<template>
  <div class="theme-info">
    <el-descriptions :column="2" border>
      <el-descriptions-item label="主题名称">
        {{ theme.name || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="主题标识">
        {{ theme.slug }}
      </el-descriptions-item>
      <el-descriptions-item label="版本">
        <span class="current-version">{{ theme.version || '-' }}</span>
        <span v-if="checkResult?.has_update" class="update-info">
          <el-tag type="danger" size="small">{{ checkResult.latest_version }}</el-tag>
        </span>
      </el-descriptions-item>
      <el-descriptions-item label="作者">
        {{ theme.author || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="仓库">
        <el-link v-if="theme.repo" :href="theme.repo" target="_blank" type="primary">
          {{ theme.repo }}
        </el-link>
        <span v-else>-</span>
      </el-descriptions-item>
      <el-descriptions-item label="许可证">
        {{ theme.license || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="描述" :span="2">
        {{ theme.description || '-' }}
      </el-descriptions-item>
    </el-descriptions>

    <div class="features-section">
      <h3>功能支持</h3>
      <div class="features-grid">
        <div
          v-for="feature in featureList"
          :key="feature.key"
          class="feature-item"
          :class="{ disabled: !feature.enabled }"
        >
          <i :class="feature.icon"></i>
          <div class="feature-info">
            <span>{{ feature.label }}</span>
            <el-tag :type="feature.enabled ? 'success' : 'info'" size="small">
              {{ feature.enabled ? '已启用' : '未启用' }}
            </el-tag>
          </div>
        </div>
      </div>
    </div>

    <div class="schema-section">
      <h3>Schema</h3>
      <pre class="schema-viewer">{{ schemaText }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue';
import { checkThemeUpdate } from '@/api/theme';
import type { ThemeResponse, ThemeUpdateCheckResponse } from '@/types/theme';

const props = defineProps<{
  theme: ThemeResponse;
}>();

const checkResult = ref<ThemeUpdateCheckResponse | null>(null);

onMounted(async () => {
  try {
    const raw = localStorage.getItem(`tcc_${props.theme.slug}`);
    if (raw) {
      const { v, t } = JSON.parse(raw);
      if (Date.now() < t) {
        checkResult.value = {
          has_update: true,
          latest_version: v,
          current_version: '',
          release_url: '',
        };
        return;
      }
      localStorage.removeItem(`tcc_${props.theme.slug}`);
    }
  } catch {}

  try {
    const result = await checkThemeUpdate(props.theme.slug);
    if (result.has_update) {
      localStorage.setItem(
        `tcc_${props.theme.slug}`,
        JSON.stringify({ v: result.latest_version, t: Date.now() + 300_000 })
      );
      checkResult.value = result;
    }
  } catch (_error) {}
});

const featureLabels: Record<string, { label: string; icon: string }> = {
  moments: { label: '动态', icon: 'ri-chat-3-line' },
  feedback: { label: '反馈投诉', icon: 'ri-feedback-line' },
  oauth: { label: 'OAuth 配置', icon: 'ri-shield-keyhole-line' },
  site_subscribe: { label: '本站订阅', icon: 'ri-mail-send-line' },
};

const schemaText = computed(() => JSON.stringify(props.theme.schema || {}, null, 2));

const featureList = computed(() => {
  const raw = (props.theme.schema || {}).$features;
  const features =
    raw && typeof raw === 'object' && !Array.isArray(raw) ? (raw as Record<string, unknown>) : {};

  return Object.entries(featureLabels).map(([key, meta]) => ({
    key,
    label: meta.label,
    icon: meta.icon,
    enabled: features[key] !== false && features[key] !== undefined,
  }));
});
</script>

<style scoped lang="scss">
.features-section {
  margin-top: 20px;

  h3 {
    margin: 0 0 12px;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
  }
}

.schema-section {
  margin-top: 20px;

  h3 {
    margin: 0 0 12px;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
  }
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 12px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  background: #fff;

  &.disabled {
    opacity: 0.55;
  }

  > i {
    width: 32px;
    height: 32px;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    border-radius: 6px;
    background: #f5f7fa;
    color: #606266;
    font-size: 16px;
    flex-shrink: 0;
  }
}

.feature-info {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;

  span {
    color: #303133;
    font-size: 14px;
  }
}

.schema-viewer {
  font-family: Consolas, 'Courier New', monospace;
  background: #f5f7fa;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 12px;
  font-size: 13px;
  line-height: 1.5;
  margin: 0;
}

:deep(.el-descriptions__body .el-descriptions__table) {
  table-layout: fixed;
}

:deep(.el-descriptions__body .el-descriptions__label) {
  width: 200px;
}

.current-version {
  margin-right: 8px;
}

.update-info {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
</style>
