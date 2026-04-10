import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import { resolve } from 'path';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());
  const apiTarget = env.VITE_API_TARGET || 'http://localhost:8000';

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src'),
      },
    },
    server: {
      host: true,
      port: Number(env.VITE_PORT) || 8002,
      proxy: {
        '/api': {
          target: apiTarget,
          changeOrigin: true,
        },
      },
      allowedHosts: [
      '80tiger.dpdns.org',
      // 如果你还有其他域名，也一并添加
    ],
    },
    build: {
      outDir: 'dist',
      target: 'es2015',
    },
  };
});
