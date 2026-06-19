<script setup lang="ts">
import { TooltipTrigger } from 'radix-vue'
import type { TooltipTriggerProps } from 'radix-vue'

const props = defineProps<TooltipTriggerProps & { class?: string }>()
</script>

<template>
  <TooltipTrigger v-bind="props" as-child>
    <slot />
  </TooltipTrigger>
</template>
