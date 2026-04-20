<template>
  <div
    class="codemirror-editor-wrapper"
    :class="{ 'is-fullscreen': isBrowserFullscreen || isPageFullscreen }"
  >
    <!-- 工具栏 -->
    <div class="editor-toolbar">
      <template v-for="(item, index) in toolbarItems" :key="index">
        <div v-if="item.type === 'divider'" class="toolbar-divider"></div>

        <!-- 弹性空间 -->
        <div v-else-if="item.type === 'spacer'" class="toolbar-spacer"></div>

        <!-- 下载在线图片按钮 -->
        <template v-else-if="item.title === '下载在线图片'">
          <el-popover :width="350" trigger="click" placement="bottom">
            <template #reference>
              <button :title="item.title" class="toolbar-btn">
                <i :class="item.icon"></i>
              </button>
            </template>
            <div style="padding: 8px 0">
              <el-input
                v-model="onlineImageUrl"
                placeholder="输入图片URL，按回车下载"
                size="small"
                clearable
                @keyup.enter="handleOnlineImageDownload"
                style="width: 100%"
              >
                <template #append>
                  <el-button
                    @click="handleOnlineImageDownload"
                    :loading="downloadingImage"
                    :disabled="!onlineImageUrl.trim()"
                    size="small"
                  >
                    下载
                  </el-button>
                </template>
              </el-input>
            </div>
          </el-popover>
        </template>
        <!-- 表情选择器按钮 -->
        <template v-else-if="item.title === '表情'">
          <el-popover
            :width="320"
            trigger="click"
            placement="bottom"
            v-model:visible="emojiState.visible"
            @show="handleEmojiPickerShow"
          >
            <template #reference>
              <button
                :title="item.title"
                class="toolbar-btn"
                :class="{ active: emojiState.visible }"
              >
                <i :class="item.icon"></i>
              </button>
            </template>
            <!-- 表情内容 -->
            <div class="emoji-wrap">
              <div class="emoji-bar">
                <button
                  v-for="(group, index) in emojiState.groups"
                  :key="index"
                  class="emoji-tab"
                  :class="{ active: emojiState.activeTab === index }"
                  @click="emojiState.activeTab = index"
                >
                  {{ group.name }}
                </button>
              </div>
              <div class="emoji-list">
                <div
                  v-for="(group, index) in emojiState.groups"
                  v-show="emojiState.activeTab === index"
                  :key="index"
                  class="emoji-group"
                  :class="{ 'emoji-text': group.type === 'emoticon' }"
                >
                  <button
                    v-for="item in group.items"
                    :key="item.key"
                    class="emoji-btn"
                    :title="item.key"
                    @click="selectEmoji(item, group.type)"
                  >
                    <img v-if="group.type === 'image'" :src="item.val" :alt="item.key" />
                    <span v-else>{{ item.val }}</span>
                  </button>
                </div>
              </div>
            </div>
          </el-popover>
        </template>
        <!-- 提示框按钮 -->
        <template v-else-if="item.title === '提示框'">
          <el-popover
            :width="300"
            trigger="click"
            placement="bottom"
            v-model:visible="noteDialog.visible"
          >
            <template #reference>
              <button
                :title="item.title"
                class="toolbar-btn"
                :class="{ active: noteDialog.visible }"
              >
                <i :class="item.icon"></i>
              </button>
            </template>
            <div class="note-dialog-wrap">
              <div class="note-form-item">
                <el-radio-group v-model="noteDialog.type" size="small">
                  <el-radio-button value="info">Info</el-radio-button>
                  <el-radio-button value="warning">Warning</el-radio-button>
                  <el-radio-button value="success">Success</el-radio-button>
                  <el-radio-button value="error">Error</el-radio-button>
                </el-radio-group>
              </div>
              <div class="note-form-item">
                <el-input
                  v-model="noteDialog.title"
                  placeholder="标题"
                  size="small"
                  clearable
                  @keyup.enter="handleInsertNote"
                />
              </div>
              <div class="note-form-actions">
                <el-button type="primary" size="small" @click="handleInsertNote"> 插入 </el-button>
              </div>
            </div>
          </el-popover>
        </template>
        <!-- 标签页按钮 -->
        <template v-else-if="item.title === '标签页'">
          <el-popover
            :width="320"
            trigger="click"
            placement="bottom"
            v-model:visible="tabsDialog.visible"
          >
            <template #reference>
              <button
                :title="item.title"
                class="toolbar-btn"
                :class="{ active: tabsDialog.visible }"
              >
                <i :class="item.icon"></i>
              </button>
            </template>
            <div class="tabs-dialog-wrap">
              <div class="tabs-list">
                <div v-for="(tab, index) in tabsDialog.tabs" :key="index" class="tabs-item">
                  <el-input v-model="tabsDialog.tabs[index]" size="small" placeholder="标签名称" />
                  <el-button
                    v-if="tabsDialog.tabs.length > 1"
                    type="danger"
                    size="small"
                    text
                    @click="removeTabsDialogTab(index)"
                  >
                    <i class="ri-close-line"></i>
                  </el-button>
                </div>
              </div>
              <div class="tabs-footer">
                <el-button
                  type="primary"
                  size="small"
                  @click="addTabsDialogTab"
                  :disabled="tabsDialog.tabs.length >= 10"
                >
                  添加标签
                </el-button>
                <el-button type="primary" size="small" @click="handleInsertTabs"> 插入 </el-button>
              </div>
            </div>
          </el-popover>
        </template>
        <!-- 折叠面板按钮 -->
        <template v-else-if="item.title === '折叠面板'">
          <el-popover
            :width="300"
            trigger="click"
            placement="bottom"
            v-model:visible="foldDialog.visible"
          >
            <template #reference>
              <button
                :title="item.title"
                class="toolbar-btn"
                :class="{ active: foldDialog.visible }"
              >
                <i :class="item.icon"></i>
              </button>
            </template>
            <div class="fold-dialog-wrap">
              <div class="fold-form-item">
                <el-input
                  v-model="foldDialog.title"
                  placeholder="标题"
                  size="small"
                  clearable
                  @keyup.enter="handleInsertFold"
                />
              </div>
              <div class="fold-form-item">
                <el-checkbox v-model="foldDialog.open" size="small">默认展开</el-checkbox>
              </div>
              <div class="fold-form-actions">
                <el-button type="primary" size="small" @click="handleInsertFold"> 插入 </el-button>
              </div>
            </div>
          </el-popover>
        </template>
        <!-- 链接卡片按钮 -->
        <template v-else-if="item.title === '链接卡片'">
          <el-popover
            :width="320"
            trigger="click"
            placement="bottom"
            v-model:visible="linkDialog.visible"
          >
            <template #reference>
              <button
                :title="item.title"
                class="toolbar-btn"
                :class="{ active: linkDialog.visible }"
              >
                <i :class="item.icon"></i>
              </button>
            </template>
            <div class="link-dialog-wrap">
              <div class="link-form-item">
                <el-input v-model="linkDialog.title" placeholder="标题" size="small" clearable />
              </div>
              <div class="link-form-item">
                <el-input v-model="linkDialog.url" placeholder="https://" size="small" clearable />
              </div>
              <div class="link-form-item">
                <el-input
                  v-model="linkDialog.description"
                  placeholder="描述（可选）"
                  size="small"
                  clearable
                />
              </div>
              <div class="link-form-actions">
                <el-button type="primary" size="small" @click="handleInsertLink"> 插入 </el-button>
              </div>
            </div>
          </el-popover>
        </template>
        <!-- 照片墙按钮 -->
        <template v-else-if="item.title === '照片墙'">
          <el-popover
            :width="520"
            trigger="click"
            placement="bottom"
            v-model:visible="photoDialog.visible"
          >
            <template #reference>
              <button
                :title="item.title"
                class="toolbar-btn"
                :class="{ active: photoDialog.visible }"
              >
                <i :class="item.icon"></i>
              </button>
            </template>
            <div class="photo-dialog-wrap">
              <div class="photo-rows">
                <div v-for="(row, rowIndex) in photoDialog.rows" :key="rowIndex" class="photo-row">
                  <div class="photo-row-header">
                    <div class="photo-row-actions">
                      <el-button
                        :disabled="rowIndex === 0"
                        size="small"
                        text
                        @click="movePhotoRowUp(rowIndex)"
                      >
                        <i class="ri-arrow-up-line"></i>
                      </el-button>
                      <el-button
                        :disabled="rowIndex === photoDialog.rows.length - 1"
                        size="small"
                        text
                        @click="movePhotoRowDown(rowIndex)"
                      >
                        <i class="ri-arrow-down-line"></i>
                      </el-button>
                    </div>
                    <el-button
                      v-if="photoDialog.rows.length > 1"
                      type="danger"
                      size="small"
                      text
                      @click="removePhotoDialogRow(rowIndex)"
                    >
                      <i class="ri-close-line"></i>
                    </el-button>
                  </div>
                  <div class="photo-images">
                    <div v-for="(img, imgIndex) in row" :key="imgIndex" class="photo-image-item">
                      <el-input
                        :model-value="getPhotoImageUrl(rowIndex, imgIndex)"
                        placeholder="图片URL"
                        size="small"
                        @update:model-value="
                          (val: string) => setPhotoImageUrl(rowIndex, imgIndex, val)
                        "
                      >
                        <template #append>
                          <el-upload
                            :show-file-list="false"
                            :before-upload="
                              (file: File) => {
                                handlePhotoImageUpload(rowIndex, imgIndex, file);
                                return false;
                              }
                            "
                            :disabled="photoDialog.uploading"
                          >
                            <el-button :loading="photoDialog.uploading" size="small">
                              <i class="ri-upload-line"></i>
                            </el-button>
                          </el-upload>
                        </template>
                      </el-input>
                      <div class="photo-image-actions">
                        <el-button
                          :disabled="imgIndex === 0"
                          size="small"
                          text
                          @click="movePhotoImageUp(rowIndex, imgIndex)"
                        >
                          <i class="ri-arrow-up-line"></i>
                        </el-button>
                        <el-button
                          :disabled="imgIndex === row.length - 1"
                          size="small"
                          text
                          @click="movePhotoImageDown(rowIndex, imgIndex)"
                        >
                          <i class="ri-arrow-down-line"></i>
                        </el-button>
                        <el-button
                          v-if="row.length > 1"
                          type="danger"
                          size="small"
                          text
                          @click="removePhotoDialogImage(rowIndex, imgIndex)"
                        >
                          <i class="ri-close-line"></i>
                        </el-button>
                      </div>
                    </div>
                    <el-button
                      v-if="row.length < 4"
                      type="primary"
                      size="small"
                      text
                      @click="addPhotoDialogImage(rowIndex)"
                    >
                      <i class="ri-add-line"></i> 添加图片
                    </el-button>
                  </div>
                </div>
              </div>
              <div class="photo-footer">
                <el-button
                  type="primary"
                  size="small"
                  @click="addPhotoDialogRow"
                  :disabled="photoDialog.rows.length >= 6"
                >
                  添加行
                </el-button>
                <el-button type="primary" size="small" @click="handleInsertPhoto"> 插入 </el-button>
              </div>
            </div>
          </el-popover>
        </template>
        <!-- 视频按钮 -->
        <template v-else-if="item.title === '视频'">
          <el-popover
            :width="320"
            trigger="click"
            placement="bottom"
            v-model:visible="videoDialog.visible"
          >
            <template #reference>
              <button
                :title="item.title"
                class="toolbar-btn"
                :class="{ active: videoDialog.visible }"
              >
                <i :class="item.icon"></i>
              </button>
            </template>
            <div class="video-dialog-wrap">
              <div class="video-form-item">
                <el-radio-group v-model="videoDialog.type" size="small">
                  <el-radio-button value="url">在线视频</el-radio-button>
                  <el-radio-button value="upload">本地上传</el-radio-button>
                </el-radio-group>
              </div>
              <div v-if="videoDialog.type === 'url'" class="video-form-item">
                <el-input
                  v-model="videoDialog.videoUrl"
                  placeholder="输入视频链接，支持B站、YouTube等"
                  size="small"
                  clearable
                  @keyup.enter="handleInsertVideo"
                />
              </div>
              <div v-else class="video-form-item">
                <el-upload
                  :show-file-list="false"
                  :before-upload="
                    (file: File) => {
                      handleVideoUpload(file);
                      return false;
                    }
                  "
                  :disabled="videoDialog.uploading"
                >
                  <el-button
                    type="primary"
                    :loading="videoDialog.uploading"
                    size="small"
                    style="width: 100%"
                  >
                    <i v-if="!videoDialog.uploading" class="ri-upload-line"></i>
                    {{ videoDialog.videoUrl ? '已上传' : '选择视频文件' }}
                  </el-button>
                </el-upload>
                <div v-if="videoDialog.videoUrl" class="video-url-preview">
                  {{ videoDialog.videoUrl }}
                </div>
              </div>
              <div class="video-form-actions">
                <el-button
                  type="primary"
                  size="small"
                  :loading="videoDialog.loading"
                  @click="handleInsertVideo"
                >
                  插入
                </el-button>
              </div>
            </div>
          </el-popover>
        </template>
        <!-- 音乐按钮 -->
        <template v-else-if="item.title === '音乐'">
          <el-popover
            :width="320"
            trigger="click"
            placement="bottom"
            v-model:visible="audioDialog.visible"
          >
            <template #reference>
              <button
                :title="item.title"
                class="toolbar-btn"
                :class="{ active: audioDialog.visible }"
              >
                <i :class="item.icon"></i>
              </button>
            </template>
            <div class="audio-dialog-wrap">
              <div class="audio-form-item">
                <el-radio-group v-model="audioDialog.type" size="small">
                  <el-radio-button value="music">在线音乐</el-radio-button>
                  <el-radio-button value="upload">本地上传</el-radio-button>
                </el-radio-group>
              </div>
              <template v-if="audioDialog.type === 'upload'">
                <div class="audio-form-item">
                  <el-input
                    v-model="audioDialog.title"
                    placeholder="音频标题"
                    size="small"
                    clearable
                  />
                </div>
                <div class="audio-form-item">
                  <el-upload
                    :show-file-list="false"
                    :before-upload="
                      (file: File) => {
                        handleAudioUpload(file);
                        return false;
                      }
                    "
                    :disabled="audioDialog.uploading"
                  >
                    <el-button
                      type="primary"
                      :loading="audioDialog.uploading"
                      size="small"
                      style="width: 100%"
                    >
                      <i v-if="!audioDialog.uploading" class="ri-upload-line"></i>
                      {{ audioDialog.audioUrl ? '已上传' : '选择音频文件' }}
                    </el-button>
                  </el-upload>
                  <div v-if="audioDialog.audioUrl" class="audio-url-preview">
                    {{ audioDialog.audioUrl }}
                  </div>
                </div>
              </template>
              <template v-else>
                <div class="audio-form-item">
                  <el-select v-model="audioDialog.musicServer" size="small" style="width: 100%">
                    <el-option value="netease" label="网易云音乐" />
                    <el-option value="tencent" label="QQ音乐" />
                    <el-option value="kugou" label="酷狗音乐" />
                    <el-option value="xiami" label="虾米音乐" />
                    <el-option value="baidu" label="百度音乐" />
                  </el-select>
                </div>
                <div class="audio-form-item">
                  <el-input
                    v-model="audioDialog.musicId"
                    placeholder="输入音乐ID"
                    size="small"
                    clearable
                  />
                </div>
                <div class="audio-form-item">
                  <el-button
                    type="primary"
                    size="small"
                    :loading="audioDialog.loading"
                    style="width: 100%"
                    @click="handleParseMusic"
                  >
                    解析
                  </el-button>
                </div>
                <div v-if="audioDialog.musicInfo" class="music-info-preview">
                  <div class="music-info-title">{{ audioDialog.musicInfo.title }}</div>
                  <div class="music-info-artist">{{ audioDialog.musicInfo.artist }}</div>
                </div>
              </template>
              <div class="audio-form-actions">
                <el-button type="primary" size="small" @click="handleInsertAudio"> 插入 </el-button>
              </div>
            </div>
          </el-popover>
        </template>
        <!-- 普通按钮 -->
        <button
          v-else
          @click="item.action"
          :title="item.title"
          :class="{
            active: item.isActive?.(),
            'mobile-only': item.mobileOnly,
          }"
          class="toolbar-btn"
        >
          <i v-if="item.icon" :class="item.icon"></i>
          <span v-else>{{ item.label }}</span>
        </button>
      </template>

      <input
        ref="imageInputRef"
        type="file"
        accept="image/*"
        multiple
        style="display: none"
        @change="handleImageSelect"
      />
    </div>

    <!-- 编辑器主体 -->
    <div class="editor-container">
      <!-- 编辑器面板 -->
      <div
        class="editor-pane"
        :class="{
          'full-width': viewMode === 'editor',
          hidden: viewMode === 'preview',
        }"
        @mousedown="handleEditorPaneMouseDown"
      >
        <div ref="editorRef" class="cm-host"></div>
      </div>

      <!-- 预览面板 -->
      <div
        v-show="viewMode !== 'editor'"
        ref="previewPaneRef"
        class="preview-pane"
        :class="{
          'full-width': viewMode === 'preview',
          'html-mode': viewMode === 'html',
        }"
      >
        <div v-if="viewMode === 'html'" class="html-code">
          <pre><code>{{ renderedHtml }}</code></pre>
        </div>
        <div v-else class="markdown-content" v-html="renderedHtml"></div>
      </div>

      <!-- 目录面板 -->
      <div v-if="showToc" class="toc-pane">
        <div class="toc-header">
          <span>目录</span>
          <button @click="showToc = false" class="toc-close">
            <i class="ri-close-line"></i>
          </button>
        </div>
        <div class="toc-content">
          <div
            v-for="(heading, index) in tableOfContents"
            :key="index"
            :class="`toc-item toc-level-${heading.level}`"
            @click="scrollToHeading(heading)"
          >
            {{ heading.text }}
          </div>
          <div v-if="tableOfContents.length === 0" class="toc-empty">暂无目录</div>
        </div>
      </div>
    </div>

    <!-- 页脚 -->
    <div class="editor-footer">
      <div class="footer-left">
        <span class="word-count">字数：{{ wordCount }}</span>
        <span class="reading-time">阅读时长：{{ readingTime }} 分钟</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  ref,
  shallowRef,
  reactive,
  computed,
  onMounted,
  onBeforeUnmount,
  watch,
  nextTick,
} from 'vue';
import { ElMessage } from 'element-plus';
import { uploadFile } from '@/api/file';
import { parseVideo } from '@/api/tools';
import { getSettingGroup } from '@/api/sysconfig';
import {
  renderMarkdownWithSourceMap,
  renderMarkdownWithStyles,
  countWords,
  estimateReadingTime,
  extractToc,
  type TocItem,
} from '@/utils/markdown';
import { EditorView, keymap, showPanel } from '@codemirror/view';
import { EditorState, StateField, StateEffect, RangeSetBuilder } from '@codemirror/state';
import { Decoration } from '@codemirror/view';
import type { Panel, DecorationSet } from '@codemirror/view';
import { defaultKeymap, history, historyKeymap } from '@codemirror/commands';
import { markdown } from '@codemirror/lang-markdown';
import { SearchCursor } from '@codemirror/search';
import mermaid from 'mermaid';

// 简易搜索功能
const setSearchQuery = StateEffect.define<string>();
const setSearchIndex = StateEffect.define<number>();

const searchStateField = StateField.define<{
  matches: { from: number; to: number }[];
  idx: number;
}>({
  create: () => ({ matches: [], idx: 0 }),
  update: (v, tr) => {
    for (const e of tr.effects) {
      if (e.is(setSearchQuery)) {
        if (!e.value) return { matches: [], idx: 0 };
        const matches: { from: number; to: number }[] = [];
        const cursor = new SearchCursor(tr.state.doc, e.value, 0, undefined, s => s.toLowerCase());
        while (!cursor.next().done) matches.push({ from: cursor.value.from, to: cursor.value.to });
        return { matches, idx: 0 };
      }
      if (e.is(setSearchIndex)) return { ...v, idx: e.value };
    }
    return v;
  },
});

const searchDecorations = StateField.define<DecorationSet>({
  create: () => Decoration.none,
  update: (_, tr) => {
    const { matches, idx } = tr.state.field(searchStateField);
    if (!matches.length) return Decoration.none;
    const builder = new RangeSetBuilder<Decoration>();
    matches.forEach((m, i) =>
      builder.add(
        m.from,
        m.to,
        Decoration.mark({
          class: i === idx ? 'cm-searchMatch-selected' : 'cm-searchMatch',
        })
      )
    );
    return builder.finish();
  },
  provide: f => EditorView.decorations.from(f),
});

let searchPanel: { dom: HTMLElement; show: () => void } | null = null;

function createSearchPanel(view: EditorView): Panel {
  const dom = document.createElement('div');
  dom.style.cssText =
    'display:none;align-items:center;padding:8px;background:#f5f5f5;border-top:1px solid #ddd';
  dom.innerHTML = `
    <input placeholder="查找..." style="width:180px;padding:4px 8px;border:1px solid #ddd;border-radius:4px;outline:none">
    <span style="margin:0 8px;color:#666;font-size:13px"></span>
    <button style="padding:4px 8px;border:1px solid #ddd;border-radius:4px;background:#fff;cursor:pointer">↑</button>
    <button style="padding:4px 8px;border:1px solid #ddd;border-radius:4px;background:#fff;cursor:pointer;margin-left:4px">↓</button>
    <button style="padding:4px 8px;border:1px solid #ddd;border-radius:4px;background:#fff;cursor:pointer;margin-left:8px">×</button>
  `;
  const [input, count, prev, next, close] = [
    dom.querySelector('input')!,
    dom.querySelector('span')!,
    ...dom.querySelectorAll('button'),
  ] as [HTMLInputElement, HTMLSpanElement, HTMLButtonElement, HTMLButtonElement, HTMLButtonElement];

  const update = () => {
    const { matches, idx } = view.state.field(searchStateField);
    count.textContent = matches.length
      ? `${idx + 1}/${matches.length}`
      : input.value
        ? '无匹配'
        : '';
  };

  const search = () => {
    view.dispatch({ effects: setSearchQuery.of(input.value) });
    update();
  };

  const go = (d: number) => {
    const { matches, idx } = view.state.field(searchStateField);
    if (!matches.length) return;
    const i = (idx + d + matches.length) % matches.length;
    view.dispatch({
      effects: setSearchIndex.of(i),
      selection: { anchor: matches[i]!.from, head: matches[i]!.to },
      scrollIntoView: true,
    });
    update();
  };

  input.oninput = search;
  input.onkeydown = e => {
    if (e.key === 'Enter') {
      e.preventDefault();
      go(e.shiftKey ? -1 : 1);
    }
    if (e.key === 'Escape') {
      view.dispatch({ effects: setSearchQuery.of('') });
      input.value = '';
      update();
    }
  };
  prev.onclick = () => go(-1);
  next.onclick = () => go(1);
  close.onclick = () => {
    view.dispatch({ effects: setSearchQuery.of('') });
    input.value = '';
    dom.style.display = 'none';
  };

  searchPanel = {
    dom,
    show: () => {
      dom.style.display = 'flex';
      input.focus();
      input.select();
    },
  };
  return { dom, top: false };
}

function openSearchPanelCustom() {
  searchPanel?.show();
  return true;
}

// 类型定义
interface ToolbarItem {
  type?: 'divider' | 'spacer';
  icon?: string;
  label?: string;
  title?: string;
  action?: () => void;
  isActive?: () => boolean;
  mobileOnly?: boolean;
}

interface PreviewAnchor {
  startOffset: number;
  endOffset: number;
  top: number;
  height: number;
  depth: number;
  kind: 'block' | 'text';
}

type ViewMode = 'split' | 'editor' | 'preview' | 'html';

// 常量
const SCROLL_EPSILON = 1;
const PREVIEW_SYNC_DURATION = 90;

const props = withDefaults(defineProps<{ modelValue: string }>(), {
  modelValue: '',
});
const emit = defineEmits<{
  'update:modelValue': [value: string];
  save: [content: string];
}>();

// Refs
const editorRef = ref<HTMLElement>();
const previewPaneRef = ref<HTMLElement>();
const imageInputRef = ref<HTMLInputElement>();
const viewMode = ref<ViewMode>('split');
const isBrowserFullscreen = ref(false);
const isPageFullscreen = ref(false);
const showToc = ref(false);
const onlineImageUrl = ref('');
const downloadingImage = ref(false);

// 表情选择器状态
const emojiState = reactive({
  visible: false,
  groups: [] as Array<{
    name: string;
    type: 'emoji' | 'image' | 'emoticon';
    items: Array<{ key: string; val: string }>;
  }>,
  activeTab: 0,
  emojiMap: new Map<string, string>(),
});

// 提示框弹窗状态
const noteDialog = reactive({
  visible: false,
  type: 'info' as 'info' | 'warning' | 'success' | 'error',
  title: '',
});

// 标签页弹窗状态
const tabsDialog = reactive({
  visible: false,
  tabs: ['标签1', '标签2'] as string[],
});

// 折叠面板弹窗状态
const foldDialog = reactive({
  visible: false,
  title: '',
  open: false,
});

// 链接卡片弹窗状态
const linkDialog = reactive({
  visible: false,
  title: '',
  url: '',
  description: '',
});

// 视频弹窗状态
const videoDialog = reactive({
  visible: false,
  type: 'url' as 'url' | 'upload',
  videoUrl: '',
  uploading: false,
  loading: false,
});

// 音频弹窗状态
const audioDialog = reactive({
  visible: false,
  type: 'music' as 'upload' | 'music',
  title: '',
  audioUrl: '',
  uploading: false,
  loading: false,
  musicServer: 'netease',
  musicId: '',
  musicInfo: null as { title: string; artist: string; pic: string } | null,
});

// 照片墙弹窗状态
const photoDialog = reactive({
  visible: false,
  rows: [['', '']] as string[][],
  uploading: false as boolean,
});

// 安全获取照片墙图片 URL
const getPhotoImageUrl = (rowIndex: number, imgIndex: number): string => {
  const row = photoDialog.rows[rowIndex];
  return row?.[imgIndex] ?? '';
};

// 安全设置照片墙图片 URL
const setPhotoImageUrl = (rowIndex: number, imgIndex: number, url: string) => {
  const row = photoDialog.rows[rowIndex];
  if (row) {
    row[imgIndex] = url;
  }
};

// 编辑器实例
const editorViewRef = shallowRef<EditorView | null>(null);

// ==================== Mermaid 图表渲染 ====================
const initMermaid = () => {
  mermaid.initialize({
    startOnLoad: false,
    theme: 'default',
    securityLevel: 'loose',
  });
};

const renderMermaidDiagrams = async () => {
  const preview = previewPaneRef.value;
  if (!preview) return;

  const elements = preview.querySelectorAll('.mermaid:not(:has(svg))');

  for (const element of elements) {
    try {
      const { svg } = await mermaid.render(`mermaid-${Date.now()}`, element.textContent || '');
      element.innerHTML = svg;
    } catch (error) {
      console.error('Mermaid 渲染失败:', error);
    }
  }
};

// ==================== 滚动同步 ====================
let cachedPreviewAnchors: PreviewAnchor[] | null = null;
let editorScrollFrame: number | null = null;
let boundEditorScroller: HTMLElement | null = null;
let previewResizeObserver: ResizeObserver | null = null;
let previewProgrammaticScrollFrame: number | null = null;
let previewTweenFrame: number | null = null;
let previewTargetScrollTop: number | null = null;
let isProgrammaticPreviewScroll = false;
let isPreviewManualScrollActive = false;

const getEditorScroller = () => editorViewRef.value?.scrollDOM ?? null;

const invalidateScrollCache = () => {
  cachedPreviewAnchors = null;
};

const clampNumber = (value: number, min: number, max: number) =>
  Math.min(max, Math.max(min, value));

const getPreviewElementTop = (element: HTMLElement, container: HTMLElement) => {
  const elementRect = element.getBoundingClientRect();
  const containerRect = container.getBoundingClientRect();
  const paddingTop = Number.parseFloat(getComputedStyle(container).paddingTop || '0') || 0;
  return elementRect.top - containerRect.top + container.scrollTop - paddingTop;
};

const getPreviewElementDepth = (element: HTMLElement, container: HTMLElement) => {
  let depth = 0;
  let current = element.parentElement;

  while (current && current !== container) {
    depth += 1;
    current = current.parentElement;
  }

  return depth;
};

const getPreviewAnchors = (): PreviewAnchor[] => {
  if (cachedPreviewAnchors) return cachedPreviewAnchors;

  const preview = previewPaneRef.value;
  if (!preview) return [];

  const anchors = Array.from(
    preview.querySelectorAll<HTMLElement>('[data-source-start-offset][data-source-end-offset]')
  )
    .map(element => {
      const startOffset = Number.parseInt(element.dataset.sourceStartOffset || '', 10);
      const endOffset = Number.parseInt(element.dataset.sourceEndOffset || '', 10);
      if (Number.isNaN(startOffset) || Number.isNaN(endOffset)) return null;

      const rect = element.getBoundingClientRect();
      return {
        startOffset,
        endOffset,
        top: Math.max(0, getPreviewElementTop(element, preview)),
        height: Math.max(1, rect.height, element.offsetHeight),
        depth: getPreviewElementDepth(element, preview),
        kind: (element.dataset.syncKind === 'text' ? 'text' : 'block') as 'block' | 'text',
      } satisfies PreviewAnchor;
    })
    .filter((anchor): anchor is PreviewAnchor => !!anchor)
    .sort((left, right) => {
      if (left.startOffset !== right.startOffset) return left.startOffset - right.startOffset;
      const leftSourceSpan = left.endOffset - left.startOffset;
      const rightSourceSpan = right.endOffset - right.startOffset;
      if (leftSourceSpan !== rightSourceSpan) return leftSourceSpan - rightSourceSpan;
      if (left.top !== right.top) return left.top - right.top;
      if (left.depth !== right.depth) return right.depth - left.depth;
      if (left.kind === right.kind) return 0;
      return left.kind === 'text' ? -1 : 1;
    });

  cachedPreviewAnchors = anchors;
  return anchors;
};

const getEditorTopSourceOffset = () => {
  const editor = editorViewRef.value;
  const editorScroller = getEditorScroller();
  if (!editor || !editorScroller) return 0;

  const scrollerRect = editorScroller.getBoundingClientRect();
  const contentRect = editor.contentDOM.getBoundingClientRect();
  const pos = editor.posAtCoords({
    x: Math.max(contentRect.left + 4, scrollerRect.left + 4),
    y: scrollerRect.top + 2,
  });

  if (pos !== null) return pos;
  return editor.lineBlockAtHeight(editorScroller.scrollTop).from;
};

const getAnchorSourceSpan = (anchor: PreviewAnchor) =>
  Math.max(0, anchor.endOffset - anchor.startOffset);

const compareAnchorSpecificity = (left: PreviewAnchor, right: PreviewAnchor) => {
  if (left.kind !== right.kind) return left.kind === 'text' ? 1 : -1;

  const leftSourceSpan = getAnchorSourceSpan(left);
  const rightSourceSpan = getAnchorSourceSpan(right);
  if (leftSourceSpan !== rightSourceSpan) return rightSourceSpan - leftSourceSpan;

  if (left.depth !== right.depth) return left.depth - right.depth;
  if (left.height !== right.height) return right.height - left.height;
  return right.top - left.top;
};

const findBestContainingAnchor = (sourceOffset: number, anchors: PreviewAnchor[]) => {
  let bestAnchor: PreviewAnchor | null = null;
  let bestIndex = -1;

  anchors.forEach((anchor, index) => {
    if (sourceOffset < anchor.startOffset || sourceOffset > anchor.endOffset) return;

    if (!bestAnchor || compareAnchorSpecificity(anchor, bestAnchor) > 0) {
      bestAnchor = anchor;
      bestIndex = index;
    }
  });

  return bestAnchor ? { anchor: bestAnchor, index: bestIndex } : null;
};

const getAnchorVisualSpan = (
  anchor: PreviewAnchor,
  anchors: PreviewAnchor[],
  anchorIndex: number
) => {
  let visualSpan = Math.max(1, anchor.height);

  for (let index = anchorIndex + 1; index < anchors.length; index++) {
    const candidate = anchors[index]!;
    if (candidate.top + SCROLL_EPSILON < anchor.top) continue;
    if (candidate.startOffset < anchor.endOffset) continue;
    visualSpan = Math.max(visualSpan, candidate.top - anchor.top);
    break;
  }

  return visualSpan;
};

const mapWithinAnchor = (
  sourceOffset: number,
  anchor: PreviewAnchor,
  anchors: PreviewAnchor[],
  anchorIndex: number
) => {
  const sourceSpan = getAnchorSourceSpan(anchor);
  if (sourceSpan <= 0) return anchor.top;

  const progress = clampNumber((sourceOffset - anchor.startOffset) / sourceSpan, 0, 1);
  return anchor.top + getAnchorVisualSpan(anchor, anchors, anchorIndex) * progress;
};

const mapSourceOffsetToPreviewTop = (sourceOffset: number, anchors: PreviewAnchor[]) => {
  if (!anchors.length) return 0;
  const containingAnchor = findBestContainingAnchor(sourceOffset, anchors);
  if (containingAnchor) {
    return mapWithinAnchor(sourceOffset, containingAnchor.anchor, anchors, containingAnchor.index);
  }

  let previousIndex = -1;
  let nextIndex = -1;

  anchors.forEach((anchor, index) => {
    if (anchor.endOffset <= sourceOffset) previousIndex = index;
    if (nextIndex === -1 && anchor.startOffset >= sourceOffset) nextIndex = index;
  });

  if (previousIndex === -1) return mapWithinAnchor(sourceOffset, anchors[0]!, anchors, 0);
  if (nextIndex === -1)
    return mapWithinAnchor(sourceOffset, anchors[anchors.length - 1]!, anchors, anchors.length - 1);

  const previous = anchors[previousIndex]!;
  const next = anchors[nextIndex]!;
  const previousTop = mapWithinAnchor(previous.endOffset, previous, anchors, previousIndex);
  const nextTop = mapWithinAnchor(next.startOffset, next, anchors, nextIndex);
  const sourceGap = next.startOffset - previous.endOffset;
  if (sourceGap <= 0) return nextTop;

  const progress = clampNumber((sourceOffset - previous.endOffset) / sourceGap, 0, 1);
  return previousTop + (nextTop - previousTop) * progress;
};

const schedulePreviewProgrammaticUnlock = () => {
  if (previewProgrammaticScrollFrame !== null) {
    cancelAnimationFrame(previewProgrammaticScrollFrame);
  }
  previewProgrammaticScrollFrame = requestAnimationFrame(() => {
    isProgrammaticPreviewScroll = false;
    previewProgrammaticScrollFrame = null;
  });
};

const setPreviewScrollTop = (preview: HTMLElement, nextTop: number) => {
  isProgrammaticPreviewScroll = true;
  preview.scrollTop = nextTop;
  schedulePreviewProgrammaticUnlock();
};

const stopPreviewTween = () => {
  if (previewTweenFrame !== null) {
    cancelAnimationFrame(previewTweenFrame);
    previewTweenFrame = null;
  }
};

const easeOutCubic = (progress: number) => 1 - Math.pow(1 - progress, 3);

const startPreviewTween = (preview: HTMLElement, targetTop: number) => {
  const startTop = preview.scrollTop;
  if (Math.abs(startTop - targetTop) <= SCROLL_EPSILON) {
    stopPreviewTween();
    if (Math.abs(startTop - targetTop) > 0) setPreviewScrollTop(preview, targetTop);
    return;
  }

  stopPreviewTween();
  const startTime = performance.now();

  const animate = (now: number) => {
    if (viewMode.value !== 'split' || isPreviewManualScrollActive) {
      previewTweenFrame = null;
      return;
    }

    const currentPreview = previewPaneRef.value;
    const latestTargetTop = previewTargetScrollTop;
    if (!currentPreview || latestTargetTop === null) {
      previewTweenFrame = null;
      return;
    }

    const maxScrollTop = Math.max(0, currentPreview.scrollHeight - currentPreview.clientHeight);
    const clampedTargetTop = clampNumber(latestTargetTop, 0, maxScrollTop);
    const progress = clampNumber((now - startTime) / PREVIEW_SYNC_DURATION, 0, 1);
    const nextTop = startTop + (clampedTargetTop - startTop) * easeOutCubic(progress);
    setPreviewScrollTop(currentPreview, progress >= 1 ? clampedTargetTop : nextTop);

    if (progress < 1 && Math.abs(clampedTargetTop - currentPreview.scrollTop) > SCROLL_EPSILON) {
      previewTweenFrame = requestAnimationFrame(animate);
      return;
    }

    previewTweenFrame = null;
  };

  previewTweenFrame = requestAnimationFrame(animate);
};

const syncPreviewToEditorTop = () => {
  if (viewMode.value !== 'split') return;
  if (isPreviewManualScrollActive) return;

  const preview = previewPaneRef.value;
  if (!preview) return;

  const anchors = getPreviewAnchors();
  if (!anchors.length) return;

  const targetTop = mapSourceOffsetToPreviewTop(getEditorTopSourceOffset(), anchors);
  const maxScrollTop = Math.max(0, preview.scrollHeight - preview.clientHeight);
  previewTargetScrollTop = Math.max(0, Math.min(targetTop, maxScrollTop));
  startPreviewTween(preview, previewTargetScrollTop);
};

const requestPreviewSync = () => {
  if (editorScrollFrame !== null) return;
  editorScrollFrame = requestAnimationFrame(() => {
    editorScrollFrame = null;
    syncPreviewToEditorTop();
  });
};

const cancelPreviewSync = () => {
  if (editorScrollFrame !== null) {
    cancelAnimationFrame(editorScrollFrame);
    editorScrollFrame = null;
  }
  stopPreviewTween();
  previewTargetScrollTop = null;
};

const resumePreviewSync = () => {
  if (!isPreviewManualScrollActive) return;
  isPreviewManualScrollActive = false;
  requestPreviewSync();
};

const handleEditorScroll = () => {
  if (isPreviewManualScrollActive) {
    isPreviewManualScrollActive = false;
  }
  requestPreviewSync();
};

const handlePreviewInteraction = () => {
  if (isProgrammaticPreviewScroll) return;
  isPreviewManualScrollActive = true;
  cancelPreviewSync();
};

const bindPreviewObservers = () => {
  const preview = previewPaneRef.value;
  if (!preview) return;

  previewResizeObserver?.disconnect();
  previewResizeObserver = new ResizeObserver(() => {
    invalidateScrollCache();
    requestPreviewSync();
  });

  const markdownContent = preview.querySelector('.markdown-content');
  if (markdownContent) {
    previewResizeObserver.observe(markdownContent);
  }
};

const bindScrollEvents = () => {
  const editorScroller = getEditorScroller();
  const preview = previewPaneRef.value;

  if (boundEditorScroller && boundEditorScroller !== editorScroller) {
    boundEditorScroller.removeEventListener('scroll', handleEditorScroll);
    boundEditorScroller = null;
  }

  if (editorScroller && boundEditorScroller !== editorScroller) {
    editorScroller.addEventListener('scroll', handleEditorScroll, {
      passive: true,
    });
    boundEditorScroller = editorScroller;
  }

  preview?.removeEventListener('scroll', handlePreviewInteraction);
  preview?.removeEventListener('click', handlePreviewInteraction);
  preview?.removeEventListener('click', togglePreviewImage);
  preview?.removeEventListener('wheel', handlePreviewInteraction);
  preview?.addEventListener('scroll', handlePreviewInteraction, {
    passive: true,
  });
  preview?.addEventListener('click', handlePreviewInteraction, {
    passive: true,
  });
  preview?.addEventListener('click', togglePreviewImage);
  preview?.addEventListener('wheel', handlePreviewInteraction, {
    passive: true,
  });
  bindPreviewObservers();
};

const unbindScrollEvents = () => {
  boundEditorScroller?.removeEventListener('scroll', handleEditorScroll);
  boundEditorScroller = null;
  previewPaneRef.value?.removeEventListener('scroll', handlePreviewInteraction);
  previewPaneRef.value?.removeEventListener('click', handlePreviewInteraction);
  previewPaneRef.value?.removeEventListener('click', togglePreviewImage);
  previewPaneRef.value?.removeEventListener('wheel', handlePreviewInteraction);
  previewResizeObserver?.disconnect();
  previewResizeObserver = null;
  cancelPreviewSync();
  if (previewProgrammaticScrollFrame !== null) {
    cancelAnimationFrame(previewProgrammaticScrollFrame);
    previewProgrammaticScrollFrame = null;
  }
  isProgrammaticPreviewScroll = false;
  isPreviewManualScrollActive = false;
};

const isPreviewVisible = computed(() => viewMode.value !== 'editor');

const renderedHtml = computed(() => {
  if (!isPreviewVisible.value) return '';

  const html =
    viewMode.value === 'html'
      ? renderMarkdownWithStyles(props.modelValue)
      : renderMarkdownWithSourceMap(props.modelValue);

  if (emojiState.emojiMap.size > 0) {
    return html.replace(/:([^:\s]+):/g, (match, key) => {
      const url = emojiState.emojiMap.get(key);
      if (url) {
        return `<img src="${url}" alt="${key}" class="emoji-image" title="${key}" />`;
      }
      return match;
    });
  }

  return html;
});

const wordCount = computed(() => countWords(props.modelValue));

const readingTime = computed(() => estimateReadingTime(props.modelValue));

const tableOfContents = computed<TocItem[]>(() => {
  return extractToc(props.modelValue);
});

// ==================== 编辑器操作 ====================

// 保存文章
const saveArticle = () => {
  const content = editorViewRef.value?.state.doc.toString() || '';

  if (!content.trim()) {
    ElMessage.warning('文章内容不能为空');
    return;
  }
  emit('save', content);

  ElMessage.success('文章保存成功');
};

// 插入文本到光标位置
const insertText = (before: string, after = '') => {
  if (!editorViewRef.value) return;
  const { from, to } = editorViewRef.value.state.selection.main;
  const text = editorViewRef.value.state.doc.sliceString(from, to);

  // 如果有选中文本，用语法包裹；否则只插入语法，光标定位在中间
  editorViewRef.value.dispatch({
    changes: { from, to, insert: `${before}${text}${after}` },
    // 如果有选中文本，保持选中状态；否则光标定位在中间
    selection: text
      ? {
          anchor: from + before.length,
          head: from + before.length + text.length,
        }
      : { anchor: from + before.length, head: from + before.length },
  });
  editorViewRef.value.focus();
};

// 插入标题
const insertHeading = (level: string) => insertText(`${'#'.repeat(+level)} `);

// 滚动到指定标题
const scrollToHeading = (heading: TocItem) => {
  if (!editorViewRef.value) return;
  const lines = editorViewRef.value.state.doc.toString().split('\n');
  const index = lines.findIndex(line => line.includes(heading.text) && line.startsWith('#'));

  if (index !== -1) {
    const pos = editorViewRef.value.state.doc.line(index + 1).from;
    editorViewRef.value.dispatch({
      selection: { anchor: pos, head: pos },
      effects: EditorView.scrollIntoView(pos, { y: 'start' }),
    });
    editorViewRef.value.focus();
  }
};

// 工具栏配置
const toolbarItems: ToolbarItem[] = [
  // 第一组：基本文本格式
  {
    icon: 'ri-bold',
    title: '粗体 (Ctrl+B)',
    action: () => insertText('**', '**'),
  },
  {
    icon: 'ri-underline',
    title: '下划线',
    action: () => insertText('++', '++'),
  },
  {
    icon: 'ri-italic',
    title: '斜体 (Ctrl+I)',
    action: () => insertText('*', '*'),
  },
  {
    icon: 'ri-strikethrough',
    title: '删除线',
    action: () => insertText('~~', '~~'),
  },
  { type: 'divider' },

  // 第二组：标题
  { label: 'H1', title: '一级标题', action: () => insertHeading('1') },
  { label: 'H2', title: '二级标题', action: () => insertHeading('2') },
  { label: 'H3', title: '三级标题', action: () => insertHeading('3') },
  { label: 'H4', title: '四级标题', action: () => insertHeading('4') },
  { label: 'H5', title: '五级标题', action: () => insertHeading('5') },
  { label: 'H6', title: '六级标题', action: () => insertHeading('6') },
  { type: 'divider' },
  { icon: 'ri-subscript', title: '下标', action: () => insertText('~', '~') },
  { icon: 'ri-superscript', title: '上标', action: () => insertText('^', '^') },
  { icon: 'ri-double-quotes-l', title: '引用', action: () => insertText('> ') },
  {
    icon: 'ri-list-unordered',
    title: '无序列表',
    action: () => insertText('- '),
  },
  {
    icon: 'ri-list-ordered',
    title: '有序列表',
    action: () => insertText('1. '),
  },
  {
    icon: 'ri-list-check',
    title: '任务列表',
    action: () => insertText('- [ ] '),
  },
  { type: 'divider' },

  // 第三组：代码和插入元素
  {
    icon: 'ri-code-line',
    title: '行内代码',
    action: () => insertText('`', '`'),
  },
  {
    icon: 'ri-code-box-line',
    title: '块级代码',
    action: () => insertText('\n```', '\n```\n'),
  },
  {
    icon: 'ri-link',
    title: '链接',
    action: () => insertText('[', '](https://)'),
  },
  {
    icon: 'ri-image-add-line',
    title: '上传本地图片',
    action: () => imageInputRef.value?.click(),
  },
  { icon: 'ri-image-download-line', title: '下载在线图片', action: () => {} },
  { icon: 'ri-emotion-line', title: '表情', action: () => toggleEmojiPicker() },
  {
    icon: 'ri-table-2',
    title: '表格',
    action: () => insertText('\n| 列1 | 列2 | 列3 |\n|:---:|:---:|:---:|\n|  ', '  |    |    |\n'),
  },
  {
    icon: 'ri-mark-pen-line',
    title: '高亮',
    action: () => insertText('==', '=='),
  },
  {
    icon: 'ri-superscript-2',
    title: '行内公式',
    action: () => insertText('$', '$'),
  },
  {
    icon: 'ri-functions',
    title: '块级公式',
    action: () => insertText('\n$$\n', '\n$$\n'),
  },
  { type: 'divider' },

  // 第四组：自定义块
  {
    icon: 'ri-information-line',
    title: '提示框',
    action: () => toggleNoteDialog(),
  },
  {
    icon: 'ri-layout-grid-line',
    title: '标签页',
    action: () => toggleTabsDialog(),
  },
  {
    icon: 'ri-increase-decrease-line',
    title: '折叠面板',
    action: () => toggleFoldDialog(),
  },
  {
    icon: 'ri-external-link-line',
    title: '链接卡片',
    action: () => toggleLinkDialog(),
  },
  {
    icon: 'ri-multi-image-line',
    title: '照片墙',
    action: () => togglePhotoDialog(),
  },
  {
    icon: 'ri-video-line',
    title: '视频',
    action: () => toggleVideoDialog(),
  },
  {
    icon: 'ri-music-line',
    title: '音乐',
    action: () => toggleAudioDialog(),
  },

  // 弹性空间，将后续按钮推到右侧
  { type: 'spacer' },

  // 第五组：视图控制（右侧）
  {
    icon: 'ri-fullscreen-line',
    title: '浏览器全屏',
    action: () =>
      document.fullscreenElement
        ? document.exitFullscreen()
        : document.documentElement.requestFullscreen(),
    isActive: () => isBrowserFullscreen.value,
  },
  {
    icon: 'ri-picture-in-picture-2-line',
    title: '页面全屏',
    action: () => (isPageFullscreen.value = !isPageFullscreen.value),
    isActive: () => isPageFullscreen.value,
  },
  {
    icon: 'ri-code-s-slash-line',
    title: 'HTML 代码预览',
    action: () => (viewMode.value = viewMode.value === 'html' ? 'split' : 'html'),
    isActive: () => viewMode.value === 'html',
  },
  {
    icon: 'ri-eye-line',
    title: '切换预览',
    action: () => (viewMode.value = viewMode.value === 'preview' ? 'editor' : 'preview'),
    isActive: () => viewMode.value === 'preview',
    mobileOnly: true,
  },
  {
    icon: 'ri-list-unordered',
    title: '目录',
    action: () => (showToc.value = !showToc.value),
    isActive: () => showToc.value,
  },
];

const uploadArticleImages = async (files: File[], onFinally?: () => void) => {
  const imageFiles = files.filter(file => {
    if (!file.type.startsWith('image/')) {
      ElMessage.error(`${file.name} 不是图片格式`);
      return false;
    }
    return true;
  });

  if (!imageFiles.length) {
    onFinally?.();
    return [];
  }

  const loading = ElMessage.info({
    message: `正在上传 ${imageFiles.length} 张图片...`,
    duration: 0,
  });
  try {
    const results = await Promise.all(imageFiles.map(file => uploadFile(file, '文章图片')));
    insertText(results.map(result => `![图片](${result.file_url})`).join('\n'));
    ElMessage.success(`成功上传 ${imageFiles.length} 张图片`);
    return results;
  } catch (error: unknown) {
    ElMessage.error((error as Error)?.message || '图片上传失败');
    return [];
  } finally {
    loading.close();
    onFinally?.();
  }
};

// ==================== 图片上传 ====================
const handleImageSelect = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  await uploadArticleImages(Array.from(input.files || []), () => {
    input.value = '';
  });
};

// 处理粘贴图片
const handlePasteImage = async (files: File[]) => {
  await uploadArticleImages(files);
};

function base64ToFile(base64Data: string, contentType: string, fileName: string): File {
  const byteCharacters = atob(base64Data);
  const byteNumbers = new Array(byteCharacters.length);
  for (let i = 0; i < byteCharacters.length; i++) {
    byteNumbers[i] = byteCharacters.charCodeAt(i);
  }
  return new File([new Uint8Array(byteNumbers)], fileName, {
    type: contentType,
  });
}

// 处理下载在线图片
const handleOnlineImageDownload = async () => {
  if (!onlineImageUrl.value.trim()) {
    ElMessage.error('请输入图片URL');
    return;
  }

  const url = onlineImageUrl.value.trim();
  if (!url.match(/^https?:\/\/.+/)) {
    ElMessage.error('请输入有效的HTTP/HTTPS图片URL');
    return;
  }

  downloadingImage.value = true;
  try {
    const { downloadImage } = await import('@/api/tools');
    const downloadResult = await downloadImage({ url });
    const file = base64ToFile(downloadResult.data, downloadResult.content_type, 'image.jpg');
    const [uploadResult] = await uploadArticleImages([file]);

    if (uploadResult) {
      onlineImageUrl.value = '';
      ElMessage.success('图片下载并上传成功');
      document.body.click();
    }
  } catch (error: unknown) {
    ElMessage.error((error as Error)?.message || '图片下载失败');
  } finally {
    downloadingImage.value = false;
  }
};

// 表情选择器
const loadEmojis = async () => {
  if (emojiState.groups.length) return;

  const blogSettings = await getSettingGroup('blog');
  const emojisUrl = blogSettings.emojis || blogSettings['blog.emojis'] || '';
  if (!emojisUrl) return;

  const response = await fetch(emojisUrl);
  const groups = await response.json();
  emojiState.groups = groups;

  // 构建 image 类型表情映射
  for (const group of groups) {
    if (group.type === 'image') {
      for (const item of group.items) {
        emojiState.emojiMap.set(item.key, item.val);
      }
    }
  }
};

const selectEmoji = (item: { key: string; val: string }, type: string) => {
  const emoji = type === 'image' ? `:${item.key}:` : item.val;
  insertText(emoji);
  emojiState.visible = false;
};

// 表情选择器显示时加载数据
const handleEmojiPickerShow = () => {
  if (!emojiState.groups.length) {
    loadEmojis();
  }
};

const toggleEmojiPicker = () => {
  emojiState.visible = !emojiState.visible;
  if (emojiState.visible && !emojiState.groups.length) {
    loadEmojis();
  }
};

// 切换提示框弹窗显示
const toggleNoteDialog = () => {
  noteDialog.visible = !noteDialog.visible;
  if (noteDialog.visible) {
    noteDialog.type = 'info';
    noteDialog.title = '';
  }
};

// 插入提示框语法
const handleInsertNote = () => {
  const title = noteDialog.title.trim() || '标题';
  const typeLabel = noteDialog.type;
  insertText(`:::note ${typeLabel} ${title}\n内容\n:::endnote\n`);
  noteDialog.visible = false;
};

// 切换标签页弹窗显示
const toggleTabsDialog = () => {
  tabsDialog.visible = !tabsDialog.visible;
  if (tabsDialog.visible) {
    tabsDialog.tabs = ['标签1', '标签2'];
  }
};

// 插入标签页语法
const handleInsertTabs = () => {
  const tabs = tabsDialog.tabs.filter(t => t.trim());
  if (tabs.length === 0) {
    tabsDialog.tabs = ['标签1', '标签2'];
    return handleInsertTabs();
  }
  const tabBlocks = tabs.map((tab, i) => `:::tab ${tab}\n内容${i + 1}\n:::endtab`).join('\n');
  insertText(`:::tabs\n${tabBlocks}\n:::endtabs\n`);
  tabsDialog.visible = false;
};

// 添加标签页
const addTabsDialogTab = () => {
  if (tabsDialog.tabs.length < 10) {
    tabsDialog.tabs.push(`标签${tabsDialog.tabs.length + 1}`);
  }
};

// 删除标签页
const removeTabsDialogTab = (index: number) => {
  if (tabsDialog.tabs.length > 1) {
    tabsDialog.tabs.splice(index, 1);
  }
};

// 切换折叠面板弹窗显示
const toggleFoldDialog = () => {
  foldDialog.visible = !foldDialog.visible;
  if (foldDialog.visible) {
    foldDialog.title = '';
    foldDialog.open = false;
  }
};

// 插入折叠面板语法
const handleInsertFold = () => {
  const title = foldDialog.title.trim() || '点击展开';
  const openFlag = foldDialog.open ? ' open' : '';
  insertText(`:::fold ${title}${openFlag}\n内容\n:::endfold\n`);
  foldDialog.visible = false;
};

// 切换链接卡片弹窗显示
const toggleLinkDialog = () => {
  linkDialog.visible = !linkDialog.visible;
  if (linkDialog.visible) {
    linkDialog.title = '';
    linkDialog.url = '';
    linkDialog.description = '';
  }
};

// 插入链接卡片语法
const handleInsertLink = () => {
  const title = linkDialog.title.trim() || '标题';
  let url = linkDialog.url.trim();
  const description = linkDialog.description.trim();
  if (!url) {
    url = 'https://example.com';
  } else if (!url.startsWith('http://') && !url.startsWith('https://')) {
    url = 'https://' + url;
  }
  const descPart = description ? ` ${description}` : '';
  insertText(`:::link ${title} ${url}${descPart} :::\n`);
  linkDialog.visible = false;
};

// 切换视频弹窗显示
const toggleVideoDialog = () => {
  videoDialog.visible = !videoDialog.visible;
  if (videoDialog.visible) {
    videoDialog.type = 'url';
    videoDialog.videoUrl = '';
    videoDialog.uploading = false;
    videoDialog.loading = false;
  }
};

// 处理视频上传
const handleVideoUpload = async (file: File) => {
  videoDialog.uploading = true;
  try {
    const results = await uploadFile(file, '文章视频');
    videoDialog.videoUrl = results.file_url;
    ElMessage.success('视频上传成功');
  } catch (error: unknown) {
    ElMessage.error((error as Error)?.message || '视频上传失败');
  } finally {
    videoDialog.uploading = false;
  }
};

// 插入视频语法
const handleInsertVideo = () => {
  if (videoDialog.type === 'url') {
    const url = videoDialog.videoUrl.trim() || 'https://example.com/video.mp4';
    videoDialog.loading = true;
    parseVideo({ url })
      .then(info => {
        if (info.platform && info.video_id) {
          insertText(`:::video ${info.platform} ${info.video_id} :::\n`);
        } else {
          insertText(`:::video ${url} :::\n`);
        }
        videoDialog.visible = false;
      })
      .catch(() => {
        insertText(`:::video ${url} :::\n`);
        videoDialog.visible = false;
      })
      .finally(() => {
        videoDialog.loading = false;
      });
  } else {
    const url = videoDialog.videoUrl.trim() || 'https://example.com/video.mp4';
    insertText(`:::video ${url} :::\n`);
    videoDialog.visible = false;
  }
};

// 切换音频弹窗显示
const toggleAudioDialog = () => {
  audioDialog.visible = !audioDialog.visible;
  if (audioDialog.visible) {
    audioDialog.type = 'music';
    audioDialog.title = '';
    audioDialog.audioUrl = '';
    audioDialog.uploading = false;
    audioDialog.loading = false;
    audioDialog.musicServer = 'netease';
    audioDialog.musicId = '';
    audioDialog.musicInfo = null;
  }
};

// 处理音频上传
const handleAudioUpload = async (file: File) => {
  audioDialog.uploading = true;
  try {
    const results = await uploadFile(file, '文章音频');
    audioDialog.audioUrl = results.file_url;
    ElMessage.success('音频上传成功');
  } catch (error: unknown) {
    ElMessage.error((error as Error)?.message || '音频上传失败');
  } finally {
    audioDialog.uploading = false;
  }
};

// 解析音乐
const handleParseMusic = async () => {
  if (!audioDialog.musicId.trim()) {
    ElMessage.warning('请输入音乐ID');
    return;
  }
  audioDialog.loading = true;
  try {
    const apiUrl = `https://api.injahow.cn/meting/?server=${audioDialog.musicServer}&type=song&id=${audioDialog.musicId.trim()}`;
    const response = await fetch(apiUrl);
    const data = await response.json();
    if (data && data.length > 0) {
      const info = data[0];
      audioDialog.musicInfo = {
        title: info.name || info.title || '未知歌曲',
        artist: info.artist || info.author || '未知艺术家',
        pic: info.pic || info.cover || '',
      };
      audioDialog.title = `${audioDialog.musicInfo.title} - ${audioDialog.musicInfo.artist}`;
      ElMessage.success('解析成功');
    } else {
      throw new Error('未获取到音乐信息');
    }
  } catch (_error) {
    ElMessage.error('解析失败，请检查音乐ID是否正确');
    audioDialog.musicInfo = null;
  } finally {
    audioDialog.loading = false;
  }
};

// 插入音频语法
const handleInsertAudio = () => {
  if (audioDialog.type === 'upload') {
    const title = audioDialog.title.trim() || '音频';
    const url = audioDialog.audioUrl.trim() || 'https://example.com/audio.mp3';
    insertText(`:::audio ${title} ${url} :::\n`);
    audioDialog.visible = false;
  } else {
    if (!audioDialog.musicInfo) {
      insertText(
        `:::music ${audioDialog.musicServer} ${audioDialog.musicId.trim() || '音乐ID'} :::\n`
      );
    } else {
      insertText(`:::music ${audioDialog.musicServer} ${audioDialog.musicId.trim()} :::\n`);
    }
    audioDialog.visible = false;
  }
};

// 切换照片墙弹窗显示
const togglePhotoDialog = () => {
  photoDialog.visible = !photoDialog.visible;
  if (photoDialog.visible) {
    photoDialog.rows = [['', '']];
  }
};

// 添加照片墙行
const addPhotoDialogRow = () => {
  if (photoDialog.rows.length < 6) {
    photoDialog.rows.push(['', '']);
  }
};

// 删除照片墙行
const removePhotoDialogRow = (index: number) => {
  if (photoDialog.rows.length > 1) {
    photoDialog.rows.splice(index, 1);
  }
};

// 添加照片墙图片
const addPhotoDialogImage = (rowIndex: number) => {
  const row = photoDialog.rows[rowIndex];
  if (row && row.length < 4) {
    row.push('');
  }
};

// 删除照片墙图片
const removePhotoDialogImage = (rowIndex: number, imgIndex: number) => {
  const row = photoDialog.rows[rowIndex];
  if (row && row.length > 1) {
    row.splice(imgIndex, 1);
  }
};

// 处理照片墙图片上传
const handlePhotoImageUpload = async (rowIndex: number, imgIndex: number, file: File) => {
  photoDialog.uploading = true;
  try {
    const results = await uploadFile(file, '文章图片');
    const row = photoDialog.rows[rowIndex];
    if (row) {
      row[imgIndex] = results.file_url;
    }
    ElMessage.success('图片上传成功');
  } catch (error: unknown) {
    ElMessage.error((error as Error)?.message || '图片上传失败');
  } finally {
    photoDialog.uploading = false;
  }
};

// 照片墙上移
const movePhotoImageUp = (rowIndex: number, imgIndex: number) => {
  const row = photoDialog.rows[rowIndex];
  if (row && imgIndex > 0) {
    const temp = row[imgIndex] ?? '';
    const prev = row[imgIndex - 1] ?? '';
    row[imgIndex - 1] = temp as string;
    row[imgIndex] = prev as string;
  }
};

// 照片墙下移
const movePhotoImageDown = (rowIndex: number, imgIndex: number) => {
  const row = photoDialog.rows[rowIndex];
  if (row && imgIndex < row.length - 1) {
    const temp = row[imgIndex] ?? '';
    const next = row[imgIndex + 1] ?? '';
    row[imgIndex + 1] = temp as string;
    row[imgIndex] = next as string;
  }
};

// 照片墙行上移
const movePhotoRowUp = (rowIndex: number) => {
  if (rowIndex > 0) {
    const temp = photoDialog.rows[rowIndex];
    const prevRow = photoDialog.rows[rowIndex - 1];
    if (temp && prevRow) {
      photoDialog.rows[rowIndex - 1] = temp;
      photoDialog.rows[rowIndex] = prevRow;
    }
  }
};

// 照片墙行下移
const movePhotoRowDown = (rowIndex: number) => {
  if (rowIndex < photoDialog.rows.length - 1) {
    const temp = photoDialog.rows[rowIndex];
    const nextRow = photoDialog.rows[rowIndex + 1];
    if (temp && nextRow) {
      photoDialog.rows[rowIndex + 1] = temp;
      photoDialog.rows[rowIndex] = nextRow;
    }
  }
};

// 插入照片墙语法
const handleInsertPhoto = () => {
  const rows = photoDialog.rows.filter(row => row.some(img => img.trim()));
  let photoBlocks;
  if (rows.length === 0) {
    photoBlocks = '图片1\n图片2\n:::n\n图片3\n图片4';
  } else {
    photoBlocks = rows.map(row => row.filter(img => img.trim()).join('\n')).join('\n:::n\n');
  }

  insertText(`:::photo\n${photoBlocks}\n:::endphoto\n`);
  photoDialog.visible = false;
};

// ==================== 编辑器初始化 ====================
const initEditor = () => {
  if (!editorRef.value) return;

  // 创建粘贴事件处理器
  const pasteHandler = EditorView.domEventHandlers({
    paste: (event: ClipboardEvent, view) => {
      // 先检查是否有图片
      const items = event.clipboardData?.items;
      if (items) {
        const files: File[] = [];
        const textItems: DataTransferItem[] = [];

        for (let i = 0; i < items.length; i++) {
          const item = items[i];
          if (item && item.type) {
            if (item.type.startsWith('image/')) {
              const file = item.getAsFile();
              if (file) {
                files.push(file);
              }
            } else if (item.kind === 'string' && item.type === 'text/plain') {
              textItems.push(item);
            }
          }
        }

        // 如果有图片，处理图片上传
        if (files.length > 0) {
          event.preventDefault();
          handlePasteImage(files);

          // 如果还有文本，在图片处理完后再处理文本
          if (textItems.length > 0) {
            textItems.forEach(item => {
              item.getAsString(text => {
                // 使用默认的粘贴行为来正确替换选中文本
                view.dispatch({
                  changes: {
                    from: view.state.selection.main.from,
                    to: view.state.selection.main.to,
                    insert: text,
                  },
                });
              });
            });
          }
          return;
        }
      }

      // 如果没有图片，让默认行为处理（这样能正确替换选中文本）
      // 不调用 event.preventDefault()
    },
  });

  editorViewRef.value = new EditorView({
    state: EditorState.create({
      doc: props.modelValue,
      extensions: [
        history(),
        markdown(),
        searchStateField,
        searchDecorations,
        showPanel.of(createSearchPanel),
        keymap.of([
          {
            key: 'Mod-b',
            run: () => (insertText('**', '**'), true),
            preventDefault: true,
          },
          {
            key: 'Mod-i',
            run: () => (insertText('*', '*'), true),
            preventDefault: true,
          },
          {
            key: 'Mod-s',
            run: () => {
              saveArticle();
              return true;
            },
            preventDefault: true,
          },
          { key: 'Mod-f', run: openSearchPanelCustom, preventDefault: true },
          ...defaultKeymap,
          ...historyKeymap,
        ]),
        EditorView.updateListener.of(update => {
          if (update.docChanged) {
            emit('update:modelValue', update.state.doc.toString());
            invalidateScrollCache();
            requestPreviewSync();
          }
        }),
        EditorView.lineWrapping,
        pasteHandler,
      ],
    }),
    parent: editorRef.value,
  });

  // 编辑器初始化完成后，绑定滚动同步事件
  nextTick(() => {
    bindScrollEvents();
    requestPreviewSync();
  });
};

// 监听外部内容变化
watch(
  () => props.modelValue,
  newValue => {
    if (editorViewRef.value && newValue !== editorViewRef.value.state.doc.toString()) {
      editorViewRef.value.dispatch({
        changes: {
          from: 0,
          to: editorViewRef.value.state.doc.length,
          insert: newValue,
        },
      });
      invalidateScrollCache();
      requestPreviewSync();
    }
  }
);

// 监听预览区图片加载完成，使缓存失效
watch(renderedHtml, async html => {
  if (!html || !isPreviewVisible.value || viewMode.value === 'html') return;

  await nextTick();
  const preview = previewPaneRef.value;
  if (!preview) return;

  const images = preview.querySelectorAll('img');
  images.forEach(img => {
    if (img.complete) return;
    img.addEventListener(
      'load',
      () => {
        invalidateScrollCache();
        requestPreviewSync();
      },
      { once: true }
    );
  });

  invalidateScrollCache();
  await renderMermaidDiagrams();
  invalidateScrollCache();
  bindPreviewObservers();
  requestPreviewSync();
});

// 监听视图模式变化
watch(viewMode, newMode => {
  if (newMode === 'split') {
    nextTick(() => {
      invalidateScrollCache();
      bindScrollEvents();
      requestPreviewSync();
    });
  } else {
    unbindScrollEvents();
  }

  if (newMode !== 'editor') {
    loadEmojis();
  }
});

// ==================== 生命周期 ====================
const handleFullscreenChange = () => (isBrowserFullscreen.value = !!document.fullscreenElement);
const handleWindowResize = () => {
  invalidateScrollCache();
  requestPreviewSync();
};

const togglePreviewImage = (event: MouseEvent) => {
  const target = event.target as HTMLElement | null;
  const image = target?.closest('.preview-collapsible-image') as HTMLImageElement | null;
  if (!image) return;
  if (image.closest('.custom-photo-wall')) return;
  if (image.classList.contains('emoji-image')) return;

  image.classList.toggle('is-expanded');
  invalidateScrollCache();
  requestPreviewSync();
};

const handleEditorPaneMouseDown = (event: MouseEvent) => {
  if (event.button !== 0) return;
  if (!editorViewRef.value) return;

  resumePreviewSync();

  const target = event.target as HTMLElement | null;
  if (target?.closest('.cm-editor')) return;

  editorViewRef.value.focus();
};

onMounted(() => {
  initMermaid();
  initEditor();
  if (viewMode.value !== 'editor') {
    loadEmojis();
  }
  document.addEventListener('fullscreenchange', handleFullscreenChange);
  window.addEventListener('resize', handleWindowResize);

  // 移动端默认为纯编辑模式
  if (window.innerWidth <= 768) {
    viewMode.value = 'editor';
  }
});

onBeforeUnmount(() => {
  // 解绑滚动同步事件
  unbindScrollEvents();
  // 销毁编辑器实例
  editorViewRef.value?.destroy();
  window.removeEventListener('resize', handleWindowResize);
  document.removeEventListener('fullscreenchange', handleFullscreenChange);
});
</script>

<style lang="scss">
// 引入 Markdown 内容排版样式
@use '@/assets/css/prose';

// 引入代码高亮样式
@import 'highlight.js/styles/github.css';

// 引入 KaTeX 数学公式样式
@import 'katex/dist/katex.min.css';

// 搜索高亮样式
.cm-searchMatch {
  background-color: #ffeb3b80;
  border-radius: 2px;
}

.cm-searchMatch-selected {
  background-color: #ff9800;
  color: white;
}
</style>

<style scoped lang="scss">
.codemirror-editor-wrapper {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  background: #fff;
  border-radius: 4px;
  overflow: hidden;

  &.is-fullscreen {
    position: fixed;
    inset: 0;
    width: 100vw !important;
    height: 100vh !important;
    z-index: 9999;
    border-radius: 0;
  }

  .editor-toolbar {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 4px 10px;
    background: #f5f7fa;
    border-bottom: 1px solid #e4e7ed;
    flex-wrap: wrap;

    .toolbar-btn {
      display: flex;
      align-items: center;
      justify-content: center;
      min-width: 28px;
      height: 28px;
      padding: 0 6px;
      background: transparent;
      border: none;
      border-radius: 4px;
      color: #606266;
      cursor: pointer;
      font-size: 13px;
      font-weight: 600;
      transition: all 0.2s;

      i {
        font-size: 15px;
      }

      &:hover {
        background: #e4e7ed;
        color: #409eff;
      }

      &.active {
        background: #409eff;
        color: #fff;
      }

      &.mobile-only {
        display: none;

        @media (max-width: 768px) {
          display: flex;
        }
      }
    }

    .toolbar-divider {
      width: 1px;
      height: 16px;
      background: #dcdfe6;
      margin: 0 4px;
    }

    .toolbar-spacer {
      flex: 1;
      min-width: 12px;
    }
  }

  .editor-container {
    flex: 1;
    display: flex;
    position: relative;
    overflow: hidden;

    .editor-pane {
      flex: 1;
      overflow: auto;
      border-right: 1px solid #e4e7ed;
      cursor: text;
      display: flex;
      flex-direction: column;
      min-height: 0;

      &.full-width {
        border-right: none;
      }

      &.hidden {
        display: none;
      }

      .cm-host {
        flex: 1;
        min-height: 0;
        display: flex;
      }

      :deep(.cm-editor) {
        width: 100%;
        flex: 1;
        min-height: 0;
        display: flex;
        flex-direction: column;
        font-size: 14px;
        font-family: Consolas, Monaco, monospace;

        &.cm-focused {
          outline: none;
        }

        .cm-content {
          padding: 16px;
          padding-bottom: max(16px, calc(100vh - 150px));
          min-height: 100%;
          box-sizing: border-box;
        }

        .cm-line {
          line-height: 1.6;
        }

        .cm-cursor {
          border-left-color: #409eff;
        }

        .cm-selectionBackground {
          background: #409eff33 !important;
        }

        .cm-activeLine {
          background: #f5f7fa;
        }

        .cm-gutters {
          background: #fafafa;
          border-right: 1px solid #e4e7ed;
        }
      }
    }

    .preview-pane {
      flex: 1;
      overflow: auto;
      padding: 20px;
      padding-bottom: max(20px, calc(100vh - 150px));

      &.html-mode {
        padding: 0;
        background: #282c34;

        pre {
          margin: 0;
          padding: 20px;
          height: 100%;

          code {
            color: #abb2bf;
            font-family: Consolas, Monaco, monospace;
            font-size: 14px;
            line-height: 1.6;
            white-space: pre-wrap;
            word-break: break-all;
          }
        }
      }

      // Mermaid 图表样式
      :deep(.markdown-content) {
        .mermaid {
          display: flex;
          justify-content: center;
          align-items: center;
          margin: 1.5rem 0;
          padding: 1rem;
          background: #f5f7fa;
          border-radius: 8px;
          overflow-x: auto;

          svg {
            max-width: 100%;
            height: auto;
          }
        }

        .mermaid-error {
          color: #f56c6c;
          padding: 1rem;
          background: #fef0f0;
          border-radius: 4px;
          border-left: 4px solid #f56c6c;
        }

        img.preview-collapsible-image {
          max-height: 160px;
          width: auto;
          cursor: zoom-in;
          transition:
            max-height 0.2s ease,
            transform 0.2s ease;
        }

        img.preview-collapsible-image.is-expanded {
          max-height: none;
          cursor: zoom-out;
        }
      }
    }

    .toc-pane {
      position: absolute;
      right: 0;
      top: 0;
      bottom: 0;
      width: 260px;
      background: #fff;
      border-left: 1px solid #e4e7ed;
      display: flex;
      flex-direction: column;
      box-shadow: -2px 0 8px rgba(0, 0, 0, 0.1);
      z-index: 10;

      .toc-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 12px 16px;
        border-bottom: 1px solid #e4e7ed;
        background: #f5f7fa;
        font-weight: 600;
        font-size: 14px;
        color: #303133;

        .toc-close {
          display: flex;
          align-items: center;
          justify-content: center;
          width: 24px;
          height: 24px;
          border: none;
          border-radius: 4px;
          background: transparent;
          color: #909399;
          cursor: pointer;
          transition: all 0.2s;

          &:hover {
            background: #e4e7ed;
            color: #606266;
          }

          i {
            font-size: 18px;
          }
        }
      }

      .toc-content {
        flex: 1;
        overflow: auto;
        padding: 12px 0;

        .toc-item {
          padding: 8px 16px;
          cursor: pointer;
          font-size: 14px;
          line-height: 1.5;
          color: #606266;
          border-left: 3px solid transparent;
          transition: all 0.2s;

          &:hover {
            background: #f5f7fa;
            color: #409eff;
            border-left-color: #409eff;
          }

          @for $i from 1 through 6 {
            &.toc-level-#{$i} {
              padding-left: 16px + ($i - 1) * 12px;

              @if $i ==1 {
                font-weight: 600;
              }
            }
          }
        }

        .toc-empty {
          padding: 40px 16px;
          text-align: center;
          color: #909399;
          font-size: 14px;
        }
      }
    }
  }

  .editor-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 4px 12px;
    background: #fafafa;
    border-top: 1px solid #e4e7ed;
    font-size: 12px;
    color: #909399;

    .footer-left {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .word-count,
    .reading-time {
      user-select: none;
    }
  }

  // 移动端适配
  @media (max-width: 767px) {
    .editor-toolbar {
      padding: 4px 8px;
      gap: 2px;
      overflow-x: auto;
      flex-wrap: nowrap;
      -webkit-overflow-scrolling: touch;

      &::-webkit-scrollbar {
        height: 3px;
      }

      .toolbar-btn {
        min-width: 32px;
        height: 32px;
        flex-shrink: 0;
      }

      .toolbar-divider {
        flex-shrink: 0;
      }
    }

    .editor-container {
      .editor-pane {
        :deep(.cm-editor) {
          .cm-content {
            padding: 12px;
            padding-bottom: max(12px, calc(100vh - 200px));
          }
        }
      }

      .preview-pane {
        padding: 12px;
        padding-bottom: max(12px, calc(100vh - 200px));
      }

      .toc-pane {
        width: 100%;
        max-width: 280px;
      }
    }

    .editor-footer {
      padding: 4px 8px;
      font-size: 11px;

      .footer-left {
        gap: 8px;
      }
    }
  }

  // 平板端适配
  @media (min-width: 768px) and (max-width: 991px) {
    .editor-toolbar {
      padding: 4px 10px;
      gap: 3px;
    }

    .editor-container {
      .toc-pane {
        width: 240px;
      }
    }
  }
}

// 表情选择器样式
.emoji-tip {
  padding: 40px 20px;
  text-align: center;
  color: #909399;
  font-size: 0.85rem;
}

.emoji-wrap {
  display: flex;
  flex-direction: column;
  height: 200px;
}

.emoji-bar {
  display: flex;
  border-bottom: 1px solid #eee;
  flex-shrink: 0;
}

.emoji-tab {
  flex: 1;
  padding: 8px 4px;
  border: none;
  background: transparent;
  color: #666;
  font-size: 0.75rem;
  cursor: pointer;

  &:hover {
    background: #f5f5f5;
  }

  &.active {
    color: #409eff;
  }
}

.emoji-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;

  &::-webkit-scrollbar {
    width: 0;
  }
}

.emoji-group {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 2px;

  &.emoji-text {
    display: flex;
    flex-wrap: wrap;
    gap: 4px;
  }
}

.emoji-btn {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  border-radius: 4px;
  cursor: pointer;
  padding: 2px;
  overflow: hidden;

  span {
    font-size: 1.4rem;
  }

  img {
    width: 32px;
    height: 32px;
  }

  &:hover {
    background: #f0f0f0;
  }

  .emoji-text & {
    width: auto;
    height: auto;
    padding: 6px 10px;

    span {
      font-size: 0.85rem;
      white-space: nowrap;
      overflow: hidden;
      max-width: 100%;
    }
  }
}

// 提示框弹窗样式
.note-dialog-wrap {
  padding: 4px 0;

  .note-form-item {
    margin-bottom: 12px;

    &:last-of-type {
      margin-bottom: 16px;
    }
  }

  .note-form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
}

// 标签页弹窗样式
.tabs-dialog-wrap {
  padding: 4px 0;

  .tabs-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 12px;
    max-height: 240px;
    overflow-y: auto;
  }

  .tabs-item {
    display: flex;
    align-items: center;
    gap: 8px;

    .el-input {
      flex: 1;
    }

    .el-button {
      flex-shrink: 0;
      padding: 4px;
    }
  }

  .tabs-footer {
    display: flex;
    justify-content: space-between;
    gap: 8px;
  }
}

// 折叠面板弹窗样式
.fold-dialog-wrap {
  padding: 4px 0;

  .fold-form-item {
    margin-bottom: 12px;

    &:last-of-type {
      margin-bottom: 16px;
    }
  }

  .fold-form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
}

// 链接卡片弹窗样式
.link-dialog-wrap {
  padding: 4px 0;

  .link-form-item {
    margin-bottom: 12px;

    &:last-of-type {
      margin-bottom: 16px;
    }
  }

  .link-form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
}

// 视频弹窗样式
.video-dialog-wrap {
  padding: 4px 0;

  .video-form-item {
    margin-bottom: 12px;

    &:last-of-type {
      margin-bottom: 16px;
    }
  }

  .video-url-preview {
    margin-top: 8px;
    font-size: 11px;
    color: #909399;
    word-break: break-all;
    max-height: 40px;
    overflow-y: auto;
  }

  .video-form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
}

// 音频弹窗样式
.audio-dialog-wrap {
  padding: 4px 0;

  .audio-form-item {
    margin-bottom: 12px;

    &:last-of-type {
      margin-bottom: 16px;
    }
  }

  .audio-url-preview {
    margin-top: 8px;
    font-size: 11px;
    color: #909399;
    word-break: break-all;
    max-height: 40px;
    overflow-y: auto;
  }

  .music-info-preview {
    margin-top: 8px;
    padding: 8px;
    background: #f5f7fa;
    border-radius: 4px;

    .music-info-title {
      font-size: 13px;
      font-weight: 500;
      margin-bottom: 4px;
    }

    .music-info-artist {
      font-size: 12px;
      color: #909399;
    }
  }

  .audio-form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
  }
}

// 照片墙弹窗样式
.photo-dialog-wrap {
  padding: 4px 0;
  max-height: 400px;
  overflow-y: auto;

  .photo-rows {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-bottom: 12px;
  }

  .photo-row {
    border: 1px solid #eee;
    border-radius: 4px;
    padding: 8px;
  }

  .photo-row-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    font-size: 12px;
    color: #909399;
  }

  .photo-row-actions {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  .photo-images {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .photo-image-item {
    display: flex;
    align-items: center;
    gap: 6px;

    .el-input {
      flex: 1;
    }

    .el-button {
      flex-shrink: 0;
      padding: 4px;
    }
  }

  .photo-image-actions {
    display: flex;
    align-items: center;
    gap: 2px;
  }

  .photo-footer {
    display: flex;
    justify-content: space-between;
    gap: 8px;
  }
}
</style>
