<script setup lang="ts">
interface Emits {
  (e: 'toggleDrawer'): void;
}

const emit = defineEmits<Emits>();

// 登录相关
const isLoggedIn = useAuth();
const { open: openLogin } = useLoginModal();

// 用户信息
const { userAvatar, userNickname, userEmail, fetchUserInfo, clearUserInfo } = useUser();

// 通知相关
const { unreadCount, clearNotifications, fetchNotifications } = useNotifications();

let pollingTimer: number | null = null;

// 监听登录状态，自动启动/停止轮询（仅在客户端执行）
watch(
  isLoggedIn,
  loggedIn => {
    // 只在客户端执行
    if (!import.meta.client) return;

    // 清理旧定时器
    if (pollingTimer) {
      clearInterval(pollingTimer);
      pollingTimer = null;
    }

    if (loggedIn) {
      // 获取用户信息
      fetchUserInfo();
      // 30秒轮询一次未读通知数量
      pollingTimer = window.setInterval(() => {
        fetchNotifications({ page: 1, page_size: 1 });
      }, 30000);
    } else {
      clearUserInfo();
      clearNotifications();
    }
  },
  { immediate: true }
);

onUnmounted(() => {
  if (pollingTimer) {
    clearInterval(pollingTimer);
  }
});

// 搜索相关状态
const showSearchModal = ref(false);

// 打开搜索弹窗
const openSearch = () => {
  showSearchModal.value = true;
};

// 用户菜单显示状态
const showUserMenu = ref(false);
const userMenuRef = ref<HTMLElement>();

// 点击外部关闭菜单
onClickOutside(userMenuRef, () => {
  showUserMenu.value = false;
});

// 切换用户菜单
const toggleUserMenu = () => {
  showUserMenu.value = !showUserMenu.value;
};

// 退出登录
const handleLogout = () => {
  showUserMenu.value = false;
  logout();
};
</script>

<template>
  <div class="nav-button">
    <button class="brighten" aria-label="搜索" @click="openSearch">
      <i class="ri-search-line ri-xl" />
    </button>
    <!-- 主题切换按钮 - 使用 CSS 控制图标显示，避免 SSR 闪烁 -->
    <button class="brighten theme-toggle" aria-label="切换主题" @click="toggleTheme">
      <i class="ri-moon-line ri-xl theme-icon-moon" />
      <i class="ri-sun-line ri-xl theme-icon-sun" />
    </button>
    <!-- 登录按钮 - 客户端渲染避免 hydration mismatch -->
    <ClientOnly>
      <button v-if="!isLoggedIn" class="brighten login-btn" aria-label="登录" @click="openLogin">
        <i class="ri-user-line ri-xl" />
      </button>
      <div v-else ref="userMenuRef" class="user-menu">
        <button class="brighten user-btn" aria-label="用户菜单" @click="toggleUserMenu">
          <i class="ri-user-3-fill ri-xl" />
        </button>
        <Transition name="dropdown">
          <div v-show="showUserMenu" class="user-dropdown" @click.stop>
            <div class="user-info">
              <img :src="userAvatar" :alt="userNickname" class="user-avatar-large" />
              <div class="user-details">
                <span class="user-nickname">{{ userNickname }}</span>
                <span class="user-email">{{ userEmail }}</span>
              </div>
            </div>
            <a href="/profile" class="dropdown-item" @click="showUserMenu = false">
              <i class="ri-user-settings-line" />
              个人设置
            </a>
            <a
              href="/notifications"
              class="dropdown-item notification-item"
              @click="showUserMenu = false"
            >
              <i class="ri-notification-3-line" />
              <span>通知中心</span>
              <span v-if="unreadCount > 0" class="notification-badge">{{
                unreadCount > 99 ? '99+' : unreadCount
              }}</span>
            </a>
            <button class="dropdown-item" @click="handleLogout">
              <i class="ri-logout-box-line" />
              退出登录
            </button>
          </div>
        </Transition>
      </div>
      <template #fallback>
        <button class="brighten login-btn" aria-label="登录">
          <i class="ri-user-line ri-xl" />
        </button>
      </template>
    </ClientOnly>
    <button class="button-menu brighten" aria-label="打开菜单" @click="emit('toggleDrawer')">
      <i class="ri-menu-line ri-xl" />
    </button>
  </div>

  <!-- 搜索弹窗 -->
  <FeaturesModalsSearchModal v-model="showSearchModal" />
</template>

<style lang="scss" scoped>
.nav-button {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 0.5rem;

  .button-menu {
    display: none;
  }

  // 主题切换按钮图标控制
  .theme-toggle {
    position: relative;
    width: 24px;
    height: 24px;
    display: inline-flex;
    align-items: center;
    justify-content: center;

    .theme-icon-moon,
    .theme-icon-sun {
      position: absolute;
      transition: opacity 0.2s ease;
    }
  }

  .login-btn {
    position: relative;

    &:after {
      content: '';
      position: absolute;
      bottom: -8px;
      left: 50%;
      width: 90%;
      height: 1px;
      background-color: var(--flec-nav-focus);
      transform: translateX(-50%) scaleX(0);
      transform-origin: center;
      transition: transform 0.3s ease;
    }

    &:hover:after {
      transform: translateX(-50%) scaleX(1);
    }
  }

  .user-menu {
    position: relative;
    display: inline-block;

    .user-btn {
      cursor: pointer;
    }

    .user-dropdown {
      position: absolute;
      top: calc(100% + 0.75rem);
      right: 0;
      background-color: var(--flec-card-bg);
      backdrop-filter: blur(3px);
      border-radius: 0.75rem;
      box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
      min-width: 180px;
      padding: 0.5rem;
      z-index: 100;

      .user-info {
        display: flex;
        align-items: center;
        gap: 0.6rem;
        padding: 0.5rem 0.75rem;
        margin-bottom: 0.25rem;

        .user-avatar-large {
          width: 40px;
          height: 40px;
          border-radius: 50%;
          object-fit: cover;
          flex-shrink: 0;
        }

        .user-details {
          display: flex;
          flex-direction: column;
          gap: 0.15rem;
          min-width: 0;
          flex: 1;

          .user-nickname {
            font-weight: 600;
            font-size: 0.9rem;
            color: var(--font-color);
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            line-height: 1.2;
          }

          .user-email {
            font-size: 0.75rem;
            color: var(--theme-meta-color);
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            line-height: 1.2;
          }
        }
      }

      .dropdown-item {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        width: 100%;
        padding: 0.65rem 0.75rem;
        border: none;
        background: transparent;
        color: var(--font-color);
        cursor: pointer;
        border-radius: 0.5rem;
        transition: all 0.2s ease;
        text-align: left;
        text-decoration: none;
        font-size: 0.9rem;
        line-height: 1.2;
        font-family: inherit;

        i {
          font-size: 1.15rem;
          flex-shrink: 0;
          line-height: 1;
        }

        &:hover {
          background-color: var(--flec-nav-menu-bg-hover);
          color: #fff;
        }

        &.notification-item {
          position: relative;

          .notification-badge {
            margin-left: auto;
            padding: 0.15rem 0.45rem;
            background: linear-gradient(135deg, #ff6b6b, #ee5a52);
            color: white;
            border-radius: 12px;
            font-size: 0.7rem;
            font-weight: 700;
            min-width: 22px;
            text-align: center;
            box-shadow: 0 2px 8px rgba(255, 107, 107, 0.3);
            line-height: 1;
          }
        }
      }
    }
  }

  // 下拉菜单动画
  .dropdown-enter-active,
  .dropdown-leave-active {
    transition: all 0.2s ease;
  }

  .dropdown-enter-from,
  .dropdown-leave-to {
    opacity: 0;
    transform: translateY(-10px);
  }

  .dropdown-enter-to,
  .dropdown-leave-from {
    opacity: 1;
    transform: translateY(0);
  }
}

@media screen and (max-width: 900px) {
  .button-menu {
    display: inline-flex !important;
  }
}
</style>
