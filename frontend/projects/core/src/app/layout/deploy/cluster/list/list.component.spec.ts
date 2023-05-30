import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { environment } from 'projects/core/src/environments/environment'
import { RouterModule } from '@angular/router'
import { DeployClusterNodesComponent } from '../nodes/nodes.component'
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
import { DeployClusterListComponent } from './list.component'
import { ElementRef } from '@angular/core'
import { ReactiveFormsModule } from '@angular/forms'
import { DeployClusterCreateComponent } from '../create/create.component'

class MockDrawerService {
  result: boolean = false

  nzAfterClose = new Subject<any>()

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

class MockElementRef extends ElementRef {
  constructor () {
    super(null)
  }
}

class MockEnsureService {
  create () {
    return 'modal is create'
  }
}

describe('DeployClusterListComponent test', () => {
  let component: DeployClusterListComponent
  let fixture: ComponentFixture<DeployClusterListComponent>

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
            component: DeployClusterListComponent
          },
          {
            path: 'deploy/cluster/content/cert',
            component: DeployClusterListComponent
          },
          {
            path: 'deploy/cluster/create',
            component: DeployClusterCreateComponent
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
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: EoNgFeedbackDrawerService, useClass: MockDrawerService },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: EoNgFeedbackModalService, useClass: MockEnsureService }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(DeployClusterListComponent)

    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('should initial configsList and environmentList', fakeAsync(() => {
    expect(component.environmentList).not.toBe([])
    expect(component.clustersList).not.toBe([])
  }))

  it('click table btns', () => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')
    // const spyDeleteCert = jest.spyOn(component, 'deleteCert')
    // @ts-ignore
    const spyModalDelete = jest.spyOn(component, 'delete')
    expect(spyRouter).toHaveBeenCalledTimes(0)
    expect(spyModalDelete).toHaveBeenCalledTimes(0)

    const item = { key: 'test' }
    component.clustersTableBody[3].btns[0].click(item)
    expect(spyRouter).toHaveBeenCalledTimes(1)
    component.clustersTableBody[3].btns[1].click(item)
    expect(spyModalDelete).toHaveBeenCalledTimes(1)
    expect(component.clustersTableBody[3].btns[1].disabledFn()).toStrictEqual(
      false
    )
    component.nzDisabled = true
    expect(component.clustersTableBody[3].btns[1].disabledFn()).toStrictEqual(
      true
    )
  })

  it('deleteCluster', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'delete')
      .mockReturnValue(of({ code: 0, data: {} }))
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetClustersData = jest.spyOn(component, 'getClustersData')
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetClustersData).not.toHaveBeenCalled()
    component.deleteCluster({ name: 'testName' })

    expect(spyService).toHaveBeenCalled()
    expect(spyGetClustersData).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()

    const spyService2 = jest
      .spyOn(httpCommonService, 'delete')
      .mockReturnValue(of({ code: 2, data: {} }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    component.deleteCluster({ name: 'testName' })

    expect(spyService2).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('getClustersData', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: 0, data: { clusters: [1, 2, 3] } }))
    component.getClustersData()

    expect(spyService).toHaveBeenCalled()
    expect(component.clustersList).toStrictEqual([1, 2, 3])

    const spyService2 = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: 2, data: {} }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    component.getClustersData()

    expect(spyService2).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##disabledEdit', fakeAsync(() => {
    component.disabledEdit(true)
    expect(component.nzDisabled).toStrictEqual(true)
    component.disabledEdit(false)
    expect(component.nzDisabled).toStrictEqual(false)
  }))

  it('##clusterTableClick test', () => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')
    expect(spyRouter).not.toHaveBeenCalled()
    component.clusterTableClick({ data: { name: 'test' } })
    expect(spyRouter).toHaveBeenCalled()
    expect(spyRouter).toHaveBeenCalledWith(
      ['/', 'deploy', 'cluster', 'content'],
      {
        queryParams: { clusterName: 'test' }
      }
    )
  })

  it('##addCluster  test', () => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')
    expect(spyRouter).not.toHaveBeenCalled()
    component.addCluster()
    expect(spyRouter).toHaveBeenCalled()
    expect(spyRouter).toHaveBeenCalledWith(['/', 'deploy', 'cluster', 'create'])
  })
})
