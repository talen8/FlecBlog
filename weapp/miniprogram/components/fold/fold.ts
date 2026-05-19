/**
 * 折叠面板组件
 * 用于渲染 :::fold 自定义块
 */
Component({
  properties: {
    /** 标题 */
    title: {
      type: String,
      value: '点击展开',
    },
    /** 是否默认展开 */
    open: {
      type: Boolean,
      value: false,
    },
    /** 嵌套内容块数组 */
    contentBlocks: {
      type: Array,
      value: [],
    },
  },

  data: {
    isOpen: false,
  },

  lifetimes: {
    attached() {
      this.setData({ isOpen: this.properties.open });
    },
  },

  methods: {
    /**
     * 切换展开/折叠状态
     */
    toggle() {
      this.setData({ isOpen: !this.data.isOpen });
    },
  },
});
