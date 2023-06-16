import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { LogRetrievalRoutingModule } from './log-retrieval-routing.module'
import { LogRetrievalComponent } from './log-retrieval.component'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCollapseModule } from 'eo-ng-collapse'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { EoNgLogRetrievalTailComponent } from './tail-log/tail-log.component'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgCodeboxModule } from 'eo-ng-codebox'
import { DirectiveModule } from '../../directive/directive.module'
import { EoNgEmptyModule } from 'eo-ng-empty'
import { EoNgCascaderModule } from 'eo-ng-cascader'
import { NzSpinModule } from 'ng-zorro-antd/spin'

@NgModule({
  declarations: [
    LogRetrievalComponent,
    EoNgLogRetrievalTailComponent
  ],
  imports: [
    CommonModule,
    DirectiveModule,
    LogRetrievalRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    EoNgCheckboxModule,
    EoNgSelectModule,
    EoNgButtonModule,
    EoNgCollapseModule,
    EoNgApintoTableModule,
    EoNgCodeboxModule,
    EoNgEmptyModule,
    EoNgCascaderModule,
    NzSpinModule
  ]
})
export class LogRetrievalModule { }
