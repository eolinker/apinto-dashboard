/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-14 22:56:33
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-21 23:33:15
 * @FilePath: /apinto/src/app/layout/upstream/service-discovery-content/service-discovery-content.component.spec.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { ComponentFixture, fakeAsync, flush, TestBed } from '@angular/core/testing'
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
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { FormsModule } from '@angular/forms'
import { ApplicationMessageComponent } from '../message/message.component'
import { ApplicationAuthenticationComponent } from '../authentication/authentication.component'
import { ApplicationPublishComponent } from '../publish/publish.component'
import { ServiceDiscoveryMessageComponent } from '../../upstream/service-discovery/message/message.component'
import { ServiceDiscoveryPublishComponent } from '../../upstream/service-discovery/publish/publish.component'
import { ApplicationCreateComponent } from './create.component'

import { EoNgTableModule } from 'eo-ng-table'

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
describe('ApplicationCreateComponent test as editPage is false', () => {
  let component: ApplicationCreateComponent
  let fixture: ComponentFixture<ApplicationCreateComponent>
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
        EoNgTableModule,
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
      declarations: [ApplicationCreateComponent
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

    fixture = TestBed.createComponent(ApplicationCreateComponent)
    renderer2 = fixture.componentRef.injector.get<Renderer2>(Renderer2 as Type<Renderer2>)
    renderer2.removeAttribute = jest.fn().mockReturnValue('remove')

    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('click customtable btns', fakeAsync(() => {
    const item = { key: '22', value: '33' }
    component.customAttrList = [{ key: '1', value: '2' }, item, { key: '3', value: '3' }]
    component.customAttrTableBody[2].btns[0].click(item)
    fixture.detectChanges()
    flush()
    expect(component.customAttrList).toStrictEqual([{ eoKey: 123456789, key: '1', value: '2' }, { eoKey: 123456789, key: '22', value: '33' }, { eoKey: 123456789, key: '', value: '' }, { eoKey: 123456789, key: '3', value: '3' }])

    component.customAttrTableBody[3].btns[0].click(item)
    fixture.detectChanges()
    flush()
    expect(component.customAttrList).toStrictEqual([{ eoKey: 123456789, key: '1', value: '2' }, { eoKey: 123456789, key: '', value: '' }, { eoKey: 123456789, key: '3', value: '3' }])

    component.customAttrList = [{ key: '1', value: '2' }, item, { key: '3', value: '3' }]
    component.customAttrTableBody[3].btns[1].click(item)
    fixture.detectChanges()
    flush()
    expect(component.customAttrList).toStrictEqual([{ eoKey: 123456789, key: '1', value: '2' }, { eoKey: 123456789, key: '22', value: '33' }, { eoKey: 123456789, key: '', value: '' }, { eoKey: 123456789, key: '3', value: '3' }])

    component.customAttrList = [{ key: '1', value: '2' }, { key: '3', value: '3' }]
    component.customAttrTableBody[2].btns[0].click(item)
    fixture.detectChanges()
    flush()
    expect(component.customAttrList).toStrictEqual([{ eoKey: 123456789, key: '1', value: '2' }, { eoKey: 123456789, key: '3', value: '3' }])

    component.customAttrTableBody[3].btns[0].click(item)
    fixture.detectChanges()
    flush()
    expect(component.customAttrList).toStrictEqual([{ eoKey: 123456789, key: '1', value: '2' }, { eoKey: 123456789, key: '3', value: '3' }])

    component.customAttrTableBody[3].btns[1].click(item)
    fixture.detectChanges()
    flush()
    expect(component.customAttrList).toStrictEqual([{ eoKey: 123456789, key: '1', value: '2' }, { eoKey: 123456789, key: '3', value: '3' }])
  }))

  it('click extraHeader table btns', fakeAsync(() => {
    const item = { key: '22', value: '33' }
    component.extraHeaderList = [{ key: '1', value: '2' }, item, { key: '3', value: '3' }]
    component.extraHeaderTableBody[2].btns[0].click(item)
    fixture.detectChanges()
    expect(component.extraHeaderList).toStrictEqual([{ eoKey: 123456789, key: '1', value: '2' }, { eoKey: 123456789, key: '22', value: '33' }, { eoKey: 123456789, key: '', value: '' }, { eoKey: 123456789, key: '3', value: '3' }])

    component.extraHeaderTableBody[3].btns[0].click(item)
    fixture.detectChanges()
    expect(component.extraHeaderList).toStrictEqual([{ eoKey: 123456789, key: '1', value: '2' }, { eoKey: 123456789, key: '', value: '' }, { eoKey: 123456789, key: '3', value: '3' }])

    component.extraHeaderList = [{ key: '1', value: '2' }, item, { key: '3', value: '3' }]
    component.extraHeaderTableBody[3].btns[1].click(item)
    fixture.detectChanges()
    expect(component.extraHeaderList).toStrictEqual([{ eoKey: 123456789, key: '1', value: '2' }, { eoKey: 123456789, key: '22', value: '33' }, { eoKey: 123456789, key: '', value: '' }, { eoKey: 123456789, key: '3', value: '3' }])

    component.extraHeaderList = [{ key: '1', value: '2' }, { key: '3', value: '3' }]
    component.extraHeaderTableBody[2].btns[0].click(item)
    fixture.detectChanges()
    expect(component.extraHeaderList).toStrictEqual([{ key: '1', value: '2' }, { key: '3', value: '3' }])

    component.extraHeaderTableBody[3].btns[0].click(item)
    fixture.detectChanges()
    expect(component.extraHeaderList).toStrictEqual([{ key: '1', value: '2' }, { key: '3', value: '3' }])

    component.extraHeaderTableBody[3].btns[1].click(item)
    fixture.detectChanges()
    expect(component.extraHeaderList).toStrictEqual([{ key: '1', value: '2' }, { key: '3', value: '3' }])
  }))

  it('editPage is false', () => {
    const spyGetApplicationMessage = jest.spyOn(component, 'getApplicationMessage')
    const spyGetApplicationId = jest.spyOn(component, 'getApplicationId')
    expect(spyGetApplicationMessage).not.toHaveBeenCalled()
    expect(spyGetApplicationId).not.toHaveBeenCalled()
    component.ngOnInit()
    expect(spyGetApplicationMessage).not.toHaveBeenCalled()
    expect(spyGetApplicationId).toHaveBeenCalled()
  })

  it('getApplicationMessage with success return (without custom_arr and extra_header)', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { application: { name: 'testName' } } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.createApplicationForm.name).toStrictEqual('')
    component.getApplicationMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.createApplicationForm.name).toStrictEqual('testName')

    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('getApplicationMessage with success return (without custom_arr)', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { application: { name: 'testName', extra_header: 'TEST' } } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.createApplicationForm.name).toStrictEqual('')
    component.getApplicationMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.createApplicationForm.name).toStrictEqual('testName')

    expect(spyMessage).toHaveBeenCalledTimes(0)
  })

  it('getApplicationMessage with success return (without extra_header)', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { application: { name: 'testName', extra_header: 'TEST' } } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.createApplicationForm.name).toStrictEqual('')
    component.getApplicationMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.createApplicationForm.name).toStrictEqual('testName')

    expect(spyMessage).toHaveReturnedTimes(0)
  })

  it('getApplicationMessage with success return )', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { application: { name: 'testName', custom_attr: 'test', extra_header: 'test' } } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.createApplicationForm.name).toStrictEqual('')
    component.getApplicationMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.createApplicationForm.name).toStrictEqual('testName')

    expect(spyMessage).toHaveBeenCalledTimes(0)
  })

  it('getApplicationMessage with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.createApplicationForm.name).toStrictEqual('')
    component.getApplicationMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.createApplicationForm.name).toStrictEqual('')
    expect(spyMessage).toHaveBeenCalled()
  })

  it('getApplicationId with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 0, data: { id: '123' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.createApplicationForm.id).toStrictEqual('')
    component.getApplicationId()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.createApplicationForm.id).toStrictEqual('123')
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('getApplicationMessage with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(of({ code: 1, data: { msg: 'fail' } }))
    const isget = httpCommonService.get('') !== null
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyMessage).not.toHaveBeenCalled()

    expect(spyService).toHaveBeenCalledTimes(1)
    expect(isget).toStrictEqual(true)

    expect(component.createApplicationForm.id).toStrictEqual('')
    component.getApplicationMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalledTimes(2)
    expect(component.createApplicationForm.id).toStrictEqual('')
    expect(spyMessage).toHaveBeenCalled()
  })

  it('saveApplication with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyServicePost = jest.spyOn(httpCommonService, 'post').mockReturnValue(of({ code: 0, data: { msg: 'success' } }))
    const ispost = httpCommonService.post('') !== null
    const spyServicePut = jest.spyOn(httpCommonService, 'put').mockReturnValue(of({ code: 0, data: { msg: 'success' } }))
    const isput = httpCommonService.put('') !== null

    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    const spybackToList = jest.spyOn(component, 'backToList')
    expect(spybackToList).not.toHaveBeenCalled()

    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    expect(spyMessageError).not.toHaveBeenCalled()

    expect(spyServicePost).toHaveBeenCalledTimes(1)
    expect(ispost).toStrictEqual(true)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(isput).toStrictEqual(true)

    component.customAttrList = [
      { key: 'key1', value: 'val1' },
      { key: 'key2', value: '' },
      { key: '', value: 'val3' },
      { key: 'key4', value: 'val4' }
    ]
    component.editPage = false
    component.saveApplication()
    fixture.detectChanges()

    expect(component.createApplicationForm.custom_attr_list).toStrictEqual([
      { eoKey: 123456789, key: 'key1', value: 'val1' },
      { eoKey: 123456789, key: 'key4', value: 'val4' }
    ])
    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(spybackToList).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalledTimes(1)
    expect(spyMessageError).not.toHaveBeenCalled()

    component.extraHeaderList = [
      { key: '', value: 'val1' },
      { key: 'key2', value: 'val2' },
      { key: 'key3', value: 'val3' },
      { key: 'key4', value: '' }
    ]
    component.editPage = true
    component.saveApplication()
    fixture.detectChanges()

    expect(component.createApplicationForm.extra_param_list).toStrictEqual([
      { eoKey: 123456789, key: 'key2', value: 'val2' },
      { eoKey: 123456789, key: 'key3', value: 'val3' }
    ])
    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyServicePut).toHaveBeenCalledTimes(2)
    expect(spybackToList).toHaveBeenCalledTimes(2)
    expect(spyMessageSuccess).toHaveBeenCalledTimes(2)
    expect(spyMessageError).not.toHaveBeenCalled()
  })

  it('saveApplication with fail return', () => {
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

    component.editPage = false
    component.saveApplication()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalledTimes(1)

    component.editPage = true
    component.saveApplication()
    fixture.detectChanges()

    expect(spyServicePut).toHaveBeenCalledTimes(2)
    expect(spyMessage).toHaveBeenCalledTimes(2)
  })

  it('backToList', fakeAsync(() => {
    // @ts-ignore
    const spyRouter = jest.spyOn(component.router, 'navigate')
    expect(spyRouter).not.toHaveBeenCalled()

    component.backToList()
    fixture.detectChanges()

    expect(spyRouter).toHaveBeenCalled()
  }))
})

describe('ApplicationCreateComponent test as editPage is false', () => {
  let component: ApplicationCreateComponent
  let fixture: ComponentFixture<ApplicationCreateComponent>
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
        EoNgTableModule,
        RouterModule.forRoot([
          {
            path: '',
            component: ServiceDiscoveryPublishComponent
          },
          {
            path: 'message',
            component: ServiceDiscoveryMessageComponent
          }
        ]
        )
      ],
      declarations: [ApplicationCreateComponent
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

    fixture = TestBed.createComponent(ApplicationCreateComponent)
    renderer2 = fixture.componentRef.injector.get<Renderer2>(Renderer2 as Type<Renderer2>)
    renderer2.removeAttribute = jest.fn().mockReturnValue('remove')

    component = fixture.componentInstance
    component.editPage = true
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('change editPage from false to true', () => {
    const spyGetApplicationMessage = jest.spyOn(component, 'getApplicationMessage')
    const spyGetApplicationId = jest.spyOn(component, 'getApplicationId')
    expect(spyGetApplicationMessage).not.toHaveBeenCalled()
    expect(spyGetApplicationId).not.toHaveBeenCalled()
    component.ngOnInit()
    expect(spyGetApplicationMessage).toHaveBeenCalled()
    expect(spyGetApplicationId).not.toHaveBeenCalled()
  })
})
