import { fileURLToPath } from 'node:url'

import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import { loadEnv } from 'vite'
import { defineConfig } from 'vitest/config'

// The Go app owns the HTML shell, so the build input is the entry module
// instead of an index.html file and the output keeps fixed flat names.
const repoRoot = fileURLToPath(new URL('..', import.meta.url))
const publicDir = fileURLToPath(new URL('../public', import.meta.url))

export default defineConfig(({ command, mode }) => {
  // Read the single .env at the repository root so the frontend and the Go
  // backend share the same configuration.
  const env = loadEnv(mode, repoRoot, '')
  const apiTarget = `http://localhost:${env.APP_PORT || '8080'}`

  return {
    // Production serves the bundle only under /build (see internal/routes/web.go),
    // so every generated asset URL must be prefixed. The dev server keeps "/".
    base: command === 'build' ? '/build/' : '/',
    plugins: [vue(), tailwindcss()],
    envDir: repoRoot,
    publicDir,
    server: {
      port: 5173,
      proxy: {
        '/api': apiTarget,
      },
    },
    build: {
      outDir: '../public/build',
      emptyOutDir: true,
      // The Go server already serves every file in public/ (including
      // vite.svg), so copying publicDir into outDir would only duplicate it.
      copyPublicDir: false,
      rollupOptions: {
        input: 'src/main.ts',
        output: {
          entryFileNames: 'index.js',
          // Everything under build/assets is named by hash only, so the output
          // never leaks source file names. CSS keeps its fixed flat name.
          chunkFileNames: 'assets/[hash].js',
          assetFileNames: (assetInfo) => {
            return assetInfo.name?.endsWith('.css') ? 'index.css' : 'assets/[hash][extname]'
          },
        },
      },
    },
    test: {
      environment: 'jsdom',
    },
  }
})
