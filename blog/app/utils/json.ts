/**
 * 解析 JSON 字符串
 * @param jsonStr - JSON 字符串
 * @param fallback - 解析失败时的默认值
 * @returns 解析后的数据
 */
// eslint-disable-next-line @typescript-eslint/no-explicit-any -- 泛型默认值，允许调用方不指定类型
export const parseJSON = <T = any>(jsonStr: string | undefined, fallback: T): T => {
  try {
    return jsonStr ? JSON.parse(jsonStr) : fallback;
  } catch {
    return fallback;
  }
};
