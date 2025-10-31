import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { TrafficStrategyRoutingModule } from './traffic-strategy-routing.module'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgBreadcrumbModule } from 'eo-ng-breadcrumb'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCascaderModule } from 'eo-ng-cascader'
import { EoNgCheckboxModule } from 'eo-ng-checkbox'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
import { EoNgDropdownModule } from 'eo-ng-dropdown'
import { EoNgEmptyModule } from 'eo-ng-empty'
import { EoNgFeedbackModalModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgLayoutModule } from 'eo-ng-layout'
import { EoNgMenuModule } from 'eo-ng-menu'
import { EoNgRadioModule } from 'eo-ng-radio'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgTabsModule } from 'eo-ng-tabs'
import { EoNgTreeModule } from 'eo-ng-tree'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzGridModule } from 'ng-zorro-antd/grid'
import { NzInputNumberModule } from 'ng-zorro-antd/input-number'
import { NzSliderModule } from 'ng-zorro-antd/slider'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { ComponentModule } from '../../../component/component.module'
import { DirectiveModule } from '../../../directive/directive.module'
import { TrafficCreateComponent } from './create/create.component'
import { TrafficMessageComponent } from './message/message.component'
import { TrafficComponent } from './traffic.component'
import { ServGovernanceModule } from '../serv-governance.module'

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
    TrafficComponent,
    TrafficCreateComponent,
    TrafficMessageComponent],
  imports: [
    CommonModule,
    FormsModule,
    TrafficStrategyRoutingModule,
    ComponentModule,
    ReactiveFormsModule,
    DirectiveModule,
    ...sharedEoLibraryModules,
    NzFormModule,
    NzGridModule,
    NzSliderModule,
    NzInputNumberModule,
    ServGovernanceModule
  ]
})
export class TrafficStrategyModule { }
