import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { ApiRoutingModule } from './api-routing.module'
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
import { EoNgCodeboxModule } from 'eo-ng-codebox'
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
import { RouterComponent } from './router/router.component'
import { DirectiveModule } from '../../directive/directive.module'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { NzDropDownModule } from 'ng-zorro-antd/dropdown'
import { ApiListComponent } from './api-list/api-list.component'
import { ApiContentComponent } from './api-list/content/content.component'
import { ApiCreateComponent } from './api-list/create/create.component'
import { ApiManagementComponent } from './api-list/group/group.component'
import { ApiManagementListComponent } from './api-list/list/list.component'
import { ApiMessageComponent } from './api-list/message/message.component'
import { ApiPublishComponent } from './api-list/publish/single/publish.component'
import { MatchFormComponent } from './api-list/match/form/form.component'
import { MatchTableComponent } from './api-list/match/table/table.component'
import { ApiImportComponent } from './api-list/import/import.component'
import { ApiBatchPublishComponent } from './api-list/publish/batch/publish.component'
import { ApiManagementProxyComponent } from './api-list/proxy/proxy.component'
import { ApiManagementEditGroupComponent } from './api-list/group/edit-group/edit-group.component'
import { ApiPluginTemplateComponent } from './plugin/plugin.component'
import { ApiPluginTemplateContentComponent } from './plugin/content/content.component'
import { ApiPluginTemplateCreateComponent } from './plugin/create/create.component'
import { ApiPluginTemplateListComponent } from './plugin/list/list.component'
import { ApiPluginTemplateMessageComponent } from './plugin/message/message.component'
import { ApiPluginConfigFormComponent } from './plugin/config/form/form.component'
import { ApiPluginConfigTableComponent } from './plugin/config/table/table.component'
import { ApiPluginTemplatePublishComponent } from './plugin/publish/publish.component'
import { EoNgCopyModule } from 'eo-ng-copy'
import { ApiHttpCreateComponent } from './api-list/create/http-create/http-create.component'
import { ApiWebsocketCreateComponent } from './api-list/create/websocket-create/websocket-create.component'
import { NzTypographyModule } from 'ng-zorro-antd/typography';
import { ApiHttpMessageComponent } from './api-list/message/http-message/http-message.component';
import { ApiWebsocketMessageComponent } from './api-list/message/websocket-message/websocket-message.component'

@NgModule({
  declarations: [
    ApiManagementComponent,
    ApiManagementListComponent,
    ApiCreateComponent,
    ApiContentComponent,
    ApiMessageComponent,
    ApiPublishComponent,
    RouterComponent,
    MatchFormComponent,
    MatchTableComponent,
    ApiImportComponent,
    ApiBatchPublishComponent,
    ApiManagementProxyComponent,
    ApiManagementEditGroupComponent,
    ApiListComponent,
    ApiPluginTemplateComponent,
    ApiPluginTemplateListComponent,
    ApiPluginTemplateCreateComponent,
    ApiPluginTemplateContentComponent,
    ApiPluginTemplateMessageComponent,
    ApiPluginConfigFormComponent,
    ApiPluginConfigTableComponent,
    ApiPluginTemplatePublishComponent,
    ApiHttpCreateComponent,
    ApiWebsocketCreateComponent,
    ApiHttpMessageComponent,
    ApiWebsocketMessageComponent
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
    ApiRoutingModule,
    DirectiveModule,
    EoNgApintoTableModule,
    NzUploadModule,
    NzDropDownModule,
    EoNgCodeboxModule,
    EoNgCopyModule,
    NzTypographyModule
  ],
  exports: [
    MatchTableComponent,
    MatchFormComponent
  ]
})
export class ApiModule { }
