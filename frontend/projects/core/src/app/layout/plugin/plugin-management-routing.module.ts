import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { PluginListComponent } from './list/list.component'
import { PluginManagementComponent } from './plugin-management.component'
import { PluginMessageComponent } from './message/message.component'

const routes: Routes = [
  {
    path: '',
    component: PluginManagementComponent,
    data: {
      id: 10
    },
    children: [
      {
        path: 'list/:pluginGroupId',
        component: PluginListComponent
      },
      {
        path: 'list',
        component: PluginListComponent
      },
      {
        path: 'message/:pluginId',
        component: PluginMessageComponent,
        children: [
          {
            path: ':mdFileName',
            component: PluginMessageComponent
          }
        ]
      }
    ]
  }
]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class PluginManagementRoutingModule { }
