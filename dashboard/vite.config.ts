import {defineConfig, loadEnv} from 'vite'
import {svelte} from '@sveltejs/vite-plugin-svelte'

// https://vite.dev/config/
export default defineConfig(({mode}) => {
    // Load env file based on `mode`
    const env = loadEnv(mode, process.cwd(), '')

    return {
        plugins: [svelte()],
        server: {
            proxy: {
                '/api': {
                    target: env.API_URL || 'http://localhost:8080',
                    changeOrigin: true,
                    rewrite: (path) => path.replace(/^\/api/, '')
                }
            }
        }
    }
})
