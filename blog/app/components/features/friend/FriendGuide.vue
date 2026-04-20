<script lang="ts" setup>
import { applyFriend } from '@/composables/api/friend';
import type { FriendApplyRequest } from '@@/types/friend';
const activeTab = ref<'yaml' | 'json' | 'html' | 'table'>('yaml');
const activeApplyTab = ref<'form' | 'format1' | 'format2'>('format1');
const { success } = useToast();
const { blogConfig, basicConfig } = useSysConfig();
const isLoggedIn = useAuth();
const { open: openLogin } = useLoginModal();

const contactEmail = computed(() => basicConfig.value.author_email || '');

const siteConfig = computed(() => ({
  name: blogConfig.value.title,
  link: import.meta.client ? window.location.origin : '',
  avatar: blogConfig.value.favicon,
  description: blogConfig.value.slogan,
  screenshot: blogConfig.value.screenshot,
}));

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text);
    success('已复制到剪贴板！');
  } catch (err) {
    console.error('复制失败:', err);
  }
};

// 站点信息模板 - 使用配置对象动态生成（仅包含非空字段）
const siteTemplates = computed(() => {
  const { name, link, avatar, description, screenshot } = siteConfig.value;

  // YAML 格式 - 动态生成
  const yamlLines = [];
  if (name) yamlLines.push(`- name: ${name}`);
  if (link) yamlLines.push(`  link: ${link}`);
  if (avatar) yamlLines.push(`  avatar: ${avatar}`);
  if (description) yamlLines.push(`  descr: ${description}`);
  if (screenshot) yamlLines.push(`  screenshot: ${screenshot}`);

  // JSON 格式 - 动态生成
  const jsonObj: Record<string, string> = {};
  if (name) jsonObj.title = name;
  if (avatar) jsonObj.avatar = avatar;
  if (link) jsonObj.url = link;
  if (description) jsonObj.description = description;
  if (screenshot) jsonObj.screenshot = screenshot;

  return {
    yaml: yamlLines.join('\n'),
    json: JSON.stringify(jsonObj, null, 4),
    html: `<a href="${link}">${name}</a>`,
  };
});

// 申请模板
const applyTemplates = {
  yaml: `\`\`\`yml
- name: 
  link: 
  avatar: 
  descr: 
  screenshot: 
\`\`\``,
  text: `\`\`\`text
网站名称：
网站地址：
头像图片链接：
站点描述：
网站截图（可选）：
\`\`\``,
};

// 复制站点信息
const copySiteInfo = (type: 'yaml' | 'json' | 'html') => copyToClipboard(siteTemplates.value[type]);

// 快速提交到评论框
const quickSubmit = (format: 'yaml' | 'text') => fillComment(applyTemplates[format]);

// 表单申请相关状态
const showForm = ref(false);
const formSubmitting = ref(false);
const formData = ref<FriendApplyRequest>({
  name: '',
  url: '',
  description: '',
  avatar: '',
  screenshot: '',
});

// 显示申请表单
const showApplyForm = () => {
  showForm.value = true;
};

// 提交申请表单
const submitApplyForm = async () => {
  try {
    formSubmitting.value = true;
    await applyFriend(formData.value);
    success('友链申请已提交！感谢您的申请，我们会尽快处理。');

    // 关闭弹窗（表单会被watch自动重置）
    showForm.value = false;
  } catch (error: unknown) {
    console.error('申请失败:', error);
    // 这里应该用error toast，暂时用alert
    alert((error as Error).message || '申请失败，请稍后重试');
  } finally {
    formSubmitting.value = false;
  }
};

// 监听弹窗关闭，重置表单数据
watch(showForm, newValue => {
  if (!newValue) {
    // 弹窗关闭时重置表单数据
    formData.value = {
      name: '',
      url: '',
      description: '',
      avatar: '',
      screenshot: '',
    };
  }
});
</script>

<template>
  <!-- 免责声明 -->
  <details class="fold">
    <summary>友链页免责声明</summary>
    <div class="fold-content">
      <p>
        本博客遵守中华人民共和国相关法律。本页内容仅作为方便学习而产生的快速链接的链接方式，对与友情链接中存在的链接、好文推荐链接等均为其他网站。我本人能力有限无法逐个甄别每篇文章的每个字，并无法获知是否在收录后原作者是否对链接增加了违反法律甚至其他破坏用户计算机等行为。因为部分友链网站甚至没有做备案、域名并未做实名认证等，所以友链网站均可能存在风险，请你须知。
      </p>
      <p>所以在我力所能及的情况下，我会包括但不限于：</p>
      <ol>
        <li>针对收录的博客中的绝大多数内容通过标题来鉴别是否存在有风险的内容</li>
        <li>在收录的友链好文推荐中检查是否存在风险内容</li>
      </ol>
      <p>但是你在访问的时候，仍然无法避免，包括但不限于：</p>
      <ol>
        <li>作者更换了超链接的指向，替换成了其他内容</li>
        <li>作者的服务器被恶意攻击、劫持、被注入恶意内容</li>
        <li>作者的域名到期，被不法分子用作他用</li>
        <li>作者修改了文章内容，增加钓鱼网站、广告等无效信息</li>
        <li>不完善的隐私保护对用户的隐私造成了侵害、泄漏</li>
      </ol>
      <p>
        如果因为从本页跳转给你造成了损失，深表歉意，并且建议用户如果发现存在问题在本页面进行回复。通常会很快处理。如果长时间无法得到处理，建议联系
        <a v-if="contactEmail" :href="`mailto:${contactEmail}`" class="highlight">{{
          contactEmail
        }}</a>
        <span v-else>网站管理员</span>。
      </p>
    </div>
  </details>

  <!-- 互链规则 -->
  <h3>互链规则</h3>
  <ol>
    <li>
      确保贵站是 <span class="highlight">博客网站</span> 且
      <span class="highlight">原创文章在5篇以上</span>。
    </li>
    <li>为了友链相关页面的统一性和美观性，可能会对你的昵称进行缩短处理。</li>
    <li>如果贵站使用的是 <span class="highlight">免费</span> 域名，将视站点质量进行添加。</li>
    <li>
      若发现
      <span class="highlight">站点三个月以上不进行维护或存在违规内容将会删除友链</span
      >，若网站恢复正常可联系我重新添加友链。
    </li>
  </ol>

  <!-- 添加本站 -->
  <h3>添加本站</h3>

  <div class="tabs">
    <div class="tab-buttons">
      <button :class="{ active: activeTab === 'yaml' }" @click="activeTab = 'yaml'">YAML</button>
      <button :class="{ active: activeTab === 'json' }" @click="activeTab = 'json'">JSON</button>
      <button :class="{ active: activeTab === 'html' }" @click="activeTab = 'html'">HTML</button>
      <button :class="{ active: activeTab === 'table' }" @click="activeTab = 'table'">通用</button>
    </div>

    <div class="tab-content">
      <div v-show="activeTab === 'yaml'" class="code-block">
        <button class="copy-btn" @click="copySiteInfo('yaml')">复制</button>
        <pre><code>{{ siteTemplates.yaml }}</code></pre>
      </div>

      <div v-show="activeTab === 'json'" class="code-block">
        <button class="copy-btn" @click="copySiteInfo('json')">复制</button>
        <pre><code>{{ siteTemplates.json }}</code></pre>
      </div>

      <div v-show="activeTab === 'html'" class="code-block">
        <button class="copy-btn" @click="copySiteInfo('html')">复制</button>
        <pre><code>&lt;a href="{{ siteConfig.link }}"&gt;{{ siteConfig.name }}&lt;/a&gt;</code></pre>
      </div>

      <div v-show="activeTab === 'table'" class="table-content">
        <table>
          <tbody>
            <tr v-if="siteConfig.name">
              <td>站点名称</td>
              <td>{{ siteConfig.name }}</td>
            </tr>
            <tr v-if="siteConfig.link">
              <td>站点链接</td>
              <td>{{ siteConfig.link }}</td>
            </tr>
            <tr v-if="siteConfig.avatar">
              <td>站点图标</td>
              <td>{{ siteConfig.avatar }}</td>
            </tr>
            <tr v-if="siteConfig.description">
              <td>站点描述</td>
              <td>{{ siteConfig.description }}</td>
            </tr>
            <tr v-if="siteConfig.screenshot">
              <td>站点截图</td>
              <td>{{ siteConfig.screenshot }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <!-- 快速申请 -->
  <h3>快速申请</h3>

  <div class="tabs">
    <div class="tab-buttons">
      <button :class="{ active: activeApplyTab === 'format1' }" @click="activeApplyTab = 'format1'">
        格式一
      </button>
      <button :class="{ active: activeApplyTab === 'format2' }" @click="activeApplyTab = 'format2'">
        格式二
      </button>
      <button :class="{ active: activeApplyTab === 'form' }" @click="activeApplyTab = 'form'">
        表单申请
      </button>
    </div>

    <div class="tab-content">
      <div v-show="activeApplyTab === 'format1'" class="apply-format">
        <div class="code-block">
          <pre><code>- name: #网站名称（不易过长）
  link: #网站地址（要求博客地址，请勿提交个人主页）
  avatar: #头像图片链接（请提供尽可能清晰的图片，我会上传到我自己的图床）
  descr: #网站描述
  screenshot: #网站截图（可选）</code></pre>
        </div>
        <button class="submit-btn" @click="quickSubmit('yaml')">快速提交</button>
      </div>

      <div v-show="activeApplyTab === 'format2'" class="apply-format">
        <div class="code-block">
          <pre><code>网站名称（不易过长）：
网站地址（要求博客地址，请勿提交个人主页）：
头像图片链接（请提供尽可能清晰的图片，我会上传到我自己的图床）：
站点描述：
网站截图（可选）：</code></pre>
        </div>
        <button class="submit-btn" @click="quickSubmit('text')">快速提交</button>
      </div>

      <div v-show="activeApplyTab === 'form'" class="apply-format">
        <div class="form-apply-content">
          <div>
            <p class="title">推荐使用表单申请</p>
            <p>
              使用表单申请更加便捷！只需填写您网站的基本信息，系统会收到您的申请并通知管理员，减少信息填写错误，管理员会尽快处理您的申请。
            </p>
            <p>友链申请需要<b>登录后</b>才能提交，确保申请的真实性和可追踪性。</p>
          </div>
          <button v-if="isLoggedIn" class="submit-btn" @click="showApplyForm">表单提交</button>
          <button v-else class="submit-btn" @click="openLogin">注册登录</button>
        </div>
      </div>
    </div>
  </div>

  <!-- 友链申请表单弹窗 -->
  <UiBaseDialog
    v-model="showForm"
    title="友链申请"
    confirm-text="提交申请"
    :loading="formSubmitting"
    @confirm="submitApplyForm"
  >
    <form @submit.prevent="submitApplyForm">
      <div class="form-group">
        <label for="site-name">网站名称 *</label>
        <input
          id="site-name"
          v-model="formData.name"
          type="text"
          placeholder="请输入网站名称"
          required
          maxlength="50"
        />
      </div>

      <div class="form-group">
        <label for="site-url">网站地址 *</label>
        <input
          id="site-url"
          v-model="formData.url"
          type="url"
          placeholder="https://example.com"
          required
        />
      </div>

      <div class="form-group">
        <label for="site-description">网站描述 *</label>
        <textarea
          id="site-description"
          v-model="formData.description"
          placeholder="请简单描述您的网站"
          required
          maxlength="500"
          rows="3"
        />
      </div>

      <div class="form-group">
        <label for="site-avatar">网站头像/Logo *</label>
        <input
          id="site-avatar"
          v-model="formData.avatar"
          type="url"
          placeholder="请提供网站头像或Logo的链接地址"
          required
        />
      </div>

      <div class="form-group">
        <label for="site-screenshot">网站截图</label>
        <input
          id="site-screenshot"
          v-model="formData.screenshot"
          type="url"
          placeholder="可选：提供网站截图链接"
        />
      </div>
    </form>
  </UiBaseDialog>
</template>

<style lang="scss" scoped>
@use '@/assets/css/mixins' as *;

// 标题
h3 {
  font-size: 1.5rem;
  margin-top: 2rem;
  margin-bottom: 0.75rem;
  font-weight: 600;
  line-height: 1.3;
  scroll-margin-top: 80px; // 锚点跳转偏移
  position: relative;
}

// 有序列表
ol {
  list-style-type: decimal;
}

// 高亮文本
.highlight {
  color: var(--theme-color);
  font-weight: 500;
}

// 折叠面板
.fold {
  @extend .cardHover;
  margin-bottom: 25px;
  overflow: hidden;

  summary {
    padding: 8px 16px;
    font-weight: 700;
    color: var(--font-color);
    cursor: pointer;
  }

  .fold-content {
    padding: 0 16px 16px 16px;
  }
}

// 选项卡
.tabs {
  @extend .cardHover;
  padding: 12px;

  .tab-buttons {
    display: flex;
    gap: 8px;
    margin-bottom: 12px;
    padding-bottom: 12px;
    border-bottom: 1px solid var(--flec-heavy-bg);

    button {
      padding: 6px 14px;
      background: transparent;
      border: none;
      border-radius: 4px;
      color: var(--theme-meta-color);
      font-size: 0.9rem;
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        background: var(--flec-heavy-bg);
      }

      &.active {
        color: var(--font-light-color);
        background: var(--theme-color);
      }
    }
  }

  .tab-content {
    .code-block {
      position: relative;
      background: var(--flec-heavy-bg);
      border-radius: 4px;

      .copy-btn {
        position: absolute;
        top: 8px;
        right: 8px;
        padding: 4px 10px;
        background: var(--flec-card-bg);
        border-radius: 4px;
        color: var(--theme-meta-color);
        font-size: 0.85rem;
        cursor: pointer;
        transition: opacity 0.2s;

        &:hover {
          opacity: 0.8;
        }
      }

      pre {
        margin: 0;
        padding: 16px;
        padding-top: 12px;
        overflow-x: auto;
        scrollbar-width: none; // Firefox
        -ms-overflow-style: none; // IE 10+

        &::-webkit-scrollbar {
          display: none; // Chrome, Safari, Edge
        }

        code {
          color: var(--font-color);
          font-family: 'Consolas', 'Monaco', monospace;
          font-size: 0.85rem;
          line-height: 1.6;
        }
      }
    }

    .table-content {
      background: var(--flec-heavy-bg);
      border-radius: 4px;
      overflow: hidden;

      table {
        width: 100%;
        border-collapse: collapse;

        tr {
          &:not(:last-child) td {
            border-bottom: 1px solid var(--flec-card-bg);
          }
        }

        td {
          padding: 10px 12px;
          font-size: 0.9rem;

          &:first-child {
            width: 100px;
            font-weight: 600;
            color: var(--font-color);
          }

          &:last-child {
            color: var(--theme-meta-color);
            word-break: break-all;
          }
        }
      }
    }
  }
}

// 快速申请
.apply-format {
  .code-block {
    margin-bottom: 16px;
  }

  .submit-btn {
    padding: 10px 24px;
    background: var(--theme-color);
    border: none;
    border-radius: 4px;
    color: var(--font-light-color);
    font-size: 0.95rem;
    font-weight: 500;
    cursor: pointer;
    transition: opacity 0.2s;

    &:hover {
      opacity: 0.9;
    }
  }
}

// 响应式设计
@media screen and (max-width: 1024px) {
  h3 {
    font-size: 1.35rem;
    margin-top: 1.75rem;
  }

  ol {
    font-size: 0.95rem;

    li {
      margin: 0.5rem 0;
    }
  }

  .tabs {
    padding: 11px;

    .tab-buttons {
      button {
        padding: 5.5px 13px;
        font-size: 0.875rem;
      }
    }

    .tab-content {
      .code-block {
        pre {
          padding: 14px;
          scrollbar-width: none;
          -ms-overflow-style: none;

          &::-webkit-scrollbar {
            display: none;
          }

          code {
            font-size: 0.825rem;
          }
        }

        .copy-btn {
          font-size: 0.825rem;
          padding: 4px 9px;
        }
      }

      .table-content {
        table {
          td {
            padding: 9px 11px;
            font-size: 0.875rem;
          }
        }
      }
    }
  }

  .apply-format {
    .submit-btn {
      padding: 9px 22px;
      font-size: 0.925rem;
    }
  }
}

@media screen and (max-width: 768px) {
  h3 {
    font-size: 1.25rem;
    margin-top: 1.5rem;
    margin-bottom: 0.65rem;
  }

  ol {
    font-size: 0.9rem;
    padding-left: 1.5rem;

    li {
      margin: 0.4rem 0;
    }
  }

  .highlight {
    font-size: 0.9rem;
  }

  .fold {
    margin-bottom: 20px;

    summary {
      padding: 7px 14px;
      font-size: 0.95rem;
    }

    .fold-content {
      padding: 0 14px 14px 14px;
      font-size: 0.9rem;

      p,
      ol,
      li {
        font-size: 0.9rem;
      }
    }
  }

  .tabs {
    padding: 10px;

    .tab-buttons {
      flex-wrap: wrap;
      gap: 6px;
      padding-bottom: 10px;

      button {
        padding: 5px 12px;
        font-size: 0.85rem;
      }
    }

    .tab-content {
      .code-block {
        pre {
          padding: 12px;
          padding-top: 10px;
          scrollbar-width: none;
          -ms-overflow-style: none;

          &::-webkit-scrollbar {
            display: none;
          }

          code {
            font-size: 0.8rem;
            line-height: 1.5;
          }
        }

        .copy-btn {
          font-size: 0.8rem;
          padding: 3.5px 8px;
          top: 6px;
          right: 6px;
        }
      }

      .table-content {
        table {
          td {
            padding: 8px 10px;
            font-size: 0.85rem;

            &:first-child {
              width: 80px;
              font-size: 0.85rem;
            }

            &:last-child {
              font-size: 0.82rem;
            }
          }
        }
      }
    }
  }

  .apply-format {
    .code-block {
      margin-bottom: 14px;
    }

    .submit-btn {
      padding: 8px 20px;
      font-size: 0.9rem;
    }
  }
}

// 表单申请内容样式
.form-apply-content {
  .title {
    font-size: 1.3rem;
    font-weight: 600;
  }

  p {
    color: var(--text-color);
    margin-bottom: 16px;
    font-size: 1rem;
    line-height: 1.6;
  }
}

// 表单样式（用于BaseDialog内部）

.form-group {
  margin-bottom: 20px;

  label {
    display: block;
    margin-bottom: 8px;
    color: var(--font-color);
    font-weight: 500;
    font-size: 0.95rem;
  }

  input,
  textarea {
    width: 100%;
    padding: 12px;
    border: 1px solid var(--flec-border);
    border-radius: 6px;
    background: var(--flec-card-bg);
    color: var(--font-color);
    font-size: 0.95rem;
    transition: border-color 0.2s ease;
    box-sizing: border-box;

    &:focus {
      outline: none;
      border-color: var(--theme-color);
      box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.1);
    }

    &::placeholder {
      color: var(--theme-meta-color);
    }
  }

  textarea {
    resize: vertical;
    min-height: 80px;
    font-family: inherit;
  }
}
</style>
