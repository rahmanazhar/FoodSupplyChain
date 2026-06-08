<template>
  <div class="card transition-shadow hover:shadow-card-hover">
    <div class="flex items-center justify-between">
      <p class="text-sm font-medium text-slate-500 dark:text-slate-400">{{ label }}</p>
      <span class="flex h-9 w-9 items-center justify-center rounded-lg" :class="toneClass">
        <slot name="icon" />
      </span>
    </div>
    <p class="mt-3 text-3xl font-semibold tracking-tight" :class="valueClass">
      <slot>{{ value }}</slot>
    </p>
    <p v-if="hint" class="mt-1 text-xs text-slate-400">{{ hint }}</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: { type: String, required: true },
  value: { type: [String, Number], default: '' },
  hint: { type: String, default: '' },
  tone: { type: String, default: 'primary' }, // primary | emerald | red | amber | slate
})

const tones = {
  primary: 'bg-primary-100 text-primary-600 dark:bg-primary-500/15 dark:text-primary-400',
  emerald: 'bg-emerald-100 text-emerald-600 dark:bg-emerald-500/15 dark:text-emerald-400',
  red: 'bg-red-100 text-red-600 dark:bg-red-500/15 dark:text-red-400',
  amber: 'bg-amber-100 text-amber-600 dark:bg-amber-500/15 dark:text-amber-400',
  slate: 'bg-slate-100 text-slate-600 dark:bg-slate-700/50 dark:text-slate-300',
}
const toneClass = computed(() => tones[props.tone] || tones.primary)
const valueClass = computed(() =>
  props.tone === 'red' ? 'text-red-600 dark:text-red-400' : 'text-slate-900 dark:text-white'
)
</script>
