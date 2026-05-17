<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { Boxes, Clipboard, RefreshCcw, Search, QrCode } from '@lucide/vue'
import { apiFetch } from '@/shared/api/http'
import BaseSelect from '@/shared/ui/BaseSelect.vue'
import { byId, compactName, formatDateTime, normalizeCollection, statusLabel } from './clientUtils'

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const error = ref('')
const cargoItems = ref([])
const orders = ref([])
const warehouses = ref([])
const query = ref(route.query.qr || '')
const status = ref('')

const statusOptions = [
  { value: '', label: 'Все статусы' },
  { value: 'created', label: 'Создано' },
  { value: 'accepted', label: 'Принято' },
  { value: 'received', label: 'Принято на склад' },
  { value: 'stored', label: 'На хранении' },
  { value: 'assigned_to_shipping', label: 'Назначено к отгрузке' },
  { value: 'shipped', label: 'Отгружено' },
  { value: 'delivered', label: 'Доставлено' },
]

const orderMap = computed(() => byId(orders.value))
const warehouseMap = computed(() => byId(warehouses.value))
const filteredItems = computed(() => {
  const q = String(query.value || '').trim().toLowerCase()
  return cargoItems.value.filter((item) => {
    const matchesStatus = !status.value || item.status === status.value
    const order = orderMap.value[String(item.order_id)]
    const text = [item.id, item.qr_code, statusLabel(item.status, 'cargo'), item.status, item.order_id, order?.status].join(' ').toLowerCase()
    return matchesStatus && (!q || text.includes(q))
  })
})

function orderRoute(item) {
  return `/client/orders/${item.order_id}`
}

function routeText(item) {
  const order = orderMap.value[String(item.order_id)]
  if (!order) return 'Маршрут появится после загрузки заявки'
  const receiving = compactName(warehouseMap.value[String(order.receiving_warehouse_id)]?.name || `Склад #${order.receiving_warehouse_id}`)
  const destination = compactName(warehouseMap.value[String(order.destination_warehouse_id)]?.name || `Склад #${order.destination_warehouse_id}`)
  return `${receiving} → ${destination}`
}

async function loadCargoItems() {
  loading.value = true
  error.value = ''
  try {
    const [cargoPayload, ordersPayload, warehousesPayload] = await Promise.all([
      apiFetch('/cargo-items?limit=100', { auth: true }),
      apiFetch('/orders?limit=100', { auth: true }),
      apiFetch('/warehouses'),
    ])
    cargoItems.value = normalizeCollection(cargoPayload, 'cargo_items')
    orders.value = normalizeCollection(ordersPayload, 'orders')
    warehouses.value = normalizeCollection(warehousesPayload, 'warehouses')
  } catch (err) {
    error.value = err?.message || 'Не удалось загрузить грузовые места.'
  } finally {
    loading.value = false
  }
}

function copyQr(value) {
  if (navigator?.clipboard && value) navigator.clipboard.writeText(value)
}

async function checkQr() {
  const qr = String(query.value || '').trim()
  if (!qr) return
  loading.value = true
  error.value = ''
  try {
    const payload = await apiFetch(`/cargo-items/scan?qr_code=${encodeURIComponent(qr)}`, { auth: true })
    const item = payload?.cargo_item || payload?.data?.cargo_item || payload
    if (item?.qr_code && !cargoItems.value.some((current) => String(current.id) === String(item.id))) {
      cargoItems.value = [item, ...cargoItems.value]
    }
    router.replace({ path: '/client/cargo-items', query: { qr } })
  } catch (err) {
    error.value = err?.message || 'QR-код не найден или нет доступа.'
  } finally {
    loading.value = false
  }
}

onMounted(loadCargoItems)
</script>

<template>
  <section class="space-y-6">
    <div class="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 backdrop-blur sm:p-8">
      <p class="text-xs font-black uppercase tracking-[0.45em] text-[#ff9ca0]">QR и статусы</p>
      <h1 class="mt-4 text-5xl font-black tracking-[-0.06em]">Грузовые места</h1>
      <p class="mt-4 max-w-3xl text-white/65">Проверяйте QR-коды, текущие статусы, маршрут и связь грузового места с заявкой.</p>
    </div>

    <div class="grid gap-3 rounded-[1.5rem] border border-white/10 bg-[#0b1527] p-4 lg:grid-cols-[1fr_240px_auto_auto]">
      <label class="flex min-h-[56px] items-center gap-3 rounded-2xl bg-white/[0.06] px-4 py-3">
        <Search class="h-5 w-5 text-white/45" />
        <input v-model="query" class="w-full bg-transparent font-bold outline-none placeholder:text-white/35" placeholder="Введите QR или номер заявки" @keyup.enter="checkQr" />
      </label>
      <BaseSelect v-model="status" :options="statusOptions" compact />
      <button class="inline-flex min-h-[56px] items-center justify-center gap-2 rounded-2xl bg-[#ff4248] px-5 py-3 font-black text-white" @click="checkQr">
        <QrCode class="h-5 w-5" /> Проверить QR
      </button>
      <button class="inline-flex min-h-[56px] items-center justify-center gap-2 rounded-2xl border border-white/15 px-5 py-3 font-black text-white/80 transition hover:bg-white hover:text-[#07101f]" @click="loadCargoItems">
        <RefreshCcw class="h-5 w-5" /> Обновить
      </button>
    </div>

    <div v-if="error" class="rounded-3xl border border-red-400/30 bg-red-500/10 px-6 py-5 font-bold text-red-100">{{ error }}</div>

    <div v-if="loading" class="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
      <div v-for="i in 6" :key="i" class="h-80 animate-pulse rounded-[2rem] bg-white/10"></div>
    </div>

    <div v-else-if="!filteredItems.length" class="rounded-[2rem] border border-dashed border-white/15 bg-white/[0.04] p-10 text-center">
      <Boxes class="mx-auto h-12 w-12 text-white/45" />
      <h2 class="mt-5 text-3xl font-black">Грузовых мест нет</h2>
      <p class="mt-3 text-white/55">Они появятся после приёмки товара на складе или после проверки QR.</p>
    </div>

    <div v-else class="grid gap-5 md:grid-cols-2 xl:grid-cols-3">
      <article v-for="item in filteredItems" :key="item.id" class="flex min-h-[390px] flex-col rounded-[2rem] border border-white/10 bg-white p-6 text-[#07101f] shadow-[0_24px_70px_rgba(0,0,0,0.12)]">
        <div class="flex items-start justify-between gap-3">
          <div class="min-w-0">
            <p class="text-xs font-black uppercase tracking-[0.3em] text-[#ff4248]">QR</p>
            <h2 class="mt-2 break-all text-2xl font-black tracking-[-0.04em]">{{ item.qr_code }}</h2>
          </div>
          <button class="rounded-2xl bg-slate-100 p-3 text-slate-500 transition hover:bg-[#ff4248] hover:text-white" @click="copyQr(item.qr_code)">
            <Clipboard class="h-5 w-5" />
          </button>
        </div>

        <div class="mt-5 rounded-2xl bg-slate-50 p-4">
          <div class="text-xs font-black uppercase tracking-[0.25em] text-slate-400">Статус</div>
          <div class="mt-2 text-xl font-black">{{ statusLabel(item.status, 'cargo') }}</div>
        </div>

        <div class="mt-4 grid gap-2 text-sm font-bold text-slate-500">
          <div>Заявка: <RouterLink :to="orderRoute(item)" class="text-[#ff4248] hover:underline">#{{ item.order_id }}</RouterLink></div>
          <div>{{ routeText(item) }}</div>
          <div>Создано: {{ formatDateTime(item.created_at) }}</div>
          <div v-if="item.received_at">Принято: {{ formatDateTime(item.received_at) }}</div>
          <div v-if="item.shipped_at">Отгружено: {{ formatDateTime(item.shipped_at) }}</div>
        </div>

        <RouterLink :to="orderRoute(item)" class="mt-auto inline-flex w-full items-center justify-center rounded-2xl bg-[#07101f] px-5 py-4 font-black text-white transition hover:bg-[#ff4248]">
          Открыть заявку
        </RouterLink>
      </article>
    </div>
  </section>
</template>
