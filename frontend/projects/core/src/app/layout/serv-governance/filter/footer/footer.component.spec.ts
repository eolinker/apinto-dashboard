import {
  ComponentFixture,
  fakeAsync,
  TestBed
} from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { API_URL } from 'projects/core/src/app/service/api.service'
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
import { CacheCreateComponent } from '../../cache/create/create.component'
import { FuseCreateComponent } from '../../fuse/create/create.component'
import { GreyCreateComponent } from '../../grey/create/create.component'
import { TrafficCreateComponent } from '../../traffic/create/create.component'
import { VisitCreateComponent } from '../../visit/create/create.component'
import { FilterFooterComponent } from './footer.component'

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

describe('FilterFooterComponent test as editPage is false', () => {
  let component: FilterFooterComponent
  let fixture: ComponentFixture<FilterFooterComponent>
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
      declarations: [FilterFooterComponent],
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
            queryParams: of({ cluster_name: 'clus2' })
          }
        }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(FilterFooterComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('##should create', () => {
    expect(component).toBeTruthy()
  })

  it('##saveFilter test', fakeAsync(() => {
    component.filterForm = {
      name: 'api',
      title: 'API',
      label: 'API',
      text: '',
      total: 3,
      values: ['5'],
      pattern: /./,
      patternIsPass: true
    }
    component.remoteSelectList = ['1', '2', '3']
    component.remoteSelectNameList = ['11', '22', '33']
    component.filterType = 'remote'

    component.saveFilter()
    expect(component.filterForm.values).toStrictEqual(['ALL'])
    expect(component.filterForm.text).toStrictEqual('所有API')
    expect(component.filterNamesSet).toStrictEqual(new Set(['api']))
    component.filterForm = {
      name: 'api',
      title: 'API',
      label: 'API',
      text: '',
      total: 3,
      values: ['5'],
      pattern: /./,
      patternIsPass: true
    }
    component.remoteSelectList = ['3']
    component.remoteSelectNameList = ['33']
    component.filterType = 'remote'

    component.saveFilter()
    expect(component.filterForm.values).toStrictEqual(['3'])
    expect(component.filterForm.text).toStrictEqual('33')

    component.filterType = 'static'
    component.filterForm = {
      name: 'methods',
      title: 'API请求方式',
      label: 'API请求方式',
      text: '',
      allChecked: false,
      total: 3,
      values: [],
      pattern: /./,
      patternIsPass: true
    }

    component.staticsList = [
      { label: 'ALL', value: 'ALL', checked: false },
      { label: 'POST', value: 'POST', checked: true },
      { label: 'PUT', value: 'PUT', checked: false },
      { label: 'GET', value: 'GET', checked: true }
    ]

    component.saveFilter()
    expect(component.filterForm.values).toStrictEqual(['POST', 'GET'])
    expect(component.filterForm.text).toStrictEqual('POST,GET')

    component.filterForm = {
      name: 'methods',
      title: 'API请求方式',
      label: 'API请求方式',
      text: '',
      allChecked: true,
      total: 3,
      values: [],
      pattern: /./,
      patternIsPass: true
    }

    component.saveFilter()
    expect(component.filterForm.values).toStrictEqual(['ALL'])
    expect(component.filterForm.text).toStrictEqual('所有API请求方式')

    component.filterForm = {
      name: 'path',
      title: 'API路径',
      label: 'API路径',
      text: '',
      total: undefined,
      values: ['1'],
      pattern: /./,
      patternIsPass: true
    }
    component.filterType = 'pattern'
    component.saveFilter()
    expect(component.filterForm.text).toStrictEqual('1')

    const ef: any = {
      name: 'methods',
      title: 'API请求方式',
      label: 'API请求方式',
      allChecked: true,
      total: 3,
      values: [],
      pattern: /./,
      patternIsPass: true
    }

    component.editFilter = ef

    // @ts-ignore
    component.filterForm = {
      title: 'test5',
      name: 'test5',
      text: 'test5',
      values: [],
      pattern: /./,
      patternIsPass: true
    }
    // @ts-ignore
    component.filterShowList = [
      { name: 'test1', label: '1', title: '1', values: [] },
      { name: 'test3', label: '3', title: '3', values: [] },
      ef,
      { name: 'test4', label: '4', title: '4', values: [] }
    ]

    component.filterNamesSet = new Set(['methods', 'test4'])

    component.filterType = 'static'
    component.saveFilter()
    expect(component.filterShowList).toStrictEqual([
      { name: 'test1', label: '1', title: '1', values: [] },
      { name: 'test3', label: '3', title: '3', values: [] },
      { name: 'test4', label: '4', title: '4', values: [] },
      {
        title: 'test5',
        name: 'test5',
        label: 'POST,GET',
        values: ['POST', 'GET']
      }
    ])
    expect(component.filterNamesSet).toStrictEqual(new Set(['test4', 'test5']))
  }))

  it('##checkFilterToSave test', fakeAsync(() => {
    component.filterType = 'static'
    // @ts-ignore
    component.filterForm = {
      filterType: 'static',
      allChecked: false,
      values: []
    }
    let res = component.checkFilterToSave()
    expect(res).toStrictEqual(true)

    // @ts-ignore
    component.filterForm = {
      filterType: 'static',
      allChecked: true,
      values: []
    }
    res = component.checkFilterToSave()
    expect(res).toStrictEqual(false)

    component.filterType = 'pattern'
    // @ts-ignore
    component.filterForm = {
      filterType: 'pattern',
      allChecked: true,
      pattern: /\w+/,
      patternIsPass: false,
      values: [' 1']
    }
    res = component.checkFilterToSave()
    expect(res).toStrictEqual(true)

    // @ts-ignore
    component.filterForm = {
      filterType: 'pattern',
      allChecked: true,
      pattern: /\w+/,
      patternIsPass: true,
      values: [' 1']
    }

    res = component.checkFilterToSave()
    expect(res).toStrictEqual(false)

    // @ts-ignore
    component.filterForm = {
      filterType: 'pattern',
      pattern: null,
      allChecked: true,
      values: []
    }

    res = component.checkFilterToSave()
    expect(res).toStrictEqual(true)

    // @ts-ignore
    component.filterForm = {
      filterType: 'pattern',
      pattern: null,
      allChecked: true,
      values: ['111:111:111']
    }

    res = component.checkFilterToSave()
    expect(res).toStrictEqual(false)

    component.filterType = 'remote'
    component.filterForm.allChecked = false
    component.remoteSelectList = []

    res = component.checkFilterToSave()
    expect(res).toStrictEqual(true)

    component.filterForm.allChecked = false
    component.remoteSelectList = ['1']

    res = component.checkFilterToSave()
    expect(res).toStrictEqual(false)

    component.filterForm.allChecked = true
    component.remoteSelectList = []

    res = component.checkFilterToSave()
    expect(res).toStrictEqual(false)
  }))
})
