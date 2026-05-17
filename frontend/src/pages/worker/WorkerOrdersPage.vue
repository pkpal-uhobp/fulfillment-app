<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { apiFetch } from '@/shared/api/http'
import {
  cargoTitle,
  formatDateTime,
  normalizeCollection,
  safeFileName,
  statusLabel,
  statusTone,
} from './workerUtils'

const route = useRoute()
const router = useRouter()

const cargoItems = ref([])
const selectedOrderId = ref(String(route.query.order_id || ''))
const search = ref('')
const loading = ref(false)
const printing = ref(false)
const downloading = ref(false)
const error = ref('')
const notice = ref('')

const problemStatuses = ['damaged', 'lost', 'cancelled']

const counters = computed(() => ({
  orders: orderCards.value.length,
  places: cargoItems.value.length,
  accepted: cargoItems.value.filter((item) => item.status === 'accepted').length,
  stored: cargoItems.value.filter((item) => item.status === 'stored').length,
  ready: cargoItems.value.filter((item) => ['ready_to_ship', 'shipped'].includes(item.status)).length,
}))

const orderCards = computed(() => {
  const groups = new Map()

  for (const item of cargoItems.value) {
    const orderId = String(item.order_id || item.orderId || item.order?.id || item.orderInfo?.id || '')

    if (!orderId) {
      continue
    }

    if (!groups.has(orderId)) {
      groups.set(orderId, [])
    }

    groups.get(orderId).push(item)
  }

  return Array.from(groups.entries())
    .map(([orderId, items]) => {
      const last = [...items].sort(
        (a, b) =>
          new Date(b.updated_at || b.updatedAt || b.created_at || b.createdAt || 0) -
          new Date(a.updated_at || a.updatedAt || a.created_at || a.createdAt || 0),
      )[0]

      const order = last?.order || last?.order_info || last?.orderInfo || {}
      const statuses = items.reduce((acc, item) => {
        acc[item.status] = (acc[item.status] || 0) + 1
        return acc
      }, {})

      const accepted = items.filter((item) => item.status === 'accepted').length
      const stored = items.filter((item) => item.status === 'stored').length
      const ready = items.filter((item) => ['ready_to_ship', 'shipped'].includes(item.status)).length
      const problems = items.filter((item) => problemStatuses.includes(item.status)).length

      return {
        id: orderId,
        status: order.status || last?.order_status || last?.orderStatus || '',
        transferType: order.transfer_type || order.transferType || last?.transfer_type || last?.transferType || '',
        pickupDate: order.pickup_date || order.pickupDate || last?.pickup_date || last?.pickupDate || '',
        receivingWarehouse:
          order.receiving_warehouse?.name ||
          order.receivingWarehouse?.name ||
          last?.receiving_warehouse?.name ||
          last?.receivingWarehouse?.name ||
          last?.receiving_warehouse_name ||
          last?.receivingWarehouseName ||
          'Склад приёмки',
        destinationWarehouse:
          order.destination_warehouse?.name ||
          order.destinationWarehouse?.name ||
          last?.destination_warehouse?.name ||
          last?.destinationWarehouse?.name ||
          last?.destination_warehouse_name ||
          last?.destinationWarehouseName ||
          'Склад назначения',
        items,
        total: items.length,
        accepted,
        stored,
        ready,
        problems,
        statuses,
        updatedAt: last?.updated_at || last?.updatedAt || last?.created_at || last?.createdAt,
      }
    })
    .sort((a, b) => {
      const aNum = Number(a.id)
      const bNum = Number(b.id)

      if (Number.isFinite(aNum) && Number.isFinite(bNum)) {
        return bNum - aNum
      }

      return String(b.id).localeCompare(String(a.id), 'ru')
    })
})

const filteredOrders = computed(() => {
  const q = search.value.trim().toLowerCase()

  if (!q) {
    return orderCards.value
  }

  return orderCards.value.filter((order) => {
    const text = [
      order.id,
      statusLabel(order.status),
      order.receivingWarehouse,
      order.destinationWarehouse,
      order.transferType,
      order.items.map((item) => `${item.id} ${item.qr_code || item.qrCode || ''}`).join(' '),
    ]
      .join(' ')
      .toLowerCase()

    return text.includes(q)
  })
})

const selectedOrder = computed(() => {
  return orderCards.value.find((order) => String(order.id) === String(selectedOrderId.value)) || null
})

function selectOrder(orderId) {
  selectedOrderId.value = String(orderId)
  notice.value = ''
  error.value = ''
}

function qrValue(item) {
  return item?.qr_code || item?.qrCode || ''
}

function orderStatusLabel(order) {
  return order.status ? statusLabel(order.status) : 'Заявка'
}

function transferLabel(value) {
  const map = {
    pickup: 'Забор с адреса',
    self_delivery: 'Самостоятельная сдача',
  }

  return map[value] || value || '—'
}

function cleanQuery(query) {
  return Object.fromEntries(
    Object.entries(query).filter(([, value]) => value !== undefined && value !== null && value !== ''),
  )
}

function isSameQuery(nextQuery) {
  const current = cleanQuery(route.query)
  const currentKeys = Object.keys(current)
  const nextKeys = Object.keys(nextQuery)

  return (
    currentKeys.length === nextKeys.length &&
    nextKeys.every((key) => String(current[key]) === String(nextQuery[key]))
  )
}

function patchRouteQuery(patch) {
  const nextQuery = cleanQuery({
    ...route.query,
    ...patch,
  })

  if (isSameQuery(nextQuery)) {
    return
  }

  router.replace({
    path: route.path,
    query: nextQuery,
    hash: route.hash,
  })
}

async function loadData() {
  loading.value = true
  error.value = ''
  notice.value = ''

  try {
    // Важно: рабочий не запрашивает /orders, потому что этот эндпоинт может быть закрыт для роли worker.
    // Берём доступные рабочему грузовые места и группируем их по заявкам.
    const payload = await apiFetch('/cargo-items?limit=200', { auth: true })
    cargoItems.value = normalizeCollection(payload, ['cargo_items', 'cargoItems', 'items', 'data'])

    if (!selectedOrderId.value && orderCards.value.length) {
      selectedOrderId.value = String(orderCards.value[0].id)
    }

    if (selectedOrderId.value && !orderCards.value.some((order) => String(order.id) === String(selectedOrderId.value))) {
      selectedOrderId.value = orderCards.value[0] ? String(orderCards.value[0].id) : ''
    }
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить заявки рабочего'
  } finally {
    loading.value = false
  }
}

async function makeQrDataUrl(value) {
  const qrModule = await import('qrcode')
  const QRCode = qrModule.default || qrModule

  return QRCode.toDataURL(value, {
    width: 900,
    margin: 3,
    errorCorrectionLevel: 'M',
  })
}

function downloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')

  link.href = url
  link.download = filename

  document.body.appendChild(link)
  link.click()
  link.remove()

  URL.revokeObjectURL(url)
}

async function downloadOrderZip(order = selectedOrder.value) {
  if (!order?.items?.length) {
    error.value = 'Выберите заявку с грузовыми местами'
    return
  }

  downloading.value = true
  error.value = ''
  notice.value = ''

  try {
    const [{ default: JSZip }, qrModule] = await Promise.all([
      import('jszip'),
      import('qrcode'),
    ])

    const QRCode = qrModule.default || qrModule
    const zip = new JSZip()
    const folderName = `order_${safeFileName(order.id)}_qr`
    const folder = zip.folder(folderName)

    for (const item of order.items) {
      const value = qrValue(item)

      if (!value) {
        continue
      }

      const dataUrl = await QRCode.toDataURL(value, {
        width: 900,
        margin: 3,
        errorCorrectionLevel: 'M',
      })

      folder.file(`${safeFileName(value)}.png`, dataUrl.split(',')[1], {
        base64: true,
      })
    }

    folder.file(
      'README.txt',
      [
        'Fulfillment Transit',
        `QR-коды грузовых мест по заявке #${order.id}`,
        `Всего мест: ${order.items.length}`,
        '',
        ...order.items.map(
          (item) => `${qrValue(item)} — место #${item.id}, статус: ${statusLabel(item.status)}`,
        ),
      ].join('\n'),
    )

    const blob = await zip.generateAsync({ type: 'blob' })
    downloadBlob(blob, `${folderName}.zip`)
    notice.value = `QR-архив для заявки #${order.id} скачан`
  } catch (err) {
    error.value = err.message || 'Не удалось сформировать ZIP с QR-кодами'
  } finally {
    downloading.value = false
  }
}

async function printOrderQr(order = selectedOrder.value) {
  if (!order?.items?.length) {
    error.value = 'Выберите заявку с грузовыми местами'
    return
  }

  const printWindow = window.open('', '_blank', 'width=1100,height=780')

  if (!printWindow) {
    error.value = 'Браузер заблокировал окно печати. Разрешите всплывающие окна для сайта.'
    return
  }

  printing.value = true
  error.value = ''
  notice.value = ''

  printWindow.document.write(`
    <!doctype html>
    <html lang="ru">
      <head>
        <meta charset="UTF-8" />
        <title>QR заявки #${order.id}</title>
        <style>
          body { margin: 0; padding: 32px; font-family: Arial, sans-serif; color: #061126; }
          .loading { font-size: 24px; font-weight: 800; }
        </style>
      </head>
      <body><div class="loading">Готовим QR для печати…</div></body>
    </html>
  `)
  printWindow.document.close()

  try {
    const labels = []

    for (const item of order.items) {
      const value = qrValue(item)

      if (!value) {
        continue
      }

      labels.push({
        id: item.id,
        status: statusLabel(item.status),
        qr: value,
        title: cargoTitle(item),
        image: await makeQrDataUrl(value),
      })
    }

    const html = `
      <!doctype html>
      <html lang="ru">
        <head>
          <meta charset="UTF-8" />
          <title>QR заявки #${order.id}</title>
          <style>
            * { box-sizing: border-box; }
            body {
              margin: 0;
              padding: 22px;
              font-family: Arial, sans-serif;
              color: #061126;
              background: #fff;
            }
            .header {
              margin-bottom: 18px;
              padding-bottom: 14px;
              border-bottom: 2px solid #061126;
            }
            .eyebrow {
              margin: 0 0 6px;
              color: #ff3f4d;
              font-size: 12px;
              font-weight: 900;
              letter-spacing: .22em;
              text-transform: uppercase;
            }
            h1 {
              margin: 0;
              font-size: 30px;
              line-height: 1.1;
            }
            .meta {
              margin-top: 8px;
              font-size: 13px;
              font-weight: 700;
              color: #4b5563;
            }
            .grid {
              display: grid;
              grid-template-columns: repeat(2, minmax(0, 1fr));
              gap: 12px;
            }
            .label {
              min-height: 310px;
              border: 2px solid #061126;
              border-radius: 18px;
              padding: 16px;
              break-inside: avoid;
              display: grid;
              grid-template-columns: 160px minmax(0, 1fr);
              gap: 16px;
              align-items: center;
            }
            .label img {
              width: 160px;
              height: 160px;
              object-fit: contain;
            }
            .label small {
              display: block;
              margin-bottom: 8px;
              color: #ff3f4d;
              font-size: 12px;
              font-weight: 900;
              letter-spacing: .18em;
              text-transform: uppercase;
            }
            .label strong {
              display: block;
              margin-bottom: 10px;
              font-size: 22px;
              line-height: 1.15;
              overflow-wrap: anywhere;
            }
            .label span {
              display: block;
              margin-top: 6px;
              font-size: 14px;
              font-weight: 800;
              color: #4b5563;
            }
            @media print {
              body { padding: 12mm; }
              .grid { gap: 8mm; }
              .label { border-radius: 8mm; page-break-inside: avoid; }
            }
          </style>
        </head>
        <body>
          <section class="header">
            <p class="eyebrow">Fulfillment Transit</p>
            <h1>QR-коды заявки #${order.id}</h1>
            <div class="meta">
              Всего мест: ${labels.length} · ${order.receivingWarehouse} → ${order.destinationWarehouse}
            </div>
          </section>

          <section class="grid">
            ${labels
              .map(
                (label) => `
                  <article class="label">
                    <img src="${label.image}" alt="${label.qr}" />
                    <div>
                      <small>QR / место #${label.id}</small>
                      <strong>${label.qr}</strong>
                      <span>Заявка #${order.id}</span>
                      <span>Статус: ${label.status}</span>
                    </div>
                  </article>
                `,
              )
              .join('')}
          </section>

          <script>
            window.addEventListener('load', () => {
              setTimeout(() => window.print(), 300)
            })
          <\/script>
        </body>
      </html>
    `

    printWindow.document.open()
    printWindow.document.write(html)
    printWindow.document.close()

    notice.value = `Открыта печать QR для заявки #${order.id}`
  } catch (err) {
    printWindow.close()
    error.value = err.message || 'Не удалось подготовить QR для печати'
  } finally {
    printing.value = false
  }
}

watch(selectedOrderId, (orderId) => {
  patchRouteQuery({
    order_id: orderId || undefined,
  })
})

watch(
  () => route.query.order_id,
  (orderId) => {
    if (orderId && String(orderId) !== String(selectedOrderId.value)) {
      selectedOrderId.value = String(orderId)
    }
  },
)

onMounted(loadData)
</script>

<template>
  <section class="orders-page">
    <header class="page-head">
      <div>
        <p class="eyebrow">Заявки склада</p>
        <h1>Заявки</h1>
        <span>
          Выберите заявку слева — справа можно распечатать все QR-коды её грузовых мест
          или скачать их одним ZIP-архивом.
        </span>
      </div>

      <button class="dark-btn" type="button" :disabled="loading" @click="loadData">
        {{ loading ? 'Загрузка…' : 'Обновить' }}
      </button>
    </header>

    <div v-if="error" class="alert error">{{ error }}</div>
    <div v-if="notice" class="alert success">{{ notice }}</div>

    <section class="mini-stats">
      <article>
        <span>Заявок</span>
        <strong>{{ counters.orders }}</strong>
      </article>

      <article>
        <span>Грузовых мест</span>
        <strong>{{ counters.places }}</strong>
      </article>

      <article>
        <span>Принято</span>
        <strong>{{ counters.accepted }}</strong>
      </article>

      <article>
        <span>К отгрузке</span>
        <strong>{{ counters.ready }}</strong>
      </article>
    </section>

    <section class="orders-workspace">
      <aside class="panel orders-panel">
        <div class="panel-head compact">
          <div>
            <p class="eyebrow">Выбор заявки</p>
            <h2>Заявки</h2>
          </div>

          <span class="count-badge">{{ filteredOrders.length }}</span>
        </div>

        <label class="search-field small">
          <span>Поиск</span>
          <input v-model.trim="search" type="text" placeholder="Номер заявки, QR или склад" />
        </label>

        <div v-if="loading" class="empty">Загружаем заявки…</div>
        <div v-else-if="!filteredOrders.length" class="empty">Заявки не найдены.</div>

        <div v-else class="orders-list">
          <button
            v-for="order in filteredOrders"
            :key="order.id"
            type="button"
            class="order-card"
            :class="{ active: String(selectedOrderId) === String(order.id) }"
            @click="selectOrder(order.id)"
          >
            <span class="order-card__top">
              <strong>Заявка #{{ order.id }}</strong>
              <em>{{ order.total }} QR</em>
            </span>

            <span class="order-card__route">
              {{ order.receivingWarehouse }} → {{ order.destinationWarehouse }}
            </span>

            <span class="order-card__meta">
              <b>Принято: {{ order.accepted }}</b>
              <b>Хранение: {{ order.stored }}</b>
              <b>К отгрузке: {{ order.ready }}</b>
              <b v-if="order.problems">Проблемы: {{ order.problems }}</b>
            </span>

            <small>Обновлено: {{ formatDateTime(order.updatedAt) }}</small>
          </button>
        </div>
      </aside>

      <section class="panel selected-panel">
        <template v-if="selectedOrder">
          <div class="selected-head">
            <div>
              <p class="eyebrow">Печать QR</p>
              <h2>Заявка #{{ selectedOrder.id }}</h2>
              <span>
                В этой заявке {{ selectedOrder.total }} грузовых мест. Нажмите «Распечатать QR»,
                чтобы открыть печатный лист со всеми QR-кодами заявки.
              </span>
            </div>

            <em :class="statusTone(selectedOrder.status)">
              {{ orderStatusLabel(selectedOrder) }}
            </em>
          </div>

          <div class="selected-actions">
            <button type="button" class="red-btn" :disabled="printing" @click="printOrderQr(selectedOrder)">
              {{ printing ? 'Готовим…' : 'Распечатать QR' }}
            </button>

            <button type="button" class="dark-btn" :disabled="downloading" @click="downloadOrderZip(selectedOrder)">
              {{ downloading ? 'Готовим ZIP…' : 'Скачать ZIP' }}
            </button>
          </div>

          <div class="qr-summary">
            <article>
              <span>Всего QR</span>
              <strong>{{ selectedOrder.total }}</strong>
            </article>

            <article>
              <span>Принято</span>
              <strong>{{ selectedOrder.accepted }}</strong>
            </article>

            <article>
              <span>На хранении</span>
              <strong>{{ selectedOrder.stored }}</strong>
            </article>

            <article>
              <span>К отгрузке</span>
              <strong>{{ selectedOrder.ready }}</strong>
            </article>
          </div>

          <div class="order-info">
            <div>
              <span>Способ передачи</span>
              <strong>{{ transferLabel(selectedOrder.transferType) }}</strong>
            </div>

            <div>
              <span>Склад приёмки</span>
              <strong>{{ selectedOrder.receivingWarehouse }}</strong>
            </div>

            <div>
              <span>Склад назначения</span>
              <strong>{{ selectedOrder.destinationWarehouse }}</strong>
            </div>

            <div>
              <span>Дата</span>
              <strong>{{ formatDateTime(selectedOrder.pickupDate || selectedOrder.updatedAt) }}</strong>
            </div>
          </div>

          <div class="qr-preview">
            <article v-for="item in selectedOrder.items" :key="item.id" class="qr-row">
              <div>
                <strong>{{ cargoTitle(item) }}</strong>
                <span>Место #{{ item.id }} · заявка #{{ selectedOrder.id }}</span>
              </div>

              <em :class="statusTone(item.status)">{{ statusLabel(item.status) }}</em>
            </article>
          </div>
        </template>

        <template v-else>
          <div class="empty large">Выберите заявку слева, чтобы распечатать её QR-коды.</div>
        </template>
      </section>
    </section>
  </section>
</template>

<style scoped>
.orders-page {
  display: grid;
  gap: 26px;
  color: #061126;
}

.page-head,
.panel,
.mini-stats article {
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

h1 {
  margin: 0;
  font-size: clamp(56px, 7vw, 104px);
  line-height: .88;
  font-weight: 950;
  letter-spacing: -.07em;
}

h2 {
  margin: 0;
  font-size: clamp(30px, 3vw, 48px);
  line-height: 1;
  font-weight: 950;
  letter-spacing: -.05em;
}

.page-head span,
.selected-head span {
  display: block;
  margin-top: 14px;
  max-width: 820px;
  color: #5d6d83;
  font-size: 18px;
  line-height: 1.55;
  font-weight: 750;
}

.red-btn,
.dark-btn {
  min-height: 58px;
  border: 0;
  border-radius: 20px;
  padding: 0 24px;
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

.red-btn:disabled,
.dark-btn:disabled {
  opacity: .6;
  cursor: wait;
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

.empty.large {
  min-height: 360px;
  display: grid;
  place-items: center;
  text-align: center;
}

.mini-stats {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.mini-stats article {
  display: grid;
  gap: 8px;
  min-height: 130px;
  align-content: center;
}

.mini-stats span,
.search-field span,
.qr-summary span,
.order-info span {
  color: #97a5bb;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .22em;
  text-transform: uppercase;
}

.mini-stats strong {
  font-size: 44px;
  line-height: 1;
  font-weight: 950;
  letter-spacing: -.05em;
}

.orders-workspace {
  --workspace-height: clamp(620px, calc(100vh - 180px), 760px);
  display: grid;
  grid-template-columns: minmax(360px, .72fr) minmax(0, 1.28fr);
  gap: 22px;
  align-items: stretch;
}

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  margin-bottom: 18px;
}

.panel-head.compact {
  align-items: flex-start;
}

.count-badge {
  min-width: 52px;
  height: 52px;
  border-radius: 18px;
  background: #061126;
  color: #fff;
  display: grid;
  place-items: center;
  font-weight: 950;
}

.search-field {
  display: grid;
  gap: 10px;
  margin-bottom: 16px;
}

.search-field input {
  width: 100%;
  min-height: 58px;
  border: 1px solid #dbe4ef;
  border-radius: 20px;
  background: #f8fbff;
  color: #061126;
  padding: 0 18px;
  font-size: 16px;
  font-weight: 900;
  box-sizing: border-box;
  outline: none;
}

.search-field input:focus {
  border-color: #ff3f4d;
  box-shadow: 0 0 0 5px rgba(255, 63, 77, .12);
  background: #fff;
}

.orders-list {
  display: grid;
  gap: 12px;
  min-height: 0;
  overflow-y: auto;
  padding: 4px 8px 4px 0;
  scroll-padding-top: 4px;
}

.order-card {
  width: 100%;
  border: 1px solid #dbe4ef;
  border-radius: 24px;
  background: #f8fbff;
  color: #061126;
  padding: 18px;
  text-align: left;
  cursor: pointer;
  display: grid;
  gap: 12px;
  box-sizing: border-box;
  transition: background .18s ease, border-color .18s ease, box-shadow .18s ease;
}

.order-card:hover:not(.active) {
  border-color: rgba(255, 63, 77, .45);
}

.order-card:focus-visible {
  outline: 3px solid rgba(255, 63, 77, .32);
  outline-offset: 2px;
}

.order-card.active {
  background: #061126;
  color: #fff;
  border-color: rgba(255, 255, 255, .28);
  box-shadow:
    inset 0 0 0 1px rgba(255, 255, 255, .18),
    0 18px 44px rgba(6, 17, 38, .18);
}

.order-card__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.order-card strong {
  font-size: 22px;
  font-weight: 950;
}

.order-card em,
em {
  border-radius: 999px;
  padding: 10px 14px;
  font-style: normal;
  font-weight: 950;
  white-space: nowrap;
  text-align: center;
}

.order-card em {
  background: #e8fff5;
  color: #047857;
}

.order-card.active em {
  background: #ff3f4d;
  color: #fff;
}

.order-card__route {
  color: #5d6d83;
  font-weight: 850;
  line-height: 1.4;
}

.order-card.active .order-card__route {
  color: #dbeafe;
}

.order-card__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.order-card__meta b {
  border-radius: 999px;
  padding: 7px 10px;
  background: #eef3f9;
  color: #5d6d83;
  font-size: 12px;
  font-weight: 950;
}

.order-card.active .order-card__meta b {
  background: rgba(255, 255, 255, .1);
  color: #dbeafe;
}

.order-card small {
  color: #6b7a90;
  font-weight: 850;
}

.order-card.active small {
  color: #a9b8ca;
}

.orders-panel,
.selected-panel {
  height: var(--workspace-height);
  min-height: 0;
  box-sizing: border-box;
  overflow: hidden;
}

.orders-panel {
  display: grid;
  grid-template-rows: auto auto minmax(0, 1fr);
}

.selected-panel {
  display: grid;
  grid-template-rows: auto auto auto auto minmax(0, 1fr);
  gap: 22px;
}

.selected-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 18px;
}

.selected-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.qr-summary {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.qr-summary article,
.order-info div,
.qr-row {
  border-radius: 22px;
  background: #f6f9fd;
  padding: 18px;
}

.qr-summary article {
  display: grid;
  gap: 7px;
}

.qr-summary strong {
  font-size: 30px;
  font-weight: 950;
}

.order-info {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.order-info div {
  display: grid;
  gap: 8px;
}

.order-info strong {
  font-size: 17px;
  line-height: 1.35;
  font-weight: 950;
}

.qr-preview {
  display: grid;
  gap: 10px;
  min-height: 0;
  overflow-y: auto;
  padding-right: 6px;
}

.qr-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 12px;
  align-items: center;
}

.qr-row strong {
  display: block;
  font-size: 18px;
  font-weight: 950;
  overflow-wrap: anywhere;
}

.qr-row span {
  display: block;
  margin-top: 5px;
  color: #66758a;
  font-weight: 800;
}

em.green { background: #dcfce7; color: #047857; }
em.red { background: #ffe4e6; color: #be123c; }
em.blue { background: #dbeafe; color: #1d4ed8; }
em.amber { background: #fef3c7; color: #b45309; }
em.violet { background: #ede9fe; color: #6d28d9; }
em.dark { background: #061126; color: #fff; }
em.gray { background: #e2e8f0; color: #475569; }

@media (max-width: 1180px) {
  .orders-workspace,
  .order-info {
    grid-template-columns: 1fr;
  }

  .orders-panel,
  .selected-panel {
    height: auto;
    max-height: none;
    overflow: visible;
  }

  .qr-preview,
  .orders-list {
    max-height: none;
    overflow: visible;
  }
}

@media (max-width: 760px) {
  .page-head,
  .selected-head {
    flex-direction: column;
    align-items: stretch;
  }

  .mini-stats,
  .qr-summary {
    grid-template-columns: 1fr;
  }

  .red-btn,
  .dark-btn {
    width: 100%;
  }

  .qr-row {
    grid-template-columns: 1fr;
  }
}
</style>
