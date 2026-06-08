import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import PaginationBar from '../PaginationBar.vue'

describe('PaginationBar', () => {
  it('renders the result range', () => {
    const wrapper = mount(PaginationBar, { props: { total: 42, limit: 20, offset: 20 } })
    expect(wrapper.text()).toContain('21–40 of 42')
    expect(wrapper.text()).toContain('Page 2 / 3')
  })

  it('disables Prev on the first page', () => {
    const wrapper = mount(PaginationBar, { props: { total: 10, limit: 20, offset: 0 } })
    const prev = wrapper.findAll('button')[0]
    expect(prev.attributes('disabled')).toBeDefined()
  })

  it('emits the next offset', async () => {
    const wrapper = mount(PaginationBar, { props: { total: 42, limit: 20, offset: 0 } })
    await wrapper.findAll('button')[1].trigger('click')
    expect(wrapper.emitted('change')[0]).toEqual([20])
  })
})
