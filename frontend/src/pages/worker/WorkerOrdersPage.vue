<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { apiFetch } from '@/shared/api/http'

import {
  cargoTitle,
  formatDateTime,
  gateTitle,
  normalizeCollection,
  safeFileName,
  statusLabel,
  zoneTitle,
} from './workerUtils'

const route = useRoute()
const router = useRouter()

const orders = ref([])
const cargoItems = ref([])

const selectedStatus = ref(String(route.query.status || 'all'))
const selectedOrderId = ref(String(route.query.order_id || ''))
const orderSearch = ref('')
const orderSort = ref('updated_desc')

const loading = ref(false)
const error = ref('')
const notice = ref('')
const printLoading = ref(false)
const archiveLoading = ref(false)

const statusOptions = [
  { value: 'all', label: 'Все статусы' },
  { value: 'created', label: 'Создана' },
  { value: 'waiting_pickup', label: 'Ожидает забора' },
  { value: 'waiting_delivery', label: 'Ожидает сдачи на склад' },
  { value: 'received', label: 'Принята на склад' },
  { value: 'stored', label: 'На хранении' },
  { value: 'assigned_to_shipping', label: 'Назначена к отгрузке' },
  { value: 'shipped', label: 'Отгружена' },
  { value: 'delivered', label: 'Доставлена' },
  { value: 'cancelled', label: 'Отменена' },
]

const sortOptions = [
  { value: 'updated_desc', label: 'Сначала новые' },
  { value: 'updated_asc', label: 'Сначала старые' },
  { value: 'order_desc', label: 'Номер заявки ↓' },
  { value: 'order_asc', label: 'Номер заявки ↑' },
  { value: 'qr_desc', label: 'Больше QR' },
  { value: 'qr_asc', label: 'Меньше QR' },
]

const orderStatusLabels = {
  created: 'Создана',
  waiting_pickup: 'Ожидает забора',
  waiting_delivery: 'Ожидает сдачи на склад',
  received: 'Принята на склад',
  pending_pickup: 'Ожидает забора',
  pending_self_delivery: 'Ожидает сдачи',
  accepted: 'Принята',
  stored: 'На хранении',
  ready_to_ship: 'К отгрузке',
  assigned_to_shipping: 'Назначена к отгрузке',
  assigned_to_shipment: 'Назначена к отгрузке',
  shipped: 'Отгружена',
  delivered: 'Доставлена',
  cancelled: 'Отменена',
}

const orderStatusTones = {
  created: 'gray',
  waiting_pickup: 'amber',
  waiting_delivery: 'amber',
  received: 'blue',
  pending_pickup: 'amber',
  pending_self_delivery: 'amber',
  accepted: 'blue',
  stored: 'green',
  ready_to_ship: 'violet',
  assigned_to_shipping: 'violet',
  assigned_to_shipment: 'violet',
  shipped: 'dark',
  delivered: 'green',
  cancelled: 'gray',
}

const cargoByOrder = computed(() => {
  const map = new Map()

  for (const item of cargoItems.value) {
    const orderId = getCargoOrderId(item)
    if (!orderId) continue

    if (!map.has(orderId)) map.set(orderId, [])
    map.get(orderId).push(item)
  }

  return map
})

const orderCards = computed(() => {
  const map = new Map()

  for (const order of orders.value) {
    const id = getOrderId(order)
    if (!id) continue

    map.set(id, {
      id,
      order,
      cargo: cargoByOrder.value.get(id) || [],
    })
  }

  for (const [id, cargo] of cargoByOrder.value.entries()) {
    if (map.has(id)) continue

    map.set(id, {
      id,
      order: cargo[0]?.order || cargo[0]?.order_info || {},
      cargo,
    })
  }

  return Array.from(map.values())
})

const filteredOrderCards = computed(() => {
  const q = orderSearch.value.trim().toLowerCase()

  const result = orderCards.value.filter((card) => {
    const status = orderStatusValue(card)
    const statusOk = selectedStatus.value === 'all' || status === selectedStatus.value

    const searchableText = [
      card.id,
      orderStatusLabel(card),
      orderRouteTitle(card),
      handoverTitle(card),
      orderClientTitle(card),
      orderAddressTitle(card),
      card.order?.status,
      ...card.cargo.map((item) => [qrValue(item), cargoTitle(item), statusLabel(item.status), zoneTitle(item), gateTitle(item)].join(' ')),
    ]
      .join(' ')
      .toLowerCase()

    return statusOk && (!q || searchableText.includes(q))
  })

  return result.sort(sortOrderCards)
})

const selectedOrder = computed(() => {
  if (selectedOrderId.value) {
    return orderCards.value.find((card) => String(card.id) === String(selectedOrderId.value)) || null
  }

  return filteredOrderCards.value[0] || null
})

const selectedOrderCargo = computed(() => selectedOrder.value?.cargo || [])
const selectedOrderQrCargo = computed(() => selectedOrderCargo.value.filter((item) => qrValue(item)))

const selectedOrderStats = computed(() => {
  const list = selectedOrderCargo.value

  return {
    total: list.length,
    qr: selectedOrderQrCargo.value.length,
    accepted: list.filter((item) => item.status === 'accepted').length,
    stored: list.filter((item) => item.status === 'stored').length,
    ship: list.filter((item) => ['ready_to_ship', 'assigned_to_shipping', 'assigned_to_shipment', 'shipped'].includes(item.status)).length,
  }
})

function getOrderId(order) {
  return String(order?.id || order?.order_id || order?.orderId || '')
}

function getCargoOrderId(item) {
  return String(item?.order_id || item?.orderId || item?.order?.id || item?.order_info?.id || '')
}

function qrValue(item) {
  return item?.qr_code || item?.qrCode || ''
}

function orderStatusValue(card) {
  const status = card.order?.status
  if (status) return status

  const cargo = card.cargo || []
  if (!cargo.length) return 'created'
  if (cargo.every((item) => item.status === 'shipped')) return 'shipped'
  if (cargo.some((item) => ['ready_to_ship', 'assigned_to_shipping', 'assigned_to_shipment'].includes(item.status))) return 'ready_to_ship'
  if (cargo.some((item) => item.status === 'stored')) return 'stored'
  if (cargo.some((item) => item.status === 'accepted')) return 'accepted'

  return 'created'
}

function orderStatusLabel(card) {
  const status = orderStatusValue(card)
  return orderStatusLabels[status] || statusLabel(status)
}

function orderStatusToneClass(card) {
  return orderStatusTones[orderStatusValue(card)] || 'gray'
}

function orderDate(card) {
  return (
    card.order?.updated_at ||
    card.order?.updatedAt ||
    card.order?.created_at ||
    card.order?.createdAt ||
    card.cargo?.[0]?.updated_at ||
    card.cargo?.[0]?.updatedAt ||
    card.cargo?.[0]?.created_at ||
    card.cargo?.[0]?.createdAt ||
    ''
  )
}

function warehouseName(value, keys, fallback = '') {
  for (const key of keys) {
    const field = value?.[key]
    if (typeof field === 'string' && field.trim()) return field
    if (field?.name) return field.name
  }

  return fallback
}

function orderRouteTitle(card) {
  const order = card.order || {}
  const firstCargo = card.cargo?.[0] || {}
  const nestedOrder = firstCargo.order || firstCargo.order_info || {}

  const from =
    warehouseName(order, ['receiving_warehouse', 'receivingWarehouse', 'receiving_warehouse_name', 'receivingWarehouseName']) ||
    warehouseName(firstCargo, ['receiving_warehouse', 'receivingWarehouse', 'receiving_warehouse_name', 'receivingWarehouseName']) ||
    warehouseName(nestedOrder, ['receiving_warehouse', 'receivingWarehouse', 'receiving_warehouse_name', 'receivingWarehouseName']) ||
    'Склад приёмки'

  const to =
    warehouseName(order, ['destination_warehouse', 'destinationWarehouse', 'destination_warehouse_name', 'destinationWarehouseName']) ||
    warehouseName(firstCargo, ['destination_warehouse', 'destinationWarehouse', 'destination_warehouse_name', 'destinationWarehouseName']) ||
    warehouseName(nestedOrder, ['destination_warehouse', 'destinationWarehouse', 'destination_warehouse_name', 'destinationWarehouseName']) ||
    'Склад назначения'

  return `${from} → ${to}`
}

function handoverTitle(card) {
  const type = card.order?.handover_type || card.order?.handoverType

  if (type === 'pickup') return 'Забор с адреса'
  if (type === 'self_delivery') return 'Сдача на склад'

  return type || '—'
}

function orderClientTitle(card) {
  const order = card.order || {}
  const client = order.client || order.user || order.customer || {}

  return (
    order.client_name ||
    order.clientName ||
    order.customer_name ||
    order.customerName ||
    client.name ||
    client.full_name ||
    client.fullName ||
    client.email ||
    '—'
  )
}

function orderAddressTitle(card) {
  const order = card.order || {}

  return (
    order.pickup_address ||
    order.pickupAddress ||
    order.address ||
    order.delivery_address ||
    order.deliveryAddress ||
    '—'
  )
}

function sortOrderCards(a, b) {
  if (orderSort.value === 'order_asc') return Number(a.id) - Number(b.id)
  if (orderSort.value === 'order_desc') return Number(b.id) - Number(a.id)
  if (orderSort.value === 'qr_asc') return a.cargo.length - b.cargo.length
  if (orderSort.value === 'qr_desc') return b.cargo.length - a.cargo.length

  const dateA = new Date(orderDate(a)).getTime() || 0
  const dateB = new Date(orderDate(b)).getTime() || 0

  if (orderSort.value === 'updated_asc') return dateA - dateB
  return dateB - dateA
}

function cleanQuery(query) {
  return Object.fromEntries(
    Object.entries(query).filter(([, value]) => value !== undefined && value !== null && value !== ''),
  )
}

function sameQuery(nextQuery) {
  const current = cleanQuery(route.query)
  const currentKeys = Object.keys(current)
  const nextKeys = Object.keys(nextQuery)

  return currentKeys.length === nextKeys.length && nextKeys.every((key) => String(current[key]) === String(nextQuery[key]))
}

function patchRouteQuery(patch) {
  const nextQuery = cleanQuery({
    ...route.query,
    ...patch,
  })

  if (sameQuery(nextQuery)) return

  router.replace({
    path: route.path,
    query: nextQuery,
    hash: route.hash,
  })
}

function selectOrder(card) {
  selectedOrderId.value = String(card.id)
  patchRouteQuery({ order_id: card.id })
}

async function loadData() {
  loading.value = true
  error.value = ''
  notice.value = ''

  try {
    const [ordersPayload, cargoPayload] = await Promise.all([
      apiFetch('/orders?limit=300', { auth: true }),
      apiFetch('/cargo-items?limit=1000', { auth: true }),
    ])

    orders.value = normalizeCollection(ordersPayload)
    cargoItems.value = normalizeCollection(cargoPayload)

    if (!selectedOrderId.value && filteredOrderCards.value.length) {
      selectedOrderId.value = String(filteredOrderCards.value[0].id)
      patchRouteQuery({ order_id: selectedOrderId.value })
    }
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить заявки'
  } finally {
    loading.value = false
  }
}

async function qrDataUrl(value) {
  const qrModule = await import('qrcode')
  const QRCode = qrModule.default || qrModule

  return QRCode.toDataURL(value, {
    width: 900,
    margin: 3,
    errorCorrectionLevel: 'M',
  })
}

function dataUrlToBlob(dataUrl) {
  const [meta, body] = dataUrl.split(',')
  const mime = meta.match(/:(.*?);/)?.[1] || 'image/png'
  const binary = atob(body)
  const bytes = new Uint8Array(binary.length)

  for (let i = 0; i < binary.length; i += 1) {
    bytes[i] = binary.charCodeAt(i)
  }

  return new Blob([bytes], { type: mime })
}

function downloadDataUrl(dataUrl, filename) {
  const link = document.createElement('a')
  link.href = dataUrl
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
}

function escapeHtml(value) {
  return String(value ?? '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#039;')
}

async function downloadOneQr(item) {
  const value = qrValue(item)

  if (!value) {
    error.value = 'У грузового места нет QR-кода'
    return
  }

  try {
    error.value = ''
    const dataUrl = await qrDataUrl(value)
    downloadDataUrl(dataUrl, `${safeFileName(value)}.png`)
    notice.value = `QR ${value} скачан`
  } catch (err) {
    error.value = err.message || 'Не удалось скачать QR-код'
  }
}

async function downloadSelectedOrderArchive() {
  if (!selectedOrder.value) {
    error.value = 'Сначала выберите заявку'
    return
  }

  if (!selectedOrderQrCargo.value.length) {
    error.value = 'В выбранной заявке нет грузовых мест с QR-кодами'
    return
  }

  archiveLoading.value = true
  error.value = ''
  notice.value = ''

  try {
    const zipModule = await import('jszip')
    const JSZip = zipModule.default || zipModule
    const zip = new JSZip()

    for (const item of selectedOrderQrCargo.value) {
      const value = qrValue(item)
      const dataUrl = await qrDataUrl(value)
      zip.file(`${safeFileName(value)}.png`, dataUrlToBlob(dataUrl))
    }

    const blob = await zip.generateAsync({ type: 'blob' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')

    link.href = url
    link.download = `order-${safeFileName(selectedOrder.value.id)}-qr.zip`
    document.body.appendChild(link)
    link.click()
    link.remove()

    setTimeout(() => URL.revokeObjectURL(url), 1000)
    notice.value = `Архив QR-кодов заявки #${selectedOrder.value.id} скачан`
  } catch (err) {
    error.value = err.message || 'Не удалось скачать архив с QR-кодами'
  } finally {
    archiveLoading.value = false
  }
}

async function printSelectedOrderQr() {
  if (!selectedOrder.value) {
    error.value = 'Сначала выберите заявку'
    return
  }

  if (!selectedOrderQrCargo.value.length) {
    error.value = 'В выбранной заявке нет грузовых мест с QR-кодами'
    return
  }

  printLoading.value = true
  error.value = ''
  notice.value = ''

  try {
    const labels = []

    for (const item of selectedOrderQrCargo.value) {
      const value = qrValue(item)
      const dataUrl = await qrDataUrl(value)

      labels.push({
        qr: value,
        dataUrl,
        orderId: getCargoOrderId(item),
        status: statusLabel(item.status),
        zone: zoneTitle(item),
        gate: gateTitle(item),
      })
    }

    const printWindow = window.open('', '_blank', 'width=1200,height=900')

    if (!printWindow) {
      error.value = 'Браузер заблокировал окно печати. Разрешите всплывающие окна.'
      return
    }

    printWindow.document.write(`
      <!doctype html>
      <html lang="ru">
        <head>
          <meta charset="utf-8">
          <title>QR заявки #${escapeHtml(selectedOrder.value.id)}</title>
          <style>
            * { box-sizing: border-box; }
            body { margin: 0; padding: 24px; font-family: Arial, sans-serif; color: #061126; }
            h1 { margin: 0 0 20px; font-size: 28px; }
            .grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18px; }
            .label { border: 2px solid #061126; border-radius: 18px; padding: 18px; min-height: 260px; page-break-inside: avoid; display: grid; grid-template-columns: 160px 1fr; gap: 18px; align-items: center; }
            img { width: 160px; height: 160px; }
            .eyebrow { margin: 0 0 8px; color: #ff3f4d; font-size: 12px; font-weight: 900; letter-spacing: .22em; text-transform: uppercase; }
            .qr { margin: 0 0 12px; font-size: 22px; font-weight: 900; word-break: break-word; }
            dl { margin: 0; display: grid; gap: 6px; font-size: 14px; }
            dt { color: #64748b; font-weight: 700; }
            dd { margin: 0 0 6px; font-weight: 900; }
            @media print { body { padding: 12px; } .label { break-inside: avoid; } }
          </style>
        </head>
        <body>
          <h1>QR-коды заявки #${escapeHtml(selectedOrder.value.id)}</h1>
          <section class="grid">
            ${labels
              .map(
                (label) => `
                  <article class="label">
                    <img src="${label.dataUrl}" alt="${escapeHtml(label.qr)}">
                    <div>
                      <p class="eyebrow">Fulfillment Transit</p>
                      <p class="qr">${escapeHtml(label.qr)}</p>
                      <dl>
                        <dt>Заявка</dt><dd>#${escapeHtml(label.orderId)}</dd>
                        <dt>Статус</dt><dd>${escapeHtml(label.status)}</dd>
                        <dt>Зона</dt><dd>${escapeHtml(label.zone)}</dd>
                        <dt>Гейт</dt><dd>${escapeHtml(label.gate)}</dd>
                      </dl>
                    </div>
                  </article>
                `,
              )
              .join('')}
          </section>
          <script>
            window.addEventListener('load', () => setTimeout(() => window.print(), 300))
          <\/script>
        </body>
      </html>
    `)

    printWindow.document.close()
    notice.value = `Печатный лист заявки #${selectedOrder.value.id} открыт`
  } catch (err) {
    error.value = err.message || 'Не удалось подготовить печать QR-кодов'
  } finally {
    printLoading.value = false
  }
}

watch(selectedStatus, (status) => {
  patchRouteQuery({ status: status === 'all' ? undefined : status })
})

watch(
  () => route.query.order_id,
  (value) => {
    selectedOrderId.value = String(value || '')
  },
)

watch(
  () => route.query.status,
  (value) => {
    selectedStatus.value = String(value || 'all')
  },
)

watch(filteredOrderCards, (cards) => {
  if (!cards.length) return
  if (selectedOrderId.value && cards.some((card) => String(card.id) === String(selectedOrderId.value))) return

  selectedOrderId.value = String(cards[0].id)
  patchRouteQuery({ order_id: selectedOrderId.value })
})

onMounted(loadData)
</script>

<template>
  <section class="worker-orders-page">
    <header class="page-head">
      <div>
        <p class="eyebrow">Заявки склада</p>
        <h1>Обработка заявок</h1>
        <span>Два рабочих виджета: слева поиск нужной заявки, справа информация по выбранной заявке и действия с QR-кодами.</span>
      </div>

      <div class="head-actions">
        <RouterLink class="red-btn" to="/worker/scan">Открыть QR-сканер</RouterLink>
        <button class="light-btn" type="button" :disabled="loading" @click="loadData">
          {{ loading ? 'Загрузка…' : 'Обновить' }}
        </button>
      </div>
    </header>

    <div v-if="error" class="alert error">{{ error }}</div>
    <div v-if="notice" class="alert success">{{ notice }}</div>

    <section class="orders-widgets">
      <article class="widget search-widget">
        <div class="widget-head">
          <div>
            <p class="eyebrow">Виджет 1</p>
            <h2>Поиск заявки</h2>
            <span>Ищет по номеру заявки, QR-коду, складу, клиенту, адресу и статусу.</span>
          </div>
          <strong class="count-badge">{{ filteredOrderCards.length }}</strong>
        </div>

        <div class="order-controls">
          <label class="field search-wide">
            <span>Поиск</span>
            <input v-model.trim="orderSearch" type="search" placeholder="Например: 104, QR-0007, склад, клиент" />
          </label>

          <label class="field">
            <span>Статус</span>
            <select v-model="selectedStatus">
              <option v-for="option in statusOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </label>

          <label class="field">
            <span>Сортировка</span>
            <select v-model="orderSort">
              <option v-for="option in sortOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </label>
        </div>

        <div v-if="loading" class="empty">Загружаем заявки…</div>
        <div v-else-if="!filteredOrderCards.length" class="empty">Заявки не найдены. Попробуйте изменить поиск или статус.</div>

        <div v-else class="order-list">
          <button
            v-for="card in filteredOrderCards"
            :key="card.id"
            type="button"
            class="order-card"
            :class="{ active: String(selectedOrder?.id) === String(card.id) }"
            @click="selectOrder(card)"
          >
            <span class="order-card-main">
              <strong>Заявка #{{ card.id }}</strong>
              <em>{{ orderRouteTitle(card) }}</em>
              <small>{{ handoverTitle(card) }} · {{ formatDateTime(orderDate(card)) }}</small>
            </span>

            <span class="order-card-meta">
              <span class="status-chip" :class="orderStatusToneClass(card)">{{ orderStatusLabel(card) }}</span>
              <span class="qr-count">{{ card.cargo.length }} QR</span>
            </span>
          </button>
        </div>
      </article>

      <article class="widget process-widget">
        <template v-if="selectedOrder">
          <div class="widget-head process-head">
            <div>
              <p class="eyebrow">Виджет 2</p>
              <h2>Заявка #{{ selectedOrder.id }}</h2>
              <span>Проверьте данные заявки, распечатайте QR-коды или скачайте их архивом.</span>
            </div>
            <span class="status-chip big" :class="orderStatusToneClass(selectedOrder)">{{ orderStatusLabel(selectedOrder) }}</span>
          </div>

          <div class="actions-row">
            <button class="red-btn" type="button" :disabled="printLoading || !selectedOrderQrCargo.length" @click="printSelectedOrderQr">
              {{ printLoading ? 'Готовим печать…' : 'Распечатать QR-коды' }}
            </button>

            <button class="dark-btn" type="button" :disabled="archiveLoading || !selectedOrderQrCargo.length" @click="downloadSelectedOrderArchive">
              {{ archiveLoading ? 'Собираем архив…' : 'Скачать ZIP с QR' }}
            </button>
          </div>

          <section class="stats-grid">
            <article>
              <span>Всего мест</span>
              <strong>{{ selectedOrderStats.total }}</strong>
            </article>
            <article>
              <span>С QR</span>
              <strong>{{ selectedOrderStats.qr }}</strong>
            </article>
            <article>
              <span>Принято</span>
              <strong>{{ selectedOrderStats.accepted }}</strong>
            </article>
            <article>
              <span>К отгрузке</span>
              <strong>{{ selectedOrderStats.ship }}</strong>
            </article>
          </section>

          <section class="info-grid">
            <article>
              <span>Маршрут</span>
              <strong>{{ orderRouteTitle(selectedOrder) }}</strong>
            </article>
            <article>
              <span>Способ передачи</span>
              <strong>{{ handoverTitle(selectedOrder) }}</strong>
            </article>
            <article>
              <span>Клиент</span>
              <strong>{{ orderClientTitle(selectedOrder) }}</strong>
            </article>
            <article>
              <span>Адрес</span>
              <strong>{{ orderAddressTitle(selectedOrder) }}</strong>
            </article>
            <article>
              <span>Дата обновления</span>
              <strong>{{ formatDateTime(orderDate(selectedOrder)) }}</strong>
            </article>
            <article>
              <span>Статус</span>
              <strong>{{ orderStatusLabel(selectedOrder) }}</strong>
            </article>
          </section>

          <section class="cargo-block">
            <div class="cargo-block-head">
              <h3>QR-коды грузовых мест</h3>
              <small>{{ selectedOrderQrCargo.length }} из {{ selectedOrderCargo.length }} доступны для печати и скачивания</small>
            </div>

            <div v-if="!selectedOrderCargo.length" class="empty">В этой заявке пока нет грузовых мест.</div>
            <div v-else class="cargo-list">
              <article v-for="item in selectedOrderCargo" :key="item.id || qrValue(item)" class="cargo-row">
                <span class="cargo-main">
                  <strong>{{ cargoTitle(item) }}</strong>
                  <small>{{ statusLabel(item.status) }} · {{ zoneTitle(item) }} · {{ gateTitle(item) }}</small>
                </span>

                <div class="cargo-actions">
                  <RouterLink
                    v-if="qrValue(item)"
                    class="soft-btn small"
                    :to="`/cargo-items/by-qr/${encodeURIComponent(qrValue(item))}`"
                  >
                    Открыть
                  </RouterLink>
                  <button class="soft-btn small" type="button" :disabled="!qrValue(item)" @click="downloadOneQr(item)">
                    QR PNG
                  </button>
                </div>
              </article>
            </div>
          </section>
        </template>

        <div v-else class="empty big-empty">Выберите заявку в первом виджете, чтобы увидеть информацию и действия с QR-кодами.</div>
      </article>
    </section>
  </section>
</template>

<style scoped>
.worker-orders-page {
  display: grid;
  gap: 26px;
  color: #061126;
}

.page-head,
.widget {
  background: #fff;
  border-radius: 34px;
  padding: 30px;
  box-shadow: 0 18px 62px rgba(15, 23, 42, .08);
}

.page-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
}

.eyebrow {
  margin: 0 0 10px;
  color: #ff3f4d;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .28em;
  text-transform: uppercase;
}

h1,
h2,
h3 {
  margin: 0;
  font-weight: 950;
  letter-spacing: -.05em;
}

h1 {
  font-size: clamp(42px, 6vw, 78px);
  line-height: .9;
}

h2 {
  font-size: clamp(30px, 3vw, 46px);
  line-height: .98;
}

h3 {
  font-size: 24px;
  line-height: 1.1;
}

.page-head span,
.widget-head span {
  display: block;
  margin-top: 14px;
  max-width: 780px;
  color: #5d6d83;
  font-size: 17px;
  line-height: 1.55;
  font-weight: 750;
}

.head-actions,
.actions-row,
.cargo-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 14px;
  flex-wrap: wrap;
}

.red-btn,
.dark-btn,
.light-btn,
.soft-btn {
  min-height: 58px;
  border: 0;
  border-radius: 20px;
  padding: 0 22px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  text-decoration: none;
  font-size: 16px;
  font-weight: 950;
  cursor: pointer;
  white-space: nowrap;
}

.red-btn {
  background: #ff3f4d;
  color: #fff;
  box-shadow: 0 18px 42px rgba(255, 63, 77, .24);
}

.dark-btn {
  background: #061126;
  color: #fff;
}

.light-btn,
.soft-btn {
  background: #eef3f9;
  color: #061126;
}

.soft-btn.small {
  min-height: 44px;
  border-radius: 16px;
  font-size: 14px;
}

button:disabled {
  opacity: .55;
  cursor: not-allowed;
}

.alert,
.empty {
  padding: 18px 22px;
  border-radius: 22px;
  font-weight: 900;
}

.alert.error {
  background: #fff0f1;
  color: #be123c;
}

.alert.success {
  background: #e8fff5;
  color: #047857;
}

.empty {
  background: #f6f9fd;
  color: #64748b;
}

.big-empty {
  min-height: 420px;
  display: grid;
  place-items: center;
  text-align: center;
}

.orders-widgets {
  display: grid;
  grid-template-columns: minmax(360px, .95fr) minmax(520px, 1.35fr);
  gap: 26px;
  align-items: start;
}

.widget-head,
.process-head,
.cargo-block-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 18px;
  margin-bottom: 18px;
}

.count-badge,
.status-chip,
.qr-count {
  border-radius: 999px;
  background: #eef3f9;
  color: #061126;
  padding: 11px 16px;
  font-weight: 950;
  white-space: nowrap;
}

.status-chip.big {
  padding: 14px 20px;
  font-size: 16px;
}

.status-chip.green {
  background: #d7f9e4;
  color: #047857;
}

.status-chip.blue {
  background: #dbeafe;
  color: #1d4ed8;
}

.status-chip.violet {
  background: #ede9fe;
  color: #6d28d9;
}

.status-chip.dark {
  background: #061126;
  color: #fff;
}

.status-chip.red {
  background: #ffe4e6;
  color: #be123c;
}

.status-chip.amber {
  background: #fef3c7;
  color: #92400e;
}

.status-chip.gray {
  background: #eef3f9;
  color: #475569;
}

.order-controls {
  display: grid;
  grid-template-columns: 1fr;
  gap: 14px;
  margin-bottom: 18px;
}

.field {
  display: grid;
  gap: 10px;
}

.field span,
.stats-grid span,
.info-grid span {
  color: #97a5bb;
  font-size: 12px;
  font-weight: 950;
  letter-spacing: .2em;
  text-transform: uppercase;
}

.field input,
.field select {
  width: 100%;
  min-height: 58px;
  border: 1px solid #dbe4ef;
  border-radius: 20px;
  background: #f8fbff;
  color: #061126;
  padding: 0 18px;
  font-size: 17px;
  font-weight: 900;
  box-sizing: border-box;
  outline: none;
}

.field input:focus,
.field select:focus {
  border-color: #ff3f4d;
  box-shadow: 0 0 0 5px rgba(255, 63, 77, .12);
  background: #fff;
}

.order-list {
  display: grid;
  gap: 12px;
  max-height: 720px;
  overflow: auto;
  padding-right: 6px;
}

.order-card {
  width: 100%;
  border: 1px solid #dbe4ef;
  border-radius: 24px;
  background: #f8fbff;
  padding: 18px;
  color: #061126;
  text-align: left;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 16px;
  cursor: pointer;
  transition: transform .18s ease, box-shadow .18s ease, border-color .18s ease;
}

.order-card:hover {
  transform: translateY(-2px);
  border-color: #cbd5e1;
  box-shadow: 0 14px 34px rgba(15, 23, 42, .08);
}

.order-card.active {
  background: #061126;
  color: #fff;
  border-color: #061126;
  box-shadow: 0 18px 46px rgba(6, 17, 38, .22);
}

.order-card-main,
.order-card-meta,
.cargo-main {
  display: grid;
  gap: 8px;
}

.order-card-main strong {
  font-size: 24px;
  font-weight: 950;
}

.order-card-main em {
  color: #5d6d83;
  font-style: normal;
  font-weight: 900;
}

.order-card-main small,
.cargo-row small,
.cargo-block-head small {
  color: #64748b;
  font-weight: 850;
}

.order-card.active .order-card-main em,
.order-card.active .order-card-main small {
  color: rgba(255, 255, 255, .72);
}

.order-card.active .qr-count {
  background: #ff3f4d;
  color: #fff;
}

.order-card-meta {
  justify-items: end;
  align-content: start;
}

.process-widget {
  display: grid;
  gap: 18px;
}

.actions-row {
  justify-content: flex-start;
}

.stats-grid,
.info-grid {
  display: grid;
  gap: 14px;
}

.stats-grid {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.info-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.stats-grid article,
.info-grid article {
  background: #f6f9fd;
  border-radius: 24px;
  padding: 22px;
  display: grid;
  gap: 12px;
}

.stats-grid strong {
  font-size: 34px;
  line-height: 1;
  font-weight: 950;
}

.info-grid strong {
  font-size: 19px;
  line-height: 1.3;
  font-weight: 950;
}

.cargo-block {
  display: grid;
  gap: 14px;
}

.cargo-block-head {
  margin-bottom: 0;
}

.cargo-list {
  display: grid;
  gap: 10px;
  max-height: 420px;
  overflow: auto;
  padding-right: 6px;
}

.cargo-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  border-radius: 20px;
  background: #f6f9fd;
  padding: 16px 18px;
}

.cargo-row strong {
  font-size: 18px;
  font-weight: 950;
}

@media (max-width: 1180px) {
  .orders-widgets,
  .stats-grid,
  .info-grid {
    grid-template-columns: 1fr;
  }

  .order-list,
  .cargo-list {
    max-height: none;
  }
}

@media (max-width: 720px) {
  .page-head,
  .widget-head,
  .process-head,
  .cargo-block-head,
  .cargo-row,
  .order-card {
    grid-template-columns: 1fr;
    flex-direction: column;
    align-items: stretch;
  }

  .head-actions,
  .actions-row,
  .cargo-actions,
  .order-card-meta {
    justify-content: stretch;
    justify-items: stretch;
  }

  .red-btn,
  .dark-btn,
  .light-btn,
  .soft-btn {
    width: 100%;
  }
}
</style>
