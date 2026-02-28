<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { Menu, X, Upload, User, LogOut, LayoutDashboard } from 'lucide-vue-next'
import { ref } from 'vue'

const authStore = useAuthStore()
const isMenuOpen = ref(false)

function handleLogout() {
  authStore.logout()
  window.location.href = '/'
}
</script>

<template>
  <nav class="border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
    <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="flex h-16 items-center justify-between">
        <div class="flex items-center">
          <RouterLink to="/" class="text-xl font-bold text-primary">
            YSM 模型站
          </RouterLink>
        </div>

        <div class="hidden md:flex md:items-center md:space-x-4">
          <RouterLink to="/" class="text-sm font-medium text-muted-foreground hover:text-foreground">
            首页
          </RouterLink>

          <template v-if="authStore.isAuthenticated">
            <RouterLink to="/upload">
              <Button variant="outline" size="sm">
                <Upload class="mr-2 h-4 w-4" />
                上传模型
              </Button>
            </RouterLink>

            <RouterLink to="/profile">
              <Button variant="ghost" size="icon">
                <Avatar class="h-8 w-8">
                  <AvatarImage :src="authStore.user?.avatar_url || undefined">
                    <User class="h-4 w-4" />
                  </AvatarImage>
                </Avatar>
              </Button>
            </RouterLink>

            <RouterLink v-if="authStore.isAdmin" to="/admin">
              <Button variant="ghost" size="icon">
                <LayoutDashboard class="h-4 w-4" />
              </Button>
            </RouterLink>

            <Button variant="ghost" size="icon" @click="handleLogout">
              <LogOut class="h-4 w-4" />
            </Button>
          </template>

          <template v-else>
            <RouterLink to="/login">
              <Button variant="outline" size="sm">登录</Button>
            </RouterLink>
          </template>
        </div>

        <div class="md:hidden">
          <Button variant="ghost" size="icon" @click="isMenuOpen = !isMenuOpen">
            <Menu v-if="!isMenuOpen" class="h-5 w-5" />
            <X v-else class="h-5 w-5" />
          </Button>
        </div>
      </div>

      <div v-if="isMenuOpen" class="md:hidden border-t py-4">
        <div class="flex flex-col space-y-4">
          <RouterLink to="/" class="text-sm font-medium" @click="isMenuOpen = false">
            首页
          </RouterLink>

          <template v-if="authStore.isAuthenticated">
            <RouterLink to="/upload" class="text-sm font-medium" @click="isMenuOpen = false">
              上传模型
            </RouterLink>
            <RouterLink to="/profile" class="text-sm font-medium" @click="isMenuOpen = false">
              个人中心
            </RouterLink>
            <RouterLink v-if="authStore.isAdmin" to="/admin" class="text-sm font-medium" @click="isMenuOpen = false">
              管理后台
            </RouterLink>
            <Button variant="outline" size="sm" @click="handleLogout">
              退出登录
            </Button>
          </template>

          <template v-else>
            <RouterLink to="/login" @click="isMenuOpen = false">
              <Button variant="outline" size="sm">登录</Button>
            </RouterLink>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>
