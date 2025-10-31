import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

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
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { DeployEnvironmentCreateComponent } from './create/create.component'
import { DeployEnvironmentDetailComponent } from './detail/detail.component'
import { DeployEnvironmentComponent } from './environment.component'
import { DeployEnvironmentListComponent } from './list/list.component'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { NzFormModule } from 'ng-zorro-antd/form'
import { ComponentModule } from '../../component/component.module'
import { DirectiveModule } from '../../directive/directive.module'
import { GlobalEnvVarRoutingModule } from './global-env-var-routing.module'

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
    DeployEnvironmentListComponent,
    DeployEnvironmentCreateComponent,
    DeployEnvironmentComponent,
    DeployEnvironmentDetailComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    DirectiveModule,
    ComponentModule,
    GlobalEnvVarRoutingModule,
    ...sharedEoLibraryModules,
    NzFormModule
  ]
})
export class GlobalEnvVarModule { }
