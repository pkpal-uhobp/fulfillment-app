<template>
  <section class="page">
    <div v-if="error" class="alert">{{ error }}</div>
    <div v-if="success" class="success">{{ success }}</div>

    <div class="create-panel">
      <div>
        <p>Отгрузки</p>
        <h2>Создать отгрузку</h2>
      </div>
      <div class="create-grid">
        <label>
          <span>Склад назначения</span>
          <BaseSelect v-model="form.destinationWarehouseId" :options="destinationWarehouseOptions" placeholder="Выберите склад" />
        </label>
        <label>
          <span>Гейт</span>
          <BaseSelect v-model="form.gateId" :options="gateOptions" placeholder="Выберите гейт" />
        </label>
        <label>
          <span>Дата и время отправки</span>
          <input v-model="form.plannedDepartureAt" type="datetime-local" />
        </label>
        <button type="button" @click="createShipment">Создать</button>
      </div>
    </div>

    <div class="toolbar panel">
      <div>
        <p>Список</p>
        <h2>Партии к отправке</h2>
      </div>
      <div class="filters">
        <BaseSelect v-model="filters.status" :options="shipmentStatusOptions" placeholder="Все статусы" @change="loadShipments" />
        <button type="button" @click="loadShipments">Обновить</button>
      </div>
    </div>

    <div class="shipment-grid">
      <article v-for="shipment in shipments" :key="shipment.id" class="shipment-card">
        <div class="card-top">
          <div>
            <span>Отгрузка #{{ shipment.id }}</span>
            <h3>{{ labelFromMap(shipmentStatusLabels, shipment.status) }}</h3>
          </div>
          <RouterLink :to="`/logist/shipments/${shipment.id}`">Открыть</RouterLink>
        </div>
        <div class="meta-grid">
          <div><small>Назначение</small><b>{{ warehouseName(shipment.destination_warehouse_id, warehouses) }}</b></div>
          <div><small>Гейт</small><b>{{ gateName(shipment.gate_id, gates) }}</b></div>
          <div><small>План отправки</small><b>{{ formatDateTime(shipment.planned_departure_at) }}</b></div>
          <div><small>Мест</small><b>{{ shipment.items?.length || 0 }}</b></div>
        </div>
        <div class="action-box">
          <BaseSelect
            v-model="forms[shipment.id].status"
            :options="shipmentStatusOptions.filter((item) => item.value)"
            placeholder="Выберите статус"
          />
          <button type="button" :disabled="loadingId === shipment.id" @click="updateStatus(shipment)">Обновить статус</button>
        </div>
      </article>
    </div>

    <div v-if="!shipments.length && !error" class="empty">Отгрузки не найдены</div>
  </section>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'
import BaseSelect from '@/shared/ui/BaseSelect.vue'
import { apiFetch } from '@/shared/api/http'
import { formatDateTime, gateName, labelFromMap, shipmentStatusLabels, shipmentStatusOptions, unwrapList, unwrapOne, warehouseName } from './logistUtils'

const shipments = ref([])
const warehouses = ref([])
const gates = ref([])
const forms = reactive({})
const filters = reactive({ status: '' })
const form = reactive({ destinationWarehouseId: '', gateId: '', plannedDepartureAt: '' })
const error = ref('')
const success = ref('')
const loadingId = ref(null)

const destinationWarehouses = computed(() => warehouses.value.filter((warehouse) => ['destination', 'both'].includes(warehouse.warehouse_type)))
const destinationWarehouseOptions = computed(() => destinationWarehouses.value.map((warehouse) => ({
  value: String(warehouse.id),
  label: `${warehouseName(warehouse.id, warehouses.value)} · ${warehouse.city}`,
  description: warehouse.address,
})))
const gateOptions = computed(() => gates.value.map((gate) => ({ value: String(gate.id), label: gate.name, description: gate.warehouse_name || '' })))

function ensureForm(shipment) {
  if (!forms[shipment.id]) forms[shipment.id] = { status: shipment.status || 'planned' }
}

async function loadCatalogs() {
  const [warehousesPayload, gatesPayload] = await Promise.all([
    apiFetch('/warehouses'),
    apiFetch('/gates', { auth: true }),
  ])
  warehouses.value = unwrapList(warehousesPayload, 'warehouses')
  gates.value = unwrapList(gatesPayload, 'gates')
}

async function loadShipments() {
  error.value = ''
  success.value = ''
  const params = new URLSearchParams()
  if (filters.status) params.set('status', filters.status)
  try {
    const payload = await apiFetch(`/shipments${params.toString() ? `?${params}` : ''}`, { auth: true })
    shipments.value = unwrapList(payload, 'shipments')
    shipments.value.forEach(ensureForm)
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить отгрузки'
  }
}

async function createShipment() {
  error.value = ''
  success.value = ''
  if (!form.destinationWarehouseId || !form.gateId || !form.plannedDepartureAt) {
    error.value = 'Выберите склад, гейт и дату отправки'
    return
  }
  try {
    await apiFetch('/shipments', {
      auth: true,
      method: 'POST',
      body: {
        destination_warehouse_id: Number(form.destinationWarehouseId),
        gate_id: Number(form.gateId),
        planned_departure_at: new Date(form.plannedDepartureAt).toISOString(),
      },
    })
    success.value = 'Отгрузка создана'
    form.destinationWarehouseId = ''
    form.gateId = ''
    form.plannedDepartureAt = ''
    await loadShipments()
  } catch (err) {
    error.value = err.message || 'Не удалось создать отгрузку'
  }
}

async function updateStatus(shipment) {
  error.value = ''
  success.value = ''
  loadingId.value = shipment.id
  try {
    const payload = await apiFetch(`/shipments/${shipment.id}/status`, {
      auth: true,
      method: 'PATCH',
      body: { status: forms[shipment.id].status },
    })
    const updated = unwrapOne(payload, 'shipment') || { ...shipment, status: forms[shipment.id].status }
    const index = shipments.value.findIndex((item) => item.id === shipment.id)
    if (index !== -1) shipments.value[index] = updated
    ensureForm(updated)
    success.value = `Статус отгрузки #${shipment.id} обновлён`
  } catch (err) {
    error.value = err.message || 'Не удалось обновить статус отгрузки'
  } finally {
    loadingId.value = null
  }
}

onMounted(async () => {
  await Promise.all([loadCatalogs(), loadShipments()])
})
</script>

<style scoped>
.page{display:grid;gap:20px}.alert,.success{padding:16px 18px;border-radius:18px;font-weight:950}.alert{background:#fee2e2;color:#991b1b}.success{background:#d1fae5;color:#065f46}.panel,.create-panel,.shipment-card,.empty{background:white;border-radius:32px;padding:24px;box-shadow:0 18px 42px rgba(7,16,31,.08)}.create-panel{display:grid;gap:18px}p{margin:0 0 8px;color:#ff3f4d;letter-spacing:.22em;text-transform:uppercase;font-size:12px;font-weight:950}h2,h3{margin:0;letter-spacing:-.04em}.toolbar{display:flex;justify-content:space-between;gap:18px;align-items:end}.toolbar h2,.create-panel h2{font-size:34px}.create-grid{display:grid;grid-template-columns:1.4fr 1fr 1.2fr auto;gap:14px;align-items:end}.create-grid label{display:grid;gap:8px}.create-grid label>span{color:#94a3b8;letter-spacing:.22em;text-transform:uppercase;font-size:12px;font-weight:950}.filters{display:grid;grid-template-columns:minmax(200px,280px) auto;gap:12px}input{min-height:58px;border:1px solid #dbe3ef;border-radius:18px;padding:0 18px;background:#f6f8fb;color:#07101f;font-weight:950;font-family:inherit}button,.card-top a{min-height:58px;border:0;border-radius:18px;padding:0 20px;background:#ff3f4d;color:white;font-weight:950;cursor:pointer;font-family:inherit;text-decoration:none;display:inline-flex;align-items:center;justify-content:center}.filters button{background:#07101f}.shipment-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:20px}.shipment-card{display:grid;gap:18px}.card-top{display:flex;justify-content:space-between;gap:14px}.card-top span{color:#94a3b8;letter-spacing:.22em;text-transform:uppercase;font-size:12px;font-weight:950}.card-top h3{margin-top:8px;font-size:28px}.meta-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:12px}.meta-grid div,.action-box{background:#f6f8fb;border-radius:22px;padding:16px}.meta-grid small{display:block;color:#64748b;font-weight:900;margin-bottom:6px}.action-box{display:grid;grid-template-columns:1fr auto;gap:12px}.action-box button{min-height:58px}.empty{text-align:center;color:#64748b;font-weight:950}@media(max-width:1180px){.shipment-grid{grid-template-columns:1fr}.create-grid,.toolbar{grid-template-columns:1fr;display:grid}.filters{grid-template-columns:1fr}}@media(max-width:680px){.meta-grid,.action-box{grid-template-columns:1fr}}
</style>
