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
import NotFoundPage from '@/pages/errors/NotFoundPage.vue'
import { getAccessToken } from '@/shared/api/http'

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
    meta: { requiresAuth: true },
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
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFoundPage,
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to) {
    if (to.hash) {
      return { el: to.hash, top: 96, behavior: 'smooth' }
    }
    return { top: 0 }
  },
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !getAccessToken()) {
    return {
      name: 'login',
      query: { redirect: to.fullPath },
    }
  }
  return true
})

export default router
