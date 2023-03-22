import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { ApiRoutingModule } from './api-routing.module'
import { ApiContentComponent } from './content/content.component'
import { ApiCreateComponent } from './create/create.component'
import { ApiManagementComponent } from './group/group.component'
import { ApiManagementListComponent } from './list/list.component'
import { ApiMessageComponent } from './message/message.component'
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
import { RouterComponent } from './router/router.component'
import { DirectiveModule } from '../../directive/directive.module'
import { MatchFormComponent } from './match/form/form.component'
import { MatchTableComponent } from './match/table/table.component'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { ApiImportComponent } from './import/import.component'
import { NzUploadModule } from 'ng-zorro-antd/upload'
import { NzDropDownModule } from 'ng-zorro-antd/dropdown'
import { ApiBatchPublishComponent } from './publish/batch/publish.component'
import { ApiPublishComponent } from './publish/single/publish.component'
import { ApiManagementProxyComponent } from './proxy/proxy.component'
import { ApiManagementEditGroupComponent } from './group/edit-group/edit-group.component'
import { EoNgCopyModule } from 'eo-ng-copy'

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
    ApiManagementEditGroupComponent
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
    EoNgCopyModule
  ],
  exports: [
    MatchTableComponent,
    MatchFormComponent
  ]
})
export class ApiModule { }
