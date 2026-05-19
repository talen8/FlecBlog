/**
 * 音频组件
 * 用于渲染 :::audio 和 :::music 自定义块
 */
Component({
  properties: {
    /** 音频标题 */
    title: {
      type: String,
      value: '',
    },
    /** 音频URL */
    audioUrl: {
      type: String,
      value: '',
    },
    /** 音乐平台 (用于在线音乐) */
    server: {
      type: String,
      value: '',
    },
    /** 音乐ID (用于在线音乐) */
    musicId: {
      type: String,
      value: '',
    },
  },

  data: {
    isPlaying: false,
    currentTime: '0:00',
    duration: '0:00',
    progress: 0,
    displayTitle: '',
    innerAudio: null as WechatMiniprogram.InnerAudioContext | null,
  },

  lifetimes: {
    attached() {
      this.initAudio();
    },

    detached() {
      this.destroyAudio();
    },
  },

  methods: {
    initAudio() {
      const { title, audioUrl, server, musicId } = this.properties;

      if (audioUrl && (audioUrl.startsWith('http://') || audioUrl.startsWith('https://'))) {
        this.setData({ displayTitle: title || '音频' });
      } else if (server && musicId) {
        const serverNames: Record<string, string> = {
          netease: '网易云音乐',
          tencent: 'QQ音乐',
          kugou: '酷狗音乐',
          xiami: '虾米音乐',
          baidu: '百度音乐',
        };
        this.setData({ displayTitle: `${serverNames[server] || server} - 在线音乐` });
      }
    },

    createAudioContext() {
      if (this.data.innerAudio) return this.data.innerAudio;

      const innerAudio = wx.createInnerAudioContext();
      this.setData({ innerAudio });

      innerAudio.onPlay(() => {
        this.setData({ isPlaying: true });
      });

      innerAudio.onPause(() => {
        this.setData({ isPlaying: false });
      });

      innerAudio.onEnded(() => {
        this.setData({ isPlaying: false, progress: 0, currentTime: '0:00' });
      });

      innerAudio.onTimeUpdate(() => {
        const current = innerAudio.currentTime;
        const total = innerAudio.duration;
        this.setData({
          currentTime: this.formatTime(current),
          duration: this.formatTime(total),
          progress: (current / total) * 100,
        });
      });

      innerAudio.onError((err) => {
        console.error('音频播放错误:', err);
        wx.showToast({ title: '播放失败', icon: 'none' });
      });

      return innerAudio;
    },

    formatTime(seconds: number): string {
      if (!seconds || !isFinite(seconds)) return '0:00';
      const mins = Math.floor(seconds / 60);
      const secs = Math.floor(seconds % 60);
      return `${mins}:${secs.toString().padStart(2, '0')}`;
    },

    /**
     * 播放/暂停
     */
    togglePlay() {
      const { audioUrl, server, musicId } = this.properties;

      if (server && musicId) {
        wx.showToast({ title: '在线音乐暂不支持', icon: 'none' });
        return;
      }

      if (!audioUrl) {
        wx.showToast({ title: '音频地址无效', icon: 'none' });
        return;
      }

      const innerAudio = this.createAudioContext();

      if (!innerAudio.src) {
        innerAudio.src = audioUrl;
      }

      if (this.data.isPlaying) {
        innerAudio.pause();
      } else {
        innerAudio.play();
      }
    },

    /**
     * 销毁音频实例
     */
    destroyAudio() {
      if (this.data.innerAudio) {
        this.data.innerAudio.destroy();
        this.setData({ innerAudio: null });
      }
    },
  },
});
