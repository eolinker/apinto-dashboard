import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing'
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
import { MockRenderer, MockMessageService, MockEnsureService, MockEmptySuccessResponse, MockGetCommonProviderService, MockPluginTemplateEnum, MockRouterGroups, MockAccessList, MockModuleList, MockApiHttpMessage, MockApiHttpMessage2 } from 'projects/core/src/app/constant/spec-test'
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
import { ApiHttpCreateComponent } from './http-create.component'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

describe('#init ApiHttpCreateComponent', () => {
  let component:ApiHttpCreateComponent
  let fixture: ComponentFixture<ApiHttpCreateComponent>
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
      declarations: [ApiHttpCreateComponent
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

    fixture = TestBed.createComponent(ApiHttpCreateComponent)
    component = fixture.componentInstance

    fixture.detectChanges()

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
            return of(MockApiHttpMessage)
          default:
            return of(MockEmptySuccessResponse)
        }
      }
    )
  })

  it('should create and init component as editing an api with Get method', fakeAsync(() => {
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
    expect(component.createApiForm.method).toEqual(['GET'])
    expect(component.createApiForm.match).toEqual(MockApiHttpMessage.data.api.match)
    expect(component.createApiForm.proxyHeader).toEqual(MockApiHttpMessage.data.api.proxyHeader)
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
    for (const m of component.methodList) {
      if (m.value !== 'GET') {
        expect(m.checked).toEqual(false)
      } else {
        expect(m.checked).toEqual(true)
      }
    }
    expect(component.allChecked).toEqual(false)
  }))

  it('should create and init component as editing an api with All method', fakeAsync(() => {
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
            return of(MockApiHttpMessage2)
          default:
            return of(MockEmptySuccessResponse)
        }
      })
    component.ngOnInit()
    component.ngAfterViewInit()
    fixture.detectChanges()
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
    expect(component.createApiForm.method).toEqual([
      'GET',
      'POST',
      'PUT',
      'DELETE',
      'PATCH',
      'HEAD'])
    expect(component.createApiForm.match).toEqual(MockApiHttpMessage2.data.api.match)
    expect(component.createApiForm.proxyHeader).not.toEqual([])
    expect(component.hostsList[0].key).toEqual('')
    expect(component.hostsList.length).toEqual(1)
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
    for (const m of component.methodList) {
      expect(m.checked).toEqual(true)
    }
    expect(component.allChecked).toEqual(true)
  }))

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
    component.saveApi('http')
    fixture.detectChanges()

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(1)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)
    expect(component.showCheckboxGroupValid).toEqual(true)
    component.saveApi('http')
    fixture.detectChanges()

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(2)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.validateForm.patchValue({
      groupUuid: 'test'
    })

    component.saveApi('http')
    fixture.detectChanges()
    expect(component.showCheckboxGroupValid).toEqual(true)

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(3)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.validateForm.patchValue({
      name: 'test'
    })

    component.saveApi('http')
    fixture.detectChanges()
    expect(component.showCheckboxGroupValid).toEqual(true)

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(4)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.validateForm.patchValue({
      requestPath: 'test'
    })

    component.saveApi('http')
    fixture.detectChanges()
    expect(component.showCheckboxGroupValid).toEqual(true)

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(5)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.validateForm.patchValue({
      service: 'test'
    })

    component.saveApi('http')
    fixture.detectChanges()
    expect(component.showCheckboxGroupValid).toEqual(true)

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(5)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.createApiForm.method = []
    component.saveApi('http')
    fixture.detectChanges()

    expect(spyFormServiceDirty).toHaveBeenCalledTimes(5)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)

    component.createApiForm.method = []
    component.allChecked = true
    component.saveApi('http')
    fixture.detectChanges()

    expect(component.showCheckboxGroupValid).toEqual(false)
    expect(spyFormServiceDirty).toHaveBeenCalledTimes(5)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).toHaveBeenCalledTimes(1)
    expect(component.submitButtonLoading).toEqual(false)

    component.createApiForm.method = ['PUT']
    component.allChecked = false
    component.saveApi('http')
    fixture.detectChanges()

    expect(component.showCheckboxGroupValid).toEqual(false)
    expect(spyFormServiceDirty).toHaveBeenCalledTimes(5)
    expect(spyPutApiService).not.toHaveBeenCalled()
    expect(spyPostApiService).toHaveBeenCalledTimes(2)
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
            return of(MockApiHttpMessage)
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
    component.saveApi('http')
    fixture.detectChanges()

    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(spyPutApiService).toHaveBeenCalledTimes(1)
    expect(component.createApiForm.method).toEqual(['GET']) // http
    expect(component.submitButtonLoading).toEqual(false)

    expect(component.showCheckboxGroupValid).toEqual(false)

    component.allChecked = false
    component.createApiForm.method = []
    component.saveApi('http')
    fixture.detectChanges()

    expect(spyFormServiceDirty).not.toHaveBeenCalled()
    expect(spyPutApiService).toHaveBeenCalledTimes(1)
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(component.submitButtonLoading).toEqual(false)
  })
})
