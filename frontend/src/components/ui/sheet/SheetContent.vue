<script setup lang="ts">
import {
  DialogContent,
  DialogOverlay,
  DialogPortal,
} from 'radix-vue'
import { X } from 'lucide-vue-next'
import { cn } from '@/lib/utils'
import type { DialogContentProps } from 'radix-vue'

interface Props extends DialogContentProps {
  class?: string
  side?: 'top' | 'right' | 'bottom' | 'left'
}

const props = withDefaults(defineProps<Props>(), {
  side: 'right',
})

const sideClasses = {
  top: 'inset-x-0 top-0 h-auto max-h-[80vh] border-b rounded-b-xl',
  bottom: 'inset-x-0 bottom-0 h-auto max-h-[90vh] border-t rounded-t-xl',
  left: 'inset-y-0 left-0 h-full w-3/4 max-w-sm border-r rounded-r-xl',
  right: 'inset-y-0 right-0 h-full w-3/4 max-w-sm border-l rounded-l-xl',
}
</script>

<template>
  <DialogPortal>
    <DialogOverlay class="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm data-[state=open]:animate-fade-in data-[state=closed]:animate-fade-in" />
    <DialogContent
      v-bind="props"
      :class="cn(
        'fixed z-50 flex flex-col gap-4 bg-background p-6 shadow-xl transition ease-in-out data-[state=open]:animate-slide-up data-[state=closed]:animate-fade-in',
        sideClasses[side],
        props.class
      )"
    >
      <slot />
      <button
        class="absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none"
      >
        <X class="h-4 w-4" />
        <span class="sr-only">关闭</span>
      </button>
    </DialogContent>
  </DialogPortal>
</template>
