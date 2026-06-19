<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/utils'

interface Props {
  modelValue?: number
  max?: number
  class?: string
  showValue?: boolean
  size?: 'sm' | 'default' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: 0,
  max: 100,
  showValue: false,
  size: 'default',
})

const percentage = computed(() => {
  const value = Math.min(Math.max(props.modelValue, 0), props.max)
  return props.max === 0 ? 0 : Math.round((value / props.max) * 100)
})

const sizeClasses = {
  sm: 'h-1.5',
  default: 'h-2',
  lg: 'h-3',
}
</script>

<template>
  <div :class="cn('w-full', props.class)">
    <div
      v-if="showValue"
      class="mb-1.5 flex justify-between text-xs text-muted-foreground"
    >
      <slot />
      <span>{{ percentage }}%</span>
    </div>
    <div
      role="progressbar"
      :aria-valuenow="modelValue"
      :aria-valuemax="max"
      aria-valuemin="0"
      :class="cn('relative w-full overflow-hidden rounded-full bg-muted', sizeClasses[size])"
    >
      <div
        class="h-full bg-primary transition-all duration-300 ease-out"
        :style="{ width: percentage + '%' }"
      />
    </div>
  </div>
</template>
