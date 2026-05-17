import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { fileURLToPath, URL } from 'node:url'

// Если frontend работает в Docker, target должен быть host.docker.internal.
// Если frontend запускается локально npm run dev, можно поставить VITE_DEV_PROXY_TARGET=http://127.0.0.1:8080.
const backendTarget = process.env.VITE_DEV_PROXY_TARGET || 'http://host.docker.internal:8080'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    strictPort: true,
    proxy: {
      '/api': {
        target: backendTarget,
        changeOrigin: true,
        secure: false,
      },
    },
  },
})
