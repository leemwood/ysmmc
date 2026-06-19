<script setup lang="ts">
import { DropdownMenuTrigger } from 'radix-vue'
import type { DropdownMenuTriggerProps } from 'radix-vue'

const props = defineProps<DropdownMenuTriggerProps & { class?: string }>()
</script>

<template>
  <DropdownMenuTrigger v-bind="props" :class="props.class" as-child>
    <slot />
  </DropdownMenuTrigger>
</template>
