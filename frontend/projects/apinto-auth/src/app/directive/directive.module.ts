import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { UserAccessDirective } from './user-access.directive'

@NgModule({
  declarations: [UserAccessDirective],
  imports: [
    CommonModule
  ],
  exports: [
    UserAccessDirective
  ]
})
export class DirectiveModule { }
