<script setup lang="ts">
interface EmojiItem {
  key: string;
  val: string;
}

interface EmojiGroup {
  name: string;
  type: 'emoji' | 'image' | 'emoticon';
  items: EmojiItem[];
}

const emit = defineEmits<{
  select: [emoji: string];
}>();

const { blogConfig } = useSysConfig();
const emojiGroups = ref<EmojiGroup[]>([]);
const activeTab = ref(0);
const loading = ref(true);
const errorMsg = ref('');

// 加载表情包数据
const loadEmojis = async () => {
  const emojisUrl = blogConfig.value.emojis;
  if (!emojisUrl) {
    errorMsg.value = '未配置表情包';
    loading.value = false;
    return;
  }

  try {
    const response = await fetch(emojisUrl);
    if (!response.ok) throw new Error('加载表情包失败');
    emojiGroups.value = await response.json();
  } catch (error: unknown) {
    const err = error as Error;
    errorMsg.value = err.message || '加载表情包失败';
  } finally {
    loading.value = false;
  }
};

// 选择表情
const selectEmoji = (item: EmojiItem, type: string) => {
  emit('select', type === 'image' ? `:${item.key}:` : item.val);
};

onMounted(loadEmojis);
</script>

<template>
  <div class="emoji-picker">
    <div v-if="loading" class="emoji-state">
      <i class="ri-loader-4-line rotating" />
      <span>加载中...</span>
    </div>

    <div v-else-if="errorMsg" class="emoji-state">
      <i class="ri-emotion-unhappy-line" />
      <span>{{ errorMsg }}</span>
    </div>

    <template v-else>
      <div class="emoji-tabs">
        <button
          v-for="(group, index) in emojiGroups"
          :key="index"
          class="emoji-tab"
          :class="{ active: activeTab === index }"
          @click="activeTab = index"
        >
          {{ group.name }}
        </button>
      </div>

      <div class="emoji-content" @wheel.stop>
        <div
          v-for="(group, index) in emojiGroups"
          v-show="activeTab === index"
          :key="index"
          class="emoji-group"
          :class="{
            'emoji-group-image': group.type === 'image',
            'emoji-group-emoticon': group.type === 'emoticon',
          }"
        >
          <button
            v-for="item in group.items"
            :key="item.key"
            class="emoji-item"
            :class="{
              'emoji-image': group.type === 'image',
              'emoji-emoticon': group.type === 'emoticon',
            }"
            :title="item.key"
            @click="selectEmoji(item, group.type)"
          >
            <img v-if="group.type === 'image'" :src="item.val" :alt="item.key" />
            <span v-else>{{ item.val }}</span>
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<style lang="scss" scoped>
@keyframes rotating {
  to {
    transform: rotate(360deg);
  }
}

.rotating {
  animation: rotating 1s linear infinite;
}

.emoji-picker {
  width: 420px;
  height: 240px;
  background: var(--flec-card-bg);
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.emoji-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px 20px;
  color: var(--theme-meta-color);

  i {
    font-size: 2rem;
  }

  span {
    font-size: 0.9rem;
  }
}

.emoji-tabs {
  display: flex;
  border-bottom: 1px solid var(--flec-border-color);
  background: var(--flec-heavy-bg);
  overflow-x: auto;

  &::-webkit-scrollbar {
    height: 2px;
  }
}

.emoji-tab {
  flex-shrink: 0;
  padding: 10px 16px;
  border: none;
  background: transparent;
  color: var(--theme-meta-color);
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s;
  border-bottom: 2px solid transparent;

  &:hover {
    color: var(--font-color);
    background: rgba(0, 0, 0, 0.03);
  }

  &.active {
    color: var(--theme-color);
    border-bottom-color: var(--theme-color);
  }
}

.emoji-content {
  flex: 1;
  overflow-y: auto;
  padding: 12px;

  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-thumb {
    background: var(--flec-border-color);
    border-radius: 3px;
  }
}

.emoji-group {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(36px, 1fr));
  gap: 4px;

  &.emoji-group-image {
    grid-template-columns: repeat(auto-fill, minmax(56px, 1fr));
  }

  &.emoji-group-emoticon {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
  }
}

.emoji-item {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;
  padding: 0;
  flex-shrink: 0;

  span {
    font-size: 1.5rem;
    line-height: 1;
  }

  &.emoji-emoticon {
    width: auto;
    height: auto;
    padding: 6px 10px;

    span {
      font-size: 0.85rem;
      white-space: nowrap;
    }
  }

  &.emoji-image {
    width: 56px;
    height: 56px;

    img {
      width: 48px;
      height: 48px;
      object-fit: contain;
    }
  }

  &:hover {
    background: var(--flec-heavy-bg);
    transform: scale(1.1);
  }
}

@media screen and (max-width: 768px) {
  .emoji-picker {
    width: 100%;
    max-width: 320px;
  }

  .emoji-tab {
    padding: 8px 12px;
    font-size: 0.8rem;
  }

  .emoji-group {
    grid-template-columns: repeat(auto-fill, minmax(32px, 1fr));

    &.emoji-group-image {
      grid-template-columns: repeat(auto-fill, minmax(48px, 1fr));
    }
  }

  .emoji-item {
    width: 32px;
    height: 32px;

    span {
      font-size: 1.3rem;
    }

    &.emoji-emoticon {
      padding: 5px 8px;

      span {
        font-size: 0.75rem;
      }
    }

    &.emoji-image {
      width: 48px;
      height: 48px;

      img {
        width: 40px;
        height: 40px;
      }
    }
  }
}
</style>
