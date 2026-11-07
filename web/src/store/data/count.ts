import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

/// Counter example store. Uses the setup syntax so state, getters and actions
/// stay grouped the same way across every store in the app.
export const useCountStore = defineStore('count', () => {
  const count = ref(0)
  const double = computed(() => count.value * 2)

  const increment = (step = 1) => {
    count.value += step
  }

  const decrement = (step = 1) => {
    count.value -= step
  }

  const reset = () => {
    count.value = 0
  }

  return { count, double, increment, decrement, reset }
})
