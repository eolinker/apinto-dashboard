/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-14 22:56:33
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-03 20:10:16
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
import { of } from 'rxjs'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { FormsModule } from '@angular/forms'
import { UpstreamMessageComponent } from '../message/message.component'
import { ServiceDiscoveryMessageComponent } from '../../service-discovery/message/message.component'
import { ServiceDiscoveryPublishComponent } from '../../service-discovery/publish/publish.component'
import { UpstreamCreateComponent } from './create.component'

import { EoNgSelectModule } from 'eo-ng-select'
import { LayoutModule } from 'projects/core/src/app/layout/layout.module'
import { EoNgTableModule } from 'eo-ng-table'

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

describe('UpstreamCreateComponent test as editPage is false', () => {
  let component: UpstreamCreateComponent
  let fixture: ComponentFixture<UpstreamCreateComponent>
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
            component: UpstreamCreateComponent
          },
          {
            path: 'message',
            component: UpstreamMessageComponent
          }
        ]
        )
      ],
      declarations: [UpstreamCreateComponent
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

    fixture = TestBed.createComponent(UpstreamCreateComponent)

    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('change editPage from false to true', () => {
    const spyGetUpstreamMessage = jest.spyOn(component, 'getUpstreamMessage')
    const spyGetDiscovery = jest.spyOn(component, 'getDiscovery')
    expect(spyGetUpstreamMessage).not.toHaveBeenCalled()
    expect(spyGetDiscovery).not.toHaveBeenCalled()
    component.ngOnInit()
    expect(spyGetUpstreamMessage).not.toHaveBeenCalled()
    expect(spyGetDiscovery).toHaveBeenCalled()
  })

  it('getDiscovery with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { discoveries: [{ name: 'test1', render: 'render1', driver: 'static' }, { name: 'test2', render: 'render2', driver: 'driver2' }] } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.discoveryList).toStrictEqual([])
    component.createUpstreamForm.discoveryName = 'test1'
    component.getDiscovery()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.discoveryList).toStrictEqual(
      [{ label: '静态节点', value: 'test1', render: 'render1' },
        { label: 'test2[driver2]', value: 'test2', render: 'render2' }
      ]
    )
    expect(component.baseData).toStrictEqual('render1')
    expect(spyMessage).not.toHaveBeenCalled()

    component.discoveryList = []
    component.createUpstreamForm.discoveryName = 'test2'
    component.getDiscovery()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(3)
    expect(component.discoveryList).toStrictEqual(
      [{ label: '静态节点', value: 'test1', render: 'render1' },
        { label: 'test2[driver2]', value: 'test2', render: 'render2' }
      ]
    )
    expect(component.baseData).toStrictEqual('render2')
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('getDiscovery with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.discoveryList).toStrictEqual([])
    component.getDiscovery()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.discoveryList).toStrictEqual([])
    expect(spyMessage).toHaveBeenCalled()
  })

  it('changeBasedata', fakeAsync(() => {
    component.discoveryList = [{ label: '静态节点', value: 'test1', render: 'render1' },
      { label: 'test2[driver2]', value: 'test2', render: 'render2' }
    ]
    component.baseData = 'render'
    component.createUpstreamForm.discoveryName = 'test'
    expect(component.baseData).toStrictEqual('render')
    component.changeBasedata()
    fixture.detectChanges()
    expect(component.baseData).toStrictEqual('render')

    component.baseData = 'render'
    component.createUpstreamForm.discoveryName = 'test1'
    expect(component.baseData).toStrictEqual('render')
    component.changeBasedata()
    fixture.detectChanges()
    expect(component.baseData).toStrictEqual('render1')
  }))

  it('getUpstreamMessage with success return', () => {
    const mockValue:any = {
      code: 0,
      data: {
        service: {
          name: 'test1',
          driver: 'static',
          desc: 'description',
          config: {
            addrs: ['test1', 'test2'],
            params: [{
              key: 'key',
              value: 'value'
            }]
          },
          render: 'render'
        }
      }
    }
    const mockRes:any = {
      name: 'test1',
      driver: 'static',
      desc: 'description',
      config: {
        serviceName: '',
        useVariable: false,
        addrsVariable: '',
        staticConf: [{
          addr: '',
          weight: null
        }],
        addrs: ['test1', 'test2'],
        params: [{
          key: 'key',
          value: 'value'
        }]
      },
      render: 'render'
    }
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of(mockValue))
    const isget = httpCommonService.get('') !== null
    const spyGetDiscovery = jest.spyOn(component, 'getDiscovery')

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(spyGetDiscovery).toHaveBeenCalledTimes(0)
    component.getUpstreamMessage()
    fixture.detectChanges()
    expect(component.createUpstreamForm).toStrictEqual(mockRes)
    expect(spyService).toHaveBeenCalledTimes(3)
    expect(spyGetDiscovery).toHaveBeenCalledTimes(1)

    const mockValue2 = {
      code: 0,
      data: {
        service: {
          name: 'test1',
          driver: 'static',
          desc: 'description',
          config: {
            staticConf: [{ addr: 'addr1', weight: 'weight1' }]
          }
        }
      }
    }

    const mockRes2 = {
      name: 'test1',
      driver: 'static',
      desc: 'description',
      config: {
        serviceName: '',
        useVariable: false,
        addrsVariable: '',
        staticConf: [{ addr: 'addr1', weight: 'weight1' }]
      }
    }
    const spyService2 = jest.spyOn(httpCommonService, 'get').mockReturnValue(of(mockValue2))
    component.getUpstreamMessage()
    fixture.detectChanges()
    expect(spyService2).toHaveBeenCalledTimes(5)
    expect(component.createUpstreamForm).toStrictEqual(mockRes2)
    expect(spyGetDiscovery).toHaveBeenCalledTimes(2)
  })

  it('getUpstreamMessage with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.getUpstreamMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalled()
  })

  it('getDataFromDynamicComponent when discoveryName is static', () => {
    const testVal1 = {
      discoveryName: 'static',
      config: {
        serviceName: '',
        useVariable: true,
        addrsVariable: '',
        staticConf: [{ weight: 0, addr: '' }]
      }
    }

    expect(component.canBeSave).toStrictEqual(false)
    component.getDataFromDynamicComponent(testVal1)
    fixture.detectChanges()

    expect(component.canBeSave).toStrictEqual(false)

    testVal1.config.addrsVariable = 'test'
    expect(component.canBeSave).toStrictEqual(false)
    component.getDataFromDynamicComponent(testVal1)
    fixture.detectChanges()

    expect(component.canBeSave).toStrictEqual(true)

    testVal1.config.useVariable = false
    expect(component.canBeSave).toStrictEqual(true)
    component.getDataFromDynamicComponent(testVal1)
    fixture.detectChanges()

    expect(component.canBeSave).toStrictEqual(false)

    testVal1.config.staticConf = [{ addr: 'testAddr', weight: 0 }]

    testVal1.config.useVariable = false
    expect(component.canBeSave).toStrictEqual(false)
    component.getDataFromDynamicComponent(testVal1)
    fixture.detectChanges()

    expect(component.canBeSave).toStrictEqual(false)

    testVal1.config.staticConf = [{ addr: '', weight: 3 }]

    testVal1.config.useVariable = false
    expect(component.canBeSave).toStrictEqual(false)
    component.getDataFromDynamicComponent(testVal1)
    fixture.detectChanges()

    expect(component.canBeSave).toStrictEqual(false)

    testVal1.config.staticConf = [{ addr: 'testAddr', weight: 3 }]

    testVal1.config.useVariable = false
    expect(component.canBeSave).toStrictEqual(false)
    component.getDataFromDynamicComponent(testVal1)
    fixture.detectChanges()

    expect(component.canBeSave).toStrictEqual(true)
  })

  it('getDataFromDynamicComponent when discoveryName is static', () => {
    const testVal2 = {
      discoveryName: 'test',
      config: {
        serviceName: '',
        useVariable: true,
        addrsVariable: '',
        staticConf: [{ weight: 0, addr: '' }]
      }
    }

    expect(component.canBeSave).toStrictEqual(false)
    component.getDataFromDynamicComponent(testVal2)
    fixture.detectChanges()

    expect(component.canBeSave).toStrictEqual(false)

    testVal2.config.serviceName = 'test'

    expect(component.canBeSave).toStrictEqual(false)
    component.getDataFromDynamicComponent(testVal2)
    fixture.detectChanges()

    expect(component.canBeSave).toStrictEqual(true)
  })

  it('saveUpstream with success return', () => {
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

    component.createUpstreamForm = {
      name: '',
      desc: '',
      scheme: 'HTTP',
      balance: 'round-robin',
      discoveryName: 'static',
      timeout: 100,
      config: {
        addrsVariable: '',
        useVariable: false,
        serviceName: '',
        staticConf: []
      }
    }

    component.editPage = false
    component.saveUpstream()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(spybackToList).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
    expect(spyMessageError).not.toHaveBeenCalled()

    component.editPage = true
    component.saveUpstream()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyServicePut).toHaveBeenCalledTimes(2)
    expect(spybackToList).toHaveBeenCalledTimes(2)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(2)
    expect(spyMessageError).not.toHaveBeenCalled()
  })

  it('saveUpstream with fail return', () => {
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
    component.saveUpstream()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalledTimes(1)

    component.editPage = true
    component.saveUpstream()
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

describe('UpstreamCreateComponent test as editPage is false', () => {
  let component: UpstreamCreateComponent
  let fixture: ComponentFixture<UpstreamCreateComponent>
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
      declarations: [UpstreamCreateComponent
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

    fixture = TestBed.createComponent(UpstreamCreateComponent)
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
    const spyGetUpstreamMessage = jest.spyOn(component, 'getUpstreamMessage')
    const spyGetDiscovery = jest.spyOn(component, 'getDiscovery')
    expect(spyGetUpstreamMessage).not.toHaveBeenCalled()
    expect(spyGetDiscovery).not.toHaveBeenCalled()
    component.ngOnInit()
    expect(spyGetUpstreamMessage).toHaveBeenCalled()
    expect(spyGetDiscovery).not.toHaveBeenCalled()
  })
})
