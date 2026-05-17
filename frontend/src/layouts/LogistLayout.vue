<template>
  <div class="logist-shell">
    <aside class="logist-sidebar">
      <RouterLink class="brand" to="/logist">
        <span class="brand-icon">FT</span>
        <span>
          <strong>Fulfillment Transit</strong>
          <small>панель логиста</small>
        </span>
      </RouterLink>

      <nav class="nav-list">
        <RouterLink v-for="item in navItems" :key="item.to" :to="item.to" class="nav-link">
          <span>{{ item.icon }}</span>
          {{ item.label }}
        </RouterLink>
      </nav>

      <div class="sidebar-footer">
        <RouterLink class="user-card" to="/logist/profile">
          <span class="avatar">{{ initials }}</span>
          <span>
            <strong>{{ currentUser?.full_name || currentUser?.email || 'Логист' }}</strong>
            <small>{{ roleLabel }}</small>
          </span>
        </RouterLink>
        <button class="logout" type="button" @click="logout">Выйти</button>
      </div>
    </aside>

    <main class="logist-main">
      <header class="topbar">
        <button class="menu-toggle" type="button" @click="sidebarOpen = !sidebarOpen">☰</button>
        <div>
          <p>Операционное управление</p>
          <h1>{{ pageTitle }}</h1>
        </div>
        <RouterLink class="quick-action" to="/logist/shipments">Создать отгрузку</RouterLink>
      </header>

      <div v-if="sidebarOpen" class="mobile-panel">
        <RouterLink v-for="item in navItems" :key="`m-${item.to}`" :to="item.to" class="nav-link" @click="sidebarOpen = false">
          <span>{{ item.icon }}</span>
          {{ item.label }}
        </RouterLink>
      </div>

      <RouterView />
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { clearAuth, getCurrentUser, loadMe } from '@/shared/api/http'

const router = useRouter()
const route = useRoute()
const currentUser = ref(getCurrentUser())
const sidebarOpen = ref(false)

const navItems = [
  { to: '/logist', label: 'Сводка', icon: '▦' },
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

onMounted(() => {
  refreshUser()
  window.addEventListener('auth:changed', refreshUser)
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
  padding: 24px;
  background: linear-gradient(180deg, #0a1425 0%, #111a2b 100%);
  border-right: 1px solid rgba(255,255,255,.08);
  display: flex;
  flex-direction: column;
  gap: 28px;
}
.brand, .user-card, .nav-link { text-decoration: none; }
.brand {
  display: flex;
  align-items: center;
  gap: 14px;
  color: white;
}
.brand-icon {
  width: 54px;
  height: 54px;
  border-radius: 18px;
  display: grid;
  place-items: center;
  background: #ff3f4d;
  color: white;
  font-weight: 900;
  box-shadow: 0 18px 40px rgba(255,63,77,.35);
}
.brand strong { display:block; font-size: 18px; line-height: 1.1; }
.brand small { display:block; margin-top: 4px; color: #ff9da5; letter-spacing: .24em; font-size: 10px; text-transform: uppercase; }
.nav-list { display: grid; gap: 10px; }
.nav-link {
  display: flex;
  align-items: center;
  gap: 12px;
  color: #b8c3d6;
  padding: 14px 16px;
  border-radius: 18px;
  font-weight: 800;
  transition: .2s ease;
}
.nav-link:hover, .nav-link.router-link-active {
  background: rgba(255,255,255,.08);
  color: white;
}
.nav-link.router-link-active { box-shadow: inset 4px 0 0 #ff3f4d; }
.sidebar-footer { margin-top: auto; display: grid; gap: 12px; }
.user-card {
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 14px;
  border-radius: 22px;
  background: rgba(255,255,255,.06);
  color: white;
}
.avatar {
  width: 42px;
  height: 42px;
  border-radius: 14px;
  background: rgba(255,63,77,.16);
  color: #ff6470;
  display:grid;
  place-items:center;
  font-weight: 900;
}
.user-card strong { display:block; font-size: 14px; }
.user-card small { color:#95a5bd; }
.logout, .quick-action, .menu-toggle {
  border: 0;
  cursor: pointer;
  font-weight: 900;
}
.logout {
  height: 48px;
  border-radius: 16px;
  background: rgba(255,255,255,.08);
  color: white;
}
.logist-main {
  min-width: 0;
  padding: 28px;
  background: #f2f6fb;
}
.topbar {
  display:flex;
  align-items:center;
  justify-content:space-between;
  gap: 18px;
  margin-bottom: 24px;
}
.topbar p { margin:0 0 4px; color:#ff3f4d; letter-spacing:.24em; text-transform: uppercase; font-weight:900; font-size: 12px; }
.topbar h1 { margin:0; font-size: clamp(28px, 4vw, 46px); line-height:.95; letter-spacing:-.05em; }
.quick-action {
  color: white;
  background:#ff3f4d;
  padding: 16px 22px;
  border-radius: 18px;
  text-decoration:none;
  box-shadow: 0 18px 38px rgba(255,63,77,.25);
  white-space: nowrap;
}
.menu-toggle { display:none; width: 48px; height:48px; border-radius:16px; background:#07101f; color:white; }
.mobile-panel { display:none; }
@media (max-width: 980px) {
  .logist-shell { grid-template-columns: 1fr; }
  .logist-sidebar { display:none; }
  .logist-main { padding: 18px; }
  .menu-toggle { display:block; }
  .mobile-panel {
    display:grid;
    gap:8px;
    margin-bottom:18px;
    padding:12px;
    border-radius:22px;
    background:#07101f;
  }
  .quick-action { display:none; }
}
</style>
