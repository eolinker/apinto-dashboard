/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-07-30 00:40:51
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-03 23:20:22
 * @FilePath: /apinto/src/app/layout/basic-layout/basic-layout.component.spec.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { BasicLayoutComponent } from './basic-layout.component'
import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'

describe('BasicLayoutComponent test', () => {
  let component: BasicLayoutComponent
  let fixture: ComponentFixture<BasicLayoutComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
      ],
      declarations: [
      ],
      providers: [
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(BasicLayoutComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })
  // })
  // describe('menulist initial', () => {
  //   let component: BasicLayoutComponent
  //   let fixture: ComponentFixture<BasicLayoutComponent>

  //   beforeEach(async () => {
  //     fixture = TestBed.createComponent(BasicLayoutComponent)
  //     component = fixture.componentInstance
  //     fixture.detectChanges()
  //   })

  it('should initial menuList and every submenu has open attribute', fakeAsync(() => {
    // const exp:any[] = [
    //   {
    //     title: '基础设施管理',
    //     routerLink: 'deploy/cluster',
    //     matchRouter: true,
    //     matchRouterExact: true,
    //     children: [
    //       {
    //         title: '集群管理',
    //         routerLink: 'deploy/cluster',
    //         menuIndex: 0,
    //         level: 1,
    //         matchRouter: true,
    //         matchRouterExact: false
    //       },
    //       {
    //         title: '环境变量',
    //         routerLink: 'deploy/environment',
    //         menuIndex: 0,
    //         level: 1,
    //         matchRouter: true,
    //         matchRouterExact: false
    //       }
    //     ],
    //     open: true,
    //     openChange: undefined
    //   }]
    // expect(component.openMap).toStrictEqual({ 上游服务: false, 基础设施管理: true })
    // expect(component.sideMenuOptions[1].children).toStrictEqual(exp[0].children)
  }))

  // it('shouldnt initial menuList when params is []', fakeAsync(() => {
  //   component.getInitMenuList([])
  //   fixture.detectChanges()
  //   tick(150)
  //   expect(component.openMap).toStrictEqual({})
  //   expect(component.sideMenuOptions).toStrictEqual([])
  // }))

  // it('openMap should be change when subtitle open status is open: openHandler is called', fakeAsync(() => {
  //   component.openMap = { test1: true, test2: false }
  //   expect(component.openMap).toStrictEqual({ test1: true, test2: false })
  //   component.openHandler('test2')
  //   fixture.detectChanges()
  //   expect(component.openMap).toStrictEqual({ test1: false, test2: true })
  // }))
})
