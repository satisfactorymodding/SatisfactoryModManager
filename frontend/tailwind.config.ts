import { join } from 'path';

import { skeleton } from '@skeletonlabs/tw-plugin';
import containerQueries from '@tailwindcss/container-queries';
import type { Config } from 'tailwindcss';

import { myCustomTheme } from './smmTheme';

const config = {
  darkMode: 'class',
  content: [
    './src/**/*.{html,js,svelte,ts}',
    join(require.resolve(
      '@skeletonlabs/skeleton'),
    '../**/*.{html,js,svelte,ts}',
    ),
  ],
  theme: {
    extend: {
      screens: {
        'h-md': { raw: '(min-height: 875px)' },
        'h-lg': { raw: '(min-height: 950px)' },
      },
    },
  },
  plugins: [
    containerQueries,
    skeleton({
      themes: {
        custom: [
          myCustomTheme,
        ],
      },
    }),
  ],
} satisfies Config;

export default config;
