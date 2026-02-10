<script setup lang="ts">
import { Loader2 } from 'lucide-vue-next'

defineProps<{
  message?: string
}>()
</script>

<template>
  <div class="loading-container" role="status" aria-live="polite">
    <Loader2 class="animate-spin" :size="48" aria-hidden="true" />
    <p v-if="message">{{ message }}</p>
  </div>
</template>

<style lang="scss" scoped>
@use '../styles/themes/variables' as *;

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $spacing-3xl 0;
  color: var(--color-text-muted);
  gap: $spacing-md;

  .animate-spin {
    animation: spin 1s linear infinite;
    color: var(--color-primary);
  }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
