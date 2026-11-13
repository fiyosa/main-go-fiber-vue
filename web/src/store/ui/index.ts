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
  const toast = ref<{ message: string; type: 'error' | 'success' } | null>(null)
  let toastTimer: any = null

  const confirmModal = ref<{
    isOpen: boolean
    title: string
    message: string
    confirmText: string
    cancelText: string
    type: 'danger' | 'warning' | 'info'
    onConfirm: () => void
    onCancel: () => void
  }>({
    isOpen: false,
    title: '',
    message: '',
    confirmText: 'Confirm',
    cancelText: 'Cancel',
    type: 'warning',
    onConfirm: () => {},
    onCancel: () => {},
  })

  applyTheme(theme.value)

  const showToast = (message: string, type: 'error' | 'success' = 'error') => {
    if (toastTimer) clearTimeout(toastTimer)
    toast.value = { message, type }
    toastTimer = setTimeout(() => {
      toast.value = null
    }, 4000)
  }

  const confirm = (options: {
    title: string
    message: string
    confirmText?: string
    cancelText?: string
    type?: 'danger' | 'warning' | 'info'
  }): Promise<boolean> => {
    return new Promise((resolve) => {
      confirmModal.value = {
        isOpen: true,
        title: options.title,
        message: options.message,
        confirmText: options.confirmText ?? 'Confirm',
        cancelText: options.cancelText ?? 'Cancel',
        type: options.type ?? 'warning',
        onConfirm: () => {
          confirmModal.value.isOpen = false
          resolve(true)
        },
        onCancel: () => {
          confirmModal.value.isOpen = false
          resolve(false)
        },
      }
    })
  }

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

  return { sidebarOpen, theme, toast, confirmModal, toggleSidebar, setTheme, toggleTheme, showToast, confirm }
})
