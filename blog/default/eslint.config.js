import { createConfigForNuxt } from '@nuxt/eslint-config/flat';
import eslintPluginPrettier from 'eslint-plugin-prettier/recommended';

export default createConfigForNuxt({
  // Nuxt 项目配置
}).append(eslintPluginPrettier);
