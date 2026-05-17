<script setup>
defineProps({ message: { type: String, default: '' } })
</script>

<template>
  <p v-if="message" class="mt-2 text-sm font-bold text-[#e11d48]">{{ message }}</p>
</template>
