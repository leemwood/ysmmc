<script setup lang="ts">
import { Box } from 'lucide-vue-next'
import { ref } from 'vue'

withDefaults(defineProps<{
  cover: string
  title: string
  iconClass?: string
  eager?: boolean
}>(), { iconClass: 'size-9' })

const failedSrc = ref('')
</script>

<template>
  <img
    v-if="cover && cover !== failedSrc"
    :key="cover"
    :src="cover"
    :alt="title"
    :loading="eager ? 'eager' : 'lazy'"
    decoding="async"
    class="resource-cover-img"
    @error="failedSrc = cover"
  />
  <div v-else class="grid h-full place-items-center text-muted-foreground" role="img" aria-label="无预览图">
    <Box :class="iconClass" class="opacity-60" aria-hidden="true" />
  </div>
</template>
