import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { PluginTemplateRoutingModule } from './plugin-template-routing.module'
import { ApiPluginConfigFormComponent } from './config/form/form.component'
import { ApiPluginTemplateContentComponent } from './content/content.component'
import { ApiPluginTemplateCreateComponent } from './create/create.component'
import { ApiPluginTemplateListComponent } from './list/list.component'
import { ApiPluginTemplateMessageComponent } from './message/message.component'
import { ApiPluginTemplateComponent } from './plugin.component'
import { ApiPluginTemplatePublishComponent } from './publish/publish.component'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCascaderModule } from 'eo-ng-cascader'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgCodeboxModule } from 'eo-ng-codebox'
import { EoNgCopyModule } from 'eo-ng-copy'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { EoNgFeedbackTooltipModule, EoNgFeedbackModalModule } from 'eo-ng-feedback'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgTabsModule } from 'eo-ng-tabs'
import { EoNgTreeModule } from 'eo-ng-tree'
import { NzHighlightModule } from 'ng-zorro-antd/core/highlight'
import { NzDropDownModule } from 'ng-zorro-antd/dropdown'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzGridModule } from 'ng-zorro-antd/grid'
import { NzSpinModule } from 'ng-zorro-antd/spin'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { ComponentModule } from '../../component/component.module'
import { DirectiveModule } from '../../directive/directive.module'

const sharedEoLibraryModules = [
  EoNgInputModule,
  EoNgButtonModule,
  EoNgSelectModule,
  EoNgFeedbackTooltipModule,
  EoNgCheckboxModule,
  EoNgDropdownModule,
  EoNgDatePickerModule,
  EoNgSwitchModule,
  EoNgCopyModule,
  EoNgApintoTableModule,
  EoNgCascaderModule,
  EoNgCodeboxModule,
  EoNgDatePickerModule,
  EoNgTreeModule,
  EoNgTabsModule,
  EoNgFeedbackModalModule
]
@NgModule({
  declarations: [
    ApiPluginTemplateComponent,
    ApiPluginTemplateListComponent,
    ApiPluginTemplateCreateComponent,
    ApiPluginTemplateContentComponent,
    ApiPluginTemplateMessageComponent,
    ApiPluginConfigFormComponent,
    ApiPluginTemplatePublishComponent],
  imports: [
    CommonModule,
    PluginTemplateRoutingModule,
    FormsModule,
    ComponentModule,
    DirectiveModule,
    ReactiveFormsModule,
    NzFormModule,
    NzGridModule,
    NzHighlightModule,
    NzUploadModule,
    NzDropDownModule,
    NzSpinModule,
    ...sharedEoLibraryModules
  ]
})
export class PluginTemplateModule { }
