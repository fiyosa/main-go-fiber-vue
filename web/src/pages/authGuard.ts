import type { RouteLocationNormalized, RouteLocationNormalizedLoaded, NavigationGuardReturn } from 'vue-router'
import { getMe } from '../api/auth'

export const authGuard = async (
  to: RouteLocationNormalized,
  _from: RouteLocationNormalizedLoaded
): Promise<NavigationGuardReturn> => {
  const isGuestOnly = to.path === '/login' || to.path === '/register'
  const requiresAuth = to.matched.some((record) => record.meta?.requiresAuth)

  if (!isGuestOnly && !requiresAuth) {
    return true
  }

  let authenticated = false
  try {
    await getMe()
    authenticated = true
  } catch {
    authenticated = false
  }

  if (isGuestOnly && authenticated) {
    return { name: 'admin' }
  }

  if (requiresAuth && !authenticated) {
    return { name: 'login' }
  }

  return true
}
