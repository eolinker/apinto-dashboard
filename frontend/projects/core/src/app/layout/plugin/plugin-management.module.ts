import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { PluginManagementRoutingModule } from './plugin-management-routing.module'
import { PluginManagementComponent } from './plugin-management.component'
import { PluginListComponent } from './list/list.component'
import { PluginMessageComponent } from './message/message.component'
import { ComponentModule } from '../../component/component.module'
import { EoNgTreeModule } from 'eo-ng-tree'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgButtonModule } from 'eo-ng-button'
import { NzFormModule } from 'ng-zorro-antd/form'
import { EoNgFeedbackModalModule } from 'eo-ng-feedback'
import { EoNgRadioModule } from 'eo-ng-radio'
import { PluginCreateComponent } from './create/create.component'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { MarkdownModule } from 'ngx-markdown'
import { PluginConfigComponent } from './config/config.component'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgAutoCompleteModule } from 'eo-ng-auto-complete'
import { GroupComponent } from './group/group.component'
import { EoNgEmptyModule } from 'eo-ng-empty'
import { DirectiveModule } from '../../directive/directive.module'

@NgModule({
  declarations: [
    PluginManagementComponent,
    PluginListComponent,
    PluginMessageComponent,
    PluginCreateComponent,
    PluginConfigComponent,
    GroupComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    PluginManagementRoutingModule,
    ComponentModule,
    DirectiveModule,
    EoNgTreeModule,
    EoNgInputModule,
    EoNgButtonModule,
    NzFormModule,
    EoNgFeedbackModalModule,
    EoNgRadioModule,
    NzUploadModule,
    EoNgAutoCompleteModule,
    MarkdownModule.forChild(),
    EoNgApintoTableModule,
    EoNgSelectModule,
    EoNgEmptyModule
  ],
  exports: [
    PluginManagementComponent
  ]
})
export class PluginManagementModule { }
