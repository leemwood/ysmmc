<script setup lang="ts">
import { onMounted, onUnmounted, computed } from 'vue'
import { CheckCircle, XCircle, AlertTriangle, Info, X } from 'lucide-vue-next'
import type { Toast } from '@/composables/useToast'
import { useToast } from '@/composables/useToast'

const props = defineProps<{ toast: Toast }>()

const { removeToast } = useToast()

let timer: ReturnType<typeof setTimeout> | null = null

function startTimer() {
  if (props.toast.duration > 0) {
    timer = setTimeout(() => {
      removeToast(props.toast.id)
    }, props.toast.duration)
  }
}

function clearTimer() {
  if (timer) {
    clearTimeout(timer)
    timer = null
  }
}

onMounted(startTimer)
onUnmounted(clearTimer)

const iconComponent = computed(() => {
  switch (props.toast.type) {
    case 'success':
      return CheckCircle
    case 'error':
      return XCircle
    case 'warning':
      return AlertTriangle
    case 'info':
      return Info
  }
})

const iconColor = computed(() => {
  switch (props.toast.type) {
    case 'success':
      return 'text-green-500'
    case 'error':
      return 'text-destructive'
    case 'warning':
      return 'text-yellow-500'
    case 'info':
      return 'text-blue-500'
  }
})
</script>

<template>
  <div
    class="flex items-start gap-3 w-full max-w-sm rounded-lg border bg-background p-4 shadow-lg"
    role="alert"
  >
    <component :is="iconComponent" class="h-5 w-5 flex-shrink-0 mt-0.5" :class="iconColor" />
    <div class="flex-1 min-w-0">
      <p v-if="toast.title" class="text-sm font-semibold text-foreground">{{ toast.title }}</p>
      <p class="text-sm text-foreground break-words" :class="toast.title ? 'mt-0.5' : ''">
        {{ toast.message }}
      </p>
    </div>
    <button
      type="button"
      class="flex-shrink-0 rounded-md p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
      aria-label="关闭"
      @click="removeToast(toast.id)"
    >
      <X class="h-4 w-4" />
    </button>
  </div>
</template>
