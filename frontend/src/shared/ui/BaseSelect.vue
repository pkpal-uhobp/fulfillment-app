<script setup>
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue'

const props = defineProps({
  modelValue: { type: [String, Number, Boolean, null], default: null },
  options: { type: Array, default: () => [] },
  placeholder: { type: String, default: 'Выберите значение' },
  label: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
  error: { type: String, default: '' },
  emptyText: { type: String, default: 'Нет вариантов' },
})

const emit = defineEmits(['update:modelValue', 'change'])
const rootRef = ref(null)
const open = ref(false)
const direction = ref('down')
const dropdownStyle = ref({})

const normalizedOptions = computed(() => props.options.map((item) => {
  if (item && typeof item === 'object') {
    return {
      value: item.value ?? item.id ?? item.key,
      label: item.label ?? item.name ?? item.title ?? String(item.value ?? item.id ?? ''),
      description: item.description ?? item.subtitle ?? '',
      badge: item.badge ?? '',
      disabled: Boolean(item.disabled),
      raw: item,
    }
  }
  return { value: item, label: String(item), description: '', badge: '', disabled: false, raw: item }
}))

const selected = computed(() => normalizedOptions.value.find((item) => String(item.value) === String(props.modelValue)))
const selectedLabel = computed(() => selected.value?.label || props.placeholder)

function lockScroll() {
  document.documentElement.classList.add('ft-select-lock')
  document.body.classList.add('ft-select-lock')
}

function unlockScroll() {
  if (!document.querySelector('.ft-select.is-open')) {
    document.documentElement.classList.remove('ft-select-lock')
    document.body.classList.remove('ft-select-lock')
  }
}

function updateDirection() {
  const rect = rootRef.value?.getBoundingClientRect()
  if (!rect) return
  direction.value = window.innerHeight - rect.bottom < 320 && rect.top > 320 ? 'up' : 'down'
  dropdownStyle.value = {
    width: `${rect.width}px`,
    left: `${rect.left}px`,
    top: direction.value === 'down' ? `${rect.bottom + 10}px` : 'auto',
    bottom: direction.value === 'up' ? `${window.innerHeight - rect.top + 10}px` : 'auto',
  }
}

function toggle() {
  if (props.disabled) return
  open.value = !open.value
}

function close() {
  open.value = false
}

function choose(option) {
  if (option.disabled) return
  emit('update:modelValue', option.value)
  emit('change', option.raw)
  close()
}

function onDocumentClick(event) {
  if (!rootRef.value?.contains(event.target)) close()
}

function onKeydown(event) {
  if (event.key === 'Escape') close()
}

watch(open, async (value) => {
  if (value) {
    await nextTick()
    updateDirection()
    lockScroll()
    document.addEventListener('click', onDocumentClick)
    window.addEventListener('resize', updateDirection)
    window.addEventListener('scroll', updateDirection, true)
    document.addEventListener('keydown', onKeydown)
  } else {
    document.removeEventListener('click', onDocumentClick)
    window.removeEventListener('resize', updateDirection)
    window.removeEventListener('scroll', updateDirection, true)
    document.removeEventListener('keydown', onKeydown)
    nextTick(unlockScroll)
  }
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocumentClick)
  window.removeEventListener('resize', updateDirection)
  window.removeEventListener('scroll', updateDirection, true)
  document.removeEventListener('keydown', onKeydown)
  open.value = false
  unlockScroll()
})
</script>

<template>
  <div ref="rootRef" class="ft-select" :class="[{ 'is-open': open, 'has-error': error, 'is-disabled': disabled }, direction]">
    <span v-if="label" class="ft-select__label">{{ label }}</span>
    <button type="button" class="ft-select__button" :disabled="disabled" @click.stop="toggle">
      <span class="ft-select__value" :class="{ muted: !selected }">{{ selectedLabel }}</span>
      <span class="ft-select__chevron">⌄</span>
    </button>
    <Teleport to="body">
      <div
        v-if="open"
        class="ft-select__dropdown"
        :class="direction"
:style="dropdownStyle"
      >
        <div v-if="!normalizedOptions.length" class="ft-select__empty">{{ emptyText }}</div>
        <button
          v-for="option in normalizedOptions"
          :key="String(option.value)"
          type="button"
          class="ft-select__option"
          :class="{ active: String(option.value) === String(modelValue), disabled: option.disabled }"
          :disabled="option.disabled"
          @click="choose(option)"
        >
          <span>
            <strong>{{ option.label }}</strong>
            <small v-if="option.description">{{ option.description }}</small>
          </span>
          <em v-if="option.badge">{{ option.badge }}</em>
          <b v-if="String(option.value) === String(modelValue)">✓</b>
        </button>
      </div>
    </Teleport>
    <small v-if="error" class="ft-select__error">{{ error }}</small>
  </div>
</template>

<style>
.ft-select-lock {
  overflow: hidden !important;
  overscroll-behavior: none !important;
}
</style>

<style scoped>
.ft-select { position: relative; width: 100%; }
.ft-select__label {
  display: block;
  margin: 0 0 10px;
  color: #97a5bb;
  font-size: 13px;
  font-weight: 900;
  letter-spacing: .22em;
  text-transform: uppercase;
}
.ft-select__button {
  width: 100%; min-height: 64px; border: 1px solid #dbe4ef; border-radius: 22px;
  background: #f8fbff; padding: 0 22px; display: flex; align-items: center; justify-content: space-between;
  gap: 14px; color: #050b1a; font-weight: 900; font-size: 18px; text-align: left;
  transition: border-color .2s, box-shadow .2s, background .2s;
}
.ft-select__button:hover, .ft-select.is-open .ft-select__button { border-color: #ff3f4d; box-shadow: 0 0 0 5px rgba(255,63,77,.12); background: #fff; }
.ft-select.has-error .ft-select__button { border-color: #ff3f4d; box-shadow: 0 0 0 5px rgba(255,63,77,.10); }
.ft-select.is-disabled { opacity: .55; pointer-events: none; }
.ft-select__value.muted { color: #94a3b8; }
.ft-select__chevron { color: #66758a; font-size: 22px; transition: transform .2s; }
.ft-select.is-open .ft-select__chevron { transform: rotate(180deg); }
.ft-select__error { display: block; margin-top: 8px; color: #ff3f4d; font-weight: 800; }
.ft-select__dropdown {
  position: fixed; z-index: 9999; max-height: 320px; overflow: auto; padding: 10px; border-radius: 22px;
  background: #fff; border: 1px solid #dbe4ef; box-shadow: 0 24px 60px rgba(5,11,26,.20);
}
.ft-select__dropdown::-webkit-scrollbar { width: 8px; }
.ft-select__dropdown::-webkit-scrollbar-thumb { border-radius: 999px; background: #cbd5e1; }
.ft-select__empty { padding: 18px; color: #64748b; font-weight: 800; }
.ft-select__option {
  width: 100%; border: 0; border-radius: 16px; background: transparent; padding: 15px 16px;
  display: flex; align-items: center; justify-content: space-between; gap: 12px; color: #071022; text-align: left; cursor: pointer;
}
.ft-select__option:hover { background: #f2f6fb; }
.ft-select__option.active { background: #ff3f4d; color: #fff; }
.ft-select__option.disabled { opacity: .45; cursor: not-allowed; }
.ft-select__option strong { display: block; font-size: 16px; font-weight: 900; }
.ft-select__option small { display: block; margin-top: 4px; color: #687890; font-size: 12px; font-weight: 800; }
.ft-select__option.active small { color: rgba(255,255,255,.82); }
.ft-select__option em { font-style: normal; font-size: 11px; font-weight: 900; padding: 6px 9px; border-radius: 999px; background: #eef4fb; color: #475569; }
.ft-select__option.active em { background: rgba(255,255,255,.16); color: #fff; }
.ft-select__option b { font-size: 18px; }
</style>
