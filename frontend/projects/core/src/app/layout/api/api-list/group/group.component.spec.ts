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
import { MockRenderer, MockMessageService, MockEnsureService, MockEmptySuccessResponse, MockRouterGroups } from 'projects/core/src/app/constant/spec-test'
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
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { LayoutModule } from '../../../layout.module'
import { routes } from '../../api-routing.module'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { ApiManagementGroupComponent } from './group.component'
import { NzHighlightModule } from 'ng-zorro-antd/core/highlight'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

describe('#init ApiManagementGroupComponent', () => {
  let component:ApiManagementGroupComponent
  let fixture: ComponentFixture<ApiManagementGroupComponent>
  let httpCommonService:any
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  let spyDeleteApiService:jest.SpyInstance<any>
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
        EoNgFeedbackTooltipModule, EoNgDropdownModule, NzHighlightModule
      ],
      declarations: [ApiManagementGroupComponent
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

    fixture = TestBed.createComponent(ApiManagementGroupComponent)
    component = fixture.componentInstance

    fixture.detectChanges()

    httpCommonService = fixture.debugElement.injector.get(ApiService)

    spyApiService = jest.spyOn(httpCommonService, 'get').mockImplementation(
      (...args) => {
        switch (args[0]) {
          case 'router/groups':
            return of(MockRouterGroups)
          default:
            return of(MockEmptySuccessResponse)
        }
      }
    )

    spyDeleteApiService = jest.spyOn(httpCommonService, 'delete').mockReturnValue(
      of(MockEmptySuccessResponse)
    )
  })

  it('should create and init component', fakeAsync(() => {
    expect(component).toBeTruthy()
    expect(component.nodesList).toEqual([])
    expect(component.showAll).toEqual(true)
    // @ts-ignore
    jest.replaceProperty(fixture.debugElement.injector.get(BaseInfoService), '_allParams', {
      apiGroupId: 'mockApiGroupId'
    })
    component.queryName = 'test'
    component.ngOnInit()
    fixture.detectChanges()
    tick(500)
    expect(component.nodesList).toEqual([
      {
        uuid: '50458642-5a9f-4136-9ff1-e30d647297e8',
        key: '50458642-5a9f-4136-9ff1-e30d647297e8',
        name: 'test1',
        title: 'test1',
        expanded: false,
        children: [
          {
            uuid: '35938ae4-1a62-4e22-ad8c-3691e111820e',
            key: '35938ae4-1a62-4e22-ad8c-3691e111820e',
            name: 'test1-c1',
            title: 'test1-c1',
            children: [],
            isDelete: false
          },
          {
            uuid: 'b238751a-dbfb-4610-8f40-a599737ac4e5',
            key: 'b238751a-dbfb-4610-8f40-a599737ac4e5',
            name: 'test1-c2',
            title: 'test1-c2',
            children: [],
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
        isDelete: false
      }

    ])
    // @ts-ignore
    expect(component.baseInfo.allParamsInfo.apiGroupId).toEqual('mockApiGroupId')
    expect(component.queryName).toEqual('')
    expect(component.groupUuid).toEqual('mockApiGroupId')

    // @ts-ignore
    jest.replaceProperty(fixture.debugElement.injector.get(BaseInfoService), '_allParams', {
      apiGroupId: 'uuid2'
    })
    eventsSub.next(new NavigationEnd(1, 'home/', 'home/'))
    fixture.detectChanges()

    expect(component.groupUuid).toEqual('uuid2')
    discardPeriodicTasks()
  }))

  it('test deleteGroup', fakeAsync(() => {
    // @ts-ignore
    const spyModalService = jest.spyOn(component.modalService, 'create')
    const spyMenuList = jest.spyOn(component, 'getMenuList')
    const spyViewAllApis = jest.spyOn(component, 'viewAllApis')
    // delete a group when view all apis
    expect(component).toBeTruthy()
    expect(component.nodesList).toEqual([])
    expect(component.showAll).toEqual(true)
    expect(spyModalService).not.toHaveBeenCalled()
    // @ts-ignore
    jest.replaceProperty(fixture.debugElement.injector.get(BaseInfoService), '_allParams', {
      apiGroupId: 'mockApiGroupId'
    })
    component.queryName = 'test'

    component.ngOnInit()
    fixture.detectChanges()
    tick(500)

    expect(component.showAll).toEqual(true)
    expect(component.groupModal).toBeUndefined()
    expect(spyModalService).not.toHaveBeenCalled()

    component.deleteGroupModal('test', 'uuid')
    fixture.detectChanges()

    expect(component.showAll).toEqual(true)
    expect(component.groupModal).not.toBeUndefined()
    expect(spyModalService).toHaveBeenCalled()
    expect(spyDeleteApiService).not.toHaveBeenCalled()
    expect(spyMenuList).toHaveBeenCalledTimes(1)
    expect(spyViewAllApis).not.toHaveBeenCalled()

    component.deleteGroup(component.nodesList[1].key, component.nodesList[1].title)
    fixture.detectChanges()

    expect(component.showAll).toEqual(true)
    expect(spyDeleteApiService).toHaveBeenCalled()
    expect(spyMenuList).toHaveBeenCalledTimes(2)
    expect(spyViewAllApis).not.toHaveBeenCalled()
    expect(document.getElementsByTagName('eo-ng-tree-default-node')[0].checkVisibility).toEqual(true)
    // delete selected group
    // delete other group when selected one
    discardPeriodicTasks()
  }))
})
