<template>
  <div class="flex items-center justify-between border-t border-slate-100 px-4 py-3 text-sm text-slate-500 dark:border-slate-800 dark:text-slate-400">
    <span>{{ rangeLabel }}</span>
    <div class="flex items-center gap-2">
      <button class="btn-secondary btn-sm" :disabled="offset === 0" @click="$emit('change', Math.max(0, offset - limit))">
        Prev
      </button>
      <span class="tabular-nums">Page {{ page }} / {{ pages }}</span>
      <button class="btn-secondary btn-sm" :disabled="offset + limit >= total" @click="$emit('change', offset + limit)">
        Next
      </button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  total: { type: Number, default: 0 },
  limit: { type: Number, default: 20 },
  offset: { type: Number, default: 0 },
})
defineEmits(['change'])

const page = computed(() => Math.floor(props.offset / props.limit) + 1)
const pages = computed(() => Math.max(1, Math.ceil(props.total / props.limit)))
const rangeLabel = computed(() => {
  if (!props.total) return '0 results'
  const from = props.offset + 1
  const to = Math.min(props.offset + props.limit, props.total)
  return `${from}–${to} of ${props.total}`
})
</script>
