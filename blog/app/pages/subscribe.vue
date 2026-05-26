<script lang="ts" setup>
import type { ApiResponse } from '@@/types/request';

definePageMeta({
  showSidebar: false,
});

useSeoMeta({
  title: '订阅本站',
  description: '通过公众号、邮件或RSS订阅本站，第一时间获取最新文章更新',
});

const route = useRoute();
const {
  public: { apiUrl },
} = useRuntimeConfig();
const { success, error: errorToast } = useToast();
const { blogConfig } = useSysConfig();
const { copy: copyToClipboard } = useClipboard();

// 公众号订阅
const showWechatDialog = ref(false);

// RSS订阅
const showRSSDialog = ref(false);
const rss2Url = computed(() => `${import.meta.client ? window.location.origin : ''}/rss.xml`);
const atomUrl = computed(() => `${import.meta.client ? window.location.origin : ''}/atom.xml`);

// 订阅
const showSubscribeDialog = ref(false);
const email = ref('');
const subscribeLoading = ref(false);

// 退订
const showUnsubscribeDialog = ref(false);
const unsubscribeLoading = ref(false);
const unsubscribeMessage = ref('');
const unsubscribeSuccess = ref(false);

// 检查退订操作
onMounted(() => {
  const { action, token } = route.query;
  if (action === 'unsubscribe' && token) {
    showUnsubscribeDialog.value = true;
    handleUnsubscribe(token as string);
  }
});

const handleCopy = async (text: string) => {
  try {
    await copyToClipboard(text);
    success('已复制到剪贴板');
  } catch {
    errorToast('复制失败');
  }
};

const handleSubscribe = async () => {
  if (!email.value) return;

  subscribeLoading.value = true;
  try {
    const res = await $fetch<ApiResponse>(`${apiUrl}/subscribe`, {
      method: 'POST',
      body: { email: email.value },
    });

    if (res.code === 0) {
      showSubscribeDialog.value = false;
      email.value = '';
      success('订阅成功！');
    } else {
      errorToast(res.message || '订阅失败');
    }
  } catch (error: unknown) {
    const err = error as Error & { data?: { message?: string } };
    errorToast(err.data?.message || '订阅失败');
  } finally {
    subscribeLoading.value = false;
  }
};

// 退订
const handleUnsubscribe = async (token: string) => {
  unsubscribeLoading.value = true;
  try {
    const res = await $fetch<ApiResponse>(`${apiUrl}/subscribe/unsubscribe?token=${token}`);
    unsubscribeSuccess.value = res.code === 0;
    unsubscribeMessage.value = res.code === 0 ? '退订成功！' : res.message || '退订失败';
  } catch (error: unknown) {
    const err = error as Error & { data?: { message?: string } };
    unsubscribeSuccess.value = false;
    unsubscribeMessage.value = err.data?.message || '退订失败';
  } finally {
    unsubscribeLoading.value = false;
  }
};

const closeUnsubscribeDialog = () => {
  showUnsubscribeDialog.value = false;
  navigateTo('/subscribe', { replace: true });
};
</script>

<template>
  <div id="subscribe-page">
    <h1 class="page-title">订阅本站</h1>
    <div class="page-subtitle">选择您喜欢的订阅方式，随时获取最新更新</div>

    <!-- 订阅方式卡片 -->
    <div class="subscribe-list">
      <!-- 公众号订阅 -->
      <button
        type="button"
        class="subscribe-item subscribe-wechat"
        title="公众号"
        @click="showWechatDialog = true"
      >
        <div class="subscribe-description">推送精选文章<br />推送全文</div>
        <div class="subscribe-info-group">
          <div class="subscribe-title">公众号订阅</div>
          <div class="subscribe-info">推荐的订阅方式</div>
          <i class="ri-wechat-fill subscribe-icon" />
        </div>
      </button>

      <!-- 邮件订阅 -->
      <button
        type="button"
        class="subscribe-item subscribe-mail"
        title="邮件订阅"
        @click="
          showSubscribeDialog = true;
          email = '';
        "
      >
        <div class="subscribe-description">推送全部文章<br />推送简介</div>
        <div class="subscribe-info-group">
          <div class="subscribe-title">邮件订阅</div>
          <div class="subscribe-info">推荐的订阅方式</div>
          <i class="ri-mail-fill subscribe-icon" />
        </div>
      </button>

      <!-- RSS 订阅 -->
      <button
        type="button"
        class="subscribe-item subscribe-rss"
        title="RSS"
        @click="showRSSDialog = true"
      >
        <div class="subscribe-description">推送全部文章<br />推送简介</div>
        <div class="subscribe-info-group">
          <div class="subscribe-title">RSS</div>
          <div class="subscribe-info">备用订阅方式</div>
          <i class="ri-rss-fill subscribe-icon" />
        </div>
      </button>
    </div>

    <!-- 公众号二维码弹窗 -->
    <UiBaseDialog
      v-model="showWechatDialog"
      :title="blogConfig.wechat_qrcode ? '关注公众号' : '公众号订阅'"
      confirm-text=""
    >
      <div class="wechat-dialog-content">
        <template v-if="blogConfig.wechat_qrcode">
          <p class="wechat-desc">
            扫描二维码关注公众号<br />
            获取最新文章推送
          </p>
          <div class="qrcode-wrapper">
            <NuxtImg
              :src="blogConfig.wechat_qrcode"
              :alt="blogConfig.wechat_name || '公众号二维码'"
              class="qrcode-image"
              loading="lazy"
            />
          </div>
          <p v-if="blogConfig.wechat_name" class="wechat-name">
            {{ blogConfig.wechat_name }}
          </p>
        </template>
        <template v-else>
          <p class="wechat-desc">
            公众号订阅功能暂未配置<br />
            建议您使用邮件订阅或 RSS 订阅获取最新文章
          </p>
        </template>
      </div>
    </UiBaseDialog>

    <!-- RSS订阅弹窗 -->
    <UiBaseDialog v-model="showRSSDialog" title="RSS 订阅" confirm-text="">
      <div class="rss-dialog-content">
        <p class="rss-desc">
          复制下方地址到您的 RSS 阅读器<br />
          或使用以下客户端快速订阅
        </p>

        <!-- RSS地址复制区域 -->
        <div class="rss-url-item">
          <span class="url-label">RSS 2.0</span>
          <div class="url-copy-box">
            <code class="url-text">{{ rss2Url }}</code>
            <button class="copy-btn" @click="handleCopy(rss2Url)">复制</button>
          </div>
        </div>
        <div class="rss-url-item">
          <span class="url-label">Atom</span>
          <div class="url-copy-box">
            <code class="url-text">{{ atomUrl }}</code>
            <button class="copy-btn" @click="handleCopy(atomUrl)">复制</button>
          </div>
        </div>

        <!-- 客户端订阅按钮 -->
        <div class="rss-url-item">
          <span class="url-label">Follow</span>
          <div class="url-copy-box">
            <code class="url-text">follow://add?url={{ atomUrl }}</code>
            <a
              :href="`follow://add?url=${encodeURIComponent(atomUrl)}`"
              target="_blank"
              class="copy-btn"
            >
              订阅
            </a>
          </div>
        </div>
      </div>
    </UiBaseDialog>

    <!-- 订阅弹窗 -->
    <UiBaseDialog
      v-model="showSubscribeDialog"
      title="邮件订阅"
      confirm-text="确认订阅"
      :loading="subscribeLoading"
      @confirm="handleSubscribe"
    >
      <div class="dialog-content">
        <p class="dialog-desc">
          订阅后将收到本站最新文章推送<br />
          可随时通过邮件中的退订链接取消订阅
        </p>

        <div class="input-group">
          <label for="email">邮箱地址</label>
          <input
            id="email"
            v-model="email"
            type="email"
            placeholder="请输入您的邮箱地址"
            :disabled="subscribeLoading"
            @keyup.enter="handleSubscribe"
          />
        </div>
      </div>
    </UiBaseDialog>

    <!-- 退订弹窗 -->
    <UiBaseDialog
      v-model="showUnsubscribeDialog"
      title="退订确认"
      confirm-text="确定"
      :close-on-click-outside="!unsubscribeLoading"
      @confirm="closeUnsubscribeDialog"
    >
      <div class="dialog-content unsubscribe-content">
        <div v-if="unsubscribeLoading" class="loading-state">
          <i class="ri-loader-4-line loading-icon" />
          <p>{{ unsubscribeMessage }}</p>
        </div>

        <div v-else class="result-state">
          <i
            :class="
              unsubscribeSuccess
                ? 'ri-checkbox-circle-line success-icon'
                : 'ri-close-circle-line error-icon'
            "
          />
          <h3>{{ unsubscribeMessage }}</h3>
        </div>
      </div>
    </UiBaseDialog>
  </div>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

#subscribe-page {
  @extend .cardHover;
  width: 100%;
  padding: 40px;
}

.page-title {
  margin: 0 0 10px;
  font-weight: bold;
  font-size: 2rem;
}

.page-subtitle {
  margin-bottom: 30px;
  color: var(--theme-meta-color);
  font-size: 1rem;
}

// 订阅卡片列表
.subscribe-list {
  display: flex;
  width: 100%;
  flex-direction: row;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 20px;
}

.subscribe-item {
  appearance: none;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 240px;
  height: 240px;
  overflow: hidden;
  text-decoration: none;
  width: calc(100% / 3 - 8px);
  padding: 0;
  border: none;
  font: inherit;
  text-align: left;
  transition: all 0.4s ease-in-out;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  position: relative;
  cursor: pointer;

  &.subscribe-wechat {
    background: var(--flec-subscribe-wechat);
  }

  &.subscribe-mail {
    background: var(--flec-subscribe-mail);
  }

  &.subscribe-rss {
    background: var(--flec-subscribe-rss);
  }

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15);

    .subscribe-icon {
      transform: translate(-3px, -3px) scale(1.05);
      opacity: 0.6;
      filter: blur(1px);
    }
  }
}

.subscribe-description {
  font-size: 16px;
  color: white;
  margin: 26px 0 0 30px;
  line-height: 1.6;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: relative;
  z-index: 2;
}

.subscribe-info-group {
  position: relative;
  margin: 0 0 26px 30px;
  color: white;
  z-index: 2;
}

.subscribe-title {
  font-size: 36px;
  font-weight: 700;
  line-height: 1;
  margin-bottom: 8px;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.subscribe-info {
  opacity: 0.9;
  font-size: 14px;
}

.subscribe-icon {
  position: absolute;
  bottom: -125px;
  right: -25px;
  font-size: 140px;
  user-select: none;
  transition: all 0.8s cubic-bezier(0.39, 0.575, 0.565, 1);
  transform-origin: bottom right;
  filter: blur(3px);
  opacity: 0.3;
  z-index: 1;
  color: rgba(255, 255, 255, 0.3);
}

.rss-desc,
.wechat-desc,
.dialog-desc {
  color: var(--theme-meta-color);
  line-height: 1.6;
}

// RSS弹窗内容
.rss-dialog-content {
  text-align: center;
  padding: 10px 0;

  .rss-desc {
    margin-bottom: 24px;
  }

  .rss-url-item {
    margin-bottom: 16px;
    text-align: left;

    .url-label {
      font-size: 14px;
      font-weight: 500;
      margin-bottom: 8px;
    }

    .url-copy-box {
      display: flex;
      gap: 8px;
      min-width: 0;

      .url-text {
        flex: 1;
        min-width: 0;
        padding: 10px 12px;
        background: var(--flec-card-bg);
        border: 1px solid var(--flec-border);
        border-radius: 6px;
        font-size: 13px;
        font-family: monospace;
        word-break: break-all;
      }

      .copy-btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        padding: 10px 16px;
        background: var(--theme-color);
        color: white;
        border: none;
        border-radius: 6px;
        font-size: 13px;
        cursor: pointer;
        text-decoration: none;

        &:hover {
          opacity: 0.9;
        }
      }
    }
  }
}

// 公众号弹窗内容
.wechat-dialog-content {
  text-align: center;
  padding: 10px 0;

  .wechat-desc {
    margin-bottom: 20px;
  }

  .qrcode-wrapper {
    display: flex;
    justify-content: center;
    margin-bottom: 16px;

    .qrcode-image {
      width: 200px;
      height: 200px;
      border-radius: 8px;
      object-fit: contain;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
    }
  }

  .wechat-name {
    font-size: 16px;
    font-weight: 500;
    color: var(--font-color);
    margin: 0;
  }
}

// 弹窗内容样式
.dialog-content {
  .dialog-desc {
    margin-bottom: 20px;
    text-align: center;
  }

  .input-group {
    margin-bottom: 16px;

    label {
      display: block;
      margin-bottom: 8px;
      color: var(--font-color);
      font-size: 0.95rem;
      font-weight: 500;
    }

    input {
      width: 100%;
      height: 44px;
      padding: 0 16px;
      border: 2px solid var(--flec-border);
      border-radius: 8px;
      background: var(--flec-card-bg);
      color: var(--font-color);
      font-size: 0.95rem;
      transition: all 0.2s;

      &:focus {
        outline: none;
        border-color: var(--theme-color);
        box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
      }

      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }

      &::placeholder {
        color: var(--theme-meta-color);
      }
    }
  }
}

// 退订弹窗内容
.unsubscribe-content {
  text-align: center;
  padding: 20px 0;

  .loading-state {
    .loading-icon {
      font-size: 48px;
      color: var(--theme-color);
      animation: spin 1s linear infinite;
      margin-bottom: 16px;
    }

    p {
      color: var(--theme-meta-color);
      font-size: 1rem;
    }
  }

  .result-state {
    .success-icon,
    .error-icon {
      font-size: 64px;
      margin-bottom: 16px;
    }

    .success-icon {
      color: #4caf50;
    }

    .error-icon {
      color: #f44336;
    }

    h3 {
      margin: 0;
      font-size: 1.25rem;
      color: var(--font-color);
    }
  }
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

// 响应式设计
@media screen and (max-width: 1024px) {
  #subscribe-page {
    padding: 30px;

    .page-title {
      font-size: 1.75rem;
    }

    .page-subtitle {
      font-size: 0.95rem;
      margin-bottom: 25px;
    }
  }

  .subscribe-item {
    width: calc(50% - 6px);
    height: 220px;

    .subscribe-description {
      font-size: 15px;
      margin: 22px 0 0 25px;
    }

    .subscribe-info-group {
      margin: 0 0 22px 25px;
    }

    .subscribe-title {
      font-size: 32px;
    }

    .subscribe-icon {
      font-size: 120px;
      bottom: -110px;
    }
  }
}

@media screen and (max-width: 768px) {
  #subscribe-page {
    padding: 18px;

    .page-title {
      font-size: 1.4rem;
    }

    .page-subtitle {
      font-size: 0.875rem;
      margin-bottom: 20px;
    }
  }

  .subscribe-list {
    flex-direction: column;
    gap: 10px;
    margin-bottom: 16px;
  }

  .subscribe-item {
    min-width: 0;
    width: 100%;
    height: 180px;

    .subscribe-description {
      font-size: 14px;
      margin: 20px 0 0 20px;
    }

    .subscribe-info-group {
      margin: 0 0 20px 20px;
    }

    .subscribe-title {
      font-size: 28px;
    }

    .subscribe-info {
      font-size: 0.8rem;
    }

    .subscribe-icon {
      font-size: 100px;
      bottom: -90px;
    }
  }

  .rss-dialog-content {
    .rss-url-item {
      .url-copy-box {
        flex-direction: column;

        .url-text,
        .copy-btn {
          width: 100%;
        }
      }
    }
  }

  .wechat-dialog-content {
    .qrcode-wrapper {
      .qrcode-image {
        width: 70vw;
        max-width: 200px;
        height: auto;
        aspect-ratio: 1;
      }
    }
  }
}
</style>
