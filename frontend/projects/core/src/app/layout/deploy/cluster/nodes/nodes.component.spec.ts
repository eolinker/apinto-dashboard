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
import { DeployClusterNodesComponent } from './nodes.component'
import { ReactiveFormsModule } from '@angular/forms'

class MockDrawerService {
  result: boolean = false;

  nzAfterClose = new Subject<any>();

  create () {
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

describe('DeployClusterNodesComponent test', () => {
  let component: DeployClusterNodesComponent
  let fixture: ComponentFixture<DeployClusterNodesComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
        BidiModule,
        NoopAnimationsModule,
        NzNoAnimationModule,
        NzDrawerModule,
        NzOutletModule,
        ReactiveFormsModule,
        HttpClientModule,
        RouterModule.forRoot([
          {
            path: '',
            component: DeployClusterNodesComponent
          },
          {
            path: 'deploy/cluster/content/cert',
            component: DeployClusterNodesComponent
          },
          {
            path: 'nodes',
            component: DeployClusterNodesComponent
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
        {
          provide: EoNgFeedbackModalService,
          useClass: MockEnsureService
        }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(DeployClusterNodesComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('should initial configsList', fakeAsync(() => {
    expect(component.nodesForms.nodes).not.toBe([])
  }))

  it('getNodeslist', () => {
    component.clusterName = 'gd_pro'
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: 0, data: { variables: ['test'] } }))
    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    component.getNodeslist()

    expect(spyService).toHaveBeenCalled()
    expect(component.nodesForms).toStrictEqual({ variables: ['test'] })

    const spyService2 = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: 2, data: {} }))

    component.getNodeslist()

    expect(spyMessageError).toBeCalledTimes(1)
    expect(spyService2).toHaveBeenCalled()
  })

  it('updateNodes', () => {
    component.clusterName = 'gd_pro'
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(of({ code: 0, data: { variables: ['test'] } }))
    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')

    const spyGetNodesList = jest.spyOn(component, 'getNodeslist')
    component.updateNodes()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyGetNodesList).toHaveBeenCalled()

    const spyService2 = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(of({ code: 2, data: {} }))

    component.updateNodes()

    expect(spyMessageError).toBeCalledTimes(1)
    expect(spyService2).toHaveBeenCalled()
  })

  it('openDrawer without http request', fakeAsync(() => {
    expect(component.nodesDrawerRef).toBeUndefined()
    component.openDrawer('resetNodes')
    fixture.detectChanges()
    flush()
    expect(component.nodesDrawerRef).not.toBeUndefined()
  }))

  it('openDrawer, form value reset when drawer closed', fakeAsync(() => {
    expect(component.nodesDrawerRef).toBeUndefined()
    component.openDrawer('resetNodes')
    fixture.detectChanges()
    flush()
    expect(component.nodesDrawerRef).not.toBeUndefined()
  }))

  it('save', () => {
    component.clusterName = 'gd_pro'
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'post')
      .mockReturnValue(of({ code: 0, data: {} }))
    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    expect(spyService).toHaveBeenCalledTimes(0)

    component.validateResetNodeForm.controls['clusterAddr'].setValue('http')
    component.save()
    expect(spyService).toHaveBeenCalledTimes(0)

    component.validateResetNodeForm.controls['clusterAddr'].setValue('http://123.123.123:123')
    component.save()
    expect(spyService).toHaveBeenCalledTimes(1)

    const spyService2 = jest
      .spyOn(httpCommonService, 'post')
      .mockReturnValue(of({ code: -1, data: {} }))

    expect(spyMessageError).not.toHaveBeenCalled()
    component.save()
    expect(spyService2).toHaveBeenCalled()
    expect(spyMessageError).toHaveBeenCalled()
  })

  it('testCluster', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(
        of({
          code: 0,
          data: { nodes: [], source: 'testSource', isUpdate: true }
        })
      )
    expect(spyService).toHaveBeenCalledTimes(0)

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')

    component.validateResetNodeForm.controls['clusterAddr'].setValue(null)

    component.testCluster()
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalledTimes(0)

    component.validateResetNodeForm.controls['clusterAddr'].setValue('http://123.123.123:123')

    component.testCluster()
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(component.nodesTestList).toStrictEqual([])
    expect(component.clusterCanBeCreated).toStrictEqual(true)
    expect(component.testPassAddr).toStrictEqual('http://123.123.123:123')
    expect(component.nodesTestTableShow).toStrictEqual(false)

    const spyService2 = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(
        of({
          code: 0,
          data: { nodes: [1, 2, 3], source: 'testSource', isUpdate: true }
        })
      )
    expect(spyService2).toHaveBeenCalledTimes(1)

    expect(component.testPassAddr).toStrictEqual('http://123.123.123:123')

    component.testCluster()
    fixture.detectChanges()
    expect(spyService2).toHaveBeenCalledTimes(2)
    expect(component.nodesTestList).toStrictEqual([1, 2, 3])
    expect(component.clusterCanBeCreated).toStrictEqual(true)
    expect(component.nodesTestTableShow).toStrictEqual(true)

    expect(spyMessageError).not.toHaveBeenCalled()

    const spyService3 = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, data: {} }))
    expect(spyService3).toHaveBeenCalledTimes(2)

    expect(component.testPassAddr).toStrictEqual('http://123.123.123:123')

    component.testCluster()
    fixture.detectChanges()
    expect(spyService3).toHaveBeenCalledTimes(3)
    expect(spyMessageError).toHaveBeenCalled()
  })
})
