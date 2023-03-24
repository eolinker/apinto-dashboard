/*
 * @Author: maggieyyy im.ymj@hotmail.com
 * @Date: 2022-07-11 23:26:46
 * @LastEditors: maggieyyy im.ymj@hotmail.com
 * @LastEditTime: 2022-07-11 23:26:55
 * @FilePath: /apinto/jest-playwright.config.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
module.exports = {
  launchOptions: {
    headless: !process.argv.includes('--headless=false'),
    viewport: { width: 1440, height: 740 }
  },
  serverOptions: {
    command: 'npm run start',
    port: 4200,
    launchTimeout: 1000000,
    debug: true,
    usedPortAction: 'ignore'
  },
  browsers: ['chromium'],
  devices: [],
  collectCoverage: true,
  config: {
    use: {
      viewport: { width: 1440, height: 720 }
    }
  }
}
