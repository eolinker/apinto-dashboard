/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-14 22:56:33
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-20 22:02:10
 * @FilePath: /apinto/src/app/layout/upstream/service-discovery-content/service-discovery-content.component.spec.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { ComponentFixture, fakeAsync, flush, TestBed } from '@angular/core/testing'
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
import { ApiCreateComponent } from './create.component'
import { ApiMessageComponent } from '../message/message.component'
import { ApiManagementListComponent } from '../list/list.component'
import { ApiPublishComponent } from '../publish/publish.component'

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

describe('ApiCreateComponent test as editPage is false', () => {
  let component: ApiCreateComponent
  let fixture: ComponentFixture<ApiCreateComponent>
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
      declarations: [ApiCreateComponent
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

    fixture = TestBed.createComponent(ApiCreateComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('click table btns', () => {
    // @ts-ignore
    const spyOpenDrawer = jest.spyOn(component, 'openDrawer')
    expect(spyOpenDrawer).toHaveBeenCalledTimes(0)

    const item = { key: 'test' }
    component.matchTableBody[4].btns[0].click(item)
    expect(spyOpenDrawer).toHaveBeenCalledTimes(1)

    component.createApiForm.match = [
      { position: 'test1', matchType: 'test1', key: 'test1', pattern: 'test1' },
      { position: 'test2', matchType: 'test2', key: 'test2', pattern: 'test2' },
      { position: 'test3', matchType: 'test3', key: 'test3', pattern: 'test3' }
    ]
    component.matchTableBody[4].btns[1].click({ position: 'test2', matchType: 'test2', key: 'test2', pattern: 'test2' })
    fixture.detectChanges()
    component.createApiForm.match = [
      { position: 'test1', matchType: 'test1', key: 'test1', pattern: 'test1' },
      { position: 'test3', matchType: 'test3', key: 'test3', pattern: 'test3' }
    ]

    component.proxyHeaderTableBody[3].btns[0].click(item)
    expect(spyOpenDrawer).toHaveBeenCalledTimes(2)

    component.createApiForm.proxyHeader = [
      { optType: 'test1', key: 'test1', value: 'test1' },
      { optType: 'test2', key: 'test2', value: 'test2' },
      { optType: 'test3', key: 'test3', value: 'test3' }
    ]
    component.proxyHeaderTableBody[3].btns[1].click({ optType: 'test3', key: 'test3', value: 'test3' })
    fixture.detectChanges()
    component.createApiForm.proxyHeader = [
      { optType: 'test1', key: 'test1', value: 'test1' },
      { optType: 'test2', key: 'test2', value: 'test2' }
    ]
  })

  it('openDrawer test', () => {
    expect(component.drawerMatchRef).toBeUndefined()
    component.openDrawer('match')
    fixture.detectChanges()
    expect(component.drawerMatchRef).not.toBeUndefined()

    expect(component.drawerProxyRef).toBeUndefined()
    component.openDrawer('proxyHeader')
    fixture.detectChanges()
    expect(component.drawerProxyRef).not.toBeUndefined()
  })

  it('closeDrawer test', fakeAsync(() => {
    component.openDrawer('match')
    fixture.detectChanges()
    expect(component.drawerMatchRef).not.toBeUndefined()

    // @ts-ignore
    const spyDrawerMatchRef = jest.spyOn(component.drawerMatchRef, 'close')

    component.closeDrawer('match')
    fixture.detectChanges()
    flush()
    expect(spyDrawerMatchRef).toHaveBeenCalledTimes(1)

    component.openDrawer('proxyHeader')
    fixture.detectChanges()
    flush()
    expect(component.drawerProxyRef).not.toBeUndefined()

    // @ts-ignore
    const spyOfflineDrawerRef = jest.spyOn(component.drawerProxyRef, 'close')

    component.closeDrawer('proxyHeader')
    fixture.detectChanges()
    flush()
    expect(spyOfflineDrawerRef).toHaveBeenCalledTimes(1)
  }))

  it('getApiMessage with success return (method is null)', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { api: { groupUuid: '123456', match: [], proxyHeader: [] } } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyUpdateAllChecked = jest.spyOn(component, 'updateAllChecked')
    const spyInitCheckbox = jest.spyOn(component, 'initCheckbox')
    const spyGetHeaderList = jest.spyOn(component, 'getHeaderList')
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyUpdateAllChecked).not.toHaveBeenCalled()
    expect(spyInitCheckbox).not.toHaveBeenCalled()
    expect(spyGetHeaderList).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.allChecked = false
    component.getApiMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(3)
    expect(component.createApiForm).toStrictEqual(
      {
        groupUuid: ['123456'],
        match: [],
        proxyHeader: [],
        method: ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD']
      }
    )
    expect(component.allChecked).toStrictEqual(true)
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyUpdateAllChecked).toHaveBeenCalled()
    expect(spyInitCheckbox).not.toHaveBeenCalled()
    expect(spyGetHeaderList).toHaveBeenCalled()
  })

  it('getApiMessage with success return (method is null)', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { api: { groupUuid: 123456, match: [], proxyHeader: [], method: ['PUT', 'DELETE'] } } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyUpdateAllChecked = jest.spyOn(component, 'updateAllChecked')
    const spyInitCheckbox = jest.spyOn(component, 'initCheckbox')
    const spyGetHeaderList = jest.spyOn(component, 'getHeaderList')
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyUpdateAllChecked).not.toHaveBeenCalled()
    expect(spyInitCheckbox).not.toHaveBeenCalled()
    expect(spyGetHeaderList).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.allChecked = false
    component.getApiMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(3)
    expect(component.createApiForm).toStrictEqual(
      {
        groupUuid: [123456],
        method: ['PUT', 'DELETE'],
        match: [],
        proxyHeader: []
      }
    )
    expect(component.allChecked).toStrictEqual(false)
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyUpdateAllChecked).not.toHaveBeenCalled()
    expect(spyInitCheckbox).toHaveBeenCalled()
    expect(spyGetHeaderList).toHaveBeenCalled()
  })

  it('getApiMessage with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)
    expect(component.createApiForm.groupUuid).toStrictEqual([])

    component.getApiMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.createApiForm.groupUuid).toStrictEqual([])
    expect(spyMessage).toHaveBeenCalled()
  })

  it('getHeaderList with success return', () => {
    const mockValue:any = {
      code: 0,
      data: {
        root: {
          groups: [
            { uuid: '111', children: [] },
            { uuid: '222', children: [{ uuid: '2223' }] }
          ]
        }
      }
    }

    const mockResList:Array<any> = [
      { uuid: '111', disabled: true, isLeaf: true, children: [] },
      {
        uuid: '222',
        children: [
          { uuid: '2223', isLeaf: true }
        ]
      }
    ]

    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of(mockValue))
    const isget = httpCommonService.get('') !== null
    const spyFindGroup = jest.spyOn(component, 'findGroup')

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(spyFindGroup).toHaveBeenCalledTimes(0)
    component.getHeaderList()
    component.createApiForm.groupUuid = ['1']
    fixture.detectChanges()
    expect(component.headerList).toStrictEqual(mockResList)
    expect(spyService).toHaveBeenCalledTimes(2)
    // expect(spyFindGroup).toHaveBeenCalledTimes(1)
  })

  it('getHeaderList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.getHeaderList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalled()
  })

  it('getServiceList with success return', () => {
    const mockValue:any = {
      code: 0,
      data: {
        list: [
          'test1', 'test2', 'test3'
        ]
      }
    }

    const mockResList:Array<any> = [
      { label: 'test1', value: 'test1' },
      { label: 'test2', value: 'test2' },
      { label: 'test3', value: 'test3' }
    ]

    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of(mockValue))
    const isget = httpCommonService.get('') !== null

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.getServiceList()
    fixture.detectChanges()
    expect(component.serviceList).toStrictEqual(mockResList)
    expect(spyService).toHaveBeenCalledTimes(2)
  })

  it('getServiceList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.getServiceList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalled()
  })

  it('updateSingleChecked test', fakeAsync(() => {
    component.allChecked = false
    component.methodList = [
      { label: '1', value: '1', checked: true },
      { label: '2', value: '2', checked: false },
      { label: '3', value: '3', checked: true }
    ]
    component.updateAllChecked()
    expect(component.methodList).toStrictEqual([
      { label: '1', value: '1', checked: false },
      { label: '2', value: '2', checked: false },
      { label: '3', value: '3', checked: false }
    ])
  }))

  it('updateSingleChecked test', fakeAsync(() => {
    component.allChecked = true
    component.methodList = [
      { label: '1', value: '1', checked: true },
      { label: '2', value: '2', checked: false },
      { label: '3', value: '3', checked: true }
    ]
    component.updateSingleChecked()
    fixture.detectChanges()
    flush()
    expect(component.allChecked).toStrictEqual(false)
    expect(component.createApiForm.method).toStrictEqual(['1', '3'])

    component.methodList = [
      { label: '1', value: '1', checked: true },
      { label: '2', value: '2', checked: true },
      { label: '3', value: '3', checked: true }
    ]
    component.updateSingleChecked()
    fixture.detectChanges()
    flush()
    expect(component.allChecked).toStrictEqual(true)
    expect(component.createApiForm.method).toStrictEqual(['1', '2', '3'])

    component.methodList = [
      { label: '1', value: '1', checked: false },
      { label: '2', value: '2', checked: false },
      { label: '3', value: '3', checked: false }
    ]
    component.updateSingleChecked()
    fixture.detectChanges()
    flush()
    expect(component.allChecked).toStrictEqual(false)
    expect(component.createApiForm.method).toStrictEqual([])
  }))

  it('saveApi with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyServicePost = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: 0, data: { msg: 'success' } }))
    const ispost = httpCommonService.post('') !== null
    const spyServicePut = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 0, data: { msg: 'success' } }))
    const isput = httpCommonService.put('') !== null

    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    const spybackToList = jest.spyOn(component, 'backToList')
    expect(spybackToList).not.toHaveBeenCalled()

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    expect(spyMessageError).not.toHaveBeenCalled()

    expect(spyServicePost).toHaveBeenCalledTimes(1)
    expect(ispost).toStrictEqual(true)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(isput).toStrictEqual(true)

    component.createApiForm = {
      name: '',
      desc: '',
      uuid: 'test',
      groupUuid: '123',
      requestPath: '123',
      service: '123',
      method: ['get'],
      proxyPath: '123',
      retry: 1,
      timeout: 100,
      proxyHeader: [{
        key: '123', value: '123', optType: '123'
      }],
      match: [{
        position: '123', matchType: '123', key: '123', pattern: '123'
      }]
    }

    component.editPage = false
    component.saveApi()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(spybackToList).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
    expect(spyMessageError).not.toHaveBeenCalled()

    component.editPage = true
    component.saveApi()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyServicePut).toHaveBeenCalledTimes(2)
    expect(spybackToList).toHaveBeenCalledTimes(2)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(2)
    expect(spyMessageError).not.toHaveBeenCalled()
  })

  it('saveApi with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyServicePost = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const ispost = httpCommonService.post('') !== null
    const spyServicePut = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isput = httpCommonService.put('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyServicePost).toHaveBeenCalledTimes(1)
    expect(ispost).toStrictEqual(true)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(isput).toStrictEqual(true)

    component.editPage = false
    component.saveApi()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalledTimes(1)

    component.editPage = true
    component.saveApi()
    fixture.detectChanges()

    expect(spyServicePut).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalledTimes(2)
  })

  it('backToList test', fakeAsync(() => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')
    expect(spyRouter).not.toHaveBeenCalled()

    const spyChangeToList = jest.spyOn(component.changeToList, 'emit')
    expect(spyChangeToList).not.toHaveBeenCalled()
    component.fromList = false
    component.backToList()
    fixture.detectChanges()

    expect(spyRouter).toHaveBeenCalled()

    component.fromList = true
    component.backToList()
    fixture.detectChanges()

    expect(spyRouter).toHaveBeenCalledTimes(1)
    expect(spyChangeToList).toHaveBeenCalledTimes(1)
  }))

  it('requestPathChange test', fakeAsync(() => {
    component.createApiForm.proxyPath = ''
    component.createApiForm.requestPath = 'test'
    component.requestPathChange()
    fixture.detectChanges()
    expect(component.createApiForm.proxyPath).toStrictEqual('test')

    component.createApiForm.requestPath = 'test2'
    component.requestPathChange()
    fixture.detectChanges()
    expect(component.createApiForm.proxyPath).toStrictEqual('test')
  }))

  it('checkTimeout test', fakeAsync(() => {
    component.createApiForm.timeout = -100
    component.checkTimeout()
    fixture.detectChanges()
    expect(component.createApiForm.timeout).toStrictEqual(1)

    component.checkTimeout()
    fixture.detectChanges()
    expect(component.createApiForm.timeout).toStrictEqual(1)
  }))

  it('clearMatchForm test', fakeAsync(() => {
    component.matchForm = {
      key: '',
      matchType: '1',
      pattern: '1',
      position: '1'
    }
    component.clearMatchForm()
    expect(component.matchForm).toStrictEqual({
      key: '',
      matchType: '',
      pattern: '',
      position: ''
    })
  }))

  it('clearProxyRequestForm test', fakeAsync(() => {
    component.proxyHeaderForm = {
      key: '1',
      value: '',
      optType: '123'
    }
    component.clearProxyRequestForm()
    expect(component.proxyHeaderForm).toStrictEqual({
      key: '',
      value: '',
      optType: ''
    })
  }))
})
