/* eslint-disable dot-notation */
import {
  ComponentFixture,
  fakeAsync,
  TestBed,
  tick
} from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { ActivatedRoute, RouterModule } from '@angular/router'
import { ElementRef, Renderer2, ChangeDetectorRef } from '@angular/core'
import { APP_BASE_HREF } from '@angular/common'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { environment } from 'projects/core/src/environments/environment'
import { BidiModule } from '@angular/cdk/bidi'
import { Overlay } from '@angular/cdk/overlay'
import { of } from 'rxjs'
import {
  EoNgFeedbackMessageService,
  EoNgFeedbackModalService
} from 'eo-ng-feedback'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'

import { EoNgSelectModule } from 'eo-ng-select'
import { CacheCreateComponent } from '../../cache/create/create.component'
import { FuseCreateComponent } from '../../fuse/create/create.component'
import { GreyCreateComponent } from '../../grey/create/create.component'
import { VisitCreateComponent } from '../../visit/create/create.component'
import { TrafficCreateComponent } from './create.component'

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

describe('TrafficCreateComponent test as editPage is false', () => {
  let component: TrafficCreateComponent
  let fixture: ComponentFixture<TrafficCreateComponent>
  class MockElementRef extends ElementRef {
    constructor () {
      super(null)
    }
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
        FormsModule,
        ReactiveFormsModule,
        EoNgSelectModule,
        BidiModule,
        NoopAnimationsModule,
        NzNoAnimationModule,
        NzDrawerModule,
        NzOutletModule,
        HttpClientModule,
        RouterModule.forRoot([
          {
            path: 'serv-governance/traffic/create',
            component: TrafficCreateComponent
          },
          {
            path: 'serv-governance/fuse/create',
            component: FuseCreateComponent
          },
          {
            path: 'serv-governance/grey/create',
            component: GreyCreateComponent
          },
          {
            path: 'serv-governance/visit/create',
            component: VisitCreateComponent
          },
          {
            path: 'serv-governance/cache/create',
            component: CacheCreateComponent
          }
        ])
      ],
      declarations: [TrafficCreateComponent],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: Renderer2, useClass: MockRenderer },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: EoNgFeedbackModalService, useClass: MockEnsureService },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef },
        {
          provide: ActivatedRoute,
          useValue: {
            queryParams: of({ clusterName: 'clus2', strategy_uuid: 'uuid' })
          }
        }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(TrafficCreateComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('##should create', () => {
    expect(component).toBeTruthy()
  })

  it('##ngOnInit should call getStrategyMessage()', () => {
    const spyGetStrategyMessage = jest.spyOn(component, 'getStrategyMessage')
    const spyGetMetricsList = jest.spyOn(component, 'getMetricsList')
    expect(spyGetStrategyMessage).not.toHaveBeenCalled()
    expect(spyGetMetricsList).not.toHaveBeenCalled()
    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetStrategyMessage).not.toHaveBeenCalled()
    expect(spyGetMetricsList).toHaveBeenCalledTimes(1)
  })

  it('##ngOnChanges test', () => {
    component.strategyUuid = ''
    component.createStrategyForm.uuid = '1'
    component.ngOnChanges()
    expect(component.createStrategyForm.uuid).toStrictEqual('1')

    component.strategyUuid = '123'
    component.ngOnChanges()
    expect(component.createStrategyForm.uuid).toStrictEqual('123')
  })

  it('##getStrategyMessage with success return (without optional attr)', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const mockValue1: any = {
      code: 0,
      data: {
        strategy: {
          name: 'default_value',
          uuid: 'default_value',
          desc: 'default_value',
          priority: '',
          config: {
            response: {
              statusCode: '',
              contentType: 'default_value',
              charset: 'default_value',
              header: [
                {
                  key: 'default_value',
                  value: 'default_value'
                }
              ],
              body: 'default_value'
            }
          }
        }
      }
    }

    const expectCreateStrategyForm: any = {
      name: 'default_value',
      uuid: 'default_value',
      desc: 'default_value',
      priority: '',
      filters: [],
      config: {
        metrics: [],
        response: {
          statusCode: '',
          contentType: 'default_value',
          charset: 'default_value',
          header: [
            {
              key: 'default_value',
              value: 'default_value'
            }
          ],
          body: 'default_value'
        }
      }
    }

    const expValidateForm: any = {
      name: 'default_value',
      desc: 'default_value',
      priority: '',
      limitQuerySecond: 0,
      limitQueryMinute: 0,
      limitQueryHour: 0,
      limitTrafficSecond: 0,
      limitTrafficMinute: 0,
      limitTrafficHour: 0
    }

    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of(mockValue1))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.filterNamesSet = new Set()
    component.filterShowList = []
    component.getStrategyMessage()

    expect(component.createStrategyForm).toStrictEqual(expectCreateStrategyForm)
    expect(component.validateForm.value).toStrictEqual(expValidateForm)
    expect(component.filterNamesSet).toStrictEqual(new Set())
    expect(component.filterShowList).toStrictEqual([])
    expect(component.responseHeaderList).toStrictEqual([
      { key: 'default_value', value: 'default_value' }
    ])
    expect(component.responseForm.value).toStrictEqual({
      statusCode: 200,
      contentType: 'default_value',
      charset: 'default_value',
      header: null,
      body: 'default_value'
    })

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getStrategyMessage with success return (with null response)', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)

    const mockValue2: any = {
      code: 0,
      data: {
        strategy: {
          name: 'default_value',
          uuid: 'default_value',
          desc: 'default_value',
          priority: '',
          filters: [],
          config: {
            metrics: [],
            response: {
              statusCode: '',
              contentType: '',
              charset: '',
              header: [
              ],
              body: ''
            }
          }
        }
      }
    }
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of(mockValue2))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.filterNamesSet = new Set()
    component.filterShowList = []
    component.getStrategyMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(component.responseHeaderList).toStrictEqual([{ key: '', value: '' }])
    expect(component.responseForm.value).toStrictEqual({
      statusCode: 200,
      contentType: 'application/json',
      charset: 'UTF-8',
      header: null,
      body: ''
    })
    expect(component.responseForm.controls['contentType'].value).toStrictEqual('application/json')
    expect(component.responseForm.controls['charset'].value).toStrictEqual('UTF-8')
    expect(component.responseForm.controls['body'].value).toStrictEqual('')
  })

  it('##getStrategyMessage with success return (with optional attr)', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)

    const mockValue2: any = {
      code: 0,
      data: {
        strategy: {
          name: 'testName',
          uuid: 'testUuid',
          desc: 'testDesc',
          priority: 5,
          filters: [{ name: 'testF', uuid: '123' }],
          config: {
            metrics: ['test1', 'test2'],
            query: { second: 0, minute: 0, hour: 0 },
            traffic: { second: 0, minute: 0, hour: 0 }
          }
        }
      }
    }
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of(mockValue2))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.filterNamesSet = new Set()
    component.filterShowList = []
    component.getStrategyMessage()
    fixture.detectChanges()

    expect(component.createStrategyForm).toStrictEqual({
      name: 'testName',
      uuid: 'testUuid',
      desc: 'testDesc',
      priority: 5,
      filters: [{ name: 'testF', uuid: '123' }],
      config: {
        metrics: ['test1', 'test2'],
        query: { second: 0, minute: 0, hour: 0 },
        traffic: { second: 0, minute: 0, hour: 0 }
      }
    })
    expect(component.filterNamesSet).toStrictEqual(new Set(['testF']))
    expect(component.filterShowList).toStrictEqual([
      { name: 'testF', uuid: '123' }
    ])
    expect(spyService).toHaveBeenCalled()
    expect(component.responseHeaderList).toStrictEqual([{ key: '', value: '' }])
    expect(component.responseForm.value).toStrictEqual({
      statusCode: 200,
      contentType: 'application/json',
      charset: 'UTF-8',
      header: null,
      body: ''
    })
    expect(component.responseForm.controls['contentType'].value).toStrictEqual('application/json')
    expect(component.responseForm.controls['charset'].value).toStrictEqual('UTF-8')
    expect(component.responseForm.controls['body'].value).toStrictEqual('')
  })

  it('##getStrategyMessage with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getStrategyMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getStrategyMessage with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getStrategyMessage()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalledWith('获取数据失败!')
  })

  it('##disabledEdit test', fakeAsync(() => {
    component.nzDisabled = false
    component.disabledEdit(true)
    expect(component.nzDisabled).toStrictEqual(true)
  }))

  it('##getMetricsList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: {
          options: [
            { title: 'test1', name: 'testName1' },
            { title: 'test2', name: 'testName2' },
            { title: 'test3', name: 'testName3' },
            { title: 'test4', name: 'testName4' }
          ]
        }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.metricsList = [{ label: 't', value: 'n' }]
    component.getMetricsList()
    fixture.detectChanges()

    expect(component.metricsList).toStrictEqual([
      { label: 'test1', value: 'testName1' },
      { label: 'test2', value: 'testName2' },
      { label: 'test3', value: 'testName3' },
      { label: 'test4', value: 'testName4' }
    ])
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getMetricsList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getMetricsList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getMetricsList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: '' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getMetricsList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取数据失败!')
  })

  it('## saveStrategy with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyServicePut = jest.spyOn(httpCommonService, 'put').mockReturnValue(
      of({
        code: 0,
        data: {
          strategies: [1, 2, 3],
          isPublish: true,
          source: '123',
          versionName: 'test1',
          unpublishMsg: 'unpublish_test'
        }
      })
    )
    const spyServicePost = jest
      .spyOn(httpCommonService, 'post')
      .mockReturnValue(
        of({
          code: 0,
          data: {
            strategies: [1, 2, 3],
            isPublish: true,
            source: '123',
            versionName: 'test1',
            unpublishMsg: 'unpublish_test'
          }
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyBackToList = jest.spyOn(component, 'backToList')
    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    component.filterShowList = [{ title: 'test1', label: 'test1', name: 'test1', values: [''] }, { title: 'test2', label: 'test2', name: 'test2', values: [''] }, { title: 'test3', label: 'test3', name: 'test3', values: [] }]

    component.createStrategyForm.config.metrics = ['test']
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(2222)
    component.validateForm.controls['limitQuerySecond'].setValue(2)
    component.validateForm.controls['limitQueryMinute'].setValue(2)
    component.validateForm.controls['limitQueryHour'].setValue(2)
    component.validateForm.controls['limitTrafficSecond'].setValue(2)
    component.validateForm.controls['limitTrafficMinute'].setValue(2)
    component.validateForm.controls['limitTrafficHour'].setValue(2)
    component.responseForm.controls['statusCode'].setValue('')
    component.responseForm.controls['contentType'].setValue('')
    component.responseForm.controls['charset'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.validateForm.controls['priority'].setValue(222)
    component.responseForm.controls['statusCode'].setValue('201')
    component.responseForm.controls['contentType'].setValue('test')
    component.responseForm.controls['charset'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.validateForm.controls['priority'].setValue(222)
    component.responseForm.controls['statusCode'].setValue('201')
    component.responseForm.controls['contentType'].setValue('test')
    component.responseForm.controls['charset'].setValue('test')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component.createStrategyForm.filters).toStrictEqual([
      { name: 'test1', values: [''] },
      { name: 'test2', values: [''] },
      { name: 'test3', values: [] }
    ])

    component.editPage = true
    component.saveStrategy()
    fixture.detectChanges()
    expect(spyServicePost).toHaveBeenCalledTimes(1)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(spyBackToList).toHaveBeenCalledTimes(2)
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##saveStrategy with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyServicePut = jest
      .spyOn(httpCommonService, 'put')
      .mockReturnValue(of({ code: -1, data: {} }))
    const spyServicePost = jest
      .spyOn(httpCommonService, 'post')
      .mockReturnValue(of({ code: -1, data: {} }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.createStrategyForm.config.metrics = ['test']
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(22)
    component.validateForm.controls['limitQuerySecond'].setValue(2)
    component.validateForm.controls['limitQueryMinute'].setValue(2)
    component.validateForm.controls['limitQueryHour'].setValue(2)
    component.validateForm.controls['limitTrafficSecond'].setValue(2)
    component.validateForm.controls['limitTrafficMinute'].setValue(2)
    component.validateForm.controls['limitTrafficHour'].setValue(2)
    component.responseForm.controls['statusCode'].setValue(200)
    component.responseForm.controls['contentType'].setValue('test')
    component.responseForm.controls['charset'].setValue('try')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('创建失败!')

    component.editPage = true
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalledTimes(1)
    expect(spyServicePut).toHaveBeenCalledTimes(1)
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##backToList test', fakeAsync(() => {
    // @ts-ignore
    const spyBack = jest.spyOn(component.location, 'back')
    expect(spyBack).not.toHaveBeenCalled()

    component.backToList()
    tick(100)
    expect(spyBack).toHaveBeenCalled()
  }))

  it('##checkMetricOrder test', fakeAsync(() => {
    component.createStrategyForm.config.metrics = ['test3', 'test2', 'test1']
    component.metricsList = [
      { label: 'test1', value: 'test1' },
      { label: 'test2', value: 'test2' },
      { label: 'test3', value: 'test3' },
      { label: 'test4', value: 'test4' }
    ]
    component.checkMetricOrder()
    expect(component.createStrategyForm.config.metrics).toStrictEqual([
      'test1',
      'test2',
      'test3'
    ])
  }))

  it('##getContentTypeList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: {
          options: [
            { title: 'test1', name: 'testName1' },
            { title: 'test2', name: 'testName2' },
            { title: 'test3', name: 'testName3' },
            { title: 'test4', name: 'testName4' }
          ]
        }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.metricsList = [{ label: 't', value: 'n' }]
    component.getMetricsList()
    fixture.detectChanges()

    expect(component.metricsList).toStrictEqual([
      { label: 'test1', value: 'testName1' },
      { label: 'test2', value: 'testName2' },
      { label: 'test3', value: 'testName3' },
      { label: 'test4', value: 'testName4' }
    ])
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getMetricsList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getMetricsList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getMetricsList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: '' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getMetricsList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取数据失败!')
  })
})

describe('TrafficCreateComponent test as editPage is false', () => {
  let component: TrafficCreateComponent
  let fixture: ComponentFixture<TrafficCreateComponent>
  class MockElementRef extends ElementRef {
    constructor () {
      super(null)
    }
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
        FormsModule,
        ReactiveFormsModule,
        EoNgSelectModule,
        BidiModule,
        NoopAnimationsModule,
        NzNoAnimationModule,
        NzDrawerModule,
        NzOutletModule,
        HttpClientModule,
        RouterModule.forRoot([
          {
            path: 'serv-governance/traffic/create',
            component: TrafficCreateComponent
          },
          {
            path: 'serv-governance/fuse/create',
            component: FuseCreateComponent
          },
          {
            path: 'serv-governance/grey/create',
            component: GreyCreateComponent
          },
          {
            path: 'serv-governance/visit/create',
            component: VisitCreateComponent
          },
          {
            path: 'serv-governance/cache/create',
            component: CacheCreateComponent
          }
        ])
      ],
      declarations: [TrafficCreateComponent],
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

    fixture = TestBed.createComponent(TrafficCreateComponent)
    component = fixture.componentInstance
    component.editPage = true
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('ngOnInit should call getStrategyMessage()', () => {
    const spyGetStrategyMessage = jest.spyOn(component, 'getStrategyMessage')
    const spyGetMetricsList = jest.spyOn(component, 'getMetricsList')
    expect(spyGetStrategyMessage).not.toHaveBeenCalled()
    expect(spyGetMetricsList).not.toHaveBeenCalled()
    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetStrategyMessage).toHaveBeenCalled()
    expect(spyGetMetricsList).toHaveBeenCalledTimes(1)
  })
})
