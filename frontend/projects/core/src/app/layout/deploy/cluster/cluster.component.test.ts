describe('网关集群 e2e test', () => {
  it('初始化页面，点击基础设施-网关集群菜单，进入列表页', async () => {
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
    // 布局
    const listHeader = await page.locator('.list-header')
    const listHeaderMT = await listHeader.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-top'))
    expect(listHeaderMT).toStrictEqual('12px')

    // 新建集群的按钮样式
    const createBtn = await page.getByRole('button', { name: '新建集群' })
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
  it('点击新建集群，检查页面样式，输入集群地址，点击测试，出现表格，检查样式', async () => {
    await page.getByRole('button', { name: '新建集群' }).click()

    // 集群名称输入框样式
    const nameInput = await page.getByPlaceholder('英文数字下划线组合，首字母必须为英文')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(nameInputW).toStrictEqual('346px')
    await expect(nameInputH).toStrictEqual('32px')

    // 环境输入框样式
    const envInput = await page.locator('eo-ng-select')
    const envInputW = await envInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const envInputH = await envInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(envInputW).toStrictEqual('346px')
    await expect(envInputH).toStrictEqual('32px')

    // 描述输入框样式
    const descInput = await page.locator('textarea')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(descInputW).toStrictEqual('346px')
    await expect(descInputH).toStrictEqual('68px')

    // 集群地址输入框样式
    const addrInput = await page.locator('nz-form-control').filter({ hasText: '测试' }).getByPlaceholder('请输入')
    const addrInputW = await addrInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const addrInputH = await addrInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(addrInputW).toStrictEqual('346px')
    await expect(addrInputH).toStrictEqual('32px')

    await page.locator('nz-form-control').filter({ hasText: '测试' }).getByPlaceholder('请输入').click()
    await page.locator('nz-form-control').filter({ hasText: '测试' }).getByPlaceholder('请输入').fill('http://172.17.0.1:9400')
    await page.getByRole('button', { name: '测试' }).click()
  })
  it('点击保存；逐一填入必填项，点击保存，修改集群地址，点击保存，只有通过测试的集群地址才可保存，如出现消息弹窗且未保存成功，则点击取消，返回列表页', async () => {
    await page.getByRole('button', { name: '测试' }).click()
    await page.locator('nz-form-control').filter({ hasText: '测试 必填项' }).getByPlaceholder('请输入').click()
    await page.locator('nz-form-control').filter({ hasText: '测试 必填项' }).getByPlaceholder('请输入').fill('http://172.0.0.1:9400')
    await page.getByRole('button', { name: '测试' }).click()
    await page.getByText('node addr http://172.0.0.1:9400 can not be connected').click()
    await page.locator('nz-form-control').filter({ hasText: '测试' }).getByPlaceholder('请输入').click()
    await page.locator('nz-form-control').filter({ hasText: '测试' }).getByPlaceholder('请输入').fill('http://172.17.0.1')
    await page.locator('nz-form-control').filter({ hasText: '测试 集群地址输入错误，请重新输入' }).getByPlaceholder('请输入').fill('http://172.17.0.1:9400')
    await page.getByRole('button', { name: '测试' }).click()
    await page.getByRole('columnheader', { name: '名称' }).click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByPlaceholder('英文数字下划线组合，首字母必须为英文').click()
    await page.getByPlaceholder('英文数字下划线组合，首字母必须为英文').fill('TESTA')
    await page.locator('nz-form-control').filter({ hasText: '测试' }).getByPlaceholder('请输入').click()
    await page.locator('nz-form-control').filter({ hasText: '测试' }).getByPlaceholder('请输入').fill('http://172.17.0.1:940')
    await page.getByText('集群地址需要通过测试').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('nz-form-control').filter({ hasText: '测试 集群地址需要通过测试' }).getByPlaceholder('请输入').click()
    await page.locator('nz-form-control').filter({ hasText: '测试 集群地址需要通过测试' }).getByPlaceholder('请输入').fill('http://172.17.0.1')
    await page.getByText('集群地址输入错误，请重新输入').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('nz-form-control').filter({ hasText: '测试 集群地址输入错误，请重新输入' }).getByPlaceholder('请输入').click()
    await page.locator('nz-form-control').filter({ hasText: '测试 集群地址输入错误，请重新输入' }).getByPlaceholder('请输入').fill('http://172.17.0.1:9400')
    await page.getByRole('button', { name: '保存' }).click()
    if (await page.getByText('集群已有这个节点信息') && await page.getByText('集群已有这个节点信息').isVisible()) {
      await page.getByRole('button', { name: '取消' }).click()
      await page.getByText('状态').click()
    }
    await page.getByText('状态').click()
  })
  it('点击某一个集群，进入信息详情页，当前处于环境变量tab页,检查样式', async () => {
    await page.locator('eo-ng-apinto-table tr').last().locator('td').first().click()
    const envTab = await page.getByRole('tab', { name: '环境变量' }).getByRole('link', { name: '环境变量' })
    const envTabC = await envTab.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(envTabC).toStrictEqual('rgb(34, 84, 157)')

    const certTab = await page.getByRole('tab', { name: '证书管理' })
    const certTabC = await certTab.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(certTabC).toStrictEqual('rgb(0, 0, 0)')

    const nodeTab = await page.getByRole('tab', { name: '网关节点' })
    const nodeTabC = await nodeTab.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(nodeTabC).toStrictEqual('rgb(0, 0, 0)')

    const confTab = await page.getByRole('tab', { name: '配置管理' })
    const confTabC = await confTab.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(confTabC).toStrictEqual('rgb(0, 0, 0)')

    const createEnvBtn = await page.getByRole('button', { name: '新建配置' })
    const createEnvBtnM = await createEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const createEnvBtnP = await createEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const createEnvBtnBC = await createEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const createEnvBtnBGC = await createEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const createEnvBtnC = await createEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(createEnvBtnM).toStrictEqual('0px 0px 0px 12px')
    expect(createEnvBtnP).toStrictEqual('0px 12px')
    expect(createEnvBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(createEnvBtnBGC).toStrictEqual('rgb(34, 84, 157)')
    expect(createEnvBtnC).toStrictEqual('rgb(255, 255, 255)')

    const publishEnvBtn = await page.getByRole('button', { name: '发布' })
    const publishEnvBtnM = await publishEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const publishEnvBtnP = await publishEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const publishEnvBtnBC = await publishEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const publishEnvBtnBGC = await publishEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const publishEnvBtnC = await publishEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(publishEnvBtnM).toStrictEqual('0px 0px 0px 12px')
    expect(publishEnvBtnP).toStrictEqual('0px 12px')
    expect(publishEnvBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(publishEnvBtnBGC).toStrictEqual('rgb(34, 84, 157)')
    expect(publishEnvBtnC).toStrictEqual('rgb(255, 255, 255)')

    const updateEnvBtn = await page.getByRole('button', { name: '同步配置' })
    const updateEnvBtnM = await updateEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const updateEnvBtnP = await updateEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const updateEnvBtnBC = await updateEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const updateEnvBtnBGC = await updateEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const updateEnvBtnC = await updateEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(updateEnvBtnM).toStrictEqual('0px 0px 0px 12px')
    expect(updateEnvBtnP).toStrictEqual('0px 12px')
    expect(updateEnvBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(updateEnvBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(updateEnvBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')

    const updateHistoryBtn = await page.getByRole('button', { name: '更改历史' })
    const updateHistoryBtnM = await updateHistoryBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const updateHistoryBtnP = await updateHistoryBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const updateHistoryBtnBC = await updateHistoryBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const updateHistoryBtnBGC = await updateHistoryBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const updateHistoryBtnC = await updateHistoryBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(updateHistoryBtnM).toStrictEqual('0px')
    expect(updateHistoryBtnP).toStrictEqual('0px 12px')
    expect(updateHistoryBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(updateHistoryBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(updateHistoryBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')

    const publishHistoryBtn = await page.getByRole('button', { name: '发布历史' })
    const publishHistoryBtnM = await publishHistoryBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const publishHistoryBtnP = await publishHistoryBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const publishHistoryBtnBC = await publishHistoryBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const publishHistoryBtnBGC = await publishHistoryBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const publishHistoryBtnC = await publishHistoryBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(publishHistoryBtnM).toStrictEqual('0px 24px 0px 12px')
    expect(publishHistoryBtnP).toStrictEqual('0px 12px')
    expect(publishHistoryBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(publishHistoryBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(publishHistoryBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')

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
  it('点击新建配置，检查样式，逐一填入必填项直至保存成功；点击新建配置，点击取消', async () => {
    // key输入框样式
    const nameInput = await page.getByPlaceholder('英文、下划线及圆点组合')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const nameInputD = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(nameInputW).toStrictEqual('346px')
    await expect(nameInputH).toStrictEqual('32px')
    await expect(nameInputD).toStrictEqual('false')

    // value输入框样式
    const envInput = await page.locator('.ant-drawer-body textarea')
    const envInputW = await envInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const envInputH = await envInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const envInputD = await envInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(envInputW).toStrictEqual('346px')
    await expect(envInputH).toStrictEqual('68px')
    await expect(envInputD).toStrictEqual('false')

    // 描述输入框样式
    const descInput = await page.locator('textarea')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const descInputD = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(descInputW).toStrictEqual('346px')
    await expect(descInputH).toStrictEqual('68px')
    await expect(descInputD).toStrictEqual('false')

    const publishEnvBtn = await page.locator('.ant-drawer-body').getByRole('button', { name: '保存' })
    const publishEnvBtnM = await publishEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const publishEnvBtnP = await publishEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const publishEnvBtnBC = await publishEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const publishEnvBtnBGC = await publishEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const publishEnvBtnC = await publishEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(publishEnvBtnM).toStrictEqual('0px')
    expect(publishEnvBtnP).toStrictEqual('0px')
    expect(publishEnvBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(publishEnvBtnBGC).toStrictEqual('rgb(34, 84, 157)')
    expect(publishEnvBtnC).toStrictEqual('rgb(255, 255, 255)')

    const updateEnvBtn = await page.locator('.ant-drawer-body').getByRole('button', { name: '取消' })
    const updateEnvBtnM = await updateEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const updateEnvBtnP = await updateEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const updateEnvBtnBC = await updateEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const updateEnvBtnBGC = await updateEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const updateEnvBtnC = await updateEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(updateEnvBtnM).toStrictEqual('0px 0px 0px 12px')
    expect(updateEnvBtnP).toStrictEqual('0px 12px')
    expect(updateEnvBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(updateEnvBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(updateEnvBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')

    await page.locator('nz-form-control').filter({ hasText: '注意：隐藏字符(空格、换行符、制表符Tab)容易导致配置出错，如果需要检测Value中隐藏字符，请点击检测隐藏字符' }).getByPlaceholder('请输入').click()
    await page.locator('nz-form-control').filter({ hasText: '注意：隐藏字符(空格、换行符、制表符Tab)容易导致配置出错，如果需要检测Value中隐藏字符，请点击检测隐藏字符' }).getByPlaceholder('请输入').fill('   test  ')
    await page.locator('nz-form-control').filter({ hasText: '注意：隐藏字符(空格、换行符、制表符Tab)容易导致配置出错，如果需要检测Value中隐藏字符，请点击检测隐藏字符' }).getByPlaceholder('请输入').press('Enter')
    await page.locator('nz-form-control').filter({ hasText: '注意：隐藏字符(空格、换行符、制表符Tab)容易导致配置出错，如果需要检测Value中隐藏字符，请点击检测隐藏字符' }).getByPlaceholder('请输入').fill('   test  \n123  rfgasd')
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByText('检测隐藏字符').click()
    await page.getByText('#空格##空格##空格#test#空格##空格##换行符#123#空格##空格#rfgasd').click()
    const detectTip = await page.getByText('#空格##空格##空格#test#空格##空格##换行符#123#空格##空格#rfgasd')
    const detectTipBGC = await detectTip.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const detectTipC = await detectTip.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const detectTipFS = await detectTip.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    const detectTipLH = await detectTip.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const detectTipM = await detectTip.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const detectTipP = await detectTip.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const detectTipW = await detectTip.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const detectTipWB = await detectTip.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('word-break'))

    expect(detectTipBGC).toStrictEqual('rgb(248, 248, 250)')
    expect(detectTipC).toStrictEqual('rgb(153, 153, 153)')
    expect(detectTipFS).toStrictEqual('12px')
    expect(detectTipLH).toStrictEqual('20px')
    expect(detectTipM).toStrictEqual('0px')
    expect(detectTipP).toStrictEqual('0px 0px 0px 12px')
    expect(detectTipW).toStrictEqual('346px')
    expect(detectTipWB).toStrictEqual('break-all')

    const detectedText = await page.locator('.detected-symbol').first()
    const detectedTextC = await detectedText.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(detectedTextC).toStrictEqual('rgb(248, 248, 250)')

    await page.getByPlaceholder('英文、下划线及圆点组合').click()
    await page.getByPlaceholder('英文、下划线及圆点组合').fill('testForeee')
    await page.locator('nz-form-item').filter({ hasText: '描述' }).getByPlaceholder('请输入').click()
    await page.locator('nz-form-item').filter({ hasText: '描述' }).getByPlaceholder('请输入').fill('desc')
    await page.getByRole('button', { name: '保存' }).click()
  })
  it('点击列表中某一个环境变量，检查KEY与描述为不可修改，修改value并提交', async () => {
    await page.getByRole('cell', { name: 'testForeee' }).click()
    // key输入框样式
    const nameInput = await page.getByPlaceholder('英文、下划线及圆点组合')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const nameInputD = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(nameInputW).toStrictEqual('346px')
    await expect(nameInputH).toStrictEqual('32px')
    await expect(nameInputD).toStrictEqual('true')

    // value输入框样式
    const envInput = await page.locator('.ant-drawer-body textarea')
    const envInputW = await envInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const envInputH = await envInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const envInputD = await envInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(envInputW).toStrictEqual('346px')
    await expect(envInputH).toStrictEqual('68px')
    await expect(envInputD).toStrictEqual('false')

    // 描述输入框样式
    const descInput = await page.locator('textarea')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const descInputD = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(descInputW).toStrictEqual('346px')
    await expect(descInputH).toStrictEqual('68px')
    await expect(descInputD).toStrictEqual('true')

    await page.locator('nz-form-control').filter({ hasText: '注意：隐藏字符(空格、换行符、制表符Tab)容易导致配置出错，如果需要检测Value中隐藏字符，请点击检测隐藏字符' }).getByPlaceholder('请输入').click()
    await page.locator('nz-form-control').filter({ hasText: '注意：隐藏字符(空格、换行符、制表符Tab)容易导致配置出错，如果需要检测Value中隐藏字符，请点击检测隐藏字符' }).getByPlaceholder('请输入').fill('test  \n123  rfgasd 1')
    await page.getByRole('button', { name: '提交' }).click()
  })
  it('点击列表中某个环境变量的编辑icon，点击取消返回列表页；点击删除，该行value被清空，发布状态变为缺失', async () => {
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').last().locator('button >> nth = 0').click()
    await page.getByRole('button', { name: '取消' }).click()
    if (await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').last().locator('button >> nth = 1') && await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').last().locator('button >> nth = 1').isVisible()) {
      await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').last().locator('button >> nth = 1').click()
      await page.getByRole('button', { name: '取消' }).click()
      expect(page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 3').innerText()).not.toStrictEqual('缺失')

      await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').last().locator('button >> nth = 1').click()
      await page.getByText('该数据删除后将无法找回，请确认是否删除？').click()
      await page.getByRole('button', { name: '确定' }).click()
    }
    expect(page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 1').innerText()).toStrictEqual('')
    expect(page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 3').innerText()).toStrictEqual('缺失')
  })
  it('点击发布，检查样式；删除必填项则不可提交，填入必填项后，如表格下方有红色信息，也不可提交，点击取消返回页面', async () => {
    await page.getByRole('button', { name: '发布' }).first().click()
    // key输入框样式
    const nameInput = await page.locator('input')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(nameInputW).toStrictEqual('346px')
    await expect(nameInputH).toStrictEqual('32px')

    // 描述输入框样式
    const descInput = await page.locator('input')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(descInputW).toStrictEqual('346px')
    await expect(descInputH).toStrictEqual('68px')

    // 布局
    const formItem0 = await page.locator('.ant-drawer-body .ant-form-item-control-input-content >> nth = 0')
    const formItem0PL = await formItem0.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const formItem1 = await page.locator('.ant-drawer-body .ant-form-item-control-input-content >> nth = 1')
    const formItem1PL = await formItem1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const formItem2 = await page.locator('.ant-drawer-body .ant-form-item-control-input-content >> nth = 2')
    const formItem2PL = await formItem2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const formItem3 = await page.locator('.ant-drawer-body .ant-form-item-control-input-content >> nth = 3')
    const formItem3PL = await formItem3.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const formItem4 = await page.locator('.ant-drawer-body .ant-form-item-control-input-content >> nth = 4')
    const formItem4PL = await formItem4.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))

    // 按钮
    const submitBtn = await page.locator('.ant-drawer-body').getByRole('button', { name: '保存' })
    const submitBtnM = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const submitBtnP = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const submitBtnBC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const submitBtnBGC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const submitBtnC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(submitBtnM).toStrictEqual('0px')
    expect(submitBtnP).toStrictEqual('0px')
    expect(submitBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(submitBtnBGC).toStrictEqual('rgb(34, 84, 157)')
    expect(submitBtnC).toStrictEqual('rgb(255, 255, 255)')

    const cancleBtn = await page.locator('.ant-drawer-body').getByRole('button', { name: '取消' })
    const cancleBtnM = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const cancleBtnP = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const cancleBtnBC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const cancleBtnBGC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const cancleBtnC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(cancleBtnM).toStrictEqual('0px 0px 0px 12px')
    expect(cancleBtnP).toStrictEqual('0px 12px')
    expect(cancleBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(cancleBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(cancleBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')

    expect(formItem0PL).toStrictEqual('12px')
    expect(formItem1PL).toStrictEqual('12px')
    expect(formItem2PL).toStrictEqual('12px')
    expect(formItem3PL).toStrictEqual('12px')
    expect(formItem4PL).toStrictEqual('12px')

    await page.locator('input').click()
    await page.locator('input').fill('')
    await page.getByText('必填项').click()
    await page.getByRole('button', { name: '提交' }).click()
    await page.locator('input').click()
    await page.locator('input').fill('test')
    if (await page.locator('.ant-drawer-body .ant-form-item-explain-error') && await page.locator('.ant-drawer-body .ant-form-item-explain-error').isVisible()) {
      await page.getByRole('button', { name: '提交' }).click()
      await page.getByRole('button', { name: '取消' }).click()
    } else {
      await page.getByRole('button', { name: '提交' }).click()
    }
  })
  it('点击同步配置，检查样式和操作', async () => {
    await page.getByRole('button', { name: '同步配置' }).click()

    const nameInput = await page.locator('input.ant-input')
    const nameInputD = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const nameInputP = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const nameInputM = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))

    expect(nameInputD).toStrictEqual('true')
    expect(nameInputW).toStrictEqual('346px')
    expect(nameInputH).toStrictEqual('32px')
    expect(nameInputP).toStrictEqual('4px 11px')
    expect(nameInputM).toStrictEqual('0px 0px 0px 12px')

    const listContent1 = await page.locator('.ant-drawer-body .list-content').first()
    const listContent1P = await listContent1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const listContent1W = await listContent1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(listContent1P).toStrictEqual('0px 0px 20px 0px')
    expect(listContent1W).toStrictEqual('568px')

    const listContent2 = await page.locator('.ant-drawer-body .list-content').last()
    const listContent2P = await listContent2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const listContent2W = await listContent2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(listContent2P).toStrictEqual('0px 0px 20px 0px')
    expect(listContent2W).toStrictEqual('568px')

    // 按钮
    const submitBtn = await page.locator('.ant-drawer-body').getByRole('button', { name: '保存' })
    const submitBtnM = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const submitBtnP = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const submitBtnBC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const submitBtnBGC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const submitBtnC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(submitBtnM).toStrictEqual('0px')
    expect(submitBtnP).toStrictEqual('0px')
    expect(submitBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(submitBtnBGC).toStrictEqual('rgb(34, 84, 157)')
    expect(submitBtnC).toStrictEqual('rgb(255, 255, 255)')

    const cancleBtn = await page.locator('.ant-drawer-body').getByRole('button', { name: '取消' })
    const cancleBtnM = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const cancleBtnP = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const cancleBtnBC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const cancleBtnBGC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const cancleBtnC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(cancleBtnM).toStrictEqual('0px 0px 0px 12px')
    expect(cancleBtnP).toStrictEqual('0px 12px')
    expect(cancleBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(cancleBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(cancleBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')

    await page.getByRole('row', { name: '集群名称 所在环境' }).getByLabel('').check()
    await page.getByRole('button', { name: '提交' }).click()
    await page.getByRole('row', { name: 'KEY VALUE 更新时间' }).getByLabel('').check()
    await page.getByRole('button', { name: '取消' }).click()
    await page.getByRole('button', { name: '同步配置' }).click()
    await page.getByRole('row', { name: '集群名称 所在环境' }).getByRole('columnheader').first().click()
    await page.getByRole('row', { name: 'KEY VALUE 更新时间' }).getByLabel('').check()
    await page.getByRole('button', { name: '提交' }).click()
    await page.getByText('success').click()
  })
  it('点击更改历史，检查样式和操作', async () => {
    await page.getByRole('button', { name: '更改历史' }).click()

    const listContent = await page.locator('.ant-drawer-body .list-content').first()
    const listContentP = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const listContentW = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(listContentP).toStrictEqual('0px 0px 0px 0px')
    expect(listContentW).toStrictEqual('568px')
    await page.getByRole('button', { name: 'Close' }).click()
  })
  it('点击发布历史，检查样式和操作', async () => {
    await page.getByRole('button', { name: '发布历史' }).click()

    const listContent = await page.locator('.ant-drawer-body .list-content').first()
    const listContentP = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const listContentW = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(listContentP).toStrictEqual('0px 0px 0px 0px')
    expect(listContentW).toStrictEqual('568px')
    await page.getByRole('button', { name: 'Close' }).click()
    await page.locator('eo-ng-select-top-control').click()
    await page.getByText('20 条/页').click()
    await page.locator('.ant-drawer-body .ant-table-row-expand-icon').first().click()
    await page.locator('.ant-drawer-body .ant-table-row-expand-icon.ant-table-row-expand-icon-expanded').first().click()
    await page.getByRole('button', { name: 'Close' }).click()
  })
  it('从面包屑进入集群列表，点击表格进入环境变量列表，切换至证书管理tab；检查证书列表的样式', async () => {
    await page.locator('nz-breadcrumb').getByRole('link', { name: '网关集群' }).click()
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td').first().click()
    await page.getByRole('link', { name: '证书管理' }).click()

    const createEnvBtn = await page.getByRole('button', { name: '新建证书' })
    const createEnvBtnM = await createEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const createEnvBtnP = await createEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const createEnvBtnBC = await createEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const createEnvBtnBGC = await createEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const createEnvBtnC = await createEnvBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(createEnvBtnM).toStrictEqual('0px 0px 0px 12px')
    expect(createEnvBtnP).toStrictEqual('0px 12px')
    expect(createEnvBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(createEnvBtnBGC).toStrictEqual('rgb(34, 84, 157)')
    expect(createEnvBtnC).toStrictEqual('rgb(255, 255, 255)')

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
  it('点击新建证书，检查样式和操作，保存；再次点击新建证书，点击取消', async () => {
    await page.getByRole('button', { name: '新建证书' }).click()

    // 上传密钥按钮
    const uploadKeyBtn = await page.locator('.ant-drawer-body label.ant-btn').first()
    const uploadKeyBtnM = await uploadKeyBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const uploadKeyBtnP = await uploadKeyBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const uploadKeyBtnBC = await uploadKeyBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const uploadKeyBtnBGC = await uploadKeyBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const uploadKeyBtnC = await uploadKeyBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(uploadKeyBtnM).toStrictEqual('0px')
    expect(uploadKeyBtnP).toStrictEqual('7px 12px')
    expect(uploadKeyBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(uploadKeyBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(uploadKeyBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')

    const uploadKeyBtnIcon = await page.locator('.ant-drawer-body label.ant-btn').first().locator('iconfont')
    const uploadKeyBtnIconFS = await uploadKeyBtnIcon.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(uploadKeyBtnIconFS).toStrictEqual('14px')

    // 上传密钥输入框
    const uploadKeyInput = await page.locator('.ant-drawer-body textarea').first()
    const uploadKeyInputW = await uploadKeyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const uploadKeyInputH = await uploadKeyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const uploadKeyInputBGC = await uploadKeyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const uploadKeyInputBC = await uploadKeyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const uploadKeyInputC = await uploadKeyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const uploadKeyInputM = await uploadKeyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))

    await expect(uploadKeyInputW).toStrictEqual('346px')
    await expect(uploadKeyInputH).toStrictEqual('68px')
    await expect(uploadKeyInputBGC).toStrictEqual('rgb(51, 51, 51)')
    await expect(uploadKeyInputBC).toStrictEqual('rgb(215, 215, 215)')
    await expect(uploadKeyInputC).toStrictEqual('rgb(187, 187, 187)')
    await expect(uploadKeyInputM).toStrictEqual('12px 0px 0px 0px')

    // 上传证书按钮
    const uploadCertBtn = await page.locator('.ant-drawer-body label.ant-btn').last()
    const uploadCertBtnM = await uploadCertBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const uploadCertBtnP = await uploadCertBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const uploadCertBtnBC = await uploadCertBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const uploadCertBtnBGC = await uploadCertBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const uploadCertBtnC = await uploadCertBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(uploadCertBtnM).toStrictEqual('0px')
    expect(uploadCertBtnP).toStrictEqual('7px 12px')
    expect(uploadCertBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(uploadCertBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(uploadCertBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')

    const uploadCertBtnIcon = await page.locator('.ant-drawer-body label.ant-btn').last().locator('iconfont')
    const uploadCertBtnIconFS = await uploadCertBtnIcon.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(uploadCertBtnIconFS).toStrictEqual('14px')
    // 上传证书输入框
    const uploadCertInput = await page.locator('.ant-drawer-body textarea').first()
    const uploadCertInputW = await uploadCertInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const uploadCertInputH = await uploadCertInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const uploadCertInputBGC = await uploadCertInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const uploadCertInputBC = await uploadCertInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const uploadCertInputC = await uploadCertInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const uploadCertInputM = await uploadCertInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))

    await expect(uploadCertInputW).toStrictEqual('346px')
    await expect(uploadCertInputH).toStrictEqual('68px')
    await expect(uploadCertInputBGC).toStrictEqual('rgb(51, 51, 51)')
    await expect(uploadCertInputBC).toStrictEqual('rgb(215, 215, 215)')
    await expect(uploadCertInputC).toStrictEqual('rgb(187, 187, 187)')
    await expect(uploadCertInputM).toStrictEqual('12px 0px 0px 0px')

    // 保存、取消按钮

    const submitBtn = await page.locator('.ant-drawer-body').getByRole('button', { name: '保存' })
    const submitBtnM = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const submitBtnP = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const submitBtnBC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const submitBtnBGC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const submitBtnC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(submitBtnM).toStrictEqual('0px')
    expect(submitBtnP).toStrictEqual('0px 12px')
    expect(submitBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(submitBtnBGC).toStrictEqual('rgb(34, 84, 157)')
    expect(submitBtnC).toStrictEqual('rgb(255, 255, 255)')

    const cancleBtn = await page.locator('.ant-drawer-body').getByRole('button', { name: '取消' })
    const cancleBtnM = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const cancleBtnP = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const cancleBtnBC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const cancleBtnBGC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const cancleBtnC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(cancleBtnM).toStrictEqual('0px 0px 0px 12px')
    expect(cancleBtnP).toStrictEqual('0px 12px')
    expect(cancleBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(cancleBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(cancleBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')

    await page.getByRole('button', { name: '保存' }).click()
    await page.getByPlaceholder('密钥文件的后缀名一般为.key的文件内容').click()
    await page.getByPlaceholder('密钥文件的后缀名一般为.key的文件内容').fill('test')
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByText('上传证书').click()
    await page.setInputFiles('label:has-text("上传证书") input', './test/apinto-2027.cert')
    await page.getByPlaceholder('证书文件的后缀名一般为.crt或.pem的文件内容').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('.ant-message')
    await page.getByRole('button', { name: '取消' }).click()
  })
  it('检查面包屑，从面包屑进入集群列表，再次点击表格进入环境变量列表，测试tab，切换至网关节点tab', async () => {
    await page.locator('nz-breadcrumb').getByRole('link', { name: '网关集群' }).click()
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 1')
    await page.getByRole('link', { name: '证书管理' }).click()
    await page.getByRole('link', { name: '网关节点' }).click()
    await page.getByRole('link', { name: '配置管理' }).click()
    await page.getByRole('link', { name: '网关节点' }).click()
    await page.getByRole('link', { name: '证书管理' }).click()
    await page.getByRole('tab', { name: '环境变量' }).getByRole('link', { name: '环境变量' }).click()
    await page.getByText('环境变量证书管理网关节点配置管理').click()
    await page.getByRole('link', { name: '网关节点' }).click()
    await page.getByRole('link', { name: '配置管理' }).click()
    await page.getByRole('link', { name: '证书管理' }).click()
    await page.getByRole('tab', { name: '环境变量' }).getByRole('link', { name: '环境变量' }).click()
    await page.getByRole('link', { name: '网关节点' }).click()
    await page.getByRole('link', { name: '配置管理' }).click()
    await page.getByRole('tab', { name: '环境变量' }).getByRole('link', { name: '环境变量' }).click()
    await page.getByRole('link', { name: '配置管理' }).click()
    await page.getByRole('link', { name: '证书管理' }).click()
    await page.getByRole('link', { name: '网关节点' }).click()
  })
  it('检查网关节点列表的样式，并测试重置配置', async () => {
    const updateBtn = await page.getByRole('button', { name: '新建配置' })

    const updateBtnD = await updateBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    const updateBtnM = await updateBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const updateBtnP = await updateBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const updateBtnBC = await updateBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const updateBtnBGC = await updateBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const updateBtnC = await updateBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(updateBtnM).toStrictEqual('0px 0px 0px 12px')
    expect(updateBtnP).toStrictEqual('0px 12px')

    if (updateBtnD !== 'true') {
      expect(updateBtnBC).toStrictEqual('rgb(34, 84, 157)')
      expect(updateBtnBGC).toStrictEqual('rgb(34, 84, 157)')
      expect(updateBtnC).toStrictEqual('rgb(255, 255, 255)')
      await updateBtn.click()
    } else {
      expect(updateBtnBC).toStrictEqual('rgb(217, 217, 217)')
      expect(updateBtnBGC).toStrictEqual('rgb(245, 245, 245)')
      expect(updateBtnC).toStrictEqual('rgba(0, 0, 0, 0.25)')
    }

    const resetBtn = await page.getByRole('button', { name: '同步配置' })
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
    expect(resetBtnM).toStrictEqual('0px 0px 0px 12px')
    expect(resetBtnP).toStrictEqual('0px 12px')
    expect(resetBtnBC).toStrictEqual('rgb(217, 217, 217)')
    expect(resetBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    expect(resetBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')

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

    await page.getByRole('button', { name: '重置配置' }).click()
    await page.locator('.ant-drawer-header').getByText('重置配置').click()
    await page.locator('nz-form-control div').nth(1).click()
    const formControlDiv = page.locator('nz-form-control div').nth(1)
    const formControlDivM = await formControlDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const formControlDivP = await formControlDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(formControlDivM).toStrictEqual('0px')
    expect(formControlDivP).toStrictEqual('0px 0px 0px 12px')

    await page.getByRole('button', { name: '测试' }).click()
    await page.getByText('必填项').click()
    await page.getByPlaceholder('请输入').click()
    await page.getByPlaceholder('请输入').fill('http:11')
    await page.getByRole('button', { name: '测试' }).click()
    await page.getByText('集群地址输入错误，请重新输入').click()
    await page.getByPlaceholder('请输入').click()
    await page.getByPlaceholder('请输入').fill('http://172.17.0.1:9')
    await page.getByRole('button', { name: '测试' }).click()
    await page.getByText('集群地址需要通过测试').click()
    await page.getByPlaceholder('请输入').click()
    await page.getByPlaceholder('请输入').fill('http://172.17.0.1:9400')
    await page.getByRole('button', { name: '测试' }).click()
    await page.locator('#cdk-drop-list-1').getByRole('columnheader', { name: '名称' }).click()
    await page.getByPlaceholder('请输入').click()
    await page.getByPlaceholder('请输入').fill('http://172.17.0.1:94')
    await page.getByRole('button', { name: '测试' }).click()
    await page.getByPlaceholder('请输入').click()
    await page.getByRole('button', { name: '提交' }).click()
    await page.getByPlaceholder('请输入').click()
    await page.getByPlaceholder('请输入').fill('http://172.17.0.1:940')
    await page.getByRole('button', { name: '提交' }).click()
    await page.getByPlaceholder('请输入').click()
    await page.getByPlaceholder('请输入').fill('http://172.17.0.1:9400')
    await page.getByRole('button', { name: '提交' }).click()
    await page.locator('.ant-message').click()
    await page.getByRole('columnheader', { name: '名称' }).click()
  })
  it('点击面包屑进入集群列表，点击表格并切换至配置管理tab', async () => {
    await page.locator('nz-breadcrumb').getByRole('link', { name: '网关集群' }).click()
    await page.getByRole('cell', { name: 'zzytest' }).click()
    await page.getByRole('link', { name: '配置管理' }).click()
  })
  it('检查redis配置的样式与操作', async () => {
    const mgRedis = await page.locator('mg_layout').first()
    const mgRedisM = await mgRedis.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const mgRedisP = await mgRedis.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(mgRedisM).toStrictEqual('12px 16px 0px 16px')
    expect(mgRedisP).toStrictEqual('0px')

    const redisBox = await page.locator('deploy-cluster-redis-box').first()
    const redisBoxM = await redisBox.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const redisBoxB = await redisBox.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border'))
    let redisBoxBB = await redisBox.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-bottom'))
    expect(redisBoxM).toStrictEqual('0px')
    expect(redisBoxB).toStrictEqual('1px solid rgb(232, 232, 232)')
    expect(redisBoxBB).toStrictEqual('none')

    const redisBoxArrowE = await redisBox.locator('.title-box span.icon-zhankai').first()
    redisBoxArrowE.click()

    redisBoxBB = await redisBox.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-bottom'))
    expect(redisBoxBB).toStrictEqual('1px solid rgb(232, 232, 232)')

    const redisBoxArrowC = await redisBox.locator('.title-box span.icon-shouqi').first()
    redisBoxArrowC.click()

    redisBoxBB = await redisBox.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-bottom'))
    expect(redisBoxBB).toStrictEqual('none')

    const listTable = await page.locator('eo-ng-table')
    const listTableMT = await listTable.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))

    expect(listTableMT).toStrictEqual('0px')

    const listTableTh1 = await page.locator('eo-ng-table tr >> nth = 0 >> th >> nth = 0')
    const listTableTh1Padding = await listTableTh1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTh1Padding).toStrictEqual('0px 12px')

    const listTableTh2 = await page.locator('eo-ng-table tr >> nth = 0 >> th >> nth = 1')
    const listTableTh2Padding = await listTableTh2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTh2Padding).toStrictEqual('0px 12px')

    const listTableTd1 = await page.locator('eo-ng-table tr >> nth = 2 >> th >> nth = 0')
    const listTableTd1Padding = await listTableTd1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTd1Padding).toStrictEqual('0px')

    const listTableTd2 = await page.locator('eo-ng-table tr >> nth = 2 >> th >> nth = 1')
    const listTableTd2Padding = await listTableTd2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTd2Padding).toStrictEqual('0px')

    const listTableTd3 = await page.locator('eo-ng-table tr >> nth = 2 >> th >> nth = 2')
    const listTableTd3Padding = await listTableTd3.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTd3Padding).toStrictEqual('0px')

    const listTableTd4 = await page.locator('eo-ng-table tr >> nth = 2 >> th >> nth = 3')
    const listTableTd4Padding = await listTableTd4.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableTd4Padding).toStrictEqual('0px 12px')

    const listTableIconTh = await page.locator('eo-ng-table tr >> nth = 1 >> td >> nth = 4')
    const listTableIconThPadding = await listTableIconTh.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listTableIconThPadding).toStrictEqual('0px 24px 0px 12px')

    const listTableIcon1 = await page.locator('eo-ng-table tr >> nth = 1 >> td >> nth = 4 >> button >> nth = 0')
    const listTableIcon1PL = await listTableIcon1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    const listTableIcon1PR = await listTableIcon1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    expect(listTableIcon1PL).toStrictEqual('0px')
    expect(listTableIcon1PR).toStrictEqual('0px')

    let iconColor = await page.locator('.depoly-cluster-redis-box .eo-ng-table-btns button').first().evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    if (await page.locator('eo-ng-table eo-ng-switch') && await page.locator('eo-ng-table eo-ng-switch').isVisible()) {
      expect(iconColor).toStrictEqual('rgb(34, 84, 157)')

      await page.locator('eo-ng-switch').getByRole('button').click()
      await page.locator('eo-ng-switch').getByRole('button').click()
    } else {
      expect(iconColor).toStrictEqual('rgb(187, 187, 187)')
    }

    await page.getByPlaceholder('请输入域名/IP：端口，多个以逗号分隔').click()
    await page.getByPlaceholder('请输入域名/IP：端口，多个以逗号分隔').fill('')
    iconColor = await page.locator('.depoly-cluster-redis-box .eo-ng-table-btns button').first().evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    expect(iconColor).toStrictEqual('rgb(187, 187, 187)')

    await page.getByPlaceholder('请输入域名/IP：端口，多个以逗号分隔').click()
    await page.getByPlaceholder('请输入域名/IP：端口，多个以逗号分隔').fill('123123.a2')
    iconColor = await page.locator('.depoly-cluster-redis-box .eo-ng-table-btns button').first().evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    expect(iconColor).toStrictEqual('rgb(187, 187, 187)')

    await page.getByPlaceholder('请输入域名/IP：端口，多个以逗号分隔').click()
    await page.getByPlaceholder('请输入域名/IP：端口，多个以逗号分隔').fill('111.111.1.1:1')
    iconColor = await page.locator('.depoly-cluster-redis-box .eo-ng-table-btns button').first().evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    expect(iconColor).toStrictEqual('rgb(34, 84, 157)')

    await page.locator('.depoly-cluster-redis-box .eo-ng-table-btns button').first().click()

    await page.getByPlaceholder('请输入用户名').click()
    await page.getByPlaceholder('请输入用户名').fill('test')
    await page.getByPlaceholder('请输入密码').click()
    await page.getByPlaceholder('请输入密码').fill('aaa')
    await page.locator('.depoly-cluster-redis-box .eo-ng-table-btns button').first().click()
  })
  it('检查集群描述的样式与操作', async () => {
    const descBlock = await page.locator('.cluster-desc-block')
    const descBlockPadding = await descBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const descBlockBGC = await descBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(descBlockPadding).toStrictEqual('16px')
    expect(descBlockBGC).toStrictEqual('rgb(248, 248, 250)')

    await page.locator('.cluster-desc-block .icon-bianji').click()
    let descText = await page.locator('.cluster-desc-block span.cluster-desc-text').innerText()
    await page.locator('.cluster-desc-block input').click()
    await page.locator('.cluster-desc-block input').fill('测试内容')
    await page.locator('.cluster-desc-block .icon-guanbi').click()
    expect(await page.locator('.cluster-desc-block span.cluster-desc-text').innerText()).toStrictEqual(descText)
    await page.locator('.cluster-desc-block input').click()
    await page.locator('.cluster-desc-block input').fill('超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容超长测试内容')
    await page.locator('.cluster-desc-block .icon-chenggong').click()
    expect(await page.locator('.cluster-desc-block span.cluster-desc-text').innerText()).not.toStrictEqual(descText)

    const textH = await page.locator('.cluster-desc-block span.cluster-desc-text').evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    const sectionH = await page.locator('.cluster-desc-block section').first().evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    expect(textH).toStrictEqual('22px')
    expect(sectionH).toStrictEqual('22px')
    await page.locator('.cluster-desc-block .icon-bianji').click()

    const inputW = await page.locator('.cluster-desc-block input').first().evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(inputW).not.toStrictEqual('346px')
    await page.locator('.cluster-desc-block input').click()
    await page.locator('.cluster-desc-block input').fill('测试内容')
    await page.locator('.cluster-desc-block .icon-chenggong').click()
    descText = await page.locator('.cluster-desc-block span.cluster-desc-text').innerText()
    expect(descText).not.toStrictEqual('测试内容')
  })
})
