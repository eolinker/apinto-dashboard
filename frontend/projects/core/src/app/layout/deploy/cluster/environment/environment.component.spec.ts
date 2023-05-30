/* eslint-disable dot-notation */
import { ComponentFixture, fakeAsync, flush, TestBed } from '@angular/core/testing'
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
import { EoNgFeedbackDrawerService, EoNgFeedbackModalService, EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { Subject } from 'rxjs/internal/Subject'
import { of } from 'rxjs'
import { DeployClusterEnvironmentComponent } from './environment.component'
import { ReactiveFormsModule } from '@angular/forms'

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

class MockEnsureService {
  create () {
    return 'modal is create'
  }
}

describe('DeployClusterEnvironmentComponent test', () => {
  let component: DeployClusterEnvironmentComponent
  let fixture: ComponentFixture<DeployClusterEnvironmentComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule,
        NzDrawerModule, NzOutletModule, HttpClientModule, ReactiveFormsModule,
        RouterModule.forRoot([
          {
            path: '',
            component: DeployClusterEnvironmentComponent
          },
          {
            path: 'deploy/cluster/content/cert',
            component: DeployClusterEnvironmentComponent
          },
          {
            path: 'nodes',
            component: DeployClusterNodesComponent
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
        { provide: EoNgFeedbackDrawerService, useClass: MockDrawerService },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: EoNgFeedbackModalService, useClass: MockEnsureService }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(DeployClusterEnvironmentComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('##should create', () => {
    expect(component).toBeTruthy()
  })

  it('##btn click', () => {
    const spyOpenDrawer = jest.spyOn(component, 'openDrawer')
    const spyDelete = jest.spyOn(component, 'delete')
    expect(component.configsTableBody[0].styleFn({ publish: 'DEFECT' })).toStrictEqual('background:#f9f9f9')
    expect(component.configsTableBody[0].styleFn({ publish: 'test' })).toStrictEqual('')

    expect(component.configsTableBody[1].styleFn({ publish: 'DEFECT' })).toStrictEqual('background:#f9f9f9')
    expect(component.configsTableBody[1].styleFn({ publish: 'test' })).toStrictEqual('')

    expect(component.configsTableBody[2].styleFn({ publish: 'DEFECT' })).toStrictEqual('background:#f9f9f9')
    expect(component.configsTableBody[2].styleFn({ publish: 'test' })).toStrictEqual('')

    expect(component.configsTableBody[3].styleFn({ publish: 'DEFECT' })).toStrictEqual('background:#f9f9f9; text-align:center')
    expect(component.configsTableBody[3].styleFn({ publish: 'test' })).toStrictEqual('')

    expect(component.configsTableBody[4].styleFn({ publish: 'DEFECT' })).toStrictEqual('background:#f9f9f9')
    expect(component.configsTableBody[4].styleFn({ publish: 'test' })).toStrictEqual('')

    expect(component.configsTableBody[5].styleFn({ publish: 'DEFECT' })).toStrictEqual('background:#f9f9f9')
    expect(component.configsTableBody[5].styleFn({ publish: 'test' })).toStrictEqual('')

    expect(spyDelete).not.toHaveBeenCalled()
    expect(spyOpenDrawer).not.toHaveBeenCalled()
    expect(component.configsTableBody[6].showFn({ publish: 'DEFECT' })).toStrictEqual(false)
    expect(component.configsTableBody[6].showFn({ publish: 'test' })).toStrictEqual(true)
    expect(component.configsTableBody[6].btns[0].click({ publish: 'test' }))
    expect(spyDelete).not.toHaveBeenCalled()
    expect(spyOpenDrawer).toHaveBeenCalledTimes(1)
    expect(component.configsTableBody[6].btns[1].click({ publish: 'test' }))
    expect(spyDelete).toHaveBeenCalledTimes(1)
    expect(spyOpenDrawer).toHaveBeenCalledTimes(1)

    expect(component.configsTableBody[7].showFn({ publish: 'DEFECT' })).toStrictEqual(true)
    expect(component.configsTableBody[7].btns[0].click({ publish: 'test' }))
    expect(spyDelete).toHaveBeenCalledTimes(1)
    expect(spyOpenDrawer).toHaveBeenCalledTimes(2)
    expect(component.configsTableBody[7].showFn({ publish: 'test' })).toStrictEqual(false)
  })

  it('##should initial configsList', fakeAsync(() => {
    expect(component.configsList).not.toBe([])
  }))

  it('##ngOnDestroy test', () => {
    // @ts-ignore
    const spySubscription = jest.spyOn(component.subscription, 'unsubscribe')
    // @ts-ignore
    expect(spySubscription).not.toHaveBeenCalled()
    component.ngOnDestroy()
    expect(spySubscription).toHaveBeenCalled()
  })

  // it('detect', fakeAsync(() => {
  //   const form = { value: '', detect: { value: '', show: false } }
  //   expect(form.value).toBe('')
  //   expect(form.detect.value).toBe('')
  //   component.detect(form)
  //   tick(100)
  //   flush()
  //   expect(form.value).toBe('')
  //   expect(form.detect.value).toBe('')

  //   const form2 = { value: 'test', detect: { value: '', show: false } }
  //   expect(form2.value).toBe('test')
  //   expect(form2.detect.value).toBe('')
  //   component.detect(form2)
  //   tick(100)
  //   flush()
  //   expect(form2.value).toBe('test')
  //   expect(form2.detect.value).toBe('{test}')

  //   const form3 = { value: 'test test', detect: { value: '', show: false } }
  //   expect(form3.value).toBe('test test')
  //   expect(form3.detect.value).toBe('')
  //   component.detect(form3)
  //   tick(100)
  //   flush()
  //   expect(form3.value).toBe('test test')
  //   expect(form3.detect.value).toBe('{test<span class="detected-symbol">#空格#</span>test}')

  //   const form4 = { value: 'test t\nest', detect: { value: '', show: false } }
  //   component.detect(form4)
  //   tick(100)
  //   flush()
  //   expect(form4.detect.value).toBe('{test<span class="detected-symbol">#空格#</span>t<span class="detected-symbol">#换行符#</span>est}')
  // }))

  it('getConfigsList', () => {
    component.clusterName = 'gd_pro'
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { variables: ['test'] } }))
    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    component.getConfigsList()

    expect(spyService).toHaveBeenCalled()
    expect(component.configsList).toStrictEqual(['test'])

    const spyService2 = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 2, data: { } }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')

    component.getConfigsList()

    expect(spyMessageError).toBeCalledTimes(1)
    expect(spyService2).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('deleteConfig', () => {
    component.clusterName = 'gd_pro'
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: 0, data: { variables: ['test'] } }))
    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const item = { key: '123' }
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    component.deleteConfig(item)
    expect(spyService).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)

    const spyService2 = jest.spyOn(httpCommonService, 'delete').mockReturnValue(of({ code: 2, data: { } }))

    component.deleteConfig(item)

    expect(spyService2).toHaveBeenCalled()
    expect(spyMessageError).toHaveBeenCalledTimes(1)
  })

  it('click table btns', () => {
    const spyOpendrawer = jest.spyOn(component, 'openDrawer')
    // const spyDeleteCert = jest.spyOn(component, 'deleteCert')
    // @ts-ignore
    const spyModal = jest.spyOn(component.modalService, 'create')
    expect(spyOpendrawer).toHaveBeenCalledTimes(0)
    expect(spyModal).toHaveBeenCalledTimes(0)

    const item = { key: 'test' }
    component.configsTableBody[6].btns[0].click(item)
    expect(spyOpendrawer).toHaveBeenCalledTimes(1)
    component.configsTableBody[6].btns[1].click(item)
    expect(spyModal).toHaveBeenCalledTimes(1)
    component.configsTableBody[7].btns[0].click(item)
    expect(spyOpendrawer).toHaveBeenCalledTimes(2)
  })

  it('openDrawer without http request', fakeAsync(() => {
    expect(component.drawerRef).toBeUndefined()
    component.openDrawer('addConfig')
    fixture.detectChanges()
    flush()
    expect(component.drawerRef).not.toBeUndefined()

    component.drawerRef = undefined
    component.openDrawer('editConfig', { id: 0 })
    fixture.detectChanges()
    flush()
    expect(component.drawerRef).not.toBeUndefined()
  }))

  it('openDrawer with success http request', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get')
      .mockReturnValue(of({
        code: 0,
        data: {
          info: {
            clusters: [{ test: 'test1' }],
            variables: [{ test: 'test2' }]
          },
          version_name: 'test_version_name',
          source: 'test_source',
          unpublish_msg: 'unpublish_msg'
        }
      }))

    const isget = httpCommonService.get('') !== null

    expect(component.drawerRef).toBeUndefined()
    component.openDrawer('updateConfig', { id: 0 })
    expect(spyService).toHaveBeenCalledTimes(2)
    expect(isget).toBe(true)
    fixture.detectChanges()
    expect(component.drawerRef).not.toBeUndefined()
    expect(component.clustersList).toStrictEqual([{ test: 'test1' }])
    expect(component.updateConfigsList).toStrictEqual([{ test: 'test2' }])

    component.drawerRef = undefined
    component.openDrawer('operateRecords', { id: 0 })
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalledTimes(3)
    expect(component.drawerRef).not.toBeUndefined()
    expect(component.operateRecordsData).toStrictEqual({
      info: {
        clusters: [{ test: 'test1' }],
        variables: [{ test: 'test2' }]
      },
      version_name: 'test_version_name',
      source: 'test_source',
      unpublish_msg: 'unpublish_msg'
    })

    component.drawerRef = undefined
    component.openDrawer('publishRecords', { id: 0 })
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalledTimes(4)
    expect(component.drawerRef).not.toBeUndefined()
    expect(component.publishRecordsData).toStrictEqual({
      info: {
        clusters: [{ test: 'test1' }],
        variables: [{ test: 'test2' }]
      },
      version_name: 'test_version_name',
      source: 'test_source',
      unpublish_msg: 'unpublish_msg'
    })

    component.drawerRef = undefined
    component.openDrawer('publish', { id: 0 })
    expect(spyService).toHaveBeenCalledTimes(5)
    fixture.detectChanges()
    expect(component.drawerRef).not.toBeUndefined()
    expect(component.publishData).toStrictEqual({
      info: {
        clusters: [{ test: 'test1' }],
        variables: [{ test: 'test2' }]
      },
      version_name: 'test_version_name',
      source: 'test_source',
      variables: [],
      unpublish_msg: 'unpublish_msg'
    })
    expect(component.validatePublishForm.controls['version_name'].value).toStrictEqual('test_version_name')
    expect(component.unpublish_msg).toStrictEqual('unpublish_msg')
  })

  it('openDrawer with fail http request', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get')
      .mockReturnValue(of({
        code: -1,
        data: {}
      }))

    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')

    expect(spyMessageError).not.toHaveBeenCalled()

    expect(component.drawerRef).toBeUndefined()
    component.openDrawer('updateConfig', { id: 0 })
    expect(spyService).toHaveBeenCalledTimes(2)
    expect(isget).toBe(true)
    fixture.detectChanges()
    expect(component.drawerRef).not.toBeUndefined()
    expect(spyMessageError).toBeCalledTimes(1)

    component.drawerRef = undefined
    component.openDrawer('operateRecords', { id: 0 })
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalledTimes(3)
    expect(component.drawerRef).not.toBeUndefined()
    expect(spyMessageError).toBeCalledTimes(2)

    component.drawerRef = undefined
    component.openDrawer('publishRecords', { id: 0 })
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalledTimes(4)
    expect(component.drawerRef).not.toBeUndefined()
    expect(spyMessageError).toBeCalledTimes(3)

    component.drawerRef = undefined
    component.openDrawer('publish', { id: 0 })
    expect(spyService).toHaveBeenCalledTimes(5)
    fixture.detectChanges()
    expect(component.drawerRef).not.toBeUndefined()
    expect(spyMessageError).toBeCalledTimes(4)
  })

  it('save with success', () => {
    component.clusterName = 'gd_pro'
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: 0, data: { certicates: ['test'] } }))
    const spyService2 = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 0, data: { certicates: ['test'] } }))
    const spyGetConfigsList = jest.spyOn(component, 'getConfigsList')

    component.openDrawer('addConfig', { id: 0 })
    component.drawerRef!.close = () => {
      return 'drawer is close'
    }

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessageError).not.toHaveBeenCalled()

    const spyFn = jest.spyOn(component.drawerRef!, 'close')
    expect(spyFn).not.toHaveBeenCalled()
    component.validateAddConfigForm.controls['key'].setValue('test_')
    component.save('addConfig')
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
    expect(spyFn).toHaveBeenCalled()
    expect(spyGetConfigsList).toHaveBeenCalledTimes(1)

    component.clustersList = [{ name: 'c1', env: 'default_value', status: 'test', id: 0, checked: true }, { name: 'c2', env: 'default_value', id: 0, status: 'test', checked: true }, { name: 'x3', env: 'default_value', id: 0, status: 'test', checked: true }]
    component.updateConfigsList = [{ key: 'c1', value: 'default_value', variable_id: 0, publish: 'test', status: 'test', desc: 'test', operator: 'test', updateTime: 'test', createTime: 'test', id: 0, checked: true }, { key: 'c2', value: 'default_value', variable_id: 0, publish: 'test', id: 0, status: 'test', desc: 'test', operator: 'test', updateTime: 'test', createTime: 'test', checked: true }, { key: 'x3', value: 'default_value', variable_id: 0, publish: 'test', status: 'test', desc: 'test', operator: 'test', updateTime: 'test', createTime: 'test', id: 0, checked: true }]

    component.save('updateConfig')
    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(2)

    component.validatePublishForm.controls['version_name'].setValue('test')
    component.publishData.is_publish = true
    component.save('publish')
    expect(spyMessageSuccess).toHaveBeenCalledTimes(3)
    expect(spyService).toHaveBeenCalledTimes(3)
    expect(spyGetConfigsList).toHaveBeenCalledTimes(2)

    component.validateAddConfigForm.controls['key'].setValue('test_')
    component.save('editConfig')
    expect(spyService2).toHaveBeenCalledTimes(1)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(4)
    expect(spyGetConfigsList).toHaveBeenCalledTimes(3)

    component.save('test')
    expect(spyService2).toHaveBeenCalledTimes(1)
    expect(spyService).toHaveBeenCalledTimes(3)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(4)
    expect(spyGetConfigsList).toHaveBeenCalledTimes(3)
  })

  it('save with fail', () => {
    component.clusterName = 'gd_pro'
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: 1, data: { certicates: ['test'] } }))
    const spyService2 = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 1, data: { certicates: ['test'] } }))

    component.validateAddConfigForm.controls['key'].setValue('test_')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    component.save('addConfig')
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalledTimes(1)

    component.clustersList = [{ name: 'c1', env: 'default_value', status: 'test', id: 0, checked: true }, { name: 'c2', env: 'default_value', id: 0, status: 'test', checked: true }, { name: 'x3', env: 'default_value', id: 0, status: 'test', checked: true }]
    component.updateConfigsList = [{ key: 'c1', value: 'default_value', variable_id: 0, publish: 'test', status: 'test', desc: 'test', operator: 'test', updateTime: 'test', createTime: 'test', id: 0, checked: true }, { key: 'c2', value: 'default_value', variable_id: 0, publish: 'test', id: 0, status: 'test', desc: 'test', operator: 'test', updateTime: 'test', createTime: 'test', checked: true }, { key: 'x3', value: 'default_value', variable_id: 0, publish: 'test', status: 'test', desc: 'test', operator: 'test', updateTime: 'test', createTime: 'test', id: 0, checked: true }]

    component.save('updateConfig')
    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalledTimes(2)

    component.validatePublishForm.controls['version_name'].setValue('test')
    component.publishData.is_publish = true
    component.save('publish')
    expect(spyService).toHaveBeenCalledTimes(3)
    expect(spyMessage).toHaveBeenCalledTimes(3)

    component.save('editConfig')
    expect(spyService2).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalledTimes(4)

    component.save('test')
    expect(spyService2).toHaveBeenCalledTimes(1)
    expect(spyService).toHaveBeenCalledTimes(3)
    expect(spyMessage).toHaveBeenCalledTimes(4)
  })

  it('clearForm', fakeAsync(() => {
    component.updateConfigForm.clusters = [{ name: 'test', env: 'test', id: 0 }]
    component.updateConfigForm.variables = [{ key: 'test', value: 'test', variable_id: 0, id: 0 }]

    component.clearUpdateForm()
    fixture.detectChanges()

    expect(component.updateConfigForm.clusters).toStrictEqual([])
    expect(component.updateConfigForm.variables).toStrictEqual([])
    for (const index in component.clustersList) {
      expect(component.clustersList[index].checked).toStrictEqual(false)
    }
    for (const index in component.updateConfigsList) {
      expect(component.updateConfigsList[index].checked).toStrictEqual(false)
    }
  }))
})
