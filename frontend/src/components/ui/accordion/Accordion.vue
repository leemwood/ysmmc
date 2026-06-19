<script setup lang="ts">
import { AccordionRoot, useForwardPropsEmits } from 'radix-vue'
import type { AccordionRootProps } from 'radix-vue'

const props = defineProps<AccordionRootProps>()
const emits = defineEmits<{
  (e: 'update:modelValue', value: string | string[]): void
}>()

const forwarded = useForwardPropsEmits(props, emits)
</script>

<template>
  <AccordionRoot v-bind="forwarded">
    <slot />
  </AccordionRoot>
</template>
