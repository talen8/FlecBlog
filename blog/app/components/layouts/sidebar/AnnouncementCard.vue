<script setup lang="ts">
import DOMPurify from 'isomorphic-dompurify';

const ALLOWED_TAGS = [
  'a',
  'br',
  'strong',
  'b',
  'em',
  'i',
  'u',
  's',
  'del',
  'ins',
  'mark',
  'span',
  'sub',
  'sup',
  'p',
];
const ALLOWED_ATTR = ['style', 'href', 'target'];
const DANGEROUS_STYLE_RE = /(?:expression\s*\(|url\s*\(|javascript:|vbscript:|data:)/i;
const COLOR_VALUE_RE = /^(?:#[0-9a-f]{3,8}|(?:rgb|hsl)a?\(\s*[-\d.%\s,]+\)|[a-z]+)$/i;

const isAllowedStyleValue = (property: string, value: string) => {
  if (!value || DANGEROUS_STYLE_RE.test(value)) return false;
  return property === 'color' && COLOR_VALUE_RE.test(value);
};

const sanitizeInlineStyle = (styleText: string) =>
  styleText
    .split(';')
    .map(rule => rule.trim())
    .filter(Boolean)
    .map(rule => {
      const [property, value] = rule.split(':').map(s => s.trim());
      return !property || !value || !isAllowedStyleValue(property.toLowerCase(), value)
        ? ''
        : `${property.toLowerCase()}: ${value}`;
    })
    .filter(Boolean)
    .join('; ');

const renderAnnouncementHtml = (content: string) => {
  if (!content) return '';
  const sanitized = DOMPurify.sanitize(content, {
    ALLOWED_TAGS,
    ALLOWED_ATTR,
    ALLOW_DATA_ATTR: false,
  });
  return sanitized.replace(
    /\sstyle=(['"])(.*?)\1/gi,
    (_match, quote: string, styleText: string) => {
      const safeStyle = sanitizeInlineStyle(styleText);
      return safeStyle ? ` style=${quote}${safeStyle}${quote}` : '';
    }
  );
};

const hasVisibleAnnouncementContent = (html: string) =>
  html
    .replace(/<br\s*\/?>/gi, '\n')
    .replace(/<[^>]+>/g, '')
    .replace(/&nbsp;/gi, ' ')
    .trim().length > 0;

const { blogConfig } = useSysConfig();

const announcementHtml = computed(() =>
  renderAnnouncementHtml(blogConfig.value.announcement?.trim() || '')
);

const showAnnouncement = computed(() => hasVisibleAnnouncementContent(announcementHtml.value));
</script>

<template>
  <div v-if="showAnnouncement" class="card-widget card-announcement">
    <div class="item-headline">
      <i class="ri-megaphone-line announcement-icon" />
      <span>公告</span>
    </div>
    <!-- eslint-disable-next-line vue/no-v-html -- 公告内容来自后台配置，为可信内容 -->
    <div class="announcement-content" v-html="announcementHtml" />
  </div>
</template>

<style lang="scss" scoped>
.card-announcement {
  .announcement-icon {
    color: #e03131;
  }

  .announcement-content {
    line-height: 1.8;
    color: var(--font-color);
    opacity: 0.85;
    white-space: pre-line;
    word-break: break-word;

    :deep(p) {
      margin: 0;
    }

    :deep(p + p) {
      margin-top: 0.45rem;
    }

    :deep(a) {
      color: var(--theme-color);
      text-decoration: underline;
      text-underline-offset: 2px;
    }
  }
}
</style>
