/* eslint-disable dot-notation */
import {
  ComponentFixture,
  fakeAsync,
  flush,
  TestBed
} from '@angular/core/testing'
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
import {
  EoNgFeedbackDrawerService,
  EoNgFeedbackModalService,
  EoNgFeedbackMessageService
} from 'eo-ng-feedback'
import { Subject } from 'rxjs/internal/Subject'
import { of } from 'rxjs'
import { ReactiveFormsModule } from '@angular/forms'
import { DeployEnvironmentCreateComponent } from './create.component'
class MockDrawerService {
  result: boolean = false

  nzAfterClose = new Subject<any>()

  create() {
    return {
      afterClose: {
        subscribe: () => {
          of(this.result)
        }
      },
      close: () => {
        return 'drawer is close'
      }
    }
  }
}

class MockMessageService {
  error() {
    return 'error'
  }

  success() {
    return 'success'
  }
}

class MockEnsureService {
  create() {
    return 'modal is create'
  }
}

describe('DeployEnvironmentCreateComponent test', () => {
  let component: DeployEnvironmentCreateComponent
  let fixture: ComponentFixture<DeployEnvironmentCreateComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
        BidiModule,
        NoopAnimationsModule,
        NzNoAnimationModule,
        NzDrawerModule,
        NzOutletModule,
        HttpClientModule,
        ReactiveFormsModule,
        RouterModule.forRoot([
          {
            path: '',
            component: DeployEnvironmentCreateComponent
          },
          {
            path: 'deploy/cluster/content/cert',
            component: DeployEnvironmentCreateComponent
          },
          {
            path: 'nodes',
            component: DeployEnvironmentCreateComponent
          }
        ])
      ],
      declarations: [],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: EoNgFeedbackDrawerService, useClass: MockDrawerService },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: EoNgFeedbackModalService, useClass: MockEnsureService }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(DeployEnvironmentCreateComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('##ngOnDestroy test', () => {
    // @ts-ignore
    const spySubscription = jest.spyOn(component.subscription, 'unsubscribe')
    // @ts-ignore
    expect(spySubscription).not.toHaveBeenCalled()
    component.ngOnDestroy()
    expect(spySubscription).toHaveBeenCalled()
  })

  it('##save successfully', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(
      of({
        code: 0,
        data: { variables: [1, 2], total: 1 }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyRouterChange = jest.spyOn(component.router, 'navigate')

    component.save()
    fixture.detectChanges()
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyRouterChange).not.toHaveBeenCalled()
    component.validateForm.controls['key'].setValue('test_')
    component.save()
    fixture.detectChanges()

    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalledWith('新增环境变量成功！')
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyRouterChange).toHaveBeenCalled()
    expect(spyRouterChange).toHaveBeenCalledWith(['/', 'deploy', 'env'])
  })

  it('##save failed', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(
      of({
        code: -1,
        data: {},
        msg: 'faild'
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyRouterChange = jest.spyOn(component.router, 'navigate')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyRouterChange).not.toHaveBeenCalled()
    component.validateForm.controls['key'].setValue('test_')
    component.save()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyRouterChange).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('faild')
  })

  it('##save failed without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(
      of({
        code: -1,
        data: {}
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyRouterChange = jest.spyOn(component.router, 'navigate')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyRouterChange).not.toHaveBeenCalled()
    component.validateForm.controls['key'].setValue('test_')
    component.save()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyRouterChange).not.toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('新增环境变量失败！')
  })
})
