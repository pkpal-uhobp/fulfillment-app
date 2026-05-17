<template>
  <section class="dashboard">
    <div v-if="error" class="alert">{{ error }}</div>

    <div class="stats-grid">
      <article class="stat-card accent">
        <span>Заявки</span>
        <strong>{{ orders.length }}</strong>
        <small>{{ activeOrders }} в работе</small>
      </article>
      <article class="stat-card">
        <span>Грузовые места</span>
        <strong>{{ cargoItems.length }}</strong>
        <small>{{ readyToShip }} готово к отгрузке</small>
      </article>
      <article class="stat-card">
        <span>Отгрузки</span>
        <strong>{{ shipments.length }}</strong>
        <small>{{ plannedShipments }} запланировано</small>
      </article>
      <article class="stat-card">
        <span>Склады</span>
        <strong>{{ warehouses.length }}</strong>
        <small>{{ activeWarehouses }} активных</small>
      </article>
    </div>

    <div class="content-grid">
      <article class="panel">
        <div class="panel-head">
          <div>
            <p>Операции</p>
            <h2>Ближайшие заявки</h2>
          </div>
          <RouterLink to="/logist/orders">Все заявки</RouterLink>
        </div>
        <div class="list">
          <div v-for="order in latestOrders" :key="order.id" class="row-card">
            <div>
              <strong>Заявка #{{ order.id }}</strong>
              <span>{{ labelFromMap(orderStatusLabels, order.status) }} · {{ handoverLabels[order.handover_type] || order.handover_type }}</span>
            </div>
            <small>{{ formatDateTime(order.created_at) }}</small>
          </div>
          <div v-if="!latestOrders.length" class="empty">Заявок пока нет</div>
        </div>
      </article>

      <article class="panel dark">
        <div class="panel-head">
          <div>
            <p>QR</p>
            <h2>Грузы требуют внимания</h2>
          </div>
          <RouterLink to="/logist/cargo-items">Открыть</RouterLink>
        </div>
        <div class="list">
          <div v-for="cargo in attentionCargo" :key="cargo.id" class="row-card dark-row">
            <div>
              <strong>{{ cargo.qr_code }}</strong>
              <span>{{ labelFromMap(cargoStatusLabels, cargo.status) }}</span>
            </div>
            <small>Заказ #{{ cargo.order_id }}</small>
          </div>
          <div v-if="!attentionCargo.length" class="empty dark-empty">Нет проблемных грузов</div>
        </div>
      </article>
    </div>

    <article class="panel wide">
      <div class="panel-head">
        <div>
          <p>Маршруты</p>
          <h2>План отгрузок</h2>
        </div>
        <RouterLink to="/logist/shipments">Управлять</RouterLink>
      </div>
      <div class="shipment-grid">
        <div v-for="shipment in shipments.slice(0, 4)" :key="shipment.id" class="shipment-card">
          <b>Отгрузка #{{ shipment.id }}</b>
          <span>{{ labelFromMap(shipmentStatusLabels, shipment.status) }}</span>
          <small>{{ warehouseName(shipment.destination_warehouse_id, warehouses) }}</small>
          <em>{{ formatDateTime(shipment.planned_departure_at) }}</em>
        </div>
        <div v-if="!shipments.length" class="empty">Отгрузки не запланированы</div>
      </div>
    </article>
  </section>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { apiFetch } from '@/shared/api/http'
import {
  cargoStatusLabels,
  formatDateTime,
  handoverLabels,
  labelFromMap,
  orderStatusLabels,
  shipmentStatusLabels,
  unwrapList,
  warehouseName,
} from './logistUtils'

const orders = ref([])
const cargoItems = ref([])
const shipments = ref([])
const warehouses = ref([])
const error = ref('')

const activeOrders = computed(() => orders.value.filter((item) => !['cancelled', 'delivered'].includes(item.status)).length)
const readyToShip = computed(() => cargoItems.value.filter((item) => item.status === 'ready_to_ship').length)
const plannedShipments = computed(() => shipments.value.filter((item) => ['planned', 'loading'].includes(item.status)).length)
const activeWarehouses = computed(() => warehouses.value.filter((item) => item.is_active).length)
const latestOrders = computed(() => [...orders.value].sort((a, b) => new Date(b.created_at) - new Date(a.created_at)).slice(0, 6))
const attentionCargo = computed(() => cargoItems.value.filter((item) => ['damaged', 'lost', 'ready_to_ship'].includes(item.status)).slice(0, 6))

async function load() {
  error.value = ''
  try {
    const [ordersPayload, cargoPayload, shipmentsPayload, warehousesPayload] = await Promise.all([
      apiFetch('/orders?limit=100', { auth: true }),
      apiFetch('/cargo-items?limit=100', { auth: true }),
      apiFetch('/shipments', { auth: true }),
      apiFetch('/warehouses'),
    ])
    orders.value = unwrapList(ordersPayload, 'orders')
    cargoItems.value = unwrapList(cargoPayload, 'cargo_items')
    shipments.value = unwrapList(shipmentsPayload, 'shipments')
    warehouses.value = unwrapList(warehousesPayload, 'warehouses')
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить данные логиста'
  }
}

onMounted(load)
</script>

<style scoped>
.dashboard { display:grid; gap:24px; }
.alert { padding:16px 18px; border-radius:18px; color:#991b1b; background:#fee2e2; font-weight:800; }
.stats-grid { display:grid; grid-template-columns: repeat(4, minmax(0,1fr)); gap:18px; }
.stat-card { border-radius:30px; padding:26px; background:white; box-shadow:0 18px 42px rgba(7,16,31,.08); display:grid; gap:10px; min-height:150px; }
.stat-card span, .panel p { color:#ff3f4d; text-transform:uppercase; letter-spacing:.22em; font-size:12px; font-weight:900; margin:0; }
.stat-card strong { font-size:46px; letter-spacing:-.06em; }
.stat-card small { color:#607087; font-weight:800; }
.stat-card.accent { background:#07101f; color:white; }
.content-grid { display:grid; grid-template-columns: 1.1fr .9fr; gap:24px; }
.panel { background:white; border-radius:34px; padding:28px; box-shadow:0 18px 42px rgba(7,16,31,.08); }
.panel.dark { background:#07101f; color:white; }
.panel-head { display:flex; align-items:flex-start; justify-content:space-between; gap:16px; margin-bottom:20px; }
.panel h2 { margin:8px 0 0; font-size:30px; letter-spacing:-.04em; }
.panel a { color:#ff3f4d; font-weight:900; text-decoration:none; }
.list { display:grid; gap:12px; }
.row-card { display:flex; justify-content:space-between; gap:16px; padding:18px; border-radius:22px; background:#f3f6fa; }
.row-card strong, .shipment-card b { display:block; font-size:18px; }
.row-card span, .row-card small, .shipment-card small, .shipment-card em { color:#62728a; font-weight:700; font-style:normal; }
.dark-row { background:rgba(255,255,255,.08); }
.dark-row span, .dark-row small, .dark-empty { color:#b8c3d6; }
.shipment-grid { display:grid; grid-template-columns: repeat(4, minmax(0,1fr)); gap:16px; }
.shipment-card { display:grid; gap:10px; padding:20px; border-radius:24px; background:#f3f6fa; }
.shipment-card span { width:max-content; padding:8px 12px; border-radius:999px; color:#ff3f4d; background:#ffe6e8; font-weight:900; }
.empty { padding:20px; border-radius:20px; background:#f3f6fa; color:#687890; font-weight:800; }
@media (max-width: 1100px) { .stats-grid, .shipment-grid { grid-template-columns: repeat(2,1fr); } .content-grid { grid-template-columns: 1fr; } }
@media (max-width: 620px) { .stats-grid, .shipment-grid { grid-template-columns: 1fr; } .row-card { flex-direction:column; } }
</style>
