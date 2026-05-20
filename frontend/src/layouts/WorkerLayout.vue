<script setup>
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router'
import { clearAuth, getCurrentUser } from '@/shared/api/http'

const route = useRoute()
const router = useRouter()
const menuOpen = ref(false)
const user = ref(getCurrentUser())

const navItems = [
  { to: '/worker', label: 'Сводка', icon: '▦', exact: true },
  { to: '/worker/orders', label: 'Заявки', icon: '▥' },
  { to: '/worker/scan', label: 'QR-сканер', icon: '▣' },
  { to: '/worker/cargo-items', label: 'Грузовые места', icon: '▤' },
]

const roleLabels = {
  admin: 'Администратор',
  client: 'Клиент',
  logist: 'Логист',
  logistician: 'Логист',
  worker: 'Рабочий',
  warehouse_worker: 'Рабочий',
}

const displayName = computed(() => user.value?.full_name || user.value?.email || 'Рабочий склада')
const roleLabel = computed(() => roleLabels[String(user.value?.role || '').toLowerCase()] || 'Рабочий')
const initials = computed(() => {
  return (
    displayName.value
      .split(/\s+/)
      .filter(Boolean)
      .map((part) => part[0])
      .join('')
      .slice(0, 2)
      .toUpperCase() || 'РС'
  )
})

function isActive(item) {
  if (item.exact) return route.path === item.to
  return route.path === item.to || route.path.startsWith(`${item.to}/`)
}

function refreshMe() {
  user.value = getCurrentUser()
}

function logout() {
  clearAuth()
  router.push({ name: 'landing' })
}

onMounted(() => {
  refreshMe()
  window.addEventListener('auth:changed', refreshMe)
})

onUnmounted(() => {
  window.removeEventListener('auth:changed', refreshMe)
})
</script>

<template>
  <div class="worker-shell">
    <div v-if="menuOpen" class="worker-backdrop" @click="menuOpen = false"></div>

    <aside class="worker-sidebar" :class="{ open: menuOpen }">
      <RouterLink class="worker-brand" to="/worker" @click="menuOpen = false">
        <span class="worker-logo">FT</span>
        <span>
          <strong>Fulfillment Transit</strong>
          <em>панель рабочего</em>
        </span>
      </RouterLink>

      <nav class="worker-nav" aria-label="Навигация рабочего">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="worker-nav__item"
          :class="{ active: isActive(item) }"
          @click="menuOpen = false"
        >
          <i>{{ item.icon }}</i>
          <span>{{ item.label }}</span>
        </RouterLink>
      </nav>

      <div class="worker-sidebar__bottom">
        <RouterLink class="worker-user-card" to="/worker/profile" aria-label="Открыть профиль" @click="menuOpen = false">
          <span class="worker-user-card__avatar">{{ initials }}</span>
          <span class="worker-user-card__info">
            <strong>{{ displayName }}</strong>
            <em>{{ roleLabel }}</em>
          </span>
        </RouterLink>
        <button type="button" class="worker-logout" @click="logout">Выйти</button>
      </div>
    </aside>

    <main class="worker-main">
      <header class="worker-mobile-head">
        <button type="button" @click="menuOpen = true">☰</button>
        <strong>Fulfillment Transit</strong>
        <span>{{ initials }}</span>
      </header>

      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.worker-shell {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  background: #edf3f9;
  color: #061126;
}

.worker-sidebar {
  position: sticky;
  top: 0;
  height: 100vh;
  padding: 28px 28px 24px;
  background:
    radial-gradient(circle at 18% 0%, rgba(255, 63, 77, .18), transparent 26%),
    linear-gradient(180deg, #071222 0%, #081525 100%);
  color: #fff;
  display: flex;
  flex-direction: column;
  gap: 34px;
  box-shadow: 20px 0 70px rgba(5, 11, 26, .25);
  z-index: 50;
}

.worker-brand {
  color: #fff;
  text-decoration: none;
  display: flex;
  align-items: center;
  gap: 16px;
}

.worker-logo {
  width: 64px;
  height: 64px;
  border-radius: 22px;
  background: #ff3f4d;
  display: grid;
  place-items: center;
  font-size: 22px;
  font-weight: 950;
  letter-spacing: -.03em;
  box-shadow: 0 20px 46px rgba(255, 63, 77, .34);
}

.worker-brand strong {
  display: block;
  font-size: 20px;
  font-weight: 950;
  letter-spacing: -.02em;
}

.worker-brand em {
  display: block;
  margin-top: 6px;
  color: #ff9ca5;
  font-style: normal;
  font-size: 12px;
  font-weight: 950;
  letter-spacing: .22em;
  text-transform: uppercase;
}

.worker-nav {
  display: grid;
  gap: 12px;
}

.worker-nav__item {
  position: relative;
  min-height: 64px;
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 0 16px;
  border-radius: 22px;
  color: rgba(255, 255, 255, .72);
  text-decoration: none;
  font-size: 18px;
  font-weight: 950;
  transition: background .16s ease, color .16s ease, transform .16s ease;
}

.worker-nav__item:hover,
.worker-nav__item.active {
  background: rgba(255, 255, 255, .1);
  color: #fff;
  transform: translateX(4px);
}

.worker-nav__item.active::before {
  content: '';
  position: absolute;
  left: -28px;
  top: 16px;
  bottom: 16px;
  width: 5px;
  border-radius: 999px;
  background: #ff3f4d;
}

.worker-nav__item i {
  width: 32px;
  height: 32px;
  display: grid;
  place-items: center;
  border-radius: 12px;
  background: rgba(255, 255, 255, .08);
  font-style: normal;
}

.worker-sidebar__bottom {
  margin-top: auto;
  display: grid;
  gap: 14px;
}

.worker-user-card {
  border-radius: 24px;
  padding: 16px;
  background: rgba(255, 255, 255, .08);
  display: flex;
  align-items: center;
  gap: 12px;
  color: #fff;
  text-decoration: none;
  transition: background .16s ease, transform .16s ease;
}

.worker-user-card:hover {
  background: rgba(255, 255, 255, .14);
  transform: translateY(-1px);
}

.worker-user-card__avatar {
  width: 46px;
  height: 46px;
  border-radius: 16px;
  background: #fff;
  color: #061126;
  display: grid;
  place-items: center;
  font-weight: 950;
}

.worker-user-card__info {
  min-width: 0;
  display: grid;
  gap: 4px;
}

.worker-user-card__info strong,
.worker-user-card__info em {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.worker-user-card__info em {
  color: #ff9ca5;
  font-style: normal;
  font-size: 12px;
  font-weight: 900;
  letter-spacing: .14em;
  text-transform: uppercase;
}

.worker-logout {
  min-height: 52px;
  border: 0;
  border-radius: 18px;
  background: rgba(255, 255, 255, .1);
  color: #fff;
  font-weight: 950;
  cursor: pointer;
}

.worker-main {
  min-width: 0;
  padding: 36px;
}

.worker-mobile-head {
  display: none;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 18px;
  border-radius: 24px;
  padding: 14px;
  background: #fff;
  box-shadow: 0 14px 44px rgba(15, 23, 42, .08);
}

.worker-mobile-head button,
.worker-mobile-head span {
  width: 44px;
  height: 44px;
  border: 0;
  border-radius: 16px;
  display: grid;
  place-items: center;
  background: #eef3f9;
  color: #061126;
  font-weight: 950;
}

.worker-backdrop {
  display: none;
}

@media (max-width: 980px) {
  .worker-shell {
    grid-template-columns: 1fr;
  }

  .worker-sidebar {
    position: fixed;
    inset: 0 auto 0 0;
    width: min(320px, 88vw);
    transform: translateX(-110%);
    transition: transform .18s ease;
  }

  .worker-sidebar.open {
    transform: translateX(0);
  }

  .worker-backdrop {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(6, 17, 38, .42);
    z-index: 40;
  }

  .worker-main {
    padding: 18px;
  }

  .worker-mobile-head {
    display: flex;
  }
}
</style>
