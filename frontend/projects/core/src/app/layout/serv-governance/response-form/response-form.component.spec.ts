/* eslint-disable dot-notation */
import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { ActivatedRoute, RouterModule } from '@angular/router'
import { ElementRef, Renderer2, ChangeDetectorRef, SimpleChange, SimpleChanges } from '@angular/core'
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
import { FormsModule, ReactiveFormsModule, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms'

import { EoNgSelectModule } from 'eo-ng-select'
import { CacheCreateComponent } from '../cache/create/create.component'
import { FuseCreateComponent } from '../fuse/create/create.component'
import { GreyCreateComponent } from '../grey/create/create.component'
import { TrafficCreateComponent } from '../traffic/create/create.component'
import { VisitCreateComponent } from '../visit/create/create.component'
import { ResponseFormComponent } from './response-form.component'

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

describe('ResponseFormComponent test as editPage is false', () => {
  let component: ResponseFormComponent
  let fixture: ComponentFixture<ResponseFormComponent>
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
      declarations: [ResponseFormComponent],
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

    fixture = TestBed.createComponent(ResponseFormComponent)
    component = fixture.componentInstance
    fixture.detectChanges()

    component._responseForm = new UntypedFormGroup({
      statusCode: new UntypedFormControl([200, [Validators.required, Validators.pattern(/^[1-9]{1}\d{2}$/)]]),
      contentType: new UntypedFormControl(['application/json', [Validators.required]]),
      charset: new UntypedFormControl(['UTF-8', [Validators.required]]),
      header: new UntypedFormControl([]),
      body: new UntypedFormControl([])
    })
  })
  it('##should create', () => {
    expect(component).toBeTruthy()
  })

  it('##set responseHeaderList', () => {
    const spyResonseListChange = jest.spyOn(component.responseHeaderListChange, 'emit')
    expect(spyResonseListChange).not.toHaveBeenCalled()

    component.responseHeaderList = [{ key: 'key', value: 'value', test: 'test' }]
    fixture.detectChanges()

    expect(spyResonseListChange).toHaveBeenCalledTimes(1)
    // @ts-ignore
    expect(component._responseHeaderList).toStrictEqual([{ key: 'key', value: 'value', test: 'test' }])
  })

  it('##responseHeaderTableBody test', fakeAsync(() => {
    expect(component.responseHeaderTableBody[0].disabledFn()).toStrictEqual(false)
    expect(component.responseHeaderTableBody[1].disabledFn()).toStrictEqual(false)
    expect(component.responseHeaderTableBody[2].btns[0].disabledFn()).toStrictEqual(false)
    expect(component.responseHeaderTableBody[3].btns[0].disabledFn()).toStrictEqual(false)
    expect(component.responseHeaderTableBody[3].btns[1].disabledFn()).toStrictEqual(false)
  }))

  it('##ngOnInit should call getContentTypeList && getCharsetList', () => {
    const spyGetContentTypeList = jest.spyOn(component, 'getContentTypeList')
    const spyGetCharsetList = jest.spyOn(component, 'getCharsetList')
    expect(spyGetContentTypeList).not.toHaveBeenCalled()
    expect(spyGetCharsetList).not.toHaveBeenCalled()
    component.responseForm = new UntypedFormGroup({
      statusCode: new UntypedFormControl([200, [Validators.required, Validators.pattern(/^[1-9]{1}\d{2}$/)]]),
      contentType: new UntypedFormControl(['application/json', [Validators.required]]),
      charset: new UntypedFormControl(['UTF-8', [Validators.required]]),
      header: new UntypedFormControl([]),
      body: new UntypedFormControl([])
    })
    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetContentTypeList).toHaveBeenCalledTimes(1)
    expect(spyGetCharsetList).toHaveBeenCalledTimes(1)
  })

  it('##ngOnChanges and disabled test', () => {
    const spyFormDisabled = jest.spyOn(component.responseForm, 'disable')
    expect(spyFormDisabled).not.toHaveBeenCalled()
    component.disabled = true
    const changes:SimpleChanges = { disabled: new SimpleChange(false, true, false) }
    component.ngOnChanges(changes)
    expect(spyFormDisabled).toHaveBeenCalled()
  })

  it('##getContentTypeList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: { items: [{ contentType: 'test1', body: 'test1' }, { contentType: 'test2', body: 'test2' }, { contentType: 'test3', body: 'test3' }] }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component._responseForm.controls['contentType'].setValue('test1')
    component.getContentTypeList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component._contentTypeList).toStrictEqual([{ label: 'test1', value: 'test1' }, { label: 'test2', value: 'test2' }, { label: 'test3', value: 'test3' }])
    expect(component.contentTypeList).toStrictEqual([{ label: 'test1', value: 'test1' }, { label: 'test2', value: 'test2' }, { label: 'test3', value: 'test3' }])
    expect(component.contentTypeMap.get('test1')).toStrictEqual('test1')
    expect(component.contentTypeMap.get('test2')).toStrictEqual('test2')
    expect(component.contentTypeMap.get('test3')).toStrictEqual('test3')

    expect(component._responseForm.controls['body'].value).toStrictEqual('test1')

    component.editPage = true
    component._responseForm.controls['body'].setValue('')

    component.getContentTypeList()
    fixture.detectChanges()

    expect(component._responseForm.controls['body'].value).toStrictEqual('')
  })

  it('##getContentTypeList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getContentTypeList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getContentTypeList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getContentTypeList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取数据失败!')
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getCharsetList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: { items: ['test1', 'test2', 'test3'] }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getCharsetList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(component.charsetList).toStrictEqual([{ label: 'test1', value: 'test1' }, { label: 'test2', value: 'test2' }, { label: 'test3', value: 'test3' }])
  })

  it('##getCharsetList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getCharsetList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getCharsetList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getCharsetList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取数据失败!')
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##changeContentType test', fakeAsync(() => {
    component.contentTypeList = []
    component._contentTypeList = [{ label: 'test1', value: 'test1' }, { label: 'test2', value: 'test2' }, { label: 'test3', value: 'test3' }]

    component.changeContentType('testChange')

    expect(component.contentTypeList).toStrictEqual([{ label: 'testChange', value: 'testChange' }, { label: 'test1', value: 'test1' }, { label: 'test2', value: 'test2' }, { label: 'test3', value: 'test3' }])

    component.changeContentType('test3')
    expect(component.contentTypeList).toStrictEqual([{ label: 'test1', value: 'test1' }, { label: 'test2', value: 'test2' }, { label: 'test3', value: 'test3' }])
  }))
})
