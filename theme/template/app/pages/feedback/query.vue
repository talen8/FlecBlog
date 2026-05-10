<script lang="ts" setup>
import { getFeedbackByTicketNo } from '@/composables/useApi';
import type { Feedback } from '@@/types';

const toast = useToast();

const ticketNo = ref('');
const loading = ref(false);
const feedback = ref<Feedback | null>(null);
const notFound = ref(false);

const statusMap: Record<string, string> = {
  pending: '处理中',
  resolved: '已解决',
  closed: '已关闭',
};

const handleQuery = async () => {
  if (!ticketNo.value.trim()) {
    toast.error('请输入工单号');
    return;
  }

  loading.value = true;
  notFound.value = false;
  feedback.value = null;

  try {
    feedback.value = await getFeedbackByTicketNo(ticketNo.value.trim());
    // eslint-disable-next-line @typescript-eslint/no-explicit-any -- 错误处理，需要访问 response.status
  } catch (error: any) {
    if (error?.response?.status === 404) {
      notFound.value = true;
    } else {
      console.error('查询失败:', error);
      toast.error('查询失败，请稍后重试');
    }
  } finally {
    loading.value = false;
  }
};

useSeoMeta({
  title: '查询反馈',
});
</script>

<template>
  <div id="page">
    <h1 class="page-title">查询反馈</h1>
    <div class="query-form">
      <input
        v-model="ticketNo"
        type="text"
        placeholder="请输入工单号"
        class="ticket-input"
        @keyup.enter="handleQuery"
      />
      <button class="query-btn" :disabled="loading" @click="handleQuery">
        {{ loading ? '查询中...' : '查询' }}
      </button>
    </div>

    <div v-if="feedback" class="feedback-detail">
      <div class="detail-row">
        <span class="label">工单号</span>
        <span class="value">{{ feedback.ticket_no }}</span>
      </div>
      <div class="detail-row">
        <span class="label">状态</span>
        <span class="value status" :class="feedback.status">{{
          statusMap[feedback.status] || feedback.status
        }}</span>
      </div>
      <div class="detail-row">
        <span class="label">反馈类型</span>
        <span class="value">{{ feedback.report_type }}</span>
      </div>
      <div class="detail-row">
        <span class="label">提交时间</span>
        <span class="value">{{ formatDate(feedback.feedback_time) }}</span>
      </div>
      <div v-if="feedback.admin_reply" class="detail-block">
        <span class="label">管理员回复</span>
        <p class="reply-content">{{ feedback.admin_reply }}</p>
      </div>
    </div>

    <div v-else-if="notFound" class="empty-state">未找到该工单号的反馈记录</div>

    <NuxtLink to="/feedback" class="back-link">返回提交反馈</NuxtLink>
  </div>
</template>

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
  margin: 0 0 20px;
  font-size: 2rem;
  font-weight: 700;
}

.query-form {
  display: flex;
  gap: 12px;
  margin-bottom: 30px;
}

.ticket-input {
  flex: 1;
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

.query-btn {
  padding: 12px 24px;
  background: var(--theme-color);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 0.95rem;
  cursor: pointer;

  &:hover:not(:disabled) {
    opacity: 0.9;
  }
  &:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
}

.feedback-detail {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 20px;
  background: var(--flec-heavy-bg);
  border-radius: 8px;
  margin-bottom: 20px;
}

.detail-row {
  display: flex;
  gap: 16px;

  .label {
    color: var(--theme-meta-color);
    min-width: 80px;
  }
  .value {
    color: var(--font-color);
  }

  .status.resolved {
    color: #52c41a;
  }
  .status.closed {
    color: var(--theme-meta-color);
  }
  .status.pending {
    color: #faad14;
  }
}

.detail-block {
  .label {
    color: var(--theme-meta-color);
    display: block;
    margin-bottom: 8px;
  }
}

.reply-content {
  margin: 0;
  padding: 12px;
  background: var(--flec-card-bg);
  border-radius: 6px;
  line-height: 1.6;
  color: var(--font-color);
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: var(--theme-meta-color);
}

.back-link {
  display: inline-block;
  margin-top: 16px;
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
  .query-form {
    flex-direction: column;
  }
}
</style>
