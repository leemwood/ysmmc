<script setup lang="ts">
import { TabsRoot, useForwardPropsEmits } from 'radix-vue'
import type { TabsRootProps } from 'radix-vue'

const props = defineProps<TabsRootProps>()
const emits = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const forwarded = useForwardPropsEmits(props, emits)
</script>

<template>
  <TabsRoot v-bind="forwarded">
    <slot />
  </TabsRoot>
</template>
