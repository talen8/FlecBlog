<template>
  <div class="admin-layout">
    <el-container class="layout-container">
      <el-aside
        :width="sidebarWidth"
        class="layout-sidebar"
        :class="{ 'is-mobile-open': mobileSidebarVisible }"
      >
        <Sidebar :is-collapse="sidebarCollapsed" @menu-click="handleMenuClick" />
      </el-aside>

      <div v-show="mobileSidebarVisible" class="sidebar-overlay" @click="closeMobileSidebar"></div>

      <el-container>
        <el-header>
          <Header
            :layout-mode="layoutMode"
            :sidebar-collapsed="sidebarCollapsed"
            @toggle-sidebar="toggleSidebar"
            @open-mobile-sidebar="openMobileSidebar"
          />
        </el-header>
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import Header from '@/components/layouts/Header.vue';
import Sidebar from '@/components/layouts/Sidebar.vue';

const route = useRoute();

const sidebarCollapsed = ref(false);
const mobileSidebarVisible = ref(false);

const sidebarWidth = computed(() => {
  return sidebarCollapsed.value ? '64px' : '200px';
});

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value;
};

const layoutMode = 'fixed';

const openMobileSidebar = () => {
  mobileSidebarVisible.value = true;
};

const closeMobileSidebar = () => {
  mobileSidebarVisible.value = false;
};

const handleMenuClick = () => {
  closeMobileSidebar();
};

watch(
  () => route.fullPath,
  () => {
    closeMobileSidebar();
  }
);
</script>

<style scoped lang="scss">
.admin-layout {
  height: 100vh;
  position: relative;
}

.layout-container {
  height: 100%;
}

.layout-sidebar {
  background-color: #304156;
  overflow: hidden;
  transition:
    width 0.3s,
    left 0.3s;

  @media (max-width: 768px) {
    position: fixed;
    left: -200px;
    top: 0;
    bottom: 0;
    width: 200px !important;
    z-index: 2000;

    &.is-mobile-open {
      left: 0;
    }
  }
}

.sidebar-overlay {
  display: none;

  @media (max-width: 768px) {
    display: block;
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 1999;
    opacity: 1;
    cursor: pointer;
  }
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  padding: 0;
  height: 60px;
}

.el-container {
  height: 100%;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow: auto;

  @media (max-width: 768px) {
    padding: 12px;
  }
}
</style>
