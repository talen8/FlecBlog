<template>
  <div class="sidebar" :class="{ 'is-collapse': isCollapse }">
    <div class="logo">
      <span v-show="!isCollapse">Flec 管理系统</span>
    </div>
    <el-menu
      ref="menuRef"
      :default-active="route.path"
      :collapse="isCollapse"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409eff"
      router
      @select="handleMenuSelect"
    >
      <el-menu-item index="/">
        <i class="ri-dashboard-2-line ri-lg"></i>
        <template #title><span>仪表盘</span></template>
      </el-menu-item>

      <el-sub-menu index="content">
        <template #title>
          <i class="ri-layout-2-line ri-lg"></i>
          <span>内容管理</span>
        </template>
        <el-menu-item index="/articles">
          <i class="ri-article-line ri-lg"></i>
          <template #title>文章管理</template>
        </el-menu-item>
        <el-menu-item index="/moments">
          <i class="ri-chat-3-line ri-lg"></i>
          <template #title>动态管理</template>
        </el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="interaction">
        <template #title>
          <i class="ri-chat-2-line ri-lg"></i>
          <span>互动管理</span>
        </template>
        <el-menu-item index="/friends">
          <i class="ri-links-line ri-lg"></i>
          <template #title>友链管理</template>
        </el-menu-item>
        <el-menu-item index="/comments">
          <i class="ri-message-3-line ri-lg"></i>
          <template #title>评论管理</template>
        </el-menu-item>
        <el-menu-item index="/rssfeed">
          <i class="ri-rss-line ri-lg"></i>
          <template #title>RSS订阅</template>
        </el-menu-item>
        <el-menu-item index="/feedback">
          <i class="ri-feedback-line ri-lg"></i>
          <template #title>反馈投诉</template>
        </el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="management">
        <template #title>
          <i class="ri-admin-line ri-lg"></i>
          <span>系统管理</span>
        </template>
        <el-menu-item index="/users">
          <i class="ri-team-line ri-lg"></i>
          <template #title>用户管理</template>
        </el-menu-item>
        <el-menu-item index="/files">
          <i class="ri-folder-image-line ri-lg"></i>
          <template #title>文件管理</template>
        </el-menu-item>
        <el-menu-item index="/menus">
          <i class="ri-menu-line ri-lg"></i>
          <template #title>菜单管理</template>
        </el-menu-item>
        <el-menu-item index="/visits">
          <i class="ri-file-list-3-line ri-lg"></i>
          <template #title>访问日志</template>
        </el-menu-item>
        <el-menu-item index="/systems">
          <i class="ri-information-line ri-lg"></i>
          <template #title>系统信息</template>
        </el-menu-item>
        <el-menu-item index="/settings">
          <i class="ri-settings-3-line ri-lg"></i>
          <template #title>系统设置</template>
        </el-menu-item>
      </el-sub-menu>
    </el-menu>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router';
import { ref, nextTick } from 'vue';

const route = useRoute();
const menuRef = ref();

// 子菜单索引列表
const subMenus = ['content', 'interaction', 'management'];

// 路由到父菜单的映射
const routeToParent: Record<string, string> = {
  '/articles': 'content',
  '/moments': 'content',
  '/friends': 'interaction',
  '/comments': 'interaction',
  '/rssfeed': 'interaction',
  '/feedback': 'interaction',
  '/users': 'management',
  '/files': 'management',
  '/menus': 'management',
  '/visits': 'management',
  '/systems': 'management',
  '/settings': 'management',
};

defineProps<{ isCollapse: boolean }>();
const emit = defineEmits(['menu-click']);

// 选择菜单项时关闭其他子菜单
const handleMenuSelect = (index: string) => {
  const parent = routeToParent[index];
  if (parent) {
    nextTick(() => {
      subMenus.forEach(key => {
        if (key !== parent) menuRef.value?.close(key);
      });
    });
  }
  emit('menu-click');
};
</script>

<style scoped lang="scss">
.sidebar {
  height: 100%;
  display: flex;
  flex-direction: column;

  .logo {
    height: 60px;
    padding: 10px 0;
    margin-bottom: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;

    span {
      color: #fff;
      font-size: 16px;
      font-weight: 600;
    }
  }

  &.is-collapse {
    .logo {
      padding: 10px;

      span {
        display: none;
      }
    }
  }

  :deep(.el-menu) {
    border-right: none;
    flex: 1;
    overflow-y: auto;

    &::-webkit-scrollbar {
      width: 4px;
    }

    &::-webkit-scrollbar-thumb {
      background-color: rgba(255, 255, 255, 0.2);
      border-radius: 2px;
    }

    &::-webkit-scrollbar-track {
      background-color: transparent;
    }

    .el-menu-item,
    .el-sub-menu__title {
      font-size: 16px;

      i {
        margin-right: 10px;
      }
    }
  }
}
</style>
