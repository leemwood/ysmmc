<script setup lang="ts">
import { useRoute } from 'vue-router'
import { LayoutDashboard, Package, Users, Megaphone } from 'lucide-vue-next'

defineProps<{
  title?: string
  description?: string
}>()

const route = useRoute()

const navItems = [
  { label: '概览', to: '/admin', icon: LayoutDashboard },
  { label: '模型管理', to: '/admin/models', icon: Package },
  { label: '用户管理', to: '/admin/users', icon: Users },
  { label: '公告管理', to: '/admin/announcements', icon: Megaphone },
]

function isActive(to: string) {
  return route.path === to
}
</script>

<template>
  <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 sm:py-8 lg:px-8 space-y-6">
    <div v-if="title" class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
      <div class="min-w-0">
        <h1 class="text-2xl font-bold tracking-tight sm:text-3xl">{{ title }}</h1>
        <p v-if="description" class="mt-1 text-sm text-muted-foreground">{{ description }}</p>
      </div>
      <div class="flex w-full flex-col gap-2 sm:w-auto sm:flex-row sm:items-center">
        <slot name="actions"></slot>
      </div>
    </div>

    <nav class="-mx-4 px-4 overflow-x-auto">
      <div class="flex min-w-max gap-1 border-b sm:gap-4">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="flex min-h-[44px] min-w-[44px] items-center justify-center gap-2 border-b-2 px-3 py-2 text-sm font-medium transition-colors whitespace-nowrap sm:min-w-0 sm:justify-start sm:px-4"
          :class="isActive(item.to) ? 'border-primary text-primary' : 'border-transparent text-muted-foreground hover:text-foreground'"
        >
          <component :is="item.icon" class="h-4 w-4 flex-shrink-0" />
          <span class="hidden sm:inline">{{ item.label }}</span>
        </RouterLink>
      </div>
    </nav>

    <slot></slot>
  </div>
</template>
