<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { apiFetch, getCurrentUser, loadMe } from '@/shared/api/http'
import {
  formatDateTime,
  handoverLabel,
  normalizeCollection,
  roleLabel,
  statusLabel,
} from './clientUtils'

const user = ref(getCurrentUser() || {})
const orders = ref([])
const loading = ref(false)
const errorMessage = ref('')

const userName = computed(() => user.value?.full_name || user.value?.fullName || user.value?.name || 'Пользователь')
const userEmail = computed(() => user.value?.email || '—')
const userPhone = computed(() => user.value?.phone || '—')
const userRole = computed(() => roleLabel(user.value?.role || 'client'))
const userInitials = computed(() => {
  const value = userName.value
    .split(' ')
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0])
    .join('')
  return value || 'FT'
})

const sortedOrders = computed(() => {
  return [...orders.value].sort((a, b) => {
    const left = new Date(a.created_at || a.createdAt || a.pickup_date || 0).getTime()
    const right = new Date(b.created_at || b.createdAt || b.pickup_date || 0).getTime()
    return right - left
  })
})

const recentOrders = computed(() => sortedOrders.value.slice(0, 5))

const stats = computed(() => {
  const done = new Set(['completed', 'closed', 'delivered', 'shipped'])
  const stopped = new Set(['cancelled', 'canceled', 'rejected'])
  const total = orders.value.length
  const completed = orders.value.filter((order) => done.has(String(order.status || '').toLowerCase())).length
  const cancelled = orders.value.filter((order) => stopped.has(String(order.status || '').toLowerCase())).length
  const active = Math.max(total - completed - cancelled, 0)

  return { total, active, completed }
})

function orderId(order) {
  return order?.id || order?.order_id || order?.orderId
}

function orderNumber(order) {
  const id = orderId(order)
  return id ? `#${id}` : '—'
}

function warehouseName(value) {
  if (!value) return ''
  if (typeof value === 'string') return value
  return value.name || value.title || value.label || ''
}

function routeText(order) {
  const receiving =
    warehouseName(order.receiving_warehouse) ||
    warehouseName(order.receivingWarehouse) ||
    order.receiving_warehouse_name ||
    order.receivingWarehouseName ||
    order.from_warehouse_name ||
    order.from ||
    'Склад приёмки'

  const destination =
    warehouseName(order.destination_warehouse) ||
    warehouseName(order.destinationWarehouse) ||
    order.destination_warehouse_name ||
    order.destinationWarehouseName ||
    order.to_warehouse_name ||
    order.to ||
    'Склад назначения'

  return `${receiving} → ${destination}`
}

function createdDate(order) {
  return formatDateTime(order.created_at || order.createdAt || order.pickup_date || order.pickupDate)
}

async function refreshProfile() {
  loading.value = true
  errorMessage.value = ''

  try {
    const me = await loadMe().catch(() => null)
    if (me && typeof me === 'object') {
      user.value = me
    } else {
      user.value = getCurrentUser() || user.value || {}
    }

    const payload = await apiFetch('/orders', { auth: true })
    orders.value = normalizeCollection(payload, 'orders')
  } catch (error) {
    errorMessage.value = error?.message || 'Не удалось загрузить профиль и историю заявок.'
  } finally {
    loading.value = false
  }
}

onMounted(refreshProfile)
</script>

<template>
  <section class="client-profile-page">
    <div class="profile-shell">
      <aside class="profile-illustration-card">
        <div class="brand-pill">
          <span class="brand-mark">FT</span>
          <span>
            <strong>Fulfillment Transit</strong>
            <small>личный профиль</small>
          </span>
        </div>

        <div class="vector-scene" aria-hidden="true">
          <svg viewBox="0 0 560 390" role="img">
            <defs>
              <linearGradient id="sceneGlow" x1="0" y1="0" x2="1" y2="1">
                <stop offset="0" stop-color="#ff4050" stop-opacity="0.86" />
                <stop offset="1" stop-color="#00b3c8" stop-opacity="0.8" />
              </linearGradient>
              <linearGradient id="cardFill" x1="0" y1="0" x2="1" y2="1">
                <stop offset="0" stop-color="#263344" />
                <stop offset="1" stop-color="#111c2d" />
              </linearGradient>
              <filter id="softShadow" x="-20%" y="-20%" width="140%" height="140%">
                <feDropShadow dx="0" dy="22" stdDeviation="18" flood-color="#000" flood-opacity="0.35" />
              </filter>
            </defs>

            <rect x="28" y="30" width="504" height="300" rx="34" fill="#07111f" stroke="#253246" />
            <circle cx="280" cy="186" r="92" fill="url(#sceneGlow)" opacity="0.2" />
            <path d="M84 278H476" stroke="#738199" stroke-width="10" stroke-linecap="round" opacity="0.38" />

            <g filter="url(#softShadow)">
              <rect x="82" y="86" width="190" height="108" rx="24" fill="url(#cardFill)" stroke="#3b495d" />
              <rect x="112" y="116" width="48" height="48" rx="12" fill="#ff4050" />
              <rect x="178" y="122" width="70" height="14" rx="7" fill="#d5dbe6" opacity="0.74" />
              <rect x="178" y="150" width="92" height="12" rx="6" fill="#8d98aa" opacity="0.7" />
              <rect x="112" y="174" width="132" height="10" rx="5" fill="#738199" opacity="0.66" />
            </g>

            <g filter="url(#softShadow)">
              <rect x="312" y="188" width="172" height="88" rx="22" fill="url(#cardFill)" stroke="#3b495d" />
              <rect x="340" y="218" width="88" height="14" rx="7" fill="#d5dbe6" opacity="0.78" />
              <rect x="340" y="246" width="58" height="12" rx="6" fill="#8d98aa" opacity="0.72" />
            </g>

            <g filter="url(#softShadow)">
              <rect x="342" y="72" width="78" height="90" rx="24" fill="#17344a" stroke="#315168" />
              <rect x="366" y="100" width="14" height="14" rx="2" fill="#f7fbff" />
              <rect x="390" y="100" width="14" height="14" rx="2" fill="#f7fbff" />
              <rect x="366" y="130" width="14" height="14" rx="2" fill="#f7fbff" />
              <rect x="390" y="130" width="14" height="14" rx="2" fill="#f7fbff" />
            </g>

            <g filter="url(#softShadow)">
              <circle cx="280" cy="184" r="58" fill="#222b3b" stroke="#566275" />
              <path d="M280 126 330 154v58l-50 30-50-30v-58l50-28Z" fill="none" stroke="#fff" stroke-width="8" stroke-linejoin="round" />
              <path d="M230 154 280 184l50-30M280 184v58" fill="none" stroke="#fff" stroke-width="8" stroke-linejoin="round" opacity="0.9" />
            </g>

            <g filter="url(#softShadow)">
              <rect x="94" y="218" width="84" height="72" rx="20" fill="#331c2f" stroke="#5d3148" />
              <path d="M118 262v-34l20-12 20 12v34h-40Z" fill="none" stroke="#ffdce0" stroke-width="6" stroke-linejoin="round" />
              <path d="M128 250h20" stroke="#ffdce0" stroke-width="5" stroke-linecap="round" />
            </g>
          </svg>
        </div>

        <div class="visual-metrics">
          <div>
            <strong>{{ stats.total }}</strong>
            <span>заявок</span>
          </div>
          <div>
            <strong>{{ stats.active }}</strong>
            <span>активно</span>
          </div>
          <div>
            <strong>{{ stats.completed }}</strong>
            <span>завершено</span>
          </div>
        </div>
      </aside>

      <main class="profile-content-card">
        <div class="profile-heading">
          <div class="avatar">{{ userInitials }}</div>
          <div>
            <p class="eyebrow">Данные пользователя</p>
            <h1>{{ userName }}</h1>
            <p class="verified-line">Проверенный аккаунт Fulfillment Transit</p>
          </div>
          <button class="refresh-button" type="button" :disabled="loading" @click="refreshProfile">
            {{ loading ? 'Обновляем...' : 'Обновить' }}
          </button>
        </div>

        <p v-if="errorMessage" class="error-box">{{ errorMessage }}</p>

        <div class="profile-rows" aria-label="Данные профиля">
          <div class="profile-row">
            <span class="row-label">ФИО</span>
            <strong>{{ userName }}</strong>
          </div>
          <div class="profile-row">
            <span class="row-label">Email</span>
            <strong>{{ userEmail }}</strong>
          </div>
          <div class="profile-row">
            <span class="row-label">Телефон</span>
            <strong>{{ userPhone }}</strong>
          </div>
          <div class="profile-row">
            <span class="row-label">Роль</span>
            <strong>{{ userRole }}</strong>
          </div>
        </div>

        <section class="orders-history">
          <div class="history-head">
            <div>
              <p class="eyebrow">История заявок</p>
              <h2>Последние операции</h2>
            </div>
            <RouterLink class="all-orders-link" :to="{ name: 'client-orders' }">Все заявки</RouterLink>
          </div>

          <div v-if="loading && !orders.length" class="empty-state">Загружаем историю заявок...</div>
          <div v-else-if="!recentOrders.length" class="empty-state">
            Пока нет заявок. Создайте первую заявку, чтобы видеть историю здесь.
          </div>

          <div v-else class="history-list">
            <article v-for="order in recentOrders" :key="orderId(order) || order.created_at" class="history-row">
              <div class="history-main">
                <span class="order-number">Заявка {{ orderNumber(order) }}</span>
                <strong>{{ routeText(order) }}</strong>
                <small>{{ handoverLabel(order.handover_type || order.handoverType) }} · {{ createdDate(order) }}</small>
              </div>
              <span class="status-pill">{{ statusLabel(order.status) }}</span>
              <RouterLink
                v-if="orderId(order)"
                class="open-order-link"
                :to="{ name: 'client-order-details', params: { id: orderId(order) } }"
              >
                Открыть
              </RouterLink>
            </article>
          </div>
        </section>
      </main>
    </div>
  </section>
</template>

<style scoped>
.client-profile-page {
  min-height: calc(100vh - 96px);
  padding: clamp(20px, 3vw, 44px);
  background:
    radial-gradient(circle at 12% 0%, rgba(255, 64, 80, 0.18), transparent 34%),
    radial-gradient(circle at 100% 10%, rgba(0, 179, 200, 0.18), transparent 30%),
    #06101f;
}

.profile-shell {
  width: min(1440px, 100%);
  margin: 0 auto;
  display: grid;
  grid-template-columns: minmax(340px, 0.9fr) minmax(520px, 1.25fr);
  gap: 28px;
  align-items: stretch;
}

.profile-illustration-card,
.profile-content-card {
  border-radius: 34px;
  overflow: hidden;
}

.profile-illustration-card {
  position: relative;
  padding: 34px;
  color: #fff;
  background:
    linear-gradient(145deg, rgba(255, 64, 80, 0.26), transparent 36%),
    linear-gradient(155deg, #0a1020 0%, #081526 55%, #073340 100%);
  border: 1px solid rgba(255, 255, 255, 0.12);
  box-shadow: 0 30px 80px rgba(0, 0, 0, 0.35);
}

.brand-pill {
  width: fit-content;
  display: inline-flex;
  align-items: center;
  gap: 14px;
  padding: 14px 18px;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.13);
  background: rgba(255, 255, 255, 0.08);
}

.brand-mark,
.avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  width: 64px;
  height: 64px;
  border-radius: 20px;
  color: #fff;
  font-weight: 900;
  background: #ff4050;
  box-shadow: 0 18px 38px rgba(255, 64, 80, 0.35);
}

.brand-pill strong {
  display: block;
  font-size: 18px;
  line-height: 1.1;
}

.brand-pill small {
  display: block;
  margin-top: 4px;
  color: #ffadb5;
  font-size: 10px;
  font-weight: 900;
  letter-spacing: 0.42em;
  text-transform: uppercase;
}

.vector-scene {
  margin-top: 38px;
  padding: 12px;
  border-radius: 30px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.12);
}

.vector-scene svg {
  display: block;
  width: 100%;
  height: auto;
}

.visual-metrics {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-top: 18px;
}

.visual-metrics div {
  min-height: 98px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.1);
}

.visual-metrics strong {
  font-size: 30px;
  line-height: 1;
}

.visual-metrics span {
  margin-top: 8px;
  color: #cdd7e7;
  font-size: 12px;
}

.profile-content-card {
  padding: clamp(26px, 4vw, 52px);
  background: #fff;
  color: #07101f;
  box-shadow: 0 30px 80px rgba(0, 0, 0, 0.22);
}

.profile-heading {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr) auto;
  gap: 20px;
  align-items: center;
}

.avatar {
  width: 84px;
  height: 84px;
  border-radius: 26px;
  font-size: 25px;
}

.eyebrow {
  margin: 0 0 8px;
  color: #ff4050;
  font-size: 13px;
  font-weight: 900;
  letter-spacing: 0.42em;
  text-transform: uppercase;
}

.profile-heading h1,
.orders-history h2 {
  margin: 0;
  font-size: clamp(32px, 4vw, 56px);
  line-height: 0.95;
  letter-spacing: -0.05em;
}

.orders-history h2 {
  font-size: clamp(28px, 3vw, 42px);
}

.verified-line {
  margin: 12px 0 0;
  color: #63718a;
  font-weight: 800;
}

.refresh-button,
.all-orders-link,
.open-order-link {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: 0;
  text-decoration: none;
  cursor: pointer;
  font-weight: 900;
  transition: 0.18s ease;
}

.refresh-button {
  min-height: 58px;
  padding: 0 26px;
  border-radius: 20px;
  color: #14233a;
  background: #edf2f8;
}

.refresh-button:disabled {
  opacity: 0.65;
  cursor: progress;
}

.profile-rows {
  margin-top: 34px;
  display: grid;
  border: 1px solid #dce4ef;
  border-radius: 26px;
  overflow: hidden;
  background: #f7f9fc;
}

.profile-row {
  display: grid;
  grid-template-columns: 190px minmax(0, 1fr);
  gap: 24px;
  align-items: center;
  min-height: 74px;
  padding: 20px 24px;
  border-bottom: 1px solid #dce4ef;
}

.profile-row:last-child {
  border-bottom: 0;
}

.row-label {
  color: #94a3ba;
  font-size: 13px;
  font-weight: 900;
  letter-spacing: 0.32em;
  text-transform: uppercase;
}

.profile-row strong {
  min-width: 0;
  color: #07101f;
  font-size: clamp(18px, 1.7vw, 24px);
  line-height: 1.18;
  word-break: break-word;
}

.orders-history {
  margin-top: 30px;
  padding-top: 30px;
  border-top: 1px solid #dce4ef;
}

.history-head {
  display: flex;
  justify-content: space-between;
  gap: 18px;
  align-items: flex-start;
}

.all-orders-link {
  min-height: 54px;
  padding: 0 22px;
  border-radius: 18px;
  color: #fff;
  background: #07101f;
  white-space: nowrap;
}

.history-list {
  display: grid;
  gap: 12px;
  margin-top: 20px;
}

.history-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto auto;
  gap: 14px;
  align-items: center;
  padding: 16px 18px;
  border: 1px solid #dce4ef;
  border-radius: 20px;
  background: #f7f9fc;
}

.history-main {
  min-width: 0;
  display: grid;
  gap: 4px;
}

.order-number {
  color: #ff4050;
  font-size: 12px;
  font-weight: 900;
  letter-spacing: 0.22em;
  text-transform: uppercase;
}

.history-main strong {
  overflow: hidden;
  color: #07101f;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.history-main small {
  color: #687894;
  font-weight: 700;
}

.status-pill {
  padding: 11px 14px;
  border-radius: 999px;
  color: #07101f;
  background: #fff;
  border: 1px solid #dce4ef;
  font-size: 13px;
  font-weight: 900;
  white-space: nowrap;
}

.open-order-link {
  min-height: 44px;
  padding: 0 18px;
  border-radius: 14px;
  color: #fff;
  background: #ff4050;
  white-space: nowrap;
}

.error-box,
.empty-state {
  margin-top: 20px;
  padding: 18px 20px;
  border-radius: 20px;
  font-weight: 800;
}

.error-box {
  color: #b4232b;
  background: #fff1f2;
  border: 1px solid #fecdd3;
}

.empty-state {
  color: #687894;
  background: #f7f9fc;
  border: 1px dashed #dce4ef;
}

@media (max-width: 1180px) {
  .profile-shell {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 760px) {
  .client-profile-page {
    padding: 14px;
  }

  .profile-illustration-card,
  .profile-content-card {
    border-radius: 26px;
  }

  .profile-heading {
    grid-template-columns: 1fr;
  }

  .refresh-button {
    width: 100%;
  }

  .profile-row,
  .history-row {
    grid-template-columns: 1fr;
  }

  .history-head {
    flex-direction: column;
  }

  .all-orders-link,
  .open-order-link {
    width: 100%;
  }

  .visual-metrics {
    grid-template-columns: 1fr;
  }
}
</style>
