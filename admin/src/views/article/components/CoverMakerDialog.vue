<template>
  <Teleport to="body">
    <Transition name="fade">
      <div v-if="visible" class="picprose-overlay" @click="handleClose">
        <div class="picprose-editor" @click.stop>
          <header class="editor-header">
            <span class="app-title">制作封面</span>
            <div class="header-right">
              <el-dropdown
                trigger="click"
                @command="handleExportCommand"
                split-button
                @click="handleExportCommand('apply')"
                :popper-class="'export-dropdown'"
              >
                应用
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="apply-export">应用并导出</el-dropdown-item>
                    <el-dropdown-item command="export-only">仅导出</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <el-button size="small" :icon="Close" @click="handleClose" circle />
            </div>
          </header>

          <main class="editor-main">
            <aside class="toolbar">
              <div class="toolbar-top">
                <el-tabs v-model="imageSource" size="small">
                  <el-tab-pane label="Unsplash" name="unsplash" />
                  <el-tab-pane label="Pixabay" name="pixabay" />
                  <el-tab-pane label="Pexels" name="pexels" />
                  <el-tab-pane label="上传" name="upload" />
                </el-tabs>
              </div>

              <div class="toolbar-middle">
                <template v-if="imageSource !== 'upload'">
                  <div class="photos-grid" ref="photosGridRef">
                    <div v-if="loadingPhotos" class="loading-state">
                      <el-icon class="is-loading">
                        <Loading />
                      </el-icon>
                      <span>加载中...</span>
                    </div>
                    <div v-else class="photos-list">
                      <div
                        v-for="photo in displayedPhotos"
                        :key="photo.id"
                        class="photo-item"
                        :class="{ active: selectedPhoto?.id === photo.id }"
                        @click="selectPhoto(photo)"
                      >
                        <img :src="photo.thumbnail" alt="" loading="lazy" />
                      </div>
                    </div>
                    <div v-if="loadingMore" class="loading-more">
                      <el-icon class="is-loading">
                        <Loading />
                      </el-icon>
                      <span>加载更多...</span>
                    </div>
                  </div>
                </template>
                <template v-else>
                  <div class="upload-section">
                    <el-upload
                      :auto-upload="false"
                      :show-file-list="false"
                      accept="image/*"
                      :on-change="handleImageUpload"
                      drag
                    >
                      <div class="upload-area">
                        <el-icon class="upload-icon">
                          <Plus />
                        </el-icon>
                        <span>拖拽或点击上传</span>
                      </div>
                    </el-upload>
                  </div>
                </template>
              </div>

              <div v-if="imageSource !== 'upload'" class="toolbar-bottom">
                <el-input
                  v-model="searchQuery"
                  placeholder="搜索图片..."
                  @keyup.enter="searchPhotos"
                  clearable
                >
                  <template #append>
                    <el-button :icon="Search" @click="searchPhotos" />
                  </template>
                </el-input>
              </div>
            </aside>

            <section class="canvas-area">
              <div class="canvas-container" ref="canvasContainerRef">
                <div v-if="!imageLoaded" class="empty-state">
                  <el-icon class="empty-icon">
                    <Picture />
                  </el-icon>
                  <p>上传图片开始创作</p>
                </div>
                <div v-else class="image-canvas" ref="canvasRef" :style="canvasStyle">
                  <img :src="imageUrl" alt="" class="base-image" />
                  <div class="image-overlay" :style="{ opacity: overlayOpacity / 100 }"></div>
                  <div v-if="isDragging" class="guide-lines">
                    <div class="guide-line vertical-center"></div>
                    <div class="guide-line horizontal-center"></div>
                  </div>
                  <div
                    v-if="textElements.subtitle.text"
                    class="text-element subtitle-element"
                    :style="getSubtitleStyle()"
                    @mousedown="startDragElement('subtitle', $event)"
                  >
                    {{ textElements.subtitle.text }}
                  </div>
                  <div
                    v-if="textElements.title.text"
                    class="text-element title-element"
                    :style="getTitleStyle()"
                    @mousedown="startDragElement('title', $event)"
                  >
                    {{ textElements.title.text }}
                  </div>
                  <div
                    v-if="textElements.author.text"
                    class="text-element author-element"
                    :style="getAuthorStyle()"
                    @mousedown="startDragElement('author', $event)"
                  >
                    {{ textElements.author.text }}
                  </div>
                  <div
                    v-if="textElements.avatar.src"
                    class="avatar-element"
                    :style="getAvatarStyle()"
                    @mousedown="startDragElement('avatar', $event)"
                  >
                    <img :src="textElements.avatar.src" :style="getAvatarImageStyle()" />
                  </div>
                </div>
              </div>
            </section>

            <aside class="properties-panel">
              <div class="property-group">
                <label>遮罩层浓度</label>
                <div style="display: flex; align-items: center; gap: 12px">
                  <el-slider
                    v-model="overlayOpacity"
                    :min="0"
                    :max="100"
                    :step="5"
                    style="flex: 1"
                  />
                  <span style="min-width: 45px; text-align: right; color: #666"
                    >{{ overlayOpacity }}%</span
                  >
                </div>
              </div>
              <div class="property-group">
                <label>标题</label>
                <el-input v-model="textElements.title.text" placeholder="请输入标题" />
              </div>
              <div class="property-group">
                <label>副标题</label>
                <el-input v-model="textElements.subtitle.text" placeholder="请输入副标题" />
              </div>
              <div class="property-group">
                <label>作者</label>
                <el-input v-model="textElements.author.text" placeholder="请输入作者名称" />
              </div>
              <div class="property-group">
                <label>头像</label>
                <el-upload
                  action="#"
                  :show-file-list="false"
                  :before-upload="handleAvatarUpload"
                  accept="image/*"
                >
                  <el-button type="primary">选择图片</el-button>
                </el-upload>
              </div>
            </aside>
          </main>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onUnmounted } from 'vue';
import { ElMessage } from 'element-plus';
import { Close, Plus, Picture, Search, Loading } from '@element-plus/icons-vue';

interface TextElement {
  text: string;
  x: number;
  y: number;
  fontSize: number;
  fontFamily: string;
  color: string;
}

interface AvatarElement {
  src: string;
  x: number;
  y: number;
  size: number;
}

interface PlatformPhoto {
  id: string;
  url: string;
  thumbnail: string;
}

interface Props {
  modelValue: boolean;
  title?: string;
  author?: string;
  avatar?: string;
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void;
  (e: 'save', imageUrl: string): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

const visible = ref(false);
const canvasRef = ref<HTMLElement>();
const photosGridRef = ref<HTMLElement>();
const imageSource = ref<'unsplash' | 'pixabay' | 'pexels' | 'upload'>('unsplash');
const searchQuery = ref('');
const platformPhotos = ref<PlatformPhoto[]>([]);
const selectedPhoto = ref<PlatformPhoto | null>(null);
const loadingPhotos = ref(false);
const loadingMore = ref(false);
const hasMorePhotos = ref(true);
const currentPage = ref(1);
const perPage = ref(20);
const imageUrl = ref('');
const imageLoaded = ref(false);
const overlayOpacity = ref(60);

const textElements = ref({
  title: {
    text: '',
    x: 50,
    y: 40,
    fontSize: 148,
    fontFamily: 'SlidefontKBK',
    color: '#ffffff',
  } as TextElement,
  subtitle: {
    text: '',
    x: 50,
    y: 25,
    fontSize: 78,
    fontFamily: 'SlidefontKBK',
    color: '#ffffff',
  } as TextElement,
  author: {
    text: '',
    x: 50,
    y: 57.5,
    fontSize: 95,
    fontFamily: 'SlidefontKBK',
    color: '#ffffff',
  } as TextElement,
  avatar: {
    src: '',
    x: 50,
    y: 75,
    size: 120,
  } as AvatarElement,
});

const dragStart = ref({ x: 0, y: 0, elementX: 0, elementY: 0 });
const isDragging = ref(false);
const currentDragElement = ref<'title' | 'subtitle' | 'author' | 'avatar' | null>(null);
const canvasContainerRef = ref<HTMLElement>();

const displayedPhotos = computed(() => {
  return platformPhotos.value.slice(0, currentPage.value * perPage.value);
});

// 计算画布缩放比例
const canvasScale = computed(() => {
  if (!canvasContainerRef.value) return 1;
  const container = canvasContainerRef.value;
  const containerWidth = container.clientWidth;
  const containerHeight = container.clientHeight;
  const scaleX = containerWidth / 1920;
  const scaleY = containerHeight / 1080;
  return Math.min(scaleX, scaleY);
});

// 获取画布样式（包含缩放）
const canvasStyle = computed(() => {
  return {
    transform: `translate(-50%, -50%) scale(${canvasScale.value})`,
  };
});

const IMAGE_API_URL = 'https://pixhub.flec.top';

function handleClose() {
  visible.value = false;
  emit('update:modelValue', false);
}

function handleAvatarUpload(file: File) {
  const reader = new FileReader();
  reader.onload = e => {
    textElements.value.avatar.src = e.target?.result as string;
  };
  reader.readAsDataURL(file);
  return false;
}

function getCanvasSize() {
  // 固定画布逻辑尺寸为 1920x1080，确保在所有设备上元素位置一致
  return { width: 1920, height: 1080 };
}

function getTitleStyle() {
  const title = textElements.value.title;
  const canvasSize = getCanvasSize();
  return {
    position: 'absolute' as const,
    left: `${(title.x / 100) * canvasSize.width}px`,
    top: `${(title.y / 100) * canvasSize.height}px`,
    fontSize: `${title.fontSize}px`,
    fontFamily: title.fontFamily,
    color: '#ffffff',
    cursor: 'move' as const,
    userSelect: 'none' as const,
    transform: 'translate(-50%, -50%)',
    textAlign: 'center' as const,
    whiteSpace: 'nowrap' as const,
  };
}

function getSubtitleStyle() {
  const subtitle = textElements.value.subtitle;
  const canvasSize = getCanvasSize();
  return {
    position: 'absolute' as const,
    left: `${(subtitle.x / 100) * canvasSize.width}px`,
    top: `${(subtitle.y / 100) * canvasSize.height}px`,
    fontSize: `${subtitle.fontSize}px`,
    fontFamily: subtitle.fontFamily,
    color: '#ffffff',
    cursor: 'move' as const,
    userSelect: 'none' as const,
    transform: 'translate(-50%, -50%)',
    textAlign: 'center' as const,
    whiteSpace: 'nowrap' as const,
  };
}

function getAuthorStyle() {
  const author = textElements.value.author;
  const canvasSize = getCanvasSize();
  return {
    position: 'absolute' as const,
    left: `${(author.x / 100) * canvasSize.width}px`,
    top: `${(author.y / 100) * canvasSize.height}px`,
    fontSize: `${author.fontSize}px`,
    fontFamily: author.fontFamily,
    color: '#ffffff',
    cursor: 'move' as const,
    userSelect: 'none' as const,
    transform: 'translate(-50%, -50%)',
    textAlign: 'center' as const,
    whiteSpace: 'nowrap' as const,
  };
}

function getAvatarStyle() {
  const avatar = textElements.value.avatar;
  const canvasSize = getCanvasSize();
  return {
    position: 'absolute' as const,
    left: `${(avatar.x / 100) * canvasSize.width}px`,
    top: `${(avatar.y / 100) * canvasSize.height}px`,
    width: `${avatar.size}px`,
    height: `${avatar.size}px`,
    cursor: 'move' as const,
    transform: 'translate(-50%, -50%)',
  };
}

function getAvatarImageStyle() {
  return {
    width: '100%',
    height: '100%',
    objectFit: 'cover' as const,
    borderRadius: '50%',
    border: '4px solid #ffffff',
  };
}

function startDragElement(element: 'title' | 'subtitle' | 'author' | 'avatar', e: MouseEvent) {
  e.preventDefault();
  e.stopPropagation();
  isDragging.value = true;
  currentDragElement.value = element;
  const el = textElements.value[element];
  dragStart.value = {
    x: e.clientX,
    y: e.clientY,
    elementX: el.x,
    elementY: el.y,
  };
  document.addEventListener('mousemove', handleDragElement);
  document.addEventListener('mouseup', stopDragElement);
}

function handleDragElement(e: MouseEvent) {
  if (!isDragging.value || !currentDragElement.value) return;
  const canvasSize = getCanvasSize();
  const scale = canvasScale.value;
  const dx = (e.clientX - dragStart.value.x) / scale;
  const dy = (e.clientY - dragStart.value.y) / scale;
  const dxPercent = (dx / canvasSize.width) * 100;
  const dyPercent = (dy / canvasSize.height) * 100;
  const element = textElements.value[currentDragElement.value];
  const gridSize = 2.5;
  let newX = dragStart.value.elementX + dxPercent;
  let newY = dragStart.value.elementY + dyPercent;
  newX = Math.round(newX / gridSize) * gridSize;
  newY = Math.round(newY / gridSize) * gridSize;
  element.x = Math.max(0, Math.min(100, newX));
  element.y = Math.max(0, Math.min(100, newY));
}

function stopDragElement() {
  isDragging.value = false;
  currentDragElement.value = null;
  document.removeEventListener('mousemove', handleDragElement);
  document.removeEventListener('mouseup', stopDragElement);
}

async function searchPhotos() {
  loadingPhotos.value = true;
  currentPage.value = 1;
  platformPhotos.value = [];
  try {
    const params = new URLSearchParams({
      platform: imageSource.value,
      page_size: String(perPage.value),
      page: '1',
      ...(searchQuery.value.trim() && { query: searchQuery.value }),
    });
    const response = await fetch(`${IMAGE_API_URL}/?${params}`);
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const data = await response.json();
    const photos = data.results || [];
    if (photos.length > 0) {
      platformPhotos.value = photos;
      hasMorePhotos.value = photos.length >= perPage.value;
      selectPhoto(photos[0]);
    } else {
      ElMessage.warning('没有找到相关图片');
    }
  } catch (_error) {
    ElMessage.error('搜索图片失败');
  } finally {
    loadingPhotos.value = false;
  }
}

async function loadMorePhotos() {
  if (loadingMore.value || !hasMorePhotos.value) return;
  loadingMore.value = true;
  const nextPage = currentPage.value + 1;
  try {
    const params = new URLSearchParams({
      platform: imageSource.value,
      page_size: String(perPage.value),
      page: String(nextPage),
      ...(searchQuery.value.trim() && { query: searchQuery.value }),
    });
    const response = await fetch(`${IMAGE_API_URL}/?${params}`);
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const data = await response.json();
    const photos = data.results || [];
    if (photos.length > 0) {
      platformPhotos.value.push(...photos);
      currentPage.value = nextPage;
      hasMorePhotos.value = photos.length >= perPage.value;
    }
  } catch (_error) {
    ElMessage.error('加载更多图片失败');
  } finally {
    loadingMore.value = false;
  }
}

function selectPhoto(photo: PlatformPhoto) {
  selectedPhoto.value = photo;
  imageUrl.value = photo.url || '';
  imageLoaded.value = true;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
function handleImageUpload(file: any) {
  const reader = new FileReader();
  reader.onload = e => {
    imageUrl.value = e.target?.result as string;
    imageLoaded.value = true;
    selectedPhoto.value = null;
  };
  reader.readAsDataURL(file.raw);
}

async function generateImageDataUrl() {
  if (!canvasRef.value || !imageLoaded.value) {
    ElMessage.warning('请先上传图片');
    return null;
  }
  try {
    const previewCanvas = canvasRef.value;
    const titleElement = previewCanvas.querySelector('.title-element') as HTMLElement;
    const subtitleElement = previewCanvas.querySelector('.subtitle-element') as HTMLElement;
    const authorElement = previewCanvas.querySelector('.author-element') as HTMLElement;
    const avatarElement = previewCanvas.querySelector('.avatar-element') as HTMLElement;
    const canvas = document.createElement('canvas');
    canvas.width = 1920;
    canvas.height = 1080;
    const ctx = canvas.getContext('2d');
    if (!ctx) throw new Error('无法创建 Canvas 上下文');
    // 预览画布现在是固定的 1920x1080，所以缩放比例是 1:1
    const img = new Image();
    img.crossOrigin = 'anonymous';
    await new Promise((resolve, reject) => {
      img.onload = resolve;
      img.onerror = reject;
      img.src = imageUrl.value;
    });
    const imgAspect = img.width / img.height;
    const canvasAspect = canvas.width / canvas.height;
    let drawWidth, drawHeight, drawX, drawY;
    if (imgAspect > canvasAspect) {
      drawHeight = canvas.height;
      drawWidth = drawHeight * imgAspect;
      drawX = (canvas.width - drawWidth) / 2;
      drawY = 0;
    } else {
      drawWidth = canvas.width;
      drawHeight = drawWidth / imgAspect;
      drawX = 0;
      drawY = (canvas.height - drawHeight) / 2;
    }
    ctx.drawImage(img, drawX, drawY, drawWidth, drawHeight);
    ctx.fillStyle = `rgba(31, 41, 55, ${overlayOpacity.value / 100})`;
    ctx.fillRect(0, 0, canvas.width, canvas.height);
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    if (textElements.value.title.text && titleElement) {
      const computedStyle = window.getComputedStyle(titleElement);
      const actualFontSize = parseFloat(computedStyle.fontSize);
      ctx.font = `${actualFontSize}px 'SlidefontKBK'`;
      ctx.fillStyle = '#ffffff';
      const titleRect = titleElement.getBoundingClientRect();
      const previewRect = previewCanvas.getBoundingClientRect();
      const titleX = (titleRect.left + titleRect.width / 2 - previewRect.left) / canvasScale.value;
      const titleY = (titleRect.top + titleRect.height / 2 - previewRect.top) / canvasScale.value;
      ctx.fillText(textElements.value.title.text, titleX, titleY);
    }
    if (textElements.value.subtitle.text && subtitleElement) {
      const computedStyle = window.getComputedStyle(subtitleElement);
      const actualFontSize = parseFloat(computedStyle.fontSize);
      ctx.font = `${actualFontSize}px 'SlidefontKBK'`;
      ctx.fillStyle = '#ffffff';
      const subtitleRect = subtitleElement.getBoundingClientRect();
      const previewRect = previewCanvas.getBoundingClientRect();
      const subtitleX =
        (subtitleRect.left + subtitleRect.width / 2 - previewRect.left) / canvasScale.value;
      const subtitleY =
        (subtitleRect.top + subtitleRect.height / 2 - previewRect.top) / canvasScale.value;
      ctx.fillText(textElements.value.subtitle.text, subtitleX, subtitleY);
    }
    if (textElements.value.author.text && authorElement) {
      const computedStyle = window.getComputedStyle(authorElement);
      const actualFontSize = parseFloat(computedStyle.fontSize);
      ctx.font = `${actualFontSize}px 'SlidefontKBK'`;
      ctx.fillStyle = '#ffffff';
      const authorRect = authorElement.getBoundingClientRect();
      const previewRect = previewCanvas.getBoundingClientRect();
      const authorX =
        (authorRect.left + authorRect.width / 2 - previewRect.left) / canvasScale.value;
      const authorY =
        (authorRect.top + authorRect.height / 2 - previewRect.top) / canvasScale.value;
      ctx.fillText(textElements.value.author.text, authorX, authorY);
    }
    if (textElements.value.avatar.src && avatarElement) {
      const avatarImg = new Image();
      avatarImg.crossOrigin = 'anonymous';
      await new Promise(resolve => {
        avatarImg.onload = resolve;
        avatarImg.onerror = () => resolve(null);
        avatarImg.src = textElements.value.avatar.src;
      });
      if (avatarImg.complete) {
        const avatarRect = avatarElement.getBoundingClientRect();
        const previewRect = previewCanvas.getBoundingClientRect();
        const avatarX =
          (avatarRect.left + avatarRect.width / 2 - previewRect.left) / canvasScale.value;
        const avatarY =
          (avatarRect.top + avatarRect.height / 2 - previewRect.top) / canvasScale.value;
        const avatarRadius = avatarRect.width / 2 / canvasScale.value;
        ctx.save();
        ctx.beginPath();
        ctx.arc(avatarX, avatarY, avatarRadius, 0, Math.PI * 2);
        ctx.closePath();
        ctx.clip();
        ctx.drawImage(
          avatarImg,
          avatarX - avatarRadius,
          avatarY - avatarRadius,
          avatarRadius * 2,
          avatarRadius * 2
        );
        ctx.restore();
        ctx.beginPath();
        ctx.arc(avatarX, avatarY, avatarRadius, 0, Math.PI * 2);
        ctx.strokeStyle = '#ffffff';
        ctx.lineWidth = 2;
        ctx.stroke();
      }
    }
    return canvas.toDataURL('image/png', 0.95);
  } catch (_error) {
    ElMessage.error('生成图片失败');
    return null;
  }
}

const handleExportCommand = async (command: 'apply' | 'apply-export' | 'export-only') => {
  const dataUrl = await generateImageDataUrl();
  if (!dataUrl) return;
  if (command === 'apply') {
    emit('save', dataUrl);
    handleClose();
    ElMessage.success('封面已应用到文章');
  } else if (command === 'apply-export') {
    emit('save', dataUrl);
    handleClose();
    ElMessage.success('封面已应用到文章');
    const link = document.createElement('a');
    link.href = dataUrl;
    link.download = `cover-${Date.now()}.png`;
    link.click();
    ElMessage.success('图片已导出');
  } else if (command === 'export-only') {
    const link = document.createElement('a');
    link.href = dataUrl;
    link.download = `cover-${Date.now()}.png`;
    link.click();
    ElMessage.success('图片已导出');
  }
};

watch(
  () => props.modelValue,
  val => {
    visible.value = val;
  }
);

watch(visible, val => {
  emit('update:modelValue', val);
  if (val) {
    if (props.title) textElements.value.title.text = props.title;
    if (props.author) textElements.value.author.text = props.author;
    if (props.avatar) {
      if (import.meta.env.DEV) {
        const apiBaseUrl = import.meta.env.VITE_API_URL;
        const backendBaseUrl = apiBaseUrl.replace(/\/api\/v\d+$/, '');
        if (props.avatar.startsWith(backendBaseUrl)) {
          textElements.value.avatar.src = props.avatar.replace(backendBaseUrl, '');
        } else {
          textElements.value.avatar.src = props.avatar;
        }
      } else {
        textElements.value.avatar.src = props.avatar;
      }
    }
    if (imageSource.value !== 'upload' && platformPhotos.value.length === 0) {
      searchPhotos();
    }
    nextTick(() => {
      setupScrollListener();
    });
  }
});

watch(imageSource, val => {
  if (val !== 'upload') {
    searchPhotos();
  }
});

function setupScrollListener() {
  const grid = photosGridRef.value;
  if (!grid) return;
  const handleScroll = () => {
    const { scrollTop, scrollHeight, clientHeight } = grid;
    if (scrollTop + clientHeight >= scrollHeight - 100) {
      if (hasMorePhotos.value && !loadingMore.value && !loadingPhotos.value) {
        loadMorePhotos();
      }
    }
  };
  grid.addEventListener('scroll', handleScroll, { passive: true });
}

onUnmounted(() => {
  if (textElements.value.avatar.src && textElements.value.avatar.src.startsWith('blob:')) {
    URL.revokeObjectURL(textElements.value.avatar.src);
  }
});
</script>

<style scoped lang="scss">
@font-face {
  font-family: 'SlidefontKBK';
  src: url('@/assets/font/SlidefontKBK.woff2') format('woff2');
  font-display: swap;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.picprose-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 3000;
  padding: 20px;
}

.picprose-editor {
  width: 100%;
  height: 100%;
  max-width: 1400px;
  max-height: 900px;
  background: #ffffff;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
}

.editor-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20px 24px;
  flex-shrink: 0;

  .app-title {
    font-size: 18px;
    font-weight: 500;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }
}

.editor-main {
  flex: 1;
  display: flex;
  background: #f8f9fa;
  height: 0;
}

.toolbar {
  width: 300px;
  background: #ffffff;
  border-right: 1px solid #e9ecef;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  height: 100%;

  .toolbar-top {
    flex-shrink: 0;
    padding: 6px 16px;
  }

  .toolbar-middle {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .toolbar-bottom {
    flex-shrink: 0;
    padding: 12px;
  }

  .photos-grid {
    flex: 1;
    overflow-y: auto;
    padding: 0 6px;
    margin: 0 6px;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-thumb {
      background: #c1c1c1;
      border-radius: 3px;
    }

    .loading-state {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 40px 0;
      color: #6c757d;

      .el-icon {
        font-size: 24px;
        margin-bottom: 8px;
      }

      span {
        font-size: 13px;
      }
    }

    .photos-list {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 8px;

      .photo-item {
        position: relative;
        aspect-ratio: 16/10;
        border-radius: 6px;
        overflow: hidden;
        cursor: pointer;
        transition: all 0.3s ease;
        background: #f0f0f0;

        &:hover {
          transform: scale(1.05);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        }

        &.active {
          border: 2px solid #409eff;
          box-shadow: 0 0 0 4px rgba(64, 158, 255, 0.2);
        }

        img {
          width: 100%;
          height: 100%;
          object-fit: cover;
          display: block;
        }
      }
    }

    .loading-more {
      padding: 16px;
      text-align: center;
      color: #adb5bd;
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
    }
  }

  .upload-section {
    flex: 1;
    padding: 0 24px;
    display: flex;
    align-items: center;
    justify-content: center;

    .upload-area {
      border: 2px dashed #dee2e6;
      border-radius: 8px;
      padding: 32px 16px;
      text-align: center;
      cursor: pointer;
      transition: all 0.3s ease;

      &:hover {
        border-color: #409eff;
        background: rgba(64, 158, 255, 0.05);
      }

      .upload-icon {
        font-size: 32px;
        color: #adb5bd;
        margin-bottom: 8px;
      }

      span {
        display: block;
        font-size: 14px;
        color: #6c757d;
      }
    }
  }
}

.canvas-area {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 32px;
  overflow: hidden;
}

.canvas-container {
  width: 100%;
  aspect-ratio: 16 / 9;
  background: #f8f9fa;
  border-radius: 8px;
  position: relative;
  overflow: hidden;

  @supports not (aspect-ratio: 16 / 9) {
    &::before {
      content: '';
      display: block;
      padding-top: 56.25%;
      width: 0;
    }

    > * {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
    }
  }
}

.empty-state {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  color: #adb5bd;

  .empty-icon {
    font-size: 64px;
    margin-bottom: 16px;
  }

  p {
    font-size: 16px;
    margin: 0;
  }
}

.image-canvas {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 1920px;
  height: 1080px;
  transform: translate(-50%, -50%);
  transform-origin: center center;

  .base-image {
    width: 100%;
    height: 100%;
    object-fit: cover;
    display: block;
    border-radius: 4px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  }

  .image-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: #1f2937;
    pointer-events: none;
    z-index: 1;
  }

  .guide-lines {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
    z-index: 10;

    .guide-line {
      position: absolute;
      background: rgba(255, 255, 255, 0.5);

      &.vertical-center {
        left: 50%;
        top: 0;
        bottom: 0;
        width: 1px;
        transform: translateX(-50%);
      }

      &.horizontal-center {
        top: 50%;
        left: 0;
        right: 0;
        height: 1px;
        transform: translateY(-50%);
      }
    }
  }
}

.text-element {
  position: absolute;
  cursor: move;
  transition: all 0.2s ease;
  padding: 8px 16px;
  border-radius: 4px;
  border: 2px solid transparent;
  z-index: 2;

  &:hover {
    border-color: rgba(64, 158, 255, 0.3);
    background: rgba(255, 255, 255, 0.1);
  }
}

.avatar-element {
  position: absolute;
  cursor: move;
  transition: all 0.2s ease;
  border-radius: 50%;
  overflow: hidden;
  border: 3px solid transparent;
  z-index: 2;

  &:hover {
    border-color: rgba(64, 158, 255, 0.5);
    transform: scale(1.05);
  }
}

.properties-panel {
  width: 320px;
  background: #ffffff;
  border-left: 1px solid #e9ecef;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  padding: 24px;

  .property-group {
    margin-bottom: 24px;

    &:last-child {
      margin-bottom: 0;
    }

    label {
      display: block;
      font-size: 12px;
      font-weight: 600;
      color: #172b4d;
      margin-bottom: 8px;
      text-transform: uppercase;
      letter-spacing: 0.5px;
    }
  }
}

:deep(.el-upload-dragger) {
  border: none;
  background: transparent;
  width: 100%;
  height: auto;
  padding: 0;
}

// 修复下拉菜单 z-index 问题
:global(.export-dropdown) {
  z-index: 4000 !important;
}

// 移动端适配
@media (max-width: 767px) {
  .picprose-overlay {
    padding: 0;
  }

  .picprose-editor {
    max-width: 100%;
    max-height: 100%;
    border-radius: 0;
  }

  .editor-header {
    padding: 12px 16px;

    .app-title {
      font-size: 16px;
    }
  }

  .editor-main {
    flex-direction: column;
  }

  .toolbar {
    width: 100%;
    max-height: 40%;
    border-right: none;
    border-bottom: 1px solid #e9ecef;

    .toolbar-top {
      padding: 8px 12px;
    }

    .toolbar-bottom {
      padding: 8px 12px;
    }

    .photos-grid {
      .photos-list {
        grid-template-columns: repeat(3, 1fr);
        gap: 6px;
      }
    }
  }

  .canvas-area {
    flex: 1;
    min-height: 0;

    .canvas-container {
      padding: 12px;
    }

    .empty-state {
      .empty-icon {
        font-size: 48px;
      }

      p {
        font-size: 14px;
      }
    }
  }

  .properties-panel {
    width: 100%;
    max-height: 35%;
    border-left: none;
    border-top: 1px solid #e9ecef;
    padding: 16px;
    overflow-y: auto;

    .property-group {
      margin-bottom: 16px;
    }
  }

  .text-element {
    padding: 6px 12px;
    font-size: 14px;
  }
}

// 平板端适配
@media (min-width: 768px) and (max-width: 991px) {
  .picprose-overlay {
    padding: 12px;
  }

  .picprose-editor {
    max-width: 100%;
    max-height: 100%;
  }

  .toolbar {
    width: 260px;
  }

  .properties-panel {
    width: 280px;
    padding: 20px;
  }
}
</style>
