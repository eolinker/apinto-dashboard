import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { PluginManagementRoutingModule } from './plugin-management-routing.module'
import { PluginManagementComponent } from './plugin-management.component'
import { PluginListComponent } from './list/list.component'
import { PluginMessageComponent } from './message/message.component'
import { ComponentModule } from '../../component/component.module'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { PluginCreateComponent } from './create/create.component'
import { MarkdownModule } from 'ngx-markdown'
import { PluginConfigComponent } from './config/config.component'
import { GroupComponent } from './group/group.component'
import { DirectiveModule } from '../../directive/directive.module'
import { EoNgEmptyModule } from 'eo-ng-empty'
import { NzFormModule } from 'ng-zorro-antd/form'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { NzGridModule } from 'ng-zorro-antd/grid'
import { EoNgRadioModule } from 'eo-ng-radio'
import { EoNgTreeModule } from 'eo-ng-tree'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgInputModule } from 'eo-ng-input'

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
    MarkdownModule.forChild(),
    EoNgEmptyModule,
    EoNgInputModule,
    NzFormModule,
    EoNgApintoTableModule,
    NzGridModule,
    EoNgRadioModule,
    EoNgTreeModule,
    NzUploadModule,
    EoNgButtonModule
  ],
  exports: [
    PluginManagementComponent
  ]
})
export class PluginManagementModule { }
