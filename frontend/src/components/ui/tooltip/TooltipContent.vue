<script setup lang="ts">
import { TooltipContent, TooltipPortal } from 'radix-vue'
import { cn } from '@/lib/utils'
import type { TooltipContentProps } from 'radix-vue'

const props = withDefaults(
  defineProps<TooltipContentProps & { class?: string }>(),
  {
    sideOffset: 4,
  }
)
</script>

<template>
  <TooltipPortal>
    <TooltipContent
      v-bind="props"
      :class="cn(
        'z-50 overflow-hidden rounded-md border bg-popover px-3 py-1.5 text-sm text-popover-foreground shadow-md animate-fade-in',
        props.class
      )"
    >
      <slot />
    </TooltipContent>
  </TooltipPortal>
</template>
