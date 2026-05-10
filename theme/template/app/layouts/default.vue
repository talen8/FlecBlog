<script lang="ts" setup></script>

<template>
  <div class="layout-wrapper">
    <LayoutNavbar />
    <main class="page-main">
      <div class="main-layout">
        <div class="main-content">
          <slot />
        </div>
        <LayoutSidebar />
      </div>
    </main>
    <LayoutFooter />
  </div>
</template>

<style lang="scss" scoped>
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
      flex: 1;
      min-width: 0;
    }
  }
}

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
