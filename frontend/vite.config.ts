import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import { defineConfig } from 'vitest/config'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '#wailsjs': resolve(__dirname, 'wailsjs')
    }
  },
  server: {
    port: 3000,
    host: true,
    hmr: {
      overlay: false
    },
    // CRITICAL: File watcher causes memory leak with Wails
    // Only watch src/ directory, ignore wailsjs/ which Wails constantly regenerates
    watch: {
      usePolling: false,
      ignored: [
        '**/node_modules/**',
        '**/.git/**',
        '**/dist/**',
        '**/wailsjs/**',  // Main culprit - Wails regenerates bindings constantly
        '**/backend/**',
        '**/*.go'
      ]
    }
  },
  build: {
    outDir: 'dist',
    sourcemap: false,
    minify: false,  // Disabled to debug circular dependency issue
    rollupOptions: {
      output: {
        // Only split vendor chunks - let Vite handle app code to avoid circular dependency issues
        manualChunks(id) {
          // Vue core + Pinia
          if (id.includes('node_modules/vue') || id.includes('node_modules/pinia')) {
            return 'vue-vendor'
          }
          // UI libraries
          if (id.includes('@headlessui/vue') || id.includes('@floating-ui') || id.includes('fuse.js')) {
            return 'ui-vendor'
          }
          // Utilities
          if (id.includes('@vueuse/core') || id.includes('clsx') || id.includes('tailwind-merge') || id.includes('nanoid')) {
            return 'utils'
          }
          // Icons library
          if (id.includes('lucide-vue-next')) {
            return 'icons'
          }
          // Highlight.js - large, lazy loaded
          if (id.includes('highlight.js')) {
            return 'highlighter'
          }
          // Virtual scroller
          if (id.includes('vue-virtual-scroller')) {
            return 'virtual-scroller'
          }
          // Do NOT manually chunk features - causes "Cannot access before initialization" errors
          // Let Vite/Rollup handle the dependency graph automatically
        },
        chunkFileNames: 'js/[name]-[hash].js',
        entryFileNames: 'js/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash].[ext]'
      }
    },
    chunkSizeWarningLimit: 500
  },
  optimizeDeps: {
    include: ['vue', 'pinia', '@vueuse/core', '@headlessui/vue', 'fuse.js'],
    exclude: ['#wailsjs']
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./tests/setup.ts']
  }
})
