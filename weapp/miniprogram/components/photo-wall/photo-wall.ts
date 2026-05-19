/**
 * 照片墙组件
 * 用于渲染 :::photo 自定义块
 */
Component({
  properties: {
    /** 图片数据，二维数组，每行一个数组 */
    rows: {
      type: Array,
      value: [],
    },
  },

  data: {
    processedRows: [] as string[][],
  },

  lifetimes: {
    attached() {
      this.processRows();
    },
  },

  observers: {
    'rows': function () {
      this.processRows();
    },
  },

  methods: {
    /**
     * 处理图片数据，提取URL
     */
    processRows() {
      const rows = this.properties.rows as string[][];
      if (!rows || rows.length === 0) {
        this.setData({ processedRows: [] });
        return;
      }

      const processedRows = rows.map(row => {
        return row.map(img => this.parseImageUrl(img));
      });

      this.setData({ processedRows });
    },

    /**
     * 解析图片URL
     * @param img - 图片字符串
     * @returns 图片URL
     */
    parseImageUrl(img: string): string {
      const imgMatch = img.match(/^!\[.*?\]\((.*?)\)$/);
      return imgMatch ? (imgMatch[1] || img) : img;
    },

    /**
     * 预览图片
     * @param e - 事件对象
     */
    onPreview(e: WechatMiniprogram.TapEvent) {
      const currentUrl = e.currentTarget.dataset.src as string;
      const allUrls: string[] = [];

      (this.data.processedRows as string[][]).forEach(row => {
        row.forEach(img => {
          allUrls.push(img);
        });
      });

      wx.previewImage({
        current: currentUrl,
        urls: allUrls,
      });
    },
  },
});
