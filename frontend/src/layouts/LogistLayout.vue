<template>
  <div class="logist-shell" :class="{ 'mobile-open': sidebarOpen }">
    <aside class="logist-sidebar">
      <RouterLink class="brand" to="/logist" @click="sidebarOpen = false">
        <span class="brand-icon">FT</span>
        <span>
          <strong>Fulfillment Transit</strong>
          <small>панель логиста</small>
        </span>
      </RouterLink>

      <nav class="nav-list" aria-label="Навигация логиста">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="nav-link"
          :class="{ active: isActive(item) }"
          @click="sidebarOpen = false"
        >
          <span class="nav-icon" aria-hidden="true">{{ item.icon }}</span>
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>

      <div class="sidebar-footer">
        <RouterLink class="user-card" to="/logist/profile" @click="sidebarOpen = false">
          <span class="avatar">{{ initials }}</span>
          <span>
            <strong>{{ currentUser?.full_name || currentUser?.email || 'Логист' }}</strong>
            <small>{{ roleLabel }}</small>
          </span>
        </RouterLink>
        <button class="logout" type="button" @click="logout">Выйти</button>
      </div>
    </aside>

    <button v-if="sidebarOpen" class="overlay" type="button" aria-label="Закрыть меню" @click="sidebarOpen = false"></button>

    <main class="logist-main">
      <header class="topbar">
        <button class="menu-toggle" type="button" @click="sidebarOpen = !sidebarOpen">☰</button>
        <div>
          <p>Операционное управление</p>
          <h1>{{ pageTitle }}</h1>
        </div>
        <RouterLink class="quick-action" to="/logist/shipments">Создать отгрузку</RouterLink>
      </header>
      <RouterView />
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { clearAuth, getCurrentUser, loadMe } from '@/shared/api/http'

const router = useRouter()
const route = useRoute()
const currentUser = ref(getCurrentUser())
const sidebarOpen = ref(false)

const navItems = [
  { to: '/logist', label: 'Сводка', icon: '▦', exact: true },
  { to: '/logist/orders', label: 'Заявки', icon: '□' },
  { to: '/logist/cargo-items', label: 'Грузы и QR', icon: '⌗' },
  { to: '/logist/pickup-calendar', label: 'Календарь', icon: '◷' },
  { to: '/logist/shipments', label: 'Отгрузки', icon: '↗' },
  { to: '/logist/warehouses', label: 'Склады', icon: '⌂' },
]

const titles = {
  'logist-dashboard': 'Сводка логиста',
  'logist-orders': 'Заявки клиентов',
  'logist-cargo-items': 'Грузы и QR-коды',
  'logist-pickup-calendar': 'Календарь приёмки',
  'logist-shipments': 'Отгрузки',
  'logist-shipment-details': 'Карточка отгрузки',
  'logist-warehouses': 'Склады, зоны и гейты',
  'logist-profile': 'Профиль',
}

const pageTitle = computed(() => titles[route.name] || 'Панель логиста')
const roleLabel = computed(() => currentUser.value?.role === 'admin' ? 'Администратор' : 'Логист')
const initials = computed(() => {
  const name = currentUser.value?.full_name || currentUser.value?.email || 'Л'
  return name.split(' ').filter(Boolean).slice(0, 2).map((part) => part[0]).join('').toUpperCase()
})

function isActive(item) {
  if (item.exact) return route.path === item.to
  return route.path === item.to || route.path.startsWith(`${item.to}/`)
}

async function refreshUser() {
  try {
    currentUser.value = await loadMe()
  } catch {
    currentUser.value = getCurrentUser()
  }
}

function logout() {
  clearAuth()
  router.push({ name: 'login' })
}

watch(sidebarOpen, (isOpen) => {
  document.body.style.overflow = isOpen ? 'hidden' : ''
})

onMounted(() => {
  refreshUser()
  window.addEventListener('auth:changed', refreshUser)
})

onUnmounted(() => {
  window.removeEventListener('auth:changed', refreshUser)
  document.body.style.overflow = ''
})
</script>

<style scoped>
.logist-shell {
  min-height: 100vh;
  background: #07101f;
  color: #07101f;
  display: grid;
  grid-template-columns: 290px minmax(0, 1fr);
}
.logist-sidebar {
  position: sticky;
  top: 0;
  height: 100vh;
  padding: 26px 24px;
  background:
    radial-gradient(circle at 0 0, rgba(255, 63, 77, .14), transparent 34%),
    linear-gradient(180deg, #081222 0%, #101b2d 100%);
  border-right: 1px solid rgba(255,255,255,.08);
  display: flex;
  flex-direction: column;
  gap: 34px;
  overflow-y: auto;
  z-index: 40;
}
.brand, .user-card, .nav-link { text-decoration: none; }
.brand {
  display: flex;
  align-items: center;
  gap: 14px;
  color: white;
}
.brand-icon {
  width: 64px;
  height: 64px;
  border-radius: 22px;
  display: grid;
  place-items: center;
  background: #ff3f4d;
  color: white;
  font-weight: 950;
  font-size: 20px;
  box-shadow: 0 18px 42px rgba(255,63,77,.35);
}
.brand strong { display:block; font-size: 20px; line-height: 1.08; letter-spacing: -.02em; }
.brand small { display:block; margin-top: 6px; color: #ff9da5; letter-spacing: .24em; font-size: 10px; text-transform: uppercase; font-weight: 900; }
.nav-list { display: grid; gap: 14px; }
.nav-link {
  position: relative;
  display: flex;
  align-items: center;
  gap: 13px;
  min-height: 62px;
  color: #b6c1d4;
  padding: 0 20px;
  border-radius: 20px;
  font-weight: 950;
  font-size: 18px;
  outline: none;
  border: 0;
  box-shadow: none;
  transition: background .18s ease, color .18s ease, transform .18s ease;
}
.nav-link::before {
  content: '';
  position: absolute;
  left: 0;
  top: 10px;
  bottom: 10px;
  width: 0;
  border-radius: 0 999px 999px 0;
  background: #ff3f4d;
  transition: width .18s ease;
}
.nav-link:hover {
  color: white;
  background: rgba(255,255,255,.06);
  transform: translateX(2px);
}
.nav-link.active {
  color: white;
  background: rgba(255,255,255,.09);
}
.nav-link.active::before { width: 5px; }
.nav-link.router-link-active,
.nav-link.router-link-exact-active {
  box-shadow: none !important;
  border: 0 !important;
}
.nav-icon {
  width: 22px;
  color: currentColor;
  opacity: .95;
  display: inline-grid;
  place-items: center;
}
.sidebar-footer { margin-top: auto; display: grid; gap: 12px; }
.user-card {
  display: flex;
  gap: 13px;
  align-items: center;
  padding: 15px 16px;
  border-radius: 22px;
  background: rgba(255,255,255,.08);
  color: white;
}
.avatar {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  background: rgba(255,63,77,.18);
  color: #ff6470;
  display:grid;
  place-items:center;
  font-weight: 950;
}
.user-card strong { display:block; font-size: 15px; line-height: 1.1; }
.user-card small { color:#95a5bd; }
.logout, .quick-action, .menu-toggle {
  border: 0;
  cursor: pointer;
  font-weight: 950;
  font-family: inherit;
}
.logout {
  height: 54px;
  border-radius: 18px;
  background: rgba(255,255,255,.08);
  color: white;
}
.logist-main {
  min-width: 0;
  padding: 30px 36px 48px;
  background: #f2f6fb;
}
.topbar {
  display:flex;
  align-items:center;
  justify-content:space-between;
  gap: 18px;
  margin-bottom: 24px;
}
.topbar p { margin:0 0 6px; color:#ff3f4d; letter-spacing:.24em; text-transform: uppercase; font-weight:950; font-size: 12px; }
.topbar h1 { margin:0; font-size: clamp(34px, 4vw, 54px); line-height:.9; letter-spacing:-.06em; font-weight: 500; }
.quick-action {
  display:inline-flex;
  align-items:center;
  justify-content:center;
  min-height: 58px;
  padding: 0 26px;
  border-radius: 20px;
  color: white;
  background:#ff3f4d;
  text-decoration:none;
  box-shadow: 0 18px 42px rgba(255,63,77,.28);
  white-space: nowrap;
}
.menu-toggle { display:none; width:48px; height:48px; border-radius:16px; background:#07101f; color:white; font-size: 20px; }
.overlay { display:none; }
@media (max-width: 980px) {
  .logist-shell { grid-template-columns: 1fr; }
  .logist-sidebar {
    position: fixed;
    left: 0;
    top: 0;
    width: min(320px, 88vw);
    transform: translateX(-105%);
    transition: transform .22s ease;
  }
  .mobile-open .logist-sidebar { transform: translateX(0); }
  .overlay {
    display: block;
    position: fixed;
    inset: 0;
    z-index: 30;
    border: 0;
    background: rgba(7,16,31,.55);
  }
  .logist-main { padding: 18px; }
  .menu-toggle { display:inline-grid; place-items:center; }
  .topbar { align-items:flex-start; }
  .quick-action { min-height:48px; padding:0 16px; }
}
@media (max-width: 620px) {
  .topbar { display:grid; grid-template-columns:auto 1fr; }
  .quick-action { grid-column:1 / -1; width:100%; }
}
</style>
