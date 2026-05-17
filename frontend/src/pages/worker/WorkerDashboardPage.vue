<template>
  <section class="page">
    <div class="hero"><div><p>Складская смена</p><h1>Панель рабочего</h1><span>QR-сканер, приёмка и смена статусов грузовых мест.</span></div><RouterLink to="/worker/scan">Открыть QR-сканер</RouterLink></div>
    <div class="stats"><article v-for="s in stats" :key="s.label"><b>{{ s.value }}</b><span>{{ s.label }}</span></article></div>
    <article class="card"><header><h2>Последние грузы</h2><RouterLink to="/worker/cargo-items">Все грузы</RouterLink></header><div class="rows"><div v-for="i in cargoItems.slice(0,6)" :key="i.id" class="row"><b>{{ i.qr_code || `Груз #${i.id}` }}</b><span>{{ cargoStatusLabel(i.status) }}</span><small>{{ formatDateTime(i.updated_at || i.created_at) }}</small></div><p v-if="!cargoItems.length" class="empty">Пока нет грузовых мест.</p></div></article>
  </section>
</template>
<script setup>
import { computed, onMounted, ref } from 'vue'
import { apiFetch } from '@/shared/api/http'
import { cargoStatusLabel, formatDateTime, normalizeCollection } from './workerUtils'
const cargoItems=ref([])
const stats=computed(()=>[{label:'грузовых мест',value:cargoItems.value.length},{label:'принято',value:cargoItems.value.filter(i=>['accepted','received'].includes(i.status)).length},{label:'на хранении',value:cargoItems.value.filter(i=>i.status==='stored').length},{label:'готово к отгрузке',value:cargoItems.value.filter(i=>i.status==='ready_to_ship').length}])
async function load(){const d=await apiFetch('/cargo-items',{auth:true}); cargoItems.value=normalizeCollection(d,['cargo_items'])}
onMounted(load)
</script>
<style scoped>
.page{display:grid;gap:24px}.hero,.card{border-radius:34px;background:rgba(255,255,255,.96);color:#07101f;box-shadow:0 34px 80px rgba(0,0,0,.16)}.hero{display:flex;justify-content:space-between;gap:24px;padding:36px}.hero p{margin:0 0 14px;color:#ff3f4c;font-weight:950;letter-spacing:.38em;text-transform:uppercase}.hero h1{margin:0;font-size:clamp(44px,7vw,84px);line-height:.9;font-weight:950}.hero span{display:block;max-width:780px;margin-top:22px;color:#596b84;font-size:20px;line-height:1.6}.hero a,.card header a{min-height:58px;display:inline-flex;align-items:center;align-self:flex-end;padding:0 24px;border-radius:20px;background:#ff3f4c;color:#fff;text-decoration:none;font-weight:950}.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:18px}.stats article{padding:24px;border-radius:28px;background:#101b2b;color:white}.stats b{display:block;font-size:42px;font-weight:950}.stats span{color:#b9c6d7;font-weight:800}.card{padding:28px}.card header{display:flex;justify-content:space-between;margin-bottom:20px}.card h2{margin:0;font-size:32px;font-weight:950}.rows{display:grid;gap:12px}.row{display:grid;grid-template-columns:1fr 190px auto;gap:14px;padding:18px;border-radius:20px;background:#f4f7fb}.row span{color:#ff3f4c;font-weight:950}.row small,.empty{color:#64748b;font-weight:850}@media(max-width:980px){.hero{flex-direction:column}.stats{grid-template-columns:repeat(2,1fr)}.row{grid-template-columns:1fr}}
</style>
