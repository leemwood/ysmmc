<script setup lang="ts">
import { DropdownMenuItem } from 'radix-vue'
import { cn } from '@/lib/utils'
import type { DropdownMenuItemProps } from 'radix-vue'

const props = defineProps<
  DropdownMenuItemProps & { class?: string; inset?: boolean }
>()
</script>

<template>
  <DropdownMenuItem
    v-bind="props"
    :class="cn(
      'relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
      inset && 'pl-8',
      props.class
    )"
  >
    <slot />
  </DropdownMenuItem>
</template>
