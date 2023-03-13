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
import { FuseCreateComponent } from '../../fuse/create/create.component'
import { VisitCreateComponent } from '../../visit/create/create.component'
import { GreyCreateComponent } from './create.component'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgRadioModule } from 'eo-ng-radio'
import { NzSliderModule } from 'ng-zorro-antd/slider'

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

describe('GreyCreateComponent test as editPage is false', () => {
  let component: GreyCreateComponent
  let fixture: ComponentFixture<GreyCreateComponent>
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
        EoNgSwitchModule,
        EoNgRadioModule,
        NzSliderModule,
        NzOutletModule,
        HttpClientModule,
        RouterModule.forRoot([
          {
            path: 'serv-governance/traffic/create',
            component: GreyCreateComponent
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
            component: GreyCreateComponent
          }
        ])
      ],
      declarations: [GreyCreateComponent],
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
            queryParams: of({ cluster_name: 'clus2', strategy_uuid: 'uuid' })
          }
        }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(GreyCreateComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('##should create', () => {
    expect(component).toBeTruthy()
  })

  it('##ngOnInit should call getStrategyMessage()', () => {
    const spyGetStrategyMessage = jest.spyOn(component, 'getStrategyMessage')
    expect(spyGetStrategyMessage).not.toHaveBeenCalled()
    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetStrategyMessage).not.toHaveBeenCalled()
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

  it('##nodesTableBody test', fakeAsync(() => {
    component.showNodesValid = true
    expect(component.nodesTableBody[0].check({ test: 'test' })).toStrictEqual(true)
    component.nodesTableBody[0].change()
    expect(component.showNodesValid).toStrictEqual(false)
    component.nzDisabled = true
    expect(component.nodesTableBody[0].disabledFn()).toStrictEqual(true)
    component.nzDisabled = false
    expect(component.nodesTableBody[0].disabledFn()).toStrictEqual(false)
    expect(component.nodesTableBody[1].btns[0].disabledFn()).toStrictEqual(false)
    expect(component.nodesTableBody[2].btns[0].disabledFn()).toStrictEqual(false)
    expect(component.nodesTableBody[2].btns[1].disabledFn()).toStrictEqual(false)

    component.nzDisabled = true
    expect(component.nodesTableBody[0].disabledFn()).toStrictEqual(true)
    expect(component.nodesTableBody[1].btns[0].disabledFn()).toStrictEqual(true)
    expect(component.nodesTableBody[2].btns[0].disabledFn()).toStrictEqual(true)
    expect(component.nodesTableBody[2].btns[1].disabledFn()).toStrictEqual(true)
  }))

  it('##drawerTipShowFn test', fakeAsync(() => {
    component.filterNamesSet = new Set()
    expect(component.drawerTipShowFn()).toStrictEqual(true)

    component.filterNamesSet.add('application')
    expect(component.drawerTipShowFn()).toStrictEqual(true)

    component.filterNamesSet.add('api')
    expect(component.drawerTipShowFn()).toStrictEqual(false)

    component.filterNamesSet.add('service')
    expect(component.drawerTipShowFn()).toStrictEqual(false)

    component.filterNamesSet.add('api')
    expect(component.drawerTipShowFn()).toStrictEqual(false)
  }))

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
            nodes: [null]
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
        match: [],
        nodes: [null]
      }
    }

    const expValidateForm: any = {
      name: 'default_value',
      desc: 'default_value',
      priority: null,
      keepSession: false,
      distribution: 'percent',
      percent1: 1,
      percent2: 99

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
    expect(component.nodesList).toStrictEqual([{ node: '' }])

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getStrategyMessage with success return (with optional attr - percent)', () => {
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
            keep_session: false,
            nodes: ['node1', 'node2'],
            distribution: 'percent',
            match: [],
            percent: 12
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
        keep_session: false,
        nodes: ['node1', 'node2'],
        distribution: 'percent',
        match: [],
        percent: 12
      }
    })
    expect(component.filterNamesSet).toStrictEqual(new Set(['testF']))
    expect(component.filterShowList).toStrictEqual([
      { name: 'testF', uuid: '123' }
    ])
    expect(component.nodesList).toStrictEqual([
      { node: 'node1' }, { node: 'node2' }
    ])
    expect(spyService).toHaveBeenCalled()
    expect(component.validateForm.controls['percent1'].value).toStrictEqual(0.12)
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

  it('##checkNodesList', fakeAsync(() => {
    component.showNodesValid = false
    component.nodesList = []
    component.checkNodesList()
    expect(component.showNodesValid).toStrictEqual(true)

    component.nodesList = [{ node: '2' }, { node: '' }]
    component.checkNodesList()
    expect(component.showNodesValid).toStrictEqual(false)

    component.nodesList = [{ node: '' }, { node: '' }]
    component.checkNodesList()
    expect(component.showNodesValid).toStrictEqual(true)

    component.nodesList = [{ node: '' }, { node: '2' }]
    component.checkNodesList()
    expect(component.showNodesValid).toStrictEqual(false)
  }))

  it('##saveStrategy with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyServicePut = jest.spyOn(httpCommonService, 'put').mockReturnValue(
      of({
        code: 0,
        data: {
        }
      })
    )
    const spyServicePost = jest
      .spyOn(httpCommonService, 'post')
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
    const spyBackToList = jest.spyOn(component, 'backToList')
    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    component.filterShowList = [{ title: 'test1', label: 'test1', name: 'test1', values: [''] }, { title: 'test2', label: 'test2', name: 'test2', values: [''] }, { title: 'test3', label: 'test3', name: 'test3', values: [] }]

    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(2222)
    component.validateForm.controls['keepSession'].setValue(false)
    component.validateForm.controls['distribution'].setValue('percent')
    component.validateForm.controls['percent1'].setValue(1)
    component.validateForm.controls['percent2'].setValue(99)

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()

    component.validateForm.controls['name'].setValue('123t')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['keepSession'].setValue(false)
    component.validateForm.controls['distribution'].setValue('percent')
    component.validateForm.controls['percent1'].setValue(1)
    component.validateForm.controls['percent2'].setValue(99)

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()

    component.validateForm.controls['name'].setValue('')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['keepSession'].setValue(false)
    component.validateForm.controls['distribution'].setValue('percent')
    component.validateForm.controls['percent1'].setValue(1)
    component.validateForm.controls['percent2'].setValue(99)
    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()

    component.validateForm.controls['name'].setValue('test1_')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['keepSession'].setValue(null)
    component.validateForm.controls['distribution'].setValue('percent')
    component.validateForm.controls['percent1'].setValue(1)
    component.validateForm.controls['percent2'].setValue(99)
    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()

    component.validateForm.controls['name'].setValue('test1_1')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['keepSession'].setValue(false)
    component.validateForm.controls['distribution'].setValue(null)
    component.validateForm.controls['percent1'].setValue(1)
    component.validateForm.controls['percent2'].setValue(99)
    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()

    component.validateForm.controls['name'].setValue('test__')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['keepSession'].setValue(false)
    component.validateForm.controls['distribution'].setValue('percent')
    component.validateForm.controls['percent1'].setValue(1)
    component.validateForm.controls['percent2'].setValue(99)
    component.nodesList = [{ node: 'test' }]
    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).toBeCalledWith('创建成功!')

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
    expect(spyMessageSuccess).toHaveBeenCalledTimes(2)
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

    component.validateForm.controls['name'].setValue('test__')
    component.validateForm.controls['priority'].setValue(null)
    component.validateForm.controls['keepSession'].setValue(false)
    component.validateForm.controls['distribution'].setValue('match')
    component.validateForm.controls['percent1'].setValue(1)
    component.validateForm.controls['percent2'].setValue(99)
    component.nodesList = [{ node: 'test' }]

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

  it('##backToList test', fakeAsync(() => {
    component.validateForm.controls['percent1'].setValue(2)
    component.changePercent(103, 'percent1')
    expect(component.validateForm.controls['percent1'].value).toStrictEqual(100)
    expect(component.validateForm.controls['percent2'].value).toStrictEqual(0)

    component.changePercent(-2, 'percent2')
    expect(component.validateForm.controls['percent1'].value).toStrictEqual(100)
    expect(component.validateForm.controls['percent2'].value).toStrictEqual(0)

    component.changePercent(23, 'percent1')
    expect(component.validateForm.controls['percent1'].value).toStrictEqual(23)
    expect(component.validateForm.controls['percent2'].value).toStrictEqual(77)

    component.changePercent(28, 'percent2')
    expect(component.validateForm.controls['percent1'].value).toStrictEqual(72)
    expect(component.validateForm.controls['percent2'].value).toStrictEqual(28)
  }))
})

describe('GreyCreateComponent test as editPage is false', () => {
  let component: GreyCreateComponent
  let fixture: ComponentFixture<GreyCreateComponent>
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
        EoNgSwitchModule,
        EoNgRadioModule,
        NzSliderModule,
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
            component: GreyCreateComponent
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
            component: GreyCreateComponent
          }
        ])
      ],
      declarations: [GreyCreateComponent],
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

    fixture = TestBed.createComponent(GreyCreateComponent)
    component = fixture.componentInstance
    component.editPage = true
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('ngOnInit should call getStrategyMessage()', () => {
    const spyGetStrategyMessage = jest.spyOn(component, 'getStrategyMessage')
    expect(spyGetStrategyMessage).not.toHaveBeenCalled()
    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetStrategyMessage).toHaveBeenCalled()
  })
})
