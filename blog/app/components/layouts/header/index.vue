<script lang="ts" setup>
const route = useRoute();
const type = computed(() => route.meta.typeHeader as string | undefined);
</script>

<template>
  <LayoutsHeaderHomeHeader v-if="type === 'home'" />
  <LayoutsHeaderPostHeader v-else-if="type === 'post'" />
</template>
