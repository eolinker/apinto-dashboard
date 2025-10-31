import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

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
import { DeployClusterCertComponent } from './cert/cert.component'
import { DeployClusterContentComponent } from './cert/content/content.component'
import { DeployClusterCreateComponent } from './cert/create/create.component'
import { DeployClusterCertFormComponent } from './cert/form/form.component'
import { DeployClusterComponent } from './cluster.component'
import { DeployClusterEnvironmentConfigFormComponent } from './environment/config/form/form.component'
import { DeployClusterEnvironmentConfigUpdateComponent } from './environment/config/update/update.component'
import { DeployClusterEnvironmentComponent } from './environment/environment.component'
import { DeployClusterEnvironmentHistoryChangeComponent } from './environment/history/change/change.component'
import { DeployClusterEnvironmentHistoryPublishComponent } from './environment/history/publish/publish.component'
import { DeployClusterEnvironmentPublishComponent } from './environment/publish/publish.component'
import { DeployClusterListComponent } from './list/list.component'
import { DeployClusterMessageComponent } from './message/message.component'
import { DeployClusterNodesFormComponent } from './nodes/form/form.component'
import { DeployClusterNodesComponent } from './nodes/nodes.component'
import { DeployClusterPluginConfigFormComponent } from './plugin/config/config.component'
import { DeployClusterPluginHistoryChangeComponent } from './plugin/history/change/change.component'
import { DeployClusterPluginHistoryPublishComponent } from './plugin/history/publish/publish.component'
import { DeployClusterPluginComponent } from './plugin/plugin.component'
import { DeployClusterPluginPublishComponent } from './plugin/publish/publish.component'
import { DeployClusterRoutingModule } from './deploy-cluster-routing.module'
import { DeployClusterSmcertComponent } from './smcert/smcert.component'
import { DeployClusterSmcertFormComponent } from './smcert/form/form.component'

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
    DeployClusterComponent,
    DeployClusterListComponent,
    DeployClusterContentComponent,
    DeployClusterEnvironmentComponent,
    DeployClusterCertComponent,
    DeployClusterNodesComponent,
    DeployClusterCreateComponent,
    DeployClusterCertFormComponent,
    DeployClusterEnvironmentConfigFormComponent,
    DeployClusterEnvironmentConfigUpdateComponent,
    DeployClusterEnvironmentPublishComponent,
    DeployClusterEnvironmentHistoryChangeComponent,
    DeployClusterEnvironmentHistoryPublishComponent,
    DeployClusterNodesFormComponent,
    DeployClusterPluginComponent,
    DeployClusterPluginConfigFormComponent,
    DeployClusterPluginHistoryChangeComponent,
    DeployClusterPluginHistoryPublishComponent,
    DeployClusterPluginPublishComponent,
    DeployClusterSmcertComponent,
    DeployClusterSmcertFormComponent,
    DeployClusterMessageComponent],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    DirectiveModule,
    ComponentModule,
    ...sharedEoLibraryModules,
    NzFormModule,
    DeployClusterRoutingModule
  ]
})
export class DeployClusterModule { }
