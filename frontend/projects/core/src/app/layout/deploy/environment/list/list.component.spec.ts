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
import { DeployEnvironmentListComponent } from './list.component'
import { ReactiveFormsModule } from '@angular/forms'
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

describe('DeployEnvironmentListComponent test', () => {
  let component: DeployEnvironmentListComponent
  let fixture: ComponentFixture<DeployEnvironmentListComponent>

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
            component: DeployEnvironmentListComponent
          },
          {
            path: 'deploy/cluster/content/cert',
            component: DeployEnvironmentListComponent
          },
          {
            path: 'nodes',
            component: DeployEnvironmentListComponent
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

    fixture = TestBed.createComponent(DeployEnvironmentListComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('should initial globalEnvForms', fakeAsync(() => {
    expect(component.globalEnvForms.variables).not.toBe([])
  }))

  it('resetSearch', fakeAsync(() => {
    component.searchForm = {
      key: 'test',
      status: 'test'
    }
    component.resetSearch()
    fixture.detectChanges()
    flush()
    expect(component.searchForm.key).toStrictEqual('')
    expect(component.searchForm.status).toStrictEqual('')
  }))

  it('##click table btns', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get')
    const spyDeleteModal = jest.spyOn(component, 'deleteModal')
    expect(spyService).toHaveBeenCalledTimes(0)
    expect(spyDeleteModal).toHaveBeenCalledTimes(0)

    const item = { key: 'test' }
    component.globalEnvTableBody[5].btns[0].click(item)
    expect(spyService).toHaveBeenCalledTimes(1)

    expect(component.globalEnvTableBody[5].btns[1].disabledFn()).toStrictEqual(
      false
    )
    component.globalEnvTableBody[5].btns[1].click(item)
    expect(spyDeleteModal).toHaveBeenCalledTimes(1)
  })

  it('##search', () => {
    const spyFn = jest.spyOn(component, 'getVariables')
    expect(spyFn).toHaveBeenCalledTimes(0)

    component.search()
    expect(spyFn).toHaveBeenCalledTimes(1)
  })

  it('##openDrawer without http request', fakeAsync(() => {
    expect(component.editConfigDrawerRef).toBeUndefined()
    component.openDrawer('view')
    fixture.detectChanges()
    flush()
    expect(component.editConfigDrawerRef).not.toBeUndefined()
  }))

  it('##deleteModal test', fakeAsync(() => {
    // @ts-ignore
    const spyModal = jest.spyOn(component.modalService, 'create')
    expect(spyModal).not.toHaveBeenCalled()
    component.deleteModal('test')
    expect(spyModal).toHaveBeenCalled()
  }))

  it('##getVariables successfully', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: { test: 'test', total: 1 },
        msg: 'success'
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.getVariables()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component.globalEnvForms).toStrictEqual({ test: 'test', total: 1 })
    expect(component.variablePage.total).toStrictEqual(1)
  })

  it('##getVariables failed', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: -1,
        data: {},
        msg: 'faild'
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.getVariables()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('faild')
  })

  it('##getVariables failed without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: -1,
        data: {}
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.getVariables()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取列表数据失败！')
  })

  it('##delete successfully', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(
      of({
        code: 0,
        data: { variables: [1, 2], total: 1 },
        msg: 'success'
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetVariables = jest.spyOn(component, 'getVariables')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetVariables).not.toHaveBeenCalled()
    component.delete({ key: 'test' })
    fixture.detectChanges()

    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyGetVariables).toHaveBeenCalled()
  })

  it('##delete failed', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(
      of({
        code: -1,
        data: {},
        msg: 'faild'
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyGetVariables = jest.spyOn(component, 'getVariables')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyGetVariables).not.toHaveBeenCalled()
    component.delete({ key: 'test' })
    fixture.detectChanges()

    expect(spyGetVariables).not.toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('faild')
  })

  it('##delete failed without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(
      of({
        code: -1,
        data: {}
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.delete({ key: 'test' })
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('删除失败！')
  })
})
