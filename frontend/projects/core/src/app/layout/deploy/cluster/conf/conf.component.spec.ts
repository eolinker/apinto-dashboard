/* eslint-disable camelcase */
import {
  ComponentFixture,
  fakeAsync,
  TestBed
} from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { environment } from 'projects/core/src/environments/environment'
import { ActivatedRoute, RouterModule } from '@angular/router'
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
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgInputModule } from 'eo-ng-input'
import { ReactiveFormsModule } from '@angular/forms'
import { DeployClusterConfComponent } from './conf.component'

class MockDrawerService {
  result: boolean = false

  nzAfterClose = new Subject<any>()

  create () {
    return {
      afterClose: new Observable(),
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

describe('DeployClusterConfComponent test', () => {
  let component: DeployClusterConfComponent
  let fixture: ComponentFixture<DeployClusterConfComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
        BidiModule,
        EoNgSwitchModule,
        EoNgInputModule,
        ReactiveFormsModule,
        NoopAnimationsModule,
        NzNoAnimationModule,
        NzDrawerModule,
        NzOutletModule,
        HttpClientModule,
        RouterModule.forRoot([
          {
            path: '',
            component: DeployClusterEnvironmentComponent
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
        },
        {
          provide: ActivatedRoute,
          useValue: {
            queryParams: of({ cluster_name: 'clus2' })
          }
        }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(DeployClusterConfComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('##ngOninit should call getRedisList', () => {
    const spyGetRedisList = jest.spyOn(component, 'getRedisList')
    expect(spyGetRedisList).not.toHaveBeenCalled()
    component.ngOnInit()
    expect(spyGetRedisList).toHaveBeenCalledTimes(1)
  })

  it('##ngOnDestroy test', () => {
    // @ts-ignore
    const spySubscription = jest.spyOn(component.subscription, 'unsubscribe')
    // @ts-ignore
    expect(spySubscription).not.toHaveBeenCalled()
    component.ngOnDestroy()
    expect(spySubscription).toHaveBeenCalled()
  })

  it('##click table btns', () => {
    const spyCheckValueForEdit = jest.spyOn(component, 'checkValueForEdit')
    const spyTestAndSaveRedis = jest.spyOn(component, 'testAndSaveRedis')
    expect(spyCheckValueForEdit).not.toHaveBeenCalled()
    expect(spyTestAndSaveRedis).not.toHaveBeenCalled()

    const item :{
      username: string,
      addrs: string,
      password: string,
      origin_username:string,
      origin_addrs: string,
      origin_password: string,
      edit?:boolean
    } = {
      username: 'username',
      addrs: 'addr',
      password: 'password',
      origin_username: 'username',
      origin_addrs: 'addr',
      origin_password: 'password'
    }

    component.redisTableBody[0].change(item)
    expect(spyCheckValueForEdit).toHaveBeenCalledTimes(1)
    expect(spyCheckValueForEdit).toHaveBeenCalledWith(item)
    expect(spyTestAndSaveRedis).not.toHaveBeenCalled()

    component.redisTableBody[1].change(item)
    expect(spyCheckValueForEdit).toHaveBeenCalledTimes(2)
    expect(spyTestAndSaveRedis).not.toHaveBeenCalled()

    component.redisTableBody[2].change(item)
    expect(spyCheckValueForEdit).toHaveBeenCalledTimes(3)
    expect(spyTestAndSaveRedis).not.toHaveBeenCalled()

    component.redisTableBody[4].btns[0].click({ addrs: 'test' })
    expect(spyCheckValueForEdit).toHaveBeenCalledTimes(3)
    expect(spyTestAndSaveRedis).toHaveBeenCalledTimes(1)

    expect(component.redisTableBody[4].btns[0].disabledFn({ data: {} })).toStrictEqual(true)
    expect(component.redisTableBody[4].btns[0].disabledFn({ data: { addrs: 'test' } })).toStrictEqual(false)
  })

  it('##getRedisList successfully', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(
        of({
          code: 0,
          data: {
            redis: {
              addrs: 'addr1,addr2,addr3',
              username: 'username',
              password: 'password',
              enable: true,
              operator: 'operator'
            }
          }
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'

    component.getRedisList()
    fixture.detectChanges()

    expect(component.redisList).toStrictEqual([{
      addrs: 'addr1,addr2,addr3',
      username: 'username',
      password: 'password',
      enable: true,
      operator: 'operator',
      origin_addrs: 'addr1,addr2,addr3',
      origin_username: 'username',
      origin_password: 'password'
    }])
  })

  it('##getRedisList failed', () => {
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

    component.getRedisList()
    fixture.detectChanges()

    expect(component.redisList).toStrictEqual([{
      username: '', password: '', addrs: '', enable: null
    }
    ])

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('faild')
  })

  it('##getRedisList failed without msg', () => {
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
    component.getRedisList()
    fixture.detectChanges()

    expect(component.redisList).toStrictEqual([{
      username: '', password: '', addrs: '', enable: null
    }
    ])

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取配置列表失败！')
  })

  it('##startRedis successfully with msg', () => {
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
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    component.stopOrStartRedis({ enable: true })
    fixture.detectChanges()

    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toBeCalledWith('cluster/test/configuration/redis/enable')
    expect(spyMessageSuccess).toBeCalledWith('success')
  })

  it('##stopRedis successfully with msg', () => {
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
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    component.stopOrStartRedis({ enable: false })
    fixture.detectChanges()

    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toBeCalledWith('cluster/test/configuration/redis/disable')
    expect(spyMessageSuccess).toBeCalledWith('success')
  })

  it('##startRedis successfully without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(
        of({
          code: 0,
          data: {}
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    component.stopOrStartRedis({ enable: true })
    fixture.detectChanges()

    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toBeCalledWith('cluster/test/configuration/redis/enable')
    expect(spyMessageSuccess).toBeCalledWith('启用成功！')
  })

  it('##stopRedis successfully without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(
        of({
          code: 0,
          data: {}
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    component.stopOrStartRedis({ enable: false })
    fixture.detectChanges()

    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toBeCalledWith('cluster/test/configuration/redis/disable')
    expect(spyMessageSuccess).toBeCalledWith('禁用成功！')
  })

  it('##startRedis failed with msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(
        of({
          code: -1,
          data: {},
          msg: 'fail'
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    component.stopOrStartRedis({ enable: true })
    fixture.detectChanges()

    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toBeCalledWith('cluster/test/configuration/redis/enable')
    expect(spyMessage).toBeCalledWith('fail')
  })

  it('##startRedis failed without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(
        of({
          code: -1,
          data: {}
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    component.stopOrStartRedis({ enable: true })
    fixture.detectChanges()

    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toBeCalledWith('cluster/test/configuration/redis/enable')
    expect(spyMessage).toBeCalledWith('启用失败！')
  })

  it('##stopRedis failed with msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(
        of({
          code: -1,
          data: {},
          msg: 'fail'
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    component.stopOrStartRedis({ enable: false })
    fixture.detectChanges()

    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toBeCalledWith('cluster/test/configuration/redis/disable')
    expect(spyMessage).toBeCalledWith('fail')
  })

  it('##stopRedis failed without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(
        of({
          code: -1,
          data: {}
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.clusterName = 'test'
    const data = { enable: false }
    component.stopOrStartRedis(data)
    fixture.detectChanges()

    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toBeCalledWith('cluster/test/configuration/redis/disable')
    expect(spyMessage).toBeCalledWith('禁用失败！')
    expect(data.enable).toStrictEqual(true)
  })

  it('##testAndSaveRedis successfully', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
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
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetRedisList = jest.spyOn(component, 'getRedisList')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyGetRedisList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    component.clusterName = 'test'
    const body = { addrs: '111.111.111.11:111', username: 'test', password: 'test' }
    component.testAndSaveRedis(body)
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyGetRedisList).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalledWith('success')
  })

  it('##testAndSaveRedis successfully without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(
        of({
          code: 0,
          data: {
          }
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetRedisList = jest.spyOn(component, 'getRedisList')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyGetRedisList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()

    component.clusterName = 'test'
    const body = { addrs: '111.111.111.11:111', username: 'test', password: 'test' }

    component.testAndSaveRedis(body)
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyGetRedisList).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalledWith('配置成功！')
  })
  it('##testAndSaveRedis failed', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(
        of({
          code: -1,
          data: {
          },
          msg: 'failed'
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetRedisList = jest.spyOn(component, 'getRedisList')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyGetRedisList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()

    component.clusterName = 'test'
    const body = { addrs: '111.111.111.11:111', username: 'test', password: 'test' }

    component.testAndSaveRedis(body)
    fixture.detectChanges()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyGetRedisList).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalledWith('failed')
  })

  it('##testAndSaveRedis failed without msg', () => {
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
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetRedisList = jest.spyOn(component, 'getRedisList')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetRedisList).not.toHaveBeenCalled()

    component.clusterName = 'test'
    const body = { addrs: '111.111.111.11:111', username: 'test', password: 'test' }

    component.testAndSaveRedis(body)
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetRedisList).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalledWith('配置失败！')
  })

  it('##nzCheckAddRow return false', fakeAsync(() => {
    expect(component.nzCheckAddRow()).toStrictEqual(false)
  }))

  it('##checkValueForEdit test', fakeAsync(() => {
    const item :{
      username: string,
      addrs: string,
      password: string,
      origin_username:string,
      origin_addrs: string,
      origin_password: string,
      edit?:boolean
    } = {
      username: 'username',
      addrs: 'addrs',
      password: 'password',
      origin_username: 'username',
      origin_addrs: 'addrs',
      origin_password: 'password'
    }
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(false)

    item.username = 'username1'
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(true)

    item.username = 'username'
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(false)

    item.addrs = 'username'
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(true)

    item.username = 'username1'
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(true)

    item.addrs = 'addrs'
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(true)

    item.username = 'username'
    component.checkValueForEdit(item)
    expect(item.username).toStrictEqual(item.origin_username)
    expect(item.addrs).toStrictEqual(item.origin_addrs)
    expect(item.password).toStrictEqual(item.origin_password)
    expect(item.edit).toStrictEqual(false)

    item.addrs = 'addrs2'
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(true)

    item.addrs = 'addrs'
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(false)

    item.password = 'password'
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(false)

    item.password = ''
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(true)

    item.password = 'password'
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(false)

    item.password = 'addrs'
    item.addrs = 'password'
    item.username = 'password'
    component.checkValueForEdit(item)
    expect(item.edit).toStrictEqual(true)
  }))
})
