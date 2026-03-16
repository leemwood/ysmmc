<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { Button } from '@/components/ui/button'
import { Download, X } from 'lucide-vue-next'

const showPrompt = ref(false)
const deferredPrompt = ref<any>(null)
const isIOS = ref(false)
const isStandalone = ref(false)

function checkIOS() {
  const isIOSDevice = /iPad|iPhone|iPod/.test(navigator.userAgent)
  const isInStandaloneMode = ('standalone' in window.navigator) && 
    (window.navigator as any).standalone
  isIOS.value = isIOSDevice && !isInStandaloneMode
}

function checkStandalone() {
  isStandalone.value = window.matchMedia('(display-mode: standalone)').matches ||
    (window.navigator as any).standalone === true
}

function handleBeforeInstallPrompt(e: Event) {
  e.preventDefault()
  deferredPrompt.value = e
  if (!isStandalone.value) {
    showPrompt.value = true
  }
}

async function installPWA() {
  if (!deferredPrompt.value) return
  
  showPrompt.value = false
  deferredPrompt.value.prompt()
  
  const { outcome } = await deferredPrompt.value.userChoice
  
  if (outcome === 'accepted') {
    console.log('PWA installed')
  }
  
  deferredPrompt.value = null
}

function dismissPrompt() {
  showPrompt.value = false
  localStorage.setItem('pwa-prompt-dismissed', 'true')
}

onMounted(() => {
  checkIOS()
  checkStandalone()
  
  window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
  
  const dismissed = localStorage.getItem('pwa-prompt-dismissed')
  if (dismissed) {
    showPrompt.value = false
  }
})

onUnmounted(() => {
  window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
})
</script>

<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition ease-out duration-300"
      enter-from-class="transform translate-y-full opacity-0"
      enter-to-class="transform translate-y-0 opacity-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="transform translate-y-0 opacity-100"
      leave-to-class="transform translate-y-full opacity-0"
    >
      <template v-if="showPrompt || isIOS">
        <div class="fixed inset-0 z-40 bg-black/80 backdrop-blur-sm" @click="dismissPrompt" />
        <div class="fixed bottom-4 left-4 right-4 sm:left-auto sm:right-4 sm:w-80 z-50">
          <div class="bg-card/95 backdrop-blur-md border rounded-lg shadow-lg p-4">
            <div class="flex items-start justify-between gap-2">
              <div class="flex-1">
                <h3 class="font-semibold text-sm">安装 YSM 模型站</h3>
                <p class="text-xs text-muted-foreground mt-1">
                  <template v-if="isIOS">
                    点击分享按钮，然后选择"添加到主屏幕"
                  </template>
                  <template v-else>
                    安装到主屏幕，获得更好的体验
                  </template>
                </p>
              </div>
              <Button
                variant="ghost"
                size="icon"
                class="h-6 w-6 flex-shrink-0"
                @click="dismissPrompt"
              >
                <X class="h-4 w-4" />
              </Button>
            </div>
            
            <div v-if="!isIOS" class="mt-3 flex gap-2">
              <Button size="sm" class="flex-1 btn-press" @click="installPWA">
                <Download class="mr-1 h-4 w-4" />
                安装
              </Button>
              <Button size="sm" variant="outline" @click="dismissPrompt">
                稍后
              </Button>
            </div>
          </div>
        </div>
      </template>
    </Transition>
  </Teleport>
</template>
