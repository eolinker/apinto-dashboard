describe('服务发现 e2e test', () => {
  it('初始化页面，点击上游服务-服务发现菜单，进入列表页', async () => {
    await page.goto('http://localhost:4200/login')
    await page.waitForTimeout(2000)
    await page.getByPlaceholder('请输入账号').click()
    await page.getByPlaceholder('请输入账号').fill('maggie')
    await page.getByPlaceholder('请输入账号').press('Tab')
    await page.getByPlaceholder('请输入密码').fill('12345678')
    await page.getByPlaceholder('请输入密码').press('Enter')
    await page.getByText('上游服务').click()
    await page.getByRole('link', { name: '服务发现' }).click()
  })
  it('检查页面样式', async () => {
    // 新建应用的按钮样式
    const createBtn = await page.getByRole('button', { name: '新建服务' })
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

    // 搜索框的样式
    const searchInput = await page.locator('eo-ng-input-group')
    const searchInputH = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const searchInputW = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const searchInputBC = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const searchInputML = await (await page.locator('.mg-top-right >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-right'))

    expect(searchInputH).toStrictEqual('32px')
    expect(searchInputW).toStrictEqual('254px')
    expect(searchInputBC).toStrictEqual('rgb(215, 215, 215)')
    expect(searchInputML).toStrictEqual('24px')

    // 表格的样式
    const listContent = await page.locator('.list-content')
    const listContentMT = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))

    const listTable = await page.locator('eo-ng-apinto-table')
    const listTableMT = await listTable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))

    expect(listContentMT).toStrictEqual('12px')
    expect(listTableMT).toStrictEqual('0px')

    const listTableTh1 = await page.locator('eo-ng-apinto-table tr >> nth = 0 >> th >> nth = 0')
    const listTableTh1Padding = await listTableTh1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTh1Padding).toStrictEqual('0px 12px')

    const listTableTh2 = await page.locator('eo-ng-apinto-table tr >> nth = 0 >> th >> nth = 1')
    const listTableTh2Padding = await listTableTh2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTh2Padding).toStrictEqual('0px 12px')

    const listTableIconTh = await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 5')
    const listTableIconThPadding = await listTableIconTh.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableIconThPadding).toStrictEqual('0px 24px 0px 12px')

    const listTableIcon1 = await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 5 >> button >> nth = 0')
    const listTableIcon1PL = await listTableIcon1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const listTableIcon1PR = await listTableIcon1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    expect(listTableIcon1PL).toStrictEqual('0px')
    expect(listTableIcon1PR).toStrictEqual('8px')

    const listTableIcon2 = await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 5 >> button >> nth = 1')
    const listTableIcon2PL = await listTableIcon2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const listTableIcon2PR = await listTableIcon2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    expect(listTableIcon2PL).toStrictEqual('8px')
    expect(listTableIcon2PR).toStrictEqual('0px')
  })
  it('点击新建服务，检查样式，点击取消并返回列表', async () => {
    await page.getByRole('button', { name: '新建服务' }).click()

    // 注册中心输入框样式
    const nameInput = await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(nameInputW).toStrictEqual('346px')
    await expect(nameInputH).toStrictEqual('32px')

    // 描述样式
    const descInput = await page.locator('textarea#desc')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(descInputW).toStrictEqual('346px')
    await expect(descInputH).toStrictEqual('68px')

    // 服务类型样式
    const schemeInput = await page.locator('eo-ng-select#driver')
    const schemeInputW = await schemeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const schemeInputH = await schemeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(schemeInputW).toStrictEqual('346px')
    await expect(schemeInputH).toStrictEqual('32px')

    // 动态渲染
    const dynamicArea = await page.locator('dynamic-component#dynamic')
    const dynamicAreaP = await dynamicArea.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const dynamicAreaBGC = await dynamicArea.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const dynamicAreaH = await dynamicArea.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const dynamicAreaW = await dynamicArea.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    await expect(dynamicAreaP).toStrictEqual('20px')
    await expect(dynamicAreaBGC).toStrictEqual('rgb(251, 251, 251)')
    await expect(dynamicAreaH).toStrictEqual('auto')
    await expect(dynamicAreaW).toStrictEqual('auto')

    const space1 = await page.locator('nz-space >> nth = 0')
    let space1H = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    let space1W = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    await expect(space1H).toStrictEqual('48px')
    await expect(space1W).toStrictEqual('378px')

    const space1Input1 = await page.locator('nz-space >> nth = 0 >> input')
    const space1Input1H = await space1Input1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const space1Input1W = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    await expect(space1Input1H).toStrictEqual('32px')
    await expect(space1Input1W).toStrictEqual('346px')

    const space1Btn = await page.locator('nz-space >> nth = 0 >> a')
    const space1BtnH = await space1Btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const space1BtnW = await space1Btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const space1BtnLH = await space1Btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const space1BtnC = await space1Btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    await expect(space1BtnH).toStrictEqual('32px')
    await expect(space1BtnW).toStrictEqual('20px')
    await expect(space1BtnLH).toStrictEqual('32px')
    await expect(space1BtnC).toStrictEqual('rgb(34, 84, 157)')

    const space2 = await page.locator('nz-space >> nth = 1')
    let space2H = await space2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    let space2W = await space2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    await expect(space2H).toStrictEqual('48px')
    await expect(space2W).toStrictEqual('474px')

    const space2Input1 = await page.locator('nz-space >> nth = 1 >> input >> nth = 0')
    const space2Input1H = await space2Input1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const space2Input1W = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    await expect(space2Input1H).toStrictEqual('32px')
    await expect(space2Input1W).toStrictEqual('174px')

    const space2Input2 = await page.locator('nz-space >> nth = 1 >> input >> nth = 1')
    const space2Input2H = await space2Input2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const space2Input2W = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const space2Input2M = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    await expect(space2Input2H).toStrictEqual('32px')
    await expect(space2Input2W).toStrictEqual('164px')
    await expect(space2Input2M).toStrictEqual('0px 0px 0px 8px')

    const space2A = await page.locator('nz-space >> nth = 1 >> a >> nth = 0')
    const space2AH = await space2A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const space2ALH = await space2A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const space2AC = await space2A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const space2AM = await space2A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    await expect(space2AH).toStrictEqual('32px')
    await expect(space2ALH).toStrictEqual('32PX')
    await expect(space2AC).toStrictEqual('rgb(34, 84, 157)')
    await expect(space2AM).toStrictEqual('0px 0px 0px 12px')

    const space2A2 = await page.locator('nz-space >> nth = 1 >> a >> nth = 1')
    const space2A2H = await space2A2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const space2A2LH = await space2A2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const space2A2C = await space2A2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const space2A2M = await space2A2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    await expect(space2A2H).toStrictEqual('32px')
    await expect(space2A2LH).toStrictEqual('32PX')
    await expect(space2A2C).toStrictEqual('rgb(34, 84, 157)')
    await expect(space2A2M).toStrictEqual('0px 0px 0px 12px')

    await page.locator('section').filter({ hasText: '*Nacos地址引用环境变量' }).locator('a').click()
    await page.locator('#dynamic a').nth(4).click()

    space1H = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    space1W = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    await expect(space1H).toStrictEqual('48px')
    await expect(space1W).toStrictEqual('410px')

    space2H = await space2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    space2W = await space2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    await expect(space2H).toStrictEqual('48px')
    await expect(space2W).toStrictEqual('474px')

    const space3 = await page.locator('nz-space >> nth = 1')
    const space3H = await space3.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const space3W = await space3.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    await expect(space3H).toStrictEqual('48px')
    await expect(space3W).toStrictEqual('506px')

    const space4 = await page.locator('nz-space >> nth = 3')
    const space4H = await space4.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const space4W = await space4.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    await expect(space4H).toStrictEqual('48px')
    await expect(space4W).toStrictEqual('506px')

    await page.getByPlaceholder('请输入主机名或IP:端口').first().click()
    await page.getByPlaceholder('请输入主机名或IP:端口').first().fill('test')
    await page.getByPlaceholder('请输入主机名或IP:端口').first().fill('')
    await page.getByRole('button', { name: '保存' }).click()

    space1H = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    space1W = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    await expect(space1H).toStrictEqual('70px')
    await expect(space1W).toStrictEqual('410px')

    await page.getByRole('button', { name: '取消' }).click()
  })
  it('点击新建服务，逐一填入必输项并提交，返回列表页', async () => {})
  it('点击表格中的某一项，进入上线管理页，点击tab切换至服务信息页，检查样式；修改内容并提交', async () => {})
  it('点击表格中的某一项，进入上线管理页，检查上线管理页样式与操作，点击面包屑返回列表', async () => {})
  it('点击表格中的某一项，进入上线管理页，点击tab切换至服务信息页，点击取消', async () => {})
  it('点击表格中的某一项，进入上线管理页，点击tab切换至服务信息页，点击面包屑返回列表', async () => {})
})
