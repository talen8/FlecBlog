<script setup lang="ts">
/**
 * 文章过时提醒组件
 */
</script>

<template>
  <div class="outdated-notice">
    <i class="ri-error-warning-line" />
    <span>本文内容可能已过时或不再适用，请结合实际情况谨慎参考。</span>
  </div>
</template>

<style lang="scss" scoped>
.outdated-notice {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  margin-bottom: 16px;
  background: #fff3cd;
  border: 1px solid #ffc107;
  border-radius: 6px;
  font-size: 0.9rem;
  color: #856404;
}

@media (prefers-color-scheme: dark) {
  .outdated-notice {
    background: #4a3d1a;
    border-color: #6d5a00;
    color: #ffc107;
  }
}
</style>
