<script lang="ts" setup>
const { blogConfig } = useSysConfig();
const { navigationMenus } = useMenus();
const isLoggedIn = useAuth();
const { userInfo } = useUser();
const { unreadCount } = useNotifications();
const showLoginModal = useState('showLoginModal', () => false);
const showMobileMenu = ref(false);

const siteSubtitle = computed(() => blogConfig.value?.['blog.subtitle'] || 'Blog');
</script>

<template>
  <nav class="navbar">
    <div class="navbar-container">
      <NuxtLink class="navbar-logo" to="/">
        {{ siteSubtitle }}
      </NuxtLink>

      <div class="navbar-menu" :class="{ active: showMobileMenu }">
        <NuxtLink
          v-for="menu in navigationMenus"
          :key="menu.id"
          :to="menu.url"
          class="navbar-link"
          @click="showMobileMenu = false"
        >
          <i v-if="menu.icon" :class="menu.icon" />
          {{ menu.title }}
        </NuxtLink>
      </div>

      <div class="navbar-actions">
        <NuxtLink class="action-btn" to="/search" title="搜索">
          <i class="ri-search-line" />
        </NuxtLink>
        <template v-if="isLoggedIn">
          <NuxtLink class="action-btn" to="/notifications" title="通知">
            <i class="ri-notification-line" />
            <span v-if="unreadCount > 0" class="badge">{{
              unreadCount > 99 ? '99+' : unreadCount
            }}</span>
          </NuxtLink>
          <NuxtLink class="action-btn user-btn" to="/profile">
            <img :src="getAvatarUrl(userInfo || {})" alt="avatar" />
          </NuxtLink>
        </template>
        <template v-else>
          <button class="login-btn" @click="showLoginModal = true">登录</button>
        </template>
        <button class="mobile-toggle" @click="showMobileMenu = !showMobileMenu">
          <i :class="showMobileMenu ? 'ri-close-line' : 'ri-menu-line'" />
        </button>
      </div>
    </div>
  </nav>

  <LazyFeaturesModalsLoginModal />
</template>

<style lang="scss" scoped>
.navbar {
  position: sticky;
  top: 0;
  z-index: 100;
  background: var(--flec-card-bg);
  border-bottom: 1px solid var(--flec-border);
}

.navbar-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
  height: 60px;
}

.navbar-logo {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--font-color);
  text-decoration: none;
}

.navbar-menu {
  display: flex;
  gap: 24px;

  &.active {
    display: flex;
    position: absolute;
    top: 60px;
    left: 0;
    right: 0;
    flex-direction: column;
    background: var(--flec-card-bg);
    border-bottom: 1px solid var(--flec-border);
    padding: 16px;
    gap: 8px;
  }
}

.navbar-link {
  display: flex;
  align-items: center;
  gap: 6px;
  color: var(--font-color);
  text-decoration: none;
  font-size: 0.95rem;

  &:hover,
  &.router-link-active {
    color: var(--theme-color);
  }
}

.navbar-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.action-btn {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  color: var(--font-color);
  text-decoration: none;

  &:hover {
    color: var(--theme-color);
  }

  .badge {
    position: absolute;
    top: 0;
    right: 0;
    min-width: 16px;
    height: 16px;
    padding: 0 4px;
    background: #ff4d4f;
    color: #fff;
    font-size: 0.7rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

.user-btn img {
  width: 32px;
  height: 32px;
  object-fit: cover;
}

.login-btn {
  padding: 8px 20px;
  background: var(--theme-color);
  color: #fff;
  border: none;
  font-size: 0.9rem;
  cursor: pointer;
}

.mobile-toggle {
  display: none;
  width: 36px;
  height: 36px;
  align-items: center;
  justify-content: center;
  color: var(--font-color);
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.25rem;
}

@media screen and (max-width: 768px) {
  .navbar-menu {
    display: none;
  }
  .mobile-toggle {
    display: flex;
  }
}
</style>
