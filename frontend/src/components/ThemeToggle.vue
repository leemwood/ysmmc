<script setup lang="ts">
import { Sun, Moon, Monitor } from 'lucide-vue-next'
import { useTheme, type Theme } from '@/composables/useTheme'
import { Button } from '@/components/ui/button'
import { ref, onMounted, onUnmounted, nextTick } from 'vue'

const { theme, setTheme } = useTheme()
const showMenu = ref(false)
const wrapperRef = ref<HTMLElement | null>(null)
const menuStyle = ref<Record<string, string>>({})

const themes: { value: Theme; label: string; icon: typeof Sun }[] = [
  { value: 'light', label: '浅色', icon: Sun },
  { value: 'dark', label: '深色', icon: Moon },
  { value: 'system', label: '系统', icon: Monitor },
]

function updateMenuPosition() {
  if (!wrapperRef.value) return
  const rect = wrapperRef.value.getBoundingClientRect()
  const nav = wrapperRef.value.closest('nav')
  const navBottom = nav ? nav.getBoundingClientRect().bottom : 0
  menuStyle.value = {
    position: 'fixed',
    top: `${Math.max(rect.bottom + 8, navBottom + 8)}px`,
    right: `${window.innerWidth - rect.right}px`,
    zIndex: '100',
  }
}

async function toggleMenu() {
  showMenu.value = !showMenu.value
  if (showMenu.value) {
    await nextTick()
    updateMenuPosition()
  }
}

function handleSelect(t: Theme) {
  setTheme(t)
  showMenu.value = false
}

function handleClickOutside(e: MouseEvent) {
  const target = e.target as Node
  if (wrapperRef.value?.contains(target)) return
  const menu = document.getElementById('theme-menu')
  if (menu && !menu.contains(target)) {
    showMenu.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  window.addEventListener('scroll', () => { if (showMenu.value) showMenu.value = false }, true)
  window.addEventListener('resize', () => { if (showMenu.value) showMenu.value = false })
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div ref="wrapperRef">
    <Button
      variant="ghost"
      size="icon"
      @click="toggleMenu"
    >
      <Sun v-if="theme === 'light'" class="h-5 w-5" />
      <Moon v-else-if="theme === 'dark'" class="h-5 w-5" />
      <Monitor v-else class="h-5 w-5" />
    </Button>

    <Teleport to="body">
      <Transition
        enter-active-class="transition ease-out duration-100"
        enter-from-class="transform opacity-0 scale-95"
        enter-to-class="transform opacity-100 scale-100"
        leave-active-class="transition ease-in duration-75"
        leave-from-class="transform opacity-100 scale-100"
        leave-to-class="transform opacity-0 scale-95"
      >
        <div
          v-if="showMenu"
          id="theme-menu"
          class="min-w-[120px] bg-popover/95 backdrop-blur-md rounded-md shadow-lg ring-1 ring-black/5 overflow-hidden"
          :style="menuStyle"
        >
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
      </Transition>
    </Teleport>
  </div>
</template>
