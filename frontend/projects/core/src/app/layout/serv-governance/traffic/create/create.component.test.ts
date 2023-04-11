// 创建和编辑流量策略表单页
// -穿梭框改成虚拟
describe('traffic-create e2e test', () => {
  it('初始化页面, 默认第一个环境列表中的第一个集群被选中,列表页显示数据, 点击新建按钮, 进入新建策略页', async () => {
    // Go to http://localhost:4200/serv-goverance/traffic
    await page.goto('http://localhost:4200/serv-goverance/traffic')
    await page.waitForTimeout(2000)
  })

  it('提交按钮置灰,不可提交, 添加条件按钮为primary,点击后出现弹窗', async () => {
    await page.locator('button:has-text("新建策略")').click()
    await expect(page.locator('button:has-text("添加条件")').isDisabled).toStrictEqual(false)
    await expect(page.locator('button:has-text("保存")').isDisabled).toStrictEqual(true)
    await expect(page.locator('text=*属性名称')).toBeUndefined()
    await page.locator('button:has-text("添加条件")').click()
    await expect(page.locator('text=*属性名称')).not.toBeUndefined()
  })

  it('配置筛选条件弹窗中, 选择API时, 右侧穿梭框无选项时,保存按钮置灰, 页面出现穿梭框, 穿梭按钮置灰, 在左侧穿梭框勾选其中两个选项, 穿梭至右侧的按钮将变为可点击状态, 至左侧按钮仍置灰, 点击后该选项进入右侧穿梭框, 左侧穿梭框找不到该选项, 穿梭按钮置灰', async () => {
    await page.locator('text=*属性名称应用 >> svg').click()

    await page.locator('text=API').nth(1).click()
    await expect(page.locator('nz-empty-simple svg')).not.toBeUndefined()
    await expect(page.locator('button:has-text("保存")').isDisabled).toStrictEqual(true)
    await expect(page.locator('.ant-transfer-operation button >> nth =1').isDisabled).toStrictEqual(true)
    await expect(page.locator('.ant-transfer-operation button >> nth =0').isDisabled).toStrictEqual(true)

    await page.locator('input[type="checkbox"] >> nth=1').check()
    await page.locator('input[type="checkbox"] >> nth=2').check()
    await expect(page.locator('.ant-transfer-operation button >> nth =0').isDisabled).toStrictEqual(true)
    await page.locator('.ant-transfer-operation button >> nth =1').click()
    await expect(page.locator('button:has-text("保存")').isDisabled).toStrictEqual(false)

    await page.locator('eo-ng-transfer-list >> nth = 1 >> input[type=checkbox] >> nth = 1').check()
    await page.locator('.ant-transfer-operation button >> nth =0').click()
  })

  it('在穿梭框有目录筛选选择器的情况下, 点击选择器, 页面将出现全局消息提示', async () => {
    await expect(page.locator('text=筛选成功!')).toBeUndefined()
    await page.locator('eo-ng-cascader div:has-text("目录筛选")').nth(1).click()
    await page.locator('li:has-text("鉴权管理")').click()
    await expect(page.locator('text=筛选成功!')).not.toBeUndefined()
  })

  it('在搜索框内输入搜索内容, 页面选项将发生变化, 点击保存将关闭弹窗, 添加条件列表中将新增一行', async () => {
    await page.locator('nz-select-item:has-text("API")').click()
    await page.locator('text=应用').nth(1).click()
    await expect(page.locator('text=application_beta!')).not.toBeUndefined()

    await page.locator(' [placeholder="请输入搜索内容"]').fill('test')
    await expect(page.locator('text=application_beta!')).toBeUndefined()
  })

  it('属性名称选择API路径时, 输入框未有选中值时, 保存按钮置灰, 反之为primary, 点击保存将关闭弹窗, 添加条件列表中将新增一行', async () => {
    await page.locator('nz-select-item:has-text("应用")').click()
    await page.locator('text=API路径').click()
    await expect(page.locator('button:has-text("保存")').isDisabled).toStrictEqual(true)
    await page.locator('[placeholder="请输入API路径"]').fill('test')
    await page.locator('button:has-text("保存")').click()
    await expect(page.locator('text=API路径')).not.toBeUndefined()
  })
  it('点击添加条件按钮, 属性名称中不会出现API路径, 属性名称选择API请求方式时, 页面的输入框将变为一组checkbox, 当checkbox未被选中时,保存按钮置灰, 反之为primary, 点击保存将关闭弹窗 ', async () => {
    await page.locator('button:has-text("添加条件")').click()
    await page.locator('text=*属性名称应用 >> svg').click()
    await expect(page.locator('eo-ng-option-item:has-text("API路径")')).not.toBeUndefined()
    await expect(page.locator('.ant-checkbox-input')).not.toBeUndefined()
    await expect(page.locator('button:has-text("保存")').isDisabled).toStrictEqual(true)

    await page.locator('td:has-text("所有API请求方式")').click()

    await page.locator('.ant-checkbox-input').first().check()
  })
  it('点击添加条件按钮, 点击取消或关闭时, 弹窗消失, 列表行数不变', async () => {
    await expect((await page.$$('eo-ng-table tr')).length).toStrictEqual(4)

    await page.locator('button:has-text("添加条件")').click()
    await page.locator('[aria-label="Close"]').click()
    await page.locator('button:has-text("添加条件")').click()
    await page.locator('#cdk-overlay-9 button:has-text("取消")').click()

    await expect((await page.$$('eo-ng-table tr')).length).toStrictEqual(4)
  })

  it('点击列表第一列(API)中的配置, 弹窗中属性名称与列表中的属性名相同, 且穿梭框右侧有值, 全选左侧穿梭框并点击穿梭按钮, 提交后, 该列的属性值将变为全部API, 点击删除, 出现删除弹窗, 点击取消, 弹窗消失, 再次点击删除, 出现删除弹窗, 点击确认后, 弹窗关闭, 列表中该列将消失', async () => {
    await page.locator('text=配置').first().click()
    await page.locator('nz-select-item:has-text("API")').click()
    await page.locator('text=*属性名称API >> section').click()
    await page.locator('input[type="checkbox"] >> nth=0').check()
    await page.locator('.ant-transfer-operation button >> nth =1').click()
    await page.locator('button:has-text("保存")').click()

    await page.locator('text=所有API').first().click()
    await page.locator('text=删除').first().click()
    await page.locator('div[role="document"] button:has-text("取消")').click()

    await page.locator('text=删除').first().click()
    await page.locator('button:has-text("确定")').click()
    await expect((await page.$$('eo-ng-table tr')).length).toStrictEqual(3)
  })

  it('点击列表第二列(API请求方式)的配置, 弹窗中的checkbox中有部分被勾选, 点击ALL选项则全部checkbox被勾选, 提交后, 该列的属性值将变为全部请求方式', async () => {
    // Click text=配置 >> nth=0
    await page.locator('text=配置').first().click()
    // Uncheck label:nth-child(5) > .ant-checkbox > .ant-checkbox-input
    await page.locator('label:nth-child(5) > .ant-checkbox > .ant-checkbox-input').uncheck()
    // Check .ant-checkbox-input >> nth=0
    await page.locator('.ant-checkbox-input').first().check()
    // Click button:has-text("保存")
    await page.locator('button:has-text("保存")').click()
  })

  it('点击列表第一列(API路径)的配置, 弹窗中的输入框内值与列表中的属性值相同, 编辑输入框, 点击取消, 弹窗消失后列表内容不变, 再次点击配置, 改变输入框内容后提交, 弹窗消失后列表内容发生相应变化', async () => {
    await page.locator('text=配置').first().click()
    await page.locator('[placeholder="请输入API路径"]').click()
    await page.locator('#cdk-overlay-11 button:has-text("取消")').click()
    await page.locator('text=配置').first().click()
    await page.locator('[placeholder="请输入API路径"]').click()
    await page.locator('[placeholder="请输入API路径"]').fill('test1')
    await page.locator('button:has-text("保存")').click()
    await expect(page.locator('test=test1")')).not.toBeUndefined()
  })

  it('选择限流维度, 此时提交按钮将为primary, 点击提交, 页面出现消息提示, 当消息提示为成功时, 页面返回列表页, 否则停留本页', async () => {
    // Fill text=*策略名称英文数字下划线任意一种,首字母必须为英文 >> [placeholder="请输入"]
    await page.locator('text=*策略名称英文数字下划线任意一种,首字母必须为英文 >> [placeholder="请输入"]').fill('tet')
    // Click eo-ng-traffic-create:has-text("*策略名称英文数字下划线任意一种,首字母必须为英文描述优先级*筛选条件添加条件属性名称属性值操作API请求方式所有API请求方式 配置 删除 API路径test")
    await page.locator('eo-ng-traffic-create:has-text("*策略名称英文数字下划线任意一种,首字母必须为英文描述优先级*筛选条件添加条件属性名称属性值操作API请求方式所有API请求方式 配置 删除 API路径test")').click()
    // Click text=[object Text] 请选择
    await page.locator('text=[object Text] 请选择').click()
    // Check input[type="checkbox"] >> nth=0
    await page.locator('input[type="checkbox"]').first().check()
    // Click button:has-text("提交")
    await page.locator('button:has-text("提交")').click()
    // Click text=创建成功!
    await page.locator('text=创建成功!').click()
    // Click text=tet
    await page.locator('text=tet').click()
  })
  it('点击新建策略按钮, 点击左侧分组的环境节点, 该环境的集群列表将收缩或展开, 页面不变; 点击左侧的集群节点, 页面变为相应集群的列表', async () => {
    // Click button:has-text("新建策略")
    await page.locator('button:has-text("新建策略")').click()
    // Click text=英文数字下划线任意一种,首字母必须为英文
    await page.locator('text=英文数字下划线任意一种,首字母必须为英文').click()
    // Click text=DEV
    await page.locator('text=DEV').click()
    // Click text=DEV
    await page.locator('text=DEV').click()
    // Click nz-tree-node-title:has-text("liu_localhost")
    await page.locator('nz-tree-node-title:has-text("liu_localhost")').click()
    // Click text=启停
    await page.locator('text=启停').click()
  })

  it('点击列表中的查看按钮, 策略名称&筛选条件列表&限流规则不为空, 保存按钮为primary, 点击保存出现全局消息提示, 并返回列表页', async () => {
    // Click button:has-text("新建策略")
    await page.locator('button:has-text("新建策略")').click()
    // Click text=*策略名称英文数字下划线任意一种,首字母必须为英文 >> [placeholder="请输入"]
    await page.locator('text=*策略名称英文数字下划线任意一种,首字母必须为英文 >> [placeholder="请输入"]').click()
    // Fill text=*策略名称英文数字下划线任意一种,首字母必须为英文 >> [placeholder="请输入"]
    await page.locator('text=*策略名称英文数字下划线任意一种,首字母必须为英文 >> [placeholder="请输入"]').fill('test')
    // Click button:has-text("添加条件")
    await page.locator('button:has-text("添加条件")').click()
    // Click eo-ng-select-top-control:has-text("应用")
    await page.locator('eo-ng-select-top-control').click()
    // Click text=API路径
    await page.locator('text=API路径').click()

    await page.locator('text=API路径').click()
    // Click [placeholder="请输入API路径"]
    await page.locator('[placeholder="请输入API路径"]').click()
    // Fill [placeholder="请输入API路径"]
    await page.locator('[placeholder="请输入API路径"]').fill('test')
    // Click button:has-text("保存")
    await page.locator('button:has-text("保存")').click()
    // Click text=[object Text] 请选择 >> input
    await page.locator('text=[object Text] 请选择 >> input').click()
    // Click #cdk-overlay-30 >> text=API
    await page.locator('#cdk-overlay-30 >> text=API').click()
    // Click button:has-text("提交")
    await page.locator('button:has-text("提交")').click()
    // Click text=启停
    await page.locator('text=启停').click()
  })

  it('输入框大小, 按钮的大小和颜色, 表格的大小, 限流规则的背景色', async () => {
  })
  it('当筛选条件列表中属性值过长, 鼠标悬浮出现气泡框', async () => {
  })
})
