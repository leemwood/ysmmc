<script setup lang="ts">
import { AccordionItem } from 'radix-vue'
import { cn } from '@/lib/utils'
import type { AccordionItemProps } from 'radix-vue'

const props = defineProps<AccordionItemProps & { class?: string }>()
</script>

<template>
  <AccordionItem
    v-bind="props"
    :class="cn('border-b', props.class)"
  >
    <slot />
  </AccordionItem>
</template>
