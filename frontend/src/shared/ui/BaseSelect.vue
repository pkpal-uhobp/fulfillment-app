<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { ChevronDown, Check } from '@lucide/vue'

const props = defineProps({
  modelValue: { type: [String, Number, Boolean, null], default: '' },
  options: { type: Array, default: () => [] },
  placeholder: { type: String, default: 'Выберите значение' },
  disabled: { type: Boolean, default: false },
  compact: { type: Boolean, default: false },
  error: { type: String, default: '' },
  popupClass: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue', 'change', 'open', 'close'])
const root = ref(null)
const open = ref(false)
const menuStyle = ref({})
const previousOverflow = ref('')
const previousPaddingRight = ref('')

const normalizedOptions = computed(() => props.options.map((option) => {
  if (typeof option === 'object' && option !== null) return option
  return { value: option, label: String(option) }
}))

const selected = computed(() => normalizedOptions.value.find((option) => String(option.value) === String(props.modelValue)))
const selectedLabel = computed(() => selected.value?.label || props.placeholder)

function scrollbarWidth() {
  return window.innerWidth - document.documentElement.clientWidth
}

function lockScroll() {
  previousOverflow.value = document.body.style.overflow
  previousPaddingRight.value = document.body.style.paddingRight
  const width = scrollbarWidth()
  document.body.style.overflow = 'hidden'
  if (width > 0) document.body.style.paddingRight = `${width}px`
}

function unlockScroll() {
  document.body.style.overflow = previousOverflow.value
  document.body.style.paddingRight = previousPaddingRight.value
}

function updatePosition() {
  if (!root.value) return
  const rect = root.value.getBoundingClientRect()
  const gap = 8
  const margin = 16
  const width = Math.min(rect.width, window.innerWidth - margin * 2)
  const left = Math.min(Math.max(rect.left, margin), window.innerWidth - width - margin)
  const spaceBelow = window.innerHeight - rect.bottom - margin
  const spaceAbove = rect.top - margin
  const openAbove = spaceBelow < 240 && spaceAbove > spaceBelow
  const maxHeight = Math.max(180, Math.min(360, openAbove ? spaceAbove - gap : spaceBelow - gap))
  const top = openAbove ? Math.max(margin, rect.top - gap - maxHeight) : rect.bottom + gap

  menuStyle.value = {
    left: `${left}px`,
    top: `${top}px`,
    width: `${width}px`,
    maxHeight: `${maxHeight}px`,
  }
}

async function setOpen(value) {
  if (props.disabled) return
  open.value = value
  if (value) {
    await nextTick()
    updatePosition()
    emit('open')
  } else {
    emit('close')
  }
}

function toggle() {
  setOpen(!open.value)
}

function choose(option) {
  if (option.disabled) return
  emit('update:modelValue', option.value)
  emit('change', option.value)
  setOpen(false)
}

function onDocumentClick(event) {
  if (!root.value?.contains(event.target)) setOpen(false)
}

function onKeydown(event) {
  if (event.key === 'Escape') setOpen(false)
}

watch(open, (value) => {
  if (value) lockScroll()
  else unlockScroll()
})

onMounted(() => {
  document.addEventListener('click', onDocumentClick)
  document.addEventListener('keydown', onKeydown)
  window.addEventListener('resize', updatePosition)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onDocumentClick)
  document.removeEventListener('keydown', onKeydown)
  window.removeEventListener('resize', updatePosition)
  if (open.value) unlockScroll()
})
</script>

<template>
  <div ref="root" class="relative">
    <button
      type="button"
      class="group flex w-full items-center justify-between gap-3 rounded-[1.35rem] border bg-slate-50 text-left font-black text-[#07101f] outline-none transition disabled:cursor-not-allowed disabled:opacity-60"
      :class="[
        compact ? 'px-4 py-3 text-sm' : 'px-5 py-4 text-base',
        error ? 'border-[#ff4248] bg-red-50/60 shadow-[0_0_0_4px_rgba(255,66,72,0.10)]' : open ? 'border-[#ff4248] bg-white shadow-[0_0_0_4px_rgba(255,66,72,0.10)]' : 'border-slate-200 hover:border-slate-300',
      ]"
      :disabled="disabled"
      @click="toggle"
    >
      <span class="min-w-0 truncate" :class="selected ? 'text-[#07101f]' : 'text-slate-400'">
        {{ selectedLabel }}
      </span>
      <ChevronDown class="h-5 w-5 shrink-0 text-slate-500 transition" :class="open ? 'rotate-180 text-[#ff4248]' : ''" />
    </button>

    <p v-if="error" class="mt-2 text-sm font-bold text-[#e11d48]">{{ error }}</p>

    <Teleport to="body">
      <div v-if="open" class="fixed inset-0 z-[999] bg-transparent" @click.self="setOpen(false)"></div>
      <div
        v-if="open"
        class="fixed z-[1000] overflow-y-auto overscroll-contain rounded-[1.35rem] border border-slate-200 bg-white p-2 text-[#07101f] shadow-[0_30px_90px_rgba(15,23,42,0.24)]"
        :class="popupClass"
        :style="menuStyle"
      >
        <button
          v-for="option in normalizedOptions"
          :key="String(option.value)"
          type="button"
          class="flex w-full items-start justify-between gap-3 rounded-2xl px-4 py-3 text-left font-black transition"
          :class="[
            String(option.value) === String(modelValue) ? 'bg-[#ff4248] text-white hover:bg-[#ff4248]' : 'text-[#07101f] hover:bg-slate-50',
            option.disabled ? 'cursor-not-allowed opacity-45' : '',
          ]"
          :disabled="option.disabled"
          @click="choose(option)"
        >
          <span class="min-w-0">
            <span class="block leading-5">{{ option.label }}</span>
            <span v-if="option.description" class="mt-1 block text-xs font-bold opacity-65">{{ option.description }}</span>
          </span>
          <Check v-if="String(option.value) === String(modelValue)" class="mt-0.5 h-4 w-4 shrink-0" />
        </button>

        <div v-if="!normalizedOptions.length" class="rounded-2xl bg-slate-50 px-4 py-4 text-sm font-bold text-slate-500">
          Нет доступных вариантов
        </div>
      </div>
    </Teleport>
  </div>
</template>
