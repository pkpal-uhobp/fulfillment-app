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
const selectedStatus = ref(String(route.query.status || 'all'))
const search = ref('')
const loading = ref(false)
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
  accepted: cargoItems.value.filter((item) => item.status === 'accepted').length,
  stored: cargoItems.value.filter((item) => item.status === 'stored').length,
  ship: cargoItems.value.filter((item) => ['ready_to_ship', 'shipped'].includes(item.status)).length,
}))

const filteredCargoItems = computed(() => {
  const q = search.value.trim().toLowerCase()

  return cargoItems.value.filter((item) => {
    const statusOk =
      selectedStatus.value === 'all' ||
      item.status === selectedStatus.value ||
      (selectedStatus.value === 'problem' && problemStatuses.includes(item.status))

    const text = [
      item.qr_code,
      item.qrCode,
      item.id,
      item.order_id,
      item.orderId,
      statusLabel(item.status),
      zoneTitle(item),
      gateTitle(item),
    ]
      .join(' ')
      .toLowerCase()

    return statusOk && (!q || text.includes(q))
  })
})

function setFilter(value) {
  selectedStatus.value = value
}

function qrValue(item) {
  return item?.qr_code || item?.qrCode || ''
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
    const payload = await apiFetch('/cargo-items?limit=200', { auth: true })
    cargoItems.value = normalizeCollection(payload)
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить грузовые места'
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

  if (!value) {
    error.value = 'У грузового места нет QR-кода'
    return
  }

  try {
    const dataUrl = await qrDataUrl(value)
    downloadDataUrl(dataUrl, `${safeFileName(value)}.png`)
    notice.value = `QR ${value} скачан`
  } catch (err) {
    error.value = err.message || 'Не удалось скачать QR-код'
  }
}

watch(selectedStatus, (status) => {
  patchRouteQuery({
    status: status === 'all' ? undefined : status,
  })
})

onMounted(loadData)
</script>

<template>
  <section class="cargo-page">
    <header class="page-head">
      <div>
        <p class="eyebrow">Грузовые места</p>
        <h1>Грузовые места</h1>
        <span>
          Здесь остался общий список грузовых мест. Заявки вынесены в отдельную вкладку
          «Заявки» в панели рабочего.
        </span>
      </div>

      <div class="head-actions">
        <RouterLink class="red-btn" to="/worker/orders">Перейти к заявкам</RouterLink>
        <RouterLink class="dark-btn" to="/worker/scan">Проверить QR</RouterLink>
        <button class="light-btn" type="button" :disabled="loading" @click="loadData">
          {{ loading ? 'Загрузка…' : 'Обновить' }}
        </button>
      </div>
    </header>

    <div v-if="error" class="alert error">{{ error }}</div>
    <div v-if="notice" class="alert success">{{ notice }}</div>

    <section class="mini-stats">
      <article>
        <span>Грузовых мест</span>
        <strong>{{ counters.total }}</strong>
      </article>

      <article>
        <span>Принято</span>
        <strong>{{ counters.accepted }}</strong>
      </article>

      <article>
        <span>На хранении</span>
        <strong>{{ counters.stored }}</strong>
      </article>

      <article>
        <span>К отгрузке</span>
        <strong>{{ counters.ship }}</strong>
      </article>
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
          >
            {{ item.label }}
          </button>
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

      <div v-if="!filteredCargoItems.length" class="empty">
        Под выбранный фильтр ничего не найдено.
      </div>

      <div v-else class="cards-grid">
        <article v-for="item in filteredCargoItems" :key="item.id" class="cargo-card">
          <div class="card-top">
            <span class="qr-mark">QR</span>

            <em :class="statusTone(item.status)">
              {{ statusLabel(item.status) }}
            </em>
          </div>

          <h3>{{ cargoTitle(item) }}</h3>

          <dl>
            <div>
              <dt>Заявка</dt>
              <dd>#{{ item.order_id || item.orderId || '—' }}</dd>
            </div>

            <div>
              <dt>Зона</dt>
              <dd>{{ zoneTitle(item) }}</dd>
            </div>

            <div>
              <dt>Гейт</dt>
              <dd>{{ gateTitle(item) }}</dd>
            </div>

            <div>
              <dt>Обновлено</dt>
              <dd>{{ formatDateTime(item.updated_at || item.created_at) }}</dd>
            </div>
          </dl>

          <div class="card-actions">
            <RouterLink
              class="card-btn"
              :to="`/worker/scan?qr=${encodeURIComponent(item.qr_code || item.qrCode || '')}`"
            >
              Открыть
            </RouterLink>

            <RouterLink
              class="soft-btn"
              :to="`/worker/orders?order_id=${encodeURIComponent(item.order_id || item.orderId || '')}`"
            >
              Заявка
            </RouterLink>

            <button type="button" class="download-btn" @click="downloadOneQr(item)">
              QR PNG
            </button>
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

.page-head span {
  display: block;
  margin-top: 14px;
  max-width: 780px;
  color: #5d6d83;
  font-size: 18px;
  line-height: 1.55;
  font-weight: 750;
}

.head-actions {
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  gap: 14px;
  flex-wrap: wrap;
}

.red-btn,
.dark-btn,
.light-btn,
.card-btn,
.soft-btn,
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

.red-btn {
  background: #ff3f4d;
  color: #fff;
  box-shadow: 0 18px 42px rgba(255, 63, 77, .24);
}

.dark-btn,
.card-btn {
  background: #061126;
  color: #fff;
}

.light-btn,
.soft-btn,
.download-btn {
  background: #eef3f9;
  color: #061126;
}

.light-btn:disabled {
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

.panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 18px;
  margin-bottom: 18px;
}

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
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
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

.card-actions {
  margin-top: auto;
  padding-top: 22px;
  flex-wrap: wrap;
}

.qr-mark {
  color: #ff3f4d;
  font-size: 13px;
  font-weight: 950;
  letter-spacing: .22em;
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
  color: #5d6d83;
  font-weight: 900;
  overflow-wrap: anywhere;
}

@media (max-width: 1180px) {
  .filters-panel,
  .cards-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 760px) {
  .page-head,
  .head-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .mini-stats {
    grid-template-columns: 1fr;
  }

  .red-btn,
  .dark-btn,
  .light-btn,
  .card-btn,
  .soft-btn,
  .download-btn {
    width: 100%;
  }
}
</style>
