import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ChangeUserRoleComponent } from './change-user-role.component';

describe('ChangeUserRoleComponent', () => {
  let component: ChangeUserRoleComponent;
  let fixture: ComponentFixture<ChangeUserRoleComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ChangeUserRoleComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ChangeUserRoleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
