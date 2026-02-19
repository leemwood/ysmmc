import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'
import fs from 'fs'

function copyIndexToRoutes() {
  return {
    name: 'copy-index-to-routes',
    closeBundle() {
      const distDir = path.resolve(__dirname, 'dist')
      const indexPath = path.join(distDir, 'index.html')
      
      if (!fs.existsSync(indexPath)) return
      
      const indexContent = fs.readFileSync(indexPath, 'utf-8')
      
      const routes = [
        'login',
        'reset-password',
        'update-password',
        'verify-email',
        'verify-email-change',
        'upload',
        'profile',
        'admin',
        'admin/users',
      ]
      
      routes.forEach(route => {
        const routeDir = path.join(distDir, route)
        if (!fs.existsSync(routeDir)) {
          fs.mkdirSync(routeDir, { recursive: true })
        }
        fs.writeFileSync(path.join(routeDir, 'index.html'), indexContent)
      })
      
      console.log('Generated route HTML files')
    }
  }
}

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    copyIndexToRoutes()
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
          'ui-vendor': ['radix-vue', 'lucide-vue-next', 'class-variance-authority', 'clsx', 'tailwind-merge'],
        },
      },
    },
    chunkSizeWarningLimit: 1000,
  },
})
