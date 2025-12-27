import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    setupFiles: ['./tests/setup.ts'],
    // Unit & integration tests only (exclude e2e which use Playwright)
    include: ['tests/**/*.spec.ts', 'tests/**/*.test.ts'],
    exclude: ['node_modules', 'dist', 'wailsjs', 'tests/e2e/**', '**/*.e2e.spec.ts'],
    testTimeout: 1000,
    hookTimeout: 10000,
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html', 'lcov'],
      exclude: [
        'node_modules/',
        'tests/',
        '**/*.spec.ts',
        '**/*.test.ts',
        'wailsjs/',
        'dist/',
        '*.config.*',
        'src/main.ts',
        'src/App.vue'
      ],
      thresholds: {
        statements: 25,
        branches: 15,
        functions: 15,
        lines: 25
      }
    }
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
      '@tests': resolve(__dirname, './tests'),
      '#wailsjs': resolve(__dirname, './tests/mocks/wailsjs'),
      '../../../../wailsjs': resolve(__dirname, './tests/mocks/wailsjs'),
      '../../../wailsjs': resolve(__dirname, './tests/mocks/wailsjs'),
      '../../wailsjs': resolve(__dirname, './tests/mocks/wailsjs'),
      '../wailsjs': resolve(__dirname, './tests/mocks/wailsjs')
    }
  }
})