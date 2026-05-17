import { createRouter, createWebHistory } from 'vue-router'

import PublicLayout from '@/layouts/PublicLayout.vue'
import LandingPage from '@/pages/public/LandingPage.vue'
import LoginPage from '@/pages/auth/LoginPage.vue'
import RegisterPage from '@/pages/auth/RegisterPage.vue'

import ClientLayout from '@/layouts/ClientLayout.vue'
import ClientDashboardPage from '@/pages/client/ClientDashboardPage.vue'
import ClientOrdersPage from '@/pages/client/ClientOrdersPage.vue'
import ClientCreateOrderPage from '@/pages/client/ClientCreateOrderPage.vue'
import ClientOrderDetailsPage from '@/pages/client/ClientOrderDetailsPage.vue'
import ClientCargoItemsPage from '@/pages/client/ClientCargoItemsPage.vue'
import ClientProfilePage from '@/pages/client/ClientProfilePage.vue'

import LogistLayout from '@/layouts/LogistLayout.vue'
import LogistDashboardPage from '@/pages/logist/LogistDashboardPage.vue'
import LogistOrdersPage from '@/pages/logist/LogistOrdersPage.vue'
import LogistCargoItemsPage from '@/pages/logist/LogistCargoItemsPage.vue'
import LogistPickupCalendarPage from '@/pages/logist/LogistPickupCalendarPage.vue'
import LogistShipmentsPage from '@/pages/logist/LogistShipmentsPage.vue'
import LogistShipmentDetailsPage from '@/pages/logist/LogistShipmentDetailsPage.vue'
import LogistWarehousesPage from '@/pages/logist/LogistWarehousesPage.vue'
import LogistProfilePage from '@/pages/logist/LogistProfilePage.vue'

import WorkerLayout from '@/layouts/WorkerLayout.vue'
import WorkerDashboardPage from '@/pages/worker/WorkerDashboardPage.vue'
import WorkerOrdersPage from '@/pages/worker/WorkerOrdersPage.vue'
import WorkerScanPage from '@/pages/worker/WorkerScanPage.vue'
import WorkerCargoItemsPage from '@/pages/worker/WorkerCargoItemsPage.vue'

import AdminLayout from '@/layouts/AdminLayout.vue'
import AdminDashboardPage from '@/pages/admin/AdminDashboardPage.vue'
import AdminUsersPage from '@/pages/admin/AdminUsersPage.vue'
import AdminWarehousesPage from '@/pages/admin/AdminWarehousesPage.vue'

import CargoItemDetailsPage from '@/pages/cargo/CargoItemDetailsPage.vue'
import NotFoundPage from '@/pages/errors/NotFoundPage.vue'
import { getAccessToken, getCurrentUser } from '@/shared/api/http'

const routes = [
  {
    path: '/',
    component: PublicLayout,
    children: [
      { path: '', name: 'landing', component: LandingPage },
      { path: 'login', name: 'login', component: LoginPage },
      { path: 'register', name: 'register', component: RegisterPage },
    ],
  },
  {
    path: '/admin',
    component: AdminLayout,
    meta: { requiresAuth: true, roles: ['admin'] },
    children: [
      { path: '', name: 'admin-dashboard', component: AdminDashboardPage },
      { path: 'users', name: 'admin-users', component: AdminUsersPage },
      { path: 'warehouses', name: 'admin-warehouses', component: AdminWarehousesPage },
    ],
  },
  {
    path: '/client',
    component: ClientLayout,
    meta: { requiresAuth: true, roles: ['client', 'admin'] },
    children: [
      { path: '', name: 'client-dashboard', component: ClientDashboardPage },
      { path: 'orders', name: 'client-orders', component: ClientOrdersPage },
      { path: 'orders/new', name: 'client-order-create', component: ClientCreateOrderPage },
      { path: 'orders/:id', name: 'client-order-details', component: ClientOrderDetailsPage, props: true },
      { path: 'cargo-items', name: 'client-cargo-items', component: ClientCargoItemsPage },
      { path: 'profile', name: 'client-profile', component: ClientProfilePage },
    ],
  },
  {
    path: '/logist',
    component: LogistLayout,
    meta: { requiresAuth: true, roles: ['logist', 'admin'] },
    children: [
      { path: '', name: 'logist-dashboard', component: LogistDashboardPage },
      { path: 'orders', name: 'logist-orders', component: LogistOrdersPage },
      { path: 'cargo-items', name: 'logist-cargo-items', component: LogistCargoItemsPage },
      { path: 'pickup-calendar', name: 'logist-pickup-calendar', component: LogistPickupCalendarPage },
      { path: 'shipments', name: 'logist-shipments', component: LogistShipmentsPage },
      { path: 'shipments/:id', name: 'logist-shipment-details', component: LogistShipmentDetailsPage, props: true },
      { path: 'warehouses', name: 'logist-warehouses', component: LogistWarehousesPage },
      { path: 'profile', name: 'logist-profile', component: LogistProfilePage },
    ],
  },
  {
    path: '/worker',
    component: WorkerLayout,
    meta: { requiresAuth: true, roles: ['worker', 'admin'] },
    children: [
      { path: '', name: 'worker-dashboard', component: WorkerDashboardPage },
      { path: 'orders', name: 'worker-orders', component: WorkerOrdersPage },
      { path: 'scan', name: 'worker-scan', component: WorkerScanPage },
      { path: 'cargo-items', name: 'worker-cargo-items', component: WorkerCargoItemsPage },
      { path: 'profile', name: 'worker-profile', redirect: '/worker' },
    ],
  },
  {
    path: '/cargo-items/by-qr/:qr',
    name: 'cargo-by-qr',
    component: CargoItemDetailsPage,
    props: true,
    meta: { requiresAuth: true },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFoundPage,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    }

    if (to.hash) {
      return {
        el: to.hash,
        top: 96,
        behavior: 'smooth',
      }
    }

    if (to.path === from.path) {
      return false
    }

    return { top: 0 }
  },
})

function homeForRole(role) {
  if (role === 'admin') return '/admin'
  if (role === 'worker' || role === 'warehouse_worker') return '/worker'
  if (role === 'logist' || role === 'logistician') return '/logist'

  return '/client'
}

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !getAccessToken()) {
    return { name: 'landing' }
  }

  const allowed = to.meta.roles

  if (allowed?.length) {
    const role = String(getCurrentUser()?.role || '').toLowerCase()

    if (role && !allowed.includes(role)) {
      return homeForRole(role)
    }
  }

  return true
})

export default router
