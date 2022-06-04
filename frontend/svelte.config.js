import adapter from '@sveltejs/adapter-static';
import preprocess from 'svelte-preprocess';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  // Consult https://github.com/sveltejs/svelte-preprocess
  // for more information about preprocessors
  preprocess: preprocess({
    postcss: true,
  }),
  kit: {
    adapter: adapter({
      out: 'build',
      fallback: 'index.html',
    }),
    vite: {
      server: {
        fs: {
          allow: ['wailsjs'],
        },
      },
      optimizeDeps: {
        exclude: ['@urql/svelte'],
      },
    },
  },
};

export default config;
