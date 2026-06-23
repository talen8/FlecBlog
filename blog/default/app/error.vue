<script setup lang="ts">
import type { NuxtError } from '#app';

const props = defineProps<{
  error: NuxtError;
}>();

useSeoMeta({
  title: props.error.status === 404 ? '页面未找到' : `错误 ${props.error.status}`,
  description:
    props.error.status === 404 ? '抱歉，您访问的页面不存在或已被删除' : props.error.message,
});

const handleError = () => {
  clearError({ redirect: '/' });
};

const goBack = () => {
  clearError();
  useRouter().back();
};

const errorInfo = computed(() => {
  if (props.error.status === 404) {
    return {
      code: '404',
      title: '页面未找到',
      description: '抱歉，您访问的页面不存在或已被删除。',
    };
  }

  if (props.error.status === 500) {
    return {
      code: '500',
      title: '服务器错误',
      description: '抱歉，服务器出现了错误，请稍后再试。',
    };
  }

  return {
    code: String(props.error.status),
    title: '发生错误',
    description: props.error.message || '抱歉，发生了未知错误。',
  };
});
</script>

<template>
  <div class="error-page-wrapper">
    <div class="error-card">
      <div class="error-code">{{ errorInfo.code }}</div>
      <h1 class="error-title">{{ errorInfo.title }}</h1>
      <p class="error-description">{{ errorInfo.description }}</p>
      <div class="error-actions">
        <button class="btn btn-secondary" @click="goBack">
          <i class="ri-arrow-left-line" />
          返回上一页
        </button>
        <button class="btn btn-primary" @click="handleError">
          <i class="ri-home-line" />
          返回首页
        </button>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.error-page-wrapper {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fafafa;
  padding: 20px;
}

.error-card {
  background: #fff;
  border-radius: 12px;
  padding: 80px 60px;
  text-align: center;
  max-width: 540px;
  width: 100%;
  box-shadow:
    0 1px 3px rgba(0, 0, 0, 0.05),
    0 4px 12px rgba(0, 0, 0, 0.03);
  border: 1px solid rgba(0, 0, 0, 0.04);
}

.error-code {
  font-size: 7rem;
  font-weight: 300;
  line-height: 1;
  color: #49b1f5;
  margin-bottom: 1.5rem;
  letter-spacing: -0.02em;
}

.error-title {
  font-size: 1.75rem;
  color: #333;
  margin: 0 0 1rem;
  font-weight: 500;
}

.error-description {
  font-size: 1rem;
  color: #888;
  margin: 0 0 2.5rem;
  line-height: 1.7;
}

.error-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: center;
  flex-wrap: wrap;

  .btn {
    display: inline-flex;
    align-items: center;
    gap: 0.4rem;
    padding: 10px 20px;
    border: none;
    border-radius: 6px;
    font-size: 0.9rem;
    cursor: pointer;
    transition: all 0.2s ease;
    font-weight: 500;

    i {
      font-size: 1rem;
    }

    &.btn-primary {
      background: #49b1f5;
      color: white;

      &:hover {
        background: #3da0e3;
      }
    }

    &.btn-secondary {
      background: transparent;
      color: #666;
      border: 1px solid #e0e0e0;

      &:hover {
        border-color: #49b1f5;
        color: #49b1f5;
      }
    }
  }
}

@media screen and (max-width: 768px) {
  .error-card {
    padding: 50px 35px;
  }

  .error-code {
    font-size: 5rem;
  }

  .error-title {
    font-size: 1.5rem;
  }

  .error-description {
    font-size: 0.95rem;
  }

  .error-actions {
    flex-direction: column;

    .btn {
      width: 100%;
      justify-content: center;
    }
  }
}

@media (prefers-color-scheme: dark) {
  .error-page-wrapper {
    background: #0a0a0a;
  }

  .error-card {
    background: #141414;
    border-color: rgba(255, 255, 255, 0.06);
    box-shadow: none;
  }

  .error-code {
    color: #49b1f5;
  }

  .error-title {
    color: #e0e0e0;
  }

  .error-description {
    color: #888;
  }

  .btn-secondary {
    color: #aaa;
    border-color: #333;

    &:hover {
      border-color: #49b1f5;
      color: #49b1f5;
    }
  }
}
</style>

<style lang="scss">
[data-theme='dark'] {
  .error-page-wrapper {
    background: #0a0a0a;
  }

  .error-card {
    background: #141414;
    border-color: rgba(255, 255, 255, 0.06);
    box-shadow: none;
  }

  .error-code {
    color: #49b1f5;
  }

  .error-title {
    color: #e0e0e0;
  }

  .error-description {
    color: #888;
  }

  .btn-secondary {
    color: #aaa;
    border-color: #333;

    &:hover {
      border-color: #49b1f5;
      color: #49b1f5;
    }
  }
}
</style>
