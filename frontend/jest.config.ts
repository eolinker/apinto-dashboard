/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-11 23:25:53
 * @LastEditors: maggieyyy im.ymj@hotmail.com
 * @LastEditTime: 2022-07-11 23:26:03
 * @FilePath: /apinto/jest.config.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
export default {
  clearMocks: true,
  collectCoverage: true,
  coverageDirectory: './unit_coverage', // 单元测试报告生成目录
  coverageProvider: 'v8',
  coverageThreshold: {
    // 所有文件總的覆蓋率要求
    global: {
      branches: 0,
      functions: 0,
      lines: 0,
      statements: 0
    },
    './src/**/*[^.old].{ts,tsx,js,jsx}': {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80
    }
  },
  modulePaths: [
    '<rootDir>',
    '/node_modules'
  ],
  moduleNameMapper: { // alias，目录映射关系对象
    '~src/(.*)$': '<rootDir>/src/$1',
    uuid: require.resolve('uuid'),
    '^lodash-es$': 'lodash'
  },
  globals: {
    'ts-jest': {
      tsconfig: 'tsconfig.json'
    }
  },
  preset: 'jest-preset-angular',
  roots: ['<rootDir>'],
  setupFilesAfterEnv: ['<rootDir>/test/setup-jest.ts'],
  testEnvironment: 'jsdom',

  testMatch: ['**/**/*.spec.{js,jsx,ts,tsx}'],
  testPathIgnorePatterns: ['/node_modules/', '/dist/'],
  transform: {
    '^.+\\.(js|jsx)$': 'babel-jest'
  }
}
