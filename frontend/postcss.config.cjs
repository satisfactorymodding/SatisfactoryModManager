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

    tailwindCSSNesting(),
    tailwindCSS(),
  ],
};
