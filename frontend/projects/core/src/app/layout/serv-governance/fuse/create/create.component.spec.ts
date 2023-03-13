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
import { FuseCreateComponent } from './create.component'
import { GreyCreateComponent } from '../../grey/create/create.component'
import { VisitCreateComponent } from '../../visit/create/create.component'

class MockRenderer {
  removeAttribute(element: any, cssClass: string) {
    return cssClass + 'is removed from' + element
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

jest.mock('uuid', () => {
  return {
    v4: () => 123456789
  }
})

describe('FuseCreateComponent test as editPage is false', () => {
  let component: FuseCreateComponent
  let fixture: ComponentFixture<FuseCreateComponent>
  class MockElementRef extends ElementRef {
    constructor() {
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
            component: FuseCreateComponent
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
      declarations: [FuseCreateComponent],
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

    fixture = TestBed.createComponent(FuseCreateComponent)
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

  it('##getStrategyMessage with success return (without optional attr)', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const mockValue1: any = {
      code: 0,
      data: {
        strategy: {
          name: 'test',
          uuid: '1111',
          desc: 'testdescc',
          priority: 2,
          filters: [
            {
              name: 'ip',
              values: ['111.111.111.111', '111.111.111.112'],
              type: 'pattern',
              label: '111.111.111.111,111.111.111.112',
              title: 'IP'
            }
          ],
          config: {
            metric: 'service',
            fuse_condition: {
              status_codes: [300, 400],
              count: ''
            },
            fuse_time: {
              time: '',
              max_time: ''
            },
            recover_condition: {
              status_codes: [100, 200],
              count: ''
            },
            response: {
              status_code: 999,
              content_type: 'application/json',
              charset: 'UTF-8',
              header: [
                {
                  key: '1',
                  value: '2'
                }
              ],
              body: 'body~'
            }
          }
        }
      }
    }

    const expectCreateStrategyForm: any = {
      name: 'test',
      uuid: '1111',
      desc: 'testdescc',
      priority: 2,
      filters: [
        {
          name: 'ip',
          values: ['111.111.111.111', '111.111.111.112'],
          type: 'pattern',
          label: '111.111.111.111,111.111.111.112',
          title: 'IP'
        }
      ],
      config: {
        metric: 'service',
        fuse_condition: {
          status_codes: [300, 400],
          count: ''
        },
        fuse_time: {
          time: '',
          max_time: ''
        },
        recover_condition: {
          status_codes: [100, 200],
          count: ''
        },
        response: {
          status_code: 999,
          content_type: 'application/json',
          charset: 'UTF-8',
          header: [
            {
              key: '1',
              value: '2'
            }
          ],
          body: 'body~'
        }
      }
    }

    const expValidateForm: any = {
      name: 'test',
      desc: 'testdescc',
      priority: 2,
      configFuseCount: 3,
      configFuseTime: 2,
      configFuseMaxTime: 300,
      configRecoverCount: 3
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
    expect(component.filterNamesSet).toStrictEqual(new Set(['ip']))
    expect(component.filterShowList).toStrictEqual([
      {
        label: '111.111.111.111,111.111.111.112',
        name: 'ip',
        title: 'IP',
        type: 'pattern',
        values: ['111.111.111.111', '111.111.111.112']
      }
    ])
    expect(component.responseHeaderList).toStrictEqual([
      { key: '1', value: '2' }
    ])
    expect(component.responseForm.value).toStrictEqual({
      status_code: 999,
      content_type: 'application/json',
      charset: 'UTF-8',
      header: null,
      body: 'body~'
    })

    expect(component.responseForm.value).toStrictEqual({
      body: 'body~',
      charset: 'UTF-8',
      content_type: 'application/json',
      header: null,
      status_code: 999
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
          config: {
            metrics: [],
            response: {
              status_code: '',
              content_type: '',
              charset: '',
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
      status_code: 200,
      content_type: 'application/json',
      charset: 'UTF-8',
      header: null,
      body: ''
    })
    expect(component.responseForm.controls['content_type'].value).toStrictEqual(
      'application/json'
    )
    expect(component.responseForm.controls['charset'].value).toStrictEqual(
      'UTF-8'
    )
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
            metric: '{service}',
            fuse_condition: { status_codes: [100, 200], count: 3 },
            fuse_time: { time: 5, max_time: 5 },
            recover_condition: { status_codes: [300, 400], count: 3 },
            response: {
              status_code: '',
              content_type: '',
              charset: '',
              header: [],
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

    expect(component.createStrategyForm).toStrictEqual({
      name: 'testName',
      uuid: 'testUuid',
      desc: 'testDesc',
      priority: 5,
      filters: [{ name: 'testF', uuid: '123' }],
      config: {
        metric: '{service}',
        fuse_condition: { status_codes: [100, 200], count: 3 },
        fuse_time: { time: 5, max_time: 5 },
        recover_condition: { status_codes: [300, 400], count: 3 },
        response: {
          status_code: '',
          content_type: '',
          charset: '',
          header: [],
          body: ''
        }
      }
    })
    expect(component.filterNamesSet).toStrictEqual(new Set(['testF']))
    expect(component.filterShowList).toStrictEqual([
      { name: 'testF', uuid: '123' }
    ])
    expect(spyService).toHaveBeenCalled()
    expect(component.responseHeaderList).toStrictEqual([{ key: '', value: '' }])
    expect(component.responseForm.value).toStrictEqual({
      status_code: 200,
      content_type: 'application/json',
      charset: 'UTF-8',
      header: null,
      body: ''
    })
    expect(component.responseForm.controls['content_type'].value).toStrictEqual(
      'application/json'
    )
    expect(component.responseForm.controls['charset'].value).toStrictEqual(
      'UTF-8'
    )
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

  it('## saveStrategy with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyServicePut = jest.spyOn(httpCommonService, 'put').mockReturnValue(
      of({
        code: 0,
        data: {
          strategies: [1, 2, 3],
          is_publish: true,
          source: '123',
          version_name: 'test1',
          unpublish_msg: 'unpublish_test'
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
            is_publish: true,
            source: '123',
            version_name: 'test1',
            unpublish_msg: 'unpublish_test'
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
    component.filterShowList = [
      { title: 'test1', label: 'test1', name: 'test1', values: [''] },
      { title: 'test2', label: 'test2', name: 'test2', values: [''] },
      { title: 'test3', label: 'test3', name: 'test3', values: [] }
    ]

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [1]
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(2)
    component.validateForm.controls['configFuseTime'].setValue(2)
    component.validateForm.controls['configFuseMaxTime'].setValue(2)
    component.validateForm.controls['configRecoverCount'].setValue(2)
    component.responseForm.controls['status_code'].setValue(2)
    component.responseForm.controls['content_type'].setValue('application/json')
    component.responseForm.controls['charset'].setValue('UTF-8')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')
    component.responseHeaderList = [
      {
        key: '1',
        value: '2'
      },
      {
        key: '1',
        value: ''
      },
      {
        key: '',
        value: '2'
      }
    ]
    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [100]
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(2)
    component.validateForm.controls['configFuseTime'].setValue(2)
    component.validateForm.controls['configFuseMaxTime'].setValue(2)
    component.validateForm.controls['configRecoverCount'].setValue(2)
    component.responseForm.controls['status_code'].setValue(2)
    component.responseForm.controls['content_type'].setValue('application/json')
    component.responseForm.controls['charset'].setValue('UTF-8')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [100]
    component.validateForm.controls['name'].setValue('')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(2)
    component.validateForm.controls['configFuseTime'].setValue(2)
    component.validateForm.controls['configFuseMaxTime'].setValue(2)
    component.validateForm.controls['configRecoverCount'].setValue(2)
    component.responseForm.controls['status_code'].setValue(2)
    component.responseForm.controls['content_type'].setValue('application/json')
    component.responseForm.controls['charset'].setValue('UTF-8')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [100]
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(null)
    component.validateForm.controls['configFuseTime'].setValue(2)
    component.validateForm.controls['configFuseMaxTime'].setValue(2)
    component.validateForm.controls['configRecoverCount'].setValue(2)
    component.responseForm.controls['status_code'].setValue(2)
    component.responseForm.controls['content_type'].setValue('application/json')
    component.responseForm.controls['charset'].setValue('UTF-8')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [100]
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(2)
    component.validateForm.controls['configFuseTime'].setValue(null)
    component.validateForm.controls['configFuseMaxTime'].setValue(2)
    component.validateForm.controls['configRecoverCount'].setValue(2)
    component.responseForm.controls['status_code'].setValue(2)
    component.responseForm.controls['content_type'].setValue('application/json')
    component.responseForm.controls['charset'].setValue('UTF-8')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [100]
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(2)
    component.validateForm.controls['configFuseTime'].setValue(2)
    component.validateForm.controls['configFuseMaxTime'].setValue(null)
    component.validateForm.controls['configRecoverCount'].setValue(2)
    component.responseForm.controls['status_code'].setValue(2)
    component.responseForm.controls['content_type'].setValue('application/json')
    component.responseForm.controls['charset'].setValue('UTF-8')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [100]
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(2)
    component.validateForm.controls['configFuseTime'].setValue(2)
    component.validateForm.controls['configFuseMaxTime'].setValue(2)
    component.validateForm.controls['configRecoverCount'].setValue(null)
    component.responseForm.controls['status_code'].setValue(2)
    component.responseForm.controls['content_type'].setValue('application/json')
    component.responseForm.controls['charset'].setValue('UTF-8')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [100]
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(2)
    component.validateForm.controls['configFuseTime'].setValue(2)
    component.validateForm.controls['configFuseMaxTime'].setValue(2)
    component.validateForm.controls['configRecoverCount'].setValue(2)
    component.responseForm.controls['status_code'].setValue(null)
    component.responseForm.controls['content_type'].setValue('application/json')
    component.responseForm.controls['charset'].setValue('UTF-8')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [100]
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(2)
    component.validateForm.controls['configFuseTime'].setValue(2)
    component.validateForm.controls['configFuseMaxTime'].setValue(2)
    component.validateForm.controls['configRecoverCount'].setValue(2)
    component.responseForm.controls['status_code'].setValue(2002)
    component.responseForm.controls['content_type'].setValue('application/json')
    component.responseForm.controls['charset'].setValue('UTF-8')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [100]
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(2)
    component.validateForm.controls['configFuseTime'].setValue(2)
    component.validateForm.controls['configFuseMaxTime'].setValue(2)
    component.validateForm.controls['configRecoverCount'].setValue(2)
    component.responseForm.controls['status_code'].setValue(2)
    component.responseForm.controls['content_type'].setValue('')
    component.responseForm.controls['charset'].setValue('UTF-8')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [100]
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(2)
    component.validateForm.controls['configFuseTime'].setValue(2)
    component.validateForm.controls['configFuseMaxTime'].setValue(2)
    component.validateForm.controls['configRecoverCount'].setValue(2)
    component.responseForm.controls['status_code'].setValue(200)
    component.responseForm.controls['content_type'].setValue('test')
    component.responseForm.controls['charset'].setValue('')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')

    component.editPage = false
    component.saveStrategy()
    fixture.detectChanges()

    expect(spyServicePost).not.toHaveBeenCalled()
    expect(spyServicePut).not.toHaveBeenCalled()
    expect(spyBackToList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

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
    expect(component.createStrategyForm.config.response.header).toStrictEqual([
      { key: '1', value: '2' }
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

    component.createStrategyForm.config.fuse_condition.status_codes = [110, 112]
    component.createStrategyForm.config.recover_condition.status_codes = [100]
    component.validateForm.controls['name'].setValue('test')
    component.validateForm.controls['priority'].setValue(222)
    component.validateForm.controls['configFuseCount'].setValue(2)
    component.validateForm.controls['configFuseTime'].setValue(2)
    component.validateForm.controls['configFuseMaxTime'].setValue(2)
    component.validateForm.controls['configRecoverCount'].setValue(2)
    component.responseForm.controls['status_code'].setValue(200)
    component.responseForm.controls['content_type'].setValue('test')
    component.responseForm.controls['charset'].setValue('test')
    component.responseForm.controls['header'].setValue('')
    component.responseForm.controls['body'].setValue('')

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
    expect(spyMessage).toBeCalledWith('修改失败!')
  })

  it('##backToList test', fakeAsync(() => {
    // @ts-ignore
    const spyBack = jest.spyOn(component.location, 'back')
    expect(spyBack).not.toHaveBeenCalled()

    component.backToList()
    tick(100)
    expect(spyBack).toHaveBeenCalled()
  }))
})

describe('FuseCreateComponent test as editPage is false', () => {
  let component: FuseCreateComponent
  let fixture: ComponentFixture<FuseCreateComponent>
  class MockElementRef extends ElementRef {
    constructor() {
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
            component: FuseCreateComponent
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
      declarations: [FuseCreateComponent],
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

    fixture = TestBed.createComponent(FuseCreateComponent)
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
