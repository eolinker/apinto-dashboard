import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing'
import { ApiWebsocketCreateComponent } from './websocket-create.component'
import { ComponentModule } from 'projects/core/src/app/component/component.module'
import { APP_BASE_HREF } from '@angular/common'
import { HttpClientModule } from '@angular/common/http'
import { ElementRef, Renderer2, ChangeDetectorRef } from '@angular/core'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { RouterModule } from '@angular/router'
import { Overlay } from '@angular/cdk/overlay'
import { EoNgFeedbackMessageService, EoNgFeedbackModalModule, EoNgFeedbackModalService, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { LayoutModule } from '../../../../layout.module'
import { BidiModule } from '@angular/cdk/bidi'
import { routes } from '../../../api-routing.module'
import { BasicLayoutComponent } from '../../../../basic-layout/basic-layout.component'
import { MockRenderer, MockMessageService, MockEnsureService, MockEmptySuccessResponse, MockGetCommonProviderService, MockPluginTemplateEnum, MockRouterGroups, MockAccessList, MockModuleList, MockApiWsMessage, MockApiWsMessage2 } from 'projects/core/src/app/constant/spec-test'
import { of } from 'rxjs'
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
import { MatchTableComponent } from '../../match/table/table.component'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

class MockBaseInfoService {
  allParam = {

  }

  get allParamsInfo () {
    return { ...this.allParam }
  }
}

describe('#init ApiWebsocketCreateComponent', () => {
  let component:ApiWebsocketCreateComponent
  let fixture: ComponentFixture<ApiWebsocketCreateComponent>
  let fixtureBasic: ComponentFixture<BasicLayoutComponent>
  let componentBasic:BasicLayoutComponent
  let spyGetApiMessage:jest.SpyInstance<any>
  let spyGetHeaderList:jest.SpyInstance<any>
  let spyTransferHeader:jest.SpyInstance<any>
  let spyNzTreeClick:jest.SpyInstance<any>
  let spyGetServiceList:jest.SpyInstance<any>
  let spyGetPluginTemplateList:jest.SpyInstance<any>
  let spyUpdateAllChecked:jest.SpyInstance<any>
  let spyInitCheckbox:jest.SpyInstance<any>
  let spyUpdateSingleChecked:jest.SpyInstance<any>
  let spyProxyTableClick:jest.SpyInstance<any>
  let spyOpenDrawer:jest.SpyInstance<any>
  let spyBackToList:jest.SpyInstance<any>
  let spySaveApi:jest.SpyInstance<any>
  let spyRequestPathChange:jest.SpyInstance<any>
  let spyCheckTimeout:jest.SpyInstance<any>
  let spySaveProxyHeader:jest.SpyInstance<any>
  let spyMessage:jest.SpyInstance<any>
  let spyNavService:jest.SpyInstance<any>
  let httpCommonService:any
  let spyApiService:jest.SpyInstance<any>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule, FormsModule, ReactiveFormsModule, ComponentModule, LayoutModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule, NzOutletModule, HttpClientModule,
        RouterModule.forRoot(routes), NzFormModule, EoNgInputModule, EoNgTreeModule, EoNgButtonModule,
        EoNgSwitchModule, EoNgCheckboxModule, EoNgApintoTableModule, EoNgSelectModule, EoNgFeedbackModalModule,
        EoNgFeedbackTooltipModule
      ],
      declarations: [ApiWebsocketCreateComponent, MatchTableComponent
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

    fixture = TestBed.createComponent(ApiWebsocketCreateComponent)
    component = fixture.componentInstance

    fixture.detectChanges()
    spyGetApiMessage = jest.spyOn(component, 'getApiMessage')
    spyGetHeaderList = jest.spyOn(component, 'getHeaderList')
    spyTransferHeader = jest.spyOn(component, 'transferHeader')
    spyGetServiceList = jest.spyOn(component, 'getServiceList')
    spyGetPluginTemplateList = jest.spyOn(component, 'getPluginTemplateList')
    spyUpdateAllChecked = jest.spyOn(component, 'updateAllChecked')
    spyUpdateSingleChecked = jest.spyOn(component, 'updateSingleChecked')
    spyProxyTableClick = jest.spyOn(component, 'proxyTableClick')
    spyOpenDrawer = jest.spyOn(component, 'openDrawer')
    spyBackToList = jest.spyOn(component, 'backToList')
    spySaveApi = jest.spyOn(component, 'saveApi')
    spyRequestPathChange = jest.spyOn(component, 'requestPathChange')
    spyCheckTimeout = jest.spyOn(component, 'checkTimeout')
    spySaveProxyHeader = jest.spyOn(component, 'saveProxyHeader')
    spyMessage = jest.spyOn(component.message, 'success')

    httpCommonService = fixture.debugElement.injector.get(ApiService)
    spyApiService = jest.spyOn(httpCommonService, 'get').mockImplementation(
      (...args) => {
        switch (args[0]) {
          case 'common/provider/Service':
            return of(MockGetCommonProviderService)
          case 'plugin/template/enum':
            return of(MockPluginTemplateEnum)
          case 'router/groups':
            return of(MockRouterGroups)
          case 'my/access':
            return of(MockAccessList)
          case 'system/modules':
            return of(MockModuleList)
          case 'router':
            return of(MockApiWsMessage)
          default:
            return of(MockEmptySuccessResponse)
        }
      }
    )
  })

  it('should create and init component as creating a new api', fakeAsync(() => {
    expect(component).toBeTruthy()
    expect(component.serviceList).toEqual([])
    expect(component.pluginTemplateList).toEqual([])
    expect(component.groupUuid).toEqual('')
    expect(component.editPage).toEqual(false)
    expect(component.headerList).toEqual([])
    expect(component.validateForm.controls['groupUuid'].value).toEqual('')
    expect(component.nzDisabled).toEqual(false)
    // @ts-ignore
    jest.replaceProperty(fixture.debugElement.injector.get(BaseInfoService), '_allParams', {
      apiGroupId: 'mockApiGroupId',
      apiId: 'mockApiId'
    })
    component.ngOnInit()
    fixture.detectChanges()
    tick()
    expect(component.editPage).toEqual(false)
    expect(spyGetApiMessage).not.toHaveBeenCalled()
    expect(component.serviceList).toEqual([
      { label: 'test1[http]', value: 'test1cjk_service@upstream' },
      { label: 'test2[http]', value: 'test2@upstream' }
    ])
    expect(component.pluginTemplateList).toEqual([
      { label: 'test1', value: '70623690-430f-23db-ec75-763fe7c380d9' },
      { label: 'test2', value: 'aa0be463-43c8-4f67-2633-ea6bdcea9709' }
    ])
    // @ts-ignore
    expect(component.baseInfo.allParamsInfo.apiGroupId).toEqual('mockApiGroupId')
    // @ts-ignore
    expect(component.baseInfo.allParamsInfo.apiId).toEqual('mockApiId')
    expect(component.groupUuid).toEqual('mockApiGroupId')
    expect(component.editPage).toEqual(false)
    expect(component.headerList).toEqual([
      {
        title: 'test1',
        key: '50458642-5a9f-4136-9ff1-e30d647297e8',
        uuid: '50458642-5a9f-4136-9ff1-e30d647297e8',
        selected: false,
        isDelete: false,
        children: [
          {
            title: 'test1-c1',
            key: '35938ae4-1a62-4e22-ad8c-3691e111820e',
            uuid: '35938ae4-1a62-4e22-ad8c-3691e111820e',
            isDelete: false,
            selected: false,
            isLeaf: true
          },
          {
            title: 'test1-c2',
            key: 'b238751a-dbfb-4610-8f40-a599737ac4e5',
            uuid: 'b238751a-dbfb-4610-8f40-a599737ac4e5',
            isDelete: false,
            selected: false,
            isLeaf: true
          }
        ]
      },
      {
        title: 'test2',
        key: '00db4977-331f-4b7e-93be-b64648751a5f',
        uuid: '00db4977-331f-4b7e-93be-b64648751a5f',
        isDelete: false,
        selected: false,
        isLeaf: true
      }
    ])
    expect(component.validateForm.controls['groupUuid'].value).toEqual('mockApiGroupId')
    expect(component.nzDisabled).toEqual(false)
  }))

  it('should init proxyHeaderTable and hostsTable as creating a new api', fakeAsync(() => {
    expect(component).toBeTruthy()
    expect(component.editPage).toEqual(false)
    expect(component.nzDisabled).toEqual(false)
    expect(spyOpenDrawer).toHaveBeenCalledTimes(0)
    // @ts-ignore
    jest.replaceProperty(fixture.debugElement.injector.get(BaseInfoService), '_allParams', {
      apiGroupId: 'mockApiGroupId',
      apiId: 'mockApiId'
    })
    component.ngOnInit()
    fixture.detectChanges()
    tick()
    expect(component.editPage).toEqual(false)
    expect(component.nzDisabled).toEqual(false)
    expect(component.proxyHeaderTableBody[3].btns[0].disabledFn()).toEqual(false)
    expect(component.proxyHeaderTableBody[3].btns[1].disabledFn()).toEqual(false)
    expect(component.proxyHeaderTableBody[3].btns[0].click).not.toBeUndefined()
    // @ts-ignore
    expect(component.hostsTableBody[0].disabledFn()).toEqual(false)
    expect(component.hostsTableBody[1].btns[0].disabledFn()).toEqual(false)
    expect(component.hostsTableBody[1].showFn).not.toBeUndefined()
    expect(component.hostsTableBody[2].showFn).not.toBeUndefined()
    expect(component.hostsTableBody[2].btns[0].disabledFn()).toEqual(false)
    expect(component.hostsTableBody[2].btns[1].disabledFn()).toEqual(false)

    component.ngAfterViewInit()
    fixture.detectChanges()
    expect(component.proxyHeaderTableBody[0].title).not.toBeUndefined()

    component.proxyHeaderTableBody[3].btns[0].click({ data: { test: 1 } })
    fixture.detectChanges()
    expect(spyOpenDrawer).toHaveBeenCalledTimes(1)
  }))

  it('should create and init component as creating a new api', () => {
    expect(component).toBeTruthy()
    expect(component.validateForm.controls['groupUuid'].value).toEqual('')
  })

  it('should create and init component as editing an api', fakeAsync(() => {
    expect(component).toBeTruthy()
    component.editPage = true
    expect(component.validateForm.controls['groupUuid'].value).toEqual('')
    expect(component.validateForm.controls['requestPath'].value).toEqual('')
    component.ngOnInit()
    component.ngAfterViewInit()
    fixture.detectChanges()
    tick()

    expect(component.validateForm.controls['groupUuid'].value).toEqual('50458642-5a9f-4136-9ff1-e30d647297e8')
    expect(component.validateForm.controls['name'].value).toEqual('ss')
    expect(component.validateForm.controls['desc'].value).toEqual('')
    expect(component.validateForm.controls['isDisable'].value).toEqual(false)
    expect(component.validateForm.controls['requestPath'].value).toEqual('tetetetete')
    expect(component.validateForm.controls['service'].value).toEqual('testService')
    expect(component.validateForm.controls['proxyPath'].value).toEqual('tetetetete')
    expect(component.validateForm.controls['timeout'].value).toEqual(10000)
    expect(component.validateForm.controls['retry'].value).toEqual(0)
    expect(component.validateForm.controls['templateUuid'].value).toEqual('e59693df-27cc-61a7-d0f0-c17da203026a')
    expect(component.createApiForm.uuid).toEqual('569c8d47-d742-5306-c0e2-a5ae38727fa7')
    expect(component.createApiForm.method).toEqual([])
    expect(component.createApiForm.match).toEqual(MockApiWsMessage.data.api.match)
    expect(component.createApiForm.proxyHeader).toEqual(MockApiWsMessage.data.api.proxyHeader)
    expect(component.hostsList[0].key).toEqual('test1.host.addr')
    expect(component.hostsList[1].key).toEqual('test2.host.addr')
    expect(component.hostsList[2].key).toEqual('')
    expect(component.hostsList.length).toEqual(3)
    expect(component.headerList).toEqual([
      {
        title: 'test1',
        key: '50458642-5a9f-4136-9ff1-e30d647297e8',
        uuid: '50458642-5a9f-4136-9ff1-e30d647297e8',
        selected: true,
        isDelete: false,
        children: [
          {
            title: 'test1-c1',
            key: '35938ae4-1a62-4e22-ad8c-3691e111820e',
            uuid: '35938ae4-1a62-4e22-ad8c-3691e111820e',
            isDelete: false,
            isLeaf: true
          },
          {
            title: 'test1-c2',
            key: 'b238751a-dbfb-4610-8f40-a599737ac4e5',
            uuid: 'b238751a-dbfb-4610-8f40-a599737ac4e5',
            isDelete: false,
            isLeaf: true
          }
        ]
      },
      {
        title: 'test2',
        key: '00db4977-331f-4b7e-93be-b64648751a5f',
        uuid: '00db4977-331f-4b7e-93be-b64648751a5f',
        isDelete: false,
        isLeaf: true
      }
    ])
  }))

  it('requestPath should only remove the first letter / && hosts always not empty', fakeAsync(() => {
    expect(component).toBeTruthy()
    component.editPage = true
    expect(component.validateForm.controls['groupUuid'].value).toEqual('')
    expect(component.validateForm.controls['requestPath'].value).toEqual('')
    spyApiService = jest.spyOn(httpCommonService, 'get').mockImplementation(
      (...args) => {
        switch (args[0]) {
          case 'common/provider/Service':
            return of(MockGetCommonProviderService)
          case 'plugin/template/enum':
            return of(MockPluginTemplateEnum)
          case 'router/groups':
            return of(MockRouterGroups)
          case 'my/access':
            return of(MockAccessList)
          case 'system/modules':
            return of(MockModuleList)
          case 'router':
            return of(MockApiWsMessage2)
          default:
            return of(MockEmptySuccessResponse)
        }
      })
    component.ngOnInit()
    component.ngAfterViewInit()
    tick()

    expect(component.validateForm.controls['groupUuid'].value).toEqual('50458642-5a9f-4136-9ff1-e30d647297e8')
    expect(component.validateForm.controls['name'].value).toEqual('ss')
    expect(component.validateForm.controls['desc'].value).toEqual('')
    expect(component.validateForm.controls['isDisable'].value).toEqual(false)
    expect(component.validateForm.controls['requestPath'].value).toEqual('{{baseUrl}}/test')
    expect(component.validateForm.controls['service'].value).toEqual('testService')
    expect(component.validateForm.controls['proxyPath'].value).toEqual('tetetetete')
    expect(component.validateForm.controls['timeout'].value).toEqual(10000)
    expect(component.validateForm.controls['retry'].value).toEqual(0)
    expect(component.validateForm.controls['templateUuid'].value).toEqual('e59693df-27cc-61a7-d0f0-c17da203026a')
    expect(component.createApiForm.uuid).toEqual('569c8d47-d742-5306-c0e2-a5ae38727fa7')
    expect(component.createApiForm.method).toEqual([])
    expect(component.createApiForm.match).toEqual(MockApiWsMessage2.data.api.match)
    expect(component.createApiForm.proxyHeader).toEqual(MockApiWsMessage2.data.api.proxyHeader)
    expect(component.hostsList).toEqual([{ key: '' }])
  }))

  it('test proxyModal when create a new api', () => {
    expect(component).toBeTruthy()
    expect(component.validateForm.controls['groupUuid'].value).toEqual('')
    expect(component.createApiForm.proxyHeader).toEqual([])
    expect(component.editData).toBeNull()
    expect(component.modalRef).toBeUndefined()
    expect(component.proxyEdit).toEqual(false)

    component.openDrawer('proxyHeader')
    fixture.detectChanges()

    expect(component.editData).toBeNull()
    expect(component.modalRef).not.toBeUndefined()
    expect(component.proxyEdit).toEqual(false)

    component.proxyTableClick({ data: { test: 1 } })
    fixture.detectChanges()

    expect(component.editData).not.toBeNull()
    expect(component.proxyEdit).toEqual(true)

    component.modalRef?.close()
    fixture.detectChanges()

    expect(component.proxyEdit).toEqual(false)
  })

  it('test checkbox when create a new api', () => {
    expect(component).toBeTruthy()
    component.ngOnInit()
    fixture.detectChanges()

    component.initCheckbox()
    fixture.detectChanges()
    // expect()
  })
})
