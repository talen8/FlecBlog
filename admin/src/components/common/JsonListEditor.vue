<template>
  <div class="json-list-editor">
    <div v-for="(item, index) in internalValue" :key="index" class="editor-item">
      <!-- 排序按钮 -->
      <template v-if="!hideControls">
        <el-button
          :icon="ArrowUp"
          circle
          size="small"
          @click="moveUp(index)"
          :disabled="disabled || index === 0"
        />
        <el-button
          :icon="ArrowDown"
          circle
          size="small"
          @click="moveDown(index)"
          :disabled="disabled || index === internalValue.length - 1"
          style="margin-left: 0"
        />
      </template>

      <!-- 动态字段 -->
      <template v-for="field in fields" :key="field.key">
        <!-- 文本输入 -->
        <el-input
          v-if="field.type === 'text'"
          v-model="item[field.key]"
          :placeholder="field.placeholder"
          :style="field.style"
          :disabled="disabled"
          @input="emitUpdate"
        />

        <!-- 下拉选择 -->
        <el-select
          v-else-if="field.type === 'select'"
          v-model="item[field.key]"
          :placeholder="field.placeholder"
          :style="field.style"
          :disabled="disabled"
          :filterable="field.filterable"
          :allow-create="field.allowCreate"
          @change="emitUpdate"
        >
          <template v-if="field.prefix" #prefix>{{ field.prefix }}</template>
          <el-option
            v-for="option in field.options"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          >
            <template v-if="option.icon">
              <i :class="option.icon" style="margin-right: 8px; font-size: 16px"></i>
              {{ option.label }}
            </template>
          </el-option>
        </el-select>

        <!-- 颜色选择器 -->
        <el-color-picker
          v-else-if="field.type === 'color'"
          v-model="item[field.key]"
          :disabled="disabled"
          @change="emitUpdate"
        />
      </template>

      <!-- 删除按钮 -->
      <el-button
        v-if="!hideControls"
        type="danger"
        :icon="Delete"
        circle
        size="small"
        @click="removeItem(index)"
        :disabled="disabled"
      />
    </div>

    <!-- 添加按钮行 -->
    <div v-if="!hideControls" class="editor-item add-row">
      <el-button
        type="primary"
        :icon="Plus"
        circle
        size="small"
        @click="addItem"
        :disabled="disabled"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { Delete, Plus, ArrowUp, ArrowDown } from '@element-plus/icons-vue';

export interface FieldConfig {
  key: string;
  type: 'text' | 'select' | 'color';
  placeholder?: string;
  style?: string;
  prefix?: string;
  filterable?: boolean;
  allowCreate?: boolean;
  options?: Array<{ label: string; value: string; icon?: string }>;
}

export interface JsonListEditorProps {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  modelValue: any[];
  fields: FieldConfig[];
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  defaultItem?: Record<string, any>;
  disabled?: boolean;
  hideControls?: boolean;
}

const props = withDefaults(defineProps<JsonListEditorProps>(), {
  disabled: false,
  defaultItem: () => ({}),
  hideControls: false,
});

const emit = defineEmits<{
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  'update:modelValue': [value: any[]];
}>();

// 内部值（深拷贝避免直接修改 prop）
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const internalValue = ref<any[]>([]);

// 监听 modelValue 变化
watch(
  () => props.modelValue,
  newVal => {
    internalValue.value = JSON.parse(JSON.stringify(newVal || []));
  },
  { immediate: true, deep: true }
);

// 发送更新
const emitUpdate = () => {
  emit('update:modelValue', JSON.parse(JSON.stringify(internalValue.value)));
};

// 上移
const moveUp = (index: number) => {
  if (index <= 0) return;
  [internalValue.value[index], internalValue.value[index - 1]] = [
    internalValue.value[index - 1],
    internalValue.value[index],
  ];
  emitUpdate();
};

// 下移
const moveDown = (index: number) => {
  if (index >= internalValue.value.length - 1) return;
  [internalValue.value[index], internalValue.value[index + 1]] = [
    internalValue.value[index + 1],
    internalValue.value[index],
  ];
  emitUpdate();
};

// 删除项
const removeItem = (index: number) => {
  internalValue.value.splice(index, 1);
  emitUpdate();
};

// 添加项
const addItem = () => {
  const newItem = { ...props.defaultItem };
  // 确保所有字段都有默认值
  props.fields.forEach(({ key }) => {
    if (!(key in newItem)) newItem[key] = '';
  });
  internalValue.value.push(newItem);
  emitUpdate();
};
</script>

<style scoped lang="scss">
.json-list-editor {
  width: 100%;

  .editor-item {
    display: flex;
    align-items: center;
    margin-bottom: 12px;
    gap: 8px;

    &.add-row {
      justify-content: flex-end;
    }
  }
}
</style>
