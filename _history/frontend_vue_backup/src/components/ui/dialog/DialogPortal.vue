<script setup lang="ts">
import { DialogPortal } from 'radix-vue'

defineProps<{
  to?: string | HTMLElement
  forceMount?: boolean
}>()
</script>

<template>
  <DialogPortal :to="to" :force-mount="forceMount">
    <slot />
  </DialogPortal>
</template>
