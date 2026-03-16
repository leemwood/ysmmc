<script setup lang="ts">
import { Sun, Moon, Monitor } from 'lucide-vue-next'
import { useTheme, type Theme } from '@/composables/useTheme'
import { Button } from '@/components/ui/button'
import { ref } from 'vue'

const { theme, setTheme } = useTheme()
const showMenu = ref(false)

const themes: { value: Theme; label: string; icon: typeof Sun }[] = [
  { value: 'light', label: '浅色', icon: Sun },
  { value: 'dark', label: '深色', icon: Moon },
  { value: 'system', label: '系统', icon: Monitor },
]

function handleSelect(t: Theme) {
  setTheme(t)
  showMenu.value = false
}
</script>

<template>
  <div class="relative">
    <Button
      variant="ghost"
      size="icon"
      @click="showMenu = !showMenu"
      class="relative"
    >
      <Sun v-if="theme === 'light'" class="h-5 w-5" />
      <Moon v-else-if="theme === 'dark'" class="h-5 w-5" />
      <Monitor v-else class="h-5 w-5" />
    </Button>
    
    <Teleport to="body">
      <Transition
        enter-active-class="transition ease-out duration-100"
        enter-from-class="transform opacity-0"
        enter-to-class="transform opacity-100"
        leave-active-class="transition ease-in duration-75"
        leave-from-class="transform opacity-100"
        leave-to-class="transform opacity-0"
      >
        <template v-if="showMenu">
          <div class="fixed inset-0 z-40 bg-black/80 backdrop-blur-sm" @click="showMenu = false" />
          <div class="fixed z-50 bg-popover/95 backdrop-blur-md rounded-md shadow-lg ring-1 ring-black/5" 
               :style="{ top: '50%', left: '50%', transform: 'translate(-50%, -50%)' }">
            <div class="py-1 min-w-[120px]">
              <button
                v-for="t in themes"
                :key="t.value"
                class="flex w-full items-center gap-2 px-4 py-2.5 text-sm hover:bg-accent transition-colors"
                :class="theme === t.value ? 'text-primary bg-accent/50' : 'text-foreground'"
                @click="handleSelect(t.value)"
              >
                <component :is="t.icon" class="h-4 w-4" />
                {{ t.label }}
              </button>
            </div>
          </div>
        </template>
      </Transition>
    </Teleport>
  </div>
</template>
