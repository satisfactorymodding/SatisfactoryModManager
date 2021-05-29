module.exports = {
  root: true,
  parser: 'vue-eslint-parser',
  parserOptions: {
    parser: '@babel/eslint-parser',
    sourceType: 'module',
  },
  env: {
    browser: true,
    node: true,
  },
  extends: [
    'airbnb-base',
    'plugin:vue/recommended',
    'plugin:import/errors',
    'plugin:import/warnings',
  ],
  globals: {
    __static: true,
  },
  plugins: [
    'html',
  ],
  rules: {
    'linebreak-style': 0,
    'no-underscore-dangle': 0,
    'import/no-extraneous-dependencies': ['error', { devDependencies: true }],
    'no-debugger': process.env.NODE_ENV === 'production' ? 2 : 0,
    'import/extensions': ['error', 'ignorePackages', {
      js: 'never',
      vue: 'never',
    }],
    'max-len': ['error', {
      code: 200,
      ignoreComments: true,
      ignoreTrailingComments: true,
      ignoreUrls: true,
      ignoreStrings: true,
      ignoreTemplateLiterals: true,
      ignoreRegExpLiterals: true,
    }],
    'no-param-reassign': ['error', { props: false }],
  },
  settings: {
    'import/core-modules': [
      '@apollo/client/core',
      '@apollo/client/link/persisted-queries',
    ],
    'import/resolver': {
      node: {
        extensions: ['.js', '.jsx', '.vue'],
      },
      webpack: {
        config: './.electron-vue/webpack.renderer.config.js',
      },
    },
  },
};
