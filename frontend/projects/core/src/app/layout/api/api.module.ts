import { NgModule } from '@angular/core'
import { CommonModule, registerLocaleData } from '@angular/common'

import { ApiRoutingModule } from './api-routing.module'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCascaderModule } from 'eo-ng-cascader'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { EoNgFeedbackModalModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgTabsModule } from 'eo-ng-tabs'
import { EoNgTreeModule } from 'eo-ng-tree'
import { EoNgCodeboxModule } from 'eo-ng-codebox'
import { NzHighlightModule } from 'ng-zorro-antd/core/highlight'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzGridModule } from 'ng-zorro-antd/grid'
import { ComponentModule } from '../../component/component.module'
import { RouterComponent } from './router/router.component'
import { DirectiveModule } from '../../directive/directive.module'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { NzDropDownModule } from 'ng-zorro-antd/dropdown'
import { ApiListComponent } from './api-list/api-list.component'
import { ApiManagementGroupComponent } from './api-list/group/group.component'
import { ApiManagementListComponent } from './api-list/list/list.component'
import { ApiMessageComponent } from './api-list/message/message.component'
import { ApiPublishComponent } from './api-list/publish/single/publish.component'
import { ApiImportComponent } from './api-list/import/import.component'
import { ApiManagementProxyComponent } from './api-list/proxy/proxy.component'
import { ApiManagementEditGroupComponent } from './api-list/group/edit-group/edit-group.component'
import { EoNgCopyModule } from 'eo-ng-copy'
import { ApiHttpCreateComponent } from './api-list/create/http-create/http-create.component'
import { ApiWebsocketCreateComponent } from './api-list/create/websocket-create/websocket-create.component'
import { ApiHttpMessageComponent } from './api-list/message/http-message/http-message.component'
import { ApiWebsocketMessageComponent } from './api-list/message/websocket-message/websocket-message.component'
import { ApiBatchPublishResultComponent } from './api-list/publish/batch/result.component'
import { NzSpinModule } from 'ng-zorro-antd/spin'
import { ApiBatchPublishChooseClusterComponent } from './api-list/publish/batch/choose-cluster.component'
import zh from '@angular/common/locales/zh'

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

registerLocaleData(zh)
@NgModule({
  declarations: [
    ApiManagementGroupComponent,
    ApiManagementListComponent,
    ApiMessageComponent,
    ApiPublishComponent,
    RouterComponent,
    ApiImportComponent,
    ApiManagementProxyComponent,
    ApiManagementEditGroupComponent,
    ApiListComponent,
    ApiHttpCreateComponent,
    ApiWebsocketCreateComponent,
    ApiHttpMessageComponent,
    ApiWebsocketMessageComponent,
    ApiBatchPublishResultComponent,
    ApiBatchPublishChooseClusterComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    ComponentModule,
    ApiRoutingModule,
    DirectiveModule,
    ReactiveFormsModule,
    NzFormModule,
    NzGridModule,
    NzHighlightModule,
    NzUploadModule,
    NzDropDownModule,
    NzSpinModule,
    ...sharedEoLibraryModules
  ],
  exports: [
  ]
})
export class ApiModule { }
