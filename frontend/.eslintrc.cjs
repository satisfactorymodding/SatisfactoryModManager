module.exports = {
  root: true,
  extends: [
    'eslint:recommended',
  ],
  plugins: [
    '@typescript-eslint',
    'import',
  ],
  parserOptions: {
    sourceType: 'module',
    ecmaVersion: 2020,
    tsconfigRootDir: __dirname,
    project: ['./tsconfig.json'],
    extraFileExtensions: ['.svelte'],
  },
  overrides: [
    {
      files: ['*.svelte'],
      parser: 'svelte-eslint-parser',
      extends: [
        'plugin:svelte/recommended',
        'plugin:@typescript-eslint/recommended',
      ],
      parserOptions: {
        parser: '@typescript-eslint/parser',
      },
      rules: {
        'svelte/indent': ['error', { indent: 2 }],
        'svelte/valid-compile': ['error', { ignoreWarnings: true }],
        indent: 'off',
      },
    },
    {
      files: ['*.ts'],
      parser: '@typescript-eslint/parser',
      extends: [
        'plugin:@typescript-eslint/recommended',
      ],
    },
    {
      files: ['*.js', '*.cjs'],
      parser: 'espree',
    },
  ],
  env: {
    browser: true,
    es2017: true,
    node: true,
  },
  rules: {
    'no-multi-spaces': 'error',
    indent: ['error', 2],
    quotes: ['error', 'single'],
    curly: ['error', 'multi-line'],
    'no-extra-semi': 'error',
    'no-var': 'error',
    'quote-props': ['error', 'as-needed', { keywords: false, unnecessary: true, numbers: false }],
    semi: ['error', 'always'],
    'comma-dangle': ['error', 'always-multiline'],
    'object-curly-spacing': ['error', 'always'],
    'import/order': ['error', { 'newlines-between': 'always' }],
    'import/newline-after-import': ['error'],
    'import/no-duplicates': ['error'],
  },
};
