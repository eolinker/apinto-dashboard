import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { UserAccessDirective } from './user-access.directive'
import { TreeDragDirective } from './tree-drag.directive'
import { RowExpandButtonDirective } from './row-expand-button.directive'
import { AutoFocusDirective } from './auto-focus.directive'

@NgModule({
  declarations: [UserAccessDirective, TreeDragDirective, RowExpandButtonDirective, AutoFocusDirective],
  imports: [
    CommonModule
  ],
  exports: [
    UserAccessDirective, TreeDragDirective,
    RowExpandButtonDirective, AutoFocusDirective
  ]
})
export class DirectiveModule { }
