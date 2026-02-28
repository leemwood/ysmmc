<script setup lang="ts">
import { RouterLink } from 'vue-router'
import type { Model } from '@/types'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Download } from 'lucide-vue-next'

interface Props {
  model: Model
}

const props = defineProps<Props>()

function formatDate(date: string) {
  return new Date(date).toLocaleDateString('zh-CN')
}
</script>

<template>
  <RouterLink :to="`/model/${model.id}`">
    <Card class="overflow-hidden transition-shadow hover:shadow-lg">
      <div class="aspect-video w-full bg-muted">
        <img
          v-if="model.image_url"
          :src="model.image_url"
          :alt="model.title"
          class="h-full w-full object-cover"
        />
        <div v-else class="flex h-full items-center justify-center text-muted-foreground">
          无预览图
        </div>
      </div>
      <CardContent class="p-4">
        <h3 class="font-semibold line-clamp-1">{{ model.title }}</h3>
        <p class="mt-1 text-sm text-muted-foreground line-clamp-2">
          {{ model.description || '暂无描述' }}
        </p>
        <div class="mt-3 flex flex-wrap gap-1">
          <Badge v-for="tag in model.tags?.slice(0, 3)" :key="tag" variant="secondary" class="text-xs">
            {{ tag }}
          </Badge>
        </div>
        <div class="mt-3 flex items-center justify-between text-xs text-muted-foreground">
          <div class="flex items-center gap-3">
            <span class="flex items-center gap-1">
              <Download class="h-3 w-3" />
              {{ model.downloads }}
            </span>
          </div>
          <span>{{ formatDate(model.created_at) }}</span>
        </div>
      </CardContent>
    </Card>
  </RouterLink>
</template>
