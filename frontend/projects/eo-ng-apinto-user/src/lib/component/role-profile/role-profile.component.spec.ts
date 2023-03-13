import { ComponentFixture, TestBed } from '@angular/core/testing'

import { RoleProfileComponent } from './role-profile.component'

describe('RoleProfileComponent', () => {
  let component: RoleProfileComponent
  let fixture: ComponentFixture<RoleProfileComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [RoleProfileComponent]
    })
      .compileComponents()

    fixture = TestBed.createComponent(RoleProfileComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })
})
