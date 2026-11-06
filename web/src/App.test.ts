import { mount } from '@vue/test-utils'
import { describe, expect, it } from 'vitest'

import App from './App.vue'

describe('App', () => {
  it('renders hello world', () => {
    const wrapper = mount(App)

    expect(wrapper.text()).toBe('hello world')
  })
})
