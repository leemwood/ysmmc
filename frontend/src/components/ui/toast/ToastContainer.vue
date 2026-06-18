<script setup lang="ts">
import { useToast } from '@/composables/useToast'
import ToastItem from './ToastItem.vue'

const { toasts } = useToast()
</script>

<template>
  <Teleport to="body">
    <TransitionGroup
      tag="div"
      name="toast"
      class="fixed top-4 right-4 z-[100] flex flex-col gap-2 pointer-events-none sm:max-w-sm"
    >
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="pointer-events-auto"
      >
        <ToastItem :toast="toast" />
      </div>
    </TransitionGroup>
  </Teleport>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.toast-move {
  transition: transform 0.3s ease;
}
</style>
