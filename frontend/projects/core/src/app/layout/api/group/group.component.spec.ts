/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-14 22:56:33
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-20 22:02:26
 * @FilePath: /apinto/src/app/layout/upstream/service-discovery-content/service-discovery-content.component.spec.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { RouterModule } from '@angular/router'
import { ElementRef, Renderer2, ChangeDetectorRef } from '@angular/core'
import { APP_BASE_HREF } from '@angular/common'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { environment } from 'projects/core/src/environments/environment'
import { BidiModule } from '@angular/cdk/bidi'
import { Overlay } from '@angular/cdk/overlay'
import { of } from 'rxjs'
import { FormsModule } from '@angular/forms'
import { EoNgFeedbackModalService, EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { EoNgSelectModule } from 'eo-ng-select'
import { LayoutModule } from 'projects/core/src/app/layout/layout.module'
import { EoNgTableModule } from 'eo-ng-table'
import { ApiMessageComponent } from '../message/message.component'
import { ApiManagementListComponent } from '../list/list.component'
import { ApiPublishComponent } from '../publish/publish.component'
import { ApiManagementComponent } from './group.component'

class MockRenderer {
  removeAttribute (element: any, cssClass: string) {
    return cssClass + 'is removed from' + element
  }
}

class MockMessageService {
  success () {
    return 'success'
  }

  error () {
    return 'error'
  }
}

class MockEnsureService {
  create () {
    return 'modal is create'
  }
}

jest.mock('uuid', () => {
  return {
    v4: () => 123456789
  }
})

describe('ApiManagementComponent test as editPage is false', () => {
  let component: ApiManagementComponent
  let fixture: ComponentFixture<ApiManagementComponent>
  class MockElementRef extends ElementRef {
    constructor () { super(null) }
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule, FormsModule, EoNgSelectModule, LayoutModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule,
        NzDrawerModule, NzOutletModule, HttpClientModule,
        EoNgTableModule,
        RouterModule.forRoot([
          {
            path: '',
            component: ApiManagementListComponent
          },
          {
            path: 'message',
            component: ApiMessageComponent
          },
          {
            path: 'publish',
            component: ApiPublishComponent
          }
        ]
        )
      ],
      declarations: [ApiManagementComponent
      ],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: Renderer2, useClass: MockRenderer },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: EoNgFeedbackModalService, useClass: MockEnsureService },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(ApiManagementComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('ngOnInit should call getMenuList() and initial folderMenus', () => {
    const spyGetMenuList = jest.spyOn(component, 'getMenuList')
    expect(spyGetMenuList).not.toHaveBeenCalled()
    component.floderMenus = []
    component.fileMenus = []
    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetMenuList).toHaveBeenCalledTimes(1)
    expect(component.floderMenus).not.toStrictEqual([])
    expect(component.fileMenus).not.toStrictEqual([])
  })

  it('getMenuList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { api: { group_uuid: 123456 }, root: { groups: [] } } }))
    const spyApiMapTransfer = jest.spyOn(component, 'apiMapTransfer')
    const spyNodesTransfer = jest.spyOn(component, 'nodesTransfer')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyApiMapTransfer).not.toHaveBeenCalled()
    expect(spyNodesTransfer).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getMenuList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyApiMapTransfer).toHaveBeenCalled()
    expect(spyNodesTransfer).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('getMenuList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getMenuList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('apiMapTransfer should change list to map', fakeAsync(() => {
    component.apiNodesMap = new Map()
    component.apiNodesMap.set('testGroupUuid', ['testApi'])
    expect(component.apiNodesMap.size).toStrictEqual(1)
    expect(component.apiNodesMap.get('testGroupUuid')).toStrictEqual(['testApi'])

    component.apiMapTransfer([])
    fixture.detectChanges()
    expect(component.apiNodesMap.size).toStrictEqual(0)
    expect(component.apiNodesMap.get('testGroupUuid')).toBeUndefined()

    component.apiMapTransfer([
      { title: 'testApi1', name: 'testApi1', uuid: '111', group_uuid: '1111' },
      { title: 'testApi2', name: 'testApi2', uuid: '222', group_uuid: '2222' },
      { title: 'testApi3', name: 'testApi3', uuid: '333', group_uuid: '1111' }])
    fixture.detectChanges()
    expect(component.apiNodesMap.size).toStrictEqual(2)
    expect(component.apiNodesMap.get('1111')).toStrictEqual([
      { title: 'testApi1', name: 'testApi1', key: '111', isLeaf: true, uuid: '111', group_uuid: '1111' },
      { title: 'testApi3', name: 'testApi3', key: '333', isLeaf: true, uuid: '333', group_uuid: '1111' }
    ])
  }))

  it('nodesTransfer test', fakeAsync(() => {
    const val1:any = [
      { uuid: '111', name: 'name1', children: [] },
      {
        uuid: '222',
        name: 'name2',
        children: [
          { uuid: '2221', name: 'name21', children: [] },
          { uuid: '2222', name: 'name22', children: [] },
          {
            uuid: '2223',
            name: 'name23',
            children: [
              { uuid: '22231', name: 'name231', children: [] },
              { uuid: '22232', name: 'name232', children: [] }
            ]
          },
          { uuid: '2224', name: 'name24', children: [] }
        ]
      },
      { uuid: '333', name: 'name3', children: [] }
    ]

    component.firstLevelMap = new Set()
    let res = component.nodesTransfer([])
    fixture.detectChanges()
    expect(res).toStrictEqual([])

    component.expandAll = true
    component.firstLevelMap = new Set()
    component.apiNodesMap.set('2221', [
      { title: 'testApi1', name: 'testApi1', key: '111', isLeaf: true, uuid: '111', group_uuid: '1111' },
      { title: 'testApi3', name: 'testApi3', key: '333', isLeaf: true, uuid: '333', group_uuid: '1111' }
    ])
    component.apiNodesMap.set('444', [])
    component.apiNodesMap.set('333', [])
    res = component.nodesTransfer(val1, true)
    fixture.detectChanges()
    expect(res).toStrictEqual([
      { key: '111', title: 'name1', uuid: '111', name: 'name1', children: [] },
      {
        key: '222',
        title: 'name2',
        expanded: true,
        uuid: '222',
        name: 'name2',
        children: [
          {
            key: '2221',
            title: 'name21',
            expanded: true,
            uuid: '2221',
            name: 'name21',
            children: [
              { title: 'testApi1', name: 'testApi1', key: '111', isLeaf: true, uuid: '111', group_uuid: '1111' },
              { title: 'testApi3', name: 'testApi3', key: '333', isLeaf: true, uuid: '333', group_uuid: '1111' }
            ]
          },
          { key: '2222', title: 'name22', uuid: '2222', name: 'name22', children: [] },
          {
            key: '2223',
            title: 'name23',
            expanded: true,
            uuid: '2223',
            name: 'name23',
            children: [
              { key: '22231', title: 'name231', uuid: '22231', name: 'name231', children: [] },
              { key: '22232', title: 'name232', uuid: '22232', name: 'name232', children: [] }
            ]
          },
          { key: '2224', title: 'name24', uuid: '2224', name: 'name24', children: [] }
        ]
      },
      { key: '333', title: 'name3', uuid: '333', name: 'name3', children: [] }
    ])
    expect(component.firstLevelMap).toStrictEqual(new Set(['111', '222', '333']))
  }))

  it('addGroupModal test', fakeAsync(() => {
    // @ts-ignore
    const spyModalService = jest.spyOn(component.modalService, 'create')
    expect(spyModalService).not.toHaveBeenCalled()
    component.addGroupModal('root')
    fixture.detectChanges()
    expect(spyModalService).toHaveBeenCalledTimes(1)

    component.addGroupModal({ data: { uuid: '123' } })
    fixture.detectChanges()
    expect(spyModalService).toHaveBeenCalledTimes(2)
  }))

  it('addGroup with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: 0, data: { } }))
    const spyGetMenuList = jest.spyOn(component, 'getMenuList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.addGroup('test')
    fixture.detectChanges()
    expect(spyGetMenuList).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('addGroup with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: -1, data: { } }))
    const spyGetMenuList = jest.spyOn(component, 'getMenuList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.addGroup('test')
    fixture.detectChanges()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('editGroupModal test', fakeAsync(() => {
    // @ts-ignore
    const spyModalService = jest.spyOn(component.modalService, 'create')
    expect(spyModalService).not.toHaveBeenCalled()
    component.editGroupModal('root')
    fixture.detectChanges()
    expect(spyModalService).toHaveBeenCalledTimes(1)
  }))

  it('editGroup with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 0, data: { } }))
    const spyGetMenuList = jest.spyOn(component, 'getMenuList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyService).not.toHaveBeenCalled()
    component.editGroup('test')
    fixture.detectChanges()
    expect(spyGetMenuList).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('editGroup with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: -1, data: { } }))
    const spyGetMenuList = jest.spyOn(component, 'getMenuList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyService).not.toHaveBeenCalled()
    component.editGroup('test')
    fixture.detectChanges()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('deleteGroupModal test', fakeAsync(() => {
    // @ts-ignore
    const spyModalService = jest.spyOn(component.modalService, 'create')
    expect(spyModalService).not.toHaveBeenCalled()
    component.deleteGroupModal({ data: { name: 'test' } })
    fixture.detectChanges()
    expect(spyModalService).toHaveBeenCalledTimes(1)
  }))

  it('deleteGroup with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: 0, data: { } }))
    const spyGetMenuList = jest.spyOn(component, 'getMenuList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.deleteGroup('test', 'test')
    fixture.detectChanges()
    expect(spyGetMenuList).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('deleteGroup with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: -1, data: { } }))
    const spyGetMenuList = jest.spyOn(component, 'getMenuList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.deleteGroup('test', 'test')
    fixture.detectChanges()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('deleteApiModal test', fakeAsync(() => {
    // @ts-ignore
    const spyModalService = jest.spyOn(component.modalService, 'create')
    expect(spyModalService).not.toHaveBeenCalled()
    component.deleteApiModal('root')
    fixture.detectChanges()
    expect(spyModalService).toHaveBeenCalledTimes(1)
  }))

  it('deleteApi with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: 0, data: { } }))
    const spyGetMenuList = jest.spyOn(component, 'getMenuList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.deleteApi({ data: { uuid: '123' } })
    fixture.detectChanges()
    expect(spyGetMenuList).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('deleteApi with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: -1, data: { } }))
    const spyGetMenuList = jest.spyOn(component, 'getMenuList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.deleteApi({ data: { uuid: '123' } })
    fixture.detectChanges()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('openFolder test', () => {
    const data:any = { node: { id: '123', isExpanded: false }, isExpanded: true }
    component.openFolder(data)
    expect(data.node.isExpanded).toStrictEqual(true)
  })

  it('activeNode test', () => {
    component.showList = false
    component.showApiPage = false
    component.editPage = true
    component.apiUuid = ''
    component.groupUuid = ''
    component.activatedNode = undefined

    const data = { node: { isExpanded: false, origin: { uuid: '123' } } }
    component.activeNode(data)
    fixture.detectChanges()
    expect(data.node.isExpanded).toStrictEqual(true)
    expect(component.showList).toStrictEqual(true)
    expect(component.showApiPage).toStrictEqual(false)
    expect(component.groupUuid).toStrictEqual('123')
    expect(component.activatedNode).toStrictEqual(data.node)

    component.showList = false
    component.showApiPage = false
    component.editPage = false
    component.apiUuid = ''
    component.groupUuid = ''
    component.activatedNode = undefined

    const data2 = { node: { isExpanded: false, origin: { group_uuid: '1234', uuid: '123' } } }
    component.activeNode(data2)
    fixture.detectChanges()
    expect(data2.node.isExpanded).toStrictEqual(true)
    expect(component.showList).toStrictEqual(false)
    expect(component.showApiPage).toStrictEqual(true)
    expect(component.editPage).toStrictEqual(true)
    expect(component.apiUuid).toStrictEqual('123')
    expect(component.groupUuid).toStrictEqual('')
    expect(component.activatedNode).toStrictEqual(data2.node)
  })

  it('changeToList & viewAllApis & addApi & editApi test', () => {
    const spyGetMenuList = jest.spyOn(component, 'getMenuList')
    component.showList = false
    component.showApiPage = false
    component.editPage = false
    component.apiUuid = ''
    component.groupUuid = ''
    expect(spyGetMenuList).not.toHaveBeenCalled()

    component.changeToList('123')
    fixture.detectChanges()
    expect(component.showList).toStrictEqual(true)
    expect(component.showApiPage).toStrictEqual(false)
    expect(component.groupUuid).toStrictEqual('123')
    expect(spyGetMenuList).toHaveBeenCalled()

    component.showList = false
    component.showApiPage = false
    component.editPage = false
    component.apiUuid = ''
    component.groupUuid = ''
    component.query_name = '123'
    expect(spyGetMenuList).toHaveBeenCalledTimes(1)

    component.viewAllApis()
    fixture.detectChanges()
    expect(component.showList).toStrictEqual(true)
    expect(component.showApiPage).toStrictEqual(false)
    expect(component.groupUuid).toStrictEqual('')
    expect(component.query_name).toStrictEqual('')
    expect(spyGetMenuList).toHaveBeenCalledTimes(2)

    component.showList = false
    component.showApiPage = false
    component.editPage = false
    component.apiUuid = ''
    component.groupUuid = ''
    expect(spyGetMenuList).toHaveBeenCalledTimes(2)

    component.addApi({ data: { uuid: '123' } })
    fixture.detectChanges()
    expect(component.showList).toStrictEqual(false)
    expect(component.showApiPage).toStrictEqual(true)
    expect(component.editPage).toStrictEqual(false)
    expect(component.groupUuid).toStrictEqual('123')
    expect(spyGetMenuList).toHaveBeenCalledTimes(2)

    component.showList = false
    component.showApiPage = false
    component.editPage = false
    component.apiUuid = ''
    component.groupUuid = ''
    expect(spyGetMenuList).toHaveBeenCalledTimes(2)

    component.editApi({ data: { uuid: '123' } })
    fixture.detectChanges()
    expect(component.showList).toStrictEqual(false)
    expect(component.showApiPage).toStrictEqual(true)
    expect(component.editPage).toStrictEqual(false)
    expect(component.apiUuid).toStrictEqual('123')
    expect(spyGetMenuList).toHaveBeenCalledTimes(2)
  })
})
