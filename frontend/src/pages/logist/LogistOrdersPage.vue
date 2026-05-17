<template>
  <section class="page">
    <div v-if="error" class="alert">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <div class="toolbar">
      <div>
        <p>Заявки</p>
        <h2>Управление статусами</h2>
      </div>
      <div class="filters">
        <select v-model="filters.status" @change="loadOrders">
          <option v-for="option in orderStatusOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
        </select>
        <select v-model="filters.handoverType" @change="loadOrders">
          <option value="">Все способы передачи</option>
          <option value="self_delivery">Сдача на склад</option>
          <option value="pickup">Забор с адреса</option>
        </select>
        <button type="button" @click="loadOrders">Обновить</button>
      </div>
    </div>

    <div class="orders-grid">
      <article v-for="order in orders" :key="order.id" class="order-card">
        <div class="card-top">
          <div>
            <span>Заявка #{{ order.id }}</span>
            <h3>{{ labelFromMap(orderStatusLabels, order.status) }}</h3>
          </div>
          <em>{{ handoverLabels[order.handover_type] || order.handover_type }}</em>
        </div>

        <div class="meta-grid">
          <div><small>Приёмка</small><b>{{ warehouseName(order.receiving_warehouse_id, warehouses) }}</b></div>
          <div><small>Назначение</small><b>{{ warehouseName(order.destination_warehouse_id, warehouses) }}</b></div>
          <div><small>Дата</small><b>{{ order.self_delivery_date || order.pickup?.pickup_date || formatDate(order.created_at) }}</b></div>
          <div><small>Мест</small><b>{{ totalPlaces(order) }}</b></div>
        </div>

        <div class="action-box">
          <label>
            Новый статус
            <select v-model="forms[order.id].status">
              <option v-for="option in orderStatusOptions.filter((item) => item.value)" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </label>
          <label>
            Комментарий
            <input v-model.trim="forms[order.id].comment" placeholder="Например: принято логистом" />
          </label>
          <button type="button" :disabled="loadingId === order.id" @click="updateStatus(order)">
            {{ loadingId === order.id ? 'Сохраняем...' : 'Обновить статус' }}
          </button>
        </div>

        <details class="history" @toggle="onHistoryToggle($event, order.id)">
          <summary>История заявки</summary>
          <div v-if="historyLoadingId === order.id" class="muted">Загружаем историю...</div>
          <div v-else-if="history[order.id]?.length" class="history-list">
            <div v-for="item in history[order.id]" :key="item.id" class="history-item">
              <b>{{ labelFromMap(orderStatusLabels, item.new_status) }}</b>
              <span>{{ formatDateTime(item.changed_at) }}</span>
              <small>{{ translateComment(item.comment) }}</small>
            </div>
          </div>
          <div v-else class="muted">Истории пока нет</div>
        </details>
      </article>
    </div>

    <div v-if="!orders.length && !error" class="empty">Заявки не найдены</div>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { apiFetch } from '@/shared/api/http'
import {
  formatDate,
  formatDateTime,
  handoverLabels,
  labelFromMap,
  orderStatusLabels,
  orderStatusOptions,
  unwrapList,
  warehouseName,
} from './logistUtils'

const orders = ref([])
const warehouses = ref([])
const forms = reactive({})
const history = reactive({})
const filters = reactive({ status: '', handoverType: '' })
const error = ref('')
const success = ref('')
const loadingId = ref(null)
const historyLoadingId = ref(null)

function ensureForm(order) {
  if (!forms[order.id]) forms[order.id] = { status: order.status || 'created', comment: '' }
}

function totalPlaces(order) {
  return (order.cargo_places || []).reduce((sum, item) => sum + Number(item.quantity || 0), 0)
}

function translateComment(comment) {
  if (!comment) return 'Без комментария'
  return String(comment)
    .replaceAll('created', 'создана')
    .replaceAll('received', 'принята')
    .replaceAll('stored', 'на хранении')
    .replaceAll('shipped', 'отгружена')
}

async function loadWarehouses() {
  const payload = await apiFetch('/warehouses')
  warehouses.value = unwrapList(payload, 'warehouses')
}

async function loadOrders() {
  error.value = ''
  success.value = ''
  const params = new URLSearchParams({ limit: '100' })
  if (filters.status) params.set('status', filters.status)
  if (filters.handoverType) params.set('handover_type', filters.handoverType)
  try {
    const payload = await apiFetch(`/orders?${params}`, { auth: true })
    orders.value = unwrapList(payload, 'orders')
    orders.value.forEach(ensureForm)
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить заявки'
  }
}

async function updateStatus(order) {
  error.value = ''
  success.value = ''
  loadingId.value = order.id
  try {
    const form = forms[order.id]
    const payload = await apiFetch(`/orders/${order.id}/status`, {
      auth: true,
      method: 'PATCH',
      body: {
        status: form.status,
        comment: form.comment || undefined,
      },
    })
    const updated = payload.order || payload.data || payload
    const index = orders.value.findIndex((item) => item.id === order.id)
    if (index !== -1) orders.value[index] = updated
    ensureForm(updated)
    forms[order.id].comment = ''
    success.value = `Статус заявки #${order.id} обновлён`
    delete history[order.id]
  } catch (err) {
    error.value = err.message || 'Не удалось обновить статус заявки'
  } finally {
    loadingId.value = null
  }
}

async function onHistoryToggle(event, orderId) {
  if (!event.target.open || history[orderId]) return
  historyLoadingId.value = orderId
  try {
    const payload = await apiFetch(`/orders/${orderId}/history`, { auth: true })
    history[orderId] = unwrapList(payload, 'history')
  } catch (err) {
    history[orderId] = []
  } finally {
    historyLoadingId.value = null
  }
}

onMounted(async () => {
  await Promise.all([loadWarehouses(), loadOrders()])
})
</script>

<style scoped>
.page { display:grid; gap:20px; }
.alert, .success { padding:16px 18px; border-radius:18px; font-weight:900; }
.alert { background:#fee2e2; color:#991b1b; }
.success { background:#d1fae5; color:#065f46; }
.toolbar { display:flex; justify-content:space-between; gap:18px; align-items:end; padding:26px; border-radius:32px; background:white; box-shadow:0 18px 42px rgba(7,16,31,.08); }
.toolbar p { margin:0 0 8px; color:#ff3f4d; text-transform:uppercase; letter-spacing:.22em; font-weight:900; font-size:12px; }
.toolbar h2 { margin:0; font-size:34px; letter-spacing:-.04em; }
.filters { display:flex; gap:12px; flex-wrap:wrap; justify-content:flex-end; }
select, input { height:48px; border:1px solid #dbe3ef; border-radius:16px; padding:0 14px; background:#f6f8fb; font-weight:800; color:#07101f; }
button { height:48px; border:0; border-radius:16px; padding:0 18px; font-weight:900; cursor:pointer; background:#07101f; color:white; }
button:disabled { opacity:.55; cursor:wait; }
.orders-grid { display:grid; grid-template-columns: repeat(2, minmax(0,1fr)); gap:20px; }
.order-card { display:flex; flex-direction:column; gap:18px; padding:24px; border-radius:32px; background:white; box-shadow:0 18px 42px rgba(7,16,31,.08); }
.card-top { display:flex; justify-content:space-between; gap:14px; align-items:start; }
.card-top span { color:#ff3f4d; letter-spacing:.22em; text-transform:uppercase; font-size:12px; font-weight:900; }
.card-top h3 { margin:8px 0 0; font-size:28px; letter-spacing:-.04em; }
.card-top em { font-style:normal; padding:10px 12px; border-radius:999px; background:#ffe6e8; color:#ff3f4d; font-weight:900; white-space:nowrap; }
.meta-grid { display:grid; grid-template-columns: repeat(2, minmax(0,1fr)); gap:12px; }
.meta-grid div { padding:16px; border-radius:20px; background:#f6f8fb; display:grid; gap:5px; }
small, .muted { color:#6b7b91; font-weight:800; }
.meta-grid b { word-break: break-word; }
.action-box { margin-top:auto; display:grid; grid-template-columns: 1fr 1fr; gap:12px; padding:16px; border-radius:24px; background:#f6f8fb; }
.action-box label { display:grid; gap:8px; font-size:12px; text-transform:uppercase; letter-spacing:.14em; color:#8b9ab0; font-weight:900; }
.action-box button { grid-column:1 / -1; background:#ff3f4d; }
.history { border-top:1px solid #edf1f7; padding-top:14px; }
summary { cursor:pointer; font-weight:900; color:#07101f; }
.history-list { display:grid; gap:10px; margin-top:12px; }
.history-item { padding:14px; border-radius:18px; background:#f6f8fb; display:grid; gap:4px; }
.empty { padding:26px; border-radius:26px; background:white; font-weight:900; color:#63738a; }
@media (max-width:1100px) { .orders-grid { grid-template-columns: 1fr; } .toolbar { flex-direction:column; align-items:stretch; } .filters { justify-content:stretch; } .filters > * { flex:1; } }
@media (max-width:640px) { .meta-grid, .action-box { grid-template-columns:1fr; } .card-top { flex-direction:column; } }
</style>
