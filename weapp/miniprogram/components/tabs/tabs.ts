/**
 * 标签页组件
 * 用于渲染 :::tabs 自定义块
 */
Component({
  properties: {
    /** 标签数据 */
    tabs: {
      type: Array,
      value: [],
    },
    /** 默认激活的标签名 */
    activeTab: {
      type: String,
      value: '',
    },
  },

  data: {
    currentTab: '',
  },

  lifetimes: {
    attached() {
      const tabs = this.properties.tabs as Array<{ name: string; contentBlocks: object[] }>;
      const activeTab = this.properties.activeTab;
      if (tabs && tabs.length > 0) {
        const defaultTab = activeTab || (tabs[0] && tabs[0].name) || '';
        this.setData({ currentTab: defaultTab });
      }
    },
  },

  methods: {
    /**
     * 切换标签
     * @param e - 事件对象
     */
    switchTab(e: WechatMiniprogram.TapEvent) {
      const tabName = e.currentTarget.dataset.tab as string;
      this.setData({ currentTab: tabName });
    },
  },

  observers: {
    'tabs, activeTab': function (tabs: Array<{ name: string }>, activeTab: string) {
      if (tabs && tabs.length > 0 && !this.data.currentTab) {
        const defaultTab = activeTab || (tabs[0] && tabs[0].name) || '';
        this.setData({ currentTab: defaultTab });
      }
    },
  },
});
