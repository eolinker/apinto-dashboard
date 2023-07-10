import { ComponentFixture, TestBed, discardPeriodicTasks, fakeAsync } from '@angular/core/testing'
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
import { MatchFormComponent } from './form.component'

export class MockElementRef extends ElementRef {
  constructor () { super(null) }
}

describe('#init MatchFormComponent', () => {
  let component:MatchFormComponent
  let fixture: ComponentFixture<MatchFormComponent>
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
      declarations: [MatchFormComponent
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

    fixture = TestBed.createComponent(MatchFormComponent)
    component = fixture.componentInstance

    fixture.detectChanges()
  })

  it('should create and init component, test open modal', fakeAsync(() => {
    expect(component).toBeTruthy()
    expect(component.accessUrl).toEqual('')

    routerStub.url = 'host/router/'
    component.ngOnInit()
    fixture.detectChanges()

    expect(component.accessUrl).toEqual('router/api')

    routerStub.url = 'host/serv-governance/'
    component.data = { position: 'header', key: 'key', matchType: 'match' }

    component.ngOnInit()
    fixture.detectChanges()

    expect(component.accessUrl).toEqual('serv-governance/grey')

    component.data = { position: 'header', key: 'key', matchType: 'match' }

    expect(component.validateMatchForm.value).toEqual({
      position: 'header', key: 'key', matchType: 'match', pattern: ''
    })

    discardPeriodicTasks()
  }))
  it('should create and init component, test open modal', () => {
    const spyMatchType = jest.spyOn(component.validateMatchForm.controls['matchType'], 'markAsDirty')
    const spyCloseDrawer = jest.spyOn(component.eoNgCloseDrawer, 'emit')
    expect(spyMatchType).not.toHaveBeenCalled()
    expect(spyCloseDrawer).not.toHaveBeenCalled()

    component.validateMatchForm.controls['pattern'].setValue('test')
    component.saveMatch()
    fixture.detectChanges()

    expect(component.validateMatchForm.controls['pattern'].value).toEqual('test')
    expect(spyMatchType).toHaveBeenCalled()

    component.validateMatchForm.controls['matchType'].setValue('NULL')
    component.validateMatchForm.controls['position'].setValue('HEADER')
    component.validateMatchForm.controls['key'].setValue('key')

    expect(!component.data).toEqual(true)
    expect(component.matchHeaderSet.size).toEqual(0)
    expect(component.matchList).toEqual([])
    expect(spyCloseDrawer).not.toHaveBeenCalled()

    component.saveMatch()
    fixture.detectChanges()

    expect(component.matchHeaderSet.size).toEqual(1)
    expect(component.matchList.length).toEqual(1)
    expect(component.matchList[0].key).toEqual('key')
    expect(component.validateMatchForm.controls['pattern'].value).toEqual('')
    expect(spyCloseDrawer).toHaveBeenCalled()

    component.validateMatchForm.controls['key'].setValue('key')
    component.validateMatchForm.controls['pattern'].setValue('test2')
    component.validateMatchForm.controls['matchType'].setValue('EXIST')
    component.saveMatch()
    fixture.detectChanges()

    expect(component.matchHeaderSet.size).toEqual(1)
    expect(component.matchList.length).toEqual(1)
    expect(component.matchList[0].matchType).toEqual('EXIST')
    expect(component.validateMatchForm.controls['pattern'].value).toEqual('')
    expect(spyCloseDrawer).toHaveBeenCalled()

    component.validateMatchForm.controls['key'].setValue('key3')
    component.validateMatchForm.controls['pattern'].setValue('test2')
    component.validateMatchForm.controls['matchType'].setValue('UNEXIST')
    component.data = { position: 'header', matchType: 'n', key: 'k' }
    component.editData = { position: 'HEADER', matchType: 'test2', key: 'key2', pattern: '' }
    component.saveMatch()
    fixture.detectChanges()

    expect(component.matchHeaderSet.size).toEqual(2)
    expect(component.matchList.length).toEqual(2)
    expect(component.matchList[0].key).toEqual('key3')
    expect(component.validateMatchForm.controls['pattern'].value).toEqual('')
    expect(spyCloseDrawer).toHaveBeenCalled()

    component.editData = { position: 'HEADER', matchType: 'UNEXIST', key: 'key3', pattern: '' }
    component.validateMatchForm.controls['key'].setValue('key4')
    component.validateMatchForm.controls['pattern'].setValue('test2')
    component.validateMatchForm.controls['matchType'].setValue('ANY')
    expect(component.editData).toEqual(component.matchList[0])
    component.saveMatch()
    fixture.detectChanges()

    expect(component.matchHeaderSet.size).toEqual(2)
    expect(component.matchList.length).toEqual(2)
    expect(component.matchList[0].key).toEqual('key4')
    expect(component.validateMatchForm.controls['pattern'].value).toEqual('')
    expect(spyCloseDrawer).toHaveBeenCalled()
  })
})
