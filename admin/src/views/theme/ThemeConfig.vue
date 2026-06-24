<template>
  <div class="theme-config">
    <el-card shadow="never">
      <div class="toolbar">
        <h2>主题配置</h2>
        <div class="actions">
          <el-button
            type="primary"
            :loading="saving"
            :disabled="loading || !canEditTheme || !currentTheme || !hasSchema"
            @click="handleSave"
          >
            保存配置
          </el-button>
          <el-button :disabled="loading || !currentTheme" @click="() => loadCurrentTheme()">
            重置
          </el-button>
        </div>
      </div>

      <el-alert
        v-if="!canEditTheme"
        title="仅超级管理员可修改主题配置"
        type="warning"
        show-icon
        :closable="false"
        class="permission-alert"
      />

      <el-skeleton v-if="loading && !currentTheme" :rows="10" animated />

      <template v-else-if="currentTheme">
        <el-tabs v-model="activeTab" class="theme-tabs">
          <el-tab-pane label="主题信息" name="info">
            <ThemeInfo :theme="currentTheme" />
          </el-tab-pane>

          <template v-if="hasSchema">
            <el-tab-pane
              v-for="group in schemaGroups"
              :key="group.name"
              :label="group.label || '配置项'"
              :name="groupTabName(group)"
            >
              <el-form label-position="top" :disabled="formDisabled">
                <el-form-item v-for="(field, key) in group.fields" :key="key">
                  <template #label>
                    <span :class="{ 'field-modified': isFieldModified(String(key)) }">{{
                      field.title || String(key)
                    }}</span>
                    <div v-if="field.description" class="field-desc">
                      {{ field.description }}
                    </div>
                  </template>

                  <el-switch v-if="field.type === 'boolean'" v-model="configValues[String(key)]" />

                  <el-input-number
                    v-else-if="field.type === 'number' || field.type === 'integer'"
                    v-model="configValues[String(key)]"
                    :min="field.min"
                    :max="field.max"
                    :step="field.type === 'integer' ? 1 : 0.1"
                  />

                  <el-select
                    v-else-if="field.type === 'string' && field.enum"
                    v-model="configValues[String(key)]"
                    :placeholder="field.placeholder || '请选择'"
                    clearable
                    filterable
                  >
                    <el-option
                      v-for="opt in field.enum"
                      :key="typeof opt === 'object' ? opt.value : String(opt)"
                      :label="typeof opt === 'object' ? opt.label : String(opt)"
                      :value="typeof opt === 'object' ? opt.value : opt"
                    />
                  </el-select>

                  <el-color-picker
                    v-else-if="field.format === 'color'"
                    v-model="configValues[String(key)]"
                  />

                  <el-date-picker
                    v-else-if="field.format === 'date'"
                    v-model="configValues[String(key)]"
                    type="date"
                    value-format="YYYY-MM-DD"
                    placeholder="选择日期"
                  />

                  <el-time-picker
                    v-else-if="field.format === 'time'"
                    v-model="configValues[String(key)]"
                    value-format="HH:mm:ss"
                    placeholder="选择时间"
                  />

                  <el-date-picker
                    v-else-if="field.format === 'date-time'"
                    v-model="configValues[String(key)]"
                    type="datetime"
                    value-format="YYYY-MM-DDTHH:mm:ss"
                    placeholder="选择日期时间"
                  />

                  <ImageUploader
                    v-else-if="field.format === 'image'"
                    :ref="(el: unknown) => setImageUploaderRef(String(key), el)"
                    v-model="configValues[String(key)]"
                    upload-type="主题图片"
                    :width="(field.width || 120) + 'px'"
                    :height="(field.height || 120) + 'px'"
                    :disabled="formDisabled"
                  />

                  <el-input
                    v-else-if="field.format === 'upload'"
                    v-model="configValues[String(key)]"
                    :placeholder="field.placeholder || field.description || '文件URL'"
                  >
                    <template #append>
                      <el-upload
                        :show-file-list="false"
                        :http-request="
                          (opts: UploadRequestOptions) => handleSimpleUpload(String(key), opts)
                        "
                        accept="image/*"
                        :disabled="formDisabled"
                      >
                        <el-button
                          :icon="Upload"
                          :type="pendingUploadFiles[String(key)] ? 'success' : 'default'"
                        />
                      </el-upload>
                    </template>
                  </el-input>

                  <JsonListEditor
                    v-else-if="field.type === 'array' && field['x-item-fields']"
                    v-model="configValues[String(key)]"
                    :fields="buildItemFields(field['x-item-fields'])"
                    :default-item="buildDefaultItem(field['x-item-fields'])"
                    :disabled="formDisabled"
                  />

                  <el-input
                    v-else-if="field.format === 'textarea'"
                    v-model="configValues[String(key)]"
                    type="textarea"
                    :rows="4"
                    :placeholder="field.placeholder || ''"
                  />

                  <el-input
                    v-else
                    v-model="configValues[String(key)]"
                    :placeholder="field.placeholder || ''"
                    clearable
                  />
                </el-form-item>
              </el-form>
            </el-tab-pane>
          </template>

          <el-tab-pane v-if="!hasSchema" label="配置项" name="config-empty">
            <el-empty description="当前主题未提供配置 schema" />
          </el-tab-pane>

          <el-tab-pane label="主题菜单" name="menus">
            <ThemeMenu
              :theme-slug="currentTheme.slug"
              :schema="currentTheme.schema"
              :menus="currentTheme.menus"
              :disabled="formDisabled"
              @refresh="loadCurrentTheme"
            />
          </el-tab-pane>
        </el-tabs>
      </template>

      <el-empty v-else description="暂无可配置主题" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { ElMessage, type UploadRequestOptions } from 'element-plus';
import { Upload } from '@element-plus/icons-vue';
import JsonListEditor from '@/components/common/JsonListEditor.vue';
import type { FieldConfig } from '@/components/common/JsonListEditor.vue';
import ImageUploader from '@/components/common/ImageUploader.vue';
import ThemeInfo from './components/ThemeInfo.vue';
import ThemeMenu from './components/ThemeMenu.vue';
import { getTheme, getThemes, updateThemeConfig } from '@/api/theme';
import { uploadFile } from '@/api/file';
import type { ThemeResponse, SchemaField, SchemaGroup } from '@/types/theme';
import { isSuperAdmin } from '@/utils/auth';

const route = useRoute();
const activeTab = ref(route.query.tab === 'menus' ? 'menus' : 'info');
const loading = ref(false);
const saving = ref(false);
const currentTheme = ref<ThemeResponse | null>(null);
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const configValues = ref<Record<string, any>>({});
const pendingUploadFiles = ref<Record<string, File>>({});
const pendingPreviews = ref<Record<string, string>>({});
interface ImageUploaderExposed {
  uploadPendingFile: () => Promise<string | null>;
  getPendingCount: () => number;
}
const imageUploaderRefs = ref<Record<string, ImageUploaderExposed>>({});

const setImageUploaderRef = (key: string, el: unknown) => {
  if (el) {
    imageUploaderRefs.value[key] = el as ImageUploaderExposed;
  } else {
    delete imageUploaderRefs.value[key];
  }
};

const canEditTheme = computed(() => isSuperAdmin());
const formDisabled = computed(() => loading.value || saving.value || !canEditTheme.value);

const originalConfigValues = ref<Record<string, unknown>>({});

const schema = computed<Record<string, unknown>>(() => {
  const raw = (currentTheme.value?.schema || {}) as Record<string, unknown>;
  const nested = raw.config as Record<string, unknown> | undefined;
  if (nested && typeof nested === 'object') return nested;
  return raw;
});

const hasSchema = computed(() => {
  const s = schema.value;
  return Object.keys(s).some(k => !k.startsWith('$'));
});

const schemaGroups = computed<SchemaGroup[]>(() => {
  const s = schema.value;
  const entries = Object.entries(s).filter(([k]) => !k.startsWith('$'));
  if (entries.length === 0) return [];

  const first = entries[0]![1];
  if (first && typeof first === 'object' && 'type' in (first as Record<string, unknown>)) {
    return [{ name: '_default', label: '配置项', fields: s as Record<string, SchemaField> }];
  }

  return entries
    .filter(([, v]) => v && typeof v === 'object' && !('type' in (v as Record<string, unknown>)))
    .map(([name, val]) => ({
      name,
      label: name,
      fields: val as Record<string, SchemaField>,
    }));
});

const groupTabName = (group: SchemaGroup) => `config:${group.name}`;

const isFieldModified = (key: string): boolean => {
  if (loading.value) return false;
  const current = configValues.value[key];
  const original = originalConfigValues.value[key];
  if (current === original) return false;
  if (current === undefined && original === undefined) return false;
  if (current === null && original === null) return false;
  return JSON.stringify(current) !== JSON.stringify(original);
};

const buildItemFields = (itemFields: Array<string | Record<string, unknown>>): FieldConfig[] => {
  if (!itemFields || itemFields.length === 0) return [];
  if (typeof itemFields[0] === 'string') {
    return (itemFields as string[]).map(key => ({
      key,
      type: 'text' as const,
      placeholder: key,
    }));
  }
  return itemFields as unknown as FieldConfig[];
};

const buildDefaultItem = (
  itemFields: Array<string | Record<string, unknown>>
): Record<string, unknown> => {
  const obj: Record<string, unknown> = {};
  if (!itemFields || itemFields.length === 0) return obj;
  if (typeof itemFields[0] === 'string') {
    (itemFields as string[]).forEach(k => {
      obj[k] = '';
    });
  } else {
    (itemFields as Array<Record<string, unknown>>).forEach(f => {
      obj[f.key as string] = f.default ?? '';
    });
  }
  return obj;
};

const collectDefaults = (obj: Record<string, unknown>, out: Record<string, unknown>) => {
  for (const [k, v] of Object.entries(obj)) {
    if (k.startsWith('$')) continue;
    if (v && typeof v === 'object' && 'type' in (v as Record<string, unknown>)) {
      const field = v as SchemaField;
      const key = String(k);
      if (out[key] === undefined && field.default !== undefined) {
        out[key] = field.default;
      }
    } else if (v && typeof v === 'object') {
      collectDefaults(v as Record<string, unknown>, out);
    }
  }
};

const handleSimpleUpload = (key: string, opts: UploadRequestOptions): Promise<void> => {
  const file = opts.file as File;
  if (!file.type.startsWith('image/')) {
    ElMessage.error('请选择图片文件');
    return Promise.resolve();
  }
  if (pendingPreviews.value[key]) {
    URL.revokeObjectURL(pendingPreviews.value[key]);
  }
  const blobUrl = URL.createObjectURL(file);
  pendingPreviews.value[key] = blobUrl;
  pendingUploadFiles.value[key] = file;
  configValues.value[key] = blobUrl;
  return Promise.resolve();
};

const loadCurrentTheme = async (slug = currentTheme.value?.slug) => {
  if (!slug) return;
  loading.value = true;
  try {
    const theme = await getTheme(slug);
    currentTheme.value = theme;
    const rawConfig = (theme.config || {}) as Record<string, unknown>;

    const defaults: Record<string, unknown> = {};
    const s = schema.value;
    collectDefaults(s, defaults);

    for (const url of Object.values(pendingPreviews.value)) {
      URL.revokeObjectURL(url);
    }
    pendingUploadFiles.value = {};
    pendingPreviews.value = {};
    const mergedConfig = { ...defaults, ...rawConfig };
    configValues.value = mergedConfig;
    originalConfigValues.value = { ...mergedConfig };
  } catch (error) {
    ElMessage.error((error as Error)?.message || '获取主题配置失败');
  } finally {
    loading.value = false;
  }
};

const handleSave = async () => {
  if (!canEditTheme.value) {
    ElMessage.warning('仅超级管理员可修改主题配置');
    return;
  }
  if (!currentTheme.value || !hasSchema.value) return;

  saving.value = true;
  try {
    for (const [key, file] of Object.entries(pendingUploadFiles.value)) {
      try {
        const result = await uploadFile(file, '主题图片');
        configValues.value[key] = result.file_url;
      } catch (_e) {
        ElMessage.error(`「${file.name}」上传失败`);
        return;
      }
    }
    for (const url of Object.values(pendingPreviews.value)) {
      URL.revokeObjectURL(url);
    }
    pendingUploadFiles.value = {};
    pendingPreviews.value = {};

    const uploadPromises: Promise<string | null>[] = [];
    for (const uploader of Object.values(imageUploaderRefs.value)) {
      if (uploader?.getPendingCount()) {
        uploadPromises.push(uploader.uploadPendingFile());
      }
    }
    if (uploadPromises.length > 0) {
      const results = await Promise.allSettled(uploadPromises);
      const failed = results.filter(r => r.status === 'rejected');
      if (failed.length > 0) {
        ElMessage.error(`${failed.length} 个图片上传失败，请重试`);
        return;
      }
    }

    const updatedConfig = await updateThemeConfig(currentTheme.value.slug, {
      ...configValues.value,
    });
    currentTheme.value = {
      ...currentTheme.value,
      config: updatedConfig as Record<string, unknown>,
    };
    configValues.value = { ...(updatedConfig as Record<string, unknown>) };
    originalConfigValues.value = { ...(updatedConfig as Record<string, unknown>) };
    ElMessage.success('主题配置保存成功');
  } catch (error) {
    ElMessage.error((error as Error)?.message || '主题配置保存失败');
  } finally {
    saving.value = false;
  }
};

watch(
  () => route.query.tab,
  tab => {
    if (tab === 'config') {
      activeTab.value = schemaGroups.value[0]
        ? groupTabName(schemaGroups.value[0])
        : 'config-empty';
      return;
    }
    if (tab === 'info' || tab === 'menus') {
      activeTab.value = tab;
    }
  }
);

watch(schemaGroups, groups => {
  if (activeTab.value === 'config' || activeTab.value === 'config-empty') {
    activeTab.value = groups[0] ? groupTabName(groups[0]) : 'config-empty';
  }
});

onMounted(async () => {
  loading.value = true;
  try {
    const list = await getThemes();
    const activeTheme = list.find(theme => theme.is_active) || list[0];
    if (activeTheme) {
      await loadCurrentTheme(activeTheme.slug);
    }
  } catch (error) {
    ElMessage.error((error as Error)?.message || '获取主题列表失败');
  } finally {
    loading.value = false;
  }
});
</script>

<style lang="scss" scoped>
.theme-config {
  height: 100%;

  :deep(.el-card) {
    height: 100%;
    display: flex;
    flex-direction: column;

    .el-card__body {
      flex: 1;
      display: flex;
      flex-direction: column;
      overflow: hidden;
    }
  }

  .toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;

    h2 {
      margin: 0;
      font-size: 20px;
      font-weight: 500;
    }

    .actions {
      display: flex;
      gap: 12px;

      .el-button {
        margin: 0;
      }
    }
  }

  .permission-alert {
    margin-bottom: 16px;
  }

  .theme-tabs {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;

    :deep(.el-tabs__content) {
      flex: 1;
      overflow: hidden;

      .el-tab-pane {
        height: 100%;
        overflow-y: auto;

        .el-form,
        .theme-info,
        .theme-menu-panel {
          max-width: 95%;
          margin: 0 auto;
        }
      }
    }
  }

  .field-desc {
    font-size: 12px;
    color: #909399;
    font-weight: normal;
    line-height: 1.4;
    margin-top: 2px;
  }

  .field-modified {
    color: #e6a23c;
    font-weight: 600;
  }

  :deep(.el-form-item__content) {
    align-items: flex-start;
  }
}

@media (max-width: 768px) {
  .theme-config {
    .toolbar {
      flex-direction: column;

      .actions {
        width: 100%;
        justify-content: flex-start;
      }
    }

    :deep(.el-form-item__label) {
      width: 100px !important;
      font-size: 13px;
    }
  }
}
</style>
