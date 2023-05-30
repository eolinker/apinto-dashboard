/* eslint-disable dot-notation */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-07-30 00:40:51
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-03 23:20:22
 * @FilePath: /apinto/src/app/layout/basic-layout/basic-layout.component.spec.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { APP_BASE_HREF } from '@angular/common'
import { HttpClientModule } from '@angular/common/http'
import { ElementRef, Renderer2, ChangeDetectorRef } from '@angular/core'
import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'
import { RouterModule } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { environment } from 'projects/core/src/environments/environment'
import { of } from 'rxjs'
import { API_URL } from '../../../service/api.service'
import { Overlay } from '@angular/cdk/overlay'
import { PasswordComponent } from './password.component'
import { ReactiveFormsModule } from '@angular/forms'
import { NzFormModule } from 'ng-zorro-antd/form'
import { CryptoService } from '../../../service/crypto.service'

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

class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

describe('PasswordComponent test', () => {
  let component: PasswordComponent
  let fixture: ComponentFixture<PasswordComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        HttpClientModule,
        ReactiveFormsModule,
        NzFormModule,
        RouterModule.forRoot([
          {
            path: '**',
            component: PasswordComponent
          }
        ]
        )
      ],
      declarations: [
      ],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: Renderer2, useClass: MockRenderer },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: CryptoService, useClass: CryptoService },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef },
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(PasswordComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('initial validateForm in ngOnInit and focus on password input', fakeAsync(() => {
    const spyService = jest.spyOn(component.autoFocusInput.nativeElement, 'focus')
    component.validateForm.controls['username'].setValue('test')
    component.validateForm.controls['password'].setValue('test')

    expect(spyService).not.toHaveBeenCalled()
    component.ngOnInit()

    expect(component.validateForm.controls['username'].value).toStrictEqual(null)
    expect(component.validateForm.controls['password'].value).toStrictEqual(null)

    component.ngAfterViewInit()

    expect(spyService).toHaveBeenCalled()
  }))

  it('validateForm is valid and login is success without callbackUrl',fakeAsync(() => {
    // @ts-ignore
    const spyLogin = jest.spyOn(component.api, 'login').mockReturnValue({ code: 0, data: {} })
    // @ts-ignore
    const spyReqFlashMenu = jest.spyOn(component.appConfig, 'reqFlashMenu')
    // @ts-ignore
    const spyNavigate = jest.spyOn(component.router, 'navigate')
    // @ts-ignore
    component.validateForm.controls['username'].setValue('admin')
    component.validateForm.controls['password'].setValue('12345678')
    component.validateForm.updateValueAndValidity()
    expect(component.validateForm.valid).toStrictEqual(true)
    expect(spyLogin).not.toHaveBeenCalled()
    expect(spyReqFlashMenu).not.toHaveBeenCalled()
    expect(spyNavigate).not.toHaveBeenCalled()
    // @ts-ignore
    component.route.snapshot.queryParams['callback'] = null
    component.login()
    expect(spyLogin).toHaveBeenCalled()
    expect(spyLogin).toHaveBeenCalledWith({
      username:'admin', 
      password:'XebnGrBSaQuxIODzaGTGZw==',
      client: 1,
      type: 1,
      app_type: 4})
    expect(spyReqFlashMenu).toHaveBeenCalled()
    expect(spyNavigate).toHaveBeenCalled()
  }))

  it('validateForm is valid and login is success with callbackUrl', () => {

  })

  it('validateForm is valid and login is fail', () => {
  })

  it('validateForm is unvalid and login', () => {
  })
})
