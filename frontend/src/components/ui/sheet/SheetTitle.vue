<script setup lang="ts">
import { DialogTitle } from 'radix-vue'
import { cn } from '@/lib/utils'
import type { DialogTitleProps } from 'radix-vue'

const props = defineProps<DialogTitleProps & { class?: string }>()
</script>

<template>
  <DialogTitle
    v-bind="props"
    :class="cn('text-lg font-semibold text-foreground', props.class)"
  >
    <slot />
  </DialogTitle>
</template>
