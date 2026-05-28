import type { IAppOption } from '../../app';

Component({
  data: {
    showNotice: true,
    noticeText: '',
  },

  lifetimes: {
    attached() {
      const app = getApp<IAppOption>();
      const blogUrl = app.globalData.siteConfig.blog_url || '';
      if (blogUrl) {
        this.setData({ noticeText: `小程序仅用于内容预览，完整功能请访问 ${blogUrl}` });
      }
    },
  },

  methods: {
    closeNotice() {
      this.setData({ showNotice: false });
    },
  },
});
