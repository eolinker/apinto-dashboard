import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { NodePluginRoutingModule } from './node-plugin-routing.module'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgBreadcrumbModule } from 'eo-ng-breadcrumb'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCascaderModule } from 'eo-ng-cascader'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgCodeboxModule } from 'eo-ng-codebox'
import { EoNgCollapseModule } from 'eo-ng-collapse'
import { EoNgCopyModule } from 'eo-ng-copy'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { EoNgFeedbackModalModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgLayoutModule } from 'eo-ng-layout'
import { EoNgMenuModule } from 'eo-ng-menu'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgTabsModule } from 'eo-ng-tabs'
import { EoNgTreeModule } from 'eo-ng-tree'
import { NzFormModule } from 'ng-zorro-antd/form'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { ComponentModule } from '../../component/component.module'
import { DirectiveModule } from '../../directive/directive.module'
import { DeployPluginCreateComponent } from './create/create.component'
import { DeployPluginComponent } from './deploy-plugin.component'
import { DeployPluginListComponent } from './list/list.component'
import { DeployPluginMessageComponent } from './message/message.component'

const sharedEoLibraryModules = [
  EoNgButtonModule,
  EoNgBreadcrumbModule,
  EoNgLayoutModule,
  EoNgMenuModule,
  EoNgSelectModule,
  EoNgFeedbackModalModule,
  EoNgFeedbackTooltipModule,
  EoNgTabsModule,
  EoNgCheckboxModule,
  EoNgInputModule,
  EoNgCascaderModule,
  EoNgTreeModule,
  EoNgDropdownModule,
  EoNgDatePickerModule,
  EoNgSwitchModule,
  EoNgCollapseModule,
  EoNgApintoTableModule,
  EoNgCopyModule,
  EoNgCodeboxModule
]
@NgModule({
  declarations: [
    DeployPluginComponent,
    DeployPluginListComponent,
    DeployPluginCreateComponent,
    DeployPluginMessageComponent],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    DirectiveModule,
    ComponentModule,
    ...sharedEoLibraryModules,
    NzFormModule,
    NodePluginRoutingModule
  ]
})
export class NodePluginModule { }
