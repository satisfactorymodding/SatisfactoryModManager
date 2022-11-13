/** @type {import('tailwindcss').Config} */
module.exports = {
  mode: 'jit',
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      screens: {
        md: { raw: '(min-height: 825px)' },
        lg: { raw: '(min-height: 900px)' },
        'w-md': { raw: '(min-width: 1500px)' },
      },
    },
  },
};
