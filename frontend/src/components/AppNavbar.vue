<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { Menu, X, Upload, User, LogOut, LayoutDashboard } from 'lucide-vue-next'
import { ref } from 'vue'
import { getAvatarUrl } from '@/utils/image'
import ThemeToggle from '@/components/ThemeToggle.vue'

const authStore = useAuthStore()
const isMenuOpen = ref(false)

function handleLogout() {
  authStore.logout()
  window.location.href = '/'
}

function closeMenu() {
  isMenuOpen.value = false
}
</script>

<template>
  <nav class="sticky top-0 z-50 border-b bg-background/80 backdrop-blur-md supports-[backdrop-filter]:bg-background/60">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="flex h-16 items-center justify-between">
        <div class="flex items-center">
          <RouterLink to="/" class="text-xl font-bold gradient-text">
            YSM 模型站
          </RouterLink>
        </div>

        <div class="hidden md:flex md:items-center md:gap-2">
          <RouterLink 
            to="/" 
            class="rounded-md px-3 py-2 text-sm font-medium text-muted-foreground transition-colors hover:text-foreground hover:bg-accent"
          >
            首页
          </RouterLink>

          <template v-if="authStore.isAuthenticated">
            <RouterLink to="/upload">
              <Button variant="outline" size="sm" class="btn-press">
                <Upload class="mr-2 h-4 w-4" />
                上传模型
              </Button>
            </RouterLink>

            <RouterLink to="/profile">
              <Button variant="ghost" size="icon" class="btn-press">
                <Avatar class="h-8 w-8">
                  <AvatarImage :src="getAvatarUrl(authStore.user?.avatar_id, authStore.user?.avatar_url) || undefined">
                    <User class="h-4 w-4" />
                  </AvatarImage>
                </Avatar>
              </Button>
            </RouterLink>

            <RouterLink v-if="authStore.isAdmin" to="/admin">
              <Button variant="ghost" size="icon" class="btn-press">
                <LayoutDashboard class="h-4 w-4" />
              </Button>
            </RouterLink>

            <Button variant="ghost" size="icon" class="btn-press" @click="handleLogout">
              <LogOut class="h-4 w-4" />
            </Button>
          </template>

          <template v-else>
            <RouterLink to="/login">
              <Button variant="outline" size="sm" class="btn-press">登录</Button>
            </RouterLink>
          </template>
          
          <ThemeToggle />
        </div>

        <div class="flex items-center gap-2 md:hidden">
          <ThemeToggle />
          <Button variant="ghost" size="icon" @click="isMenuOpen = !isMenuOpen" class="btn-press">
            <Menu v-if="!isMenuOpen" class="h-5 w-5" />
            <X v-else class="h-5 w-5" />
          </Button>
        </div>
      </div>
    </div>

    <Transition
      enter-active-class="transition ease-out duration-200"
      enter-from-class="opacity-0 -translate-y-2"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition ease-in duration-150"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-2"
    >
      <div v-if="isMenuOpen" class="md:hidden border-t bg-background/95 backdrop-blur">
        <div class="mx-auto max-w-7xl px-4 py-4 space-y-1">
          <RouterLink 
            to="/" 
            class="block rounded-md px-3 py-2 text-base font-medium text-foreground hover:bg-accent transition-colors"
            @click="closeMenu"
          >
            首页
          </RouterLink>

          <template v-if="authStore.isAuthenticated">
            <RouterLink 
              to="/upload" 
              class="block rounded-md px-3 py-2 text-base font-medium text-foreground hover:bg-accent transition-colors"
              @click="closeMenu"
            >
              上传模型
            </RouterLink>
            <RouterLink 
              to="/profile" 
              class="block rounded-md px-3 py-2 text-base font-medium text-foreground hover:bg-accent transition-colors"
              @click="closeMenu"
            >
              个人中心
            </RouterLink>
            <RouterLink 
              v-if="authStore.isAdmin" 
              to="/admin" 
              class="block rounded-md px-3 py-2 text-base font-medium text-foreground hover:bg-accent transition-colors"
              @click="closeMenu"
            >
              管理后台
            </RouterLink>
            <button 
              class="w-full text-left rounded-md px-3 py-2 text-base font-medium text-destructive hover:bg-destructive/10 transition-colors"
              @click="handleLogout(); closeMenu()"
            >
              退出登录
            </button>
          </template>

          <template v-else>
            <RouterLink 
              to="/login" 
              class="block rounded-md px-3 py-2 text-base font-medium text-primary hover:bg-primary/10 transition-colors"
              @click="closeMenu"
            >
              登录
            </RouterLink>
          </template>
        </div>
      </div>
    </Transition>
  </nav>
</template>
