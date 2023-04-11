/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-14 22:56:33
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-21 23:33:15
 * @FilePath: /apinto/src/app/layout/upstream/service-discovery-content/service-discovery-content.component.spec.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { RouterModule } from '@angular/router'
import { ElementRef, Renderer2, ChangeDetectorRef, Type } from '@angular/core'
import { APP_BASE_HREF } from '@angular/common'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { environment } from 'projects/core/src/environments/environment'
import { BidiModule } from '@angular/cdk/bidi'
import { Overlay } from '@angular/cdk/overlay'
import { of } from 'rxjs'
import { FormsModule } from '@angular/forms'
import { ApplicationMessageComponent } from '../message/message.component'
import { ApplicationAuthenticationComponent } from './authentication.component'
import { ApplicationPublishComponent } from '../publish/publish.component'
import { EoNgFeedbackModalService, EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { EoNgTableModule } from 'eo-ng-table'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgSelectModule } from 'eo-ng-select'
import { ComponentModule } from 'projects/core/src/app/component/component.module'
import { EoNgSwitchModule } from 'eo-ng-switch'

class MockRenderer {
  removeAttribute (element: any, cssClass: string) {
    return cssClass + 'is removed from' + element
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

jest.mock('uuid', () => {
  return {
    v4: () => 123456789
  }
})
describe('ApplicationAuthenticationComponent test', () => {
  let component: ApplicationAuthenticationComponent
  let fixture: ComponentFixture<ApplicationAuthenticationComponent>
  let renderer2: Renderer2
  class MockElementRef extends ElementRef {
    constructor () { super(null) }
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule, FormsModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule,
        NzDrawerModule, NzOutletModule, HttpClientModule,
        EoNgTableModule, EoNgDatePickerModule, EoNgInputModule,
        EoNgSelectModule, ComponentModule, EoNgSwitchModule,
        RouterModule.forRoot([
          {
            path: '',
            component: ApplicationMessageComponent
          },
          {
            path: 'auth',
            component: ApplicationAuthenticationComponent
          },
          {
            path: 'publish',
            component: ApplicationPublishComponent
          },
          {
            path: 'application',
            component: ApplicationPublishComponent
          }
        ]
        )
      ],
      declarations: [ApplicationAuthenticationComponent
      ],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: Renderer2, useClass: MockRenderer },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: EoNgFeedbackModalService, useClass: MockEnsureService },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(ApplicationAuthenticationComponent)
    renderer2 = fixture.componentRef.injector.get<Renderer2>(Renderer2 as Type<Renderer2>)
    renderer2.removeAttribute = jest.fn().mockReturnValue('remove')

    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('click authTable btns', () => {
    const spyGetAuthMessage = jest.spyOn(component, 'getAuthMessage')
    // @ts-ignore
    const spyModal = jest.spyOn(component.modalService, 'create')
    expect(spyGetAuthMessage).not.toHaveBeenCalled()
    expect(spyModal).not.toHaveBeenCalled()

    // 查看
    component.authenticationTableBody[6].btns[0].click({})
    fixture.detectChanges()
    expect(spyModal).not.toHaveBeenCalled()
    expect(spyGetAuthMessage).toHaveBeenCalledTimes(1)

    // 删除
    component.authenticationTableBody[6].btns[1].click({})
    fixture.detectChanges()
    expect(spyModal).toHaveBeenCalledTimes(1)
    expect(spyGetAuthMessage).toHaveBeenCalledTimes(1)
  })

  it('ngOnInit', () => {
    const spyGetAuthsData = jest.spyOn(component, 'getAuthsData')
    expect(spyGetAuthsData).not.toHaveBeenCalled()
    component.ngOnInit()
    expect(spyGetAuthsData).toHaveBeenCalled()
  })

  it('getDriversList with success return )', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { drivers: [{ name: 'apikey', render: 'render1' }, { name: 'aksk', render: 'render2' }] } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.driverList).toStrictEqual([])
    component.createAuthForm.driver = 'aksk'
    component.driverList = [{ label: 'test', value: 'test', render: 'test' }]
    component.getDriversList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.driverList).toStrictEqual([
      { label: 'ApiKey', value: 'apikey', render: 'render1' },
      { label: 'AkSk', value: 'aksk', render: 'render2' }
    ])
    expect(component.baseData).toStrictEqual('render2')
    expect(spyMessage).toHaveBeenCalledTimes(0)

    component.createAuthForm.driver = 'test'
    component.getDriversList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(3)
    expect(component.driverList).toStrictEqual([
      { label: 'ApiKey', value: 'apikey', render: 'render1' },
      { label: 'AkSk', value: 'aksk', render: 'render2' }
    ])
    expect(component.baseData).toStrictEqual('render2')
    expect(spyMessage).toHaveBeenCalledTimes(0)
  })

  it('getDriversList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.getDriversList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalled()
  })

  it('openDrawer', () => {
    const spyGetDriversList = jest.spyOn(component, 'getDriversList')
    expect(component.drawerAuthRef).toBeUndefined()
    expect(spyGetDriversList).not.toHaveBeenCalled()
    component.openDrawer()
    fixture.detectChanges()
    expect(component.drawerAuthRef).not.toBeUndefined()
    expect(spyGetDriversList).toHaveBeenCalled()
  })

  it('getAuthMessage with success return (expireTime = 0  & label != {})', () => {
    const val1:any = {
      code: 0,
      data: {
        auth: {
          uuid: '771c3b5e-547b-a80f-8321-6ecb5801bf5d',
          driver: 'jwt',
          expireTime: 0,
          operator: '',
          position: 'header',
          tokenName: '',
          updateTime: '2022-08-29 06:58:23',
          isTransparent: false,
          config: {
            userName: '1',
            password: '1',
            apikey: '1',
            ak: '1',
            sk: '1',
            iss: '1',
            algorithm: 'HS256',
            secret: '1',
            publicKey: '1',
            label: { 1: '2', 2: '1' }
          }
        }
      },
      msg: 'success'
    }

    const res1:any = {
      uuid: '771c3b5e-547b-a80f-8321-6ecb5801bf5d',
      driver: 'jwt',
      expireTime: 0,
      operator: '',
      position: 'header',
      tokenName: '',
      updateTime: '2022-08-29 06:58:23',
      isTransparent: false,
      config: {
        userName: '1',
        password: '1',
        apikey: '1',
        ak: '1',
        sk: '1',
        iss: '1',
        algorithm: 'HS256',
        secret: '1',
        publicKey: '1',
        label: [{ key: '1', value: '2' }, { key: '2', value: '1' }]
      }
    }
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of(val1))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()
    const spyOpenDrawer = jest.spyOn(component, 'openDrawer')
    expect(spyOpenDrawer).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.getAuthMessage({ test: 'test' })
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(3)
    expect(component.createAuthForm).toStrictEqual(res1)
    expect(spyOpenDrawer).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('getAuthMessage with success return ( expireTime != 0 & config = {})', () => {
    // label不为空对象且有时间戳,但config字段不存在
    const val2 = {
      code: 0,
      data: {
        auth: {
          uuid: '1dfdaaa0-6fb3-b91b-ca15-ab0bbcf6be12',
          driver: 'basic',
          expireTime: 1661961599,
          operator: '',
          position: 'header',
          tokenName: '1',
          updateTime: '2022-08-29 07:16:48',
          isTransparent: false,
          config: {
          }
        }
      },
      msg: 'success'
    }

    const res2 = {
      uuid: '1dfdaaa0-6fb3-b91b-ca15-ab0bbcf6be12',
      driver: 'basic',
      expireTime: 1661961599,
      expireTimeDate: new Date(1661961599000),
      operator: '',
      position: 'header',
      tokenName: '1',
      updateTime: '2022-08-29 07:16:48',
      isTransparent: false,
      config: {
        userName: '',
        password: '',
        apikey: '',
        ak: '',
        sk: '',
        iss: '',
        algorithm: '',
        secret: '',
        publicKey: '',
        label: [{ key: '', value: '' }]
      }
    }
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of(val2))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()
    const spyOpenDrawer = jest.spyOn(component, 'openDrawer')
    expect(spyOpenDrawer).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    component.getAuthMessage({ test: 'test' })
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(3)
    expect(component.createAuthForm).toStrictEqual(res2)
    expect(spyOpenDrawer).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('getAuthMessage with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)
    component.editAuth = false
    component.getAuthMessage({ test: 'test' })
    fixture.detectChanges()

    expect(component.editAuth).toStrictEqual(true)
    expect(spyService).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalled()
  })

  it('getAuthsData with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyServiceGet = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({
      code: 0,
      data: {
        auths: [
          {
            uuid: 'f7e1a642-6e33-4ab8-9c87-c1144c7180fb',
            info: 'apikey-fdsafjaklfasjkl',
            driver: 'apikey',
            expireTime: 0,
            operator: '',
            updateTime: '2022-08-29 14:31:44',
            ruleInfo: 'apikey-',
            isTransparent: true
          },
          {
            uuid: '274a7855-30bf-3775-661f-e878f7151e69',
            info: 'basic-1',
            driver: 'basic',
            expireTime: 0,
            operator: '',
            updateTime: '2022-08-30 11:07:26',
            ruleInfo: 'basic-',
            isTransparent: false
          },
          {
            uuid: '2fa8e791-37fc-1dfc-76fb-011e03822a89',
            info: 'basic-zhangzeyi',
            driver: 'basic',
            expireTime: 0,
            operator: '',
            updateTime: '2022-08-31 14:26:53',
            ruleInfo: 'basic-',
            isTransparent: true
          },
          {
            uuid: 'd0c7a6b6-221a-51fd-7325-8b97b5105241',
            info: 'jwt-1',
            driver: 'jwt',
            expireTime: 1662652799,
            operator: '',
            updateTime: '2022-08-31 14:27:01',
            ruleInfo: 'jwt-',
            isTransparent: false
          }
        ]
      },
      msg: 'success'
    }))
    expect(spyServiceGet).toHaveBeenCalledTimes(0)

    component.getAuthsData()
    expect(spyServiceGet).toHaveBeenCalledTimes(1)
    expect(component.authenticationList).toStrictEqual(
      [
        {
          uuid: 'f7e1a642-6e33-4ab8-9c87-c1144c7180fb',
          info: 'apikey-fdsafjaklfasjkl',
          driver: 'ApiKey',
          expireTime: '永不过期',
          operator: '',
          updateTime: '2022-08-29 14:31:44',
          ruleInfo: 'apikey-',
          isTransparent: '是'
        },
        {
          uuid: '274a7855-30bf-3775-661f-e878f7151e69',
          info: 'basic-1',
          driver: 'Basic',
          expireTime: '永不过期',
          operator: '',
          updateTime: '2022-08-30 11:07:26',
          ruleInfo: 'basic-',
          isTransparent: '否'
        },
        {
          uuid: '2fa8e791-37fc-1dfc-76fb-011e03822a89',
          info: 'basic-zhangzeyi',
          driver: 'Basic',
          expireTime: '永不过期',
          operator: '',
          updateTime: '2022-08-31 14:26:53',
          ruleInfo: 'basic-',
          isTransparent: '是'
        },
        {
          uuid: 'd0c7a6b6-221a-51fd-7325-8b97b5105241',
          info: 'jwt-1',
          driver: 'Jwt',
          expireTime: '2022-09-08 23:59:59',
          operator: '',
          updateTime: '2022-08-31 14:27:01',
          ruleInfo: 'jwt-',
          isTransparent: '否'
        }
      ]
    )
  })

  it('getAuthsData with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyServiceGet = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()
    component.getAuthsData()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyServiceGet).toHaveBeenCalled()
  })

  it('getDateInList', () => {
    let res = component.getDateInList(0)
    expect(res).toStrictEqual('1970-01-01 08:00:00')

    // @ts-ignore
    res = component.getDateInList('-1s')
    expect(res).toStrictEqual('日期数据格式有误')

    // @ts-ignore
    res = component.getDateInList('1111111111')
    expect(res).toStrictEqual('2005-03-18 09:58:31')
  })

  it('getAuthDriver', () => {
    let res = component.getAuthDriver('basic')
    expect(res).toStrictEqual('Basic')

    res = component.getAuthDriver('apikey')
    expect(res).toStrictEqual('ApiKey')

    res = component.getAuthDriver('aksk')
    expect(res).toStrictEqual('AkSk')

    res = component.getAuthDriver('jwt')
    expect(res).toStrictEqual('Jwt')

    res = component.getAuthDriver('test')
    expect(res).toStrictEqual('test')
  })

  it('getDataFromDynamicComponent when driver is basic', () => {
    let val1 = { driver: 'basic', config: { userName: '', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'basic', config: { userName: '1', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'basic', config: { userName: '1', password: '1' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(true)
  })

  it('getDataFromDynamicComponent when driver is apikey', () => {
    let val1:any = { driver: 'apikey', config: { userName: '', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'apikey', config: { apikey: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'apikey', config: { apikey: '1' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(true)
  })

  it('getDataFromDynamicComponent when driver is aksk', () => {
    let val1:any = { driver: 'aksk', config: { userName: '', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'aksk', config: { ak: '1', sk: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'aksk', config: { ak: '1', sk: '1' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(true)
  })

  it('getDataFromDynamicComponent when driver is jwt', () => {
    let val1:any = { driver: 'jwt', config: { userName: '', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'jwt', config: { iss: '1', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'jwt', config: { iss: '1', secret: '1' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'jwt', config: { algorithm: 'HS123', userName: '', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'jwt', config: { algorithm: 'HS123', iss: '1', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'jwt', config: { algorithm: 'HS123', iss: '1', secret: '1' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(true)

    val1 = { driver: 'jwt', config: { algorithm: '1ES2', userName: '', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'jwt', config: { algorithm: '123ES', iss: '1', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'jwt', config: { algorithm: 'ES123', iss: '1', publicKey: '1' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(true)

    val1 = { driver: 'jwt', config: { algorithm: '123RS', userName: '', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'jwt', config: { algorithm: 'RS123', iss: '1', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'jwt', config: { algorithm: 'RS123', iss: '1', publicKey: '1' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(true)
  })

  it('getDataFromDynamicComponent when driver is basic', () => {
    let val1 = { driver: 'basic', config: { userName: '', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'basic', config: { userName: '1', password: '' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(false)

    val1 = { driver: 'basic', config: { userName: '1', password: '1' } }
    component.canBeSave = false
    component.getDataFromDynamicComponent(val1)
    fixture.detectChanges()
    expect(component.canBeSave).toStrictEqual(true)
  })

  it('changeBasedata', fakeAsync(() => {
    component.driverList = [{ label: '静态节点', value: 'test1', render: 'render1' },
      { label: 'test2[driver2]', value: 'test2', render: 'render2' }
    ]
    component.baseData = 'render'
    component.createAuthForm.driver = 'test'
    expect(component.baseData).toStrictEqual('render')
    component.changeBasedata()
    fixture.detectChanges()
    expect(component.baseData).toStrictEqual('render')

    component.baseData = 'render'
    component.createAuthForm.driver = 'test1'
    expect(component.baseData).toStrictEqual('render')
    component.changeBasedata()
    fixture.detectChanges()
    expect(component.baseData).toStrictEqual('render1')
  }))

  it('saveAuth with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyServicePost = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: 0, data: { msg: 'success' } }))
    const ispost = httpCommonService.post('') !== null
    const spyServicePut = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 0, data: { msg: 'success' } }))
    const isput = httpCommonService.put('') !== null
    const spyGetAuthsData = jest.spyOn(component, 'getAuthsData')

    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    expect(spyMessageError).not.toHaveBeenCalled()

    expect(spyServicePost).toHaveBeenCalledTimes(1)
    expect(ispost).toStrictEqual(true)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(isput).toStrictEqual(true)

    component.createAuthForm = {
      position: 'header',
      uuid: '',
      tokenName: '1',
      isTransparent: false,
      expireTime: 0,
      expireTimeDate: null,
      driver: 'basic',
      config: {
        userName: '1',
        password: '2',
        apikey: '',
        ak: '',
        sk: '',
        iss: '',
        algorithm: '',
        label: [{ key: 'key1', value: 'val1' }]
      }
    }

    const res:any = {
      position: 'header',
      uuid: '',
      tokenName: '1',
      isTransparent: false,
      expireTime: 0,
      expireTimeDate: null,
      driver: 'basic',
      config: {
        userName: '1',
        password: '2',
        apikey: '',
        ak: '',
        sk: '',
        iss: '',
        publicKey: undefined,
        secret: undefined,
        algorithm: '',
        label: { key1: 'val1' }
      }
    }
    component.editAuth = false
    component.saveAuth()
    fixture.detectChanges()

    expect(component.createAuthForm).toEqual(res)
    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(spyGetAuthsData).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
    expect(spyMessageError).not.toHaveBeenCalled()

    component.editAuth = true
    component.saveAuth()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyServicePut).toHaveBeenCalledTimes(2)
    expect(spyGetAuthsData).toHaveBeenCalledTimes(2)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(2)
    expect(spyMessageError).not.toHaveBeenCalled()
  })

  it('saveAuth with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyServicePost = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const ispost = httpCommonService.post('') !== null
    const spyServicePut = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isput = httpCommonService.put('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyServicePost).toHaveBeenCalledTimes(1)
    expect(ispost).toStrictEqual(true)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(isput).toStrictEqual(true)

    component.createAuthForm.expireTimeDate = new Date()
    component.editAuth = false
    component.saveAuth()
    fixture.detectChanges()
    expect(component.createAuthForm.expireTime).not.toStrictEqual(0)
    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalledTimes(1)

    component.editAuth = true
    component.saveAuth()
    fixture.detectChanges()

    expect(spyServicePut).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalledTimes(2)
  })

  it('saveAuth with fail return', () => {
    component.createAuthForm = {
      position: 'header',
      uuid: '',
      tokenName: '1',
      isTransparent: false,
      expireTime: 0,
      expireTimeDate: null,
      driver: 'basic',
      config: {
        userName: '1',
        password: '2',
        apikey: '',
        ak: '',
        sk: '',
        iss: '',
        algorithm: '',
        label: [{ key: 'key1', value: 'val1' }]
      }
    }
    component.clearAuthForm()
    fixture.detectChanges()
    expect(component.createAuthForm).toStrictEqual(
      {
        position: 'header',
        uuid: '',
        tokenName: '',
        isTransparent: false,
        expireTime: 0,
        expireTimeDate: null,
        driver: 'basic',
        config: {
          userName: '',
          password: '',
          apikey: '',
          ak: '',
          sk: '',
          iss: '',
          algorithm: '',
          secret: '',
          publicKey: '',
          label: [{ key: '', value: '' }]
        }
      }
    )
  })
})
