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
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { ApiManagementProxyComponent } from '../../proxy/proxy.component'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

describe('#init ApiWebsocketCreateComponent', () => {
  let component:ApiWebsocketCreateComponent
  let fixture: ComponentFixture<ApiWebsocketCreateComponent>
  let spyGetApiMessage:jest.SpyInstance<any>
  let spyOpenDrawer:jest.SpyInstance<any>
  let httpCommonService:any
  let spyPutApiService:jest.SpyInstance<any>
  let spyPostApiService:jest.SpyInstance<any>
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  let spyApiService:jest.SpyInstance<any>
  global.structuredClone = (val:any) => JSON.parse(JSON.stringify(val))

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule, FormsModule, ReactiveFormsModule, ComponentModule, LayoutModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule, NzOutletModule, HttpClientModule,
        RouterModule.forRoot(routes), NzFormModule, EoNgInputModule, EoNgTreeModule, EoNgButtonModule,
        EoNgSwitchModule, EoNgCheckboxModule, EoNgApintoTableModule, EoNgSelectModule, EoNgFeedbackModalModule,
        EoNgFeedbackTooltipModule
      ],
      declarations: [ApiWebsocketCreateComponent
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
      ],
      teardown: { destroyAfterEach: false }
    }).compileComponents()

    fixture = TestBed.createComponent(ApiWebsocketCreateComponent)
    component = fixture.componentInstance

    fixture.detectChanges()
    spyGetApiMessage = jest.spyOn(component, 'getApiMessage')
    spyOpenDrawer = jest.spyOn(component, 'openDrawer')

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

  it('test proxyModal when create a new api', fakeAsync(() => {
    const fixtureProxy: ComponentFixture<ApiManagementProxyComponent> = TestBed.createComponent(ApiManagementProxyComponent)
    const componentProxy:ApiManagementProxyComponent = fixtureProxy.componentInstance

    expect(component).toBeTruthy()
    expect(component.validateForm.controls['groupUuid'].value).toEqual('')
    expect(component.createApiForm.proxyHeader).toEqual([])
    expect(component.editData).toBeNull()
    expect(component.modalRef).toBeUndefined()
    expect(component.proxyEdit).toEqual(false)

    component.openDrawer('proxyHeader')
    fixture.detectChanges()
    tick()

    expect(component.editData).toBeNull()
    expect(component.modalRef).not.toBeUndefined()
    expect(component.proxyEdit).toEqual(false)

    component.proxyTableClick({ data: { test: 1 } })
    fixture.detectChanges()
    tick()

    expect(component.editData).not.toBeNull()
    expect(component.proxyEdit).toEqual(false)

    // @ts-ignore
    expect(component.modalRef?.nzOnOk(componentProxy)).toEqual(false)
  }))

  it('test init checkbox', () => {
    expect(component).toBeTruthy()
    component.ngOnInit()
    fixture.detectChanges()

    expect(component.createApiForm.method).toEqual([])
    expect(component.methodList.length).toEqual(6) // GET, POST, PUT, DELETE, PATCH, HEAD

    component.initCheckbox()
    fixture.detectChanges()
    expect(component.createApiForm.method).toEqual([])
    for (const m of component.methodList) {
      expect(m.checked).toEqual(false)
    }
    component.createApiForm.method = ['PUT']

    component.initCheckbox()
    fixture.detectChanges()

    for (const m of component.methodList) {
      if (m.label !== 'PUT') {
        expect(m.checked).toEqual(false)
      } else {
        expect(m.checked).toEqual(true)
      }
    }
  })

  it('test update all checkbox', () => {
    expect(component).toBeTruthy()
    component.ngOnInit()
    fixture.detectChanges()

    expect(component.createApiForm.method).toEqual([])
    expect(component.methodList.length).toEqual(6) // GET, POST, PUT, DELETE, PATCH, HEAD

    component.initCheckbox()
    fixture.detectChanges()
    expect(component.createApiForm.method).toEqual([])
    for (const m of component.methodList) {
      expect(m.checked).toEqual(false)
    }
    expect(component.allChecked).toEqual(false)

    component.updateAllChecked()
    fixture.detectChanges()

    for (const m of component.methodList) {
      expect(m.checked).toEqual(false)
    }
    expect(component.createApiForm.method).toEqual([])
    expect(component.showCheckboxGroupValid).toEqual(true)

    component.allChecked = true
    component.updateAllChecked()
    fixture.detectChanges()

    for (const m of component.methodList) {
      expect(m.checked).toEqual(true)
    }
    expect(component.createApiForm.method.length).toEqual(component.methodList.length)
    expect(component.showCheckboxGroupValid).toEqual(false)

    component.allChecked = false
    component.updateAllChecked()
    fixture.detectChanges()

    for (const m of component.methodList) {
      expect(m.checked).toEqual(false)
    }
    expect(component.createApiForm.method).toEqual([])
    expect(component.showCheckboxGroupValid).toEqual(true)
  })

  it('test update single checkbox', () => {
    expect(component).toBeTruthy()
    component.ngOnInit()
    fixture.detectChanges()

    expect(component.createApiForm.method).toEqual([])
    expect(component.methodList.length).toEqual(6) // GET, POST, PUT, DELETE, PATCH, HEAD

    component.initCheckbox()
    fixture.detectChanges()
    expect(component.createApiForm.method).toEqual([])
    for (const m of component.methodList) {
      expect(m.checked).toEqual(false)
    }
    expect(component.allChecked).toEqual(false)

    component.updateSingleChecked()
    fixture.detectChanges()

    for (const m of component.methodList) {
      expect(m.checked).toEqual(false)
    }
    expect(component.allChecked).toEqual(false)
    expect(component.createApiForm.method).toEqual([])
    expect(component.showCheckboxGroupValid).toEqual(true)

    component.methodList[2].checked = true
    component.updateSingleChecked()
    fixture.detectChanges()

    for (const m of component.methodList) {
      if (component.methodList.indexOf(m) !== 2) {
        expect(m.checked).toEqual(false)
      } else {
        expect(m.checked).toEqual(true)
      }
    }
    expect(component.allChecked).toEqual(false)
    expect(component.createApiForm.method.length).toEqual(1)
    expect(component.showCheckboxGroupValid).toEqual(false)

    for (const m of component.methodList) {
      m.checked = true
    }
    component.updateSingleChecked()
    fixture.detectChanges()

    for (const m of component.methodList) {
      expect(m.checked).toEqual(true)
    }
    expect(component.allChecked).toEqual(true)
    expect(component.createApiForm.method.length).toEqual(component.methodList.length)
    expect(component.showCheckboxGroupValid).toEqual(false)
  })

  it('test default proxyPath', () => {
    expect(component).toBeTruthy()
    component.ngOnInit()
    fixture.detectChanges()

    component.validateForm.controls['requestPath'].setValue('requestPath')
    expect(!component.validateForm.controls['proxyPath'].value).toEqual(true)
    expect(!!component.validateForm.controls['requestPath'].value).toEqual(true)

    component.requestPathChange()
    fixture.detectChanges()

    expect(!!component.validateForm.controls['proxyPath'].value).toEqual(true)
    expect(component.validateForm.controls['proxyPath'].value).toEqual('/requestPath')

    component.validateForm.controls['requestPath'].setValue('requestPath2')

    component.requestPathChange()
    fixture.detectChanges()

    expect(component.validateForm.controls['proxyPath'].value).toEqual('/requestPath')
  })

  it('test timeout', () => {
    expect(component).toBeTruthy()
    component.ngOnInit()
    component.checkTimeout()
    fixture.detectChanges()
    expect(component.validateForm.controls['timeout'].value).toEqual(10000)

    component.validateForm.controls['timeout'].setValue('')
    component.checkTimeout()
    fixture.detectChanges()
    expect(component.validateForm.controls['timeout'].value).toEqual('')

    component.validateForm.controls['timeout'].setValue(-5)
    component.checkTimeout()
    fixture.detectChanges()
    expect(component.validateForm.controls['timeout'].value).toEqual(1)

    component.validateForm.controls['timeout'].setValue(2)
    component.checkTimeout()
    fixture.detectChanges()
    expect(component.validateForm.controls['timeout'].value).toEqual(2)
  })

  it('test save proxy header', () => {
    const fixtureProxy: ComponentFixture<ApiManagementProxyComponent> = TestBed.createComponent(ApiManagementProxyComponent)
    const componentProxy:ApiManagementProxyComponent = fixtureProxy.componentInstance
    const spyProxyKeyDirty = jest.spyOn(componentProxy.validateProxyHeaderForm.controls['key'], 'markAsDirty')
    expect(component).toBeTruthy()
    component.ngOnInit()
    fixture.detectChanges()
    expect(component.proxyEdit).toEqual(false)
    expect(spyProxyKeyDirty).not.toHaveBeenCalled()
    componentProxy.validateProxyHeaderForm.controls['optType'].setValue('ADD')
    componentProxy.validateProxyHeaderForm.controls['key'].setValue('')
    componentProxy.validateProxyHeaderForm.controls['value'].setValue('')

    component.saveProxyHeader(componentProxy)
    expect(spyProxyKeyDirty).toHaveBeenCalledTimes(1)

    componentProxy.validateProxyHeaderForm.controls['optType'].setValue('DELETE')
    componentProxy.validateProxyHeaderForm.controls['key'].setValue('')

    component.saveProxyHeader(componentProxy)
    expect(spyProxyKeyDirty).toHaveBeenCalledTimes(2)

    componentProxy.validateProxyHeaderForm.controls['optType'].setValue('ADD')
    componentProxy.validateProxyHeaderForm.controls['key'].setValue('KEY')
    componentProxy.validateProxyHeaderForm.controls['value'].setValue('VALUE')

    component.saveProxyHeader(componentProxy)
    expect(component.createApiForm.proxyHeader).toEqual([{ optType: 'ADD', key: 'KEY', value: 'VALUE' }])
    expect(spyProxyKeyDirty).toHaveBeenCalledTimes(2)

    componentProxy.validateProxyHeaderForm.controls['optType'].setValue('ADD')
    componentProxy.validateProxyHeaderForm.controls['key'].setValue('KEY2')
    componentProxy.validateProxyHeaderForm.controls['value'].setValue('VALUE2')

    component.saveProxyHeader(componentProxy)
    expect(component.createApiForm.proxyHeader).toEqual([{ optType: 'ADD', key: 'KEY2', value: 'VALUE2' }, { optType: 'ADD', key: 'KEY', value: 'VALUE' }])
    expect(spyProxyKeyDirty).toHaveBeenCalledTimes(2)

    component.proxyEdit = true
    component.editData = { optType: 'ADD', key: 'KEY2', value: 'VALUE2' }
    componentProxy.validateProxyHeaderForm.controls['optType'].setValue('ADD')
    componentProxy.validateProxyHeaderForm.controls['key'].setValue('KEY22')
    componentProxy.validateProxyHeaderForm.controls['value'].setValue('VALUE22')

    component.saveProxyHeader(componentProxy)
    expect(component.createApiForm.proxyHeader).toEqual([{ optType: 'ADD', key: 'KEY22', value: 'VALUE22' }, { optType: 'ADD', key: 'KEY', value: 'VALUE' }])
    expect(spyProxyKeyDirty).toHaveBeenCalledTimes(2)

    component.proxyEdit = false
    component.editData = {}
    componentProxy.validateProxyHeaderForm.controls['optType'].setValue('DELETE')
    componentProxy.validateProxyHeaderForm.controls['key'].setValue('')
    componentProxy.validateProxyHeaderForm.controls['value'].setValue('')

    component.saveProxyHeader(componentProxy)
    expect(component.createApiForm.proxyHeader).toEqual([{ optType: 'ADD', key: 'KEY22', value: 'VALUE22' }, { optType: 'ADD', key: 'KEY', value: 'VALUE' }])
    expect(spyProxyKeyDirty).toHaveBeenCalledTimes(3)

    component.proxyEdit = false
    component.editData = {}
    componentProxy.validateProxyHeaderForm.controls['optType'].setValue('DELETE')
    componentProxy.validateProxyHeaderForm.controls['key'].setValue('DELETE1')
    componentProxy.validateProxyHeaderForm.controls['value'].setValue('')

    component.saveProxyHeader(componentProxy)
    expect(component.createApiForm.proxyHeader).toEqual([{ optType: 'DELETE', key: 'DELETE1', value: '' }, { optType: 'ADD', key: 'KEY22', value: 'VALUE22' }, { optType: 'ADD', key: 'KEY', value: 'VALUE' }])
    expect(spyProxyKeyDirty).toHaveBeenCalledTimes(3)

    component.proxyEdit = true
    component.editData = { optType: 'DELETE', key: 'DELETE1', value: '' }
    componentProxy.validateProxyHeaderForm.controls['optType'].setValue('DELETE')
    componentProxy.validateProxyHeaderForm.controls['key'].setValue('DELETE12')
    componentProxy.validateProxyHeaderForm.controls['value'].setValue('')

    component.saveProxyHeader(componentProxy)
    expect(component.createApiForm.proxyHeader).toEqual([{ optType: 'DELETE', key: 'DELETE12', value: '' }, { optType: 'ADD', key: 'KEY22', value: 'VALUE22' }, { optType: 'ADD', key: 'KEY', value: 'VALUE' }])
    expect(spyProxyKeyDirty).toHaveBeenCalledTimes(3)

    component.proxyEdit = true
    component.editData = { optType: 'DELETE', key: 'DELETE12', value: '' }
    componentProxy.validateProxyHeaderForm.controls['optType'].setValue('ADD')
    componentProxy.validateProxyHeaderForm.controls['key'].setValue('KEY222')
    componentProxy.validateProxyHeaderForm.controls['value'].setValue('VALUE222')

    component.saveProxyHeader(componentProxy)
    expect(component.createApiForm.proxyHeader).toEqual([{ optType: 'ADD', key: 'KEY222', value: 'VALUE222' }, { optType: 'ADD', key: 'KEY22', value: 'VALUE22' }, { optType: 'ADD', key: 'KEY', value: 'VALUE' }])
    expect(spyProxyKeyDirty).toHaveBeenCalledTimes(3)

    component.proxyEdit = true
    component.editData = { optType: 'ADD', key: 'KEY', value: 'VALUE' }
    componentProxy.validateProxyHeaderForm.controls['optType'].setValue('DELETE')
    componentProxy.validateProxyHeaderForm.controls['key'].setValue('KEY2222')
    componentProxy.validateProxyHeaderForm.controls['value'].setValue('')

    component.saveProxyHeader(componentProxy)

    expect(component.createApiForm.proxyHeader).toEqual([{ optType: 'DELETE', key: 'KEY2222', value: '' }, { optType: 'ADD', key: 'KEY222', value: 'VALUE222' }, { optType: 'ADD', key: 'KEY22', value: 'VALUE22' }])
    expect(spyProxyKeyDirty).toHaveBeenCalledTimes(3)
  })

  it('test save api', () => {
    expect(component).toBeTruthy()
    component.ngOnInit()
    component.checkTimeout()
    fixture.detectChanges()
    expect(component.validateForm.controls['timeout'].value).toEqual(10000)

    component.validateForm.controls['timeout'].setValue('')
    component.checkTimeout()
    fixture.detectChanges()
    expect(component.validateForm.controls['timeout'].value).toEqual('')

    component.validateForm.controls['timeout'].setValue(-5)
    component.checkTimeout()
    fixture.detectChanges()
    expect(component.validateForm.controls['timeout'].value).toEqual(1)

    component.validateForm.controls['timeout'].setValue(2)
    component.checkTimeout()
    fixture.detectChanges()
    expect(component.validateForm.controls['timeout'].value).toEqual(2)
  })

  it('test save proxy header as create a new api', () => {
    spyPutApiService = jest.spyOn(httpCommonService, 'put').mockReturnValue(
      of(MockEmptySuccessResponse)
    )

    spyPostApiService = jest.spyOn(httpCommonService, 'post').mockReturnValue(
      of(MockEmptySuccessResponse)
    )
    const spyFormServiceDirty = jest.spyOn(component.validateForm.controls['service'], 'markAsDirty')
    expect(component).toBeTruthy()
    expect(spyFormServiceDirty).not.toHaveBeenCalled()
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.ngOnInit()
    component.saveApi('websocket')
    fixture.detectChanges()

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(1)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)
    component.showCheckboxGroupValid = true
    component.saveApi('websocket')
    fixture.detectChanges()

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(2)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.validateForm.patchValue({
      groupUuid: 'test'
    })

    component.saveApi('websocket')
    fixture.detectChanges()

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(3)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.validateForm.patchValue({
      name: 'test'
    })

    component.saveApi('websocket')
    fixture.detectChanges()

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(4)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.validateForm.patchValue({
      requestPath: 'test'
    })

    component.saveApi('websocket')
    fixture.detectChanges()

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(5)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.validateForm.patchValue({
      service: 'test'
    })

    component.saveApi('websocket')
    fixture.detectChanges()

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(5)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.showCheckboxGroupValid = false
    component.saveApi('websocket')
    fixture.detectChanges()

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(5)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).toHaveBeenCalledTimes(1)
    expect(component.createApiForm.method).toEqual([]) // websocket
    expect(component.submitButtonLoading).toEqual(false)
  })

  it('test save proxy header as edit  api', () => {
    httpCommonService = fixture.debugElement.injector.get(ApiService)
    spyPutApiService = jest.spyOn(httpCommonService, 'put').mockReturnValue(
      of(MockEmptySuccessResponse)
    )

    spyPostApiService = jest.spyOn(httpCommonService, 'post').mockReturnValue(
      of(MockEmptySuccessResponse)
    )
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

    const spyFormServiceDirty = jest.spyOn(component.validateForm.controls['service'], 'markAsDirty')
    expect(component).toBeTruthy()
    component.editPage = true

    expect(spyFormServiceDirty).not.toHaveBeenCalled()
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.ngOnInit()
    component.saveApi('websocket')
    fixture.detectChanges()

    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(spyPutApiService).toHaveBeenCalledTimes(1)
    expect(component.createApiForm.method).toEqual([]) // websocket
    expect(component.submitButtonLoading).toEqual(false)

    component.showCheckboxGroupValid = true
    component.saveApi('websocket')
    fixture.detectChanges()

    expect(spyFormServiceDirty).not.toHaveBeenCalled()
    expect(spyPutApiService).toHaveBeenCalledTimes(1)
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)
  })
})
