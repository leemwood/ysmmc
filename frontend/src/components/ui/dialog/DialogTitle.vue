<script setup lang="ts">
import { DialogTitle } from 'radix-vue'
import { cn } from '@/lib/utils'

defineProps<{
  class?: string
}>()
</script>

<template>
  <DialogTitle :class="cn('text-lg font-semibold leading-none tracking-tight', $props.class)">
    <slot />
  </DialogTitle>
</template>
