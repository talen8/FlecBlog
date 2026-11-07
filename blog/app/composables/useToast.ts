import { ref } from 'vue';

// Toast 类型
export type ToastType = 'success' | 'error' | 'warning' | 'info';

// 内部状态接口
interface ToastItem {
  id: number;
  message: string;
  type: ToastType;
}

// 全局状态
const toasts = ref<ToastItem[]>([]);
let nextId = 1;

/**
 * 按标识移除提示
 * @param id - 提示标识
 */
const removeById = (id: number): void => {
  const index = toasts.value.findIndex(t => t.id === id);
  if (index > -1) {
    toasts.value.splice(index, 1);
  }
};

/**
 * 创建一条提示并安排自动消失
 * @param message - 提示文案
 * @param type - 语义类型，决定配色
 * @param duration - 自定义持续时间，缺省 4.5 秒
 */
const show = (message: string, type: ToastType, duration?: number): void => {
  if (import.meta.server) return;

  const id = nextId++;
  toasts.value.push({ id, message, type });

  while (toasts.value.length > 3) {
    removeById(toasts.value[0]!.id);
  }

  // 自动移除
  setTimeout(() => removeById(id), duration ?? 4500);
};

// 便捷方法
const success = (message: string, duration?: number) => show(message, 'success', duration);
const error = (message: string, duration?: number) => show(message, 'error', duration);
const warning = (message: string, duration?: number) => show(message, 'warning', duration);
const info = (message: string, duration?: number) => show(message, 'info', duration);

/**
 * Toast 消息提示
 */
export function useToast() {
  return { toasts, success, error, warning, info };
}
