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
import { MockRenderer, MockMessageService, MockEnsureService, MockEmptySuccessResponse, MockApiOnlineInfo, MockApiOnlineInfo2 } from 'projects/core/src/app/constant/spec-test'
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
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { LayoutModule } from '../../../../layout.module'
import { routes } from '../../../api-routing.module'
import { ApiPublishComponent } from './publish.component'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

describe('#init ApiPublishComponent', () => {
  let component:ApiPublishComponent
  let fixture: ComponentFixture<ApiPublishComponent>
  let httpCommonService:any
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
      declarations: [ApiPublishComponent
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

    fixture = TestBed.createComponent(ApiPublishComponent)
    component = fixture.componentInstance

    fixture.detectChanges()

    httpCommonService = fixture.debugElement.injector.get(ApiService)

    spyApiService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of(MockApiOnlineInfo)
    )

    spyPutApiService = jest.spyOn(httpCommonService, 'put').mockReturnValue(
      of(MockEmptySuccessResponse)
    )
  })

  it('should create and init component', fakeAsync(() => {
    expect(component).toBeTruthy()

    component.ngOnInit()
    component.ngAfterViewInit()
    fixture.detectChanges()
    tick(100)

    expect(component.apiInfo).not.toBeUndefined()
    expect(component.apiInfo?.scheme).toEqual('websocket')
    expect(component.apiInfo?.method).toEqual('-')
    expect(component.publishList).not.toEqual([])
    expect(component.publishTableBody).not.toEqual([])

    spyApiService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of(MockApiOnlineInfo2)
    )

    component.getApisData()
    fixture.detectChanges()
    tick(100)

    expect(component.apiInfo?.scheme).toEqual('http')
    expect(component.apiInfo?.method).toEqual('ALL')
    expect(component.publishList).not.toEqual([])

    discardPeriodicTasks()
  }))

  it('should select cluster and online', fakeAsync(() => {
    expect(component).toBeTruthy()

    component.ngOnInit()
    component.ngAfterViewInit()
    fixture.detectChanges()
    tick(100)

    expect(component.apiInfo).not.toBeUndefined()
    expect(component.apiInfo?.scheme).toEqual('websocket')
    expect(component.apiInfo?.method).toEqual('-')
    expect(component.publishList).not.toEqual([])
    expect(component.publishTableBody).not.toEqual([])

    const tableData = { checked: false, data: component.publishList }
    component.tableClick(tableData)
    tick(1)

    expect(tableData.checked).toEqual(true)
    expect(component.publishList[0].checked).toEqual(true)
    expect(component.selectedClusters).toEqual(component.publishList[0].title)
    expect(component.selectedNum).toEqual(1)
  }))

  it('should select cluster and offline', fakeAsync(() => {
  }))
})
