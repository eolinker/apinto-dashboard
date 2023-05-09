import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { DeployRoutingModule } from './deploy-routing.module'
import { DeployClusterCertComponent } from './cluster/cert/cert.component'
import { DeployClusterComponent } from './cluster/cluster.component'
import { DeployClusterContentComponent } from './cluster/content/content.component'
import { DeployClusterCreateComponent } from './cluster/create/create.component'
import { DeployClusterEnvironmentComponent } from './cluster/environment/environment.component'
import { DeployClusterListComponent } from './cluster/list/list.component'
import { DeployClusterNodesComponent } from './cluster/nodes/nodes.component'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgBreadcrumbModule } from 'eo-ng-breadcrumb'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCascaderModule } from 'eo-ng-cascader'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { EoNgFeedbackModalModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgLayoutModule } from 'eo-ng-layout'
import { EoNgMenuModule } from 'eo-ng-menu'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgTableModule } from 'eo-ng-table'
import { EoNgTabsModule } from 'eo-ng-tabs'
import { EoNgTreeModule } from 'eo-ng-tree'
import { NzAvatarModule } from 'ng-zorro-antd/avatar'
import { NzCheckboxModule } from 'ng-zorro-antd/checkbox'
import { NzHighlightModule } from 'ng-zorro-antd/core/highlight'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzGridModule } from 'ng-zorro-antd/grid'
import { NzLayoutModule } from 'ng-zorro-antd/layout'
import { NzModalModule } from 'ng-zorro-antd/modal'
import { NzTableModule } from 'ng-zorro-antd/table'
import { NzToolTipModule } from 'ng-zorro-antd/tooltip'
import { ComponentModule } from '../../component/component.module'
import { EoNgTransferModule } from '../../component/transfer/transfer.module'
import { DirectiveModule } from '../../directive/directive.module'
import { DeployEnvironmentCreateComponent } from './environment/create/create.component'
import { DeployEnvironmentListComponent } from './environment/list/list.component'
import { DeployEnvironmentComponent } from './environment/environment.component'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { EoNgCollapseModule } from 'eo-ng-collapse'
import { EoNgCopyModule } from 'eo-ng-copy'
import { EoNgCodeboxModule } from 'eo-ng-codebox'
import { DeployClusterCertFormComponent } from './cluster/cert/form/form.component'
import { DeployClusterEnvironmentConfigFormComponent } from './cluster/environment/config/form/form.component'
import { DeployClusterEnvironmentConfigUpdateComponent } from './cluster/environment/config/update/update.component'
import { DeployClusterEnvironmentPublishComponent } from './cluster/environment/publish/publish.component'
import { DeployClusterEnvironmentHistoryChangeComponent } from './cluster/environment/history/change/change.component'
import { DeployClusterEnvironmentHistoryPublishComponent } from './cluster/environment/history/publish/publish.component'
import { DeployClusterNodesFormComponent } from './cluster/nodes/form/form.component'
import { DeployEnvironmentDetailComponent } from './environment/detail/detail.component'
import { DeployPluginComponent } from './plugin/deploy-plugin.component'
import { DeployPluginCreateComponent } from './plugin/create/create.component'
import { DeployPluginListComponent } from './plugin/list/list.component'
import { DeployPluginMessageComponent } from './plugin/message/message.component'
import { DeployClusterPluginComponent } from './cluster/plugin/plugin.component'
import { DeployClusterPluginConfigFormComponent } from './cluster/plugin/config/config.component'
import { DeployClusterPluginHistoryChangeComponent } from './cluster/plugin/history/change/change.component'
import { DeployClusterPluginHistoryPublishComponent } from './cluster/plugin/history/publish/publish.component'
import { DeployClusterPluginPublishComponent } from './cluster/plugin/publish/publish.component'
import { DeployClusterMessageComponent } from './cluster/message/message.component'

@NgModule({
  declarations: [
    DeployClusterComponent,
    DeployClusterListComponent,
    DeployClusterContentComponent,
    DeployClusterEnvironmentComponent,
    DeployClusterCertComponent,
    DeployClusterNodesComponent,
    DeployEnvironmentListComponent,
    DeployClusterCreateComponent,
    DeployEnvironmentCreateComponent,
    DeployEnvironmentComponent,
    DeployClusterCertFormComponent,
    DeployClusterEnvironmentConfigFormComponent,
    DeployClusterEnvironmentConfigUpdateComponent,
    DeployClusterEnvironmentPublishComponent,
    DeployClusterEnvironmentHistoryChangeComponent,
    DeployClusterEnvironmentHistoryPublishComponent,
    DeployClusterNodesFormComponent,
    DeployEnvironmentDetailComponent,
    DeployPluginComponent,
    DeployPluginListComponent,
    DeployPluginCreateComponent,
    DeployPluginMessageComponent,
    DeployClusterPluginComponent,
    DeployClusterPluginConfigFormComponent,
    DeployClusterPluginHistoryChangeComponent,
    DeployClusterPluginHistoryPublishComponent,
    DeployClusterPluginPublishComponent,
    DeployClusterMessageComponent
  ],
  imports: [
    EoNgLayoutModule,
    EoNgMenuModule,
    CommonModule,
    EoNgBreadcrumbModule,
    EoNgSelectModule,
    FormsModule,
    ComponentModule,
    EoNgTableModule,
    EoNgFeedbackModalModule,
    EoNgFeedbackTooltipModule,
    EoNgTabsModule,
    EoNgCheckboxModule,
    NzTableModule,
    EoNgInputModule,
    EoNgCascaderModule,
    EoNgTreeModule,
    EoNgDropdownModule,
    EoNgDatePickerModule,
    EoNgSwitchModule,
    NzAvatarModule,
    EoNgButtonModule,
    EoNgTransferModule,
    NzFormModule,
    ReactiveFormsModule,
    NzLayoutModule,
    NzToolTipModule,
    NzCheckboxModule,
    NzGridModule,
    NzModalModule,
    NzHighlightModule,
    DirectiveModule,
    DeployRoutingModule,
    EoNgApintoTableModule,
    EoNgCollapseModule,
    EoNgCopyModule,
    EoNgCodeboxModule
  ]
})
export class DeployModule { }
