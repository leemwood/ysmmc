<script setup lang="ts">
import { cn } from '@/lib/utils'

interface Props {
  orientation?: 'horizontal' | 'vertical'
  decorative?: boolean
  class?: string
}

const props = withDefaults(defineProps<Props>(), {
  orientation: 'horizontal',
  decorative: true,
})
</script>

<template>
  <div
    role="separator"
    :aria-orientation="orientation"
    :aria-hidden="decorative"
    :class="cn(
      'shrink-0 bg-border',
      orientation === 'horizontal' ? 'h-[1px] w-full' : 'h-full w-[1px]',
      props.class
    )"
  />
</template>
