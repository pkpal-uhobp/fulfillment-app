<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowRight, ClipboardList, PlusCircle, RefreshCcw, Search, XCircle } from '@lucide/vue'
import { apiFetch } from '@/shared/api/http'
import BaseSelect from '@/shared/ui/BaseSelect.vue'
import { byId, compactName, formatDateTime, handoverLabel, normalizeCollection, statusLabel } from './clientUtils'

const loading = ref(true)
const error = ref('')
const orders = ref([])
const warehouses = ref([])
const status = ref('')
const handoverType = ref('')
const search = ref('')
const statusOptions = [
  { value: '', label: 'Все статусы' },
  { value: 'created', label: 'Создана' },
  { value: 'waiting_pickup', label: 'Ожидает забора' },
  { value: 'waiting_delivery', label: 'Ожидает сдачи' },
  { value: 'accepted', label: 'Принята' },
  { value: 'received', label: 'Принята на склад' },
  { value: 'stored', label: 'На хранении' },
  { value: 'shipped', label: 'Отгружена' },
  { value: 'delivered', label: 'Доставлена' },
  { value: 'cancelled', label: 'Отменена' },
]

const handoverOptions = [
  { value: '', label: 'Любой способ' },
  { value: 'self_delivery', label: 'Сдача на склад' },
  { value: 'pickup', label: 'Забор с адреса' },
]

const warehouseMap = computed(() => byId(warehouses.value))
const filteredOrders = computed(() => {
  const query = search.value.trim().toLowerCase()
  return orders.value.filter((order) => {
    const matchesStatus = !status.value || order.status === status.value
    const matchesHandover = !handoverType.value || order.handover_type === handoverType.value
    const text = [
      order.id,
      order.status,
      order.handover_type,
      warehouseName(order.receiving_warehouse_id),
      warehouseName(order.destination_warehouse_id),
      order.comment,
    ].join(' ').toLowerCase()
    return matchesStatus && matchesHandover && (!query || text.includes(query))
  })
})

function warehouseName(id) {
  return compactName(warehouseMap.value[String(id)]?.name || `Склад #${id}`)
}

async function loadOrders() {
  loading.value = true
  error.value = ''
  try {
    const params = new URLSearchParams({ limit: '100' })
    if (status.value) params.set('status', status.value)
    if (handoverType.value) params.set('handover_type', handoverType.value)

    const [ordersPayload, warehousesPayload] = await Promise.all([
      apiFetch(`/orders?${params.toString()}`, { auth: true }),
      apiFetch('/warehouses'),
    ])
    orders.value = normalizeCollection(ordersPayload, 'orders')
    warehouses.value = normalizeCollection(warehousesPayload, 'warehouses')
  } catch (err) {
    error.value = err?.message || 'Не удалось загрузить заявки.'
  } finally {
    loading.value = false
  }
}

async function cancelOrder(order) {
  if (!confirm(`Отменить заявку #${order.id}?`)) return
  try {
    await apiFetch(`/orders/${order.id}/cancel`, {
      method: 'PATCH',
      auth: true,
      body: { comment: 'Отменено клиентом из личного кабинета' },
    })
    await loadOrders()
  } catch (err) {
    alert(err?.message || 'Не удалось отменить заявку')
  }
}

onMounted(loadOrders)
</script>

<template>
  <section class="space-y-6">
    <div class="flex flex-col gap-5 rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 backdrop-blur sm:p-8 lg:flex-row lg:items-end lg:justify-between">
      <div>
        <p class="text-xs font-black uppercase tracking-[0.45em] text-[#ff9ca0]">Заявки</p>
        <h1 class="mt-4 text-5xl font-black tracking-[-0.06em]">Мои заявки</h1>
        <p class="mt-4 max-w-2xl text-white/65">Следите за маршрутом, статусом, складом приёмки и складом назначения.</p>
      </div>
      <RouterLink to="/client/orders/new" class="inline-flex items-center justify-center gap-3 rounded-2xl bg-[#ff4248] px-6 py-4 font-black text-white shadow-[0_18px_50px_rgba(255,66,72,0.24)]">
        Создать заявку <PlusCircle class="h-5 w-5" />
      </RouterLink>
    </div>

    <div class="grid gap-3 rounded-[1.5rem] border border-white/10 bg-[#0b1527] p-4 md:grid-cols-[1fr_220px_220px_auto]">
      <label class="flex min-h-[56px] items-center gap-3 rounded-2xl bg-white/[0.06] px-4 py-3">
        <Search class="h-5 w-5 text-white/45" />
        <input v-model="search" class="w-full bg-transparent font-bold outline-none placeholder:text-white/35" placeholder="Поиск по номеру, складу, статусу" />
      </label>
      <BaseSelect v-model="status" :options="statusOptions" compact @change="loadOrders" />
      <BaseSelect v-model="handoverType" :options="handoverOptions" compact @change="loadOrders" />
      <button class="inline-flex min-h-[56px] items-center justify-center gap-2 rounded-2xl border border-white/15 px-5 py-3 font-black text-white/80 transition hover:bg-white hover:text-[#07101f]" @click="loadOrders">
        <RefreshCcw class="h-5 w-5" /> Обновить
      </button>
    </div>

    <div v-if="error" class="rounded-3xl border border-red-400/30 bg-red-500/10 px-6 py-5 font-bold text-red-100">{{ error }}</div>

    <div v-if="loading" class="grid gap-4 lg:grid-cols-2">
      <div v-for="i in 6" :key="i" class="h-48 animate-pulse rounded-[2rem] bg-white/10"></div>
    </div>

    <div v-else-if="!filteredOrders.length" class="rounded-[2rem] border border-dashed border-white/15 bg-white/[0.04] p-10 text-center">
      <ClipboardList class="mx-auto h-12 w-12 text-white/45" />
      <h2 class="mt-5 text-3xl font-black">Заявок не найдено</h2>
      <p class="mt-3 text-white/55">Измените фильтры или создайте новую заявку.</p>
      <RouterLink to="/client/orders/new" class="mt-6 inline-flex items-center gap-2 rounded-2xl bg-[#ff4248] px-6 py-4 font-black text-white">Создать заявку <ArrowRight class="h-5 w-5" /></RouterLink>
    </div>

    <div v-else class="grid gap-4 lg:grid-cols-2">
      <article v-for="order in filteredOrders" :key="order.id" class="flex min-h-[360px] flex-col rounded-[2rem] border border-white/10 bg-white p-5 text-[#07101f] shadow-[0_24px_70px_rgba(0,0,0,0.12)]">
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div>
            <div class="flex flex-wrap items-center gap-2">
              <h2 class="text-2xl font-black tracking-[-0.04em]">Заявка #{{ order.id }}</h2>
              <span class="rounded-full bg-[#ff4248]/10 px-3 py-1 text-xs font-black text-[#ff4248]">{{ statusLabel(order.status) }}</span>
            </div>
            <p class="mt-2 text-sm font-bold text-slate-500">{{ formatDateTime(order.created_at) }}</p>
          </div>
          <span class="rounded-2xl bg-slate-100 px-4 py-2 text-sm font-black text-slate-600">{{ handoverLabel(order.handover_type) }}</span>
        </div>

        <div class="mt-5 rounded-2xl bg-slate-50 p-4">
          <div class="text-xs font-black uppercase tracking-[0.25em] text-slate-400">Маршрут</div>
          <div class="mt-2 grid gap-2 text-sm font-bold text-slate-700">
            <div>{{ warehouseName(order.receiving_warehouse_id) }}</div>
            <div class="text-slate-400">↓</div>
            <div>{{ warehouseName(order.destination_warehouse_id) }}</div>
          </div>
        </div>

        <div class="mt-auto flex flex-wrap gap-3 pt-5">
          <RouterLink :to="`/client/orders/${order.id}`" class="inline-flex flex-1 items-center justify-center gap-2 rounded-2xl bg-[#07101f] px-5 py-4 font-black text-white transition hover:bg-[#ff4248]">
            Открыть <ArrowRight class="h-5 w-5" />
          </RouterLink>
          <button v-if="!['cancelled','delivered','shipped'].includes(order.status)" class="inline-flex min-h-[56px] items-center justify-center gap-2 rounded-2xl border border-red-100 px-5 py-4 font-black text-red-600 transition hover:bg-red-50" @click="cancelOrder(order)">
            <XCircle class="h-5 w-5" /> Отменить
          </button>
        </div>
      </article>
    </div>
  </section>
</template>
