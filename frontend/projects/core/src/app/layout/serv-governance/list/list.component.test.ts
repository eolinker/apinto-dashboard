// 流量策略列表页
// 优先级编辑
describe('traffic-list e2e test', () => {
  it('初始化页面, 默认第一个环境列表中的第一个集群被选中,列表页显示数据', async () => {
    // Go to http://localhost:4200/serv-goverance/traffic
    await page.goto('http://localhost:4200/serv-goverance/traffic')
    await page.waitForTimeout(2000)
  })
  it('点击优先级和更新时间标题栏的排序按钮, 列表相应排序', async () => {
    const firstP = await page.$('td:nth-child(2) >> nth = 0')

    await page.locator('nz-table-sorters:has-text("优先级")').click()

    let firstPTest = await page.$('td:nth-child(2) >> nth = 0')
    expect(firstPTest).not.toStrictEqual(firstP)

    await page.locator('nz-table-sorters:has-text("优先级")').click()
    firstPTest = await page.$('td:nth-child(2) >> nth = 0')
    expect(firstPTest).toStrictEqual(firstP)

    await page.locator('nz-table-sorters:has-text("更新时间")').click()
    const firstT = await page.$('td:nth-child(8) >> nth = 0')

    await page.locator('nz-table-sorters:has-text("更新时间")').click()
    let firstTTest = await page.$('td:nth-child(8) >> nth = 0')
    expect(firstTTest).not.toStrictEqual(firstT)

    await page.locator('nz-table-sorters:has-text("更新时间")').click()
    firstTTest = await page.$('td:nth-child(8) >> nth = 0')
    expect(firstTTest).toStrictEqual(firstT)
  })

  it('鼠标悬浮在启停标题栏, 提示[策略的限流规则是否被生效执行]', async () => {
    await expect(page.locator('text=策略的限流规则是否被生效执行')).toBeUndefined()
    await page.locator('text=启停').hover()
    await expect(page.locator('text=策略的限流规则是否被生效执行')).not.toBeUndefined()
  })

  it('修改优先级, 将优先级空置, 此时页面出现全局消息提示, 并将输入框标红', async () => {
    await expect(page.locator('div:has-text("优先级不能为空, 请填写后提交")')).toBeUndefined()
    await page.locator('td:nth-child(2) >> nth = 0').fill('')
    await expect(page.locator('div:has-text("优先级不能为空, 请填写后提交")')).not.toBeUndefined()
  })

  it('修改优先级, 将优先级修改为与其他列冲突的等级, 此时页面出现全局消息提示, 并将两个输入框标红', async () => {
    await expect(page.locator('div:has-text("修改后的优先级与test1冲突，无法自动提交")')).toBeUndefined()
    await page.locator('td:nth-child(2) >> nth = 0').fill('1')
    await page.locator('td:nth-child(2) >> nth = 1').fill('1')
    await expect(page.locator('div:has-text("修改后的优先级与test1冲突，无法自动提交")')).not.toBeUndefined()
  })

  it('修改优先级, 将优先级修改为与其他列不冲突的等级, 此时页面出现成功的全局消息提示', async () => {
    await expect(page.locator('div:has-text("修改优先级成功")')).toBeUndefined()
    await page.locator('td:nth-child(2) >> nth = 1').fill('2')
    await expect(page.locator('div:has-text("修改优先级成功")')).toBeUndefined()
  })

  it('点击启停开关, 页面出现全局消息提示(请求发送成功)', async () => {
    await expect(page.locator('div:has-text("停用策略成功")')).toBeUndefined()
    await page.locator('td:nth-child(4) >> nth = 1').click()
    await expect(page.locator('div:has-text("停用策略成功")')).not.toBeUndefined()
  })

  it('点击新建策略按钮,右侧页面将变为新建策略页,点击取消,页面返回列表页', async () => {
    await expect(page.locator('text=英文数字下划线任意一种,首字母必须为英文')).toBeUndefined()
    await page.locator('button:has-text("新建策略")').click()
    await expect(page.locator('text=英文数字下划线任意一种,首字母必须为英文')).not.toBeUndefined()

    await expect(page.locator('text=启停')).toBeUndefined()
    await page.locator('button:has-text("取消")').click()
    await expect(page.locator('text=启停')).not.toBeUndefined()
  })

  it('点击发布按钮,页面出现弹窗,清空发布名称则不可发布,点击取消,页面返回列表页', async () => {
    await expect(page.locator('text=发布名称')).toBeUndefined()
    await page.locator('button:has-text("发布")').click()
    await expect(page.locator('text=发布名称')).not.toBeUndefined()

    await expect(page.locator('button:has-text("保存")').isDisabled).toStrictEqual(false)
    await page.locator('.mg_content input').fill('')
    await expect(page.locator('button:has-text("保存")').isDisabled).toStrictEqual(true)

    await page.locator('button:has-text("取消")').click()
    await expect(page.locator('text=发布名称')).toBeUndefined()
  })

  it('点击发布按钮,页面出现弹窗,后端未返回source字段时不可发布,点击关闭,页面返回列表页', async () => {
    await page.locator('button:has-text("发布")').click()

    await page.locator('[aria-label="Close"]').click()
  })

  it('点击发布按钮,页面出现弹窗,点击保存,页面返回列表页并出现全局消息提示', async () => {
    await expect(page.locator('div:has-text("发布策略成功!")')).toBeUndefined()
    await page.locator('button:has-text("发布")').click()
    await expect(page.locator('text=发布名称')).not.toBeUndefined()
    await page.locator('button:has-text("保存")').click()
    await expect(page.locator('text=发布名称')).toBeUndefined()
    await expect(page.locator('div:has-text("发布策略成功!")')).not.toBeUndefined()
  })

  it('点击列表中的查看按钮,右侧页面变为编辑策略页,策略名称与列表中策略名相同, 点击保存, 页面返回列表页并出现全局消息提示', async () => {
    await page.locator('td:nth-child(9) >> nth=0 >>text=查看').click()
    await expect(page.locator('text=test1')).not.toBeUndefined()
    await expect(page.locator('text=英文数字下划线任意一种,首字母必须为英文')).not.toBeUndefined()

    await expect(page.locator('div:has-text("修改成功!")')).toBeUndefined()
    await page.locator('button:has-text("保存")').click()
    await expect(page.locator('div:has-text("修改成功!")')).not.toBeUndefined()
  })

  it('点击列表中的查看按钮,右侧页面变为编辑策略页,策略名称与列表中策略名相同, 点击取消, 页面返回列表页', async () => {
    await expect(page.locator('button:has-text("新建策略")')).not.toBeUndefined()
    await page.locator('td:nth-child(9) >> nth=0 >>text=查看').click()
    await expect(page.locator('text=英文数字下划线任意一种,首字母必须为英文')).not.toBeUndefined()
    await expect(page.locator('button:has-text("新建策略")')).toBeUndefined()

    await page.locator('button:has-text("取消")').click()
    await expect(page.locator('button:has-text("新建策略")')).not.toBeUndefined()
    await expect(page.locator('text=英文数字下划线任意一种,首字母必须为英文')).not.toBeUndefined()
  })

  it('点击列表中的删除按钮 页面出现全局消息提示', async () => {
    await page.locator('td:nth-child(9) >> nth=0 >>text=删除').click()
    await page.locator('button:has-text("确定")').click()
  })

  it('点击列表中的恢复按钮 页面出现全局消息提示', async () => {
    await expect(page.locator('text=恢复策略成功')).toBeUndefined()

    await page.locator('text=恢复').click()
    await expect(page.locator('text=恢复策略成功')).not.toBeUndefined()
  })

  it('列表每列的width和每行的height, 新建策略按钮和发布按钮的大小与背景色', async () => {
  })
  it('列表每栏的width, 新建策略按钮和发布按钮的大小与背景色', async () => {
  })
  it('发布状态待更新-字体红色,已上线-绿色,待删除-橘色, 未上线-黑色', async () => {
  })
  it('启停状态 启动-switch背景色为绿色, 停用-背景色为灰色', async () => {
  })
})
