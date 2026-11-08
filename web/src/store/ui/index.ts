import { ref } from 'vue'
import { defineStore } from 'pinia'

export type Theme = 'light' | 'dark'

// Keep in sync with the inline bootstrap script in web/index.html and the Go
// SPA shell: both read the same key before the app mounts.
const THEME_STORAGE_KEY = 'app.theme'

const readStoredTheme = (): Theme | null => {
  try {
    const value = localStorage.getItem(THEME_STORAGE_KEY)
    return value === 'light' || value === 'dark' ? value : null
  } catch {
    return null
  }
}

const systemTheme = (): Theme =>
  typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'

const applyTheme = (value: Theme) => {
  if (typeof document === 'undefined') {
    return
  }

  document.documentElement.classList.toggle('dark', value === 'dark')
  document.documentElement.dataset.theme = value
}

const initialTheme = (): Theme => readStoredTheme() ?? systemTheme()

/// UI-only state: layout toggles and the active theme. Keeping it separate
/// from data stores makes it clear which state is view configuration.
export const useUIStore = defineStore('ui', () => {
  const sidebarOpen = ref(false)
  const theme = ref<Theme>(initialTheme())

  applyTheme(theme.value)

  const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
  }

  const setTheme = (value: Theme) => {
    theme.value = value
    applyTheme(value)

    try {
      localStorage.setItem(THEME_STORAGE_KEY, value)
    } catch {
      // Storage can be blocked; the in-memory theme still applies.
    }
  }

  const toggleTheme = () => {
    setTheme(theme.value === 'dark' ? 'light' : 'dark')
  }

  return { sidebarOpen, theme, toggleSidebar, setTheme, toggleTheme }
})
