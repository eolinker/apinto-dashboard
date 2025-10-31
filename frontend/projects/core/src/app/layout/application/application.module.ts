import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { ApplicationRoutingModule } from './application-routing.module'
import { ApplicationManagementComponent } from './application.component'
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
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzGridModule } from 'ng-zorro-antd/grid'
import { ComponentModule } from '../../component/component.module'
import { ApplicationAuthenticationComponent } from './authentication/authentication.component'
import { ApplicationContentComponent } from './content/content.component'
import { ApplicationCreateComponent } from './create/create.component'
import { ApplicationManagementListComponent } from './list/list.component'
import { ApplicationMessageComponent } from './message/message.component'
import { ApplicationPublishComponent } from './publish/publish.component'
import { DirectiveModule } from '../../directive/directive.module'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { ApplicationAuthenticationFormComponent } from './authentication/form/form.component'
import { NzSpinModule } from 'ng-zorro-antd/spin'
import { ApplicationExtraComponent } from './extra/extra.component'
import { ApplicationAuthenticationViewComponent } from './authentication/view/view.component'
import { ApplicationExtraFormComponent } from './extra/form/form.component'

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
  EoNgTableModule,
  EoNgApintoTableModule
]
@NgModule({
  declarations: [
    ApplicationManagementComponent,
    ApplicationManagementListComponent,
    ApplicationCreateComponent,
    ApplicationContentComponent,
    ApplicationPublishComponent,
    ApplicationMessageComponent,
    ApplicationAuthenticationComponent,
    ApplicationAuthenticationFormComponent,
    ApplicationExtraComponent,
    ApplicationAuthenticationViewComponent,
    ApplicationExtraFormComponent
  ],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    ApplicationRoutingModule,
    DirectiveModule,
    ComponentModule,
    ...sharedEoLibraryModules,
    NzFormModule,
    NzGridModule,
    NzSpinModule
  ]
})
export class ApplicationModule { }
