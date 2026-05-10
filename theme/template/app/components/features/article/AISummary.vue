<script lang="ts" setup>
defineProps<{ aiSummary: string }>();
</script>

<template>
  <div class="ai-summary">
    <div class="summary-header">
      <i class="ri-robot-2-line" />
      <span>AI 摘要</span>
    </div>
    <p class="summary-text">{{ aiSummary }}</p>
  </div>
</template>

<style lang="scss" scoped>
.ai-summary {
  margin: 20px 0;
  padding: 16px 20px;
  background: var(--flec-heavy-bg);
  border: 1px solid var(--theme-color);
}

.summary-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  color: var(--theme-color);
  font-weight: 500;
  font-size: 0.95rem;
}

.summary-text {
  margin: 0;
  color: var(--font-color);
  font-size: 0.95rem;
  line-height: 1.7;
}
</style>
