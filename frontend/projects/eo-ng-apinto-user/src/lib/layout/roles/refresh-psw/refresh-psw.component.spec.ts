import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RefreshPswComponent } from './refresh-psw.component';

describe('RefreshPswComponent', () => {
  let component: RefreshPswComponent;
  let fixture: ComponentFixture<RefreshPswComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ RefreshPswComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RefreshPswComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
