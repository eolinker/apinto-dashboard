import { ComponentFixture, TestBed } from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { NzTabsModule } from 'ng-zorro-antd/tabs'
import { NzIconTestModule } from 'ng-zorro-antd/icon/testing'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { RouterModule, Router } from '@angular/router'
import { EoNgFeedbackDrawerService, EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { Subject } from 'rxjs/internal/Subject'
import { of } from 'rxjs'
import { DeployClusterContentComponent } from './content.component'
import { ElementRef, Renderer2, ChangeDetectorRef } from '@angular/core'
import { APP_BASE_HREF, CommonModule } from '@angular/common'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { environment } from 'projects/core/src/environments/environment'
import { DeployClusterEnvironmentComponent } from '../environment/environment.component'
import { DeployClusterNodesComponent } from '../nodes/nodes.component'
import { BidiModule } from '@angular/cdk/bidi'
import { Overlay } from '@angular/cdk/overlay'
import { EoNgTabsModule } from 'eo-ng-tabs'
import { RouterTestingModule } from '@angular/router/testing'

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
  error () {
    return 'error'
  }
}

class MockRenderer {
  removeAttribute (element: any, cssClass: string) {
    return cssClass + 'is removed from' + element
  }
}

describe('DeployClusterContentComponent test', () => {
  let component: DeployClusterContentComponent
  let fixture: ComponentFixture<DeployClusterContentComponent>
  let router: Router
  class MockElementRef extends ElementRef {
    constructor () { super(null) }
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule, CommonModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule,
        NzDrawerModule, NzOutletModule, HttpClientModule,
        EoNgTabsModule, NzTabsModule, NzIconTestModule,
        RouterTestingModule.withRoutes([
          {
            path: '',
            component: DeployClusterEnvironmentComponent
          },
          {
            path: 'cert',
            component: DeployClusterContentComponent
          },
          {
            path: 'nodes',
            component: DeployClusterNodesComponent
          }
        ]),
        RouterModule.forRoot([
          {
            path: '',
            component: DeployClusterEnvironmentComponent
          },
          {
            path: 'cert',
            component: DeployClusterContentComponent
          },
          {
            path: 'nodes',
            component: DeployClusterNodesComponent
          }
        ]
        )
      ],
      declarations: [DeployClusterContentComponent
      ],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: EoNgFeedbackDrawerService, useClass: MockDrawerService },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: Renderer2, useClass: MockRenderer },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(DeployClusterContentComponent)

    component = fixture.componentInstance
    fixture.detectChanges()
    router = TestBed.inject(Router)
    router.initialNavigation()
  })
  it('should create', async () => {
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

  it('##getClustersData successfully', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(
        of({
          code: 0,
          data: {
            cluster: {
              env: 'pro',
              status: 'NORMAL',
              desc: 'password',
              name: 'name',
              create_time: '2022-01-02',
              update_time: '2022-01-02'
            }
          }
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'

    component.getClustersData()
    fixture.detectChanges()

    expect(component.clusterDesc).toStrictEqual('password')
    expect(component._clusterDesc).toStrictEqual('password')
  })

  it('##getClustersData failed', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(
        of({
          code: -1,
          data: {
          },
          msg: 'faild'
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'

    component.getClustersData()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('faild')
  })

  it('##getClustersData failed without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(
        of({
          code: -1,
          data: {
          }
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    component.getClustersData()
    fixture.detectChanges()
    expect(component.clusterDesc).toStrictEqual('')
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取列表数据失败！')
  })

  it('##save successfully', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(
        of({
          code: 0,
          data: {},
          msg: 'success'
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    component._clusterDesc = 'test_desc'
    component.clusterDesc = ''
    component.disabled = false

    component.save()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component.clusterDesc).toStrictEqual('test_desc')
    expect(component.disabled).toStrictEqual(true)
  })

  it('##save failed', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(
        of({
          code: -1,
          data: {
          },
          msg: 'faild'
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    component._clusterDesc = 'test_desc'
    component.clusterDesc = ''
    component.disabled = false
    component.save()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(component.clusterDesc).toStrictEqual('')
    expect(spyMessage).toBeCalledWith('faild')
    expect(component.disabled).toStrictEqual(true)
  })

  it('##save failed without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(
        of({
          code: -1,
          data: {
          }
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    component._clusterDesc = 'test_desc'
    component.disabled = false
    component.clusterDesc = ''
    component.save()
    fixture.detectChanges()

    expect(component._clusterDesc).toStrictEqual('')

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(component.disabled).toStrictEqual(true)
    expect(spyMessage).toBeCalledWith('修改失败！')
  })
})
