<script setup>
import { computed } from 'vue'
import BaseSelect from './BaseSelect.vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
  from: { type: Number, default: 0 },
  to: { type: Number, default: 23 },
  step: { type: Number, default: 60 },
  placeholder: { type: String, default: 'Время' },
  error: { type: String, default: '' },
})

const emit = defineEmits(['update:modelValue', 'change'])

const options = computed(() => {
  const result = []
  for (let hour = props.from; hour <= props.to; hour += 1) {
    for (let minute = 0; minute < 60; minute += props.step) {
      const value = `${String(hour).padStart(2, '0')}:${String(minute).padStart(2, '0')}`
      result.push({ value, label: value })
    }
  }
  return result
})

function update(value) {
  emit('update:modelValue', value)
  emit('change', value)
}
</script>

<template>
  <BaseSelect :model-value="modelValue" :options="options" :placeholder="placeholder" :error="error" compact @update:model-value="update" @change="update" />
</template>
