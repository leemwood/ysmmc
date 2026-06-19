<script setup lang="ts">
import { DropdownMenuRoot, useForwardPropsEmits } from 'radix-vue'
import type { DropdownMenuRootProps } from 'radix-vue'

const props = defineProps<DropdownMenuRootProps>()
const emits = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const forwarded = useForwardPropsEmits(props, emits)
</script>

<template>
  <DropdownMenuRoot v-bind="forwarded">
    <slot />
  </DropdownMenuRoot>
</template>
