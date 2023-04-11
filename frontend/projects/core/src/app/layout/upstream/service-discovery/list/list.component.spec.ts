import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { environment } from 'projects/core/src/environments/environment'
import { RouterModule } from '@angular/router'
import { APP_BASE_HREF } from '@angular/common'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { BidiModule } from '@angular/cdk/bidi'
import { Overlay } from '@angular/cdk/overlay'
import { EoNgFeedbackDrawerService, EoNgFeedbackModalService, EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { Subject } from 'rxjs/internal/Subject'
import { of } from 'rxjs'
import { ElementRef } from '@angular/core'
import { ServiceDiscoveryMessageComponent } from '../message/message.component'
import { ServiceDiscoveryPublishComponent } from '../publish/publish.component'
import { ServiceDiscoveryListComponent } from './list.component'
import { ServiceDiscoveryCreateComponent } from '../create/create.component'

class MockDrawerService {
  result:boolean =false

  nzAfterClose = new Subject<any>();

  create () {
    return {
      afterClose: {
        subscribe: () => { of(this.result) }
      },
      close: () => {
        return 'drawer is close'
      }
    }
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

class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

class MockEnsureService {
  create () {
    return 'modal is create'
  }
}

describe('ServiceDiscoveryListComponent test', () => {
  let component: ServiceDiscoveryListComponent
  let fixture: ComponentFixture<ServiceDiscoveryListComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
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
            path: 'upstream/serv-discovery/create',
            component: ServiceDiscoveryCreateComponent
          }
        ]
        )
      ],
      declarations: [
      ],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: EoNgFeedbackDrawerService, useClass: MockDrawerService },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: EoNgFeedbackModalService, useClass: MockEnsureService }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(ServiceDiscoveryListComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('should initial servicesList', fakeAsync(() => {
    expect(component.servicesList).not.toBe([])
  }))

  it('click table btns', () => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')
    // const spyDeleteCert = jest.spyOn(component, 'deleteCert')
    // @ts-ignore
    const spyModal = jest.spyOn(component.modalService, 'create')
    expect(spyRouter).toHaveBeenCalledTimes(0)
    expect(spyModal).toHaveBeenCalledTimes(0)

    const item = { key: 'test' }
    component.servicesTableBody[4].btns[0].click(item)
    expect(spyRouter).toHaveBeenCalledTimes(1)
  })
  it('getServicesList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { discoveries: [{ name: 'test1', driver: 'static', desc: 'description', updateTime: 'time', isDelete: 'false' }, { name: 'test2', driver: 'static', desc: 'description', updateTime: 'time', isDelete: 'false' }] } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    component.servicesList = []
    component.serviceName = ''
    component.serviceNameForSear = 'test1'
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toBe(true)

    component.getServicesList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.servicesList).not.toBe([])
    expect(component.serviceName).toBe('test1')
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)

    component.getServicesList(true)
    fixture.detectChanges()
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
  })

  it('getServicesList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyMessageSuccess).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toBe(true)

    component.getServicesList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessageSuccess).not.toHaveBeenCalled()
  })

  it('addService', fakeAsync(() => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')
    expect(spyRouter).not.toHaveBeenCalled()

    component.addService()
    fixture.detectChanges()

    expect(spyRouter).toHaveBeenCalled()
  }))

  it('delete modal is created', () => {
    // @ts-ignore
    const spyModal = jest.spyOn(component.modalService, 'create')
    expect(spyModal).not.toHaveBeenCalled()
    component.delete('test')
    expect(spyModal).toHaveBeenCalled()
  })

  it('deleteDiscovery with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: 0, data: { msg: 'success' } }))
    const isdelete = httpCommonService.delete('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetServicesList = jest.spyOn(component, 'getServicesList')
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isdelete).toBe(true)
    expect(spyGetServicesList).not.toHaveBeenCalled()

    component.deleteDiscovery('test')
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyGetServicesList).toHaveBeenCalled()
  })

  it('deleteDiscovery with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(0)
    expect(isget).toBe(true)

    component.deleteDiscovery('test')
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
  })
})
