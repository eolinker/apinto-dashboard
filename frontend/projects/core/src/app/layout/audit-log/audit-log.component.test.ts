describe('审计日志 e2e test', () => {
  it('初始化进入审计日志页面', async () => {
    await page.goto('http://localhost:4200/login')
    await page.waitForTimeout(2000)
    await page.getByPlaceholder('请输入账号').click()
    await page.getByPlaceholder('请输入账号').fill('maggie')
    await page.getByPlaceholder('请输入账号').press('Tab')
    await page.getByPlaceholder('请输入密码').fill('12345678')
    await page.getByPlaceholder('请输入密码').press('Enter')
    await page.getByRole('link', { name: '审计日志' }).click()
  })

  it('检查面包屑、页面样式', async () => {
    // 面包屑样式 面包屑审计日志，字号14，字体颜色为主题色
    const breadcrumbItem = await page.locator('nz-breadcrumb-item').getByText('审计日志')
    const breadcrumbItemFS = await breadcrumbItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(breadcrumbItemFS).toStrictEqual('14px')
    const breadcrumbItemC = await breadcrumbItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(breadcrumbItemC).toStrictEqual('rgb(34, 84, 157)')

    // 输入框样式
    const groupInput1 = await page.getByText('操作类型 请选择 操作对象 请选择 搜索内容')
    const groupInput1ML = await groupInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    await expect(groupInput1ML).toStrictEqual('16px')

    const operTypeInput = await page.locator('eo-ng-select').first()
    const operTypeInputW = await operTypeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const operTypeInputH = await operTypeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const operTypeInputML = await operTypeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    await expect(operTypeInputW).toStrictEqual('254px')
    await expect(operTypeInputH).toStrictEqual('32px')
    await expect(operTypeInputML).toStrictEqual('12px')

    const labelTarget = await page.locator('label').filter({ hasText: '操作对象' })
    const labelTargetML = await labelTarget.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    await expect(labelTargetML).toStrictEqual('24px')

    const operTargetInput = await page.locator('eo-ng-select >> nth = 1')
    const operTargetInputW = await operTargetInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const operTargetInputH = await operTargetInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const operTargetInputML = await operTargetInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    await expect(operTargetInputW).toStrictEqual('254px')
    await expect(operTargetInputH).toStrictEqual('32px')
    await expect(operTargetInputML).toStrictEqual('12px')

    const labelSearch = await page.locator('label').filter({ hasText: '搜索内容' })
    const labelSearchML = await labelSearch.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    await expect(labelSearchML).toStrictEqual('24px')

    const searchInput = await page.locator('eo-ng-input-group').first()
    const searchInputW = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const searchInputH = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const searchInputML = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    await expect(searchInputW).toStrictEqual('254px')
    await expect(searchInputH).toStrictEqual('32px')
    await expect(searchInputML).toStrictEqual('12px')

    const groupInput2 = await page.getByText('操作类型 请选择 操作对象 请选择 搜索内容')
    const groupInput2ML = await groupInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    await expect(groupInput2ML).toStrictEqual('16px')

    const operTimeInput = await page.locator('nz-range-picker')
    const operTimeInputW = await operTimeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const operTimeInputH = await operTimeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const operTimeInputML = await operTimeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    await expect(operTimeInputW).toStrictEqual('254px')
    await expect(operTimeInputH).toStrictEqual('32px')
    await expect(operTimeInputML).toStrictEqual('12px')

    // 查询按钮样式
    const searchBtn = await page.getByRole('button', { name: '查询' })
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
    const searchBtnML = await searchBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    expect(searchBtnH).toStrictEqual('32px')
    expect(searchBtnW).toStrictEqual('54px')
    expect(searchBtnBG).toStrictEqual('rgb(34, 84, 157)')
    expect(searchBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(searchBtnFS).toStrictEqual('14px')
    expect(searchBtnML).toStrictEqual('12px')

    // 取消按钮样式
    const resetBtn = await page.getByRole('button', { name: '重置' })
    const resetBtnH = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const resetBtnW = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const resetBtnBG = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const resetBtnBC = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const resetBtnFS = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const resetBtnML = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    expect(resetBtnH).toStrictEqual('32px')
    expect(resetBtnW).toStrictEqual('54px')
    expect(resetBtnBG).toStrictEqual('rgb(255, 255, 255)')
    expect(resetBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(resetBtnFS).toStrictEqual('14px')
    expect(resetBtnML).toStrictEqual('24px')

    // 整体布局
    const header = await page.locator('.list-header')
    const headerPT = await header.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-top'))
    expect(headerPT).toStrictEqual('12px')

    const content = await page.locator('.list-content')
    const contentMT = await content.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    const contentPB = await content.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-bottom'))
    expect(contentMT).toStrictEqual('12px')
    expect(contentPB).toStrictEqual('20px')

    const pagination = await page.locator('.mg_pagination_t')
    const paginationPT = await pagination.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    expect(paginationPT).toStrictEqual('16px')

    // 表格样式
    const listTable = await page.locator('eo-ng-apinto-table')
    const listTableMT = await listTable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))

    expect(listTableMT).toStrictEqual('0px')

    const listTableTh1 = await page.getByRole('columnheader', { name: '用户名' })
    const listTableTh1Padding = await listTableTh1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTh1Padding).toStrictEqual('0px 12px')

    const listTableTh2 = await page.getByRole('columnheader', { name: '操作类型' })
    const listTableTh2Padding = await listTableTh2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTh2Padding).toStrictEqual('0px 12px')

    // const listTableIconTh = await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 5 ')
    // const listTableIconThPadding = await listTableIconTh.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('padding'))
    // expect(listTableIconThPadding).toStrictEqual('0px 24px 0px 12px')

    // const listTableIcon1 = await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 5 >> button >> nth = 0')
    // const listTableIcon1PL = await listTableIcon1.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('padding-left'))
    // const listTableIcon1PR = await listTableIcon1.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('padding-right'))
    // expect(listTableIcon1PL).toStrictEqual('0px')
    // expect(listTableIcon1PR).toStrictEqual('0px')
  })
  it('点击某一行，检查弹窗样式，点击关闭', async () => {
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 0 ').click()
    await page.getByText('日志详情').click()
    const drawerTable = await page.locator('.drawer-table')
    const drawerTableP = await drawerTable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const drawerTableM = await drawerTable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    expect(drawerTableP).toStrictEqual('0px')
    expect(drawerTableM).toStrictEqual('0px auto')

    const drawerListTable = await page.locator('.drawer-list-content')
    const drawerListTableP = await drawerListTable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(drawerListTableP).toStrictEqual('0px 0px 20px')

    // 取消按钮样式
    const resetBtn = await page.getByRole('button', { name: '取消' })
    const resetBtnH = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const resetBtnW = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const resetBtnBG = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const resetBtnBC = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const resetBtnFS = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const resetBtnM = await resetBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))

    expect(resetBtnH).toStrictEqual('32px')
    expect(resetBtnW).toStrictEqual('54px')
    expect(resetBtnBG).toStrictEqual('rgb(255, 255, 255)')
    expect(resetBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(resetBtnFS).toStrictEqual('14px')
    expect(resetBtnM).toStrictEqual('0px')

    await page.getByRole('button', { name: '取消' }).click()
  })
  it('点击某一行的查看图标，点击取消', async () => {
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 0 ').click()
    // await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 5 >> button ').click()

    await page.getByRole('button', { name: 'Close' }).click()
  })
  it('测试查询与重置功能', async () => {
    await page.locator('eo-ng-select-top-control').first().click()
    await page.locator('.cdk-overlay-pane').getByText('新建').click()
    await page.locator('nz-select-item').filter({ hasText: '新建' }).click()
    await page.locator('.cdk-overlay-pane').getByText('编辑').click()
    await page.locator('nz-select-item').filter({ hasText: '编辑' }).click()
    await page.locator('.cdk-overlay-pane').getByText('删除').click()
    await page.locator('nz-select-item').filter({ hasText: '删除' }).click()
    await page.locator('.cdk-overlay-pane').getByText('发布').click()
    await page.locator('nz-select-item').filter({ hasText: '发布' }).click()
    await page.locator('eo-ng-select').filter({ hasText: '发布' }).locator('nz-select-clear svg').click()
    await page.locator('eo-ng-select-top-control').nth(1).click()
    await page.locator('eo-ng-option-item').filter({ hasText: 'API' }).first().click()
    await page.locator('nz-select-clear svg').click()
    await page.locator('eo-ng-input-group').click()
    await page.getByPlaceholder('请输入用户名或请求IP').click()
    await page.getByPlaceholder('请输入用户名或请求IP').fill('TET')
    await page.getByPlaceholder('开始日期').click()
    await page.getByRole('row', { name: '28 29 30 1 2 3 4' }).getByText('2').nth(2).click()
    await page.getByText('00').first().click()
    await page.getByText('21').nth(3).click()
    await page.getByRole('listitem').filter({ hasText: '确定' }).click()
    await page.getByRole('row', { name: '26 27 28 29 30 31 1' }).getByText('30').click()
    await page.getByText('00').nth(1).click()
    await page.getByText('21').nth(3).click()
    await page.getByRole('button', { name: '确定' }).click()
    await page.getByRole('button', { name: '查询' }).click()

    await page.getByRole('button', { name: '重置' }).click()
    await page.locator('eo-ng-select-top-control').first().click()
    await page.locator('eo-ng-option-item').filter({ hasText: '编辑' }).click()
    await page.locator('eo-ng-select-top-control').filter({ hasText: '请选择' }).click()
    await page.locator('.cdk-overlay-pane').getByText('上游').click()
    await page.locator('eo-ng-input-group').click()
    await page.getByText('操作类型编辑操作对象上游搜索内容操作时间 重置 查询').click()
    await page.getByPlaceholder('开始日期').click()
    await page.getByRole('row', { name: '28 29 30 1 2 3 4' }).getByText('2').nth(2).click()
    await page.getByText('00').first().click()
    await page.getByText('00').nth(2).click()
    await page.getByText('00').nth(3).click()
    await page.getByText('21').nth(1).click()
    await page.getByText('01').first().click()
    await page.getByText('00').nth(1).click()
    await page.getByText('01').first().click()
    await page.getByRole('button', { name: '确定' }).click()
    await page.getByText('23').nth(1).click()
    await page.getByRole('button', { name: '确定' }).click()
    await page.getByRole('button', { name: '查询' }).click()
    await page.locator('eo-ng-select-top-control').filter({ hasText: '上游' }).click()
    await page.locator('eo-ng-option-item').filter({ hasText: '上游' }).click()
    await page.getByRole('button', { name: '查询' }).click()
  })
})
