import { ComponentFixture, TestBed } from '@angular/core/testing'
import { ComponentModule } from 'projects/core/src/app/component/component.module'
import { APP_BASE_HREF } from '@angular/common'
import { HttpClientModule } from '@angular/common/http'
import { ElementRef, ChangeDetectorRef } from '@angular/core'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { NavigationEnd, Router, RouterModule } from '@angular/router'
import { Overlay } from '@angular/cdk/overlay'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { BidiModule } from '@angular/cdk/bidi'
import { BehaviorSubject } from 'rxjs'
import { API_URL } from 'projects/core/src/app/service/api.service'
import { environment } from 'projects/core/src/environments/environment'
import { LayoutModule } from '../../../layout.module'
import { routes } from '../../api-routing.module'
import { ApiMessageComponent } from './message.component'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

describe('#init ApiMessageComponent', () => {
  let component:ApiMessageComponent
  let fixture: ComponentFixture<ApiMessageComponent>

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
        RouterModule.forRoot(routes)
      ],
      declarations: [ApiMessageComponent
      ],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef },
        { provide: Router, useValue: routerStub }
      ],
      teardown: { destroyAfterEach: false }
    }).compileComponents()

    fixture = TestBed.createComponent(ApiMessageComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create and init component', () => {
    expect(component).toBeTruthy()
    expect(component.apiUuid).toBeUndefined()

    // @ts-ignore
    jest.replaceProperty(fixture.debugElement.injector.get(BaseInfoService), '_allParams', {
      apiId: 'mockApiId'
    })
    fixture = TestBed.createComponent(ApiMessageComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
    expect(component.apiUuid).toEqual('mockApiId')
  })
})
