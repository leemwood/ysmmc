<script setup lang="ts">
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

defineProps<{ page: number; totalPages: number }>()
defineEmits<{ change: [page: number] }>()
</script>

<template>
  <nav v-if="totalPages > 1" aria-label="分页" class="mt-7 flex flex-wrap items-center justify-center gap-3 sm:gap-5">
    <button type="button" class="site-button site-button-outline" :disabled="page <= 1" @click="$emit('change', page - 1)">
      <ChevronLeft class="size-4" aria-hidden="true" /> 上一页
    </button>
    <span class="text-sm tabular-nums text-muted-foreground">{{ page }} / {{ totalPages }} 页</span>
    <button type="button" class="site-button site-button-outline" :disabled="page >= totalPages" @click="$emit('change', page + 1)">
      下一页 <ChevronRight class="size-4" aria-hidden="true" />
    </button>
  </nav>
</template>
