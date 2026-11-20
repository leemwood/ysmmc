<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, shallowRef } from 'vue'
import { Download, X } from 'lucide-vue-next'

interface InstallPromptEvent extends Event {
  prompt(): Promise<unknown>
}

const DISMISSED_KEY = 'pwa-prompt-dismissed'
const deferredPrompt = shallowRef<InstallPromptEvent>()
const isIOS = ref(false)
const dismissed = ref(false)
const standalone = ref(false)
const visible = computed(() => !standalone.value && !dismissed.value && (isIOS.value || deferredPrompt.value))

function capturePrompt(event: Event) {
  event.preventDefault()
  deferredPrompt.value = event as InstallPromptEvent
}

async function install() {
  const prompt = deferredPrompt.value
  deferredPrompt.value = undefined
  await prompt?.prompt().catch(() => {})
}

function dismiss() {
  dismissed.value = true
  try {
    localStorage.setItem(DISMISSED_KEY, 'true')
  } catch {
    // 隐私模式下仍可关闭当前提示。
  }
}

function handleInstalled() {
  standalone.value = true
  deferredPrompt.value = undefined
}

onMounted(() => {
  const iosStandalone = (navigator as Navigator & { standalone?: boolean }).standalone
  isIOS.value = /iPad|iPhone|iPod/.test(navigator.userAgent)
  standalone.value = iosStandalone === true || window.matchMedia('(display-mode: standalone)').matches
  try {
    dismissed.value = localStorage.getItem(DISMISSED_KEY) === 'true'
  } catch {
    // 无存储权限时保留当前会话的状态。
  }
  window.addEventListener('beforeinstallprompt', capturePrompt)
  window.addEventListener('appinstalled', handleInstalled)
})

onUnmounted(() => {
  window.removeEventListener('beforeinstallprompt', capturePrompt)
  window.removeEventListener('appinstalled', handleInstalled)
})
</script>

<template>
  <aside
    v-if="visible"
    class="site-panel fixed inset-x-4 bottom-[max(1rem,env(safe-area-inset-bottom))] z-40 flex items-start gap-3 p-4 shadow-lg sm:left-auto sm:w-88"
    aria-labelledby="install-title"
  >
    <div class="status-icon"><Download class="size-5" aria-hidden="true" /></div>
    <div class="min-w-0 flex-1">
      <h2 id="install-title" class="font-bold">安装 YSM 模型站</h2>
      <p v-if="isIOS" class="mt-1 text-sm text-muted-foreground">点击浏览器的共享按钮，再选择「添加到主屏幕」。</p>
      <button v-else type="button" class="site-button site-button-primary mt-2.5" @click="install">添加到主屏幕</button>
    </div>
    <button type="button" class="site-button site-button-ghost site-button-icon size-8 min-h-8" aria-label="关闭安装提示" @click="dismiss">
      <X class="size-4" aria-hidden="true" />
    </button>
  </aside>
</template>
