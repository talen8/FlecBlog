<script setup lang="ts">
const { basicConfig } = useSysConfig();
</script>

<template>
  <div class="logo-container">
    <a href="/" class="logo brighten" aria-label="返回首页">
      {{ basicConfig.subtitle }}
    </a>
  </div>
</template>

<style lang="scss" scoped>
.logo-container {
  .logo {
    font-weight: bold;
    font-size: 1.25rem;
  }
}
</style>
