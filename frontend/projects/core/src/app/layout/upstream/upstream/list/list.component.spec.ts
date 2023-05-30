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
import { UpstreamCreateComponent } from '../create/create.component'
import { UpstreamMessageComponent } from '../message/message.component'
import { UpstreamListComponent } from './list.component'

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

describe('UpstreamListComponent test', () => {
  let component: UpstreamListComponent
  let fixture: ComponentFixture<UpstreamListComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule,
        NzDrawerModule, NzOutletModule, HttpClientModule,
        RouterModule.forRoot([
          {
            path: '',
            component: UpstreamCreateComponent
          },
          {
            path: 'message',
            component: UpstreamMessageComponent
          },
          {
            path: 'upstream/upstream/create',
            component: UpstreamCreateComponent
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

    fixture = TestBed.createComponent(UpstreamListComponent)

    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('ngOnInit should call method getUpstreamsList', () => {
    const spyGetUpstreamList = jest.spyOn(component, 'getUpstreamsList')
    expect(spyGetUpstreamList).not.toHaveBeenCalled()
    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetUpstreamList).toHaveBeenCalledTimes(1)
  })

  it('click table btns', () => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')
    expect(spyRouter).toHaveBeenCalledTimes(0)

    const item = { key: 'test' }
    component.upstreamsTableBody[5].btns[0].click(item)
    expect(spyRouter).toHaveBeenCalledTimes(1)
  })

  it('getUpstreamsList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { services: { test: 'test' }, total: 10 } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    component.upstreamsList = []
    component.total = 0
    component.upstreamName = ''
    component.upstreamNameForSear = 'test1'
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toBe(true)

    component.getUpstreamsList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.upstreamsList).not.toBe([])
    expect(component.total).toBe(10)
    expect(component.upstreamName).toBe('test1')
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)

    component.getUpstreamsList(true)
    fixture.detectChanges()
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
  })

  it('getUpstreamsList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessageSuccess).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toBe(true)
    expect(spyMessage).not.toHaveBeenCalled()

    component.getUpstreamsList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
  })

  it('addUpstream', fakeAsync(() => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')
    expect(spyRouter).not.toHaveBeenCalled()

    component.addUpstream()
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

  it('deleteService with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: 0, data: { msg: 'success' } }))
    const isdelete = httpCommonService.delete('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetUpstreamsList = jest.spyOn(component, 'getUpstreamsList')
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isdelete).toBe(true)
    expect(spyGetUpstreamsList).not.toHaveBeenCalled()

    component.deleteService('test')
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyGetUpstreamsList).toHaveBeenCalled()
  })

  it('deleteService with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(0)
    expect(isget).toBe(true)

    component.deleteService('test')
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
  })
})
