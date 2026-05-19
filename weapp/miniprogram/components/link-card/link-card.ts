/**
 * 链接卡片组件
 * 用于渲染 :::link 自定义块
 */
Component({
  properties: {
    /** 标题 */
    title: {
      type: String,
      value: '',
    },
    /** 链接地址 */
    link: {
      type: String,
      value: '',
    },
    /** 描述 */
    description: {
      type: String,
      value: '',
    },
  },

  data: {
    isExternal: false,
    linkType: '',
  },

  lifetimes: {
    attached() {
      const link = this.properties.link;
      const isExternal = link.startsWith('http://') || link.startsWith('https://');
      this.setData({
        isExternal,
        linkType: isExternal ? '引用站外链接' : '站内链接',
      });
    },
  },

  methods: {
    /**
     * 点击链接
     */
    onTap() {
      const { link, isExternal } = this.data;
      if (!link) return;

      if (isExternal) {
        wx.setClipboardData({
          data: link,
          success: () => {
            wx.showToast({ title: '链接已复制', icon: 'success' });
          },
        });
      } else {
        wx.navigateTo({
          url: link,
          fail: () => {
            wx.setClipboardData({
              data: link,
              success: () => {
                wx.showToast({ title: '链接已复制', icon: 'success' });
              },
            });
          },
        });
      }
    },
  },
});
