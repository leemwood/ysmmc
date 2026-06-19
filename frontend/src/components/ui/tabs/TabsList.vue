<script setup lang="ts">
import { TabsList } from 'radix-vue'
import { cn } from '@/lib/utils'
import type { TabsListProps } from 'radix-vue'

const props = defineProps<TabsListProps & { class?: string }>()
</script>

<template>
  <TabsList
    v-bind="props"
    :class="cn(
      'inline-flex h-10 items-center justify-center rounded-lg bg-muted p-1 text-muted-foreground',
      props.class
    )"
  >
    <slot />
  </TabsList>
</template>
