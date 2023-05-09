describe('环境变量 e2e test', () => {
  it('初始化页面，点击基础设施-环境变量菜单，进入列表页', async () => {
    await page.goto('http://localhost:4200/login')
    await page.waitForTimeout(2000)
    await page.getByPlaceholder('请输入账号').click()
    await page.getByPlaceholder('请输入账号').fill('maggie')
    await page.getByPlaceholder('请输入账号').press('Tab')
    await page.getByPlaceholder('请输入密码').fill('12345678')
    await page.getByPlaceholder('请输入密码').press('Enter')
    await page.getByText('基础设施').click()
    await page.getByRole('link', { name: '环境变量' }).click()
  })
  it('检查页面样式', async () => {
    // 新建配置按钮
    const createBtn = await page.getByRole('button', { name: '新建配置' })
    const createBtnH = await createBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const createBtnW = await createBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const createBtnBG = await createBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const createBtnBC = await createBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const createBtnFS = await createBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const createBtnML = await createBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    expect(createBtnH).toStrictEqual('32px')
    expect(createBtnW).toStrictEqual('82px')
    expect(createBtnBG).toStrictEqual('rgb(34, 84, 157)')
    expect(createBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(createBtnFS).toStrictEqual('14px')
    expect(createBtnML).toStrictEqual('12px')

    // 搜索组合
    const searchGroup = await page.locator('.list-header.block_lr div >> nth = 1')
    const searchGroupH = await searchGroup.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(searchGroupH).toStrictEqual('32px')

    // key标签
    const keyLable = await page.locator('label').filter({ hasText: 'KEY' })
    const keyLableP = await keyLable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const keyLableM = await keyLable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    expect(keyLableP).toStrictEqual('0px')
    expect(keyLableM).toStrictEqual('0px')

    // key输入框
    const searchInput1 = await page.getByPlaceholder('请输入')
    const searchInput1H = await searchInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const searchInput1W = await searchInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const searchInput1BC = await searchInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const searchInput1ML = await searchInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    expect(searchInput1H).toStrictEqual('32px')
    expect(searchInput1W).toStrictEqual('216px')
    expect(searchInput1BC).toStrictEqual('rgb(215, 215, 215)')
    expect(searchInput1ML).toStrictEqual('12px')

    // key标签
    const statusLable = await page.locator('label').filter({ hasText: 'KEY' })
    const statusLableP = await statusLable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const statusLableM = await statusLable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    expect(statusLableP).toStrictEqual('0px')
    expect(statusLableM).toStrictEqual('0px 0px 0px 24px')

    // 状态输入框
    const searchInput2 = await page.getByPlaceholder('请输入')
    const searchInput2H = await searchInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const searchInput2W = await searchInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const searchInput2BC = await searchInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const searchInput2ML = await searchInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    expect(searchInput2H).toStrictEqual('32px')
    expect(searchInput2W).toStrictEqual('216px')
    expect(searchInput2BC).toStrictEqual('rgb(215, 215, 215)')
    expect(searchInput2ML).toStrictEqual('12px')

    // 重置按钮
    const resetBtn = await page.getByRole('button', { name: '发布历史' })
    const resetBtnM = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const resetBtnP = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const resetBtnBC = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const resetBtnBGC = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const resetBtnC = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(resetBtnM).toStrictEqual('0px 0px 0px 24px')
    expect(resetBtnP).toStrictEqual('0px 12px')
    expect(resetBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(resetBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(resetBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')

    // 查询按钮
    const searchBtn = await page.getByRole('button', { name: '新建配置' })
    const searchBtnH = await searchBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const searchBtnW = await searchBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const searchBtnBG = await searchBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const searchBtnBC = await searchBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const searchBtnFS = await searchBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const searchBtnM = await searchBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const searchBtnP = await searchBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))

    expect(searchBtnH).toStrictEqual('32px')
    expect(searchBtnW).toStrictEqual('54px')
    expect(searchBtnBG).toStrictEqual('rgb(34, 84, 157)')
    expect(searchBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(searchBtnFS).toStrictEqual('14px')
    expect(searchBtnM).toStrictEqual('0px 24px 0px 12px')
    expect(searchBtnP).toStrictEqual('0px 12px')
  })
  it('点击新建配置，检查样式；逐一填入必输项，直至保存后返回列表', async () => {
    await page.getByRole('button', { name: '新建配置' }).click()

    // 应用名称输入框样式
    const nameInput = await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(nameInputW).toStrictEqual('346px')
    await expect(nameInputH).toStrictEqual('32px')

    // 描述样式
    const descInput = await page.getByPlaceholder('请输入')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(descInputW).toStrictEqual('346px')
    await expect(descInputH).toStrictEqual('68px')

    // 重置按钮
    const saveBtn = await page.getByRole('button', { name: '保存' })
    const saveBtnM = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const saveBtnP = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const saveBtnBC = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const saveBtnBGC = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const saveBtnC = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(saveBtnM).toStrictEqual('0px')
    expect(saveBtnP).toStrictEqual('0px 12px')
    expect(saveBtnBGC).toStrictEqual('rgb(34, 84, 157)')
    expect(saveBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(saveBtnC).toStrictEqual('rgb(255, 255, 255)')

    // 查询按钮
    const cancleBtn = await page.getByRole('button', { name: '取消' })
    const cancleBtnH = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const cancleBtnW = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const cancleBtnBG = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const cancleBtnBC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const cancleBtnFS = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const cancleBtnM = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const cancleBtnP = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))

    expect(cancleBtnH).toStrictEqual('32px')
    expect(cancleBtnW).toStrictEqual('54px')
    expect(cancleBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(cancleBtnBG).toStrictEqual('rgb(255, 255, 255)')
    expect(cancleBtnFS).toStrictEqual('14px')
    expect(cancleBtnM).toStrictEqual('0px 0px 0px 12px')
    expect(cancleBtnP).toStrictEqual('0px 12px')

    await page.getByRole('button', { name: '保存' }).click()
    await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').click()
    await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').fill('testValue')
    await page.getByRole('button', { name: '保存' }).click()
  })
  it('点击新建配置，点击取消后返回列表', async () => {
    await page.getByRole('button', { name: '新建配置' }).click()
    await page.getByPlaceholder('请输入').click()
    await page.getByPlaceholder('请输入').fill('testdesc')
    await page.getByRole('button', { name: '取消' }).click()
  })
  it('点击新建配置，点击面包屑返回列表', async () => {
    await page.getByRole('button', { name: '新建配置' }).click()
    await page.getByText('环境变量 / 新建配置 /').click()
    await page.locator('nz-breadcrumb').getByRole('link', { name: '环境变量' }).click()
  })
  it('点击表格，出现环境变量弹窗，检查样式并关闭', async () => {
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').first().click()
    await page.getByText('查看环境变量').click()
    const drawerTable = await page.locator('.drawer-table')
    const drawerTableM = await drawerTable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const drawerTableP = await drawerTable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    expect(drawerTableM).toStrictEqual('10px auto 0px auto')
    expect(drawerTableP).toStrictEqual('0px')
    await page.getByRole('button', { name: 'Close' }).click()
  })
  it('测试表格中的图标', async () => {
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').last().locator('button').first()
    await page.getByText('查看环境变量').click()
    await page.getByRole('button', { name: 'Close' }).click()
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').last().locator('button').last()
    await page.getByText('该数据删除后将无法找回，请确认是否删除？').click()
    await page.locator('.ant-modal-footer').getByRole('button', { name: '取消' }).click()
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').last().locator('button').last()
    await page.getByText('该数据删除后将无法找回，请确认是否删除？').click()
    await page.locator('.ant-modal-footer').getByRole('button', { name: '确定' }).click()
  })
})
