import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';
import path from 'path';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  optimizeDeps: {
    exclude: ['@urql/svelte'],
  },
  publicDir: 'static',
  resolve: {
    alias: {
      $wailsjs: path.resolve('./wailsjs'),
      $lib: path.resolve('./src/lib'),
    }
  },
  build: {
    outDir: 'build'
  }
});
