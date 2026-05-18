<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'

import { apiFetch } from '@/shared/api/http'

const users = ref([])
const warehouses = ref([])
const orders = ref([])
const cargoItems = ref([])
const loading = ref(false)
const error = ref('')

function collection(payload, keys) {
  if (Array.isArray(payload)) return payload

  for (const key of keys) {
    if (Array.isArray(payload?.[key])) return payload[key]
  }

  return []
}

const stats = computed(() => {
  const activeUsers = users.value.filter((user) => user.is_active && !user.is_blocked)
  const blockedUsers = users.value.filter((user) => user.is_blocked)
  const activeWarehouses = warehouses.value.filter((warehouse) => warehouse.is_active)
  const openOrders = orders.value.filter(
    (order) => !['cancelled', 'delivered', 'completed'].includes(order.status),
  )
  const problemCargo = cargoItems.value.filter((item) => ['damaged', 'lost', 'cancelled'].includes(item.status))

  return {
    users: users.value.length,
    activeUsers: activeUsers.length,
    blockedUsers: blockedUsers.length,
    admins: users.value.filter((user) => user.role === 'admin').length,
    clients: users.value.filter((user) => user.role === 'client').length,
    logist: users.value.filter((user) => user.role === 'logist').length,
    workers: users.value.filter((user) => user.role === 'worker').length,
    warehouses: warehouses.value.length,
    activeWarehouses: activeWarehouses.length,
    orders: orders.value.length,
    openOrders: openOrders.length,
    cargo: cargoItems.value.length,
    problemCargo: problemCargo.length,
  }
})

const roleCards = computed(() => [
  { label: 'Клиенты', value: stats.value.clients },
  { label: 'Логисты', value: stats.value.logist },
  { label: 'Рабочие', value: stats.value.workers },
  { label: 'Админы', value: stats.value.admins },
])

async function loadData() {
  loading.value = true
  error.value = ''

  try {
    const [usersPayload, warehousesPayload, ordersPayload, cargoPayload] = await Promise.all([
      apiFetch('/users?limit=300', { auth: true }),
      apiFetch('/warehouses', { auth: true }),
      apiFetch('/orders?limit=300', { auth: true }),
      apiFetch('/cargo-items?limit=500', { auth: true }),
    ])

    users.value = collection(usersPayload, ['users', 'items', 'data'])
    warehouses.value = collection(warehousesPayload, ['warehouses', 'items', 'data'])
    orders.value = collection(ordersPayload, ['orders', 'items', 'data'])
    cargoItems.value = collection(cargoPayload, ['cargo_items', 'cargoItems', 'items', 'data'])
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить данные администратора'
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<template>
  <section class="admin-page">
    <header class="hero-card">
      <div>
        <p class="eyebrow">Админ-панель</p>
        <h1>Управление системой</h1>
        <span>
          Здесь собраны административные функции: пользователи, склады, зоны хранения и гейты.
        </span>
      </div>

      <button class="dark-btn" type="button" :disabled="loading" @click="loadData">
        {{ loading ? 'Загрузка…' : 'Обновить' }}
      </button>
    </header>

    <div v-if="error" class="alert error">{{ error }}</div>

    <section class="stats-grid">
      <article>
        <span>Пользователей</span>
        <strong>{{ stats.users }}</strong>
        <small>активных: {{ stats.activeUsers }} · заблокировано: {{ stats.blockedUsers }}</small>
      </article>

      <article>
        <span>Складов</span>
        <strong>{{ stats.warehouses }}</strong>
        <small>активных: {{ stats.activeWarehouses }}</small>
      </article>

      <article>
        <span>Заявок</span>
        <strong>{{ stats.orders }}</strong>
        <small>открытых: {{ stats.openOrders }}</small>
      </article>

      <article>
        <span>Грузовых мест</span>
        <strong>{{ stats.cargo }}</strong>
        <small>проблемных: {{ stats.problemCargo }}</small>
      </article>
    </section>

    <section class="dashboard-layout">
      <article class="panel-card role-panel">
        <p class="eyebrow">Роли</p>
        <h2>Пользователи по ролям</h2>

        <div class="role-list">
          <div v-for="role in roleCards" :key="role.label">
            <span>{{ role.label }}</span>
            <strong>{{ role.value }}</strong>
          </div>
        </div>
      </article>
    </section>

    <section class="action-grid">
      <RouterLink to="/admin/users" class="action-card">
        <span>01</span>
        <strong>Пользователи</strong>
        <small>создание аккаунтов, блокировка, роли и статусы</small>
      </RouterLink>

      <RouterLink to="/admin/warehouses" class="action-card">
        <span>02</span>
        <strong>Склады и структура</strong>
        <small>склады, зоны хранения и гейты отгрузки</small>
      </RouterLink>
    </section>
  </section>
</template>

<style scoped>
.admin-page {
  display: grid;
  gap: 26px;
}

.hero-card,
.panel-card,
.stats-grid article,
.action-card {
  background: #fff;
  border-radius: 34px;
  box-shadow: 0 18px 62px rgba(15, 23, 42, .08);
}

.hero-card {
  padding: 34px;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
}

.eyebrow {
  margin: 0 0 12px;
  color: #ff3f4d;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .28em;
  text-transform: uppercase;
}

h1,
h2 {
  margin: 0;
  color: #061126;
  font-weight: 950;
  letter-spacing: -.06em;
}

h1 {
  font-size: clamp(48px, 7vw, 84px);
  line-height: .9;
}

h2 {
  font-size: clamp(34px, 4vw, 58px);
  line-height: .95;
}

.hero-card span,
.action-card small {
  display: block;
  margin-top: 14px;
  color: #5d6d83;
  font-size: 18px;
  font-weight: 800;
  line-height: 1.55;
}

.dark-btn {
  min-height: 58px;
  border: 0;
  border-radius: 20px;
  padding: 0 24px;
  background: #061126;
  color: #fff;
  font-size: 16px;
  font-weight: 950;
  cursor: pointer;
  white-space: nowrap;
}

.dark-btn:disabled {
  opacity: .65;
  cursor: wait;
}

.alert {
  padding: 18px 22px;
  border-radius: 22px;
  font-weight: 900;
}

.alert.error {
  background: #fff0f1;
  color: #be123c;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 18px;
}

.stats-grid article {
  min-height: 150px;
  padding: 24px;
  display: grid;
  align-content: center;
}

.stats-grid span,
.role-list span {
  color: #97a5bb;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .2em;
  text-transform: uppercase;
}

.stats-grid strong {
  margin-top: 12px;
  color: #061126;
  font-size: 52px;
  line-height: 1;
  font-weight: 950;
  letter-spacing: -.06em;
}

.stats-grid small {
  color: #5d6d83;
  font-weight: 800;
}

.dashboard-layout {
  display: grid;
  grid-template-columns: 1fr;
  gap: 18px;
  align-items: stretch;
}

.panel-card {
  padding: 34px;
}

.role-panel {
  width: 100%;
}

.role-list {
  margin-top: 28px;
  display: grid;
  grid-template-columns: 1fr;
  gap: 14px;
}

.role-list div {
  min-height: 96px;
  border-radius: 24px;
  background: #f6f9fd;
  padding: 22px 26px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
}

.role-list strong {
  color: #061126;
  font-size: clamp(34px, 4vw, 56px);
  line-height: 1;
  font-weight: 950;
  letter-spacing: -.06em;
}

.action-grid {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.action-card {
  padding: 30px;
  min-height: 170px;
  color: #061126;
  text-decoration: none;
  font-weight: 950;
}

.action-card span {
  color: #ff3f4d;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .2em;
}

.action-card strong {
  display: block;
  margin-top: 16px;
  font-size: 30px;
  line-height: 1;
}

@media (max-width: 1180px) {
  .stats-grid,
  .action-grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 760px) {
  .hero-card {
    flex-direction: column;
    align-items: stretch;
  }

  .stats-grid,
  .action-grid {
    grid-template-columns: 1fr;
  }

  .role-list div {
    align-items: flex-start;
    flex-direction: column;
  }

  .dark-btn {
    width: 100%;
  }
}
</style>
