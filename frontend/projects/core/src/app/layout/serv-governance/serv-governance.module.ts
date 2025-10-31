import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

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
import { EoNgTabsModule } from 'eo-ng-tabs'
import { EoNgTreeModule } from 'eo-ng-tree'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzGridModule } from 'ng-zorro-antd/grid'
import { NzTimePickerModule } from 'ng-zorro-antd/time-picker'
import { ComponentModule } from '../../component/component.module'
import { DirectiveModule } from '../../directive/directive.module'
import { GroupComponent } from './group/group.component'
import { ListComponent } from './list/list.component'
import { NzSliderModule } from 'ng-zorro-antd/slider'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { ResponseFormComponent } from './response-form/response-form.component'
import { FilterFooterComponent } from './filter/footer/footer.component'
import { FilterFormComponent } from './filter/form/form.component'
import { FilterTableComponent } from './filter/table/table.component'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { ServiceGovernancePublishComponent } from './publish/publish.component'
import { EoNgEmptyModule } from 'eo-ng-empty'
import { RouterOutlet } from '@angular/router'

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
  EoNgEmptyModule,
  EoNgApintoTableModule,
  EoNgRadioModule
]

@NgModule({
  declarations: [
    GroupComponent,
    ListComponent,
    ResponseFormComponent,
    FilterFooterComponent,
    FilterFormComponent,
    FilterTableComponent,
    ServiceGovernancePublishComponent
  ],
  imports: [
    CommonModule,
    RouterOutlet,
    FormsModule,
    ComponentModule,
    ReactiveFormsModule,
    DirectiveModule,
    ...sharedEoLibraryModules,
    NzFormModule,
    NzGridModule,
    NzTimePickerModule,
    NzSliderModule,
    NzInputNumberModule
  ],
  exports: [
    GroupComponent,
    ListComponent,
    ResponseFormComponent,
    FilterFooterComponent,
    FilterFormComponent,
    FilterTableComponent,
    ServiceGovernancePublishComponent
  ]
})
export class ServGovernanceModule { }
