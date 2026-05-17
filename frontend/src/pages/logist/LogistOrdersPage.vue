<template>
  <section class="page">
    <div v-if="error" class="alert">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <div class="toolbar panel">
      <div>
        <p>Заявки</p>
        <h2>Управление статусами</h2>
      </div>

      <div class="filters filters--orders">
        <BaseSelect
          v-model="filters.status"
          :options="orderStatusOptions"
          placeholder="Все статусы"
          @change="loadOrders"
        />

        <BaseSelect
          v-model="filters.handoverType"
          :options="handoverOptions"
          placeholder="Все способы передачи"
          @change="loadOrders"
        />

        <BaseSelect
          v-model="filters.warehouseId"
          :options="warehouseOptions"
          placeholder="Все склады"
        />

        <button type="button" @click="loadOrders">Обновить</button>
      </div>
    </div>

    <div class="orders-grid">
      <article v-for="order in filteredOrders" :key="order.id" class="order-card">
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
            <span>Новый статус</span>
            <BaseSelect
              v-model="forms[order.id].status"
              :options="orderStatusOptions.filter((item) => item.value)"
              placeholder="Выберите статус"
            />
          </label>

          <label>
            <span>Комментарий</span>
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

    <div v-if="!filteredOrders.length && !error" class="empty">Заявки не найдены</div>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import BaseSelect from '@/shared/ui/BaseSelect.vue'
import { apiFetch } from '@/shared/api/http'
import {
  formatDate,
  formatDateTime,
  handoverLabels,
  labelFromMap,
  orderStatusLabels,
  orderStatusOptions,
  unwrapList,
  unwrapOne,
  warehouseName,
} from './logistUtils'

const handoverOptions = [
  { value: '', label: 'Все способы передачи' },
  { value: 'self_delivery', label: 'Сдача на склад' },
  { value: 'pickup', label: 'Забор с адреса' },
]

const orders = ref([])
const warehouses = ref([])
const forms = reactive({})
const history = reactive({})
const filters = reactive({ status: '', handoverType: '', warehouseId: '' })
const error = ref('')
const success = ref('')
const loadingId = ref(null)
const historyLoadingId = ref(null)

const warehouseOptions = computed(() => [
  { value: '', label: 'Все склады' },
  ...warehouses.value.map((warehouse) => ({
    value: String(warehouse.id),
    label: warehouseName(warehouse.id, warehouses.value),
    description: warehouse.address || warehouse.city || '',
  })),
])

const filteredOrders = computed(() => {
  return orders.value.filter((order) => {
    if (!filters.warehouseId) return true

    const id = Number(filters.warehouseId)
    return Number(order.receiving_warehouse_id) === id || Number(order.destination_warehouse_id) === id
  })
})

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
  loadingId.value = order.id
  error.value = ''
  success.value = ''

  const form = forms[order.id]

  try {
    const payload = await apiFetch(`/orders/${order.id}/status`, {
      auth: true,
      method: 'PATCH',
      body: { status: form.status, comment: form.comment || undefined },
    })

    const updated = unwrapOne(payload, 'order') || { ...order, status: form.status }
    const index = orders.value.findIndex((item) => item.id === order.id)

    if (index !== -1) orders.value[index] = updated

    forms[order.id] = { status: updated.status, comment: '' }
    delete history[order.id]
    success.value = `Статус заявки #${order.id} обновлён`
  } catch (err) {
    error.value = err.message || 'Не удалось обновить статус'
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
  } catch {
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
.alert, .success { padding:16px 18px; border-radius:18px; font-weight:950; }
.alert { background:#fee2e2; color:#991b1b; }
.success { background:#d1fae5; color:#065f46; }
.panel, .order-card, .empty { background:white; border-radius:32px; padding:24px; box-shadow:0 18px 42px rgba(7,16,31,.08); }
.toolbar { display:flex; align-items:end; justify-content:space-between; gap:18px; }
.toolbar p { margin:0 0 8px; color:#ff3f4d; letter-spacing:.22em; text-transform:uppercase; font-size:12px; font-weight:950; }
.toolbar h2 { margin:0; font-size:34px; letter-spacing:-.04em; }
.filters { display:grid; gap:12px; align-items:end; }
.filters--orders { grid-template-columns: minmax(180px, 240px) minmax(200px, 270px) minmax(190px, 280px) auto; }
.filters > button { min-height:58px; border:0; border-radius:18px; padding:0 20px; background:#07101f; color:white; font-weight:950; cursor:pointer; font-family:inherit; }
.orders-grid { display:grid; grid-template-columns: repeat(2, minmax(0,1fr)); gap:20px; }
.order-card { display:flex; flex-direction:column; gap:18px; }
.card-top { display:flex; justify-content:space-between; gap:14px; align-items:start; }
.card-top span, .history summary, .action-box label > span { color:#94a3b8; letter-spacing:.22em; text-transform:uppercase; font-size:12px; font-weight:950; }
.card-top h3 { margin:8px 0 0; font-size:28px; letter-spacing:-.04em; }
.card-top em { display:inline-flex; align-items:center; min-height:36px; padding:0 12px; border-radius:999px; background:#ffe4e6; color:#e11d48; font-style:normal; font-weight:950; }
.meta-grid { display:grid; grid-template-columns:repeat(2, minmax(0,1fr)); gap:12px; }
.meta-grid div, .action-box { background:#f6f8fb; border-radius:22px; padding:16px; }
.meta-grid small { display:block; color:#64748b; font-weight:900; margin-bottom:6px; }
.meta-grid b { display:block; font-size:16px; }
.action-box { display:grid; grid-template-columns: 1fr 1fr; gap:14px; align-items:end; }
.action-box label { display:grid; gap:8px; }
input { min-height:58px; border:1px solid #dbe3ef; border-radius:18px; padding:0 18px; background:white; color:#07101f; font-weight:950; font-family:inherit; }
.action-box button { grid-column:1 / -1; min-height:58px; border:0; border-radius:18px; background:#ff3f4d; color:white; font-weight:950; cursor:pointer; font-family:inherit; }
.action-box button:disabled { opacity:.55; cursor:wait; }
.history { background:#f8fafc; border-radius:22px; padding:16px; }
.history summary { cursor:pointer; color:#07101f; }
.history-list { display:grid; gap:10px; margin-top:14px; }
.history-item { display:grid; gap:4px; padding:12px; border-radius:16px; background:white; }
.history-item span, .history-item small, .muted { color:#64748b; font-weight:800; }
.empty { text-align:center; color:#64748b; font-weight:950; }
@media (max-width: 1280px) {
  .filters--orders { grid-template-columns: repeat(2, minmax(0,1fr)); }
  .filters--orders > button { grid-column: 1 / -1; }
}
@media (max-width: 1180px) {
  .orders-grid { grid-template-columns:1fr; }
  .toolbar { display:grid; }
  .filters { grid-template-columns:1fr; }
}
@media (max-width: 680px) {
  .meta-grid, .action-box { grid-template-columns:1fr; }
}
</style>
