import Lenis from 'lenis';

export default defineNuxtPlugin({
  setup() {
    if (!import.meta.client) return;

    const lenis = new Lenis({
      duration: 1.2,
      easing: t => Math.min(1, 1.001 - Math.pow(2, -10 * t)),
    });

    function raf(time: number) {
      lenis.raf(time);
      requestAnimationFrame(raf);
    }

    requestAnimationFrame(raf);

    return { provide: { lenis } };
  },
});
