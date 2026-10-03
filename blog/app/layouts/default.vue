<script lang="ts" setup>
import Navbar from '@/components/layouts/navbar/index.vue';
import Header from '@/components/layouts/header/index.vue';
import Sidebar from '@/components/layouts/sidebar/index.vue';
import Footer from '@/components/layouts/footer/index.vue';
import FloatButton from '@/components/ui/FloatButton.vue';
import MomentWidget from '@/components/features/moment/MomentWidget.vue';

const route = useRoute();
const showSidebar = computed(() => (route.meta.showSidebar as boolean | undefined) ?? true);
const showMomentWidget = computed(() => route.path == '/');
</script>

<template>
  <div class="layout-wrapper">
    <Navbar />
    <Header />
    <main class="page-main">
      <MomentWidget v-if="showMomentWidget" />
      <div class="main-layout">
        <div class="main-content" :class="!showSidebar ? 'full-width' : ''">
          <slot />
        </div>
        <Sidebar v-if="showSidebar" />
      </div>
    </main>
    <Footer />
    <FloatButton />
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

.layout-wrapper {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.page-main {
  background: var(--flec-page-bg);
  width: 100%;
  padding: 40px 0;
  flex: 1;

  .main-layout {
    position: relative;
    display: flex;
    flex: 1 auto;
    margin: 0 auto;
    padding: 0 15px;
    max-width: 1200px;
    width: 100%;

    .main-content {
      width: 74%;
      transition: width 0.3s ease;

      &.full-width {
        width: 100%;
      }
    }
  }
}

// 响应式设计
@media screen and (max-width: 900px) {
  .page-main {
    .main-layout {
      flex-direction: column;
      padding: 0 12px;

      .main-content {
        width: 100%;
      }
    }
  }
}
</style>
