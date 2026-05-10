<template>
  <div id="page">
    <h1 class="page-title">反馈</h1>

    <div v-if="submittedTicket" class="success-state">
      <p>反馈提交成功！</p>
      <p>
        您的工单号为：<strong>{{ submittedTicket }}</strong>
      </p>
      <p class="hint">请保存工单号，您可以使用它查询反馈处理进度。</p>
      <div class="success-actions">
        <button class="action-btn" @click="submittedTicket = ''">继续反馈</button>
        <NuxtLink to="/feedback/query" class="action-btn">查询进度</NuxtLink>
      </div>
    </div>

    <div v-else class="feedback-form">
      <div class="form-group">
        <label>反馈类型</label>
        <select v-model="form.reportType">
          <option value="suggestion">建议</option>
          <option value="summary">内容纠错</option>
          <option value="copyright">版权投诉</option>
          <option value="inappropriate">不当内容</option>
        </select>
      </div>
      <div class="form-group">
        <label>相关链接</label>
        <input v-model="form.reportUrl" type="url" placeholder="相关页面链接（可选）" />
      </div>
      <div class="form-group">
        <label>邮箱</label>
        <input v-model="form.email" type="email" placeholder="您的邮箱（可选）" />
      </div>
      <div class="form-group">
        <label>详细描述</label>
        <textarea v-model="form.description" rows="5" placeholder="请详细描述您的反馈..." />
      </div>
      <button class="submit-btn" :disabled="submitting" @click="handleSubmit">
        {{ submitting ? '提交中...' : '提交反馈' }}
      </button>
      <NuxtLink to="/feedback/query" class="query-link">查询反馈进度</NuxtLink>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { submitFeedback } from '@/composables/useApi';

const toast = useToast();

const form = reactive({
  reportType: 'suggestion' as 'suggestion' | 'summary' | 'copyright' | 'inappropriate',
  reportUrl: '',
  email: '',
  description: '',
});

const submitting = ref(false);
const submittedTicket = ref('');

const handleSubmit = async () => {
  if (!form.description.trim()) {
    toast.error('请填写详细描述');
    return;
  }

  submitting.value = true;
  try {
    const result = await submitFeedback({
      reportType: form.reportType,
      reportUrl: form.reportUrl || undefined,
      email: form.email || undefined,
      description: form.description.trim(),
    });
    submittedTicket.value = result.ticket_no || '';
    toast.success('反馈提交成功');
    form.description = '';
    form.reportUrl = '';
    form.email = '';
  } catch (error) {
    console.error('提交失败:', error);
    toast.error('提交失败，请稍后重试');
  } finally {
    submitting.value = false;
  }
};

useSeoMeta({
  title: '反馈',
});
</script>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border-radius: 12px;
  border: 1px solid var(--flec-border);
  padding: 40px;
  max-width: 600px;
  margin: 0 auto;
}

.page-title {
  margin: 0 0 30px;
  font-size: 2rem;
  font-weight: 700;
}

.success-state {
  text-align: center;

  p {
    margin-bottom: 8px;
  }
  strong {
    color: var(--theme-color);
    font-size: 1.2rem;
  }
  .hint {
    color: var(--theme-meta-color);
    font-size: 0.85rem;
    margin-bottom: 20px;
  }
}

.success-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
}

.action-btn {
  padding: 10px 20px;
  border: 1px solid var(--flec-border);
  border-radius: 8px;
  background: var(--flec-card-bg);
  color: var(--font-color);
  font-size: 0.9rem;
  cursor: pointer;
  text-decoration: none;

  &:hover {
    border-color: var(--theme-color);
    color: var(--theme-color);
  }
}

.feedback-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;

  label {
    font-weight: 500;
    color: var(--font-color);
  }

  input,
  select,
  textarea {
    padding: 12px;
    border: 1px solid var(--flec-border);
    border-radius: 8px;
    background: var(--flec-card-bg);
    color: var(--font-color);
    font-size: 0.95rem;

    &:focus {
      outline: none;
      border-color: var(--theme-color);
    }
  }

  textarea {
    resize: vertical;
  }
}

.submit-btn {
  padding: 12px 24px;
  background: var(--theme-color);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  cursor: pointer;

  &:hover:not(:disabled) {
    opacity: 0.9;
  }
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.query-link {
  text-align: center;
  color: var(--theme-color);
  text-decoration: none;
  font-size: 0.9rem;
}

@media screen and (max-width: 768px) {
  #page {
    padding: 20px;
  }
  .page-title {
    font-size: 1.5rem;
  }
}
</style>
