// 升级时删除旧版本保存过的 API 响应；即使应用页面尚未打开也执行。
self.addEventListener('activate', (event) => {
  event.waitUntil(caches.delete('api-cache'))
})
