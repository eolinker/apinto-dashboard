import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { DeployPluginCreateComponent } from './create/create.component'
import { DeployPluginComponent } from './deploy-plugin.component'
import { DeployPluginListComponent } from './list/list.component'
import { DeployPluginMessageComponent } from './message/message.component'

const routes: Routes = [{
  path: '',
  component: DeployPluginComponent,
  children: [
    {
      path: '',
      component: DeployPluginListComponent
    },
    {
      path: 'create',
      component: DeployPluginCreateComponent
    },
    {
      path: 'message/:pluginName',
      component: DeployPluginMessageComponent
    }
  ]
}
]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class NodePluginRoutingModule { }
