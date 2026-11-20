<script setup lang="ts">
import { Check, Monitor, Moon, Sun } from 'lucide-vue-next'
import { useTheme } from '@/composables/useTheme'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu'

const { theme, setTheme } = useTheme()
const options = [
  { value: 'light', label: '浅色', icon: Sun },
  { value: 'dark', label: '深色', icon: Moon },
  { value: 'system', label: '跟随系统', icon: Monitor },
] as const
</script>

<template>
  <DropdownMenu>
    <DropdownMenuTrigger>
      <button type="button" class="site-button site-button-ghost site-button-icon" aria-label="切换外观主题">
        <Sun v-if="theme === 'light'" class="size-[1.125rem]" aria-hidden="true" />
        <Moon v-else-if="theme === 'dark'" class="size-[1.125rem]" aria-hidden="true" />
        <Monitor v-else class="size-[1.125rem]" aria-hidden="true" />
      </button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="end" :side-offset="8" class="enter-content w-36 rounded-xl p-1.5">
      <DropdownMenuItem v-for="option in options" :key="option.value" class="gap-2 max-sm:min-h-11" @select="setTheme(option.value)">
        <component :is="option.icon" class="size-4" aria-hidden="true" />
        {{ option.label }}
        <Check v-if="theme === option.value" class="ml-auto size-4 text-primary" aria-label="当前主题" />
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>
