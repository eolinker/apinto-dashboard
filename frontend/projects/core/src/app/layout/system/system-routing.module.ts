import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { RolesListComponent } from 'projects/eo-ng-apinto-user/src/public-api'
import { SystemEmailConfigComponent } from './email/config/config.component'
import { SystemEmailComponent } from './email/system-email.component'
import { ExternalAppCreateComponent } from './external-app/create/create.component'
import { ExternalAppComponent } from './external-app/external-app.component'
import { ExternalAppListComponent } from './external-app/list/list.component'
import { ExternalAppMessageComponent } from './external-app/message/message.component'
import { SystemRoleComponent } from './role/role.component'
import { SystemComponent } from './system.component'
import { SystemWebhookListComponent } from './webhook/list/list.component'
import { SystemWebhookComponent } from './webhook/webhook.component'

const routes: Routes = [
  {
    path: '',
    component: SystemComponent,
    data: {
      id: '6'
    },
    children: [{
      path: 'role',
      component: SystemRoleComponent,
      data: {
        id: '601'
      },
      children: [
        {
          path: ':roleId',
          component: RolesListComponent
        },
        {
          path: '',
          component: RolesListComponent
        }
      ]
    },
    {
      path: 'ext-app',
      component: ExternalAppComponent,
      data: {
        id: '602'
      },
      children: [
        {
          path: '',
          component: ExternalAppListComponent
        },
        {
          path: 'create',
          component: ExternalAppCreateComponent
        },
        {
          path: 'message/:extAppId',
          component: ExternalAppMessageComponent
        }
      ]
    },
    {
      path: 'email',
      component: SystemEmailComponent,
      data: {
        id: '603'
      },
      children: [
        {
          path: '',
          component: SystemEmailConfigComponent
        }
      ]
    },
    {
      path: 'webhook',
      component: SystemWebhookComponent,
      data: {
        id: '604'
      },
      children: [
        {
          path: '',
          component: SystemWebhookListComponent
        }
      ]
    }
    ]
  }]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class SystemRoutingModule { }
