import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { DeployClusterCertComponent } from './cert/cert.component'
import { DeployClusterContentComponent } from './cert/content/content.component'
import { DeployClusterCreateComponent } from './cert/create/create.component'
import { DeployClusterComponent } from './cluster.component'
import { DeployClusterEnvironmentComponent } from './environment/environment.component'
import { DeployClusterListComponent } from './list/list.component'
import { DeployClusterMessageComponent } from './message/message.component'
import { DeployClusterNodesComponent } from './nodes/nodes.component'
import { DeployClusterPluginComponent } from './plugin/plugin.component'
import { DeployClusterSmcertComponent } from './smcert/smcert.component'

const routes: Routes = [{
  path: '',
  component: DeployClusterComponent,
  data: {
    id: '1'
  },
  children: [
    {
      path: '',
      component: DeployClusterListComponent
    },
    {
      path: 'create',
      component: DeployClusterCreateComponent
    },
    {
      path: 'content/:clusterName',
      component: DeployClusterContentComponent,
      children: [
        {
          path: 'env',
          component: DeployClusterEnvironmentComponent
        },
        {
          path: 'cert',
          component: DeployClusterCertComponent
        },
        {
          path: '',
          component: DeployClusterNodesComponent
        },
        {
          path: 'plugin',
          component: DeployClusterPluginComponent
        },
        {
          path: 'message',
          component: DeployClusterMessageComponent
        },
        {
          path: 'smcert',
          component: DeployClusterSmcertComponent
        }
      ]
    }
  ]
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class DeployClusterRoutingModule { }
