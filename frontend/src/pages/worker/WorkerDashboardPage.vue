<script setup>
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { apiFetch } from '@/shared/api/http'
import { cargoTitle, formatDateTime, gateTitle, normalizeCollection, statusLabel, statusTone, zoneTitle } from './workerUtils'

const cargoItems = ref([])
const loading = ref(false)
const error = ref('')

const stats = computed(() => {
  const items = cargoItems.value
  return {
    total: items.length,
    receive: items.filter((item) => item.status === 'accepted').length,
    storage: items.filter((item) => item.status === 'stored').length,
    shipment: items.filter((item) => ['ready_to_ship', 'shipped'].includes(item.status)).length,
  }
})

const latest = computed(() => cargoItems.value.slice(0, 8))

async function loadData() {
  loading.value = true
  error.value = ''
  try {
    const payload = await apiFetch('/cargo-items', { auth: true })
    cargoItems.value = normalizeCollection(payload)
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить рабочую сводку'
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<template>
  <section class="worker-page">
    <header class="worker-hero">
      <div>
        <p class="eyebrow">Складская смена</p>
        <h1>Рабочая панель</h1>
        <span>Быстрый доступ к QR-проверке, приемке, размещению и подготовке товара к отгрузке.</span>
      </div>
      <div class="hero-actions">
        <RouterLink class="red-btn" to="/worker/scan">Проверить QR</RouterLink>
        <button type="button" class="ghost-btn" :disabled="loading" @click="loadData">{{ loading ? 'Загрузка…' : 'Обновить' }}</button>
      </div>
    </header>

    <div v-if="error" class="alert error">{{ error }}</div>

    <section class="quick-grid">
      <article class="metric-card main">
        <small>Всего грузовых мест</small>
        <strong>{{ stats.total }}</strong>
        <span>В рабочем списке склада</span>
      </article>
      <RouterLink class="metric-card" to="/worker/cargo-items?status=accepted">
        <small>Приемка</small><strong>{{ stats.receive }}</strong><span>Грузы, которые нужно проверить и принять.</span>
      </RouterLink>
      <RouterLink class="metric-card" to="/worker/cargo-items?status=stored">
        <small>Хранение</small><strong>{{ stats.storage }}</strong><span>Товары, размещенные в зонах хранения.</span>
      </RouterLink>
      <RouterLink class="metric-card" to="/worker/cargo-items?status=ready_to_ship">
        <small>Отгрузка</small><strong>{{ stats.shipment }}</strong><span>Грузы для переноса к гейту и отправки.</span>
      </RouterLink>
    </section>

    <section class="operation-strip">
      <div>
        <p class="eyebrow">Основное действие</p>
        <h2>Сканируйте QR и фиксируйте фактический статус</h2>
      </div>
      <RouterLink class="dark-btn" to="/worker/scan">Открыть сканер</RouterLink>
    </section>

    <section class="panel">
      <div class="panel-head">
        <div><p class="eyebrow">Последние грузовые места</p><h2>Текущая очередь</h2></div>
        <RouterLink class="ghost-link" to="/worker/cargo-items">Все места</RouterLink>
      </div>
      <div v-if="!latest.length" class="empty">Пока нет грузовых мест для отображения.</div>
      <div v-else class="cargo-list">
        <article v-for="item in latest" :key="item.id" class="cargo-row">
          <div class="cargo-main">
            <strong>{{ cargoTitle(item) }}</strong>
            <span>Заявка #{{ item.order_id || '—' }} · {{ formatDateTime(item.updated_at || item.created_at) }}</span>
          </div>
          <div class="cargo-meta">
            <span>{{ zoneTitle(item) }}</span>
            <span>{{ gateTitle(item) }}</span>
          </div>
          <em :class="statusTone(item.status)">{{ statusLabel(item.status) }}</em>
        </article>
      </div>
    </section>
  </section>
</template>

<style scoped>
.worker-page { display: grid; gap: 26px; color: #061126; }
.worker-hero, .panel, .operation-strip, .metric-card { background: #fff; box-shadow: 0 18px 62px rgba(15, 23, 42, .08); }
.worker-hero { border-radius: 34px; padding: 34px; display: flex; align-items: flex-start; justify-content: space-between; gap: 24px; overflow: hidden; position: relative; }
.worker-hero::after { content: ''; position: absolute; right: -120px; top: -120px; width: 300px; height: 300px; border-radius: 999px; background: rgba(255, 63, 77, .11); }
.eyebrow { margin: 0 0 10px; color: #ff3f4d; font-size: 13px; font-weight: 950; letter-spacing: .28em; text-transform: uppercase; }
h1 { margin: 0; max-width: 780px; font-size: clamp(46px, 7vw, 88px); line-height: .88; font-weight: 950; letter-spacing: -.06em; }
h2 { margin: 0; font-size: clamp(26px, 3vw, 42px); line-height: 1; font-weight: 950; letter-spacing: -.04em; }
.worker-hero span { display: block; margin-top: 18px; max-width: 720px; color: #5d6d83; font-size: 20px; line-height: 1.55; }
.hero-actions { display: flex; gap: 12px; align-items: center; flex-wrap: wrap; position: relative; z-index: 1; }
.red-btn, .dark-btn, .ghost-btn, .ghost-link { min-height: 60px; border: 0; border-radius: 20px; padding: 0 24px; display: inline-flex; align-items: center; justify-content: center; text-decoration: none; font-size: 17px; font-weight: 950; cursor: pointer; white-space: nowrap; }
.red-btn { background: #ff3f4d; color: #fff; box-shadow: 0 18px 42px rgba(255, 63, 77, .26); }
.dark-btn { background: #061126; color: #fff; }
.ghost-btn, .ghost-link { background: #eef3f9; color: #10223a; }
.quick-grid { display: grid; grid-template-columns: 1.15fr repeat(3, 1fr); gap: 18px; }
.metric-card { min-height: 176px; border-radius: 30px; padding: 26px; color: #061126; text-decoration: none; display: grid; align-content: space-between; transition: transform .18s ease, box-shadow .18s ease; }
.metric-card:hover { transform: translateY(-3px); box-shadow: 0 24px 70px rgba(15, 23, 42, .12); }
.metric-card small { color: #94a3b8; font-size: 13px; font-weight: 950; letter-spacing: .22em; text-transform: uppercase; }
.metric-card strong { display: block; margin-top: 12px; font-size: 56px; line-height: 1; font-weight: 950; letter-spacing: -.06em; }
.metric-card span { margin-top: 14px; color: #5d6d83; font-weight: 800; line-height: 1.45; }
.metric-card.main { background: #061126; color: #fff; }
.metric-card.main small, .metric-card.main span { color: #a9b8ca; }
.operation-strip { border-radius: 30px; padding: 28px; display: flex; align-items: center; justify-content: space-between; gap: 24px; background: linear-gradient(135deg, #071222, #123247); color: #fff; }
.operation-strip .eyebrow { color: #ff9ca5; }
.operation-strip h2 { max-width: 780px; }
.panel { border-radius: 34px; padding: 30px; }
.panel-head { display: flex; align-items: center; justify-content: space-between; gap: 18px; margin-bottom: 18px; }
.cargo-list { display: grid; gap: 12px; }
.cargo-row { min-height: 82px; border-radius: 24px; background: #f6f9fd; padding: 16px 18px; display: grid; grid-template-columns: minmax(220px, 1.2fr) minmax(240px, 1fr) auto; gap: 16px; align-items: center; }
.cargo-main strong { display: block; font-size: 19px; font-weight: 950; overflow-wrap: anywhere; }
.cargo-main span, .cargo-meta span { display: block; margin-top: 6px; color: #66758a; font-weight: 800; }
.cargo-meta { display: grid; gap: 4px; }
em { font-style: normal; padding: 10px 14px; border-radius: 999px; font-weight: 950; white-space: nowrap; text-align: center; }
em.green { background: #dcfce7; color: #047857; } em.red { background: #ffe4e6; color: #be123c; } em.blue { background: #dbeafe; color: #1d4ed8; } em.amber { background: #fef3c7; color: #b45309; } em.violet { background: #ede9fe; color: #6d28d9; } em.dark { background: #061126; color: #fff; } em.gray { background: #e2e8f0; color: #475569; }
.alert, .empty { padding: 18px 22px; border-radius: 22px; font-weight: 900; }
.alert.error { background: #fff0f1; color: #be123c; }
.empty { background: #f6f9fd; color: #64748b; }
@media (max-width: 1180px) { .quick-grid { grid-template-columns: repeat(2, 1fr); } .cargo-row { grid-template-columns: 1fr; align-items: start; } }
@media (max-width: 760px) { .worker-hero, .operation-strip, .panel-head { flex-direction: column; align-items: stretch; } .quick-grid { grid-template-columns: 1fr; } .red-btn, .dark-btn, .ghost-btn, .ghost-link { width: 100%; } }
</style>
