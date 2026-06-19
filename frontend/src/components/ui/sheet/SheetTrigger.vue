<script setup lang="ts">
import { DialogTrigger } from 'radix-vue'
import type { DialogTriggerProps } from 'radix-vue'

const props = defineProps<DialogTriggerProps & { class?: string }>()
</script>

<template>
  <DialogTrigger v-bind="props" as-child>
    <slot />
  </DialogTrigger>
</template>
