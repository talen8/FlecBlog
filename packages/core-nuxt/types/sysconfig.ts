// 配置分组类型
export type SettingGroupType = 'basic' | 'oauth' | 'upload';

// useSysConfig 返回的配置数据结构
export interface SysConfigData {
  basic: Record<string, string>;
  upload: Record<string, string>;
  oauth: Record<string, string>;
}
