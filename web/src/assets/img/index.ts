import viteUrl from './vite.svg?url'
import vueUrl from './vue.svg?url'

/// Central image registry so components never hard-code asset paths.
export const img = {
  vite: viteUrl,
  vue: vueUrl,
} as const
