<script setup lang="ts">
import { cn } from '@/lib/utils'

interface Props {
  class?: string
  src?: string
  alt?: string
}

const props = defineProps<Props>()
</script>

<template>
  <img
    v-if="src"
    :src="src"
    :alt="alt"
    :class="cn('aspect-square h-full w-full', props.class)"
  />
  <span v-else :class="cn('flex h-full w-full items-center justify-center rounded-full bg-muted', props.class)">
    <slot />
  </span>
</template>
