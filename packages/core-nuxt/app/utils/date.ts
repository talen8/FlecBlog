import dayjs from 'dayjs/esm';
import 'dayjs/esm/locale/zh-cn';
import relativeTime from 'dayjs/esm/plugin/relativeTime';
import customParseFormat from 'dayjs/esm/plugin/customParseFormat';

// 配置 dayjs
dayjs.locale('zh-cn');
dayjs.extend(relativeTime);
dayjs.extend(customParseFormat);

// 后端统一格式
export const BACKEND_DATE_FORMAT = 'YYYY-MM-DD HH:mm:ss';

/**
 * 格式化日期为完整的日期时间
 * @param date 日期字符串或 Date 对象
 * @returns 格式化后的字符串，如 "2025-10-03 13:46:59"
 */
export function formatDateTime(date: string | Date | null | undefined): string {
  if (!date) return '-';
  return dayjs(date).format(BACKEND_DATE_FORMAT);
}

/**
 * 格式化日期为日期（不含时间）
 * @param date 日期字符串或 Date 对象
 * @returns 格式化后的字符串，如 "2025-10-03"
 */
export function formatDate(date: string | Date | null | undefined): string {
  if (!date) return '-';
  return dayjs(date).format('YYYY-MM-DD');
}

/**
 * 格式化为相对时间（如：2小时前）
 * @param date 日期字符串或 Date 对象
 * @returns 相对时间字符串
 */
export function formatRelativeTime(date: string | Date | null | undefined): string {
  if (!date) return '-';
  return dayjs(date).fromNow();
}

/**
 * 解析后端日期格式为 Date 对象
 * @param dateString 后端返回的日期字符串
 * @returns Date 对象或 null
 */
export function parseBackendDate(dateString: string | null | undefined): Date | null {
  if (!dateString?.trim()) return null;
  const parsed = dayjs(dateString, BACKEND_DATE_FORMAT, true);
  return parsed.isValid() ? parsed.toDate() : null;
}

/**
 * 将 Date 对象格式化为后端需要的格式
 * @param date Date 对象
 * @returns 后端格式的日期字符串
 */
export function formatForBackend(date: Date | null | undefined): string {
  if (!date) return '';
  return dayjs(date).format(BACKEND_DATE_FORMAT);
}

/**
 * 判断日期是否有效
 * @param date 日期字符串或 Date 对象
 * @returns 是否有效
 */
export function isValidDate(date: string | Date | null | undefined): boolean {
  return date ? dayjs(date).isValid() : false;
}

/**
 * 格式化为友好的显示格式
 * @param date 日期字符串或 Date 对象
 * @returns 友好的日期字符串，如 "2025年10月3日"
 */
export function formatFriendly(date: string | Date | null | undefined): string {
  if (!date) return '-';

  const target = dayjs(date);

  // 始终显示完整的年月日
  return target.format('YYYY年M月D日');
}

/**
 * 格式化为动态友好时间
 * @param date 日期字符串或 Date 对象
 * @returns 友好时间字符串：n小时前（24小时内）、n天前（3天内）、几月几日（本年）、几年几月几日（非本年）
 */
export function formatMomentTime(date: string | Date | null | undefined): string {
  if (!date) return '-';

  const now = dayjs();
  const target = dayjs(date);
  const diffHours = now.diff(target, 'hour');
  const diffDays = now.diff(target, 'day');

  // 24小时内显示小时
  if (diffHours < 24) {
    if (diffHours < 1) {
      const diffMinutes = now.diff(target, 'minute');
      return diffMinutes < 1 ? '刚刚' : `${diffMinutes}分钟前`;
    }
    return `${diffHours}小时前`;
  }

  // 3天内显示天数
  if (diffDays < 3) {
    return `${diffDays}天前`;
  }

  // 今年的日期显示月日
  if (now.year() === target.year()) {
    return target.format('M月D日');
  }

  // 其他年份显示年月日
  return target.format('YYYY年M月D日');
}
