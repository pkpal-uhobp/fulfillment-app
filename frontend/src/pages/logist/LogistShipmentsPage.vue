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
          Склад назначения
          <select v-model="form.destinationWarehouseId">
            <option value="">Выберите склад</option>
            <option v-for="warehouse in destinationWarehouses" :key="warehouse.id" :value="warehouse.id">
              {{ warehouseName(warehouse.id, warehouses) }} · {{ warehouse.city }}
            </option>
          </select>
        </label>
        <label>
          Гейт
          <select v-model="form.gateId">
            <option value="">Выберите гейт</option>
            <option v-for="gate in gates" :key="gate.id" :value="gate.id">{{ gate.name }}</option>
          </select>
        </label>
        <label>
          Дата и время отправки
          <input v-model="form.plannedDepartureAt" type="datetime-local" />
        </label>
        <button type="button" @click="createShipment">Создать</button>
      </div>
    </div>

    <div class="toolbar">
      <div>
        <p>Список</p>
        <h2>Партии к отправке</h2>
      </div>
      <div class="filters">
        <select v-model="filters.status" @change="loadShipments">
          <option v-for="option in shipmentStatusOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
        </select>
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
          <select v-model="forms[shipment.id].status">
            <option v-for="option in shipmentStatusOptions.filter((item) => item.value)" :key="option.value" :value="option.value">{{ option.label }}</option>
          </select>
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
    const updated = unwrapOne(payload, 'shipment')
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
.page { display:grid; gap:20px; }
.alert,.success{padding:16px 18px;border-radius:18px;font-weight:900}.alert{background:#fee2e2;color:#991b1b}.success{background:#d1fae5;color:#065f46}
.create-panel,.toolbar,.shipment-card{background:white;border-radius:34px;padding:26px;box-shadow:0 18px 42px rgba(7,16,31,.08)}
.create-panel,.toolbar{display:grid;gap:20px}.toolbar{display:flex;justify-content:space-between;align-items:end}.create-panel p,.toolbar p,.card-top span{margin:0 0 8px;color:#ff3f4d;text-transform:uppercase;letter-spacing:.22em;font-size:12px;font-weight:900}
h2,h3{margin:0;letter-spacing:-.04em}h2{font-size:34px}h3{font-size:28px}
.create-grid{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:12px;align-items:end}.filters,.action-box{display:flex;gap:12px;flex-wrap:wrap}
label{display:grid;gap:8px;color:#8b9ab0;font-size:12px;text-transform:uppercase;letter-spacing:.14em;font-weight:900}
select,input{height:52px;border:1px solid #dbe3ef;border-radius:18px;padding:0 16px;background:#f6f8fb;color:#07101f;font-weight:800}
button,.card-top a{height:52px;border:0;border-radius:18px;padding:0 18px;background:#ff3f4d;color:white;font-weight:900;cursor:pointer;text-decoration:none;display:inline-flex;align-items:center;justify-content:center}.card-top a{height:44px;background:#07101f}
.shipment-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:20px}.shipment-card{display:grid;gap:18px}.card-top{display:flex;justify-content:space-between;gap:14px;align-items:start}
.meta-grid{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:12px}.meta-grid div{padding:16px;border-radius:20px;background:#f6f8fb;display:grid;gap:5px}small{color:#6b7b91;font-weight:800}.empty{padding:26px;border-radius:26px;background:white;font-weight:900;color:#63738a}
@media(max-width:1100px){.create-grid,.shipment-grid{grid-template-columns:1fr}.toolbar{flex-direction:column;align-items:stretch}}@media(max-width:640px){.meta-grid{grid-template-columns:1fr}.card-top{flex-direction:column}}
</style>
