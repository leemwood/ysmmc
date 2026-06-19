<script setup lang="ts">
import { DropdownMenuContent, DropdownMenuPortal } from 'radix-vue'
import { cn } from '@/lib/utils'
import type { DropdownMenuContentProps } from 'radix-vue'

const props = withDefaults(
  defineProps<DropdownMenuContentProps & { class?: string }>(),
  {
    sideOffset: 4,
  }
)
</script>

<template>
  <DropdownMenuPortal>
    <DropdownMenuContent
      v-bind="props"
      :class="cn(
        'z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md animate-fade-in',
        props.class
      )"
    >
      <slot />
    </DropdownMenuContent>
  </DropdownMenuPortal>
</template>
