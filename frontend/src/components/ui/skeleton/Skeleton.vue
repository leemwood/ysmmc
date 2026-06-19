<script setup lang="ts">
import { cn } from '@/lib/utils'

interface Props {
  class?: string
  variant?: 'default' | 'shimmer'
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
})
</script>

<template>
  <div
    :class="cn(
      'rounded-md',
      props.variant === 'shimmer' ? 'animate-shimmer' : 'animate-pulse bg-primary/10',
      props.class
    )"
  />
</template>
