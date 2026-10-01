<script setup lang="ts">
const route = useRoute();
const showHeader = computed(() => !!route.meta.typeHeader);

const { y } = useWindowScroll();
const isFixed = computed(() => y.value > 0);
const lastScrollY = ref(0);
const isScrollingDown = ref(false);
const showDrawer = ref(false);

watch(y, newY => {
  isScrollingDown.value = newY > lastScrollY.value;
  lastScrollY.value = newY;
});

// 切换抽屉
const toggleDrawer = () => {
  showDrawer.value = !showDrawer.value;
};
</script>

<template>
  <nav id="navbar" :class="{ fixed: isFixed, 'no-header': !showHeader }">
    <div class="nav-left">
      <LayoutsNavbarAggregate />
      <LayoutsNavbarLogo />
    </div>
    <LayoutsNavbarMenu :is-scrolling-down="isScrollingDown" :is-fixed="isFixed" />
    <LayoutsNavbarButtons @toggle-drawer="toggleDrawer" />
  </nav>

  <LayoutsNavbarMobileDrawer v-model="showDrawer" />
</template>

<style lang="scss">
@mixin text-color {
  color: var(--flec-nav-font);

  &:hover {
    color: var(--flec-nav-font-hover);
  }
}

@mixin fixed-text-color {
  color: var(--flec-nav-fixed-font);

  &:hover {
    color: var(--flec-nav-fixed-font-hover);
  }
}

#navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 6.5rem;
  height: 4rem;
  font-size: 1.1rem;
  position: sticky;
  top: 0;
  left: 0;
  right: 0;
  z-index: 50;
  background-color: transparent;
  transition:
    background-color 0.5s ease,
    opacity 0.5s ease-out;

  .nav-left {
    flex: 1;
    display: flex;
    align-items: center;
  }

  &.no-header {
    background-color: var(--flec-nav-bg);

    .brighten {
      @include fixed-text-color;
    }
  }

  .brighten {
    @include text-color;
    transition: color 0.3s ease;
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
      will-change: transform;
    }

    &:hover:not(.no-after):after {
      transform: translateX(-50%) scaleX(1);
    }
  }

  &.fixed {
    background-color: var(--flec-nav-bg);
    opacity: 1;

    .brighten {
      @include fixed-text-color;
    }
  }
}

@media screen and (max-width: 1024px) {
  #navbar {
    padding: 0 3rem;
  }
}

@media screen and (max-width: 900px) {
  #navbar {
    padding: 0 1rem;
  }
}
</style>
