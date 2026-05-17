<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { getCurrentUser } from '@/shared/api/http'

const route = useRoute()
const open = ref(false)
const user = ref(getCurrentUser())
const root = ref(null)

const panels = [
  { to: '/admin', label: 'Админ-панель', description: 'пользователи, склады, зоны и гейты', icon: 'АД' },
  { to: '/client', label: 'Клиент', description: 'кабинет клиента и заявки', icon: 'КЛ' },
  { to: '/logist', label: 'Логист', description: 'заявки, календарь и отгрузки', icon: 'ЛГ' },
  { to: '/worker', label: 'Рабочий', description: 'склад, QR и грузовые места', icon: 'РБ' },
]

const isAdmin = computed(() => String(user.value?.role || '').toLowerCase() === 'admin')
const isPublicPage = computed(() => ['landing', 'login', 'register'].includes(String(route.name || '')))
const visible = computed(() => isAdmin.value && !isPublicPage.value)
const currentPanel = computed(() => panels.find((panel) => isActive(panel.to)) || panels[0])

function isActive(path) {
  return route.path === path || route.path.startsWith(`${path}/`)
}

function close() {
  open.value = false
}

function refreshFromStorage() {
  user.value = getCurrentUser()
}

function refreshUser() {
  refreshFromStorage()
}

function onDocumentClick(event) {
  if (!root.value?.contains(event.target)) {
    close()
  }
}

function onKeydown(event) {
  if (event.key === 'Escape') {
    close()
  }
}

onMounted(() => {
  refreshUser()
  window.addEventListener('auth:changed', refreshFromStorage)
  window.addEventListener('storage', refreshFromStorage)
  document.addEventListener('click', onDocumentClick)
  document.addEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  window.removeEventListener('auth:changed', refreshFromStorage)
  window.removeEventListener('storage', refreshFromStorage)
  document.removeEventListener('click', onDocumentClick)
  document.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <div v-if="visible" ref="root" class="panel-switcher" :class="{ open }">
    <button class="panel-switcher__button" type="button" @click.stop="open = !open">
      <span class="panel-switcher__mark">{{ currentPanel.icon }}</span>
      <span class="panel-switcher__caption">
        <b>{{ currentPanel.label }}</b>
        <small>переключить панель</small>
      </span>
      <i>⌄</i>
    </button>

    <nav v-if="open" class="panel-switcher__menu" aria-label="Переключение панелей администратора">
      <RouterLink
        v-for="panel in panels"
        :key="panel.to"
        :to="panel.to"
        class="panel-switcher__item"
        :class="{ active: isActive(panel.to) }"
        @click="close"
      >
        <span>{{ panel.icon }}</span>
        <em>
          <b>{{ panel.label }}</b>
          <small>{{ panel.description }}</small>
        </em>
      </RouterLink>
    </nav>
  </div>
</template>

<style scoped>
.panel-switcher {
  position: fixed;
  right: 22px;
  bottom: 22px;
  z-index: 9999;
  color: #061126;
  font-family: inherit;
}

.panel-switcher__button {
  min-height: 62px;
  border: 1px solid rgba(6, 17, 38, .08);
  border-radius: 24px;
  padding: 8px 14px 8px 10px;
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(255, 255, 255, .94);
  color: #061126;
  cursor: pointer;
  box-shadow: 0 22px 70px rgba(6, 17, 38, .22);
  backdrop-filter: blur(18px);
}

.panel-switcher__button i {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  display: grid;
  place-items: center;
  background: #eef3f9;
  color: #64748b;
  font-size: 18px;
  font-style: normal;
  font-weight: 950;
  transition: transform .18s ease;
}

.panel-switcher.open .panel-switcher__button i {
  transform: rotate(180deg);
}

.panel-switcher__mark,
.panel-switcher__item span {
  width: 46px;
  height: 46px;
  border-radius: 16px;
  display: grid;
  place-items: center;
  background: #ff3f4d;
  color: #fff;
  font-size: 13px;
  font-weight: 950;
  box-shadow: 0 12px 32px rgba(255, 63, 77, .24);
}

.panel-switcher__caption {
  min-width: 130px;
  display: grid;
  text-align: left;
}

.panel-switcher__caption b,
.panel-switcher__item b {
  font-size: 15px;
  line-height: 1.1;
  font-weight: 950;
}

.panel-switcher__caption small,
.panel-switcher__item small {
  margin-top: 3px;
  color: #64748b;
  font-size: 11px;
  line-height: 1.1;
  font-weight: 850;
}

.panel-switcher__menu {
  position: absolute;
  right: 0;
  bottom: calc(100% + 12px);
  width: min(340px, calc(100vw - 32px));
  border: 1px solid rgba(6, 17, 38, .08);
  border-radius: 28px;
  padding: 10px;
  display: grid;
  gap: 8px;
  background: rgba(255, 255, 255, .97);
  box-shadow: 0 26px 80px rgba(6, 17, 38, .25);
  backdrop-filter: blur(18px);
}

.panel-switcher__item {
  min-height: 70px;
  border-radius: 22px;
  padding: 10px;
  display: flex;
  align-items: center;
  gap: 12px;
  color: #061126;
  text-decoration: none;
  transition: background .16s ease, transform .16s ease;
}

.panel-switcher__item:hover {
  background: #f3f7fc;
  transform: translateY(-1px);
}

.panel-switcher__item.active {
  background: #061126;
  color: #fff;
}

.panel-switcher__item.active span {
  background: #ff3f4d;
  color: #fff;
}

.panel-switcher__item.active small {
  color: #a9b8ca;
}

.panel-switcher__item em {
  display: grid;
  font-style: normal;
}

@media (max-width: 640px) {
  .panel-switcher {
    right: 12px;
    bottom: 12px;
    left: 12px;
  }

  .panel-switcher__button {
    width: 100%;
    justify-content: flex-start;
  }

  .panel-switcher__caption {
    flex: 1;
  }

  .panel-switcher__button i {
    margin-left: auto;
  }

  .panel-switcher__menu {
    left: 0;
    right: 0;
    width: 100%;
  }
}
</style>
