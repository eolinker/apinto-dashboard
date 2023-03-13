import { NgModule } from '@angular/core'
import { CommonModule } from '@angular/common'

import { SystemRoutingModule } from './system-routing.module'
import { EoNgApintoTableModule } from 'projects/eo-ng-apinto-table/src/public-api'
import { ComponentModule } from '../../component/component.module'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgSwitchModule } from 'eo-ng-switch'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgApintoUserModule } from 'projects/eo-ng-apinto-user/src/public-api'
import { ExternalAppCreateComponent } from './external-app/create/create.component'
import { ExternalAppComponent } from './external-app/external-app.component'
import { ExternalAppListComponent } from './external-app/list/list.component'
import { ExternalAppMessageComponent } from './external-app/message/message.component'
import { SystemRoleComponent } from './role/role.component'
import { SystemComponent } from './system.component'
import { EoNgButtonModule } from 'eo-ng-button'
import { NzFormModule } from 'ng-zorro-antd/form'
import { DirectiveModule } from '../../directive/directive.module'
import { EoNgCopyModule } from 'eo-ng-copy'
import { EoNgFeedbackAlertModule, EoNgFeedbackTooltipModule } from 'eo-ng-feedback'
import { SystemEmailComponent } from './email/system-email.component'
import { SystemWebhookComponent } from './webhook/webhook.component'
import { SystemWebhookListComponent } from './webhook/list/list.component'
import { SystemEmailConfigComponent } from './email/config/config.component'
import { EoNgSelectModule } from 'eo-ng-select'
import { SystemWebhookConfigComponent } from './webhook/config/config.component'

@NgModule({
  declarations: [
    SystemComponent,
    SystemRoleComponent,
    ExternalAppComponent,
    ExternalAppListComponent,
    ExternalAppCreateComponent,
    ExternalAppMessageComponent,
    SystemEmailComponent,
    SystemWebhookComponent,
    SystemWebhookListComponent,
    SystemEmailConfigComponent,
    SystemWebhookConfigComponent
  ],
  imports: [
    CommonModule,
    SystemRoutingModule,
    EoNgApintoTableModule,
    EoNgApintoUserModule,
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
