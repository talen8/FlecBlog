<template>
  <el-form :model="form" label-width="120px" class="setting-form">
    <el-divider content-position="left">基础配置</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('base_url') }">API 端点</span>
      </template>
      <el-input
        v-model="form.base_url"
        placeholder="例如 https://api.deepseek.com"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('api_key') }">API 密钥</span>
      </template>
      <el-input
        v-model="form.api_key"
        type="password"
        show-password
        placeholder="输入 API Key"
        :disabled="loading"
        autocomplete="off"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('model') }">模型名称</span>
      </template>
      <el-input
        v-model="form.model"
        placeholder="例如 deepseek-chat，留空使用内置模型"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item label=" ">
      <el-button :loading="testing" @click="handleTest">测试连接</el-button>
    </el-form-item>

    <el-divider content-position="left">提示词配置</el-divider>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('summary_prompt') }">文章摘要提示词</span>
      </template>
      <el-input
        v-model="form.summary_prompt"
        type="textarea"
        :rows="5"
        placeholder="用于生成文章摘要，留空时使用系统默认提示词"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('ai_summary_prompt') }"
          >AI 总结提示词</span
        >
      </template>
      <el-input
        v-model="form.ai_summary_prompt"
        type="textarea"
        :rows="5"
        placeholder="用于生成 AI 总结，留空时使用系统默认提示词"
        :disabled="loading"
      />
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('title_prompt') }">标题提示词</span>
      </template>
      <el-input
        v-model="form.title_prompt"
        type="textarea"
        :rows="5"
        placeholder="用于生成标题，留空时使用系统默认提示词"
        :disabled="loading"
      />
    </el-form-item>

    <el-divider content-position="left">MCP</el-divider>

    <el-form-item label="认证方式">
      <span>{{ mcpAuthModeLabel }}</span>
    </el-form-item>

    <el-form-item label="OAuth">
      <span>{{ mcpOAuthLabel }}</span>
    </el-form-item>

    <el-form-item v-if="mcpAuthStatus?.mode !== 'oauth'">
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('mcp_secret') }">Secret</span>
      </template>
      <el-input
        v-model="form.mcp_secret"
        type="password"
        show-password
        readonly
        placeholder="系统会自动生成 MCP Secret"
      >
        <template #append>
          <el-button
            type="warning"
            plain
            :disabled="loading || resetting"
            :loading="resetting"
            @click="resetSecret"
            >重置</el-button
          >
        </template>
      </el-input>
    </el-form-item>

    <el-form-item v-if="mcpAuthStatus?.mode !== 'oauth'">
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('mcp_static_operator_user_id') }"
          >管理身份</span
        >
      </template>
      <div class="mcp-operator-field">
        <el-select
          v-model="form.mcp_static_operator_user_id"
          filterable
          :loading="loadingOperators"
          :disabled="loading"
          placeholder="自动选择"
        >
          <el-option label="自动选择" value="" />
          <el-option
            v-if="missingConfiguredOperator"
            :label="'当前配置 ID ' + form.mcp_static_operator_user_id + '（不可用或未加载）'"
            :value="form.mcp_static_operator_user_id"
            disabled
          />
          <el-option
            v-for="user in operatorOptions"
            :key="user.id"
            :label="formatOperatorLabel(user)"
            :value="String(user.id)"
          />
        </el-select>
        <div v-if="operatorLoadError" class="mcp-operator-warning">
          管理员列表加载失败；当前保存值不会被自动清空。
        </div>
      </div>
    </el-form-item>

    <el-form-item>
      <template #label>
        <span :class="{ 'field-modified': isFieldModified('mcp_admin_tools_enabled') }"
          >管理员工具</span
        >
      </template>
      <div class="mcp-admin-tools-field">
        <el-switch
          v-model="form.mcp_admin_tools_enabled"
          active-value="true"
          inactive-value="false"
          :disabled="loading"
        />
        <div class="mcp-admin-tools-hint">
          启用后允许调用用户管理工具；OAuth 连接需重新授权后生效。
        </div>
      </div>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { testAIConfig } from '@/api/ai';
import { getMCPAuthStatus, resetMCPSecret } from '@/api/sysconfig';
import type { MCPAuthStatusResponse } from '@/api/sysconfig';
import { getUsers } from '@/api/user';
import type { User } from '@/types/user';

interface AIForm {
  base_url: string;
  api_key: string;
  model: string;
  summary_prompt: string;
  ai_summary_prompt: string;
  title_prompt: string;
  mcp_secret: string;
  mcp_static_operator_user_id: string;
  mcp_admin_tools_enabled: string;
}

const form = defineModel<AIForm>('form', { required: true });

defineProps<{
  loading?: boolean;
  isFieldModified: (key: string) => boolean;
}>();

const testing = ref(false);
const resetting = ref(false);
const loadingOperators = ref(false);
const operatorLoadError = ref(false);
const operatorOptions = ref<User[]>([]);
const mcpAuthStatus = ref<MCPAuthStatusResponse | null>(null);
const mcpAuthStatusFailed = ref(false);

const mcpAuthModeLabel = computed(() => {
  if (mcpAuthStatusFailed.value) return '未知';
  if (!mcpAuthStatus.value) return '加载中';
  return { static: '密钥', oauth: 'OAuth', hybrid: '兼容' }[mcpAuthStatus.value.mode];
});

const mcpOAuthLabel = computed(() => {
  if (mcpAuthStatusFailed.value) return '未知';
  if (!mcpAuthStatus.value) return '加载中';
  return { disabled: '未启用', embedded: '内置', external: '外部' }[mcpAuthStatus.value.oauth];
});

const missingConfiguredOperator = computed(() => {
  const configuredID = form.value.mcp_static_operator_user_id;
  if (!configuredID || configuredID === '0') return false;
  return !operatorOptions.value.some(user => String(user.id) === configuredID);
});

function formatOperatorLabel(user: User): string {
  const displayName = user.nickname || user.email;
  return displayName + ' · ' + user.email + ' · ' + user.role + ' · ID ' + user.id;
}

async function loadMCPAuthStatus() {
  mcpAuthStatusFailed.value = false;
  try {
    mcpAuthStatus.value = await getMCPAuthStatus();
  } catch {
    mcpAuthStatus.value = null;
    mcpAuthStatusFailed.value = true;
  }
}

async function loadOperatorOptions() {
  loadingOperators.value = true;
  operatorLoadError.value = false;
  try {
    const [superAdmins, admins] = await Promise.all([
      getUsers({
        page: 1,
        page_size: 100,
        role: 'super_admin',
        is_enabled: true,
        is_deleted: false,
      }),
      getUsers({ page: 1, page_size: 100, role: 'admin', is_enabled: true, is_deleted: false }),
    ]);
    const uniqueUsers = new Map<number, User>();
    for (const user of [...superAdmins.list, ...admins.list]) {
      if (user.is_enabled && !user.deleted_at) uniqueUsers.set(user.id, user);
    }
    operatorOptions.value = Array.from(uniqueUsers.values()).sort((a, b) => {
      if (a.role !== b.role) return a.role === 'super_admin' ? -1 : 1;
      return a.id - b.id;
    });
  } catch {
    operatorLoadError.value = true;
  } finally {
    loadingOperators.value = false;
  }
}

onMounted(() => {
  void loadOperatorOptions();
  void loadMCPAuthStatus();
});

async function handleTest() {
  testing.value = true;
  try {
    await testAIConfig({
      base_url: form.value.base_url,
      api_key: form.value.api_key,
      model: form.value.model,
    });
    ElMessage.success('连接成功，配置可用');
  } catch (error: unknown) {
    ElMessage.error((error as Error)?.message || '连接失败，请检查配置');
  } finally {
    testing.value = false;
  }
}

async function resetSecret() {
  try {
    await ElMessageBox.confirm('重置后现有客户端会立刻失效，确定继续吗？', '重置 MCP Secret', {
      type: 'warning',
      confirmButtonText: '确认重置',
      cancelButtonText: '取消',
    });
  } catch {
    return;
  }

  resetting.value = true;
  try {
    const data = await resetMCPSecret();
    form.value.mcp_secret = data.secret || '';
    ElMessage.success('MCP Secret 已重置');
  } catch (error: unknown) {
    ElMessage.error((error as Error)?.message || '重置失败');
  } finally {
    resetting.value = false;
  }
}
</script>

<style lang="scss" scoped>
.mcp-operator-field,
.mcp-admin-tools-field {
  width: 100%;
}

.mcp-operator-field :deep(.el-select) {
  width: 100%;
}

.mcp-operator-warning,
.mcp-admin-tools-hint {
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.5;
  color: var(--el-text-color-secondary);
}

.mcp-operator-warning {
  color: var(--el-color-warning);
}

// 移动端适配
@media (max-width: 768px) {
  :deep(.el-form-item__label) {
    width: 100px !important;
    font-size: 13px;
  }
}
</style>
