import { ref } from 'vue'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

export interface Toast {
  id: number
  type: ToastType
  message: string
  title?: string
  duration: number
}

const toasts = ref<Toast[]>([])
let toastId = 0

const MAX_TOASTS = 5
const DEFAULT_DURATION = 3000

function addToast(type: ToastType, message: string, options?: { title?: string; duration?: number }) {
  const id = ++toastId
  const duration = options?.duration ?? DEFAULT_DURATION
  const toast: Toast = {
    id,
    type,
    message,
    title: options?.title,
    duration,
  }

  toasts.value.push(toast)

  // 超出最大数量时移除最早的
  if (toasts.value.length > MAX_TOASTS) {
    toasts.value.shift()
  }

  return id
}

function removeToast(id: number) {
  const index = toasts.value.findIndex(t => t.id === id)
  if (index !== -1) {
    toasts.value.splice(index, 1)
  }
}

export function useToast() {
  const toast = {
    success(message: string, options?: { title?: string; duration?: number }) {
      return addToast('success', message, options)
    },
    error(message: string, options?: { title?: string; duration?: number }) {
      return addToast('error', message, options)
    },
    warning(message: string, options?: { title?: string; duration?: number }) {
      return addToast('warning', message, options)
    },
    info(message: string, options?: { title?: string; duration?: number }) {
      return addToast('info', message, options)
    },
    remove: removeToast,
  }

  return {
    toasts,
    toast,
    removeToast,
  }
}
