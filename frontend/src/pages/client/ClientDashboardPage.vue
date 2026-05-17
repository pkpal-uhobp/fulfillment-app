<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { ArrowRight, Boxes, CalendarCheck, ClipboardList, PackageCheck, PlusCircle, QrCode, Truck } from '@lucide/vue'
import { apiFetch, getCurrentUser } from '@/shared/api/http'
import { byId, compactName, formatDateTime, normalizeCollection, statusLabel } from './clientUtils'

const loading = ref(true)
const error = ref('')
const orders = ref([])
const cargoItems = ref([])
const warehouses = ref([])
const user = ref(getCurrentUser())

const warehouseMap = computed(() => byId(warehouses.value))
const activeOrders = computed(() => orders.value.filter((order) => !['delivered', 'cancelled'].includes(order.status)))
const completedOrders = computed(() => orders.value.filter((order) => ['delivered', 'cancelled'].includes(order.status)))
const cargoInWork = computed(() => cargoItems.value.filter((item) => !['delivered', 'cancelled'].includes(item.status)))
const latestOrders = computed(() => orders.value.slice(0, 5))

function warehouseName(id) {
  return compactName(warehouseMap.value[String(id)]?.name || `Склад #${id}`)
}

async function loadDashboard() {
  loading.value = true
  error.value = ''
  try {
    const [ordersPayload, cargoPayload, warehousesPayload] = await Promise.all([
      apiFetch('/orders?limit=20', { auth: true }),
      apiFetch('/cargo-items?limit=20', { auth: true }).catch(() => ({ cargo_items: [] })),
      apiFetch('/warehouses'),
    ])
    orders.value = normalizeCollection(ordersPayload, 'orders')
    cargoItems.value = normalizeCollection(cargoPayload, 'cargo_items')
    warehouses.value = normalizeCollection(warehousesPayload, 'warehouses')
  } catch (err) {
    error.value = err?.message || 'Не удалось загрузить данные кабинета.'
  } finally {
    loading.value = false
  }
}

onMounted(loadDashboard)
</script>

<template>
  <section class="space-y-6">
    <div class="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
      <div class="overflow-hidden rounded-[2rem] border border-white/10 bg-white/[0.06] p-6 shadow-[0_30px_90px_rgba(0,0,0,0.2)] backdrop-blur sm:p-8">
        <div class="flex flex-col gap-6 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <p class="text-xs font-black uppercase tracking-[0.45em] text-[#ff9ca0]">Личный кабинет</p>
            <h1 class="mt-4 max-w-3xl text-4xl font-black leading-[0.95] tracking-[-0.06em] sm:text-5xl lg:text-6xl">
              Здравствуйте, {{ user?.full_name || 'клиент' }}
            </h1>
            <p class="mt-5 max-w-2xl text-lg leading-8 text-white/68">
              Здесь можно создавать заявки, отслеживать грузовые места по QR, смотреть историю статусов и планировать сдачу товара на склад.
            </p>
          </div>
          <RouterLink to="/client/orders/new" class="inline-flex items-center justify-center gap-3 rounded-2xl bg-[#ff4248] px-6 py-4 font-black text-white shadow-[0_18px_50px_rgba(255,66,72,0.28)] transition hover:-translate-y-0.5">
            Новая заявка <PlusCircle class="h-5 w-5" />
          </RouterLink>
        </div>
      </div>

      <div class="rounded-[2rem] border border-white/10 bg-[#0b1527] p-6 sm:p-8">
        <p class="text-xs font-black uppercase tracking-[0.35em] text-white/45">Быстрые действия</p>
        <div class="mt-5 grid gap-3">
          <RouterLink to="/client/orders" class="group flex items-center justify-between rounded-2xl bg-white/[0.06] px-5 py-4 font-black text-white transition hover:bg-white hover:text-[#07101f]">
            Мои заявки <ArrowRight class="h-5 w-5 transition group-hover:translate-x-1" />
          </RouterLink>
          <RouterLink to="/client/cargo-items" class="group flex items-center justify-between rounded-2xl bg-white/[0.06] px-5 py-4 font-black text-white transition hover:bg-white hover:text-[#07101f]">
            Грузовые места и QR <ArrowRight class="h-5 w-5 transition group-hover:translate-x-1" />
          </RouterLink>
          <RouterLink to="/client/profile" class="group flex items-center justify-between rounded-2xl bg-white/[0.06] px-5 py-4 font-black text-white transition hover:bg-white hover:text-[#07101f]">
            Профиль <ArrowRight class="h-5 w-5 transition group-hover:translate-x-1" />
          </RouterLink>
        </div>
      </div>
    </div>

    <div v-if="error" class="rounded-3xl border border-red-400/30 bg-red-500/10 px-6 py-5 font-bold text-red-100">{{ error }}</div>

    <div class="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
      <div class="rounded-[1.6rem] border border-white/10 bg-white p-5 text-[#07101f]">
        <ClipboardList class="h-7 w-7 text-[#ff4248]" />
        <div class="mt-5 text-4xl font-black">{{ orders.length }}</div>
        <div class="mt-1 text-sm font-bold text-slate-500">Всего заявок</div>
      </div>
      <div class="rounded-[1.6rem] border border-white/10 bg-white p-5 text-[#07101f]">
        <CalendarCheck class="h-7 w-7 text-[#ff4248]" />
        <div class="mt-5 text-4xl font-black">{{ activeOrders.length }}</div>
        <div class="mt-1 text-sm font-bold text-slate-500">В работе</div>
      </div>
      <div class="rounded-[1.6rem] border border-white/10 bg-white p-5 text-[#07101f]">
        <Boxes class="h-7 w-7 text-[#ff4248]" />
        <div class="mt-5 text-4xl font-black">{{ cargoInWork.length }}</div>
        <div class="mt-1 text-sm font-bold text-slate-500">Грузовых мест</div>
      </div>
      <div class="rounded-[1.6rem] border border-white/10 bg-white p-5 text-[#07101f]">
        <PackageCheck class="h-7 w-7 text-[#ff4248]" />
        <div class="mt-5 text-4xl font-black">{{ completedOrders.length }}</div>
        <div class="mt-1 text-sm font-bold text-slate-500">Завершено</div>
      </div>
    </div>

    <div class="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
      <div class="rounded-[2rem] border border-white/10 bg-white p-6 text-[#07101f] shadow-[0_24px_70px_rgba(0,0,0,0.13)]">
        <div class="flex items-center justify-between gap-4">
          <div>
            <p class="text-xs font-black uppercase tracking-[0.35em] text-[#ff4248]">Последние заявки</p>
            <h2 class="mt-2 text-3xl font-black tracking-[-0.05em]">Маршруты товара</h2>
          </div>
          <RouterLink to="/client/orders" class="hidden rounded-2xl bg-slate-100 px-4 py-3 text-sm font-black text-slate-700 transition hover:bg-[#ff4248] hover:text-white sm:inline-flex">Все заявки</RouterLink>
        </div>

        <div v-if="loading" class="mt-6 grid gap-3">
          <div v-for="i in 4" :key="i" class="h-24 animate-pulse rounded-2xl bg-slate-100"></div>
        </div>
        <div v-else-if="!latestOrders.length" class="mt-6 rounded-3xl border border-dashed border-slate-200 bg-slate-50 p-8 text-center">
          <p class="font-bold text-slate-600">У вас пока нет заявок.</p>
          <RouterLink to="/client/orders/new" class="mt-4 inline-flex items-center gap-2 rounded-2xl bg-[#ff4248] px-5 py-3 font-black text-white">Создать первую <ArrowRight class="h-4 w-4" /></RouterLink>
        </div>
        <div v-else class="mt-6 grid gap-3">
          <RouterLink v-for="order in latestOrders" :key="order.id" :to="`/client/orders/${order.id}`" class="grid gap-4 rounded-2xl border border-slate-200 bg-slate-50 p-4 transition hover:border-[#ff4248]/40 hover:bg-white sm:grid-cols-[1fr_auto]">
            <div>
              <div class="flex flex-wrap items-center gap-2">
                <span class="text-lg font-black">Заявка #{{ order.id }}</span>
                <span class="rounded-full bg-[#ff4248]/10 px-3 py-1 text-xs font-black text-[#ff4248]">{{ statusLabel(order.status) }}</span>
              </div>
              <p class="mt-2 text-sm font-bold text-slate-500">{{ warehouseName(order.receiving_warehouse_id) }} → {{ warehouseName(order.destination_warehouse_id) }}</p>
            </div>
            <div class="text-sm font-bold text-slate-500 sm:text-right">
              {{ formatDateTime(order.created_at) }}
            </div>
          </RouterLink>
        </div>
      </div>

      <div class="rounded-[2rem] border border-white/10 bg-[#0b1527] p-6">
        <p class="text-xs font-black uppercase tracking-[0.35em] text-[#ff9ca0]">QR-контроль</p>
        <h2 class="mt-3 text-3xl font-black tracking-[-0.05em]">Отследить грузовое место</h2>
        <p class="mt-3 leading-7 text-white/62">Перейдите в раздел грузовых мест, чтобы проверить QR-код, статус, склад, зону хранения и историю операций.</p>
        <div class="mt-6 grid gap-3">
          <RouterLink to="/client/cargo-items" class="flex items-center justify-between rounded-2xl bg-white px-5 py-4 font-black text-[#07101f] transition hover:-translate-y-0.5">
            Открыть QR-раздел <QrCode class="h-5 w-5" />
          </RouterLink>
          <RouterLink to="/client/orders/new" class="flex items-center justify-between rounded-2xl border border-white/15 px-5 py-4 font-black text-white/80 transition hover:bg-white/10">
            Оформить новую заявку <Truck class="h-5 w-5" />
          </RouterLink>
        </div>
      </div>
    </div>
  </section>
</template>
