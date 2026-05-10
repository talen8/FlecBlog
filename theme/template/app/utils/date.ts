/**
 * 日期格式化工具函数
 *
 * 基于 dayjs 提供日期格式化功能，支持相对时间显示。
 *
 * @module utils/date
 */

import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import relativeTime from 'dayjs/plugin/relativeTime';

dayjs.locale('zh-cn');
dayjs.extend(relativeTime);

/**
 * 格式化日期为标准格式
 *
 * 将日期格式化为 YYYY-MM-DD 格式。
 *
 * @param date - 日期字符串或 Date 对象
 * @returns 格式化后的日期字符串，无效日期返回 '-'
 *
 * @example
 * ```ts
 * formatDate('2024-01-15T10:30:00Z'); // '2024-01-15'
 * formatDate(new Date()); // 当前日期
 * formatDate(null); // '-'
 * ```
 */
export function formatDate(date: string | Date | null | undefined): string {
  if (!date) return '-';
  return dayjs(date).format('YYYY-MM-DD');
}

/**
 * 格式化动态时间
 *
 * 根据时间距离现在的长短，智能选择显示格式：
 * - 1 分钟内：显示"刚刚"
 * - 1 小时内：显示"X分钟前"
 * - 24 小时内：显示"X小时前"
 * - 3 天内：显示"X天前"
 * - 今年内：显示"M月D日"
 * - 往年：显示"YYYY年M月D日"
 *
 * @param date - 日期字符串或 Date 对象
 * @returns 格式化后的时间字符串，无效日期返回 '-'
 *
 * @example
 * ```ts
 * formatMomentTime(new Date()); // '刚刚'
 * formatMomentTime('2024-01-15T10:30:00Z'); // 根据时间差显示不同格式
 * ```
 */
export function formatMomentTime(date: string | Date | null | undefined): string {
  if (!date) return '-';

  const now = dayjs();
  const target = dayjs(date);
  const diffHours = now.diff(target, 'hour');
  const diffDays = now.diff(target, 'day');

  if (diffHours < 24) {
    if (diffHours < 1) {
      const diffMinutes = now.diff(target, 'minute');
      return diffMinutes < 1 ? '刚刚' : `${diffMinutes}分钟前`;
    }
    return `${diffHours}小时前`;
  }

  if (diffDays < 3) {
    return `${diffDays}天前`;
  }

  if (now.year() === target.year()) {
    return target.format('M月D日');
  }

  return target.format('YYYY年M月D日');
}
