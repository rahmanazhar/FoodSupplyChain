import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import InventoryView from '../views/InventoryView.vue'
import ShipmentsView from '../views/ShipmentsView.vue'
import CatalogView from '../views/CatalogView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: DashboardView
    },
    {
      path: '/inventory',
      name: 'inventory',
      component: InventoryView
    },
    {
      path: '/shipments',
      name: 'shipments',
      component: ShipmentsView
    },
    {
      path: '/catalog',
      name: 'catalog',
      component: CatalogView
    }
  ]
})

export default router
