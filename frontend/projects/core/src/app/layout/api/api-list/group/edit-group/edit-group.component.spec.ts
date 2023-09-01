import { ComponentFixture, TestBed, discardPeriodicTasks, fakeAsync } from '@angular/core/testing'
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
import { MockRenderer, MockMessageService, MockEnsureService, MockEmptySuccessResponse } from 'projects/core/src/app/constant/spec-test'
import { BehaviorSubject, of } from 'rxjs'
import { API_URL, ApiService } from 'projects/core/src/app/service/api.service'
import { environment } from 'projects/core/src/environments/environment'
import { NzFormModule } from 'ng-zorro-antd/form'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgButtonModule } from 'eo-ng-button'
import { LayoutModule } from '../../../../layout.module'
import { routes } from '../../../api-routing.module'
import { ApiManagementEditGroupComponent } from './edit-group.component'
import { EoNgCopyModule } from 'eo-ng-copy'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

describe('#init ApiManagementEditGroupComponent', () => {
  let component:ApiManagementEditGroupComponent
  let fixture: ComponentFixture<ApiManagementEditGroupComponent>
  let httpCommonService:any
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  let spyPutApiService:jest.SpyInstance<any>
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  let spyPostApiService:jest.SpyInstance<any>
  let spyCloseModal:jest.SpyInstance<any>
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
        RouterModule.forRoot(routes), NzFormModule, EoNgInputModule, EoNgButtonModule,
        EoNgFeedbackModalModule, EoNgFeedbackTooltipModule, EoNgCopyModule
      ],
      declarations: [ApiManagementEditGroupComponent
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

    fixture = TestBed.createComponent(ApiManagementEditGroupComponent)
    component = fixture.componentInstance
    fixture.detectChanges()

    httpCommonService = fixture.debugElement.injector.get(ApiService)

    spyPutApiService = jest.spyOn(httpCommonService, 'put').mockReturnValue(
      of(MockEmptySuccessResponse)
    )

    spyPostApiService = jest.spyOn(httpCommonService, 'post').mockReturnValue(
      of(MockEmptySuccessResponse)
    )
  })

  it('should create and init component', () => {
    expect(component).toBeTruthy()
    expect(component.groupName).toEqual('')
    expect(component.validateApiGroupForm.controls['groupName'].value).toEqual('')

    component.groupName = 'test'
    component.ngOnInit()
    fixture.detectChanges()

    expect(component.validateApiGroupForm.controls['groupName'].value).toEqual('test')
  })

  it('test add group', fakeAsync(() => {
    expect(component).toBeTruthy()
    expect(component.groupName).toEqual('')
    expect(component.validateApiGroupForm.controls['groupName'].value).toEqual('')
    component.closeModal = (args:any) => { console.log(args) }
    spyCloseModal = jest.spyOn(component, 'closeModal')

    component.ngOnInit()
    fixture.detectChanges()

    expect(component.validateApiGroupForm.controls['groupName'].value).toEqual('')
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(spyPutApiService).not.toHaveBeenCalled()
    component.addGroup('root')
    fixture.detectChanges()

    expect(spyPostApiService).not.toHaveBeenCalled()
    component.validateApiGroupForm.controls['groupName'].setValue('group1')
    expect(spyCloseModal).not.toHaveBeenCalled()

    component.addGroup('root')
    fixture.detectChanges()

    expect(spyCloseModal).toHaveBeenCalledWith(MockEmptySuccessResponse)
    expect(spyPostApiService).toHaveBeenCalledTimes(1)
    expect(spyPutApiService).not.toHaveBeenCalled()

    discardPeriodicTasks()
  }))

  it('test edit group', fakeAsync(() => {
    expect(component).toBeTruthy()
    expect(component.groupName).toEqual('')
    expect(component.validateApiGroupForm.controls['groupName'].value).toEqual('')
    component.closeModal = (args:any) => { console.log(args) }
    spyCloseModal = jest.spyOn(component, 'closeModal')

    component.ngOnInit()
    fixture.detectChanges()

    expect(component.validateApiGroupForm.controls['groupName'].value).toEqual('')
    expect(spyPostApiService).not.toHaveBeenCalled()
    expect(spyPutApiService).not.toHaveBeenCalled()

    component.editGroup('root')
    fixture.detectChanges()

    expect(spyPutApiService).not.toHaveBeenCalled()
    component.validateApiGroupForm.controls['groupName'].setValue('group1')
    expect(spyCloseModal).not.toHaveBeenCalled()

    component.editGroup('root')
    fixture.detectChanges()

    expect(spyCloseModal).toHaveBeenCalledWith(MockEmptySuccessResponse)
    expect(spyPutApiService).toHaveBeenCalledTimes(1)
    expect(spyPostApiService).not.toHaveBeenCalled()

    discardPeriodicTasks()
  }))
})
