/** @type {import('tailwindcss').Config} */
module.exports = {
  mode: 'jit',
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      screens: {
        'h-md': { raw: '(min-height: 825px)' },
        'h-lg': { raw: '(min-height: 900px)' },
      },
    },
  },
  plugins: [
    require('@tailwindcss/container-queries'),
  ],
};
