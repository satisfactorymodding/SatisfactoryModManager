import path from 'path';

import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    svelte({
      hot: true,
      onwarn: (warning, defaultHandler) => {
        if (warning.code === 'a11y-click-events-have-key-events') {
          return;
        }
        if (defaultHandler) {
          defaultHandler(warning);
        }
      },
    }),
  ],
  optimizeDeps: {
    exclude: ['@urql/svelte'],
    include: ['lodash.get', 'lodash.isequal', 'lodash.clonedeep'],
  },
  publicDir: 'static',
  resolve: {
    alias: {
      $wailsjs: path.resolve('./src/lib/generated/wailsjs'),
      $lib: path.resolve('./src/lib'),
    },
  },
  build: {
    outDir: 'build',
  },
  server: {
    port: 3000,
    strictPort: true,
  },
});
