<script setup lang="ts">
import { cn } from '@/lib/utils'

interface Props {
  class?: string
  placeholder?: string
  disabled?: boolean
  rows?: number
  modelValue?: string
}

const props = withDefaults(defineProps<Props>(), {
  rows: 4,
})

const emit = defineEmits(['update:modelValue'])

function handleInput(event: Event) {
  emit('update:modelValue', (event.target as HTMLTextAreaElement).value)
}
</script>

<template>
  <textarea
    :class="cn('flex min-h-[60px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50', props.class)"
    :placeholder="placeholder"
    :disabled="disabled"
    :rows="rows"
    :value="modelValue"
    @input="handleInput"
  />
</template>
