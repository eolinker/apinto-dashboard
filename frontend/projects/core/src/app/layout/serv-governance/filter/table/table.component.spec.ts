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
import { FilterTableComponent } from './table.component'

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

describe('FilterTableComponent test as editPage is false', () => {
  let component: FilterTableComponent
  let fixture: ComponentFixture<FilterTableComponent>
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
      declarations: [FilterTableComponent],
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

    fixture = TestBed.createComponent(FilterTableComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('##should create', () => {
    expect(component).toBeTruthy()
  })

  it('##click table btns', () => {
    const spyOpenDrawer = jest.spyOn(component, 'openDrawer')
    // @ts-ignore
    const spyModalService = jest.spyOn(component.modalService, 'create')
    expect(spyOpenDrawer).toHaveBeenCalledTimes(0)
    expect(spyModalService).toHaveBeenCalledTimes(0)

    const item = { key: 'test' }
    component.filterTableBody[2].btns[0].click(item)
    expect(spyOpenDrawer).toHaveBeenCalledTimes(1)
    expect(spyModalService).toHaveBeenCalledTimes(0)

    component.filterTableBody[2].btns[1].click(item)
    expect(spyOpenDrawer).toHaveBeenCalledTimes(1)
    expect(spyModalService).toHaveBeenCalledTimes(1)

    expect(component.filterTableBody[2].btns[0].disabledFn()).toStrictEqual(false)
    expect(component.filterTableBody[2].btns[1].disabledFn()).toStrictEqual(false)
  })

  it('##filterTableClick test', fakeAsync(() => {
    const spyOpenDrawer = jest.spyOn(component, 'openDrawer')
    expect(spyOpenDrawer).not.toBeCalled()

    component.filterTableClick({ data: { value: ['ALL'] } })
    expect(spyOpenDrawer).toHaveBeenCalledTimes(1)
  }))

  it('##filterDelete test', () => {
    component.filterNamesSet = new Set(['test1', 'test2', 'test3'])
    const deleteF: any = { name: 'test2', label: '2', title: '2', values: [] }
    const filterList: any = [
      { name: 'test1', label: '1', title: '1', values: [] },
      deleteF,
      { name: 'test3', label: '3', title: '3', values: [] },
      { name: 'test4', label: '4', title: '4', values: [] }
    ]
    component.filterShowList = filterList
    component.filterDelete({ name: 'test' })
    expect(component.filterNamesSet).toStrictEqual(
      new Set(['test1', 'test2', 'test3'])
    )
    expect(component.filterShowList).toStrictEqual(filterList)

    component.filterShowList = filterList
    component.filterDelete({ name: 'test1' })
    expect(component.filterNamesSet).toStrictEqual(new Set(['test2', 'test3']))
    expect(component.filterShowList).toStrictEqual(filterList)

    component.filterShowList = filterList
    component.filterDelete(deleteF)
    expect(component.filterNamesSet).toStrictEqual(new Set(['test3']))
    expect(component.filterShowList).toStrictEqual([
      { name: 'test1', label: '1', title: '1', values: [] },
      { name: 'test3', label: '3', title: '3', values: [] },
      { name: 'test4', label: '4', title: '4', values: [] }
    ])
  })

  it('##openDrawer & cancelDrawer test', () => {
    expect(component.drawerFilterRef).toBeUndefined()
    component.openDrawer('addFilter')
    fixture.detectChanges()
    expect(component.drawerFilterRef).not.toBeUndefined()
    component.drawerFilterRef = undefined
    component.openDrawer('editFilter', {
      name: 'static',
      title: 'API请求方式',
      label: 'API请求方式',
      total: 3,
      values: ['ALL']
    })
    fixture.detectChanges()
    expect(component.drawerFilterRef).not.toBeUndefined()
    expect(component.filterForm).toStrictEqual({
      name: 'static',
      title: 'API请求方式',
      label: 'API请求方式',
      total: 3,
      allChecked: true,
      values: ['ALL']
    })
  })

  it('##openDrawer & cancelDrawer test', () => {
    expect(component.drawerFilterRef).toBeUndefined()
    component.openDrawer('addFilter')
    fixture.detectChanges()
    expect(component.drawerFilterRef).not.toBeUndefined()
    component.drawerFilterRef = undefined
    component.openDrawer('editFilter', {
      name: 'ip',
      title: 'API请求方式',
      label: 'API请求方式',
      total: 3,
      values: ['ttttt', 'eeeeeee', 'ssssssss', 'ttttttt']
    })
    fixture.detectChanges()
    expect(component.drawerFilterRef).not.toBeUndefined()
    expect(component.filterForm).toStrictEqual({
      name: 'ip',
      title: 'API请求方式',
      label: 'API请求方式',
      total: 3,
      values: ['ttttt\neeeeeee\nssssssss\nttttttt']
    })
    expect(component.editFilter).toStrictEqual({
      name: 'ip',
      title: 'API请求方式',
      label: 'API请求方式',
      total: 3,
      values: ['ttttt\neeeeeee\nssssssss\nttttttt']
    })

    // @ts-ignore
    component.drawerFilterRef.close = () => {
      return 'drawer is close'
    }
    // @ts-ignore
    const spyFn = jest.spyOn(component.drawerFilterRef, 'close')
    expect(spyFn).not.toHaveBeenCalled()
    component.drawerClose(false)
    expect(spyFn).not.toHaveBeenCalled()

    component.drawerClose(true)
    expect(spyFn).toHaveBeenCalledTimes(1)
  })

  it('##cleanFilterForm test', fakeAsync(() => {
    component.filterForm = {
      name: 'test',
      title: '',
      values: ['1'],
      label: '',
      text: 'test',
      allChecked: false,
      showAll: false,
      total: 0,
      groupUuid: 'test',
      pattern: null,
      patternIsPass: false
    }
    component.editFilter = { test: 'test' }
    component.remoteSelectList = ['2']
    component.remoteSelectNameList = ['test2']
    component.staticsList = ['test2']

    component.cleanFilterForm()

    expect(component.filterForm).toStrictEqual(
      {
        name: '',
        title: '',
        values: [],
        label: '',
        text: '',
        allChecked: false,
        showAll: false,
        total: 0,
        groupUuid: '',
        pattern: null,
        patternIsPass: true
      }
    )
    expect(component.editFilter).toStrictEqual(null)
    expect(component.remoteSelectList).toStrictEqual([])
    expect(component.remoteSelectNameList).toStrictEqual([])
    expect(component.staticsList).toStrictEqual([])
  }))
})
