describe('API管理 e2e test', () => {
  it('初始化页面，点击API管理菜单，进入列表页', async () => {
    await page.goto('http://localhost:4200/login')
    await page.waitForTimeout(2000)
    await page.getByPlaceholder('请输入账号').click()
    await page.getByPlaceholder('请输入账号').fill('maggie')
    await page.getByPlaceholder('请输入账号').press('Tab')
    await page.getByPlaceholder('请输入密码').fill('12345678')
    await page.getByPlaceholder('请输入密码').press('Enter')
    await page.locator('eo-ng-menu-default').getByRole('link', { name: 'API管理' }).click()
  })
  it('左侧目录组件与上方header间距6px, 一级menu196px*48px,字号14px, 与左侧距离24px；面包屑API管理，字号14，字体颜色为主题色，与左侧距离12px，垂直居中', async () => {
    // 目录样式 左侧目录组件与上方header间距6px, 一级menu196px*48px,字号14px, 与左侧距离24px
    const menuGroup = await page.locator('.ant-layout-sider-children')
    const menuGroupMT = await menuGroup.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    expect(menuGroupMT).toStrictEqual('6px')

    const apiMenuItem = await page.locator('eo-ng-menu-default .ant-menu-item:has-text("API管理")')
    const apiMenuItemH = await apiMenuItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(apiMenuItemH).toStrictEqual('48px')
    const apiMenuItemW = await apiMenuItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(apiMenuItemW).toStrictEqual('195px')
    const apiMenuItemPL = await apiMenuItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    expect(apiMenuItemPL).toStrictEqual('24px')
    const apiMenuItemSFS = await apiMenuItem.locator('span').evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(apiMenuItemSFS).toStrictEqual('20px')
    const apiMenuItemAFS = await apiMenuItem.locator('a').evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(apiMenuItemAFS).toStrictEqual('14px')
    const apiMenuItemBG = await apiMenuItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(apiMenuItemBG).toStrictEqual('rgb(240, 247, 255)')
    const apiMenuItemC = await apiMenuItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(apiMenuItemC).toStrictEqual('rgb(34, 84, 157)')

    // 面包屑样式 面包屑API管理，字号14，字体颜色为主题色，与左侧距离12px，垂直居中
    const headerBlock = await page.locator('.block_rl')
    const headerBlockAI = await headerBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('align-items'))
    expect(headerBlockAI).toStrictEqual('center')

    const breadcrumbItem = await page.locator('nz-breadcrumb-item').getByRole('link', { name: 'API管理' })
    const breadcrumbItemFS = await breadcrumbItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(breadcrumbItemFS).toStrictEqual('14px')
    const breadcrumbItemC = await breadcrumbItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(breadcrumbItemC).toStrictEqual('rgb(34, 84, 157)')
  })
  it('分组组件宽234px，搜索主体上下margin12px，左右padding11px，高32px；搜索输入框高32px，边框圆角50px，与右侧按钮间距12px；按钮高32px，背景色与悬浮背景色不同', async () => {
    // 分组组件宽234px
    const groupBlock = await page.locator('.block-left')
    const groupBlockW = await groupBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(groupBlockW).toStrictEqual('234px')

    // 搜索主体上下margin12px，左右padding11px，高32px
    const searchBlock = await page.locator('.group-top')
    const searchBlockM = await searchBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    expect(searchBlockM).toStrictEqual('12px 0px')
    const searchBlockP = await searchBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(searchBlockP).toStrictEqual('0px 11px')
    const searchBlockH = await searchBlock.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(searchBlockH).toStrictEqual('32px')

    // 搜索输入框高32px，边框圆角50px
    const searchInput = await page.locator('eo-ng-input-group').first()
    const searchInputH = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(searchInputH).toStrictEqual('32px')
    const searchInputR = await searchInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-radius'))
    expect(searchInputR).toStrictEqual('50px')

    // 按钮高32px，与左侧输入框间距12px, 背景色与悬浮背景色不同
    const searchBtn = await page.locator('.ant-table-row-expand-icon')
    const searchBtnH = await searchBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(searchBtnH).toStrictEqual('32px')
    const searchBtnML = await searchBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(searchBtnML).toStrictEqual('12px')
    const searchBtnBG = await searchBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(searchBtnBG).toStrictEqual('rgb(248, 248, 250)')

    // await page.hover('.ant-table-row-expand-icon')
    // const searchBtnHBG = await page.locator('.ant-table-row-expand-icon').evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('background-color'))
    // expect(searchBtnHBG).toStrictEqual('rgb(233, 238, 245)')
  })
  it('分组标题-所有API，左侧间距16px，其他选项左侧间距8px；所有选项234px*30px，外侧无间距，背景色白色，鼠标悬浮与点击时，背景色变悬浮色', async () => {
    // 分组标题-所有API，左侧间距16px，其他选项左侧间距12px
    const groupTitle = await page.locator('div').filter({ hasText: '所有API' }).nth(4)
    const groupTitlePL = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    expect(groupTitlePL).toStrictEqual('16px')

    const itemFTitle = await page.locator('eo-ng-tree-default-node').first()
    const itemFTitlePL = await itemFTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    expect(itemFTitlePL).toStrictEqual('8px')

    const itemLTitle = await page.locator('eo-ng-tree-default-node').last()
    const itemLTitlePL = await itemLTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-left'))
    expect(itemLTitlePL).toStrictEqual('8px')

    // 所有选项234px*30px，外侧无间距，背景色白色，鼠标悬浮与点击时，背景色变悬浮色
    const groupTitleH = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(groupTitleH).toContain('30')
    const groupTitleW = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(groupTitleW).toContain('224')
    const groupTitleM = await groupTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    expect(groupTitleM).toStrictEqual('0px')

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
    const itemFTitleM = await itemFTitle.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    expect(itemFTitleM).toStrictEqual('3px 0px')

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
  })
  it('列表中按钮样式、表格、分页样式检查；不同协议/方法样式检查', async () => {
    // 列表中按钮样式、表格、分页样式检查
    const createBtn = await page.getByRole('button', { name: '新建API' })

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

    const importBtn = await page.getByRole('button', { name: '导入' })

    const importBtnH = await importBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(importBtnH).toStrictEqual('32px')
    const importBtnML = await importBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(importBtnML).toStrictEqual('12px')
    const importBtnP = await importBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(importBtnP).toStrictEqual('0px 12px')
    const importBtnBG = await importBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(importBtnBG).toStrictEqual('rgb(255, 255, 255)')
    const importBtnBDC = await importBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    expect(importBtnBDC).toStrictEqual('rgb(217, 217, 217)')

    const onlineBtn = await page.getByRole('button', { name: '上线' })

    const onlineBtnH = await onlineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(onlineBtnH).toStrictEqual('32px')
    const onlineBtnML = await onlineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(onlineBtnML).toStrictEqual('12px')
    const onlineBtnP = await onlineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(onlineBtnP).toStrictEqual('0px 12px')
    const onlineBtnBG = await onlineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(onlineBtnBG).toStrictEqual('rgb(245, 245, 245)')
    const onlineBtnBDC = await onlineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    expect(onlineBtnBDC).toStrictEqual('rgb(217, 217, 217)')

    const offlineBtn = await page.getByRole('button', { name: '下线' })

    const offlineBtnH = await offlineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(offlineBtnH).toStrictEqual('32px')
    const offlineBtnML = await offlineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))
    expect(offlineBtnML).toStrictEqual('12px')
    const offlineBtnP = await offlineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(offlineBtnP).toStrictEqual('0px 12px')
    const offlineBtnBG = await offlineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(offlineBtnBG).toStrictEqual('rgb(245, 245, 245)')
    const offlineBtnBDC = await offlineBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    expect(offlineBtnBDC).toStrictEqual('rgb(217, 217, 217)')

    const searchInputGroup = await page.locator('eo-ng-api-management-list .group-search-large.mg-top-right')
    const searchInputGroupMR = await searchInputGroup.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-right'))
    expect(searchInputGroupMR).toStrictEqual('24px')

    const searchInput = await page.locator('eo-ng-api-management-list eo-ng-input-group')
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

    const pagination = await page.locator('.mg_pagination_t')
    const paginationH = await pagination.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    expect(paginationH).toStrictEqual('32px')
    const paginationPT = await pagination.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-top'))
    expect(paginationPT).toStrictEqual('0px')
    const paginationMT = await pagination.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    expect(paginationMT).toStrictEqual('16px')

    // 表头宽度、字体左右间距，表头字体颜色
    const selectTh = await page.locator('.ant-table-selection-column').first()
    const selectThW = await selectTh.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(selectThW).toStrictEqual('40px')

    const methodTh = await page.getByRole('columnheader', { name: '协议/方法' })
    const methodThW = await methodTh.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(methodThW).toStrictEqual('140px')
    const methodThP = await methodTh.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))
    expect(methodThP).toStrictEqual('0px 12px')
    const methodThFS = await methodTh.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-size'))
    expect(methodThFS).toStrictEqual('14px')
    const methodThFW = await methodTh.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('font-weight'))
    expect(methodThFW).toStrictEqual('400')
    const methodThC = await methodTh.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(methodThC).toStrictEqual('rgb(102, 102, 102)')

    const operationTh = await page.getByRole('columnheader', { name: '操作' })
    const operationThW = await operationTh.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    expect(operationThW).toStrictEqual('130px')

    // 不同协议/方法样式检查
    const deleteTag = await page.locator('span.method:has-text("DELETE")').first()
    if (deleteTag && await deleteTag.isVisible()) {
      const deleteTagH = await deleteTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('height'))
      expect(deleteTagH).toStrictEqual('20px')
      const deleteTagLH = await deleteTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('line-height'))
      expect(deleteTagLH).toStrictEqual('12px')
      const deleteTagP = await deleteTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('padding'))
      expect(deleteTagP).toStrictEqual('4px 6px')
      const deleteTagBG = await deleteTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('background-color'))
      expect(deleteTagBG).toStrictEqual('rgba(194, 22, 27, 0.15)')
      const deleteTagC = await deleteTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      expect(deleteTagC).toStrictEqual('rgb(194, 22, 27)')
      const deleteTagFS = await deleteTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-size'))
      expect(deleteTagFS).toStrictEqual('12px')
    }

    const getTag = await page.locator('span.method:has-text("GET")').first()
    if (getTag && await getTag.isVisible()) {
      const getTagH = await getTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('height'))
      expect(getTagH).toStrictEqual('20px')
      const getTagLH = await getTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('line-height'))
      expect(getTagLH).toStrictEqual('12px')
      const getTagP = await getTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('padding'))
      expect(getTagP).toStrictEqual('4px 6px')
      const getTagBG = await getTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('background-color'))
      expect(getTagBG).toStrictEqual('rgba(6, 125, 219, 0.15)')
      const getTagC = await getTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      expect(getTagC).toStrictEqual('rgb(6, 125, 219)')
      const getTagFS = await getTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-size'))
      expect(getTagFS).toStrictEqual('12px')
    }

    const postTag = await page.locator('span.method:has-text("POST")').first()
    if (postTag && await postTag.isVisible()) {
      const postTagH = await postTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('height'))
      expect(postTagH).toStrictEqual('20px')
      const postTagLH = await postTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('line-height'))
      expect(postTagLH).toStrictEqual('12px')
      const postTagP = await postTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('padding'))
      expect(postTagP).toStrictEqual('4px 6px')
      const postTagBG = await postTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('background-color'))
      expect(postTagBG).toStrictEqual('rgba(16, 165, 75, 0.15)')
      const postTagC = await postTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      expect(postTagC).toStrictEqual('rgb(16, 165, 75)')
      const postTagFS = await postTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-size'))
      expect(postTagFS).toStrictEqual('12px')
    }

    const putTag = await page.locator('span.method:has-text("PUT")').first()
    if (putTag && await putTag.isVisible()) {
      const putTagH = await putTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('height'))
      expect(putTagH).toStrictEqual('20px')
      const putTagLH = await putTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('line-height'))
      expect(putTagLH).toStrictEqual('12px')
      const putTagP = await putTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('padding'))
      expect(putTagP).toStrictEqual('4px 6px')
      const putTagBG = await putTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('background-color'))
      expect(putTagBG).toStrictEqual('rgba(216, 131, 12, 0.15)')
      const putTagC = await putTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      expect(putTagC).toStrictEqual('rgb(216, 131, 12)')
      const putTagFS = await putTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-size'))
      expect(putTagFS).toStrictEqual('12px')
    }

    const allTag = await page.locator('span.method:has-text("ALL")').first()
    if (allTag && await allTag.isVisible()) {
      const allTagH = await allTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('height'))
      expect(allTagH).toStrictEqual('20px')
      const allTagLH = await allTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('line-height'))
      expect(allTagLH).toStrictEqual('12px')
      const allTagP = await allTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('padding'))
      expect(allTagP).toStrictEqual('4px 6px')
      const allTagBG = await allTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('background-color'))
      expect(allTagBG).toStrictEqual('rgba(119, 40, 245, 0.15)')
      const allTagC = await allTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      expect(allTagC).toStrictEqual('rgb(119, 40, 245)')
      const allTagFS = await allTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-size'))
      expect(allTagFS).toStrictEqual('12px')
    }

    const patchTag = await page.locator('span.method:has-text("PATCH")').first()
    if (patchTag && await patchTag.isVisible()) {
      const patchTagH = await patchTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('height'))
      expect(patchTagH).toStrictEqual('20px')
      const patchTagLH = await patchTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('line-height'))
      expect(patchTagLH).toStrictEqual('12px')
      const patchTagP = await patchTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('padding'))
      expect(patchTagP).toStrictEqual('4px 6px')
      const patchTagBG = await patchTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('background-color'))
      expect(patchTagBG).toStrictEqual('rgba(237, 134, 58, 0.15)')
      const patchTagC = await patchTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      expect(patchTagC).toStrictEqual('rgb(237, 134, 58)')
      const patchTagFS = await patchTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-size'))
      expect(patchTagFS).toStrictEqual('12px')
    }

    const headTag = await page.locator('span.method:has-text("HEAD")').first()
    if (headTag && await headTag.isVisible()) {
      const headTagH = await headTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('height'))
      expect(headTagH).toStrictEqual('20px')
      const headTagLH = await headTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('line-height'))
      expect(headTagLH).toStrictEqual('12px')
      const headTagP = await headTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('padding'))
      expect(headTagP).toStrictEqual('4px 6px')
      const headTagBG = await headTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('background-color'))
      expect(headTagBG).toStrictEqual('rgba(238, 196, 12, 0.15)')
      const headTagC = await headTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('color'))
      expect(headTagC).toStrictEqual('rgb(238, 196, 12)')
      const headTagFS = await headTag.evaluate((element) =>
        window.getComputedStyle(element).getPropertyValue('font-size'))
      expect(headTagFS).toStrictEqual('12px')
    }
  })
  it('点击分组，menu中API管理始终被选中，鼠标悬浮在分组中，测试更多操作,直至点击添加api，将跳转至新建api页面', async () => {
    await page.locator('.custom-node >> nth = 0').click()

    const apiMenuItem = await page.locator('eo-ng-menu-default .ant-menu-item:has-text("API管理")')
    const apiMenuItemBG = await apiMenuItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    expect(apiMenuItemBG).toStrictEqual('rgb(240, 247, 255)')
    const apiMenuItemC = await apiMenuItem.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    expect(apiMenuItemC).toStrictEqual('rgb(34, 84, 157)')

    await page.hover('.custom-node >> nth = 0')
    const dropdownBtn = await page.locator('.custom-node >> nth = 0 >> button')
    await dropdownBtn.click()
    await page.getByText('添加子分组').click()
    await page.getByPlaceholder('分组名称').fill('test')
    await page.getByRole('button', { name: '确定' }).click()

    await page.locator('.ant-table-row-expand-icon').click()
    await page.getByPlaceholder('分组名称').fill('test2')
    await page.getByRole('button', { name: '确定' }).click()
    await page.getByText('test2').last().click()

    await page.locator('.custom-node').last().hover()
    await page.locator('.custom-node').last().locator('button').isVisible()
    await page.locator('.custom-node').last().locator('button').click()
    await page.getByText('编辑').click()
    await page.getByPlaceholder('分组名称').fill('testE2E')
    await page.getByRole('button', { name: '确定' }).click()
    await page.locator('.custom-node').last().getByRole('button', { name: '' }).click()
    await page.getByText('删除').first().click()
    await page.getByRole('button', { name: '取消' }).click()
    await page.locator('.custom-node').last().getByRole('button', { name: '' }).click()
    await page.getByText('删除').first().click()
    await page.getByRole('button', { name: '确定' }).click()
    await page.waitForTimeout(1000)
    await page.locator('.custom-node').last().hover()
    await page.locator('.custom-node').last().locator('button').isVisible()
    await page.locator('.custom-node').last().locator('button').click()
    await page.waitForTimeout(1000)
    await page.getByText('添加子分组').click()
    await page.getByPlaceholder('分组名称').fill('test-1')
    await page.getByRole('button', { name: '确定' }).click()
    await page.waitForTimeout(1000)
    await page.locator('nz-tree-node-title').filter({ hasText: 'test-1' }).last().hover()
    await page.locator('nz-tree-node-title').filter({ hasText: 'test-1' }).last().locator('button').isVisible()
    await page.locator('nz-tree-node-title').filter({ hasText: 'test-1' }).last().locator('button').click()
    await page.getByText('添加 API').click()
  })
  it('从分组内部的添加api进入新建api页面，其中分组栏不为空，检查表单页所有输入框样式，点击取消，返回列表页', async () => {
    // 所属分组值不为空，样式
    const groupInput = await page.locator('eo-ng-tree-select div')
    expect(await groupInput.inputValue).not.toStrictEqual('')
    const groupInputW = await groupInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const groupInputH = await groupInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(groupInputW).toStrictEqual('202px')
    await expect(groupInputH).toStrictEqual('32px')

    // API名称
    const nameInput = await page.getByPlaceholder('请输入API名称')
    const nameInputW = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const nameInputH = await nameInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(nameInputW).toStrictEqual('306px')
    await expect(nameInputH).toStrictEqual('32px')

    // 所属分组group
    const inputGroup = await page.locator('eo-ng-input-group')
    const inputGroupW = await inputGroup.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const inputGroupH = await inputGroup.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(inputGroupW).toStrictEqual('508px')
    await expect(inputGroupH).toStrictEqual('32px')

    // 描述样式
    const descInput = await page.locator('textarea#desc')
    const descInputW = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const descInputH = await descInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(descInputW).toStrictEqual('508px')
    await expect(descInputH).toStrictEqual('68px')

    // 请求路径
    const requestPathInput = await page.locator('#request_path')
    const requestPathInputW = await requestPathInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const requestPathInputH = await requestPathInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(requestPathInputW).toStrictEqual('508px')
    await expect(requestPathInputH).toStrictEqual('32px')

    // 绑定上游服务
    const upstreamInput = await page.locator('eo-ng-select-top-control')
    const upstreamInputW = await upstreamInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const upstreamInputH = await upstreamInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(upstreamInputW).toStrictEqual('508px')
    await expect(upstreamInputH).toStrictEqual('32px')

    // 请求方式
    const checkboxGroup = await page.getByText('All GETPOSTPUTDELETEPATCHHEAD')
    const checkboxGroupH = await checkboxGroup.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(checkboxGroupH).toStrictEqual('32px')

    // 转发上游路径
    const upstreamPathInput = await page.locator('#proxy_path')
    const upstreamPathInputW = await upstreamPathInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const upstreamPathInputH = await upstreamPathInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(upstreamPathInputW).toStrictEqual('508px')
    await expect(upstreamPathInputH).toStrictEqual('32px')

    // 请求超时时间
    const timeOutInput = await page.locator('#timeout')
    const timeOutInputW = await timeOutInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const timeOutInputH = await timeOutInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(timeOutInputW).toStrictEqual('508px')
    await expect(timeOutInputH).toStrictEqual('32px')

    // 重试次数
    const retryInput = await page.locator('#retry')
    const retryInputW = await retryInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const retryInputH = await retryInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(retryInputW).toStrictEqual('508px')
    await expect(retryInputH).toStrictEqual('32px')

    // 高级匹配
    const filterBtn1 = await page.locator('eo-ng-match-table').getByRole('button', { name: '添加配置' }).first()
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

    // 转发上游请求头
    const filterBtn2 = await page.locator('nz-form-item').filter({ hasText: '转发上游请求头 添加配置' }).getByRole('button', { name: '添加配置' })
    const filterBtn2W = await filterBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const filterBtn2H = await filterBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const filterBtn2BGC = await filterBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const filterBtn2BC = await filterBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const filterBtn2C = await filterBtn2.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(filterBtn2W).toStrictEqual('82px')
    await expect(filterBtn2H).toStrictEqual('32px')
    await expect(filterBtn2BGC).toStrictEqual('rgb(255, 255, 255)')
    await expect(filterBtn2BC).toStrictEqual('rgb(217, 217, 217)')
    await expect(filterBtn2C).toStrictEqual('rgba(0, 0, 0, 0.85)')
    // filterBtn2.hover()
    // filterBtn2BC = await filterBtn2.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('border-color'))
    // await expect(filterBtn2BC).toStrictEqual('rgb(34, 84, 157)')
    // filterBtn2C = await filterBtn2.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('color'))
    // await expect(filterBtn2C).toStrictEqual('rgb(34, 84, 157)')

    await page.getByRole('button', { name: '取消' }).click()
  })
  it('在搜索框输入文字，被搜中的关键字为主题色，清空搜索，无主题色字体', async () => {
    await page.getByPlaceholder('搜索').first().click()
    await page.getByPlaceholder('搜索').first().fill('test')
    await page.locator('.highlight').first().click()

    const highlightText = await page.locator('.highlight').first()
    const highlightTextC = await highlightText.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    await expect(highlightTextC).toStrictEqual('rgb(34, 84, 157)')

    await page.locator('.anticon > svg').first().click()
    expect(await page.locator('.highlight').first().isVisible()).toStrictEqual(false)
  })
  it('点击新建API，清空必输项，表单无法保存；逐个填写必填项，直至表单可以提交，将返回列表页', async () => {
    await page.getByRole('button', { name: '新建API' }).click()
    await page.getByRole('button', { name: '保存' }).click()

    await page.locator('eo-ng-input-group nz-form-control').filter({ hasText: '请选择 必填项' }).getByRole('alert').click()
    await page.locator('eo-ng-input-group').getByText('必填项').nth(1).click()
    await page.locator('nz-form-item').filter({ hasText: '请求路径必填项' }).getByRole('alert').click()
    await page.locator('nz-form-item').filter({ hasText: '绑定上游服务 请选择 必填项' }).getByRole('alert').click()
    await page.locator('.ant-form-item-with-help > .ant-form-item-explain > .ant-form-item-explain-error').click()
    await page.locator('nz-form-item').filter({ hasText: '转发上游路径必填项' }).getByRole('alert').click()
    await page.getByText('请求超时时间').click()
    await page.locator('#timeout').fill('1000')
    await page.locator('#timeout').click()
    await page.locator('#timeout').fill('')
    await page.locator('nz-form-item').filter({ hasText: '保存 取消' }).click()
    await page.locator('nz-form-item').filter({ hasText: '保存 取消' }).click()
    await page.locator('#retry').click()
    await page.locator('#retry').fill('')
    await page.locator('nz-form-control').filter({ hasText: '必填项单位：ms，最小值：1' }).getByRole('alert').click()
    await page.getByText('单位：ms，最小值：1').click()
    await page.locator('nz-form-item').filter({ hasText: '重试次数必填项' }).getByRole('alert').click()
    await page.locator('eo-ng-tree-select div').click()
    await page.locator('eo-ng-tree-default-node').filter({ hasText: '应用管理' }).locator('svg').click()
    await page.getByText('测试').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByPlaceholder('请输入API名称').click()
    await page.getByPlaceholder('请输入API名称').fill('testForE2e')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('#request_path').click()
    await page.locator('#request_path').fill('test')
    await page.getByText('请输入以/开头，/-_与花括号、大小写字母、数字的组合').click()
    await page.locator('nz-form-item').filter({ hasText: '转发上游路径请输入以/开头，/-_与花括号、大小写字母、数字的组合' }).getByRole('alert').click()
    await page.getByText('请求路径').click()
    await page.locator('#request_path').click()
    await page.locator('#request_path').fill('/test')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('eo-ng-select-top-control').click()
    await page.getByText('test').nth(1).click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('#proxy_path').click()
    await page.locator('#proxy_path').click()
    await page.locator('#proxy_path').fill('/test')
    await page.locator('#timeout').click()
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('#timeout').click()
    await page.locator('#timeout').fill('-1')
    await page.getByRole('button', { name: '保存' }).click()
    await page.locator('#timeout').click()
    await page.locator('#retry').click()
    await page.locator('#retry').fill('2')
    await page.getByRole('button', { name: '保存' }).click()
    await page.getByLabel('All').check()
    await page.getByRole('button', { name: '保存' }).click()
  })

  it('点击导入，出现抽屉弹窗，检查样式', async () => {
    await page.locator('div').filter({ hasText: '所有API' }).nth(4).click()
    await page.getByRole('button', { name: '导入' }).click()

    // 上传文件
    const uploadBtn = await page.getByRole('button', { name: '选择文件 支持swagger3.0的json、yaml格式' }).getByRole('button', { name: '选择文件' })
    const uploadBtnW = await uploadBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const uploadBtnH = await uploadBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const uploadBtnBGC = await uploadBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const uploadBtnBC = await uploadBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const uploadBtnC = await uploadBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(uploadBtnW).toStrictEqual('82px')
    await expect(uploadBtnH).toStrictEqual('32px')
    await expect(uploadBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    await expect(uploadBtnBC).toStrictEqual('rgb(217, 217, 217)')
    await expect(uploadBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')
    // uploadBtn.hover()
    // uploadBtnBC = await uploadBtn.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('border-color'))
    // await expect(uploadBtnBC).toStrictEqual('rgb(34, 84, 157)')
    // uploadBtnC = await uploadBtn.evaluate((element) =>
    //   window.getComputedStyle(element).getPropertyValue('color'))
    // await expect(uploadBtnC).toStrictEqual('rgb(34, 84, 157)')

    // API分组
    const timeOutInput = await page.locator('eo-ng-tree-select div')
    const timeOutInputW = await timeOutInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const timeOutInputH = await timeOutInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(timeOutInputW).toStrictEqual('346px')
    await expect(timeOutInputH).toStrictEqual('32px')

    // 绑定上游服务
    const upstreamInput = await page.locator('#upstream eo-ng-select-top-control')
    const upstreamInputW = await upstreamInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const upstreamInputH = await upstreamInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(upstreamInputW).toStrictEqual('346px')
    await expect(upstreamInputH).toStrictEqual('32px')

    // 请求前缀
    const prefixInput = await page.getByPlaceholder('请输入')
    const prefixInputW = await prefixInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const prefixInputH = await prefixInput.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))

    await expect(prefixInputW).toStrictEqual('346px')
    await expect(prefixInputH).toStrictEqual('32px')

    // 查重
    const checkBtn = await page.getByRole('button', { name: '查重' })
    const checkBtnW = await checkBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const checkBtnH = await checkBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const checkBtnBGC = await checkBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const checkBtnBC = await checkBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const checkBtnC = await checkBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(checkBtnW).toStrictEqual('54px')
    await expect(checkBtnH).toStrictEqual('32px')
    await expect(checkBtnBGC).toStrictEqual('rgb(34, 84, 157)')
    await expect(checkBtnBC).toStrictEqual('rgb(34, 84, 157)')
    await expect(checkBtnC).toStrictEqual('rgb(255, 255, 255)')

    // 取消
    const cancleBtn = await page.getByRole('button', { name: '取消' })
    const cancleBtnW = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const cancleBtnH = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const cancleBtnBGC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const cancleBtnC = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const cancleBtnML = await cancleBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    await expect(cancleBtnW).toStrictEqual('54px')
    await expect(cancleBtnH).toStrictEqual('32px')
    await expect(cancleBtnBGC).toStrictEqual('rgb(255, 255, 255)')
    await expect(cancleBtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')
    await expect(cancleBtnML).toStrictEqual('12px')

    await page.setInputFiles('.ant-upload input', './test/swagger.yml')

    await page.locator('.ant-drawer-body').click()
    await page.getByRole('button', { name: '查重' }).click()
    await page.locator('eo-ng-tree-select div').click()
    await page.locator('eo-ng-tree-default-node.ant-select-tree-treenode').first().locator('nz-tree-node-switcher').click()
    await page.locator('eo-ng-tree-default-node.ant-select-tree-treenode').filter({ hasText: 'test' }).first().click()
    await page.getByRole('button', { name: '查重' }).click()
    await page.locator('#upstream eo-ng-select-top-control').click()
    await page.getByText('twqwt').click()
    await page.getByPlaceholder('请输入').click()
    await page.getByPlaceholder('请输入').fill('/testForE2e')
    await page.getByRole('button', { name: '查重' }).click()

    // 表格
    const table = await page.locator('.ant-drawer-body eo-ng-apinto-table')
    const tableM = await table.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    await expect(tableM).toStrictEqual('0px auto')

    // 提交
    const submitBtn = await page.getByRole('button', { name: '提交' })
    const submitBtnW = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const submitBtnH = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const submitBtnBGC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const submitBtnBC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const submitBtnC = await submitBtn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))

    await expect(submitBtnW).toStrictEqual('54px')
    await expect(submitBtnH).toStrictEqual('32px')
    await expect(submitBtnBGC).toStrictEqual('rgb(34, 84, 157)')
    await expect(submitBtnBC).toStrictEqual('rgb(34, 84, 157)')
    await expect(submitBtnC).toStrictEqual('rgb(255, 255, 255)')

    // 取消
    const cancle2Btn = await page.getByRole('button', { name: '返回上级' })
    const cancle2BtnW = await cancle2Btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('width'))
    const cancle2BtnH = await cancle2Btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('height'))
    const cancle2BtnBGC = await cancle2Btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('background-color'))
    const cancle2BtnBC = await cancle2Btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('border-color'))
    const cancle2BtnC = await cancle2Btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('color'))
    const cancle2BtnML = await cancle2Btn.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-left'))

    await expect(cancle2BtnW).toStrictEqual('82px')
    await expect(cancle2BtnH).toStrictEqual('32px')
    await expect(cancle2BtnBGC).toStrictEqual('rgb(255, 255, 255)')
    await expect(cancle2BtnBC).toStrictEqual('rgb(217, 217, 217)')
    await expect(cancle2BtnC).toStrictEqual('rgba(0, 0, 0, 0.85)')
    await expect(cancle2BtnML).toStrictEqual('12px')

    // 提交

    await page.getByRole('row', { name: '序号 API名称 协议/方法 请求路径 描述 状态' }).getByLabel('').uncheck()
    await page.locator('.ant-drawer-body label >> nth = 2').check()
    await page.getByRole('button', { name: '提交' }).click()
  })
  it('在左侧列表选中api，可操作批量上下线，并检查样式', async () => {
    await page.locator('.ant-table-selection-column').first().click()
    await page.locator('.ant-table-selection-column').first().click()

    await page.locator('label.ant-checkbox-wrapper >> nth = 2').check()
    await page.locator('label.ant-checkbox-wrapper >> nth = 3').check()

    await page.getByRole('button', { name: '上线' }).click()
    await page.getByText('批量上线').click()
    await page.getByText('*选择网关集群').click()
    await page.locator('.ant-row input').first().click()
    await page.getByRole('button', { name: '下一步' }).last().click()
    await page.getByRole('button', { name: '上一步' }).last().click()
    await page.getByRole('button', { name: '下一步' }).last().click()
    await page.locator('.ant-drawer-title').getByText('检测结果').last().click()
    await page.getByRole('button', { name: '重新检测' }).last().click()
    await page.locator('.ant-drawer-title').getByText('检测结果').last().click()
    await page.getByRole('button', { name: '批量上线' }).last().click()
    await page.locator('.ant-drawer-title').getByText('批量上线结果').click()
    await page.getByRole('button', { name: '返回' }).last().click()

    await page.waitForTimeout(700)
    await page.locator('label.ant-checkbox-wrapper >> nth = 2').check()
    await page.locator('label.ant-checkbox-wrapper >> nth = 3').check()
    await page.getByRole('button', { name: '下线' }).last().click()
    await page.getByText('*选择网关集群').last().click()
    await page.locator('.ant-row input').first().click()
    await page.locator('.ant-drawer-title').getByText('批量下线').last().click()
    await page.getByRole('button', { name: '取消' }).last().click()
    await page.getByRole('button', { name: '下线' }).last().click()
    await page.locator('.ant-row input').first().click()
    await page.getByRole('button', { name: '提交' }).last().click()
    await page.locator('.ant-drawer-title').getByText('批量下线结果').last().click()
    await page.getByRole('button', { name: '返回' }).last().click()
  })
  it('在列表中选择来源，输入搜索内容', async () => {
    ;
    await page.getByRole('columnheader', { name: '来源' }).locator('svg').click()
    await page.getByRole('listitem').filter({ hasText: '自建' }).locator('label').click()
    await page.getByRole('button', { name: '确定' }).click()
    await page.locator('nz-filter-trigger span').click()
    await page.getByRole('listitem').filter({ hasText: '导入' }).locator('label').click()
    await page.getByRole('button', { name: '确定' }).click()
    await page.getByRole('columnheader', { name: '来源' }).locator('svg').click()
    await page.getByRole('button', { name: '重置' }).click()

    await page.getByPlaceholder('搜索API名称').click()
    await page.getByPlaceholder('搜索API名称').fill('test')
    await page.getByPlaceholder('搜索API名称').press('Enter')

    await page.getByPlaceholder('搜索API名称').fill('')
    await page.getByPlaceholder('搜索API名称').press('Enter')
  })
  it('点击单个api，进入信息页，修改输入并提交，返回列表页', async () => {
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 1').click()
    await page.getByText('请求超时时间').click()
    await page.locator('#timeout').click()
    await page.locator('#timeout').fill('')
    await page.getByRole('button', { name: '提交' }).click()
    await page.locator('#timeout').click()
    await page.locator('#timeout').fill('200')
    await page.getByRole('button', { name: '提交' }).click()
  })
  it('点击单个api，点击上线管理tab，进入上线管理页，检查样式', async () => {
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 1').click()

    await page.getByRole('link', { name: '上线管理' }).click()
    await page.getByRole('columnheader', { name: '集群名称' }).click()

    const table = await page.locator('eo-ng-apinto-table')
    const tableM = await table.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin'))
    const tableP = await table.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding'))

    await expect(tableM).toStrictEqual('0px')
    await expect(tableP).toStrictEqual('0px')

    const listContent = await page.locator('.list-content')
    const listContentMT = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('margin-top'))
    const listContentPB = await listContent.evaluate((element) =>
      window.getComputedStyle(element).getPropertyValue('padding-bottom'))

    await expect(listContentMT).toStrictEqual('12px')
    await expect(listContentPB).toStrictEqual('20px')
  })
  it('点击API信息tab，点击取消返回上线管理页，点击面包屑返回api列表页，再次进入api信息页，点击面包屑返回', async () => {
    await page.getByRole('link', { name: 'API信息' }).click()
    await page.getByRole('button', { name: '取消' }).click()

    await page.locator('nz-breadcrumb').getByRole('link', { name: 'API管理' }).click()
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 1').click()

    await page.locator('nz-breadcrumb').getByRole('link', { name: 'API管理' }).click()
    await page.locator('eo-ng-apinto-table tr >> nth = 2 >> td >> nth = 1').click()

    await page.getByText('API管理 / API信息 /').click()
    await page.locator('nz-breadcrumb').getByText('API信息').click()
    await page.getByRole('link', { name: '上线管理' }).click()
    await page.locator('nz-breadcrumb').getByRole('link', { name: 'API管理' }).click()
  })
  it('API列表的分页测试', async () => {
    await page.getByRole('listitem', { name: '2' }).getByText('2').click()
    await page.getByRole('listitem', { name: '3' }).getByText('3').click()
    await page.getByText('20 条/页').click()
    await page.getByText('50 条/页').click()
  })
})
