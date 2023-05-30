describe('商业授权 e2e test', () => {
  it('登陆前进入auth页面，检查样式；上传证书后进入auth信息查询页，检查样式', async () => {
    await page.goto('http://localhost:4200/auth')

    const activationBlock = await page.locator('.activation-block')
    const activationBlockW = await activationBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const activationBlockBGC = await activationBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const activationBlockMT = await activationBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    const activationBlockP = await activationBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(activationBlockW).toStrictEqual('760px')
    expect(activationBlockBGC).toStrictEqual('rgb(248, 248, 250)')
    expect(activationBlockMT).toStrictEqual('10vh')
    expect(activationBlockP).toStrictEqual('54px 0px')

    const title = await page.getByRole('heading', { name: '激活系统' })
    const titleFrontSize = await title.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const titleFrontWeight = await title.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-weight'))
    const titleFrontTA = await title.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('text-align'))
    expect(titleFrontSize).toStrictEqual('24px')
    expect(titleFrontWeight).toStrictEqual('700')
    expect(titleFrontTA).toStrictEqual('center')

    const activationContent = await page.locator('.activation-content')
    const activationContentM = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const activationContentC = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const activationContentFS = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const activationContentFW = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-weight'))
    const activationContentLH = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const activationContentLST = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('list-style-type'))
    expect(activationContentM).toStrictEqual('40px auto 0px auto')
    expect(activationContentC).toStrictEqual('rgb(51, 51, 51)')
    expect(activationContentFS).toStrictEqual('14px')
    expect(activationContentFW).toStrictEqual('500')
    expect(activationContentLH).toStrictEqual('26px')
    expect(activationContentLST).toStrictEqual('none')

    const listContent = await page.locator('.list-content')
    const listContentML = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(listContentML).toStrictEqual('28px')

    const activationButtonDisabled = await (await page.getByRole('button', { name: '激活系统' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    expect(activationButtonDisabled).toStrictEqual('true')

    await page.getByRole('button', { name: '复制' }).click()
    await page.getByText('复制成功').click()
    await page.setInputFiles('button:has-text("上传证书")', './test/apinto-2027.cert')
    await page.getByRole('button', { name: '激活系统' }).click()

    await page.getByRole('heading', { name: '标准版授权' }).click()

    const activationInfo = await page.locator('.activation-info')
    const activationInfoC = await activationInfo.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(activationInfoC).toStrictEqual('rgb(102, 102, 102)')

    const goToLogin = await page.locator('.auth-a')
    const goToLoginM = await goToLogin.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const goToLoginTL = await goToLogin.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('text-align'))
    expect(goToLoginM).toStrictEqual('rgb(24px auto 0px auto)')
    expect(goToLoginTL).toStrictEqual('center')
    await page.getByText('跳转至登录页').click()
  })
  it('点击登录，进入登录页，从左侧下方进入授权页面', async () => {
    if (await page.getByRole('heading', { name: '欢迎来到 Apinto' }) && await page.getByRole('heading', { name: '欢迎来到 Apinto' }).isVisible()) {
      await page.getByPlaceholder('请输入账号').click()
      await page.getByPlaceholder('请输入账号').fill('maggie')
      await page.getByPlaceholder('请输入账号').press('Tab')
      await page.getByPlaceholder('请输入密码').fill('12345678')
      await page.getByRole('button', { name: '登录' }).click()
      await page.getByText('版权所有 © 银云信息技术').click()
    }

    const authInfoBlock = await page.locator('.auth-box')
    const authInfoBlockPo = await authInfoBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('position'))
    const authInfoBlockP = await authInfoBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const authInfoBlockB = await authInfoBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('bottom'))
    const authInfoBlockW = await authInfoBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const authInfoBlockH = await authInfoBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const authInfoBlockBGC = await authInfoBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const authInfoBlockC = await authInfoBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const authInfoBlockFS = await authInfoBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const authInfoBlockLH = await authInfoBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const authInfoBlockTA = await authInfoBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('text-align'))
    expect(authInfoBlockPo).toStrictEqual('absolute')
    expect(authInfoBlockP).toStrictEqual('12px 8px')
    expect(authInfoBlockB).toStrictEqual('0px')
    expect(authInfoBlockW).toStrictEqual('100%')
    expect(authInfoBlockH).toStrictEqual('68px')
    expect(authInfoBlockBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(authInfoBlockC).toStrictEqual('rgb(153, 153, 153)')
    expect(authInfoBlockFS).toStrictEqual('12px')
    expect(authInfoBlockLH).toStrictEqual('20px')
    expect(authInfoBlockTA).toStrictEqual('center')

    const modalActiveInfo = await page.locator('.auth-modal-header .auth-info-modal ol')
    const modalActiveInfoC = await modalActiveInfo.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const modalActiveInfoFS = await modalActiveInfo.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const modalActiveInfoFW = await modalActiveInfo.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-weight'))
    const modalActiveInfoLH = await modalActiveInfo.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const modalActiveInfoLST = await modalActiveInfo.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('list-style-type'))
    const modalActiveInfoMB = await modalActiveInfo.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(modalActiveInfoC).toStrictEqual('rgb(102, 102, 102)')
    expect(modalActiveInfoFS).toStrictEqual('14px')
    expect(modalActiveInfoFW).toStrictEqual('500')
    expect(modalActiveInfoLH).toStrictEqual('26px')
    expect(modalActiveInfoLST).toStrictEqual('none')
    expect(modalActiveInfoMB).toStrictEqual('12px')
    await page.getByText('更新授权').click()
    await page.getByRole('heading', { name: '更新授权' }).click()
  })
  it('更新授权，检查授权信息页', async () => {
    const activationBlock = await page.locator('.activation-block')
    const activationBlockW = await activationBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const activationBlockBGC = await activationBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const activationBlockMT = await activationBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    const activationBlockP = await activationBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(activationBlockW).toStrictEqual('760px')
    expect(activationBlockBGC).toStrictEqual('rgb(248, 248, 250)')
    expect(activationBlockMT).toStrictEqual('10vh')
    expect(activationBlockP).toStrictEqual('54px 0px')

    const title = await page.getByRole('heading', { name: '更新授权' })
    const titleFrontSize = await title.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const titleFrontWeight = await title.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-weight'))
    const titleFrontTA = await title.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('text-align'))
    expect(titleFrontSize).toStrictEqual('24px')
    expect(titleFrontWeight).toStrictEqual('700')
    expect(titleFrontTA).toStrictEqual('center')

    const activationContent = await page.locator('.activation-content')
    const activationContentM = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const activationContentC = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const activationContentFS = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const activationContentFW = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-weight'))
    const activationContentLH = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const activationContentLST = await activationContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('list-style-type'))
    expect(activationContentM).toStrictEqual('40px auto 0px auto')
    expect(activationContentC).toStrictEqual('rgb(51, 51, 51)')
    expect(activationContentFS).toStrictEqual('14px')
    expect(activationContentFW).toStrictEqual('500')
    expect(activationContentLH).toStrictEqual('26px')
    expect(activationContentLST).toStrictEqual('none')

    const listContent = await page.locator('.list-content')
    const listContentML = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(listContentML).toStrictEqual('28px')

    const activationButtonDisabled = await (await page.getByRole('button', { name: '更新授权' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    expect(activationButtonDisabled).toStrictEqual('true')

    await page.getByRole('button', { name: '复制' }).click()
    await page.getByText('复制成功').click()
    await page.setInputFiles('button:has-text("上传证书")', './test/apinto-2027.cert')
    await page.getByRole('button', { name: '更新授权' }).click()

    await page.getByRole('heading', { name: '标准版授权' }).click()

    const activationInfo = await page.locator('.activation-info')
    const activationInfoC = await activationInfo.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(activationInfoC).toStrictEqual('rgb(102, 102, 102)')

    const goToLogin = await page.locator('.auth-a')
    const goToLoginM = await goToLogin.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const goToLoginTL = await goToLogin.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('text-align'))
    expect(goToLoginM).toStrictEqual('rgb(24px auto 0px auto)')
    expect(goToLoginTL).toStrictEqual('center')
    await page.getByText('更新授权').click()
  })
})
