export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig();
  const version = config.public.appVersion;

  console.log(
    `%c FlecBlog %c v${version} %c github.com/talen8/FlecBlog `,
    'background: #49b1f5; color: #fff; padding: 4px 6px; border-radius: 4px 0 0 4px; font-weight: bold;',
    'background: #3a8fd4; color: #fff; padding: 4px 6px;',
    'background: #2d7ab8; color: #fff; padding: 4px 6px; border-radius: 0 4px 4px 0;'
  );
});
