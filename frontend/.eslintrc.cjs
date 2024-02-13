module.exports = {
  root: true,
  extends: [
    'eslint:recommended',
    'plugin:@stylistic/disable-legacy',
  ],
  plugins: [
    '@typescript-eslint',
    'import',
    '@stylistic',
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
        /// linting
        'svelte/valid-compile': ['error', { ignoreWarnings: true }],
        'svelte/no-extra-reactive-curlies': 'warn',
        'svelte/html-self-closing': 'error',
        '@typescript-eslint/no-unused-vars': 'warn',
        /// formatting
        '@stylistic/indent': 'off',
        'svelte/indent': ['error', { indent: 2 }],
        'svelte/first-attribute-linebreak': 'warn',
        'svelte/max-attributes-per-line': ['warn', { singleline: 3 }],
        'svelte/mustache-spacing': 'warn',
        'svelte/no-spaces-around-equal-signs-in-attribute': 'error',
        'svelte/sort-attributes': 'warn',
      },
    },
    {
      files: ['*.ts'],
      parser: '@typescript-eslint/parser',
      extends: [
        'plugin:@typescript-eslint/recommended',
      ],
      rules: {
        '@typescript-eslint/no-unused-vars': 'warn',
      },
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
    /// linting
    curly: ['error', 'multi-line'],
    'no-var': 'error',
    /// formatting
    // general
    '@stylistic/no-multi-spaces': 'error',
    '@stylistic/brace-style': ['error', '1tbs', { allowSingleLine: true }],
    '@stylistic/indent': ['error', 2],
    '@stylistic/quotes': ['error', 'single'],
    '@stylistic/no-extra-semi': 'error',
    '@stylistic/quote-props': ['error', 'as-needed', { keywords: false, unnecessary: true, numbers: false }],
    '@stylistic/semi': ['error', 'always'],
    '@stylistic/comma-dangle': ['error', 'always-multiline'],
    '@stylistic/comma-spacing': 'error',
    '@stylistic/object-curly-spacing': ['error', 'always'],
    // imports
    'sort-imports': ['error', { ignoreDeclarationSort: true }],
    'import/order': ['error', { 'newlines-between': 'always', alphabetize: { order: 'asc' } }],
    'import/newline-after-import': ['error'],
    'import/no-duplicates': ['error'],
  },
};
