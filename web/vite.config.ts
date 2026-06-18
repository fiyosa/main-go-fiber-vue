import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'

// https://vite.dev/config/
export default defineConfig(({ command }) => ({
  plugins: [vue(), tailwindcss()],

  base: command === 'build' ? '/build/' : '/',
  publicDir: '../public',
  envDir: '..',

  server: {
    port: 3000,
    host: 'localhost',
  },

  build: {
    emptyOutDir: true,
    manifest: false,
    copyPublicDir: false,
    outDir: '../public/build',
    rolldownOptions: {
      input: './src/main.ts',
      output: {
        codeSplitting: true,
        entryFileNames: 'main.js',
        chunkFileNames: 'assets/[hash].js',
        assetFileNames: (assetInfo) => {
          if (assetInfo.names.includes('main.css')) return 'main.css'
          return 'assets/[hash].[ext]'
        },
      },
    },
  },
}))
