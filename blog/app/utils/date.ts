import dayjs from 'dayjs/esm';
import 'dayjs/esm/locale/zh-cn';

// 配置 dayjs
dayjs.locale('zh-cn');

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
 * 计算自建站日期起运行的天数
 * @param established 建站日期（YYYY-MM-DD）
 * @returns 运行天数
 */
export function getRunningDays(established: string | null | undefined): number | null {
  const start = established ? new Date(established).getTime() : NaN;
  if (!Number.isFinite(start)) return null;
  return Math.max(0, Math.floor((Date.now() - start) / 86400000));
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
 * 格式化音频时长为 mm:ss
 * @param seconds 秒数
 * @returns 格式化后的时长字符串，非法值返回 "00:00"
 */
export function formatDuration(seconds: number | null | undefined): string {
  if (typeof seconds !== 'number' || !Number.isFinite(seconds) || Number.isNaN(seconds)) {
    return '00:00';
  }
  const total = Math.max(0, Math.floor(seconds));
  const mins = Math.floor(total / 60);
  const secs = total % 60;
  return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
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
