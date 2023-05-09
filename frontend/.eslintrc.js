module.exports = {
  env: {
    browser: true,
    es2021: true
  },
  extends: ['standard'],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 12,
    sourceType: 'module'
  },
  plugins: ['@typescript-eslint'],
  rules: {
    'no-unused-vars': 'off',
    '@typescript-eslint/no-unused-vars': 2,
    'no-useless-constructor': 'off',
    '@typescript-eslint/no-useless-constructor': 'warn'
  },
  globals: {
    angular: true,
    page: true,
    describe: true,
    it: true,
    test: true,
    expect: true,
    beforeEach: true,
    jest: true
  }
}
