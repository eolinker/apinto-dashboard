describe('login e2e test', () => {
  it('login 左侧图片宽360px,高度100%', async () => {
    // Go to http://localhost:4200/login
    await page.goto('http://localhost:4200/login')
    await page.waitForTimeout(2000)
    const pic = await page.locator('nz-sider div >> nth=2')
    const width = await pic.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))

    expect(width).toStrictEqual('360px')
  })

  it('标题字号32px,不加粗,标题与下方次标题间距44px', async () => {
    const title = await page.locator('text=欢迎来到 Apinto')
    const titleFrontSize = await title.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(titleFrontSize).toStrictEqual('32px')
    const titleFrontWeight = await title.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-weight'))
    expect(titleFrontWeight).toStrictEqual('400')

    const titleMarginBottom = await title.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(titleMarginBottom).toStrictEqual('44px')
  })

  it('次标题字16px,不加粗,次标题底部96px*1px', async () => {
    const title = await page.locator('text=账号密码登录')
    const titleFrontSize = await title.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(titleFrontSize).toStrictEqual('16px')
    const titleFrontWeight = await title.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-weight'))
    expect(titleFrontWeight).toStrictEqual('400')

    const titleDiv = await page.locator('div:has-text("账号密码登录") >> nth=2')
    const titleDivWidth = await titleDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(titleDivWidth).toStrictEqual('96px')

    const titleDivBorderWidth = await titleDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-bottom-width'))
    expect(titleDivBorderWidth).toStrictEqual('1px')
  })

  it('次标题与输入框间距28px，输入框372px*40px, 登录按钮372px*42px,输入框之间间距20px,输入框与按钮间距49px', async () => {
    // 次标题与输入框间距28px
    const formDiv = await page.locator('form:has-text("登录")')
    const formDivPaddingTop = await formDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-top'))
    expect(formDivPaddingTop).toStrictEqual('28px')

    // 输入框372px*40px, 输入框之间间距20px
    const accountInput = await page.locator('[placeholder="请输入账号"]')
    const accInputW = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(accInputW).toStrictEqual('372px')

    const accInputH = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(accInputH).toStrictEqual('40px')

    const accInputDiv = await page.locator('nz-form-item').nth(0)
    const accInputDivMarginBottom = await accInputDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(accInputDivMarginBottom).toStrictEqual('20px')

    const pswInput = await page.locator('[placeholder="请输入密码"]')
    const pswInputW = await pswInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(pswInputW).toStrictEqual('372px')

    const pswInputH = await pswInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(pswInputH).toStrictEqual('40px')

    const pswInputDiv = await page.locator('nz-form-item').nth(1)
    const pswInputDivMarginBottom = await pswInputDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(pswInputDivMarginBottom).toStrictEqual('20px')

    // 登录按钮372px*42px，输入框与按钮间距49px
    const loginButton = await page.locator('button:has-text("登录")')
    const loginButtonMarginTop = await loginButton.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    expect(loginButtonMarginTop).toStrictEqual('29px')

    const loginButtonW = await loginButton.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(loginButtonW).toStrictEqual('372px')

    const loginButtonH = await loginButton.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(loginButtonH).toStrictEqual('42px')
  })

  it('次标题与按钮背景色为主题色', async () => {
    const titleDiv = await page.locator('.login-tab')
    const titleDivBgC = await titleDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(titleDivBgC).toStrictEqual('rgb(34, 84, 157)')
    const titleDivBoderC = await titleDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-bottom-color'))
    expect(titleDivBoderC).toStrictEqual('rgb(34, 84, 157)')
  })

  it('未输入账号密码，点击登录，则输入框下方出现提示语，提示语下方存在20px间距', async () => {
    const accInputDiv = await page.locator('nz-form-item').nth(0)
    const accInputDivMarginBottom = await accInputDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(accInputDivMarginBottom).toStrictEqual('20px')

    const pswInputDiv = await page.locator('nz-form-item').nth(1)
    const pswInputDivMarginBottom = await pswInputDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(pswInputDivMarginBottom).toStrictEqual('20px')
  })

  it('输入错误账号密码,提示登录失败;输入正确账号密码,跳转进入项目', async () => {
    // Click [placeholder="请输入账号"]
    await page.locator('[placeholder="请输入账号"]').click()
    // Fill [placeholder="请输入账号"]
    await page.locator('[placeholder="请输入账号"]').fill('maggie')
    // Click [placeholder="请输入密码"]
    await page.locator('[placeholder="请输入密码"]').click()
    // Fill [placeholder="请输入密码"]
    await page.locator('[placeholder="请输入密码"]').fill('11111')
    // Click button:has-text("登录")
    await page.locator('button:has-text("登录")').click()
    // Click #cdk-overlay-0 div:has-text("密码错误") >> nth=3
    await page.locator('#cdk-overlay-0 div:has-text("密码错误")').nth(3).click()

    // Click [placeholder="请输入密码"]
    await page.locator('[placeholder="请输入密码"]').click()
    // Fill [placeholder="请输入密码"]
    await page.locator('[placeholder="请输入密码"]').fill('12345678')
    // Click button:has-text("登录")
    await page.locator('button:has-text("登录")').click()
    // Click #cdk-overlay-0 div:has-text("登录成功") >> nth=3
    await page.locator('#cdk-overlay-0 div:has-text("登录成功")').nth(3).click()
    // Click .logo
    await page.locator('.logo').click()

    await page.locator('.avatar.mg-top-right').click()

    await page.getByText('退出登录').click()

    await page.getByRole('heading', { name: '欢迎来到 Apinto' }).click()
  })
  it('输入无权限账号密码,跳转进入项目,主页出现权限提示提示框.提示框436px*186px, 标题与字体间距24px,字体为14px,标题字体加粗,管理员三字为蓝色, 右下角按钮间间距12px,按钮与右侧间距20px,下方间距24px,右侧按钮为主题色', async () => {
    // 输入无权限账号，出现提示框
    await page.locator('input[placeholder="请输入账号"]').click()
    await page.locator('input[placeholder="请输入账号"]').fill('testNoAccess')
    await page.locator('input[placeholder="请输入密码"]').click()
    await page.locator('input[placeholder="请输入密码"]').fill('12345678')

    await page.getByRole('button', { name: '登录' }).click()
    await page.locator('.ant-modal').click()

    // 提示框样式
    const tipBox = await page.locator('.ant-modal-body')
    const tipBoxHeight = await tipBox.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(tipBoxHeight).toStrictEqual('186px')
    const tipBoxWidth = await tipBox.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(tipBoxWidth).toStrictEqual('436px')

    const tipBoxPadding = await tipBox.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(tipBoxPadding).toStrictEqual('20px')

    const tipBoxTitleIcon = await page.locator('.anticon.anticon-exclamation-circle')

    const tipBoxTitleIconFontSize = await tipBoxTitleIcon.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(tipBoxTitleIconFontSize).toStrictEqual('20px')

    const tipBoxTitleIconMarginRight = await tipBoxTitleIcon.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-right'))
    expect(tipBoxTitleIconMarginRight).toStrictEqual('8px')

    const tipBoxTitleText = await page.locator('.ant-modal-confirm-title')

    const tipBoxTitleTextFontSize = await tipBoxTitleText.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(tipBoxTitleTextFontSize).toStrictEqual('14px')

    const tipBoxTitleTextMarginRight = await tipBoxTitleText.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-weight'))
    expect(tipBoxTitleTextMarginRight).toStrictEqual('500')

    const tipBoxContent = await page.locator('.modal-header').last()

    const tipBoxContentMarginTop = await tipBoxContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    expect(tipBoxContentMarginTop).toStrictEqual('24px')

    const tipBoxContentMarginRight = await tipBoxTitleIcon.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-right'))
    expect(tipBoxContentMarginRight).toStrictEqual('8px')

    const tipBoxManager = await page.locator('text=管理员')

    const tipBoxManagerColor = await tipBoxManager.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(tipBoxManagerColor).toStrictEqual('rgb(0, 0, 255)')

    const tipBoxBtns = await page.locator('.ant-modal-confirm-btns')

    const tipBoxBtnsMarginTop = await tipBoxBtns.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    expect(tipBoxBtnsMarginTop).toStrictEqual('24px')

    const tipBoxBtnsFloat = await tipBoxBtns.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('float'))
    expect(tipBoxBtnsFloat).toStrictEqual('right')

    const tipBoxBtnsRight = await page.locator('role=button[name="确定"]')

    const tipBoxBtnsRightMarginLeft = await tipBoxBtnsRight.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(tipBoxBtnsRightMarginLeft).toStrictEqual('12px')

    const tipBoxBtnsRightColor = await tipBoxBtnsRight.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(tipBoxBtnsRightColor).toStrictEqual('rgb(34, 84, 157)')

    await page.getByRole('button', { name: '确定' }).click()

    await page.locator('div:has-text("网关控制台")').click()

    await page.locator('.avatar.mg-top-right').click()

    await page.getByText('退出登录').click()

    await page.getByRole('heading', { name: '欢迎来到 Apinto' }).click()
  })
  it('点击取消,则关闭提示框,页面不做跳转', async () => {
    await page.locator('input[placeholder="请输入账号"]').click()
    await page.locator('input[placeholder="请输入账号"]').fill('testNoAccess')
    await page.locator('input[placeholder="请输入密码"]').click()
    await page.locator('input[placeholder="请输入密码"]').fill('12345678')
    await page.getByRole('button', { name: '登录' }).click()
    await page.locator('.ant-modal-body').click()

    await page.getByRole('button', { name: '取消' }).click()

    await page.locator('div:has-text("网关控制台")').click()

    await page.locator('.avatar.mg-top-right').click()

    await page.getByText('退出登录').click()

    await page.getByRole('heading', { name: '欢迎来到 Apinto' }).click()
  })

  it('点击关闭,则关闭提示框,页面不做跳转', async () => {
    await page.locator('input[placeholder="请输入账号"]').click()
    await page.locator('input[placeholder="请输入账号"]').fill('testNoAccess')
    await page.locator('input[placeholder="请输入密码"]').click()
    await page.locator('input[placeholder="请输入密码"]').fill('12345678')
    await page.getByRole('button', { name: '登录' }).click()
    await page.locator('.ant-modal-close-x').click()

    await page.locator('div:has-text("网关控制台")').click()

    await page.locator('.avatar.mg-top-right').click()

    await page.getByText('退出登录').click()

    await page.getByRole('heading', { name: '欢迎来到 Apinto' }).click()
  })
})
