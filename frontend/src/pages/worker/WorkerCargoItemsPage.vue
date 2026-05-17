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
  statusTone,
  zoneTitle,
} from './workerUtils'

const route = useRoute()
const router = useRouter()

const cargoItems = ref([])
const selectedOrderId = ref(String(route.query.order_id || ''))
const selectedStatus = ref(String(route.query.status || 'all'))
const search = ref('')
const loading = ref(false)
const downloading = ref(false)
const error = ref('')
const notice = ref('')

const filterOptions = [
  { value: 'all', label: 'Все' },
  { value: 'accepted', label: 'Принято' },
  { value: 'stored', label: 'Хранение' },
  { value: 'ready_to_ship', label: 'К отгрузке' },
  { value: 'shipped', label: 'Отгружено' },
  { value: 'problem', label: 'Проблемы' },
]

const problemStatuses = ['damaged', 'lost', 'cancelled']

const counters = computed(() => ({
  total: cargoItems.value.length,
  orders: orderCards.value.length,
  accepted: cargoItems.value.filter((item) => item.status === 'accepted').length,
  stored: cargoItems.value.filter((item) => item.status === 'stored').length,
  ship: cargoItems.value.filter((item) => ['ready_to_ship', 'shipped'].includes(item.status)).length,
}))

const filteredCargoItems = computed(() => {
  const q = search.value.trim().toLowerCase()
  return cargoItems.value.filter((item) => {
    const statusOk = selectedStatus.value === 'all'
      || item.status === selectedStatus.value
      || (selectedStatus.value === 'problem' && problemStatuses.includes(item.status))
    const text = [item.qr_code, item.id, item.order_id, statusLabel(item.status), zoneTitle(item), gateTitle(item)]
      .join(' ')
      .toLowerCase()
    return statusOk && (!q || text.includes(q))
  })
})

const orderCards = computed(() => {
  const groups = new Map()
  for (const item of cargoItems.value) {
    const orderId = String(item.order_id || item.orderId || 'Без заявки')
    if (!groups.has(orderId)) groups.set(orderId, [])
    groups.get(orderId).push(item)
  }

  return Array.from(groups.entries())
    .map(([orderId, items]) => {
      const last = [...items].sort((a, b) => new Date(b.updated_at || b.created_at || 0) - new Date(a.updated_at || a.created_at || 0))[0]
      const statuses = items.reduce((acc, item) => {
        acc[item.status] = (acc[item.status] || 0) + 1
        return acc
      }, {})
      const ready = items.filter((item) => ['ready_to_ship', 'shipped'].includes(item.status)).length
      const problems = items.filter((item) => problemStatuses.includes(item.status)).length
      const accepted = items.filter((item) => item.status === 'accepted').length
      const stored = items.filter((item) => item.status === 'stored').length
      return {
        id: orderId,
        items,
        total: items.length,
        ready,
        problems,
        accepted,
        stored,
        statuses,
        updatedAt: last?.updated_at || last?.created_at,
      }
    })
    .sort((a, b) => Number(b.id) - Number(a.id))
})

const selectedOrder = computed(() => orderCards.value.find((order) => String(order.id) === String(selectedOrderId.value)) || null)
const selectedOrderItems = computed(() => selectedOrder.value?.items || [])

function selectOrder(orderId) {
  selectedOrderId.value = String(orderId)
  notice.value = ''
  error.value = ''
}

function setFilter(value) {
  selectedStatus.value = value
}

function qrValue(item) {
  return item?.qr_code || item?.qrCode || ''
}

async function loadData() {
  loading.value = true
  error.value = ''
  notice.value = ''
  try {
    const payload = await apiFetch('/cargo-items?limit=200', { auth: true })
    cargoItems.value = normalizeCollection(payload)

    if (!selectedOrderId.value && orderCards.value.length) {
      selectedOrderId.value = String(orderCards.value[0].id)
    }
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить грузовые места'
  } finally {
    loading.value = false
  }
}

async function qrDataUrl(value) {
  const qrModule = await import('qrcode')
  const QRCode = qrModule.default || qrModule
  return QRCode.toDataURL(value, { width: 900, margin: 3, errorCorrectionLevel: 'M' })
}

function downloadDataUrl(dataUrl, filename) {
  const link = document.createElement('a')
  link.href = dataUrl
  link.download = filename
  document.body.appendChild(link)
  link.click()
  link.remove()
}

async function downloadOneQr(item) {
  const value = qrValue(item)
  if (!value) return
  const dataUrl = await qrDataUrl(value)
  downloadDataUrl(dataUrl, `${safeFileName(value)}.png`)
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
    const [{ default: JSZip }, qrModule] = await Promise.all([import('jszip'), import('qrcode')])
    const QRCode = qrModule.default || qrModule
    const zip = new JSZip()
    const folderName = `order_${safeFileName(order.id)}_qr`
    const folder = zip.folder(folderName)

    for (const item of order.items) {
      const value = qrValue(item)
      if (!value) continue
      const dataUrl = await QRCode.toDataURL(value, { width: 900, margin: 3, errorCorrectionLevel: 'M' })
      const base64 = dataUrl.split(',')[1]
      folder.file(`${safeFileName(value)}.png`, base64, { base64: true })
    }

    folder.file(
      'README.txt',
      [
        `Fulfillment Transit`,
        `QR-коды грузовых мест по заявке #${order.id}`,
        `Всего мест: ${order.items.length}`,
        ``,
        ...order.items.map((item) => `${qrValue(item)} — место #${item.id}, статус: ${statusLabel(item.status)}`),
      ].join('\n'),
    )

    const blob = await zip.generateAsync({ type: 'blob' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${folderName}.zip`
    document.body.appendChild(link)
    link.click()
    link.remove()
    URL.revokeObjectURL(url)

    notice.value = `QR-архив для заявки #${order.id} скачан`
  } catch (err) {
    error.value = err.message || 'Не удалось сформировать ZIP с QR-кодами'
  } finally {
    downloading.value = false
  }
}

watch(selectedOrderId, (orderId) => {
  router.replace({ query: { ...route.query, order_id: orderId || undefined } })
})

watch(selectedStatus, (status) => {
  router.replace({ query: { ...route.query, status: status === 'all' ? undefined : status } })
})

onMounted(loadData)
</script>

<template>
  <section class="cargo-page">
    <header class="page-head">
      <div>
        <p class="eyebrow">QR и грузовые места</p>
        <h1>Заявки склада</h1>
        <span>Выберите заявку из списка — справа появятся все её грузовые места, QR-коды и кнопка скачивания архива.</span>
      </div>
      <div class="head-actions">
        <RouterLink class="red-btn" to="/worker/scan">Проверить QR</RouterLink>
        <button class="dark-btn" type="button" :disabled="loading" @click="loadData">
          {{ loading ? 'Загрузка…' : 'Обновить' }}
        </button>
      </div>
    </header>

    <div v-if="error" class="alert error">{{ error }}</div>
    <div v-if="notice" class="alert success">{{ notice }}</div>

    <section class="mini-stats">
      <article><span>Заявок</span><strong>{{ counters.orders }}</strong></article>
      <article><span>Грузовых мест</span><strong>{{ counters.total }}</strong></article>
      <article><span>Принято</span><strong>{{ counters.accepted }}</strong></article>
      <article><span>К отгрузке</span><strong>{{ counters.ship }}</strong></article>
    </section>

    <section class="orders-workspace">
      <aside class="panel orders-panel">
        <div class="panel-head compact">
          <div>
            <p class="eyebrow">Выбор заявки</p>
            <h2>Заявки</h2>
          </div>
          <span class="count-badge">{{ orderCards.length }}</span>
        </div>

        <div v-if="!orderCards.length" class="empty">Пока нет заявок с грузовыми местами.</div>
        <div v-else class="orders-list">
          <button
            v-for="order in orderCards"
            :key="order.id"
            type="button"
            class="order-card"
            :class="{ active: String(selectedOrderId) === String(order.id) }"
            @click="selectOrder(order.id)"
          >
            <span class="order-card__top">
              <strong>Заявка #{{ order.id }}</strong>
              <em>{{ order.total }} мест</em>
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
              <p class="eyebrow">QR для заявки</p>
              <h2>Заявка #{{ selectedOrder.id }}</h2>
              <span>Все грузовые места этой заявки доступны для скачивания отдельными PNG или одним ZIP-архивом.</span>
            </div>
            <button type="button" class="red-btn" :disabled="downloading" @click="downloadOrderZip(selectedOrder)">
              {{ downloading ? 'Готовим…' : 'Скачать все QR' }}
            </button>
          </div>

          <div class="qr-summary">
            <article><span>Всего мест</span><strong>{{ selectedOrder.total }}</strong></article>
            <article><span>Принято</span><strong>{{ selectedOrder.accepted }}</strong></article>
            <article><span>На хранении</span><strong>{{ selectedOrder.stored }}</strong></article>
            <article><span>К отгрузке</span><strong>{{ selectedOrder.ready }}</strong></article>
          </div>

          <div class="order-qr-list">
            <article v-for="item in selectedOrderItems" :key="item.id" class="qr-row">
              <div>
                <strong>{{ cargoTitle(item) }}</strong>
                <span>Место #{{ item.id }} · заявка #{{ item.order_id || selectedOrder.id }}</span>
              </div>
              <em :class="statusTone(item.status)">{{ statusLabel(item.status) }}</em>
              <RouterLink class="row-link" :to="`/worker/scan?qr=${encodeURIComponent(item.qr_code || '')}`">Открыть</RouterLink>
              <button type="button" class="png-btn" @click="downloadOneQr(item)">PNG</button>
            </article>
          </div>
        </template>
        <template v-else>
          <div class="empty large">Выберите заявку слева, чтобы увидеть QR-коды её грузовых мест.</div>
        </template>
      </section>
    </section>

    <section class="panel filters-panel">
      <div>
        <p class="field-label">Статус</p>
        <div class="filter-pills">
          <button
            v-for="item in filterOptions"
            :key="item.value"
            type="button"
            :class="{ active: selectedStatus === item.value }"
            @click="setFilter(item.value)"
          >{{ item.label }}</button>
        </div>
      </div>
      <label class="search-field">
        <span>Поиск</span>
        <input v-model.trim="search" type="text" placeholder="QR, ID места или номер заявки" />
      </label>
    </section>

    <section class="panel">
      <div class="panel-head">
        <div>
          <p class="eyebrow">Очередь склада</p>
          <h2>{{ filteredCargoItems.length }} грузовых мест</h2>
        </div>
      </div>

      <div v-if="!filteredCargoItems.length" class="empty">Под выбранный фильтр ничего не найдено.</div>
      <div v-else class="cards-grid">
        <article v-for="item in filteredCargoItems" :key="item.id" class="cargo-card">
          <div class="card-top">
            <span class="qr-mark">QR</span>
            <em :class="statusTone(item.status)">{{ statusLabel(item.status) }}</em>
          </div>
          <h3>{{ cargoTitle(item) }}</h3>
          <dl>
            <div><dt>Заявка</dt><dd>#{{ item.order_id || '—' }}</dd></div>
            <div><dt>Зона</dt><dd>{{ zoneTitle(item) }}</dd></div>
            <div><dt>Гейт</dt><dd>{{ gateTitle(item) }}</dd></div>
            <div><dt>Обновлено</dt><dd>{{ formatDateTime(item.updated_at || item.created_at) }}</dd></div>
          </dl>
          <div class="card-actions">
            <RouterLink class="card-btn" :to="`/worker/scan?qr=${encodeURIComponent(item.qr_code || '')}`">Открыть</RouterLink>
            <button type="button" class="download-btn" @click="downloadOneQr(item)">QR PNG</button>
          </div>
        </article>
      </div>
    </section>
  </section>
</template>

<style scoped>
.cargo-page {
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
  font-size: clamp(44px, 6vw, 82px);
  line-height: .9;
  font-weight: 950;
  letter-spacing: -.06em;
}
h2 {
  margin: 0;
  font-size: clamp(28px, 3vw, 42px);
  line-height: 1;
  font-weight: 950;
  letter-spacing: -.04em;
}
h3 {
  margin: 14px 0 18px;
  font-size: 25px;
  line-height: 1.08;
  font-weight: 950;
  overflow-wrap: anywhere;
}
.page-head span,
.selected-head span,
.panel-text {
  display: block;
  margin-top: 14px;
  max-width: 760px;
  color: #5d6d83;
  font-size: 18px;
  line-height: 1.55;
  font-weight: 750;
}
.head-actions,
.selected-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14px;
  flex-wrap: wrap;
}
.red-btn,
.dark-btn,
.row-link,
.png-btn,
.card-btn,
.download-btn {
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
.red-btn { background: #ff3f4d; color: #fff; box-shadow: 0 18px 42px rgba(255, 63, 77, .24); }
.dark-btn, .card-btn { background: #061126; color: #fff; }
.row-link, .download-btn { background: #eef3f9; color: #061126; }
.png-btn { min-height: 48px; border-radius: 16px; background: #061126; color: #fff; }
.red-btn:disabled,
.dark-btn:disabled { opacity: .6; cursor: wait; }
.alert,
.empty {
  padding: 18px 22px;
  border-radius: 22px;
  font-weight: 900;
}
.alert.error { background: #fff0f1; color: #be123c; }
.alert.success { background: #e8fff5; color: #047857; }
.empty { background: #f6f9fd; color: #64748b; }
.empty.large { min-height: 320px; display: grid; place-items: center; text-align: center; }
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
.field-label,
.search-field span {
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
  display: grid;
  grid-template-columns: minmax(360px, .72fr) minmax(0, 1.28fr);
  gap: 22px;
  align-items: start;
}
.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  margin-bottom: 18px;
}
.panel-head.compact { align-items: flex-start; }
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
.orders-list {
  display: grid;
  gap: 12px;
  max-height: 690px;
  overflow: auto;
  padding-right: 4px;
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
  transition: transform .18s ease, background .18s ease, border-color .18s ease, box-shadow .18s ease;
}
.order-card:hover {
  transform: translateY(-2px);
  border-color: rgba(255, 63, 77, .45);
}
.order-card.active {
  background: #061126;
  color: #fff;
  border-color: #061126;
  box-shadow: 0 18px 44px rgba(6, 17, 38, .18);
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
.order-card em {
  border-radius: 999px;
  padding: 8px 10px;
  background: #e8fff5;
  color: #047857;
  font-style: normal;
  font-weight: 950;
  white-space: nowrap;
}
.order-card.active em {
  background: #ff3f4d;
  color: #fff;
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
.order-card.active .order-card__meta b { background: rgba(255,255,255,.1); color: #dbeafe; }
.order-card small { color: #6b7a90; font-weight: 850; }
.order-card.active small { color: #a9b8ca; }
.selected-panel { display: grid; gap: 22px; }
.qr-summary {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}
.qr-summary article {
  border-radius: 22px;
  background: #f6f9fd;
  padding: 18px;
  display: grid;
  gap: 7px;
}
.qr-summary span {
  color: #97a5bb;
  font-size: 12px;
  font-weight: 950;
  letter-spacing: .16em;
  text-transform: uppercase;
}
.qr-summary strong {
  font-size: 30px;
  font-weight: 950;
}
.order-qr-list,
.cards-grid {
  display: grid;
  gap: 12px;
}
.qr-row {
  min-height: 78px;
  border-radius: 24px;
  background: #f6f9fd;
  padding: 14px 16px;
  display: grid;
  grid-template-columns: minmax(220px, 1fr) auto auto auto;
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
em {
  font-style: normal;
  padding: 10px 14px;
  border-radius: 999px;
  font-weight: 950;
  white-space: nowrap;
  text-align: center;
}
em.green { background: #dcfce7; color: #047857; }
em.red { background: #ffe4e6; color: #be123c; }
em.blue { background: #dbeafe; color: #1d4ed8; }
em.amber { background: #fef3c7; color: #b45309; }
em.violet { background: #ede9fe; color: #6d28d9; }
em.dark { background: #061126; color: #fff; }
em.gray { background: #e2e8f0; color: #475569; }
.filters-panel {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(280px, 420px);
  gap: 18px;
  align-items: end;
}
.filter-pills {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 10px;
}
.filter-pills button {
  min-height: 48px;
  border: 0;
  border-radius: 999px;
  padding: 0 18px;
  background: #eef3f9;
  color: #5d6d83;
  font-weight: 950;
  cursor: pointer;
}
.filter-pills button.active {
  background: #ff3f4d;
  color: #fff;
}
.search-field {
  display: grid;
  gap: 10px;
}
.search-field input {
  width: 100%;
  min-height: 62px;
  border: 1px solid #dbe4ef;
  border-radius: 22px;
  background: #f8fbff;
  color: #061126;
  padding: 0 20px;
  font-size: 18px;
  font-weight: 900;
  box-sizing: border-box;
  outline: none;
}
.search-field input:focus {
  border-color: #ff3f4d;
  box-shadow: 0 0 0 5px rgba(255, 63, 77, .12);
  background: #fff;
}
.cards-grid {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}
.cargo-card {
  border-radius: 28px;
  background: #f8fbff;
  padding: 22px;
  min-height: 310px;
  display: flex;
  flex-direction: column;
}
.card-top,
.card-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.qr-mark {
  color: #ff3f4d;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .22em;
}
dl {
  margin: 0;
  display: grid;
  gap: 10px;
}
dl div {
  display: grid;
  gap: 3px;
}
dt {
  color: #97a5bb;
  font-size: 12px;
  font-weight: 950;
  letter-spacing: .18em;
  text-transform: uppercase;
}
dd {
  margin: 0;
  color: #53647b;
  font-weight: 900;
  overflow-wrap: anywhere;
}
.card-actions { margin-top: auto; padding-top: 20px; }
.card-btn,
.download-btn { flex: 1; }
@media (max-width: 1240px) {
  .orders-workspace,
  .filters-panel { grid-template-columns: 1fr; }
  .cards-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
}
@media (max-width: 760px) {
  .page-head,
  .selected-head { flex-direction: column; align-items: stretch; }
  .head-actions,
  .card-actions { flex-direction: column; align-items: stretch; }
  .mini-stats,
  .qr-summary,
  .cards-grid { grid-template-columns: 1fr; }
  .qr-row { grid-template-columns: 1fr; align-items: stretch; }
  .red-btn,
  .dark-btn,
  .row-link,
  .png-btn,
  .card-btn,
  .download-btn { width: 100%; }
}
</style>
