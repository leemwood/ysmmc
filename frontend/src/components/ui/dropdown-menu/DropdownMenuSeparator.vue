<script setup lang="ts">
import { DropdownMenuSeparator } from 'radix-vue'
import { cn } from '@/lib/utils'
import type { DropdownMenuSeparatorProps } from 'radix-vue'

const props = defineProps<DropdownMenuSeparatorProps & { class?: string }>()
</script>

<template>
  <DropdownMenuSeparator
    v-bind="props"
    :class="cn('-mx-1 my-1 h-px bg-muted', props.class)"
  />
</template>
