import { ComponentFixture, TestBed } from '@angular/core/testing'

import { RolesGroupComponent } from './group.component'

describe('RolesGroupComponent', () => {
  let component: RolesGroupComponent
  let fixture: ComponentFixture<RolesGroupComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [RolesGroupComponent]
    })
      .compileComponents()

    fixture = TestBed.createComponent(RolesGroupComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })
})
