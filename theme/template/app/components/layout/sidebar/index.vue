<template>
  <aside class="sidebar">
    <LayoutSidebarAuthorCard />
    <LayoutSidebarAnnouncementCard />
    <LayoutSidebarWebInfoCard />
  </aside>
</template>

<style lang="scss" scoped>
.sidebar {
  width: 260px;
  flex-shrink: 0;
  padding-left: 20px;
}

.sidebar-card {
  padding: 16px;
  margin-bottom: 16px;
  background: var(--flec-card-bg);
  border: 1px solid var(--flec-border);
}

@media screen and (max-width: 900px) {
  .sidebar {
    width: 100%;
    padding-left: 0;
    margin-top: 16px;
  }
}
</style>
