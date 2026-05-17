<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Boxes, Clock3, FileText, PackageCheck, RefreshCcw, XCircle } from '@lucide/vue'
import { apiFetch } from '@/shared/api/http'
import { byId, compactName, formatDateTime, handoverLabel, historyCommentLabel, normalizeCollection, statusLabel } from './clientUtils'

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const error = ref('')
const order = ref(null)
const history = ref([])
const cargoItems = ref([])
const warehouses = ref([])
const productTypes = ref([])
const cargoPlaceTypes = ref([])

const warehouseMap = computed(() => byId(warehouses.value))
const productTypeMap = computed(() => byId(productTypes.value))
const cargoPlaceTypeMap = computed(() => byId(cargoPlaceTypes.value))

function warehouseName(id) {
  return compactName(warehouseMap.value[String(id)]?.name || `Склад #${id}`)
}

function productTypeName(id) {
  return productTypeMap.value[String(id)]?.name || `Тип #${id}`
}

function cargoPlaceTypeName(id) {
  return cargoPlaceTypeMap.value[String(id)]?.name || `Тип места #${id}`
}

async function loadOrder() {
  loading.value = true
  error.value = ''
  try {
    const orderId = route.params.id
    const [orderPayload, historyPayload, cargoPayload, warehousesPayload, productTypesPayload, cargoPlaceTypesPayload] = await Promise.all([
      apiFetch(`/orders/${orderId}`, { auth: true }),
      apiFetch(`/orders/${orderId}/history`, { auth: true }).catch(() => ({ history: [] })),
      apiFetch(`/cargo-items?order_id=${orderId}&limit=100`, { auth: true }).catch(() => ({ cargo_items: [] })),
      apiFetch('/warehouses'),
      apiFetch('/product-types'),
      apiFetch('/cargo-place-types'),
    ])
    order.value = orderPayload?.order || orderPayload?.data?.order || orderPayload
    history.value = normalizeCollection(historyPayload, 'history')
    cargoItems.value = normalizeCollection(cargoPayload, 'cargo_items')
    warehouses.value = normalizeCollection(warehousesPayload, 'warehouses')
    productTypes.value = normalizeCollection(productTypesPayload, 'product_types')
    cargoPlaceTypes.value = normalizeCollection(cargoPlaceTypesPayload, 'cargo_place_types')
  } catch (err) {
    error.value = err?.message || 'Не удалось загрузить заявку.'
  } finally {
    loading.value = false
  }
}

async function cancelOrder() {
  if (!order.value || !confirm(`Отменить заявку #${order.value.id}?`)) return
  try {
    await apiFetch(`/orders/${order.value.id}/cancel`, {
      method: 'PATCH',
      auth: true,
      body: { comment: 'Отменено клиентом из личного кабинета' },
    })
    await loadOrder()
  } catch (err) {
    alert(err?.message || 'Не удалось отменить заявку')
  }
}

onMounted(loadOrder)
</script>

<template>
  <section class="space-y-6">
    <button class="inline-flex items-center gap-2 rounded-2xl border border-white/15 px-5 py-3 font-black text-white/75 transition hover:bg-white hover:text-[#07101f]" @click="router.back()">
      <ArrowLeft class="h-5 w-5" /> Назад
    </button>

    <div v-if="loading" class="grid gap-6 lg:grid-cols-[1fr_380px]">
      <div class="h-96 animate-pulse rounded-[2rem] bg-white/10"></div>
      <div class="h-96 animate-pulse rounded-[2rem] bg-white/10"></div>
    </div>

    <div v-else-if="error" class="rounded-3xl border border-red-400/30 bg-red-500/10 px-6 py-5 font-bold text-red-100">{{ error }}</div>

    <template v-else-if="order">
      <div class="grid gap-6 lg:grid-cols-[1fr_380px]">
        <div class="rounded-[2rem] border border-white/10 bg-white p-6 text-[#07101f] shadow-[0_24px_70px_rgba(0,0,0,0.12)] sm:p-8">
          <div class="flex flex-wrap items-start justify-between gap-4">
            <div>
              <p class="text-xs font-black uppercase tracking-[0.35em] text-[#ff4248]">Заявка</p>
              <h1 class="mt-3 text-5xl font-black tracking-[-0.06em]">#{{ order.id }}</h1>
              <p class="mt-3 text-slate-500">Создана {{ formatDateTime(order.created_at) }}</p>
            </div>
            <span class="rounded-full bg-[#ff4248]/10 px-4 py-2 text-sm font-black text-[#ff4248]">{{ statusLabel(order.status) }}</span>
          </div>

          <div class="mt-8 grid gap-4 md:grid-cols-2">
            <div class="rounded-3xl bg-slate-50 p-5">
              <div class="text-xs font-black uppercase tracking-[0.25em] text-slate-400">Приёмка</div>
              <div class="mt-2 text-xl font-black">{{ warehouseName(order.receiving_warehouse_id) }}</div>
            </div>
            <div class="rounded-3xl bg-slate-50 p-5">
              <div class="text-xs font-black uppercase tracking-[0.25em] text-slate-400">Назначение</div>
              <div class="mt-2 text-xl font-black">{{ warehouseName(order.destination_warehouse_id) }}</div>
            </div>
            <div class="rounded-3xl bg-slate-50 p-5">
              <div class="text-xs font-black uppercase tracking-[0.25em] text-slate-400">Тип товара</div>
              <div class="mt-2 text-xl font-black">{{ productTypeName(order.product_type_id) }}</div>
            </div>
            <div class="rounded-3xl bg-slate-50 p-5">
              <div class="text-xs font-black uppercase tracking-[0.25em] text-slate-400">Передача</div>
              <div class="mt-2 text-xl font-black">{{ handoverLabel(order.handover_type) }}</div>
            </div>
          </div>

          <div class="mt-6 rounded-3xl border border-slate-200 p-5">
            <div class="flex items-center gap-3">
              <FileText class="h-6 w-6 text-[#ff4248]" />
              <h2 class="text-2xl font-black tracking-[-0.04em]">Состав заявки</h2>
            </div>
            <div class="mt-5 grid gap-3">
              <div v-for="place in order.cargo_places" :key="place.id" class="rounded-2xl bg-slate-50 p-4">
                <div class="flex flex-wrap items-center justify-between gap-3">
                  <div class="font-black">{{ cargoPlaceTypeName(place.cargo_place_type_id) }}</div>
                  <div class="rounded-full bg-white px-3 py-1 text-sm font-black text-slate-600">{{ place.quantity }} шт.</div>
                </div>
                <div class="mt-3 grid gap-2 text-sm text-slate-500 sm:grid-cols-4">
                  <div>Вес: <b>{{ place.weight_per_place_kg || '—' }}</b> кг</div>
                  <div>Длина: <b>{{ place.length_cm || '—' }}</b> см</div>
                  <div>Ширина: <b>{{ place.width_cm || '—' }}</b> см</div>
                  <div>Высота: <b>{{ place.height_cm || '—' }}</b> см</div>
                </div>
              </div>
            </div>
          </div>

          <div v-if="order.pickup" class="mt-6 rounded-3xl border border-slate-200 p-5">
            <h2 class="text-2xl font-black tracking-[-0.04em]">Забор с адреса</h2>
            <div class="mt-4 grid gap-3 text-sm font-bold text-slate-600 md:grid-cols-2">
              <div>Адрес: {{ order.pickup.pickup_address }}</div>
              <div>Дата: {{ order.pickup.pickup_date }}</div>
              <div>Время: {{ order.pickup.pickup_time_from || '—' }} — {{ order.pickup.pickup_time_to || '—' }}</div>
              <div>Статус: {{ statusLabel(order.pickup.status) }}</div>
            </div>
          </div>
        </div>

        <aside class="space-y-6">
          <div class="rounded-[2rem] border border-white/10 bg-[#0b1527] p-6">
            <div class="flex items-center justify-between gap-4">
              <h2 class="text-2xl font-black tracking-[-0.04em]">Действия</h2>
              <button class="rounded-2xl border border-white/15 p-3 text-white/70 transition hover:bg-white hover:text-[#07101f]" @click="loadOrder"><RefreshCcw class="h-5 w-5" /></button>
            </div>
            <div class="mt-5 grid gap-3">
              <RouterLink to="/client/cargo-items" class="flex items-center justify-between rounded-2xl bg-white px-5 py-4 font-black text-[#07101f]">
                Грузовые места <Boxes class="h-5 w-5" />
              </RouterLink>
              <button v-if="!['cancelled','delivered','shipped'].includes(order.status)" class="flex items-center justify-between rounded-2xl border border-red-400/30 bg-red-500/10 px-5 py-4 font-black text-red-100" @click="cancelOrder">
                Отменить заявку <XCircle class="h-5 w-5" />
              </button>
            </div>
          </div>

          <div class="rounded-[2rem] border border-white/10 bg-white/[0.06] p-6">
            <div class="flex items-center gap-3">
              <Clock3 class="h-6 w-6 text-[#ff9ca0]" />
              <h2 class="text-2xl font-black tracking-[-0.04em]">История</h2>
            </div>
            <div class="mt-5 grid gap-3">
              <div v-for="event in history" :key="event.id" class="rounded-2xl bg-white/[0.06] p-4">
                <div class="font-black">{{ statusLabel(event.new_status) }}</div>
                <div class="mt-1 text-sm text-white/45">{{ formatDateTime(event.changed_at) }}</div>
                <div v-if="event.comment" class="mt-2 text-sm text-white/65">{{ historyCommentLabel(event.comment) }}</div>
              </div>
              <div v-if="!history.length" class="rounded-2xl bg-white/[0.06] p-4 text-sm text-white/45">История пока пустая.</div>
            </div>
          </div>
        </aside>
      </div>

      <div class="rounded-[2rem] border border-white/10 bg-white p-6 text-[#07101f]">
        <div class="flex items-center gap-3">
          <PackageCheck class="h-6 w-6 text-[#ff4248]" />
          <h2 class="text-2xl font-black tracking-[-0.04em]">Грузовые места по заявке</h2>
        </div>
        <div v-if="!cargoItems.length" class="mt-5 rounded-3xl bg-slate-50 p-6 text-center font-bold text-slate-500">Грузовые места появятся после приёмки товара на складе.</div>
        <div v-else class="mt-5 grid gap-3 md:grid-cols-2 xl:grid-cols-3">
          <RouterLink v-for="item in cargoItems" :key="item.id" :to="`/client/cargo-items?qr=${item.qr_code}`" class="rounded-3xl border border-slate-200 bg-slate-50 p-5 transition hover:border-[#ff4248]/40 hover:bg-white">
            <div class="text-xs font-black uppercase tracking-[0.25em] text-slate-400">QR</div>
            <div class="mt-2 break-all text-xl font-black">{{ item.qr_code }}</div>
            <div class="mt-4 rounded-full bg-[#ff4248]/10 px-3 py-1 text-sm font-black text-[#ff4248]">{{ statusLabel(item.status, 'cargo') }}</div>
          </RouterLink>
        </div>
      </div>
    </template>
  </section>
</template>
