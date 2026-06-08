<template>
  <div class="flex items-center gap-6">
    <div class="relative shrink-0" :style="{ width: size + 'px', height: size + 'px' }">
      <div class="absolute inset-0 rounded-full" :style="ringStyle"></div>
      <div class="absolute rounded-full bg-white dark:bg-slate-900" style="inset: 22%"></div>
      <div class="absolute inset-0 flex flex-col items-center justify-center">
        <span class="text-xl font-bold text-slate-900 dark:text-white">{{ total }}</span>
        <span class="text-[10px] uppercase tracking-wide text-slate-400">{{ centerLabel }}</span>
      </div>
    </div>
    <ul class="flex-1 space-y-2 text-sm">
      <li v-for="seg in segments" :key="seg.label" class="flex items-center gap-2">
        <span class="h-2.5 w-2.5 rounded-full" :style="{ backgroundColor: seg.color }"></span>
        <span class="text-slate-600 dark:text-slate-300">{{ seg.label }}</span>
        <span class="ml-auto font-medium tabular-nums text-slate-900 dark:text-white">{{ seg.value }}</span>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  segments: { type: Array, default: () => [] }, // [{ label, value, color }]
  size: { type: Number, default: 120 },
  centerLabel: { type: String, default: 'total' },
})

const total = computed(() => props.segments.reduce((s, x) => s + (x.value || 0), 0))

const ringStyle = computed(() => {
  if (!total.value) {
    return { background: 'var(--donut-empty, #e2e8f0)' }
  }
  let acc = 0
  const stops = props.segments
    .filter((s) => s.value > 0)
    .map((s) => {
      const start = (acc / total.value) * 100
      acc += s.value
      const end = (acc / total.value) * 100
      return `${s.color} ${start}% ${end}%`
    })
  return { background: `conic-gradient(${stops.join(', ')})` }
})
</script>
