<script setup lang="ts">
import { cn } from '@/lib/utils'

interface Props {
  class?: string
  type?: string
  placeholder?: string
  disabled?: boolean
  modelValue?: string
  error?: boolean
}

const props = defineProps<Props>()
const emit = defineEmits(['update:modelValue'])

function handleInput(event: Event) {
  emit('update:modelValue', (event.target as HTMLInputElement).value)
}
</script>

<template>
  <div class="relative flex items-center">
    <span v-if="$slots.prefix" class="absolute left-3 text-muted-foreground">
      <slot name="prefix" />
    </span>
    <input
      :class="cn(
        'flex h-9 w-full rounded-md border bg-transparent px-3 py-1 text-sm shadow-sm transition-colors file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 disabled:cursor-not-allowed disabled:opacity-50',
        props.error
          ? 'border-destructive focus-visible:ring-destructive text-destructive placeholder:text-destructive/60'
          : 'border-input focus-visible:ring-ring',
        $slots.prefix && 'pl-9',
        $slots.suffix && 'pr-9',
        props.class
      )"
      :type="type || 'text'"
      :placeholder="placeholder"
      :disabled="disabled"
      :value="modelValue"
      @input="handleInput"
    />
    <span v-if="$slots.suffix" class="absolute right-3 text-muted-foreground">
      <slot name="suffix" />
    </span>
  </div>
</template>
