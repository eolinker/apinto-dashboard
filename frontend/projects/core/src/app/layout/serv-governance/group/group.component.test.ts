// 流量策略左侧分组组件
// 需要补充:刷新停留本页
describe('flow-control-group e2e test', () => {
  it('初始化页面, 默认第一个环境列表中的第一个集群被选中', async () => {
    // Go to http://localhost:4200/serv-goverance/flow-control
    await page.goto('http://localhost:4200/serv-goverance/flow-control')
    await page.waitForTimeout(2000)

    await expect(page.locator('nz-tree-indent + nz-tree-node-title').nth(1).getAttribute('class')).toContain('ant-tree-node-selected')
    await expect(page.locator('nz-tree-indent + nz-tree-node-title').nth(2).getAttribute('class')).not.toContain('ant-tree-node-selected')
  })
  it('点击分组中的一级节点(环境), 一级节点展开/收缩,右侧页面中的列表数据不变', async () => {
    const rightSide = await page.$('.block_right')
    let rightSideTest = await page.$('.block_right')
    // 第一个环境中的集群可见
    await expect(page.locator('nz-tree-node-title:has-text("liu_localhost")')).not.toBeNull()

    // Click nz-tree-node-title:has-text("PRO")
    await page.locator('nz-tree-node-title:has-text("PRO")').click()
    rightSideTest = await page.$('.block_right')
    // 第一个环境中的集群不可见
    await expect(page.locator('nz-tree-node-title:has-text("liu_localhost")')).toBeNull()
    await expect(page.locator('nz-tree-node-title:has-text("apinto")')).toBeNull()
    await expect(rightSide).toStrictEqual(rightSideTest)

    // Click nz-tree-node-title:has-text("DEV")
    await page.locator('nz-tree-node-title:has-text("DEV")').click()
    rightSideTest = await page.$('.block_right')
    await expect(page.locator('nz-tree-node-title:has-text("apinto")')).not.toBeNull()
    await expect(rightSide).toStrictEqual(rightSideTest)
  })
  it('点击分组中的二级节点(集群), 右侧页面中的列表数据发生变化', async () => {
    await page.locator('nz-tree-indent + nz-tree-node-title').nth(1).click()
    const rightSide = await page.$('.block_right')
    let rightSideTest = await page.$('.block_right')

    await page.locator('nz-tree-indent + nz-tree-node-title').nth(2).click()
    rightSideTest = await page.$('.block_right')
    await expect(rightSide).not.toStrictEqual(rightSideTest)

    await page.locator('nz-tree-indent + nz-tree-node-title').nth(3).click()
    rightSideTest = await page.$('.block_right')
    await expect(rightSide).not.toStrictEqual(rightSideTest)
  })
  it('点击分组中的二级节点(集群), 被选中的节点背景色变化为--HOVER_BG', async () => {
    await page.locator('nz-tree-indent + nz-tree-node-title').nth(2).click()
    await expect(page.locator('nz-tree-indent + nz-tree-node-title').nth(1).getAttribute('styles')).not.toContain('{background-color:var(--MAIN_THEME_COLOR)}')
    await expect(page.locator('nz-tree-indent + nz-tree-node-title').nth(2).getAttribute('styles')).toContain('{background-color:var(--MAIN_THEME_COLOR)}')
  })
  it('鼠标悬浮的节点背景色变化为--HOVER_BG', async () => {
    await page.locator('nz-tree-indent + nz-tree-node-title').nth(1).hover()
    await expect(page.locator('nz-tree-indent + nz-tree-node-title').nth(1).getAttribute('styles')).not.toContain('{background-color:var(--MAIN_THEME_COLOR)}')
  })
})
