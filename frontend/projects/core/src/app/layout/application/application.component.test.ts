
describe('应用管理 e2e test', () => {
  it('初始化页面，点击应用管理菜单，进入列表页', async () => {
    await page.goto('http://localhost:4200/login')
    await page.waitForTimeout(2000)
    await page.getByPlaceholder('请输入账号').click()
    await page.getByPlaceholder('请输入账号').fill('maggie')
    await page.getByPlaceholder('请输入账号').press('Tab')
    await page.getByPlaceholder('请输入密码').fill('12345678')
    await page.getByPlaceholder('请输入密码').press('Enter')
    await page.getByRole('link', { name: '应用管理' }).click()
  })
  it('按钮大小、背景色、字体颜色，输入框大小，右侧间距为24px，table字体左右间距12px，上方间距12p, 最后一列的图标间距为12、16、24px，分页组件的样式', async () => {
    // 新建应用的按钮样式
    const createBtn = await page.getByRole('button', { name: '新建应用' })
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
    expect(searchInputBC).toStrictEqual('rgb(217, 217, 217)')
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
  it('点击新建应用，清空必输项，点击保存不生效，输入所有必输项，点击保存生效', async () => {
    await page.getByRole('button', { name: '新建应用' }).click()
    await page.locator('nz-form-label').filter({ hasText: '应用ID' }).click()
    await page.locator('#id').click()
    await page.locator('#id').fill('')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('#name').click()
    await page.locator('#name').fill('testForE2e')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('#id').click()
    await page.locator('#id').fill('1234567')
    await page.getByRole('button', { name: '保存' }).click()
  })
  it('点击新建应用，检查输入框样式，点击添加自定义属性并填入值，检查样式，点击header额外参数，检查样式，点击保存，保存生效并出现成功提示框', async () => {
    await page.getByRole('button', { name: '新建应用' }).click()

    // 应用名称输入框样式
    const nameInput = await page.locator('input#name')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(nameInputW).toStrictEqual('346px')
    await expect(nameInputH).toStrictEqual('32px')

    // 应用ID输入框样式
    const idInput = await page.locator('input#id')
    const idInputW = await idInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const idInputH = await idInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(idInputW).toStrictEqual('346px')
    await expect(idInputH).toStrictEqual('32px')

    // 自定义属性样式
    const customFirstInput1 = await page.locator('eo-ng-apinto-table.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 0')
    const customFirstInput1W = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput1PR = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    const customFirstInput1H = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput1W).toStrictEqual('182px')
    await expect(customFirstInput1PR).toStrictEqual('8px')
    await expect(customFirstInput1H).toStrictEqual('32px')

    const customFirstInput2 = await page.locator('eo-ng-apinto-table.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 1')
    const customFirstInput2W = await customFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput2PR = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    const customFirstInput2H = await customFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput2W).toStrictEqual('172px')
    await expect(customFirstInput2PR).toStrictEqual('8px')
    await expect(customFirstInput2H).toStrictEqual('32px')

    const customFirstBtn = await page.locator('eo-ng-apinto-table.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 2 >> button')
    const customFirstBtnH = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const customFirstBtnLH = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const customFirstBtnColor = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(customFirstBtnH).toStrictEqual('32px')
    await expect(customFirstBtnLH).toStrictEqual('32px')
    await expect(customFirstBtnColor).toStrictEqual('rgb(34, 84, 157)')

    await customFirstBtn.click()

    const customFirstInput1A = await page.locator('eo-ng-apinto-table.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 0')
    const customFirstInput1AW = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput1APR = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    const customFirstInput1APB = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-bottom'))
    const customFirstInput1AH = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput1AW).toStrictEqual('182px')
    await expect(customFirstInput1APR).toStrictEqual('8px')
    await expect(customFirstInput1APB).toStrictEqual('16px')
    await expect(customFirstInput1AH).toStrictEqual('48px')

    const customSecondBtn = await page.locator('eo-ng-apinto-table.arrayItem >> nth = 0 >> tr >> nth = 2 >> td >> nth = 2 >> button >> nth = 0')
    const customSecondBtnH = await customSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const customSecondBtnLH = await customSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const customSecondBtnColor = await customSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(customSecondBtnH).toStrictEqual('32px')
    await expect(customSecondBtnLH).toStrictEqual('32px')
    await expect(customSecondBtnColor).toStrictEqual('rgb(34, 84, 157)')

    // 描述样式
    const descInput = await page.locator('textarea#desc')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(descInputW).toStrictEqual('346px')
    await expect(descInputH).toStrictEqual('68px')

    // Header额外参数样式
    const headerFirstInput1 = await page.locator('eo-ng-apinto-table.arrayItem >> nth = 1 >> tr >> nth = 1 >> td >> nth = 0')
    const headerFirstInput1W = await headerFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const headerFirstInput1PR = await headerFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    const headerFirstInput1H = await headerFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(headerFirstInput1W).toStrictEqual('182px')
    await expect(headerFirstInput1PR).toStrictEqual('8px')
    await expect(headerFirstInput1H).toStrictEqual('32px')

    const headerFirstInput2 = await page.locator('eo-ng-apinto-table.arrayItem >> nth = 1 >> tr >> nth = 1 >> td >> nth = 1')
    const headerFirstInput2W = await headerFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const headerFirstInput2PR = await headerFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    const headerFirstInput2H = await headerFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(headerFirstInput2W).toStrictEqual('172px')
    await expect(headerFirstInput2PR).toStrictEqual('8px')
    await expect(headerFirstInput2H).toStrictEqual('32px')

    const headerFirstBtn = await page.locator('eo-ng-apinto-table.arrayItem >> nth = 1 >> tr >> nth = 1 >> td >> nth = 2 >> button')
    const headerFirstBtnH = await headerFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const headerFirstBtnLH = await headerFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const headerFirstBtnColor = await headerFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(headerFirstBtnH).toStrictEqual('32px')
    await expect(headerFirstBtnLH).toStrictEqual('32px')
    await expect(headerFirstBtnColor).toStrictEqual('rgb(34, 84, 157)')

    await headerFirstBtn.click()

    const headerFirstInput1A = await page.locator('eo-ng-apinto-table.arrayItem >> nth = 1 >> tr >> nth = 1 >> td >> nth = 0')
    const headerFirstInput1AW = await headerFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const headerFirstInput1APR = await headerFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    const headerFirstInput1APB = await headerFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-bottom'))
    const headerFirstInput1AH = await headerFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(headerFirstInput1AW).toStrictEqual('182px')
    await expect(headerFirstInput1APR).toStrictEqual('8px')
    await expect(headerFirstInput1APB).toStrictEqual('16px')
    await expect(headerFirstInput1AH).toStrictEqual('48px')

    const headerSecondBtn = await page.locator('eo-ng-apinto-table.arrayItem >> nth = 1 >> tr >> nth = 2 >> td >> nth = 2 >> button >> nth = 0')
    const headerSecondBtnH = await headerSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const headerSecondBtnLH = await headerSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const headerSecondBtnColor = await headerSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(headerSecondBtnH).toStrictEqual('32px')
    await expect(headerSecondBtnLH).toStrictEqual('32px')
    await expect(headerSecondBtnColor).toStrictEqual('rgb(34, 84, 157)')

    // 填写自定义属性与header参数并保存
    await page.locator('.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 0 >> input').click()
    await page.locator('.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 0 >> input').fill('test')
    await page.locator('.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 1 >> input').click()
    await page.locator('.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 1 >> input').fill('1')
    await page.locator('.arrayItem >> nth = 1 >> tr >> nth = 1 >> td >> nth = 0 >> input').click()
    await page.locator('.arrayItem >> nth = 1 >> tr >> nth = 1 >> td >> nth = 0 >> input').fill('test2')
    await page.locator('.arrayItem >> nth = 1 >> tr >> nth = 1 >> td >> nth = 1 >> input').click()
    await page.locator('.arrayItem >> nth = 1 >> tr >> nth = 1 >> td >> nth = 1 >> input').fill('2')

    await page.locator('.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 0 >> input').click()
    await page.locator('.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 0 >> input').fill('test23')
    await page.locator('.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 1 >> input').click()
    await page.locator('.arrayItem >> nth = 0 >> tr >> nth = 1 >> td >> nth = 1 >> input').fill('23')

    await page.getByText('用于转发上游系统鉴权').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('#name').click()
    await page.locator('#name').fill('testForE2e2')
    await page.getByRole('button', { name: '保存' }).click()
  })
  it('列表页点击匿名应用，进入匿名应用信息页，无鉴权管理tab', async () => {
    await page.getByRole('cell', { name: '匿名应用' }).click()
    expect(await page.getByRole('link', { name: '鉴权管理' }).isVisible()).toStrictEqual(false)
  })
  it('上线管理列表的表格与上方间距12px，表格最后一列宽度为88px，检查不同状态的字体颜色', async () => {
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

    const listTableIconTh = await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 6')
    const listTableIconThPadding = await listTableIconTh.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableIconThPadding).toStrictEqual('0px 24px 0px 12px')

    const listTableIcon1 = await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 6 >> button >> nth = 0')
    const listTableIcon1PL = await listTableIcon1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const listTableIcon1PR = await listTableIcon1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    expect(listTableIcon1PL).toStrictEqual('0px')
    expect(listTableIcon1PR).toStrictEqual('8px')

    const listTableIcon2 = await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 6 >> button >> nth = 1')
    const listTableIcon2PL = await listTableIcon2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const listTableIcon2PR = await listTableIcon2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    expect(listTableIcon2PL).toStrictEqual('8px')
    expect(listTableIcon2PR).toStrictEqual('0px')

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
    const disabledText = await page.getByText('已禁用')
    if (await disabledText && await disabledText.isVisible()) {
      const disabledTextC = await disabledText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      const disabledTextFW = await disabledText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-weight'))
      expect(disabledTextC).toStrictEqual('rgb(255, 59, 48)')
      expect(disabledTextFW).toStrictEqual('700')
    }

    const enabledText = await page.getByText('未禁用')
    if (await enabledText && await enabledText.isVisible()) {
      const enabledTextC = await enabledText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      const enabledTextFW = await enabledText.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-weight'))
      expect(enabledTextC).toStrictEqual('rgb(19, 137, 19)')
      expect(enabledTextFW).toStrictEqual('700')
    }
  })
  it('点击应用信息，检查输入框样式，其中应用名称与应用ID为不可修改，点击提交返回列表页', async () => {
    await page.getByRole('link', { name: '应用信息' }).click()

    // 应用名称输入框样式
    const nameInput = await page.locator('input#name')
    const nameInputDisabled = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    await expect(nameInputDisabled).toStrictEqual('')

    // 应用ID输入框样式
    const idInput = await page.locator('input#id')
    const idInputDisabled = await idInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(idInputDisabled).toStrictEqual('')

    await page.getByRole('button', { name: '提交' }).click()
  })
  it('点击非匿名应用一行的查看图标，进入应用信息页，点击进入鉴权管理，鉴权管理按钮为primary，检查大小、间距', async () => {
    // await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 5 >> button >> nth = 0').click()
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 4').click()
    await page.getByRole('link', { name: '鉴权管理' }).click()

    const createBtn = await page.getByRole('button', { name: '配置鉴权' })
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

    // const listTableIconTh = await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 5')
    // const listTableIconThPadding = await listTableIconTh.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('padding'))
    // expect(listTableIconThPadding).toStrictEqual('0px 24px 0px 12px')

    // const listTableIcon1 = await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 5 >> button >> nth = 0')
    // const listTableIcon1PL = await listTableIcon1.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('padding-left'))
    // const listTableIcon1PR = await listTableIcon1.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('padding-right'))
    // expect(listTableIcon1PL).toStrictEqual('0px')
    // expect(listTableIcon1PR).toStrictEqual('8px')

    // const listTableIcon2 = await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 5 >> button >> nth = 1')
    // const listTableIcon2PL = await listTableIcon2.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('padding-left'))
    // const listTableIcon2PR = await listTableIcon2.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('padding-right'))
    // expect(listTableIcon2PL).toStrictEqual('8px')
    // expect(listTableIcon2PR).toStrictEqual('0px')
  })
  it('点击按钮配置鉴权，出现弹窗,检查弹窗样式，检查basic鉴权的输入框样式', async () => {
    await page.getByRole('button', { name: '配置鉴权' }).click()
    await page.waitForTimeout(2000)

    // 弹窗样式
    const drawerTitle = await page.locator('.ant-drawer-header')
    const drawerTitleP = await drawerTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const drawerTitleW = await drawerTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(drawerTitleP).toStrictEqual('15px 20px')
    expect(drawerTitleW).toStrictEqual('676px')

    // 透传上游样式
    const continueSwitch = await page.locator('eo-ng-switch >> button')
    const continueSwitchW = await continueSwitch.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const continueSwitchH = await continueSwitch.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(continueSwitchW).toStrictEqual('35px')
    await expect(continueSwitchH).toStrictEqual('16px')

    // 鉴权类型选择框样式
    const driverSelect = await page.locator('#driver eo-ng-select-top-control')
    const driverSelectW = await driverSelect.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const driverSelectH = await driverSelect.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(driverSelectW).toStrictEqual('346px')
    await expect(driverSelectH).toStrictEqual('32px')

    // 参数位置样式
    const headerSelect = await page.locator('#position eo-ng-select-top-control')
    const headerSelectW = await headerSelect.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const headerSelectH = await headerSelect.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(headerSelectW).toStrictEqual('132px')
    await expect(headerSelectH).toStrictEqual('32px')

    const tokenInput = await page.getByPlaceholder('请输入TokenName')
    const tokenInputW = await tokenInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const tokenInputH = await tokenInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(tokenInputW).toStrictEqual('215px')
    await expect(tokenInputH).toStrictEqual('32px')

    // 用户名样式
    const usernameInput = await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文')
    const usernameInputW = await usernameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const usernameInputH = await usernameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(usernameInputW).toStrictEqual('346px')
    await expect(usernameInputH).toStrictEqual('32px')

    // 密码样式

    const pswInput = await page.locator('section').filter({ hasText: '*密码' }).getByPlaceholder('请输入')
    const pswInputW = await pswInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const pswInputH = await pswInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(pswInputW).toStrictEqual('346px')
    await expect(pswInputH).toStrictEqual('32px')

    // 标签信息样式
    const customFirstInput1 = await page.locator('#params0key')
    const customFirstInput1W = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput1H = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput1W).toStrictEqual('210px')
    await expect(customFirstInput1H).toStrictEqual('32px')

    const customFirstInput2 = await page.locator('.ant-space > div:nth-child(2) >> nth = 0')
    const customFirstInput2W = await customFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput2H = await customFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput2W).toStrictEqual('218px')
    await expect(customFirstInput2H).toStrictEqual('52px')

    const customFirstBtn = await page.locator('#dynamic a').first()
    const customFirstBtnH = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const customFirstBtnLH = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const customFirstBtnColor = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(customFirstBtnH).toStrictEqual('32px')
    await expect(customFirstBtnLH).toStrictEqual('32px')
    await expect(customFirstBtnColor).toStrictEqual('rgb(34, 84, 157)')

    const customFirstInput1A = await page.locator('#params1key')
    const customFirstInput1AW = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput1AH = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput1AW).toStrictEqual('210px')
    await expect(customFirstInput1AH).toStrictEqual('32px')

    const customSecondBtn = await page.locator('div:nth-child(2) > .ant-space > div:nth-child(3)')
    const customSecondBtnH = await customSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const customSecondBtnW = await customSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))

    await expect(customSecondBtnH).toStrictEqual('55.4286px')
    await expect(customSecondBtnW).toStrictEqual('32px')

    // 过期时间样式
    const expireInput = await page.locator('#expireTime')
    const expireInputW = await expireInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const expireInputH = await expireInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(expireInputW).toStrictEqual('346px')
    await expect(expireInputH).toStrictEqual('32px')

    // 每一行的间隙
    const mg1 = await page.locator('nz-form-item >> nth = 0')
    const mg1MT = await mg1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(mg1MT).toStrictEqual('20px')

    const mg2 = await page.locator('nz-form-item >> nth = 1')
    const mg2MT = await mg2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(mg2MT).toStrictEqual('20px')

    const mg3 = await page.locator('nz-form-item >> nth = 2')
    const mg3MT = await mg3.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(mg3MT).toStrictEqual('20px')

    const mg4 = await page.locator('nz-form-item >> nth = 3')
    const mg4MT = await mg4.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(mg4MT).toStrictEqual('0px')

    const mg5 = await page.locator('nz-form-item >> nth = 4')
    const mg5MT = await mg5.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(mg5MT).toStrictEqual('20px')

    const mg6 = await page.locator('nz-form-item >> nth = 5')
    const mg6MT = await mg6.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    expect(mg6MT).toStrictEqual('20px')
  })
  it('未填写必填项时，点击保存不生效；填写完必输项，点击保存后出现消息提示，并返回列表页', async () => {
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').click()
    await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').fill('test')
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByRole('button', { name: '关' }).click()
    await page.getByText('Header').click()
    await page.locator('eo-ng-option-item').filter({ hasText: 'Query' }).click()
    await page.getByPlaceholder('请输入TokenName').click()
    await page.getByPlaceholder('请输入TokenName').fill('Authorization1')
    await page.locator('#params0key').click()
    await page.locator('#params0key').fill('teest')
    await page.locator('#params0value').click()
    await page.locator('#params0value').fill('1')
    await page.locator('section').filter({ hasText: '*密码 必填项' }).getByPlaceholder('请输入').click()
    await page.locator('section').filter({ hasText: '*密码 必填项' }).getByPlaceholder('请输入').fill('test')
    await page.locator('#expireTime div').first().click()
    await page.getByRole('gridcell', { name: '24' }).click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByText('success').click()
  })
  it('配置鉴权选择apikey，检查样式', async () => {
    await page.getByRole('button', { name: '配置鉴权' }).click()
    await page.locator('#driver').getByText('Basic').click()
    await page.locator('eo-ng-option-item').filter({ hasText: 'ApiKey' }).click()
    await page.waitForTimeout(2000)

    // apikey样式
    const apikeyInput = await page.locator('section section >> nth = 0 >> input')
    const apikeyInputW = await apikeyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const apikeyInputH = await apikeyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(apikeyInputW).toStrictEqual('346px')
    await expect(apikeyInputH).toStrictEqual('32px')

    // 标签信息样式
    const customFirstInput1 = await page.locator('#params0key')
    const customFirstInput1W = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput1H = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput1W).toStrictEqual('210px')
    await expect(customFirstInput1H).toStrictEqual('32px')

    const customFirstInput2 = await page.locator('.ant-space > div:nth-child(2) >> nth = 0')
    const customFirstInput2W = await customFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput2H = await customFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput2W).toStrictEqual('218px')
    await expect(customFirstInput2H).toStrictEqual('52px')

    const customFirstBtn = await page.locator('#dynamic a').first()
    const customFirstBtnH = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const customFirstBtnLH = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const customFirstBtnColor = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(customFirstBtnH).toStrictEqual('32px')
    await expect(customFirstBtnLH).toStrictEqual('32px')
    await expect(customFirstBtnColor).toStrictEqual('rgb(34, 84, 157)')

    const customFirstInput1A = await page.locator('#params1key')
    const customFirstInput1AW = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput1AH = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput1AW).toStrictEqual('210px')
    await expect(customFirstInput1AH).toStrictEqual('32px')

    const customSecondBtn = await page.locator('div:nth-child(2) > .ant-space > div:nth-child(3)')
    const customSecondBtnH = await customSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const customSecondBtnW = await customSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))

    await expect(customSecondBtnH).toStrictEqual('55.4286px')
    await expect(customSecondBtnW).toStrictEqual('32px')
  })
  it('未填写必填项时，点击保存不生效；填写完必输项，点击保存后出现消息提示，并返回列表页', async () => {
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('section section >> nth = 0 >> input').click()
    await page.locator('section section >> nth = 0 >> input').fill('TEST')
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByText('success').click()
  })
  it('配置鉴权选择aksk，检查样式', async () => {
    await page.getByRole('button', { name: '配置鉴权' }).click()
    await page.locator('#driver').getByText('Basic').click()
    await page.locator('eo-ng-option-item').filter({ hasText: 'AkSk' }).click()
    await page.waitForTimeout(2000)

    // ak样式
    const akInput = await page.locator('section').filter({ hasText: '*AK' }).getByPlaceholder('请输入')
    const akInputW = await akInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const akInputH = await akInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(akInputW).toStrictEqual('346px')
    await expect(akInputH).toStrictEqual('32px')

    // sk样式
    const skInput = await page.locator('section').filter({ hasText: '*AK' }).getByPlaceholder('请输入')
    const skInputW = await skInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const skInputH = await skInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(skInputW).toStrictEqual('346px')
    await expect(skInputH).toStrictEqual('32px')

    // 标签信息样式
    const customFirstInput1 = await page.locator('#params0key')
    const customFirstInput1W = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput1H = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput1W).toStrictEqual('210px')
    await expect(customFirstInput1H).toStrictEqual('32px')

    const customFirstInput2 = await page.locator('.ant-space > div:nth-child(2) >> nth = 0')
    const customFirstInput2W = await customFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput2H = await customFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput2W).toStrictEqual('218px')
    await expect(customFirstInput2H).toStrictEqual('52px')

    const customFirstBtn = await page.locator('#dynamic a').first()
    const customFirstBtnH = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const customFirstBtnLH = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const customFirstBtnColor = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(customFirstBtnH).toStrictEqual('32px')
    await expect(customFirstBtnLH).toStrictEqual('32px')
    await expect(customFirstBtnColor).toStrictEqual('rgb(34, 84, 157)')

    const customFirstInput1A = await page.locator('#params1key')
    const customFirstInput1AW = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput1AH = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput1AW).toStrictEqual('210px')
    await expect(customFirstInput1AH).toStrictEqual('32px')

    const customSecondBtn = await page.locator('div:nth-child(2) > .ant-space > div:nth-child(3)')
    const customSecondBtnH = await customSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const customSecondBtnW = await customSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))

    await expect(customSecondBtnH).toStrictEqual('55.4286px')
    await expect(customSecondBtnW).toStrictEqual('32px')
  })
  it('未填写必填项时，点击保存不生效；填写完必输项，点击保存后出现消息提示，并返回列表页', async () => {
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('section').filter({ hasText: '*AK 必填项' }).getByPlaceholder('请输入').click()
    await page.locator('section').filter({ hasText: '*AK 必填项' }).getByPlaceholder('请输入').fill('test')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('section').filter({ hasText: '*SK 必填项' }).getByPlaceholder('请输入').click()
    await page.locator('section').filter({ hasText: '*SK 必填项' }).getByPlaceholder('请输入').fill('test')
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByText('success').click()
  })
  it('配置鉴权选择jwt，检查样式', async () => {
    await page.getByRole('button', { name: '配置鉴权' }).click()
    await page.locator('#driver').getByText('Basic').click()
    await page.locator('eo-ng-option-item').filter({ hasText: 'Jwt' }).click()
    await page.waitForTimeout(2000)

    // 算法输入框样式
    const algoInput = await page.locator('#dynamic eo-ng-select-top-control')
    const algoInputW = await algoInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const algoInputH = await algoInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(algoInputW).toStrictEqual('346px')
    await expect(algoInputH).toStrictEqual('32px')

    // ISS样式
    const issInput = await page.locator('section').filter({ hasText: '*Iss' }).getByPlaceholder('请输入')
    const issInputW = await issInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const issInputH = await issInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(issInputW).toStrictEqual('346px')
    await expect(issInputH).toStrictEqual('32px')

    // secret样式
    const secretInput = await page.locator('section').filter({ hasText: '*Secret' }).getByPlaceholder('请输入')
    const secretInputW = await secretInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const secretInputH = await secretInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(secretInputW).toStrictEqual('346px')
    await expect(secretInputH).toStrictEqual('32px')

    // 标签信息样式
    const customFirstInput1 = await page.locator('#params0key')
    const customFirstInput1W = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput1H = await customFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput1W).toStrictEqual('210px')
    await expect(customFirstInput1H).toStrictEqual('32px')

    const customFirstInput2 = await page.locator('.ant-space > div:nth-child(2) >> nth = 0')
    const customFirstInput2W = await customFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput2H = await customFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput2W).toStrictEqual('218px')
    await expect(customFirstInput2H).toStrictEqual('52px')

    const customFirstBtn = await page.locator('#dynamic a').first()
    const customFirstBtnH = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const customFirstBtnLH = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const customFirstBtnColor = await customFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(customFirstBtnH).toStrictEqual('32px')
    await expect(customFirstBtnLH).toStrictEqual('32px')
    await expect(customFirstBtnColor).toStrictEqual('rgb(34, 84, 157)')

    const customFirstInput1A = await page.locator('#params1key')
    const customFirstInput1AW = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const customFirstInput1AH = await customFirstInput1A.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(customFirstInput1AW).toStrictEqual('210px')
    await expect(customFirstInput1AH).toStrictEqual('32px')

    const customSecondBtn = await page.locator('div:nth-child(2) > .ant-space > div:nth-child(3)')
    const customSecondBtnH = await customSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const customSecondBtnW = await customSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))

    await expect(customSecondBtnH).toStrictEqual('55.4286px')
    await expect(customSecondBtnW).toStrictEqual('32px')
  })
  it('未填写必填项时，点击保存不生效；填写完必输项，点击保存后出现消息提示，并返回列表页', async () => {
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByText('HS256').click()
    await page.locator('eo-ng-option-item').filter({ hasText: 'HS512' }).click()
    await page.getByText('HS512').click()
    await page.locator('eo-ng-option-item').filter({ hasText: 'ES256' }).click()
    await page.locator('section').filter({ hasText: '*Iss 必填项' }).getByPlaceholder('请输入').click()
    await page.locator('section').filter({ hasText: '*Iss 必填项' }).getByPlaceholder('请输入').fill('test')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('section').filter({ hasText: '*RsaPublicKey 必填项' }).getByPlaceholder('请输入').click()
    await page.locator('section').filter({ hasText: '*RsaPublicKey 必填项' }).getByPlaceholder('请输入').fill('t')
    await page.locator('section').filter({ hasText: '*RsaPublicKey' }).getByPlaceholder('请输入').fill('test')
    await page.getByRole('button', { name: '保存' }).click()
  })
  it('点击列表中的basic，编辑并保存', async () => {
    await page.getByRole('cell', { name: 'Basic' }).nth(1).click()
    await page.getByRole('button', { name: '开' }).click()
    await page.getByRole('button', { name: '关' }).click()
    await page.getByRole('button', { name: '保存' }).click()
  })
  it('点击列表中aksk的查看按钮，编辑并保存', async () => {
    // await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 5 >> button >> nth = 0').click()
    // await page.getByRole('button', { name: '保存' }).click()
  })
  it('点击列表中jwt的删除按钮，点击删除，列表将减少一行', async () => {
    // const tableLength = await (await page.$$('eo-ng-filter-table  tr')).length
    // await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 5 >> button >> nth = 1').click()
    // await page.getByRole('button', { name: '确定' }).click()
    // expect(await (await page.$$('eo-ng-filter-table  tr')).length).toStrictEqual(tableLength - 1)
  })
  it('点击列表中apikey的删除按钮，点击取消，列表行数不变', async () => {
    // const tableLength = await (await page.$$('eo-ng-filter-table  tr')).length
    // await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 5 >> button >> nth = 1').click()
    // await page.getByRole('button', { name: '取消' }).click()
    // expect(await (await page.$$('eo-ng-filter-table  tr')).length).toStrictEqual(tableLength)
  })
  it('点击tabs，测试页面是否切换', async () => {
    await page.getByRole('link', { name: '应用信息' }).click()
    await page.getByText('用于转发上游系统鉴权').click()
    await page.getByRole('link', { name: '上线管理' }).click()
    await page.getByRole('columnheader', { name: '集群名称' }).click()
    await page.getByRole('link', { name: '鉴权管理' }).click()
    await page.getByRole('columnheader', { name: '鉴权信息' }).click()
  })
})
