/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-14 22:56:33
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-14 23:02:19
 * @FilePath: /apinto/src/app/layout/upstream/service-discovery-content/service-discovery-content.component.spec.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { API_URL } from 'projects/core/src/app/service/api.service'
import { RouterModule } from '@angular/router'
import { ElementRef, Renderer2, ChangeDetectorRef, Type } from '@angular/core'
import { APP_BASE_HREF } from '@angular/common'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { environment } from 'projects/core/src/environments/environment'
import { BidiModule } from '@angular/cdk/bidi'
import { Overlay } from '@angular/cdk/overlay'
import { UpstreamContentComponent } from './content.component'
import { UpstreamListComponent } from '../list/list.component'
import { UpstreamMessageComponent } from '../message/message.component'
import { UpstreamPublishComponent } from '../publish/publish.component'
class MockRenderer {
  removeAttribute (element: any, cssClass: string) {
    return cssClass + 'is removed from' + element
  }
}

describe('UpstreamContentComponent test', () => {
  let component: UpstreamContentComponent
  let fixture: ComponentFixture<UpstreamContentComponent>
  let renderer2: Renderer2
  class MockElementRef extends ElementRef {
    constructor () { super(null) }
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
        BidiModule, NoopAnimationsModule, NzNoAnimationModule,
        NzDrawerModule, NzOutletModule, HttpClientModule,
        RouterModule.forRoot([
          {
            path: '',
            component: UpstreamListComponent
          },
          {
            path: 'message',
            component: UpstreamMessageComponent
          },
          {
            path: 'publish',
            component: UpstreamPublishComponent
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
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(UpstreamContentComponent)
    renderer2 = fixture.componentRef.injector.get<Renderer2>(Renderer2 as Type<Renderer2>)
    renderer2.removeAttribute = jest.fn().mockReturnValue('remove')

    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('should create', () => {
    expect(component).toBeTruthy()
  })

  it('should initial clusterDesc and tabOptions after ngAfterViewInit', fakeAsync(() => {
    component.tabOptions = []
    component.ngAfterViewInit()
    fixture.detectChanges()
    expect(component.tabOptions).not.toStrictEqual([])
  }))

  it('should remove hidden attribute from tabs ngAfterViewChecked', fakeAsync(() => {
    // @ts-ignore
    const spyRemoveAttr = jest.spyOn(component.renderer, 'removeAttribute')
    expect(spyRemoveAttr).toHaveBeenCalledTimes(1)
    component.ngAfterViewChecked()
    fixture.detectChanges()
    expect(spyRemoveAttr).toHaveBeenCalledTimes(3)
  }))
})
