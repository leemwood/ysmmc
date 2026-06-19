<script setup lang="ts">
import { DialogContent, DialogPortal, DialogOverlay } from 'radix-vue'
import { cn } from '@/lib/utils'

interface Props {
  class?: string
  forceMount?: boolean
  fullScreen?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  fullScreen: false,
})
</script>

<template>
  <DialogPortal>
    <DialogOverlay
      :class="cn(
        'fixed inset-0 z-50 bg-black/80 backdrop-blur-sm transition-opacity animate-fade-in',
        props.fullScreen && 'bg-black/90'
      )"
    />
    <DialogContent
      :class="cn(
        'fixed z-50 grid w-full gap-4 border bg-background p-6 text-foreground shadow-2xl transition-all animate-fade-in',
        props.fullScreen
          ? 'inset-0 h-screen w-screen rounded-none'
          : 'left-[50%] top-[50%] max-w-lg translate-x-[-50%] translate-y-[-50%] rounded-lg',
        props.class
      )"
    >
      <slot />
    </DialogContent>
  </DialogPortal>
</template>
