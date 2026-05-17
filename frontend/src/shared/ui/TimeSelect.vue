<template>
  <BaseSelect
    :model-value="modelValue"
    :options="options"
    :label="label"
    :placeholder="placeholder"
    :error="error"
    :disabled="disabled"
    empty-text="Нет доступного времени"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script setup>
import { computed } from 'vue'
import BaseSelect from './BaseSelect.vue'

const props = defineProps({
  modelValue: { type: String, default: '' },
  label: { type: String, default: 'Время' },
  placeholder: { type: String, default: 'Выберите время' },
  error: { type: String, default: '' },
  disabled: { type: Boolean, default: false },
  startHour: { type: Number, default: 9 },
  endHour: { type: Number, default: 21 },
  stepMinutes: { type: Number, default: 60 },
})

defineEmits(['update:modelValue'])

const options = computed(() => {
  const result = []
  for (let total = props.startHour * 60; total <= props.endHour * 60; total += props.stepMinutes) {
    const value = `${String(Math.floor(total / 60)).padStart(2, '0')}:${String(total % 60).padStart(2, '0')}`
    result.push({ value, label: value, description: total < 12 * 60 ? 'утренний слот' : total < 17 * 60 ? 'дневной слот' : 'вечерний слот' })
  }
  return result
})
</script>
