<script setup lang="ts">
import { searchArticles } from '@/composables/api/article';
import type { Article } from '@@/types/article';

const props = defineProps<{ modelValue: boolean }>();
const emit = defineEmits<{ 'update:modelValue': [value: boolean] }>();

const keyword = ref('');
const articles = ref<Article[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = 5;
const loading = ref(false);
const inputRef = ref<HTMLInputElement | null>(null);

const totalPages = computed(() => Math.ceil(total.value / pageSize));
const hasSearched = computed(() => keyword.value.trim().length > 0);

// 高亮关键词
const highlight = (text: string) => {
  const kw = keyword.value.trim();
  if (!kw || !text) return text;

  // 转义正则特殊字符
  const escaped = kw.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
  const regex = new RegExp(`(${escaped})`, 'gi');

  return text.replace(regex, '<mark>$1</mark>');
};

const close = () => {
  emit('update:modelValue', false);
  setTimeout(() => {
    keyword.value = '';
    articles.value = [];
    total.value = 0;
    page.value = 1;
  }, 300);
};

const search = async (newPage = 1) => {
  const searchTerm = keyword.value.trim();
  if (!searchTerm) {
    articles.value = [];
    total.value = 0;
    return;
  }

  loading.value = true;
  page.value = newPage;

  try {
    const data = await searchArticles(searchTerm, {
      page: newPage,
      page_size: pageSize,
    });
    articles.value = data.list;
    total.value = data.total;
  } catch {
    articles.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
};

// 防抖搜索（500ms）
const debouncedSearch = useDebounceFn(() => search(1), 500);

// 监听关键词变化，自动触发搜索
watch(keyword, () => {
  page.value = 1;
  debouncedSearch();
});

const prevPage = () => {
  if (page.value > 1) search(page.value - 1);
};

const nextPage = () => {
  if (page.value < totalPages.value) search(page.value + 1);
};

watch(
  () => props.modelValue,
  async open => {
    if (open) {
      await nextTick();
      inputRef.value?.focus();
    }
  }
);
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="modelValue" class="modal" @click="close">
        <div class="box" @click.stop>
          <!-- 标题栏 -->
          <div class="header">
            <span class="title">搜索</span>
            <button class="close" @click="close"><i class="ri-close-line" /></button>
          </div>

          <!-- 搜索栏 -->
          <div class="search">
            <input
              ref="inputRef"
              v-model="keyword"
              placeholder="输入关键词搜索..."
              @keyup.esc="close"
            />
            <i v-if="loading" class="ri-loader-4-line spin loading-icon" />
          </div>

          <!-- 结果区域 -->
          <div v-if="hasSearched" class="results">
            <!-- 加载中 -->
            <div v-if="loading" class="loading">
              <i class="ri-loader-4-line spin" />
            </div>
            <!-- 有结果 -->
            <template v-else-if="articles.length > 0">
              <NuxtLink
                v-for="item in articles"
                :key="item.id"
                :to="item.url"
                class="item"
                @click="close"
              >
                <NuxtImg v-if="item.cover" :src="item.cover" :alt="item.title" loading="lazy" />
                <div class="info">
                  <!-- eslint-disable-next-line vue/no-v-html -- 高亮关键词已转义特殊字符 -->
                  <h3 v-html="highlight(item.title)" />
                  <!-- eslint-disable-next-line vue/no-v-html -- 高亮关键词已转义特殊字符 -->
                  <p v-if="item.excerpt" class="excerpt" v-html="highlight(item.excerpt)" />
                  <div class="meta">
                    <span>{{ formatDate(item.publish_time) }}</span>
                    <span v-if="item.category">{{ item.category.name }}</span>
                  </div>
                </div>
              </NuxtLink>

              <!-- 分页 -->
              <div v-if="totalPages > 1" class="pagination">
                <button :disabled="page <= 1" @click="prevPage">
                  <i class="ri-arrow-left-s-line" />
                </button>
                <span>{{ page }} / {{ totalPages }}</span>
                <button :disabled="page >= totalPages" @click="nextPage">
                  <i class="ri-arrow-right-s-line" />
                </button>
              </div>
            </template>
            <!-- 无结果 -->
            <div v-else class="empty">
              <i class="ri-search-line" />
              <p>未找到相关文章</p>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style lang="scss" scoped>
.modal {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: flex-start;
  padding: 100px 20px;
  z-index: 10000;
  overflow-y: auto;
}

.box {
  background: var(--flec-card-bg);
  border-radius: 12px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  width: 100%;
  max-width: 600px;
  max-height: calc(100vh - 120px);
  padding: 20px;
  margin: 0 auto;
  position: relative;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;

  .title {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--font-color);
  }
}

.close {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: transparent;
  border: none;
  color: var(--theme-meta-color);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;

  &:hover {
    background: var(--flec-heavy-bg);
    transform: rotate(90deg);
  }
}

.search {
  position: relative;
  display: flex;
  align-items: center;
  padding: 8px 12px;
  background: var(--flec-page-bg);
  border: 2px solid transparent;
  border-radius: 8px;
  transition: border-color 0.3s;

  &:focus-within {
    border-color: var(--theme-color);
  }

  input {
    flex: 1;
    border: none;
    background: none;
    outline: none;
    font-size: 0.95rem;
    color: var(--font-color);
    padding-right: 30px;

    &::placeholder {
      color: var(--theme-meta-color);
    }
  }

  .loading-icon {
    position: absolute;
    right: 12px;
    font-size: 1.1rem;
    color: var(--theme-color);
  }
}

.results {
  margin-top: 16px;
  overflow-y: auto;
  flex: 1;
}

.loading,
.empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: var(--theme-meta-color);

  i {
    font-size: 2rem;
    margin-bottom: 8px;
    opacity: 0.5;
  }

  p {
    margin: 0;
    font-size: 0.9rem;
  }
}

.item {
  display: flex;
  gap: 10px;
  padding: 10px;
  background: var(--flec-page-bg);
  border-radius: 6px;
  margin-bottom: 8px;
  transition: background 0.2s;
  color: inherit;
  text-decoration: none;

  &:hover {
    background: var(--flec-heavy-bg);
  }

  img {
    width: 60px;
    height: 60px;
    border-radius: 4px;
    object-fit: cover;
    flex-shrink: 0;
  }

  .info {
    flex: 1;
    min-width: 0;
  }

  h3 {
    font-size: 0.9rem;
    font-weight: 600;
    margin: 0 0 4px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: var(--font-color);
  }

  .excerpt {
    font-size: 0.8rem;
    color: var(--theme-meta-color);
    margin: 0 0 6px;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    line-height: 1.4;
  }

  h3,
  .excerpt {
    :deep(mark) {
      background: rgba(255, 235, 59, 0.4);
      color: inherit;
      padding: 0 2px;
      border-radius: 2px;
    }
  }

  .meta {
    display: flex;
    gap: 8px;
    font-size: 0.75rem;
    color: var(--theme-meta-color);

    span {
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 12px 0 4px;
  font-size: 0.85rem;
  color: var(--font-color);

  button {
    width: 32px;
    height: 32px;
    border-radius: 6px;
    border: none;
    background: var(--flec-page-bg);
    color: var(--font-color);
    cursor: pointer;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover:not(:disabled) {
      background: var(--flec-btn-hover);
      color: #fff;
    }

    &:disabled {
      opacity: 0.3;
      cursor: not-allowed;
    }
  }
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

@media screen and (max-width: 768px) {
  .modal {
    padding: 80px 12px;
  }

  .box {
    padding: 16px;
  }

  .close {
    top: 10px;
    right: 10px;
  }

  .item {
    img {
      width: 50px;
      height: 50px;
    }

    h3 {
      font-size: 0.85rem;
    }

    .meta {
      font-size: 0.7rem;
      flex-direction: column;
      gap: 2px;
    }
  }
}
</style>
