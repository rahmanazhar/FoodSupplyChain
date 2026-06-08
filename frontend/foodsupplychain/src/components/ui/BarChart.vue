<template>
  <div class="space-y-3">
    <div v-for="item in items" :key="item.label">
      <div class="mb-1 flex items-center justify-between text-sm">
        <span class="capitalize text-slate-600 dark:text-slate-300">{{ item.label }}</span>
        <span class="font-medium tabular-nums text-slate-900 dark:text-white">{{ item.value }}</span>
      </div>
      <div class="h-2 overflow-hidden rounded-full bg-slate-100 dark:bg-slate-800">
        <div
          class="h-full rounded-full transition-all duration-500"
          :style="{ width: pct(item.value) + '%', backgroundColor: item.color || defaultColor }"
        ></div>
      </div>
    </div>
    <p v-if="!items.length" class="py-6 text-center text-sm text-slate-400">No data</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  items: { type: Array, default: () => [] }, // [{ label, value, color? }]
  defaultColor: { type: String, default: '#0ea5e9' },
})

const max = computed(() => Math.max(1, ...props.items.map((i) => i.value || 0)))
const pct = (v) => Math.round(((v || 0) / max.value) * 100)
</script>
