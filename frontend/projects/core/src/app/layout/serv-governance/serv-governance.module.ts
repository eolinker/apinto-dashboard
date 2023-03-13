import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { ServGovernanceRoutingModule } from './serv-governance-routing.module'
import { TrafficCreateComponent } from './traffic/create/create.component'
import { TrafficComponent } from './traffic/traffic.component'
import { TrafficMessageComponent } from './traffic/message/message.component'
import { ServiceGovernanceComponent } from './service-governance.component'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgBreadcrumbModule } from 'eo-ng-breadcrumb'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCascaderModule } from 'eo-ng-cascader'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { EoNgRadioModule } from 'eo-ng-radio'
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
import { EoNgApintoUserModule } from 'projects/eo-ng-apinto-user/src/public-api'
import { ComponentModule } from '../../component/component.module'
import { EoNgTransferModule } from '../../component/transfer/transfer.module'
import { DirectiveModule } from '../../directive/directive.module'
import { GroupComponent } from './group/group.component'
import { ListComponent } from './list/list.component'
import { FuseComponent } from './fuse/fuse.component'
import { CacheComponent } from './cache/cache.component'
import { VisitComponent } from './visit/visit.component'
import { GreyComponent } from './grey/grey.component'
import { CacheCreateComponent } from './cache/create/create.component'
import { CacheMessageComponent } from './cache/message/message.component'
import { FuseCreateComponent } from './fuse/create/create.component'
import { FuseMessageComponent } from './fuse/message/message.component'
import { GreyCreateComponent } from './grey/create/create.component'
import { GreyMessageComponent } from './grey/message/message.component'
import { VisitCreateComponent } from './visit/create/create.component'
import { VisitMessageComponent } from './visit/message/message.component'
import { EoNgAutoCompleteModule } from 'eo-ng-auto-complete'
import { NzSliderModule } from 'ng-zorro-antd/slider'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { ApiModule } from '../api/api.module'
import { ResponseFormComponent } from './response-form/response-form.component'
import { FilterFooterComponent } from './filter/footer/footer.component'
import { FilterFormComponent } from './filter/form/form.component'
import { FilterTableComponent } from './filter/table/table.component'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { ServiceGovernancePublishComponent } from './publish/publish.component'
@NgModule({
  declarations: [
    ServiceGovernanceComponent,
    TrafficComponent,
    TrafficCreateComponent,
    TrafficMessageComponent,
    GroupComponent,
    ListComponent,
    FuseComponent,
    CacheComponent,
    VisitComponent,
    GreyComponent,
    CacheCreateComponent,
    CacheMessageComponent,
    FuseCreateComponent,
    FuseMessageComponent,
    GreyCreateComponent,
    GreyMessageComponent,
    VisitCreateComponent,
    VisitMessageComponent,
    ResponseFormComponent,
    FilterFooterComponent,
    FilterFormComponent,
    FilterTableComponent,
    ServiceGovernancePublishComponent
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
    EoNgApintoUserModule,
    NzFormModule,
    ReactiveFormsModule,
    NzLayoutModule,
    NzCheckboxModule,
    NzGridModule,
    NzModalModule,
    NzHighlightModule,
    DirectiveModule,
    ServGovernanceRoutingModule,
    EoNgAutoCompleteModule,
    EoNgRadioModule,
    NzSliderModule,
    NzInputNumberModule,
    ApiModule,
    EoNgApintoTableModule
  ]
})
export class ServGovernanceModule { }
