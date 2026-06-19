<script setup lang="ts">
import { ref, watch, nextTick, onUnmounted } from 'vue'
import { Dialog, DialogContent } from '@/components/ui/dialog'
import { X, ChevronLeft, ChevronRight } from 'lucide-vue-next'
import { fileApi } from '@/lib/api'
import type { ModelImage } from '@/types'

interface Props {
  open: boolean
  images: ModelImage[]
  initialIndex?: number
}

const props = withDefaults(defineProps<Props>(), {
  initialIndex: 0,
})

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'update:index', value: number): void
}>()

const currentIndex = ref(props.initialIndex)
const closeButtonRef = ref<HTMLButtonElement | null>(null)

function getImageUrl(fileId: string) {
  return fileApi.getUrl(fileId)
}

function close() {
  emit('update:open', false)
}

function prev() {
  if (currentIndex.value > 0) {
    currentIndex.value--
  } else {
    currentIndex.value = props.images.length - 1
  }
  emitIndex()
}

function next() {
  if (currentIndex.value < props.images.length - 1) {
    currentIndex.value++
  } else {
    currentIndex.value = 0
  }
  emitIndex()
}

function emitIndex() {
  emit('update:index', currentIndex.value)
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'ArrowLeft') {
    e.preventDefault()
    prev()
  } else if (e.key === 'ArrowRight') {
    e.preventDefault()
    next()
  } else if (e.key === 'Escape') {
    e.preventDefault()
    close()
  }
}

let touchStartX = 0
let touchEndX = 0

function onTouchStart(e: TouchEvent) {
  touchStartX = e.changedTouches[0].screenX
}

function onTouchEnd(e: TouchEvent) {
  touchEndX = e.changedTouches[0].screenX
  const diff = touchEndX - touchStartX
  if (Math.abs(diff) > 50) {
    if (diff > 0) {
      prev()
    } else {
      next()
    }
  }
}

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      currentIndex.value = props.initialIndex
      window.addEventListener('keydown', onKeydown)
      nextTick(() => {
        closeButtonRef.value?.focus()
      })
    } else {
      window.removeEventListener('keydown', onKeydown)
    }
  }
)

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent
      fullScreen
      class="border-0 bg-black/95 p-0"
      @touchstart="onTouchStart"
      @touchend="onTouchEnd"
    >
      <button
        ref="closeButtonRef"
        type="button"
        aria-label="关闭"
        class="focus-ring touch-target absolute right-4 top-4 z-50 inline-flex h-10 w-10 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
        @click="close"
      >
        <X class="h-5 w-5" />
      </button>

      <button
        v-if="images.length > 1"
        type="button"
        aria-label="上一张"
        class="focus-ring touch-target absolute left-4 top-1/2 z-50 -translate-y-1/2 inline-flex h-12 w-12 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
        @click="prev"
      >
        <ChevronLeft class="h-6 w-6" />
      </button>

      <button
        v-if="images.length > 1"
        type="button"
        aria-label="下一张"
        class="focus-ring touch-target absolute right-4 top-1/2 z-50 -translate-y-1/2 inline-flex h-12 w-12 items-center justify-center rounded-full bg-white/10 text-white backdrop-blur-sm transition-colors hover:bg-white/20"
        @click="next"
      >
        <ChevronRight class="h-6 w-6" />
      </button>

      <div class="flex h-full w-full items-center justify-center p-4">
        <img
          v-if="images[currentIndex]"
          :src="getImageUrl(images[currentIndex].file_id)"
          :alt="`图片 ${currentIndex + 1} / ${images.length}`"
          class="motion-reduce max-h-full max-w-full object-contain animate-scale-in"
        />
      </div>

      <div class="pb-safe absolute bottom-6 left-0 right-0 z-50 text-center text-sm text-white/80">
        {{ currentIndex + 1 }} / {{ images.length }}
      </div>
    </DialogContent>
  </Dialog>
</template>
