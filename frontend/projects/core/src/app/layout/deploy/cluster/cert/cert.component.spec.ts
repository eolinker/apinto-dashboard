/* eslint-disable dot-notation */
import { DeployClusterCertComponent } from './cert.component'
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
import { DeployClusterEnvironmentComponent } from '../environment/environment.component'
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
import { Observable, of } from 'rxjs'
import { ReactiveFormsModule } from '@angular/forms'
import { EoNgButtonModule } from 'eo-ng-button'

class MockDrawerService {
  result: boolean = false

  nzAfterClose = new Subject<any>()

  create() {
    return {
      afterClose: new Observable(),
      close: () => {
        return 'drawer is close'
      }
    }
  }
}

class MockMessageService {
  success() {
    return 'success'
  }

  error() {
    return 'error'
  }
}

class MockEnsureService {
  create() {
    return 'modal is create'
  }
}

describe('DeployClusterCertComponent test', () => {
  let component: DeployClusterCertComponent
  let fixture: ComponentFixture<DeployClusterCertComponent>

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
        EoNgButtonModule,
        RouterModule.forRoot([
          {
            path: '',
            component: DeployClusterEnvironmentComponent
          },
          {
            path: 'cert',
            component: DeployClusterCertComponent
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

    fixture = TestBed.createComponent(DeployClusterCertComponent)
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

  it('##getCertsList', () => {
    component.clusterName = 'gd_pro'
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: 0, data: { certicates: ['test'] } }))
    component.getCertsList()

    expect(spyService).toHaveBeenCalled()
    expect(component.certsList).not.toStrictEqual([])

    const spyService2 = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: 2, data: {} }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')

    component.getCertsList()

    expect(spyService2).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##click table btns', () => {
    const spyOpendrawer = jest.spyOn(component, 'openDrawer')
    // @ts-ignore
    const spyModal = jest.spyOn(component.modalService, 'create')
    expect(spyOpendrawer).toHaveBeenCalledTimes(0)
    expect(spyModal).toHaveBeenCalledTimes(0)

    const item = { key: 'test' }
    component.certsTableBody[4].btns[0].click(item)
    expect(spyOpendrawer).toHaveBeenCalledTimes(1)
    component.certsTableBody[4].btns[1].click(item)
    expect(spyModal).toHaveBeenCalledTimes(1)
  })

  it('openDrawer', fakeAsync(() => {
    component.openDrawer('test')

    expect(component.drawerRef).toBeUndefined()
    component.openDrawer('addCert')
    fixture.detectChanges()
    flush()
    expect(component.drawerRef).not.toBeUndefined()

    component.drawerRef = undefined
    component.openDrawer('editCert', { id: 0 })
    fixture.detectChanges()
    flush()
    expect(component.drawerRef).not.toBeUndefined()
  }))

  it('##save test', () => {
    component.clusterName = 'gd_pro'
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'post')
      .mockReturnValue(of({ code: 0, data: { certicates: ['test'] } }))
    const spyService2 = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(of({ code: 0, data: { certicates: ['test'] } }))
    const spyFn = jest.spyOn(component, 'getCertsList')

    component.save('test')

    expect(spyFn).not.toHaveBeenCalled()
    expect(spyService).not.toHaveBeenCalled()

    component.save('addCert')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyFn).not.toHaveBeenCalled()

    component.save('editCert')
    expect(spyService2).not.toHaveBeenCalled()
    expect(spyFn).not.toHaveBeenCalled()

    component.validateForm.controls['key'].setValue('key')
    component.validateForm.controls['pem'].setValue('pem')
    component.save('addCert')
    expect(spyService).toHaveBeenCalled()
    expect(spyFn).toHaveBeenCalledTimes(1)

    component.save('editCert')
    expect(spyService2).toHaveBeenCalled()
    expect(spyFn).toHaveBeenCalledTimes(2)
  })

  it('save with fail request', () => {
    component.clusterName = 'test'
    component.validateForm.controls['key'].setValue('key')
    component.validateForm.controls['pem'].setValue('pem')
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'post')
      .mockReturnValue(of({ code: 1, data: { certicates: ['test'] } }))
    const spyService3 = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(of({ code: 1, data: { certicates: ['test'] } }))
    expect(spyService).not.toHaveBeenCalled()
    expect(spyService3).not.toHaveBeenCalled()

    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    component.save('addCert')
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalledWith('添加证书失败！')

    component.save('editCert')
    expect(spyService3).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalledTimes(2)
  })

  it('##deleteCert with success request', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'delete')
      .mockReturnValue(of({ code: 0, data: {} }))
    const spyGetCertList = jest.spyOn(component, 'getCertsList')

    expect(spyGetCertList).not.toHaveBeenCalled()
    // @ts-ignore
    component.deleteCert({ row: {} })
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalledWith('cluster//certificate/')

    component.clusterName = 'gd_pro'

    component.deleteCert({
      id: 0,
      name: '',
      operator: '',
      update_time: '',
      valid_time: '',
      create_time: ''
    })
    expect(spyService).toBeCalledTimes(2)
    expect(spyGetCertList).toHaveBeenCalled()
  })

  it('deleteCert with fail request', () => {
    component.clusterName = 'gd_pro'

    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'delete')
      .mockReturnValue(of({ code: -1, data: {} }))
    const isdelete = httpCommonService.delete('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')

    expect(spyMessage).not.toHaveBeenCalled()
    component.deleteCert({
      id: 0,
      name: '',
      operator: '',
      update_time: '',
      valid_time: '',
      create_time: ''
    })
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(isdelete).toBe(true)
  })

  it('##encode', fakeAsync(() => {
    expect(component.encode('test')).toStrictEqual(
      Buffer.from('test').toString('base64')
    )
    expect(component.encode('AAA')).not.toStrictEqual('AAA')
  }))

  it('##disabledEdit', fakeAsync(() => {
    component.disabledEdit(true)
    expect(component.nzDisabled).toStrictEqual(true)
    component.disabledEdit(false)
    expect(component.nzDisabled).toStrictEqual(false)
  }))

  it('##certTableClick', () => {
    const spyOpendrawer = jest.spyOn(component, 'openDrawer')
    expect(spyOpendrawer).not.toHaveBeenCalled()

    component.certTableClick({ data: 'test' })
    expect(spyOpendrawer).toHaveBeenCalledTimes(1)
    expect(spyOpendrawer).toHaveBeenCalledWith('editCert', 'test')
  })
})
