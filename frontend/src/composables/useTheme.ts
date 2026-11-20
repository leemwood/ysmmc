import { watchEffect } from 'vue'
import { usePreferredColorScheme, useStorage } from '@vueuse/core'

export type Theme = 'light' | 'dark' | 'system'

export function useTheme() {
  const theme = useStorage<Theme>('ysm-theme', 'system')
  const preferredColorScheme = usePreferredColorScheme()

  watchEffect(() => {
    const dark = theme.value === 'dark' || (theme.value === 'system' && preferredColorScheme.value === 'dark')
    document.documentElement.classList.toggle('dark', dark)
  })

  return {
    theme,
    setTheme: (value: Theme) => { theme.value = value },
  }
}
