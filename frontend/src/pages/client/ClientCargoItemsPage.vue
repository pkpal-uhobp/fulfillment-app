<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { Boxes, Clipboard, Keyboard, Loader2, QrCode, RefreshCcw, Search, Video, X } from '@lucide/vue'
import { Html5Qrcode } from 'html5-qrcode'

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

const qrModalOpen = ref(false)
const qrMode = ref('choice')
const manualQr = ref('')
const qrError = ref('')
const scannerRunning = ref(false)
const scannerMessage = ref('')
let html5QrCode = null

const statusOptions = [
  { value: '', label: 'Все статусы' },
  { value: 'created', label: 'Создано' },
  { value: 'accepted', label: 'Принято' },
  { value: 'received', label: 'Принято на склад' },
  { value: 'stored', label: 'На хранении' },
  { value: 'assigned_to_shipping', label: 'Назначено к отгрузке' },
  { value: 'assigned_to_gate', label: 'Назначено к гейту' },
  { value: 'ready_to_ship', label: 'Готово к отгрузке' },
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
    const text = [item.id, item.qr_code, statusLabel(item.status, 'cargo'), item.status, item.order_id, order?.status]
      .join(' ')
      .toLowerCase()
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

async function openCargoByQr(qr) {
  const value = String(qr || '').trim()
  if (!value) {
    qrError.value = 'Введите или отсканируйте QR-код.'
    return
  }

  qrError.value = ''
  query.value = value
  await stopScanner()
  qrModalOpen.value = false
  qrMode.value = 'choice'
  router.push({ name: 'cargo-by-qr', params: { qrCode: value } })
}

function openQrModal() {
  manualQr.value = String(query.value || '')
  qrError.value = ''
  scannerMessage.value = ''
  qrMode.value = 'choice'
  qrModalOpen.value = true
}

async function closeQrModal() {
  await stopScanner()
  qrModalOpen.value = false
  qrMode.value = 'choice'
  qrError.value = ''
}

function chooseManual() {
  qrMode.value = 'manual'
  manualQr.value = String(query.value || '')
  qrError.value = ''
}

async function chooseScan() {
  qrMode.value = 'scan'
  qrError.value = ''
  scannerMessage.value = 'Запрашиваем доступ к камере...'
  await nextTick()
  await startScanner()
}

async function startScanner() {
  try {
    await stopScanner()
    html5QrCode = new Html5Qrcode('client-qr-reader')
    await html5QrCode.start(
      { facingMode: 'environment' },
      { fps: 10, qrbox: { width: 260, height: 260 } },
      async (decodedText) => {
        if (!decodedText) return
        scannerMessage.value = 'QR найден, открываем карточку...'
        await openCargoByQr(decodedText)
      },
      () => {},
    )
    scannerRunning.value = true
    scannerMessage.value = 'Наведите камеру на QR-код грузового места.'
  } catch (err) {
    scannerRunning.value = false
    scannerMessage.value = ''
    qrError.value = 'Не удалось включить камеру. Разрешите доступ к камере или введите код вручную.'
  }
}

async function stopScanner() {
  if (!html5QrCode) return

  try {
    if (scannerRunning.value) await html5QrCode.stop()
  } catch {
    // Камера могла уже остановиться браузером.
  }

  try {
    await html5QrCode.clear()
  } catch {
    // ignore
  }

  html5QrCode = null
  scannerRunning.value = false
}

async function checkQr() {
  const qr = String(query.value || '').trim()
  if (!qr) {
    openQrModal()
    return
  }
  await openCargoByQr(qr)
}

onMounted(loadCargoItems)
onBeforeUnmount(stopScanner)
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

      <button class="inline-flex min-h-[56px] items-center justify-center gap-2 rounded-2xl bg-[#ff4248] px-5 py-3 font-black text-white shadow-[0_18px_45px_rgba(255,66,72,.25)]" type="button" @click="openQrModal">
        <QrCode class="h-5 w-5" /> Проверить QR
      </button>

      <button class="inline-flex min-h-[56px] items-center justify-center gap-2 rounded-2xl border border-white/15 px-5 py-3 font-black text-white/80 transition hover:bg-white hover:text-[#07101f]" type="button" @click="loadCargoItems">
        <RefreshCcw class="h-5 w-5" /> Обновить
      </button>
    </div>

    <p v-if="error" class="rounded-2xl border border-red-400/30 bg-red-500/10 p-4 font-bold text-red-100">{{ error }}</p>

    <div v-if="loading" class="flex min-h-[260px] items-center justify-center rounded-[2rem] border border-white/10 bg-white/[0.05] text-white/70">
      <Loader2 class="mr-3 h-6 w-6 animate-spin text-[#ff4248]" /> Загружаем грузовые места...
    </div>

    <div v-else class="grid gap-5 xl:grid-cols-3">
      <article v-for="item in filteredItems" :key="item.id" class="rounded-[2rem] bg-white p-6 text-[#061126] shadow-xl shadow-slate-950/5">
        <div class="flex items-start justify-between gap-4">
          <div>
            <p class="text-xs font-black uppercase tracking-[0.35em] text-[#ff4248]">QR</p>
            <h2 class="mt-3 text-2xl font-black tracking-[-0.04em]">{{ item.qr_code || `Место #${item.id}` }}</h2>
          </div>

          <button class="rounded-2xl bg-slate-100 p-3 text-slate-500 transition hover:bg-slate-200" type="button" title="Скопировать QR" @click="copyQr(item.qr_code)">
            <Clipboard class="h-5 w-5" />
          </button>
        </div>

        <div class="mt-6 rounded-2xl bg-slate-50 p-5">
          <p class="text-xs font-black uppercase tracking-[0.3em] text-slate-400">Статус</p>
          <p class="mt-2 text-xl font-black">{{ statusLabel(item.status, 'cargo') }}</p>
        </div>

        <dl class="mt-5 grid gap-2 text-sm font-bold text-slate-600">
          <div><dt class="inline">Заявка: </dt><dd class="inline text-[#ff4248]">#{{ item.order_id }}</dd></div>
          <div><dt class="inline">Маршрут: </dt><dd class="inline">{{ routeText(item) }}</dd></div>
          <div><dt class="inline">Создано: </dt><dd class="inline">{{ formatDateTime(item.created_at) }}</dd></div>
          <div v-if="item.received_at"><dt class="inline">Принято: </dt><dd class="inline">{{ formatDateTime(item.received_at) }}</dd></div>
          <div v-if="item.shipped_at"><dt class="inline">Отгружено: </dt><dd class="inline">{{ formatDateTime(item.shipped_at) }}</dd></div>
        </dl>

        <RouterLink :to="orderRoute(item)" class="mt-6 inline-flex w-full items-center justify-center gap-2 rounded-2xl bg-[#061126] px-5 py-4 font-black text-white transition hover:bg-[#ff4248]">
          <Boxes class="h-5 w-5" /> Открыть заявку
        </RouterLink>
      </article>

      <div v-if="!filteredItems.length" class="xl:col-span-3 rounded-[2rem] border border-white/10 bg-white/[0.05] p-8 text-center text-white/70">
        Грузовые места не найдены. Измените фильтры или нажмите «Проверить QR».
      </div>
    </div>

    <div v-if="qrModalOpen" class="fixed inset-0 z-50 grid place-items-center bg-[#061126]/80 p-4 backdrop-blur" @click.self="closeQrModal">
      <div class="w-full max-w-3xl rounded-[2rem] bg-white p-6 text-[#061126] shadow-2xl">
        <div class="flex items-start justify-between gap-4">
          <div>
            <p class="text-xs font-black uppercase tracking-[0.35em] text-[#ff4248]">QR-контроль</p>
            <h2 class="mt-2 text-4xl font-black tracking-[-0.05em]">Проверить грузовое место</h2>
            <p class="mt-2 max-w-xl font-bold text-slate-500">Выберите способ проверки: отсканируйте QR камерой или введите код вручную.</p>
          </div>
          <button class="rounded-2xl bg-slate-100 p-3 text-slate-500 hover:bg-slate-200" type="button" @click="closeQrModal">
            <X class="h-5 w-5" />
          </button>
        </div>

        <div v-if="qrMode === 'choice'" class="mt-6 grid gap-4 md:grid-cols-2">
          <button class="rounded-[1.5rem] border border-slate-200 bg-slate-50 p-6 text-left transition hover:border-[#ff4248] hover:bg-white" type="button" @click="chooseScan">
            <span class="inline-grid h-14 w-14 place-items-center rounded-2xl bg-[#ff4248] text-white"><Video class="h-7 w-7" /></span>
            <strong class="mt-5 block text-2xl font-black">Сканировать QR</strong>
            <span class="mt-2 block font-bold text-slate-500">Откроется камера браузера. Разрешите доступ и наведите камеру на QR-код.</span>
          </button>

          <button class="rounded-[1.5rem] border border-slate-200 bg-slate-50 p-6 text-left transition hover:border-[#ff4248] hover:bg-white" type="button" @click="chooseManual">
            <span class="inline-grid h-14 w-14 place-items-center rounded-2xl bg-[#061126] text-white"><Keyboard class="h-7 w-7" /></span>
            <strong class="mt-5 block text-2xl font-black">Ввести код</strong>
            <span class="mt-2 block font-bold text-slate-500">Подходит для кода вида QR-TPRO-DEMO-028-01.</span>
          </button>
        </div>

        <div v-else-if="qrMode === 'manual'" class="mt-6 rounded-[1.5rem] bg-slate-50 p-5">
          <label class="block text-xs font-black uppercase tracking-[0.3em] text-slate-400">QR-код</label>
          <input v-model="manualQr" class="mt-3 w-full rounded-2xl border border-slate-200 bg-white px-5 py-4 text-lg font-black outline-none focus:border-[#ff4248]" placeholder="QR-TPRO-DEMO-028-01" @keyup.enter="openCargoByQr(manualQr)" />
          <div class="mt-4 flex flex-col gap-3 sm:flex-row">
            <button class="rounded-2xl bg-[#ff4248] px-6 py-4 font-black text-white" type="button" @click="openCargoByQr(manualQr)">Открыть карточку</button>
            <button class="rounded-2xl bg-slate-200 px-6 py-4 font-black" type="button" @click="qrMode = 'choice'">Назад</button>
          </div>
        </div>

        <div v-else class="mt-6 grid gap-4 lg:grid-cols-[1fr_260px]">
          <div class="overflow-hidden rounded-[1.5rem] bg-[#061126] p-3">
            <div id="client-qr-reader" class="min-h-[320px]"></div>
          </div>
          <div class="rounded-[1.5rem] bg-slate-50 p-5">
            <p class="font-black">{{ scannerMessage || 'Подготовка камеры...' }}</p>
            <p class="mt-3 font-bold text-slate-500">Если камера не включилась, используйте ручной ввод.</p>
            <button class="mt-5 w-full rounded-2xl bg-slate-200 px-5 py-4 font-black" type="button" @click="chooseManual">Ввести код вручную</button>
          </div>
        </div>

        <p v-if="qrError" class="mt-5 rounded-2xl border border-red-200 bg-red-50 p-4 font-bold text-red-700">{{ qrError }}</p>
      </div>
    </div>
  </section>
</template>
