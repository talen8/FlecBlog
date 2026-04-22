<template>
  <el-dialog
    v-model="visible"
    :title="dialogTitle"
    width="90%"
    style="max-width: 800px"
    :close-on-click-modal="false"
    top="8vh"
  >
    <div class="moment-editor">
      <!-- 功能工具栏 -->
      <div class="toolbar">
        <!-- 链接按钮 -->
        <el-button
          :type="formData.content.link?.url ? 'primary' : 'default'"
          :icon="Link"
          circle
          @click="linkDialogVisible = true"
        />

        <!-- 图片按钮 -->
        <el-button
          :type="imageItems.length ? 'primary' : 'default'"
          :icon="Picture"
          circle
          @click="imageDialogVisible = true"
        />

        <!-- 音乐按钮 -->
        <el-button
          :type="formData.content.music?.id ? 'primary' : 'default'"
          :icon="Headset"
          circle
          @click="musicDialogVisible = true"
        />

        <!-- 视频按钮 -->
        <el-button
          :type="videoItem ? 'primary' : 'default'"
          :icon="VideoPlay"
          circle
          @click="videoDialogVisible = true"
        />
      </div>

      <!-- 文本输入区域 -->
      <div class="content-area">
        <el-input
          v-model="formData.content.text"
          type="textarea"
          placeholder="分享你的想法..."
          :rows="6"
          maxlength="1000"
          show-word-limit
          resize="none"
        />
      </div>

      <!-- 图片预览区域 -->
      <div v-if="imageItems.length" class="images-preview">
        <div v-for="(item, index) in imageItems" :key="item.id" class="image-item">
          <img :src="item.url" @error="handleImageError" />
          <el-button size="small" circle text @click="removeImage(index)">
            <el-icon>
              <Close />
            </el-icon>
          </el-button>
        </div>
      </div>

      <!-- 其他内容预览区域 -->
      <div v-if="otherContentPreviews.length" class="content-preview">
        <div
          v-for="preview in otherContentPreviews"
          :key="preview.type"
          :class="[`${preview.type}-preview`, preview.type === 'video' ? 'video-preview' : '']"
        >
          <!-- 链接类型：显示图标、标题和URL -->
          <template v-if="preview.type === 'link'">
            <div class="content-icon">
              <img
                v-if="preview.favicon"
                :src="preview.favicon"
                alt="网站图标"
                class="favicon-icon"
              />
            </div>
            <div class="content-info">
              <div class="preview-title">{{ preview.title }}</div>
              <div class="preview-url">{{ preview.url }}</div>
            </div>
            <el-button size="small" text @click="removeContent(preview.type)">
              <el-icon>
                <Close />
              </el-icon>
            </el-button>
          </template>
          <!-- 音乐类型：显示封面、歌名和艺术家 -->
          <template v-else-if="preview.type === 'music'">
            <div class="content-icon">
              <img v-if="preview.cover" :src="preview.cover" alt="音乐封面" class="music-cover" />
              <div v-else class="music-cover-placeholder">
                <el-icon>
                  <Headset />
                </el-icon>
              </div>
            </div>
            <div class="content-info">
              <div class="preview-title">
                {{ preview.title }}
                <el-tag size="small" style="margin-left: 8px">{{ preview.musicType }}</el-tag>
              </div>
              <div class="preview-artist">{{ preview.artist }}</div>
            </div>
            <el-button size="small" text @click="removeContent(preview.type)">
              <el-icon>
                <Close />
              </el-icon>
            </el-button>
          </template>
          <!-- 视频类型：使用iframe显示网络视频，使用video标签显示本地视频 -->
          <template v-else-if="preview.type === 'video'">
            <!-- 网络视频预览（iframe） -->
            <div class="video-preview-container" v-if="!preview.isLocal && preview.url">
              <iframe
                :src="preview.url"
                frameborder="0"
                scrolling="no"
                allowfullscreen
                allow="
                  accelerometer;
                  autoplay;
                  clipboard-write;
                  encrypted-media;
                  gyroscope;
                  picture-in-picture;
                "
                class="video-iframe"
              ></iframe>
            </div>
            <!-- 本地视频预览（video标签） -->
            <div class="video-local-preview-container" v-else-if="preview.isLocal && preview.url">
              <video :src="preview.url" controls class="video-local"></video>
            </div>
            <el-button size="small" text @click="removeContent(preview.type)">
              <el-icon>
                <Close />
              </el-icon>
            </el-button>
          </template>
          <!-- 其他类型 -->
          <template v-else>
            <div class="content-info">
              {{ preview.text }}
            </div>
            <el-button size="small" text @click="removeContent(preview.type)">
              <el-icon>
                <Close />
              </el-icon>
            </el-button>
          </template>
        </div>
      </div>

      <!-- 底部工具栏 -->
      <div class="bottom-toolbar">
        <!-- 位置按钮 -->
        <el-button
          :type="formData.content.location ? 'primary' : 'default'"
          :icon="Location"
          text
          @click="locationDialogVisible = true"
        >
          {{ formData.content.location || '发布位置' }}
        </el-button>

        <!-- 标签按钮 -->
        <el-button
          :type="formData.content.tags ? 'primary' : 'default'"
          :icon="PriceTag"
          text
          @click="tagDialogVisible = true"
        >
          {{ formData.content.tags || '分类标签' }}
        </el-button>

        <!-- 发布时间按钮 -->
        <el-button
          :type="publishTime ? 'primary' : 'default'"
          :icon="Timer"
          text
          @click="timeDialogVisible = true"
        >
          {{ publishTime ? formatTime(publishTime) : '发布时间' }}
        </el-button>

        <div class="publish-status">
          <el-radio-group v-model="formData.is_publish" size="small">
            <el-radio-button :value="false">草稿</el-radio-button>
            <el-radio-button :value="true">发布</el-radio-button>
          </el-radio-group>
        </div>
      </div>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">
          {{ formData.is_publish ? '发布' : '保存' }}
        </el-button>
      </span>
    </template>
  </el-dialog>

  <!-- 链接Dialog -->
  <el-dialog v-model="linkDialogVisible" title="网站分享" width="90%" style="max-width: 400px">
    <div class="link-form">
      <div style="display: flex; gap: 8px; margin-bottom: 12px">
        <el-input
          v-model="formData.content.link!.url"
          placeholder="请输入网站地址"
          style="flex: 1"
        />
        <el-button type="primary" :loading="fetchingLink" @click="handleParseLink">
          {{ fetchingLink ? '解析中...' : '解析' }}
        </el-button>
      </div>
      <el-input
        v-model="formData.content.link!.title"
        placeholder="网站标题"
        style="margin-bottom: 12px"
      />
      <div style="display: flex; align-items: center; gap: 8px">
        <el-input v-model="formData.content.link!.favicon" placeholder="网站图标" style="flex: 1" />
        <div class="favicon-preview">
          <img
            v-if="formData.content.link!.favicon"
            :src="formData.content.link!.favicon"
            alt="图标"
          />
          <span v-else>图标</span>
        </div>
      </div>
    </div>
  </el-dialog>

  <!-- 图片Dialog -->
  <el-dialog v-model="imageDialogVisible" title="动态配图" width="90%" style="max-width: 400px">
    <div class="image-form">
      <div style="display: flex; gap: 8px; margin-bottom: 16px">
        <el-input
          v-model="imageUrlInput"
          placeholder="输入图片链接或点击右侧上传"
          @keyup.enter="addImageUrl"
          style="flex: 1"
        />
        <el-button type="primary" @click="addImageUrl" :disabled="!imageUrlInput.trim()"
          >添加</el-button
        >
        <el-button type="primary" @click="handleImageUpload">上传</el-button>
      </div>

      <div v-if="imageItems.length" class="image-url-list">
        <div v-for="(item, index) in imageItems" :key="item.id" class="url-item">
          <el-input :value="item.url" readonly style="flex: 1" />
          <el-button type="danger" size="small" @click="removeImage(index)">删除</el-button>
        </div>
      </div>
    </div>
  </el-dialog>

  <!-- 音乐Dialog -->
  <el-dialog v-model="musicDialogVisible" title="动态音乐" width="90%" style="max-width: 400px">
    <div class="music-form">
      <el-select
        v-model="formData.content.music!.server"
        placeholder="音乐平台"
        style="width: 100%; margin-bottom: 12px"
      >
        <el-option label="网易云音乐" value="netease" />
        <el-option label="QQ音乐" value="tencent" />
        <el-option label="酷狗音乐" value="kugou" />
        <el-option label="虾米音乐" value="xiami" />
        <el-option label="百度音乐" value="baidu" />
        <el-option label="酷我音乐" value="kuwo" />
      </el-select>
      <el-select
        v-model="formData.content.music!.type"
        placeholder="类型"
        style="width: 100%; margin-bottom: 12px"
      >
        <el-option label="单曲" value="song" />
        <el-option label="歌单" value="playlist" />
        <el-option label="艺术家" value="artist" />
        <el-option label="专辑" value="album" />
        <el-option label="搜索" value="search" />
      </el-select>
      <div style="display: flex; gap: 8px; margin-bottom: 12px">
        <el-input v-model="formData.content.music!.id" placeholder="音乐ID" style="flex: 1" />
        <el-button type="primary" :loading="fetchingMusic" @click="handleParseMusic">
          {{ fetchingMusic ? '解析中...' : '解析' }}
        </el-button>
      </div>

      <!-- 音乐信息预览 -->
      <div v-if="musicInfo" class="music-info-preview">
        <img v-if="musicInfo.pic" :src="musicInfo.pic" alt="音乐封面" class="music-preview-cover" />
        <div class="music-preview-info">
          <div class="music-preview-title">
            {{ musicInfo.title }}
            <el-tag size="small" style="margin-left: 8px">
              {{
                {
                  netease: '网易云',
                  tencent: 'QQ音乐',
                  kugou: '酷狗',
                  xiami: '虾米',
                  baidu: '百度',
                  kuwo: '酷我',
                }[musicInfo.server]
              }}
              ·
              {{
                {
                  search: '搜索',
                  song: '单曲',
                  album: '专辑',
                  artist: '艺术家',
                  playlist: '歌单',
                }[musicInfo.type]
              }}
            </el-tag>
          </div>
          <div class="music-preview-artist">{{ musicInfo.artist }}</div>
        </div>
        <el-button size="small" text @click="removeContent('music')">
          <el-icon>
            <Close />
          </el-icon>
        </el-button>
      </div>
    </div>
  </el-dialog>

  <!-- 视频Dialog -->
  <el-dialog v-model="videoDialogVisible" title="动态视频" width="90%" style="max-width: 400px">
    <div class="video-form">
      <div style="display: flex; gap: 8px">
        <!-- 未添加视频时：显示输入框和解析/上传按钮 -->
        <template v-if="!videoItem">
          <el-input
            v-model="videoUrlInput"
            placeholder="输入视频链接或点击右侧上传"
            @keyup.enter="addVideoUrl"
            style="flex: 1"
          />
          <el-button
            type="primary"
            @click="addVideoUrl"
            :disabled="!videoUrlInput.trim()"
            :loading="fetchingVideo"
          >
            {{ fetchingVideo ? '解析中...' : '解析' }}
          </el-button>
          <el-button type="primary" @click="handleVideoUpload">上传</el-button>
        </template>

        <!-- 已添加视频时：显示只读输入框和删除按钮 -->
        <template v-else>
          <el-input
            :value="
              videoItem.platform && videoItem.video_id
                ? getVideoIframeSrc(videoItem.platform, videoItem.video_id)
                : videoItem.url
            "
            readonly
            style="flex: 1"
          />
          <el-button type="danger" @click="removeVideo">删除</el-button>
        </template>
      </div>
    </div>
  </el-dialog>

  <!-- 标签Dialog -->
  <el-dialog v-model="tagDialogVisible" title="分类标签" width="90%" style="max-width: 400px">
    <div class="tags-form">
      <el-input
        v-model="formData.content.tags"
        placeholder="输入标签名称"
        maxlength="20"
        show-word-limit
      />
      <div style="margin-top: 12px; font-size: 12px; color: #999">
        输入一个标签，如：日常、学习、技术等
      </div>
    </div>
  </el-dialog>

  <!-- 位置Dialog -->
  <el-dialog v-model="locationDialogVisible" title="发布位置" width="90%" style="max-width: 400px">
    <div class="location-form">
      <el-input
        v-model="formData.content.location"
        placeholder="输入位置信息"
        maxlength="100"
        show-word-limit
      />
      <div style="margin-top: 12px; font-size: 12px; color: #999">
        可以输入具体地址、城市或地标名称
      </div>
    </div>
  </el-dialog>

  <!-- 时间Dialog -->
  <el-dialog v-model="timeDialogVisible" title="发布时间" width="90%" style="max-width: 400px">
    <div class="time-form">
      <el-date-picker
        v-model="publishTime"
        type="datetime"
        placeholder="选择发布时间"
        format="YYYY-MM-DD HH:mm"
        value-format="YYYY-MM-DD HH:mm:ss"
        style="width: 100%"
      />
      <div style="margin-top: 12px; font-size: 12px; color: #999">
        设置动态显示的发布时间，不设置则使用当前时间
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue';
import { ElMessage } from 'element-plus';
import {
  Location,
  Link,
  Picture,
  Headset,
  VideoPlay,
  PriceTag,
  Close,
  Timer,
} from '@element-plus/icons-vue';
import type { CreateMomentRequest, UpdateMomentRequest, Moment } from '@/types/moment';
import { createMoment, updateMoment } from '@/api/moment';
import { fetchLinkInfo, parseVideo } from '@/api/tools';
import { uploadFile } from '@/api/file';
import { formatForBackend } from '@/utils/date';
const props = defineProps<{
  modelValue: boolean;
  editMoment?: Moment | null;
}>();

const emit = defineEmits(['update:modelValue', 'success']);

// 基础状态
const visible = computed({
  get: () => props.modelValue,
  set: val => emit('update:modelValue', val),
});
const isEdit = computed(() => !!props.editMoment);
const dialogTitle = computed(() => (isEdit.value ? '编辑动态' : '新增动态'));
const submitLoading = ref(false);
const fetchingLink = ref(false);
const fetchingVideo = ref(false);

// Dialog 显示状态
const linkDialogVisible = ref(false);
const imageDialogVisible = ref(false);
const musicDialogVisible = ref(false);
const videoDialogVisible = ref(false);
const tagDialogVisible = ref(false);
const locationDialogVisible = ref(false);
const timeDialogVisible = ref(false);

// 其他状态
const imageUrlInput = ref('');
const videoUrlInput = ref('');

// 音乐信息
interface MusicInfo {
  title: string;
  artist: string;
  pic: string;
  type: 'search' | 'song' | 'album' | 'artist' | 'playlist'; // 音乐类型
  server: 'netease' | 'tencent' | 'kugou' | 'xiami' | 'baidu' | 'kuwo'; // 平台
}
const musicInfo = ref<MusicInfo | null>(null);
const fetchingMusic = ref(false);

// 图片数据项
interface ImageItem {
  id: string;
  type: 'file' | 'url';
  file?: File;
  url: string;
}
const imageItems = ref<ImageItem[]>([]);

// 视频数据项
interface VideoItem {
  type: 'file' | 'url';
  file?: File;
  url: string;
  platform?: string; // 视频平台
  video_id?: string; // 视频ID
}
const videoItem = ref<VideoItem | null>(null);

// 表单数据
const formData = reactive<CreateMomentRequest>({
  content: {
    text: '',
    tags: '',
    images: [],
    video: { url: '', platform: '', video_id: '' },
    music: { server: 'netease', type: 'song', id: '' },
    link: { url: '', title: '', favicon: '' },
    location: '',
  },
  is_publish: true,
});
const publishTime = ref('');

// 计算属性
const otherContentPreviews = computed(() => {
  const previews = [];
  if (formData.content.link?.url) {
    previews.push({
      type: 'link',
      favicon: formData.content.link.favicon || null,
      title: formData.content.link.title || formData.content.link.url,
      url: formData.content.link.url,
    });
  }
  if (formData.content.music?.id) {
    const MUSIC_LABELS = {
      type: {
        search: '搜索',
        song: '单曲',
        album: '专辑',
        artist: '艺术家',
        playlist: '歌单',
      },
      server: {
        netease: '网易云',
        tencent: 'QQ音乐',
        kugou: '酷狗',
        xiami: '虾米',
        baidu: '百度',
        kuwo: '酷我',
      },
    };

    const music = formData.content.music;
    const info = musicInfo.value;
    const typeLabel = info?.type ? MUSIC_LABELS.type[info.type] : music.type;
    const serverLabel = info?.server ? MUSIC_LABELS.server[info.server] : music.server;

    previews.push({
      type: 'music',
      title: info?.title || '未知歌曲',
      artist: info?.artist || `${music.server} - ${music.type}`,
      cover: info?.pic || '',
      musicType: info ? `${serverLabel} · ${typeLabel}` : `${music.server} - ${music.type}`,
    });
  }
  if (videoItem.value) {
    const { platform, video_id, url, file } = videoItem.value;
    previews.push({
      type: 'video',
      url: platform && video_id ? getVideoIframeSrc(platform, video_id) : url,
      isLocal: !platform,
      text: file?.name || '视频',
    });
  }
  return previews;
});

// 工具方法
const formatTime = (time: string) => {
  if (!time) return '';
  const date = new Date(time);
  const now = new Date();
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
  const targetDay = new Date(date.getFullYear(), date.getMonth(), date.getDate());

  if (targetDay.getTime() === today.getTime()) {
    return `今天 ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
  } else {
    return `${date.getMonth() + 1}-${date.getDate()} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
  }
};

const resetForm = () => {
  // 清理所有 Blob URLs
  imageItems.value.forEach(cleanupBlobUrl);
  imageItems.value = [];
  cleanupVideoBlob();
  videoItem.value = null;

  Object.assign(formData, {
    content: {
      text: '',
      tags: '',
      images: [],
      video: { url: '', platform: '', video_id: '' },
      music: { server: 'netease', type: 'song', id: '' },
      link: { url: '', title: '', favicon: '' },
      location: '',
    },
    is_publish: true,
  });
  publishTime.value = '';
  imageUrlInput.value = '';
  videoUrlInput.value = '';
  musicInfo.value = null;

  // 关闭所有dialog
  linkDialogVisible.value = false;
  imageDialogVisible.value = false;
  musicDialogVisible.value = false;
  videoDialogVisible.value = false;
  tagDialogVisible.value = false;
  locationDialogVisible.value = false;
  timeDialogVisible.value = false;
};

// 链接解析函数
const handleParseLink = async () => {
  if (!formData.content.link?.url?.trim()) {
    ElMessage.warning('请先输入链接地址');
    return;
  }

  fetchingLink.value = true;
  try {
    const result = await fetchLinkInfo({
      url: formData.content.link.url.trim(),
    });
    if (formData.content.link) {
      formData.content.link.title = result.title || '';
      formData.content.link.favicon = result.favicon || '';
    }
    ElMessage.success('解析成功');
  } catch (_error) {
    ElMessage.error('解析失败，请手动填写信息');
  } finally {
    fetchingLink.value = false;
  }
};

// 音乐解析函数
const handleParseMusic = async () => {
  if (
    !formData.content.music?.server ||
    !formData.content.music?.type ||
    !formData.content.music?.id
  ) {
    ElMessage.warning('请先填写音乐平台、类型和ID');
    return;
  }

  fetchingMusic.value = true;
  try {
    const { server, type, id } = formData.content.music;
    const apiUrl = `https://api.injahow.cn/meting/?server=${server}&type=${type}&id=${id}`;

    const response = await fetch(apiUrl);
    const data = await response.json();

    if (data && data.length > 0) {
      const info = data[0];
      musicInfo.value = {
        title: info.name || info.title || '未知歌曲',
        artist: info.artist || info.author || '未知艺术家',
        pic: info.pic || info.cover || '',
        type: type as 'search' | 'song' | 'album' | 'artist' | 'playlist',
        server: server as 'netease' | 'tencent' | 'kugou' | 'xiami' | 'baidu' | 'kuwo',
      };
      ElMessage.success('解析成功');
    } else {
      throw new Error('未获取到音乐信息');
    }
  } catch (_error) {
    ElMessage.error('解析失败，请检查音乐ID是否正确');
    musicInfo.value = null;
  } finally {
    fetchingMusic.value = false;
  }
};

// 本地上传图片
const handleImageUpload = () => {
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = 'image/*';
  input.multiple = true;
  input.onchange = e => {
    const files = (e.target as HTMLInputElement).files;
    if (!files) return;

    Array.from(files).forEach(file => {
      imageItems.value.push({
        id: `${Date.now()}-${Math.random()}`,
        type: 'file',
        file,
        url: URL.createObjectURL(file),
      });
    });
  };
  input.click();
};

// 添加网络图片
const addImageUrl = () => {
  const url = imageUrlInput.value.trim();
  if (!url) return;

  imageItems.value.push({
    id: `${Date.now()}-${Math.random()}`,
    type: 'url',
    url,
  });
  imageUrlInput.value = '';
};

// 视频平台配置
const VIDEO_PLATFORMS = {
  bilibili: {
    name: '哔哩哔哩',
    getIframeSrc: (id: string) => `//player.bilibili.com/player.html?isOutside=true&bvid=${id}&p=1`,
  },
  youtube: {
    name: 'YouTube',
    getIframeSrc: (id: string) => `//www.youtube.com/embed/${id}`,
  },
} as const;

const getPlatformName = (platform: string) =>
  VIDEO_PLATFORMS[platform as keyof typeof VIDEO_PLATFORMS]?.name || platform;
const getVideoIframeSrc = (platform: string, videoId: string) =>
  VIDEO_PLATFORMS[platform as keyof typeof VIDEO_PLATFORMS]?.getIframeSrc(videoId) || '';

// 清理旧的视频Blob URL
const cleanupVideoBlob = () => {
  if (videoItem.value?.type === 'file' && videoItem.value.url.startsWith('blob:')) {
    URL.revokeObjectURL(videoItem.value.url);
  }
};

// 本地上传视频
const handleVideoUpload = () => {
  const input = document.createElement('input');
  input.type = 'file';
  input.accept = 'video/*';
  input.onchange = e => {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;

    cleanupVideoBlob();
    videoItem.value = {
      type: 'file',
      file,
      url: URL.createObjectURL(file),
    };
  };
  input.click();
};

// 添加网络视频
const addVideoUrl = async () => {
  const url = videoUrlInput.value.trim();
  if (!url) return;

  fetchingVideo.value = true;
  try {
    const result = await parseVideo({ url });
    cleanupVideoBlob();

    videoItem.value = {
      type: 'url',
      url,
      platform: result.platform,
      video_id: result.video_id,
    };
    videoUrlInput.value = '';
    ElMessage.success(`已识别：${getPlatformName(result.platform)} - ${result.video_id}`);
  } catch (_error) {
    ElMessage.error('无法识别的视频链接，请检查URL格式（支持B站、YouTube）');
  } finally {
    fetchingVideo.value = false;
  }
};

// 删除视频
const removeVideo = () => {
  cleanupVideoBlob();
  videoItem.value = null;
};

// 清理 Blob URL
const cleanupBlobUrl = (item: ImageItem) => {
  if (item.type === 'file' && item.url.startsWith('blob:')) {
    URL.revokeObjectURL(item.url);
  }
};

// 删除图片（Dialog和预览区域共用）
const removeImage = (index: number) => {
  const item = imageItems.value[index];
  if (!item) return;
  cleanupBlobUrl(item);
  imageItems.value.splice(index, 1);
};

const handleImageError = (e: Event) => {
  const target = e.target as HTMLImageElement;
  const parent = target.parentElement!;
  target.style.display = 'none';
  parent.style.display = 'flex';
  parent.style.alignItems = 'center';
  parent.style.justifyContent = 'center';
  parent.innerHTML = '<span style="color: #c0c4cc; font-size: 12px;">加载失败</span>';
};

// 移除内容
const removeContent = (type: string) => {
  const removeMap = {
    link: () => {
      formData.content.link = { url: '', title: '', favicon: '' };
    },
    images: () => {
      formData.content.images = [];
    },
    music: () => {
      formData.content.music = { server: 'netease', type: 'song', id: '' };
      musicInfo.value = null;
    },
    video: () => {
      removeVideo();
    },
  };
  removeMap[type as keyof typeof removeMap]?.();
};

// 监听器
watch(
  () => props.editMoment,
  moment => {
    resetForm();
    if (moment) {
      Object.assign(formData.content, moment.content);
      formData.is_publish = moment.is_publish;
      publishTime.value = moment.publish_time || '';

      // 加载已有图片（编辑时）
      if (moment.content.images?.length) {
        imageItems.value = moment.content.images.map(url => ({
          id: `${Date.now()}-${Math.random()}`,
          type: 'url' as const,
          url: url,
        }));
      }

      // 加载已有视频（编辑时）
      if (moment.content.video?.url) {
        videoItem.value = {
          type: 'url' as const,
          url: moment.content.video.url,
          platform: moment.content.video.platform,
          video_id: moment.content.video.video_id,
        };
      }
    }
  },
  { immediate: true }
);

// 上传图片（本地文件上传，网络图片直接使用）
const uploadImages = async (): Promise<string[]> => {
  const uploadPromises = imageItems.value.map(async item => {
    if (item.type === 'file' && item.file) {
      const result = await uploadFile(item.file, '动态配图');
      return result.file_url;
    }
    return item.url;
  });

  return Promise.all(uploadPromises);
};

// 上传视频（本地文件上传，网络视频直接使用）
const uploadVideo = async (): Promise<string | null> => {
  if (!videoItem.value) return null;

  if (videoItem.value.type === 'file' && videoItem.value.file) {
    const result = await uploadFile(videoItem.value.file, '动态视频');
    return result.file_url;
  }

  return videoItem.value.url;
};

// 提交表单
const handleCancel = () => {
  visible.value = false;
  resetForm();
};

const handleSubmit = async () => {
  submitLoading.value = true;
  try {
    // 上传图片
    const uploadedImages = imageItems.value.length ? await uploadImages() : [];

    // 上传视频
    const uploadedVideo = await uploadVideo();

    // 清理数据，只传递有值的字段
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const content: any = {};
    if (formData.content.text?.trim()) content.text = formData.content.text.trim();
    if (formData.content.tags?.trim()) content.tags = formData.content.tags.trim();
    if (uploadedImages.length) content.images = uploadedImages;
    if (formData.content.location?.trim()) content.location = formData.content.location.trim();

    // 复杂对象只在有主要字段时才添加
    if (uploadedVideo?.trim()) {
      content.video = { url: uploadedVideo.trim() };
      // 如果是网络视频且已解析出平台和ID，一并发送
      if (videoItem.value?.platform && videoItem.value?.video_id) {
        content.video.platform = videoItem.value.platform;
        content.video.video_id = videoItem.value.video_id;
      }
    }
    if (formData.content.music?.id?.trim()) {
      content.music = {
        server: formData.content.music.server,
        type: formData.content.music.type,
        id: formData.content.music.id.trim(),
      };
    }
    if (formData.content.link?.url?.trim()) {
      content.link = { url: formData.content.link.url.trim() };
      if (formData.content.link.title?.trim())
        content.link.title = formData.content.link.title.trim();
      if (formData.content.link.favicon?.trim())
        content.link.favicon = formData.content.link.favicon.trim();
    }

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const saveData: any = { content, is_publish: formData.is_publish };
    // 如果没有设置发布时间，使用当前时间
    saveData.publish_time = publishTime.value || formatForBackend(new Date());

    if (isEdit.value && props.editMoment) {
      await updateMoment(props.editMoment.id, saveData as UpdateMomentRequest);
      ElMessage.success('更新成功');
    } else {
      await createMoment(saveData as CreateMomentRequest);
      ElMessage.success('创建成功');
    }

    visible.value = false;
    resetForm();
    emit('success');
  } catch (error) {
    if (error instanceof Error) ElMessage.error(error.message);
  } finally {
    submitLoading.value = false;
  }
};
</script>

<style scoped lang="scss">
.moment-editor {
  .toolbar {
    display: flex;
    gap: 12px;
    padding: 16px 0;
    border-bottom: 1px solid #f0f0f0;
    margin-bottom: 16px;
  }

  .content-area {
    margin-bottom: 16px;

    :deep(.el-textarea__inner) {
      border: none;
      box-shadow: none;
      resize: none;
      padding: 0;
      font-size: 16px;
      line-height: 1.6;

      &:focus {
        outline: none;
      }
    }
  }

  // 图片预览区域样式
  .images-preview {
    margin-bottom: 16px;
    display: flex;
    flex-wrap: wrap;
    gap: 8px;

    .image-item {
      position: relative;
      width: 100px;
      height: 100px;
      border-radius: 4px;
      overflow: hidden;
      border: 1px solid #e4e7ed;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }

      .el-button {
        position: absolute;
        top: 4px;
        right: 4px;
        width: 20px;
        height: 20px;
        background-color: rgba(0, 0, 0, 0.7);
        color: white;
      }
    }
  }

  // 其他内容预览区域样式
  .content-preview {
    margin-bottom: 16px;
    padding: 12px;
    background-color: #f8f9fa;
    border-radius: 8px;

    > div {
      display: flex;
      align-items: center;
      gap: 12px;

      .content-icon {
        display: flex;
        align-items: center;
        justify-content: center;

        .favicon-icon {
          width: 50px;
          height: 50px;
          object-fit: contain;
          border-radius: 2px;
        }

        .music-cover {
          width: 50px;
          height: 50px;
          object-fit: cover;
          border-radius: 4px;
        }

        .music-cover-placeholder,
        .video-icon-placeholder {
          width: 50px;
          height: 50px;
          display: flex;
          align-items: center;
          justify-content: center;
          background-color: #e4e7ed;
          border-radius: 4px;
          font-size: 24px;
          color: #909399;
        }
      }

      .content-info {
        flex: 1;

        .preview-title {
          font-weight: 500;
          margin-bottom: 4px;
          display: flex;
          align-items: center;
        }

        .preview-url,
        .preview-artist {
          font-size: 12px;
          color: #999;
        }
      }

      // 视频预览容器
      .video-preview-container {
        flex: 1;
        position: relative;
        width: 100%;
        padding-bottom: 56.25%; // 16:9 宽高比
        border-radius: 8px;
        overflow: hidden;
        background-color: #000;

        .video-iframe {
          position: absolute;
          top: 0;
          left: 0;
          width: 100%;
          height: 100%;
          border: none;
          border-radius: 8px;
        }
      }
    }

    // 视频预览项特殊样式
    .video-preview {
      flex-direction: column;
      align-items: stretch;

      .video-preview-container {
        margin-bottom: 8px;
      }

      // 本地视频预览容器
      .video-local-preview-container {
        margin-bottom: 8px;
        border-radius: 8px;
        overflow: hidden;
        background-color: #000;

        .video-local {
          width: 100%;
          height: auto;
          display: block;
        }
      }

      .el-button {
        align-self: flex-end;
      }
    }
  }

  .bottom-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 16px;
    border-top: 1px solid #f0f0f0;

    .el-button {
      justify-content: center;
      margin-left: 0;
    }

    .el-button.active {
      color: var(--el-color-primary);
    }

    .publish-status {
      display: flex;
      align-items: center;
    }

    // 移动端适配：两行两列布局
    @media screen and (max-width: 600px) {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 12px;

      .publish-status {
        justify-content: center;
      }
    }
  }
}

.el-button.active {
  background-color: var(--el-color-primary);
  color: white;
  border-color: var(--el-color-primary);
}

// Dialog内容样式
.link-form,
.image-form,
.music-form,
.video-form,
.tags-form,
.location-form,
.time-form {
  .favicon-preview {
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid #e4e7ed;
    border-radius: 6px;
    background-color: #fafafa;
    flex-shrink: 0;

    img {
      width: 100%;
      height: auto;
      object-fit: contain;
      border-radius: 2px;
    }

    span {
      font-size: 10px;
      color: #c0c4cc;
      text-align: center;
    }
  }

  .el-button {
    margin-left: 0;
  }
}

.location-form,
.time-form {
  padding: 8px 0;
}

.music-form {
  .music-info-preview {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    background-color: #f8f9fa;
    border-radius: 8px;

    .music-preview-cover {
      width: 60px;
      height: 60px;
      object-fit: cover;
      border-radius: 4px;
    }

    .music-preview-info {
      flex: 1;

      .music-preview-title {
        font-weight: 500;
        margin-bottom: 4px;
        display: flex;
        align-items: center;
      }

      .music-preview-artist {
        font-size: 12px;
        color: #999;
      }
    }
  }
}

.image-form {
  .image-url-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 12px;

    .url-item {
      display: flex;
      gap: 8px;
      align-items: center;
    }
  }
}
</style>
