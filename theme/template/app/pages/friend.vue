<script lang="ts" setup>
import { getFriends, applyFriend } from '@/composables/useApi';
import type { FriendGroup } from '@@/types';

const groups = ref<FriendGroup[]>([]);
const isLoggedIn = useAuth();
const showLoginModal = useState('showLoginModal', () => false);
const toast = useToast();

const showApplyForm = ref(false);
const applyForm = reactive({
  name: '',
  url: '',
  description: '',
  avatar: '',
  screenshot: '',
});
const applying = ref(false);

const { data: initialData } = await useAsyncData('friends-list', async () => {
  const { groups: friendGroups } = await getFriends();
  return { groups: friendGroups };
});

if (initialData.value) {
  groups.value = initialData.value.groups || [];
}

const handleApply = async () => {
  if (!isLoggedIn.value) {
    showLoginModal.value = true;
    return;
  }
  if (!applyForm.name || !applyForm.url || !applyForm.description || !applyForm.avatar) {
    toast.error('请填写完整信息');
    return;
  }
  applying.value = true;
  try {
    await applyFriend(applyForm);
    toast.success('申请已提交，等待审核');
    showApplyForm.value = false;
    applyForm.name = '';
    applyForm.url = '';
    applyForm.description = '';
    applyForm.avatar = '';
    applyForm.screenshot = '';
  } catch (error) {
    console.error('申请失败:', error);
    toast.error('申请失败');
  } finally {
    applying.value = false;
  }
};

useSeoMeta({
  title: '友情链接',
  description: '我的友情链接',
});
</script>

<template>
  <div id="page">
    <div class="page-header">
      <h1 class="page-title">友情链接</h1>
      <button class="apply-btn" @click="showApplyForm = !showApplyForm">
        <i class="ri-add-line" /> 申请友链
      </button>
    </div>

    <div v-if="showApplyForm" class="apply-form">
      <div class="form-group">
        <label>名称 *</label>
        <input v-model="applyForm.name" placeholder="网站名称" />
      </div>
      <div class="form-group">
        <label>链接 *</label>
        <input v-model="applyForm.url" type="url" placeholder="https://" />
      </div>
      <div class="form-group">
        <label>描述 *</label>
        <input v-model="applyForm.description" placeholder="网站描述" />
      </div>
      <div class="form-group">
        <label>头像 *</label>
        <input v-model="applyForm.avatar" type="url" placeholder="头像链接" />
      </div>
      <div class="form-group">
        <label>截图（可选）</label>
        <input v-model="applyForm.screenshot" type="url" placeholder="网站截图链接" />
      </div>
      <button class="submit-btn" :disabled="applying" @click="handleApply">
        {{ applying ? '提交中...' : '提交申请' }}
      </button>
    </div>

    <div v-if="groups.length > 0" class="friend-list">
      <div v-for="group in groups" :key="group.type_id ?? 'default'" class="friend-group">
        <h2 v-if="group.type_name" class="group-title">{{ group.type_name }}</h2>
        <div class="friend-grid">
          <a
            v-for="friend in group.friends"
            :key="friend.id"
            :href="friend.url"
            target="_blank"
            rel="noopener noreferrer"
            class="friend-card"
            :class="{ invalid: friend.is_invalid }"
          >
            <img :src="friend.avatar" :alt="friend.name" class="friend-avatar" />
            <div class="friend-info">
              <span class="friend-name">{{ friend.name }}</span>
              <span class="friend-desc">{{ friend.description }}</span>
            </div>
          </a>
        </div>
      </div>
    </div>
    <div v-else class="empty-state">暂无友链</div>
  </div>
</template>

<style lang="scss" scoped>
#page {
  background: var(--flec-card-bg);
  border: 1px solid var(--flec-border);
  padding: 40px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-title {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
}

.apply-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: var(--theme-color);
  color: #fff;
  border: none;
  font-size: 0.9rem;
  cursor: pointer;
}

.apply-form {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 30px;
  padding: 20px;
  background: var(--flec-heavy-bg);

  .submit-btn {
    grid-column: span 2;
    padding: 10px;
    background: var(--theme-color);
    color: #fff;
    border: none;
    cursor: pointer;
    &:disabled {
      opacity: 0.6;
      cursor: not-allowed;
    }
  }
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
  label {
    font-size: 0.85rem;
    color: var(--theme-meta-color);
  }
  input {
    padding: 10px 12px;
    border: 1px solid var(--flec-border);
    background: var(--flec-card-bg);
    color: var(--font-color);
    font-size: 0.9rem;
    &:focus {
      outline: none;
      border-color: var(--theme-color);
    }
  }
}

.friend-group {
  margin-bottom: 30px;
  &:last-child {
    margin-bottom: 0;
  }
}

.group-title {
  font-size: 1.2rem;
  font-weight: 600;
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--flec-border);
}

.friend-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.friend-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: var(--flec-heavy-bg);
  text-decoration: none;
  &.invalid {
    opacity: 0.5;
  }
}

.friend-avatar {
  width: 48px;
  height: 48px;
  object-fit: cover;
  flex-shrink: 0;
}

.friend-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.friend-name {
  font-weight: 500;
  color: var(--font-color);
}

.friend-desc {
  font-size: 0.85rem;
  color: var(--theme-meta-color);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--theme-meta-color);
}

@media screen and (max-width: 768px) {
  #page {
    padding: 20px;
  }
  .page-title {
    font-size: 1.5rem;
  }
  .apply-form {
    grid-template-columns: 1fr;
    .submit-btn {
      grid-column: span 1;
    }
  }
}
</style>
