import { ref } from 'vue';

// Toast 类型
type ToastType = 'success' | 'error' | 'warning' | 'info';

// 内部状态接口
interface ToastItem {
  id: number;
  message: string;
  type: ToastType;
  show: boolean;
}

// 全局状态
const toasts = ref<ToastItem[]>([]);
let nextId = 1;

// 默认持续时间配置
const DEFAULT_DURATION: Record<ToastType, number> = {
  success: 3000,
  error: 4000,
  warning: 3500,
  info: 3000,
};

/**
 * Toast 消息提示
 */
export function useToast() {
  const show = (message: string, type: ToastType = 'info', duration?: number) => {
    const id = nextId++;
    const toast: ToastItem = { id, message, type, show: true };

    toasts.value.push(toast);

    // 自动移除
    const delay = duration ?? DEFAULT_DURATION[type];
    setTimeout(() => {
      const index = toasts.value.findIndex(t => t.id === id);
      if (index > -1 && toasts.value[index]) {
        toasts.value[index]!.show = false;
        setTimeout(() => toasts.value.splice(index, 1), 300);
      }
    }, delay);

    return id;
  };

  // 便捷方法
  const success = (message: string, duration?: number) => show(message, 'success', duration);
  const error = (message: string, duration?: number) => show(message, 'error', duration);
  const warning = (message: string, duration?: number) => show(message, 'warning', duration);
  const info = (message: string, duration?: number) => show(message, 'info', duration);

  return { toasts, success, error, warning, info };
}
