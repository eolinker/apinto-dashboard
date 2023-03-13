describe('用户角色 e2e test', () => {
  it('初始化进入用户角色页面', async () => {
    await page.goto('http://localhost:4200/login')
    await page.waitForTimeout(2000)
    await page.getByPlaceholder('请输入账号').click()
    await page.getByPlaceholder('请输入账号').fill('maggie')
    await page.getByPlaceholder('请输入账号').press('Tab')
    await page.getByPlaceholder('请输入密码').fill('12345678')
    await page.getByPlaceholder('请输入密码').press('Enter')
    await page.locator('eo-ng-menu-default').getByText('系统管理').click()
    await page.locator('eo-ng-menu-default').getByRole('link', { name: '用户角色' }).click()
  })
  it('分组样式，页面样式检查', async () => {
    // 分组组件宽234px
    const groupBlock = await page.locator('.block-left')
    const groupBlockW = await groupBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(groupBlockW).toStrictEqual('234px')

    // 面包屑样式 面包屑API管理，字号14，字体颜色为主题色，与左侧距离12px，垂直居中
    const headerBlock = await page.locator('.block_rl')
    const headerBlockAI = await headerBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('align-items'))
    expect(headerBlockAI).toStrictEqual('center')

    const breadcrumbItem = await page.locator('nz-breadcrumb-item').getByRole('link', { name: '用户角色' })
    const breadcrumbItemFS = await breadcrumbItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(breadcrumbItemFS).toStrictEqual('14px')
    const breadcrumbItemC = await breadcrumbItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(breadcrumbItemC).toStrictEqual('rgb(34, 84, 157)')

    const itemFTitle = await page.locator('eo-ng-tree-default-node').first().locator('nz-tree-node-title')
    const itemFTitleML = await itemFTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    const itemFTitleP = await itemFTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(itemFTitleML).toStrictEqual('12px')
    expect(itemFTitleP).toStrictEqual('0px 4px')

    const itemLTitle = await page.locator('eo-ng-tree-default-node').last().locator('nz-tree-node-title')
    const itemLTitleML = await itemLTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    const itemLTitleP = await itemFTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(itemLTitleP).toStrictEqual('0px 4px')
    expect(itemLTitleML).toStrictEqual('12px')

    // 所有选项234px*30px，外侧无间距，背景色白色，鼠标悬浮与点击时，背景色变悬浮色
    const groupTitle = await page.locator('.group-title')
    const groupTitleH = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(groupTitleH).toContain('30')
    const groupTitleW = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(groupTitleW).toContain('224')
    const groupTitleM = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    expect(groupTitleM).toStrictEqual('0px')
    const groupTitleP = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(groupTitleP).toStrictEqual('4px 16px')

    await itemFTitle.click()

    const groupTitleBG = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(groupTitleBG).toStrictEqual('rgba(0, 0, 0, 0)')
    await groupTitle.hover()
    const groupTitleHBG = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(groupTitleHBG).toStrictEqual('rgb(240, 247, 255)')
    await groupTitle.click()
    const groupTitleCBG = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(groupTitleCBG).toStrictEqual('rgb(240, 247, 255)')

    const itemFTitleH = await itemFTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(itemFTitleH).toContain('30px')
    const itemFTitleW = await itemFTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(itemFTitleW).toContain('224')
    const itemFTitleBG = await itemFTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(itemFTitleBG).toStrictEqual('rgba(0, 0, 0, 0)')
    await itemFTitle.hover()
    const itemFTitleHBG = await itemFTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(itemFTitleHBG).toStrictEqual('rgb(240, 247, 255)')
    await itemFTitle.click()
    const itemFTitleCBG = await itemFTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(itemFTitleCBG).toStrictEqual('rgb(240, 247, 255)')

    const createRoleDiv = await page.locator('div.center.blue')
    const createRoleDivPL = await createRoleDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    expect(createRoleDivPL).toStrictEqual('4px')
    const createRoleDivML = await createRoleDiv.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(createRoleDivML).toStrictEqual('12px')

    const createRoleDivA = await page.locator('div.center.blue a')
    const createRoleDivAC = await createRoleDivA.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(createRoleDivAC).toStrictEqual('rgb(34, 84, 157)')

    await createRoleDivA.hover()
    expect(await createRoleDivA.evaluate((element) => window.getComputedStyle(element).getPropertyValue('color')))
      .not.toStrictEqual('rgb(34, 84, 157)')
  })

  it('列表页面样式检查；', async () => {
    const createBtn = await page.getByRole('button', { name: '新建用户' })

    const createBtnH = await createBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(createBtnH).toStrictEqual('32px')
    const createBtnML = await createBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(createBtnML).toStrictEqual('12px')
    const createBtnP = await createBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(createBtnP).toStrictEqual('0px 12px')
    const createBtnBG = await createBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(createBtnBG).toStrictEqual('rgb(34, 84, 157)')
    const createBtnBDC = await createBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    expect(createBtnBDC).toStrictEqual('rgb(34, 84, 157)')

    const deleteBtn = await page.getByRole('button', { name: '删除' })

    const deleteBtnH = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(deleteBtnH).toStrictEqual('32px')
    const deleteBtnML = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(deleteBtnML).toStrictEqual('12px')
    const deleteBtnP = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(deleteBtnP).toStrictEqual('0px 12px')
    const deleteBtnDisabled = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    expect(deleteBtnDisabled).toStrictEqual('true')
    const deleteBtnBG = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(deleteBtnBG).toStrictEqual('rgb(245, 245, 245)')
    const deleteBtnBDC = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    expect(deleteBtnBDC).toStrictEqual('rgb(217, 217, 217)')
    const deleteBtnC = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(deleteBtnC).toStrictEqual('rgba(0, 0, 0, 0.25)')

    const searchInputGroup = await page.locator('.group-search-large.mg-top-right')
    const searchInputGroupMR = await searchInputGroup.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-right'))
    expect(searchInputGroupMR).toStrictEqual('24px')

    const searchInput = await page.locator('.group-search-large.mg-top-right eo-ng-input-group')
    const searchInputH = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(searchInputH).toStrictEqual('32px')
    const searchInputW = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(searchInputW).toStrictEqual('254px')
    const searchInputML = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(searchInputML).toStrictEqual('12px')

    const listHeader = await page.locator('.list-header')
    const listHeaerPT = await listHeader.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-top'))
    expect(listHeaerPT).toStrictEqual('12px')

    const listContent = await page.locator('.list-content')
    const listContentPT = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-top'))
    expect(listContentPT).toStrictEqual('0px')
    const listContentMT = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    expect(listContentMT).toStrictEqual('12px')
    const listContentPB = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-bottom'))
    expect(listContentPB).toStrictEqual('20px')
  })

  it('点击分组，点击新建自定义角色，检查样式；', async () => {
    await page.getByText('新建自定义角色').click()
    // API名称
    const nameInput = await page.getByPlaceholder('请输入角色名称')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(nameInputW).toStrictEqual('346px')
    await expect(nameInputH).toStrictEqual('32px')

    // 描述样式
    const descInput = await page.locator('textarea')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(descInputW).toStrictEqual('346px')
    await expect(descInputH).toStrictEqual('68px')

    // 布局
    const drawerHeader = await page.locator('.drawer-list-header')
    const drawerHeaderH = await drawerHeader.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    await expect(drawerHeaderH).toStrictEqual('140px')

    const drawerContent = await page.locator('.drawer-list-content')
    const drawerContentPB = await drawerContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-bottom'))

    await expect(drawerContentPB).toStrictEqual('20px')

    const drawerFooter = await page.locator('.drawer-list-content')
    const drawerFooterP = await drawerFooter.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    const drawerFooterM = await drawerFooter.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))

    await expect(drawerFooterP).toStrictEqual('0px')
    await expect(drawerFooterM).toStrictEqual('0px auto')
  })
  it('新建自定义角色，在分组中删除该角色', async () => {
    await page.getByText('新建自定义角色').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByPlaceholder('请输入角色名称').click()
    await page.getByPlaceholder('请输入角色名称').fill('testForE2e')
    await page.getByRole('button', { name: '保存' }).click()

    await page.locator('nz-tree-node-title').filter({ hasText: 'testForE2e(0)' }).locator('div').nth(1).click()

    const groupTitle = await page.locator('.group-title')
    const groupTitleBG = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(groupTitleBG).toStrictEqual('rgba(0, 0, 0, 0)')
    await groupTitle.hover()
    const groupTitleHBG = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(groupTitleHBG).toStrictEqual('rgb(240, 247, 255)')

    await page.locator('nz-tree-node-title').filter({ hasText: 'testForE2e(0)' }).locator('div').nth(1).click()

    await page.locator('nz-tree-node-title').filter({ hasText: 'testForE2e(0)' }).getByRole('button', { name: '' }).click()
    await page.getByText('编辑').click()
    await page.getByPlaceholder('请输入角色名称').click()
    await page.getByPlaceholder('请输入角色名称').fill('testForE2e2')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('nz-tree-node-title').filter({ hasText: 'testForE2e2(0)' }).locator('div').nth(1).click()
    await page.locator('nz-tree-node-title').filter({ hasText: 'testForE2e2(0)' }).getByRole('button', { name: '' }).click()
    await page.getByText('删除').click()
    await page.getByText('如果该角色关联了用户，删除后关联的用户将失去该角色权限并会成为未分配角色用户，请确认是否删除？').click()
    await page.getByRole('button', { name: '确定' }).click()

    const groupTitleCBG = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(groupTitleCBG).toStrictEqual('rgb(240, 247, 255)')
  })
  it('点击新建用户，检查样式', async () => {
    await page.getByRole('button', { name: '新建用户' }).click()
    // 账号
    const accountInput = await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文')
    const accountInputW = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const accountInputH = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(accountInputW).toStrictEqual('346px')
    await expect(accountInputH).toStrictEqual('32px')

    // 名称
    const nameInput = await page.getByPlaceholder('请输入名称')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(nameInputW).toStrictEqual('346px')
    await expect(nameInputH).toStrictEqual('32px')

    // 邮箱
    const emailInput = await page.getByPlaceholder('请输入邮箱')
    const emailInputW = await emailInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const emailInputH = await emailInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(emailInputW).toStrictEqual('346px')
    await expect(emailInputH).toStrictEqual('32px')

    // 角色
    const roleInput = await page.locator('.ant-modal-body eo-ng-select-top-control')
    const roleInputW = await roleInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const roleInputH = await roleInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(roleInputW).toStrictEqual('346px')
    await expect(roleInputH).toStrictEqual('32px')

    // 描述样式
    const descInput = await page.locator('textarea#desc')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(descInputW).toStrictEqual('346px')
    await expect(descInputH).toStrictEqual('68px')
  })
  it('逐一填入必输项，检查是否可以提交；提交后返回列表；点击该用户，选择删除', async () => {
    await page.getByRole('button', { name: '新建用户' }).click()
    await page.getByRole('button', { name: '提交' }).click()
    await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').click()
    await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').fill('test2')
    await page.getByRole('button', { name: '提交' }).click()
    await page.getByPlaceholder('请输入名称').click()
    await page.getByPlaceholder('请输入名称').fill('testForE2e')
    await page.getByRole('button', { name: '提交' }).click()
    await page.getByPlaceholder('请输入邮箱').click()
    await page.getByPlaceholder('请输入邮箱').fill('test')
    await page.getByRole('button', { name: '提交' }).click()
    await page.getByPlaceholder('请输入邮箱').click()
    await page.getByPlaceholder('请输入邮箱').fill('test@1')
    await page.getByRole('button', { name: '提交' }).click()

    // await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 8 >> button >> nth = 2').click()

    // await page.getByText('该数据删除后将无法找回，请确认是否删除？').click()
    // await page.getByRole('button', { name: '确定' }).click()
  })
  // it('点击新建用户，检查样式，填入所有必填项，提交', async () => {
  //   await page.getByRole('button', { name: '新建用户' }).click();
  //   await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').click();
  //   await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文').fill('test2');
  //   await page.getByPlaceholder('请输入名称').click();
  //   await page.getByPlaceholder('请输入名称').fill('est2');
  //   await page.getByPlaceholder('请输入邮箱').click();
  //   await page.getByPlaceholder('请输入邮箱').fill('test@1');
  //   await page.locator('nz-select-item').click();
  //   await page.locator('#cdk-overlay-20').getByText('测试角色').click();
  //   await page.locator('textarea').click();
  //   await page.locator('textarea').fill('testdesc');
  //   await page.getByRole('button', { name: '提交' }).click();
  // })
  it('选择新增用户，点击删除，取消删除；点击用户，出现用户详情弹窗,其中账号不可修改；修改名称和邮箱，点击保存', async () => {
    const deleteBtn = await page.getByRole('button', { name: '删除' })

    const deleteBtnH = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(deleteBtnH).toStrictEqual('32px')
    const deleteBtnML = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(deleteBtnML).toStrictEqual('12px')
    const deleteBtnP = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(deleteBtnP).toStrictEqual('0px 12px')
    const deleteBtnBG = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(deleteBtnBG).toStrictEqual('rgb(255, 255, 255)')
    const deleteBtnBDC = await deleteBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    expect(deleteBtnBDC).toStrictEqual('rgb(217, 217, 217)')
  })
  it('选择修改后的用户，账号不可修改，修改名称后选择取消，返回列表', async () => {
    await page.locator('eo-ng-apinto-table').getByRole('cell', { name: 'test' }).first()
    const accountInput = await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文')
    const accountInputDisabled = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    expect(accountInputDisabled).toStrictEqual(true)

    await page.getByPlaceholder('请输入名称').click()
    await page.getByPlaceholder('请输入名称').fill('testForE2e')
    await page.locator('nz-select-item').click()
    await page.locator('eo-ng-option-item').filter({ hasText: '未分配' }).click()
    await page.locator('textarea').click()
    await page.locator('textarea').fill('test')
    await page.getByRole('button', { name: '保存' }).click()
  })
  it('选择超级管理员，检查样式', async () => {
    await page.locator('eo-ng-tree-default-node').filter({ hasText: '超级管理员-内置' }).click()
    expect(await page.getByRole('button', { name: '删除' }).isVisible()).toStrictEqual(false)
    expect(await page.getByRole('button', { name: '新建用户' }).isVisible()).toStrictEqual(false)
    expect(await page.getByRole('button', { name: '添加用户' }).isVisible()).toStrictEqual(false)
    expect(await page.getByRole('button', { name: '移除' }).isVisible()).toStrictEqual(false)
    expect(await page.locator('.list-content').evaluate((element) => window.getComputedStyle(element).getPropertyValue('padding-top'))).toStrictEqual('0px')
    expect(await page.locator('.list-content').evaluate((element) => window.getComputedStyle(element).getPropertyValue('margin-top'))).toStrictEqual('0px')
    expect(await (await page.$$('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 6 >> button')).length).toStrictEqual(0)
  })
  it('新建自定义角色，在分组中取消删除该角色；点击新建的角色，选择添加用户，查看样式', async () => {
    await page.getByText('新建自定义角色').click()
    await page.getByPlaceholder('请输入角色名称').click()
    await page.getByPlaceholder('请输入角色名称').fill('test')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('nz-tree-node-title').filter({ hasText: 'test(0)' }).last().locator('div').nth(1).click()
    await page.getByRole('button', { name: '添加用户' }).click()

    const searchBlock = await page.locator('.drawer-table > .block_lr > .floatR')
    const searchBlockF = await searchBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('float'))
    expect(searchBlockF).toStrictEqual('right')
    const searchBlockP = await searchBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(searchBlockP).toStrictEqual('0px')
    const searchBlockM = await searchBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    expect(searchBlockM).toStrictEqual('0px auto')

    const searchInput = await page.locator('.ant-drawer-body eo-ng-input-group')
    const searchInputM = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    expect(searchInputM).toStrictEqual('0px 0px 0px 12px')

    const listContent = await page.locator('.list-content')
    const listContentM = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    expect(listContentM).toStrictEqual('12px auto 0px auto')
    const listContentP = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(listContentP).toStrictEqual('0px 0px 20px 0px')

    // 保存按钮样式
    const saveBtn = await page.getByRole('button', { name: '保存' }).last()
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
    const saveBtnML = await saveBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    expect(saveBtnH).toStrictEqual('32px')
    expect(saveBtnW).toStrictEqual('54px')
    expect(saveBtnBG).toStrictEqual('rgb(34, 84, 157)')
    expect(saveBtnBC).toStrictEqual('rgb(34, 84, 157)')
    expect(saveBtnFS).toStrictEqual('14px')
    expect(saveBtnML).toStrictEqual('12px')

    // 取消按钮样式
    const cancleBtn = await page.getByRole('button', { name: '取消' }).last()
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
    expect(cancleBtnML).toStrictEqual('24px')
  })
  it('在新建的角色里添加用户，再选择该用户，选择移除；编辑新建的角色权限为只读，将test2移入该角色，等候测试无编辑权限用', async () => {
    await page.locator('nz-tree-node-title').filter({ hasText: 'test(0)' }).locator('div').nth(2).click()
    await page.getByRole('button', { name: '添加用户' }).click()
    await page.getByRole('row', { name: 'test2 testForE2e test@1 未分配' }).getByLabel('').check()
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByRole('button', { name: '' }).click()
    await page.getByText('请确认是否移除指定用户的角色权限？').click()
    await page.getByRole('button', { name: '确定' }).click()
    await page.getByRole('cell', { name: '暂无数据' }).locator('div').filter({ hasText: '暂无数据' }).click()
    await page.getByRole('button', { name: '添加用户' }).click()
    await page.getByRole('row', { name: 'test2 testForE2e test@1 未分配' }).getByLabel('').check()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('nz-select-item').click()
    await page.locator('eo-ng-option-item').filter({ hasText: '测试角色' }).click()
    await page.getByRole('cell', { name: '暂无数据' }).locator('svg').click()
    await page.getByText('测试角色(2)').click()
    await page.getByRole('cell', { name: 'test2' }).getByText('test2').click()
    await page.locator('.ant-drawer-body eo-ng-option-item').getByText('测试角色').click()
    await page.getByRole('button', { name: '取消' }).click()

    await page.locator('nz-tree-node-title').filter({ hasText: 'test(0)' }).locator('div').nth(1).click()
    await page.locator('nz-tree-node-title').filter({ hasText: 'test(0)' }).getByRole('button', { name: '' }).click()
    await page.getByText('编辑').click()
    await page.getByRole('row', { name: 'API管理 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '服务发现 查看 编辑' }).getByLabel('查看').uncheck()
    await page.getByRole('row', { name: '基础设施 网关集群 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '服务发现 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '环境变量 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '服务治理 流量策略 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '熔断策略 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '访问策略 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '缓存策略 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '灰度策略 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '系统管理 用户角色 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '外部应用 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '审计日志 查看' }).getByLabel('查看').check()
    await page.getByRole('row', { name: '商业授权 查看 编辑' }).getByLabel('查看').check()
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByRole('button', { name: '添加用户' }).click()
    await page.getByRole('row', { name: 'test2 testForE2e test@1 测试角色' }).getByLabel('').check()
    await page.getByRole('button', { name: '保存' }).click()
  })
  it('用户头像，用户设置的样式与操作', async () => {
    await page.getByRole('button', { name: ' maggie' }).click()
    await page.getByText('用户设置').click()

    // 账号
    const accountInput = await page.getByPlaceholder('英文数字下划线任意一种，首字母必须为英文')
    const accountInputDisabled = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    await expect(accountInputDisabled).toStrictEqual('true')

    // 名称
    const nameInput = await page.getByPlaceholder('请输入名称')
    const nameInputDisabled = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    await expect(nameInputDisabled).toStrictEqual('false')

    // 邮箱
    const emailInput = await page.getByPlaceholder('请输入邮箱')
    const emailInputDisabled = await emailInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    await expect(emailInputDisabled).toStrictEqual('false')

    // 角色
    const roleInput = await page.locator('.ant-modal-body eo-ng-select-top-control')
    const roleInputDisabled = await roleInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    await expect(roleInputDisabled).toStrictEqual('true')

    // 描述样式
    const descInput = await page.locator('textarea#desc')
    const descInputDisabled = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))
    await expect(descInputDisabled).toStrictEqual('false')

    await page.getByPlaceholder('请输入名称').click()
    await page.getByPlaceholder('请输入邮箱').click()
    await page.getByRole('button', { name: '保存' }).click()
  })

  it('用户头像，修改密码的样式与操作', async () => {
    await page.getByRole('button', { name: ' maggie' }).click()
    await page.getByText('修改密码').click()

    // 账号
    const accountInput = await page.locator('nz-form-item').filter({ hasText: '账号' }).getByRole('textbox')
    const accountInputW = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const accountInputH = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const accountInputDisabled = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(accountInputW).toStrictEqual('346px')
    await expect(accountInputH).toStrictEqual('32px')
    await expect(accountInputDisabled).toStrictEqual('true')

    // 旧密码
    const oldPswInput = await page.locator('nz-form-item').filter({ hasText: '旧密码' }).getByPlaceholder('请输入6-32位字符')
    const oldPswInputW = await oldPswInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const oldPswInputH = await oldPswInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const oldPswInputDisabled = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(oldPswInputW).toStrictEqual('346px')
    await expect(oldPswInputH).toStrictEqual('32px')
    await expect(oldPswInputDisabled).toStrictEqual('false')

    // 新密码
    const newPswInput = await page.locator('nz-form-item').filter({ hasText: '旧密码' }).getByPlaceholder('请输入6-32位字符')
    const newPswInputW = await newPswInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const newPswInputH = await newPswInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const newPswInputDisabled = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(newPswInputW).toStrictEqual('346px')
    await expect(newPswInputH).toStrictEqual('32px')
    await expect(newPswInputDisabled).toStrictEqual('false')

    // 确认新密码
    const newPswConfirmInput = await page.locator('nz-form-item').filter({ hasText: '旧密码' }).getByPlaceholder('请输入6-32位字符')
    const newPswConfirmInputW = await newPswConfirmInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const newPswConfirmInputH = await newPswConfirmInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const newPswConfirmInputDisabled = await accountInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))

    await expect(newPswConfirmInputW).toStrictEqual('346px')
    await expect(newPswConfirmInputH).toStrictEqual('32px')
    await expect(newPswConfirmInputDisabled).toStrictEqual('false')

    // 操作
    await page.locator('nz-form-item').filter({ hasText: '旧密码' }).getByPlaceholder('请输入6-32位字符').click()
    await page.locator('nz-form-item').filter({ hasText: '旧密码' }).getByPlaceholder('请输入6-32位字符').fill('12345678')
    await page.getByPlaceholder('请输入6-32位字符').nth(1).click()
    await page.getByPlaceholder('请输入6-32位字符').nth(1).fill('12345678')
    await page.locator('nz-form-item').filter({ hasText: '确认新密码' }).getByPlaceholder('请输入6-32位字符').click()
    await page.locator('nz-form-item').filter({ hasText: '确认新密码' }).getByPlaceholder('请输入6-32位字符').fill('12345678')
    await page.getByText('密码强度：弱，建议使用英文、数字、特殊字符组合').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('nz-form-control').filter({ hasText: '密码强度：弱，建议使用英文、数字、特殊字符组合' }).getByPlaceholder('请输入6-32位字符').click()
    await page.locator('nz-form-control').filter({ hasText: '密码强度：弱，建议使用英文、数字、特殊字符组合' }).getByPlaceholder('请输入6-32位字符').fill('12345678abv')
    await page.locator('nz-form-control').filter({ hasText: '新密码与确认新密码不一致' }).getByPlaceholder('请输入6-32位字符').click()
    await page.locator('nz-form-control').filter({ hasText: '新密码与确认新密码不一致' }).getByPlaceholder('请输入6-32位字符').fill('12345678abv')
    await page.getByRole('button', { name: '保存' }).click()

    if (await page.getByRole('button', { name: '取消' }) && await page.getByRole('button', { name: '取消' }).isVisible()) {
      await page.getByRole('button', { name: '取消' }).click()
    }

    await page.locator('eo-ng-apinto-table tr:has-text("testForE2e") >> td >> nth = 8 >> button >> nth = 0').first().click()
    await page.getByText('系统默认重置后密码为12345678，请确认是否重置？').click()
    await page.getByRole('button', { name: '确定' }).click()
    await page.locator('eo-ng-apinto-table tr:has-text("testForE2e")').locator('eo-ng-switch').getByRole('button').first().click()
    await page.locator('eo-ng-apinto-table tr:has-text("testForE2e")').locator('eo-ng-switch').getByRole('button').first().click()
  })
  it('退出登录，登录无编辑权限角色，测试所有页面的编辑，再退出登录', async () => {
    // 登录无编辑权限角色
    await page.getByRole('button', { name: ' maggie' }).click()
    await page.getByText('退出登录').click()
    await page.getByRole('heading', { name: '欢迎来到 Apinto' }).click()
    await page.getByPlaceholder('请输入账号').click()
    await page.getByPlaceholder('请输入账号').fill('test2')
    await page.getByPlaceholder('请输入账号').press('Tab')
    await page.getByPlaceholder('请输入密码').fill('12345678')
    await page.getByPlaceholder('请输入密码').press('Enter')

    // 测试API页面
    await page.locator('eo-ng-menu-default').getByRole('link', { name: 'API管理' }).click()
    await page.locator('nz-breadcrumb-item').getByRole('link', { name: 'API管理' }).click()

    const apiCreateBtn = await page.getByRole('button', { name: '新建API' })
    expect(await apiCreateBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const apiImportBtn = await page.getByRole('button', { name: '导入' })
    expect(await apiImportBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const apiBatchOnlineBtn = await page.getByRole('button', { name: '上线' })
    expect(await apiBatchOnlineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const apiBatchOfflineBtn = await page.getByRole('button', { name: '下线' })
    expect(await apiBatchOfflineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const addGroupBtn = await page.locator('button.ant-table-row-expand-icon')
    expect(await addGroupBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const dropDownBtn = await page.locator('.custom-node button').first()
    expect(await dropDownBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('visibility'))).toStrictEqual('hidden')

    expect(await (await page.locator('.eo-table-btn-td button >> nth = 2').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 6 >> button >> nth = 0').click()

    const onlineIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shangxianicon') }).first()
    if (onlineIconBtn && await onlineIconBtn.isVisible()) {
      expect(await onlineIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const offlineIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-xiaxianicon') }).first()
    if (offlineIconBtn && await offlineIconBtn.isVisible()) {
      expect(await offlineIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const updateIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shuaxinjiankongzhuangtai') }).first()
    if (updateIconBtn && await updateIconBtn.isVisible()) {
      expect(await updateIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const disabledIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-zanting') }).first()
    if (disabledIconBtn && await disabledIconBtn.isVisible()) {
      expect(await disabledIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const enableIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-zhihang') }).first()
    if (enableIconBtn && await enableIconBtn.isVisible()) {
      expect(await enableIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.getByRole('link', { name: 'API信息' }).click()

    expect(await (await page.locator('eo-ng-tree-select div')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('#desc')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 5')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 6')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 7')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 8')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 9')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 10')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 11')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 12')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 13')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '添加配置' }).first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '添加配置' }).last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
  })

  it('无权限角色浏览应用管理页面', async () => {
    await page.getByRole('link', { name: '应用管理' }).click()
    expect(await (await page.getByRole('button', { name: '新建应用' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.eo-table-btn-td button >> nth = 1').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('cell', { name: '匿名应用' }).click()

    const onlineIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shangxianicon') }).first()
    if (onlineIconBtn && await onlineIconBtn.isVisible()) {
      expect(await onlineIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const offlineIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-xiaxianicon') }).first()
    if (offlineIconBtn && await offlineIconBtn.isVisible()) {
      expect(await offlineIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const updateIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shuaxinjiankongzhuangtai') }).first()
    if (updateIconBtn && await updateIconBtn.isVisible()) {
      expect(await updateIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const disabledIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-zanting') }).first()
    if (disabledIconBtn && await disabledIconBtn.isVisible()) {
      expect(await disabledIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const enableIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-zhihang') }).first()
    if (enableIconBtn && await enableIconBtn.isVisible()) {
      expect(await enableIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.getByRole('link', { name: '应用信息' }).click()

    expect(await (await page.locator('#desc')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 5')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('table button').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('table button').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('cell', { name: 'weweeww' }).click()
    await page.getByRole('link', { name: '鉴权管理' }).click()

    expect(await (await page.getByRole('button', { name: '配置鉴权' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('cell', { name: 'basic-test' }).click()

    expect(await (await page.locator('eo-ng-switch').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 5')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 6')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body nz-date-picker')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body').getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
    await page.getByRole('cell', { name: 'apikey-test' }).click()
    expect(await (await page.locator('eo-ng-switch').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 5')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body nz-date-picker')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body').getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
    await page.getByRole('cell', { name: 'aksk-test' }).click()

    expect(await (await page.locator('eo-ng-switch').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 5')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 6')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body nz-date-picker')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body').getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
    await page.getByRole('cell', { name: 'jwt-test2' }).click()

    expect(await (await page.locator('eo-ng-switch').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select >> nth = 1').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 5')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 6')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body nz-date-picker')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body').getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
    await page.getByRole('cell', { name: 'jwt-test3' }).click()

    expect(await (await page.locator('eo-ng-switch').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select >> nth = 1').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body eo-ng-select').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 5')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body input.ant-input >> nth = 6')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body nz-date-picker')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body').getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
  })
  it('无权限角色浏览上游服务-上游管理页面', async () => {
    await page.locator('eo-ng-menu-default div').filter({ hasText: '上游服务' }).click()
    await page.getByRole('link', { name: '上游管理' }).click()

    expect(await (await page.getByRole('button', { name: '新建上游' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.eo-table-btn-td button >> nth = 1').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 0').click()

    const onlineIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shangxianicon') }).first()
    if (onlineIconBtn && await onlineIconBtn.isVisible()) {
      expect(await onlineIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const offlineIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-xiaxianicon') }).first()
    if (offlineIconBtn && await offlineIconBtn.isVisible()) {
      expect(await offlineIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const updateIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shuaxinjiankongzhuangtai') }).first()
    if (updateIconBtn && await updateIconBtn.isVisible()) {
      expect(await updateIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.getByRole('link', { name: '上游信息' }).click()

    expect(await (await page.locator('#desc')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-select >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-select >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-select >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
  })
  it('无权限角色浏览上游服务-服务发现页面', async () => {
    await page.getByRole('link', { name: '服务发现' }).click()

    expect(await (await page.getByRole('button', { name: '新建服务' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.eo-table-btn-td button >> nth = 1').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 0').click()

    const onlineIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shangxianicon') }).first()
    if (onlineIconBtn && await onlineIconBtn.isVisible()) {
      expect(await onlineIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const offlineIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-xiaxianicon') }).first()
    if (offlineIconBtn && await offlineIconBtn.isVisible()) {
      expect(await offlineIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const updateIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shuaxinjiankongzhuangtai') }).first()
    if (updateIconBtn && await updateIconBtn.isVisible()) {
      expect(await updateIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.getByRole('link', { name: '服务信息' }).click()

    expect(await (await page.locator('#desc')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-select >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-select >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 5')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 6')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 7')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
  })
  it('无权限角色浏览基础设施-网关集群页面 ', async () => {
    await page.getByText('基础设施').click()
    await page.getByRole('link', { name: '网关集群' }).click()

    expect(await (await page.getByRole('button', { name: '新建集群' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.eo-table-btn-td button >> nth = 1').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 0').click()

    expect(await (await page.getByRole('button', { name: '新建配置' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '发布' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '同步配置' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '更改历史' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('false')

    expect(await (await page.getByRole('button', { name: '发布历史' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('false')

    const editDescIconBtn = await page.locator('cluster-desc-content iconfont.icon-bianji')
    expect(editDescIconBtn.isVisible()).toStrictEqual(false)

    const editIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-bianji') }).first()
    if (editIconBtn && await editIconBtn.isVisible()) {
      expect(await editIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const deleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (deleteIconBtn && await deleteIconBtn.isVisible()) {
      expect(await deleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 0').click()

    expect(await (await page.locator('.ant-drawer-body input')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body textarea >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body textarea >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body ').getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.locator('.ant-drawer-body').getByRole('button', { name: '取消' }).click()

    await page.getByRole('link', { name: '证书管理' }).click()

    expect(await (await page.getByRole('button', { name: '新建证书' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const editCertIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-bianji') }).first()
    if (editCertIconBtn && await editCertIconBtn.isVisible()) {
      expect(await editCertIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const deleteCertIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (deleteCertIconBtn && await deleteCertIconBtn.isVisible()) {
      expect(await deleteCertIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 0').click()

    expect(await (await page.locator('label').filter({ hasText: '上传密钥' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('cursor'))).toStrictEqual('not-allowed')

    expect(await (await page.locator('label').filter({ hasText: '上传密钥' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))).toStrictEqual('rgb(245, 245, 245)')

    expect(await (await page.locator('label').filter({ hasText: '上传密钥' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))).toStrictEqual('rgba(0, 0, 0, 0.25)')

    expect(await (await page.locator('label').filter({ hasText: '上传密钥' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))).toStrictEqual('rgb(217, 217, 217)')

    expect(await (await page.locator('label').filter({ hasText: '上传证书' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('cursor'))).toStrictEqual('not-allowed')

    expect(await (await page.locator('label').filter({ hasText: '上传证书' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))).toStrictEqual('rgb(245, 245, 245)')

    expect(await (await page.locator('label').filter({ hasText: '上传证书' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))).toStrictEqual('rgba(0, 0, 0, 0.25)')

    expect(await (await page.locator('label').filter({ hasText: '上传证书' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))).toStrictEqual('rgb(217, 217, 217)')

    expect(await (await page.locator('.ant-drawer-body textarea >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body textarea >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('.ant-drawer-body ').getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
    await page.getByRole('link', { name: '网关节点' }).click()

    expect(await (await page.getByRole('button', { name: '更新配置' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '重置配置' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('false')

    await page.getByRole('link', { name: '配置管理' }).click()

    expect(await (await page.locator('input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-switch')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const testIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-kuaisuceshi-1') }).first()
    if (testIconBtn && await testIconBtn.isVisible()) {
      expect(await testIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }
  })

  it('无权限角色浏览基础设施-环境变量页面 ', async () => {
    await page.locator('eo-ng-menu-default').getByRole('link', { name: '环境变量' }).click()

    expect(await (await page.getByRole('button', { name: '新建配置' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const deleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (deleteIconBtn && await deleteIconBtn.isVisible()) {
      expect(await deleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }
  })

  it('无权限角色浏览服务治理-流量策略页面', async () => {
    await page.locator('eo-ng-menu-default div').filter({ hasText: '服务治理' }).click()
    await page.getByRole('link', { name: '流量策略' }).click()

    expect(await (await page.getByRole('button', { name: '新建策略' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-table input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-table eo-ng-switch >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const deleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (deleteIconBtn && await deleteIconBtn.isVisible()) {
      expect(await deleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.locator('eo-ng-apinto-table').getByRole('cell', { name: 'testForE2e' }).first()

    expect(await (await page.locator('eo-ng-select >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-select >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-select >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('#desc')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 5')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 6')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 7')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 8')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 9')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 10')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 11')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 12')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 13')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 14')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const filterEditIconBtn = await page.locator('eo-ng-filter-table button').filter({ has: page.locator('iconfont.icon-a-peizhianniu_huaban1') }).first()
    if (filterEditIconBtn && await filterEditIconBtn.isVisible()) {
      expect(await filterEditIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const filterDeleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (filterDeleteIconBtn && await filterDeleteIconBtn.isVisible()) {
      expect(await filterDeleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    expect(await (await page.getByRole('button', { name: '添加条件' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
  })

  it('无权限角色浏览服务治理-熔断策略页面', async () => {
    await page.getByRole('link', { name: '熔断策略' }).click()

    expect(await (await page.getByRole('button', { name: '新建策略' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-table input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-table eo-ng-switch >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const deleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (deleteIconBtn && await deleteIconBtn.isVisible()) {
      expect(await deleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.locator('eo-ng-apinto-table').getByRole('cell', { name: 'testForE2e' }).first()

    expect(await (await page.locator('eo-ng-select >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-select >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-select >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-select >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-select >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('#desc')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 5')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 6')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 7')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 8')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 9')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 10')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 11')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 12')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 13')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 14')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const filterEditIconBtn = await page.locator('eo-ng-filter-table button').filter({ has: page.locator('iconfont.icon-a-peizhianniu_huaban1') }).first()
    if (filterEditIconBtn && await filterEditIconBtn.isVisible()) {
      expect(await filterEditIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const filterDeleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (filterDeleteIconBtn && await filterDeleteIconBtn.isVisible()) {
      expect(await filterDeleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    expect(await (await page.getByRole('button', { name: '添加条件' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
  })

  it('无权限角色浏览服务治理-访问策略页面', async () => {
    await page.getByRole('link', { name: '访问策略' }).click()

    expect(await (await page.getByRole('button', { name: '新建策略' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-table input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-table eo-ng-switch >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const deleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (deleteIconBtn && await deleteIconBtn.isVisible()) {
      expect(await deleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.locator('eo-ng-apinto-table').getByRole('cell', { name: 'testForE2e' }).first()

    expect(await (await page.locator('eo-ng-select >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('#desc')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-switch')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const filterEditIconBtn = await page.locator('eo-ng-filter-table button').filter({ has: page.locator('iconfont.icon-a-peizhianniu_huaban1') }).first()
    if (filterEditIconBtn && await filterEditIconBtn.isVisible()) {
      expect(await filterEditIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const filterDeleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (filterDeleteIconBtn && await filterDeleteIconBtn.isVisible()) {
      expect(await filterDeleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    expect(await (await page.getByRole('button', { name: '添加条件' }).first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '添加条件' }).last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
  })

  it('无权限角色浏览服务治理-缓存策略页面', async () => {
    await page.getByRole('link', { name: '缓存策略' }).click()

    expect(await (await page.getByRole('button', { name: '新建策略' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-table input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-table eo-ng-switch >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const deleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (deleteIconBtn && await deleteIconBtn.isVisible()) {
      expect(await deleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.locator('eo-ng-apinto-table').getByRole('cell', { name: 'testForE2e' }).first()

    expect(await (await page.locator('#desc')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const filterEditIconBtn = await page.locator('eo-ng-filter-table button').filter({ has: page.locator('iconfont.icon-a-peizhianniu_huaban1') }).first()
    if (filterEditIconBtn && await filterEditIconBtn.isVisible()) {
      expect(await filterEditIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const filterDeleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (filterDeleteIconBtn && await filterDeleteIconBtn.isVisible()) {
      expect(await filterDeleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    expect(await (await page.getByRole('button', { name: '添加条件' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
  })

  it('无权限角色浏览服务治理-灰度策略页面', async () => {
    await page.getByRole('link', { name: '灰度策略' }).click()

    expect(await (await page.getByRole('button', { name: '新建策略' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-table input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-table eo-ng-switch >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const deleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (deleteIconBtn && await deleteIconBtn.isVisible()) {
      expect(await deleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.locator('eo-ng-apinto-table').getByRole('cell', { name: 'testForE2e' }).first()

    expect(await (await page.locator('#desc')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 3')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 4')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input >> nth = 5')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    if (await page.locator('input >> nth = 6') && await page.locator('input >> nth = 6').isVisible()) {
      expect(await (await page.locator('input >> nth = 6')).evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    if (await page.locator('input >> nth = 7') && await page.locator('input >> nth = 7').isVisible()) {
      expect(await (await page.locator('input >> nth = 7')).evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const filterEditIconBtn = await page.locator('eo-ng-filter-table button').filter({ has: page.locator('iconfont.icon-a-peizhianniu_huaban1') }).first()
    if (filterEditIconBtn && await filterEditIconBtn.isVisible()) {
      expect(await filterEditIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const filterDeleteIconBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (filterDeleteIconBtn && await filterDeleteIconBtn.isVisible()) {
      expect(await filterDeleteIconBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    expect(await (await page.getByRole('button', { name: '添加条件' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    if (await page.getByRole('button', { name: '添加配置' }) && await page.getByRole('button', { name: '添加配置' }).isVisible()) {
      expect(await (await page.getByRole('button', { name: '添加配置' })).evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

      const editConfigBtn = await page.locator('.limit-bg button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
      if (editConfigBtn && await editConfigBtn.isVisible()) {
        expect(await editConfigBtn.evaluate((element) =>
          window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

        const deleteConfigIconBtn = await page.locator('.limit-bg button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
        expect(await deleteConfigIconBtn.evaluate((element) =>
          window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

        await page.locator('.limit-bg eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 0').click()

        await page.getByText('配置路由规则')

        expect(await (await page.locator('eo-ng-match-form input.ant-input >> nth = 0')).evaluate((element) =>
          window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

        expect(await (await page.locator('eo-ng-match-form input.ant-input >> nth = 1')).evaluate((element) =>
          window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

        expect(await (await page.locator('eo-ng-match-form eo-ng-select >> nth = 0')).evaluate((element) =>
          window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

        expect(await (await page.locator('eo-ng-match-form eo-ng-select >> nth = 1')).evaluate((element) =>
          window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

        expect(await (await page.locator('eo-ng-match-form').getByRole('button', { name: '提交' })).evaluate((element) =>
          window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

        await page.locator('eo-ng-match-form').getByRole('button', { name: '提交' }).click()
      }
    } else {
      expect(await (await page.locator('nz-slider >> nth = 0 ')).evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
      expect(await (await page.locator('nz-slider >> nth = 1 ')).evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    expect(await (await page.getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
  })

  it('无权限角色浏览系统管理-用户角色页面', async () => {
    await page.locator('eo-ng-menu-default div').filter({ hasText: '系统管理' }).click()
    await page.getByRole('link', { name: '用户角色' }).click()

    expect(await (await page.$$('.ant-checkbox')).length).toStrictEqual(0)

    expect(await (await page.getByRole('button', { name: '新建用户' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '删除' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-switch').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const pswResetBtn = await page.locator('eo-ng-filter-table button').filter({ has: page.locator('iconfont.icon-zhongzhiicon') }).first()
    if (pswResetBtn && await pswResetBtn.isVisible()) {
      expect(await pswResetBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const editBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-bianji') }).first()
    if (editBtn && await editBtn.isVisible()) {
      expect(await editBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const deleteBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (deleteBtn && await deleteBtn.isVisible()) {
      expect(await deleteBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.locator('eo-ng-apinto-table tr ').last().locator('td').first().click()

    expect(await (await page.locator('eo-ng-apinto-user-profile input.ant-input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-user-profile input.ant-input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-user-profile input.ant-input >> nth = 2')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-user-profile eo-ng-select')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-user-profile textarea')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-user-profile').getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.locator('eo-ng-apinto-user-profile').getByRole('button', { name: '取消' })

    await page.locator('eo-ng-tree-default-node').last().locator('.ant-tree-node-content-wrapper > .custom-node > .f-row-js-ac > .ant-btn').click()

    expect(await (await page.locator('.ant-dropdown-menu').getByText('删除')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.locator('.ant-dropdown-menu').getByText('编辑').click()
    await page.getByText('编辑角色').click()

    expect(await (await page.locator('eo-ng-apinto-role-profile input.ant-input')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-role-profile textarea')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-role-profile eo-ng-apinto-table .ant-checkbox').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-role-profile eo-ng-apinto-table .ant-checkbox').last()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-role-profile').getByText('提交')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
  })
  it('无权限角色浏览系统管理-外部应用页面', async () => {
    await page.locator('eo-ng-menu-default div').filter({ hasText: '系统管理' }).click()
    await page.getByRole('link', { name: '外部应用' }).click()

    expect(await (await page.getByRole('button', { name: '新建应用' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('eo-ng-apinto-table eo-ng-switch').first()).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    const tokenResetBtn = await page.locator('eo-ng-filter-table button').filter({ has: page.locator('iconfont.icon-shuaxinjiankongzhuangtai') }).first()
    if (tokenResetBtn && await tokenResetBtn.isVisible()) {
      expect(await tokenResetBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const editBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-chakan') }).first()
    if (editBtn && await editBtn.isVisible()) {
      expect(await editBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    const deleteBtn = await page.locator('button').filter({ has: page.locator('iconfont.icon-shanchu') }).first()
    if (deleteBtn && await deleteBtn.isVisible()) {
      expect(await deleteBtn.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')
    }

    await page.locator('eo-ng-apinto-table tr').last().locator('td').first().click()

    expect(await (await page.locator('input.ant-input >> nth = 0')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('input.ant-input >> nth = 1')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.locator('textarea')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    expect(await (await page.getByRole('button', { name: '提交' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('true')

    await page.getByRole('button', { name: '取消' }).click()
  })
  it('无权限角色浏览商业授权页面', async () => {
    await page.locator('.auth-box  a').click()

    expect(await (await page.locator('.ant-modal-body a')).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('visibility'))).toStrictEqual('hidden')

    await page.getByRole('button', { name: 'Close' }).click()
  })

  it('登录无权限用户，出现弹窗', async () => {
    await page.getByRole('button', { name: ' test2' }).click()
    await page.getByText('退出登录').click()
    await page.getByRole('heading', { name: '欢迎来到 Apinto' }).click()
    await page.getByPlaceholder('请输入账号').click()
    await page.getByPlaceholder('请输入账号').fill('testNoAccess')
    await page.getByPlaceholder('请输入密码').click()
    await page.getByPlaceholder('请输入密码').fill('12345678')
    await page.getByRole('button', { name: '登录' }).click()
    await page.getByText('登录成功').click()
    await page.getByText('无法获取您当前账号的相关权限信息，请确认是否赋予权限。').click()
    await page.getByRole('button', { name: '确定' }).click()
    await page.reload()
    await page.getByText('无法获取您当前账号的相关权限信息，请确认是否赋予权限。').click()
    await page.getByRole('button', { name: '取消' }).click()
    await page.reload()
    await page.getByText('无法获取您当前账号的相关权限信息，请确认是否赋予权限。').click()
    await page.getByRole('button', { name: 'Close' }).click()
  })

  it('登录正常权限用户，首页正常', async () => {
    await page.getByRole('button', { name: ' testNoAccess' }).click()
    await page.getByPlaceholder('请输入账号').click()
    await page.getByPlaceholder('请输入账号').fill('maggie')
    await page.getByPlaceholder('请输入账号').press('Tab')
    await page.getByPlaceholder('请输入密码').fill('12345678')
    await page.getByPlaceholder('请输入密码').press('Enter')

    await page.locator('eo-ng-menu-default').getByRole('link', { name: 'API管理' }).click()

    expect(await (await page.getByRole('button', { name: '新建API' })).evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('disabled'))).toStrictEqual('false')
  })
})
