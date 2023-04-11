import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { environment } from 'projects/core/src/environments/environment'
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
import { RouterModule } from '@angular/router'
import { ApplicationAuthenticationComponent } from '../authentication/authentication.component'
import { ApplicationPublishComponent } from './publish.component'
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

  remove () {
    return 'remove'
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

describe('ApplicationPublishComponent test', () => {
  let component: ApplicationPublishComponent
  let fixture: ComponentFixture<ApplicationPublishComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule,
        NzDrawerModule, NzOutletModule, HttpClientModule,
        RouterModule.forRoot([
          {
            path: '',
            component: ApplicationPublishComponent
          },
          {
            path: 'message',
            component: ApplicationAuthenticationComponent
          },
          {
            path: 'authentication',
            component: ApplicationAuthenticationComponent
          }
        ])
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

    fixture = TestBed.createComponent(ApplicationPublishComponent)

    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('click table btns', () => {
    const spyUpdateOrOnline = jest.spyOn(component, 'updateOrOnline')
    const spyOffline = jest.spyOn(component, 'offline')
    const spyDisable = jest.spyOn(component, 'disable')
    const spyEnable = jest.spyOn(component, 'enable')
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(0)
    expect(spyOffline).toHaveBeenCalledTimes(0)
    expect(spyDisable).toHaveBeenCalledTimes(0)
    expect(spyEnable).toHaveBeenCalledTimes(0)

    const item = { key: 'test' }
    component.clustersTableBody[6].btns[0].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(1)
    expect(spyOffline).toHaveBeenCalledTimes(0)
    expect(spyDisable).toHaveBeenCalledTimes(0)
    expect(spyEnable).toHaveBeenCalledTimes(0)

    component.clustersTableBody[6].btns[1].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(1)
    expect(spyOffline).toHaveBeenCalledTimes(1)
    expect(spyDisable).toHaveBeenCalledTimes(0)
    expect(spyEnable).toHaveBeenCalledTimes(0)

    component.clustersTableBody[6].btns[2].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(1)
    expect(spyOffline).toHaveBeenCalledTimes(1)
    expect(spyDisable).toHaveBeenCalledTimes(1)
    expect(spyEnable).toHaveBeenCalledTimes(0)

    component.clustersTableBody[7].btns[0].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(2)
    expect(spyOffline).toHaveBeenCalledTimes(1)
    expect(spyDisable).toHaveBeenCalledTimes(1)
    expect(spyEnable).toHaveBeenCalledTimes(0)

    component.clustersTableBody[7].btns[1].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(2)
    expect(spyOffline).toHaveBeenCalledTimes(2)
    expect(spyDisable).toHaveBeenCalledTimes(1)
    expect(spyEnable).toHaveBeenCalledTimes(0)

    component.clustersTableBody[7].btns[2].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(2)
    expect(spyOffline).toHaveBeenCalledTimes(2)
    expect(spyDisable).toHaveBeenCalledTimes(1)
    expect(spyEnable).toHaveBeenCalledTimes(1)

    component.clustersTableBody[8].btns[0].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(2)
    expect(spyOffline).toHaveBeenCalledTimes(3)
    expect(spyDisable).toHaveBeenCalledTimes(1)
    expect(spyEnable).toHaveBeenCalledTimes(1)

    component.clustersTableBody[8].btns[1].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(2)
    expect(spyOffline).toHaveBeenCalledTimes(3)
    expect(spyDisable).toHaveBeenCalledTimes(2)
    expect(spyEnable).toHaveBeenCalledTimes(1)

    component.clustersTableBody[9].btns[0].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(2)
    expect(spyOffline).toHaveBeenCalledTimes(4)
    expect(spyDisable).toHaveBeenCalledTimes(2)
    expect(spyEnable).toHaveBeenCalledTimes(1)

    component.clustersTableBody[9].btns[1].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(2)
    expect(spyOffline).toHaveBeenCalledTimes(4)
    expect(spyDisable).toHaveBeenCalledTimes(2)
    expect(spyEnable).toHaveBeenCalledTimes(2)

    component.clustersTableBody[10].btns[0].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(3)
    expect(spyOffline).toHaveBeenCalledTimes(4)
    expect(spyDisable).toHaveBeenCalledTimes(2)
    expect(spyEnable).toHaveBeenCalledTimes(2)

    component.clustersTableBody[10].btns[1].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(3)
    expect(spyOffline).toHaveBeenCalledTimes(4)
    expect(spyDisable).toHaveBeenCalledTimes(3)
    expect(spyEnable).toHaveBeenCalledTimes(2)

    component.clustersTableBody[11].btns[0].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(4)
    expect(spyOffline).toHaveBeenCalledTimes(4)
    expect(spyDisable).toHaveBeenCalledTimes(3)
    expect(spyEnable).toHaveBeenCalledTimes(2)

    component.clustersTableBody[11].btns[1].click(item)
    expect(spyUpdateOrOnline).toHaveBeenCalledTimes(4)
    expect(spyOffline).toHaveBeenCalledTimes(4)
    expect(spyDisable).toHaveBeenCalledTimes(3)
    expect(spyEnable).toHaveBeenCalledTimes(3)
  })

  it('getClustersData with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { clusters: [1, 2, 3] } }))
    const isget = httpCommonService.get('') !== null

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.getClustersData()
    fixture.detectChanges()
    expect(component.clustersList).toStrictEqual([1, 2, 3])
    expect(spyService).toHaveBeenCalledTimes(2)
  })

  it('getClustersData with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.getClustersData()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalled()
  })

  it('updateOrOnline with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 0, data: { clusters: [1, 2, 3] } }))

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetClustersData = jest.spyOn(component, 'getClustersData')

    expect(spyService).toHaveBeenCalledTimes(0)
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
    expect(spyGetClustersData).toHaveBeenCalledTimes(0)

    expect(component.solutionRouter).toStrictEqual('')
    expect(component.solutionParam).toStrictEqual({})

    component.updateOrOnline({ name: 'test' }, '更新')
    fixture.detectChanges()
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
    expect(spyGetClustersData).toHaveBeenCalledTimes(1)
  })

  it('updateOrOnline with fail return(no router)', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: -1, data: {}, msg: 'fail' }))

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessageRemove = jest.spyOn(component.message, 'remove')
    const spyGetClustersData = jest.spyOn(component, 'getClustersData')

    expect(spyService).toHaveBeenCalledTimes(0)
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
    expect(spyMessageRemove).toHaveBeenCalledTimes(0)
    expect(spyGetClustersData).toHaveBeenCalledTimes(0)

    expect(component.solutionRouter).toStrictEqual('')
    expect(component.solutionParam).toStrictEqual({})
    component.errorMessageId = ''
    component.updateOrOnline({ name: 'test' }, '上线')
    fixture.detectChanges()

    fixture.whenStable().then(() => {
      expect(spyMessageError).toHaveBeenCalledTimes(1)
      expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
      expect(spyGetClustersData).toHaveBeenCalledTimes(0)
      expect(spyMessageRemove).toHaveBeenCalledTimes(0)

      expect(component.solutionRouter).toStrictEqual('')
      expect(component.solutionParam).toStrictEqual({})
      expect(component.type).toStrictEqual('上线')
      expect(component.failmsg).toStrictEqual('fail')
      expect(component.errorMessageId).not.toStrictEqual('')
    })

    component.updateOrOnline({ name: 'test' }, '上线')
    fixture.detectChanges()

    fixture.whenStable().then(() => {
      expect(spyMessageError).toHaveBeenCalledTimes(2)
      expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
      expect(spyGetClustersData).toHaveBeenCalledTimes(0)
      expect(spyMessageRemove).toHaveBeenCalledTimes(1)

      expect(component.solutionRouter).toStrictEqual('')
      expect(component.solutionParam).toStrictEqual({})
      expect(component.type).toStrictEqual('上线')
      expect(component.failmsg).toStrictEqual('fail')
      expect(component.errorMessageId).not.toStrictEqual('')
    })
  })

  it('updateOrOnline with fail return(no router)', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: -1, data: { router: { name: 'routerName', params: 'routerParams' } }, msg: 'fail' }))

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessageRemove = jest.spyOn(component.message, 'remove')
    const spyGetClustersData = jest.spyOn(component, 'getClustersData')

    expect(spyService).toHaveBeenCalledTimes(0)
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
    expect(spyMessageRemove).toHaveBeenCalledTimes(0)
    expect(spyGetClustersData).toHaveBeenCalledTimes(0)

    expect(component.solutionRouter).toStrictEqual('')
    expect(component.solutionParam).toStrictEqual({})
    component.errorMessageId = ''
    component.updateOrOnline({ name: 'test' }, '上线')
    fixture.detectChanges()

    fixture.whenStable().then(() => {
      expect(spyMessageError).toHaveBeenCalledTimes(1)
      expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
      expect(spyGetClustersData).toHaveBeenCalledTimes(0)
      expect(spyMessageRemove).toHaveBeenCalledTimes(0)

      expect(component.solutionRouter).toStrictEqual('routerName')
      expect(component.solutionParam).toStrictEqual('routerParams')
      expect(component.type).toStrictEqual('上线')
      expect(component.failmsg).toStrictEqual('fail')
      expect(component.errorMessageId).not.toStrictEqual('')
    })

    component.updateOrOnline({ name: 'test' }, '上线')
    fixture.detectChanges()

    fixture.whenStable().then(() => {
      expect(spyMessageError).toHaveBeenCalledTimes(2)
      expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
      expect(spyGetClustersData).toHaveBeenCalledTimes(0)
      expect(spyMessageRemove).toHaveBeenCalledTimes(1)

      expect(component.solutionRouter).toStrictEqual('routerName')
      expect(component.solutionParam).toStrictEqual('routerParams')
      expect(component.type).toStrictEqual('上线')
      expect(component.failmsg).toStrictEqual('fail')
      expect(component.errorMessageId).not.toStrictEqual('')
    })
  })

  it('offline with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 0, data: { clusters: [1, 2, 3] } }))

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetClustersData = jest.spyOn(component, 'getClustersData')

    expect(spyService).toHaveBeenCalledTimes(0)
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
    expect(spyGetClustersData).toHaveBeenCalledTimes(0)

    expect(component.solutionRouter).toStrictEqual('')
    expect(component.solutionParam).toStrictEqual({})

    component.offline({ name: 'test' })
    fixture.detectChanges()
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
    expect(spyGetClustersData).toHaveBeenCalledTimes(1)
  })

  it('offline with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: -1, data: {}, msg: 'fail' }))

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessageRemove = jest.spyOn(component.message, 'remove')
    const spyGetClustersData = jest.spyOn(component, 'getClustersData')

    expect(spyService).toHaveBeenCalledTimes(0)
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
    expect(spyMessageRemove).toHaveBeenCalledTimes(0)
    expect(spyGetClustersData).toHaveBeenCalledTimes(0)

    expect(component.solutionRouter).toStrictEqual('')
    expect(component.solutionParam).toStrictEqual({})
    component.errorMessageId = ''
    component.offline({ name: 'test' })
    fixture.detectChanges()

    fixture.whenStable().then(() => {
      expect(spyMessageError).toHaveBeenCalledTimes(1)
      expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
      expect(spyGetClustersData).toHaveBeenCalledTimes(0)
      expect(spyMessageRemove).toHaveBeenCalledTimes(0)

      expect(component.solutionRouter).toStrictEqual('')
      expect(component.solutionParam).toStrictEqual({})
      expect(component.type).toStrictEqual('下线')
      expect(component.failmsg).toStrictEqual('fail')
      expect(component.errorMessageId).not.toStrictEqual('')
    })

    component.offline({ name: 'test' })
    fixture.detectChanges()

    fixture.whenStable().then(() => {
      expect(spyMessageError).toHaveBeenCalledTimes(2)
      expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
      expect(spyGetClustersData).toHaveBeenCalledTimes(0)
      expect(spyMessageRemove).toHaveBeenCalledTimes(1)

      expect(component.solutionRouter).toStrictEqual('')
      expect(component.solutionParam).toStrictEqual({})
      expect(component.type).toStrictEqual('上线')
      expect(component.failmsg).toStrictEqual('fail')
      expect(component.errorMessageId).not.toStrictEqual('')
    })
  })

  it('enable with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 0, data: { clusters: [1, 2, 3] } }))

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetClustersData = jest.spyOn(component, 'getClustersData')

    expect(spyService).toHaveBeenCalledTimes(0)
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
    expect(spyGetClustersData).toHaveBeenCalledTimes(0)

    expect(component.solutionRouter).toStrictEqual('')
    expect(component.solutionParam).toStrictEqual({})

    component.enable({ name: 'test' })
    fixture.detectChanges()
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
    expect(spyGetClustersData).toHaveBeenCalledTimes(1)
  })

  it('enable with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: -1, data: {}, msg: 'fail' }))

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessageRemove = jest.spyOn(component.message, 'remove')
    const spyGetClustersData = jest.spyOn(component, 'getClustersData')

    expect(spyService).toHaveBeenCalledTimes(0)
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
    expect(spyMessageRemove).toHaveBeenCalledTimes(0)
    expect(spyGetClustersData).toHaveBeenCalledTimes(0)

    expect(component.solutionRouter).toStrictEqual('')
    expect(component.solutionParam).toStrictEqual({})
    component.errorMessageId = ''
    component.enable({ name: 'test' })
    fixture.detectChanges()

    fixture.whenStable().then(() => {
      expect(spyMessageError).toHaveBeenCalledTimes(1)
      expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
      expect(spyGetClustersData).toHaveBeenCalledTimes(0)
      expect(spyMessageRemove).toHaveBeenCalledTimes(0)

      expect(component.solutionRouter).toStrictEqual('')
      expect(component.solutionParam).toStrictEqual({})
      expect(component.failmsg).toStrictEqual('fail')
      expect(component.errorMessageId).not.toStrictEqual('')
    })

    component.enable({ name: 'test' })
    fixture.detectChanges()

    fixture.whenStable().then(() => {
      expect(spyMessageError).toHaveBeenCalledTimes(2)
      expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
      expect(spyGetClustersData).toHaveBeenCalledTimes(0)
      expect(spyMessageRemove).toHaveBeenCalledTimes(1)

      expect(component.solutionRouter).toStrictEqual('')
      expect(component.solutionParam).toStrictEqual({})
      expect(component.failmsg).toStrictEqual('fail')
      expect(component.errorMessageId).not.toStrictEqual('')
    })
  })

  it('disable with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 0, data: { clusters: [1, 2, 3] } }))

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetClustersData = jest.spyOn(component, 'getClustersData')

    expect(spyService).toHaveBeenCalledTimes(0)
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
    expect(spyGetClustersData).toHaveBeenCalledTimes(0)

    expect(component.solutionRouter).toStrictEqual('')
    expect(component.solutionParam).toStrictEqual({})

    component.disable({ name: 'test' })
    fixture.detectChanges()
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
    expect(spyGetClustersData).toHaveBeenCalledTimes(1)
  })

  it('disable with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: -1, data: {}, msg: 'fail' }))

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessageRemove = jest.spyOn(component.message, 'remove')
    const spyGetClustersData = jest.spyOn(component, 'getClustersData')

    expect(spyService).toHaveBeenCalledTimes(0)
    expect(spyMessageError).toHaveBeenCalledTimes(0)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
    expect(spyMessageRemove).toHaveBeenCalledTimes(0)
    expect(spyGetClustersData).toHaveBeenCalledTimes(0)

    expect(component.solutionRouter).toStrictEqual('')
    expect(component.solutionParam).toStrictEqual({})
    component.errorMessageId = ''
    component.disable({ name: 'test' })
    fixture.detectChanges()

    fixture.whenStable().then(() => {
      expect(spyMessageError).toHaveBeenCalledTimes(1)
      expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
      expect(spyGetClustersData).toHaveBeenCalledTimes(0)
      expect(spyMessageRemove).toHaveBeenCalledTimes(0)

      expect(component.solutionRouter).toStrictEqual('')
      expect(component.solutionParam).toStrictEqual({})
      expect(component.failmsg).toStrictEqual('fail')
      expect(component.errorMessageId).not.toStrictEqual('')
    })

    component.disable({ name: 'test' })
    fixture.detectChanges()

    fixture.whenStable().then(() => {
      expect(spyMessageError).toHaveBeenCalledTimes(2)
      expect(spyMessageSuccess).toHaveBeenCalledTimes(0)
      expect(spyGetClustersData).toHaveBeenCalledTimes(0)
      expect(spyMessageRemove).toHaveBeenCalledTimes(1)

      expect(component.solutionRouter).toStrictEqual('')
      expect(component.solutionParam).toStrictEqual({})
      expect(component.failmsg).toStrictEqual('fail')
      expect(component.errorMessageId).not.toStrictEqual('')
    })
  })

  it('closeErrorBtn', fakeAsync(() => {
    // @ts-ignore
    const spyMessageRemove = jest.spyOn(component.message, 'remove')

    expect(spyMessageRemove).toHaveBeenCalledTimes(0)
    component.errorMessageId = ''
    component.closeErrorBtn()
    fixture.detectChanges()
    expect(spyMessageRemove).toHaveBeenCalledTimes(1)
    expect(component.errorMessageId).toStrictEqual('')

    component.closeErrorBtn()
    fixture.detectChanges()
    expect(spyMessageRemove).toHaveBeenCalledTimes(2)
    expect(component.errorMessageId).toStrictEqual('')

    component.errorMessageId = '100'
    component.closeErrorBtn()
    fixture.detectChanges()
    expect(spyMessageRemove).toHaveBeenCalledTimes(3)
    expect(component.errorMessageId).toStrictEqual('')
  }))

  it('viewSolution', fakeAsync(() => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')

    expect(spyRouter).toHaveBeenCalledTimes(0)

    component.solutionRouter = 'test'
    component.viewSolution()
    fixture.detectChanges()
    expect(spyRouter).toHaveBeenCalledTimes(1)

    component.solutionRouter = 'test/test'
    component.viewSolution()
    fixture.detectChanges()
    expect(spyRouter).toHaveBeenCalledTimes(2)
  }))
})
