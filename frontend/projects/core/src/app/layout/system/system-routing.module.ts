import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { ExternalAppCreateComponent } from './external-app/create/create.component'
import { ExternalAppComponent } from './external-app/external-app.component'
import { ExternalAppListComponent } from './external-app/list/list.component'
import { ExternalAppMessageComponent } from './external-app/message/message.component'
import { SystemComponent } from './system.component'

const routes: Routes = [
  {
    path: '',
    component: SystemComponent,
    data: {
      id: '6'
    },
    children: [
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
      }
    ]
  }]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class SystemRoutingModule { }
