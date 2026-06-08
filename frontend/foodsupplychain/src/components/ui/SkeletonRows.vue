<template>
  <div class="divide-y divide-slate-100 dark:divide-slate-800">
    <div v-for="r in rows" :key="r" class="flex items-center gap-4 px-4 py-4">
      <div v-for="c in cols" :key="c" class="skeleton h-4" :style="{ width: widths[(c - 1) % widths.length] }"></div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  rows: { type: Number, default: 6 },
  cols: { type: Number, default: 4 },
})
const widths = ['30%', '14%', '20%', '10%', '12%', '8%']
</script>
