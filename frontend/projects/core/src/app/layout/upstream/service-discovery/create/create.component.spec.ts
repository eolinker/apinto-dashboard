/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-14 22:56:33
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-21 23:33:15
 * @FilePath: /apinto/src/app/layout/upstream/service-discovery-content/service-discovery-content.component.spec.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { RouterModule } from '@angular/router'
import { ElementRef, Renderer2, ChangeDetectorRef, Type } from '@angular/core'
import { APP_BASE_HREF } from '@angular/common'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { environment } from 'projects/core/src/environments/environment'
import { BidiModule } from '@angular/cdk/bidi'
import { Overlay } from '@angular/cdk/overlay'
import { ServiceDiscoveryPublishComponent } from '../publish/publish.component'
import { ServiceDiscoveryMessageComponent } from '../message/message.component'
import { ServiceDiscoveryCreateComponent } from './create.component'
import { of } from 'rxjs'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { FormsModule } from '@angular/forms'
import { ServiceDiscoveryListComponent } from '../list/list.component'

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

describe('ServiceDiscoveryCreateComponent test as editPage is false', () => {
  let component: ServiceDiscoveryCreateComponent
  let fixture: ComponentFixture<ServiceDiscoveryCreateComponent>
  let renderer2: Renderer2
  class MockElementRef extends ElementRef {
    constructor () { super(null) }
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule, FormsModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule,
        NzDrawerModule, NzOutletModule, HttpClientModule,
        RouterModule.forRoot([
          {
            path: '',
            component: ServiceDiscoveryPublishComponent
          },
          {
            path: 'message',
            component: ServiceDiscoveryMessageComponent
          },
          {
            path: 'upstream/discovery',
            component: ServiceDiscoveryListComponent
          }
        ]
        )
      ],
      declarations: [ServiceDiscoveryCreateComponent
      ],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: Renderer2, useClass: MockRenderer },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(ServiceDiscoveryCreateComponent)
    renderer2 = fixture.componentRef.injector.get<Renderer2>(Renderer2 as Type<Renderer2>)
    renderer2.removeAttribute = jest.fn().mockRejectedValue('remove')

    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('change editPage from false to true', () => {
    const spyGetServiceMessage = jest.spyOn(component, 'getServiceMessage')
    const spyGetDriverList = jest.spyOn(component, 'getDriverList')
    expect(spyGetServiceMessage).not.toHaveBeenCalled()
    expect(spyGetDriverList).not.toHaveBeenCalled()
    component.ngOnInit()
    expect(spyGetServiceMessage).not.toHaveBeenCalled()
    expect(spyGetDriverList).toHaveBeenCalled()
  })

  it('getDriverList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { drivers: [{ name: 'test1', render: 'render' }, { name: 'test2', render: 'render' }] } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.driverList).toStrictEqual([])
    component.getDriverList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.driverList).not.toStrictEqual([])
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('getDriverList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.driverList).toStrictEqual([])
    component.getDriverList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.driverList).toStrictEqual([])
    expect(spyMessage).toHaveBeenCalled()
  })

  it('changeBasedata', fakeAsync(() => {
    component.driverList = [{ label: 'test1', value: 'test1', render: 'render1' }, { label: 'test2', value: 'test1', render: 'render2' }]
    component.baseData = 'render'
    component.createServiceForm.driver = 'test2'
    expect(component.baseData).toStrictEqual('render')
    component.changeBasedata()
    fixture.detectChanges()
    expect(component.baseData).toStrictEqual('render2')
  }))

  it('getServiceMessage with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { discovery: { name: 'test1', driver: 'static', desc: 'description', config: { addrs: ['test1', 'test2'], params: [{ key: 'key', value: 'value' }] }, render: 'render' } } }))
    const isget = httpCommonService.get('') !== null
    const spyDriverList = jest.spyOn(component, 'getDriverList')

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(spyDriverList).toHaveBeenCalledTimes(0)
    component.getServiceMessage()
    fixture.detectChanges()
    expect(component.createServiceForm).toStrictEqual({ name: 'test1', driver: 'static', desc: 'description', config: { addrs: ['test1', 'test2'], params: [{ key: 'key', value: 'value' }] }, render: 'render' })
    expect(spyService).toHaveBeenCalledTimes(3)
    expect(spyDriverList).toHaveBeenCalledTimes(1)

    const spyService2 = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { discovery: { name: 'test1', driver: 'static', desc: 'description', config: { addrs: ['test1', 'test2'], params: [] }, render: 'render' } } }))
    component.getServiceMessage()
    fixture.detectChanges()
    expect(spyService2).toHaveBeenCalledTimes(5)
    expect(component.createServiceForm.config.params).toStrictEqual([{ key: '', value: '' }])
    expect(spyDriverList).toHaveBeenCalledTimes(2)
  })

  it('getServiceMessage with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.getServiceMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalled()
  })

  it('saveService with success return', () => {
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

    component.editPage = false
    component.saveService()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(spybackToList).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
    expect(spyMessageError).not.toHaveBeenCalled()

    component.editPage = true
    component.saveService()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyServicePut).toHaveBeenCalledTimes(2)
    expect(spybackToList).toHaveBeenCalledTimes(2)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(2)
    expect(spyMessageError).not.toHaveBeenCalled()
  })

  it('saveService with fail return', () => {
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
    component.saveService()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalledTimes(1)

    component.editPage = true
    component.saveService()
    fixture.detectChanges()

    expect(spyServicePut).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalledTimes(2)
  })

  it('backToList', fakeAsync(() => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')
    expect(spyRouter).not.toHaveBeenCalled()

    component.backToList()
    fixture.detectChanges()

    expect(spyRouter).toHaveBeenCalled()
  }))
})

describe('ServiceDiscoveryCreateComponent test as editPage is false', () => {
  let component: ServiceDiscoveryCreateComponent
  let fixture: ComponentFixture<ServiceDiscoveryCreateComponent>
  let renderer2: Renderer2
  class MockElementRef extends ElementRef {
    constructor () { super(null) }
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule, FormsModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule,
        NzDrawerModule, NzOutletModule, HttpClientModule,
        RouterModule.forRoot([
          {
            path: '',
            component: ServiceDiscoveryPublishComponent
          },
          {
            path: 'message',
            component: ServiceDiscoveryMessageComponent
          }
        ]
        )
      ],
      declarations: [ServiceDiscoveryCreateComponent
      ],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: Renderer2, useClass: MockRenderer },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(ServiceDiscoveryCreateComponent)
    renderer2 = fixture.componentRef.injector.get<Renderer2>(Renderer2 as Type<Renderer2>)
    renderer2.removeAttribute = jest.fn().mockRejectedValue('remove')

    component = fixture.componentInstance
    component.editPage = true
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('change editPage from false to true', () => {
    const spyGetServiceMessage = jest.spyOn(component, 'getServiceMessage')
    const spyGetDriverList = jest.spyOn(component, 'getDriverList')
    expect(spyGetServiceMessage).not.toHaveBeenCalled()
    expect(spyGetDriverList).not.toHaveBeenCalled()
    component.ngOnInit()
    expect(spyGetServiceMessage).toHaveBeenCalled()
    expect(spyGetDriverList).not.toHaveBeenCalled()
  })
})
