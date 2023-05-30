describe('系统管理 e2e test', () => {
  it('初始化页面，点击系统管理-外部应用菜单，进入列表页', async () => {
    await page.goto('http://localhost:4200/login')
    await page.waitForTimeout(2000)
    await page.getByPlaceholder('请输入账号').click()
    await page.getByPlaceholder('请输入账号').fill('maggie')
    await page.getByPlaceholder('请输入账号').press('Tab')
    await page.getByPlaceholder('请输入密码').fill('12345678')
    await page.getByPlaceholder('请输入密码').press('Enter')
    await page.locator('eo-ng-menu-default').getByText('系统管理').click()
    await page.locator('eo-ng-menu-default').getByRole('link', { name: '外部应用' }).click()
  })
  it('面包屑API管理，字号14，字体颜色为主题色，与左侧距离12px，垂直居中；按钮样式；表格样式', async () => {
    // 面包屑样式 面包屑API管理，字号14，字体颜色为主题色，与左侧距离12px，垂直居中
    const headerBlock = await page.locator('.block_rl')
    const headerBlockAI = await headerBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('align-items'))
    expect(headerBlockAI).toStrictEqual('center')

    const breadcrumbItem = await page.locator('nz-breadcrumb-item').getByText('外部应用')
    const breadcrumbItemFS = await breadcrumbItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(breadcrumbItemFS).toStrictEqual('14px')
    const breadcrumbItemC = await breadcrumbItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(breadcrumbItemC).toStrictEqual('rgb(34, 84, 157)')

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

    // 表格的样式
    const listContent = await page.locator('.list-content')
    const listContentMT = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))

    const listTable = await page.locator('eo-ng-apinto-table')
    const listTableMT = await listTable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))

    expect(listContentMT).toStrictEqual('12px')
    expect(listTableMT).toStrictEqual('0px')

    const listTableTh1 = await page.getByRole('columnheader', { name: '应用名称' })
    const listTableTh1Padding = await listTableTh1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTh1Padding).toStrictEqual('0px 12px')

    const listTableTh2 = await page.getByRole('columnheader', { name: '应用ID' })
    const listTableTh2Padding = await listTableTh2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTh2Padding).toStrictEqual('0px 12px')

    const listTableIconTh = await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 7 ')
    const listTableIconThPadding = await listTableIconTh.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableIconThPadding).toStrictEqual('0px 24px 0px 12px')

    const listTableIcon1 = await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 7 >> button >> nth = 0')
    const listTableIcon1PL = await listTableIcon1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const listTableIcon1PR = await listTableIcon1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    expect(listTableIcon1PL).toStrictEqual('0px')
    expect(listTableIcon1PR).toStrictEqual('8px')

    const listTableIcon2 = await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 7  >> button >> nth = 1')
    const listTableIcon2PL = await listTableIcon2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const listTableIcon2PR = await listTableIcon2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    expect(listTableIcon2PL).toStrictEqual('8px')
    expect(listTableIcon2PR).toStrictEqual('8px')

    const listTableIcon3 = await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 7 >> button').last()
    const listTableIcon3PL = await listTableIcon3.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const listTableIcon3PR = await listTableIcon3.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    expect(listTableIcon3PL).toStrictEqual('8px')
    expect(listTableIcon3PR).toStrictEqual('0px')
  })
  it('点击新建应用，检查样式和面包屑；点击保存，清空应用ID；逐项填入数据项，提交成功；点击新建应用，点击取消，返回列表页', async () => {
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

    // 描述样式
    const descInput = await page.locator('textarea#desc')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(descInputW).toStrictEqual('346px')
    await expect(descInputH).toStrictEqual('68px')

    // 表单间距
    const formItem0 = await page.locator('nz-form-item >> nth = 0')
    const formItem0MB = await formItem0.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    await expect(formItem0MB).toStrictEqual('20px')

    const formItem1 = await page.locator('nz-form-item >> nth = 1')
    const formItem1MB = await formItem1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    await expect(formItem1MB).toStrictEqual('20px')

    const formItem2 = await page.locator('nz-form-item >> nth = 2')
    const formItem2MB = await formItem2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-bottom'))
    await expect(formItem2MB).toStrictEqual('0px')

    const formItem3 = await page.locator('nz-form-item >> nth = 3')
    const formItem3P = await formItem3.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('PADDING'))
    await expect(formItem3P).toStrictEqual('20px 0px')

    // 保存按钮样式
    const saveBtn = await page.getByRole('button', { name: '保存' })
    const saveBtnH = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const saveBtnW = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const saveBtnBG = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const saveBtnBC = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const saveBtnFS = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))

    expect(saveBtnH).toStrictEqual('32px')
    expect(saveBtnW).toStrictEqual('54px')
    expect(saveBtnBG).toStrictEqual('rgb(34, 84, 157)')
    expect(saveBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(saveBtnFS).toStrictEqual('14px')

    // 取消按钮样式
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
    const cancleBtnML = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    expect(cancleBtnH).toStrictEqual('32px')
    expect(cancleBtnW).toStrictEqual('54px')
    expect(cancleBtnBG).toStrictEqual('rgb(255, 255, 255)')
    expect(cancleBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(cancleBtnFS).toStrictEqual('14px')
    expect(cancleBtnML).toStrictEqual('12px')

    // 验证操作
    await page.getByText('应用ID').click()
    const randomId:string = await page.locator('#id').inputValue()
    await page.locator('#id').click()
    await page.locator('#id').fill('')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('nz-form-item').filter({ hasText: '应用名称必填项' }).getByRole('alert').click()
    await page.locator('nz-form-item').filter({ hasText: '应用ID必填项' }).getByRole('alert').click()
    await page.locator('#name').click()
    await page.locator('#name').fill('test')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('#id').click()
    await page.locator('#id').fill(randomId)
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByRole('button', { name: '新建应用' }).click()
    await page.getByRole('button', { name: '取消' }).click()
  })
  it('点击表格中某一栏，进入编辑页面，检查面包屑和应用ID禁用；修改ID，点击提交，返回列表页', async () => {
    await page.waitForTimeout(700)

    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 0').click()

    // 检查面包屑
    await page.getByText('外部应用 / 外部应用详情 /').click()

    const breadcrumb1 = await page.locator('nz-breadcrumb-item').filter({ hasText: '外部应用 /' }).locator('a')
    const breadcrumb1C = await breadcrumb1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    const breadcrumb2 = await page.getByText('外部应用详情')
    const breadcrumb2C = await breadcrumb2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(breadcrumb1C).toStrictEqual('rgb(153, 153, 153)')
    expect(breadcrumb2C).toStrictEqual('rgb(34, 84, 157)')

    // 应用ID输入框样式
    const idInput = await page.locator('input#id')
    const idInputDisabled = await idInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(idInputDisabled).toStrictEqual('')

    // 修改名称、描述并提交
    await page.locator('#name').click()
    await page.locator('#name').fill('test11')
    await page.locator('#desc').click()
    await page.locator('#desc').fill('testdesc1')
    await page.getByRole('button', { name: '提交' }).click()
  })
  it('点击表格中某一栏的查看图标，点击取消，返回列表页', async () => {
    await page.waitForTimeout(700)
    // await page.locator('eo-ng-apinto-table tr >> nth =2 >> td >> nth = 7 >> button >> nth = 1').click()
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 0').click()
    await page.getByRole('button', { name: '取消' }).click()
  })
  it('点击禁用开关，出现消息弹窗；点击刷新token，出现消息弹窗；点击复制token，出现成功复制token；点击删除，测试删除弹窗', async () => {
    await page.locator('eo-ng-switch >> nth = 0').click()
    await page.locator('.ant-message').click()

    // await page.locator('eo-ng-apinto-table tr >> nth =2 >> td >> nth = 7 >> button >> nth = 0').click()
    // await page.locator('.ant-message').click()

    // await page.locator('eo-ng-apinto-table tr >> nth =2 >> td >> nth = 7 >> button >> nth = 1').click()
    // await page.locator('.ant-message').click()

    // const tableLength = await (await page.$$('eo-ng-apinto-table tr')).length
    // await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 7 >> button >> nth = 3').click()
    // await page.getByRole('button', { name: '确定' }).click()
    // expect(await (await page.$$('eo-ng-apinto-table  tr')).length).toStrictEqual(tableLength - 1)
    // const tableLength = await (await page.$$('eo-ng-apinto-table tr')).length
    // await page.locator('eo-ng-apinto-table tr >> nth = 1 >> td >> nth = 7 >> button >> nth = 3').click()
    // await page.getByRole('button', { name: '取消' }).click()
    // expect(await (await page.$$('eo-ng-apinto-table  tr')).length).toStrictEqual(tableLength)
  })
})
