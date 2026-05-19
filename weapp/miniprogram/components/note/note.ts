/**
 * 提示框组件
 * 用于渲染 :::note 自定义块
 */
Component({
  properties: {
    /** 提示类型: info, tip, warning, danger, success */
    type: {
      type: String,
      value: 'info',
    },
    /** 标题 */
    title: {
      type: String,
      value: '',
    },
    /** 嵌套内容块数组 */
    contentBlocks: {
      type: Array,
      value: [],
    },
  },

  data: {
    iconMap: {
      info: '💡',
      tip: '💡',
      warning: '⚠️',
      warn: '⚠️',
      danger: '🚫',
      error: '🚫',
      success: '✅',
    },
  },

  methods: {},
});
