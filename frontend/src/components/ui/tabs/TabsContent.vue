<script setup lang="ts">
import { TabsContent } from 'radix-vue'
import { cn } from '@/lib/utils'
import type { TabsContentProps } from 'radix-vue'

const props = defineProps<TabsContentProps & { class?: string }>()
</script>

<template>
  <TabsContent
    v-bind="props"
    :class="cn(
      'mt-2 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
      props.class
    )"
  >
    <slot />
  </TabsContent>
</template>
