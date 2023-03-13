import { ComponentFixture, TestBed } from '@angular/core/testing'

import { RoleAccessComponent } from './role-access.component'

describe('RoleAccessComponent', () => {
  let component: RoleAccessComponent
  let fixture: ComponentFixture<RoleAccessComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [RoleAccessComponent]
    })
      .compileComponents()

    fixture = TestBed.createComponent(RoleAccessComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })
})
