import adapter from '@sveltejs/adapter-auto';
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
      // default options are shown
      pages: "build",
      assets: "build",
		}),

		vite: {
			server: {
				fs: {
					allow: [
						'wailsjs',
					]
				}
			}
		}
	}
};

export default config;
