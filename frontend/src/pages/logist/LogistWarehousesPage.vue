<template>
  <section class="page">
    <div v-if="error" class="alert">{{ error }}</div>

    <div class="warehouse-grid">
      <article v-for="warehouse in warehouses" :key="warehouse.id" class="warehouse-card">
        <div class="card-top">
          <div>
            <span>{{ typeLabel(warehouse.warehouse_type) }}</span>
            <h3>{{ warehouse.name }}</h3>
          </div>
          <em :class="{ off: !warehouse.is_active }">{{ warehouse.is_active ? 'Активен' : 'Выключен' }}</em>
        </div>
        <p>{{ warehouse.city }} · {{ warehouse.address }}</p>
        <small v-if="warehouse.marketplace">Маркетплейс: {{ warehouse.marketplace }}</small>
      </article>
    </div>

    <div class="catalog-grid">
      <article class="panel">
        <div class="panel-head">
          <p>Зоны хранения</p>
          <button type="button" @click="loadCatalogs">Обновить</button>
        </div>
        <div class="catalog-list">
          <div v-for="zone in zones" :key="zone.id" class="catalog-row">
            <b>{{ zone.name }}</b>
            <span>{{ warehouseName(zone.warehouse_id, warehouses) }}</span>
          </div>
          <div v-if="!zones.length" class="empty">Зон пока нет</div>
        </div>
      </article>

      <article class="panel">
        <div class="panel-head">
          <p>Гейты</p>
          <button type="button" @click="loadCatalogs">Обновить</button>
        </div>
        <div class="catalog-list">
          <div v-for="gate in gates" :key="gate.id" class="catalog-row">
            <b>{{ gate.name }}</b>
            <span>{{ warehouseName(gate.warehouse_id, warehouses) }}</span>
          </div>
          <div v-if="!gates.length" class="empty">Гейтов пока нет</div>
        </div>
      </article>
    </div>
  </section>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { apiFetch } from '@/shared/api/http'
import { unwrapList, warehouseName } from './logistUtils'

const warehouses = ref([])
const zones = ref([])
const gates = ref([])
const error = ref('')

function typeLabel(value) {
  return {
    both: 'Приёмка и назначение',
    receiving: 'Склад приёмки',
    destination: 'Склад назначения',
  }[value] || value || 'Склад'
}

async function loadCatalogs() {
  error.value = ''
  try {
    const [warehousesPayload, zonesPayload, gatesPayload] = await Promise.all([
      apiFetch('/warehouses'),
      apiFetch('/storage-zones', { auth: true }),
      apiFetch('/gates', { auth: true }),
    ])
    warehouses.value = unwrapList(warehousesPayload, 'warehouses')
    zones.value = unwrapList(zonesPayload, 'storage_zones')
    gates.value = unwrapList(gatesPayload, 'gates')
  } catch (err) {
    error.value = err.message || 'Не удалось загрузить складские справочники'
  }
}

onMounted(loadCatalogs)
</script>

<style scoped>
.page{display:grid;gap:20px}.alert{padding:16px 18px;border-radius:18px;font-weight:900;background:#fee2e2;color:#991b1b}.warehouse-grid{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:18px}.warehouse-card,.panel{background:white;border-radius:32px;padding:24px;box-shadow:0 18px 42px rgba(7,16,31,.08)}.card-top,.panel-head,.catalog-row{display:flex;justify-content:space-between;gap:14px;align-items:flex-start}.card-top span,.panel p{color:#ff3f4d;text-transform:uppercase;letter-spacing:.22em;font-size:12px;font-weight:900}.card-top h3{margin:8px 0 0;font-size:28px;letter-spacing:-.04em}.card-top em{font-style:normal;padding:9px 12px;border-radius:999px;background:#dcfce7;color:#047857;font-weight:900}.card-top em.off{background:#fee2e2;color:#991b1b}.warehouse-card p{color:#334155;font-weight:800;line-height:1.5}.warehouse-card small,.catalog-row span{color:#64748b;font-weight:800}.catalog-grid{display:grid;grid-template-columns:1fr 1fr;gap:20px}.panel-head{align-items:center;margin-bottom:16px}.panel p{margin:0}.panel button{height:42px;border:0;border-radius:14px;background:#07101f;color:white;font-weight:900;padding:0 14px;cursor:pointer}.catalog-list{display:grid;gap:10px}.catalog-row{padding:16px;border-radius:20px;background:#f6f8fb}.empty{padding:22px;border-radius:20px;background:#f6f8fb;color:#64748b;font-weight:900}@media(max-width:1100px){.warehouse-grid{grid-template-columns:1fr 1fr}.catalog-grid{grid-template-columns:1fr}}@media(max-width:640px){.warehouse-grid{grid-template-columns:1fr}.catalog-row,.card-top{flex-direction:column}}
</style>
