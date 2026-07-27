import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

const apiTarget = process.env.VITE_API_TARGET ?? 'http://localhost:4300'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: { '@': '/src' },
  },
  server: {
    host: true,
    port: 5173,
    // Redundant with CHOKIDAR_USEPOLLING: guarantees polling even outside the dev container.
    watch: { usePolling: true },
    proxy: {
      '/api': {
        target: apiTarget,
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
      '/ws': {
        target: apiTarget,
        ws: true,
        rewrite: (path) => path.replace(/^\/ws/, ''),
      },
    },
  },
})
