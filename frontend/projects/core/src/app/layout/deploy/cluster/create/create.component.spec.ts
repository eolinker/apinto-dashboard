/* eslint-disable dot-notation */
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
import { DeployClusterCreateComponent } from './create.component'
import { ReactiveFormsModule } from '@angular/forms'
import { EoNgSelectModule } from 'eo-ng-select'

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

describe('DeployClusterCreateComponent test', () => {
  let component: DeployClusterCreateComponent
  let fixture: ComponentFixture<DeployClusterCreateComponent>
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
        ReactiveFormsModule, EoNgSelectModule,
        RouterTestingModule.withRoutes([
          {
            path: '',
            component: DeployClusterEnvironmentComponent
          },
          {
            path: 'cert',
            component: DeployClusterCreateComponent
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
            component: DeployClusterCreateComponent
          },
          {
            path: 'nodes',
            component: DeployClusterNodesComponent
          }
        ]
        )
      ],
      declarations: [DeployClusterCreateComponent
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

    fixture = TestBed.createComponent(DeployClusterCreateComponent)

    component = fixture.componentInstance
    fixture.detectChanges()
    router = TestBed.inject(Router)
    router.initialNavigation()
  })
  it('should create', async () => {
    expect(component).toBeTruthy()
  })

  it('##getEnvList successfully', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(
        of({
          code: 0,
          data: {
            envs: [
              { name: 'PRO', value: 'PRO' },
              { name: 'FAT', value: 'FAT' },
              { name: 'DEV', value: 'DEV' },
              { name: 'UAT', value: 'UAT' }
            ]
          }
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getEnvList()
    fixture.detectChanges()
    expect(component.environmentList).toStrictEqual([
      { label: 'PRO', value: 'PRO' },
      { label: 'FAT', value: 'FAT' },
      { label: 'DEV', value: 'DEV' },
      { label: 'UAT', value: 'UAT' }
    ])
    expect(component.validateForm.controls['envValue'].value).toStrictEqual('PRO')
  })

  it('##getEnvList failed', () => {
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

    component.getEnvList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('faild')
  })

  it('##getEnvList failed without msg', () => {
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
    component.getEnvList()
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取列表数据失败！')
  })

  it('##testCluster successfully', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(
        of({
          code: 0,
          data: {
            nodes: [
              { name: 'name1', admin_addr: 'addr1', service_addr: 'addr1', status: 'NOTRUNNING' },
              { name: 'name2', admin_addr: 'addr2', service_addr: 'addr2', status: 'NOTRUNNING' },
              { name: 'name3', admin_addr: 'addr3', service_addr: 'addr3', status: 'RUNNING' },
              { name: 'name4', admin_addr: 'addr4', service_addr: 'addr4', status: 'RUNNING' }
            ],
            source: 'tttttttttttt',
            is_update: true
          }
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyClusterMarkAsDirty = jest.spyOn(component.validateForm.controls['clusterAddr'], 'markAsDirty')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyClusterMarkAsDirty).not.toHaveBeenCalled()

    component.validateForm.controls['clusterAddr'].setValue(null)
    component.testCluster()
    fixture.detectChanges()
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyClusterMarkAsDirty).toHaveBeenCalled()

    component.validateForm.controls['clusterAddr'].setValue('123.123.123:111')
    component.testCluster()
    fixture.detectChanges()
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyClusterMarkAsDirty).toHaveBeenCalledTimes(2)

    component.validateForm.controls['clusterAddr'].setValue('http:123.123.123.123::')
    component.testCluster()
    fixture.detectChanges()
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyClusterMarkAsDirty).toHaveBeenCalledTimes(3)

    component.validateForm.controls['clusterAddr'].setValue('http://123.123.123:123')
    component.testCluster()
    fixture.detectChanges()
    expect(spyClusterMarkAsDirty).toHaveBeenCalledTimes(3)
    expect(component.nodesList).toStrictEqual([
      { name: 'name1', admin_addr: 'addr1', service_addr: 'addr1', status: 'NOTRUNNING' },
      { name: 'name2', admin_addr: 'addr2', service_addr: 'addr2', status: 'NOTRUNNING' },
      { name: 'name3', admin_addr: 'addr3', service_addr: 'addr3', status: 'RUNNING' },
      { name: 'name4', admin_addr: 'addr4', service_addr: 'addr4', status: 'RUNNING' }
    ])
    expect(component.clusterCanBeCreated).toStrictEqual(true)
    expect(component.source).toStrictEqual('tttttttttttt')
    expect(component.testPassAddr).toStrictEqual('http://123.123.123:123')
  })

  it('##testCluster failed', () => {
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
    const spyClusterMarkAsDirty = jest.spyOn(component.validateForm.controls['clusterAddr'], 'markAsDirty')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.validateForm.controls['clusterAddr'].setValue('http://123.123.123:123')

    component.testCluster()
    fixture.detectChanges()

    expect(spyClusterMarkAsDirty).not.toHaveBeenCalled()
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('faild')
  })

  it('##testCluster failed without msg', () => {
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
    const spyClusterMarkAsDirty = jest.spyOn(component.validateForm.controls['clusterAddr'], 'markAsDirty')
    component.validateForm.controls['clusterAddr'].setValue('http://123.123.123:123')
    component.testCluster()
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('操作失败！')
    expect(spyClusterMarkAsDirty).not.toHaveBeenCalled()
  })

  it('##saveCluster successfully', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'post')
      .mockReturnValue(
        of({
          code: 0,
          data: {
          },
          msg: 'success'
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyClusterMarkAsDirty = jest.spyOn(component.validateForm.controls['clusterAddr'], 'markAsDirty')
    const spyFormMarkAllAsTouched = jest.spyOn(component.validateForm, 'markAllAsTouched')
    // @ts-ignore
    const spyRouterNavigate = jest.spyOn(component.router, 'navigate')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyClusterMarkAsDirty).not.toHaveBeenCalled()
    expect(spyFormMarkAllAsTouched).not.toHaveBeenCalled()
    expect(spyRouterNavigate).not.toHaveBeenCalled()

    component.validateForm.controls['clusterAddr'].setValue(null)
    component.saveCluster()
    fixture.detectChanges()
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyRouterNavigate).not.toHaveBeenCalled()
    expect(spyFormMarkAllAsTouched).toHaveBeenCalled()

    component.validateForm.controls['clusterAddr'].setValue('123.123.123:111')
    component.saveCluster()
    fixture.detectChanges()
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyClusterMarkAsDirty).toHaveBeenCalled()
    expect(spyRouterNavigate).not.toHaveBeenCalled()
    expect(spyFormMarkAllAsTouched).toHaveBeenCalledTimes(2)

    component.validateForm.controls['clusterAddr'].setValue('http:123.123.123.123::')
    component.saveCluster()
    fixture.detectChanges()
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyClusterMarkAsDirty).toHaveBeenCalledTimes(3)
    expect(spyRouterNavigate).not.toHaveBeenCalled()
    expect(spyFormMarkAllAsTouched).toHaveBeenCalledTimes(3)

    component.validateForm.controls['clusterAddr'].setValue('http://123.123.123:123')
    component.saveCluster()
    fixture.detectChanges()
    expect(spyClusterMarkAsDirty).toHaveBeenCalledTimes(3)
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyRouterNavigate).not.toHaveBeenCalled()
    expect(spyFormMarkAllAsTouched).toHaveBeenCalledTimes(4)

    component.validateForm.controls['envValue'].setValue('PRO')
    component.validateForm.controls['clusterName'].setValue('123')
    component.saveCluster()
    fixture.detectChanges()
    expect(spyClusterMarkAsDirty).toHaveBeenCalledTimes(3)
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyRouterNavigate).not.toHaveBeenCalled()
    expect(spyFormMarkAllAsTouched).toHaveBeenCalledTimes(5)

    component.validateForm.controls['envValue'].setValue('PRO')
    component.validateForm.controls['clusterName'].setValue('test123_')
    component.validateForm.controls['envValue'].setValue('123')
    component.saveCluster()
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyFormMarkAllAsTouched).toHaveBeenCalledTimes(6)
    expect(spyRouterNavigate).toHaveBeenCalled()
    expect(spyRouterNavigate).toHaveBeenCalledWith(['/', 'deploy', 'cluster', 'content'], { queryParams: { cluster_name: 'test123_' } })
  })

  it('##saveCluster failed', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'post')
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
    const spyClusterMarkAsDirty = jest.spyOn(component.validateForm.controls['clusterAddr'], 'markAsDirty')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.validateForm.controls['envValue'].setValue('PRO')
    component.validateForm.controls['clusterName'].setValue('test123_')
    component.validateForm.controls['clusterAddr'].setValue('http://123.123.123:123')
    component.validateForm.controls['envValue'].setValue('123')
    component.saveCluster()
    fixture.detectChanges()

    expect(spyClusterMarkAsDirty).not.toHaveBeenCalled()
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('faild')
  })

  it('##saveCluster failed without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'post')
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
    const spyClusterMarkAsDirty = jest.spyOn(component.validateForm.controls['clusterAddr'], 'markAsDirty')

    component.validateForm.controls['envValue'].setValue('PRO')
    component.validateForm.controls['clusterName'].setValue('test123_')
    component.validateForm.controls['clusterAddr'].setValue('http://123.123.123:123')
    component.validateForm.controls['envValue'].setValue('123')
    component.saveCluster()
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('新建集群失败！')
    expect(spyClusterMarkAsDirty).not.toHaveBeenCalled()
  })

  it('##cancel test', () => {
    // @ts-ignore
    const spyOnRouterNavigate = jest.spyOn(component.router, 'navigate')
    expect(spyOnRouterNavigate).not.toHaveBeenCalled()

    component.cancel()
    expect(spyOnRouterNavigate).toHaveBeenCalled()
    expect(spyOnRouterNavigate).toHaveBeenCalledWith(['/', 'deploy', 'cluster'])
  })
})
