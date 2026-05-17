<template>
  <section class="page">
    <div v-if="error" class="alert">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <article v-if="shipment" class="hero-card">
      <div>
        <p>Отгрузка</p>
        <h2>#{{ shipment.id }} · {{ labelFromMap(shipmentStatusLabels, shipment.status) }}</h2>
        <span>{{ warehouseName(shipment.destination_warehouse_id, warehouses) }} · {{ gateName(shipment.gate_id, gates) }}</span>
      </div>
      <RouterLink to="/logist/shipments">Назад к списку</RouterLink>
    </article>

    <div class="grid">
      <article class="panel">
        <p>Добавить груз</p>
        <h3>QR или ID места</h3>
        <div class="inline-form">
          <input v-model.trim="cargoQuery" placeholder="QR-TPRO-MSK-240001 или 12" @keyup.enter="addItem" />
          <button type="button" @click="addItem">Добавить</button>
        </div>
        <small>Если введён QR-код, сначала будет выполнен поиск через сканирование.</small>
      </article>

      <article class="panel">
        <p>Статус</p>
        <h3>Управление отправкой</h3>
        <div class="inline-form">
          <BaseSelect
            v-model="statusForm.status"
            :options="shipmentStatusOptions.filter((item) => item.value)"
            placeholder="Выберите статус"
          />
          <button type="button" @click="updateStatus">Обновить</button>
        </div>
      </article>
    </div>

    <article class="panel wide">
      <div class="head-row">
        <div>
          <p>Состав</p>
          <h3>Грузовые места в отгрузке</h3>
        </div>
        <b>{{ shipment?.items?.length || 0 }} мест</b>
      </div>

      <div class="items-list">
        <div v-for="item in shipment?.items || []" :key="item.cargo_item_id" class="item-row">
          <div>
            <strong>{{ item.qr_code }}</strong>
            <span>{{ labelFromMap(cargoStatusLabels, item.status) }} · заявка #{{ item.order_id }}</span>
          </div>
          <button type="button" class="ghost" @click="removeItem(item.cargo_item_id)">Убрать</button>
        </div>
        <div v-if="!shipment?.items?.length" class="empty">В отгрузке пока нет грузовых мест</div>
      </div>
    </article>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import BaseSelect from '@/shared/ui/BaseSelect.vue'
import { apiFetch } from '@/shared/api/http'
import { cargoStatusLabels, gateName, labelFromMap, shipmentStatusLabels, shipmentStatusOptions, unwrapList, unwrapOne, warehouseName } from './logistUtils'

const route = useRoute()
const shipment = ref(null)
const warehouses = ref([])
const gates = ref([])
const cargoQuery = ref('')
const statusForm = reactive({ status: 'planned' })
const error = ref('')
const success = ref('')

async function loadCatalogs() {
  const [warehousesPayload, gatesPayload] = await Promise.all([
    apiFetch('/warehouses'),
    apiFetch('/gates', { auth: true }),
  ])
  warehouses.value = unwrapList(warehousesPayload, 'warehouses')
  gates.value = unwrapList(gatesPayload, 'gates')
}

async function loadShipment() {
  error.value = ''
  try {
    const payload = await apiFetch(`/shipments/${route.params.id}`, { auth: true })
    shipment.value = unwrapOne(payload, 'shipment')
    statusForm.status = shipment.value?.status || 'planned'
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить отгрузку'
  }
}

async function resolveCargoId() {
  const value = cargoQuery.value.trim()
  if (!value) throw new Error('Введите QR-код или ID грузового места')
  if (/^\d+$/.test(value)) return Number(value)
  const payload = await apiFetch(`/cargo-items/scan?qr_code=${encodeURIComponent(value)}`, { auth: true })
  const cargo = unwrapOne(payload, 'cargo_item')
  if (!cargo?.id) throw new Error('Грузовое место не найдено')
  return Number(cargo.id)
}

async function addItem() {
  error.value = ''
  success.value = ''
  try {
    const cargoItemId = await resolveCargoId()
    const payload = await apiFetch(`/shipments/${route.params.id}/items`, {
      auth: true,
      method: 'POST',
      body: { cargo_item_id: cargoItemId },
    })
    shipment.value = unwrapOne(payload, 'shipment')
    cargoQuery.value = ''
    success.value = 'Грузовое место добавлено в отгрузку'
  } catch (err) {
    error.value = err.message || 'Не удалось добавить грузовое место'
  }
}

async function removeItem(cargoItemId) {
  error.value = ''
  success.value = ''
  try {
    await apiFetch(`/shipments/${route.params.id}/items/${cargoItemId}`, { auth: true, method: 'DELETE' })
    success.value = 'Грузовое место убрано из отгрузки'
    await loadShipment()
  } catch (err) {
    error.value = err.message || 'Не удалось убрать грузовое место'
  }
}

async function updateStatus() {
  error.value = ''
  success.value = ''
  try {
    const payload = await apiFetch(`/shipments/${route.params.id}/status`, {
      auth: true,
      method: 'PATCH',
      body: { status: statusForm.status },
    })
    shipment.value = unwrapOne(payload, 'shipment')
    success.value = 'Статус отгрузки обновлён'
  } catch (err) {
    error.value = err.message || 'Не удалось обновить статус'
  }
}

onMounted(async () => {
  await Promise.all([loadCatalogs(), loadShipment()])
})
</script>

<style scoped>
.page{display:grid;gap:20px}.alert,.success{padding:16px 18px;border-radius:18px;font-weight:950}.alert{background:#fee2e2;color:#991b1b}.success{background:#d1fae5;color:#065f46}.hero-card,.panel{background:white;border-radius:34px;padding:26px;box-shadow:0 18px 42px rgba(7,16,31,.08)}.hero-card{display:flex;justify-content:space-between;gap:18px;align-items:center;background:#07101f;color:white}p{margin:0 0 8px;color:#ff9da5;text-transform:uppercase;letter-spacing:.22em;font-size:12px;font-weight:950}h2,h3{margin:0;letter-spacing:-.04em}h2{font-size:36px}h3{font-size:28px}.hero-card span{color:#b8c3d6;font-weight:800}.hero-card a{background:white;color:#07101f;border-radius:18px;padding:16px 20px;text-decoration:none;font-weight:950;white-space:nowrap}.grid{display:grid;grid-template-columns:1fr 1fr;gap:20px}.inline-form{display:grid;grid-template-columns:minmax(0,1fr) auto;gap:12px;margin-top:16px}input{height:58px;min-width:0;border:1px solid #dbe3ef;border-radius:18px;padding:0 16px;background:#f6f8fb;color:#07101f;font-weight:900;font-family:inherit}button{height:58px;border:0;border-radius:18px;padding:0 18px;background:#ff3f4d;color:white;font-weight:950;cursor:pointer;font-family:inherit}.panel small{display:block;margin-top:12px;color:#64748b;font-weight:800}.head-row,.item-row{display:flex;justify-content:space-between;gap:16px;align-items:center}.head-row{margin-bottom:18px}.head-row b{padding:10px 14px;border-radius:999px;background:#ffe6e8;color:#ff3f4d}.items-list{display:grid;gap:12px}.item-row{padding:18px;border-radius:22px;background:#f6f8fb}.item-row strong{display:block;font-size:19px}.item-row span{color:#64748b;font-weight:800}.ghost{background:#edf2f7;color:#07101f}.empty{padding:22px;border-radius:22px;background:#f6f8fb;color:#64748b;font-weight:900}@media(max-width:900px){.grid{grid-template-columns:1fr}.hero-card,.head-row,.item-row{display:grid}.inline-form{grid-template-columns:1fr}}
</style>
