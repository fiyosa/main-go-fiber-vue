import { ref } from 'vue'
import { defineStore } from 'pinia'

export type Theme = 'light' | 'dark'

/// UI-only state: layout toggles and the active theme. Keeping it separate
/// from data stores makes it clear which state is view configuration.
export const useUIStore = defineStore('ui', () => {
  const sidebarOpen = ref(true)
  const theme = ref<Theme>('light')

  const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
  }

  const setTheme = (value: Theme) => {
    theme.value = value
  }

  return { sidebarOpen, theme, toggleSidebar, setTheme }
})
