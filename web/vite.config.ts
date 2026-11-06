import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

// The Go app owns the HTML shell, so the build input is the entry module
// instead of an index.html file and the output keeps fixed flat names.
export default defineConfig({
  plugins: [vue()],
  publicDir: false,
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
  build: {
    outDir: '../public/build',
    emptyOutDir: true,
    rollupOptions: {
      input: 'src/main.ts',
      output: {
        entryFileNames: 'index.js',
        chunkFileNames: 'assets/[name]-[hash].js',
        assetFileNames: (assetInfo) => {
          return assetInfo.name?.endsWith('.css')
            ? 'index.css'
            : 'assets/[name]-[hash][extname]'
        },
      },
    },
  },
  test: {
    environment: 'jsdom',
  },
})
