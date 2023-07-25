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
import { MatchTableComponent } from './table.component'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}
const testData = { position: 'position', key: 'key', matchType: 'matchType', pattern: 'pattern' }

describe('#init MatchTableComponent', () => {
  let component:MatchTableComponent
  let fixture: ComponentFixture<MatchTableComponent>
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
        EoNgSwitchModule, EoNgCheckboxModule, EoNgApintoTableModule, EoNgSelectModule,
        EoNgDropdownModule
      ],
      declarations: [MatchTableComponent
      ],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef },
        { provide: Router, useValue: routerStub }
      ],
      teardown: { destroyAfterEach: false }
    }).compileComponents()

    fixture = TestBed.createComponent(MatchTableComponent)
    component = fixture.componentInstance

    fixture.detectChanges()
  })

  it('should create and init component, test open modal', () => {
    expect(component).toBeTruthy()
    expect(component.modalRef).toBeUndefined()
    expect(component.editData).toBeUndefined()
    expect(component.validateMatchForm.controls['position'].value).toEqual('')

    component.ngOnInit()
    component.matchTableBody[4].btns[0].click({ data: testData })

    expect(component.matchTableBody[4].btns[0].disabledFn()).toEqual(false)
    expect(component.matchTableBody[4].btns[1].disabledFn()).toEqual(false)
    expect(component.editData).toEqual(testData)
    expect(component.validateMatchForm.controls['position'].value).toEqual('position')
    expect(component.modalRef).not.toBeUndefined()

    component.modalRef = undefined

    testData.key = 'key1'
    component.matchTableClick({ data: testData })

    expect(component.editData).toEqual(testData)
    expect(component.validateMatchForm.controls['key'].value).toEqual('key1')
    expect(component.modalRef).not.toBeUndefined()

    component.openDrawer('match')
    expect(component.validateMatchForm.controls['position'].value).toEqual('')
    expect(component.validateMatchForm.controls['key'].value).toEqual('')
    expect(component.validateMatchForm.controls['matchType'].value).toEqual('')
    expect(component.validateMatchForm.controls['pattern'].value).toEqual('')
  })

  it('test close modal', () => {
    component.openDrawer('match')
    const spyCloseModal = jest.spyOn(component.modalRef!, 'close')
    expect(spyCloseModal).not.toHaveBeenCalled()

    component.closeDrawer()

    expect(spyCloseModal).toHaveBeenCalled()

    expect(component._matchList).toEqual([])
    component.matchList = [testData]

    expect(component._matchList).toEqual([testData])
  })
})
