import { NgModule } from '@angular/core'
import { RouterModule, Routes } from '@angular/router'
import { DeployClusterCertComponent } from './cluster/cert/cert.component'
import { DeployClusterComponent } from './cluster/cluster.component'
import { DeployClusterContentComponent } from './cluster/content/content.component'
import { DeployClusterCreateComponent } from './cluster/create/create.component'
import { DeployClusterEnvironmentComponent } from './cluster/environment/environment.component'
import { DeployClusterListComponent } from './cluster/list/list.component'
import { DeployClusterNodesComponent } from './cluster/nodes/nodes.component'
import { DeployEnvironmentCreateComponent } from './environment/create/create.component'
import { DeployEnvironmentComponent } from './environment/environment.component'
import { DeployEnvironmentListComponent } from './environment/list/list.component'
import { DeployPluginComponent } from './plugin/deploy-plugin.component'
import { DeployPluginCreateComponent } from './plugin/create/create.component'
import { DeployPluginListComponent } from './plugin/list/list.component'
import { DeployPluginMessageComponent } from './plugin/message/message.component'
import { DeployClusterPluginComponent } from './cluster/plugin/plugin.component'
const routes: Routes = [{
  path: '',
  component: DeployClusterComponent,
  data: {
    id: '1'
  },
  children: [
    {
      path: 'cluster',
      component: DeployClusterComponent,
      data: {
        id: '101',
        parentId: '1'
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
              path: '',
              component: DeployClusterEnvironmentComponent
            },
            {
              path: 'cert',
              component: DeployClusterCertComponent
            },
            {
              path: 'nodes',
              component: DeployClusterNodesComponent
            },
            {
              path: 'plugin',
              component: DeployClusterPluginComponent
            }
          ]
        }]
    },
    {
      path: 'env',
      component: DeployEnvironmentComponent,
      data: {
        id: '102',
        parentId: '1'
      },
      children: [
        {
          path: '',
          component: DeployEnvironmentListComponent
        },
        {
          path: 'create',
          component: DeployEnvironmentCreateComponent
        }
      ]
    },
    {
      path: 'plugin',
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
}]

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class DeployRoutingModule { }
