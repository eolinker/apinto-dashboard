import { NgModule } from '@angular/core'

import { CommonModule } from '@angular/common'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgTableModule } from 'eo-ng-table'
import { NzSpaceModule } from 'ng-zorro-antd/space'
import { DragDropModule } from '@angular/cdk/drag-drop'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { TableComponent } from './table/table.component'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgCopyModule } from 'eo-ng-copy'
@NgModule({
  declarations: [
    TableComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    EoNgTableModule,
    NzSpaceModule,
    ReactiveFormsModule,
    DragDropModule,
    EoNgButtonModule,
    EoNgFeedbackTooltipModule,
    EoNgSwitchModule,
    EoNgDropdownModule,
    EoNgCheckboxModule,
    EoNgCopyModule
  ],
  exports: [
    TableComponent
  ]
})
export class EoNgApintoTableModule { }
