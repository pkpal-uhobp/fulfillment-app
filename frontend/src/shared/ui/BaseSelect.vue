<template>
  <div ref="root" class="ft-select">
    <label v-if="label" class="ft-label">{{ label }}</label>
    <button ref="button" type="button" class="ft-trigger" :class="{ open, error: error }" :disabled="disabled" @click="toggle">
      <span :class="{ placeholder: !selected }">{{ selected ? optLabel(selected) : placeholder }}</span>
      <span class="chevron" :class="{ rotate: open }">⌄</span>
    </button>
    <p v-if="error" class="ft-error">{{ error }}</p>

    <Teleport to="body">
      <div v-show="open" ref="dropdown" class="ft-menu" :style="menuStyle" @wheel.stop @touchmove.stop>
        <input v-if="searchable" v-model="q" class="ft-search" placeholder="Поиск" />
        <button
          v-for="option in filtered"
          :key="String(optValue(option))"
          type="button"
          class="ft-option"
          :class="{ active: String(optValue(option)) === String(modelValue), disabled: option?.disabled }"
          :disabled="option?.disabled"
          @click="choose(option)"
        >
          <span>
            <strong>{{ optLabel(option) }}</strong>
            <small v-if="optDesc(option)">{{ optDesc(option) }}</small>
          </span>
          <b v-if="String(optValue(option)) === String(modelValue)">✓</b>
        </button>
        <div v-if="!filtered.length" class="ft-empty">{{ emptyText }}</div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, ref } from 'vue'

const props = defineProps({
  modelValue: { type: [String, Number, Boolean, null], default: null },
  options: { type: Array, default: () => [] },
  label: { type: String, default: '' },
  placeholder: { type: String, default: 'Выберите значение' },
  optionLabel: { type: String, default: 'label' },
  optionValue: { type: String, default: 'value' },
  optionDescription: { type: String, default: 'description' },
  error: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
  searchable: { type: Boolean, default: false },
  emptyText: { type: String, default: 'Нет значений' },
})
const emit = defineEmits(['update:modelValue', 'change'])

const root = ref(null)
const button = ref(null)
const dropdown = ref(null)
const open = ref(false)
const q = ref('')
const menuStyle = ref({})
let oldOverflow = ''
let oldPadding = ''

const selected = computed(() => props.options.find((o) => String(optValue(o)) === String(props.modelValue)))
const filtered = computed(() => {
  const term = q.value.trim().toLowerCase()
  if (!term) return props.options
  return props.options.filter((o) => `${optLabel(o)} ${optDesc(o)}`.toLowerCase().includes(term))
})

function optLabel(o) { return typeof o === 'object' ? String(o?.[props.optionLabel] ?? o?.name ?? o?.title ?? o?.value ?? '') : String(o ?? '') }
function optValue(o) { return typeof o === 'object' ? (o?.[props.optionValue] ?? o?.id ?? o?.value) : o }
function optDesc(o) { return typeof o === 'object' ? String(o?.[props.optionDescription] ?? o?.description ?? o?.subtitle ?? '') : '' }

function lock() {
  oldOverflow = document.body.style.overflow
  oldPadding = document.body.style.paddingRight
  const w = window.innerWidth - document.documentElement.clientWidth
  document.body.style.overflow = 'hidden'
  if (w > 0) document.body.style.paddingRight = `${w}px`
}
function unlock() { document.body.style.overflow = oldOverflow; document.body.style.paddingRight = oldPadding }
function position() {
  const r = button.value?.getBoundingClientRect()
  if (!r) return
  const h = Math.min(360, Math.max(90, props.options.length * 70 + (props.searchable ? 64 : 0)))
  const up = r.bottom + h + 12 > window.innerHeight
  menuStyle.value = {
    left: `${Math.max(12, Math.min(r.left, window.innerWidth - r.width - 12))}px`,
    top: `${up ? Math.max(12, r.top - h - 10) : r.bottom + 10}px`,
    width: `${r.width}px`,
    maxHeight: `${Math.min(360, up ? r.top - 24 : window.innerHeight - r.bottom - 24)}px`,
  }
}
async function show() {
  if (props.disabled || open.value) return
  open.value = true; q.value = ''; lock(); await nextTick(); position()
  window.addEventListener('resize', position)
  window.addEventListener('click', outside, true)
  window.addEventListener('keydown', esc)
}
function hide() {
  if (!open.value) return
  open.value = false; unlock()
  window.removeEventListener('resize', position)
  window.removeEventListener('click', outside, true)
  window.removeEventListener('keydown', esc)
}
function toggle() { open.value ? hide() : show() }
function outside(e) { if (!root.value?.contains(e.target) && !dropdown.value?.contains(e.target)) hide() }
function esc(e) { if (e.key === 'Escape') hide() }
function choose(o) { if (o?.disabled) return; emit('update:modelValue', optValue(o)); emit('change', o); hide() }
onBeforeUnmount(hide)
</script>

<style scoped>
.ft-label{display:block;margin:0 0 10px;color:#94a3b8;font-size:12px;font-weight:950;letter-spacing:.28em;text-transform:uppercase}.ft-trigger{width:100%;min-height:64px;display:flex;align-items:center;justify-content:space-between;gap:14px;padding:0 20px;border:1px solid #dce5f0;border-radius:20px;background:#f7faff;color:#07101f;font:inherit;font-weight:950;text-align:left;cursor:pointer;transition:.18s}.ft-trigger:hover,.ft-trigger.open{border-color:#ff3f4c;background:white;box-shadow:0 18px 42px rgba(255,63,76,.14)}.ft-trigger.error{border-color:#ff3f4c;background:#fff5f6}.placeholder{color:#8da0ba}.chevron{font-size:28px;line-height:1;color:#65758b;transition:.18s}.chevron.rotate{transform:rotate(180deg)}.ft-error{margin:9px 0 0;color:#ff3f4c;font-weight:850;font-size:13px}.ft-menu{position:fixed;z-index:9999;overflow:auto;padding:10px;border:1px solid rgba(220,229,240,.95);border-radius:24px;background:rgba(255,255,255,.98);box-shadow:0 34px 80px rgba(7,16,31,.24);backdrop-filter:blur(16px)}.ft-menu::-webkit-scrollbar{width:8px}.ft-menu::-webkit-scrollbar-thumb{background:#c8d3e0;border-radius:999px}.ft-search{width:100%;height:50px;margin-bottom:8px;padding:0 14px;border:0;border-radius:16px;background:#f3f6fa;color:#07101f;font:inherit;font-weight:850;outline:0}.ft-option{width:100%;min-height:58px;display:flex;align-items:center;justify-content:space-between;gap:14px;padding:13px 16px;border:0;border-radius:18px;background:transparent;color:#07101f;font:inherit;text-align:left;cursor:pointer}.ft-option:hover{background:#f1f5fb}.ft-option.active{background:#ff3f4c;color:white}.ft-option.disabled{opacity:.55;cursor:not-allowed}.ft-option strong{display:block;font-weight:950;line-height:1.2}.ft-option small{display:block;margin-top:6px;color:#6b7d95;font-size:13px;font-weight:800}.ft-option.active small{color:rgba(255,255,255,.82)}.ft-empty{padding:18px;color:#7b8da7;font-weight:850}
</style>
