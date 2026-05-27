/**
 * 格式化工具函数
 * 提供日期、相对时间等常用格式化功能
 */

/**
 * 兼容 iOS 日期格式
 * iOS 只支持 ISO 格式 (yyyy-MM-ddTHH:mm:ss)，将空格替换为 T
 */
function normalizeDate(date: string | Date | number): Date {
  if (typeof date === 'string' && date.includes(' ')) {
    return new Date(date.replace(' ', 'T'));
  }
  return new Date(date);
}

/**
 * 格式化日期
 * @param date - 日期字符串或 Date 对象
 * @param format - 格式化模板，默认 'YYYY-MM-DD'
 * @returns 格式化后的日期字符串
 */
function formatDate(
  date: string | Date | number,
  format = 'YYYY-MM-DD'
): string {
  const d = normalizeDate(date);

  if (isNaN(d.getTime())) {
    return '';
  }

  const year = d.getFullYear();
  const month = d.getMonth() + 1;
  const day = d.getDate();
  const hours = d.getHours();
  const minutes = d.getMinutes();
  const seconds = d.getSeconds();

  const pad = (n: number) => (n < 10 ? `0${n}` : String(n));

  return format
    .replace('YYYY', String(year))
    .replace('MM', pad(month))
    .replace('DD', pad(day))
    .replace('HH', pad(hours))
    .replace('mm', pad(minutes))
    .replace('ss', pad(seconds));
}

/**
 * 格式化相对时间
 * @param date - 日期字符串或 Date 对象
 * @returns 相对时间描述（如：刚刚、5分钟前、1小时前、昨天、3天前）
 */
export function formatRelativeTime(date: string | Date | number): string {
  const d = normalizeDate(date);
  const now = new Date();
  const diff = now.getTime() - d.getTime();

  const minute = 60 * 1000;
  const hour = 60 * minute;
  const day = 24 * hour;

  if (diff < minute) {
    return '刚刚';
  } else if (diff < hour) {
    return `${Math.floor(diff / minute)}分钟前`;
  } else if (diff < day) {
    return `${Math.floor(diff / hour)}小时前`;
  } else if (diff < 2 * day) {
    return '昨天';
  } else if (diff < 7 * day) {
    return `${Math.floor(diff / day)}天前`;
  } else {
    return formatDate(date, 'YYYY-MM-DD');
  }
}
