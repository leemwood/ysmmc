<script setup lang="ts">
import { ArrowUpRight, Box, ChevronDown, LogIn, LogOut, Search, UserRound } from 'lucide-vue-next'
import { useRoute } from 'vue-router'
import { useNexusmcStore } from '@/stores/nexusmc'
import { NEXUSMC_REGISTER_URL, NEXUSMC_SITE } from '@/lib/nexusmc'
import NexusmcAvatar from '@/components/NexusmcAvatar.vue'
import ThemeToggle from '@/components/ThemeToggle.vue'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator } from '@/components/ui/dropdown-menu'

const store = useNexusmcStore()
const route = useRoute()

function focusHomeSearch() {
  if (route.name !== 'home') return
  const input = document.getElementById('model-search')
  input?.scrollIntoView({ block: 'center' })
  input?.focus({ preventScroll: true })
}
</script>

<template>
  <header class="site-header sticky top-0 z-40 border-b border-border">
    <div class="container-app flex h-(--header-height) items-center justify-between gap-2.5 sm:gap-4">
      <RouterLink to="/" class="inline-flex shrink-0 items-center gap-2.5 rounded-md font-bold" aria-label="YSM 模型站首页">
        <span class="grid size-9 place-items-center rounded-[0.625rem] bg-primary text-primary-foreground">
          <Box class="size-[1.35rem]" :stroke-width="2.25" aria-hidden="true" />
        </span>
        <span class="text-[1.0625rem] tracking-tight sm:text-lg">YSM 模型站</span>
      </RouterLink>

      <nav aria-label="主导航" class="hidden grow items-center justify-center gap-7 text-sm font-semibold md:flex">
        <RouterLink
          to="/"
          class="rounded-sm text-muted-foreground hover:text-primary"
          :class="{ 'text-primary!': route.name === 'home' || route.name === 'nexusmc-resource' }"
        >模型库</RouterLink>
        <a :href="NEXUSMC_SITE" target="_blank" rel="noopener noreferrer" class="inline-flex items-center gap-1 rounded-sm text-muted-foreground hover:text-primary">
          NexusMC <ArrowUpRight class="size-3.5" aria-hidden="true" />
        </a>
      </nav>

      <div class="flex shrink-0 items-center gap-0.5 sm:gap-2">
        <RouterLink
          :to="{ path: '/', query: route.name === 'home' ? route.query : {}, hash: '#search' }"
          class="site-button site-button-ghost site-button-icon"
          aria-label="搜索模型"
          @click="focusHomeSearch"
        ><Search class="size-[1.125rem]" aria-hidden="true" /></RouterLink>
        <ThemeToggle />
        <DropdownMenu>
          <DropdownMenuTrigger>
            <button
              type="button"
              class="site-button site-button-outline"
              :class="store.isLoggedIn ? 'gap-1 px-1.5 py-1' : 'px-2.5 sm:px-3.5'"
              :aria-label="store.isLoggedIn ? '账号菜单' : '登录或注册'"
            >
              <template v-if="store.isLoggedIn">
                <NexusmcAvatar size="sm" />
                <ChevronDown class="size-3.5 text-muted-foreground" aria-hidden="true" />
              </template>
              <template v-else>
                <LogIn class="hidden size-4 sm:block" aria-hidden="true" />
                登录
              </template>
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" :side-offset="8" class="enter-content w-48 rounded-xl p-1.5 max-sm:[&_[role=menuitem]]:min-h-11">
            <template v-if="store.isLoggedIn">
              <DropdownMenuItem as-child>
                <RouterLink to="/me" class="gap-2"><UserRound class="size-4" />个人主页</RouterLink>
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem class="gap-2" @select="store.logout()"><LogOut class="size-4" />退出登录</DropdownMenuItem>
            </template>
            <template v-else>
              <DropdownMenuItem class="gap-2" @select="store.login(route.name === 'nexusmc-callback' ? '/' : route.fullPath)">
                <LogIn class="size-4" />通过 NexusMC 登录
              </DropdownMenuItem>
              <DropdownMenuItem as-child>
                <a :href="NEXUSMC_REGISTER_URL" target="_blank" rel="noopener noreferrer" class="gap-2">
                  <ArrowUpRight class="size-4" />注册 NexusMC 账号
                </a>
              </DropdownMenuItem>
            </template>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  </header>
</template>

<style scoped>
.site-header {
  background: hsl(var(--background) / 0.92);
  backdrop-filter: blur(12px);
}
</style>
