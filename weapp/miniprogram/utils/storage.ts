/**
 * 本地存储工具函数
 */

/**
 * 存储数据到本地
 */
export function setStorage<T>(key: string, data: T): void {
  try {
    wx.setStorageSync(key, data);
  } catch (error) {
    console.error(`Storage set error [${key}]:`, error);
  }
}

/**
 * 从本地获取数据
 */
export function getStorage<T>(key: string, defaultValue?: T): T | undefined {
  try {
    const value = wx.getStorageSync(key);
    return value !== '' ? (value as T) : defaultValue;
  } catch (error) {
    console.error(`Storage get error [${key}]:`, error);
    return defaultValue;
  }
}

/**
 * 从本地移除数据
 */
export function removeStorage(key: string): void {
  try {
    wx.removeStorageSync(key);
  } catch (error) {
    console.error(`Storage remove error [${key}]:`, error);
  }
}

