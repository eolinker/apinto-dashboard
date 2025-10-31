import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { OpenApiRoutingModule } from './open-api-routing.module'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgCopyModule } from 'eo-ng-copy'
import { EoNgFeedbackTooltipModule, EoNgFeedbackAlertModule } from 'eo-ng-feedback'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgSelectModule } from 'eo-ng-select'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { ExternalAppCreateComponent } from './create/create.component'
import { ExternalAppComponent } from './external-app.component'
import { ExternalAppListComponent } from './list/list.component'
import { ExternalAppMessageComponent } from './message/message.component'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { NzFormModule } from 'ng-zorro-antd/form'
import { ComponentModule } from '../../component/component.module'
import { DirectiveModule } from '../../directive/directive.module'

const sharedEoLibraryModules = [
  EoNgButtonModule,
  EoNgCopyModule,
  EoNgSelectModule,
  EoNgFeedbackTooltipModule,
  EoNgInputModule,
  EoNgSwitchModule,
  EoNgApintoTableModule,
  EoNgFeedbackAlertModule
]

@NgModule({
  declarations: [
    ExternalAppComponent,
    ExternalAppListComponent,
    ExternalAppCreateComponent,
    ExternalAppMessageComponent],
  imports: [
    CommonModule,
    OpenApiRoutingModule,
    ComponentModule,
    FormsModule,
    ReactiveFormsModule,
    DirectiveModule,
    ...sharedEoLibraryModules,
    NzFormModule
  ]
})
export class OpenApiModule { }
