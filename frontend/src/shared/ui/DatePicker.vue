<template>
  <div ref="root" class="date-picker">
    <label v-if="label" class="dp-label">{{ label }}</label>
    <button ref="button" type="button" class="dp-trigger" :class="{ open, error: error }" :disabled="disabled" @click="toggle">
      <span :class="{ placeholder: !modelValue }">{{ modelValue ? human(modelValue) : placeholder }}</span>
      <span>📅</span>
    </button>
    <p v-if="error" class="dp-error">{{ error }}</p>

    <Teleport to="body">
      <div v-show="open" ref="panel" class="dp-panel" :style="style" @wheel.stop @touchmove.stop>
        <header>
          <button type="button" @click="prev">‹</button>
          <strong>{{ months[month] }} {{ year }}</strong>
          <button type="button" @click="next">›</button>
        </header>
        <div class="week"><span v-for="d in week" :key="d">{{ d }}</span></div>
        <div class="grid">
          <button v-for="day in days" :key="day.key" type="button" class="day" :class="{ muted: !day.current, today: day.today, selected: day.iso===modelValue, disabled: day.disabled, closed: day.meta?.is_closed }" :disabled="day.disabled" :title="day.title" @click="select(day)">
            <b>{{ day.date.getDate() }}</b>
            <small v-if="day.meta && day.current">{{ day.meta.is_closed ? 'закрыто' : free(day.meta) }}</small>
          </button>
        </div>
        <footer><span><i></i> доступно</span><span><i class="closed-dot"></i> закрыто / лимит</span></footer>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'
const props = defineProps({ modelValue:{type:String,default:''}, label:{type:String,default:''}, placeholder:{type:String,default:'Выберите дату'}, error:{type:String,default:''}, disabled:{type:Boolean,default:false}, min:{type:String,default:''}, max:{type:String,default:''}, availabilityDays:{type:Array,default:()=>[]} })
const emit = defineEmits(['update:modelValue','change'])
const root=ref(null), button=ref(null), panel=ref(null), open=ref(false), style=ref({})
const now=new Date(); const initial=props.modelValue?parse(props.modelValue):now
const year=ref(initial.getFullYear()), month=ref(initial.getMonth())
const months=['Январь','Февраль','Март','Апрель','Май','Июнь','Июль','Август','Сентябрь','Октябрь','Ноябрь','Декабрь']; const week=['Пн','Вт','Ср','Чт','Пт','Сб','Вс']
let oldOverflow='', oldPadding=''
const map=computed(()=>{const m=new Map(); props.availabilityDays.forEach(d=>{const iso=String(d.pickup_date||d.date||d.blocked_date||'').slice(0,10); if(iso)m.set(iso,d)}); return m})
const days=computed(()=>{const first=new Date(year.value,month.value,1); const shift=(first.getDay()+6)%7; const start=new Date(year.value,month.value,1-shift); return Array.from({length:42},(_,i)=>{const date=new Date(start); date.setDate(start.getDate()+i); const iso=toIso(date); const meta=map.value.get(iso); const disabled=isDisabled(iso,meta); return {key:iso+'-'+i,date,iso,meta,disabled,current:date.getMonth()===month.value,today:iso===toIso(now),title:title(meta,disabled)}})})
function parse(v){const [y,m,d]=String(v).slice(0,10).split('-').map(Number); return new Date(y||now.getFullYear(),(m||1)-1,d||1)}
function toIso(d){return `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`}
function human(v){return new Intl.DateTimeFormat('ru-RU',{day:'2-digit',month:'long',year:'numeric'}).format(parse(v))}
function free(d){const max=Number(d.max_orders||0), cur=Number(d.current_orders||0); return max?`${Math.max(0,max-cur)}/${max}`:'без лимита'}
function isDisabled(iso,d){if(props.min&&iso<props.min)return true; if(props.max&&iso>props.max)return true; if(!d)return false; if(d.is_closed||d.closed)return true; const max=Number(d.max_orders||0), cur=Number(d.current_orders||0); return max>0&&cur>=max}
function title(d,disabled){if(!d)return disabled?'Дата недоступна':'Дата доступна'; if(d.is_closed)return d.reason||'Дата закрыта логистом'; const max=Number(d.max_orders||0), cur=Number(d.current_orders||0); if(max>0&&cur>=max)return 'Лимит заявок исчерпан'; return max?`Доступно ${Math.max(0,max-cur)} из ${max}`:'Дата доступна'}
function lock(){oldOverflow=document.body.style.overflow; oldPadding=document.body.style.paddingRight; const w=innerWidth-document.documentElement.clientWidth; document.body.style.overflow='hidden'; if(w>0)document.body.style.paddingRight=`${w}px`}
function unlock(){document.body.style.overflow=oldOverflow; document.body.style.paddingRight=oldPadding}
function position(){const r=button.value?.getBoundingClientRect(); if(!r)return; const width=Math.max(360,r.width), height=480, up=r.bottom+height+12>innerHeight; style.value={left:`${Math.min(Math.max(12,r.left),innerWidth-width-12)}px`,top:`${up?Math.max(12,r.top-height-10):r.bottom+10}px`,width:`${Math.min(width,innerWidth-24)}px`}}
async function show(){if(props.disabled||open.value)return; open.value=true; lock(); await nextTick(); position(); addEventListener('resize',position); addEventListener('click',outside,true); addEventListener('keydown',esc)}
function hide(){if(!open.value)return; open.value=false; unlock(); removeEventListener('resize',position); removeEventListener('click',outside,true); removeEventListener('keydown',esc)}
function toggle(){open.value?hide():show()} function outside(e){if(!root.value?.contains(e.target)&&!panel.value?.contains(e.target))hide()} function esc(e){if(e.key==='Escape')hide()}
function prev(){if(month.value===0){month.value=11;year.value--}else month.value--} function next(){if(month.value===11){month.value=0;year.value++}else month.value++}
function select(d){if(d.disabled)return; emit('update:modelValue',d.iso); emit('change',d); hide()}
watch(()=>props.modelValue,v=>{if(!v)return; const d=parse(v); year.value=d.getFullYear(); month.value=d.getMonth()})
onBeforeUnmount(hide)
</script>

<style scoped>
.dp-label{display:block;margin:0 0 10px;color:#94a3b8;font-size:12px;font-weight:950;letter-spacing:.28em;text-transform:uppercase}.dp-trigger{width:100%;min-height:64px;display:flex;align-items:center;justify-content:space-between;padding:0 20px;border:1px solid #dce5f0;border-radius:20px;background:#f7faff;color:#07101f;font:inherit;font-weight:950;cursor:pointer}.dp-trigger:hover,.dp-trigger.open{border-color:#ff3f4c;background:white;box-shadow:0 18px 42px rgba(255,63,76,.14)}.dp-trigger.error{border-color:#ff3f4c;background:#fff5f6}.placeholder{color:#8da0ba}.dp-error{margin:9px 0 0;color:#ff3f4c;font-weight:850;font-size:13px}.dp-panel{position:fixed;z-index:9999;padding:18px;border:1px solid rgba(220,229,240,.95);border-radius:28px;background:rgba(255,255,255,.98);box-shadow:0 34px 80px rgba(7,16,31,.24);backdrop-filter:blur(16px)}.dp-panel header{display:flex;align-items:center;justify-content:space-between;margin-bottom:16px}.dp-panel header button{width:42px;height:42px;border:0;border-radius:14px;background:#f1f5fb;color:#07101f;font-size:28px;cursor:pointer}.dp-panel header strong{font-size:18px;font-weight:950}.week,.grid{display:grid;grid-template-columns:repeat(7,1fr);gap:8px}.week{margin-bottom:8px;color:#8da0ba;font-size:12px;font-weight:950;text-align:center}.day{min-height:54px;display:grid;place-items:center;gap:2px;border:0;border-radius:16px;background:#f6f9fd;color:#07101f;cursor:pointer}.day:hover:not(:disabled){background:#eaf1fb}.day b{font-weight:950}.day small{font-size:10px;font-weight:800;color:#64748b}.day.muted{opacity:.42}.day.today{outline:2px solid rgba(255,63,76,.3)}.day.selected{background:#ff3f4c;color:white;box-shadow:0 18px 36px rgba(255,63,76,.24)}.day.selected small{color:rgba(255,255,255,.82)}.day.disabled{cursor:not-allowed;color:#98a8bc;background:#f3f5f8}.day.closed{background:#fff0f2;color:#ff3f4c}footer{display:flex;gap:18px;margin-top:16px;color:#64748b;font-size:12px;font-weight:850}footer span{display:flex;align-items:center;gap:8px}footer i{width:10px;height:10px;border-radius:999px;background:#eaf1fb}.closed-dot{background:#fff0f2!important}
</style>
