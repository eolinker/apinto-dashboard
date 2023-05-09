import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { UpstreamRoutingModule } from './upstream-routing.module'
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
import { ServiceDiscoveryContentComponent } from './service-discovery/content/content.component'
import { ServiceDiscoveryCreateComponent } from './service-discovery/create/create.component'
import { ServiceDiscoveryListComponent } from './service-discovery/list/list.component'
import { ServiceDiscoveryMessageComponent } from './service-discovery/message/message.component'
import { ServiceDiscoveryPublishComponent } from './service-discovery/publish/publish.component'
import { ServiceDiscoveryComponent } from './service-discovery/service-discovery.component'
import { UpstreamContentComponent } from './upstream/content/content.component'
import { UpstreamCreateComponent } from './upstream/create/create.component'
import { UpstreamListComponent } from './upstream/list/list.component'
import { UpstreamMessageComponent } from './upstream/message/message.component'
import { UpstreamPublishComponent } from './upstream/publish/publish.component'
import { UpstreamComponent } from './upstream/upstream.component'
import { DirectiveModule } from '../../directive/directive.module'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'

@NgModule({
  declarations: [
    UpstreamComponent,
    UpstreamListComponent,
    UpstreamContentComponent,
    UpstreamPublishComponent,
    UpstreamMessageComponent,
    UpstreamCreateComponent,
    ServiceDiscoveryComponent,
    ServiceDiscoveryListComponent,
    ServiceDiscoveryCreateComponent,
    ServiceDiscoveryContentComponent,
    ServiceDiscoveryPublishComponent,
    ServiceDiscoveryMessageComponent
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
    UpstreamRoutingModule,
    EoNgApintoTableModule
  ]
})
export class UpstreamModule { }
