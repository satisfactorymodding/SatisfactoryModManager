const importUrl = require('postcss-import-url');
const postcssPresetEnv = require('postcss-preset-env');
const tailwindCSS = require('tailwindcss');
const tailwindCSSNesting = require('tailwindcss/nesting');

module.exports = {
  plugins: [
    postcssPresetEnv({
      stage: 4,
      features: {
        'nesting-rules': true,
      },
    }),

    importUrl({
      modernBrowser: true,
    }),
    tailwindCSSNesting(),
    tailwindCSS(),
  ],
};
