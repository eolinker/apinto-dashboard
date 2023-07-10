import { ComponentFixture, TestBed } from '@angular/core/testing'
import { ComponentModule } from 'projects/core/src/app/component/component.module'
import { APP_BASE_HREF } from '@angular/common'
import { HttpClientModule } from '@angular/common/http'
import { ElementRef, ChangeDetectorRef } from '@angular/core'
import { FormsModule, ReactiveFormsModule, Validators } from '@angular/forms'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { RouterModule } from '@angular/router'
import { Overlay } from '@angular/cdk/overlay'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { BidiModule } from '@angular/cdk/bidi'
import { API_URL } from 'projects/core/src/app/service/api.service'
import { environment } from 'projects/core/src/environments/environment'
import { NzFormModule } from 'ng-zorro-antd/form'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgButtonModule } from 'eo-ng-button'
import { LayoutModule } from '../../../layout.module'
import { routes } from '../../api-routing.module'
import { ApiManagementProxyComponent } from './proxy.component'
import { EoNgSelectModule } from 'eo-ng-select'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

describe('#init ApiManagementProxyComponent', () => {
  let component:ApiManagementProxyComponent
  let fixture: ComponentFixture<ApiManagementProxyComponent>
  global.structuredClone = (val:any) => JSON.parse(JSON.stringify(val))
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule, FormsModule, ReactiveFormsModule, ComponentModule, LayoutModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule, NzOutletModule, HttpClientModule,
        RouterModule.forRoot(routes), NzFormModule, EoNgInputModule, EoNgButtonModule, EoNgSelectModule
      ],
      declarations: [ApiManagementProxyComponent
      ],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef }
      ],
      teardown: { destroyAfterEach: false }
    }).compileComponents()

    fixture = TestBed.createComponent(ApiManagementProxyComponent)
    component = fixture.componentInstance

    fixture.detectChanges()
  })

  it('should create and init component', () => {
    expect(component).toBeTruthy()

    expect(component.editPage).toEqual(false)
    expect(component.nzDisabled).toEqual(false)
    expect(component.data).toEqual({})
    expect(component.listOfOptTypes.length).not.toEqual(0)
    expect(component.validateProxyHeaderForm.controls['key'].value).toEqual('')
    expect(component.validateProxyHeaderForm.controls['value'].value).toEqual('')
    expect(component.validateProxyHeaderForm.controls['optType'].value).toEqual('')
    expect(component.validateProxyHeaderForm.controls['optType'].hasValidator(Validators.required)).toEqual(true)
    expect(component.validateProxyHeaderForm.controls['key'].hasValidator(Validators.required)).toEqual(true)
    expect(component.validateProxyHeaderForm.controls['value'].hasValidator(Validators.required)).toEqual(false)

    component.editPage = true
    const data = { key: 'testKey', value: 'testValue', optType: 'type' }
    component.data = data

    component.ngOnInit()
    component.changeValidator()
    fixture.detectChanges()

    expect(component.validateProxyHeaderForm.controls['key'].value).toEqual('')
    expect(component.validateProxyHeaderForm.controls['value'].value).toEqual('')
    expect(component.validateProxyHeaderForm.controls['optType'].value).toEqual('type')
    expect(component.validateProxyHeaderForm.controls['optType'].hasValidator(Validators.required)).toEqual(true)
    expect(component.validateProxyHeaderForm.controls['key'].hasValidator(Validators.required)).toEqual(true)
    expect(component.validateProxyHeaderForm.controls['value'].hasValidator(Validators.required)).toEqual(true)

    component.data = data
    component.ngOnInit()
    component.validateProxyHeaderForm.controls['optType'].setValue('DELETE')
    component.changeValidator()
    fixture.detectChanges()

    expect(component.validateProxyHeaderForm.controls['key'].value).toEqual('')
    expect(component.validateProxyHeaderForm.controls['value'].value).toEqual('')
    expect(component.validateProxyHeaderForm.controls['optType'].hasValidator(Validators.required)).toEqual(true)
    expect(component.validateProxyHeaderForm.controls['key'].hasValidator(Validators.required)).toEqual(true)
    expect(component.validateProxyHeaderForm.controls['value'].hasValidator(Validators.required)).toEqual(false)
  })
})
