import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import { transformWithOxc } from 'vite'

export default defineConfig({
  plugins: [
    {
      name: 'jsx-in-js',
      enforce: 'pre',
      async transform(code, id) {
        if (id.endsWith('.js') && !id.includes('/node_modules/')) {
          return transformWithOxc(code, id, {
            lang: 'jsx',
            jsx: { runtime: 'automatic' },
          })
        }
      },
    },
    react(),
  ],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./vitest.setup.js'],
  },
})
