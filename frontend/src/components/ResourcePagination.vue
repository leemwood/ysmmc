<script setup lang="ts">
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'
import { RouterLink, type RouteLocationRaw } from 'vue-router'

defineProps<{ page: number; totalPages: number; to?: (page: number) => RouteLocationRaw }>()
defineEmits<{ change: [page: number] }>()
</script>

<template>
  <nav v-if="totalPages > 1" aria-label="分页" class="mt-7 flex flex-wrap items-center justify-center gap-3 sm:gap-5">
    <RouterLink v-if="to && page > 1" :to="to(page - 1)" class="site-button site-button-outline">
      <ChevronLeft class="size-4" aria-hidden="true" /> 上一页
    </RouterLink>
    <button v-else type="button" class="site-button site-button-outline" :disabled="page <= 1" @click="$emit('change', page - 1)">
      <ChevronLeft class="size-4" aria-hidden="true" /> 上一页
    </button>
    <span class="text-sm tabular-nums text-muted-foreground">{{ page }} / {{ totalPages }} 页</span>
    <RouterLink v-if="to && page < totalPages" :to="to(page + 1)" class="site-button site-button-outline">
      下一页 <ChevronRight class="size-4" aria-hidden="true" />
    </RouterLink>
    <button v-else type="button" class="site-button site-button-outline" :disabled="page >= totalPages" @click="$emit('change', page + 1)">
      下一页 <ChevronRight class="size-4" aria-hidden="true" />
    </button>
  </nav>
</template>
