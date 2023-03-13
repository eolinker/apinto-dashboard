// 创建和编辑流量策略表单页
// -穿梭框改成虚拟
describe('服务治理-灰度策略 e2e test', () => {
  it('初始化页面, 默认第一个环境列表中的第一个集群被选中,列表页显示数据, 点击新建按钮, 进入新建策略页', async () => {
    // Go to http://localhost:4200/serv-goverance/traffic
    await page.goto('http://localhost:4200/login')
    await page.waitForTimeout(2000)

    await page.locator('[placeholder="请输入账号"]').fill('maggie')
    // Click [placeholder="请输入密码"]
    await page.locator('[placeholder="请输入密码"]').click()
    // Fill [placeholder="请输入密码"]
    await page.locator('[placeholder="请输入密码"]').fill('12345678')
    await page.locator('button:has-text("登录")').click()
    await page.locator('#cdk-overlay-0 div:has-text("登录成功")').nth(3).click()

    await page.locator('eo-ng-menu-default div:has-text("服务治理")').click()
    await page.getByRole('link', { name: '灰度策略' }).click()
  })

  it('保存按钮为primary，未添加数据时，点击保存按钮将出现错误提示 ', async () => {
    await page.getByRole('button', { name: '新建策略' }).click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('.ant-form-item-explain-error >> nth = 0').click()
  })

  it('添加条件按钮为primary,点击后出现弹窗，配置筛选条件弹窗中, 选择API时, 右侧穿梭框无选项时,保存按钮置灰, 页面出现穿梭框, 穿梭按钮置灰, 在左侧穿梭框勾选其中两个选项, 穿梭至右侧的按钮将变为可点击状态, 至左侧按钮仍置灰, 点击后该选项进入右侧穿梭框, 左侧穿梭框找不到该选项, 穿梭按钮置灰', async () => {
    await page.locator('button:has-text("添加条件")').first().click()
    await page.locator('eo-ng-select#name').last().click()

    await page.locator('text=API').nth(1).click()
    await page.waitForTimeout(200)

    await expect(page.locator('nz-empty-simple svg')).not.toBeUndefined()

    const saveBtn = await page.locator('eo-ng-filter-footer >> nth = 0 >> .ant-btn-primary')
    let saveBtnDisabled = await saveBtn.isDisabled()
    await expect(saveBtnDisabled).toStrictEqual(true)

    const transferBtn1 = await page.locator('.ant-transfer-operation button >> nth =0')
    let transferBtn1Disabled = await transferBtn1.isDisabled()

    const transferBtn2 = await page.locator('.ant-transfer-operation button >> nth =1')
    const transferBtn2Disabled = await transferBtn2.isDisabled()

    await expect(transferBtn1Disabled).toStrictEqual(true)
    await expect(transferBtn2Disabled).toStrictEqual(true)

    await page.locator('input[type="checkbox"] >> nth=1').click()
    await page.locator('input[type="checkbox"] >> nth=2').click()
    transferBtn1Disabled = await transferBtn1.isDisabled()

    await expect(transferBtn1Disabled).toStrictEqual(true)
    await page.locator('.ant-transfer-operation button >> nth =1').click()

    saveBtnDisabled = await saveBtn.isDisabled()
    await expect(saveBtnDisabled).toStrictEqual(false)

    await page.locator('eo-ng-filter-footer .ant-btn-primary').click()
  })

  // it('在搜索框内输入搜索内容, 页面选项将发生变化, 点击保存将关闭弹窗, 添加条件列表中将新增一行', async () => {
  //
  //   await page.locator('button:has-text("添加条件")').click()
  //   await page.locator('eo-ng-select#name').click()

  //   await page.locator('text=应用').nth(1).click()
  //   await page.waitForTimeout(100)

  //   await page.locator(' [placeholder="请输入搜索内容"] >> nth = 1').fill('test')
  //   await expect(await page.locator('text=匿名应用').isVisible()).toStrictEqual(false)
  // })

  it('属性名称选择API路径时, 输入框未有值时, 保存按钮置灰, 输入值不符合正则时，保存按钮置灰，反之为primary, 点击保存将关闭弹窗, 添加条件列表中将新增一行', async () => {
    await page.getByLabel('按规则').click()
    await page.locator('button:has-text("添加条件")').last().click()
    await page.locator('eo-ng-select#name').last().click()
    await page.waitForTimeout(400)
    await page.locator('text=API路径').click()
    const tableLength = await (await page.$$('eo-ng-filter-table  tr')).length

    const saveBtn = await page.locator('eo-ng-filter-footer  .ant-btn-primary')
    let saveBtnDisabled = await saveBtn.isDisabled()
    await expect(saveBtnDisabled).toStrictEqual(true)

    await page.locator('[placeholder="请输入API路径"]').fill('11 11')
    saveBtnDisabled = await saveBtn.isDisabled()
    await expect(saveBtnDisabled).toStrictEqual(true)

    await page.locator('[placeholder="请输入API路径"]').fill('test')
    await page.locator('eo-ng-filter-footer .ant-btn-primary').click()
    expect(await (await page.$$('eo-ng-filter-table  tr')).length).toStrictEqual(tableLength + 1)
  })
  it('点击添加条件按钮, 属性名称中不会出现API路径, 属性名称选择API请求方式时, 页面的输入框将变为一组checkbox, 当checkbox未被选中时,保存按钮置灰, 反之为primary, 点击保存将关闭弹窗 ', async () => {
    await page.locator('button:has-text("添加条件")').first().click()
    await page.waitForTimeout(400)
    await page.locator('eo-ng-select#name').last().click()

    await expect(await page.locator('eo-ng-option-item:has-text("API路径")').isVisible()).toStrictEqual(false)

    const saveBtn = await page.locator('eo-ng-filter-footer >> nth = 0 >> .ant-btn-primary')
    const saveBtnDisabled = await saveBtn.isDisabled()
    await expect(saveBtnDisabled).toStrictEqual(true)

    await page.getByText('API请求方式').click()
    await page.getByLabel('ALL').check()
    await page.locator('eo-ng-filter-footer  >> nth = 0 >>.ant-btn-primary').click()
    await expect(await page.locator('eo-ng-filter-table >> nth = 0 >> tbody tr td:has-text("所有API请求方式")').isVisible()).toStrictEqual(true)
  })
  it('点击添加条件按钮, 之前选择过的属性不会出现在列表中 点击取消时, 弹窗消失, 列表行数不变', async () => {
    const tableLength = await (await page.$$('eo-ng-filter-table >> nth = 1 >> tr')).length

    await page.locator('button:has-text("添加条件")').last().click()
    await page.waitForTimeout(400)
    await page.locator('eo-ng-select#name').last().click()
    await expect(await page.locator('eo-ng-option-item:has-text("API路径")').isVisible()).toStrictEqual(false)
    await expect(await page.locator('eo-ng-option-item:has-text("API请求方式")').isVisible()).toStrictEqual(false)
    await expect(await page.locator('eo-ng-option-item:has-text("API")').isVisible()).toStrictEqual(false)
    await page.locator('eo-ng-filter-footer .ml-btnbase').last().click()

    await expect((await page.$$('eo-ng-filter-table  >> nth = 1 >>  tr')).length).toStrictEqual(tableLength)
  })

  // it('点击列表第一列(API)中的配置, 弹窗中属性名称与列表中的属性名相同, 且穿梭框右侧有值, 全选左侧穿梭框并点击穿梭按钮, 提交后, 该列的属性值将变为全部API, 点击删除, 出现删除弹窗, 点击取消, 弹窗消失, 再次点击删除, 出现删除弹窗, 点击确认后, 弹窗关闭, 列表中该列将消失', async () => {
  //   const tableLength = await (await page.$$('eo-ng-filter-table tr')).length
  //   await page.getByRole('cell', { name: 'API' }).click()
  //   // await page.locator('eo-ng-filter-table tbody tr >> nth = 0  >> td >> nth =2 >> button >> nth = 0').first().click()
  //   await page.waitForTimeout(2000)

  //   await page.locator('input[type="checkbox"] >> nth=0').click()
  //   await page.locator('.ant-transfer-operation button >> nth =1').click()
  //   await page.locator('eo-ng-filter-footer .ant-btn-primary').click()

  //   await page.locator('.icon-shanchu >> nth = 2').click()
  //   await page.locator('.nz-modal-footer button:has-text("取消")').click()
  //   await expect((await page.$$('eo-ng-filter-table tr')).length).toStrictEqual(tableLength)

  //   await page.locator('.icon-shanchu >> nth = 1').click()
  //   await page.locator('.nz-modal-footer button:has-text("确定")').click()
  //   await expect((await page.$$('eo-ng-filter-table tr')).length).toStrictEqual(tableLength - 1)
  // })

  // it('点击列表第二列(API请求方式)的配置, 弹窗中的checkbox中有部分被勾选, 点击ALL选项则全部checkbox被勾选, 提交后, 该列的属性值将变为全部请求方式', async () => {
  //   await page.locator('.icon-a-peizhianniu_huaban1 >> nth = 1').click()
  //   await page.getByLabel('DELETE').uncheck()
  //   await page.getByLabel('ALL').check()
  //   // Click button:has-text("保存")
  //   await page.locator('eo-ng-filter-footer .ant-btn-primary:has-text("提交")').click()
  // })

  // it('点击列表第一列(API路径)的配置, 弹窗中的输入框内值与列表中的属性值相同, 编辑输入框, 点击取消, 弹窗消失后列表内容不变, 再次点击配置, 改变输入框内容后提交, 弹窗消失后列表内容发生相应变化', async () => {
  //   await page.getByRole('cell', { name: 'API路径' }).click()
  //   await page.locator('[placeholder="请输入API路径"]').click()
  //   await page.locator('eo-ng-filter-footer button:has-text("取消")').click()

  //   await page.getByRole('cell', { name: 'API路径' }).click()
  //   await page.locator('[placeholder="请输入API路径"]').click()
  //   await page.locator('[placeholder="请输入API路径"]').fill('test1')
  //   await page.locator('eo-ng-filter-footer button:has-text("提交")').click()
  //   await page.getByRole('cell', { name: 'test1' }).click()
  //   await page.locator('eo-ng-filter-footer button:has-text("取消")').click()
  // })

  it('点击添加条件，选择ip地址，当输入的不符合正则时，出现错误提示且无法提交', async () => {
    await page.getByRole('button', { name: '添加条件' }).first().click()
    await page.locator('eo-ng-select#name').last().click()
    await page.getByText('IP').last().click()
    await page.getByPlaceholder('请输入IP地址或CIDR范围，每条以换行分割').click()
    await page.getByPlaceholder('请输入IP地址或CIDR范围，每条以换行分割').fill('111.111')
    await page.getByText('输入的IP或CIDR不符合格式').click()
    await page.getByPlaceholder('请输入IP地址或CIDR范围，每条以换行分割').click()
    await page.getByPlaceholder('请输入IP地址或CIDR范围，每条以换行分割').nth(0).fill('111.111.111.111')
    await page.locator('eo-ng-filter-footer:has-text("保存 取消")').getByRole('button', { name: '保存' }).click()
  })

  it('填写所有必填项，点击保存，页面出现消息提示, 当消息提示为成功时, 页面返回列表页, 否则停留本页', async () => {
    await page.locator('.fix-buttom-group').getByRole('button', { name: '保存' }).click()
    await page.getByText('必填项').first().click()
    await page.locator('section').getByRole('alert').click()
    await page.getByPlaceholder('请输入首字母为英文，英文数字下划线任意一种组合').fill('testForE2e')

    await page.locator('#desc').click()
    await page.locator('#desc').fill('desc')
    await page.getByRole('button', { name: '否' }).click()
    await page.locator('.fix-buttom-group').getByRole('button', { name: '保存' }).click()
    await page.getByPlaceholder('请输入主机名或IP：端口').click()
    await page.getByPlaceholder('请输入主机名或IP：端口').fill('111')
    await page.getByRole('cell', { name: '请输入主机名或IP：端口 请输入主机名或IP：端口' }).getByText('请输入主机名或IP：端口').click()
    await page.locator('.fix-buttom-group').getByRole('button', { name: '保存' }).click()

    await page.getByPlaceholder('请输入主机名或IP：端口').nth(0).fill('111.111.111.11:11')

    await page.getByPlaceholder('请输入主机名或IP：端口').nth(1).click()
    await page.getByPlaceholder('请输入主机名或IP：端口').nth(1).fill('222')
    await page.getByRole('cell', { name: '请输入主机名或IP：端口 请输入主机名或IP：端口' }).getByText('请输入主机名或IP：端口').click()
    await page.getByPlaceholder('请输入主机名或IP：端口').nth(2).click()
    await page.locator('.fix-buttom-group').getByRole('button', { name: '保存' }).click()
    await page.getByPlaceholder('请输入主机名或IP：端口').nth(1).fill('222.222.222.22:22')
    await page.getByLabel('按百分比').click()
    await page.locator('section nz-form-item').filter({ hasText: '*灰度节点流量占比%' }).locator('nz-slider div').nth(1).click()
    await page.getByRole('spinbutton').nth(1).click()
    await page.getByRole('spinbutton').nth(1).fill('70')
    await page.getByRole('spinbutton').nth(1).press('Enter')
    await page.getByRole('spinbutton').first().click()
    await page.locator('.fix-buttom-group').getByRole('button', { name: '保存' }).click()
    await page.getByText('success').click()
  })

  it('点击列表中的查看按钮, 策略名称&筛选条件列表&限流规则不为空, 保存按钮为primary, 点击保存出现全局消息提示, 并返回列表页', async () => {
    await page.getByRole('cell', { name: 'testForE2e' }).last().click()
    await page.getByText('灰度策略 / testForE2e /').click()
    await page.getByRole('button', { name: '提交' }).click()
  })

  it('输入框大小, 按钮的大小和颜色, 表格的大小, 限流规则的背景色', async () => {
    await page.getByRole('button', { name: '新建策略' }).click()

    const nameInput = await page.locator('input#name')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(nameInputW).toStrictEqual('346px')
    await expect(nameInputH).toStrictEqual('32px')

    const descInput = await page.locator('textarea#desc')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(descInputW).toStrictEqual('346px')
    await expect(descInputH).toStrictEqual('68px')

    const priorityInput = await page.locator('input#priority')
    const priorityInputW = await priorityInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const priorityInputH = await priorityInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(priorityInputW).toStrictEqual('346px')
    await expect(priorityInputH).toStrictEqual('32px')

    const greyNodesInput = await page.getByRole('row', { name: '请输入主机名或IP：端口 ' }).getByPlaceholder('请输入主机名或IP：端口')
    const greyNodesInputW = await greyNodesInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const greyNodesInputH = await greyNodesInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(greyNodesInputW).toStrictEqual('346px')
    await expect(greyNodesInputH).toStrictEqual('32px')

    const greyNodes1btn = await page.getByRole('row', { name: '请输入主机名或IP：端口 ' }).getByRole('button', { name: '' })
    const greyNodes1btnW = await greyNodes1btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const greyNodes1btnH = await greyNodes1btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(greyNodes1btnW).toStrictEqual('20px')
    await expect(greyNodes1btnH).toStrictEqual('32px')

    await page.getByRole('row', { name: '请输入主机名或IP：端口 ' }).getByRole('button', { name: '' }).click()

    const greyNodes2btn = await page.locator('.arrayItem tr >> nth = 2 >> button >> nth = 0')
    const greyNodes2btnW = await greyNodes2btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const greyNodes2btnH = await greyNodes2btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(greyNodes2btnW).toStrictEqual('28px')
    await expect(greyNodes2btnH).toStrictEqual('32px')

    const filterBtn1 = await page.locator('nz-form-item').filter({ hasText: '筛选流量 添加条件' }).getByRole('button', { name: '添加条件' }).first()
    const filterBtn1W = await filterBtn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const filterBtn1H = await filterBtn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const filterBtn1BGC = await filterBtn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const filterBtn1BC = await filterBtn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const filterBtn1C = await filterBtn1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(filterBtn1W).toStrictEqual('82px')
    await expect(filterBtn1H).toStrictEqual('32px')
    await expect(filterBtn1BGC).toStrictEqual('rgb(255, 255, 255)')
    await expect(filterBtn1BC).toStrictEqual('rgb(217, 217, 217)')
    await expect(filterBtn1C).toStrictEqual('rgba(0, 0, 0, 0.85)')
    // filterBtn1.hover()
    // filterBtn1BC = await filterBtn1.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('border-color'))
    // await expect(filterBtn1BC).toStrictEqual('rgb(34, 84, 157)')
    // filterBtn1C = await filterBtn1.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('color'))
    // await expect(filterBtn1C).toStrictEqual('rgb(34, 84, 157)')

    const label1 = await page.getByText('按百分比')
    const label2 = await page.getByText('按规则')
    const label1C = await label1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const label2C = await label2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(label1C).toStrictEqual('rgb(0, 0, 0)')
    await expect(label2C).toStrictEqual('rgba(0, 0, 0, 0.85)')

    const sider1 = await page.locator('section nz-form-item').filter({ hasText: '*灰度节点流量占比%' }).locator('nz-slider div').nth(1)
    const sider1Input = await page.locator('eo-ng-input-group').first()
    const sider2 = await page.locator('section nz-form-item').filter({ hasText: '*上游节点流量占比%' }).locator('nz-slider div').nth(1)
    const sider2Input = await page.locator('eo-ng-input-group').last()

    const sider1W = await sider1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const sider1H = await sider1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const sider1Radius = await sider1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-radius'))
    const sider1BG = await sider1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))

    await expect(sider1W).toStrictEqual('318px')
    await expect(sider1H).toStrictEqual('10px')
    await expect(sider1Radius).toStrictEqual('16px')
    await expect(sider1BG).toStrictEqual('rgb(217, 217, 217)')

    const sider1InputW = await sider1Input.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const sider1InputH = await sider1Input.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const sider1InputBC = await sider1Input.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))

    await expect(sider1InputW).toStrictEqual('80px')
    await expect(sider1InputH).toStrictEqual('32px')
    await expect(sider1InputBC).toStrictEqual('rgb(217, 217, 217)')

    const sider2W = await sider2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const sider2H = await sider2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const sider2Radius = await sider2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-radius'))
    const sider2BG = await sider2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))

    await expect(sider2W).toStrictEqual('318px')
    await expect(sider2H).toStrictEqual('10px')
    await expect(sider2Radius).toStrictEqual('16px')
    await expect(sider2BG).toStrictEqual('rgb(217, 217, 217)')

    const sider2InputW = await sider2Input.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const sider2InputH = await sider2Input.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const sider2InputBC = await sider2Input.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))

    await expect(sider2InputW).toStrictEqual('80px')
    await expect(sider2InputH).toStrictEqual('32px')
    await expect(sider2InputBC).toStrictEqual('rgb(217, 217, 217)')

    await page.getByLabel('按规则').click()
    const filterBtn2 = await page.getByRole('button', { name: '添加配置' })
    const filterBtn2W = await filterBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const filterBtn2H = await filterBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const filterBtn2BGC = await filterBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const filterBtn2BC = await filterBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))

    await expect(filterBtn2W).toStrictEqual('82px')
    await expect(filterBtn2H).toStrictEqual('32px')
    await expect(filterBtn2BGC).toStrictEqual('rgb(255, 255, 255)')
    await expect(filterBtn2BC).toStrictEqual('rgb(217, 217, 217)')
    // filterBtn2.hover()
    // filterBtn2BC = await filterBtn2.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('border-color'))
    // await expect(filterBtn2BC).toStrictEqual('rgb(34, 84, 157)')
    // filterBtn2C = await filterBtn2.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('color'))
    // await expect(filterBtn2C).toStrictEqual('rgb(34, 84, 157)')

    const continueSwitch = await page.locator('eo-ng-switch >> button')
    const continueSwitchW = await continueSwitch.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const continueSwitchH = await continueSwitch.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(continueSwitchW).toStrictEqual('35px')
    await expect(continueSwitchH).toStrictEqual('16px')
  })
})
