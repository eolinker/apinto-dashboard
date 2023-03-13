/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-11 23:26:13
 * @LastEditors: maggieyyy im.ymj@hotmail.com
 * @LastEditTime: 2022-07-11 23:26:24
 * @FilePath: /apinto/jest.e2e.comfig.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
module.exports = {
  roots: ['<rootDir>/projects'],
  preset: 'jest-playwright-preset',
  transform: {
    '^.+\\.ts$': 'ts-jest',
    '^.+\\.js$': 'babel-jest'

  },
  testPathIgnorePatterns: ['/node_modules/', 'dist', 'lib'],
  testMatch: ['**/**/*.test.{js,jsx,ts,tsx}'],
  verbose: true,
  setupFilesAfterEnv: ['<rootDir>/test/setup-e2e.ts'],
  moduleNameMapper: {
    '~src/(.*)$': '<rootDir>/projects/core/src/$1'
    // angular 13+
    // '@angular/core/testing': '<rootDir>/node_modules/@angular/core/fesm2015/testing.mjs',
    // '@angular/platform-browser-dynamic/testing': '<rootDir>/node_modules/@angular/platform-browser-dynamic/fesm2015/testing.mjs',
    // '@angular/platform-browser/testing': '<rootDir>/node_modules/@angular/platform-browser/fesm2015/testing.mjs',
  }
}
