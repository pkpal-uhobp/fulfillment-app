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

const roleLabel = computed(() => {
  const role = String(user.value?.role || '').toLowerCase()
  return roleLabels[role] || 'Рабочий'
})

const initials = computed(() => {
  const source = displayName.value || user.value?.email || 'Рабочий склада'

  return source
    .split(/\s+/)
    .filter(Boolean)
    .map((part) => part[0])
    .join('')
    .slice(0, 2)
    .toUpperCase() || 'РС'
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
        <div class="worker-user-card" aria-label="Текущий пользователь">
          <span class="worker-user-card__avatar">{{ initials }}</span>
          <span class="worker-user-card__info">
            <strong>{{ displayName }}</strong>
            <em>{{ roleLabel }}</em>
          </span>
        </div>

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
  gap: 15px;
  padding: 0 18px;
  border-radius: 0 22px 22px 0;
  color: #aeb9c8;
  text-decoration: none;
  font-size: 18px;
  font-weight: 950;
  transition: background .18s ease, color .18s ease, transform .18s ease;
  outline: none;
}

.worker-nav__item::before {
  content: '';
  position: absolute;
  left: -28px;
  top: 12px;
  bottom: 12px;
  width: 5px;
  border-radius: 999px;
  background: transparent;
}

.worker-nav__item:hover {
  color: #fff;
  background: rgba(255, 255, 255, .06);
  transform: translateX(2px);
}

.worker-nav__item.active {
  color: #fff;
  background: #202b3d;
  box-shadow: none;
}

.worker-nav__item.active::before {
  background: #ff3f4d;
  box-shadow: 0 0 0 5px rgba(255, 63, 77, .12);
}

.worker-nav__item i {
  width: 24px;
  text-align: center;
  font-style: normal;
  color: inherit;
}

.worker-sidebar__bottom {
  margin-top: auto;
  display: grid;
  gap: 16px;
}

.worker-user-card {
  width: 100%;
  min-height: 92px;
  padding: 20px;
  border-radius: 26px;
  background: #202b3d;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, .03);
}

.worker-user-card__avatar {
  width: 56px;
  height: 56px;
  flex: 0 0 auto;
  border-radius: 16px;
  background: rgba(255, 63, 77, .25);
  color: #ff6370;
  display: grid;
  place-items: center;
  font-size: 20px;
  font-weight: 950;
  letter-spacing: -.03em;
}

.worker-user-card__info {
  min-width: 0;
  display: grid;
  gap: 6px;
}

.worker-user-card__info strong {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #fff;
  font-size: 18px;
  font-weight: 950;
  line-height: 1.1;
}

.worker-user-card__info em {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #b7c1cf;
  font-style: normal;
  font-size: 15px;
  font-weight: 700;
}

.worker-logout {
  min-height: 68px;
  border: 0;
  border-radius: 24px;
  background: #202b3d;
  color: #fff;
  font-size: 18px;
  font-weight: 950;
  cursor: pointer;
  transition: background .18s ease, transform .18s ease;
}

.worker-logout:hover {
  background: #ff3f4d;
  transform: translateY(-1px);
}

.worker-main {
  min-width: 0;
  padding: 34px;
}

.worker-mobile-head,
.worker-backdrop {
  display: none;
}

.worker-mobile-head span {
  width: 46px;
  height: 46px;
  border-radius: 15px;
  background: #ff3f4d;
  display: grid;
  place-items: center;
  font-weight: 950;
}

@media (max-width: 980px) {
  .worker-shell {
    display: block;
  }

  .worker-sidebar {
    position: fixed;
    inset: 0 auto 0 0;
    width: min(320px, 88vw);
    transform: translateX(-105%);
    transition: transform .22s ease;
  }

  .worker-sidebar.open {
    transform: translateX(0);
  }

  .worker-backdrop {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(4, 10, 24, .58);
    z-index: 40;
  }

  .worker-main {
    padding: 18px;
  }

  .worker-mobile-head {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    margin-bottom: 18px;
    padding: 12px;
    border-radius: 22px;
    background: #071222;
    color: #fff;
  }

  .worker-mobile-head button {
    width: 46px;
    height: 46px;
    border: 0;
    border-radius: 15px;
    background: #ff3f4d;
    color: #fff;
    display: grid;
    place-items: center;
    font-weight: 950;
  }
}
</style>
