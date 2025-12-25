import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    setupFiles: ['./tests/setup.ts'],
    include: ['tests/**/*.spec.ts', 'tests/**/*.test.ts'],
    exclude: ['node_modules', 'dist', 'wailsjs'],
    testTimeout: 1000,
    hookTimeout: 10000,
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'tests/',
        '**/*.spec.ts',
        '**/*.test.ts',
        'wailsjs/',
        'dist/'
      ],
      thresholds: {
        // Realistic thresholds based on current coverage
        // Current: lines 23.75%, functions 7.54%, branches 5.65%, statements 22.38%
        // TODO: Gradually increase as test coverage improves
        lines: 20,
        functions: 5,
        branches: 5,
        statements: 20
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