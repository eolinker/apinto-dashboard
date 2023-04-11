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
import { FormsModule } from '@angular/forms'

import { EoNgSelectModule } from 'eo-ng-select'
import { TransferChange } from 'ng-zorro-antd/transfer'
import { CacheCreateComponent } from '../../cache/create/create.component'
import { FuseCreateComponent } from '../../fuse/create/create.component'
import { GreyCreateComponent } from '../../grey/create/create.component'
import { TrafficCreateComponent } from '../../traffic/create/create.component'
import { VisitCreateComponent } from '../../visit/create/create.component'
import { FilterFormComponent } from './form.component'

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

describe('FilterFormComponent test as editPage is false', () => {
  let component: FilterFormComponent
  let fixture: ComponentFixture<FilterFormComponent>
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
      declarations: [FilterFormComponent],
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
            queryParams: of({ clusterName: 'clus2' })
          }
        }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(FilterFormComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('##should create', () => {
    expect(component).toBeTruthy()
  })

  it('##ngOnInit should call getFilterNamesList()', () => {
    const spyGetFilterNamesList = jest.spyOn(component, 'getFilterNamesList')
    expect(spyGetFilterNamesList).not.toHaveBeenCalled()
    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetFilterNamesList).toHaveBeenCalledTimes(1)
  })

  it('##getFilterNamesList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: {
          options: [
            { name: 'api', title: 'API', type: 'remote' },
            {
              name: 'methods',
              title: 'API请求方式',
              type: 'static',
              options: [
                'ALL',
                'POST',
                'PUT',
                'GET',
                'DELETE',
                'OPTION',
                'PATCH',
                'HEAD'
              ]
            },
            { name: 'path', title: 'API路径', type: 'pattern', pattern: '' }
          ]
        }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyGetApiGroupList = jest.spyOn(component, 'getApiGroupList')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyGetApiGroupList).not.toHaveBeenCalled()

    component.filterForm.name = ''
    component.filterNamesSet = new Set()
    component.filterTypeMap = new Map()
    component.filterNamesList = []
    component.filterType = ''

    component.getFilterNamesList()
    fixture.detectChanges()

    expect(spyGetApiGroupList).toHaveBeenCalledTimes(1)
    expect(component.filterForm).toStrictEqual({
      name: 'api',
      title: 'API',
      label: 'API',
      value: 'api',
      type: 'remote',
      patternIsPass: true,
      allChecked: false,
      total: undefined,
      values: []
    })
    expect(component.filterNamesList).toStrictEqual([
      {
        allChecked: false,
        name: 'api',
        title: 'API',
        label: 'API',
        value: 'api',
        type: 'remote',
        patternIsPass: true,
        total: undefined,
        values: []
      },
      {
        allChecked: false,
        name: 'methods',
        title: 'API请求方式',
        label: 'API请求方式',
        value: 'methods',
        patternIsPass: true,
        type: 'static',
        options: [
          'ALL',
          'POST',
          'PUT',
          'GET',
          'DELETE',
          'OPTION',
          'PATCH',
          'HEAD'
        ],
        total: 7,
        values: []
      },
      {
        allChecked: false,
        name: 'path',
        title: 'API路径',
        label: 'API路径',
        value: 'path',
        type: 'pattern',
        patternIsPass: true,
        total: 0,
        pattern: '',
        values: []
      }
    ])
    expect(component.filterTypeMap.get('api')).toStrictEqual({
      allChecked: false,
      name: 'api',
      title: 'API',
      label: 'API',
      value: 'api',
      patternIsPass: true,
      type: 'remote',
      total: undefined,
      values: []
    })
    expect(component.filterTypeMap.get('methods')).toStrictEqual({
      allChecked: false,
      name: 'methods',
      title: 'API请求方式',
      label: 'API请求方式',
      patternIsPass: true,
      value: 'methods',
      type: 'static',
      options: [
        'ALL',
        'POST',
        'PUT',
        'GET',
        'DELETE',
        'OPTION',
        'PATCH',
        'HEAD'
      ],
      total: 7,
      values: []
    })
    expect(component.filterTypeMap.get('path')).toStrictEqual({
      allChecked: false,
      name: 'path',
      title: 'API路径',
      label: 'API路径',
      value: 'path',
      patternIsPass: true,
      type: 'pattern',
      total: 0,
      pattern: '',
      values: []
    })
    expect(component.filterType).toStrictEqual('remote')
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.filterForm.name = 'methods'

    component.getFilterNamesList()

    expect(component.filterForm).toStrictEqual({
      allChecked: false,
      name: 'methods',
      title: 'API请求方式',
      label: 'API请求方式',
      value: 'methods',
      patternIsPass: true,
      type: 'static',
      showAll: true,
      options: [
        'ALL',
        'POST',
        'PUT',
        'GET',
        'DELETE',
        'OPTION',
        'PATCH',
        'HEAD'
      ],
      total: 7,
      values: []
    })
  })

  it('##getFilterNamesList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getFilterNamesList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getFilterNamesList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getFilterNamesList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取数据失败!')
  })

  it('##getRemoteList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: {
          target: 'test1',
          titles: [
            { title: 'uuid', field: 'uuid' },
            { title: 'name', field: 'name' },
            { title: 'scheme', field: 'scheme' }
          ],
          test1: [
            { uuid: '111', name: 'name1', scheme: 'test1', desc: 'desc1' },
            { uuid: '222', name: 'name2', scheme: 'test2', desc: 'desc2' },
            { uuid: '333', name: 'name3', scheme: 'test3', desc: 'desc3' }
          ],
          total: 3
        }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getRemoteList('api')
    fixture.detectChanges()
    expect(component.remoteSelectList).toStrictEqual([])
    expect(component.remoteSelectNameList).toStrictEqual([])
    expect(component.remoteList).toStrictEqual([
      {
        uuid: '111',
        key: '111',
        checked: false,
        title: 'name1',
        direction: 'left',
        name: 'name1',
        scheme: 'test1',
        desc: 'desc1'
      },
      {
        uuid: '222',
        key: '222',
        checked: false,
        title: 'name2',
        direction: 'left',
        name: 'name2',
        scheme: 'test2',
        desc: 'desc2'
      },
      {
        uuid: '333',
        key: '333',
        checked: false,
        title: 'name3',
        direction: 'left',
        name: 'name3',
        scheme: 'test3',
        desc: 'desc3'
      }
    ])
    expect(component.filterTbody).toStrictEqual([
      {
        key: 'checked',
        type: 'checkbox'
      },
      { key: 'uuid' },
      { key: 'name' },
      { key: 'scheme' }
    ])
    expect(component.filterThead).toStrictEqual([
      {
        type: 'checkbox'
      },
      { title: 'uuid' },
      { title: 'name' },
      { title: 'scheme' }
    ])
    // expect(component.filterForm.total).toStrictEqual(3)
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getRemoteList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getRemoteList('test')
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getRemoteList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getRemoteList('test')
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取数据失败!')
  })

  it('##getApiGroupList with success return, transferHeader test', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: {
          root: {
            groups: [
              {
                uuid: '123',
                name: '123',
                children: [
                  {
                    uuid: '1231',
                    name: '1231',
                    children: [{ uuid: '12311', name: '12311' }]
                  },
                  {
                    uuid: '1231',
                    name: '1231'
                  }
                ]
              },
              { uuid: '234', name: '234' },
              {
                uuid: '345',
                name: '345',
                children: [
                  {
                    uuid: '3451',
                    name: '3451'
                  }
                ]
              }
            ]
          }
        }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getApiGroupList()
    fixture.detectChanges()
    expect(component.apiGroupList).toStrictEqual([
      {
        uuid: '123',
        name: '123',
        children: [
          {
            uuid: '1231',
            name: '1231',
            children: [{ uuid: '12311', name: '12311', isLeaf: true }]
          },
          {
            uuid: '1231',
            name: '1231',
            isLeaf: true
          }
        ]
      },
      { uuid: '234', name: '234', isLeaf: true },
      {
        uuid: '345',
        name: '345',
        children: [
          {
            uuid: '3451',
            name: '3451',
            isLeaf: true
          }
        ]
      }
    ])
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getApiGroupList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getApiGroupList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getApiGroupList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getApiGroupList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取数据失败!')
  })

  it('##disabledEdit test', fakeAsync(() => {
    component.nzDisabled = false
    component.disabledEdit(true)
    expect(component.nzDisabled).toStrictEqual(true)
    component.disabledEdit(false)
    expect(component.nzDisabled).toStrictEqual(false)
  }))

  it('##getSearchRemoteList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: {
          target: 'test1',
          titles: [
            { title: 'uuid', field: 'uuid' },
            { title: 'name', field: 'name' },
            { title: 'scheme', field: 'scheme' }
          ],
          test1: [
            { uuid: '111', name: 'name1', scheme: 'test1', desc: 'desc1' },
            { uuid: '222', name: 'name2', scheme: 'test2', desc: 'desc2' },
            { uuid: '333', name: 'name3', scheme: 'test3', desc: 'desc3' }
          ],
          total: 3
        }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getSearchRemoteList()
    fixture.detectChanges()
    expect(component.remoteList).toStrictEqual([
      {
        uuid: '111',
        key: '111',
        checked: false,
        title: 'name1',
        direction: 'left',
        name: 'name1',
        scheme: 'test1',
        desc: 'desc1'
      },
      {
        uuid: '222',
        key: '222',
        checked: false,
        title: 'name2',
        direction: 'left',
        name: 'name2',
        scheme: 'test2',
        desc: 'desc2'
      },
      {
        uuid: '333',
        key: '333',
        checked: false,
        title: 'name3',
        direction: 'left',
        name: 'name3',
        scheme: 'test3',
        desc: 'desc3'
      }
    ])
    expect(component.filterForm.total).toStrictEqual(3)
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getSearchRemoteList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getSearchRemoteList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getSearchRemoteList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: '' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getSearchRemoteList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('筛选失败!')
  })

  it('##changeFilterType test', () => {
    component.filterTypeMap.set('api', {
      name: 'api',
      title: 'API',
      label: 'API',
      value: 'api',
      type: 'remote',
      options: undefined,
      total: undefined,
      allChecked: true,
      values: []
    })
    component.filterTypeMap.set('methods', {
      name: 'methods',
      title: 'API请求方式',
      label: 'API请求方式',
      value: 'methods',
      type: 'static',
      allChecked: true,
      options: [
        'ALL',
        'POST',
        'PUT',
        'GET',
        'DELETE',
        'OPTION',
        'PATCH',
        'HEAD'
      ],
      total: 7,
      values: []
    })
    component.filterTypeMap.set('path', {
      name: 'path',
      title: 'API路径',
      label: 'API路径',
      value: 'path',
      type: 'pattern',
      options: undefined,
      total: undefined,
      pattern: '1',
      values: []
    })

    const spyUpdateAllChecked = jest.spyOn(component, 'updateAllChecked')
    const spyGetRemoteList = jest.spyOn(component, 'getRemoteList')
    component.filterForm.allChecked = true
    component.changeFilterType('methods')

    expect(component.filterType).toStrictEqual('static')
    expect(component.filterForm).toStrictEqual({
      name: 'methods',
      title: 'API请求方式',
      label: 'API请求方式',
      value: 'methods',
      type: 'static',
      allChecked: true,
      options: [
        'ALL',
        'POST',
        'PUT',
        'GET',
        'DELETE',
        'OPTION',
        'PATCH',
        'HEAD'
      ],
      total: 7,
      showAll: true,
      values: []
    })
    expect(component.staticsList).toStrictEqual([
      { label: 'POST', value: 'POST', checked: true },
      { label: 'PUT', value: 'PUT', checked: true },
      { label: 'GET', value: 'GET', checked: true },
      { label: 'DELETE', value: 'DELETE', checked: true },
      { label: 'OPTION', value: 'OPTION', checked: true },
      { label: 'PATCH', value: 'PATCH', checked: true },
      { label: 'HEAD', value: 'HEAD', checked: true }
    ])
    expect(spyUpdateAllChecked).toHaveBeenCalledTimes(1)
    expect(spyGetRemoteList).not.toBeCalled()

    component.changeFilterType('api')
    fixture.detectChanges()

    expect(component.filterType).toStrictEqual('remote')
    expect(component.filterForm).toStrictEqual({
      name: 'api',
      title: 'API',
      label: 'API',
      value: 'api',
      type: 'remote',
      allChecked: true,
      options: undefined,
      total: undefined,
      values: []
    })
    expect(spyUpdateAllChecked).toHaveBeenCalledTimes(1)
    expect(spyGetRemoteList).toHaveBeenCalledTimes(1)

    component.changeFilterType('path')

    expect(component.filterType).toStrictEqual('pattern')
    expect(component.filterForm).toStrictEqual({
      name: 'path',
      title: 'API路径',
      label: 'API路径',
      value: 'path',
      type: 'pattern',
      options: undefined,
      total: undefined,
      pattern: /1/,
      values: []
    })
    expect(component.filterForm.pattern).toStrictEqual(/1/)
    expect(spyUpdateAllChecked).toHaveBeenCalledTimes(1)
    expect(spyGetRemoteList).toHaveBeenCalledTimes(1)
  })

  it('##change test', fakeAsync(() => {
    const ret: TransferChange = {
      from: 'left',
      to: 'right',
      list: [
        {
          uuid: '18a92926-8e01-9c75-897d-614d005c48de',
          name: 'APITest1',
          service: 'QQQQService',
          group: '应用管理/上线管理',
          requestPath: '2',
          key: '18a92926-8e01-9c75-897d-614d005c48de',
          checked: false,
          title: 'APITest1',
          direction: 'right',
          hide: false
        },
        {
          uuid: '74f83c4b-2199-66c6-3e1c-fc5b5f4d54a3',
          name: '1',
          service: 'test_baidu_service',
          group: 'API管理/上线管理',
          requestPath: '2',
          key: '74f83c4b-2199-66c6-3e1c-fc5b5f4d54a3',
          checked: false,
          title: '1',
          direction: 'right',
          hide: false
        }
      ]
    }

    component.remoteList = [
      {
        uuid: '18a92926-8e01-9c75-897d-614d005c48de',
        name: 'APITest1',
        service: 'QQQQService',
        group: '应用管理/上线管理',
        requestPath: '2',
        key: '18a92926-8e01-9c75-897d-614d005c48de',
        checked: false,
        title: 'APITest1',
        hide: false
      },
      {
        uuid: '74f83c4b-2199-66c6-3e1c-fc5b5f4d54a3',
        name: '1',
        service: 'test_baidu_service',
        group: 'API管理/上线管理',
        requestPath: '2',
        key: '74f83c4b-2199-66c6-3e1c-fc5b5f4d54a3',
        checked: false,
        title: '1',
        hide: false
      },
      {
        uuid: '111',
        key: '111',
        checked: false,
        title: 'name1',
        direction: 'left',
        name: 'name1',
        scheme: 'test1',
        desc: 'desc1'
      },
      {
        uuid: '222',
        key: '222',
        checked: false,
        title: 'name2',
        direction: 'left',
        name: 'name2',
        scheme: 'test2',
        desc: 'desc2'
      },
      {
        uuid: '333',
        key: '333',
        checked: false,
        title: 'name3',
        direction: 'left',
        name: 'name3',
        scheme: 'test3',
        desc: 'desc3'
      }
    ]

    component.change(ret)
    fixture.detectChanges()
    tick(300)
    expect(component.remoteSelectList).toStrictEqual([
      '18a92926-8e01-9c75-897d-614d005c48de',
      '74f83c4b-2199-66c6-3e1c-fc5b5f4d54a3'
    ])
    expect(component.remoteSelectNameList).toStrictEqual(['APITest1', '1'])

    const ret2: TransferChange = {
      from: 'right',
      to: 'left',
      list: [
        {
          uuid: '74f83c4b-2199-66c6-3e1c-fc5b5f4d54a3',
          name: '1',
          service: 'test_baidu_service',
          group: 'API管理/上线管理',
          requestPath: '2',
          key: '74f83c4b-2199-66c6-3e1c-fc5b5f4d54a3',
          checked: false,
          title: '1',
          hide: false
        }
      ]
    }
    component.change(ret2)
    fixture.detectChanges()
    tick(300)
    expect(component.remoteSelectList).toStrictEqual([
      '18a92926-8e01-9c75-897d-614d005c48de'
    ])
    expect(component.remoteSelectNameList).toStrictEqual(['APITest1'])
  }))

  it('##updateSingleChecked test', fakeAsync(() => {
    component.staticsList = [
      { value: 'test1', checked: true },
      { value: 'test2', checked: true }
    ]
    component.updateSingleChecked()
    fixture.detectChanges()

    expect(component.filterForm.allChecked).toStrictEqual(true)
    expect(component.filterForm.values).toStrictEqual(['test1', 'test2'])

    component.staticsList = [
      { value: 'test1', checked: true },
      { value: 'test2', checked: false }
    ]
    component.updateSingleChecked()
    fixture.detectChanges()

    expect(component.filterForm.allChecked).toStrictEqual(false)
    expect(component.filterForm.values).toStrictEqual(['test1'])

    component.staticsList = [
      { value: 'test1', checked: false },
      { value: 'test2', checked: false }
    ]
    component.updateSingleChecked()
    fixture.detectChanges()

    expect(component.filterForm.allChecked).toStrictEqual(false)
    expect(component.filterForm.values).toStrictEqual([])
  }))

  it('checkPattern test', fakeAsync(() => {
    component.filterForm.values = ['1']
    component.filterForm.pattern = /\w+/
    component.filterForm.patternIsPass = true

    component.checkPattern()
    expect(component.filterForm.patternIsPass).toStrictEqual(true)

    component.filterForm.values = ['  1']
    component.checkPattern()
    expect(component.filterForm.patternIsPass).toStrictEqual(true)

    component.filterForm.values = ['1']
    component.filterForm.pattern = /\W+/
    component.filterForm.patternIsPass = true

    component.checkPattern()
    expect(component.filterForm.patternIsPass).toStrictEqual(false)
  }))
})
