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
import { API_URL } from '../../service/api.service'
import { LoginComponent } from './login.component'
import { Overlay } from '@angular/cdk/overlay'

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

describe('LoginComponent test', () => {
  let component: LoginComponent
  let fixture: ComponentFixture<LoginComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        HttpClientModule,
        RouterModule.forRoot([
          {
            path: '**',
            component: LoginComponent
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
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(LoginComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('check Auth access and return success in init', () => {
    // @ts-ignore
    const spyService = jest.spyOn(component.api, 'checkAuth').mockReturnValue(of({ code: 0, data: { } }))
    // @ts-ignore
    const spyGetMenuList = jest.spyOn(component.navigationService, 'getMenuList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.ngOnInit()
    expect(spyService).toHaveBeenCalledTimes(1)
    expect(spyGetMenuList).toHaveBeenCalledTimes(1)
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('check Auth access and return fail in init', () => {
    // @ts-ignore
    const spyService = jest.spyOn(component.api, 'checkAuth').mockReturnValue(of({ code: -1, data: { } }))
    // @ts-ignore
    const spyGetMenuList = jest.spyOn(component.navigationService, 'getMenuList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyGetMenuList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.ngOnInit()
    fixture.detectChanges()

    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalledTimes(1)
  })

  it('unsubscribe in ngOnDestroy', fakeAsync(() => {
    // @ts-ignore
    const spyUnsubscribe = jest.spyOn(component.subscription, 'unsubscribe')
    expect(spyUnsubscribe).not.toHaveBeenCalled()
    component.ngOnDestroy()
    expect(spyUnsubscribe).toHaveBeenCalledTimes(1)
  }))
})
