import { NgModule } from '@angular/core'
import { CommonModule, registerLocaleData } from '@angular/common'

import { AuditLogRoutingModule } from './audit-log-routing.module'
import { AuditLogComponent } from './audit-log.component'
import { AuditLogListComponent } from './list/list.component'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgTableModule } from 'eo-ng-table'
import { ComponentModule } from '../../component/component.module'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'

import zh from '@angular/common/locales/zh'
import { AuditLogDetailComponent } from './detail/detail.component'
import { EoNgFeedbackTooltipModule } from 'eo-ng-feedback'

registerLocaleData(zh)

@NgModule({
  declarations: [
    AuditLogComponent,
    AuditLogListComponent,
    AuditLogDetailComponent
  ],
  imports: [
    CommonModule,
    AuditLogRoutingModule,
    EoNgInputModule,
    EoNgButtonModule,
    EoNgSelectModule,
    EoNgTableModule,
    FormsModule,
    ReactiveFormsModule,
    EoNgDatePickerModule,
    ComponentModule,
    EoNgApintoTableModule,
    EoNgFeedbackTooltipModule
  ]
})
export class AuditLogModule { }
