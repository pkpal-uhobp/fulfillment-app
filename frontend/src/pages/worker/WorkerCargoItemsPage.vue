<template>
  <section class="page">
    <div class="panel"><div><p>Склад</p><h1>Грузовые места</h1><span>Фильтрация и быстрая смена статуса.</span></div><button @click="load">Обновить</button></div>
    <div class="toolbar"><input v-model="q" placeholder="Поиск по QR или заявке"/><BaseSelect v-model="status" :options="statusOptions" label="Статус"/></div>
    <div class="cards"><article v-for="i in filtered" :key="i.id" class="card"><header><div><p>QR</p><h2>{{ i.qr_code || `Груз #${i.id}` }}</h2></div><span>{{ cargoStatusLabel(i.status) }}</span></header><dl><div><dt>Заявка</dt><dd>#{{ i.order_id || '—' }}</dd></div><div><dt>Зона</dt><dd>#{{ i.storage_zone_id || '—' }}</dd></div><div><dt>Гейт</dt><dd>#{{ i.gate_id || '—' }}</dd></div><div><dt>Создано</dt><dd>{{ formatDateTime(i.created_at) }}</dd></div></dl><div class="op"><BaseSelect v-model="forms[i.id].status" :options="workerCargoStatusOptions" label="Статус"/><textarea v-model="forms[i.id].comment" placeholder="Комментарий"></textarea><button @click="updateStatus(i)">Сохранить</button></div></article><p v-if="!filtered.length" class="empty">Грузы не найдены.</p></div>
  </section>
</template>
<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { apiFetch } from '@/shared/api/http'
import BaseSelect from '@/shared/ui/BaseSelect.vue'
import { cargoStatusLabel, formatDateTime, normalizeCollection, workerCargoStatusOptions } from './workerUtils'
const items=ref([]), q=ref(''), status=ref('all'), forms=reactive({})
const statusOptions=[{value:'all',label:'Все статусы'},...workerCargoStatusOptions]
const filtered=computed(()=>items.value.filter(i=>(status.value==='all'||i.status===status.value)&&(!q.value.trim()||`${i.qr_code||''} ${i.id} ${i.order_id||''}`.toLowerCase().includes(q.value.trim().toLowerCase()))))
function ensure(i){if(!forms[i.id])forms[i.id]={status:i.status||'received',comment:''}}
async function load(){const d=await apiFetch('/cargo-items',{auth:true}); items.value=normalizeCollection(d,['cargo_items']); items.value.forEach(ensure)}
async function updateStatus(i){const f=forms[i.id]; await apiFetch(`/cargo-items/${i.id}/status`,{method:'PATCH',auth:true,body:{status:f.status,comment:f.comment?.trim()||null}}); f.comment=''; await load()}
onMounted(load)
</script>
<style scoped>
.page{display:grid;gap:22px}.panel,.toolbar,.card{border-radius:32px;background:rgba(255,255,255,.96);color:#07101f;box-shadow:0 28px 70px rgba(0,0,0,.12)}.panel{display:flex;justify-content:space-between;gap:24px;padding:34px}.panel p,.card header p{margin:0 0 12px;color:#ff3f4c;font-weight:950;letter-spacing:.32em;text-transform:uppercase}.panel h1{margin:0;font-size:clamp(38px,5vw,64px);line-height:.95;font-weight:950}.panel span{display:block;margin-top:16px;color:#5b6f88;font-size:18px}.panel button,.op button{min-height:62px;border:0;border-radius:20px;background:#07101f;color:#fff;font:inherit;font-weight:950;padding:0 24px;cursor:pointer}.toolbar{display:grid;grid-template-columns:1fr 320px;gap:14px;align-items:end;padding:22px}input,textarea{width:100%;border:1px solid #dce5f0;border-radius:20px;background:#f7faff;color:#07101f;font:inherit;font-weight:900;outline:none}input{min-height:64px;padding:0 20px}textarea{min-height:64px;padding:18px;resize:vertical}.cards{display:grid;gap:18px}.card{padding:24px}.card header{display:flex;justify-content:space-between;gap:18px;margin-bottom:18px}.card h2{margin:0;font-size:30px;font-weight:950}.card header span{padding:12px 16px;border-radius:16px;background:#fff0f2;color:#ff3f4c;font-weight:950}dl{display:grid;grid-template-columns:repeat(4,1fr);gap:14px;margin:0 0 18px}dl div{padding:16px;border-radius:18px;background:#f4f7fb}dt{color:#93a4bb;font-size:12px;font-weight:950;letter-spacing:.22em;text-transform:uppercase}dd{margin:8px 0 0;font-weight:950}.op{display:grid;grid-template-columns:1fr 1fr auto;gap:14px;align-items:end;padding:18px;border-radius:24px;background:#f8fafc}.empty{color:white;font-weight:900}@media(max-width:980px){.toolbar,.op,dl{grid-template-columns:1fr}.panel{flex-direction:column}}
</style>
