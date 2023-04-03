import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { SystemRoutingModule } from './system-routing.module'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { ComponentModule } from '../../component/component.module'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgInputModule } from 'eo-ng-input'
import { ExternalAppCreateComponent } from './external-app/create/create.component'
import { ExternalAppComponent } from './external-app/external-app.component'
import { ExternalAppListComponent } from './external-app/list/list.component'
import { ExternalAppMessageComponent } from './external-app/message/message.component'
import { SystemComponent } from './system.component'
import { EoNgButtonModule } from 'eo-ng-button'
import { NzFormModule } from 'ng-zorro-antd/form'
import { DirectiveModule } from '../../directive/directive.module'
import { EoNgCopyModule } from 'eo-ng-copy'
import { EoNgFeedbackAlertModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { EoNgSelectModule } from 'eo-ng-select'

@NgModule({
  declarations: [
    SystemComponent,
    ExternalAppComponent,
    ExternalAppListComponent,
    ExternalAppCreateComponent,
    ExternalAppMessageComponent
  ],
  imports: [
    CommonModule,
    SystemRoutingModule,
    EoNgApintoTableModule,
    ComponentModule,
    FormsModule,
    ReactiveFormsModule,
    EoNgSwitchModule,
    EoNgInputModule,
    EoNgButtonModule,
    NzFormModule,
    DirectiveModule,
    EoNgCopyModule,
    EoNgSelectModule,
    EoNgFeedbackTooltipModule,
    EoNgFeedbackAlertModule
  ]
})
export class SystemModule { }
