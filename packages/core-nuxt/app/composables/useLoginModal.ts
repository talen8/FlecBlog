// 全局共享的登录弹窗状态
const showLoginModal = ref(false);

/**
 * 全局登录弹窗状态控制
 * @returns showLoginModal - 弹窗显隐状态（ref），open - 打开弹窗，close - 关闭弹窗
 */
export const useLoginModal = () => {
  const open = () => {
    showLoginModal.value = true;
  };

  const close = () => {
    showLoginModal.value = false;
  };

  return {
    showLoginModal,
    open,
    close,
  };
};
