import { ComponentFixture, TestBed, discardPeriodicTasks, fakeAsync, tick } from '@angular/core/testing'
import { ComponentModule } from 'projects/core/src/app/component/component.module'
import { APP_BASE_HREF } from '@angular/common'
import { HttpClientModule } from '@angular/common/http'
import { ElementRef, Renderer2, ChangeDetectorRef } from '@angular/core'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { NavigationEnd, Router, RouterModule } from '@angular/router'
import { Overlay } from '@angular/cdk/overlay'
import { EoNgFeedbackMessageService, EoNgFeedbackModalModule, EoNgFeedbackModalService, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { BidiModule } from '@angular/cdk/bidi'
import { MockRenderer, MockMessageService, MockEnsureService, MockEmptySuccessResponse, MockApiSource, MockApisList, MockApisList2 } from 'projects/core/src/app/constant/spec-test'
import { BehaviorSubject, of } from 'rxjs'
import { API_URL, ApiService } from 'projects/core/src/app/service/api.service'
import { environment } from 'projects/core/src/environments/environment'
import { NzFormModule } from 'ng-zorro-antd/form'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgTreeModule } from 'eo-ng-tree'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { EoNgSelectModule } from 'eo-ng-select'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { LayoutModule } from '../../../layout.module'
import { routes } from '../../api-routing.module'
import { ApiManagementListComponent } from './list.component'
import { EoNgDropdownModule } from 'eo-ng-dropdown'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

describe('#init ApiManagementListComponent', () => {
  let component:ApiManagementListComponent
  let fixture: ComponentFixture<ApiManagementListComponent>
  let httpCommonService:any
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  let spyDeleteApiService:jest.SpyInstance<any>
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  let spyApiService:jest.SpyInstance<any>
  const eventsSub = new BehaviorSubject<any>(null)
  const routerStub = {
    events: eventsSub,
    url: '',
    navigate: (...args:Array<string>) => {
      eventsSub.next(new NavigationEnd(1, args.join('/'), args.join('/')))
    }
  }
  global.structuredClone = (val:any) => JSON.parse(JSON.stringify(val))

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule, FormsModule, ReactiveFormsModule, ComponentModule, LayoutModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule, NzOutletModule, HttpClientModule,
        RouterModule.forRoot(routes), NzFormModule, EoNgInputModule, EoNgTreeModule, EoNgButtonModule,
        EoNgSwitchModule, EoNgCheckboxModule, EoNgApintoTableModule, EoNgSelectModule, EoNgFeedbackModalModule,
        EoNgFeedbackTooltipModule, EoNgDropdownModule
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
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef },
        { provide: Router, useValue: routerStub }
      ],
      teardown: { destroyAfterEach: false }
    }).compileComponents()

    fixture = TestBed.createComponent(ApiManagementListComponent)
    component = fixture.componentInstance

    fixture.detectChanges()

    httpCommonService = fixture.debugElement.injector.get(ApiService)

    spyApiService = jest.spyOn(httpCommonService, 'get').mockImplementation(
      (...args) => {
        switch (args[0]) {
          case 'routers':
            return of(MockApisList)
          case 'router/source':
            return of(MockApiSource)
          default:
            return of(MockEmptySuccessResponse)
        }
      }
    )

    spyDeleteApiService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(
      of(MockEmptySuccessResponse)
    )
  })

  it('should create and init component as creating a new api', fakeAsync(() => {
    expect(component).toBeTruthy()
    expect(component.sourcesList).toEqual([])
    expect(component.apisForm.apis).toEqual([])
    expect(component.apisForm.total).toEqual(0)
    expect(component.apisForm.pageNum).toEqual(1)
    expect(component.apisForm.pageSize).toEqual(20)
    expect(component.nzDisabled).toEqual(false)
    // @ts-ignore
    jest.replaceProperty(fixture.debugElement.injector.get(BaseInfoService), '_allParams', {
      apiGroupId: 'mockApiGroupId'
    })
    component.ngOnInit()
    component.ngAfterViewInit()
    fixture.detectChanges()
    tick(500)
    expect(component.sourcesList).toEqual([
      { text: '自建', value: 'self-build:-1:' },
      { text: '导入', value: 'import:-1:' },
      { text: 'Apikit', value: 'sync:1:' },
      { text: 'Postcat', value: 'sync:2:' }
    ])
    // @ts-ignore
    expect(component.baseInfo.allParamsInfo.apiGroupId).toEqual('mockApiGroupId')
    expect(component.nzDisabled).toEqual(false)
    expect(component.apisTableHeadName.length).toEqual(12)
    expect(component.apisTableBody.length).toEqual(12)
    expect(component.apisForm.apis).not.toEqual([])
    expect(component.apisForm.groupUuid).toEqual('mockApiGroupId')
    expect(component.apisForm.total).toEqual(7)
    expect(component.apisForm.pageNum).toEqual(1)
    expect(component.apisForm.pageSize).toEqual(20)
    expect(component.apiTableLoading).toEqual(false)
    discardPeriodicTasks()
  }))

  it('should change apisForm and get new api data when url changed', fakeAsync(() => {
    const spyGetApiData = jest.spyOn(component, 'getApisData')
    expect(component).toBeTruthy()
    expect(component.sourcesList).toEqual([])
    expect(component.apisForm.apis).toEqual([])
    expect(component.apisForm.total).toEqual(0)
    expect(component.apisForm.pageNum).toEqual(1)
    expect(component.apisForm.pageSize).toEqual(20)
    expect(component.nzDisabled).toEqual(false)
    expect(spyGetApiData).not.toHaveBeenCalled()
    // @ts-ignore
    jest.replaceProperty(fixture.debugElement.injector.get(BaseInfoService), '_allParams', {
      apiGroupId: 'mockApiGroupId'
    })
    component.ngOnInit()
    component.ngAfterViewInit()
    fixture.detectChanges()
    tick()
    // @ts-ignore
    expect(component.baseInfo.allParamsInfo.apiGroupId).toEqual('mockApiGroupId')
    expect(component.nzDisabled).toEqual(false)
    expect(component.apisTableHeadName.length).toEqual(12)
    expect(component.apisTableBody.length).toEqual(12)
    expect(component.apisForm.apis).not.toEqual([])
    expect(component.apisForm.groupUuid).toEqual('mockApiGroupId')
    expect(component.apisForm.total).toEqual(7)
    expect(component.apisForm.pageNum).toEqual(1)
    expect(component.apisForm.pageSize).toEqual(20)
    expect(component.apiTableLoading).toEqual(false)
    expect(spyGetApiData).toHaveBeenCalledTimes(1)

    spyApiService = jest.spyOn(httpCommonService, 'get').mockImplementation(
      (...args) => {
        switch (args[0]) {
          case 'routers':
            return of(MockApisList2)
          case 'router/source':
            return of(MockApiSource)
          default:
            return of(MockEmptySuccessResponse)
        }
      }
    )

    // @ts-ignore
    jest.replaceProperty(fixture.debugElement.injector.get(BaseInfoService), '_allParams', {
      apiGroupId: 'test'
    })

    component.apisSet = new Set(['549a9e3f-50ed-f004-c033-df1cc6c48df1'])
    component.apisTableClick({ data: { uuid: 'test', scheme: 'websocket' } })
    fixture.detectChanges()
    tick()
    expect(component.nzDisabled).toEqual(false)
    expect(component.apisTableHeadName.length).toEqual(9)
    expect(component.apisTableBody.length).toEqual(9)
    expect(component.apisForm.groupUuid).toEqual('test')
    expect(component.apisForm.apis.length).toEqual(1)
    expect(component.apisForm.apis[0].checked).toEqual(true)
    expect(component.apisForm.total).toEqual(7)
    expect(component.apisForm.pageNum).toEqual(1)
    expect(component.apisForm.pageSize).toEqual(20)
    expect(component.apiTableLoading).toEqual(false)
    expect(spyGetApiData).not.toHaveBeenCalledTimes(1)

    component.apisTableClick({ data: { uuid: 'test', scheme: 'http' } })
    fixture.detectChanges()
    tick()
    expect(component.nzDisabled).toEqual(false)
    expect(component.apisTableHeadName.length).toEqual(9)
    expect(component.apisTableBody.length).toEqual(9)
    expect(component.apisForm.groupUuid).toEqual('test')
    expect(component.apisForm.apis.length).toEqual(1)
    expect(component.apisForm.apis[0].checked).toEqual(true)
    expect(component.apisForm.total).toEqual(7)
    expect(component.apisForm.pageNum).toEqual(1)
    expect(component.apisForm.pageSize).toEqual(20)
    expect(component.apiTableLoading).toEqual(false)
    expect(spyGetApiData).not.toHaveBeenCalledTimes(1)
    discardPeriodicTasks()
  }))

  it('test apiSet', () => {
    expect(component).toBeTruthy()
    expect(component.apisSet.size).toEqual(0)
    expect(component.apisForm.apis.length).toEqual(0)
    // @ts-ignore
    jest.replaceProperty(fixture.debugElement.injector.get(BaseInfoService), '_allParams', {
      apiGroupId: 'mockApiGroupId'
    })

    component.ngOnInit()
    component.ngAfterViewInit()
    fixture.detectChanges()

    expect(component.apisSet.size).toEqual(0)
    expect(component.apisForm.apis.length).toEqual(7)

    component.changeApisSet(true, 'all')
    component.ngAfterViewInit()
    fixture.detectChanges()

    expect(component.apisSet.size).toEqual(7)

    component.changeApisSet(false, 'all')
    component.ngAfterViewInit()
    fixture.detectChanges()
    expect(component.apisSet.size).toEqual(0)

    component.changeApisSet({ uuid: MockApisList.data.apis[2].uuid })
    component.ngAfterViewInit()
    fixture.detectChanges()
    expect(component.apisSet.has(MockApisList.data.apis[2].uuid)).toEqual(true)
    expect(component.apisSet.size).toEqual(1)

    component.changeApisSet({ uuid: MockApisList.data.apis[3].uuid })
    component.ngAfterViewInit()
    fixture.detectChanges()
    expect(component.apisSet.has(MockApisList.data.apis[2].uuid)).toEqual(true)
    expect(component.apisSet.has(MockApisList.data.apis[3].uuid)).toEqual(true)
    expect(component.apisSet.size).toEqual(2)

    component.changeApisSet(false, 'all')
    component.ngAfterViewInit()
    fixture.detectChanges()
    expect(component.apisSet.size).toEqual(0)

    component.changeApisSet(true, 'all')
    component.ngAfterViewInit()
    fixture.detectChanges()

    expect(component.apisSet.size).toEqual(7)

    component.changeApisSet({ uuid: MockApisList.data.apis[5].uuid, checked: true })
    component.ngAfterViewInit()
    fixture.detectChanges()

    expect(component.apisSet.has(MockApisList.data.apis[2].uuid)).toEqual(true)
    expect(component.apisSet.has(MockApisList.data.apis[3].uuid)).toEqual(true)
    expect(component.apisSet.has(MockApisList.data.apis[5].uuid)).toEqual(false)
    expect(component.apisSet.size).toEqual(6)
  })

  it('test delete api', () => {
    // @ts-ignore
    const spyCreateModal = jest.spyOn(component.modalService, 'create')
    const spyGetApiData = jest.spyOn(component, 'getApisData')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'success')
    expect(component).toBeTruthy()
    expect(spyCreateModal).not.toHaveBeenCalled()
    expect(spyGetApiData).not.toHaveBeenCalled()
    expect(spyDeleteApiService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component.apisSet.size).toEqual(0)

    component.ngOnInit()
    fixture.detectChanges()

    expect(spyCreateModal).not.toHaveBeenCalled()
    expect(spyGetApiData).toHaveBeenCalledTimes(1)
    expect(component.apisForm.apis.length).not.toEqual(0)
    expect(spyDeleteApiService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component.apisSet.size).toEqual(0)

    component.deleteApiModal({ uuid: MockApisList.data.apis[0].uuid })
    fixture.detectChanges()

    expect(spyGetApiData).toHaveBeenCalledTimes(1)
    expect(spyCreateModal).toHaveBeenCalledTimes(1)
    expect(component.apisForm.apis.length).toEqual(MockApisList.data.apis.length)
    expect(spyDeleteApiService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component.apisSet.size).toEqual(0)

    component.changeApisSet(true, 'all')
    fixture.detectChanges()

    expect(spyGetApiData).toHaveBeenCalledTimes(1)
    expect(spyCreateModal).toHaveBeenCalledTimes(1)
    expect(component.apisForm.apis.length).toEqual(MockApisList.data.apis.length)
    expect(spyDeleteApiService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component.apisSet.size).toEqual(MockApisList.data.apis.length)

    component.deleteApi(MockApisList.data.apis[2])
    fixture.detectChanges()

    expect(spyGetApiData).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalledWith(MockEmptySuccessResponse.msg, { nzDuration: 1000 })
    expect(component.apisSet.size).toEqual(MockApisList.data.apis.length - 1)
    expect(spyCreateModal).toHaveBeenCalledTimes(1)
    expect(component.apisForm.apis.length).toEqual(MockApisList.data.apis.length)
    expect(spyDeleteApiService).toHaveBeenCalledTimes(1)
  })

  it('test apiList filter', () => {
    const spyGetApiData = jest.spyOn(component, 'getApisData')
    expect(component).toBeTruthy()
    expect(component.sourcesList).toEqual([])
    expect(spyGetApiData).not.toHaveBeenCalled()
    // @ts-ignore
    jest.replaceProperty(fixture.debugElement.injector.get(BaseInfoService), '_allParams', {
      apiGroupId: 'mockApiGroupId'
    })

    component.ngOnInit()
    component.ngAfterViewInit()
    fixture.detectChanges()

    expect(component.sourcesList).toEqual([
      { text: '自建', value: 'self-build:-1:' },
      { text: '导入', value: 'import:-1:' },
      { text: 'Apikit', value: 'sync:1:' },
      { text: 'Postcat', value: 'sync:2:' }
    ])
    expect(spyGetApiData).toHaveBeenCalledTimes(2)

    component.apisFilterChange({ col: { title: '来源' }, keys: [MockApiSource.data.list[2].id, MockApiSource.data.list[3].id] })
    fixture.detectChanges()

    expect(component.sourcesList[0]['byDefault']).toEqual(false)
    expect(component.sourcesList[1]['byDefault']).toEqual(false)
    expect(component.sourcesList[2]['byDefault']).toEqual(true)
    expect(component.sourcesList[3]['byDefault']).toEqual(true)
    expect(spyGetApiData).toHaveBeenCalledTimes(3)
  })

  it('test addApi', () => {
    const spyRouterChange = jest.spyOn(component.router, 'navigate')
    expect(spyRouterChange).not.toHaveBeenCalled()
    expect(component).toBeTruthy()

    component.apisForm.groupUuid = ''
    component.addApi()
    fixture.detectChanges()

    expect(spyRouterChange).toHaveBeenCalledWith(['/', 'router', 'api', 'create'])

    component.apisForm.groupUuid = 'test-ws'
    component.addApi('websocket')
    fixture.detectChanges()

    expect(spyRouterChange).toHaveBeenCalledWith(['/', 'router', 'api', 'create-ws', 'test-ws'])

    component.apisForm.groupUuid = 'http'
    component.addApi()
    fixture.detectChanges()

    expect(spyRouterChange).toHaveBeenCalledWith(['/', 'router', 'api', 'create', 'http'])

    component.apisForm.groupUuid = ''
    component.addApi('websocket')
    fixture.detectChanges()

    expect(spyRouterChange).toHaveBeenCalledWith(['/', 'router', 'api', 'create-ws'])
  })

  it('test batchPublish ', () => {
    // @ts-ignore
    const spyBatchPublishModal = jest.spyOn(component.service, 'batchPublishApiModal')
    expect(spyBatchPublishModal).not.toHaveBeenCalled()
    expect(component).toBeTruthy()

    component.apisSet = new Set()
    component.batchPublish('online')
    fixture.detectChanges()

    expect(spyBatchPublishModal).toHaveBeenCalledTimes(1)

    component.apisSet = new Set(['123', '223', '333'])
    component.batchPublish('offline')
    fixture.detectChanges()

    expect(spyBatchPublishModal).toHaveBeenCalledTimes(2)

    component.apisSet = new Set(['4'])
    component.batchPublish('online')
    fixture.detectChanges()

    expect(spyBatchPublishModal).toHaveBeenCalledTimes(3)

    component.apisSet = new Set()
    component.batchPublish('offline')
    fixture.detectChanges()

    expect(spyBatchPublishModal).toHaveBeenCalledTimes(4)
  })
})
