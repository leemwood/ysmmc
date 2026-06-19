<script setup lang="ts">
import { RouterLink, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuTrigger,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuLabel,
} from '@/components/ui/dropdown-menu'
import {
  Menu,
  X,
  Home,
  Upload,
  User,
  LogOut,
  LayoutDashboard,
  UserPlus,
} from 'lucide-vue-next'
import { ref } from 'vue'
import { getAvatarUrl } from '@/utils/image'
import ThemeToggle from '@/components/ThemeToggle.vue'

const authStore = useAuthStore()
const route = useRoute()
const isMenuOpen = ref(false)

function handleLogout() {
  authStore.logout()
  window.location.href = '/'
}

function closeMenu() {
  isMenuOpen.value = false
}

function isActive(path: string) {
  return route.path === path
}

interface NavLink {
  to: string
  label: string
  icon: typeof Home
  admin?: boolean
}

const publicLinks: NavLink[] = [
  { to: '/', label: '首页', icon: Home },
]

const authLinks: NavLink[] = [
  { to: '/upload', label: '上传模型', icon: Upload },
  { to: '/profile', label: '个人中心', icon: User },
  { to: '/admin', label: '管理后台', icon: LayoutDashboard, admin: true },
]

const mobileLinks = publicLinks.concat(authLinks)
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
            v-for="link in publicLinks"
            :key="link.to"
            :to="link.to"
            class="inline-flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors focus-ring"
            :class="
              isActive(link.to)
                ? 'bg-accent text-primary'
                : 'text-muted-foreground hover:bg-accent hover:text-foreground'
            "
          >
            <component :is="link.icon" class="h-4 w-4" />
            {{ link.label }}
          </RouterLink>

          <RouterLink
            v-if="authStore.isAuthenticated"
            to="/upload"
            class="inline-flex items-center gap-2 rounded-md px-3 py-2 text-sm font-medium transition-colors focus-ring"
            :class="
              isActive('/upload')
                ? 'bg-accent text-primary'
                : 'text-muted-foreground hover:bg-accent hover:text-foreground'
            "
          >
            <Upload class="h-4 w-4" />
            上传模型
          </RouterLink>

          <template v-if="authStore.isAuthenticated">
            <DropdownMenu>
              <DropdownMenuTrigger as-child>
                <Button
                  variant="ghost"
                  size="icon"
                  class="btn-press rounded-full focus-ring"
                  aria-label="用户菜单"
                >
                  <Avatar class="h-8 w-8">
                    <AvatarImage
                      v-if="getAvatarUrl(authStore.user?.avatar_id, authStore.user?.avatar_url)"
                      :src="getAvatarUrl(authStore.user?.avatar_id, authStore.user?.avatar_url) || undefined"
                      :alt="authStore.user?.username || '用户头像'"
                    />
                    <span v-else class="flex h-full w-full items-center justify-center rounded-full bg-muted">
                      <User class="h-4 w-4" />
                    </span>
                  </Avatar>
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end" class="w-48">
                <DropdownMenuLabel>
                  {{ authStore.user?.username || '我的账号' }}
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem @click="$router.push('/profile')">
                  <User class="mr-2 h-4 w-4" />
                  个人中心
                </DropdownMenuItem>
                <DropdownMenuItem
                  v-if="authStore.isAdmin"
                  @click="$router.push('/admin')"
                >
                  <LayoutDashboard class="mr-2 h-4 w-4" />
                  管理后台
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
            <RouterLink to="/login">
              <Button variant="outline" size="sm" class="btn-press focus-ring">
                登录
              </Button>
            </RouterLink>
            <RouterLink to="/login?mode=register">
              <Button size="sm" class="btn-press focus-ring">
                <UserPlus class="mr-1.5 h-4 w-4" />
                注册
              </Button>
            </RouterLink>
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
            <RouterLink
              v-for="link in mobileLinks"
              :key="link.to"
              v-show="
                !['/profile', '/admin', '/upload'].includes(link.to) ||
                (authStore.isAuthenticated && (!link.admin || authStore.isAdmin))
              "
              :to="link.to"
              class="flex items-center gap-3 rounded-lg px-3 py-3 text-base font-medium transition-colors focus-ring"
              :class="
                isActive(link.to)
                  ? 'bg-accent text-primary'
                  : 'text-foreground hover:bg-accent'
              "
              @click="closeMenu"
            >
              <component :is="link.icon" class="h-5 w-5" />
              {{ link.label }}
            </RouterLink>

            <template v-if="!authStore.isAuthenticated">
              <div class="mt-4 space-y-1 border-t pt-4">
                <RouterLink
                  to="/login"
                  class="flex items-center gap-3 rounded-lg px-3 py-3 text-base font-medium text-primary transition-colors hover:bg-primary/10 focus-ring"
                  @click="closeMenu"
                >
                  <User class="h-5 w-5" />
                  登录
                </RouterLink>
                <RouterLink
                  to="/login?mode=register"
                  class="flex items-center gap-3 rounded-lg px-3 py-3 text-base font-medium text-primary transition-colors hover:bg-primary/10 focus-ring"
                  @click="closeMenu"
                >
                  <UserPlus class="h-5 w-5" />
                  注册
                </RouterLink>
              </div>
            </template>

            <button
              v-else
              class="mt-4 flex w-full items-center gap-3 rounded-lg border-t px-3 py-3 text-left text-base font-medium text-destructive transition-colors hover:bg-destructive/10 focus-ring"
              @click="handleLogout(); closeMenu()"
            >
              <LogOut class="h-5 w-5" />
              退出登录
            </button>
          </div>

          <div class="border-t p-4 safe-area-inset">
            <div class="flex items-center justify-between">
              <span class="text-sm text-muted-foreground">外观主题</span>
              <ThemeToggle />
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </nav>
</template>
