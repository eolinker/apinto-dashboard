/*
 * @Date: 2023-12-12 18:57:19
 * @LastEditors: maggieyyy
 * @LastEditTime: 2023-12-14 19:37:46
 * @FilePath: \apinto\projects\core\src\app\layout\open-api\open-api-routing.module.ts
 */
import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { ExternalAppCreateComponent } from './create/create.component'
import { ExternalAppComponent } from './external-app.component'
import { ExternalAppListComponent } from './list/list.component'
import { ExternalAppMessageComponent } from './message/message.component'

const routes: Routes = [
  {
    path: '',
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

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class OpenApiRoutingModule { }
