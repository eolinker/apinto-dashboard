import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'
import { UserAccessDirective } from './user-access.directive'
import { TreeDragDirective } from './tree-drag.directive'
import { RowExpandButtonDirective } from './row-expand-button.directive'
import { TabHostDirective } from './tab-host.directive';
import { AutoFocusDirective } from './auto-focus.directive'


@NgModule({
  declarations: [UserAccessDirective, TreeDragDirective, RowExpandButtonDirective, TabHostDirective, AutoFocusDirective],
  imports: [
    CommonModule
  ],
  exports: [
    UserAccessDirective, TreeDragDirective,
    RowExpandButtonDirective,TabHostDirective,AutoFocusDirective
  ]
})
export class DirectiveModule { }
