import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { clearNexusmcUser, getNexusmcUser, type NexusmcUser } from '@/lib/nexusmc'

export const NEXUSMC_SITE = 'https://www.nexusmc.cn'
export const NEXUSMC_REGISTER_URL = `${NEXUSMC_SITE}/register`

// NexusMC 账号体系（OAuth2 登录，用户信息来自 /api/nexusmc/auth/callback）。
// 本站是 NexusMC 资源的展示与引流站，登录注册均由 NexusMC 承载。
export const useNexusmcStore = defineStore('nexusmc', () => {
  const user = ref<NexusmcUser | null>(getNexusmcUser())
  const isLoggedIn = computed(() => !!user.value)

  function login() {
    window.location.href = '/api/nexusmc/auth/start'
  }

  function logout() {
    clearNexusmcUser()
    user.value = null
  }

  return { user, isLoggedIn, login, logout }
})
