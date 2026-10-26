import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import customParseFormat from 'dayjs/plugin/customParseFormat';

// 配置 dayjs
dayjs.locale('zh-cn');
dayjs.extend(customParseFormat);

// 后端统一格式
const BACKEND_DATE_FORMAT = 'YYYY-MM-DD HH:mm:ss';

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
 * 解析后端日期格式为 Date 对象
 * @param dateString 后端返回的日期字符串
 * @returns Date 对象或 null
 */
export function parseBackendDate(dateString: string | null | undefined): Date | null {
  if (!dateString || !dateString.trim()) return null;
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
 * 获取 N 天前的日期
 * @param days 天数
 * @returns 格式化后的日期字符串 YYYY-MM-DD
 */
export function getDaysAgo(days: number): string {
  return dayjs().subtract(days, 'day').format('YYYY-MM-DD');
}

/**
 * 获取 N 月前的日期
 * @param months 月数
 * @returns 格式化后的日期字符串 YYYY-MM-DD
 */
export function getMonthsAgo(months: number): string {
  return dayjs().subtract(months, 'month').format('YYYY-MM-DD');
}

/**
 * 获取今天的日期
 * @returns 格式化后的日期字符串 YYYY-MM-DD
 */
export function getToday(): string {
  return dayjs().format('YYYY-MM-DD');
}

/**
 * 生成日期序列
 * @param startDate 起始日期 YYYY-MM-DD
 * @param endDate 结束日期 YYYY-MM-DD
 * @param unit 单位 'day' | 'month'
 * @param format 格式化模板，默认 YYYY-MM-DD
 * @param minCount 最小数量
 * @returns 日期字符串数组
 */
export function generateDateSeries(
  startDate: string,
  endDate: string,
  unit: 'day' | 'month',
  format: string = 'YYYY-MM-DD',
  minCount: number = 7
): string[] {
  const dates: string[] = [];
  const start = dayjs(startDate);
  const end = dayjs(endDate);

  let current = start;
  const seen = new Set<string>();

  while (current.isBefore(end) || current.isSame(end, 'day') || dates.length < minCount) {
    const dateStr = current.format(format);
    if (!seen.has(dateStr)) {
      dates.push(dateStr);
      seen.add(dateStr);
    }
    current = current.add(1, unit);
  }

  return dates;
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
      if (diffMinutes < 1) {
        return '刚刚';
      }
      return `${diffMinutes}分钟前`;
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
