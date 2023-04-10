describe('上游管理 e2e test', () => {
  it('初始化页面，点击基础设施-环境变量菜单，进入列表页', async () => {
    await page.goto('http://localhost:4200/login')
    await page.waitForTimeout(2000)
    await page.getByPlaceholder('请输入账号').click()
    await page.getByPlaceholder('请输入账号').fill('maggie')
    await page.getByPlaceholder('请输入账号').press('Tab')
    await page.getByPlaceholder('请输入密码').fill('12345678')
    await page.getByPlaceholder('请输入密码').press('Enter')
    await page.getByText('基础设施').click()
    await page.getByRole('link', { name: '网关集群' }).click()
  })
  it('检查页面样式', async () => {
    // 新建应用的按钮样式
    const createBtn = await page.getByRole('button', { name: '新建上游' })
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

    // 分页样式
    const paginationM = await page.locator('.mg_pagination_t')
    const paginationMT = await paginationM.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))

    const paginationMH = await paginationM.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(paginationMT).toStrictEqual('16px')
    expect(paginationMH).toStrictEqual('32px')

    const pagination = await page.locator('eo-ng-pagination')
    const paginationMR = await pagination.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-right'))
    const paginationH = await pagination.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(paginationMR).toStrictEqual('24px')
    expect(paginationH).toStrictEqual('32px')
  })
  it('点击新建上游，检查页面样式，测试添加环境变量弹窗后，点击取消返回列表', async () => {
    await page.getByRole('button', { name: '新建上游' }).click()

    // 上游名称输入框样式
    const nameInput = await page.locator('input#name')
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

    // 请求协议样式
    const schemeInput = await page.locator('eo-ng-select#scheme')
    const schemeInputW = await schemeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const schemeInputH = await schemeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(schemeInputW).toStrictEqual('346px')
    await expect(schemeInputH).toStrictEqual('32px')

    // 负载算法输入框样式
    const balanceInput = await page.locator('eo-ng-select#balance')
    const balanceInputW = await balanceInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const balanceInputH = await balanceInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(balanceInputW).toStrictEqual('346px')
    await expect(balanceInputH).toStrictEqual('32px')

    // 服务发现输入框样式
    const discoveryInput = await page.locator('eo-ng-select#discoveryName')
    const discoveryInputW = await discoveryInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const discoveryInputH = await discoveryInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(discoveryInputW).toStrictEqual('346px')
    await expect(discoveryInputH).toStrictEqual('32px')

    // 请求超时时间输入框样式
    const timeoutInput = await page.locator('input#timeout')
    const timeoutInputW = await timeoutInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const timeoutInputH = await timeoutInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(timeoutInputW).toStrictEqual('346px')
    await expect(timeoutInputH).toStrictEqual('32px')

    // 动态渲染部分
    const checkboxLabel = await page.locator('.ant-checkbox-wrapper span >> nth = 1').getByText('引用环境变量')
    const checkboxLabelC = await checkboxLabel.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const checkboxLabelPL = await checkboxLabel.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    await expect(checkboxLabelC).toStrictEqual('rgb(102, 102, 102)')
    await expect(checkboxLabelPL).toStrictEqual('12px')

    const space1 = await page.locator('.ArrayItems nz-space').first()
    let space1W = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    let space1H = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(space1W).toStrictEqual('517px')
    await expect(space1H).toStrictEqual('52px')

    const space1Input1 = await page.locator('.ArrayItems nz-space').first().locator('input').first()
    let space1Input1W = await space1Input1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    let space1Input1H = await space1Input1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(space1Input1W).toStrictEqual('346px')
    await expect(space1Input1H).toStrictEqual('32px')

    const space1Input2 = await page.locator('.ArrayItems nz-space').first().locator('input >> nth = 1')
    let space1Input2W = await space1Input2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    let space1Input2H = await space1Input2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(space1Input2W).toStrictEqual('131px')
    await expect(space1Input2H).toStrictEqual('32px')

    const space1Btn1 = await page.locator('.ArrayItems nz-space').first().locator('a >> nth = 1')
    let space1Btn1W = await space1Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    let space1Btn1H = await space1Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    let space1Btn1C = await space1Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    let space1Btn1LH = await space1Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    let space1Btn1M = await space1Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    await expect(space1Btn1W).toStrictEqual('20')
    await expect(space1Btn1H).toStrictEqual('32px')
    await expect(space1Btn1C).toStrictEqual('rgb(34, 84, 157)')
    await expect(space1Btn1LH).toStrictEqual('32px')
    await expect(space1Btn1M).toStrictEqual('0px 0px 0px 12px')

    await space1Btn1.click()

    space1W = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    space1H = await space1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(space1W).toStrictEqual('549px')
    await expect(space1H).toStrictEqual('52px')

    space1Input1W = await space1Input1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    space1Input1H = await space1Input1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(space1Input1W).toStrictEqual('346px')
    await expect(space1Input1H).toStrictEqual('32px')

    space1Input2W = await space1Input2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    space1Input2H = await space1Input2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(space1Input2W).toStrictEqual('131px')
    await expect(space1Input2H).toStrictEqual('32px')

    space1Btn1W = await space1Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    space1Btn1H = await space1Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    space1Btn1C = await space1Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    space1Btn1LH = await space1Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    space1Btn1M = await space1Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    await expect(space1Btn1W).toStrictEqual('131px')
    await expect(space1Btn1H).toStrictEqual('32px')
    await expect(space1Btn1C).toStrictEqual('rgb(34, 84, 157)')
    await expect(space1Btn1LH).toStrictEqual('32px')
    await expect(space1Btn1M).toStrictEqual('0px 0px 0px 12px')

    const space2 = await page.locator('.ArrayItems nz-space >> nth = 1')
    const space2W = await space2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const space2H = await space2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(space2W).toStrictEqual('517px')
    await expect(space2H).toStrictEqual('52px')

    const space2Input1 = await page.locator('.ArrayItems nz-space >> nth = 1').locator('input').first()
    const space2Input1W = await space2Input1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const space2Input1H = await space2Input1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(space2Input1W).toStrictEqual('346px')
    await expect(space2Input1H).toStrictEqual('32px')

    const space2Input2 = await page.locator('.ArrayItems nz-space >> nth = 1').locator('input >> nth = 1')
    const space2Input2W = await space2Input2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const space2Input2H = await space2Input2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(space2Input2W).toStrictEqual('131px')
    await expect(space2Input2H).toStrictEqual('32px')

    const space2Btn1 = await page.locator('.ArrayItems nz-space >> nth = 1').locator('a >> nth = 0')
    const space2Btn1H = await space2Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const space2Btn1C = await space2Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const space2Btn1LH = await space2Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const space2Btn1M = await space2Btn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    await expect(space2Btn1H).toStrictEqual('32px')
    await expect(space2Btn1C).toStrictEqual('rgb(34, 84, 157)')
    await expect(space2Btn1LH).toStrictEqual('32px')
    await expect(space2Btn1M).toStrictEqual('0px 0px 0px 12px')

    const space2Btn2 = await page.locator('.ArrayItems nz-space >> nth = 1').locator('a >> nth = 1')
    const space2Btn2W = await space2Btn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const space2Btn2H = await space2Btn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const space2Btn2C = await space2Btn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const space2Btn2LH = await space2Btn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const space2Btn2M = await space2Btn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    await expect(space2Btn2W).toStrictEqual('20')
    await expect(space2Btn2H).toStrictEqual('32px')
    await expect(space2Btn2C).toStrictEqual('rgb(34, 84, 157)')
    await expect(space2Btn2LH).toStrictEqual('32px')
    await expect(space2Btn2M).toStrictEqual('0px 0px 0px 12px')

    await space2Btn2.click()

    await page.getByText('引用环境变量').click()

    const space1Env = await page.locator('nz-space').first()
    const space1EnvW = await space1Env.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const space1EnvH = await space1Env.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(space1EnvW).toStrictEqual('430px')
    await expect(space1EnvH).toStrictEqual('52px')

    const space1EnvInput = await page.getByPlaceholder('请输入环境变量').first()
    const space1EnvInputW = await space1EnvInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const space1EnvInputH = await space1EnvInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(space1EnvInputW).toStrictEqual('346px')
    await expect(space1EnvInputH).toStrictEqual('32px')

    const space1EnvA = await page.locator('#dynamic a').first()
    const space1EnvAH = await space1EnvA.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const space1EnvAC = await space1EnvA.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const space1EnvALH = await space1EnvA.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const space1EnvAFS = await space1EnvA.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    await expect(space1EnvAH).toStrictEqual('32px')
    await expect(space1EnvAC).toStrictEqual('rgb(34, 84, 157)')
    await expect(space1EnvALH).toStrictEqual('32px')
    await expect(space1EnvAFS).toStrictEqual('12px')

    await space1EnvA.click()
    const drawerListHeader = await page.locator('.drawer-list-header').first()
    const drawerListHeaderH = await drawerListHeader.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const drawerListHeaderW = await drawerListHeader.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const drawerListHeaderP = await drawerListHeader.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const drawerListHeaderM = await drawerListHeader.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    await expect(drawerListHeaderH).toStrictEqual('44px')
    await expect(drawerListHeaderW).toStrictEqual('636px')
    await expect(drawerListHeaderP).toStrictEqual('0px 0px 12px 0px')
    await expect(drawerListHeaderM).toStrictEqual('0px')

    const drawerListContent = await page.locator('.drawer-list-content').first()
    const drawerListContentW = await drawerListContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const drawerListContentP = await drawerListContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const drawerListContentM = await drawerListContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    await expect(drawerListContentW).toStrictEqual('636px')
    await expect(drawerListContentP).toStrictEqual('0px 0px 20px 0px')
    await expect(drawerListContentM).toStrictEqual('0px')

    const drawerPagination = await page.locator('.drawer-pagination').first()
    const drawerPaginationW = await drawerPagination.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const drawerPaginationP = await drawerPagination.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const drawerPaginationM = await drawerPagination.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    await expect(drawerPaginationW).toStrictEqual('636px')
    await expect(drawerPaginationP).toStrictEqual('0px 0px 20px 0px')
    await expect(drawerPaginationM).toStrictEqual('16px 0px 0px 0px')

    await page.getByRole('button', { name: '添加变量' }).click()

    const keyInputTd = await page.getByRole('cell', { name: '请输入KEY' })
    const keyInputTdVA = await keyInputTd.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('vertical-align'))
    await expect(keyInputTdVA).toStrictEqual('middle')

    const keyInput = await page.getByRole('cell', { name: '请输入KEY' }).locator('input')
    const keyInputH = await keyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const keyInputW = await keyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const keyInputBC = await keyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const keyInputBGC = await keyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    await expect(keyInputH).toStrictEqual('32px')
    await expect(keyInputW).toStrictEqual('268px')
    await expect(keyInputBC).toStrictEqual('rgba(0, 0, 0, 0)')
    await expect(keyInputBGC).toStrictEqual('rgba(0, 0, 0, 0)')

    const ValueInputTd = await page.getByRole('cell', { name: '请输入描述' })
    const ValueInputTdVA = await ValueInputTd.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('vertical-align'))
    await expect(ValueInputTdVA).toStrictEqual('middle')

    const valueInput = await page.getByRole('cell', { name: '请输入KEY' }).locator('input')
    const valueInputH = await valueInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const valueInputW = await valueInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const valueInputBC = await valueInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const valueInputBGC = await valueInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    await expect(valueInputH).toStrictEqual('32px')
    await expect(valueInputW).toStrictEqual('268px')
    await expect(valueInputBC).toStrictEqual('rgba(0, 0, 0, 0)')
    await expect(valueInputBGC).toStrictEqual('rgba(0, 0, 0, 0)')

    const tdBtn1 = await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').last().locator('button').first()
    const tdBtn1H = await tdBtn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const tdBtn1W = await tdBtn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const tdBtn1BC = await tdBtn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const tdBtn1BGC = await tdBtn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const tdBtn1M = await tdBtn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const tdBtn1P = await tdBtn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    await expect(tdBtn1H).toStrictEqual('32px')
    await expect(tdBtn1W).toStrictEqual('28px')
    await expect(tdBtn1BC).toStrictEqual('rgba(0, 0, 0, 0)')
    await expect(tdBtn1BGC).toStrictEqual('rgba(0, 0, 0, 0)')
    await expect(tdBtn1M).toStrictEqual('0px')
    await expect(tdBtn1P).toStrictEqual('0px 8px 0px 0px')

    const tdBtn2 = await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').last().locator('button').last()
    const tdBtn2H = await tdBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const tdBtn2W = await tdBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const tdBtn2BC = await tdBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const tdBtn2BGC = await tdBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const tdBtn2M = await tdBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const tdBtn2P = await tdBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    await expect(tdBtn2H).toStrictEqual('32px')
    await expect(tdBtn2W).toStrictEqual('28px')
    await expect(tdBtn2BC).toStrictEqual('rgba(0, 0, 0, 0)')
    await expect(tdBtn2BGC).toStrictEqual('rgba(0, 0, 0, 0)')
    await expect(tdBtn2M).toStrictEqual('0px')
    await expect(tdBtn2P).toStrictEqual('0px 0px 0px 8px')

    await page.locator('#dynamic a').click()
    await page.getByRole('button', { name: '添加变量' }).click()
    await page.getByRole('button', { name: '添加变量' }).dblclick()
    await page.getByPlaceholder('请输入KEY').nth(2).click()
    await page.getByRole('button', { name: '' }).nth(2).click()
    await page.getByText('parameter error').click()
    await page.getByPlaceholder('请输入KEY').nth(1).click()
    await page.getByPlaceholder('请输入KEY').nth(1).fill('test')
    await page.getByRole('button', { name: '' }).nth(1).first().click()
    await page.getByRole('button', { name: 'Close' }).click()
    await page.getByRole('button', { name: '取消' }).click()
  })
  it('点击新建上游，逐一填入必输项，点击保存', async () => {
    await page.getByRole('button', { name: '新建上游' }).click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').click()
    await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').fill('test')
    await page.getByPlaceholder('请输入主机名或IP:端口').click()
    await page.getByPlaceholder('请输入主机名或IP:端口').fill('test')
    await page.getByPlaceholder('请输入权重').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByPlaceholder('请输入权重').click()
    await page.getByPlaceholder('请输入权重').fill('2')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('eo-ng-select#discoveryName').click()
    await page.locator('eo-ng-option-item >> nth = 2').click()
    await page.locator('input[type="text"]').click()
    await page.locator('input[type="text"]').fill('test')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('eo-ng-select#discoveryName').click()
    await page.getByText('静态节点').click()
    await page.getByText('引用环境变量').click()
    await page.getByPlaceholder('请输入环境变量').click()
    await page.getByPlaceholder('请输入环境变量').fill('test')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('nz-space div').filter({ hasText: '引用环境变量' }).click()
    await page.locator('eo-ng-upstream-create').click()
    await page.getByText('保存 取消').click()
    await page.locator('span').filter({ hasText: '引用环境变量' }).click()
    await page.locator('eo-ng-upstream-create').click()
    await page.locator('#dynamic a').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByRole('button', { name: '取消' }).click()
  })
  it('点击上游列表中的某一项，进入上线管理，检查样式，点击面包屑返回列表页', async () => {
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 0').click()
    await page.getByText('上游管理 / 上线管理 /').click()

    const tab1 = await page.getByRole('link', { name: '上线管理' })
    const tab1FS = await tab1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const tab1C = await tab1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    await expect(tab1FS).toStrictEqual('14px')
    await expect(tab1C).toStrictEqual('rgb(34, 84, 157)')

    const listContent = await page.locator('.list-content')
    const listContentMT = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    const listContentPB = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-bottom'))

    const listTable = await page.locator('eo-ng-apinto-table')
    const listTableMT = await listTable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))

    expect(listContentMT).toStrictEqual('12px')
    expect(listTableMT).toStrictEqual('0px')
    expect(listContentPB).toStrictEqual('20px')

    const listTableTh1 = await page.locator('eo-ng-apinto-table tr >> nth = 0 >> th >> nth = 0')
    const listTableTh1Padding = await listTableTh1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTh1Padding).toStrictEqual('0px 12px')

    const listTableTh2 = await page.locator('eo-ng-apinto-table tr >> nth = 0 >> th >> nth = 1')
    const listTableTh2Padding = await listTableTh2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTh2Padding).toStrictEqual('0px 12px')

    if (await (await page.$$('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 5')).length === 2) {
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
    } else {
      const listTableIcon2 = await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 5 >> button')
      const listTableIcon2PL = await listTableIcon2.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('padding-left'))
      const listTableIcon2PR = await listTableIcon2.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('padding-right'))
      expect(listTableIcon2PL).toStrictEqual('0px')
      expect(listTableIcon2PR).toStrictEqual('0px')
    }

    const onlineText = await page.getByText('已上线')
    if (await onlineText && await onlineText.isVisible()) {
      const onlineTextC = await onlineText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      const onlineTextFW = await onlineText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-weight'))
      expect(onlineTextC).toStrictEqual('rgb(19, 137, 19)')
      expect(onlineTextFW).toStrictEqual('700')
    }

    const notgoonlineText = await page.getByText('未上线')
    if (await notgoonlineText && await notgoonlineText.isVisible()) {
      const notgoonlineTextC = await notgoonlineText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      const notgoonlineTextFW = await notgoonlineText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-weight'))
      expect(notgoonlineTextC).toStrictEqual('rgb(143, 142, 147)')
      expect(notgoonlineTextFW).toStrictEqual('700')
    }

    const toUpdateText = await page.getByText('待更新')
    if (await toUpdateText && await toUpdateText.isVisible()) {
      const toUpdateTextC = await toUpdateText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      const toUpdateTextFW = await toUpdateText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-weight'))
      expect(toUpdateTextC).toStrictEqual('rgb(3, 169, 244)')
      expect(toUpdateTextFW).toStrictEqual('700')
    }

    const offlineText = await page.getByText('已下线')
    if (await offlineText && await offlineText.isVisible()) {
      const offlineTextC = await offlineText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      const offlineTextFW = await offlineText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-weight'))
      expect(offlineTextC).toStrictEqual('rgb(143, 142, 147)')
      expect(offlineTextFW).toStrictEqual('700')
    }
    await page.locator('nz-breadcrumb').getByRole('link', { name: '上游管理' }).click()
  })
  it('点击上游列表中的某一项，进入上游管理页，点击tab切换至上游信息页，检查样式，修改并提交', async () => {
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 1').click()
    await page.getByRole('link', { name: '上游信息' }).click()

    // 上游名称输入框样式
    const nameInput = await page.locator('input#name')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const nameInputD = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(nameInputW).toStrictEqual('346px')
    await expect(nameInputH).toStrictEqual('32px')
    await expect(nameInputD).toStrictEqual('true')

    await page.locator('#desc').click()
    await page.locator('#desc').fill('testt')
    await page.getByRole('button', { name: '提交' }).click()
  })
  it('点击上游列表中的某一项的查看图标，进入上线管理页，检查tab；通过面包屑返回列表页', async () => {
    await page.getByRole('button', { name: '' }).first().click()
    await page.getByRole('link', { name: '上游信息' }).click()
    await page.getByRole('link', { name: '上线管理' }).click()
    await page.getByRole('link', { name: '上游信息' }).click()
    await page.locator('nz-breadcrumb').getByText('上游信息').click()
    await page.getByRole('link', { name: '上线管理' }).click()
    await page.getByRole('link', { name: '上游信息' }).click()
    await page.getByRole('link', { name: '上线管理' }).click()
    await page.locator('nz-breadcrumb').getByRole('link', { name: '上游管理' }).click()
    await page.getByRole('button', { name: '' }).last().click()
    const deleteBtn = page.getByRole('button', { name: '' }).last()
    const deleteBtnD = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    if (deleteBtnD !== 'true') {
      await deleteBtn.click()
      await page.getByText('该数据删除后将无法找回，请确认是否删除？').click()
      await page.locator('ant-modal-footer').getByRole('button', { name: '取消' }).click()
      await deleteBtn.click()
      await page.getByText('该数据删除后将无法找回，请确认是否删除？').click()
      await page.locator('ant-modal-footer').getByRole('button', { name: '确认' }).click()
    }
  })
})
