import js from '@eslint/js';
import vueTsEslintConfig from '@vue/eslint-config-typescript';
import unusedImports from 'eslint-plugin-unused-imports';
import pluginVue from 'eslint-plugin-vue';

export default [
  {
    name: 'app/files-to-lint',
    files: ['**/*.{ts,mts,tsx,vue}'],
  },
  {
    name: 'app/files-to-ignore',
    ignores: [
      '**/dist/**', 
      '**/dist-ssr/**', 
      '**/coverage/**', 
      '**/node_modules/**', 
      '**/scripts/**', 
      'env.d.ts',
      '**/wailsjs/**',  // Auto-generated Wails bindings
      'tailwind.config.js',
    ],
  },
  js.configs.recommended,
  ...pluginVue.configs['flat/essential'],
  ...vueTsEslintConfig(),
  {
    plugins: {
      'unused-imports': unusedImports,
    },
    rules: {
      'vue/multi-word-component-names': 'off',
      'no-unused-vars': 'off',
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
          caughtErrorsIgnorePattern: '.*', // Allow all catch error variables
        },
      ],
      'unused-imports/no-unused-imports': 'error',
      'unused-imports/no-unused-vars': [
        'warn',
        {
          vars: 'all',
          varsIgnorePattern: '^_',
          args: 'after-used',
          argsIgnorePattern: '^_',
          caughtErrorsIgnorePattern: '.*', // Allow all catch error variables
        },
      ],
      // Relax strict rules for practical development
      '@typescript-eslint/no-explicit-any': 'warn', // Downgrade from error to warning
      '@typescript-eslint/ban-ts-comment': 'off', // Allow @ts-ignore
      '@typescript-eslint/no-require-imports': 'off', // Allow require() for dynamic imports
      'vue/require-toggle-inside-transition': 'warn', // Downgrade to warning
      'no-case-declarations': 'off', // Allow declarations in case blocks
      'no-useless-escape': 'warn', // Downgrade to warning
      'prefer-const': 'warn', // Downgrade to warning
    },
  },
  // Override for test files - more relaxed rules
  {
    files: ['**/tests/**/*.{ts,tsx}', '**/*.spec.ts', '**/*.test.ts'],
    rules: {
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unused-vars': 'off',
      'unused-imports/no-unused-vars': 'off',
    },
  },
  // Override for utility files that need any types
  {
    files: ['**/utils/**/*.ts', '**/services/**/*.ts', '**/composables/**/*.ts'],
    rules: {
      '@typescript-eslint/no-explicit-any': 'off', // These files often need any for flexibility
    },
  },
];
