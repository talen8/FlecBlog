/**
 * 视频组件
 * 用于渲染 :::video 自定义块
 * 支持B站、YouTube、直接视频链接
 */
Component({
  properties: {
    /** 平台类型: bilibili, youtube, 或直接URL */
    platform: {
      type: String,
      value: '',
    },
    /** 视频ID或URL */
    videoId: {
      type: String,
      value: '',
    },
  },

  data: {
    videoUrl: '',
    videoTitle: '',
    isDirectVideo: false,
  },

  lifetimes: {
    attached() {
      this.initVideo();
    },
  },

  methods: {
    initVideo() {
      const { platform, videoId } = this.properties;

      if (platform === 'bilibili' && videoId) {
        this.setData({
          videoUrl: `https://www.bilibili.com/video/${videoId}`,
          videoTitle: 'B站视频',
          isDirectVideo: false,
        });
      } else if (platform === 'youtube' && videoId) {
        this.setData({
          videoUrl: `https://www.youtube.com/watch?v=${videoId}`,
          videoTitle: 'YouTube视频',
          isDirectVideo: false,
        });
      } else if (platform.startsWith('http://') || platform.startsWith('https://')) {
        this.setData({
          videoUrl: platform,
          videoTitle: '在线视频',
          isDirectVideo: true,
        });
      }
    },

    /**
     * 点击视频链接
     */
    onTapLink() {
      const { videoUrl } = this.data;
      if (!videoUrl) return;

      wx.setClipboardData({
        data: videoUrl,
        success: () => {
          wx.showToast({ title: '链接已复制', icon: 'success' });
        },
      });
    },

    /**
     * 视频错误处理
     */
    onVideoError() {
      wx.showToast({ title: '视频加载失败', icon: 'none' });
    },
  },
});
