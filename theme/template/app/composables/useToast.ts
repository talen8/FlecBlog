/**
 * Toast 消息提示
 *
 * 提供全局的消息提示功能，支持成功、错误、信息、警告四种类型。
 * 消息会自动在指定时间后消失。
 *
 * @module composables/useToast
 */

type ToastType = 'success' | 'error' | 'info' | 'warning';

interface ToastMessage {
  id: number;
  type: ToastType;
  message: string;
}

const toasts = ref<ToastMessage[]>([]);
let toastId = 0;

/**
 * 添加一条 Toast 消息
 *
 * @param type - 消息类型
 * @param message - 消息内容
 * @param duration - 显示时长（毫秒），默认 3000ms，设为 0 则不自动关闭
 */
const addToast = (type: ToastType, message: string, duration = 3000) => {
  const id = ++toastId;
  toasts.value.push({ id, type, message });

  if (duration > 0) {
    setTimeout(() => {
      removeToast(id);
    }, duration);
  }
};

/**
 * 移除一条 Toast 消息
 *
 * @param id - 消息 ID
 */
const removeToast = (id: number) => {
  const index = toasts.value.findIndex(t => t.id === id);
  if (index !== -1) {
    toasts.value.splice(index, 1);
  }
};

/**
 * Toast 消息提示 Composable
 *
 * 提供全局消息提示功能，支持四种消息类型。
 *
 * @returns {Object} Toast 方法和状态
 * @returns {Ref<ToastMessage[]>} toasts - 当前显示的消息列表
 * @returns {Function} success - 显示成功消息
 * @returns {Function} error - 显示错误消息
 * @returns {Function} info - 显示信息消息
 * @returns {Function} warning - 显示警告消息
 * @returns {Function} remove - 移除指定消息
 *
 * @example
 * ```ts
 * const toast = useToast();
 *
 * // 显示成功消息
 * toast.success('操作成功！');
 *
 * // 显示错误消息，持续 5 秒
 * toast.error('操作失败，请重试', 5000);
 *
 * // 显示不自动关闭的消息
 * toast.info('请注意...', 0);
 * ```
 */
export const useToast = () => {
  return {
    toasts,
    success: (message: string, duration?: number) => addToast('success', message, duration),
    error: (message: string, duration?: number) => addToast('error', message, duration),
    info: (message: string, duration?: number) => addToast('info', message, duration),
    warning: (message: string, duration?: number) => addToast('warning', message, duration),
    remove: removeToast,
  };
};
