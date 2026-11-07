import { ref } from 'vue'
import { defineStore } from 'pinia'

/// Boolean flag store used by pages that need a simple on/off state.
export const useToggleStore = defineStore('toggle', () => {
  const active = ref(false)

  const toggle = () => {
    active.value = !active.value
  }

  const setActive = (value: boolean) => {
    active.value = value
  }

  return { active, toggle, setActive }
})
