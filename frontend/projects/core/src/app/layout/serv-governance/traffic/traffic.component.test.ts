// 创建和编辑流量策略表单页
// -穿梭框改成虚拟
describe('服务治理-流量策略 test', () => {
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
  })

  it('保存按钮为primary，未添加数据时，点击保存按钮将出现错误提示, 添加条件按钮为primary,点击后出现弹窗', async () => {
    await page.getByRole('button', { name: '新建策略' }).click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('.ant-form-item-explain-error >> nth = 0').click()
    await page.locator('.ant-form-item-explain-error >> nth = 1').click()

    const saveBtn = await page.locator('.ant-btn-primary')
    const saveBtnBg = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))

    await expect(saveBtnBg).toStrictEqual('rgb(34, 84, 157)')

    await page.locator('button:has-text("添加条件")').click()
  })

  it('配置筛选条件弹窗中, 选择API时, 右侧穿梭框无选项时,保存按钮置灰, 页面出现穿梭框, 穿梭按钮置灰, 在左侧穿梭框勾选其中两个选项, 穿梭至右侧的按钮将变为可点击状态, 至左侧按钮仍置灰, 点击后该选项进入右侧穿梭框, 左侧穿梭框找不到该选项, 穿梭按钮置灰', async () => {
    await page.locator('eo-ng-select#name').last().click()

    await page.locator('text=API').nth(1).click()
    await page.waitForTimeout(1000)

    await expect(page.locator('nz-empty-simple svg')).not.toBeUndefined()

    const saveBtn = await page.locator('eo-ng-filter-footer .ant-btn-primary')
    let saveBtnDisabled = await saveBtn.isDisabled()
    await expect(saveBtnDisabled).toStrictEqual(true)

    const transferBtn1 = await page.locator('.ant-transfer-operation button >> nth =0')
    let transferBtn1Disabled = await transferBtn1.isDisabled()

    const transferBtn2 = await page.locator('.ant-transfer-operation button >> nth =1')
    const transferBtn2Disabled = await transferBtn2.isDisabled()

    await expect(transferBtn1Disabled).toStrictEqual(true)
    await expect(transferBtn2Disabled).toStrictEqual(true)

    await page.locator('input[type="checkbox"] >> nth=1').check()

    await page.locator('input[type="checkbox"] >> nth=2').check()

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
    await page.locator('button:has-text("添加条件")').click()
    await page.locator('eo-ng-select#name').last().click()
    await page.waitForTimeout(400)
    await page.locator('text=API路径').click()
    let tableLength = await (await page.$$('eo-ng-filter-table tr')).length
    expect(tableLength).toStrictEqual(2)

    const saveBtn = await page.locator('eo-ng-filter-footer .ant-btn-primary')
    let saveBtnDisabled = await saveBtn.isDisabled()
    await expect(saveBtnDisabled).toStrictEqual(true)

    await page.locator('[placeholder="请输入API路径"]').fill('11 11')
    saveBtnDisabled = await saveBtn.isDisabled()
    await expect(saveBtnDisabled).toStrictEqual(true)

    await page.locator('[placeholder="请输入API路径"]').fill('test')
    await page.locator('eo-ng-filter-footer .ant-btn-primary').click()
    tableLength = await (await page.$$('eo-ng-filter-table tr')).length
    expect(tableLength).toStrictEqual(3)
  })
  it('点击添加条件按钮, 属性名称中不会出现API路径, 属性名称选择API请求方式时, 页面的输入框将变为一组checkbox, 当checkbox未被选中时,保存按钮置灰, 反之为primary, 点击保存将关闭弹窗 ', async () => {
    await page.locator('button:has-text("添加条件")').click()
    await page.waitForTimeout(400)
    await page.locator('eo-ng-select#name').last().click()

    await expect(await page.locator('eo-ng-option-item:has-text("API路径")').isVisible()).toStrictEqual(false)

    const saveBtn = await page.locator('eo-ng-filter-footer .ant-btn-primary')
    const saveBtnDisabled = await saveBtn.isDisabled()
    await expect(saveBtnDisabled).toStrictEqual(true)

    await page.getByText('API请求方式').click()
    await page.getByLabel('ALL').check()
    await page.locator('eo-ng-filter-footer .ant-btn-primary').click()
    await expect(await page.locator('eo-ng-filter-table tbody tr td:has-text("所有API请求方式")').isVisible()).toStrictEqual(true)
  })
  it('点击添加条件按钮, 点击取消时, 弹窗消失, 列表行数不变', async () => {
    const tableLength = await (await page.$$('eo-ng-filter-table tr')).length

    await page.locator('button:has-text("添加条件")').click()
    await page.waitForTimeout(400)
    await page.locator('eo-ng-filter-footer .ml-btnbase').last().click()

    await expect((await page.$$('eo-ng-filter-table tr')).length).toStrictEqual(tableLength)
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
    await page.getByRole('button', { name: '添加条件' }).click()
    await page.locator('eo-ng-select#name').last().click()
    await page.getByText('IP').click()
    await page.getByPlaceholder('请输入IP地址或CIDR范围，每条以换行分割').click()
    await page.getByPlaceholder('请输入IP地址或CIDR范围，每条以换行分割').fill('111.111')
    await page.getByText('输入的IP或CIDR不符合格式').click()
    await page.getByPlaceholder('请输入IP地址或CIDR范围，每条以换行分割').click()
    await page.getByPlaceholder('请输入IP地址或CIDR范围，每条以换行分割').fill('111.111.111.111')
    await page.locator('eo-ng-filter-footer:has-text("保存 取消")').getByRole('button', { name: '保存' }).click()
  })

  it('清空所有必输项，点击保存，页面不变；逐一填写必输项，直至填写完毕，点击保存，页面出现消息提示, 当消息提示为成功时, 页面返回列表页, 否则停留本页', async () => {
    await page.getByText('*每秒限流数').click()
    await page.locator('input[name="query\\.second"]').click()
    await page.locator('input[name="query\\.second"]').fill('')
    await page.getByText('*限流维度[object Text] 请选择 必填项*每秒限流数必填项*每分限流数*每小时限流数*每秒限报文总量 MB*每分限报文总量MB*每小时限报文总量 M').click()
    await page.locator('input[name="query\\.minute"]').click()
    await page.locator('input[name="query\\.minute"]').fill('')
    await page.locator('input[name="query\\.hour"]').click()
    await page.locator('input[name="query\\.hour"]').fill('')
    await page.locator('input[name="traffic\\.second"]').click()
    await page.locator('input[name="traffic\\.second"]').fill('')
    await page.locator('eo-ng-input-group:has-text("MB")').nth(1).click()
    await page.locator('input[name="traffic\\.minute"]').click()
    await page.locator('input[name="traffic\\.minute"]').fill('')
    await page.locator('input[name="traffic\\.hour"]').click()
    await page.locator('input[name="traffic\\.hour"]').fill('')
    await page.locator('eo-ng-response-form:has-text("*HTTP状态码*Content-Typeapplication/json* CharsetUTF-8Header参数Body") input[type="number"]').click()
    await page.locator('eo-ng-response-form:has-text("*HTTP状态码*Content-Typeapplication/json* CharsetUTF-8Header参数Body") input[type="number"]').fill('')
    await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').click()
    await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').fill('testForE2e')
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByText('[object Text] 请选择').click()
    await page.locator('eo-ng-option-item:has-text("应用")').click()
    await page.locator('eo-ng-option-item:has-text("API")').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('input[name="query\\.second"]').click()
    await page.locator('input[name="query\\.second"]').fill('1')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('input[name="query\\.minute"]').click()
    await page.locator('input[name="query\\.minute"]').fill('1')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('input[name="query\\.hour"]').click()
    await page.locator('input[name="query\\.hour"]').fill('1')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('input[name="traffic\\.second"]').click()
    await page.locator('input[name="traffic\\.second"]').fill('1')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('input[name="traffic\\.minute"]').click()
    await page.locator('input[name="traffic\\.minute"]').fill('1')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('input[name="traffic\\.hour"]').click()
    await page.locator('input[name="traffic\\.hour"]').fill('1')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('eo-ng-response-form:has-text("*HTTP状态码必填项*Content-Typeapplication/json* CharsetUTF-8Header参数Body") input[type="number"]').click()
    await page.locator('eo-ng-response-form:has-text("*HTTP状态码必填项*Content-Typeapplication/json* CharsetUTF-8Header参数Body") input[type="number"]').fill('100')
    await page.getByText('application/json').click()

    await page.getByText('*Content-Type').click()
    await page.getByText('UTF-8').click()
    await page.getByText('GBK').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('div:has-text("success")').nth(3).click()

    // const tableLength = await (await page.$$('eo-ng-apinto-table tr')).length
    // await page.locator('eo-ng-apinto-table tr >> nth = ' + (tableLength - 1) + '>> td >> nth = 8 >> button >> nth = 1').click()
    // await page.getByRole('button', { name: '确定' }).click()
  })

  it('点击列表中的查看按钮, 策略名称&筛选条件列表&限流规则不为空, 保存按钮为primary, 点击保存出现全局消息提示, 并返回列表页', async () => {
    await page.getByRole('cell', { name: 'testForE2e' }).last().click()
    await page.getByText('流量策略 / testForE2e /').click()
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

    const metricInput = await page.locator('eo-ng-select#metrics')
    const metricInputW = await metricInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const metricInputnH = await metricInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(metricInputW).toStrictEqual('346px')
    await expect(metricInputnH).toStrictEqual('32px')

    const querySecondInput = await page.locator('[id="query\\.second"]')
    const querySecondInputW = await querySecondInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const querySecondInputH = await querySecondInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(querySecondInputW).toStrictEqual('346px')
    await expect(querySecondInputH).toStrictEqual('32px')

    const queryMinuteInput = await page.locator('[id="query\\.minute"]')
    const queryMinuteInputW = await queryMinuteInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const queryMinuteInputH = await queryMinuteInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(queryMinuteInputW).toStrictEqual('346px')
    await expect(queryMinuteInputH).toStrictEqual('32px')

    const queryHourInput = await page.locator('[id="query\\.hour"]')
    const queryHourInputW = await queryHourInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const queryHourInputH = await queryHourInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(queryHourInputW).toStrictEqual('346px')
    await expect(queryHourInputH).toStrictEqual('32px')

    const trafficSecondInput = await page.locator('eo-ng-input-group:has-text("MB") >> nth=0')
    const trafficSecondInputW = await trafficSecondInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const trafficSecondInputH = await trafficSecondInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(trafficSecondInputW).toStrictEqual('346px')
    await expect(trafficSecondInputH).toStrictEqual('32px')

    const trafficMinuteInput = await page.locator('eo-ng-input-group:has-text("MB") >> nth=1')
    const trafficMinuteInputW = await trafficMinuteInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const trafficMinuteInputH = await trafficMinuteInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(trafficMinuteInputW).toStrictEqual('346px')
    await expect(trafficMinuteInputH).toStrictEqual('32px')

    const trafficHourInput = await page.locator('eo-ng-input-group:has-text("MB") >> nth=2')
    const trafficHourInputW = await trafficHourInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const trafficHourInputH = await trafficHourInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(trafficHourInputW).toStrictEqual('346px')
    await expect(trafficHourInputH).toStrictEqual('32px')

    const httpStatusInput = await page.locator('eo-ng-response-form').getByRole('spinbutton')
    const httpStatusInputW = await httpStatusInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const httpStatusInputH = await httpStatusInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(httpStatusInputW).toStrictEqual('346px')
    await expect(httpStatusInputH).toStrictEqual('32px')

    const contentTypeInput = await page.locator('eo-ng-select-top-control').filter({ hasText: 'application/json' })
    const contentTypeInputW = await contentTypeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const contentTypeInputH = await contentTypeInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(contentTypeInputW).toStrictEqual('346px')
    await expect(contentTypeInputH).toStrictEqual('32px')

    const charsetInput = await page.locator('eo-ng-select-top-control').filter({ hasText: 'UTF-8' })
    const charsetInputW = await charsetInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const charsetInputH = await charsetInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(charsetInputW).toStrictEqual('346px')
    await expect(charsetInputH).toStrictEqual('32px')

    const headerFirstInput1 = await page.getByRole('cell', { name: '请输入Key' })
    const headerFirstInput1W = await headerFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const headerFirstInput1PR = await headerFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    const headerFirstInput1H = await headerFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(headerFirstInput1W).toStrictEqual('182px')
    await expect(headerFirstInput1PR).toStrictEqual('8px')
    await expect(headerFirstInput1H).toStrictEqual('32px')

    const headerFirstInput2 = await page.getByRole('cell', { name: '请输入Value' })
    const headerFirstInput2W = await headerFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const headerFirstInput2PR = await headerFirstInput1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-right'))
    const headerFirstInput2H = await headerFirstInput2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(headerFirstInput2W).toStrictEqual('172px')
    await expect(headerFirstInput2PR).toStrictEqual('8px')
    await expect(headerFirstInput2H).toStrictEqual('32px')

    const headerFirstBtn = await page.getByRole('row', { name: '请输入Key 请输入Value ' }).getByRole('button', { name: '' })
    const headerFirstBtnH = await headerFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const headerFirstBtnLH = await headerFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const headerFirstBtnColor = await headerFirstBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(headerFirstBtnH).toStrictEqual('32px')
    await expect(headerFirstBtnLH).toStrictEqual('32px')
    await expect(headerFirstBtnColor).toStrictEqual('rgb(34, 84, 157)')

    await page.getByRole('button', { name: '' }).click()

    const headerFirstInput1A = await page.locator('eo-ng-apinto-table.arrayItem tr >> nth = 1 >> td >> nth = 0')
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
    await expect(headerFirstInput1APB).toStrictEqual('12px')
    await expect(headerFirstInput1AH).toStrictEqual('44px')

    const headerSecondBtn = await page.locator('eo-ng-apinto-table.arrayItem tr >> nth = 2 >> td >> nth = 2 >> button >> nth = 0')
    const headerSecondBtnH = await headerSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const headerSecondBtnLH = await headerSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('line-height'))
    const headerSecondBtnColor = await headerSecondBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(headerSecondBtnH).toStrictEqual('32px')
    await expect(headerSecondBtnLH).toStrictEqual('32px')
    await expect(headerSecondBtnColor).toStrictEqual('rgb(34, 84, 157)')

    const bodyInput = await page.locator('eo-ng-response-form nz-form-item').filter({ hasText: 'Body' }).getByPlaceholder('请输入')
    const bodyInputW = await bodyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const bodyInputH = await bodyInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(bodyInputW).toStrictEqual('346px')
    await expect(bodyInputH).toStrictEqual('32px')

    const limitSection1 = await page.locator('section.limit-bg >> nth = 0')
    const limitSection1BG = await limitSection1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const limitSection1P = await limitSection1.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))

    await expect(limitSection1BG).toStrictEqual('rgb(248, 248, 250)')
    await expect(limitSection1P).toStrictEqual('20px')

    const limitSection2 = await page.locator('section.limit-bg >> nth = 1')
    const limitSection2BG = await limitSection2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const limitSection2P = await limitSection2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))

    await expect(limitSection2BG).toStrictEqual('rgb(248, 248, 250)')
    await expect(limitSection2P).toStrictEqual('20px')
  })
})
