<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { clearAuth, getCurrentUser, loadMe } from '@/shared/api/http'

const route = useRoute()
const router = useRouter()
const menuOpen = ref(false)
const user = ref(getCurrentUser())

const adminItems = [
  { to: '/admin', label: 'Сводка', icon: '▦', exact: true },
  { to: '/admin/users', label: 'Пользователи', icon: '◉' },
  { to: '/admin/warehouses', label: 'Склады и зоны', icon: '▤' },
]

const initials = computed(() => {
  const name = user.value?.full_name || user.value?.email || 'Администратор'
  return name
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join('') || 'АД'
})

const displayName = computed(() => user.value?.full_name || user.value?.email || 'Пользователь')

function isActive(item) {
  if (item.exact) return route.path === item.to
  return route.path === item.to || route.path.startsWith(`${item.to}/`)
}

async function refreshMe() {
  try {
    user.value = await loadMe()
  } catch {
    user.value = getCurrentUser()
  }
}

function logout() {
  clearAuth()
  router.push({ name: 'login' })
}

onMounted(refreshMe)
</script>

<template>
  <div class="admin-shell">
    <div v-if="menuOpen" class="admin-backdrop" @click="menuOpen = false"></div>
    <button class="burger" type="button" @click="menuOpen = true">☰</button>

    <aside class="admin-sidebar" :class="{ open: menuOpen }">
      <RouterLink class="admin-brand" to="/admin" @click="menuOpen = false">
        <span class="admin-logo">FT</span>
        <span>
          <strong>Fulfillment Transit</strong>
          <em>панель администратора</em>
        </span>
      </RouterLink>

      <div class="nav-group-title">Администрирование</div>
      <nav class="admin-nav" aria-label="Навигация администратора">
        <RouterLink
          v-for="item in adminItems"
          :key="item.to"
          :to="item.to"
          class="admin-nav__item"
          :class="{ active: isActive(item) }"
          @click="menuOpen = false"
        >
          <i>{{ item.icon }}</i>
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>

      <RouterLink class="admin-user-card" to="/admin/profile" aria-label="Открыть профиль" @click="menuOpen = false">
        <div class="admin-avatar">{{ initials }}</div>
        <div>
          <strong>{{ displayName }}</strong>
          <span>администратор системы</span>
        </div>
      </RouterLink>

      <button class="logout-btn" type="button" @click="logout">Выйти</button>
    </aside>

    <main class="admin-main">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.admin-shell {
  min-height: 100vh;
  background:
    radial-gradient(circle at 0 0, rgba(255, 63, 77, .14), transparent 30%),
    linear-gradient(135deg, #edf3f9 0%, #dfe8f1 100%);
  display: grid;
  grid-template-columns: 300px minmax(0, 1fr);
  color: #061126;
}

.admin-sidebar {
  min-height: 100vh;
  background: #061126;
  color: #fff;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 20px;
  position: sticky;
  top: 0;
}

.admin-brand,
.admin-user-card {
  text-decoration: none;
  color: #fff;
}

.admin-brand {
  min-height: 70px;
  border-radius: 24px;
  padding: 12px;
  display: flex;
  align-items: center;
  gap: 12px;
  background: rgba(255, 255, 255, .08);
  border: 1px solid rgba(255, 255, 255, .08);
}

.admin-logo {
  width: 48px;
  height: 48px;
  border-radius: 16px;
  display: grid;
  place-items: center;
  background: #ff3f4d;
  font-weight: 950;
  box-shadow: 0 14px 36px rgba(255, 63, 77, .3);
}

.admin-brand strong,
.admin-user-card strong {
  display: block;
  font-size: 16px;
  font-weight: 950;
}

.admin-brand em,
.admin-user-card span {
  display: block;
  margin-top: 2px;
  color: #a9b8ca;
  font-size: 12px;
  font-style: normal;
  font-weight: 900;
  letter-spacing: .16em;
  text-transform: uppercase;
}

.nav-group-title {
  margin: 8px 8px -8px;
  color: #64748b;
  font-size: 11px;
  font-weight: 950;
  letter-spacing: .22em;
  text-transform: uppercase;
}

.admin-nav {
  display: grid;
  gap: 10px;
}

.admin-nav__item {
  min-height: 58px;
  border-radius: 20px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  color: #cbd5e1;
  text-decoration: none;
  font-weight: 950;
  background: transparent;
}

.admin-nav__item i {
  width: 34px;
  height: 34px;
  border-radius: 12px;
  background: rgba(255, 255, 255, .08);
  display: grid;
  place-items: center;
  font-style: normal;
}

.admin-nav__item.active,
.admin-nav__item:hover {
  background: #243044;
  color: #fff;
}

.admin-nav__item.active i {
  background: #ff3f4d;
  color: #fff;
}

.admin-user-card {
  margin-top: auto;
  border-radius: 24px;
  background: rgba(255, 255, 255, .08);
  padding: 14px;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: background .18s ease, transform .18s ease;
}

.admin-user-card:hover {
  background: rgba(255, 255, 255, .14);
  transform: translateY(-1px);
}

.admin-avatar {
  width: 46px;
  height: 46px;
  border-radius: 16px;
  display: grid;
  place-items: center;
  background: #fff;
  color: #061126;
  font-weight: 950;
}

.logout-btn {
  min-height: 58px;
  border: 0;
  border-radius: 20px;
  background: #243044;
  color: #fff;
  font-size: 16px;
  font-weight: 950;
  cursor: pointer;
}

.logout-btn:hover {
  background: #ff3f4d;
}

.admin-main {
  min-width: 0;
  padding: 30px;
  display: grid;
  gap: 26px;
}

.burger,
.admin-backdrop {
  display: none;
}

@media (max-width: 900px) {
  .admin-shell {
    grid-template-columns: 1fr;
  }

  .burger {
    display: block;
    position: fixed;
    z-index: 40;
    left: 16px;
    top: 16px;
    width: 48px;
    height: 48px;
    border: 0;
    border-radius: 16px;
    background: #ff3f4d;
    color: #fff;
    font-weight: 950;
  }

  .admin-backdrop {
    display: block;
    position: fixed;
    inset: 0;
    z-index: 45;
    background: rgba(6, 17, 38, .55);
  }

  .admin-sidebar {
    position: fixed;
    inset: 0 auto 0 0;
    z-index: 50;
    width: min(320px, calc(100vw - 40px));
    transform: translateX(-110%);
    transition: transform .2s ease;
  }

  .admin-sidebar.open {
    transform: translateX(0);
  }

  .admin-main {
    padding: 80px 16px 24px;
  }
}
</style>
