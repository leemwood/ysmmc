import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  clearNexusmcUser,
  getNexusmcUser,
  NEXUSMC_LOGIN_URL,
  type NexusmcUser,
} from '@/lib/nexusmc'

const RETURN_TO_KEY = 'nexusmc_return_to'

// NexusMC 账号体系（OAuth2 登录，用户信息来自 /api/nexusmc/auth/callback）。
// 本站是 NexusMC 资源的展示与引流站，登录注册均由 NexusMC 承载。
export const useNexusmcStore = defineStore('nexusmc', () => {
  const user = ref<NexusmcUser | null>(getNexusmcUser())
  const isLoggedIn = computed(() => !!user.value)
  // 展示用用户名，未登录或无用户名时回退为中性称呼
  const displayName = computed(() => user.value?.username || 'NexusMC 用户')

  function login(returnTo?: string) {
    // 记住登录前的目标路径，回调完成后由 CallbackView 读取并回跳
    if (returnTo && typeof window !== 'undefined') {
      window.sessionStorage.setItem(RETURN_TO_KEY, returnTo)
    }
    window.location.href = NEXUSMC_LOGIN_URL
  }

  function logout() {
    // 让服务端同时清掉 httpOnly token cookie（尽力而为，无需等待结果）
    if (typeof window !== 'undefined') {
      void fetch('/api/nexusmc/auth/logout', { method: 'POST' }).catch(() => {})
    }
    clearNexusmcUser()
    user.value = null
  }

  function consumeReturnTo(): string {
    if (typeof window === 'undefined') return ''
    const target = window.sessionStorage.getItem(RETURN_TO_KEY) || ''
    window.sessionStorage.removeItem(RETURN_TO_KEY)
    return target
  }

  return { user, isLoggedIn, displayName, login, logout, consumeReturnTo }
})
