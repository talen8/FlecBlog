/**
 * JSON 解析工具函数
 *
 * 提供安全的 JSON 解析功能，解析失败时返回默认值。
 *
 * @module utils/json
 */

/**
 * 安全解析 JSON 字符串
 *
 * 尝试解析 JSON 字符串，如果解析失败则返回提供的默认值。
 * 适用于解析可能无效的 JSON 数据，如 localStorage 存储的数据。
 *
 * @template T - 返回值类型
 * @param jsonStr - 要解析的 JSON 字符串
 * @param fallback - 解析失败时的默认值
 * @returns 解析结果或默认值
 *
 * @example
 * ```ts
 * // 解析成功
 * const data = parseJSON<{ name: string }>('{"name":"John"}', { name: '' });
 * // { name: 'John' }
 *
 * // 解析失败，返回默认值
 * const data = parseJSON('invalid json', { name: 'default' });
 * // { name: 'default' }
 *
 * // 空值，返回默认值
 * const data = parseJSON(null, []);
 * // []
 * ```
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，提供灵活性
export const parseJSON = <T = any>(jsonStr: string | undefined, fallback: T): T => {
  try {
    return jsonStr ? JSON.parse(jsonStr) : fallback;
  } catch {
    return fallback;
  }
};
