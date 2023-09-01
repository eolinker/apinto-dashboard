import { ComponentFixture, TestBed, discardPeriodicTasks, fakeAsync, tick } from '@angular/core/testing'
import { ComponentModule } from 'projects/core/src/app/component/component.module'
import { APP_BASE_HREF } from '@angular/common'
import { HttpClientModule } from '@angular/common/http'
import { ElementRef, Renderer2, ChangeDetectorRef } from '@angular/core'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { NavigationEnd, Router, RouterModule } from '@angular/router'
import { Overlay } from '@angular/cdk/overlay'
import { EoNgFeedbackMessageService, EoNgFeedbackModalModule, EoNgFeedbackModalService, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { BidiModule } from '@angular/cdk/bidi'
import { MockRenderer, MockMessageService, MockEnsureService, MockEmptySuccessResponse, MockRouterGroups, MockApiImport, MockGetCommonProviderService } from 'projects/core/src/app/constant/spec-test'
import { BehaviorSubject, of } from 'rxjs'
import { API_URL, ApiService } from 'projects/core/src/app/service/api.service'
import { environment } from 'projects/core/src/environments/environment'
import { NzFormModule } from 'ng-zorro-antd/form'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgTreeModule } from 'eo-ng-tree'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { EoNgSelectModule } from 'eo-ng-select'
import { LayoutModule } from '../../../layout.module'
import { routes } from '../../api-routing.module'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { ApiImportComponent } from './import.component'
import { NzFormatEmitEvent, NzTreeNode } from 'ng-zorro-antd/tree'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

describe('#init ApiImportComponent', () => {
  let component:ApiImportComponent
  let fixture: ComponentFixture<ApiImportComponent>
  let httpCommonService:any
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  let spyDeleteApiService:jest.SpyInstance<any>
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  let spyPostApiService:jest.SpyInstance<any>
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  let spyPutApiService:jest.SpyInstance<any>
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  let spyApiService:jest.SpyInstance<any>
  const eventsSub = new BehaviorSubject<any>(null)
  const routerStub = {
    events: eventsSub,
    url: '',
    navigate: (...args:Array<string>) => {
      eventsSub.next(new NavigationEnd(1, args.join('/'), args.join('/')))
    }
  }
  global.structuredClone = (val:any) => JSON.parse(JSON.stringify(val))

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule, FormsModule, ReactiveFormsModule, ComponentModule, LayoutModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule, NzOutletModule, HttpClientModule,
        RouterModule.forRoot(routes), NzFormModule, EoNgInputModule, EoNgTreeModule, EoNgButtonModule,
        EoNgSwitchModule, EoNgCheckboxModule, EoNgApintoTableModule, EoNgSelectModule, EoNgFeedbackModalModule,
        EoNgFeedbackTooltipModule, EoNgDropdownModule
      ],
      declarations: [ApiImportComponent
      ],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: Renderer2, useClass: MockRenderer },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: EoNgFeedbackModalService, useClass: MockEnsureService },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef },
        { provide: Router, useValue: routerStub }
      ],
      teardown: { destroyAfterEach: false }
    }).compileComponents()

    fixture = TestBed.createComponent(ApiImportComponent)
    component = fixture.componentInstance

    fixture.detectChanges()

    httpCommonService = fixture.debugElement.injector.get(ApiService)

    spyApiService = jest.spyOn(httpCommonService, 'get').mockImplementation(
      (...args) => {
        switch (args[0]) {
          case 'router/groups':
            return of(MockRouterGroups)
          case 'common/provider/Service':
            return of(MockGetCommonProviderService)
          default:
            return of(MockEmptySuccessResponse)
        }
      }
    )

    spyPostApiService = jest.spyOn(httpCommonService, 'post').mockReturnValue(
      of(MockApiImport)
    )

    spyPutApiService = jest.spyOn(httpCommonService, 'put').mockReturnValue(
      of(MockEmptySuccessResponse)
    )
  })

  it('should create and init component, test open modal', fakeAsync(() => {
    expect(component).toBeTruthy()

    component.ngAfterViewInit()
    fixture.detectChanges()

    expect(component.resultTableThead[0].click).not.toBeUndefined()
    expect(component.resultTableTbody[0].click).not.toBeUndefined()
    expect(component.resultTableTbody[2].check).not.toBeUndefined()
    expect(component.resultTableTbody[3].title).not.toBeUndefined()
    expect(component.importFormPage).toEqual(true)
    expect(component.token).toEqual('')
    expect(component.groupList).toEqual([])
    expect(component.upstreamList).toEqual([])
    expect(component.fileList).toEqual([])
    expect(component.validateForm.controls['file']).toBeUndefined()
    expect(component.modalRef).toBeUndefined()

    component.openDrawer()
    fixture.detectChanges()
    tick(100)

    expect(component.apisSet.size).toEqual(0)
    expect(component.importFormPage).toEqual(true)
    expect(component.token).toEqual('')
    expect(component.groupList).toEqual([
      {
        uuid: '50458642-5a9f-4136-9ff1-e30d647297e8',
        key: '50458642-5a9f-4136-9ff1-e30d647297e8',
        name: 'test1',
        title: 'test1',
        children: [
          {
            uuid: '35938ae4-1a62-4e22-ad8c-3691e111820e',
            key: '35938ae4-1a62-4e22-ad8c-3691e111820e',
            name: 'test1-c1',
            title: 'test1-c1',
            children: [],
            isLeaf: true,
            isDelete: false
          },
          {
            uuid: 'b238751a-dbfb-4610-8f40-a599737ac4e5',
            key: 'b238751a-dbfb-4610-8f40-a599737ac4e5',
            name: 'test1-c2',
            title: 'test1-c2',
            children: [],
            isLeaf: true,
            isDelete: false
          }
        ],
        isDelete: false
      },
      {
        uuid: '00db4977-331f-4b7e-93be-b64648751a5f',
        key: '00db4977-331f-4b7e-93be-b64648751a5f',
        name: 'test2',
        title: 'test2',
        children: [],
        isLeaf: true,
        isDelete: false
      }
    ])

    expect(component.upstreamList).toEqual([
      {
        value: 'test1cjk_service@upstream',
        label: 'test1[http]'
      },
      {
        value: 'test2@upstream',
        label: 'test2[http]'
      }])

    expect(component.fileList).toEqual([])
    expect(component.validateForm.controls['file'].value).toBeNull()
    expect(component.modalRef).not.toBeUndefined()

    component.checkConflict()
    fixture.detectChanges()
    tick(100)

    expect(spyPostApiService).not.toHaveBeenCalled()

    component.fileList.push({ uid: 'test', name: 'test' })
    component.validateForm.controls['group'].setValue(component.groupList[0].uuid)
    component.validateForm.controls['upstream'].setValue(component.upstreamList[0].value)

    component.checkConflict()
    fixture.detectChanges()
    tick(100)

    expect(component.importFormPage).toEqual(false)
    expect(component.resultList).toEqual([
      {
        id: 1,
        name: 'Returns pet inventories by status',
        method: 'GET',
        path: '/asda/store/inventory',
        desc: 'Returns a map of status codes to quantities',
        status: 'normal',
        disabled: false,
        checked: true,
        statusString: '正常'
      },
      {
        id: 2,
        name: 'Delete purchase order by ID',
        method: 'DELETE',
        path: '/asda/store/order/{orderId}',
        desc: 'For valid response try integer IDs with value < 1000. Anything above 1000 or nonintegers will generate API errors',
        status: 'conflict',
        disabled: true,
        checked: true,
        statusString: '冲突'
      },
      {
        id: 3,
        name: 'Find purchase order by ID',
        method: 'GET',
        path: '/asda/store/order/{orderId}',
        desc: 'For valid response try integer IDs with value <= 5 or > 10. Other values will generate exceptions.',
        status: 'invalid',
        disabled: true,
        checked: true,
        statusString: '无效path'
      },
      {
        id: 4,
        name: 'Returns pet inventories by status',
        method: 'GET',
        path: '/asda/store/inventory',
        desc: 'Returns a map of status codes to quantities',
        status: 'normal',
        disabled: false,
        checked: true,
        statusString: '正常'
      }
    ])
    expect(component.apisSet.size).toEqual(2)
    expect(component.apisSet.has(1)).toEqual(true)
    expect(component.apisSet.has(4)).toEqual(true)
    expect(component.token).toEqual('tokenForTest')

    component.resultTableThead[0].click && component.resultTableThead[0].click(false, 'all')
    fixture.detectChanges()

    expect(component.apisSet.size).toEqual(0)

    component.resultTableThead[0].click && component.resultTableThead[0].click(true, 'all')
    fixture.detectChanges()

    expect(component.apisSet.size).toEqual(2)

    component.resultTableTbody[0].click && component.resultTableTbody[0].click(component.resultList[3])
    fixture.detectChanges()

    expect(component.apisSet.size).toEqual(1)

    component.resultList[3].checked = false
    component.resultTableTbody[0].click && component.resultTableTbody[0].click(component.resultList[3])
    fixture.detectChanges()

    expect(component.apisSet.size).toEqual(2)

    if (component.resultTableTbody[2].check) {
      // eslint-disable-next-line prefer-const

      expect(component.resultTableTbody[2].check('')).toEqual(false)
      expect(component.resultTableTbody[2].check('test')).toEqual(true)
    }

    expect(spyPutApiService).not.toHaveBeenCalled()
    component.importApis()

    expect(spyPutApiService).toHaveBeenCalled()
    expect(component.importBtnLoading).toEqual(false)

    component.modalRef = undefined
    component.resultMap.set(1, { name: '' })
    component.importApis()

    discardPeriodicTasks()
  }))

  it('test nzTreeClick', () => {
    expect(component).toBeTruthy()
    component.openDrawer()
    fixture.detectChanges()

    const node = new NzTreeNode(component.groupList[0])
    node.origin.selectable = false
    node.origin.expanded = false
    const value:NzFormatEmitEvent = { eventName: 'test', node }

    component.nzTreeClick(value)
    expect(node.origin.expanded).toEqual(true)
  })

  it('test upload file', () => {
    expect(component.fileList).toEqual([])
    expect(component.authFile).toBeUndefined()
    expect(component.fileError).toEqual(false)

    expect(component.beforeUpload({ uid: 'test', name: 'test' })).toEqual(false)
    expect(component.fileList.length).toEqual(1)
    expect(component.authFile).toEqual({ uid: 'test', name: 'test' })
    expect(component.fileError).toEqual(false)

    expect(component.removeFile()).toEqual(true)
    expect(component.fileList).toEqual([])
    expect(component.authFile).toBeUndefined()
    expect(component.fileError).toEqual(true)
  })

  it('test disabledEdit & nzCheckAddRow ', () => {
    expect(component.nzDisabled).toEqual(false)

    component.disabledEdit(true)

    expect(component.nzDisabled).toEqual(true)

    expect(component.nzCheckAddRow()).toEqual(false)
  })
})
