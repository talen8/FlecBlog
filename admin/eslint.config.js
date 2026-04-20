import js from '@eslint/js'
import ts from 'typescript-eslint'
import pluginVue from 'eslint-plugin-vue'
import eslintPluginPrettier from 'eslint-plugin-prettier/recommended'

export default [
  {
    ignores: [
      'dist/**',
      'dev-dist/**',
      'node_modules/**',
      '*.config.js',
      '*.config.ts',
      '**/sw.js',
      '**/workbox-*.js',
      '**/cfw_*.js',
    ],
  },
  js.configs.recommended,
  ...ts.configs.recommended,
  ...pluginVue.configs['flat/essential'],
  {
    files: ['**/*.vue'],
    languageOptions: {
      parserOptions: {
        parser: ts.parser,
      },
    },
  },
  {
    rules: {
      'vue/multi-word-component-names': 'off',
      '@typescript-eslint/no-explicit-any': 'warn',
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
          caughtErrorsIgnorePattern: '^_',
        },
      ],
      'no-empty': ['error', { allowEmptyCatch: true }],
      '@typescript-eslint/ban-ts-comment': 'warn',
    },
  },
  eslintPluginPrettier,
]
