import { ComponentFixture, TestBed } from '@angular/core/testing'

import { ResetPswComponent } from './reset-psw.component'

describe('ResetPswComponent', () => {
  let component: ResetPswComponent
  let fixture: ComponentFixture<ResetPswComponent>

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ResetPswComponent]
    })
      .compileComponents()

    fixture = TestBed.createComponent(ResetPswComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('should create', () => {
    expect(component).toBeTruthy()
  })
})
