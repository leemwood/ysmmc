<script setup lang="ts">
import { DialogDescription } from 'radix-vue'
import { cn } from '@/lib/utils'
import type { DialogDescriptionProps } from 'radix-vue'

const props = defineProps<DialogDescriptionProps & { class?: string }>()
</script>

<template>
  <DialogDescription
    v-bind="props"
    :class="cn('text-sm text-muted-foreground', props.class)"
  >
    <slot />
  </DialogDescription>
</template>
