/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-14 22:56:33
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-20 22:01:52
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
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { FormsModule } from '@angular/forms'

import { EoNgSelectModule } from 'eo-ng-select'
import { LayoutModule } from 'projects/core/src/app/layout/layout.module'
import { EoNgTableModule } from 'eo-ng-table'
import { ApiMessageComponent } from '../message/message.component'
import { ApiPublishComponent } from '../publish/publish.component'
import { ApiManagementListComponent } from './list.component'
import { ApiCreateComponent } from '../create/create.component'

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

describe('ApiManagementListComponent test as editPage is false', () => {
  let component: ApiManagementListComponent
  let fixture: ComponentFixture<ApiManagementListComponent>
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
          },
          {
            path: 'router/create',
            component: ApiCreateComponent
          }
        ]
        )
      ],
      declarations: [ApiManagementListComponent
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

    Object.defineProperty(window, 'open', {
      value: jest.fn(() => { return { matches: true } })
    })

    fixture = TestBed.createComponent(ApiManagementListComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('click table btns', () => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')
    expect(spyRouter).toHaveBeenCalledTimes(0)
    const open = jest.fn()
    Object.defineProperty(window, 'open', open)
    expect(open).not.toHaveBeenCalled()

    const item = { key: 'test' }
    component.apisTableBody[6].btns[0].click(item)
    expect(spyRouter).toHaveBeenCalledTimes(1)
    component.apisTableBody[6].btns[1].click(item)
    expect(spyRouter).toHaveBeenCalledTimes(2)
    // component.onlineResultTableBody[4].btns[0].click({ solution: { params: '', name: '' } })
    // expect(open).toHaveBeenCalledTimes(1)
  })

  it('ngOnInit should call getApisData()', () => {
    const spyGetApisData = jest.spyOn(component, 'getApisData')
    expect(spyGetApisData).not.toHaveBeenCalled()
    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetApisData).toHaveBeenCalledTimes(1)
  })

  it('trStyleFn test', fakeAsync(() => {
    const res1 = component.trStyleFn({ status: true })
    expect(res1).toStrictEqual('color:green')

    const res2 = component.trStyleFn({})
    expect(res2).toStrictEqual('color:red')
  }))

  it('getApisData with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get')
      .mockReturnValue(of({
        code: 0,
        data: {
          apis: [{ group_uuid: 123456 }],
          total: 10,
          page_num: 10,
          page_size: 20
        }
      }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getApisData()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component.apisForm.apis).toStrictEqual([{ group_uuid: 123456, eoKey: 123456789, checked: false }])
    expect(component.apisForm.group_uuid).toStrictEqual('')
    expect(component.apisForm.total).toStrictEqual(10)
    expect(component.apisForm.page_num).toStrictEqual(10)
    expect(component.apisForm.page_size).toStrictEqual(20)
    expect(component.objDiffers).not.toStrictEqual([])
  })

  it('getApisData with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getApisData()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('addApi test', fakeAsync(() => {
    // @ts-ignore
    const spyRouterService = jest.spyOn(component.router, 'navigate')
    expect(spyRouterService).not.toHaveBeenCalled()
    component.addApi()
    fixture.detectChanges()
    expect(spyRouterService).toHaveBeenCalledTimes(1)
  }))

  it('openDrawer test', () => {
    expect(component.onlineDrawerRef).toBeUndefined()
    component.openDrawer('online')
    fixture.detectChanges()
    expect(component.onlineDrawerRef).not.toBeUndefined()

    expect(component.offlineDrawerRef).toBeUndefined()
    component.openDrawer('offline')
    fixture.detectChanges()
    expect(component.offlineDrawerRef).not.toBeUndefined()
  })

  it('apisOperatorResult test', () => {
    const spyOnlineApisCheck = jest.spyOn(component, 'onlineApisCheck')
    const spyOfflineApis = jest.spyOn(component, 'offlineApis')
    expect(component.onlineDrawerRef).toBeUndefined()
    expect(spyOnlineApisCheck).not.toHaveBeenCalled()
    component.apisOperatorResult('online')
    fixture.detectChanges()
    expect(component.onlineDrawerRef).not.toBeUndefined()
    expect(spyOnlineApisCheck).toHaveBeenCalledTimes(1)

    // @ts-ignore
    const spyOnlineDrawerRef = jest.spyOn(component.onlineDrawerRef, 'close')
    expect(spyOnlineDrawerRef).not.toHaveBeenCalled()
    component.apisOperatorResult('online-res')
    fixture.detectChanges()
    expect(component.onlineDrawerRef).not.toBeUndefined()
    expect(spyOnlineApisCheck).toHaveBeenCalledTimes(1)
    expect(spyOnlineDrawerRef).toHaveBeenCalled()

    expect(spyOfflineApis).not.toHaveBeenCalled()
    expect(component.offlineDrawerRef).toBeUndefined()
    component.apisOperatorResult('offline')
    fixture.detectChanges()
    expect(component.offlineDrawerRef).not.toBeUndefined()
    expect(spyOfflineApis).toHaveBeenCalled()
  })

  it('openDrawer test', () => {
    expect(component.onlineDrawerRef).toBeUndefined()
    component.openDrawer('online')
    fixture.detectChanges()
    expect(component.onlineDrawerRef).not.toBeUndefined()

    expect(component.offlineDrawerRef).toBeUndefined()
    component.openDrawer('offline')
    fixture.detectChanges()
    expect(component.offlineDrawerRef).not.toBeUndefined()
  })

  it('closeDrawer test', fakeAsync(() => {
    component.openDrawer('online')
    fixture.detectChanges()
    expect(component.onlineDrawerRef).not.toBeUndefined()

    // @ts-ignore
    const spyOnlineDrawerRef = jest.spyOn(component.onlineDrawerRef, 'close')

    component.closeDrawer('online')
    fixture.detectChanges()
    flush()
    expect(spyOnlineDrawerRef).toHaveBeenCalledTimes(1)

    component.openDrawer('offline')
    fixture.detectChanges()
    flush()
    expect(component.offlineDrawerRef).not.toBeUndefined()

    // @ts-ignore
    const spyOfflineDrawerRef = jest.spyOn(component.offlineDrawerRef, 'close')

    component.closeDrawer('offline')
    fixture.detectChanges()
    flush()
    expect(spyOfflineDrawerRef).toHaveBeenCalledTimes(1)

    component.openDrawer('online')
    component.openDrawer('offline')
    fixture.detectChanges()
    flush()
    expect(component.onlineDrawerRef).not.toBeUndefined()
    expect(component.offlineDrawerRef).not.toBeUndefined()

    component.closeDrawer('res')
    fixture.detectChanges()
    flush()
    expect(spyOnlineDrawerRef).toHaveBeenCalledTimes(2)
    expect(spyOfflineDrawerRef).toHaveBeenCalledTimes(2)
  }))

  it('onlineApisCheck with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: 0, data: { list: [1, 2, 3], online_token: '123' } }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.onlineApisCheck()
    fixture.detectChanges()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component.onlineResultList).toStrictEqual([1, 2, 3])
    expect(component.onlineToken).toStrictEqual('123')
  })

  it('onlineApisCheck with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: -1, data: { }, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.onlineApisCheck()
    fixture.detectChanges()

    expect(spyMessage).toHaveBeenCalled()
  })

  it('onlineApis with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: 0, data: { list: [1, 2, 3], online_token: '123' } }))
    const spyApisOperatorResult = jest.spyOn(component, 'apisOperatorResult')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyApisOperatorResult).not.toHaveBeenCalled()
    component.onlineApis()
    fixture.detectChanges()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyApisOperatorResult).toHaveBeenCalled()
    expect(component.resultList).toStrictEqual([1, 2, 3])
  })

  it('onlineApis with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: -1, data: { }, msg: 'fail' }))
    const spyApisOperatorResult = jest.spyOn(component, 'apisOperatorResult')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyApisOperatorResult).not.toHaveBeenCalled()
    component.onlineApis()
    fixture.detectChanges()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyApisOperatorResult).toHaveBeenCalled()
    expect(component.resultList).toBeUndefined()
  })

  it('offlineApis with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: 0, data: { list: [1, 2, 3], online_token: '123' } }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.offlineApis()
    fixture.detectChanges()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component.resultList).toStrictEqual([1, 2, 3])
  })

  it('offlineApis with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: -1, data: { } }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.offlineApis()
    fixture.detectChanges()
    expect(component.resultList).toBeUndefined()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('changeClustersSet test', fakeAsync(() => {
    const val1 = [
      { value: 'val1', checked: false },
      { value: 'val2', checked: true },
      { value: 'val3', checked: true }]
    component.clustersSet = new Set()
    component.changeClustersSet(val1)
    fixture.detectChanges()
    expect(component.clustersSet).toStrictEqual(new Set(['val2', 'val3']))

    const val2 = [
      { value: 'val1', checked: false },
      { value: 'val2', checked: false },
      { value: 'val3', checked: true }]
    component.clustersSet = new Set()
    component.changeClustersSet(val2)
    fixture.detectChanges()
    expect(component.clustersSet).toStrictEqual(new Set(['val3']))
  }))

  it('getClusterList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get')
      .mockReturnValue(of({
        code: 0,
        data: {
          clusters: [
            { name: 'name1', env: 'env1', status: 'ABNORMAL' },
            { name: 'name2', env: 'env2' },
            { name: 'name3', env: 'env3', status: 'ABNORMAL' }
          ]
        }
      }))
    const spyOpenDrawer = jest.spyOn(component, 'openDrawer')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyOpenDrawer).not.toHaveBeenCalled()

    component.getClusterList('test')
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyOpenDrawer).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('getClusterList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getClusterList('test')
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('deleteApi with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: 0, data: { } }))
    const spyGetApisData = jest.spyOn(component, 'getApisData')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyGetApisData).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.deleteApi({ uuid: '123' })
    fixture.detectChanges()
    expect(spyGetApisData).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('deleteApi with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: -1, data: { } }))
    const spyGetApisData = jest.spyOn(component, 'getApisData')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyGetApisData).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.deleteApi({ uuid: '123' })
    fixture.detectChanges()
    expect(spyGetApisData).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })
})
