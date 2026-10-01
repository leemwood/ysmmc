<script setup lang="ts">
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useNexusmcStore } from '@/stores/nexusmc'
import { NEXUSMC_REGISTER_URL, NEXUSMC_SITE } from '@/lib/nexusmc'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import NexusmcAvatar from '@/components/NexusmcAvatar.vue'
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuLabel,
} from '@/components/ui/dropdown-menu'
import { Menu, X, Home, LogOut, UserPlus, ExternalLink, Sparkles, UserRound, Search } from 'lucide-vue-next'
import { ref, watch } from 'vue'
import ThemeToggle from '@/components/ThemeToggle.vue'

const store = useNexusmcStore()
const route = useRoute()
const router = useRouter()
const isMenuOpen = ref(false)
// 顶栏搜索框与首页路由 ?q= 双向同步
const searchInput = ref(String(route.query.q || ''))

watch(() => route.query.q, (q) => {
  searchInput.value = String(q || '')
})

function submitSearch() {
  const q = searchInput.value.trim()
  // 任何页面搜索都回首页并带上关键词；首页的列表通过路由 watch 响应
  void router.push(q ? { path: '/', query: { q } } : { path: '/' })
  isMenuOpen.value = false
}

function handleLogout() {
  store.logout()
  isMenuOpen.value = false
}

function closeMenu() {
  isMenuOpen.value = false
}

function isActive(path: string) {
  return route.path === path
}
</script>

<template>
  <nav class="sticky top-0 z-40 border-b bg-background/80 backdrop-blur-md supports-[backdrop-filter]:bg-background/60">
    <div class="container-app">
      <div class="flex h-16 items-center justify-between">
        <RouterLink to="/" class="text-xl font-bold gradient-text focus-ring rounded-md">
          YSM 模型站
        </RouterLink>

        <!-- Desktop nav -->
        <div class="hidden md:flex md:items-center md:gap-1">
          <RouterLink
            to="/"
            class="inline-flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors focus-ring"
            :class="
              isActive('/')
                ? 'bg-accent text-primary'
                : 'text-muted-foreground hover:bg-accent hover:text-foreground'
            "
          >
            <Home class="h-4 w-4" />
            资源广场
          </RouterLink>

          <!-- 顶栏搜索（桌面端） -->
          <form class="relative hidden md:block" @submit.prevent="submitSearch">
            <Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              v-model="searchInput"
              placeholder="搜索 NexusMC 资源…"
              class="h-9 w-44 pl-9 lg:w-56"
              maxlength="100"
              aria-label="搜索 NexusMC 资源"
            />
          </form>

          <!-- NexusMC 登录 / 用户信息 -->
          <template v-if="store.isLoggedIn">
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <Button
                  variant="ghost"
                  class="btn-press gap-2 rounded-full pl-1 pr-3 focus-ring"
                  aria-label="NexusMC 用户菜单"
                >
                  <NexusmcAvatar size="sm" class="border" />
                  <span class="max-w-32 truncate text-sm font-medium">{{ store.user?.username }}</span>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" class="w-52">
                <DropdownMenuLabel>
                  <div class="flex flex-col">
                    <span>{{ store.user?.username }}</span>
                    <span class="text-xs font-normal text-muted-foreground">NexusMC 账号</span>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem as-child>
                  <RouterLink to="/me" class="w-full cursor-pointer">
                    <UserRound class="mr-2 h-4 w-4" />
                    个人主页
                  </RouterLink>
                </DropdownMenuItem>
                <DropdownMenuItem as-child>
                  <a
                    :href="NEXUSMC_SITE"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="w-full cursor-pointer"
                  >
                    <ExternalLink class="mr-2 h-4 w-4" />
                    打开 NexusMC
                  </a>
                </DropdownMenuItem>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                  @click="handleLogout"
                  class="text-destructive focus:text-destructive"
                >
                  <LogOut class="mr-2 h-4 w-4" />
                  退出登录
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </template>

          <template v-else>
            <a :href="NEXUSMC_REGISTER_URL" target="_blank" rel="noopener noreferrer">
              <Button variant="outline" size="sm" class="btn-press focus-ring">
                <UserPlus class="mr-1.5 h-4 w-4" />
                注册
              </Button>
            </a>
            <Button size="sm" class="btn-press focus-ring" @click="store.login()">
              <Sparkles class="mr-1.5 h-4 w-4" />
              NexusMC 登录
            </Button>
          </template>

          <ThemeToggle />
        </div>

        <!-- Mobile menu trigger -->
        <div class="flex items-center md:hidden">
          <Button
            variant="ghost"
            size="icon"
            class="btn-press focus-ring"
            aria-label="打开菜单"
            :aria-expanded="isMenuOpen"
            aria-controls="mobile-menu"
            @click="isMenuOpen = !isMenuOpen"
          >
            <Menu v-if="!isMenuOpen" class="h-5 w-5" />
            <X v-else class="h-5 w-5" />
          </Button>
        </div>
      </div>
    </div>

    <!-- Mobile drawer -->
    <!-- Mobile drawer：Teleport 到 body，避免 nav 的 backdrop-filter 把
         fixed 定位变成相对导航栏，导致抽屉被压进 64px 的导航条里 -->
    <Teleport to="body">
    <Transition
      enter-active-class="transition ease-out duration-300"
      enter-from-class="translate-x-full"
      enter-to-class="translate-x-0"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="translate-x-0"
      leave-to-class="translate-x-full"
    >
      <div
        v-if="isMenuOpen"
        id="mobile-menu"
        class="md:hidden fixed inset-0 z-50"
      >
        <div
          class="absolute inset-0 bg-black/50 backdrop-blur-sm"
          @click="closeMenu"
        />
        <div
          class="absolute inset-y-0 right-0 flex w-3/4 max-w-sm flex-col bg-background/95 backdrop-blur shadow-2xl"
        >
          <div class="flex h-16 items-center justify-between border-b px-4">
            <span class="font-semibold">菜单</span>
            <Button
              variant="ghost"
              size="icon"
              class="btn-press focus-ring"
              aria-label="关闭菜单"
              @click="closeMenu"
            >
              <X class="h-5 w-5" />
            </Button>
          </div>

          <div class="flex-1 overflow-y-auto p-4 space-y-1">
            <!-- 移动抽屉内的搜索 -->
            <form class="mb-3 flex gap-2" @submit.prevent="submitSearch">
              <div class="relative flex-1">
                <Search class="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  v-model="searchInput"
                  placeholder="搜索 NexusMC 资源…"
                  class="pl-9"
                  maxlength="100"
                  aria-label="搜索 NexusMC 资源"
                />
              </div>
              <Button type="submit" variant="secondary" size="sm" class="btn-press focus-ring shrink-0">搜索</Button>
            </form>

            <RouterLink
              to="/"
              class="flex items-center gap-3 rounded-lg px-3 py-3 text-base font-medium transition-colors focus-ring"
              :class="isActive('/') ? 'bg-accent text-primary' : 'text-foreground hover:bg-accent'"
              @click="closeMenu"
            >
              <Home class="h-5 w-5" />
              资源广场
            </RouterLink>

            <!-- 用户信息 -->
            <div v-if="store.isLoggedIn" class="mt-4 space-y-1 border-t pt-4">
              <div class="flex items-center gap-3 rounded-lg px-3 py-3">
                <NexusmcAvatar size="md" class="border" />
                <div class="min-w-0">
                  <div class="truncate text-sm font-medium">{{ store.user?.username }}</div>
                  <div class="text-xs text-muted-foreground">NexusMC 账号</div>
                </div>
              </div>
              <RouterLink
                to="/me"
                class="flex items-center gap-3 rounded-lg px-3 py-3 text-base font-medium transition-colors hover:bg-accent focus-ring"
                @click="closeMenu"
              >
                <UserRound class="h-5 w-5" />
                个人主页
              </RouterLink>
              <a
                :href="NEXUSMC_SITE"
                target="_blank"
                rel="noopener noreferrer"
                class="flex items-center gap-3 rounded-lg px-3 py-3 text-base font-medium transition-colors hover:bg-accent focus-ring"
                @click="closeMenu"
              >
                <ExternalLink class="h-5 w-5" />
                打开 NexusMC
              </a>
              <button
                class="flex w-full items-center gap-3 rounded-lg px-3 py-3 text-left text-base font-medium text-destructive transition-colors hover:bg-destructive/10 focus-ring"
                @click="handleLogout"
              >
                <LogOut class="h-5 w-5" />
                退出登录
              </button>
            </div>

            <template v-else>
              <div class="mt-4 space-y-1 border-t pt-4">
                <a
                  :href="NEXUSMC_REGISTER_URL"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="flex items-center gap-3 rounded-lg px-3 py-3 text-base font-medium text-primary transition-colors hover:bg-primary/10 focus-ring"
                  @click="closeMenu"
                >
                  <UserPlus class="h-5 w-5" />
                  注册 NexusMC 账号
                </a>
                <button
                  class="flex w-full items-center gap-3 rounded-lg px-3 py-3 text-left text-base font-medium text-primary transition-colors hover:bg-primary/10 focus-ring"
                  @click="store.login(); closeMenu()"
                >
                  <Sparkles class="h-5 w-5" />
                  NexusMC 登录
                </button>
              </div>
            </template>
          </div>

          <div class="border-t p-4 pb-safe">
            <div class="flex items-center justify-between">
              <span class="text-sm text-muted-foreground">外观主题</span>
              <ThemeToggle />
            </div>
          </div>
        </div>
      </div>
    </Transition>
    </Teleport>
  </nav>
</template>
