<script setup lang="ts">
import { cn } from '@/lib/utils'

interface Props {
  class?: string
  for?: string
}

const props = defineProps<Props>()
</script>

<template>
  <label
    :class="cn('text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70', props.class)"
    :for="props.for"
  >
    <slot />
  </label>
</template>
