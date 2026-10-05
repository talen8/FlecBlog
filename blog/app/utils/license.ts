import type { CopyrightLicense, LicenseKey } from '../../types/theme';

/** 协议库 */
const LICENSE_LIBRARY: Record<LicenseKey, CopyrightLicense> = {
  'cc-by': {
    short: 'CC BY 4.0',
    name: '知识共享 署名 4.0 国际许可协议',
    url: 'https://creativecommons.org/licenses/by/4.0/deed.zh',
    clauses: ['by'],
  },
  'cc-by-sa': {
    short: 'CC BY-SA 4.0',
    name: '知识共享 署名-相同方式共享 4.0 国际许可协议',
    url: 'https://creativecommons.org/licenses/by-sa/4.0/deed.zh',
    clauses: ['by', 'sa'],
  },
  'cc-by-nd': {
    short: 'CC BY-ND 4.0',
    name: '知识共享 署名-禁止演绎 4.0 国际许可协议',
    url: 'https://creativecommons.org/licenses/by-nd/4.0/deed.zh',
    clauses: ['by', 'nd'],
  },
  'cc-by-nc': {
    short: 'CC BY-NC 4.0',
    name: '知识共享 署名-非商业性使用 4.0 国际许可协议',
    url: 'https://creativecommons.org/licenses/by-nc/4.0/deed.zh',
    clauses: ['by', 'nc'],
  },
  'cc-by-nc-sa': {
    short: 'CC BY-NC-SA 4.0',
    name: '知识共享 署名-非商业性使用-相同方式共享 4.0 国际许可协议',
    url: 'https://creativecommons.org/licenses/by-nc-sa/4.0/deed.zh',
    clauses: ['by', 'nc', 'sa'],
  },
  'cc-by-nc-nd': {
    short: 'CC BY-NC-ND 4.0',
    name: '知识共享 署名-非商业性使用-禁止演绎 4.0 国际许可协议',
    url: 'https://creativecommons.org/licenses/by-nc-nd/4.0/deed.zh',
    clauses: ['by', 'nc', 'nd'],
  },
  cc0: {
    short: 'CC0 1.0',
    name: '知识共享 公有领域贡献 1.0 通用',
    url: 'https://creativecommons.org/publicdomain/zero/1.0/deed.zh',
    clauses: ['cc0'],
  },
};

/**
 * 按主题配置取值获取协议完整信息，未配置或取值非法时回退为 CC BY-NC-SA 4.0
 * @param value - 主题配置 copyright_license 的原始值
 * @returns 协议完整信息（简称、名称、条款页地址、条款标识）
 */
export const getLicense = (value: unknown): CopyrightLicense =>
  typeof value === 'string' && value in LICENSE_LIBRARY
    ? LICENSE_LIBRARY[value as LicenseKey]
    : LICENSE_LIBRARY['cc-by-nc-sa'];
