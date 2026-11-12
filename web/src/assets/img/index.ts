import loginDoodleUrl from './login-doodle.svg?url'
import loginSecureUrl from './login-secure.svg?url'
import viteUrl from './vite.svg?url'
import vueUrl from './vue.svg?url'

/// Central image registry so components never hard-code asset paths.
export const img = {
  loginDoodle: loginDoodleUrl,
  loginSecure: loginSecureUrl,
  vite: viteUrl,
  vue: vueUrl,
} as const
