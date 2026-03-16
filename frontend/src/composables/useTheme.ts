import { ref, watch } from 'vue'
import { usePreferredColorScheme, useStorage } from '@vueuse/core'

export type Theme = 'light' | 'dark' | 'system'

const THEME_KEY = 'ysm-theme'

const storedTheme = useStorage<Theme>(THEME_KEY, 'system')
const preferredColorScheme = usePreferredColorScheme()

function applyTheme(theme: Theme) {
  const root = document.documentElement
  const isDark = theme === 'dark' || (theme === 'system' && preferredColorScheme.value === 'dark')
  
  if (isDark) {
    root.classList.add('dark')
  } else {
    root.classList.remove('dark')
  }
}

export function useTheme() {
  const theme = ref<Theme>(storedTheme.value)

  function setTheme(newTheme: Theme) {
    theme.value = newTheme
    storedTheme.value = newTheme
    applyTheme(newTheme)
  }

  function toggleTheme() {
    const currentIsDark = document.documentElement.classList.contains('dark')
    setTheme(currentIsDark ? 'light' : 'dark')
  }

  watch(preferredColorScheme, () => {
    if (theme.value === 'system') {
      applyTheme('system')
    }
  })

  return {
    theme,
    setTheme,
    toggleTheme,
  }
}
