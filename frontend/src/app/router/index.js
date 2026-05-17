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
import WorkerScanPage from '@/pages/worker/WorkerScanPage.vue'
import WorkerCargoItemsPage from '@/pages/worker/WorkerCargoItemsPage.vue'
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
    path: '/client',
    component: ClientLayout,
    meta: { requiresAuth: true, roles: ['client'] },
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
      { path: 'scan', name: 'worker-scan', component: WorkerScanPage },
      { path: 'cargo-items', name: 'worker-cargo-items', component: WorkerCargoItemsPage },
      { path: 'profile', redirect: '/worker' },
    ],
  },
  {
    path: '/cargo-items/by-qr/:qr',
    name: 'cargo-by-qr',
    component: CargoItemDetailsPage,
    props: true,
    meta: { requiresAuth: true },
  },
  { path: '/:pathMatch(.*)*', name: 'not-found', component: NotFoundPage },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to) {
    if (to.hash) return { el: to.hash, top: 96, behavior: 'smooth' }
    return { top: 0 }
  },
})

function homeForRole(role) {
  if (role === 'worker') return '/worker'
  if (role === 'logist' || role === 'admin') return '/logist'
  return '/client'
}

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !getAccessToken()) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  const allowed = to.meta.roles
  if (allowed?.length) {
    const role = getCurrentUser()?.role
    if (role && !allowed.includes(role)) return homeForRole(role)
  }

  return true
})

export default router
